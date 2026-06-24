package e2e_test

import (
	"context"
	"database/sql"
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

// TestHomepageHeroCarousel_RendersAllActiveHeroesAsSlides verifies that the public
// homepage renders EVERY active hero as a slide in a rotating carousel — not just the
// single hero returned by GetActiveHero.
//
// The shared setupApp uses a stub renderer, so this test builds a local Echo with the
// REAL renderer (mirroring production) and asserts on the rendered home page markup.
func TestHomepageHeroCarousel_RendersAllActiveHeroesAsSlides(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	// Seed three ACTIVE heroes with distinct headlines and explicit display order.
	heroes := []struct {
		headline string
		order    int64
	}{
		{"First Active Hero Slide", 1},
		{"Second Active Hero Slide", 2},
		{"Third Active Hero Slide", 3},
	}
	for _, hh := range heroes {
		_, err := queries.CreateHero(ctx, sqlc.CreateHeroParams{
			Headline:        hh.headline,
			Subheadline:     "Subheadline for " + hh.headline,
			PrimaryCtaText:  "Explore",
			PrimaryCtaUrl:   "/products",
			BackgroundImage: sql.NullString{String: "/uploads/products/bj-ifp65.jpg", Valid: true},
			IsActive:        1,
			DisplayOrder:    hh.order,
		})
		if err != nil {
			t.Fatalf("create hero %q: %v", hh.headline, err)
		}
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
		t.Fatalf("expected 200 from /, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	// All three active heroes must be rendered as slides.
	for _, hh := range heroes {
		if !strings.Contains(body, hh.headline) {
			t.Errorf("expected active hero %q to render on the homepage, but it was absent", hh.headline)
		}
	}

	// The hero section must be a carousel with one slide per active hero.
	if !strings.Contains(body, "data-hero-carousel") {
		t.Errorf("expected a hero carousel container (data-hero-carousel) on the homepage")
	}
	// Count the attribute-with-value form so we match rendered elements only, not the
	// inline JS selector ("[data-hero-slide]") or the dots container ("data-hero-dots").
	if got := strings.Count(body, `data-hero-slide="`); got != 3 {
		t.Errorf("expected 3 hero slides (data-hero-slide), got %d", got)
	}

	// With more than one slide, navigation controls (dots) must be present.
	if got := strings.Count(body, `data-hero-dot="`); got != 3 {
		t.Errorf("expected 3 carousel dots (data-hero-dot), got %d", got)
	}
}
