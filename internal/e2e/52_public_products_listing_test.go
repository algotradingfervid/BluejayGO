package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicProductsListingPage(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Detectors",
		Slug:        "detectors",
		Description: "Detection equipment",
		Icon:        "radar",
		SortOrder:   1,
	})

	queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Sensors",
		Slug:        "sensors",
		Description: "Sensor equipment",
		Icon:        "sensor",
		SortOrder:   2,
	})

	t.Run("products page loads", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("products page returns HTML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Result().ContentLength == 0 {
			t.Error("expected non-empty response")
		}
	})
}

func TestPublicProductsCategoryRouting(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Industrial Sensors",
		Slug:        "industrial-sensors",
		Description: "Industrial sensor products",
		Icon:        "sensor",
		SortOrder:   1,
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SENS-001",
		Slug:        "temperature-sensor",
		Name:        "Temperature Sensor",
		Description: "High precision temperature sensor",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	t.Run("category page loads", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/industrial-sensors", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("nonexistent category returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/nonexistent-category", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})
}

func TestPublicProductsPagination(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Test Category",
		Slug:        "test-category",
		Description: "Test products",
		Icon:        "test",
		SortOrder:   1,
	})

	for i := 1; i <= 15; i++ {
		queries.CreateProduct(ctx, sqlc.CreateProductParams{
			Sku:         "TEST-" + string(rune(i)),
			Slug:        "test-product-" + string(rune(i)),
			Name:        "Test Product " + string(rune(i)),
			Description: "Product description",
			CategoryID:  cat.ID,
			Status:      "published",
		})
	}

	t.Run("first page loads", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/test-category?page=1", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("second page loads", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/test-category?page=2", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestPublicProductsSearch(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Search Test",
		Slug:        "search-test",
		Description: "Search test category",
		Icon:        "search",
		SortOrder:   1,
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SEARCH-001",
		Slug:        "searchable-product",
		Name:        "Searchable Product",
		Description: "A searchable product",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	t.Run("search endpoint exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/search?q=searchable", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code == http.StatusNotFound {
			t.Error("search endpoint should exist")
		}
	})
}
