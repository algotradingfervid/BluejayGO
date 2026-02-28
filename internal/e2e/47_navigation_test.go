package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestNavigationList(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("list navigation menus", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/navigation", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestNavigationCreate(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("create navigation menu", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/navigation", strings.NewReader(url.Values{
			"name":     {"Test Menu"},
			"location": {"header"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		menus, _ := queries.ListNavigationMenus(context.Background())
		if len(menus) == 0 {
			t.Error("expected menu to be created")
		}
	})
}

func TestNavigationEdit(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	menu, _ := queries.CreateNavigationMenu(ctx, sqlc.CreateNavigationMenuParams{
		Name:     "Edit Test Menu",
		Location: "footer",
	})

	t.Run("view menu editor", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/navigation/%d", menu.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("update menu settings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/navigation/%d/settings", menu.ID), strings.NewReader(url.Values{
			"name":     {"Updated Menu"},
			"location": {"sidebar"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		updated, _ := queries.GetNavigationMenu(ctx, menu.ID)
		if updated.Name != "Updated Menu" {
			t.Errorf("expected 'Updated Menu', got %q", updated.Name)
		}
	})
}

func TestNavigationAddItem(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	menu, _ := queries.CreateNavigationMenu(ctx, sqlc.CreateNavigationMenuParams{
		Name:     "Item Test Menu",
		Location: "header",
	})

	t.Run("add page link item", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/navigation/%d/items", menu.ID), strings.NewReader(url.Values{
			"label":           {"Products"},
			"link_type":       {"page"},
			"page_identifier": {"Products"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		items, _ := queries.ListNavigationItems(ctx, menu.ID)
		if len(items) == 0 {
			t.Error("expected item to be created")
		}
	})

	t.Run("add url link item", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/navigation/%d/items", menu.ID), strings.NewReader(url.Values{
			"label":     {"External"},
			"link_type": {"url"},
			"url":       {"https://example.com"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}
	})
}

func TestNavigationDeleteItem(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	menu, _ := queries.CreateNavigationMenu(ctx, sqlc.CreateNavigationMenuParams{
		Name:     "Delete Test Menu",
		Location: "header",
	})

	item, _ := queries.CreateNavigationItem(ctx, sqlc.CreateNavigationItemParams{
		MenuID:     menu.ID,
		Label:      "Delete Me",
		LinkType:   "url",
		Url:        sql.NullString{String: "/test", Valid: true},
		IsActive:   sql.NullInt64{Int64: 1, Valid: true},
		SortOrder:  sql.NullInt64{Int64: 0, Valid: true},
	})

	t.Run("delete navigation item", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/navigation/items/%d", item.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		_, err := queries.GetNavigationItem(ctx, item.ID)
		if err != sql.ErrNoRows {
			t.Error("expected item to be deleted")
		}
	})
}

func TestNavigationReorder(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	menu, _ := queries.CreateNavigationMenu(ctx, sqlc.CreateNavigationMenuParams{
		Name:     "Reorder Test Menu",
		Location: "header",
	})

	item1, _ := queries.CreateNavigationItem(ctx, sqlc.CreateNavigationItemParams{
		MenuID:    menu.ID,
		Label:     "Item 1",
		LinkType:  "url",
		Url:       sql.NullString{String: "/1", Valid: true},
		IsActive:  sql.NullInt64{Int64: 1, Valid: true},
		SortOrder: sql.NullInt64{Int64: 0, Valid: true},
	})

	item2, _ := queries.CreateNavigationItem(ctx, sqlc.CreateNavigationItemParams{
		MenuID:    menu.ID,
		Label:     "Item 2",
		LinkType:  "url",
		Url:       sql.NullString{String: "/2", Valid: true},
		IsActive:  sql.NullInt64{Int64: 1, Valid: true},
		SortOrder: sql.NullInt64{Int64: 1, Valid: true},
	})

	t.Run("reorder items", func(t *testing.T) {
		reorderData := []map[string]interface{}{
			{"id": item2.ID, "parent_id": nil, "order": 0},
			{"id": item1.ID, "parent_id": nil, "order": 1},
		}
		body, _ := json.Marshal(reorderData)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/navigation/%d/reorder", menu.ID), strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestNavigationDeleteMenu(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	menu, _ := queries.CreateNavigationMenu(ctx, sqlc.CreateNavigationMenuParams{
		Name:     "Delete Menu Test",
		Location: "footer",
	})

	t.Run("delete menu", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/navigation/%d", menu.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		_, err := queries.GetNavigationMenu(ctx, menu.ID)
		if err != sql.ErrNoRows {
			t.Error("expected menu to be deleted")
		}
	})
}
