package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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
	return &PartnersHandler{queries: queries, logger: logger, cache: cache}
}

func (h *PartnersHandler) List(c echo.Context) error {
	partners, err := h.queries.ListAllPartners(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list partners", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	tiers, _ := h.queries.ListPartnerTiers(c.Request().Context())

	// Filter parameters
	search := strings.TrimSpace(c.QueryParam("search"))
	tierFilter := c.QueryParam("tier")
	status := c.QueryParam("status")
	hasFilters := search != "" || tierFilter != "" || status != ""

	var tierFilterID int64
	if tierFilter != "" {
		tierFilterID, _ = strconv.ParseInt(tierFilter, 10, 64)
	}

	// Apply in-memory filtering
	var filtered []sqlc.ListAllPartnersRow
	for _, p := range partners {
		if search != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
			continue
		}
		if tierFilterID > 0 && p.TierID != tierFilterID {
			continue
		}
		if status == "active" && p.IsActive == 0 {
			continue
		}
		if status == "inactive" && p.IsActive != 0 {
			continue
		}
		filtered = append(filtered, p)
	}

	return c.Render(http.StatusOK, "admin/pages/partners_list.html", map[string]interface{}{
		"Title":      "Partners",
		"Items":      filtered,
		"Tiers":      tiers,
		"Search":     search,
		"TierFilter": tierFilterID,
		"Status":     status,
		"HasFilters": hasFilters,
	})
}

func (h *PartnersHandler) New(c echo.Context) error {
	tiers, _ := h.queries.ListPartnerTiers(c.Request().Context())
	return c.Render(http.StatusOK, "admin/pages/partners_form.html", map[string]interface{}{
		"Title":      "New Partner",
		"FormAction": "/admin/partners",
		"Item":       nil,
		"Tiers":      tiers,
	})
}

func (h *PartnersHandler) Create(c echo.Context) error {
	tierID, _ := strconv.ParseInt(c.FormValue("tier_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	_, err := h.queries.CreatePartner(c.Request().Context(), sqlc.CreatePartnerParams{
		Name:         c.FormValue("name"),
		TierID:       tierID,
		LogoUrl:      sql.NullString{String: c.FormValue("logo_url"), Valid: c.FormValue("logo_url") != ""},
		Icon:         sql.NullString{String: c.FormValue("icon"), Valid: c.FormValue("icon") != ""},
		WebsiteUrl:   sql.NullString{String: c.FormValue("website_url"), Valid: c.FormValue("website_url") != ""},
		Description:  sql.NullString{String: c.FormValue("description"), Valid: c.FormValue("description") != ""},
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create partner", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:partners")
	logActivity(c, "created", "partner", 0, c.FormValue("name"), "Created Partner '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/partners")
}

func (h *PartnersHandler) Edit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetPartner(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	tiers, _ := h.queries.ListPartnerTiers(c.Request().Context())
	return c.Render(http.StatusOK, "admin/pages/partners_form.html", map[string]interface{}{
		"Title":      "Edit Partner",
		"FormAction": "/admin/partners/" + c.Param("id"),
		"Item":       item,
		"Tiers":      tiers,
	})
}

func (h *PartnersHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	tierID, _ := strconv.ParseInt(c.FormValue("tier_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	isActive := int64(1)
	if c.FormValue("is_active") == "" || c.FormValue("is_active") == "0" {
		isActive = 0
	}
	_, err := h.queries.UpdatePartner(c.Request().Context(), sqlc.UpdatePartnerParams{
		Name:         c.FormValue("name"),
		TierID:       tierID,
		LogoUrl:      sql.NullString{String: c.FormValue("logo_url"), Valid: c.FormValue("logo_url") != ""},
		Icon:         sql.NullString{String: c.FormValue("icon"), Valid: c.FormValue("icon") != ""},
		WebsiteUrl:   sql.NullString{String: c.FormValue("website_url"), Valid: c.FormValue("website_url") != ""},
		Description:  sql.NullString{String: c.FormValue("description"), Valid: c.FormValue("description") != ""},
		DisplayOrder: order,
		IsActive:     isActive,
		ID:           id,
	})
	if err != nil {
		h.logger.Error("failed to update partner", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:partners")
	logActivity(c, "updated", "partner", id, c.FormValue("name"), "Updated Partner '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/partners")
}

func (h *PartnersHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeletePartner(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete partner", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:partners")
	logActivity(c, "deleted", "partner", id, "", "Deleted Partner #%d", id)
	return c.NoContent(http.StatusOK)
}

// Testimonials
func (h *PartnersHandler) TestimonialsList(c echo.Context) error {
	items, err := h.queries.ListActiveTestimonials(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list testimonials", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/testimonials_list.html", map[string]interface{}{
		"Title": "Partner Testimonials",
		"Items": items,
	})
}

func (h *PartnersHandler) TestimonialNew(c echo.Context) error {
	partners, _ := h.queries.ListAllPartners(c.Request().Context())
	return c.Render(http.StatusOK, "admin/pages/testimonials_form.html", map[string]interface{}{
		"Title":      "New Testimonial",
		"FormAction": "/admin/partners/testimonials",
		"Item":       nil,
		"Partners":   partners,
	})
}

func (h *PartnersHandler) TestimonialCreate(c echo.Context) error {
	partnerID, _ := strconv.ParseInt(c.FormValue("partner_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	_, err := h.queries.CreateTestimonial(c.Request().Context(), sqlc.CreateTestimonialParams{
		PartnerID:    partnerID,
		Quote:        c.FormValue("quote"),
		AuthorName:   c.FormValue("author_name"),
		AuthorTitle:  c.FormValue("author_title"),
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create testimonial", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:partners")
	logActivity(c, "created", "partner_testimonial", 0, c.FormValue("author_name"), "Created Partner Testimonial by '%s'", c.FormValue("author_name"))
	return c.Redirect(http.StatusSeeOther, "/admin/partners/testimonials")
}

func (h *PartnersHandler) TestimonialEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetTestimonial(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	partners, _ := h.queries.ListAllPartners(c.Request().Context())
	return c.Render(http.StatusOK, "admin/pages/testimonials_form.html", map[string]interface{}{
		"Title":      "Edit Testimonial",
		"FormAction": "/admin/partners/testimonials/" + c.Param("id"),
		"Item":       item,
		"Partners":   partners,
	})
}

func (h *PartnersHandler) TestimonialUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	partnerID, _ := strconv.ParseInt(c.FormValue("partner_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	isActive := int64(1)
	if c.FormValue("is_active") == "" {
		isActive = 0
	}
	_, err := h.queries.UpdateTestimonial(c.Request().Context(), sqlc.UpdateTestimonialParams{
		PartnerID:    partnerID,
		Quote:        c.FormValue("quote"),
		AuthorName:   c.FormValue("author_name"),
		AuthorTitle:  c.FormValue("author_title"),
		DisplayOrder: order,
		IsActive:     isActive,
		ID:           id,
	})
	if err != nil {
		h.logger.Error("failed to update testimonial", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:partners")
	logActivity(c, "updated", "partner_testimonial", id, c.FormValue("author_name"), "Updated Partner Testimonial by '%s'", c.FormValue("author_name"))
	return c.Redirect(http.StatusSeeOther, "/admin/partners/testimonials")
}

func (h *PartnersHandler) TestimonialDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteTestimonial(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete testimonial", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:partners")
	logActivity(c, "deleted", "partner_testimonial", id, "", "Deleted Partner Testimonial #%d", id)
	return c.NoContent(http.StatusOK)
}
