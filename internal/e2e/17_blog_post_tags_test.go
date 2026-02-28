package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestBlogPostTags_CreateWithTags(t *testing.T) {
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
	tag1, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "JavaScript", Slug: "javascript"})
	tag2, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Go", Slug: "go"})

	formData := url.Values{
		"title":       {"Post with Tags"},
		"slug":        {"post-with-tags"},
		"excerpt":     {"Test excerpt"},
		"body":        {"Test body"},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}
	formData.Add("tag_ids", fmt.Sprintf("%d", tag1.ID))
	formData.Add("tag_ids", fmt.Sprintf("%d", tag2.ID))

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/posts", strings.NewReader(formData.Encode()))
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

	postTags, _ := queries.GetPostTagsByPostID(ctx, posts[0].ID)
	if len(postTags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(postTags))
	}
}

func TestBlogPostTags_UpdateTags(t *testing.T) {
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
	tag1, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "JavaScript", Slug: "javascript"})
	tag2, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Go", Slug: "go"})
	tag3, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Python", Slug: "python"})

	post, _ := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title: "Original", Slug: "original", Excerpt: "exc", Body: "body", CategoryID: cat.ID, AuthorID: author.ID, Status: "draft",
	})
	queries.AddTagToPost(ctx, sqlc.AddTagToPostParams{BlogPostID: post.ID, BlogTagID: tag1.ID})

	formData := url.Values{
		"title":       {"Updated Post"},
		"slug":        {"updated-post"},
		"excerpt":     {"Updated excerpt"},
		"body":        {"Updated body"},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}
	formData.Add("tag_ids", fmt.Sprintf("%d", tag2.ID))
	formData.Add("tag_ids", fmt.Sprintf("%d", tag3.ID))

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/blog/posts/%d", post.ID), strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	postTags, _ := queries.GetPostTagsByPostID(ctx, post.ID)
	if len(postTags) != 2 {
		t.Errorf("expected 2 tags after update, got %d", len(postTags))
	}

	foundTag1 := false
	for _, pt := range postTags {
		if pt.ID == tag1.ID {
			foundTag1 = true
		}
	}
	if foundTag1 {
		t.Error("expected tag1 to be removed after update")
	}
}

func TestBlogPostTags_RemoveAllTags(t *testing.T) {
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
	tag1, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "JavaScript", Slug: "javascript"})

	post, _ := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title: "Original", Slug: "original", Excerpt: "exc", Body: "body", CategoryID: cat.ID, AuthorID: author.ID, Status: "draft",
	})
	queries.AddTagToPost(ctx, sqlc.AddTagToPostParams{BlogPostID: post.ID, BlogTagID: tag1.ID})

	formData := url.Values{
		"title":       {"Updated Post"},
		"slug":        {"updated-post"},
		"excerpt":     {"Updated excerpt"},
		"body":        {"Updated body"},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/blog/posts/%d", post.ID), strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	postTags, _ := queries.GetPostTagsByPostID(ctx, post.ID)
	if len(postTags) != 0 {
		t.Errorf("expected 0 tags after removing all, got %d", len(postTags))
	}
}

func TestBlogPostTags_DeletePostClearsTags(t *testing.T) {
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
	tag1, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "JavaScript", Slug: "javascript"})

	post, _ := queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title: "To Delete", Slug: "to-delete", Excerpt: "exc", Body: "body", CategoryID: cat.ID, AuthorID: author.ID, Status: "draft",
	})
	queries.AddTagToPost(ctx, sqlc.AddTagToPostParams{BlogPostID: post.ID, BlogTagID: tag1.ID})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/blog/posts/%d", post.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	postTags, _ := queries.GetPostTagsByPostID(ctx, post.ID)
	if len(postTags) != 0 {
		t.Errorf("expected 0 tags after post deletion, got %d", len(postTags))
	}
}

func TestBlogPostTags_TagSearch(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "JavaScript", Slug: "javascript"})
	queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Java", Slug: "java"})

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags/search?_tag_search=java", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

}

func TestBlogPostTags_QuickCreateFromPostForm(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags/quick-create", strings.NewReader(url.Values{
		"name": {"NewTag"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	tags, _ := queries.ListAllBlogTags(ctx)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Name != "NewTag" {
		t.Errorf("expected tag name 'NewTag', got %q", tags[0].Name)
	}
}

func TestBlogPostTags_MultipleTagsOrder(t *testing.T) {
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
	tag1, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Tag1", Slug: "tag1"})
	tag2, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Tag2", Slug: "tag2"})
	tag3, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Tag3", Slug: "tag3"})

	formData := url.Values{
		"title":       {"Multi Tag Post"},
		"slug":        {"multi-tag-post"},
		"excerpt":     {"Test excerpt"},
		"body":        {"Test body"},
		"category_id": {fmt.Sprintf("%d", cat.ID)},
		"author_id":   {fmt.Sprintf("%d", author.ID)},
		"status":      {"draft"},
	}
	formData.Add("tag_ids", fmt.Sprintf("%d", tag1.ID))
	formData.Add("tag_ids", fmt.Sprintf("%d", tag2.ID))
	formData.Add("tag_ids", fmt.Sprintf("%d", tag3.ID))

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/posts", strings.NewReader(formData.Encode()))
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

	postTags, _ := queries.GetPostTagsByPostID(ctx, posts[0].ID)
	if len(postTags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(postTags))
	}
}
