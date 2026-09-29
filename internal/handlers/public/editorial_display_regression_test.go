package public_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestProductInitialImageSelectionAndFallback(t *testing.T) {
	e, q, cache, cleanup := setupProductsHandler(t)
	defer cleanup()
	e.Renderer = templates.NewRenderer("../../../templates")
	ctx := context.Background()
	cat := createTestCategory(t, q, "Gallery", "gallery")
	p := createTestProduct(t, q, "GALLERY", "gallery-product", "Gallery Product", cat.ID, "published")
	updateMedia := func(path string) {
		t.Helper()
		err := q.UpdateProduct(ctx, sqlc.UpdateProductParams{ID: p.ID, Name: p.Name, Sku: p.Sku, Slug: p.Slug, CategoryID: cat.ID, Status: "published", Description: "Fixture", PrimaryImage: sql.NullString{String: path, Valid: path != ""}})
		if err != nil {
			t.Fatal(err)
		}
	}
	render := func(want string) {
		t.Helper()
		cache.DeleteByPrefix("page:products")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/products/gallery/gallery-product", nil))
		if rec.Code != 200 {
			t.Fatalf("render %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if want == "" {
			if !strings.Contains(body, `id="main-image-placeholder"`) {
				t.Fatal("no-image fallback missing")
			}
		} else if !strings.Contains(body, `id="main-image" src="`+want+`"`) {
			t.Fatalf("initial photo did not use %s", want)
		}
	}
	render("")
	first, err := q.CreateProductImage(ctx, sqlc.CreateProductImageParams{ProductID: p.ID, ImagePath: "/first.png", DisplayOrder: 1})
	if err != nil {
		t.Fatal(err)
	}
	render("/first.png")
	updateMedia("/media.png")
	render("/media.png")
	selected, err := q.CreateProductImage(ctx, sqlc.CreateProductImageParams{ProductID: p.ID, ImagePath: "/selected.png", DisplayOrder: 3, IsThumbnail: true, AltText: sql.NullString{String: "Selected product photo", Valid: true}})
	if err != nil {
		t.Fatal(err)
	}
	render("/selected.png")
	if _, err = q.SetPrimaryProductImage(ctx, sqlc.SetPrimaryProductImageParams{ProductID: p.ID, ImageID: first.ID}); err != nil {
		t.Fatal(err)
	}
	render("/first.png")
	if err = q.DeleteProductImage(ctx, first.ID); err != nil {
		t.Fatal(err)
	}
	render("/media.png")
	updateMedia("")
	render(selected.ImagePath)
}

func TestPreviewBannerReflectsSavedProductBlogAndCaseStudyStatus(t *testing.T) {
	db, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	cache := services.NewCache()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("X-Test-Admin") == "yes" {
				c.Set("session", &middleware.Session{UserID: 1})
			}
			return next(c)
		}
	})
	e.GET("/products/:category/:slug", public.NewProductsHandler(q, logger, services.NewProductService(q), cache).ProductDetail)
	e.GET("/blog/:slug", public.NewBlogHandler(q, logger, cache).BlogPost)
	e.GET("/case-studies/:slug", public.NewCaseStudiesHandler(q, logger, cache).CaseStudyDetail)
	cat := createTestCategory(t, q, "Preview", "preview")
	bc, err := q.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{Name: "Preview", Slug: "preview", ColorHex: "#000000"})
	if err != nil {
		t.Fatal(err)
	}
	author, err := q.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{Name: "Preview author", Slug: "preview-author"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := db.Exec(`INSERT INTO industries(name,slug,icon,description) VALUES('Preview','preview','test','test')`)
	if err != nil {
		t.Fatal(err)
	}
	industryID, _ := result.LastInsertId()
	for _, status := range []string{"draft", "published"} {
		slug := "preview-" + status
		createTestProduct(t, q, slug, slug, slug, cat.ID, status)
		_, err = q.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{Title: slug, Slug: slug, Excerpt: "Preview", Body: "<p>Preview story</p>", CategoryID: bc.ID, AuthorID: author.ID, Status: status, PublishedAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: status == "published"}})
		if err != nil {
			t.Fatal(err)
		}
		published := int64(0)
		if status == "published" {
			published = 1
		}
		_, err = q.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{Title: slug, Slug: slug, ClientName: "Fixture", IndustryID: industryID, Summary: "Preview", IsPublished: published})
		if err != nil {
			t.Fatal(err)
		}
		for _, path := range []string{"/products/preview/" + slug, "/blog/" + slug, "/case-studies/" + slug} {
			req := httptest.NewRequest(http.MethodGet, path+"?preview=true", nil)
			req.Header.Set("X-Test-Admin", "yes")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Saved status: "+status) {
				t.Fatalf("%s preview status missing: HTTP %d %s", path, rec.Code, rec.Body.String())
			}
			if strings.Contains(rec.Body.String(), "This content is not published yet") {
				t.Fatal("old unconditional unpublished banner remains")
			}
			if status == "draft" {
				rec = httptest.NewRecorder()
				e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path+"?preview=true", nil))
				if rec.Code != 404 {
					t.Fatalf("unauthenticated draft preview %s returned %d", path, rec.Code)
				}
			}
		}
	}
}
