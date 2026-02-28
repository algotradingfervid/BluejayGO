package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

// NOTE: Industries routes are not registered in setupApp() yet.
// These tests will fail until routes are added to setupApp().

func TestIndustriesList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name:        "Healthcare",
		Slug:        "healthcare",
		Icon:        "medical-cross",
		Description: "Medical industry",
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/industries", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestIndustriesCreate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/industries", strings.NewReader(url.Values{
		"name":        {"Finance"},
		"icon":        {"dollar-sign"},
		"description": {"Financial services"},
		"sort_order":  {"2"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}

	items, _ := queries.ListIndustries(context.Background())
	if len(items) != 1 {
		t.Fatalf("expected 1 industry, got %d", len(items))
	}
	if items[0].Name != "Finance" {
		t.Errorf("expected 'Finance', got %q", items[0].Name)
	}
	if items[0].Slug != "finance" {
		t.Errorf("expected slug 'finance', got %q", items[0].Slug)
	}
}

func TestIndustriesUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	industry, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name:        "Tech",
		Slug:        "tech",
		Icon:        "laptop",
		Description: "Technology sector",
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/industries/%d", industry.ID), strings.NewReader(url.Values{
		"name":        {"Technology"},
		"icon":        {"computer"},
		"description": {"Updated description"},
		"sort_order":  {"5"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetIndustry(ctx, industry.ID)
	if updated.Name != "Technology" {
		t.Errorf("expected 'Technology', got %q", updated.Name)
	}
	if updated.Slug != "technology" {
		t.Errorf("expected slug 'technology', got %q", updated.Slug)
	}
}

func TestIndustriesDelete_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	industry, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name:        "Delete Me",
		Slug:        "delete-me",
		Icon:        "trash",
		Description: "To be deleted",
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/industries/%d", industry.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	items, _ := queries.ListIndustries(ctx)
	if len(items) != 0 {
		t.Errorf("expected 0 industries after delete, got %d", len(items))
	}
}
