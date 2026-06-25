package e2e_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHeaderSettingsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/header", strings.NewReader(url.Values{
		"header_logo_path":       {"/images/logo.png"},
		"header_logo_alt":        {"Bluejay Logo"},
		"header_cta_enabled":     {"on"},
		"header_cta_text":        {"Get Started"},
		"header_cta_url":         {"/signup"},
		"header_cta_style":       {"primary"},
		"header_show_phone":      {"on"},
		"header_show_email":      {"on"},
		"header_show_social":     {"on"},
		"header_social_style":    {"icons"},
		"show_nav_products":      {"on"},
		"show_nav_solutions":     {"on"},
		"show_nav_case_studies":  {"on"},
		"show_nav_about":         {"on"},
		"show_nav_blog":          {"on"},
		"show_nav_whitepapers":   {"on"},
		"show_nav_partners":      {"on"},
		"show_nav_contact":       {"on"},
		"nav_label_products":     {"Our Products"},
		"nav_label_solutions":    {"Solutions"},
		"nav_label_case_studies": {"Case Studies"},
		"nav_label_about":        {"About Us"},
		"nav_label_blog":         {"Blog"},
		"nav_label_whitepapers":  {"Resources"},
		"nav_label_partners":     {"Partners"},
		"nav_label_contact":      {"Contact"},
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
	if settings.HeaderLogoPath != "/images/logo.png" {
		t.Errorf("expected header_logo_path '/images/logo.png', got %q", settings.HeaderLogoPath)
	}
	if !settings.HeaderCtaEnabled {
		t.Error("expected header_cta_enabled true")
	}
	if settings.NavLabelProducts != "Our Products" {
		t.Errorf("expected nav_label_products 'Our Products', got %q", settings.NavLabelProducts)
	}
}

func TestHeaderSettingsToggleAllNav_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/header", strings.NewReader(url.Values{
		"header_logo_path":      {"/logo.png"},
		"show_nav_products":     {""},
		"show_nav_solutions":    {""},
		"show_nav_case_studies": {""},
		"show_nav_about":        {""},
		"show_nav_blog":         {""},
		"show_nav_whitepapers":  {""},
		"show_nav_partners":     {""},
		"show_nav_contact":      {""},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	settings, _ := queries.GetSettings(ctx)
	if settings.ShowNavProducts {
		t.Error("expected show_nav_products false")
	}
	if settings.ShowNavAbout {
		t.Error("expected show_nav_about false")
	}
}

func TestHeaderCTAStyles_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	styles := []string{"primary", "secondary"}
	for _, style := range styles {
		req := httptest.NewRequest(http.MethodPost, "/admin/header", strings.NewReader(url.Values{
			"header_logo_path":   {"/logo.png"},
			"header_cta_enabled": {"on"},
			"header_cta_text":    {"Click"},
			"header_cta_url":     {"/"},
			"header_cta_style":   {style},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("update with style %q: expected 303, got %d", style, rec.Code)
		}

		settings, _ := queries.GetSettings(ctx)
		if settings.HeaderCtaStyle != style {
			t.Errorf("expected header_cta_style %q, got %q", style, settings.HeaderCtaStyle)
		}
	}
}

// TestHeaderLogoUpload_E2E verifies that uploading a logo file via the header
// settings form persists the file to disk and stores its public /uploads/branding/
// path as the header logo, overriding any pasted text path.
func TestHeaderLogoUpload_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// Build a multipart form with an uploaded logo file plus a (different) pasted
	// path; the uploaded file must win.
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("header_logo_path", "/should/be/overridden.png"); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteField("header_logo_alt", "Bluejay Logo"); err != nil {
		t.Fatal(err)
	}
	part, err := w.CreateFormFile("header_logo_file", "my-logo.svg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/header", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("upload: expected 303, got %d (body: %s)", rec.Code, rec.Body.String())
	}

	settings, _ := queries.GetSettings(ctx)
	if !strings.HasPrefix(settings.HeaderLogoPath, "/uploads/branding/") {
		t.Errorf("expected uploaded logo path under /uploads/branding/, got %q", settings.HeaderLogoPath)
	}
	if !strings.HasSuffix(settings.HeaderLogoPath, "_my-logo.svg") {
		t.Errorf("expected saved path to keep original filename, got %q", settings.HeaderLogoPath)
	}
	if settings.HeaderLogoAlt != "Bluejay Logo" {
		t.Errorf("expected header_logo_alt 'Bluejay Logo', got %q", settings.HeaderLogoAlt)
	}
}
