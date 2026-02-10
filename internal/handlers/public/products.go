package public

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type ProductsHandler struct {
	queries    *sqlc.Queries
	logger     *slog.Logger
	productSvc *services.ProductService
	cache      *services.Cache
}

func NewProductsHandler(queries *sqlc.Queries, logger *slog.Logger, productSvc *services.ProductService, cache *services.Cache) *ProductsHandler {
	return &ProductsHandler{
		queries:    queries,
		logger:     logger,
		productSvc: productSvc,
		cache:      cache,
	}
}

func (h *ProductsHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}
	if cats := c.Get("footer_categories"); cats != nil {
		data["FooterCategories"] = cats
	}
	if sols := c.Get("footer_solutions"); sols != nil {
		data["FooterSolutions"] = sols
	}
	if res := c.Get("footer_resources"); res != nil {
		data["FooterResources"] = res
	}
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		h.logger.Error("template render failed", "template", templateName, "error", err)
		return err
	}
	html := buf.String()
	h.cache.Set(cacheKey, html, ttlSeconds)
	return c.HTML(statusCode, html)
}

// GET /products
func (h *ProductsHandler) ProductsList(c echo.Context) error {
	cacheKey := "page:products"
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	categories, err := h.queries.ListProductCategories(ctx)
	if err != nil {
		h.logger.Error("failed to list product categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	type categoryWithCount struct {
		Category sqlc.ProductCategory
		Count    int64
	}

	var categoriesWithCount []categoryWithCount
	for _, cat := range categories {
		count, _ := h.queries.CountProductsByCategory(ctx, cat.ID)
		categoriesWithCount = append(categoriesWithCount, categoryWithCount{
			Category: cat,
			Count:    count,
		})
	}

	heroSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products", SectionKey: "hero"})
	categoriesSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products", SectionKey: "categories_section"})
	ctaSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products", SectionKey: "cta"})

	data := map[string]interface{}{
		"Title":             "Products",
		"Categories":        categoriesWithCount,
		"TotalCount":        len(categories),
		"PageHero":          heroSection,
		"CategoriesSection": categoriesSection,
		"PageCTA":           ctaSection,
	}

	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/products.html", data)
}

// GET /products/:category
func (h *ProductsHandler) ProductsByCategory(c echo.Context) error {
	ctx := c.Request().Context()
	categorySlug := c.Param("category")

	cacheKey := fmt.Sprintf("page:products:%s", categorySlug)
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	category, err := h.queries.GetProductCategoryBySlug(ctx, categorySlug)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}
	if err != nil {
		h.logger.Error("failed to load category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	limit := int64(12)
	offset := int64((page - 1)) * limit

	products, err := h.queries.ListProductsByCategory(ctx, sqlc.ListProductsByCategoryParams{
		CategoryID: category.ID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		h.logger.Error("failed to load products", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	total, _ := h.queries.CountProductsByCategory(ctx, category.ID)
	totalPages := int((total + limit - 1) / limit)

	categoryHero, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products_category", SectionKey: "hero"})
	emptyState, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products_category", SectionKey: "empty_state"})

	data := map[string]interface{}{
		"Title":        fmt.Sprintf("%s | Products", category.Name),
		"Category":     category,
		"Products":     products,
		"TotalCount":   total,
		"CurrentPage":  page,
		"TotalPages":   totalPages,
		"CategoryHero": categoryHero,
		"EmptyState":   emptyState,
	}

	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/products_category.html", data)
}

// GET /products/:category/:slug
func (h *ProductsHandler) ProductDetail(c echo.Context) error {
	categorySlug := c.Param("category")
	productSlug := c.Param("slug")
	preview := isPreviewRequest(c)

	if !preview {
		cacheKey := fmt.Sprintf("page:products:%s:%s", categorySlug, productSlug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	detail, err := h.productSvc.GetProductDetail(c.Request().Context(), productSlug)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	if err != nil {
		h.logger.Error("failed to load product detail", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if detail.Category.Slug != categorySlug {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found in this category")
	}

	specSections := groupSpecsBySection(detail.Specs)

	ctx := c.Request().Context()
	detailCTA, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "product_detail", SectionKey: "cta"})
	// Replace placeholders in CTA
	replacer := strings.NewReplacer("{product_name}", detail.Product.Name, "{product_sku}", detail.Product.Sku)
	detailCTA.Heading = replacer.Replace(detailCTA.Heading)
	detailCTA.Description = replacer.Replace(detailCTA.Description)
	detailCTA.PrimaryButtonUrl = replacer.Replace(detailCTA.PrimaryButtonUrl)
	detailCTA.SecondaryButtonUrl = replacer.Replace(detailCTA.SecondaryButtonUrl)

	sections, _ := h.queries.ListPageSections(ctx, "product_detail")
	sectionMap := make(map[string]sqlc.PageSection)
	for _, s := range sections {
		sectionMap[s.SectionKey] = s
	}

	metaTitle := ""
	if detail.Product.MetaTitle.Valid && detail.Product.MetaTitle.String != "" {
		metaTitle = detail.Product.MetaTitle.String
	}
	metaDesc := ""
	if detail.Product.MetaDescription.Valid && detail.Product.MetaDescription.String != "" {
		metaDesc = detail.Product.MetaDescription.String
	}

	data := map[string]interface{}{
		"Title":           fmt.Sprintf("%s | Products", detail.Product.Name),
		"MetaTitle":       metaTitle,
		"MetaDescription": metaDesc,
		"OGImage":         detail.Product.OgImage,
		"CanonicalURL":    fmt.Sprintf("/products/%s/%s", detail.Category.Slug, detail.Product.Slug),
		"Product":         detail.Product,
		"Category":        detail.Category,
		"Images":          detail.Images,
		"Features":        detail.Features,
		"SpecSections":    specSections,
		"Certifications":  detail.Certifications,
		"Downloads":       detail.Downloads,
		"DetailCTA":       detailCTA,
		"Sections":        sectionMap,
	}

	if preview {
		data["IsPreview"] = true
		data["EditURL"] = fmt.Sprintf("/admin/products/%d/edit", detail.Product.ID)
		return h.renderAndCache(c, "preview:product:"+productSlug, 0, http.StatusOK, "public/pages/product_detail.html", data)
	}

	cacheKey := fmt.Sprintf("page:products:%s:%s", categorySlug, productSlug)
	return h.renderAndCache(c, cacheKey, 1800, http.StatusOK, "public/pages/product_detail.html", data)
}

// GET /products/search?q=...
func (h *ProductsHandler) ProductSearch(c echo.Context) error {
	ctx := c.Request().Context()
	q := c.QueryParam("q")

	var products []sqlc.Product
	if q != "" {
		wildcard := "%" + q + "%"
		var err error
		products, err = h.queries.SearchProducts(ctx, sqlc.SearchProductsParams{
			Name:        wildcard,
			Description: wildcard,
			Tagline:     sql.NullString{String: wildcard, Valid: true},
			Limit:       24,
			Offset:      0,
		})
		if err != nil {
			h.logger.Error("failed to search products", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	data := map[string]interface{}{
		"Title":    fmt.Sprintf("Search: %s | Products", q),
		"Products": products,
		"Query":    q,
	}

	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "public/partials/product_search_results.html", data)
	}

	return c.Render(http.StatusOK, "public/pages/products.html", data)
}

func groupSpecsBySection(specs []sqlc.ProductSpec) map[string][]sqlc.ProductSpec {
	sections := make(map[string][]sqlc.ProductSpec)
	for _, spec := range specs {
		sections[spec.SectionName] = append(sections[spec.SectionName], spec)
	}
	return sections
}
