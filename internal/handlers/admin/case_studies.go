package admin

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

const caseStudiesPerPage = 15

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

// List displays all case studies with filtering and pagination
func (h *CaseStudiesHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	search := c.QueryParam("search")
	status := c.QueryParam("status")
	pageStr := c.QueryParam("page")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * caseStudiesPerPage)

	caseStudies, err := h.queries.AdminListCaseStudiesFiltered(ctx, sqlc.AdminListCaseStudiesFilteredParams{
		FilterSearch: search,
		FilterStatus: status,
		PageLimit:    caseStudiesPerPage,
		PageOffset:   offset,
	})
	if err != nil {
		h.logger.Error("Failed to list case studies", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load case studies")
	}

	total, err := h.queries.CountCaseStudiesAdminFiltered(ctx, sqlc.CountCaseStudiesAdminFilteredParams{
		FilterSearch: search,
		FilterStatus: status,
	})
	if err != nil {
		h.logger.Error("Failed to count case studies", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to count case studies")
	}

	totalPages := int(math.Ceil(float64(total) / float64(caseStudiesPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	showFrom := offset + 1
	showTo := offset + int64(len(caseStudies))
	if total == 0 {
		showFrom = 0
	}

	hasFilters := search != "" || status != ""

	return c.Render(http.StatusOK, "admin/pages/case_studies_list.html", map[string]interface{}{
		"Title":       "Case Studies",
		"CaseStudies": caseStudies,
		"Search":      search,
		"Status":      status,
		"HasFilters":  hasFilters,
		"Page":        page,
		"TotalPages":  totalPages,
		"Pages":       pages,
		"Total":       total,
		"ShowFrom":    showFrom,
		"ShowTo":      showTo,
	})
}

// New displays the form for creating a new case study
func (h *CaseStudiesHandler) New(c echo.Context) error {
	industries, err := h.queries.ListIndustries(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list industries", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load industries")
	}

	return c.Render(http.StatusOK, "admin/pages/case_studies_form.html", map[string]interface{}{
		"Title":      "New Case Study",
		"FormAction": "/admin/case-studies",
		"Item":       nil,
		"Industries": industries,
		"IsNew":      true,
	})
}

// Create handles case study creation
func (h *CaseStudiesHandler) Create(c echo.Context) error {
	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	clientName := c.FormValue("client_name")
	summary := c.FormValue("summary")
	heroImageUrl := c.FormValue("hero_image_url")
	challengeTitle := c.FormValue("challenge_title")
	challengeContent := c.FormValue("challenge_content")
	solutionTitle := c.FormValue("solution_title")
	solutionContent := c.FormValue("solution_content")
	outcomeTitle := c.FormValue("outcome_title")
	outcomeContent := c.FormValue("outcome_content")
	metaTitle := c.FormValue("meta_title")
	metaDescription := c.FormValue("meta_description")

	industryID := int64(0)
	if v := c.FormValue("industry_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			industryID = parsed
		}
	}

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	isPublished := int64(0)
	if v := c.FormValue("is_published"); v == "1" || v == "on" {
		isPublished = 1
	}

	// Convert challenge_bullets comma-separated to JSON array
	challengeBulletsRaw := c.FormValue("challenge_bullets")
	challengeBulletsJSON := sql.NullString{}
	if challengeBulletsRaw != "" {
		bullets := strings.Split(challengeBulletsRaw, ",")
		trimmedBullets := make([]string, 0, len(bullets))
		for _, bullet := range bullets {
			trimmed := strings.TrimSpace(bullet)
			if trimmed != "" {
				trimmedBullets = append(trimmedBullets, trimmed)
			}
		}
		if len(trimmedBullets) > 0 {
			jsonBytes, err := json.Marshal(trimmedBullets)
			if err == nil {
				challengeBulletsJSON = sql.NullString{String: string(jsonBytes), Valid: true}
			}
		}
	}

	params := sqlc.AdminCreateCaseStudyParams{
		Slug:              slug,
		Title:             title,
		ClientName:        clientName,
		IndustryID:        industryID,
		HeroImageUrl:      sql.NullString{String: heroImageUrl, Valid: heroImageUrl != ""},
		ChallengeBullets:  challengeBulletsJSON,
		MetaTitle:         sql.NullString{String: metaTitle, Valid: metaTitle != ""},
		MetaDescription:   sql.NullString{String: metaDescription, Valid: metaDescription != ""},
		Summary:           summary,
		ChallengeTitle:    challengeTitle,
		ChallengeContent:  challengeContent,
		SolutionTitle:     solutionTitle,
		SolutionContent:   solutionContent,
		OutcomeTitle:      outcomeTitle,
		OutcomeContent:    outcomeContent,
		IsPublished:       isPublished,
		DisplayOrder:      displayOrder,
	}

	_, err := h.queries.AdminCreateCaseStudy(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create case study", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to create case study")
	}

	h.cache.DeleteByPrefix("page:case-studies")
	logActivity(c, "created", "case_study", 0, c.FormValue("title"), "Created Case Study '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/case-studies")
}

// Edit displays the form for editing a case study
func (h *CaseStudiesHandler) Edit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid case study ID")
	}

	caseStudy, err := h.queries.AdminGetCaseStudy(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get case study", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load case study")
	}

	industries, err := h.queries.ListIndustries(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list industries", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load industries")
	}

	caseStudyProducts, err := h.queries.AdminListCaseStudyProducts(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get case study products", "error", err)
		caseStudyProducts = []sqlc.AdminListCaseStudyProductsRow{}
	}

	metrics, err := h.queries.AdminListMetrics(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get case study metrics", "error", err)
		metrics = []sqlc.AdminListMetricsRow{}
	}

	allProducts, err := h.queries.ListProducts(c.Request().Context(), sqlc.ListProductsParams{Limit: 1000, Offset: 0})
	if err != nil {
		h.logger.Error("Failed to list all products", "error", err)
		allProducts = []sqlc.Product{}
	}

	return c.Render(http.StatusOK, "admin/pages/case_studies_form.html", map[string]interface{}{
		"Title":       "Edit Case Study",
		"FormAction":  "/admin/case-studies/" + c.Param("id"),
		"Item":        caseStudy,
		"Industries":  industries,
		"Products":    caseStudyProducts,
		"Metrics":     metrics,
		"AllProducts": allProducts,
		"IsNew":       false,
	})
}

