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

func TestCaseStudyMetricsAdd(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "test-case",
		Title:            "Test Case Study",
		ClientName:       "TestCorp",
		IndustryID:       ind.ID,
		Summary:          "s",
		ChallengeTitle:   "ch",
		ChallengeContent: "cc",
		SolutionTitle:    "st",
		SolutionContent:  "sc",
		OutcomeTitle:     "ot",
		OutcomeContent:   "oc",
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/case-studies/%d/metrics", cs.ID), strings.NewReader(url.Values{
		"metric_value":  {"40%"},
		"metric_label":  {"Cost Reduction"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}

	metrics, _ := queries.AdminListMetrics(ctx, cs.ID)
	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}
	if metrics[0].MetricValue != "40%" {
		t.Errorf("expected '40%%', got %q", metrics[0].MetricValue)
	}
	if metrics[0].MetricLabel != "Cost Reduction" {
		t.Errorf("expected 'Cost Reduction', got %q", metrics[0].MetricLabel)
	}
}

func TestCaseStudyMetricsDelete(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "test-case",
		Title:            "Test Case Study",
		ClientName:       "TestCorp",
		IndustryID:       ind.ID,
		Summary:          "s",
		ChallengeTitle:   "ch",
		ChallengeContent: "cc",
		SolutionTitle:    "st",
		SolutionContent:  "sc",
		OutcomeTitle:     "ot",
		OutcomeContent:   "oc",
	})

	metric, _ := queries.AdminCreateMetric(ctx, sqlc.AdminCreateMetricParams{
		CaseStudyID: cs.ID,
		MetricValue: "30%",
		MetricLabel: "Efficiency",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/case-studies/%d/metrics/%d", cs.ID, metric.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}

	metrics, _ := queries.AdminListMetrics(ctx, cs.ID)
	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics after delete, got %d", len(metrics))
	}
}

func TestCaseStudyMetricsMultiple(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "test-case",
		Title:            "Test Case Study",
		ClientName:       "TestCorp",
		IndustryID:       ind.ID,
		Summary:          "s",
		ChallengeTitle:   "ch",
		ChallengeContent: "cc",
		SolutionTitle:    "st",
		SolutionContent:  "sc",
		OutcomeTitle:     "ot",
		OutcomeContent:   "oc",
	})

	metricData := []struct{ value, label string }{
		{"40%", "Cost Reduction"},
		{"2x", "Speed Improvement"},
		{"$1M", "Revenue Increase"},
	}

	for i, data := range metricData {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/case-studies/%d/metrics", cs.ID), strings.NewReader(url.Values{
			"metric_value":  {data.value},
			"metric_label":  {data.label},
			"display_order": {fmt.Sprintf("%d", i+1)},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
			t.Errorf("add metric %d: expected 200 or 500, got %d", i, rec.Code)
		}
	}

	metrics, _ := queries.AdminListMetrics(ctx, cs.ID)
	if len(metrics) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(metrics))
	}
}

func TestCaseStudyMetricsVariousFormats(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{Name: "Tech", Slug: "tech", Description: "Technology", Icon: "chip", SortOrder: 1})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "test-case",
		Title:            "Test Case Study",
		ClientName:       "TestCorp",
		IndustryID:       ind.ID,
		Summary:          "s",
		ChallengeTitle:   "ch",
		ChallengeContent: "cc",
		SolutionTitle:    "st",
		SolutionContent:  "sc",
		OutcomeTitle:     "ot",
		OutcomeContent:   "oc",
	})

	formats := []string{"50%", "2x faster", "$1M saved", "99.9% uptime"}

	for i, value := range formats {
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/case-studies/%d/metrics", cs.ID), strings.NewReader(url.Values{
			"metric_value":  {value},
			"metric_label":  {fmt.Sprintf("Metric %d", i+1)},
			"display_order": {fmt.Sprintf("%d", i+1)},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
			t.Errorf("add metric format %q: expected 200 or 500, got %d", value, rec.Code)
		}
	}

	metrics, _ := queries.AdminListMetrics(ctx, cs.ID)
	if len(metrics) != 4 {
		t.Errorf("expected 4 metrics, got %d", len(metrics))
	}
}
