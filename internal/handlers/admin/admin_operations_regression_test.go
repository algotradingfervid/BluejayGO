package admin_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
	"golang.org/x/net/html"
)

func operationEcho() *echo.Echo {
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	return e
}

func operationGET(t *testing.T, e *echo.Echo, handler echo.HandlerFunc, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	if err := handler(e.NewContext(httptest.NewRequest(http.MethodGet, path, nil), rec)); err != nil {
		t.Fatal(err)
	}
	return rec
}

func TestContactInboxRealRendererPaginationAndDetail(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	var id int64
	for i := 0; i < 27; i++ {
		row, err := q.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{Name: fmt.Sprintf("Inbox fixture %02d", i), Email: "fixture@example.test", Message: "A persisted inquiry for the inbox"})
		if err != nil {
			t.Fatal(err)
		}
		id = row.ID
	}
	e := operationEcho()
	h := admin.NewAdminContactHandler(q, logger, nil)
	for _, page := range []string{"1", "2"} {
		rec := operationGET(t, e, h.ListSubmissions, "/admin/contact/submissions?page="+page)
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Inbox fixture") || !strings.Contains(rec.Body.String(), "27 submissions found") {
			t.Fatalf("page %s did not render inbox: %d", page, rec.Code)
		}
		target := "?page=2"
		if page == "2" {
			target = "?page=1"
		}
		if !strings.Contains(rec.Body.String(), target) {
			t.Fatalf("missing pagination %s", target)
		}
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest("GET", "/admin/contact/submissions/1", nil), rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(id))
	if err := h.ViewSubmission(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "A persisted inquiry for the inbox") {
		t.Fatal("stored inquiry not readable")
	}
}

func TestActivityRealRendererWithPaginationFiltersAndEmptyState(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	e := operationEcho()
	h := admin.NewActivityHandler(q, logger)
	empty := operationGET(t, e, h.List, "/admin/activity")
	if !strings.Contains(empty.Body.String(), "No Activity Recorded Yet") {
		t.Fatal("missing empty state")
	}
	for i := 0; i < 51; i++ {
		if err := q.CreateActivityLog(context.Background(), sqlc.CreateActivityLogParams{Action: "updated", ResourceType: "settings", Description: fmt.Sprintf("Saved settings fixture %02d", i)}); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"/admin/activity", "/admin/activity?page=2", "/admin/activity?action=updated&search=fixture"} {
		rec := operationGET(t, e, h.List, path)
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Saved settings fixture") {
			t.Fatalf("activity failed %s", path)
		}
	}
	filtered := operationGET(t, e, h.List, "/admin/activity?search=nomatch")
	if !strings.Contains(filtered.Body.String(), "No matching activity") || !strings.Contains(filtered.Body.String(), "clear the filters") {
		t.Fatal("filtered empty state gives no recovery")
	}
}

func TestAdminOperationsDatabaseFailureKeepsRecoveryNavigation(t *testing.T) {
	db, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	e := operationEcho()
	db.Close()
	for _, tc := range []struct {
		path    string
		handler echo.HandlerFunc
	}{
		{"/admin/contact/submissions", admin.NewAdminContactHandler(q, logger, nil).ListSubmissions},
		{"/admin/activity", admin.NewActivityHandler(q, logger).List},
		{"/admin/settings", admin.NewSettingsHandler(q, logger, nil).Edit},
	} {
		rec := operationGET(t, e, tc.handler, tc.path)
		if rec.Code != 500 || !strings.Contains(rec.Body.String(), "Try again") || !strings.Contains(rec.Body.String(), "Back to dashboard") {
			t.Fatalf("unrecoverable failure %s: %d", tc.path, rec.Code)
		}
		if strings.Contains(rec.Body.String(), "database is closed") {
			t.Fatal("internal error leaked")
		}
	}
}

