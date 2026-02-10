package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type BlogCategoriesHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewBlogCategoriesHandler(queries *sqlc.Queries, logger *slog.Logger) *BlogCategoriesHandler {
	return &BlogCategoriesHandler{queries: queries, logger: logger}
}

func (h *BlogCategoriesHandler) List(c echo.Context) error {
	items, err := h.queries.ListBlogCategories(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list blog categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/blog_categories_list.html", map[string]interface{}{
		"Title": "Blog Categories",
		"Items": items,
	})
}

func (h *BlogCategoriesHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/blog_categories_form.html", map[string]interface{}{
		"Title":      "New Blog Category",
		"FormAction": "/admin/blog-categories",
		"Item":       nil,
	})
}

func (h *BlogCategoriesHandler) Create(c echo.Context) error {
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	desc := c.FormValue("description")
	_, err := h.queries.CreateBlogCategory(c.Request().Context(), sqlc.CreateBlogCategoryParams{
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		ColorHex:    c.FormValue("color_hex"),
		Description: sql.NullString{String: desc, Valid: desc != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to create blog category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "blog_category", 0, c.FormValue("name"), "Created blog_category '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/blog-categories")
}

func (h *BlogCategoriesHandler) Edit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetBlogCategory(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}
	return c.Render(http.StatusOK, "admin/pages/blog_categories_form.html", map[string]interface{}{
		"Title":      "Edit Blog Category",
		"FormAction": "/admin/blog-categories/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *BlogCategoriesHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	desc := c.FormValue("description")
	_, err := h.queries.UpdateBlogCategory(c.Request().Context(), sqlc.UpdateBlogCategoryParams{
		ID:          id,
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		ColorHex:    c.FormValue("color_hex"),
		Description: sql.NullString{String: desc, Valid: desc != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to update blog category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "blog_category", id, c.FormValue("name"), "Updated blog_category '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/blog-categories")
}

func (h *BlogCategoriesHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteBlogCategory(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete blog category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "blog_category", id, "", "Deleted blog_category #%d", id)
	return c.NoContent(http.StatusOK)
}
