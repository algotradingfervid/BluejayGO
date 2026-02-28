package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicCaseStudiesListing(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Manufacturing",
		Slug: "manufacturing",
	})

	_, err := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "success-story",
		Title:            "Success Story",
		ClientName:       "Acme Corp",
		IndustryID:       ind.ID,
		HeroImageUrl:     sql.NullString{},
		Summary:          "Great results achieved",
		ChallengeTitle:   "Challenge",
		ChallengeContent: "Business challenge",
		ChallengeBullets: sql.NullString{},
		SolutionTitle:    "Solution",
		SolutionContent:  "Our solution",
		OutcomeTitle:     "Outcome",
		OutcomeContent:   "Positive outcome",
		MetaTitle:        sql.NullString{},
		MetaDescription:  sql.NullString{},
		IsPublished:      1,
		DisplayOrder:     0,
	})
	if err != nil {
		t.Fatalf("create case study: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/case-studies", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("case studies listing route not found")
	}
}

func TestPublicCaseStudiesListing_IndustryFilter(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Healthcare",
		Slug: "healthcare",
	})

	_, _ = queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "hospital-success",
		Title:            "Hospital Success",
		ClientName:       "City Hospital",
		IndustryID:       ind.ID,
		HeroImageUrl:     sql.NullString{},
		Summary:          "Improved patient care",
		ChallengeTitle:   "Challenge",
		ChallengeContent: "Challenge text",
		ChallengeBullets: sql.NullString{},
		SolutionTitle:    "Solution",
		SolutionContent:  "Solution text",
		OutcomeTitle:     "Outcome",
		OutcomeContent:   "Outcome text",
		MetaTitle:        sql.NullString{},
		MetaDescription:  sql.NullString{},
		IsPublished:      1,
		DisplayOrder:     0,
	})

	req := httptest.NewRequest(http.MethodGet, "/case-studies?industry="+strconv.FormatInt(ind.ID, 10), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("case studies industry filter route not found")
	}
}

func TestPublicCaseStudiesListing_InvalidIndustry(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/case-studies?industry=invalid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid industry, got %d", rec.Code)
	}
}

func TestPublicCaseStudyDetail(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Retail",
		Slug: "retail",
	})

	_, err := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "retail-transformation",
		Title:            "Retail Transformation",
		ClientName:       "Super Mart",
		IndustryID:       ind.ID,
		HeroImageUrl:     sql.NullString{},
		Summary:          "Digital transformation",
		ChallengeTitle:   "Challenge",
		ChallengeContent: "Legacy systems",
		ChallengeBullets: sql.NullString{},
		SolutionTitle:    "Solution",
		SolutionContent:  "Modern platform",
		OutcomeTitle:     "Outcome",
		OutcomeContent:   "Increased sales",
		MetaTitle:        sql.NullString{},
		MetaDescription:  sql.NullString{},
		IsPublished:      1,
		DisplayOrder:     0,
	})
	if err != nil {
		t.Fatalf("create case study: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/case-studies/retail-transformation", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("case study detail route not found")
	}
}

func TestPublicCaseStudyDetail_NotFound(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/case-studies/nonexistent", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestPublicCaseStudyDetail_WithMetrics(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Logistics",
		Slug: "logistics",
	})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "logistics-optimization",
		Title:            "Logistics Optimization",
		ClientName:       "Global Shipping",
		IndustryID:       ind.ID,
		HeroImageUrl:     sql.NullString{},
		Summary:          "Optimized operations",
		ChallengeTitle:   "Challenge",
		ChallengeContent: "Inefficient routing",
		ChallengeBullets: sql.NullString{},
		SolutionTitle:    "Solution",
		SolutionContent:  "Smart routing",
		OutcomeTitle:     "Outcome",
		OutcomeContent:   "Cost savings",
		MetaTitle:        sql.NullString{},
		MetaDescription:  sql.NullString{},
		IsPublished:      1,
		DisplayOrder:     0,
	})

	_, err := queries.AdminCreateMetric(ctx, sqlc.AdminCreateMetricParams{
		CaseStudyID:  cs.ID,
		MetricValue:  "25%",
		MetricLabel:  "Cost Reduction",
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create metric: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/case-studies/logistics-optimization", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("case study with metrics route not found")
	}
}

