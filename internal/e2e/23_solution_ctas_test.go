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

func TestSolutionCTAsAdd(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/ctas", sol.ID), strings.NewReader(url.Values{
		"heading":              {"Get Started Today"},
		"subheading":           {"Transform your business"},
		"primary_button_text":  {"Contact Sales"},
		"primary_button_url":   {"/contact"},
		"secondary_button_text": {"Learn More"},
		"secondary_button_url":  {"/learn"},
		"phone_number":         {"1-800-123-4567"},
		"section_name":         {"hero-cta"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}

	ctas, _ := queries.GetSolutionCTAs(ctx, sol.ID)
	if len(ctas) != 1 {
		t.Fatalf("expected 1 CTA, got %d", len(ctas))
	}
	if ctas[0].Heading != "Get Started Today" {
		t.Errorf("expected 'Get Started Today', got %q", ctas[0].Heading)
	}
	if ctas[0].SectionName != "hero-cta" {
		t.Errorf("expected 'hero-cta', got %q", ctas[0].SectionName)
	}
}

func TestSolutionCTAsDelete(t *testing.T) {
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

	cta, _ := queries.CreateSolutionCTA(ctx, sqlc.CreateSolutionCTAParams{
		SolutionID:  sol.ID,
		Heading:     "Test CTA",
		SectionName: "test-section",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/solutions/%d/ctas/%d", sol.ID, cta.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	ctas, _ := queries.GetSolutionCTAs(ctx, sol.ID)
	if len(ctas) != 0 {
		t.Errorf("expected 0 CTAs after delete, got %d", len(ctas))
	}
}

func TestSolutionCTAsMultiple(t *testing.T) {
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

	ctaData := []struct{ heading, section string }{
		{"Hero CTA", "hero-cta"},
		{"Footer CTA", "footer-cta"},
	}

	for _, data := range ctaData {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/ctas", sol.ID), strings.NewReader(url.Values{
			"heading":      {data.heading},
			"section_name": {data.section},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
			t.Errorf("add CTA: expected 200 or 500, got %d", rec.Code)
		}
	}

	ctas, _ := queries.GetSolutionCTAs(ctx, sol.ID)
	if len(ctas) != 2 {
		t.Errorf("expected 2 CTAs, got %d", len(ctas))
	}
}

func TestSolutionCTAsOptionalFields(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/ctas", sol.ID), strings.NewReader(url.Values{
		"heading":      {"Minimal CTA"},
		"section_name": {"minimal"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}

	ctas, _ := queries.GetSolutionCTAs(ctx, sol.ID)
	if len(ctas) != 1 {
		t.Fatalf("expected 1 CTA, got %d", len(ctas))
	}
}
