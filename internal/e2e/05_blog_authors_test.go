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

func TestBlogAuthorCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog-authors", strings.NewReader(url.Values{
		"name":         {"John Doe"},
		"title":        {"Senior Editor"},
		"bio":          {"Expert in technology"},
		"avatar_url":   {"https://example.com/avatar.jpg"},
		"linkedin_url": {"https://linkedin.com/in/johndoe"},
		"email":        {"john@example.com"},
		"sort_order":   {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d; body: %s", rec.Code, rec.Body.String())
	}

	authors, _ := queries.ListBlogAuthors(context.Background())
	if len(authors) != 1 {
		t.Fatalf("expected 1 author, got %d", len(authors))
	}

	if authors[0].Name != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", authors[0].Name)
	}

	if authors[0].Slug != "john-doe" {
		t.Errorf("expected slug 'john-doe', got %q", authors[0].Slug)
	}

	if authors[0].Title != "Senior Editor" {
		t.Errorf("expected title 'Senior Editor', got %q", authors[0].Title)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/blog-authors/%d", authors[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	authors, _ = queries.ListBlogAuthors(context.Background())
	if len(authors) != 0 {
		t.Errorf("expected 0 authors after delete, got %d", len(authors))
	}
}

func TestBlogAuthorUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	author, _ := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{
		Name:        "Jane Smith",
		Slug:        "jane-smith",
		Title:       "Editor",
		Bio:         sql.NullString{String: "Writer", Valid: true},
		AvatarUrl:   sql.NullString{String: "https://example.com/jane.jpg", Valid: true},
		LinkedinUrl: sql.NullString{String: "", Valid: false},
		Email:       sql.NullString{String: "jane@example.com", Valid: true},
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/blog-authors/%d", author.ID), strings.NewReader(url.Values{
		"name":         {"Jane Smith Updated"},
		"title":        {"Senior Editor"},
		"bio":          {"Updated bio"},
		"avatar_url":   {"https://example.com/jane-new.jpg"},
		"linkedin_url": {"https://linkedin.com/in/janesmith"},
		"email":        {"jane.new@example.com"},
		"sort_order":   {"2"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetBlogAuthor(ctx, author.ID)
	if updated.Name != "Jane Smith Updated" {
		t.Errorf("expected 'Jane Smith Updated', got %q", updated.Name)
	}

	if updated.Slug != "jane-smith-updated" {
		t.Errorf("expected slug 'jane-smith-updated', got %q", updated.Slug)
	}

	if updated.Title != "Senior Editor" {
		t.Errorf("expected title 'Senior Editor', got %q", updated.Title)
	}

	if !updated.LinkedinUrl.Valid || updated.LinkedinUrl.String != "https://linkedin.com/in/janesmith" {
		t.Errorf("expected linkedin_url 'https://linkedin.com/in/janesmith', got %v", updated.LinkedinUrl)
	}
}

func TestBlogAuthorOptionalFields_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog-authors", strings.NewReader(url.Values{
		"name":       {"Minimal Author"},
		"title":      {"Writer"},
		"sort_order": {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	authors, _ := queries.ListBlogAuthors(context.Background())
	if len(authors) != 1 {
		t.Fatalf("expected 1 author, got %d", len(authors))
	}

	if authors[0].Bio.Valid {
		t.Errorf("expected bio to be NULL, got %v", authors[0].Bio)
	}

	if authors[0].AvatarUrl.Valid {
		t.Errorf("expected avatar_url to be NULL, got %v", authors[0].AvatarUrl)
	}

	if authors[0].LinkedinUrl.Valid {
		t.Errorf("expected linkedin_url to be NULL, got %v", authors[0].LinkedinUrl)
	}

	if authors[0].Email.Valid {
		t.Errorf("expected email to be NULL, got %v", authors[0].Email)
	}
}
