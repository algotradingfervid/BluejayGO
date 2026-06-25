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

// TestProductCategory_CardImageNotCropped verifies that product CARDS on the public
// category listing page display the FULL uploaded product image rather than center-
// cropping/zooming it inside the fixed square box. Product images come in varying
// aspect ratios; the card image must use object-contain so the whole image shows,
// not object-cover which crops anything that isn't square.
//
// The shared setupApp uses a stub renderer that cannot see template bugs, so this test
// builds a local Echo with the REAL renderer (mirroring production) and asserts on the
// rendered markup.
func TestProductCategory_CardImageNotCropped(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Microphones",
		Slug:        "microphones",
		Description: "desc",
		Icon:        "icon",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:          "MIC-100",
		Slug:         "studio-mic",
		Name:         "Studio Mic",
		Description:  "A tall portrait product image must not be cropped on the card.",
		CategoryID:   cat.ID,
		Status:       "published",
		PrimaryImage: sql.NullString{String: "/uploads/products/mic-portrait.png", Valid: true},
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
	e.GET("/products/:category", h.ProductsByCategory)

	req := httptest.NewRequest(http.MethodGet, "/products/microphones", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	// Sanity: the product card image must actually render.
	if !strings.Contains(body, `/uploads/products/mic-portrait.png`) {
		t.Fatalf("expected the product card image to render, but it was absent; body: %s", body)
	}

	// Isolate the product card <img> (the one pointing at the uploaded product image).
	cardImgRe := regexp.MustCompile(`<img[^>]*src="/uploads/products/mic-portrait\.png"[^>]*>`)
	cardImg := cardImgRe.FindString(body)
	if cardImg == "" {
		t.Fatalf("could not locate the product card image element in rendered page; body: %s", body)
	}

	if strings.Contains(cardImg, "object-cover") {
		t.Errorf("product card image uses object-cover, which crops/zooms non-square images and hides the full image: %s", cardImg)
	}
	if !strings.Contains(cardImg, "object-contain") {
		t.Errorf("expected product card image to use object-contain so the full image displays responsively, got: %s", cardImg)
	}
}

// TestProductDetail_GalleryThumbnailsNotCropped verifies the gallery THUMBNAIL images
// on the product detail page also show the full uploaded image (object-contain), not a
// center-cropped object-cover. The main image was already fixed; the thumbnails were
// missed and still crop non-square images.
func TestProductDetail_GalleryThumbnailsNotCropped(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Microphones",
		Slug:        "microphones",
		Description: "desc",
		Icon:        "icon",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	prod, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:          "MIC-200",
		Slug:         "studio-mic-2",
		Name:         "Studio Mic 2",
		Description:  "Gallery thumbnails must show full image.",
		CategoryID:   cat.ID,
		Status:       "published",
		PrimaryImage: sql.NullString{String: "/uploads/products/mic-main.png", Valid: true},
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	_, err = queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID:    prod.ID,
		ImagePath:    "/uploads/products/mic-thumb.png",
		AltText:      sql.NullString{String: "thumb", Valid: true},
		DisplayOrder: 1,
		IsThumbnail:  false,
	})
	if err != nil {
		t.Fatalf("create product image: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	productSvc := services.NewProductService(queries)
	cache := services.NewCache()
	h := publicHandlers.NewProductsHandler(queries, logger, productSvc, cache)
	e.GET("/products/:category/:slug", h.ProductDetail)

	req := httptest.NewRequest(http.MethodGet, "/products/microphones/studio-mic-2", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	if !strings.Contains(body, `/uploads/products/mic-thumb.png`) {
		t.Fatalf("expected the gallery thumbnail image to render, but it was absent; body: %s", body)
	}

	thumbImgRe := regexp.MustCompile(`<img[^>]*src="/uploads/products/mic-thumb\.png"[^>]*>`)
	thumbImg := thumbImgRe.FindString(body)
	if thumbImg == "" {
		t.Fatalf("could not locate the gallery thumbnail image element; body: %s", body)
	}

	if strings.Contains(thumbImg, "object-cover") {
		t.Errorf("gallery thumbnail uses object-cover, which crops/zooms non-square images: %s", thumbImg)
	}
	if !strings.Contains(thumbImg, "object-contain") {
		t.Errorf("expected gallery thumbnail to use object-contain so the full image displays, got: %s", thumbImg)
	}
}
