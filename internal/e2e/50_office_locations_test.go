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

func TestOfficeLocationsList(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("list office locations", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/contact/offices", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestOfficeLocationCreate(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("create office location", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/contact/offices", strings.NewReader(url.Values{
			"name":          {"Bangalore HQ"},
			"address_line1": {"123 Tech Park"},
			"city":          {"Bangalore"},
			"state":         {"Karnataka"},
			"country":       {"India"},
			"phone":         {"+91-1234567890"},
			"email":         {"bangalore@example.com"},
			"is_primary":    {"on"},
			"is_active":     {"on"},
			"display_order": {"1"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		offices, _ := queries.ListAllOfficeLocations(context.Background())
		if len(offices) == 0 {
			t.Error("expected office to be created")
		}
	})
}

func TestOfficeLocationUpdate(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	office, _ := queries.CreateOfficeLocation(ctx, sqlc.CreateOfficeLocationParams{
		Name:         "Mumbai Office",
		AddressLine1: "456 Business Center",
		AddressLine2: sql.NullString{},
		City:         "Mumbai",
		State:        "Maharashtra",
		PostalCode:   "",
		Country:      "India",
		Phone:        sql.NullString{String: "+91-9876543210", Valid: true},
		Email:        sql.NullString{String: "mumbai@example.com", Valid: true},
		IsPrimary:    0,
		IsActive:     1,
		DisplayOrder: 2,
	})

	t.Run("update office location", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/contact/offices/%d", office.ID), strings.NewReader(url.Values{
			"name":          {"Mumbai HQ"},
			"address_line1": {"789 Updated Center"},
			"city":          {"Mumbai"},
			"state":         {"Maharashtra"},
			"country":       {"India"},
			"phone":         {"+91-9999999999"},
			"email":         {"mumbai-hq@example.com"},
			"is_primary":    {},
			"is_active":     {"on"},
			"display_order": {"1"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		updated, _ := queries.GetOfficeLocationByID(ctx, office.ID)
		if updated.Name != "Mumbai HQ" {
			t.Errorf("expected 'Mumbai HQ', got %q", updated.Name)
		}
	})
}

func TestOfficePrimaryEnforcement(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	office1, _ := queries.CreateOfficeLocation(ctx, sqlc.CreateOfficeLocationParams{
		Name:         "Office 1",
		AddressLine1: "Address 1",
		AddressLine2: sql.NullString{},
		City:         "City 1",
		State:        "",
		PostalCode:   "",
		Country:      "India",
		Phone:        sql.NullString{},
		Email:        sql.NullString{},
		IsPrimary:    1,
		IsActive:     1,
		DisplayOrder: 0,
	})

	office2, _ := queries.CreateOfficeLocation(ctx, sqlc.CreateOfficeLocationParams{
		Name:         "Office 2",
		AddressLine1: "Address 2",
		AddressLine2: sql.NullString{},
		City:         "City 2",
		State:        "",
		PostalCode:   "",
		Country:      "India",
		Phone:        sql.NullString{},
		Email:        sql.NullString{},
		IsPrimary:    0,
		IsActive:     1,
		DisplayOrder: 0,
	})

	t.Run("set new primary office", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/contact/offices/%d", office2.ID), strings.NewReader(url.Values{
			"name":          {"Office 2"},
			"address_line1": {"Address 2"},
			"city":          {"City 2"},
			"country":       {"India"},
			"is_primary":    {"on"},
			"is_active":     {"on"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		updated1, _ := queries.GetOfficeLocationByID(ctx, office1.ID)
		updated2, _ := queries.GetOfficeLocationByID(ctx, office2.ID)

		if updated1.IsPrimary != 0 {
			t.Error("expected office 1 to no longer be primary")
		}
		if updated2.IsPrimary != 1 {
			t.Error("expected office 2 to be primary")
		}
	})
}

func TestOfficeLocationDelete(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	office, _ := queries.CreateOfficeLocation(ctx, sqlc.CreateOfficeLocationParams{
		Name:         "Delete Test Office",
		AddressLine1: "Delete Address",
		AddressLine2: sql.NullString{},
		City:         "Delete City",
		State:        "",
		PostalCode:   "",
		Country:      "India",
		Phone:        sql.NullString{},
		Email:        sql.NullString{},
		IsPrimary:    0,
		IsActive:     1,
		DisplayOrder: 0,
	})

	t.Run("delete office location", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/contact/offices/%d", office.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
			t.Errorf("expected 200 or 204, got %d", rec.Code)
		}

		_, err := queries.GetOfficeLocationByID(ctx, office.ID)
		if err != sql.ErrNoRows {
			t.Error("expected office to be deleted")
		}
	})
}
