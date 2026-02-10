package admin

import (
	"database/sql"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

const solutionsPerPage = 15

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

// List displays all solutions with filtering and pagination
func (h *SolutionsHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	search := c.QueryParam("search")
	status := c.QueryParam("status")
	pageStr := c.QueryParam("page")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * solutionsPerPage)

	solutions, err := h.queries.ListSolutionsAdminFiltered(ctx, sqlc.ListSolutionsAdminFilteredParams{
		FilterStatus: status,
		FilterSearch: search,
		PageLimit:    solutionsPerPage,
		PageOffset:   offset,
	})
	if err != nil {
		h.logger.Error("Failed to list solutions", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load solutions")
	}

	total, err := h.queries.CountSolutionsAdminFiltered(ctx, sqlc.CountSolutionsAdminFilteredParams{
		FilterStatus: status,
		FilterSearch: search,
	})
	if err != nil {
		h.logger.Error("Failed to count solutions", "error", err)
		total = 0
	}

	totalPages := int(math.Ceil(float64(total) / float64(solutionsPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	showFrom := offset + 1
	showTo := offset + int64(len(solutions))
	if total == 0 {
		showFrom = 0
	}

	return c.Render(http.StatusOK, "admin/pages/solutions_list.html", map[string]interface{}{
		"Title":      "Solutions",
		"Solutions":  solutions,
		"Search":     search,
		"Status":     status,
		"HasFilters": search != "" || status != "",
		"Page":       page,
		"TotalPages": totalPages,
		"Pages":      pages,
		"Total":      total,
		"ShowFrom":   showFrom,
		"ShowTo":     showTo,
	})
}

// New displays the form for creating a new solution
func (h *SolutionsHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/solutions_form.html", map[string]interface{}{
		"Title":      "New Solution",
		"FormAction": "/admin/solutions",
		"Item":       nil,
	})
}

// Create handles solution creation
func (h *SolutionsHandler) Create(c echo.Context) error {
	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	icon := c.FormValue("icon")
	shortDescription := c.FormValue("short_description")
	heroImageUrl := c.FormValue("hero_image_url")
	heroTitle := c.FormValue("hero_title")
	heroDescription := c.FormValue("hero_description")
	overviewContent := c.FormValue("overview_content")
	metaDescription := c.FormValue("meta_description")
	referenceCode := c.FormValue("reference_code")
	isPublished := c.FormValue("is_published") == "1" || c.FormValue("is_published") == "on"

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	params := sqlc.CreateSolutionParams{
		Title:            title,
		Slug:             slug,
		Icon:             icon,
		ShortDescription: shortDescription,
		HeroImageUrl:     sql.NullString{String: heroImageUrl, Valid: heroImageUrl != ""},
		HeroTitle:        sql.NullString{String: heroTitle, Valid: heroTitle != ""},
		HeroDescription:  sql.NullString{String: heroDescription, Valid: heroDescription != ""},
		OverviewContent:  sql.NullString{String: overviewContent, Valid: overviewContent != ""},
		MetaDescription:  sql.NullString{String: metaDescription, Valid: metaDescription != ""},
		ReferenceCode:    sql.NullString{String: referenceCode, Valid: referenceCode != ""},
		IsPublished:      sql.NullBool{Bool: isPublished, Valid: true},
		DisplayOrder:     sql.NullInt64{Int64: displayOrder, Valid: true},
	}

	_, err := h.queries.CreateSolution(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create solution", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to create solution")
	}

	h.cache.DeleteByPrefix("page:solutions")
	logActivity(c, "created", "solution", 0, c.FormValue("title"), "Created Solution '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/solutions")
}

// Edit displays the form for editing a solution
func (h *SolutionsHandler) Edit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
	}

	solution, err := h.queries.GetSolutionByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get solution", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load solution")
	}

	stats, err := h.queries.GetSolutionStats(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get solution stats", "error", err)
		stats = []sqlc.SolutionStat{}
	}

	challenges, err := h.queries.GetSolutionChallenges(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get solution challenges", "error", err)
		challenges = []sqlc.SolutionChallenge{}
	}

	products, err := h.queries.GetSolutionProducts(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get solution products", "error", err)
		products = []sqlc.GetSolutionProductsRow{}
	}

	ctas, err := h.queries.GetSolutionCTAs(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get solution CTAs", "error", err)
		ctas = []sqlc.SolutionCta{}
	}

	return c.Render(http.StatusOK, "admin/pages/solutions_form.html", map[string]interface{}{
		"Title":      "Edit Solution",
		"FormAction": "/admin/solutions/" + c.Param("id"),
		"Item":       solution,
		"Stats":      stats,
		"Challenges": challenges,
		"Products":   products,
		"CTAs":       ctas,
	})
}

