package public

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type PartnersHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewPartnersHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *PartnersHandler {
	return &PartnersHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

func (h *PartnersHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
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
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		h.logger.Error("template render failed", "template", templateName, "error", err)
		return err
	}
	html := buf.String()
	h.cache.Set(cacheKey, html, ttlSeconds)
	return c.HTML(statusCode, html)
}

func (h *PartnersHandler) PartnersPage(c echo.Context) error {
	cacheKey := "page:partners"
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	allPartners, err := h.queries.ListPartnersByTier(ctx)
	if err != nil {
		h.logger.Error("failed to load partners", "error", err)
		allPartners = []sqlc.ListPartnersByTierRow{}
	}

	tiers, err := h.queries.ListPartnerTiers(ctx)
	if err != nil {
		h.logger.Error("failed to load partner tiers", "error", err)
		tiers = []sqlc.PartnerTier{}
	}

	testimonials, err := h.queries.ListActiveTestimonials(ctx)
	if err != nil {
		h.logger.Error("failed to load testimonials", "error", err)
		testimonials = []sqlc.ListActiveTestimonialsRow{}
	}

	// Group partners by tier name
	partnersByTier := make(map[string][]sqlc.ListPartnersByTierRow)
	for _, p := range allPartners {
		partnersByTier[p.TierName] = append(partnersByTier[p.TierName], p)
	}

	data := map[string]interface{}{
		"Title":          "Partners",
		"CurrentPage":    "partners",
		"PartnersByTier": partnersByTier,
		"Tiers":          tiers,
		"Testimonials":   testimonials,
	}

	return h.renderAndCache(c, cacheKey, 300, http.StatusOK, "public/pages/partners.html", data)
}
