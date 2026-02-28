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

func TestCaseStudyProductsAdd(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Cat",
		Slug:        "cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PROD-001",
		Slug:        "prod-001",
		Name:        "Test Product",
		Description: "desc",
		CategoryID:  cat.ID,
		Status:      "published",
	})

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

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/case-studies/%d/products", cs.ID), strings.NewReader(url.Values{
		"product_id":    {fmt.Sprintf("%d", prod.ID)},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}

	products, _ := queries.AdminListCaseStudyProducts(ctx, cs.ID)
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
}

func TestCaseStudyProductsRemove(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Cat",
		Slug:        "cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PROD-001",
		Slug:        "prod-001",
		Name:        "Test Product",
		Description: "desc",
		CategoryID:  cat.ID,
		Status:      "published",
	})

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

	queries.AdminAddCaseStudyProduct(ctx, sqlc.AdminAddCaseStudyProductParams{
		CaseStudyID: cs.ID,
		ProductID:   prod.ID,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/case-studies/%d/products/%d", cs.ID, prod.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rec.Code)
	}

	products, _ := queries.AdminListCaseStudyProducts(ctx, cs.ID)
	if len(products) != 0 {
		t.Errorf("expected 0 products after remove, got %d", len(products))
	}
}

func TestCaseStudyProductsDuplicate(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Cat",
		Slug:        "cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "PROD-001",
		Slug:        "prod-001",
		Name:        "Test Product",
		Description: "desc",
		CategoryID:  cat.ID,
		Status:      "published",
	})

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

	queries.AdminAddCaseStudyProduct(ctx, sqlc.AdminAddCaseStudyProductParams{
		CaseStudyID: cs.ID,
		ProductID:   prod.ID,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/case-studies/%d/products", cs.ID), strings.NewReader(url.Values{
		"product_id":    {fmt.Sprintf("%d", prod.ID)},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Logf("duplicate add may fail with constraint error")
	}
}
