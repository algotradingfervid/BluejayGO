package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestContactPage_Loads(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/contact", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestContactForm_Submit_Success(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	formData := url.Values{
		"name":         {"John Doe"},
		"email":        {"john@example.com"},
		"phone":        {"555-1234"},
		"company":      {"ACME Corp"},
		"message":      {"Test inquiry"},
		"inquiry_type": {"Sales"},
	}

	req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Thank you") {
		t.Error("expected success message")
	}

	submissions, _ := queries.ListContactSubmissions(ctx, sqlc.ListContactSubmissionsParams{
		Limit:  100,
		Offset: 0,
	})
	if len(submissions) != 1 {
		t.Errorf("expected 1 submission, got %d", len(submissions))
	}
}

func TestContactForm_Submit_MissingFields(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	formData := url.Values{
		"name": {"John Doe"},
	}

	req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "required") {
		t.Error("expected validation error message")
	}
}

func TestContactForm_RateLimit(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	formData := url.Values{
		"name":    {"John Doe"},
		"email":   {"john@example.com"},
		"phone":   {"555-1234"},
		"company": {"ACME"},
		"message": {"Test"},
	}

	for i := 0; i < 6; i++ {
		req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if i < 5 {
			if rec.Code != http.StatusOK {
				t.Errorf("submission %d: expected 200, got %d", i+1, rec.Code)
			}
		} else {
			if rec.Code != http.StatusTooManyRequests {
				t.Errorf("submission %d: expected 429, got %d", i+1, rec.Code)
			}
		}
	}
}

func TestContactForm_WithOfficeLocations(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, _ = queries.CreateOfficeLocation(ctx, sqlc.CreateOfficeLocationParams{
		Name:         "HQ Office",
		AddressLine1: "123 Main St",
		AddressLine2: sql.NullString{},
		City:         "New York",
		State:        "NY",
		PostalCode:   "10001",
		Country:      "USA",
		Phone:        sql.NullString{String: "555-0000", Valid: true},
		Email:        sql.NullString{},
		IsPrimary:    1,
	})

	req := httptest.NewRequest(http.MethodGet, "/contact", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
