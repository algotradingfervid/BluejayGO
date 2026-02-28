package e2e_test

import (
	"context"
	"database/sql"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestSitemap_Loads(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Header().Get("Content-Type"), "xml") {
		t.Errorf("expected XML content type, got %s", rec.Header().Get("Content-Type"))
	}
}

func TestSitemap_ValidXML(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var urlset struct {
		XMLName xml.Name `xml:"urlset"`
	}
	if err := xml.Unmarshal(rec.Body.Bytes(), &urlset); err != nil {
		t.Errorf("invalid XML: %v", err)
	}
}

func TestSitemap_ContainsStaticPages(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	body := rec.Body.String()

	staticPages := []string{
		"<loc>https://bluejaylabs.com/</loc>",
		"<loc>https://bluejaylabs.com/products</loc>",
		"<loc>https://bluejaylabs.com/solutions</loc>",
		"<loc>https://bluejaylabs.com/blog</loc>",
		"<loc>https://bluejaylabs.com/contact</loc>",
	}

	for _, page := range staticPages {
		if !strings.Contains(body, page) {
			t.Errorf("sitemap missing static page: %s", page)
		}
	}
}

func TestSitemap_ContainsDynamicContent(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		Icon:             "icon",
		ShortDescription: "Test",
		HeroImageUrl:     sql.NullString{},
		HeroTitle:        sql.NullString{},
		HeroDescription:  sql.NullString{},
		OverviewContent:  sql.NullString{},
		MetaDescription:  sql.NullString{},
		ReferenceCode:    sql.NullString{},
		IsPublished:      sql.NullBool{Bool: true, Valid: true},
		DisplayOrder:     sql.NullInt64{Int64: 1, Valid: true},
	})

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	expectedURL := "<loc>https://bluejaylabs.com/solutions/" + sol.Slug + "</loc>"
	if !strings.Contains(rec.Body.String(), expectedURL) {
		t.Errorf("sitemap missing solution URL: %s", expectedURL)
	}
}

func TestRobotsTxt_Loads(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "text/plain; charset=UTF-8" {
		t.Errorf("expected text/plain content type, got %s", rec.Header().Get("Content-Type"))
	}
}

func TestRobotsTxt_Content(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	body := rec.Body.String()

	requiredDirectives := []string{
		"User-agent: *",
		"Disallow: /admin/",
		"Sitemap:",
	}

	for _, directive := range requiredDirectives {
		if !strings.Contains(body, directive) {
			t.Errorf("robots.txt missing directive: %s", directive)
		}
	}
}
