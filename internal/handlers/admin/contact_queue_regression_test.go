package admin_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
	"golang.org/x/net/html"
)

func contactQueueLinks(t *testing.T, body string) map[string]string {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	links := map[string]string{}
	var text func(*html.Node) string
	text = func(n *html.Node) string {
		if n.Type == html.TextNode {
			return n.Data
		}
		var out string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			out += text(c)
		}
		return out
	}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links[strings.Join(strings.Fields(text(n)), " ")] = a.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return links
}

func TestContactQueueCombinedFiltersEncodedPaginationAndNavigation(t *testing.T) {
	db, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	term := "Queue A&B + # /?"
	ids := []int64{}
	for i := 0; i < 30; i++ {
		row, err := q.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{Name: fmt.Sprintf("%s %02d", term, i), Email: fmt.Sprintf("queue-%02d@example.test", i), Message: "Queue fixture"})
		if err != nil {
			t.Fatal(err)
		}
		status, kind := "closed", "contact"
		if i == 27 {
			status = "new"
		}
		if i == 28 {
			kind = "rfq"
		}
		name := fmt.Sprintf("%s %02d", term, i)
		if i == 29 {
			name = "Unrelated"
		}
		if _, err = db.Exec("UPDATE contact_submissions SET status=?,submission_type=?,name=?,created_at='2026-01-01 00:00:00' WHERE id=?", status, kind, name, row.ID); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, row.ID)
	}
	e := operationEcho()
	h := admin.NewAdminContactHandler(q, logger, nil)
	filters := url.Values{"search": {term}, "status": {"closed"}, "type": {"contact"}, "page": {"2"}}
	returnTo := "/admin/contact/submissions?" + filters.Encode()
	rec := operationGET(t, e, h.ListSubmissions, returnTo)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "27 submissions found") {
		t.Fatalf("combined filters/count failed: %d", rec.Code)
	}
	links := contactQueueLinks(t, rec.Body.String())
	previous, err := url.Parse(links["← Previous"])
	if err != nil {
		t.Fatal(err)
	}
	if previous.Query().Get("search") != term || previous.Query().Get("status") != "closed" || previous.Query().Get("type") != "contact" || previous.Query().Get("page") != "1" {
		t.Fatalf("pagination lost/changed scope: %s", previous)
	}
	detailURL := links[fmt.Sprintf("%s 00", term)]
	u, err := url.Parse(detailURL)
	if err != nil {
		t.Fatal(err)
	}
	if u.Query().Get("return_to") != returnTo {
		t.Fatalf("detail lost encoded scope: %s", detailURL)
	}
	if !strings.Contains(rec.Body.String(), `value="closed" selected`) {
		t.Fatal("closed filter selection missing")
	}
	// Identical timestamps must still support deterministic adjacent navigation.
	detail := func(id int64, queue string) string {
		r := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest("GET", fmt.Sprintf("/admin/contact/submissions/%d?%s", id, url.Values{"return_to": {queue}}.Encode()), nil), r)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(id))
		if err := h.ViewSubmission(c); err != nil {
			t.Fatal(err)
		}
		return r.Body.String()
	}
	body := detail(ids[1], returnTo)
	nav := contactQueueLinks(t, body)
	for label, id := range map[string]int64{"← Previous": ids[2], "Next →": ids[0]} {
		u, err := url.Parse(nav[label])
		if err != nil {
			t.Fatal(err)
		}
		if u.Path != fmt.Sprintf("/admin/contact/submissions/%d", id) || u.Query().Get("return_to") != returnTo {
			t.Fatalf("%s left filtered queue: %s", label, nav[label])
		}
	}
	if nav["← Back to Submissions"] != returnTo {
		t.Fatalf("Back lost scope: %s", nav["← Back to Submissions"])
	}
	// A single matching enquiry must never offer unrelated next/previous items.
	single := url.Values{"search": {"queue-05@example.test"}, "status": {"closed"}, "type": {"contact"}}
	one := contactQueueLinks(t, detail(ids[5], "/admin/contact/submissions?"+single.Encode()))
	if one["← Previous"] != "" || one["Next →"] != "" {
		t.Fatal("single-result queue contains adjacent navigation")
	}
	// Status changes retain context even when the item leaves that status queue.
	r, c := postForm(e, fmt.Sprintf("/admin/contact/submissions/%d/status", ids[1]), url.Values{"status": {"read"}, "notes": {"Reviewed fixture"}, "return_to": {returnTo}})
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(ids[1]))
	if err := h.UpdateSubmissionStatus(c); err != nil {
		t.Fatal(err)
	}
	target, err := url.Parse(r.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if r.Code != 303 || target.Query().Get("return_to") != returnTo {
		t.Fatal("status update lost queue")
	}
	persisted, err := q.GetContactSubmissionByID(ctx, ids[1])
	if err != nil {
		t.Fatal(err)
	}
	if persisted.Status != "read" || persisted.Notes.String != "Reviewed fixture" {
		t.Fatal("status/notes not persisted")
	}
}

func TestContactQueueReturnContextRejectsExternalOrUnexpectedPaths(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	row, err := q.CreateContactSubmission(context.Background(), sqlc.CreateContactSubmissionParams{Name: "Return context fixture", Email: "return@example.test", Message: "Test"})
	if err != nil {
		t.Fatal(err)
	}
	e := operationEcho()
	h := admin.NewAdminContactHandler(q, logger, nil)
	for _, raw := range []string{"https://example.test/phishing", "//example.test/admin/contact/submissions", "/admin/settings", "/admin/contact/submissions#bad", "/admin/contact/submissions/../settings", "%"} {
		t.Run(raw, func(t *testing.T) {
			r, c := postForm(e, fmt.Sprintf("/admin/contact/submissions/%d/status", row.ID), url.Values{"status": {"read"}, "return_to": {raw}})
			c.SetParamNames("id")
			c.SetParamValues(fmt.Sprint(row.ID))
			if err := h.UpdateSubmissionStatus(c); err != nil {
				t.Fatal(err)
			}
			target, err := url.Parse(r.Header().Get("Location"))
			if err != nil {
				t.Fatal(err)
			}
			if target.IsAbs() || target.Host != "" || target.Query().Get("return_to") != "/admin/contact/submissions?page=1" {
				t.Fatalf("unsafe return context: %s", target)
			}
		})
	}
	// Unsupported fields and extreme page values cannot enter generated navigation.
	raw := "/admin/contact/submissions?status=bogus&type=bogus&page=9223372036854775807&next=https://example.test"
	r, c := postForm(e, fmt.Sprintf("/admin/contact/submissions/%d/status", row.ID), url.Values{"status": {"read"}, "return_to": {raw}})
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(row.ID))
	if err := h.UpdateSubmissionStatus(c); err != nil {
		t.Fatal(err)
	}
	target, _ := url.Parse(r.Header().Get("Location"))
	if target.Query().Get("return_to") != "/admin/contact/submissions?page=1" {
		t.Fatalf("unexpected fields retained: %s", target)
	}
}
