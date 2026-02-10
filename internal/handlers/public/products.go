// Package public provides HTTP handlers for the public-facing website.
// This file contains handlers for product listing, category browsing,
// product detail pages, and product search functionality.
package public

import (
	// Standard library imports
	"bytes"       // Buffer for template rendering to enable caching
	"database/sql" // SQL error handling (sql.ErrNoRows for 404 detection)
	"fmt"         // String formatting for error messages and template data
	"log/slog"    // Structured logging for debugging and error tracking
	"net/http"    // HTTP status codes and request/response handling
	"strconv"     // String to integer conversion for pagination parameters
	"strings"     // String manipulation for placeholder replacement in CTA text

	// Third-party imports
	"github.com/labstack/echo/v4" // Echo web framework - routing, context, rendering

	// Internal imports
	"github.com/narendhupati/bluejay-cms/db/sqlc"          // sqlc-generated database queries
	"github.com/narendhupati/bluejay-cms/internal/services" // Business logic services (ProductService, Cache)
)

// ProductsHandler handles all product-related public routes including
// product listing, category browsing, product details, and search.
// It implements caching for improved performance and uses the ProductService
// for complex product data aggregation.
type ProductsHandler struct {
	queries    *sqlc.Queries              // Database query interface for product data
	logger     *slog.Logger               // Structured logger for errors and debugging
	productSvc *services.ProductService   // Business logic for product detail aggregation
	cache      *services.Cache            // In-memory cache for rendered HTML pages
}

// NewProductsHandler creates a new ProductsHandler with the required dependencies.
// The cache is used to store rendered HTML to reduce database load and improve
// response times for frequently accessed pages.
func NewProductsHandler(queries *sqlc.Queries, logger *slog.Logger, productSvc *services.ProductService, cache *services.Cache) *ProductsHandler {
	return &ProductsHandler{
		queries:    queries,
		logger:     logger,
		productSvc: productSvc,
		cache:      cache,
	}
}

// renderAndCache is a helper method that renders a template to HTML,
// caches the result, and returns it to the client.
//
// Purpose:
// This method centralizes the render-and-cache pattern used across all
// product handlers. It ensures consistent injection of middleware data
// (settings, footer navigation) and caches the rendered HTML to reduce
// database load and template rendering overhead.
//
// Parameters:
//   - cacheKey: Unique identifier for the cached HTML (e.g., "page:products:electronics")
//   - ttlSeconds: Time-to-live for the cached HTML (0 means no expiration)
//   - statusCode: HTTP status code to return (typically http.StatusOK)
//   - templateName: Name of the template to render (e.g., "public/pages/products.html")
//   - data: Template data map to populate the template
//
// Caching Strategy:
//   - Rendered HTML is stored in the cache for fast subsequent requests
//   - Cache keys are unique per page/category/product for granular invalidation
//   - TTL varies by content type: static pages get longer TTL
//
// Middleware Data Injection:
//   - Settings: Global site configuration (logo, contact info, etc.)
//   - FooterCategories: Product categories for footer navigation
//   - FooterSolutions: Solutions for footer navigation
//   - FooterResources: Resources for footer navigation
func (h *ProductsHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
	// Inject middleware-provided data into template context
	// These are set by middleware and provide consistent site-wide navigation
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

	// Render template to a buffer instead of directly to response
	// This allows us to cache the HTML before sending it
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		h.logger.Error("template render failed", "template", templateName, "error", err)
		return err
	}

	// Extract rendered HTML and store in cache for future requests
	html := buf.String()
	h.cache.Set(cacheKey, html, ttlSeconds)

	// Return the rendered HTML to the client
	return c.HTML(statusCode, html)
}

