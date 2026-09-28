package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicBlogListing(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:        "News",
		Slug:        "news",
		ColorHex:    "#FF5733",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:        "John Doe",
		Slug:        "john-doe",
		Title:       "Writer",
		Bio:         sql.NullString{},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{String: "john@test.com", Valid: true},
		SortOrder:   0,
	})

	_, err := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:              "First Post",
		Slug:               "first-post",
		Excerpt:            "This is the first post",
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
	if err != nil {
		t.Fatalf("create blog post: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/blog", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("blog listing route not found")
	}
}

func TestPublicBlogListing_CategoryFilter(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:        "Tech",
		Slug:        "tech",
		ColorHex:    "#0000FF",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:        "Jane Doe",
		Slug:        "jane-doe",
		Title:       "Tech Writer",
		Bio:         sql.NullString{String: "Tech writer", Valid: true},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{String: "jane@test.com", Valid: true},
		SortOrder:   0,
	})

	_, _ = queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:              "Tech Post",
		Slug:               "tech-post",
		Excerpt:            "Tech content",
		Body:               "Full tech content",
		FeaturedImageUrl:   sql.NullString{},
		FeaturedImageAlt:   sql.NullString{},
		CategoryID:         cat.ID,
		AuthorID:           author.ID,
		MetaDescription:    sql.NullString{},
		ReadingTimeMinutes: sql.NullInt64{},
		Status:             "published",
		PublishedAt:        sql.NullTime{},
	})

	req := httptest.NewRequest(http.MethodGet, "/blog?category=tech", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("blog category filter route not found")
	}
}

func TestPublicBlogListing_Pagination(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:        "Updates",
		Slug:        "updates",
		ColorHex:    "#00FF00",
		Description: sql.NullString{},
		SortOrder:   0,
	})

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:        "Test Author",
		Slug:        "test-author",
		Title:       "Author",
		Bio:         sql.NullString{},
		AvatarUrl:   sql.NullString{},
		LinkedinUrl: sql.NullString{},
		Email:       sql.NullString{},
		SortOrder:   0,
	})

	for i := 1; i <= 15; i++ {
		queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
			Title:              "Post " + string(rune(i)),
			Slug:               "post-" + string(rune(i)),
			Excerpt:            "Excerpt",
			Body:               "Content",
			FeaturedImageUrl:   sql.NullString{},
			FeaturedImageAlt:   sql.NullString{},
			CategoryID:         cat.ID,
			AuthorID:           author.ID,
			MetaDescription:    sql.NullString{},
			ReadingTimeMinutes: sql.NullInt64{},
			Status:             "published",
			PublishedAt:        sql.NullTime{},
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/blog?page=2", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("blog pagination route not found")
	}
}
