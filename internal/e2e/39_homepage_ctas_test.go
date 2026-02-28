package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHomepageCTACRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/homepage/cta", strings.NewReader(url.Values{
		"headline":           {"Get Started Today"},
		"description":        {"Join thousands of satisfied customers"},
		"primary_cta_text":   {"Sign Up"},
		"primary_cta_url":    {"/signup"},
		"secondary_cta_text": {"Learn More"},
		"secondary_cta_url":  {"/about"},
		"background_style":   {"gradient"},
		"is_active":          {"on"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	ctas, _ := queries.ListAllCTAs(ctx)
	if len(ctas) != 1 {
		t.Fatalf("expected 1 CTA, got %d", len(ctas))
	}
	if ctas[0].Headline != "Get Started Today" {
		t.Errorf("expected headline 'Get Started Today', got %q", ctas[0].Headline)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/homepage/cta/%d", ctas[0].ID), strings.NewReader(url.Values{
		"headline":           {"Updated CTA"},
		"description":        {""},
		"primary_cta_text":   {"Contact Us"},
		"primary_cta_url":    {"/contact"},
		"secondary_cta_text": {""},
		"secondary_cta_url":  {""},
		"background_style":   {"dark"},
		"is_active":          {""},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("update: expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetCTA(ctx, ctas[0].ID)
	if updated.Headline != "Updated CTA" {
		t.Errorf("expected headline 'Updated CTA', got %q", updated.Headline)
	}
	if updated.IsActive != 0 {
		t.Errorf("expected is_active 0, got %d", updated.IsActive)
	}
	if updated.BackgroundStyle.Valid && updated.BackgroundStyle.String != "dark" {
		t.Errorf("expected background_style 'dark', got %q", updated.BackgroundStyle.String)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/homepage/cta/%d", ctas[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	ctas, _ = queries.ListAllCTAs(ctx)
	if len(ctas) != 0 {
		t.Errorf("expected 0 CTAs after delete, got %d", len(ctas))
	}
}

func TestHomepageCTABackgroundStyles_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	styles := []string{"light", "dark", "gradient", "primary"}
	for _, style := range styles {
		req := httptest.NewRequest(http.MethodPost, "/admin/homepage/cta", strings.NewReader(url.Values{
			"headline":         {fmt.Sprintf("CTA %s", style)},
			"primary_cta_text": {"Click"},
			"primary_cta_url":  {"/"},
			"background_style": {style},
			"is_active":        {"on"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("create with style %q: expected 303, got %d", style, rec.Code)
		}
	}

	ctas, _ := queries.ListAllCTAs(ctx)
	if len(ctas) != len(styles) {
		t.Errorf("expected %d CTAs, got %d", len(styles), len(ctas))
	}
}
