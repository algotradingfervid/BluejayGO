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
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestHome_PopularProductsImageNotCropped verifies that the homepage
// "Discover our most popular interactive solutions" (Featured Products) card images
// display the FULL uploaded product image rather than center-cropping/zooming it. This
// is the section in the user's screenshot. Product images have varying aspect ratios;
// the card image must use object-contain so the whole image shows, not object-cover
// which crops anything that does not match the aspect-video box.
//
// Uses the REAL renderer (the shared setupApp stub renderer cannot see template bugs).
func TestHome_PopularProductsImageNotCropped(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Interactive",
		Slug:        "interactive",
		Description: "desc",
		Icon:        "icon",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:           "IFP-100",
		Slug:          "interactive-panel",
		Name:          "Interactive Panel",
		Description:   "A non-square product image must not be cropped on the homepage.",
		CategoryID:    cat.ID,
		Status:        "published",
		IsFeatured:    true,
		FeaturedOrder: sql.NullInt64{Int64: 1, Valid: true},
		PrimaryImage:  sql.NullString{String: "/uploads/products/panel-tall.png", Valid: true},
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	h := publicHandlers.NewHomeHandler(queries, logger)
	e.GET("/", h.ShowHomePage)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	if !strings.Contains(body, `/uploads/products/panel-tall.png`) {
		t.Fatalf("expected the popular-products card image to render, but it was absent; body: %s", body)
	}

	cardImgRe := regexp.MustCompile(`<img[^>]*src="/uploads/products/panel-tall\.png"[^>]*>`)
	cardImg := cardImgRe.FindString(body)
	if cardImg == "" {
		t.Fatalf("could not locate the homepage popular-products card image element; body: %s", body)
	}

	if strings.Contains(cardImg, "object-cover") {
		t.Errorf("homepage popular-products card image uses object-cover, which crops/zooms non-square images and hides the full image: %s", cardImg)
	}
	if !strings.Contains(cardImg, "object-contain") {
		t.Errorf("expected homepage popular-products card image to use object-contain so the full image displays, got: %s", cardImg)
	}
}
