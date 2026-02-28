package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRateLimit_ContactForm_WithinLimit(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	formData := url.Values{
		"name":    {"John Doe"},
		"email":   {"john@example.com"},
		"phone":   {"555-1234"},
		"company": {"ACME"},
		"message": {"Test"},
	}

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, rec.Code)
		}
	}
}

func TestRateLimit_ContactForm_ExceedsLimit(t *testing.T) {
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
				t.Errorf("request %d: expected 200, got %d", i+1, rec.Code)
			}
		} else {
			if rec.Code != http.StatusTooManyRequests {
				t.Errorf("request %d: expected 429, got %d", i+1, rec.Code)
			}
			if !strings.Contains(rec.Body.String(), "Too many requests") {
				t.Error("expected rate limit error message")
			}
		}
	}
}

func TestRateLimit_InvalidSubmissionsCountTowardLimit(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	invalidFormData := url.Values{
		"name": {"John"},
	}

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(invalidFormData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusTooManyRequests {
			if rec.Code != http.StatusBadRequest {
				t.Errorf("request %d: expected 400 (invalid) or 429 (rate limited), got %d", i+1, rec.Code)
			}
		}
	}
}

func TestRateLimit_GETRequestsNotLimited(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/contact", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, rec.Code)
		}
	}
}
