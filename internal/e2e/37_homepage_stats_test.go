package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHomepageStatsCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/homepage/stats", strings.NewReader(url.Values{
		"stat_value":    {"500+"},
		"stat_label":    {"Happy Clients"},
		"display_order": {"1"},
		"is_active":     {"on"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	stats, _ := queries.ListAllStats(ctx)
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat, got %d", len(stats))
	}
	if stats[0].StatValue != "500+" {
		t.Errorf("expected stat_value '500+', got %q", stats[0].StatValue)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/homepage/stats/%d", stats[0].ID), strings.NewReader(url.Values{
		"stat_value":    {"99%"},
		"stat_label":    {"Satisfaction Rate"},
		"display_order": {"2"},
		"is_active":     {""},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("update: expected 303, got %d", rec.Code)
	}

	updatedStat, _ := queries.GetStat(ctx, stats[0].ID)
	if updatedStat.StatValue != "99%" {
		t.Errorf("expected stat_value '99%%', got %q", updatedStat.StatValue)
	}
	if updatedStat.IsActive != 0 {
		t.Errorf("expected is_active 0, got %d", updatedStat.IsActive)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/homepage/stats/%d", stats[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	stats, _ = queries.ListAllStats(ctx)
	if len(stats) != 0 {
		t.Errorf("expected 0 stats after delete, got %d", len(stats))
	}
}

func TestHomepageStatsVariousFormats_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	testCases := []struct {
		value string
		label string
	}{
		{"10+", "Years Experience"},
		{"98.5%", "Uptime"},
		{"1000-5000", "Range"},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodPost, "/admin/homepage/stats", strings.NewReader(url.Values{
			"stat_value":    {tc.value},
			"stat_label":    {tc.label},
			"display_order": {"1"},
			"is_active":     {"on"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("create %q: expected 303, got %d", tc.value, rec.Code)
		}
	}

	stats, _ := queries.ListAllStats(ctx)
	if len(stats) != len(testCases) {
		t.Errorf("expected %d stats, got %d", len(testCases), len(stats))
	}
}
