package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type BlogAuthorsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewBlogAuthorsHandler(queries *sqlc.Queries, logger *slog.Logger) *BlogAuthorsHandler {
	return &BlogAuthorsHandler{queries: queries, logger: logger}
}

func (h *BlogAuthorsHandler) List(c echo.Context) error {
	items, err := h.queries.ListBlogAuthors(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list blog authors", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/blog_authors_list.html", map[string]interface{}{
		"Title": "Blog Authors",
		"Items": items,
	})
}

func (h *BlogAuthorsHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/blog_authors_form.html", map[string]interface{}{
		"Title":      "New Blog Author",
		"FormAction": "/admin/blog-authors",
		"Item":       nil,
	})
}

func (h *BlogAuthorsHandler) Create(c echo.Context) error {
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	bio := c.FormValue("bio")
	avatarUrl := c.FormValue("avatar_url")
	linkedinUrl := c.FormValue("linkedin_url")
	email := c.FormValue("email")

	_, err := h.queries.CreateBlogAuthor(c.Request().Context(), sqlc.CreateBlogAuthorParams{
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Title:       c.FormValue("title"),
		Bio:         sql.NullString{String: bio, Valid: bio != ""},
		AvatarUrl:   sql.NullString{String: avatarUrl, Valid: avatarUrl != ""},
		LinkedinUrl: sql.NullString{String: linkedinUrl, Valid: linkedinUrl != ""},
		Email:       sql.NullString{String: email, Valid: email != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to create blog author", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "blog_author", 0, c.FormValue("name"), "Created blog_author '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/blog-authors")
}

func (h *BlogAuthorsHandler) Edit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetBlogAuthor(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Author not found")
	}
	return c.Render(http.StatusOK, "admin/pages/blog_authors_form.html", map[string]interface{}{
		"Title":      "Edit Blog Author",
		"FormAction": "/admin/blog-authors/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *BlogAuthorsHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	bio := c.FormValue("bio")
	avatarUrl := c.FormValue("avatar_url")
	linkedinUrl := c.FormValue("linkedin_url")
	email := c.FormValue("email")

	_, err := h.queries.UpdateBlogAuthor(c.Request().Context(), sqlc.UpdateBlogAuthorParams{
		ID:          id,
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Title:       c.FormValue("title"),
		Bio:         sql.NullString{String: bio, Valid: bio != ""},
		AvatarUrl:   sql.NullString{String: avatarUrl, Valid: avatarUrl != ""},
		LinkedinUrl: sql.NullString{String: linkedinUrl, Valid: linkedinUrl != ""},
		Email:       sql.NullString{String: email, Valid: email != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to update blog author", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "blog_author", id, c.FormValue("name"), "Updated blog_author '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/blog-authors")
}

func (h *BlogAuthorsHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteBlogAuthor(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete blog author", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "blog_author", id, "", "Deleted blog_author #%d", id)
	return c.NoContent(http.StatusOK)
}
