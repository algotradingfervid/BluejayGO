package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicSolutionDetail(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	sol, err := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Energy Management",
		Slug:             "energy-management",
		ShortDescription: "Optimize energy usage",
		Icon:             "energy",
		HeroImageUrl:     sql.NullString{},
		HeroTitle:        sql.NullString{},
		HeroDescription:  sql.NullString{},
		OverviewContent:  sql.NullString{String: "Comprehensive energy management solution", Valid: true},
		MetaDescription:  sql.NullString{},
		ReferenceCode:    sql.NullString{},
		IsPublished:      sql.NullBool{Bool: true, Valid: true},
		DisplayOrder:     sql.NullInt64{},
	})
	if err != nil {
		t.Fatalf("create solution: %v", err)
	}

	_, err = queries.CreateSolutionStat(ctx, sqlc.CreateSolutionStatParams{
		SolutionID: sol.ID,
		Label:      "Energy Saved",
		Value:      "30%",
		DisplayOrder: sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		t.Fatalf("create stat: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/solutions/energy-management", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("solution detail route not found")
	}
}

func TestPublicSolutionDetail_NotFound(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/solutions/nonexistent", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestPublicSolutionDetail_PreviewMode(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Draft Solution",
		Slug:             "draft-solution",
		ShortDescription: "Draft only",
		Icon:             "draft",
		HeroImageUrl:     sql.NullString{},
		HeroTitle:        sql.NullString{},
		HeroDescription:  sql.NullString{},
		OverviewContent:  sql.NullString{String: "Draft content", Valid: true},
		MetaDescription:  sql.NullString{},
		ReferenceCode:    sql.NullString{},
		IsPublished:      sql.NullBool{Bool: false, Valid: true},
		DisplayOrder:     sql.NullInt64{},
	})
	if err != nil {
		t.Fatalf("create draft solution: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/solutions/draft-solution", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("draft without preview should return 404, got %d", rec.Code)
	}

	// Preview mode requires authentication + ?preview=true
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req = httptest.NewRequest(http.MethodGet, "/solutions/draft-solution?preview=true", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("preview mode should show draft solutions")
	}
}

func TestPublicSolutionDetail_WithProducts(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Sensors",
		Slug:        "sensors",
		Description: "Sensor products",
		Icon:        "sensor",
		ImageUrl:     sql.NullString{},
		SortOrder:    1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "SEN-001",
		Slug:        "temp-sensor",
		Name:        "Temperature Sensor",
		Description: "Measures temperature",
		CategoryID:  cat.ID,
		Status:      "published",
	})

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Cold Chain",
		Slug:             "cold-chain",
		ShortDescription: "Monitor temperature",
		Icon:             "cold",
		HeroImageUrl:     sql.NullString{},
		HeroTitle:        sql.NullString{},
		HeroDescription:  sql.NullString{},
		OverviewContent:  sql.NullString{String: "Full cold chain monitoring", Valid: true},
		MetaDescription:  sql.NullString{},
		ReferenceCode:    sql.NullString{},
		IsPublished:      sql.NullBool{Bool: true, Valid: true},
		DisplayOrder:     sql.NullInt64{},
	})

	err := queries.AddProductToSolution(ctx, sqlc.AddProductToSolutionParams{
		SolutionID: sol.ID,
		ProductID:  prod.ID,
		DisplayOrder: sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		t.Fatalf("link product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/solutions/cold-chain", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("solution with products route not found")
	}
}

func TestPublicSolutionDetail_MetaDescription(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		ShortDescription: "Short desc",
		Icon:             "test",
		HeroImageUrl:     sql.NullString{},
		HeroTitle:        sql.NullString{},
		HeroDescription:  sql.NullString{},
		OverviewContent:  sql.NullString{String: "Full desc", Valid: true},
		MetaDescription:  sql.NullString{String: "SEO description", Valid: true},
		ReferenceCode:    sql.NullString{},
		IsPublished:      sql.NullBool{Bool: true, Valid: true},
		DisplayOrder:     sql.NullInt64{},
	})
	if err != nil {
		t.Fatalf("create solution with meta: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/solutions/test-solution", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("solution detail route not found")
	}
}

func TestPublicSolutionDetail_Challenges(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	sol, _ := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Predictive Maintenance",
		Slug:             "predictive-maintenance",
		ShortDescription: "Predict failures",
		Icon:             "maintenance",
		HeroImageUrl:     sql.NullString{},
		HeroTitle:        sql.NullString{},
		HeroDescription:  sql.NullString{},
		OverviewContent:  sql.NullString{String: "Predictive maintenance solution", Valid: true},
		MetaDescription:  sql.NullString{},
		ReferenceCode:    sql.NullString{},
		IsPublished:      sql.NullBool{Bool: true, Valid: true},
		DisplayOrder:     sql.NullInt64{},
	})

	_, err := queries.CreateSolutionChallenge(ctx, sqlc.CreateSolutionChallengeParams{
		SolutionID:  sol.ID,
		Title:       "Equipment Downtime",
		Description: "Unexpected equipment failures",
		Icon:        "alert",
		DisplayOrder: sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		t.Fatalf("create challenge: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/solutions/%s", sol.Slug), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("solution with challenges route not found")
	}
}
