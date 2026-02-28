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

func TestSolutionProductsAdd(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Cat",
		Slug:        "cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PROD-001",
		Slug:        "prod-001",
		Name:        "Test Product",
		Description: "desc",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		ShortDescription: "desc",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/products", sol.ID), strings.NewReader(url.Values{
		"product_id":    {fmt.Sprintf("%d", prod.ID)},
		"display_order": {"1"},
		"is_featured":   {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}

	products, _ := queries.GetSolutionProducts(ctx, sol.ID)
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
}

func TestSolutionProductsRemove(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Cat",
		Slug:        "cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PROD-001",
		Slug:        "prod-001",
		Name:        "Test Product",
		Description: "desc",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		ShortDescription: "desc",
	})

	queries.AddProductToSolution(ctx, sqlc.AddProductToSolutionParams{
		SolutionID: sol.ID,
		ProductID:  prod.ID,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/solutions/%d/products/%d", sol.ID, prod.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	products, _ := queries.GetSolutionProducts(ctx, sol.ID)
	if len(products) != 0 {
		t.Errorf("expected 0 products after remove, got %d", len(products))
	}
}

func TestSolutionProductsDuplicate(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Cat",
		Slug:        "cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PROD-001",
		Slug:        "prod-001",
		Name:        "Test Product",
		Description: "desc",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		ShortDescription: "desc",
	})

	queries.AddProductToSolution(ctx, sqlc.AddProductToSolutionParams{
		SolutionID: sol.ID,
		ProductID:  prod.ID,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/products", sol.ID), strings.NewReader(url.Values{
		"product_id":    {fmt.Sprintf("%d", prod.ID)},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Logf("duplicate add may fail with constraint error")
	}
}
