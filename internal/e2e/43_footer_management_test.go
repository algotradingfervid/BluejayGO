package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestFooterSettingsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	formData := url.Values{
		"footer_columns":      {"3"},
		"footer_bg_style":     {"dark"},
		"footer_show_social":  {"on"},
		"footer_social_style": {"icons"},
		"footer_copyright":    {"© 2024 Bluejay Labs"},
		"col_0_heading":       {"Quick Links"},
		"col_0_type":          {"links"},
		"col_1_heading":       {"About Us"},
		"col_1_type":          {"text"},
		"col_1_content":       {"Company description text"},
		"col_2_heading":       {"Contact"},
		"col_2_type":          {"contact"},
	}
	formData["col_0_link_label[]"] = []string{"Home", "Products"}
	formData["col_0_link_url[]"] = []string{"/", "/products"}
	formData["legal_link_label[]"] = []string{"Privacy", "Terms"}
	formData["legal_link_url[]"] = []string{"/privacy", "/terms"}

	req := httptest.NewRequest(http.MethodPost, "/admin/footer", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	settings, _ := queries.GetSettings(ctx)
	if settings.FooterColumns != 3 {
		t.Errorf("expected footer_columns 3, got %d", settings.FooterColumns)
	}
	if settings.FooterShowSocial != 1 {
		t.Errorf("expected footer_show_social 1, got %d", settings.FooterShowSocial)
	}

	columnItems, _ := queries.ListFooterColumnItems(ctx)
	if len(columnItems) != 3 {
		t.Errorf("expected 3 column items, got %d", len(columnItems))
	}

	links, _ := queries.ListAllFooterLinks(ctx)
	if len(links) != 2 {
		t.Errorf("expected 2 footer links, got %d", len(links))
	}

	legalLinks, _ := queries.ListFooterLegalLinks(ctx)
	if len(legalLinks) != 2 {
		t.Errorf("expected 2 legal links, got %d", len(legalLinks))
	}
}

func TestFooterColumnCountChange_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	formData := url.Values{
		"footer_columns":  {"4"},
		"footer_bg_style": {"light"},
		"footer_copyright": {"Test"},
		"col_0_heading":   {"Col 0"},
		"col_0_type":      {"text"},
		"col_1_heading":   {"Col 1"},
		"col_1_type":      {"text"},
		"col_2_heading":   {"Col 2"},
		"col_2_type":      {"text"},
		"col_3_heading":   {"Col 3"},
		"col_3_type":      {"text"},
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/footer", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create 4 columns: expected 303, got %d", rec.Code)
	}

	columnItems, _ := queries.ListFooterColumnItems(ctx)
	if len(columnItems) != 4 {
		t.Fatalf("expected 4 column items, got %d", len(columnItems))
	}

	formData.Set("footer_columns", "2")
	req = httptest.NewRequest(http.MethodPost, "/admin/footer", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	columnItems, _ = queries.ListFooterColumnItems(ctx)
	if len(columnItems) != 2 {
		t.Errorf("expected 2 column items after reduction, got %d", len(columnItems))
	}
}

func TestFooterLinksArray_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	formData := url.Values{
		"footer_columns":  {"2"},
		"footer_bg_style": {"dark"},
		"footer_copyright": {"Test"},
		"col_0_heading":   {"Links"},
		"col_0_type":      {"links"},
		"col_1_heading":   {"Text"},
		"col_1_type":      {"text"},
	}
	formData["col_0_link_label[]"] = []string{"Link 1", "Link 2", "Link 3"}
	formData["col_0_link_url[]"] = []string{"/1", "/2", "/3"}

	req := httptest.NewRequest(http.MethodPost, "/admin/footer", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	links, _ := queries.ListAllFooterLinks(ctx)
	if len(links) != 3 {
		t.Errorf("expected 3 links, got %d", len(links))
	}
}
