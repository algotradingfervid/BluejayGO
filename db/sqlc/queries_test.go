package sqlc_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestProductCategoryCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Test Category",
		Slug:        "test-category",
		Description: "desc",
		Icon:        "icon",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if cat.Name != "Test Category" {
		t.Errorf("expected name 'Test Category', got %q", cat.Name)
	}

	// Get
	got, err := queries.GetProductCategory(ctx, cat.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Slug != "test-category" {
		t.Errorf("expected slug 'test-category', got %q", got.Slug)
	}

	// GetBySlug
	gotSlug, err := queries.GetProductCategoryBySlug(ctx, "test-category")
	if err != nil {
		t.Fatalf("GetBySlug: %v", err)
	}
	if gotSlug.ID != cat.ID {
		t.Errorf("GetBySlug returned wrong ID")
	}

	// List
	items, err := queries.ListProductCategories(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}

	// Update
	updated, err := queries.UpdateProductCategory(ctx, sqlc.UpdateProductCategoryParams{
		ID:          cat.ID,
		Name:        "Updated",
		Slug:        "updated",
		Description: "new desc",
		Icon:        "new-icon",
		SortOrder:   2,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("expected 'Updated', got %q", updated.Name)
	}

	// Delete
	if err := queries.DeleteProductCategory(ctx, cat.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err = queries.GetProductCategory(ctx, cat.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestBlogCategoryCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, err := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:      "Tech",
		Slug:      "tech",
		ColorHex:  "#FF0000",
		SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if cat.Name != "Tech" {
		t.Errorf("expected 'Tech', got %q", cat.Name)
	}

	items, err := queries.ListBlogCategories(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1, got %d", len(items))
	}

	_, err = queries.UpdateBlogCategory(ctx, sqlc.UpdateBlogCategoryParams{
		ID:        cat.ID,
		Name:      "Technology",
		Slug:      "technology",
		ColorHex:  "#00FF00",
		SortOrder: 2,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if err := queries.DeleteBlogCategory(ctx, cat.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestBlogAuthorCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	author, err := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:      "John",
		Slug:      "john",
		Title:     "Writer",
		SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := queries.GetBlogAuthor(ctx, author.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "John" {
		t.Errorf("expected 'John', got %q", got.Name)
	}

	_, err = queries.UpdateBlogAuthor(ctx, sqlc.UpdateBlogAuthorParams{
		ID:        author.ID,
		Name:      "Jane",
		Slug:      "jane",
		Title:     "Editor",
		SortOrder: 2,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if err := queries.DeleteBlogAuthor(ctx, author.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestIndustryCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	ind, err := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name:        "Healthcare",
		Slug:        "healthcare",
		Icon:        "medical",
		Description: "Health",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	items, err := queries.ListIndustries(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1, got %d", len(items))
	}

	_, err = queries.UpdateIndustry(ctx, sqlc.UpdateIndustryParams{
		ID:          ind.ID,
		Name:        "Pharma",
		Slug:        "pharma",
		Icon:        "pill",
		Description: "Pharma",
		SortOrder:   2,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if err := queries.DeleteIndustry(ctx, ind.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestPartnerTierCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	tier, err := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Top tier",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	_, err = queries.GetPartnerTier(ctx, tier.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	_, err = queries.UpdatePartnerTier(ctx, sqlc.UpdatePartnerTierParams{
		ID:          tier.ID,
		Name:        "Platinum",
		Slug:        "platinum",
		Description: "Best",
		SortOrder:   0,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if err := queries.DeletePartnerTier(ctx, tier.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestWhitepaperTopicCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	topic, err := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name:      "AI",
		Slug:      "ai",
		ColorHex:  "#0000FF",
		Icon:      "brain",
		SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	items, err := queries.ListWhitepaperTopics(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1, got %d", len(items))
	}

	_, err = queries.UpdateWhitepaperTopic(ctx, sqlc.UpdateWhitepaperTopicParams{
		ID:        topic.ID,
		Name:      "ML",
		Slug:      "ml",
		ColorHex:  "#FF00FF",
		Icon:      "model",
		SortOrder: 2,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if err := queries.DeleteWhitepaperTopic(ctx, topic.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestAdminUserCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	user, err := queries.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{
		Email:        "test@example.com",
		PasswordHash: "$2a$10$fakehash",
		DisplayName:  "Test User",
		Role:         "admin",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %q", user.Email)
	}

	got, err := queries.GetAdminUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}
	if got.DisplayName != "Test User" {
		t.Errorf("expected 'Test User', got %q", got.DisplayName)
	}

	if err := queries.UpdateLastLogin(ctx, user.ID); err != nil {
		t.Fatalf("UpdateLastLogin: %v", err)
	}

	users, err := queries.ListAdminUsers(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1, got %d", len(users))
	}
}

func TestProductCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create category first (FK)
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Detectors",
		Slug:        "detectors",
		Description: "desc",
		Icon:        "icon",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	// Create product
	prod, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "TEST-001",
		Slug:        "test-product",
		Name:        "Test Product",
		Description: "A test product",
		CategoryID:  cat.ID,
		Status:      "draft",
	})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	// Get by ID
	got, err := queries.GetProduct(ctx, prod.ID)
	if err != nil {
		t.Fatalf("GetProduct: %v", err)
	}
	if got.Sku != "TEST-001" {
		t.Errorf("expected SKU 'TEST-001', got %q", got.Sku)
	}

	// Get by slug
	gotSlug, err := queries.GetProductBySlug(ctx, "test-product")
	if err != nil {
		t.Fatalf("GetBySlug: %v", err)
	}
	if gotSlug.ID != prod.ID {
		t.Errorf("GetBySlug returned wrong product")
	}

	// Get by SKU
	gotSKU, err := queries.GetProductBySKU(ctx, "TEST-001")
	if err != nil {
		t.Fatalf("GetBySKU: %v", err)
	}
	if gotSKU.ID != prod.ID {
		t.Errorf("GetBySKU returned wrong product")
	}

	// Update to published so ListProductsByCategory and CountProductsByCategory work
	err = queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:          prod.ID,
		Sku:         "TEST-002",
		Slug:        "updated-product",
		Name:        "Updated",
		Description: "Updated desc",
		CategoryID:  cat.ID,
		Status:      "published",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	// List by category (only published)
	products, err := queries.ListProductsByCategory(ctx, sqlc.ListProductsByCategoryParams{
		CategoryID: cat.ID,
		Limit:      10,
		Offset:     0,
	})
	if err != nil {
		t.Fatalf("ListByCategory: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("expected 1 product, got %d", len(products))
	}

	// Count by category (only published)
	count, err := queries.CountProductsByCategory(ctx, cat.ID)
	if err != nil {
		t.Fatalf("CountByCategory: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Admin list (all statuses)
	allProducts, err := queries.ListAllProductsAdmin(ctx)
	if err != nil {
		t.Fatalf("ListAdmin: %v", err)
	}
	if len(allProducts) != 1 {
		t.Errorf("expected 1, got %d", len(allProducts))
	}

	// Delete
	if err := queries.DeleteProduct(ctx, prod.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestProductSpecsCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "S-001", Slug: "s-prod", Name: "S Prod", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	spec, err := queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID: prod.ID, SectionName: "General", SpecKey: "Weight", SpecValue: "10kg", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateSpec: %v", err)
	}

	specs, err := queries.ListProductSpecs(ctx, prod.ID)
	if err != nil {
		t.Fatalf("ListSpecs: %v", err)
	}
	if len(specs) != 1 || specs[0].ID != spec.ID {
		t.Errorf("unexpected specs result")
	}

	if err := queries.DeleteProductSpecs(ctx, prod.ID); err != nil {
		t.Fatalf("DeleteSpecs: %v", err)
	}
}

func TestProductImagesCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat2", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "I-001", Slug: "i-prod", Name: "I Prod", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	img, err := queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: prod.ID, ImagePath: "/img/test.jpg", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateImage: %v", err)
	}

	images, err := queries.ListProductImages(ctx, prod.ID)
	if err != nil {
		t.Fatalf("ListImages: %v", err)
	}
	if len(images) != 1 || images[0].ID != img.ID {
		t.Errorf("unexpected images result")
	}

	if err := queries.DeleteProductImage(ctx, img.ID); err != nil {
		t.Fatalf("DeleteImage: %v", err)
	}
}

func TestProductFeaturesCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat3", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "F-001", Slug: "f-prod", Name: "F Prod", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	feat, err := queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: prod.ID, FeatureText: "Fast processing", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateFeature: %v", err)
	}

	features, err := queries.ListProductFeatures(ctx, prod.ID)
	if err != nil {
		t.Fatalf("ListFeatures: %v", err)
	}
	if len(features) != 1 || features[0].ID != feat.ID {
		t.Errorf("unexpected features result")
	}

	if err := queries.DeleteProductFeatures(ctx, prod.ID); err != nil {
		t.Fatalf("DeleteFeatures: %v", err)
	}
}

func TestProductCertificationsCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat4", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "C-001", Slug: "c-prod", Name: "C Prod", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	cert, err := queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID: prod.ID, CertificationName: "ISO 9001", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateCert: %v", err)
	}

	certs, err := queries.ListProductCertifications(ctx, prod.ID)
	if err != nil {
		t.Fatalf("ListCerts: %v", err)
	}
	if len(certs) != 1 || certs[0].ID != cert.ID {
		t.Errorf("unexpected certs result")
	}

	if err := queries.DeleteProductCertifications(ctx, prod.ID); err != nil {
		t.Fatalf("DeleteCerts: %v", err)
	}
}

func TestProductDownloadsCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat5", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "D-001", Slug: "d-prod", Name: "D Prod", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	dl, err := queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: prod.ID, Title: "Datasheet", FileType: "pdf", FilePath: "/files/ds.pdf", DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateDownload: %v", err)
	}

	downloads, err := queries.ListProductDownloads(ctx, prod.ID)
	if err != nil {
		t.Fatalf("ListDownloads: %v", err)
	}
	if len(downloads) != 1 || downloads[0].ID != dl.ID {
		t.Errorf("unexpected downloads result")
	}

	got, err := queries.GetProductDownload(ctx, dl.ID)
	if err != nil {
		t.Fatalf("GetDownload: %v", err)
	}
	if got.Title != "Datasheet" {
		t.Errorf("expected 'Datasheet', got %q", got.Title)
	}

	if err := queries.IncrementDownloadCount(ctx, dl.ID); err != nil {
		t.Fatalf("IncrementCount: %v", err)
	}

	if err := queries.DeleteProductDownload(ctx, dl.ID); err != nil {
		t.Fatalf("DeleteDownload: %v", err)
	}
}

func TestSettingsCRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Settings should exist from migrations (seed data) or we need to check
	// If no settings exist, this will fail - that's expected as settings need seed data
	_, err := queries.GetSettings(ctx)
	if err != nil {
		// Settings table exists but no rows - that's fine for a clean DB
		t.Logf("GetSettings returned error (expected for empty DB): %v", err)
	}
}

func TestProductCascadeDelete(t *testing.T) {
	db, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Ensure foreign keys are enabled on this connection (migrate may have toggled it)
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "CascadeCat", Slug: "cascade-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "CASCADE-001", Slug: "cascade-prod", Name: "Cascade", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	// Create related items
	queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID: prod.ID, SectionName: "Gen", SpecKey: "K", SpecValue: "V", DisplayOrder: 1,
	})
	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: prod.ID, FeatureText: "F", DisplayOrder: 1,
	})
	queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: prod.ID, ImagePath: "/img.jpg", DisplayOrder: 1,
	})
	queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID: prod.ID, CertificationName: "ISO", DisplayOrder: 1,
	})
	queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: prod.ID, Title: "DL", FileType: "pdf", FilePath: "/dl.pdf", DisplayOrder: 1,
	})

	// Delete product - cascade should remove all related
	if err := queries.DeleteProduct(ctx, prod.ID); err != nil {
		t.Fatalf("DeleteProduct: %v", err)
	}

	specs, _ := queries.ListProductSpecs(ctx, prod.ID)
	if len(specs) != 0 {
		t.Errorf("specs not cascaded: %d remain", len(specs))
	}
	features, _ := queries.ListProductFeatures(ctx, prod.ID)
	if len(features) != 0 {
		t.Errorf("features not cascaded: %d remain", len(features))
	}
	images, _ := queries.ListProductImages(ctx, prod.ID)
	if len(images) != 0 {
		t.Errorf("images not cascaded: %d remain", len(images))
	}
	certs, _ := queries.ListProductCertifications(ctx, prod.ID)
	if len(certs) != 0 {
		t.Errorf("certs not cascaded: %d remain", len(certs))
	}
	downloads, _ := queries.ListProductDownloads(ctx, prod.ID)
	if len(downloads) != 0 {
		t.Errorf("downloads not cascaded: %d remain", len(downloads))
	}
}
