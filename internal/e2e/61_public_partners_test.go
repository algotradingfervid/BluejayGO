package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPublicPartnersPage(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold tier partners",
		SortOrder:   1,
	})

	_, err := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Partner Co",
		TierID:       tier.ID,
		LogoUrl:      sql.NullString{},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{String: "https://partner.com", Valid: true},
		Description:  sql.NullString{String: "Leading technology partner", Valid: true},
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create partner: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page route not found")
	}
}

func TestPublicPartnersPage_MultipleTiers(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	platinum, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Platinum",
		Slug:        "platinum",
		Description: "Platinum partners",
		SortOrder:   1,
	})

	gold, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Gold",
		Slug:        "gold",
		Description: "Gold partners",
		SortOrder:   2,
	})

	silver, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Silver",
		Slug:        "silver",
		Description: "Silver partners",
		SortOrder:   3,
	})

	_, _ = queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Platinum Partner",
		TierID:       platinum.ID,
		LogoUrl:      sql.NullString{},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{},
		Description:  sql.NullString{String: "Top tier partner", Valid: true},
		DisplayOrder: 1,
	})

	_, _ = queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Gold Partner",
		TierID:       gold.ID,
		LogoUrl:      sql.NullString{},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{},
		Description:  sql.NullString{String: "Gold tier partner", Valid: true},
		DisplayOrder: 2,
	})

	_, _ = queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Silver Partner",
		TierID:       silver.ID,
		LogoUrl:      sql.NullString{},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{},
		Description:  sql.NullString{String: "Silver tier partner", Valid: true},
		DisplayOrder: 3,
	})

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page with multiple tiers route not found")
	}
}

func TestPublicPartnersPage_WithLogos(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Premier",
		Slug:        "premier",
		Description: "Premier partners",
		SortOrder:   1,
	})

	_, err := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Logo Partner",
		TierID:       tier.ID,
		LogoUrl:      sql.NullString{String: "/images/partner-logo.png", Valid: true},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{},
		Description:  sql.NullString{String: "Partner with logo", Valid: true},
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create partner with logo: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page with logos route not found")
	}
}

func TestPublicPartnersPage_WithTestimonials(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Elite",
		Slug:        "elite",
		Description: "Elite partners",
		SortOrder:   1,
	})

	partner, _ := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Elite Partner",
		TierID:       tier.ID,
		LogoUrl:      sql.NullString{},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{},
		Description:  sql.NullString{String: "Elite partnership", Valid: true},
		DisplayOrder: 1,
	})

	_, err := queries.CreateTestimonial(ctx, sqlc.CreateTestimonialParams{
		PartnerID:    partner.ID,
		Quote:        "Excellent partnership experience",
		AuthorName:   "John CEO",
		AuthorTitle:  "CEO",
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create testimonial: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page with testimonials route not found")
	}
}

func TestPublicPartnersPage_WithWebsiteLinks(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Standard",
		Slug:        "standard",
		Description: "Standard partners",
		SortOrder:   1,
	})

	_, err := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Web Partner",
		TierID:       tier.ID,
		LogoUrl:      sql.NullString{},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{String: "https://webpartner.com", Valid: true},
		Description:  sql.NullString{String: "Partner with website", Valid: true},
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create partner with website: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page with website links route not found")
	}
}

func TestPublicPartnersPage_EmptyPartners(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page should handle empty partners gracefully")
	}
}

func TestPublicPartnersPage_ManyPartners(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Bronze",
		Slug:        "bronze",
		Description: "Bronze partners",
		SortOrder:   1,
	})

	for i := 1; i <= 15; i++ {
		_, _ = queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
			Name:         "Partner " + string(rune(i+'0')),
			TierID:       tier.ID,
			LogoUrl:      sql.NullString{},
			Icon:         sql.NullString{},
			WebsiteUrl:   sql.NullString{},
			Description:  sql.NullString{String: "Partner description", Valid: true},
			DisplayOrder: int64(i),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page with many partners route not found")
	}
}

func TestPublicPartnersPage_CompleteData(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	tier, _ := queries.CreatePartnerTier(ctx, sqlc.CreatePartnerTierParams{
		Name:        "Premium",
		Slug:        "premium",
		Description: "Premium partners get exclusive benefits",
		SortOrder:   1,
	})

	partner, _ := queries.CreatePartner(ctx, sqlc.CreatePartnerParams{
		Name:         "Complete Partner",
		TierID:       tier.ID,
		LogoUrl:      sql.NullString{String: "/images/complete-logo.png", Valid: true},
		Icon:         sql.NullString{},
		WebsiteUrl:   sql.NullString{String: "https://complete.com", Valid: true},
		Description:  sql.NullString{String: "Full partner profile", Valid: true},
		DisplayOrder: 1,
	})

	_, _ = queries.CreateTestimonial(ctx, sqlc.CreateTestimonialParams{
		PartnerID:    partner.ID,
		Quote:        "Outstanding collaboration",
		AuthorName:   "Jane Smith",
		AuthorTitle:  "CTO",
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("partners page with complete data route not found")
	}
}
