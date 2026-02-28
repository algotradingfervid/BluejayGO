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

func TestContactSubmissionsList(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:           "John Doe",
		Email:          "john@example.com",
		Phone:          "1234567890",
		Company:        "Test Co",
		Message:        "Test message",
		InquiryType:    sql.NullString{String: "contact", Valid: true},
		IpAddress:      sql.NullString{},
		UserAgent:      sql.NullString{},
	})

	t.Run("list all submissions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/contact/submissions", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("filter by type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/contact/submissions?type=contact", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/contact/submissions?status=new", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("search submissions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/contact/submissions?search=john", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestContactSubmissionDetail(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	submission, _ := queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:           "Jane Doe",
		Email:          "jane@example.com",
		Phone:          "0987654321",
		Company:        "Jane Co",
		Message:        "Test detail message",
		InquiryType:    sql.NullString{String: "rfq", Valid: true},
		IpAddress:      sql.NullString{},
		UserAgent:      sql.NullString{},
	})

	t.Run("view submission detail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/contact/submissions/%d", submission.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestContactSubmissionStatusUpdate(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	submission, _ := queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:           "Status Test",
		Email:          "status@example.com",
		Phone:          "1111111111",
		Company:        "Status Co",
		Message:        "Status test message",
		InquiryType:    sql.NullString{String: "contact", Valid: true},
		IpAddress:      sql.NullString{},
		UserAgent:      sql.NullString{},
	})

	t.Run("update status to reviewed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/contact/submissions/%d/status", submission.ID), strings.NewReader(url.Values{
			"status": {"reviewed"},
			"notes":  {"Reviewed by admin"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		updated, _ := queries.GetContactSubmissionByID(ctx, submission.ID)
		if updated.Status != "reviewed" {
			t.Errorf("expected status 'reviewed', got %q", updated.Status)
		}
	})
}

func TestContactSubmissionDelete(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	submission, _ := queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:        "Delete Test",
		Email:       "delete@example.com",
		Phone:       "2222222222",
		Company:     "Delete Co",
		Message:     "Delete test message",
		InquiryType: sql.NullString{String: "contact", Valid: true},
		IpAddress:   sql.NullString{},
		UserAgent:   sql.NullString{},
	})

	t.Run("delete submission", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/contact/submissions/%d", submission.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
			t.Errorf("expected 200 or 204, got %d", rec.Code)
		}

		_, err := queries.GetContactSubmissionByID(ctx, submission.ID)
		if err != sql.ErrNoRows {
			t.Error("expected submission to be deleted")
		}
	})
}

func TestContactBulkMarkReviewed(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:           "Bulk Test 1",
		Email:          "bulk1@example.com",
		Phone:          "3333333333",
		Company:        "Bulk Co",
		Message:        "Bulk test 1",
		InquiryType:    sql.NullString{String: "contact", Valid: true},
		IpAddress:      sql.NullString{},
		UserAgent:      sql.NullString{},
	})

	queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:           "Bulk Test 2",
		Email:          "bulk2@example.com",
		Phone:          "4444444444",
		Company:        "Bulk Co",
		Message:        "Bulk test 2",
		InquiryType:    sql.NullString{String: "contact", Valid: true},
		IpAddress:      sql.NullString{},
		UserAgent:      sql.NullString{},
	})

	t.Run("bulk mark as reviewed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/contact/submissions/bulk-mark-read", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther && rec.Code != http.StatusOK {
			t.Errorf("expected 303 or 200, got %d", rec.Code)
		}

		newCount, _ := queries.CountContactSubmissionsByStatus(ctx, "new")
		if newCount != 0 {
			t.Errorf("expected 0 new submissions, got %d", newCount)
		}
	})
}