// Update handles solution updates
func (h *SolutionsHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
	}

	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	icon := c.FormValue("icon")
	shortDescription := c.FormValue("short_description")
	heroImageUrl := c.FormValue("hero_image_url")
	heroTitle := c.FormValue("hero_title")
	heroDescription := c.FormValue("hero_description")
	overviewContent := c.FormValue("overview_content")
	metaDescription := c.FormValue("meta_description")
	referenceCode := c.FormValue("reference_code")
	isPublished := c.FormValue("is_published") == "1" || c.FormValue("is_published") == "on"

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	params := sqlc.UpdateSolutionParams{
		ID:               id,
		Title:            title,
		Slug:             slug,
		Icon:             icon,
		ShortDescription: shortDescription,
		HeroImageUrl:     sql.NullString{String: heroImageUrl, Valid: heroImageUrl != ""},
		HeroTitle:        sql.NullString{String: heroTitle, Valid: heroTitle != ""},
		HeroDescription:  sql.NullString{String: heroDescription, Valid: heroDescription != ""},
		OverviewContent:  sql.NullString{String: overviewContent, Valid: overviewContent != ""},
		MetaDescription:  sql.NullString{String: metaDescription, Valid: metaDescription != ""},
		ReferenceCode:    sql.NullString{String: referenceCode, Valid: referenceCode != ""},
		IsPublished:      sql.NullBool{Bool: isPublished, Valid: true},
		DisplayOrder:     sql.NullInt64{Int64: displayOrder, Valid: true},
	}

	err = h.queries.UpdateSolution(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to update solution", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update solution")
	}

	h.cache.DeleteByPrefix("page:solutions")
	logActivity(c, "updated", "solution", id, c.FormValue("title"), "Updated Solution '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/solutions")
}

// Delete handles solution deletion
func (h *SolutionsHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
	}

	err = h.queries.DeleteSolution(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete solution", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete solution")
	}

	h.cache.DeleteByPrefix("page:solutions")
	logActivity(c, "deleted", "solution", id, "", "Deleted Solution #%d", id)
	return c.NoContent(http.StatusNoContent)
}

// AddStat adds a stat to a solution
func (h *SolutionsHandler) AddStat(c echo.Context) error {
	solutionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
	}

	value := c.FormValue("value")
	label := c.FormValue("label")

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	params := sqlc.CreateSolutionStatParams{
		SolutionID:   solutionID,
		Value:        value,
		Label:        label,
		DisplayOrder: sql.NullInt64{Int64: displayOrder, Valid: true},
	}

	_, err = h.queries.CreateSolutionStat(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create solution stat", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add stat")
	}

	h.cache.DeleteByPrefix("page:solutions")

	stats, err := h.queries.GetSolutionStats(c.Request().Context(), solutionID)
	if err != nil {
		h.logger.Error("Failed to get solution stats", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load stats")
	}

	logActivity(c, "updated", "solution", solutionID, "", "Updated Solution #%d sub-resources", solutionID)
	return c.Render(http.StatusOK, "admin/partials/solution_stats.html", map[string]interface{}{
		"Stats": stats,
	})
}

