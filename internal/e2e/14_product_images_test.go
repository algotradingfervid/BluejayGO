package e2e_test

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestProductImages_ListEmpty(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-001", Slug: "img-001", Name: "Image Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/images", prod.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProductImages_AddImage(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-002", Slug: "img-002", Name: "Image Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "product.png")
	part.Write([]byte("fake image content"))
	writer.WriteField("alt_text", "Product image")
	writer.WriteField("caption", "Main view")
	writer.WriteField("is_thumbnail", "1")
	writer.WriteField("display_order", "1")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/images", prod.ID), strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	images, _ := queries.ListProductImages(ctx, prod.ID)
	if len(images) != 1 {
		t.Errorf("expected 1 image, got %d", len(images))
	}
	if !images[0].IsThumbnail {
		t.Error("expected image to be marked as thumbnail")
	}
}

func TestProductImages_DeleteImage(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-003", Slug: "img-003", Name: "Image Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})
	img, _ := queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: prod.ID, ImagePath: "/test.png", DisplayOrder: 1, IsThumbnail: false,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/products/%d/images/%d", prod.ID, img.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	images, _ := queries.ListProductImages(ctx, prod.ID)
	if len(images) != 0 {
		t.Errorf("expected 0 images after delete, got %d", len(images))
	}
}

func TestProductImages_RequiresAuth(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-004", Slug: "img-004", Name: "Image Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/images", prod.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
}

func TestProductImages_AddWithoutFile(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-005", Slug: "img-005", Name: "Image Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	writer.WriteField("alt_text", "Test")
	writer.WriteField("display_order", "1")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/images", prod.ID), strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Errorf("expected error status, got 200")
	}
}

func TestProductImages_DisplayOrder(t *testing.T) {
	_, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-006", Slug: "img-006", Name: "Image Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: prod.ID, ImagePath: "/test2.png", DisplayOrder: 2, IsThumbnail: false,
	})
	queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: prod.ID, ImagePath: "/test1.png", DisplayOrder: 1, IsThumbnail: true,
	})

	images, _ := queries.ListProductImages(ctx, prod.ID)
	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(images))
	}
	if images[0].DisplayOrder != 1 {
		t.Errorf("expected first image display_order 1, got %d", images[0].DisplayOrder)
	}
	if !images[0].IsThumbnail {
		t.Error("expected first image to be thumbnail")
	}
}

func TestProductImages_OptionalFields(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-007", Slug: "img-007", Name: "Image Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "product.png")
	part.Write([]byte("fake image content"))
	writer.WriteField("display_order", "1")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/images", prod.ID), strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	images, _ := queries.ListProductImages(ctx, prod.ID)
	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}
	if images[0].AltText.Valid {
		t.Error("expected alt_text to be null when not provided")
	}
	if images[0].Caption.Valid {
		t.Error("expected caption to be null when not provided")
	}
}

func TestProductImagesEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "ImgEditCat", Slug: "img-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-EDIT-1", Slug: "img-edit", Name: "Image Edit Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	img, _ := queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: product.ID, ImagePath: "/uploads/products/img.jpg", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/images?edit=%d", product.ID, img.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/images/%d"`, product.ID, img.ID)) {
		t.Errorf("expected inline edit form for image %d, body: %s", img.ID, body)
	}
	if !strings.Contains(body, `name="alt_text"`) {
		t.Errorf("expected alt_text input in edit form, body: %s", body)
	}
}

func TestProductImagesUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "ImgUpdCat", Slug: "img-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-UPD-1", Slug: "img-upd", Name: "Image Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	img, _ := queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: product.ID, ImagePath: "/uploads/products/img.jpg", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/images/%d", product.ID, img.ID), strings.NewReader(url.Values{
		"alt_text":      {"Front view"},
		"caption":       {"Product front"},
		"display_order": {"7"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	images, _ := queries.ListProductImages(ctx, product.ID)
	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}
	if !images[0].AltText.Valid || images[0].AltText.String != "Front view" {
		t.Errorf("expected alt text 'Front view', got %+v", images[0].AltText)
	}
	if images[0].DisplayOrder != 7 {
		t.Errorf("expected display_order 7, got %d", images[0].DisplayOrder)
	}
	// Metadata-only edit must preserve the uploaded image path.
	if images[0].ImagePath != "/uploads/products/img.jpg" {
		t.Errorf("expected image path preserved, got %q", images[0].ImagePath)
	}
}
