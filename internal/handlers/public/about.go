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

// AboutHandler handles HTTP requests for the About Us page and related content.
// It manages company overview, mission/vision/values, core values, milestones, and certifications.
type AboutHandler struct {
	queries *sqlc.Queries    // Database query interface for fetching about page data
	logger  *slog.Logger     // Structured logger for debugging and error tracking
	cache   *services.Cache  // In-memory cache for rendered HTML to improve response times
}

// NewAboutHandler constructs a new AboutHandler with required dependencies.
// This constructor is called during application initialization to wire up dependencies.
func NewAboutHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *AboutHandler {
	return &AboutHandler{
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
func (h *AboutHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// AboutPage handles GET requests to /about
// Renders the complete About Us page with company information, mission/vision/values,
// core values, company milestones, and industry certifications.
//
// Route: GET /about
// Template: templates/public/pages/about.html (full page, not HTMX fragment)
// Cache: 300 seconds (5 minutes) - about page content rarely changes
//
// HTMX Behavior: This endpoint returns a full HTML page, not an HTMX fragment.
// It is designed for direct browser navigation, not HTMX swaps.
//
// Returns: HTTP 200 with rendered about.html template
func (h *AboutHandler) AboutPage(c echo.Context) error {
	// Define cache key for this page
	cacheKey := "page:about"

	// Check if cached version exists and return it immediately to improve performance
	// This avoids database queries and template rendering for repeated requests
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Fetch company overview (single row with company description, history, etc.)
	// This is optional content - if not configured in admin, we gracefully continue
	overview, err := h.queries.GetCompanyOverview(ctx)
	if err != nil {
		h.logger.Debug("no company overview found", "error", err)
	}

	// Fetch mission, vision, and values statement (single row)
	// Also optional - company may not have configured this yet
	mvv, err := h.queries.GetMissionVisionValues(ctx)
	if err != nil {
		h.logger.Debug("no mission/vision/values found", "error", err)
	}

	// Fetch list of core values (e.g., "Integrity", "Innovation", "Customer Focus")
	// Each core value has a title, description, and optional icon
	coreValues, err := h.queries.ListCoreValues(ctx)
	if err != nil {
		h.logger.Error("failed to load core values", "error", err)
		coreValues = []sqlc.CoreValue{} // Default to empty slice to prevent template errors
	}

	// Fetch company milestones (e.g., "Founded 2015", "Reached 1M users 2020")
	// Displayed in a timeline format on the about page
	milestones, err := h.queries.ListMilestones(ctx)
	if err != nil {
		h.logger.Error("failed to load milestones", "error", err)
		milestones = []sqlc.Milestone{} // Default to empty slice
	}

	// Fetch industry certifications and accreditations (e.g., ISO 9001, SOC 2)
	// Displayed as badges or cards to build trust with visitors
	certs, err := h.queries.ListCertifications(ctx)
	if err != nil {
		h.logger.Error("failed to load certifications", "error", err)
		certs = []sqlc.Certification{} // Default to empty slice
	}

	// Build template data map with all fetched content
	data := map[string]interface{}{
		"Title":       "About Us",           // Page title for <title> tag and H1
		"CurrentPage": "about",              // Used by navigation to highlight active link
		"Overview":    nil,                  // Will be set below if data exists
		"MVV":         nil,                  // Will be set below if data exists
		"CoreValues":  coreValues,           // Array of core value objects
		"Milestones":  milestones,           // Array of milestone objects
		"Certs":       certs,                // Array of certification objects
	}

	// Only include overview if it was successfully fetched (ID > 0 means valid row)
	// This prevents template from trying to render nil/empty struct
	if overview.ID > 0 {
		data["Overview"] = overview
	}
	// Only include MVV if it was successfully fetched
	if mvv.ID > 0 {
		data["MVV"] = mvv
	}

	// Render template and cache for 5 minutes, return HTML to client
	return h.renderAndCache(c, cacheKey, 300, http.StatusOK, "public/pages/about.html", data)
}