func TestPublicCaseStudyDetail_WithProducts(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Energy",
		Slug: "energy",
	})

	cs, _ := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "energy-efficiency",
		Title:            "Energy Efficiency",
		ClientName:       "Power Plant",
		IndustryID:       ind.ID,
		HeroImageUrl:     sql.NullString{},
		Summary:          "Reduced energy consumption",
		ChallengeTitle:   "Challenge",
		ChallengeContent: "High energy costs",
		ChallengeBullets: sql.NullString{},
		SolutionTitle:    "Solution",
		SolutionContent:  "Monitoring system",
		OutcomeTitle:     "Outcome",
		OutcomeContent:   "30% savings",
		MetaTitle:        sql.NullString{},
		MetaDescription:  sql.NullString{},
		IsPublished:      1,
		DisplayOrder:     0,
	})

	prodCat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Monitors",
		Slug:        "monitors",
		Description: "Monitoring systems",
		Icon:        "monitor",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "MON-001",
		Slug:        "energy-monitor",
		Name:        "Energy Monitor",
		Description: "Real-time energy monitoring",
		CategoryID:  prodCat.ID,
		Status:      "published",
	})

	_, err := queries.AdminAddCaseStudyProduct(ctx, sqlc.AdminAddCaseStudyProductParams{
		CaseStudyID:  cs.ID,
		ProductID:    prod.ID,
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("link product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/case-studies/energy-efficiency", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("case study with products route not found")
	}
}

func TestPublicCaseStudyDetail_WithChallengeBullets(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Finance",
		Slug: "finance",
	})

	bullets := []string{"Compliance requirements", "Legacy infrastructure", "Data security"}
	bulletsJSON, _ := json.Marshal(bullets)

	_, err := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "financial-compliance",
		Title:            "Financial Compliance",
		ClientName:       "Bank Corp",
		IndustryID:       ind.ID,
		HeroImageUrl:     sql.NullString{},
		Summary:          "Achieved compliance",
		ChallengeTitle:   "Challenge",
		ChallengeContent: "Regulatory challenges",
		ChallengeBullets: sql.NullString{String: string(bulletsJSON), Valid: true},
		SolutionTitle:    "Solution",
		SolutionContent:  "Compliance platform",
		OutcomeTitle:     "Outcome",
		OutcomeContent:   "Full compliance",
		MetaTitle:        sql.NullString{},
		MetaDescription:  sql.NullString{},
		IsPublished:      1,
		DisplayOrder:     0,
	})
	if err != nil {
		t.Fatalf("create case study with bullets: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/case-studies/financial-compliance", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("case study with challenge bullets route not found")
	}
}

func TestPublicCaseStudyDetail_PreviewMode(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	ind, _ := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Technology",
		Slug: "technology",
	})

	_, err := queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "draft-case-study",
		Title:            "Draft Case Study",
		ClientName:       "Tech Startup",
		IndustryID:       ind.ID,
		HeroImageUrl:     sql.NullString{},
		Summary:          "Draft summary",
		ChallengeTitle:   "Challenge",
		ChallengeContent: "Draft challenge",
		ChallengeBullets: sql.NullString{},
		SolutionTitle:    "Solution",
		SolutionContent:  "Draft solution",
		OutcomeTitle:     "Outcome",
		OutcomeContent:   "Draft outcome",
		MetaTitle:        sql.NullString{},
		MetaDescription:  sql.NullString{},
		IsPublished:      0,
		DisplayOrder:     0,
	})
	if err != nil {
		t.Fatalf("create draft case study: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/case-studies/draft-case-study", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("draft should return 404, got %d", rec.Code)
	}

	// Preview mode requires authentication + ?preview=true
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req = httptest.NewRequest(http.MethodGet, "/case-studies/draft-case-study?preview=true", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("preview mode should show draft case studies")
	}
}
