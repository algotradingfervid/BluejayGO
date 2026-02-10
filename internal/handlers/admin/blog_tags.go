package admin

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type BlogTagsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewBlogTagsHandler(queries *sqlc.Queries, logger *slog.Logger) *BlogTagsHandler {
	return &BlogTagsHandler{queries: queries, logger: logger}
}

func (h *BlogTagsHandler) List(c echo.Context) error {
	tags, err := h.queries.ListAllBlogTags(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list blog tags", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/blog_tags_list.html", map[string]interface{}{
		"Title": "Blog Tags",
		"Items": tags,
	})
}

func (h *BlogTagsHandler) Create(c echo.Context) error {
	name := c.FormValue("name")
	_, err := h.queries.CreateBlogTag(c.Request().Context(), sqlc.CreateBlogTagParams{
		Name: name,
		Slug: makeSlug(name),
	})
	if err != nil {
		h.logger.Error("failed to create blog tag", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "blog_tag", 0, name, "Created blog_tag '%s'", name)
	return c.Redirect(http.StatusSeeOther, "/admin/blog/tags")
}

func (h *BlogTagsHandler) Search(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("_tag_search"))
	if query == "" {
		return c.Render(http.StatusOK, "admin/partials/tag_suggestions.html", map[string]interface{}{
			"Tags":  nil,
			"Query": "",
		})
	}
	tags, _ := h.queries.SearchBlogTags(c.Request().Context(), "%"+query+"%")
	return c.Render(http.StatusOK, "admin/partials/tag_suggestions.html", map[string]interface{}{
		"Tags":  tags,
		"Query": query,
	})
}

func (h *BlogTagsHandler) QuickCreate(c echo.Context) error {
	name := strings.TrimSpace(c.FormValue("name"))
	if name == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	tag, err := h.queries.CreateBlogTag(c.Request().Context(), sqlc.CreateBlogTagParams{
		Name: name,
		Slug: makeSlug(name),
	})
	if err != nil {
		h.logger.Error("failed to quick-create blog tag", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "blog_tag", 0, name, "Created blog_tag '%s'", name)
	return c.Render(http.StatusOK, "admin/partials/tag_chip.html", map[string]interface{}{
		"ID":   tag.ID,
		"Name": tag.Name,
	})
}

func (h *BlogTagsHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteBlogTag(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete blog tag", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "blog_tag", id, "", "Deleted blog_tag #%d", id)
	return c.NoContent(http.StatusOK)
}
