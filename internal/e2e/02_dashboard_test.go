package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestDashboardAccess_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusSeeOther {
		t.Errorf("expected authenticated access, got redirect to %s", rec.Header().Get("Location"))
	}
}

func TestDashboardStats_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Test Cat",
		Slug:        "test-cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "P1",
		Slug:        "p1",
		Name:        "Product 1",
		Description: "d",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "P2",
		Slug:        "p2",
		Name:        "Product 2",
		Description: "d",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	pubCount, _ := queries.CountProducts(ctx)
	if pubCount != 1 {
		t.Errorf("expected 1 published product, got %d", pubCount)
	}

	draftCount, _ := queries.CountDraftProducts(ctx)
	if draftCount != 1 {
		t.Errorf("expected 1 draft product, got %d", draftCount)
	}
}
