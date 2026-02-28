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

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPartnerTestimonialsCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold tier partners",
		SortOrder:   1,
	})

	partner, _ := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Acme Corp",
		TierID:       tier.ID,
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/partners/testimonials", strings.NewReader(url.Values{
		"partner_id":    {fmt.Sprintf("%d", partner.ID)},
		"quote":         {"This partnership has been transformative for our business."},
		"author_name":   {"John Doe"},
		"author_title":  {"CTO"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	testimonials, _ := queries.ListActiveTestimonials(ctx)
	if len(testimonials) != 1 {
		t.Fatalf("expected 1 testimonial, got %d", len(testimonials))
	}
	if testimonials[0].AuthorName != "John Doe" {
		t.Errorf("expected 'John Doe', got %q", testimonials[0].AuthorName)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/partners/testimonials/%d", testimonials[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	_, err := queries.GetTestimonial(ctx, testimonials[0].ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestPartnerTestimonialsList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold tier partners",
		SortOrder:   1,
	})

	partner, _ := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Acme Corp",
		TierID:       tier.ID,
		DisplayOrder: 1,
	})

	for i := 1; i <= 3; i++ {
		queries.CreateTestimonial(ctx, sqlc.CreateTestimonialParams{
			PartnerID:    partner.ID,
			Quote:        fmt.Sprintf("Quote %d", i),
			AuthorName:   fmt.Sprintf("Author %d", i),
			AuthorTitle:  "Manager",
			DisplayOrder: int64(i),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/partners/testimonials", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partner testimonials list route not found")
	}
}

func TestPartnerTestimonialEdit_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold tier partners",
		SortOrder:   1,
	})

	partner, _ := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Acme Corp",
		TierID:       tier.ID,
		DisplayOrder: 1,
	})

	testimonial, _ := queries.CreateTestimonial(ctx, sqlc.CreateTestimonialParams{
		PartnerID:    partner.ID,
		Quote:        "Original quote",
		AuthorName:   "John Doe",
		AuthorTitle:  "CTO",
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/partners/testimonials/%d", testimonial.ID), strings.NewReader(url.Values{
		"partner_id":    {fmt.Sprintf("%d", partner.ID)},
		"quote":         {"Updated quote"},
		"author_name":   {"Jane Smith"},
		"author_title":  {"CEO"},
		"display_order": {"2"},
		"is_active":     {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetTestimonial(ctx, testimonial.ID)
	if updated.AuthorName != "Jane Smith" {
		t.Errorf("expected 'Jane Smith', got %q", updated.AuthorName)
	}
}