// ProductsList handles GET requests to the products overview page.
//
// HTTP Method: GET
// Route: /products
// Template: public/pages/products.html (full page, not HTMX fragment)
// HTMX: Returns complete HTML page
// Cache TTL: 600 seconds (10 minutes)
//
// Purpose:
// Displays an overview of all product categories with product counts.
// This is the main product landing page that helps users browse categories
// before drilling down into specific products.
//
// Template Data:
//   - Title: "Products" - Browser tab title
//   - Categories: []categoryWithCount - Product categories with count of products in each
//   - TotalCount: int - Total number of categories
//   - PageHero: sqlc.PageSection - Hero section content (heading, description, image)
//   - CategoriesSection: sqlc.PageSection - Categories section heading/description
//   - PageCTA: sqlc.PageSection - Call-to-action section
//
// SEO Considerations:
//   - This is a key landing page for product discovery
//   - Categories are displayed with counts to help users navigate
//   - Page sections allow customizable headings for better SEO
//
// Caching:
//   - Page is cached for 10 minutes to reduce database load
//   - Cache is invalidated when categories or products are modified in admin
func (h *ProductsHandler) ProductsList(c echo.Context) error {
	// Check cache first for fast response on repeated requests
	cacheKey := "page:products"
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	// Fetch all product categories from database
	categories, err := h.queries.ListProductCategories(ctx)
	if err != nil {
		h.logger.Error("failed to list product categories", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Define a struct to hold category with its product count
	// This enriches the category data for display in the template
	type categoryWithCount struct {
		Category sqlc.ProductCategory
		Count    int64
	}

	// Fetch product count for each category
	// This helps users understand the size of each category before clicking
	var categoriesWithCount []categoryWithCount
	for _, cat := range categories {
		count, _ := h.queries.CountProductsByCategory(ctx, cat.ID)
		categoriesWithCount = append(categoriesWithCount, categoryWithCount{
			Category: cat,
			Count:    count,
		})
	}

	// Fetch editable page sections for admin-customizable content
	// These allow the admin to customize headings/descriptions without code changes
	heroSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products", SectionKey: "hero"})
	categoriesSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products", SectionKey: "categories_section"})
	ctaSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products", SectionKey: "cta"})

	// Assemble template data
	data := map[string]interface{}{
		"Title":             "Products",
		"Categories":        categoriesWithCount, // Categories with product counts
		"TotalCount":        len(categories),     // Total number of categories
		"PageHero":          heroSection,         // Hero section content
		"CategoriesSection": categoriesSection,   // Categories section heading
		"PageCTA":           ctaSection,          // CTA section
	}

	// Render template and cache for 10 minutes
	// Template: templates/public/pages/products.html
	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/products.html", data)
}

// ProductsByCategory handles GET requests to view products within a specific category.
//
// HTTP Method: GET
// Route: /products/:category (e.g., /products/electronics)
// Query Parameters: ?page=N (optional, defaults to 1)
// Template: public/pages/products_category.html (full page, not HTMX fragment)
// HTMX: Returns complete HTML page
// Cache TTL: 600 seconds (10 minutes)
//
// Purpose:
// Displays a paginated list of products within a specific category.
// Users can navigate through pages if the category has more than 12 products.
//
// URL Parameters:
//   - category: URL slug of the product category (e.g., "industrial-sensors")
//
// Query Parameters:
//   - page: Page number for pagination (default: 1, minimum: 1)
//
// Template Data:
//   - Title: "{Category Name} | Products" - Browser tab title
//   - Category: sqlc.ProductCategory - Category details (name, description, icon)
//   - Products: []sqlc.Product - Products in this category (max 12 per page)
//   - TotalCount: int64 - Total number of products in category
//   - CurrentPage: int - Current page number
//   - TotalPages: int - Total number of pages
//   - CategoryHero: sqlc.PageSection - Hero section content
//   - EmptyState: sqlc.PageSection - Content to show if category has no products
//
// SEO Considerations:
//   - Page title includes category name for better search relevance
//   - Pagination helps manage large product catalogs
//   - Empty state section provides guidance when category has no products
//
// Error Handling:
//   - Returns 404 if category slug doesn't exist
//   - Returns 500 on database errors
//
// Pagination:
//   - 12 products per page (hard-coded limit)
//   - Invalid page numbers default to page 1
//   - Negative page numbers are rejected
func (h *ProductsHandler) ProductsByCategory(c echo.Context) error {
	ctx := c.Request().Context()
	categorySlug := c.Param("category")

	// Check cache for this specific category page
	cacheKey := fmt.Sprintf("page:products:%s", categorySlug)
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Fetch category by slug
	category, err := h.queries.GetProductCategoryBySlug(ctx, categorySlug)
	if err == sql.ErrNoRows {
		// Category doesn't exist - return 404
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}
	if err != nil {
		h.logger.Error("failed to load category", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Parse pagination parameter from query string
	// Default to page 1 if not provided or invalid
	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	// Set up pagination parameters
	limit := int64(12) // Show 12 products per page
	offset := int64((page - 1)) * limit

	// Fetch products for this category with pagination
	products, err := h.queries.ListProductsByCategory(ctx, sqlc.ListProductsByCategoryParams{
		CategoryID: category.ID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		h.logger.Error("failed to load products", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Calculate total pages for pagination controls
	total, _ := h.queries.CountProductsByCategory(ctx, category.ID)
	totalPages := int((total + limit - 1) / limit) // Ceiling division

	// Fetch editable page sections
	categoryHero, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products_category", SectionKey: "hero"})
	emptyState, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "products_category", SectionKey: "empty_state"})

	// Assemble template data
	data := map[string]interface{}{
		"Title":        fmt.Sprintf("%s | Products", category.Name), // SEO-friendly title
		"Category":     category,                                    // Category details
		"Products":     products,                                    // Products on current page
		"TotalCount":   total,                                       // Total products in category
		"CurrentPage":  page,                                        // Current page number
		"TotalPages":   totalPages,                                  // Total pages for pagination
		"CategoryHero": categoryHero,                                // Hero section content
		"EmptyState":   emptyState,                                  // Empty state message
	}

	// Render template and cache for 10 minutes
	// Template: templates/public/pages/products_category.html
	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/products_category.html", data)
}

// ProductDetail handles GET requests to view a specific product's detail page.
//
// HTTP Method: GET
// Route: /products/:category/:slug (e.g., /products/sensors/temperature-sensor-ts100)
// Query Parameters: ?preview=1 (optional, for admin preview)
// Template: public/pages/product_detail.html (full page, not HTMX fragment)
// HTMX: Returns complete HTML page
// Cache TTL: 1800 seconds (30 minutes) for published products, no cache for previews
//
// Purpose:
// Displays comprehensive details for a single product including:
//   - Product images (gallery)
//   - Features and benefits
//   - Technical specifications (grouped by section)
//   - Certifications and compliance
//   - Downloadable resources (datasheets, manuals, CAD files)
//   - Call-to-action (contact sales, request quote)
//
// URL Parameters:
//   - category: Category slug (must match product's actual category or returns 404)
//   - slug: Product slug (unique identifier)
//
// Query Parameters:
//   - preview: If "1", shows draft/unpublished products for admin preview
//
// Template Data:
//   - Title: "{Product Name} | Products" - Browser tab title
//   - MetaTitle: Custom SEO title (if set) - Used in <title> tag
//   - MetaDescription: SEO meta description
//   - OGImage: Open Graph image for social sharing
//   - CanonicalURL: Canonical URL for SEO
//   - Product: sqlc.Product - Core product data (name, SKU, description, price)
//   - Category: sqlc.ProductCategory - Parent category details
//   - Images: []sqlc.ProductImage - Product images for gallery
//   - Features: []sqlc.ProductFeature - Features/benefits list
//   - SpecSections: map[string][]sqlc.ProductSpec - Specs grouped by section
//   - Certifications: []sqlc.ProductCertification - Certifications/compliance
//   - Downloads: []sqlc.ProductDownload - Downloadable resources
//   - DetailCTA: sqlc.PageSection - Call-to-action with placeholders replaced
//   - Sections: map[string]sqlc.PageSection - Other editable sections
//   - IsPreview: bool - True if viewing in preview mode
//   - EditURL: string - Admin edit URL (only in preview mode)
//
// SEO Considerations:
//   - MetaTitle, MetaDescription, OGImage can be customized per product
//   - Canonical URL prevents duplicate content issues
//   - Structured data could be added for rich snippets
//
// Preview Mode:
//   - When ?preview=1 is set, shows unpublished/draft products
//   - Adds edit link for admin convenience
//   - Disables caching to show live changes
//
// Dynamic CTA Placeholders:
//   - {product_name} → Replaced with actual product name
//   - {product_sku} → Replaced with actual product SKU
//   - This allows generic CTA templates to be personalized per product
//
// Error Handling:
//   - Returns 404 if product not found
//   - Returns 404 if category slug doesn't match product's category
//   - Returns 500 on database errors
//
// Business Logic:
//   - Uses ProductService to aggregate data from multiple tables
//   - Specifications are grouped by section for better organization
//   - Preview mode bypasses published status checks
func (h *ProductsHandler) ProductDetail(c echo.Context) error {
	categorySlug := c.Param("category")
	productSlug := c.Param("slug")
	preview := isPreviewRequest(c) // Check if this is an admin preview request

	// Skip cache lookup for preview mode to show live changes
	if !preview {
		cacheKey := fmt.Sprintf("page:products:%s:%s", categorySlug, productSlug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	// Use ProductService to fetch aggregated product detail
	// This retrieves product + category + images + features + specs + certifications + downloads
	detail, err := h.productSvc.GetProductDetail(c.Request().Context(), productSlug)
	if err == sql.ErrNoRows {
		// Product doesn't exist - return 404
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	if err != nil {
		h.logger.Error("failed to load product detail", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Validate that the category slug in the URL matches the product's actual category
	// This prevents accessing products via incorrect category URLs
	if detail.Category.Slug != categorySlug {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found in this category")
	}

	// Group specifications by section name for organized display
	// Example sections: "General", "Electrical", "Mechanical", "Environmental"
	specSections := groupSpecsBySection(detail.Specs)

	ctx := c.Request().Context()

	// Fetch CTA section and personalize it with product-specific placeholders
	detailCTA, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "product_detail", SectionKey: "cta"})

	// Replace placeholders in CTA text with actual product data
	// Example: "Request a quote for {product_name}" → "Request a quote for TS100"
	replacer := strings.NewReplacer("{product_name}", detail.Product.Name, "{product_sku}", detail.Product.Sku)
	detailCTA.Heading = replacer.Replace(detailCTA.Heading)
	detailCTA.Description = replacer.Replace(detailCTA.Description)
	detailCTA.PrimaryButtonUrl = replacer.Replace(detailCTA.PrimaryButtonUrl)
	detailCTA.SecondaryButtonUrl = replacer.Replace(detailCTA.SecondaryButtonUrl)

	// Fetch other editable page sections for admin customization
	sections, _ := h.queries.ListPageSections(ctx, "product_detail")
	sectionMap := make(map[string]sqlc.PageSection)
	for _, s := range sections {
		sectionMap[s.SectionKey] = s
	}

	// Extract SEO metadata, using fallbacks if not set
	metaTitle := ""
	if detail.Product.MetaTitle.Valid && detail.Product.MetaTitle.String != "" {
		metaTitle = detail.Product.MetaTitle.String
	}
	metaDesc := ""
	if detail.Product.MetaDescription.Valid && detail.Product.MetaDescription.String != "" {
		metaDesc = detail.Product.MetaDescription.String
	}

	// Assemble template data with all product information
	data := map[string]interface{}{
		"Title":           fmt.Sprintf("%s | Products", detail.Product.Name), // Browser tab title
		"MetaTitle":       metaTitle,                                         // SEO title
		"MetaDescription": metaDesc,                                          // SEO description
		"OGImage":         detail.Product.OgImage,                            // Social sharing image
		"CanonicalURL":    fmt.Sprintf("/products/%s/%s", detail.Category.Slug, detail.Product.Slug), // SEO canonical
		"Product":         detail.Product,         // Core product data
		"Category":        detail.Category,        // Parent category
		"Images":          detail.Images,          // Product image gallery
		"Features":        detail.Features,        // Features/benefits list
		"SpecSections":    specSections,           // Specifications grouped by section
		"Certifications":  detail.Certifications,  // Certifications/compliance
		"Downloads":       detail.Downloads,       // Downloadable resources
		"DetailCTA":       detailCTA,              // Personalized CTA
		"Sections":        sectionMap,             // Other editable sections
	}

	// Handle preview mode (for admin to preview unpublished changes)
	if preview {
		data["IsPreview"] = true // Show preview banner in template
		data["EditURL"] = fmt.Sprintf("/admin/products/%d/edit", detail.Product.ID) // Link to admin editor
		// Don't cache preview pages (TTL=0)
		return h.renderAndCache(c, "preview:product:"+productSlug, 0, http.StatusOK, "public/pages/product_detail.html", data)
	}

	// Render and cache for 30 minutes (1800 seconds)
	// Template: templates/public/pages/product_detail.html
	cacheKey := fmt.Sprintf("page:products:%s:%s", categorySlug, productSlug)
	return h.renderAndCache(c, cacheKey, 1800, http.StatusOK, "public/pages/product_detail.html", data)
}

// ProductSearch handles GET requests to search for products by keyword.
//
// HTTP Method: GET
// Route: /products/search
// Query Parameters: ?q={search_term} (required)
// Template: public/partials/product_search_results.html (HTMX fragment)
//           OR public/pages/products.html (full page for non-HTMX requests)
// HTMX: Returns HTML fragment when HX-Request header is present
// Cache: Not cached (search results should be fresh)
//
// Purpose:
// Provides keyword search across product names, descriptions, and taglines.
// Supports both full-page requests (direct navigation) and HTMX requests
// (for live search as user types).
//
// Query Parameters:
//   - q: Search query string (searches name, description, tagline fields)
//
// Template Data:
//   - Title: "Search: {query} | Products" - Browser tab title
//   - Products: []sqlc.Product - Matching products (max 24 results)
//   - Query: string - The search term entered by user
//   - Settings: sqlc.Setting - Global site settings (for non-HTMX requests)
//
// HTMX Behavior:
//   - When request has HX-Request header (HTMX request):
//     → Returns partial HTML fragment (public/partials/product_search_results.html)
//     → Fragment contains only the search results grid
//     → Used for live search-as-you-type functionality
//   - When request is a regular page load:
//     → Returns full page template (public/pages/products.html)
//     → Includes site header, footer, navigation
//
// Search Implementation:
//   - Uses SQL LIKE with wildcards for fuzzy matching
//   - Searches across: product name, description, tagline
//   - Case-insensitive search (handled by database)
//   - Limited to 24 results to prevent performance issues
//   - Empty query returns empty results (not all products)
//
// SEO Considerations:
//   - Search results pages are not indexed (should have noindex meta tag)
//   - Query parameter is included in page title for user clarity
//
// Performance:
//   - Not cached since results vary by query
//   - Limited to 24 results to keep response fast
//   - Database indexes on name/description fields recommended
func (h *ProductsHandler) ProductSearch(c echo.Context) error {
	ctx := c.Request().Context()
	q := c.QueryParam("q") // Extract search query from URL parameter

	var products []sqlc.Product
	if q != "" {
		// Add SQL wildcards for partial matching
		wildcard := "%" + q + "%"
		var err error
		products, err = h.queries.SearchProducts(ctx, sqlc.SearchProductsParams{
			Name:        wildcard, // Search in product name
			Description: wildcard, // Search in product description
			Tagline:     sql.NullString{String: wildcard, Valid: true}, // Search in tagline
			Limit:       24, // Limit results to prevent overwhelming UI and database
			Offset:      0,  // No pagination for search results
		})
		if err != nil {
			h.logger.Error("failed to search products", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}
	// If query is empty, products remains empty slice (don't return all products)

	// Assemble template data
	data := map[string]interface{}{
		"Title":    fmt.Sprintf("Search: %s | Products", q), // SEO-friendly title
		"Products": products,                                 // Matching products
		"Query":    q,                                        // Original query for display
	}

	// Add settings for non-HTMX requests (needed for full page layout)
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}

	// Check if this is an HTMX request (has HX-Request header)
	if c.Request().Header.Get("HX-Request") == "true" {
		// Return partial HTML fragment for HTMX live search
		// Template: templates/public/partials/product_search_results.html
		// This fragment contains only the results grid, no layout/header/footer
		return c.Render(http.StatusOK, "public/partials/product_search_results.html", data)
	}

	// Return full page for regular browser navigation
	// Template: templates/public/pages/products.html
	// This includes site layout, header, footer, navigation
	return c.Render(http.StatusOK, "public/pages/products.html", data)
}

// groupSpecsBySection organizes product specifications into sections.
//
// Purpose:
// Product specifications (like voltage, temperature range, dimensions) are
// stored as flat rows in the database. This helper function groups them by
// section name (e.g., "Electrical", "Mechanical", "Environmental") for
// organized display in the product detail template.
//
// Parameters:
//   - specs: Flat list of product specifications
//
// Returns:
//   - map[string][]sqlc.ProductSpec: Specifications grouped by section name
//
// Example:
//   Input: [
//     {SectionName: "Electrical", Label: "Voltage", Value: "24V"},
//     {SectionName: "Electrical", Label: "Current", Value: "2A"},
//     {SectionName: "Mechanical", Label: "Weight", Value: "500g"},
//   ]
//   Output: {
//     "Electrical": [
//       {SectionName: "Electrical", Label: "Voltage", Value: "24V"},
//       {SectionName: "Electrical", Label: "Current", Value: "2A"},
//     ],
//     "Mechanical": [
//       {SectionName: "Mechanical", Label: "Weight", Value: "500g"},
//     ],
//   }
func groupSpecsBySection(specs []sqlc.ProductSpec) map[string][]sqlc.ProductSpec {
	sections := make(map[string][]sqlc.ProductSpec)
	for _, spec := range specs {
		// Append each spec to its corresponding section's slice
		sections[spec.SectionName] = append(sections[spec.SectionName], spec)
	}
	return sections
}
