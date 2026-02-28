package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHomepageSettings_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/homepage/settings", strings.NewReader(url.Values{
		"homepage_show_heroes":        {"on"},
		"homepage_show_stats":         {"on"},
		"homepage_show_testimonials":  {"on"},
		"homepage_show_cta":           {"on"},
		"homepage_max_heroes":         {"5"},
		"homepage_max_stats":          {"6"},
		"homepage_max_testimonials":   {"3"},
		"homepage_hero_autoplay":      {"on"},
		"homepage_hero_interval":      {"8"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "saved=1") {
		t.Errorf("expected redirect with saved=1, got %q", loc)
	}

	settings, _ := queries.GetSettings(ctx)
	if settings.HomepageShowHeroes != 1 {
		t.Errorf("expected homepage_show_heroes 1, got %d", settings.HomepageShowHeroes)
	}
	if settings.HomepageMaxStats != 6 {
		t.Errorf("expected homepage_max_stats 6, got %d", settings.HomepageMaxStats)
	}
	if settings.HomepageHeroInterval != 8 {
		t.Errorf("expected homepage_hero_interval 8, got %d", settings.HomepageHeroInterval)
	}
}

func TestHomepageSettingsToggleAllOff_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/homepage/settings", strings.NewReader(url.Values{
		"homepage_show_heroes":       {""},
		"homepage_show_stats":        {""},
		"homepage_show_testimonials": {""},
		"homepage_show_cta":          {""},
		"homepage_max_heroes":        {"0"},
		"homepage_max_stats":         {"0"},
		"homepage_max_testimonials":  {"0"},
		"homepage_hero_autoplay":     {""},
		"homepage_hero_interval":     {"5"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	settings, _ := queries.GetSettings(ctx)
	if settings.HomepageShowHeroes != 0 {
		t.Errorf("expected homepage_show_heroes 0, got %d", settings.HomepageShowHeroes)
	}
	if settings.HomepageShowStats != 0 {
		t.Errorf("expected homepage_show_stats 0, got %d", settings.HomepageShowStats)
	}
	if settings.HomepageHeroAutoplay != 0 {
		t.Errorf("expected homepage_hero_autoplay 0, got %d", settings.HomepageHeroAutoplay)
	}
}

func TestHomepageSettingsEmptyFields_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/homepage/settings", strings.NewReader(url.Values{
		"homepage_max_heroes": {""},
		"homepage_max_stats":  {""},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}
}
