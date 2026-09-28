package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

var productSlugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// Validate unique fields before uploading files; database constraints still arbitrate races.
func (h *ProductsHandler) productConflict(c echo.Context, id int64) (field string, existingID int64, err error) {
	if !productSlugPattern.MatchString(productFormSlug(c)) {
		return "slug", 0, nil
	}
	checks := []struct {
		field, value string
		find         func(string) (sqlc.Product, error)
	}{
		{"sku", c.FormValue("sku"), func(v string) (sqlc.Product, error) { return h.queries.GetProductBySKU(c.Request().Context(), v) }},
		{"slug", productFormSlug(c), func(v string) (sqlc.Product, error) { return h.queries.GetProductBySlug(c.Request().Context(), v) }},
	}
	for _, check := range checks {
		item, e := check.find(check.value)
		if e == nil && item.ID != id {
			return check.field, item.ID, nil
		}
		if e != nil && !errors.Is(e, sql.ErrNoRows) {
			return "", 0, e
		}
	}
	return "", 0, nil
}

func productFormSlug(c echo.Context) string {
	if slug := strings.TrimSpace(c.FormValue("slug")); slug != "" {
		return slug
	}
	return makeSlug(c.FormValue("name"))
}

func (h *ProductsHandler) renderProductConflict(c echo.Context, item sqlc.Product, field string, existingID int64) error {
	categories, err := h.queries.ListProductCategories(c.Request().Context())
	if err != nil {
		return err
	}
	previewURL := ""
	if item.ID > 0 {
		if category, err := h.queries.GetProductCategory(c.Request().Context(), item.CategoryID); err == nil {
			previewURL = fmt.Sprintf("/products/%s/%s?preview=true", category.Slug, item.Slug)
		}
	}
	item.Name = c.FormValue("name")
	item.Sku = c.FormValue("sku")
	item.Slug = productFormSlug(c)
	item.Description = c.FormValue("description")
	item.Status = c.FormValue("status")
	item.CategoryID, _ = strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	item.IsFeatured = c.FormValue("is_featured") == "1"
	order, _ := strconv.ParseInt(c.FormValue("featured_order"), 10, 64)
	item.FeaturedOrder = sql.NullInt64{Int64: order, Valid: c.FormValue("featured_order") != ""}
	for name, target := range map[string]*sql.NullString{"tagline": &item.Tagline, "overview": &item.Overview, "meta_title": &item.MetaTitle, "meta_description": &item.MetaDescription, "video_url": &item.VideoUrl} {
		value := c.FormValue(name)
		*target = sql.NullString{String: value, Valid: value != ""}
	}
	title, action := "New Product", "/admin/products"
	if item.ID > 0 {
		title = "Edit Product"
		action = fmt.Sprintf("/admin/products/%d", item.ID)
	}
	categorySlug := ""
	for _, cat := range categories {
		if cat.ID == item.CategoryID {
			categorySlug = cat.Slug
		}
	}
	label := "SKU"
	if field == "slug" {
		label = "Slug"
	}
	message := label + " already exists. Choose a different value."
	if existingID == 0 {
		message = "Slug must start with a lowercase letter or number and contain only lowercase letters, numbers and hyphens."
	}
	return c.Render(http.StatusUnprocessableEntity, "admin/pages/products_form.html", map[string]interface{}{
		"Title": title, "FormAction": action, "Item": item, "Categories": categories, "CategorySlug": categorySlug,
		"PreviewURL": previewURL, "ConflictField": field, "ConflictID": existingID, "FormError": message,
	})
}