func TestDefaultOGImageChooseReplaceRemoveAndPublicFallback(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	e := operationEcho()
	h := admin.NewSettingsHandler(q, logger, nil)
	for _, name := range []string{"first.png", "replacement.jpg", "brochure.pdf"} {
		mime := "image/png"
		if strings.HasSuffix(name, ".pdf") {
			mime = "application/pdf"
		}
		_, err := q.CreateMediaFile(ctx, sqlc.CreateMediaFileParams{Filename: name, OriginalFilename: name, FilePath: "/uploads/" + name, MimeType: mime, FileSize: 10})
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"/uploads/first.png", "/uploads/replacement.jpg", ""} {
		form := settingsFormFixture()
		form.Set("default_og_image", path)
		rec, c := postForm(e, "/admin/settings", form)
		if err := h.Update(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != 303 {
			t.Fatalf("save status %d", rec.Code)
		}
		settings, err := q.GetSettings(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if settings.DefaultOgImage != path {
			t.Fatal("image selection not persisted")
		}
		settingsPage := operationGET(t, e, h.Edit, "/admin/settings?tab=seo").Body.String()
		if path != "" && (!strings.Contains(settingsPage, `value="`+path+`" selected`) || !strings.Contains(settingsPage, `src="`+path+`"`)) {
			t.Fatal("persisted selection/preview missing")
		}
		if strings.Contains(settingsPage, `value="/uploads/brochure.pdf"`) {
			t.Fatal("non-image offered as image")
		}
		for _, override := range []string{"", "/uploads/page-specific.png"} {
			rr := httptest.NewRecorder()
			cc := e.NewContext(httptest.NewRequest("GET", "/contact", nil), rr)
			if err := cc.Render(200, "public/pages/contact.html", map[string]interface{}{"Settings": settings, "OGImage": override}); err != nil {
				t.Fatal(err)
			}
			want := path
			if override != "" {
				want = override
			}
			if want != "" && !strings.Contains(rr.Body.String(), `property="og:image" content="`+want+`"`) {
				t.Fatalf("public sharing image %q missing", want)
			}
			if want == "" && strings.Contains(rr.Body.String(), `property="og:image"`) {
				t.Fatal("removed fallback still advertised")
			}
		}
	}
	for _, invalid := range []string{"/uploads/brochure.pdf", "javascript:alert(1)", "/uploads/missing.png"} {
		form := settingsFormFixture()
		form.Set("default_og_image", invalid)
		_, c := postForm(e, "/admin/settings", form)
		err := h.Update(c)
		if httpErr, ok := err.(*echo.HTTPError); !ok || httpErr.Code != 400 {
			t.Fatalf("accepted invalid image %q: %v", invalid, err)
		}
	}
}

func TestSettingsVisibleControlsHaveAssociatedLabels(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	body := operationGET(t, operationEcho(), admin.NewSettingsHandler(q, logger, nil).Edit, "/admin/settings").Body.String()
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	labels := map[string]bool{}
	var controls []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "label" {
				for _, a := range n.Attr {
					if a.Key == "for" {
						labels[a.Val] = true
					}
				}
			}
			if n.Data == "input" || n.Data == "select" || n.Data == "textarea" {
				controls = append(controls, n)
			}
		}
		for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
			walk(ch)
		}
	}
	walk(doc)
	for _, n := range controls {
		id, name, kind := "", "", ""
		for _, a := range n.Attr {
			switch a.Key {
			case "id":
				id = a.Val
			case "name":
				name = a.Val
			case "type":
				kind = a.Val
			}
		}
		if name == "" || kind == "hidden" {
			continue
		}
		if id == "" || !labels[id] {
			t.Errorf("%s lacks associated visible label", name)
		}
	}
}

func TestLoginRecoveryExplainsAdministratorRoute(t *testing.T) {
	e := operationEcho()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest("GET", "/admin/login", nil), rec)
	if err := c.Render(200, "admin/pages/login.html", map[string]interface{}{}); err != nil {
		t.Fatal(err)
	}
	body := rec.Body.String()
	for _, want := range []string{"Contact your site administrator", "server operator with deployment access"} {
		if !strings.Contains(body, want) {
			t.Errorf("missing recovery guidance %q", want)
		}
	}
	if strings.Contains(body, "./server reset-password") {
		t.Fatal("unsupported CLI recovery advertised")
	}
}
