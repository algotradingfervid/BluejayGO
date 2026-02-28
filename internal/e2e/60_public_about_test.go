package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicAboutPage(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.UpsertCompanyOverview(ctx, sqlc.UpsertCompanyOverviewParams{
		Headline:             "Leading Innovator",
		Tagline:              "Innovation at its best",
		DescriptionMain:      "We are a leading company in the industry",
		DescriptionSecondary: sql.NullString{},
		DescriptionTertiary:  sql.NullString{},
		HeroImageUrl:         sql.NullString{},
		CompanyImageUrl:      sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create company overview: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about page route not found")
	}
}

func TestPublicAboutPage_WithMissionVisionValues(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.UpsertMissionVisionValues(ctx, sqlc.UpsertMissionVisionValuesParams{
		Mission:       "To innovate and deliver excellence",
		Vision:        "A world transformed by technology",
		ValuesSummary: sql.NullString{String: "Integrity, Innovation, Customer Focus", Valid: true},
		MissionIcon:   "target",
		VisionIcon:    "eye",
		ValuesIcon:    "heart",
	})
	if err != nil {
		t.Fatalf("create MVV: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about page with MVV route not found")
	}
}

func TestPublicAboutPage_WithCoreValues(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.CreateCoreValue(ctx, sqlc.CreateCoreValueParams{
		Title:        "Innovation",
		Description:  "We constantly innovate",
		Icon:         "lightbulb",
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create core value: %v", err)
	}

	_, err = queries.CreateCoreValue(ctx, sqlc.CreateCoreValueParams{
		Title:        "Integrity",
		Description:  "We operate with integrity",
		Icon:         "shield",
		DisplayOrder: 2,
	})
	if err != nil {
		t.Fatalf("create core value 2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about page with core values route not found")
	}
}

func TestPublicAboutPage_WithMilestones(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		Year:         2015,
		Title:        "Company Founded",
		Description:  "Started operations",
		IsCurrent:    0,
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create milestone: %v", err)
	}

	_, err = queries.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		Year:         2020,
		Title:        "Global Expansion",
		Description:  "Expanded to 20 countries",
		IsCurrent:    0,
		DisplayOrder: 2,
	})
	if err != nil {
		t.Fatalf("create milestone 2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about page with milestones route not found")
	}
}

func TestPublicAboutPage_WithCertifications(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.CreateCertification(ctx, sqlc.CreateCertificationParams{
		Name:         "ISO 9001",
		Abbreviation: "ISO",
		Description:  sql.NullString{String: "Quality management certification", Valid: true},
		Icon:         sql.NullString{},
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create certification: %v", err)
	}

	_, err = queries.CreateCertification(ctx, sqlc.CreateCertificationParams{
		Name:         "SOC 2",
		Abbreviation: "AICPA",
		Description:  sql.NullString{String: "Security compliance", Valid: true},
		Icon:         sql.NullString{},
		DisplayOrder: 2,
	})
	if err != nil {
		t.Fatalf("create certification 2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about page with certifications route not found")
	}
}

func TestPublicAboutPage_EmptyContent(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about page should handle empty content gracefully")
	}
}

func TestPublicAboutPage_CompleteData(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, _ = queries.UpsertCompanyOverview(ctx, sqlc.UpsertCompanyOverviewParams{
		Headline:             "Industry Leader",
		Tagline:              "Excellence in every detail",
		DescriptionMain:      "Comprehensive company description",
		DescriptionSecondary: sql.NullString{},
		DescriptionTertiary:  sql.NullString{},
		HeroImageUrl:         sql.NullString{},
		CompanyImageUrl:      sql.NullString{},
	})

	_, _ = queries.UpsertMissionVisionValues(ctx, sqlc.UpsertMissionVisionValuesParams{
		Mission:       "Deliver excellence",
		Vision:        "Transform the industry",
		ValuesSummary: sql.NullString{String: "Quality, Trust, Innovation", Valid: true},
		MissionIcon:   "target",
		VisionIcon:    "eye",
		ValuesIcon:    "heart",
	})

	_, _ = queries.CreateCoreValue(ctx, sqlc.CreateCoreValueParams{
		Title:        "Quality",
		Description:  "Quality in everything we do",
		Icon:         "check",
		DisplayOrder: 1,
	})

	_, _ = queries.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		Year:         2018,
		Title:        "Major Achievement",
		Description:  "Reached important milestone",
		IsCurrent:    0,
		DisplayOrder: 1,
	})

	_, _ = queries.CreateCertification(ctx, sqlc.CreateCertificationParams{
		Name:         "ISO 27001",
		Abbreviation: "ISO",
		Description:  sql.NullString{},
		Icon:         sql.NullString{},
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about page with complete data route not found")
	}
}
