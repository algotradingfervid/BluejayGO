package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicHomepage(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	queries.CreateHero(ctx, sqlc.CreateHeroParams{
		Headline:         "Welcome to Bluejay",
		Subheadline:      "Industry-leading solutions",
		BadgeText:        sql.NullString{String: "Featured", Valid: true},
		PrimaryCtaText:   "",
		PrimaryCtaUrl:    "",
		SecondaryCtaText: sql.NullString{},
		SecondaryCtaUrl:  sql.NullString{},
		BackgroundImage:  sql.NullString{},
		IsActive:         1,
		DisplayOrder:     0,
	})

	queries.CreateStat(ctx, sqlc.CreateStatParams{
		StatLabel:    "Customers",
		StatValue:    "500+",
		DisplayOrder: 0,
		IsActive:     1,
	})

	t.Run("homepage loads successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("homepage returns HTML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		contentType := rec.Header().Get("Content-Type")
		if contentType != "text/html; charset=UTF-8" && rec.Result().ContentLength == 0 {
			t.Error("expected HTML content")
		}
	})
}

func TestHomepageSections(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "Electronic products",
		Icon:        "chip",
		SortOrder:   1,
	})

	queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PROD-001",
		Slug:        "test-product",
		Name:        "Test Product",
		Description: "A test product",
		CategoryID:  cat.ID,
		Status:      "published",
		IsFeatured:  true,
	})

	t.Run("homepage with featured products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestHomepageWithTestimonials(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	queries.CreateTestimonialHomepage(ctx, sqlc.CreateTestimonialHomepageParams{
		Quote:         "Great product!",
		AuthorName:    "John Doe",
		AuthorTitle:   sql.NullString{},
		AuthorCompany: sql.NullString{String: "Test Co", Valid: true},
		AuthorImage:   sql.NullString{},
		Rating:        5,
		DisplayOrder:  0,
		IsActive:      1,
	})

	t.Run("homepage with testimonials", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestHomepageWithBlogPosts(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	ctx := context.Background()

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name: "Test Author",
		Slug: "test-author",
	})

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "News",
		Slug: "news",
	})

	queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:              "Latest News",
		Slug:               "latest-news",
		Excerpt:            "Read about our latest news",
		Body:               "Full content here",
		FeaturedImageUrl:   sql.NullString{},
		FeaturedImageAlt:   sql.NullString{},
		CategoryID:         cat.ID,
		AuthorID:           author.ID,
		MetaDescription:    sql.NullString{},
		ReadingTimeMinutes: sql.NullInt64{},
		Status:             "published",
		PublishedAt:        sql.NullTime{},
	})

	t.Run("homepage with blog posts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}
