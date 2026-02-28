package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPartnersCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold tier partners",
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/partners", strings.NewReader(url.Values{
		"name":          {"Acme Corp"},
		"tier_id":       {fmt.Sprintf("%d", tier.ID)},
		"logo_url":      {"https://example.com/logo.png"},
		"icon":          {"business"},
		"website_url":   {"https://acme.com"},
		"description":   {"Leading solutions provider"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	partners, _ := queries.ListAllPartners(ctx)
	if len(partners) != 1 {
		t.Fatalf("expected 1 partner, got %d", len(partners))
	}
	if partners[0].Name != "Acme Corp" {
		t.Errorf("expected 'Acme Corp', got %q", partners[0].Name)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/partners/%d", partners[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	_, err := queries.GetPartner(ctx, partners[0].ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestPartnersList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold tier partners",
		SortOrder:   1,
	})

	for i := 1; i <= 3; i++ {
		queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
			Name:         fmt.Sprintf("Partner %d", i),
			TierID:       tier.ID,
			DisplayOrder: int64(i),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/partners", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners list route not found")
	}
}

func TestPartnerEdit_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold tier partners",
		SortOrder:   1,
	})

	partner, _ := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Original Name",
		TierID:       tier.ID,
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/partners/%d", partner.ID), strings.NewReader(url.Values{
		"name":          {"Updated Name"},
		"tier_id":       {fmt.Sprintf("%d", tier.ID)},
		"display_order": {"2"},
		"is_active":     {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetPartner(ctx, partner.ID)
	if updated.Name != "Updated Name" {
		t.Errorf("expected 'Updated Name', got %q", updated.Name)
	}
}
