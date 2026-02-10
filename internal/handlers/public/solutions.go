// Package public provides HTTP handlers for the public-facing website.
// This file contains handlers for solutions listing and individual solution
// detail pages. Solutions represent industry-specific or use-case-specific
// offerings that combine multiple products and services.
package public

import (
	// Standard library imports
	"bytes"       // Buffer for template rendering to enable HTML caching
	"database/sql" // SQL error handling (sql.ErrNoRows for 404 detection)
	"fmt"         // String formatting for cache keys and template data
	"log/slog"    // Structured logging for debugging and error tracking
	"net/http"    // HTTP status codes and request/response handling
	"strings"     // String manipulation for placeholder replacement in sections

	// Third-party imports
	"github.com/labstack/echo/v4" // Echo web framework - routing, context, rendering

	// Internal imports
	"github.com/narendhupati/bluejay-cms/db/sqlc"          // sqlc-generated database queries
	"github.com/narendhupati/bluejay-cms/internal/services" // Cache service for HTML caching
)

// SolutionsHandler handles all solution-related public routes including
// the solutions overview page and individual solution detail pages.
// Solutions represent industry/use-case specific offerings (e.g., "Smart Factory",
// "Cold Chain Monitoring", "Energy Management"). It implements caching for
// improved performance on these content-heavy pages.
type SolutionsHandler struct {
	queries *sqlc.Queries    // Database query interface for solutions and related data
	logger  *slog.Logger     // Structured logger for errors and debugging
	cache   *services.Cache  // In-memory cache for rendered HTML pages
}

// NewSolutionsHandler creates a new SolutionsHandler with the required dependencies.
// The cache is used to store rendered HTML to reduce database queries and
// template rendering overhead. Solutions pages are relatively stable and
// benefit significantly from caching.
func NewSolutionsHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *SolutionsHandler {
	return &SolutionsHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

// renderAndCache is a helper method that renders a template to HTML,
// caches the result, and returns it to the client.
//
// Purpose:
// This method centralizes the render-and-cache pattern used across all
// solutions handlers. It ensures consistent injection of middleware data
// (settings, footer navigation) and caches the rendered HTML to reduce
// database load and template rendering overhead.
//
// Parameters:
//   - cacheKey: Unique identifier for the cached HTML (e.g., "page:solutions:smart-factory")
//   - ttlSeconds: Time-to-live for the cached HTML (0 means no expiration)
//   - statusCode: HTTP status code to return (typically http.StatusOK)
//   - templateName: Name of the template to render (e.g., "public/pages/solutions_list.html")
//   - data: Template data map to populate the template
//
// Caching Strategy:
//   - Solutions listing: 600 seconds (10 minutes) - moderate TTL
//   - Solution detail: 1800 seconds (30 minutes) - longer TTL for stable content
//   - Preview mode: No cache (TTL=0) - to show live changes
//
// Middleware Data Injection:
//   - Settings: Global site configuration
//   - FooterCategories: Product categories for footer navigation
//   - FooterSolutions: Solutions for footer navigation
//   - FooterResources: Resources for footer navigation
func (h *SolutionsHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
	// Inject middleware-provided data into template context
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

	// Render template to buffer for caching
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		h.logger.Error("template render failed", "template", templateName, "error", err)
		return err
	}

	// Cache the rendered HTML and return to client
	html := buf.String()
	h.cache.Set(cacheKey, html, ttlSeconds)
	return c.HTML(statusCode, html)
}

