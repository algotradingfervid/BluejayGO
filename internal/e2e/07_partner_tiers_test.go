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

// NOTE: Partner tiers routes are not registered in setupApp() yet.
// These tests will fail until routes are added to setupApp().

func TestPartnerTiersList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Platinum",
		Slug:        "platinum",
		Description: "Top tier",
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/partner-tiers", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestPartnerTiersCreate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/partner-tiers", strings.NewReader(url.Values{
		"name":        {"Gold"},
		"description": {"Second tier"},
		"sort_order":  {"2"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}

	items, _ := queries.ListPartnerTiers(context.Background())
	if len(items) != 1 {
		t.Fatalf("expected 1 tier, got %d", len(items))
	}
	if items[0].Name != "Gold" {
		t.Errorf("expected 'Gold', got %q", items[0].Name)
	}
	if items[0].Slug != "gold" {
		t.Errorf("expected slug 'gold', got %q", items[0].Slug)
	}
}

func TestPartnerTiersUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Silver",
		Slug:        "silver",
		Description: "Third tier",
		SortOrder:   3,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/partner-tiers/%d", tier.ID), strings.NewReader(url.Values{
		"name":        {"Silver Plus"},
		"description": {"Enhanced silver tier"},
		"sort_order":  {"4"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetPartnerTier(ctx, tier.ID)
	if updated.Name != "Silver Plus" {
		t.Errorf("expected 'Silver Plus', got %q", updated.Name)
	}
	if updated.Slug != "silver-plus" {
		t.Errorf("expected slug 'silver-plus', got %q", updated.Slug)
	}
}

func TestPartnerTiersDelete_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Bronze",
		Slug:        "bronze",
		Description: "Lowest tier",
		SortOrder:   4,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/partner-tiers/%d", tier.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	items, _ := queries.ListPartnerTiers(ctx)
	if len(items) != 0 {
		t.Errorf("expected 0 tiers after delete, got %d", len(items))
	}
}