// DeleteStat deletes a stat from a solution
func (h *SolutionsHandler) DeleteStat(c echo.Context) error {
	statID, err := strconv.ParseInt(c.Param("statId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid stat ID")
	}

	err = h.queries.DeleteSolutionStat(c.Request().Context(), statID)
	if err != nil {
		h.logger.Error("Failed to delete solution stat", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete stat")
	}

	h.cache.DeleteByPrefix("page:solutions")
	// Note: solutionID not available in DeleteStat, use statID as reference
	logActivity(c, "updated", "solution", statID, "", "Updated Solution #%d sub-resources", statID)
	return c.NoContent(http.StatusOK)
}

// AddChallenge adds a challenge to a solution
func (h *SolutionsHandler) AddChallenge(c echo.Context) error {
	solutionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	icon := c.FormValue("icon")

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	params := sqlc.CreateSolutionChallengeParams{
		SolutionID:   solutionID,
		Title:        title,
		Description:  description,
		Icon:         icon,
		DisplayOrder: sql.NullInt64{Int64: displayOrder, Valid: true},
	}

	_, err = h.queries.CreateSolutionChallenge(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create solution challenge", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add challenge")
	}

	h.cache.DeleteByPrefix("page:solutions")

	challenges, err := h.queries.GetSolutionChallenges(c.Request().Context(), solutionID)
	if err != nil {
		h.logger.Error("Failed to get solution challenges", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load challenges")
	}

	logActivity(c, "updated", "solution", solutionID, "", "Updated Solution #%d sub-resources", solutionID)
	return c.Render(http.StatusOK, "admin/partials/solution_challenges.html", map[string]interface{}{
		"Challenges": challenges,
	})
}

// DeleteChallenge deletes a challenge from a solution
func (h *SolutionsHandler) DeleteChallenge(c echo.Context) error {
	challengeID, err := strconv.ParseInt(c.Param("challengeId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid challenge ID")
	}

	err = h.queries.DeleteSolutionChallenge(c.Request().Context(), challengeID)
	if err != nil {
		h.logger.Error("Failed to delete solution challenge", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete challenge")
	}

	h.cache.DeleteByPrefix("page:solutions")
	// Note: solutionID not available in DeleteChallenge, use challengeID as reference
	logActivity(c, "updated", "solution", challengeID, "", "Updated Solution #%d sub-resources", challengeID)
	return c.NoContent(http.StatusOK)
}

// AddProduct links a product to a solution
func (h *SolutionsHandler) AddProduct(c echo.Context) error {
	solutionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
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

	isFeatured := c.FormValue("is_featured") == "1"

	params := sqlc.AddProductToSolutionParams{
		SolutionID:   solutionID,
		ProductID:    productID,
		DisplayOrder: sql.NullInt64{Int64: displayOrder, Valid: true},
		IsFeatured:   sql.NullBool{Bool: isFeatured, Valid: true},
	}

	err = h.queries.AddProductToSolution(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to add product to solution", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add product")
	}

	h.cache.DeleteByPrefix("page:solutions")

	products, err := h.queries.GetSolutionProducts(c.Request().Context(), solutionID)
	if err != nil {
		h.logger.Error("Failed to get solution products", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load products")
	}

	logActivity(c, "updated", "solution", solutionID, "", "Updated Solution #%d sub-resources", solutionID)
	return c.Render(http.StatusOK, "admin/partials/solution_products.html", map[string]interface{}{
		"Products": products,
	})
}

// RemoveProduct removes a product from a solution
func (h *SolutionsHandler) RemoveProduct(c echo.Context) error {
	solutionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
	}

	productID, err := strconv.ParseInt(c.Param("productId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid product ID")
	}

	params := sqlc.RemoveProductFromSolutionParams{
		SolutionID: solutionID,
		ProductID:  productID,
	}

	err = h.queries.RemoveProductFromSolution(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to remove product from solution", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to remove product")
	}

	h.cache.DeleteByPrefix("page:solutions")
	logActivity(c, "updated", "solution", solutionID, "", "Updated Solution #%d sub-resources", solutionID)
	return c.NoContent(http.StatusOK)
}

// AddCTA adds a CTA to a solution
func (h *SolutionsHandler) AddCTA(c echo.Context) error {
	solutionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid solution ID")
	}

	heading := c.FormValue("heading")
	subheading := c.FormValue("subheading")
	primaryButtonText := c.FormValue("primary_button_text")
	primaryButtonUrl := c.FormValue("primary_button_url")
	secondaryButtonText := c.FormValue("secondary_button_text")
	secondaryButtonUrl := c.FormValue("secondary_button_url")
	phoneNumber := c.FormValue("phone_number")
	sectionName := c.FormValue("section_name")

	params := sqlc.CreateSolutionCTAParams{
		SolutionID:          solutionID,
		Heading:             heading,
		Subheading:          sql.NullString{String: subheading, Valid: subheading != ""},
		PrimaryButtonText:   sql.NullString{String: primaryButtonText, Valid: primaryButtonText != ""},
		PrimaryButtonUrl:    sql.NullString{String: primaryButtonUrl, Valid: primaryButtonUrl != ""},
		SecondaryButtonText: sql.NullString{String: secondaryButtonText, Valid: secondaryButtonText != ""},
		SecondaryButtonUrl:  sql.NullString{String: secondaryButtonUrl, Valid: secondaryButtonUrl != ""},
		PhoneNumber:         sql.NullString{String: phoneNumber, Valid: phoneNumber != ""},
		SectionName:         sectionName,
	}

	_, err = h.queries.CreateSolutionCTA(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create solution CTA", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to add CTA")
	}

	h.cache.DeleteByPrefix("page:solutions")

	ctas, err := h.queries.GetSolutionCTAs(c.Request().Context(), solutionID)
	if err != nil {
		h.logger.Error("Failed to get solution CTAs", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load CTAs")
	}

	logActivity(c, "updated", "solution", solutionID, "", "Updated Solution #%d sub-resources", solutionID)
	return c.Render(http.StatusOK, "admin/partials/solution_ctas.html", map[string]interface{}{
		"CTAs": ctas,
	})
}

// DeleteCTA deletes a CTA from a solution
func (h *SolutionsHandler) DeleteCTA(c echo.Context) error {
	ctaID, err := strconv.ParseInt(c.Param("ctaId"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid CTA ID")
	}

	err = h.queries.DeleteSolutionCTA(c.Request().Context(), ctaID)
	if err != nil {
		h.logger.Error("Failed to delete solution CTA", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete CTA")
	}

	h.cache.DeleteByPrefix("page:solutions")
	// Note: solutionID not available in DeleteCTA, use ctaID as reference
	logActivity(c, "updated", "solution", ctaID, "", "Updated Solution #%d sub-resources", ctaID)
	return c.NoContent(http.StatusOK)
}
