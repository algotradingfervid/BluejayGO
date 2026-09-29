package admin_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestCaseStudyPublicationCreateUpdateAndCachedVisibility(t *testing.T) {
	db, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	result, err := db.Exec(`INSERT INTO industries(name,slug,icon,description) VALUES('Publication','publication','test','Fixture')`)
	if err != nil {
		t.Fatal(err)
	}
	industryID, _ := result.LastInsertId()
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	cache := services.NewCache()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("X-Test-Admin") == "yes" {
				c.Set("session", &middleware.Session{UserID: 1})
			}
			return next(c)
		}
	})
	ah := admin.NewCaseStudiesHandler(q, logger, cache)
	ph := public.NewCaseStudiesHandler(q, logger, cache)
	e.POST("/admin/case-studies", ah.Create)
	e.POST("/admin/case-studies/:id", ah.Update)
	e.GET("/admin/case-studies/:id/edit", ah.Edit)
	e.GET("/case-studies/:slug", ph.CaseStudyDetail)
	form := url.Values{"title": {"Publication fixture"}, "slug": {"publication-fixture"}, "client_name": {"Synthetic client"}, "industry_id": {fmt.Sprint(industryID)}, "summary": {"Keep summary"}, "challenge_title": {"Challenge"}, "challenge_content": {"<p>Keep <strong>challenge</strong>.</p>"}, "solution_title": {"Solution"}, "solution_content": {"<p>Keep solution.</p>"}, "outcome_title": {"Outcome"}, "outcome_content": {"<p>Keep outcome.</p>"}, "is_published": {"1"}}
	post := func(path string) {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != 303 {
			t.Fatalf("POST %s: %d %s", path, rec.Code, rec.Body.String())
		}
	}
	get := func(path string, adminSession bool, want int) string {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if adminSession {
			req.Header.Set("X-Test-Admin", "yes")
		}
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Fatalf("GET %s: got%d want%d %s", path, rec.Code, want, rec.Body.String())
		}
		return rec.Body.String()
	}
	// Creation and updates use the real select values, while legacy checkbox requests remain valid.
	post("/admin/case-studies")
	var id int64
	if err = db.QueryRow(`SELECT id FROM case_studies WHERE slug='publication-fixture'`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	updatePath := fmt.Sprintf("/admin/case-studies/%d", id)
	assertState := func(published bool) {
		t.Helper()
		item, err := q.AdminGetCaseStudy(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}
		wantValue, wantStatus, publicCode := "0", "draft", 404
		if published {
			wantValue, wantStatus, publicCode = "1", "published", 200
		}
		if (item.IsPublished == 1) != published {
			t.Fatalf("saved publication state %d", item.IsPublished)
		}
		body := get(updatePath+"/edit", false, 200)
		if !strings.Contains(body, `value="`+wantValue+`" selected`) {
			t.Fatalf("edit does not restore status %s", wantValue)
		}
		if item.Summary != "Keep summary" || item.ChallengeContent != "<p>Keep <strong>challenge</strong>.</p>" {
			t.Fatal("status transition changed story fields")
		}
		get("/case-studies/publication-fixture", false, publicCode)
		preview := get("/case-studies/publication-fixture?preview=true", true, 200)
		if !strings.Contains(preview, "Saved status: "+wantStatus) {
			t.Fatal("preview status disagrees with saved state")
		}
		if !published {
			get("/case-studies/publication-fixture?preview=true", false, 404)
		}
	}
	assertState(true) // Also primes the public page cache before unpublishing.
	for _, state := range []struct {
		value     string
		published bool
	}{{"0", false}, {"1", true}, {"", false}, {"on", true}, {"0", false}} {
		if state.value == "" {
			form.Del("is_published")
		} else {
			form.Set("is_published", state.value)
		}
		post(updatePath)
		assertState(state.published)
	}
}
