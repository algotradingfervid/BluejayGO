package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestSearchPage_Loads(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSearchPage_WithQuery(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Detectors",
		Slug:        "detectors",
		Description: "Detection equipment",
		Icon:        "radar",
		SortOrder:   1,
	})

	_, _ = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "DET-001",
		Slug:        "alpha-detector",
		Name:        "Alpha Detector",
		Description: "Detects alpha particles",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	req := httptest.NewRequest(http.MethodGet, "/search?q=alpha", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSearchPage_EmptyQuery(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/search?q=", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSearchSuggest_WithQuery(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Detectors",
		Slug:        "detectors",
		Description: "Detection equipment",
		Icon:        "radar",
		SortOrder:   1,
	})

	_, _ = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "DET-001",
		Slug:        "test-product",
		Name:        "Test Product",
		Description: "Test description",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	req := httptest.NewRequest(http.MethodGet, "/search/suggest?q=test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSearchSuggest_EmptyQuery(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/search/suggest?q=", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSearch_SpecialCharacters(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/search?q=%3Cscript%3E", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
