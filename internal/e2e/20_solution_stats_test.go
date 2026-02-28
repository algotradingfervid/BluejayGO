package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestSolutionStatsAdd(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		ShortDescription: "desc",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/stats", sol.ID), strings.NewReader(url.Values{
		"value":         {"99%"},
		"label":         {"Uptime"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}

	stats, _ := queries.GetSolutionStats(ctx, sol.ID)
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat, got %d", len(stats))
	}
	if stats[0].Value != "99%" {
		t.Errorf("expected '99%%', got %q", stats[0].Value)
	}
	if stats[0].Label != "Uptime" {
		t.Errorf("expected 'Uptime', got %q", stats[0].Label)
	}
}

func TestSolutionStatsDelete(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		ShortDescription: "desc",
	})

	stat, _ := queries.CreateSolutionStat(ctx, sqlc.CreateSolutionStatParams{
		SolutionID: sol.ID,
		Value:      "50%",
		Label:      "Efficiency",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/solutions/%d/stats/%d", sol.ID, stat.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	stats, _ := queries.GetSolutionStats(ctx, sol.ID)
	if len(stats) != 0 {
		t.Errorf("expected 0 stats after delete, got %d", len(stats))
	}
}

func TestSolutionStatsMultiple(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		ShortDescription: "desc",
	})

	for i, data := range []struct{ value, label string }{
		{"99%", "Uptime"},
		{"10K+", "Users"},
		{"50ms", "Response Time"},
	} {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/stats", sol.ID), strings.NewReader(url.Values{
			"value":         {data.value},
			"label":         {data.label},
			"display_order": {fmt.Sprintf("%d", i+1)},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
			t.Errorf("add stat %d: expected 200 or 500, got %d", i, rec.Code)
		}
	}

	stats, _ := queries.GetSolutionStats(ctx, sol.ID)
	if len(stats) != 3 {
		t.Errorf("expected 3 stats, got %d", len(stats))
	}
}
