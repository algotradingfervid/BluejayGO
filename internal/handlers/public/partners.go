// Package public provides HTTP handlers for public-facing pages of the Bluejay CMS.
// These handlers serve the front-end website content to visitors (non-admin users).
package public

import (
	// bytes provides buffer operations for building HTML output before sending to client
	"bytes"
	// log/slog is the structured logging library used for debug and error logging
	"log/slog"
	// net/http provides HTTP constants and status codes
	"net/http"

	// echo is the web framework used for routing and request/response handling
	"github.com/labstack/echo/v4"
	// sqlc provides type-safe database query interfaces generated from SQL files
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	// services provides business logic components like caching
	"github.com/narendhupati/bluejay-cms/internal/services"
)

// PartnersHandler handles HTTP requests for the Partners page.
// It manages displaying partner companies organized by tier (e.g., Platinum, Gold, Silver)
// and partner testimonials to build credibility.
type PartnersHandler struct {
	queries *sqlc.Queries    // Database query interface for fetching partners, tiers, and testimonials
	logger  *slog.Logger     // Structured logger for debugging and error tracking
	cache   *services.Cache  // In-memory cache for rendered HTML to improve response times
}

// NewPartnersHandler constructs a new PartnersHandler with required dependencies.
// This constructor is called during application initialization to wire up dependencies.
func NewPartnersHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *PartnersHandler {
	return &PartnersHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

// renderAndCache is a utility method that renders a template to HTML and caches the result.
// It injects global data (settings, footer navigation) from middleware into the template data,
// then renders the template to a buffer, stores the HTML in cache, and returns it to the client.
//
// Parameters:
//   - cacheKey: Unique key for storing this rendered HTML in cache
//   - ttlSeconds: Time-to-live in seconds; cache expires after this duration
//   - statusCode: HTTP status code to return (typically 200 OK)
//   - templateName: Path to the Go html/template file to render
//   - data: Template variables to pass to the template
//
// Returns HTML response with the specified status code, or error if rendering fails.
func (h *PartnersHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
	// Inject global settings from middleware context into template data
	// Settings include site name, contact info, social media links, etc.
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}
	// Inject footer navigation data from middleware
	// These are populated by middleware that runs before handlers
	if cats := c.Get("footer_categories"); cats != nil {
		data["FooterCategories"] = cats
	}
	if sols := c.Get("footer_solutions"); sols != nil {
		data["FooterSolutions"] = sols
	}
	if res := c.Get("footer_resources"); res != nil {
		data["FooterResources"] = res
	}

	// Render template to an in-memory buffer instead of directly to the response
	// This allows us to cache the rendered HTML before sending it
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		h.logger.Error("template render failed", "template", templateName, "error", err)
		return err
	}

	// Extract rendered HTML from buffer and store in cache for future requests
	html := buf.String()
	h.cache.Set(cacheKey, html, ttlSeconds)

	// Return the rendered HTML to the client with specified status code
	return c.HTML(statusCode, html)
}

// PartnersPage handles GET requests to /partners
// Renders the Partners page with partner companies grouped by tier level and partner testimonials.
//
// Route: GET /partners
// Template: templates/public/pages/partners.html (full page, not HTMX fragment)
// Cache: 300 seconds (5 minutes) - partner data changes occasionally
//
// HTMX Behavior: This endpoint returns a full HTML page, not an HTMX fragment.
// It is designed for direct browser navigation, not HTMX swaps.
//
// Business Logic: Partners are organized by tier (e.g., Platinum, Gold, Silver, Bronze).
// Each tier has different benefits and is displayed in priority order on the page.
//
// Returns: HTTP 200 with rendered partners.html template
func (h *PartnersHandler) PartnersPage(c echo.Context) error {
	// Define cache key for this page
	cacheKey := "page:partners"

	// Check if cached version exists and return it immediately to improve performance
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Fetch all partners sorted by tier display order
	// Query returns partners with tier information joined, ordered by tier priority
	allPartners, err := h.queries.ListPartnersByTier(ctx)
	if err != nil {
		h.logger.Error("failed to load partners", "error", err)
		allPartners = []sqlc.ListPartnersByTierRow{} // Default to empty slice to prevent template errors
	}

	// Fetch partner tier definitions (Platinum, Gold, Silver, etc.)
	// Each tier has a name, description, display order, and optional color/icon
	tiers, err := h.queries.ListPartnerTiers(ctx)
	if err != nil {
		h.logger.Error("failed to load partner tiers", "error", err)
		tiers = []sqlc.PartnerTier{} // Default to empty slice
	}

	// Fetch active partner testimonials for credibility
	// Testimonials from partner companies help build trust with potential partners
	testimonials, err := h.queries.ListActiveTestimonials(ctx)
	if err != nil {
		h.logger.Error("failed to load testimonials", "error", err)
		testimonials = []sqlc.ListActiveTestimonialsRow{} // Default to empty slice
	}

	// Group partners by tier name for easier template rendering
	// This creates a map like: {"Platinum": [partner1, partner2], "Gold": [partner3, partner4]}
	// Template can then iterate over tiers and render partners within each tier section
	partnersByTier := make(map[string][]sqlc.ListPartnersByTierRow)
	for _, p := range allPartners {
		partnersByTier[p.TierName] = append(partnersByTier[p.TierName], p)
	}

	// Build template data map with all fetched content
	data := map[string]interface{}{
		"Title":          "Partners",        // Page title for <title> tag and H1
		"CurrentPage":    "partners",        // Used by navigation to highlight active link
		"PartnersByTier": partnersByTier,    // Map of tier name to array of partner objects
		"Tiers":          tiers,             // Array of tier definitions for rendering tier sections
		"Testimonials":   testimonials,      // Array of testimonial objects
	}

	// Render template and cache for 5 minutes, return HTML to client
	return h.renderAndCache(c, cacheKey, 300, http.StatusOK, "public/pages/partners.html", data)
}
