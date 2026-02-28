package e2e_test

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
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
