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

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

// NOTE: Product details routes (specs, features, certifications) are not registered in setupApp() yet.
// These tests will fail until routes are added to setupApp().

func TestProductCertificationsList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Industrial",
		Slug:        "industrial",
		Description: "Industrial products",
		Icon:        "factory",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "CERT-001",
		Slug:        "certification-product",
		Name:        "Product with Certifications",
		Description: "Test product",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         product.ID,
		CertificationName: "CE Certified",
		CertificationCode: sql.NullString{String: "CE-12345", Valid: true},
		IconType:          sql.NullString{String: "badge", Valid: true},
		IconPath:          sql.NullString{String: "/icons/ce.png", Valid: true},
		DisplayOrder:      1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/certifications", product.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProductCertificationsAdd_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Safety",
		Slug:        "safety",
		Description: "Safety products",
		Icon:        "shield",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "CERT-002",
		Slug:        "cert-add-product",
		Name:        "Add Certification Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/certifications", product.ID), strings.NewReader(url.Values{
		"certification_name": {"UL Listed"},
		"certification_code": {"UL-98765"},
		"icon_type":          {"image"},
		"icon_path":          {"/icons/ul.png"},
		"display_order":      {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	certs, _ := queries.ListProductCertifications(ctx, product.ID)
	if len(certs) != 1 {
		t.Fatalf("expected 1 certification, got %d", len(certs))
	}
	if certs[0].CertificationName != "UL Listed" {
		t.Errorf("expected 'UL Listed', got %q", certs[0].CertificationName)
	}
	if !certs[0].CertificationCode.Valid || certs[0].CertificationCode.String != "UL-98765" {
		t.Errorf("expected code 'UL-98765', got %q", certs[0].CertificationCode.String)
	}
}

func TestProductCertificationsAddWithoutOptionalFields_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Consumer",
		Slug:        "consumer",
		Description: "Consumer products",
		Icon:        "cart",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "CERT-003",
		Slug:        "cert-optional-product",
		Name:        "Optional Fields Certification Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/certifications", product.ID), strings.NewReader(url.Values{
		"certification_name": {"ISO 9001"},
		"certification_code": {""},
		"icon_type":          {""},
		"icon_path":          {""},
		"display_order":      {"2"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	certs, _ := queries.ListProductCertifications(ctx, product.ID)
	if len(certs) != 1 {
		t.Fatalf("expected 1 certification, got %d", len(certs))
	}
	if certs[0].CertificationName != "ISO 9001" {
		t.Errorf("expected 'ISO 9001', got %q", certs[0].CertificationName)
	}
	if certs[0].CertificationCode.Valid {
		t.Errorf("expected null certification code, got %q", certs[0].CertificationCode.String)
	}
}

func TestProductCertificationsDeleteAll_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Medical",
		Slug:        "medical",
		Description: "Medical products",
		Icon:        "heart",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "CERT-004",
		Slug:        "cert-delete-product",
		Name:        "Delete Certification Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         product.ID,
		CertificationName: "FDA Approved",
		CertificationCode: sql.NullString{String: "FDA-2024", Valid: true},
		IconType:          sql.NullString{String: "badge", Valid: true},
		IconPath:          sql.NullString{String: "/icons/fda.png", Valid: true},
		DisplayOrder:      1,
	})

	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         product.ID,
		CertificationName: "ISO 13485",
		CertificationCode: sql.NullString{String: "ISO-13485-2024", Valid: true},
		IconType:          sql.NullString{String: "image", Valid: true},
		IconPath:          sql.NullString{String: "/icons/iso-13485.png", Valid: true},
		DisplayOrder:      2,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/products/%d/certifications", product.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	certs, _ := queries.ListProductCertifications(ctx, product.ID)
	if len(certs) != 0 {
		t.Errorf("expected 0 certifications after delete, got %d", len(certs))
	}
}

func TestProductCertificationsDisplayOrder_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	_ = loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Electrical",
		Slug:        "electrical",
		Description: "Electrical products",
		Icon:        "bolt",
		SortOrder:   1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "CERT-005",
		Slug:        "ordered-cert-product",
		Name:        "Ordered Certifications Product",
		Description: "Test",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         product.ID,
		CertificationName: "Third Cert",
		CertificationCode: sql.NullString{String: "CERT-003", Valid: true},
		DisplayOrder:      3,
	})

	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         product.ID,
		CertificationName: "First Cert",
		CertificationCode: sql.NullString{String: "CERT-001", Valid: true},
		DisplayOrder:      1,
	})

	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         product.ID,
		CertificationName: "Second Cert",
		CertificationCode: sql.NullString{String: "CERT-002", Valid: true},
		DisplayOrder:      2,
	})

	certs, _ := queries.ListProductCertifications(ctx, product.ID)
	if len(certs) != 3 {
		t.Fatalf("expected 3 certifications, got %d", len(certs))
	}

	if certs[0].CertificationName != "First Cert" {
		t.Errorf("expected first cert to be 'First Cert', got %q", certs[0].CertificationName)
	}
	if certs[1].CertificationName != "Second Cert" {
		t.Errorf("expected second cert to be 'Second Cert', got %q", certs[1].CertificationName)
	}
	if certs[2].CertificationName != "Third Cert" {
		t.Errorf("expected third cert to be 'Third Cert', got %q", certs[2].CertificationName)
	}
}
