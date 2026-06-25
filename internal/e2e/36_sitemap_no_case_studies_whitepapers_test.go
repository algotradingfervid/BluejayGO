package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

// TestSitemap_OmitsCaseStudiesAndWhitepapers verifies that /sitemap.xml no longer
// emits any /case-studies or /whitepapers URLs (neither the static index entries nor
// the dynamic per-item URLs), while leaving the rest of the sitemap intact.
func TestSitemap_OmitsCaseStudiesAndWhitepapers(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	// Seed an industry + a published case study so a dynamic /case-studies/{slug}
	// URL would be emitted if the case-studies loop were still present.
	ind, err := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Manufacturing", Slug: "manufacturing",
	})
	if err != nil {
		t.Fatalf("create industry: %v", err)
	}
	_, err = queries.AdminCreateCaseStudy(ctx, sqlc.AdminCreateCaseStudyParams{
		Slug:             "acme-case-study",
		Title:            "Acme Case Study",
		ClientName:       "Acme",
		IndustryID:       ind.ID,
		Summary:          "s",
		ChallengeContent: "c",
		SolutionContent:  "x",
		OutcomeContent:   "o",
		IsPublished:      1,
	})
	if err != nil {
		t.Fatalf("create case study: %v", err)
	}

	// Seed a topic + a published whitepaper so a dynamic /whitepapers/{slug} URL
	// would be emitted if the whitepapers loop were still present.
	topic, err := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "Security", Slug: "security", ColorHex: "#000000", Icon: "shield",
	})
	if err != nil {
		t.Fatalf("create whitepaper topic: %v", err)
	}
	_, err = queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Acme Whitepaper",
		Slug:           "acme-whitepaper",
		Description:    "d",
		TopicID:        topic.ID,
		PdfFilePath:    "/files/acme.pdf",
		FileSizeBytes:  1024,
		PublishedDate:  "2026-01-01",
		IsPublished:    1,
		CoverColorFrom: "#111111",
		CoverColorTo:   "#222222",
	})
	if err != nil {
		t.Fatalf("create whitepaper: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /sitemap.xml, got %d", rec.Code)
	}

	body := rec.Body.String()

	if strings.Contains(body, "/case-studies") {
		t.Errorf("sitemap should NOT contain any /case-studies URL, but it does; body:\n%s", body)
	}
	if strings.Contains(body, "/whitepapers") {
		t.Errorf("sitemap should NOT contain any /whitepapers URL, but it does; body:\n%s", body)
	}

	// Sanity: other entries must remain so we know we didn't gut the sitemap.
	if !strings.Contains(body, "/blog") {
		t.Errorf("sitemap should still contain /blog index entry; body:\n%s", body)
	}
	if !strings.Contains(body, "/solutions") {
		t.Errorf("sitemap should still contain /solutions index entry; body:\n%s", body)
	}
}
