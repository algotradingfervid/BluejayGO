package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type WhitepaperTopicsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewWhitepaperTopicsHandler(queries *sqlc.Queries, logger *slog.Logger) *WhitepaperTopicsHandler {
	return &WhitepaperTopicsHandler{queries: queries, logger: logger}
}

func (h *WhitepaperTopicsHandler) List(c echo.Context) error {
	items, err := h.queries.ListWhitepaperTopicsWithCount(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list whitepaper topics", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/whitepaper_topics_list.html", map[string]interface{}{
		"Title": "Whitepaper Topics",
		"Items": items,
	})
}

func (h *WhitepaperTopicsHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/whitepaper_topics_form.html", map[string]interface{}{
		"Title":      "New Whitepaper Topic",
		"FormAction": "/admin/whitepaper-topics",
		"Item":       nil,
	})
}

func (h *WhitepaperTopicsHandler) Create(c echo.Context) error {
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	desc := c.FormValue("description")
	_, err := h.queries.CreateWhitepaperTopic(c.Request().Context(), sqlc.CreateWhitepaperTopicParams{
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		ColorHex:    c.FormValue("color_hex"),
		Icon:        c.FormValue("icon"),
		Description: sql.NullString{String: desc, Valid: desc != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to create whitepaper topic", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "whitepaper_topic", 0, c.FormValue("name"), "Created Whitepaper Topic '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/whitepaper-topics")
}

func (h *WhitepaperTopicsHandler) Edit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetWhitepaperTopic(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Topic not found")
	}
	return c.Render(http.StatusOK, "admin/pages/whitepaper_topics_form.html", map[string]interface{}{
		"Title":      "Edit Whitepaper Topic",
		"FormAction": "/admin/whitepaper-topics/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *WhitepaperTopicsHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	desc := c.FormValue("description")
	_, err := h.queries.UpdateWhitepaperTopic(c.Request().Context(), sqlc.UpdateWhitepaperTopicParams{
		ID:          id,
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		ColorHex:    c.FormValue("color_hex"),
		Icon:        c.FormValue("icon"),
		Description: sql.NullString{String: desc, Valid: desc != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to update whitepaper topic", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "whitepaper_topic", id, c.FormValue("name"), "Updated Whitepaper Topic '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/whitepaper-topics")
}

func (h *WhitepaperTopicsHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteWhitepaperTopic(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete whitepaper topic", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "whitepaper_topic", id, "", "Deleted Whitepaper Topic #%d", id)
	return c.NoContent(http.StatusOK)
}