// Update handles case study updates
func (h *CaseStudiesHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid case study ID")
	}

	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	clientName := c.FormValue("client_name")
	summary := c.FormValue("summary")
	heroImageUrl := c.FormValue("hero_image_url")
	challengeTitle := c.FormValue("challenge_title")
	challengeContent := c.FormValue("challenge_content")
	solutionTitle := c.FormValue("solution_title")
	solutionContent := c.FormValue("solution_content")
	outcomeTitle := c.FormValue("outcome_title")
	outcomeContent := c.FormValue("outcome_content")
	metaTitle := c.FormValue("meta_title")
	metaDescription := c.FormValue("meta_description")

	industryID := int64(0)
	if v := c.FormValue("industry_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			industryID = parsed
		}
	}

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	isPublished := int64(0)
	if c.FormValue("is_published") == "on" {
		isPublished = 1
	}

	// Convert challenge_bullets comma-separated to JSON array
	challengeBulletsRaw := c.FormValue("challenge_bullets")
	challengeBulletsJSON := sql.NullString{}
	if challengeBulletsRaw != "" {
		bullets := strings.Split(challengeBulletsRaw, ",")
		trimmedBullets := make([]string, 0, len(bullets))
		for _, bullet := range bullets {
			trimmed := strings.TrimSpace(bullet)
			if trimmed != "" {
				trimmedBullets = append(trimmedBullets, trimmed)
			}
		}
		if len(trimmedBullets) > 0 {
			jsonBytes, err := json.Marshal(trimmedBullets)
			if err == nil {
				challengeBulletsJSON = sql.NullString{String: string(jsonBytes), Valid: true}
			}
		}
	}

	params := sqlc.AdminUpdateCaseStudyParams{
		Slug:             slug,
		Title:            title,
		ClientName:       clientName,
		IndustryID:       industryID,
		HeroImageUrl:     sql.NullString{String: heroImageUrl, Valid: heroImageUrl != ""},
		ChallengeBullets: challengeBulletsJSON,
		MetaTitle:        sql.NullString{String: metaTitle, Valid: metaTitle != ""},
		MetaDescription:  sql.NullString{String: metaDescription, Valid: metaDescription != ""},
		Summary:          summary,
		ChallengeTitle:   challengeTitle,
		ChallengeContent: challengeContent,
		SolutionTitle:    solutionTitle,
		SolutionContent:  solutionContent,
		OutcomeTitle:     outcomeTitle,
		OutcomeContent:   outcomeContent,
		IsPublished:      isPublished,
		DisplayOrder:     displayOrder,
		ID:               id,
	}

	_, err = h.queries.AdminUpdateCaseStudy(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to update case study", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update case study")
	}

	h.cache.DeleteByPrefix("page:case-studies")
	logActivity(c, "updated", "case_study", id, c.FormValue("title"), "Updated Case Study '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/case-studies")
}

