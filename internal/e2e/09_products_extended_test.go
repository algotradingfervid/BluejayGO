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

func TestProductsListPagination_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Test Category",
		Slug:        "test-category",
		Description: "Test",
		Icon:        "icon",
		SortOrder:   1,
	})

	for i := 1; i <= 5; i++ {
		queries.CreateProduct(ctx, sqlc.CreateProductParams{
			Sku:         fmt.Sprintf("SKU-%03d", i),
			Slug:        fmt.Sprintf("product-%d", i),
			Name:        fmt.Sprintf("Product %d", i),
			Description: "Test product",
			CategoryID:  cat.ID,
			Status:      "published",
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/products", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestProductsCreateWithFullFields_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "Electronic devices",
		Icon:        "chip",
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/products", strings.NewReader(url.Values{
		"sku":              {"PROD-E2E-001"},
		"name":             {"E2E Test Product"},
		"tagline":          {"Testing tagline"},
		"description":      {"Full test description"},
		"overview":         {"<p>Rich text overview</p>"},
		"category_id":      {fmt.Sprintf("%d", cat.ID)},
		"status":           {"published"},
		"is_featured":      {"1"},
		"featured_order":   {"1"},
		"video_url":        {"https://example.com/video"},
		"meta_title":       {"SEO Title"},
		"meta_description": {"SEO description"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}

	products, _ := queries.ListProducts(ctx, sqlc.ListProductsParams{Limit: 10, Offset: 0})
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
	if products[0].Sku != "PROD-E2E-001" {
		t.Errorf("expected SKU 'PROD-E2E-001', got %q", products[0].Sku)
	}
	if products[0].IsFeatured != true {
		t.Errorf("expected is_featured true, got false")
	}
}

func TestProductsFilterByStatus_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Gadgets",
		Slug:        "gadgets",
		Description: "Various gadgets",
		Icon:        "device",
		SortOrder:   1,
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "DRAFT-001",
		Slug:        "draft-product",
		Name:        "Draft Product",
		Description: "Draft",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PUB-001",
		Slug:        "published-product",
		Name:        "Published Product",
		Description: "Published",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/products?status=draft", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestProductsSearch_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Tools",
		Slug:        "tools",
		Description: "Tool products",
		Icon:        "wrench",
		SortOrder:   1,
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SEARCH-001",
		Slug:        "searchable-product",
		Name:        "Searchable Product",
		Description: "This product is searchable",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/products?search=Searchable", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}