// SolutionsList handles GET requests to the solutions overview page.
//
// HTTP Method: GET
// Route: /solutions
// Template: public/pages/solutions_list.html (full page, not HTMX fragment)
// HTMX: Returns complete HTML page
// Cache TTL: 600 seconds (10 minutes)
//
// Purpose:
// Displays an overview of all published solutions with:
//   - Hero section with page introduction
//   - Grid of solution cards (each with icon, title, description, link)
//   - Solution page features (benefits of the solutions approach)
//   - Call-to-action section
//
// Solutions represent industry-specific or use-case-specific offerings that
// combine products and services. Examples:
//   - Smart Factory Automation
//   - Cold Chain Monitoring
//   - Energy Management Systems
//   - Predictive Maintenance
//
// Template Data:
//   - Title: "Solutions" - Browser tab title
//   - Solutions: []sqlc.Solution - All published solutions
//   - Features: []sqlc.SolutionPageFeature - Solution page features/benefits
//   - CTA: sqlc.SolutionCTA - Call-to-action section
//   - CurrentPage: "solutions" - For navigation highlighting
//   - PageHero: sqlc.PageSection - Hero section content
//   - GridSection: sqlc.PageSection - Solutions grid section heading
//   - FeaturesSection: sqlc.PageSection - Features section heading
//
// SEO Considerations:
//   - Solutions page is a key landing page for industry-specific searches
//   - Each solution links to detailed solution pages
//   - Page sections allow customizable headings for SEO optimization
//   - Features section highlights value propositions
//
// Page Structure:
//   1. Hero section - Introduction to solutions approach
//   2. Solutions grid - Cards for each solution
//   3. Features section - Benefits of solutions (e.g., "Tailored to Your Industry")
//   4. CTA section - Contact sales or request consultation
//
// Error Handling:
//   - Critical errors (solutions list) return HTTP 500
//   - Non-critical content (features, CTA) fail gracefully with empty arrays
//
// Caching:
//   - Cached for 10 minutes to reduce database load
//   - Cache invalidated when solutions are published/unpublished
func (h *SolutionsHandler) SolutionsList(c echo.Context) error {
	// Check cache for fast response on repeated requests
	cacheKey := "page:solutions"
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	// Fetch all published solutions (critical - must succeed)
	solutions, err := h.queries.ListPublishedSolutions(ctx)
	if err != nil {
		h.logger.Error("failed to list published solutions", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Fetch solution page features (benefits like "Industry Expertise", "Proven Results")
	// Non-critical - gracefully degrade to empty array on error
	features, err := h.queries.ListSolutionPageFeatures(ctx)
	if err != nil {
		h.logger.Error("failed to list solution page features", "error", err)
		features = []sqlc.SolutionPageFeature{} // Default to empty array
	}

	// Fetch active CTA for the solutions listing page
	// Non-critical - can be missing
	cta, err := h.queries.GetActiveSolutionsListingCTA(ctx)
	if err != nil && err != sql.ErrNoRows {
		// Log error only if it's not just a "no rows" case
		h.logger.Error("failed to get active solutions listing CTA", "error", err)
	}

	// Fetch editable page sections for admin customization
	// These allow the admin to customize headings/descriptions without code changes
	heroSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "solutions", SectionKey: "hero"})
	gridSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "solutions", SectionKey: "grid_section"})
	featuresSection, _ := h.queries.GetPageSection(ctx, sqlc.GetPageSectionParams{PageKey: "solutions", SectionKey: "features_section"})

	// Assemble template data
	data := map[string]interface{}{
		"Title":           "Solutions",      // Browser tab title
		"Solutions":       solutions,        // All published solutions
		"Features":        features,         // Solution page features
		"CTA":             cta,              // Call-to-action section
		"CurrentPage":     "solutions",      // For nav highlighting
		"PageHero":        heroSection,      // Hero section content
		"GridSection":     gridSection,      // Grid section heading
		"FeaturesSection": featuresSection,  // Features section heading
	}

	// Render template and cache for 10 minutes
	// Template: templates/public/pages/solutions_list.html
	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/solutions_list.html", data)
}

