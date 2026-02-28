package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestActivityLogList(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	queries.CreateActivityLog(ctx, sqlc.CreateActivityLogParams{
		UserID:        sql.NullInt64{},
		Action:        "create",
		ResourceType:  "test",
		ResourceID:    sql.NullInt64{},
		ResourceTitle: sql.NullString{String: "Test Resource", Valid: true},
		Description:   "Created test resource",
	})

	t.Run("load activity log", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/activity", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("filter by action", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/activity?action=create", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("search activity logs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/activity?search=test", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("paginate activity logs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/activity?page=1", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestActivityLogFilters(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	queries.CreateActivityLog(ctx, sqlc.CreateActivityLogParams{
		UserID:        sql.NullInt64{},
		Action:        "update",
		ResourceType:  "product",
		ResourceID:    sql.NullInt64{},
		ResourceTitle: sql.NullString{String: "Alpha Product", Valid: true},
		Description:   "Updated product",
	})

	queries.CreateActivityLog(ctx, sqlc.CreateActivityLogParams{
		UserID:        sql.NullInt64{},
		Action:        "delete",
		ResourceType:  "product",
		ResourceID:    sql.NullInt64{},
		ResourceTitle: sql.NullString{String: "Beta Product", Valid: true},
		Description:   "Deleted product",
	})

	t.Run("filter update actions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/activity?action=update", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("filter delete actions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/activity?action=delete", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("combined filters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/activity?action=create&search=admin", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestActivityLogReadOnly(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("verify no create endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/activity", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 404 or 405, got %d", rec.Code)
		}
	})

	t.Run("verify no delete endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/activity/1", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 404 or 405, got %d", rec.Code)
		}
	})
}
