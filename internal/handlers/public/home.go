package public

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type HomeHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewHomeHandler(queries *sqlc.Queries, logger *slog.Logger) *HomeHandler {
	return &HomeHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *HomeHandler) ShowHomePage(c echo.Context) error {
	ctx := c.Request().Context()

	settings, err := h.queries.GetSettings(ctx)
	if err != nil {
		h.logger.Error("failed to get settings", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Homepage-specific tables
	hero, _ := h.queries.GetActiveHero(ctx)
	stats, _ := h.queries.ListActiveStats(ctx)
	testimonials, _ := h.queries.ListActiveTestimonialsHomepage(ctx)
	cta, _ := h.queries.GetActiveCTA(ctx)

	// Existing content tables
	featuredProducts, _ := h.queries.ListFeaturedProducts(ctx, 6)
	solutions, _ := h.queries.ListPublishedSolutions(ctx)
	featuredPartners, _ := h.queries.ListFeaturedPartners(ctx, 12)
	latestPosts, _ := h.queries.ListLatestPublishedPosts(ctx, 3)

	// Page sections for editable labels/headings
	sections, _ := h.queries.ListPageSections(ctx, "home")
	sectionMap := make(map[string]sqlc.PageSection)
	for _, s := range sections {
		sectionMap[s.SectionKey] = s
	}

	data := map[string]interface{}{
		"Title":            settings.SiteName,
		"Settings":         settings,
		"Hero":             hero,
		"Stats":            stats,
		"Testimonials":     testimonials,
		"CTA":              cta,
		"FeaturedProducts": featuredProducts,
		"Solutions":        solutions,
		"FeaturedPartners": featuredPartners,
		"LatestPosts":      latestPosts,
		"Sections":         sectionMap,
	}

	if cats := c.Get("footer_categories"); cats != nil {
		data["FooterCategories"] = cats
	}
	if sols := c.Get("footer_solutions"); sols != nil {
		data["FooterSolutions"] = sols
	}
	if res := c.Get("footer_resources"); res != nil {
		data["FooterResources"] = res
	}

	return c.Render(http.StatusOK, "public/pages/home.html", data)
}
