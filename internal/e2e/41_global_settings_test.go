package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGlobalSettingsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/settings", strings.NewReader(url.Values{
		"active_tab":          {"general"},
		"site_name":           {"Bluejay Labs Test"},
		"site_tagline":        {"Innovation in Testing"},
		"contact_email":       {"test@bluejay.com"},
		"contact_phone":       {"555-1234"},
		"address":             {"123 Test St"},
		"business_hours":      {"Mon-Fri 9-5"},
		"meta_description":    {"Test meta description"},
		"meta_keywords":       {"test, keywords"},
		"google_analytics_id": {"UA-12345"},
		"social_facebook":     {"https://facebook.com/bluejay"},
		"social_twitter":      {"https://twitter.com/bluejay"},
		"social_linkedin":     {"https://linkedin.com/company/bluejay"},
		"social_instagram":    {"https://instagram.com/bluejay"},
		"social_youtube":      {"https://youtube.com/bluejay"},
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
	if !strings.Contains(loc, "tab=general") {
		t.Errorf("expected redirect with tab=general, got %q", loc)
	}

	settings, _ := queries.GetSettings(ctx)
	if settings.SiteName != "Bluejay Labs Test" {
		t.Errorf("expected site_name 'Bluejay Labs Test', got %q", settings.SiteName)
	}
	if settings.ContactEmail != "test@bluejay.com" {
		t.Errorf("expected contact_email 'test@bluejay.com', got %q", settings.ContactEmail)
	}
	if settings.GoogleAnalyticsID != "UA-12345" {
		t.Errorf("expected google_analytics_id 'UA-12345', got %q", settings.GoogleAnalyticsID)
	}
}

func TestGlobalSettingsTabPersistence_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	tabs := []string{"general", "contact", "social", "seo"}
	for _, tab := range tabs {
		req := httptest.NewRequest(http.MethodPost, "/admin/settings", strings.NewReader(url.Values{
			"active_tab":       {tab},
			"site_name":        {"Test"},
			"contact_email":    {"test@example.com"},
			"meta_description": {"Test"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		loc := rec.Header().Get("Location")
		if !strings.Contains(loc, "tab="+tab) {
			t.Errorf("expected redirect with tab=%s, got %q", tab, loc)
		}
	}
}

func TestGlobalSettingsEmptyFields_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/settings", strings.NewReader(url.Values{
		"site_name":           {"Required Field"},
		"site_tagline":        {""},
		"social_facebook":     {""},
		"google_analytics_id": {""},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("settings route not found")
	}
}
