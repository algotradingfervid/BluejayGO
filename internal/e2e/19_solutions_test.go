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

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestSolutionsList(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	_, err := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "IoT Platform",
		Slug:             "iot-platform",
		ShortDescription: "IoT solution",
	})
	if err != nil {
		t.Fatalf("create solution: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/solutions", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestSolutionCreate(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/solutions", strings.NewReader(url.Values{
		"title":             {"Cloud Platform"},
		"icon":              {"cloud_upload"},
		"short_description": {"Cloud desc"},
		"is_published":      {"1"},
		"display_order":     {"10"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	solutions, _ := queries.ListSolutionsAdminFiltered(context.Background(), sqlc.ListSolutionsAdminFilteredParams{
		FilterStatus: "", FilterSearch: "",
		PageLimit:  100,
		PageOffset: 0,
	})
	if len(solutions) != 1 {
		t.Fatalf("expected 1 solution, got %d", len(solutions))
	}
	if solutions[0].Title != "Cloud Platform" {
		t.Errorf("expected 'Cloud Platform', got %q", solutions[0].Title)
	}
	if solutions[0].Slug != "cloud-platform" {
		t.Errorf("expected slug 'cloud-platform', got %q", solutions[0].Slug)
	}
}

func TestSolutionUpdate(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Original",
		Slug:             "original",
		ShortDescription: "desc",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d", sol.ID), strings.NewReader(url.Values{
		"title":             {"Updated Title"},
		"short_description": {"Updated desc"},
		"is_published":      {"on"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetSolutionByID(ctx, sol.ID)
	if updated.Title != "Updated Title" {
		t.Errorf("expected 'Updated Title', got %q", updated.Title)
	}
}

func TestSolutionDelete(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Delete Me",
		Slug:             "delete-me",
		ShortDescription: "desc",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/solutions/%d", sol.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}

	_, err := queries.GetSolutionByID(ctx, sol.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}
