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
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestProductSearch_ViewDetailsLinksToDetailPage verifies that each product card
// rendered by the HTMX product search fragment links to the product detail page
// at /products/{category-slug}/{product-slug}.
//
// The shared setupApp uses a stub renderer (which would hide rendered markup), so
// this test builds a local Echo with the REAL renderer, mirroring production wiring
// for /products/search. TestMain has already chdir'd to the project root, so the
// "templates" path resolves correctly.
func TestProductSearch_ViewDetailsLinksToDetailPage(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	// Seed a published product in a category with known slugs we can assert on.
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Detectors", Slug: "detectors", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "BJ-D300", Slug: "bj-d300", Name: "BJ-D300 Detector",
		Description: "Detects 300 series particles", CategoryID: cat.ID, Status: "published",
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	productSvc := services.NewProductService(queries)
	appCache := services.NewCache()
	h := publicHandlers.NewProductsHandler(queries, logger, productSvc, appCache)
	e.GET("/products/search", h.ProductSearch)

	req := httptest.NewRequest(http.MethodGet, "/products/search?q=300", nil)
	req.Header.Set("HX-Request", "true") // exercise the HTMX partial path
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /products/search, got %d; body: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()

	// Sanity: the matching product should appear in the results.
	if !strings.Contains(body, "BJ-D300 Detector") {
		t.Fatalf("expected product 'BJ-D300 Detector' in search results; body:\n%s", body)
	}

	// The "View Details" action must be a real link to the detail page built from
	// the category slug + product slug.
	wantHref := `href="/products/detectors/bj-d300"`
	if !strings.Contains(body, wantHref) {
		t.Errorf("expected View Details to link to detail page (%s), but it was absent; body:\n%s", wantHref, body)
	}
}
