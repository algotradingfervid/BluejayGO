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

// TestProductDetail_MainImageNotCropped verifies that the public product detail page
// displays the FULL uploaded product image rather than center-cropping/zooming it to a
// fixed square box. Product images come in varying aspect ratios (portrait, landscape,
// square); the main gallery image must use object-contain so the whole image shows,
// not object-cover which crops anything that isn't square.
//
// The shared setupApp uses a stub renderer that cannot see template bugs, so this test
// builds a local Echo with the REAL renderer (mirroring production) and asserts on the
// rendered product detail markup.
func TestProductDetail_MainImageNotCropped(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Sensors",
		Slug:        "sensors",
		Description: "desc",
		Icon:        "icon",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:          "SENS-100",
		Slug:         "sensor-x100",
		Name:         "Sensor X100",
		Description:  "A tall portrait product image must not be cropped.",
		CategoryID:   cat.ID,
		Status:       "published",
		PrimaryImage: sql.NullString{String: "/uploads/products/portrait.png", Valid: true},
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	productSvc := services.NewProductService(queries)
	cache := services.NewCache()
	h := publicHandlers.NewProductsHandler(queries, logger, productSvc, cache)
	e.GET("/products/:category/:slug", h.ProductDetail)

	req := httptest.NewRequest(http.MethodGet, "/products/sensors/sensor-x100", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	// Sanity: the main product image must actually render.
	if !strings.Contains(body, `/uploads/products/portrait.png`) {
		t.Fatalf("expected the product image to render, but it was absent; body: %s", body)
	}

	// Isolate the main image element (id="main-image"), which is the primary product
	// image display. It must NOT use object-cover (which center-crops non-square images).
	mainImgRe := regexp.MustCompile(`<img[^>]*id="main-image"[^>]*>`)
	mainImg := mainImgRe.FindString(body)
	if mainImg == "" {
		t.Fatalf("could not locate the main product image element in rendered page; body: %s", body)
	}

	if strings.Contains(mainImg, "object-cover") {
		t.Errorf("main product image uses object-cover, which crops/zooms non-square images and hides the full image: %s", mainImg)
	}
	if !strings.Contains(mainImg, "object-contain") {
		t.Errorf("expected main product image to use object-contain so the full image displays responsively, got: %s", mainImg)
	}
}
