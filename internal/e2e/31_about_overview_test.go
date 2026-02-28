package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestAboutOverview_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	queries.UpsertCompanyOverview(ctx, sqlc.UpsertCompanyOverviewParams{
		Headline:             "Leading Security Solutions",
		Tagline:              "Protecting what matters most",
		DescriptionMain:      "We provide comprehensive security solutions.",
		DescriptionSecondary: sql.NullString{String: "With over 20 years of experience.", Valid: true},
		DescriptionTertiary:  sql.NullString{String: "Trusted by organizations worldwide.", Valid: true},
		HeroImageUrl:         sql.NullString{String: "https://example.com/hero.jpg", Valid: true},
		CompanyImageUrl:      sql.NullString{String: "https://example.com/company.jpg", Valid: true},
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/about/overview", strings.NewReader(url.Values{
		"headline":              {"Updated Security Solutions"},
		"tagline":               {"Excellence in security"},
		"description_main":      {"We offer industry-leading security."},
		"description_secondary": {"Serving clients globally since 2000."},
		"description_tertiary":  {"Innovation at our core."},
		"hero_image_url":        {"https://example.com/new-hero.jpg"},
		"company_image_url":     {"https://example.com/new-company.jpg"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	overview, _ := queries.GetCompanyOverview(ctx)
	if overview.Headline != "Updated Security Solutions" {
		t.Errorf("expected 'Updated Security Solutions', got %q", overview.Headline)
	}
	if overview.Tagline != "Excellence in security" {
		t.Errorf("expected 'Excellence in security', got %q", overview.Tagline)
	}
}

func TestAboutOverviewLoad_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	queries.UpsertCompanyOverview(ctx, sqlc.UpsertCompanyOverviewParams{
		Headline:        "Test Headline",
		Tagline:         "Test Tagline",
		DescriptionMain: "Test Description",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/about/overview", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about overview route not found")
	}
}
