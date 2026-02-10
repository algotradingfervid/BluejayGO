package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type ProductsHandler struct {
	queries   *sqlc.Queries
	logger    *slog.Logger
	uploadSvc *services.UploadService
	cache     *services.Cache
}

func NewProductsHandler(queries *sqlc.Queries, logger *slog.Logger, uploadSvc *services.UploadService, cache *services.Cache) *ProductsHandler {
	return &ProductsHandler{
		queries:   queries,
		logger:    logger,
		uploadSvc: uploadSvc,
		cache:     cache,
	}
}

const productsPerPage = 15

func (h *ProductsHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	search := c.QueryParam("search")
	status := c.QueryParam("status")
	categoryStr := c.QueryParam("category")
	pageStr := c.QueryParam("page")

	var categoryID int64
	if categoryStr != "" {
		categoryID, _ = strconv.ParseInt(categoryStr, 10, 64)
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * productsPerPage)

	filterParams := sqlc.ListProductsAdminFilteredParams{
		FilterStatus:   status,
		FilterCategory: categoryID,
		FilterSearch:   search,
		PageLimit:      productsPerPage,
		PageOffset:     offset,
	}

	products, err := h.queries.ListProductsAdminFiltered(ctx, filterParams)
	if err != nil {
		h.logger.Error("failed to list products", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	total, err := h.queries.CountProductsAdminFiltered(ctx, sqlc.CountProductsAdminFilteredParams{
		FilterStatus:   status,
		FilterCategory: categoryID,
		FilterSearch:   search,
	})
	if err != nil {
		h.logger.Error("failed to count products", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	categories, err := h.queries.ListProductCategories(ctx)
	if err != nil {
		h.logger.Error("failed to list categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	totalPages := int(math.Ceil(float64(total) / float64(productsPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Build page numbers
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	showFrom := offset + 1
	showTo := offset + int64(len(products))
	if total == 0 {
		showFrom = 0
	}

	hasFilters := search != "" || status != "" || categoryStr != ""

	// Build category map for display
	categoryMap := make(map[int64]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	return c.Render(http.StatusOK, "admin/pages/products_list.html", map[string]interface{}{
		"Title":       "Manage Products",
		"Products":    products,
		"Categories":  categories,
		"CategoryMap": categoryMap,
		"Search":      search,
		"Status":      status,
		"CategoryID":  categoryID,
		"HasFilters":  hasFilters,
		"Page":        page,
		"TotalPages":  totalPages,
		"Pages":       pages,
		"Total":       total,
		"ShowFrom":    showFrom,
		"ShowTo":      showTo,
	})
}

func (h *ProductsHandler) New(c echo.Context) error {
	categories, err := h.queries.ListProductCategories(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/products_form.html", map[string]interface{}{
		"Title":      "New Product",
		"FormAction": "/admin/products",
		"Item":       nil,
		"Categories": categories,
	})
}

func (h *ProductsHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	isFeatured := c.FormValue("is_featured") == "1"
	featuredOrder, _ := strconv.ParseInt(c.FormValue("featured_order"), 10, 64)

	tagline := c.FormValue("tagline")
	overview := c.FormValue("overview")
	metaTitle := c.FormValue("meta_title")
	metaDesc := c.FormValue("meta_description")
	videoURL := c.FormValue("video_url")

	var imagePath sql.NullString
	if fileHeader, err := c.FormFile("primary_image"); err == nil {
		path, err := h.uploadSvc.UploadProductImage(fileHeader)
		if err != nil {
			h.logger.Error("failed to upload image", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Failed to upload image: "+err.Error())
		}
		imagePath = sql.NullString{String: path, Valid: true}
	}

	status := c.FormValue("status")
	var publishedAt sql.NullTime
	if status == "published" {
		publishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	_, err := h.queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:             c.FormValue("sku"),
		Slug:            makeSlug(c.FormValue("name")),
		Name:            c.FormValue("name"),
		Tagline:         sql.NullString{String: tagline, Valid: tagline != ""},
		Description:     c.FormValue("description"),
		Overview:        sql.NullString{String: overview, Valid: overview != ""},
		CategoryID:      categoryID,
		Status:          status,
		IsFeatured:      isFeatured,
		FeaturedOrder:   sql.NullInt64{Int64: featuredOrder, Valid: featuredOrder > 0},
		MetaTitle:       sql.NullString{String: metaTitle, Valid: metaTitle != ""},
		MetaDescription: sql.NullString{String: metaDesc, Valid: metaDesc != ""},
		PrimaryImage:    imagePath,
		VideoUrl:        sql.NullString{String: videoURL, Valid: videoURL != ""},
		PublishedAt:     publishedAt,
	})
	if err != nil {
		h.logger.Error("failed to create product", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	h.cache.DeleteByPrefix("page:products")

	logActivity(c, "created", "product", 0, c.FormValue("name"), "Created Product '%s'", c.FormValue("name"))

	return c.Redirect(http.StatusSeeOther, "/admin/products")
}

func (h *ProductsHandler) Edit(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	product, err := h.queries.GetProduct(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	categories, err := h.queries.ListProductCategories(ctx)
	if err != nil {
		h.logger.Error("failed to list categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Get category slug for preview URL
	var categorySlug string
	category, err := h.queries.GetProductCategory(ctx, product.CategoryID)
	if err == nil {
		categorySlug = category.Slug
	}

	return c.Render(http.StatusOK, "admin/pages/products_form.html", map[string]interface{}{
		"Title":        "Edit Product",
		"FormAction":   fmt.Sprintf("/admin/products/%d", id),
		"Item":         product,
		"Categories":   categories,
		"CategorySlug": categorySlug,
	})
}

func (h *ProductsHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	existing, err := h.queries.GetProduct(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	isFeatured := c.FormValue("is_featured") == "1"
	featuredOrder, _ := strconv.ParseInt(c.FormValue("featured_order"), 10, 64)

	tagline := c.FormValue("tagline")
	overview := c.FormValue("overview")
	metaTitle := c.FormValue("meta_title")
	metaDesc := c.FormValue("meta_description")
	videoURL := c.FormValue("video_url")

	imagePath := existing.PrimaryImage
	if fileHeader, err := c.FormFile("primary_image"); err == nil {
		path, err := h.uploadSvc.UploadProductImage(fileHeader)
		if err != nil {
			h.logger.Error("failed to upload image", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Failed to upload image")
		}
		imagePath = sql.NullString{String: path, Valid: true}
	}

	status := c.FormValue("status")
	publishedAt := existing.PublishedAt
	if status == "published" && !existing.PublishedAt.Valid {
		publishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	err = h.queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		Sku:             c.FormValue("sku"),
		Slug:            makeSlug(c.FormValue("name")),
		Name:            c.FormValue("name"),
		Tagline:         sql.NullString{String: tagline, Valid: tagline != ""},
		Description:     c.FormValue("description"),
		Overview:        sql.NullString{String: overview, Valid: overview != ""},
		CategoryID:      categoryID,
		Status:          status,
		IsFeatured:      isFeatured,
		FeaturedOrder:   sql.NullInt64{Int64: featuredOrder, Valid: featuredOrder > 0},
		MetaTitle:       sql.NullString{String: metaTitle, Valid: metaTitle != ""},
		MetaDescription: sql.NullString{String: metaDesc, Valid: metaDesc != ""},
		PrimaryImage:    imagePath,
		VideoUrl:        sql.NullString{String: videoURL, Valid: videoURL != ""},
		PublishedAt:     publishedAt,
		ID:              id,
	})
	if err != nil {
		h.logger.Error("failed to update product", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	h.cache.DeleteByPrefix("page:products")

	logActivity(c, "updated", "product", id, c.FormValue("name"), "Updated Product '%s'", c.FormValue("name"))

	return c.Redirect(http.StatusSeeOther, "/admin/products")
}

func (h *ProductsHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteProduct(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete product", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:products")

	logActivity(c, "deleted", "product", id, "", "Deleted Product #%d", id)

	return c.NoContent(http.StatusOK)
}
