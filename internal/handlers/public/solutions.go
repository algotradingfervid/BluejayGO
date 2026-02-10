package public

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type SolutionsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewSolutionsHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *SolutionsHandler {
	return &SolutionsHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

func (h *SolutionsHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// GET /solutions
func (h *SolutionsHandler) SolutionsList(c echo.Context) error {
	cacheKey := "page:solutions"
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	solutions, err := h.queries.ListPublishedSolutions(ctx)
	if err != nil {
		h.logger.Error("failed to list published solutions", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	features, err := h.queries.ListSolutionPageFeatures(ctx)
	if err != nil {
		h.logger.Error("failed to list solution page features", "error", err)
		features = []sqlc.SolutionPageFeature{}
	}

	cta, err := h.queries.GetActiveSolutionsListingCTA(ctx)
	if err != nil && err != sql.ErrNoRows {
		h.logger.Error("failed to get active solutions listing CTA", "error", err)
	}

	heroSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "solutions", SectionKey: "hero"})
	gridSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "solutions", SectionKey: "grid_section"})
	featuresSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "solutions", SectionKey: "features_section"})

	data := map[string]interface{}{
		"Title":           "Solutions",
		"Solutions":       solutions,
		"Features":        features,
		"CTA":             cta,
		"CurrentPage":     "solutions",
		"PageHero":        heroSection,
		"GridSection":     gridSection,
		"FeaturesSection": featuresSection,
	}

	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/solutions_list.html", data)
}

// GET /solutions/:slug
func (h *SolutionsHandler) SolutionDetail(c echo.Context) error {
	slug := c.Param("slug")
	preview := isPreviewRequest(c)

	if !preview {
		cacheKey := fmt.Sprintf("page:solutions:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	ctx := c.Request().Context()

	var solution sqlc.Solution
	var err error
	if preview {
		solution, err = h.queries.GetSolutionBySlugIncludeDrafts(ctx, slug)
	} else {
		solution, err = h.queries.GetSolutionBySlug(ctx, slug)
	}
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Solution not found")
	}
	if err != nil {
		h.logger.Error("failed to load solution", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	stats, err := h.queries.GetSolutionStats(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution stats", "error", err)
		stats = []sqlc.SolutionStat{}
	}

	challenges, err := h.queries.GetSolutionChallenges(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution challenges", "error", err)
		challenges = []sqlc.SolutionChallenge{}
	}

	products, err := h.queries.GetSolutionProducts(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution products", "error", err)
		products = []sqlc.GetSolutionProductsRow{}
	}

	ctas, err := h.queries.GetSolutionCTAs(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution CTAs", "error", err)
		ctas = []sqlc.SolutionCta{}
	}

	allSolutions, err := h.queries.ListPublishedSolutions(ctx)
	if err != nil {
		h.logger.Error("failed to list published solutions", "error", err)
		allSolutions = []sqlc.ListPublishedSolutionsRow{}
	}

	var otherSolutions []sqlc.ListPublishedSolutionsRow
	for _, s := range allSolutions {
		if s.ID != solution.ID {
			otherSolutions = append(otherSolutions, s)
		}
	}

	metaDesc := ""
	if solution.MetaDescription.Valid {
		metaDesc = solution.MetaDescription.String
	}

	sections, _ := h.queries.ListPageSections(ctx, "solution_detail")
	replacer := strings.NewReplacer("{solution_title}", solution.Title)
	sectionMap := make(map[string]sqlc.PageSection)
	for _, s := range sections {
		s.Heading = replacer.Replace(s.Heading)
		s.Description = replacer.Replace(s.Description)
		sectionMap[s.SectionKey] = s
	}

	data := map[string]interface{}{
		"Title":           solution.Title,
		"MetaTitle":       solution.MetaTitle,
		"MetaDescription": metaDesc,
		"MetaDesc":        metaDesc,
		"OGImage":         solution.OgImage,
		"CanonicalURL":    fmt.Sprintf("/solutions/%s", solution.Slug),
		"Solution":        solution,
		"Stats":           stats,
		"Challenges":      challenges,
		"Products":        products,
		"CTAs":            ctas,
		"OtherSolutions":  otherSolutions,
		"CurrentPage":     "solutions",
		"Sections":        sectionMap,
	}

	if preview {
		data["IsPreview"] = true
		data["EditURL"] = fmt.Sprintf("/admin/solutions/%d/edit", solution.ID)
		return h.renderAndCache(c, "preview:solution:"+slug, 0, http.StatusOK, "public/pages/solution_detail.html", data)
	}

	return h.renderAndCache(c, fmt.Sprintf("page:solutions:%s", slug), 1800, http.StatusOK, "public/pages/solution_detail.html", data)
}
