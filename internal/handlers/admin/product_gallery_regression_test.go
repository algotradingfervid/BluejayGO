package admin_test

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestGalleryPrimaryActionOwnershipUploadAndCache(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	t.Chdir("../../..") // Product detail partials are loaded from the application root.
	ctx := context.Background()
	cat, err := q.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{Name: "Gallery", Slug: "gallery", Description: "Fixture", Icon: "image"})
	if err != nil {
		t.Fatal(err)
	}
	makeProduct := func(sku string) sqlc.Product {
		p, err := q.CreateProduct(ctx, sqlc.CreateProductParams{Name: sku, Sku: sku, Slug: sku, Description: "Fixture", Status: "published", CategoryID: cat.ID})
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	p, other := makeProduct("gallery-one"), makeProduct("gallery-other")
	makeImage := func(productID int64, path string, primary bool) sqlc.ProductImage {
		img, err := q.CreateProductImage(ctx, sqlc.CreateProductImageParams{ProductID: productID, ImagePath: path, IsThumbnail: primary})
		if err != nil {
			t.Fatal(err)
		}
		return img
	}
	first := makeImage(p.ID, "/first.png", true)
	second := makeImage(p.ID, "/second.png", false)
	foreign := makeImage(other.ID, "/foreign.png", true)
	cache := services.NewCache()
	h := admin.NewProductDetailsHandler(q, logger, services.NewUploadService(t.TempDir()), cache)
	e := echo.New()
	e.POST("/products/:id/images/:image_id/primary", h.SetPrimaryImage)
	e.POST("/products/:id/images", h.AddImage)
	request := func(path string, want int) {
		t.Helper()
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, nil))
		if rec.Code != want {
			t.Fatalf("%s: %d: %s", path, rec.Code, rec.Body.String())
		}
	}
	assertPrimary := func(productID, wantID int64) {
		t.Helper()
		imgs, err := q.ListProductImages(ctx, productID)
		if err != nil {
			t.Fatal(err)
		}
		var ids []int64
		for _, img := range imgs {
			if img.IsThumbnail {
				ids = append(ids, img.ID)
			}
		}
		if len(ids) != 1 || ids[0] != wantID {
			t.Fatalf("product %d primary IDs %v, want %d", productID, ids, wantID)
		}
	}
	cache.Set("page:products:gallery:gallery-one", "stale", 1800)
	request(fmt.Sprintf("/products/%d/images/%d/primary", p.ID, second.ID), 200)
	assertPrimary(p.ID, second.ID)
	if _, ok := cache.Get("page:products:gallery:gallery-one"); ok {
		t.Fatal("primary selection left stale public cache")
	}
	for _, id := range []string{fmt.Sprint(foreign.ID), "9999999", "bad", "-1"} {
		request(fmt.Sprintf("/products/%d/images/%s/primary", p.ID, id), 404)
		assertPrimary(p.ID, second.ID)
		assertPrimary(other.ID, foreign.ID)
	}
	request(fmt.Sprintf("/products/bad/images/%d/primary", first.ID), 404)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("is_thumbnail", "1")
	part, err := writer.CreateFormFile("image", "primary.png")
	if err != nil {
		t.Fatal(err)
	}
	if err = png.Encode(part, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	writer.Close()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/products/%d/images", p.ID), &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("upload: %d %s", rec.Code, rec.Body.String())
	}
	imgs, err := q.ListProductImages(ctx, p.ID)
	if err != nil {
		t.Fatal(err)
	}
	var uploadedID int64
	for _, img := range imgs {
		if strings.Contains(img.ImagePath, "primary.png") {
			uploadedID = img.ID
		}
	}
	if uploadedID == 0 {
		t.Fatal("uploaded image missing")
	}
	assertPrimary(p.ID, uploadedID)
	assertPrimary(other.ID, foreign.ID)
}
