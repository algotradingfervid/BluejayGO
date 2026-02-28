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

func TestHomepageHeroesCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/homepage/heroes", strings.NewReader(url.Values{
		"headline":           {"Test Hero"},
		"subheadline":        {"Test Subheadline"},
		"badge_text":         {"New"},
		"primary_cta_text":   {"Learn More"},
		"primary_cta_url":    {"/about"},
		"secondary_cta_text": {"Contact Us"},
		"secondary_cta_url":  {"/contact"},
		"background_image":   {"/images/hero.jpg"},
		"is_active":          {"on"},
		"display_order":      {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	heroes, _ := queries.ListAllHeroes(ctx)
	if len(heroes) != 1 {
		t.Fatalf("expected 1 hero, got %d", len(heroes))
	}
	if heroes[0].Headline != "Test Hero" {
		t.Errorf("expected headline 'Test Hero', got %q", heroes[0].Headline)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/homepage/heroes/%d", heroes[0].ID), strings.NewReader(url.Values{
		"headline":           {"Updated Hero"},
		"subheadline":        {"Updated Subheadline"},
		"badge_text":         {""},
		"primary_cta_text":   {"Get Started"},
		"primary_cta_url":    {"/signup"},
		"secondary_cta_text": {""},
		"secondary_cta_url":  {""},
		"background_image":   {""},
		"is_active":          {""},
		"display_order":      {"2"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("update: expected 303, got %d", rec.Code)
	}

	updatedHero, _ := queries.GetHero(ctx, heroes[0].ID)
	if updatedHero.Headline != "Updated Hero" {
		t.Errorf("expected headline 'Updated Hero', got %q", updatedHero.Headline)
	}
	if updatedHero.IsActive != 0 {
		t.Errorf("expected is_active 0, got %d", updatedHero.IsActive)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/homepage/heroes/%d", heroes[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	heroes, _ = queries.ListAllHeroes(ctx)
	if len(heroes) != 0 {
		t.Errorf("expected 0 heroes after delete, got %d", len(heroes))
	}
}

func TestHomepageHeroesValidation_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/homepage/heroes", strings.NewReader(url.Values{
		"headline":         {""},
		"primary_cta_text": {"Click Me"},
		"primary_cta_url":  {""},
		"display_order":    {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("hero create route not found")
	}
}
