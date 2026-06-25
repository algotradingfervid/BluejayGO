package e2e_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestProductSearch_ResultImageNotCropped verifies that product cards in the HTMX
// product search results fragment display the FULL uploaded product image rather than
// center-cropping/zooming it. The card image must use object-contain (whole image
// shown inside the aspect-[4/3] box) not object-cover (which crops non-matching ratios).
//
// Uses the REAL renderer via the production /products/search wiring (pattern from
// 38_product_search_view_details_link_test.go).
func TestProductSearch_ResultImageNotCropped(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Detectors", Slug: "detectors", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "BJ-D400", Slug: "bj-d400", Name: "BJ-D400 Detector",
		Description:  "Detects 400 series particles",
		CategoryID:   cat.ID,
		Status:       "published",
		PrimaryImage: sql.NullString{String: "/uploads/products/detector-wide.png", Valid: true},
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

	req := httptest.NewRequest(http.MethodGet, "/products/search?q=400", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /products/search, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	if !strings.Contains(body, `/uploads/products/detector-wide.png`) {
		t.Fatalf("expected the search-result card image to render, but it was absent; body: %s", body)
	}

	cardImgRe := regexp.MustCompile(`<img[^>]*src="/uploads/products/detector-wide\.png"[^>]*>`)
	cardImg := cardImgRe.FindString(body)
	if cardImg == "" {
		t.Fatalf("could not locate the search-result card image element; body: %s", body)
	}

	if strings.Contains(cardImg, "object-cover") {
		t.Errorf("search-result card image uses object-cover, which crops/zooms non-square images: %s", cardImg)
	}
	if !strings.Contains(cardImg, "object-contain") {
		t.Errorf("expected search-result card image to use object-contain so the full image displays, got: %s", cardImg)
	}
}
