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

func TestBlogCategoryCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/blog-categories", strings.NewReader(url.Values{
		"name":        {"Industry News"},
		"color_hex":   {"#1E88E5"},
		"description": {"Latest industry updates"},
		"sort_order":  {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d; body: %s", rec.Code, rec.Body.String())
	}

	cats, _ := queries.ListBlogCategories(context.Background())
	if len(cats) != 1 {
		t.Fatalf("expected 1 category, got %d", len(cats))
	}

	if cats[0].Name != "Industry News" {
		t.Errorf("expected 'Industry News', got %q", cats[0].Name)
	}

	if cats[0].Slug != "industry-news" {
		t.Errorf("expected slug 'industry-news', got %q", cats[0].Slug)
	}

	if cats[0].ColorHex != "#1E88E5" {
		t.Errorf("expected color_hex '#1E88E5', got %q", cats[0].ColorHex)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/blog-categories/%d", cats[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	cats, _ = queries.ListBlogCategories(context.Background())
	if len(cats) != 0 {
		t.Errorf("expected 0 categories after delete, got %d", len(cats))
	}
}

func TestBlogCategoryUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name:        "Old Name",
		Slug:        "old-name",
		ColorHex:    "#FF0000",
		Description: sql.NullString{String: "Old desc", Valid: true},
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/blog-categories/%d", cat.ID), strings.NewReader(url.Values{
		"name":        {"Updated Name"},
		"color_hex":   {"#00FF00"},
		"description": {"Updated desc"},
		"sort_order":  {"2"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetBlogCategory(ctx, cat.ID)
	if updated.Name != "Updated Name" {
		t.Errorf("expected 'Updated Name', got %q", updated.Name)
	}

	if updated.Slug != "updated-name" {
		t.Errorf("expected slug 'updated-name', got %q", updated.Slug)
	}

	if updated.ColorHex != "#00FF00" {
		t.Errorf("expected color_hex '#00FF00', got %q", updated.ColorHex)
	}
}