// Delete handles case study deletion
func (h *CaseStudiesHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid case study ID")
	}

	err = h.queries.AdminDeleteCaseStudy(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete case study", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete case study")
	}

	h.cache.DeleteByPrefix("page:case-studies")
	logActivity(c, "deleted", "case_study", id, "", "Deleted Case Study #%d", id)
	return c.NoContent(http.StatusNoContent)
}

// AddProduct links a product to a case study
func (h *CaseStudiesHandler) AddProduct(c echo.Context) error {
	caseStudyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid case study ID")
	}

	productID, err := strconv.ParseInt(c.FormValue("product_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid product ID")
	}

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	params := sqlc.AdminAddCaseStudyProductParams{
		CaseStudyID:  caseStudyID,
		ProductID:    productID,
		DisplayOrder: displayOrder,
	}

	_, err = h.queries.AdminAddCaseStudyProduct(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to add product to case study", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add product")
	}

	h.cache.DeleteByPrefix("page:case-studies")

	caseStudyProducts, err := h.queries.AdminListCaseStudyProducts(c.Request().Context(), caseStudyID)
	if err != nil {
		h.logger.Error("Failed to get case study products", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load products")
	}

	logActivity(c, "updated", "case_study", caseStudyID, "", "Updated Case Study #%d sub-resources", caseStudyID)
	return c.Render(http.StatusOK, "admin/partials/case_study_products.html", map[string]interface{}{
		"Products": caseStudyProducts,
	})
}

// RemoveProduct removes a product from a case study
func (h *CaseStudiesHandler) RemoveProduct(c echo.Context) error {
	caseStudyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid case study ID")
	}

	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid product ID")
	}

	params := sqlc.AdminRemoveCaseStudyProductParams{
		CaseStudyID: caseStudyID,
		ProductID:   productID,
	}

	err = h.queries.AdminRemoveCaseStudyProduct(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to remove product from case study", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to remove product")
	}

	h.cache.DeleteByPrefix("page:case-studies")
	logActivity(c, "updated", "case_study", caseStudyID, "", "Updated Case Study #%d sub-resources", caseStudyID)
	return c.NoContent(http.StatusNoContent)
}

// AddMetric adds a metric to a case study
func (h *CaseStudiesHandler) AddMetric(c echo.Context) error {
	caseStudyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid case study ID")
	}

	metricValue := c.FormValue("metric_value")
	metricLabel := c.FormValue("metric_label")

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	params := sqlc.AdminCreateMetricParams{
		CaseStudyID:  caseStudyID,
		MetricValue:  metricValue,
		MetricLabel:  metricLabel,
		DisplayOrder: displayOrder,
	}

	_, err = h.queries.AdminCreateMetric(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create metric", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add metric")
	}

	h.cache.DeleteByPrefix("page:case-studies")

	metrics, err := h.queries.AdminListMetrics(c.Request().Context(), caseStudyID)
	if err != nil {
		h.logger.Error("Failed to get case study metrics", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load metrics")
	}

	logActivity(c, "updated", "case_study", caseStudyID, "", "Updated Case Study #%d sub-resources", caseStudyID)
	return c.Render(http.StatusOK, "admin/partials/case_study_metrics.html", map[string]interface{}{
		"Metrics": metrics,
	})
}

// DeleteMetric deletes a metric from a case study
func (h *CaseStudiesHandler) DeleteMetric(c echo.Context) error {
	metricID, err := strconv.ParseInt(c.Param("metricId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid metric ID")
	}

	err = h.queries.AdminDeleteMetric(c.Request().Context(), metricID)
	if err != nil {
		h.logger.Error("Failed to delete metric", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete metric")
	}

	h.cache.DeleteByPrefix("page:case-studies")
	// Note: caseStudyID not available in DeleteMetric, use metricID as reference
	logActivity(c, "updated", "case_study", metricID, "", "Updated Case Study #%d sub-resources", metricID)
	return c.NoContent(http.StatusNoContent)
}
