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

func TestBlogTags_List(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestBlogTags_Create(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags", strings.NewReader(url.Values{
		"name": {"Technology"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	tags, _ := queries.ListAllBlogTags(ctx)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Name != "Technology" {
		t.Errorf("expected name 'Technology', got %q", tags[0].Name)
	}
	if tags[0].Slug != "technology" {
		t.Errorf("expected slug 'technology', got %q", tags[0].Slug)
	}
}

func TestBlogTags_Delete(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	tag, _ := queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{
		Name: "ToDelete", Slug: "todelete",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/blog/tags/%d", tag.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	_, err := queries.GetBlogTag(ctx, tag.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestBlogTags_SearchEmpty(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags/search?_tag_search=", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestBlogTags_SearchWithResults(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "JavaScript", Slug: "javascript"})
	queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Java", Slug: "java"})
	queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "Python", Slug: "python"})

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags/search?_tag_search=java", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

}

func TestBlogTags_SearchNoResults(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	queries.CreateBlogTag(ctx, sqlc.CreateBlogTagParams{Name: "JavaScript", Slug: "javascript"})

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags/search?_tag_search=xyz123nonexistent", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestBlogTags_QuickCreate(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags/quick-create", strings.NewReader(url.Values{
		"name": {"DevOps"},
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
	if tags[0].Name != "DevOps" {
		t.Errorf("expected name 'DevOps', got %q", tags[0].Name)
	}

}

func TestBlogTags_QuickCreateEmpty(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags/quick-create", strings.NewReader(url.Values{
		"name": {""},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Error("expected error status for empty tag name")
	}
}

func TestBlogTags_RequiresAuth(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	paths := []string{
		"/admin/blog/tags",
		"/admin/blog/tags/search",
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

func TestBlogTags_AutoSlugGeneration(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags", strings.NewReader(url.Values{
		"name": {"Machine Learning"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	tags, _ := queries.ListAllBlogTags(ctx)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Slug != "machine-learning" {
		t.Errorf("expected slug 'machine-learning', got %q", tags[0].Slug)
	}
}
