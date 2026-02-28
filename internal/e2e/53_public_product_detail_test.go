package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicProductDetail(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Detectors",
		Slug:        "detectors",
		Description: "Detection equipment",
		Icon:        "radar",
		ImageUrl:  sql.NullString{},
		SortOrder: 1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "DET-001",
		Slug:        "alpha-detector",
		Name:        "Alpha Detector",
		Description: "High-precision alpha particle detector",
		CategoryID:  cat.ID,
		Status:      "published",
		Tagline:     sql.NullString{String: "Precision Detection", Valid: true},
	})

	t.Run("product detail page loads", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/detectors/alpha-detector", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("wrong category slug returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/wrong-category/alpha-detector", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})

	t.Run("nonexistent product returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/detectors/nonexistent-product", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rec.Code)
		}
	})

	_ = product
}

func TestProductDetailWithFeatures(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Sensors",
		Slug:        "sensors",
		Description: "Sensor products",
		Icon:        "sensor",
		ImageUrl:  sql.NullString{},
		SortOrder: 1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SENS-001",
		Slug:        "temp-sensor",
		Name:        "Temperature Sensor",
		Description: "Precision temperature sensor",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:   product.ID,
		FeatureText:  "High accuracy",
		DisplayOrder: 1,
	})

	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:   product.ID,
		FeatureText:  "Wide temperature range",
		DisplayOrder: 2,
	})

	t.Run("product detail with features", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/sensors/temp-sensor", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestProductDetailWithSpecifications(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Equipment",
		Slug:        "equipment",
		Description: "Industrial equipment",
		Icon:        "tool",
		ImageUrl:  sql.NullString{},
		SortOrder: 1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "EQUIP-001",
		Slug:        "analyzer",
		Name:        "Spectrum Analyzer",
		Description: "High-performance spectrum analyzer",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID: product.ID,
		SectionName: "Technical Specifications",
		SpecKey:   "Frequency Range",
		SpecValue: "10 Hz - 50 GHz",
		DisplayOrder:1,
	})

	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID: product.ID,
		SectionName: "Technical Specifications",
		SpecKey:   "Resolution",
		SpecValue: "1 Hz",
		DisplayOrder:2,
	})

	t.Run("product detail with specifications", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/equipment/analyzer", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestProductDetailWithImages(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Tools",
		Slug:        "tools",
		Description: "Professional tools",
		Icon:        "wrench",
		ImageUrl:  sql.NullString{},
		SortOrder: 1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "TOOL-001",
		Slug:        "multimeter",
		Name:        "Digital Multimeter",
		Description: "Professional digital multimeter",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: product.ID,
		ImagePath:    "/uploads/products/multimeter-1.jpg",
		AltText:      sql.NullString{String: "Multimeter front view", Valid: true},
		Caption:      sql.NullString{},
		DisplayOrder: 1,
		IsThumbnail:  false,
	})

	queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: product.ID,
		ImagePath:    "/uploads/products/multimeter-2.jpg",
		AltText:      sql.NullString{String: "Multimeter side view", Valid: true},
		Caption:      sql.NullString{},
		DisplayOrder: 2,
		IsThumbnail:  false,
	})

	t.Run("product detail with images", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/tools/multimeter", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestProductDetailWithCertifications(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Certified Products",
		Slug:        "certified-products",
		Description: "Certified industrial products",
		Icon:        "cert",
		ImageUrl:  sql.NullString{},
		SortOrder: 1,
	})

	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "CERT-001",
		Slug:        "certified-sensor",
		Name:        "Certified Sensor",
		Description: "Industry certified sensor",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         product.ID,
		CertificationName: "ISO 9001",
		CertificationCode: sql.NullString{},
		IconType:          sql.NullString{String: "image", Valid: true},
		IconPath:          sql.NullString{String: "/uploads/certs/iso9001.jpg", Valid: true},
		DisplayOrder:      1,
	})

	t.Run("product detail with certifications", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/certified-products/certified-sensor", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}
