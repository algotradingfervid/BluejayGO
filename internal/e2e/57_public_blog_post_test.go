package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicBlogPost(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:  "Announcements",
		Slug:  "announcements",
		ColorHex: "#FF0000",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:  "Alice Johnson",
		Slug:        "alice",
		Title:       "Writer",
		Bio:         sql.NullString{String: "Product manager", Valid: true},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{String: "alice@test.com", Valid: true},
		SortOrder:   0,
	})

	_, err := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:      "Product Launch",
		Slug:       "product-launch",
		Excerpt:    "New product announcement",
		Body:    "<p>We're excited to announce our new product!</p>",
		CategoryID: cat.ID,
		AuthorID:   author.ID,
		Status:      "published",
		PublishedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		t.Fatalf("create blog post: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/product-launch", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("blog post route not found")
	}
}

func TestPublicBlogPost_NotFound(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/blog/nonexistent-post", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestPublicBlogPost_WithTags(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:  "Tutorials",
		Slug:  "tutorials",
		ColorHex: "#00FF00",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:  "Tutorial Author",
		Slug:        "tutorial",
		Title:       "Writer",
		Bio:         sql.NullString{String: "Teacher", Valid: true},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{String: "tutorial@test.com", Valid: true},
		SortOrder:   0,
	})

	post, _ := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:      "How to Tutorial",
		Slug:       "how-to-tutorial",
		Excerpt:    "Learn the basics",
		Body:    "<p>Tutorial content</p>",
		CategoryID: cat.ID,
		AuthorID:   author.ID,
		Status:      "published",
		PublishedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	tag, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{
		Name: "beginner",
		Slug: "beginner",
	})

	err := queries.AddTagToPost(ctx, sqlc.AddTagToPostParams{
		BlogPostID: post.ID,
		BlogTagID:  tag.ID,
	})
	if err != nil {
		t.Fatalf("link tag: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/how-to-tutorial", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("blog post with tags route not found")
	}
}

func TestPublicBlogPost_WithRelatedProducts(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:  "Reviews",
		Slug:  "reviews",
		ColorHex: "#0000FF",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:  "Reviewer",
		Slug:        "review",
		Title:       "Writer",
		Bio:         sql.NullString{String: "Product reviewer", Valid: true},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{String: "review@test.com", Valid: true},
		SortOrder:   0,
	})

	post, _ := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:      "Product Review",
		Slug:       "product-review",
		Excerpt:    "Detailed review",
		Body:    "<p>Review content</p>",
		CategoryID: cat.ID,
		AuthorID:   author.ID,
		Status:      "published",
		PublishedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	prodCat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Devices",
		Slug:        "devices",
		Description: "Device products",
		Icon:        "device",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "DEV-001",
		Slug:        "smart-device",
		Name:        "Smart Device",
		Description: "Smart device description",
		CategoryID:  prodCat.ID,
		Status:      "published",
	})

	err := queries.AddProductToPost(ctx, sqlc.AddProductToPostParams{
		BlogPostID:   post.ID,
		ProductID:    prod.ID,
		DisplayOrder: sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		t.Fatalf("link product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/product-review", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("blog post with related products route not found")
	}
}

func TestPublicBlogPost_PreviewMode(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:  "Drafts",
		Slug:  "drafts",
		ColorHex: "#CCCCCC",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:  "Draft Author",
		Slug:        "draft",
		Title:       "Writer",
		Bio:         sql.NullString{String: "Draft writer", Valid: true},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{String: "draft@test.com", Valid: true},
		SortOrder:   0,
	})

	_, err := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:      "Draft Post",
		Slug:       "draft-post",
		Excerpt:    "Draft excerpt",
		Body:    "<p>Draft content</p>",
		CategoryID: cat.ID,
		AuthorID:   author.ID,
		Status:     "draft",
	})
	if err != nil {
		t.Fatalf("create draft post: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/draft-post", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("draft post should return 404, got %d", rec.Code)
	}

	// Preview mode requires authentication + ?preview=true
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req = httptest.NewRequest(http.MethodGet, "/blog/draft-post?preview=true", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("preview mode should show draft posts")
	}
}

func TestPublicBlogPost_MetaDescription(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:  "SEO",
		Slug:  "seo",
		ColorHex: "#123456",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:  "SEO Writer",
		Slug:        "seo",
		Title:       "Writer",
		Bio:         sql.NullString{String: "SEO expert", Valid: true},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{String: "seo@test.com", Valid: true},
		SortOrder:   0,
	})

	_, err := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:           "SEO Optimized Post",
		Slug:            "seo-post",
		Excerpt:         "SEO excerpt",
		Body:         "<p>SEO content</p>",
		CategoryID:      cat.ID,
		AuthorID:        author.ID,
		Status:          "published",
		PublishedAt:      sql.NullTime{Time: time.Now(), Valid: true},
		MetaDescription: sql.NullString{String: "Custom SEO meta description", Valid: true},
	})
	if err != nil {
		t.Fatalf("create post with meta: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/blog/seo-post", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("blog post with meta route not found")
	}
}
