package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicSolutionsListing(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	_, err := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Smart Factory",
		Slug:             "smart-factory",
		Icon:             "factory",
		ShortDescription: "Automate your manufacturing",
		HeroImageUrl:     sql.NullString{},
		HeroTitle:        sql.NullString{},
		HeroDescription:  sql.NullString{},
		OverviewContent:  sql.NullString{String: "Full description", Valid: true},
		MetaDescription:  sql.NullString{},
		ReferenceCode:    sql.NullString{},
		IsPublished:      sql.NullBool{Bool: true, Valid: true},
		DisplayOrder:     sql.NullInt64{},
	})
	if err != nil {
		t.Fatalf("create solution: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/solutions", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("solutions route not found")
	}
}

func TestPublicSolutionsListing_Empty(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/solutions", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("solutions route not found with empty list")
	}
}
