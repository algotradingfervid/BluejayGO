package admin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type IndustriesHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewIndustriesHandler(queries *sqlc.Queries, logger *slog.Logger) *IndustriesHandler {
	return &IndustriesHandler{queries: queries, logger: logger}
}

func (h *IndustriesHandler) List(c echo.Context) error {
	items, err := h.queries.ListIndustries(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list industries", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/industries_list.html", map[string]interface{}{
		"Title": "Industries",
		"Items": items,
	})
}

func (h *IndustriesHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/industries_form.html", map[string]interface{}{
		"Title":      "New Industry",
		"FormAction": "/admin/industries",
		"Item":       nil,
	})
}

func (h *IndustriesHandler) Create(c echo.Context) error {
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	_, err := h.queries.CreateIndustry(c.Request().Context(), sqlc.CreateIndustryParams{
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Icon:        c.FormValue("icon"),
		Description: c.FormValue("description"),
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to create industry", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "industry", 0, c.FormValue("name"), "Created Industry '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/industries")
}

func (h *IndustriesHandler) Edit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetIndustry(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Industry not found")
	}
	return c.Render(http.StatusOK, "admin/pages/industries_form.html", map[string]interface{}{
		"Title":      "Edit Industry",
		"FormAction": "/admin/industries/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *IndustriesHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	_, err := h.queries.UpdateIndustry(c.Request().Context(), sqlc.UpdateIndustryParams{
		ID:          id,
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Icon:        c.FormValue("icon"),
		Description: c.FormValue("description"),
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to update industry", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "industry", id, c.FormValue("name"), "Updated Industry '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/industries")
}

func (h *IndustriesHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteIndustry(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete industry", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "industry", id, "", "Deleted Industry #%d", id)
	return c.NoContent(http.StatusOK)
}
