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

func TestProductSpecsList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Hardware",
		Slug:        "hardware",
		Description: "Hardware products",
		Icon:        "chip",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SPEC-001",
		Slug:        "spec-product",
		Name:        "Product with Specs",
		Description: "Test product",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID:    product.ID,
		SectionName:  "Dimensions",
		SpecKey:      "Weight",
		SpecValue:    "5 kg",
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/specs", product.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProductSpecsAdd_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Equipment",
		Slug:        "equipment",
		Description: "Equipment products",
		Icon:        "tool",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SPEC-002",
		Slug:        "spec-add-product",
		Name:        "Add Spec Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/specs", product.ID), strings.NewReader(url.Values{
		"section_name":  {"Power"},
		"spec_key":      {"Voltage"},
		"spec_value":    {"220V"},
		"display_order": {"2"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	specs, _ := queries.ListProductSpecs(ctx, product.ID)
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}
	if specs[0].SpecKey != "Voltage" {
		t.Errorf("expected 'Voltage', got %q", specs[0].SpecKey)
	}
	if specs[0].SpecValue != "220V" {
		t.Errorf("expected '220V', got %q", specs[0].SpecValue)
	}
	if specs[0].SectionName != "Power" {
		t.Errorf("expected section 'Power', got %q", specs[0].SectionName)
	}
}

func TestProductSpecsDeleteAll_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Devices",
		Slug:        "devices",
		Description: "Device products",
		Icon:        "device",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SPEC-003",
		Slug:        "spec-delete-product",
		Name:        "Delete Spec Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID:    product.ID,
		SectionName:  "General",
		SpecKey:      "Color",
		SpecValue:    "Black",
		DisplayOrder: 1,
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID:    product.ID,
		SectionName:  "General",
		SpecKey:      "Material",
		SpecValue:    "Aluminum",
		DisplayOrder: 2,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/products/%d/specs", product.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	specs, _ := queries.ListProductSpecs(ctx, product.ID)
	if len(specs) != 0 {
		t.Errorf("expected 0 specs after delete, got %d", len(specs))
	}
}

func TestProductSpecsGroupedBySections_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	_ = loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Machinery",
		Slug:        "machinery",
		Description: "Machine products",
		Icon:        "gear",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SPEC-004",
		Slug:        "grouped-specs-product",
		Name:        "Grouped Specs Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID:    product.ID,
		SectionName:  "Physical",
		SpecKey:      "Height",
		SpecValue:    "50 cm",
		DisplayOrder: 1,
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID:    product.ID,
		SectionName:  "Physical",
		SpecKey:      "Width",
		SpecValue:    "30 cm",
		DisplayOrder: 2,
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID:    product.ID,
		SectionName:  "Electrical",
		SpecKey:      "Power",
		SpecValue:    "500W",
		DisplayOrder: 1,
	})

	specs, _ := queries.ListProductSpecs(ctx, product.ID)
	if len(specs) != 3 {
		t.Errorf("expected 3 specs, got %d", len(specs))
	}

	physicalCount := 0
	electricalCount := 0
	for _, spec := range specs {
		if spec.SectionName == "Physical" {
			physicalCount++
		}
		if spec.SectionName == "Electrical" {
			electricalCount++
		}
	}

	if physicalCount != 2 {
		t.Errorf("expected 2 Physical specs, got %d", physicalCount)
	}
	if electricalCount != 1 {
		t.Errorf("expected 1 Electrical spec, got %d", electricalCount)
	}
}
