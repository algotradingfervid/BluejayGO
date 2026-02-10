// Package public provides HTTP handlers for the public-facing website.
// These handlers serve the customer-facing pages including home, products,
// blog, and solutions. All handlers in this package are read-only and do not
// require authentication.
package public

import (
	// Standard library imports
	"log/slog" // Structured logging for handler errors and debug info
	"net/http" // HTTP status codes and request/response handling

	// Third-party imports
	"github.com/labstack/echo/v4" // Echo web framework - handles routing, context, rendering

	// Internal imports
	"github.com/narendhupati/bluejay-cms/db/sqlc" // sqlc-generated database queries
)

// HomeHandler handles requests for the homepage route.
// It aggregates data from multiple content tables (hero, stats, testimonials,
// products, solutions, partners, blog posts) to compose the homepage view.
type HomeHandler struct {
	queries *sqlc.Queries // Database query interface for fetching homepage content
	logger  *slog.Logger  // Structured logger for error tracking
}

// NewHomeHandler creates a new HomeHandler instance with the provided dependencies.
// This constructor is used during application initialization to wire up the handler
// with the database queries and logger.
func NewHomeHandler(queries *sqlc.Queries, logger *slog.Logger) *HomeHandler {
	return &HomeHandler{
		queries: queries,
		logger:  logger,
	}
}

// ShowHomePage handles GET requests to the homepage route (/).
//
// HTTP Method: GET
// Route: /
// Template: public/pages/home.html (full page template, not an HTMX fragment)
// HTMX: This endpoint always returns a complete HTML page, not a fragment
//
// Purpose:
// Renders the main homepage by aggregating content from multiple database tables:
//   - Site settings (site name, logo, contact info)
//   - Hero section (main headline, subheading, CTA buttons, background image)
//   - Statistics section (numerical metrics like "500+ Customers")
//   - Testimonials (customer quotes and reviews)
//   - Call-to-action banner
//   - Featured products (up to 6 products marked as featured)
//   - Published solutions (industry/use-case specific solutions)
//   - Featured partners (logos/links of partner companies)
//   - Latest blog posts (3 most recent published posts)
//   - Page sections (editable headings/labels for each homepage section)
//
// SEO Considerations:
//   - Page title is set to site name from settings
//   - This is the primary landing page, so it should load quickly
//   - Footer categories, solutions, and resources are injected via middleware
//     for consistent site-wide navigation and internal linking
//
// Error Handling:
//   - Critical errors (like missing settings) return HTTP 500
//   - Non-critical content (hero, stats, etc.) silently fail to allow graceful degradation
//   - Missing sections will simply not appear on the page
//
// Template Data Structure:
//   - Title: string - Browser tab title
//   - Settings: sqlc.Setting - Global site configuration
//   - Hero: sqlc.Hero - Main hero section content
//   - Stats: []sqlc.Stat - List of statistics for stats section
//   - Testimonials: []sqlc.Testimonial - Customer testimonials
//   - CTA: sqlc.CTA - Call-to-action banner
//   - FeaturedProducts: []sqlc.Product - Featured products
//   - Solutions: []sqlc.Solution - Published solutions
//   - FeaturedPartners: []sqlc.Partner - Featured partner logos
//   - LatestPosts: []sqlc.BlogPost - Recent blog posts
//   - Sections: map[string]sqlc.PageSection - Editable section content by key
//   - FooterCategories: []sqlc.ProductCategory - For footer navigation
//   - FooterSolutions: []sqlc.Solution - For footer navigation
//   - FooterResources: []sqlc.Resource - For footer navigation
func (h *HomeHandler) ShowHomePage(c echo.Context) error {
	// Extract request context for database operations
	ctx := c.Request().Context()

	// Fetch global site settings - this is critical and must succeed
	settings, err := h.queries.GetSettings(ctx)
	if err != nil {
		h.logger.Error("failed to get settings", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Homepage-specific tables
	// These are non-critical - errors are ignored to allow graceful degradation
	hero, _ := h.queries.GetActiveHero(ctx)                          // Main hero section with headline and CTA
	stats, _ := h.queries.ListActiveStats(ctx)                       // Statistics section (e.g., "500+ Customers")
	testimonials, _ := h.queries.ListActiveTestimonialsHomepage(ctx) // Customer testimonials for homepage
	cta, _ := h.queries.GetActiveCTA(ctx)                            // Call-to-action banner

	// Existing content tables
	// Fetch related content to create a comprehensive homepage
	featuredProducts, _ := h.queries.ListFeaturedProducts(ctx, 6)        // Up to 6 featured products
	solutions, _ := h.queries.ListPublishedSolutions(ctx)                // All published solutions
	featuredPartners, _ := h.queries.ListFeaturedPartners(ctx, 12)       // Up to 12 partner logos
	latestPosts, _ := h.queries.ListLatestPublishedPosts(ctx, 3)         // 3 most recent blog posts

	// Page sections for editable labels/headings
	// Allows admin to customize section headings without code changes
	sections, _ := h.queries.ListPageSections(ctx, "home")
	sectionMap := make(map[string]sqlc.PageSection)
	for _, s := range sections {
		// Build a map keyed by section_key for easy template lookup
		sectionMap[s.SectionKey] = s
	}

	// Assemble template data with all homepage content
	data := map[string]interface{}{
		"Title":            settings.SiteName, // Browser tab title
		"Settings":         settings,          // Global site settings
		"Hero":             hero,              // Hero section content
		"Stats":            stats,             // Statistics section
		"Testimonials":     testimonials,      // Customer testimonials
		"CTA":              cta,               // Call-to-action banner
		"FeaturedProducts": featuredProducts,  // Featured products grid
		"Solutions":        solutions,         // Solutions section
		"FeaturedPartners": featuredPartners,  // Partner logos
		"LatestPosts":      latestPosts,       // Recent blog posts
		"Sections":         sectionMap,        // Editable section content
	}

	// Inject footer navigation data set by middleware
	// These provide consistent site-wide navigation links in the footer
	if cats := c.Get("footer_categories"); cats != nil {
		data["FooterCategories"] = cats // Product categories for footer nav
	}
	if sols := c.Get("footer_solutions"); sols != nil {
		data["FooterSolutions"] = sols // Solutions for footer nav
	}
	if res := c.Get("footer_resources"); res != nil {
		data["FooterResources"] = res // Resources for footer nav
	}

	// Render the full homepage template with aggregated data
	// Template: templates/public/pages/home.html
	// Layout: templates/public/layouts/public-layout.html
	return c.Render(http.StatusOK, "public/pages/home.html", data)
}
