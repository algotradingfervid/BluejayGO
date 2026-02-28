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

func TestSolutionChallengesAdd(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/challenges", sol.ID), strings.NewReader(url.Values{
		"title":         {"Data Silos"},
		"description":   {"Legacy systems prevent data sharing"},
		"icon":          {"storage"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}

	challenges, _ := queries.GetSolutionChallenges(ctx, sol.ID)
	if len(challenges) != 1 {
		t.Fatalf("expected 1 challenge, got %d", len(challenges))
	}
	if challenges[0].Title != "Data Silos" {
		t.Errorf("expected 'Data Silos', got %q", challenges[0].Title)
	}
}

func TestSolutionChallengesDelete(t *testing.T) {
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

	challenge, _ := queries.CreateSolutionChallenge(ctx, sqlc.CreateSolutionChallengeParams{
		SolutionID:  sol.ID,
		Title:       "Test Challenge",
		Description: "Test description",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/solutions/%d/challenges/%d", sol.ID, challenge.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	challenges, _ := queries.GetSolutionChallenges(ctx, sol.ID)
	if len(challenges) != 0 {
		t.Errorf("expected 0 challenges after delete, got %d", len(challenges))
	}
}

func TestSolutionChallengesMultiple(t *testing.T) {
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

	challengeData := []struct{ title, desc, icon string }{
		{"Data Silos", "Legacy systems", "storage"},
		{"Integration", "System compatibility", "link"},
		{"Scalability", "Growth concerns", "trending_up"},
	}

	for i, data := range challengeData {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/challenges", sol.ID), strings.NewReader(url.Values{
			"title":         {data.title},
			"description":   {data.desc},
			"icon":          {data.icon},
			"display_order": {fmt.Sprintf("%d", i+1)},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
			t.Errorf("add challenge %d: expected 200 or 500, got %d", i, rec.Code)
		}
	}

	challenges, _ := queries.GetSolutionChallenges(ctx, sol.ID)
	if len(challenges) != 3 {
		t.Errorf("expected 3 challenges, got %d", len(challenges))
	}
}
