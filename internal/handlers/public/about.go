package public

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type AboutHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewAboutHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *AboutHandler {
	return &AboutHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

func (h *AboutHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

func (h *AboutHandler) AboutPage(c echo.Context) error {
	cacheKey := "page:about"
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	overview, err := h.queries.GetCompanyOverview(ctx)
	if err != nil {
		h.logger.Debug("no company overview found", "error", err)
	}

	mvv, err := h.queries.GetMissionVisionValues(ctx)
	if err != nil {
		h.logger.Debug("no mission/vision/values found", "error", err)
	}

	coreValues, err := h.queries.ListCoreValues(ctx)
	if err != nil {
		h.logger.Error("failed to load core values", "error", err)
		coreValues = []sqlc.CoreValue{}
	}

	milestones, err := h.queries.ListMilestones(ctx)
	if err != nil {
		h.logger.Error("failed to load milestones", "error", err)
		milestones = []sqlc.Milestone{}
	}

	certs, err := h.queries.ListCertifications(ctx)
	if err != nil {
		h.logger.Error("failed to load certifications", "error", err)
		certs = []sqlc.Certification{}
	}

	data := map[string]interface{}{
		"Title":       "About Us",
		"CurrentPage": "about",
		"Overview":    nil,
		"MVV":         nil,
		"CoreValues":  coreValues,
		"Milestones":  milestones,
		"Certs":       certs,
	}

	if overview.ID > 0 {
		data["Overview"] = overview
	}
	if mvv.ID > 0 {
		data["MVV"] = mvv
	}

	return h.renderAndCache(c, cacheKey, 300, http.StatusOK, "public/pages/about.html", data)
}
