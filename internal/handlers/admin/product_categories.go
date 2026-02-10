package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

var slugRegexp = regexp.MustCompile(`[^a-z0-9-]+`)

func makeSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRegexp.ReplaceAllString(s, "")
	return s
}

type ProductCategoriesHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewProductCategoriesHandler(queries *sqlc.Queries, logger *slog.Logger) *ProductCategoriesHandler {
	return &ProductCategoriesHandler{queries: queries, logger: logger}
}

func (h *ProductCategoriesHandler) List(c echo.Context) error {
	items, err := h.queries.ListProductCategories(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list product categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/product_categories_list.html", map[string]interface{}{
		"Title": "Product Categories",
		"Items": items,
	})
}

func (h *ProductCategoriesHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/product_categories_form.html", map[string]interface{}{
		"Title":      "New Product Category",
		"FormAction": "/admin/product-categories",
		"Item":       nil,
	})
}

func (h *ProductCategoriesHandler) Create(c echo.Context) error {
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	imageUrl := c.FormValue("image_url")
	_, err := h.queries.CreateProductCategory(c.Request().Context(), sqlc.CreateProductCategoryParams{
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Description: c.FormValue("description"),
		Icon:        c.FormValue("icon"),
		ImageUrl:    sql.NullString{String: imageUrl, Valid: imageUrl != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to create product category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "product_category", 0, c.FormValue("name"), "Created Product Category '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/product-categories")
}

func (h *ProductCategoriesHandler) Edit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetProductCategory(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}
	return c.Render(http.StatusOK, "admin/pages/product_categories_form.html", map[string]interface{}{
		"Title":      "Edit Product Category",
		"FormAction": "/admin/product-categories/" + c.Param("id"),
		"Item":       item,
	})
}

func (h *ProductCategoriesHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	sortOrder, _ := strconv.ParseInt(c.FormValue("sort_order"), 10, 64)
	imageUrl := c.FormValue("image_url")
	_, err := h.queries.UpdateProductCategory(c.Request().Context(), sqlc.UpdateProductCategoryParams{
		ID:          id,
		Name:        c.FormValue("name"),
		Slug:        makeSlug(c.FormValue("name")),
		Description: c.FormValue("description"),
		Icon:        c.FormValue("icon"),
		ImageUrl:    sql.NullString{String: imageUrl, Valid: imageUrl != ""},
		SortOrder:   sortOrder,
	})
	if err != nil {
		h.logger.Error("failed to update product category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "product_category", id, c.FormValue("name"), "Updated Product Category '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/product-categories")
}

func (h *ProductCategoriesHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteProductCategory(c.Request().Context(), id); err != nil {
		if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			return c.String(http.StatusConflict, "Cannot delete this category because it still has products assigned to it. Please reassign or remove those products first.")
		}
		h.logger.Error("failed to delete product category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "product_category", id, "", "Deleted Product Category #%d", id)
	return c.NoContent(http.StatusOK)
}
