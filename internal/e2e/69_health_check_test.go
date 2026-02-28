package e2e_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthCheck_ReturnsOK(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHealthCheck_ReturnsJSON(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("expected JSON content type, got %s", contentType)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", resp["status"])
	}
}

func TestHealthCheck_ContainsTimestamp(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["time"] == "" {
		t.Error("expected time field to be present")
	}

	_, err := time.Parse(time.RFC3339, resp["time"])
	if err != nil {
		t.Errorf("time field not in RFC3339 format: %v", err)
	}
}

func TestHealthCheck_MultipleRequests(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, rec.Code)
		}

		var resp map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Errorf("request %d: failed to decode JSON: %v", i+1, err)
		}

		if resp["status"] != "ok" {
			t.Errorf("request %d: expected status 'ok', got %q", i+1, resp["status"])
		}
	}
}

func TestHealthCheck_NoAuthRequired(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Error("health check should not require authentication")
	}
}

func TestHealthCheck_TimeIsRecent(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	before := time.Now()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	after := time.Now()

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	healthTime, err := time.Parse(time.RFC3339, resp["time"])
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}

	tolerance := 1 * time.Second
	if healthTime.Before(before.Add(-tolerance)) || healthTime.After(after.Add(tolerance)) {
		t.Errorf("health check time %v not between %v and %v (with 1s tolerance)", healthTime, before, after)
	}
}
