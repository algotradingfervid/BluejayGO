package e2e_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// renderHome renders the public homepage with the REAL renderer (the shared setupApp
// uses a stub renderer, so we build a local Echo to assert on the actual markup).
func renderHome(t *testing.T, queries *sqlc.Queries) string {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	h := publicHandlers.NewHomeHandler(queries, logger)
	e.GET("/", h.ShowHomePage)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /, got %d; body: %s", rec.Code, rec.Body.String())
	}
	return rec.Body.String()
}

// TestHome_SectionHeadingsAreSeededAndEditable proves the end-to-end requirement:
//
//  1. After migrations (which now include 037 seeding the home page_sections), the
//     homepage renders the SEEDED section headings from page_sections, not just the
//     template fallbacks. We assert seeded headings render and that each seeded section
//     has an editable row (so the admin /admin/page-sections editor can edit them).
//  2. Editing a section heading/subheading via the existing UpdatePageSection query
//     (the same path the admin POST /admin/page-sections/:id uses) is reflected on the
//     homepage, and the old default for that section disappears.
//
// We target the Solutions section for the edit check because it is backed by a seeded
// row and renders once a published solution exists.
func TestHome_SectionHeadingsAreSeededAndEditable(t *testing.T) {
	db, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Seed content so the relevant homepage sections render.
	// Published solution -> Solutions section renders.
	if _, err := db.ExecContext(ctx,
		`INSERT INTO solutions (title, slug, icon, short_description, is_published, display_order)
		 VALUES (?, ?, ?, ?, 1, 1)`,
		"Education", "education", "school", "Solutions for education"); err != nil {
		t.Fatalf("seed solution: %v", err)
	}
	// Active testimonial -> Testimonials section renders.
	if _, err := db.ExecContext(ctx,
		`INSERT INTO homepage_testimonials (quote, author_name, rating, display_order, is_active)
		 VALUES (?, ?, 5, 1, 1)`,
		"Outstanding product.", "Jane Doe"); err != nil {
		t.Fatalf("seed testimonial: %v", err)
	}

	// ---- Part 1: seeded page_section headings render on the homepage ----
	body := renderHome(t, queries)

	// Seeded headings from migration 037 / 013 must render (they come from page_sections).
	for _, want := range []string{
		"Solutions By Industry", // 037 updated solutions_section
		"What Our Clients Say",  // 037 seeded testimonials_section
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected seeded section heading %q to render on the homepage, but it was absent", want)
		}
	}

	// Every homepage section must now have an editable page_sections row (the admin
	// editor only edits existing rows, so missing rows = not editable).
	for _, key := range []string{
		"products_section", "solutions_section", "stats_section",
		"testimonials_section", "partners_section", "blog_section",
	} {
		if _, err := queries.GetPageSection(ctx, sqlc.GetPageSectionParams{
			PageKey: "home", SectionKey: key,
		}); err != nil {
			t.Errorf("expected an editable home page_section row for %q, but GetPageSection failed: %v", key, err)
		}
	}

	// ---- Part 2: editing a section heading is reflected on the homepage ----
	const customHeading = "ZZZ Custom Solutions Heading"
	const customSubheading = "ZZZ Custom Solutions Subheading"
	const oldDefaultHeading = "Solutions By Industry"

	section, err := queries.GetPageSection(ctx, sqlc.GetPageSectionParams{
		PageKey: "home", SectionKey: "solutions_section",
	})
	if err != nil {
		t.Fatalf("get solutions_section: %v", err)
	}

	// Same write the admin POST /admin/page-sections/:id performs.
	if err := queries.UpdatePageSection(ctx, sqlc.UpdatePageSectionParams{
		Heading:    customHeading,
		Subheading: customSubheading,
		IsActive:   true,
		ID:         section.ID,
	}); err != nil {
		t.Fatalf("update solutions_section: %v", err)
	}

	body = renderHome(t, queries)

	if !strings.Contains(body, customHeading) {
		t.Errorf("expected edited section heading %q to render on the homepage, but it was absent", customHeading)
	}
	if !strings.Contains(body, customSubheading) {
		t.Errorf("expected edited section subheading %q to render on the homepage, but it was absent", customSubheading)
	}
	if strings.Contains(body, oldDefaultHeading) {
		t.Errorf("expected old default heading %q NOT to render once a custom value is set", oldDefaultHeading)
	}
}
