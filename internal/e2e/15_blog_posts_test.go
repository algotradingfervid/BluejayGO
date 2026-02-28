package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestBlogPosts_List(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/posts", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestBlogPosts_CreateAndList(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "Tech", Slug: "tech", ColorHex: "#000000", SortOrder: 1,
	})
	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name: "Author", Slug: "author", Title: "Writer", SortOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/posts", strings.NewReader(url.Values{
		"title":       {"Test Post"},
		"slug":        {"test-post"},
		"excerpt":     {"Test excerpt"},
		"body":        {"Test body content"},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	posts, _ := queries.ListBlogPostsAdminFiltered(ctx, sqlc.ListBlogPostsAdminFilteredParams{
		FilterStatus: "", FilterCategory: int64(0), FilterAuthor: int64(0), FilterSearch: "",
		PageLimit: 15, PageOffset: 0,
	})
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if posts[0].Title != "Test Post" {
		t.Errorf("expected title 'Test Post', got %q", posts[0].Title)
	}
}

func TestBlogPosts_Update(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "Tech", Slug: "tech", ColorHex: "#000000", SortOrder: 1,
	})
	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name: "Author", Slug: "author", Title: "Writer", SortOrder: 1,
	})
	post, _ := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title: "Original", Slug: "original", Excerpt: "exc", Body: "body", CategoryID: cat.ID, AuthorID: author.ID, Status: "draft",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/blog/posts/%d", post.ID), strings.NewReader(url.Values{
		"title":       {"Updated Post"},
		"slug":        {"updated-post"},
		"excerpt":     {"Updated excerpt"},
		"body":        {"Updated body"},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetBlogPost(ctx, post.ID)
	if updated.Title != "Updated Post" {
		t.Errorf("expected title 'Updated Post', got %q", updated.Title)
	}
}

func TestBlogPosts_Delete(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "Tech", Slug: "tech", ColorHex: "#000000", SortOrder: 1,
	})
	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name: "Author", Slug: "author", Title: "Writer", SortOrder: 1,
	})
	post, _ := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title: "To Delete", Slug: "to-delete", Excerpt: "exc", Body: "body", CategoryID: cat.ID, AuthorID: author.ID, Status: "draft",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/blog/posts/%d", post.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	_, err := queries.GetBlogPost(ctx, post.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestBlogPosts_RequiresAuth(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	paths := []string{
		"/admin/blog/posts",
		"/admin/blog/posts/new",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != http.StatusSeeOther {
				t.Errorf("expected 303 redirect, got %d", rec.Code)
			}
		})
	}
}

func TestBlogPosts_ListWithFilters(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "Tech", Slug: "tech", ColorHex: "#000000", SortOrder: 1,
	})
	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name: "Author", Slug: "author", Title: "Writer", SortOrder: 1,
	})
	queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title: "Draft Post", Slug: "draft-post", Excerpt: "exc", Body: "body", CategoryID: cat.ID, AuthorID: author.ID, Status: "draft",
	})
	queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title: "Published Post", Slug: "published-post", Excerpt: "exc", Body: "body", CategoryID: cat.ID, AuthorID: author.ID, Status: "published",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/posts?status=draft", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestBlogPosts_AutoSlugGeneration(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "Tech", Slug: "tech", ColorHex: "#000000", SortOrder: 1,
	})
	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name: "Author", Slug: "author", Title: "Writer", SortOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/posts", strings.NewReader(url.Values{
		"title":       {"My Blog Post"},
		"excerpt":     {"Test excerpt"},
		"body":        {"Test body"},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	posts, _ := queries.ListBlogPostsAdminFiltered(ctx, sqlc.ListBlogPostsAdminFilteredParams{
		FilterStatus: "", FilterCategory: int64(0), FilterAuthor: int64(0), FilterSearch: "",
		PageLimit: 15, PageOffset: 0,
	})
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if posts[0].Slug != "my-blog-post" {
		t.Errorf("expected slug 'my-blog-post', got %q", posts[0].Slug)
	}
}

func TestBlogPosts_ReadingTimeCalculation(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "Tech", Slug: "tech", ColorHex: "#000000", SortOrder: 1,
	})
	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name: "Author", Slug: "author", Title: "Writer", SortOrder: 1,
	})

	longBody := strings.Repeat("word ", 250)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/posts", strings.NewReader(url.Values{
		"title":       {"Long Post"},
		"excerpt":     {"Test excerpt"},
		"body":        {longBody},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	posts, _ := queries.ListBlogPostsAdminFiltered(ctx, sqlc.ListBlogPostsAdminFilteredParams{
		FilterStatus: "", FilterCategory: int64(0), FilterAuthor: int64(0), FilterSearch: "",
		PageLimit: 15, PageOffset: 0,
	})
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if !posts[0].ReadingTimeMinutes.Valid || posts[0].ReadingTimeMinutes.Int64 < 1 {
		t.Error("expected reading time to be calculated")
	}
}
