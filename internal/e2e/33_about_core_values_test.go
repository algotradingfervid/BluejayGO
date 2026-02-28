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

func TestCoreValuesCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	req := httptest.NewRequest(http.MethodPost, "/admin/about/values", strings.NewReader(url.Values{
		"title":         {"Integrity"},
		"description":   {"We operate with honesty and transparency in all our interactions."},
		"icon":          {"verified"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	values, _ := queries.ListCoreValues(ctx)
	if len(values) != 1 {
		t.Fatalf("expected 1 core value, got %d", len(values))
	}
	if values[0].Title != "Integrity" {
		t.Errorf("expected 'Integrity', got %q", values[0].Title)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/about/values/%d", values[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	_, err := queries.GetCoreValue(ctx, values[0].ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestCoreValuesList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	for i := 1; i <= 3; i++ {
		queries.CreateCoreValue(ctx, sqlc.CreateCoreValueParams{
			Title:        fmt.Sprintf("Value %d", i),
			Description:  fmt.Sprintf("Description %d", i),
			Icon:         "icon",
			DisplayOrder: int64(i),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/about/values", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("core values list route not found")
	}
}

func TestCoreValueEdit_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	value, _ := queries.CreateCoreValue(ctx, sqlc.CreateCoreValueParams{
		Title:        "Original Title",
		Description:  "Original description",
		Icon:         "shield",
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/about/values/%d", value.ID), strings.NewReader(url.Values{
		"title":         {"Updated Title"},
		"description":   {"Updated description"},
		"icon":          {"verified"},
		"display_order": {"2"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetCoreValue(ctx, value.ID)
	if updated.Title != "Updated Title" {
		t.Errorf("expected 'Updated Title', got %q", updated.Title)
	}
}
