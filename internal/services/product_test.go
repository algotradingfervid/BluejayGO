package services_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestGetProductDetail(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Detectors", Slug: "detectors", Description: "desc", Icon: "icon", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	prod, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DET-001", Slug: "detector-one", Name: "Detector One",
		Description: "A detector", CategoryID: cat.ID, Status: "published",
	})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	// Add related data
	_, err = queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID: prod.ID, SectionName: "General", SpecKey: "Weight", SpecValue: "5kg", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateProductSpec: %v", err)
	}

	_, err = queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: prod.ID, FeatureText: "High accuracy", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateProductFeature: %v", err)
	}

	_, err = queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: prod.ID, ImagePath: "/img/det.jpg", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateProductImage: %v", err)
	}

	_, err = queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID: prod.ID, CertificationName: "ISO 9001", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateProductCertification: %v", err)
	}

	_, err = queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: prod.ID, Title: "Datasheet", FileType: "pdf", FilePath: "/ds.pdf", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateProductDownload: %v", err)
	}

	svc := services.NewProductService(queries)
	detail, err := svc.GetProductDetail(ctx, "detector-one")
	if err != nil {
		t.Fatalf("GetProductDetail: %v", err)
	}

	if detail.Product.Name != "Detector One" {
		t.Errorf("expected product 'Detector One', got %q", detail.Product.Name)
	}
	if detail.Category.Name != "Detectors" {
		t.Errorf("expected category 'Detectors', got %q", detail.Category.Name)
	}
	if len(detail.Specs) != 1 {
		t.Errorf("expected 1 spec, got %d", len(detail.Specs))
	}
	if len(detail.Features) != 1 {
		t.Errorf("expected 1 feature, got %d", len(detail.Features))
	}
	if len(detail.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(detail.Images))
	}
	if len(detail.Certifications) != 1 {
		t.Errorf("expected 1 cert, got %d", len(detail.Certifications))
	}
	if len(detail.Downloads) != 1 {
		t.Errorf("expected 1 download, got %d", len(detail.Downloads))
	}
}

func TestGetProductDetail_NotFound(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	svc := services.NewProductService(queries)
	_, err := svc.GetProductDetail(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent product")
	}
}

func TestGetProductDetail_EmptyRelated(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "EMPTY-001", Slug: "empty-prod", Name: "Empty", Description: "d", CategoryID: cat.ID, Status: "draft",
	})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	svc := services.NewProductService(queries)
	detail, err := svc.GetProductDetail(ctx, "empty-prod")
	if err != nil {
		t.Fatalf("GetProductDetail: %v", err)
	}

	if len(detail.Specs) != 0 {
		t.Errorf("expected 0 specs, got %d", len(detail.Specs))
	}
	if len(detail.Features) != 0 {
		t.Errorf("expected 0 features, got %d", len(detail.Features))
	}
	if len(detail.Images) != 0 {
		t.Errorf("expected 0 images, got %d", len(detail.Images))
	}
}

// Ensure sql import is used (for nullable fields in CreateProductParams)
var _ = sql.NullString{}
