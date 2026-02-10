package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

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
	return &AboutHandler{queries: queries, logger: logger, cache: cache}
}

// Company Overview
func (h *AboutHandler) OverviewEdit(c echo.Context) error {
	overview, _ := h.queries.GetCompanyOverview(c.Request().Context())
	return c.Render(http.StatusOK, "admin/pages/about_overview_form.html", map[string]interface{}{
		"Title": "Company Overview",
		"Item":  overview,
	})
}

func (h *AboutHandler) OverviewUpdate(c echo.Context) error {
	_, err := h.queries.UpsertCompanyOverview(c.Request().Context(), sqlc.UpsertCompanyOverviewParams{
		Headline:             c.FormValue("headline"),
		Tagline:              c.FormValue("tagline"),
		DescriptionMain:      c.FormValue("description_main"),
		DescriptionSecondary: sql.NullString{String: c.FormValue("description_secondary"), Valid: c.FormValue("description_secondary") != ""},
		DescriptionTertiary:  sql.NullString{String: c.FormValue("description_tertiary"), Valid: c.FormValue("description_tertiary") != ""},
		HeroImageUrl:         sql.NullString{String: c.FormValue("hero_image_url"), Valid: c.FormValue("hero_image_url") != ""},
		CompanyImageUrl:      sql.NullString{String: c.FormValue("company_image_url"), Valid: c.FormValue("company_image_url") != ""},
	})
	if err != nil {
		h.logger.Error("failed to upsert company overview", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "updated", "about", 0, "Overview", "Updated About Overview")
	return c.Redirect(http.StatusSeeOther, "/admin/about/overview")
}

// Mission Vision Values
func (h *AboutHandler) MVVEdit(c echo.Context) error {
	mvv, _ := h.queries.GetMissionVisionValues(c.Request().Context())
	return c.Render(http.StatusOK, "admin/pages/about_mvv_form.html", map[string]interface{}{
		"Title": "Mission, Vision & Values",
		"Item":  mvv,
	})
}

func (h *AboutHandler) MVVUpdate(c echo.Context) error {
	_, err := h.queries.UpsertMissionVisionValues(c.Request().Context(), sqlc.UpsertMissionVisionValuesParams{
		Mission:       c.FormValue("mission"),
		Vision:        c.FormValue("vision"),
		ValuesSummary: sql.NullString{String: c.FormValue("values_summary"), Valid: c.FormValue("values_summary") != ""},
		MissionIcon:   c.FormValue("mission_icon"),
		VisionIcon:    c.FormValue("vision_icon"),
		ValuesIcon:    c.FormValue("values_icon"),
	})
	if err != nil {
		h.logger.Error("failed to upsert mvv", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "updated", "about", 0, "MVV", "Updated About Mission/Vision/Values")
	return c.Redirect(http.StatusSeeOther, "/admin/about/mvv")
}

// Core Values
func (h *AboutHandler) CoreValuesList(c echo.Context) error {
	items, err := h.queries.ListCoreValues(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list core values", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/core_values_list.html", map[string]interface{}{
		"Title": "Core Values",
		"Items": items,
	})
}

func (h *AboutHandler) CoreValueNew(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/core_values_form.html", map[string]interface{}{
		"Title":      "New Core Value",
		"FormAction": "/admin/about/values",
		"Item":       nil,
	})
}

