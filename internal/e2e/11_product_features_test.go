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

// NOTE: Product details routes (specs, features, certifications) are not registered in setupApp() yet.
// These tests will fail until routes are added to setupApp().

func TestProductFeaturesList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Software",
		Slug:        "software",
		Description: "Software products",
		Icon:        "code",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "FEAT-001",
		Slug:        "feature-product",
		Name:        "Product with Features",
		Description: "Test product",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    product.ID,
		FeatureText:  "High performance processing",
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/features", product.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProductFeaturesAdd_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Applications",
		Slug:        "applications",
		Description: "Application products",
		Icon:        "app",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "FEAT-002",
		Slug:        "feature-add-product",
		Name:        "Add Feature Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/features", product.ID), strings.NewReader(url.Values{
		"feature_text":  {"Advanced security encryption"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	features, _ := queries.ListProductFeatures(ctx, product.ID)
	if len(features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(features))
	}
	if features[0].FeatureText != "Advanced security encryption" {
		t.Errorf("expected 'Advanced security encryption', got %q", features[0].FeatureText)
	}
}

func TestProductFeaturesDeleteAll_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Platforms",
		Slug:        "platforms",
		Description: "Platform products",
		Icon:        "platform",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "FEAT-003",
		Slug:        "feature-delete-product",
		Name:        "Delete Feature Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    product.ID,
		FeatureText:  "Real-time analytics",
		DisplayOrder: 1,
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    product.ID,
		FeatureText:  "Cloud-based storage",
		DisplayOrder: 2,
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    product.ID,
		FeatureText:  "Multi-platform support",
		DisplayOrder: 3,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/products/%d/features", product.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	features, _ := queries.ListProductFeatures(ctx, product.ID)
	if len(features) != 0 {
		t.Errorf("expected 0 features after delete, got %d", len(features))
	}
}

func TestProductFeaturesEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "FeatEditCat", Slug: "feat-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "FEAT-EDIT-1", Slug: "feature-edit-form", Name: "Edit Form Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	feat, _ := queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: product.ID, FeatureText: "Original feature text", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/features?edit=%d", product.ID, feat.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/features/%d"`, product.ID, feat.ID)) {
		t.Errorf("expected inline edit form for feature %d, body: %s", feat.ID, body)
	}
	if !strings.Contains(body, `value="Original feature text"`) {
		t.Errorf("expected pre-filled feature text, body: %s", body)
	}
}

func TestProductFeaturesUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "FeatUpdCat", Slug: "feat-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "FEAT-UPD-1", Slug: "feature-update", Name: "Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	f1, _ := queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: product.ID, FeatureText: "Alpha", DisplayOrder: 1,
	})
	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: product.ID, FeatureText: "Bravo", DisplayOrder: 2,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/features/%d", product.ID, f1.ID), strings.NewReader(url.Values{
		"feature_text":  {"Alpha Updated"},
		"display_order": {"5"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	features, _ := queries.ListProductFeatures(ctx, product.ID)
	if len(features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(features))
	}
	if features[0].FeatureText != "Bravo" {
		t.Errorf("expected 'Bravo' first after reorder, got %q", features[0].FeatureText)
	}
	if features[1].FeatureText != "Alpha Updated" {
		t.Errorf("expected 'Alpha Updated' second, got %q", features[1].FeatureText)
	}
	if features[1].DisplayOrder != 5 {
		t.Errorf("expected display_order 5, got %d", features[1].DisplayOrder)
	}
}

func TestProductFeaturesDisplayOrder_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	_ = loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Services",
		Slug:        "services",
		Description: "Service products",
		Icon:        "service",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "FEAT-004",
		Slug:        "ordered-features-product",
		Name:        "Ordered Features Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    product.ID,
		FeatureText:  "Third feature",
		DisplayOrder: 3,
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    product.ID,
		FeatureText:  "First feature",
		DisplayOrder: 1,
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    product.ID,
		FeatureText:  "Second feature",
		DisplayOrder: 2,
	})

	features, _ := queries.ListProductFeatures(ctx, product.ID)
	if len(features) != 3 {
		t.Fatalf("expected 3 features, got %d", len(features))
	}

	if features[0].FeatureText != "First feature" {
		t.Errorf("expected first feature to be 'First feature', got %q", features[0].FeatureText)
	}
	if features[1].FeatureText != "Second feature" {
		t.Errorf("expected second feature to be 'Second feature', got %q", features[1].FeatureText)
	}
	if features[2].FeatureText != "Third feature" {
		t.Errorf("expected third feature to be 'Third feature', got %q", features[2].FeatureText)
	}
}