// SolutionDetail handles GET requests to view a specific solution's detail page.
//
// HTTP Method: GET
// Route: /solutions/:slug (e.g., /solutions/smart-factory-automation)
// Query Parameters: ?preview=1 (optional, for admin preview)
// Template: public/pages/solution_detail.html (full page, not HTMX fragment)
// HTMX: Returns complete HTML page
// Cache TTL: 1800 seconds (30 minutes) for published solutions, no cache for previews
//
// Purpose:
// Displays comprehensive details for a single solution including:
//   - Solution overview (title, description, hero image)
//   - Statistics/metrics (e.g., "30% cost reduction", "500+ installations")
//   - Challenges addressed (problems this solution solves)
//   - Featured products (products included in this solution)
//   - Call-to-action sections (request demo, contact sales, download whitepaper)
//   - Other solutions (cross-linking for discovery)
//   - SEO metadata (custom title, description, OG image)
//
// URL Parameters:
//   - slug: Unique slug identifier for the solution (e.g., "cold-chain-monitoring")
//
// Query Parameters:
//   - preview: If "1", shows draft/unpublished solutions for admin preview
//
// Template Data:
//   - Title: Solution title - Browser tab title
//   - MetaTitle: Custom SEO title (if set)
//   - MetaDescription: SEO meta description
//   - MetaDesc: Duplicate of MetaDescription (for template compatibility)
//   - OGImage: Open Graph image for social sharing
//   - CanonicalURL: Canonical URL for SEO
//   - Solution: sqlc.Solution - Core solution data (title, description, content)
//   - Stats: []sqlc.SolutionStat - Metrics/statistics for this solution
//   - Challenges: []sqlc.SolutionChallenge - Problems/challenges addressed
//   - Products: []sqlc.Product - Products included in this solution
//   - CTAs: []sqlc.SolutionCTA - Call-to-action sections
//   - OtherSolutions: []sqlc.Solution - Other published solutions (for cross-linking)
//   - CurrentPage: "solutions" - For navigation highlighting
//   - Sections: map[string]sqlc.PageSection - Editable sections with placeholders replaced
//   - IsPreview: bool - True if viewing in preview mode (only in preview)
//   - EditURL: string - Admin edit URL (only in preview mode)
//
// SEO Considerations:
//   - Solution pages are key landing pages for industry/use-case searches
//   - Custom meta titles and descriptions improve search visibility
//   - OG images enhance social media sharing
//   - Canonical URLs prevent duplicate content issues
//   - Cross-linking to other solutions and products improves internal linking
//
// Preview Mode:
//   - When ?preview=1 is set, shows unpublished/draft solutions
//   - Bypasses published status check for admin review
//   - Adds edit link for admin convenience
//   - Disables caching to show live changes immediately
//
// Dynamic Placeholders:
//   - {solution_title} → Replaced with actual solution title in sections
//   - Example: "Learn more about {solution_title}" → "Learn more about Smart Factory"
//   - This allows generic section templates to be personalized per solution
//
// Page Structure:
//   1. Hero section - Solution title, tagline, hero image
//   2. Overview - Detailed description of the solution
//   3. Statistics - Quantifiable results/metrics
//   4. Challenges - Problems this solution addresses
//   5. Products - Featured products in this solution
//   6. CTAs - Multiple calls-to-action (demo, contact, download)
//   7. Other solutions - Cross-linking for discovery
//
// Error Handling:
//   - Returns 404 if solution slug doesn't exist
//   - Returns 500 on database errors
//   - Gracefully handles missing stats/challenges/products/CTAs with empty arrays
//
// Caching:
//   - Published solutions cached for 30 minutes (longer than listing)
//   - Preview mode has no cache (TTL=0) for immediate updates
//   - Cache is invalidated when solution is updated in admin
func (h *SolutionsHandler) SolutionDetail(c echo.Context) error {
	slug := c.Param("slug")
	preview := isPreviewRequest(c) // Check if this is an admin preview request

	// Skip cache lookup for preview mode to show live changes
	if !preview {
		cacheKey := fmt.Sprintf("page:solutions:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	ctx := c.Request().Context()

	// Fetch solution by slug (different query for preview vs published)
	var solution sqlc.Solution
	var err error
	if preview {
		// Preview mode: include draft/unpublished solutions
		solution, err = h.queries.GetSolutionBySlugIncludeDrafts(ctx, slug)
	} else {
		// Normal mode: only show published solutions
		solution, err = h.queries.GetSolutionBySlug(ctx, slug)
	}
	if err == sql.ErrNoRows {
		// Solution doesn't exist or is not published
		return echo.NewHTTPError(http.StatusNotFound, "Solution not found")
	}
	if err != nil {
		h.logger.Error("failed to load solution", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Fetch associated data for this solution
	// All are non-critical - gracefully degrade to empty arrays on error

	// Statistics/metrics (e.g., "30% cost reduction", "99.9% uptime")
	stats, err := h.queries.GetSolutionStats(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution stats", "error", err)
		stats = []sqlc.SolutionStat{}
	}

	// Challenges addressed by this solution
	challenges, err := h.queries.GetSolutionChallenges(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution challenges", "error", err)
		challenges = []sqlc.SolutionChallenge{}
	}

	// Products included in this solution
	products, err := h.queries.GetSolutionProducts(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution products", "error", err)
		products = []sqlc.GetSolutionProductsRow{}
	}

	// Call-to-action sections (request demo, download whitepaper, etc.)
	ctas, err := h.queries.GetSolutionCTAs(ctx, solution.ID)
	if err != nil {
		h.logger.Error("failed to load solution CTAs", "error", err)
		ctas = []sqlc.SolutionCta{}
	}

	// Fetch all published solutions for "Other Solutions" section
	allSolutions, err := h.queries.ListPublishedSolutions(ctx)
	if err != nil {
		h.logger.Error("failed to list published solutions", "error", err)
		allSolutions = []sqlc.ListPublishedSolutionsRow{}
	}

	// Filter out the current solution from the list
	// This creates a cross-linking opportunity for users to explore other solutions
	var otherSolutions []sqlc.ListPublishedSolutionsRow
	for _, s := range allSolutions {
		if s.ID != solution.ID {
			otherSolutions = append(otherSolutions, s)
		}
	}

	// Extract meta description with null-safety
	metaDesc := ""
	if solution.MetaDescription.Valid {
		metaDesc = solution.MetaDescription.String
	}

	// Fetch editable page sections and replace placeholders
	sections, _ := h.queries.ListPageSections(ctx, "solution_detail")
	replacer := strings.NewReplacer("{solution_title}", solution.Title)
	sectionMap := make(map[string]sqlc.PageSection)
	for _, s := range sections {
		// Replace placeholders in section content with actual solution data
		s.Heading = replacer.Replace(s.Heading)
		s.Description = replacer.Replace(s.Description)
		sectionMap[s.SectionKey] = s
	}

	// Assemble template data
	data := map[string]interface{}{
		"Title":           solution.Title,                          // Browser tab title
		"MetaTitle":       solution.MetaTitle,                      // SEO title
		"MetaDescription": metaDesc,                                // SEO description
		"MetaDesc":        metaDesc,                                // Duplicate for template compatibility
		"OGImage":         solution.OgImage,                        // Social sharing image
		"CanonicalURL":    fmt.Sprintf("/solutions/%s", solution.Slug), // SEO canonical URL
		"Solution":        solution,                                // Core solution data
		"Stats":           stats,                                   // Statistics/metrics
		"Challenges":      challenges,                              // Challenges addressed
		"Products":        products,                                // Featured products
		"CTAs":            ctas,                                    // Call-to-action sections
		"OtherSolutions":  otherSolutions,                          // Other solutions for cross-linking
		"CurrentPage":     "solutions",                             // For nav highlighting
		"Sections":        sectionMap,                              // Editable sections
	}

	// Handle preview mode
	if preview {
		data["IsPreview"] = true // Show preview banner in template
		data["EditURL"] = fmt.Sprintf("/admin/solutions/%d/edit", solution.ID) // Link to admin editor
		// Don't cache preview pages (TTL=0)
		return h.renderAndCache(c, "preview:solution:"+slug, 0, http.StatusOK, "public/pages/solution_detail.html", data)
	}

	// Render and cache for 30 minutes (1800 seconds)
	// Template: templates/public/pages/solution_detail.html
	return h.renderAndCache(c, fmt.Sprintf("page:solutions:%s", slug), 1800, http.StatusOK, "public/pages/solution_detail.html", data)
}
