package e2e_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// captureHandler is a minimal slog.Handler that records every log message it sees,
// allowing tests to assert on which log lines were (or were not) emitted.
type captureHandler struct {
	mu   *sync.Mutex
	msgs *[]string
}

func (h captureHandler) Enabled(context.Context, slog.Level) bool { return true }
func (h captureHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	*h.msgs = append(*h.msgs, r.Message)
	return nil
}
func (h captureHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h captureHandler) WithGroup(string) slog.Handler      { return h }

// TestSearch_OmitsCaseStudies verifies that the public search no longer queries the
// case studies FTS index (so case studies never appear in search results), while
// product search continues to work.
//
// The shared setupApp uses a stub renderer (which would hide rendered results), so
// this test builds a local Echo with the REAL renderer and the SettingsLoader
// middleware, mirroring production wiring for /search.
func TestSearch_OmitsCaseStudies(t *testing.T) {
	db, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	var mu sync.Mutex
	var logs []string
	logger := slog.New(captureHandler{mu: &mu, msgs: &logs})

	// Seed a published product whose name contains a distinctive token we can search.
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Detectors", Slug: "detectors", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "ZED-1", Slug: "zorptron-detector", Name: "Zorptron Detector",
		Description: "Detects zorptron particles", CategoryID: cat.ID, Status: "published",
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	// Seed a published case study that ALSO matches the search token, to prove it is
	// intentionally excluded from search results.
	ind, err := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Defense", Slug: "defense"})
	if err != nil {
		t.Fatalf("create industry: %v", err)
	}
	_, err = queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug: "zorptron-case", Title: "Zorptron Case Study", ClientName: "Acme",
		IndustryID: ind.ID, Summary: "s", ChallengeContent: "c", SolutionContent: "x", OutcomeContent: "o",
		IsPublished: 1,
	})
	if err != nil {
		t.Fatalf("create case study: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")
	searchHandler := publicHandlers.NewSearchHandler(db, logger)
	pub := e.Group("", customMiddleware.SettingsLoader(queries))
	pub.GET("/search", searchHandler.SearchPage)

	req := httptest.NewRequest(http.MethodGet, "/search?q=Zorptron", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /search, got %d; body: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()

	// Product search must still work.
	if !strings.Contains(body, "Zorptron Detector") {
		t.Errorf("expected product 'Zorptron Detector' in search results; body:\n%s", body)
	}

	// Case studies must not appear as a result type, nor link.
	if strings.Contains(body, "Case Study") {
		t.Errorf("search results should NOT include 'Case Study' type; body:\n%s", body)
	}
	if strings.Contains(body, "/case-studies/") {
		t.Errorf("search results should NOT link to /case-studies/; body:\n%s", body)
	}

	// The case-studies FTS query must no longer be attempted. Previously it errored
	// on every search (the case_studies table has no 'status' column), logging
	// "case_studies fts query failed". After removal, that log must be gone.
	mu.Lock()
	defer mu.Unlock()
	for _, m := range logs {
		if strings.Contains(m, "case_studies fts query failed") {
			t.Errorf("search should not query case_studies FTS, but it logged: %q", m)
		}
	}
}
