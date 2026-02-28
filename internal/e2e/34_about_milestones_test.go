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

func TestMilestonesCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	req := httptest.NewRequest(http.MethodPost, "/admin/about/milestones", strings.NewReader(url.Values{
		"year":          {"2008"},
		"title":         {"Company Founded"},
		"description":   {"Started with a vision to revolutionize security."},
		"is_current":    {""},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	milestones, _ := queries.ListMilestones(ctx)
	if len(milestones) != 1 {
		t.Fatalf("expected 1 milestone, got %d", len(milestones))
	}
	if milestones[0].Title != "Company Founded" {
		t.Errorf("expected 'Company Founded', got %q", milestones[0].Title)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/about/milestones/%d", milestones[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	_, err := queries.GetMilestone(ctx, milestones[0].ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestMilestonesList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	for i := 1; i <= 3; i++ {
		queries.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
			Year:         int64(2008 + i),
			Title:        fmt.Sprintf("Milestone %d", i),
			Description:  fmt.Sprintf("Description %d", i),
			IsCurrent:    int64(0),
			DisplayOrder: int64(i),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/about/milestones", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("milestones list route not found")
	}
}

func TestMilestoneEdit_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	milestone, _ := queries.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		Year:         2020,
		Title:        "Original Title",
		Description:  "Original description",
		IsCurrent:    0,
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/about/milestones/%d", milestone.ID), strings.NewReader(url.Values{
		"year":          {"2021"},
		"title":         {"Updated Title"},
		"description":   {"Updated description"},
		"is_current":    {"on"},
		"display_order": {"2"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetMilestone(ctx, milestone.ID)
	if updated.Title != "Updated Title" {
		t.Errorf("expected 'Updated Title', got %q", updated.Title)
	}
	if updated.IsCurrent != 1 {
		t.Errorf("expected is_current=1, got %d", updated.IsCurrent)
	}
}
