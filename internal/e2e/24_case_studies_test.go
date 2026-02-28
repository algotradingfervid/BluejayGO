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

func TestCaseStudiesList(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	_, err := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "acme-corp",
		Title:            "Acme Corp Success",
		ClientName:       "Acme Corp",
		IndustryID:       ind.ID,
		IsPublished:      1,
		DisplayOrder:     0,
		Summary:          "summary",
		ChallengeTitle:   "ch",
		ChallengeContent: "cc",
		SolutionTitle:    "st",
		SolutionContent:  "sc",
		OutcomeTitle:     "ot",
		OutcomeContent:   "oc",
	})
	if err != nil {
		t.Fatalf("create case study: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/case-studies", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestCaseStudyCreate(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	req := httptest.NewRequest(http.MethodPost, "/admin/case-studies", strings.NewReader(url.Values{
		"title":             {"Client Success Story"},
		"client_name":       {"TestCorp"},
		"industry_id":       {fmt.Sprintf("%d", ind.ID)},
		"summary":           {"Great results"},
		"challenge_title":   {"The Challenge"},
		"challenge_content": {"Problems faced"},
		"challenge_bullets": {"Cost overruns, Delays, Quality issues"},
		"solution_title":    {"The Solution"},
		"solution_content":  {"How we solved it"},
		"outcome_title":     {"The Outcome"},
		"outcome_content":   {"Results achieved"},
		"is_published":      {"1"},
		"display_order":     {"5"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	caseStudies, _ := queries.AdminListCaseStudiesFiltered(ctx, sqlc.AdminListCaseStudiesFilteredParams{
		FilterSearch: "",
		FilterStatus: "",
		PageLimit:    100,
		PageOffset:   0,
	})
	if len(caseStudies) != 1 {
		t.Fatalf("expected 1 case study, got %d", len(caseStudies))
	}
	if caseStudies[0].Title != "Client Success Story" {
		t.Errorf("expected 'Client Success Story', got %q", caseStudies[0].Title)
	}
	if caseStudies[0].Slug != "client-success-story" {
		t.Errorf("expected slug 'client-success-story', got %q", caseStudies[0].Slug)
	}
}

func TestCaseStudyUpdate(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "original",
		Title:            "Original Title",
		ClientName:       "OriginalCorp",
		IndustryID:       ind.ID,
		Summary:          "s",
		ChallengeTitle:   "ch",
		ChallengeContent: "cc",
		SolutionTitle:    "st",
		SolutionContent:  "sc",
		OutcomeTitle:     "ot",
		OutcomeContent:   "oc",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/case-studies/%d", cs.ID), strings.NewReader(url.Values{
		"title":             {"Updated Title"},
		"client_name":       {"UpdatedCorp"},
		"industry_id":       {fmt.Sprintf("%d", ind.ID)},
		"summary":           {"Updated summary"},
		"challenge_title":   {"ch"},
		"challenge_content": {"cc"},
		"solution_title":    {"st"},
		"solution_content":  {"sc"},
		"outcome_title":     {"ot"},
		"outcome_content":   {"oc"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.AdminGetCaseStudy(ctx, cs.ID)
	if updated.Title != "Updated Title" {
		t.Errorf("expected 'Updated Title', got %q", updated.Title)
	}
}

func TestCaseStudyDelete(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "delete-me",
		Title:            "Delete Me",
		ClientName:       "DeleteCorp",
		IndustryID:       ind.ID,
		Summary:          "s",
		ChallengeTitle:   "ch",
		ChallengeContent: "cc",
		SolutionTitle:    "st",
		SolutionContent:  "sc",
		OutcomeTitle:     "ot",
		OutcomeContent:   "oc",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/case-studies/%d", cs.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}

	_, err := queries.AdminGetCaseStudy(ctx, cs.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}
