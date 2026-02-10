package admin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type PartnerTiersHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewPartnerTiersHandler(queries *sqlc.Queries, logger *slog.Logger) *PartnerTiersHandler {
	return &PartnerTiersHandler{queries: queries, logger: logger}
}

func (h *PartnerTiersHandler) List(c echo.Context) error {
	items, err := h.queries.ListPartnerTiers(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list partner tiers", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/partner_tiers_list.html", map[string]interface{}{
		"Title": "Partner Tiers",
		"Items": items,
	})
}

func (h *PartnerTiersHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/partner_tiers_form.html", map[string]interface{}{
		"Title":      "New Partner Tier",
		"FormAction": "/admin/partner-tiers",
		"Item":       nil,
	})
}

func (h *PartnerTiersHandler) Create(c echo.Context) error {
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	_, err := h.queries.CreatePartnerTier(c.Request().Context(), sqlc.CreatePartnerTierParams{
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Description: c.FormValue("description"),
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to create partner tier", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "partner_tier", 0, c.FormValue("name"), "Created Partner Tier '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/partner-tiers")
}

func (h *PartnerTiersHandler) Edit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetPartnerTier(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Partner tier not found")
	}
	return c.Render(http.StatusOK, "admin/pages/partner_tiers_form.html", map[string]interface{}{
		"Title":      "Edit Partner Tier",
		"FormAction": "/admin/partner-tiers/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *PartnerTiersHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	_, err := h.queries.UpdatePartnerTier(c.Request().Context(), sqlc.UpdatePartnerTierParams{
		ID:          id,
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Description: c.FormValue("description"),
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to update partner tier", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "partner_tier", id, c.FormValue("name"), "Updated Partner Tier '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/partner-tiers")
}

func (h *PartnerTiersHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeletePartnerTier(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete partner tier", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "partner_tier", id, "", "Deleted Partner Tier #%d", id)
	return c.NoContent(http.StatusOK)
}
