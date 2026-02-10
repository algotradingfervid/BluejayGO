package public

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type CaseStudiesHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewCaseStudiesHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *CaseStudiesHandler {
	return &CaseStudiesHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

func (h *CaseStudiesHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// GET /case-studies
func (h *CaseStudiesHandler) CaseStudiesList(c echo.Context) error {
	ctx := c.Request().Context()

	// Get optional industry filter
	industryParam := c.QueryParam("industry")
	var selectedIndustryID int64
	var caseStudies interface{}
	var totalCount int64
	var err error

	if industryParam != "" {
		selectedIndustryID, err = strconv.ParseInt(industryParam, 10, 64)
		if err != nil {
			h.logger.Error("invalid industry parameter", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid industry parameter")
		}
	}

	// Determine cache key based on filter
	cacheKey := "page:case-studies"
	if selectedIndustryID > 0 {
		cacheKey = fmt.Sprintf("page:case-studies:industry:%d", selectedIndustryID)
	}

	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Fetch case studies based on filter
	if selectedIndustryID > 0 {
		caseStudies, err = h.queries.ListCaseStudiesByIndustry(ctx, selectedIndustryID)
		if err != nil {
			h.logger.Error("failed to list case studies by industry", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		totalCount, err = h.queries.CountCaseStudiesByIndustry(ctx, selectedIndustryID)
		if err != nil {
			h.logger.Error("failed to count case studies by industry", "error", err)
			totalCount = 0
		}
	} else {
		caseStudies, err = h.queries.ListCaseStudies(ctx)
		if err != nil {
			h.logger.Error("failed to list case studies", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		totalCount, err = h.queries.CountCaseStudies(ctx)
		if err != nil {
			h.logger.Error("failed to count case studies", "error", err)
			totalCount = 0
		}
	}

	// Get industries for filter dropdown
	industries, err := h.queries.ListIndustries(ctx)
	if err != nil {
		h.logger.Error("failed to list industries", "error", err)
		industries = []sqlc.Industry{}
	}

	data := map[string]interface{}{
		"Title":              "Case Studies",
		"CaseStudies":        caseStudies,
		"Industries":         industries,
		"SelectedIndustryID": selectedIndustryID,
		"TotalCount":         totalCount,
		"CurrentPage":        "case-studies",
	}

	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/case_studies.html", data)
}

// GET /case-studies/:slug
func (h *CaseStudiesHandler) CaseStudyDetail(c echo.Context) error {
	slug := c.Param("slug")
	preview := isPreviewRequest(c)

	if !preview {
		cacheKey := fmt.Sprintf("page:case-studies:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	ctx := c.Request().Context()

	var csID int64
	var csTitle, csSlug, csOgImage string
	var csMetaTitle, csMetaDesc, csBullets sql.NullString
	var caseStudyObj interface{}

	if preview {
		cs, err := h.queries.GetCaseStudyBySlugIncludeDrafts(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Case study not found")
		}
		if err != nil {
			h.logger.Error("failed to load case study", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		csID, csTitle, csSlug, csOgImage = cs.ID, cs.Title, cs.Slug, cs.OgImage
		csMetaTitle, csMetaDesc, csBullets = cs.MetaTitle, cs.MetaDescription, cs.ChallengeBullets
		caseStudyObj = cs
	} else {
		cs, err := h.queries.GetCaseStudyBySlug(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Case study not found")
		}
		if err != nil {
			h.logger.Error("failed to load case study", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		csID, csTitle, csSlug, csOgImage = cs.ID, cs.Title, cs.Slug, cs.OgImage
		csMetaTitle, csMetaDesc, csBullets = cs.MetaTitle, cs.MetaDescription, cs.ChallengeBullets
		caseStudyObj = cs
	}

	products, err := h.queries.GetCaseStudyProducts(ctx, csID)
	if err != nil {
		h.logger.Error("failed to load case study products", "error", err)
		products = []sqlc.GetCaseStudyProductsRow{}
	}

	metrics, err := h.queries.GetCaseStudyMetrics(ctx, csID)
	if err != nil {
		h.logger.Error("failed to load case study metrics", "error", err)
		metrics = []sqlc.GetCaseStudyMetricsRow{}
	}

	// Parse ChallengeBullets JSON array
	var challengeBullets []string
	if csBullets.Valid {
		if err := json.Unmarshal([]byte(csBullets.String), &challengeBullets); err != nil {
			h.logger.Error("failed to parse challenge bullets", "error", err)
			challengeBullets = []string{}
		}
	} else {
		challengeBullets = []string{}
	}

	metaDesc := ""
	if csMetaDesc.Valid {
		metaDesc = csMetaDesc.String
	}
	metaTitle := ""
	if csMetaTitle.Valid && csMetaTitle.String != "" {
		metaTitle = csMetaTitle.String
	}

	data := map[string]interface{}{
		"Title":            csTitle,
		"MetaTitle":        metaTitle,
		"MetaDescription":  metaDesc,
		"MetaDesc":         metaDesc,
		"OGImage":          csOgImage,
		"CanonicalURL":     fmt.Sprintf("/case-studies/%s", csSlug),
		"CaseStudy":        caseStudyObj,
		"Products":         products,
		"Metrics":          metrics,
		"ChallengeBullets": challengeBullets,
		"CurrentPage":      "case-studies",
	}

	if preview {
		data["IsPreview"] = true
		data["EditURL"] = fmt.Sprintf("/admin/case-studies/%d/edit", csID)
		return h.renderAndCache(c, "preview:case-study:"+slug, 0, http.StatusOK, "public/pages/case_study_detail.html", data)
	}

	return h.renderAndCache(c, fmt.Sprintf("page:case-studies:%s", slug), 1800, http.StatusOK, "public/pages/case_study_detail.html", data)
}