func (h *AboutHandler) CoreValueCreate(c echo.Context) error {
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	_, err := h.queries.CreateCoreValue(c.Request().Context(), sqlc.CreateCoreValueParams{
		Title:        c.FormValue("title"),
		Description:  c.FormValue("description"),
		Icon:         c.FormValue("icon"),
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create core value", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "created", "core_value", 0, c.FormValue("title"), "Created Core Value '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/about/values")
}

func (h *AboutHandler) CoreValueEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetCoreValue(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.Render(http.StatusOK, "admin/pages/core_values_form.html", map[string]interface{}{
		"Title":      "Edit Core Value",
		"FormAction": "/admin/about/values/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *AboutHandler) CoreValueUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	_, err := h.queries.UpdateCoreValue(c.Request().Context(), sqlc.UpdateCoreValueParams{
		Title:        c.FormValue("title"),
		Description:  c.FormValue("description"),
		Icon:         c.FormValue("icon"),
		DisplayOrder: order,
		ID:           id,
	})
	if err != nil {
		h.logger.Error("failed to update core value", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "updated", "core_value", id, c.FormValue("title"), "Updated Core Value '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/about/values")
}

func (h *AboutHandler) CoreValueDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteCoreValue(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete core value", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "deleted", "core_value", id, "", "Deleted Core Value #%d", id)
	return c.NoContent(http.StatusOK)
}

// Milestones
func (h *AboutHandler) MilestonesList(c echo.Context) error {
	items, err := h.queries.ListMilestones(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list milestones", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/milestones_list.html", map[string]interface{}{
		"Title": "Milestones",
		"Items": items,
	})
}

func (h *AboutHandler) MilestoneNew(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/milestones_form.html", map[string]interface{}{
		"Title":      "New Milestone",
		"FormAction": "/admin/about/milestones",
		"Item":       nil,
	})
}

func (h *AboutHandler) MilestoneCreate(c echo.Context) error {
	year, _ := strconv.ParseInt(c.FormValue("year"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	isCurrent := int64(0)
	if c.FormValue("is_current") == "on" {
		isCurrent = 1
	}
	_, err := h.queries.CreateMilestone(c.Request().Context(), sqlc.CreateMilestoneParams{
		Year:         year,
		Title:        c.FormValue("title"),
		Description:  c.FormValue("description"),
		IsCurrent:    isCurrent,
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create milestone", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "created", "milestone", 0, c.FormValue("title"), "Created Milestone '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/about/milestones")
}

func (h *AboutHandler) MilestoneEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetMilestone(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.Render(http.StatusOK, "admin/pages/milestones_form.html", map[string]interface{}{
		"Title":      "Edit Milestone",
		"FormAction": "/admin/about/milestones/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *AboutHandler) MilestoneUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	year, _ := strconv.ParseInt(c.FormValue("year"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	isCurrent := int64(0)
	if c.FormValue("is_current") == "on" {
		isCurrent = 1
	}
	_, err := h.queries.UpdateMilestone(c.Request().Context(), sqlc.UpdateMilestoneParams{
		Year:         year,
		Title:        c.FormValue("title"),
		Description:  c.FormValue("description"),
		IsCurrent:    isCurrent,
		DisplayOrder: order,
		ID:           id,
	})
	if err != nil {
		h.logger.Error("failed to update milestone", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "updated", "milestone", id, c.FormValue("title"), "Updated Milestone '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/about/milestones")
}

func (h *AboutHandler) MilestoneDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteMilestone(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete milestone", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "deleted", "milestone", id, "", "Deleted Milestone #%d", id)
	return c.NoContent(http.StatusOK)
}

// Certifications
func (h *AboutHandler) CertificationsList(c echo.Context) error {
	items, err := h.queries.ListCertifications(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list certifications", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/certifications_list.html", map[string]interface{}{
		"Title": "Certifications",
		"Items": items,
	})
}

func (h *AboutHandler) CertificationNew(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/certifications_form.html", map[string]interface{}{
		"Title":      "New Certification",
		"FormAction": "/admin/about/certifications",
		"Item":       nil,
	})
}

func (h *AboutHandler) CertificationCreate(c echo.Context) error {
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	_, err := h.queries.CreateCertification(c.Request().Context(), sqlc.CreateCertificationParams{
		Name:         c.FormValue("name"),
		Abbreviation: c.FormValue("abbreviation"),
		Description:  sql.NullString{String: c.FormValue("description"), Valid: c.FormValue("description") != ""},
		Icon:         sql.NullString{String: c.FormValue("icon"), Valid: c.FormValue("icon") != ""},
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create certification", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "created", "certification", 0, c.FormValue("name"), "Created Certification '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/about/certifications")
}

func (h *AboutHandler) CertificationEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetCertification(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.Render(http.StatusOK, "admin/pages/certifications_form.html", map[string]interface{}{
		"Title":      "Edit Certification",
		"FormAction": "/admin/about/certifications/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *AboutHandler) CertificationUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	_, err := h.queries.UpdateCertification(c.Request().Context(), sqlc.UpdateCertificationParams{
		Name:         c.FormValue("name"),
		Abbreviation: c.FormValue("abbreviation"),
		Description:  sql.NullString{String: c.FormValue("description"), Valid: c.FormValue("description") != ""},
		Icon:         sql.NullString{String: c.FormValue("icon"), Valid: c.FormValue("icon") != ""},
		DisplayOrder: order,
		ID:           id,
	})
	if err != nil {
		h.logger.Error("failed to update certification", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "updated", "certification", id, c.FormValue("name"), "Updated Certification '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/about/certifications")
}

func (h *AboutHandler) CertificationDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteCertification(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete certification", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:about")
	logActivity(c, "deleted", "certification", id, "", "Deleted Certification #%d", id)
	return c.NoContent(http.StatusOK)
}
