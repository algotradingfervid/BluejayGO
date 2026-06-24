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

func TestProductDownloads_ListEmpty(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-001", Slug: "dl-001", Name: "Download Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/downloads", prod.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProductDownloads_AddDownload(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-002", Slug: "dl-002", Name: "Download Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "manual.pdf")
	part.Write([]byte("fake pdf content"))
	writer.WriteField("title", "Product Manual")
	writer.WriteField("description", "Full manual")
	writer.WriteField("file_type", "PDF")
	writer.WriteField("version", "1.0")
	writer.WriteField("display_order", "1")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/downloads", prod.ID), strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	downloads, _ := queries.ListProductDownloads(ctx, prod.ID)
	if len(downloads) != 1 {
		t.Errorf("expected 1 download, got %d", len(downloads))
	}
	if downloads[0].Title != "Product Manual" {
		t.Errorf("expected title 'Product Manual', got %q", downloads[0].Title)
	}
}

func TestProductDownloads_DeleteDownload(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-003", Slug: "dl-003", Name: "Download Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})
	dl, _ := queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: prod.ID, Title: "Test Download", FileType: "PDF", FilePath: "/test.pdf", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/products/%d/downloads/%d", prod.ID, dl.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	downloads, _ := queries.ListProductDownloads(ctx, prod.ID)
	if len(downloads) != 0 {
		t.Errorf("expected 0 downloads after delete, got %d", len(downloads))
	}
}

func TestProductDownloads_RequiresAuth(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-004", Slug: "dl-004", Name: "Download Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/downloads", prod.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rec.Code)
	}
}

func TestProductDownloads_AddWithoutFile(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-005", Slug: "dl-005", Name: "Download Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Product Manual")
	writer.WriteField("display_order", "1")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/downloads", prod.ID), strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Errorf("expected error status, got 200")
	}
}

func TestProductDownloads_DisplayOrder(t *testing.T) {
	_, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test Cat", Slug: "test-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-006", Slug: "dl-006", Name: "Download Test", Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: prod.ID, Title: "Second", FileType: "PDF", FilePath: "/test2.pdf", DisplayOrder: 2,
	})
	queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: prod.ID, Title: "First", FileType: "PDF", FilePath: "/test1.pdf", DisplayOrder: 1,
	})

	downloads, _ := queries.ListProductDownloads(ctx, prod.ID)
	if len(downloads) != 2 {
		t.Fatalf("expected 2 downloads, got %d", len(downloads))
	}
	if downloads[0].Title != "First" {
		t.Errorf("expected first download to be 'First', got %q", downloads[0].Title)
	}
	if downloads[1].Title != "Second" {
		t.Errorf("expected second download to be 'Second', got %q", downloads[1].Title)
	}
}

func TestProductDownloadsEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "DlEditCat", Slug: "dl-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-EDIT-1", Slug: "dl-edit", Name: "Download Edit Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	dl, _ := queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: product.ID, Title: "Datasheet", FileType: ".pdf",
		FilePath: "/uploads/downloads/ds.pdf", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/downloads?edit=%d", product.ID, dl.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/downloads/%d"`, product.ID, dl.ID)) {
		t.Errorf("expected inline edit form for download %d, body: %s", dl.ID, body)
	}
	if !strings.Contains(body, `value="Datasheet"`) {
		t.Errorf("expected pre-filled download title, body: %s", body)
	}
}

func TestProductDownloadsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "DlUpdCat", Slug: "dl-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-UPD-1", Slug: "dl-upd", Name: "Download Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	dl, _ := queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: product.ID, Title: "Datasheet", FileType: ".pdf",
		FilePath: "/uploads/downloads/ds.pdf", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/downloads/%d", product.ID, dl.ID), strings.NewReader(url.Values{
		"title":         {"User Manual"},
		"file_type":     {".pdf"},
		"version":       {"2.0"},
		"description":   {"Updated manual"},
		"display_order": {"4"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	downloads, _ := queries.ListProductDownloads(ctx, product.ID)
	if len(downloads) != 1 {
		t.Fatalf("expected 1 download, got %d", len(downloads))
	}
	if downloads[0].Title != "User Manual" {
		t.Errorf("expected title 'User Manual', got %q", downloads[0].Title)
	}
	if downloads[0].DisplayOrder != 4 {
		t.Errorf("expected display_order 4, got %d", downloads[0].DisplayOrder)
	}
	// Metadata-only edit must preserve the uploaded file path.
	if downloads[0].FilePath != "/uploads/downloads/ds.pdf" {
		t.Errorf("expected file path preserved, got %q", downloads[0].FilePath)
	}
}
