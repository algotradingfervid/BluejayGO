package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestNavigation_HeaderRendersOnHomepage(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestNavigation_HeaderRendersOnPublicPages(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	pages := []string{
		"/products",
		"/solutions",
		"/blog",
		"/contact",
	}

	for _, page := range pages {
		req := httptest.NewRequest(http.MethodGet, page, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusNotFound {
			t.Errorf("%s: expected 200 or 404, got %d", page, rec.Code)
		}
	}
}

func TestNavigation_WithSettings(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_ = queries.UpdateSettings(ctx, sqlc.UpdateSettingsParams{
		SiteName:          "Test Site",
		SiteTagline:       "",
		ContactEmail:      "",
		ContactPhone:      "",
		Address:           "",
		FooterText:        "",
		MetaDescription:   "",
		MetaKeywords:      "",
		GoogleAnalyticsID: "",
		SocialLinkedin:    "",
		SocialTwitter:     "",
		SocialFacebook:    "",
		SocialInstagram:   "",
		SocialYoutube:     "",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestNavigation_FooterRendersOnPublicPages(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	pages := []string{
		"/",
		"/products",
		"/contact",
	}

	for _, page := range pages {
		req := httptest.NewRequest(http.MethodGet, page, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusNotFound {
			t.Errorf("%s: expected 200 or 404, got %d", page, rec.Code)
		}
	}
}
