// Package public provides HTTP handlers for public-facing pages of the Bluejay CMS.
// These handlers serve the front-end website content to visitors (non-admin users).
package public

import (
	// bytes provides buffer operations for building HTML output before sending to client
	"bytes"
	// database/sql provides sql.NullString and sql.ErrNoRows for handling nullable fields and query results
	"database/sql"
	// fmt provides string formatting for building cache keys and URLs
	"fmt"
	// log/slog is the structured logging library used for debug and error logging
	"log/slog"
	// net/http provides HTTP constants and status codes
	"net/http"
	// strconv provides string to integer conversion for parsing query parameters
	"strconv"
	// strings provides string manipulation utilities like TrimSpace for input sanitization
	"strings"

	// echo is the web framework used for routing and request/response handling
	"github.com/labstack/echo/v4"
	// sqlc provides type-safe database query interfaces generated from SQL files
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	// services provides business logic components like caching
	"github.com/narendhupati/bluejay-cms/internal/services"
)

// WhitepapersHandler handles HTTP requests for whitepaper pages.
// It manages the whitepapers listing page (with topic filtering), individual whitepaper
// detail pages (with preview mode for admins), and the gated download form submission.
type WhitepapersHandler struct {
	queries *sqlc.Queries    // Database query interface for fetching whitepapers, topics, and recording downloads
	logger  *slog.Logger     // Structured logger for debugging and error tracking
	cache   *services.Cache  // In-memory cache for rendered HTML to improve response times
}

// NewWhitepapersHandler constructs a new WhitepapersHandler with required dependencies.
// This constructor is called during application initialization to wire up dependencies.
func NewWhitepapersHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *WhitepapersHandler {
	return &WhitepapersHandler{
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
//   - ttlSeconds: Time-to-live in seconds; cache expires after this duration (0 = no cache)
//   - statusCode: HTTP status code to return (typically 200 OK)
//   - templateName: Path to the Go html/template file to render
//   - data: Template variables to pass to the template
//
// Returns HTML response with the specified status code, or error if rendering fails.
func (h *WhitepapersHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// WhitepapersList handles GET requests to /whitepapers
// Renders the whitepapers listing page with optional topic filtering.
//
// Route: GET /whitepapers (with optional ?topic=<id> query parameter)
// Template: templates/public/pages/whitepapers.html (full page, not HTMX fragment)
// Cache: 600 seconds (10 minutes) - whitepaper listings are relatively static
//
// HTMX Behavior: This endpoint returns a full HTML page, not an HTMX fragment.
// It is designed for direct browser navigation. However, the topic filter dropdown
// may use HTMX to reload the page content when selection changes.
//
// Query Parameters:
//   - topic (optional): Integer ID of topic to filter by (e.g., ?topic=2)
//
// Business Logic: If topic filter is applied, shows only whitepapers for that topic.
// Otherwise shows all published whitepapers. Each view has separate cache entries.
// Only published whitepapers are shown (drafts excluded from public view).
//
// Returns: HTTP 200 with rendered whitepapers.html template, or HTTP 400 if topic param is invalid
func (h *WhitepapersHandler) WhitepapersList(c echo.Context) error {
	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Parse optional topic filter from query string
	topicParam := c.QueryParam("topic")
	var selectedTopicID int64
	var whitepapers interface{}
	var totalCount int64
	var err error

	// Convert topic parameter to integer if provided
	if topicParam != "" {
		selectedTopicID, err = strconv.ParseInt(topicParam, 10, 64)
		if err != nil {
			// Invalid topic ID format - return 400 error
			h.logger.Error("invalid topic parameter", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid topic parameter")
		}
	}

	// Determine cache key based on whether filter is applied
	// Different cache entries for filtered vs unfiltered views
	cacheKey := "page:whitepapers"
	if selectedTopicID > 0 {
		cacheKey = fmt.Sprintf("page:whitepapers:topic:%d", selectedTopicID)
	}

	// Check if cached version exists and return it immediately
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Fetch whitepapers based on whether topic filter is applied
	// Use different query depending on filter state
	if selectedTopicID > 0 {
		// Filtered view: only published whitepapers for selected topic
		whitepapers, err = h.queries.ListPublishedWhitepapersByTopic(ctx, selectedTopicID)
		if err != nil {
			h.logger.Error("failed to list whitepapers by topic", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Get count for display (e.g., "Showing 8 whitepapers")
		totalCount, err = h.queries.CountPublishedWhitepapersByTopic(ctx, selectedTopicID)
		if err != nil {
			h.logger.Error("failed to count whitepapers by topic", "error", err)
			totalCount = 0 // Graceful degradation - continue without count
		}
	} else {
		// Unfiltered view: all published whitepapers
		whitepapers, err = h.queries.ListPublishedWhitepapers(ctx)
		if err != nil {
			h.logger.Error("failed to list whitepapers", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Get total count of all published whitepapers
		totalCount, err = h.queries.CountPublishedWhitepapers(ctx)
		if err != nil {
			h.logger.Error("failed to count whitepapers", "error", err)
			totalCount = 0 // Graceful degradation
		}
	}

	// Fetch all whitepaper topics for filter dropdown
	// Topics categorize whitepapers (e.g., "Security", "Performance", "Best Practices")
	topics, err := h.queries.ListWhitepaperTopics(ctx)
	if err != nil {
		h.logger.Error("failed to list whitepaper topics", "error", err)
		topics = []sqlc.WhitepaperTopic{} // Default to empty slice
	}

	// Build template data map
	data := map[string]interface{}{
		"Title":           "Whitepapers",      // Page title for <title> tag and H1
		"Whitepapers":     whitepapers,        // Array of whitepaper objects (filtered or all)
		"Topics":          topics,             // Array of topic objects for filter dropdown
		"SelectedTopicID": selectedTopicID,    // Currently selected topic ID (0 if none)
		"TotalCount":      totalCount,         // Number of whitepapers being displayed
		"CurrentPage":     "whitepapers",      // Used by navigation to highlight active link
	}

	// Render template and cache for 10 minutes, return HTML to client
	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/whitepapers.html", data)
}

// WhitepaperDetail handles GET requests to /whitepapers/:slug
// Renders an individual whitepaper detail page with description, learning points, and download form.
// Supports preview mode for admins to view draft whitepapers before publishing.
//
// Route: GET /whitepapers/:slug (with optional ?preview=true query parameter for admins)
// Template: templates/public/pages/whitepaper_detail.html (full page, not HTMX fragment)
// Cache: 900 seconds (15 minutes) for published, 0 seconds for preview mode
//
// HTMX Behavior: This endpoint returns a full HTML page, not an HTMX fragment.
// The page includes a download form that uses HTMX to submit (see WhitepaperDownload handler).
//
// Preview Mode: When ?preview=true is in query string (checked by isPreviewRequest helper),
// the handler fetches draft whitepapers and skips caching. This allows admins to preview
// unpublished content before making it live. Preview mode adds edit link to the page.
//
// Related Content: Fetches learning points (key takeaways) and related whitepapers
// from the same topic to encourage further engagement.
//
// Returns: HTTP 200 with rendered whitepaper_detail.html, or HTTP 404 if slug not found
func (h *WhitepapersHandler) WhitepaperDetail(c echo.Context) error {
	// Extract slug from URL path parameter (e.g., /whitepapers/cloud-security-best-practices)
	slug := c.Param("slug")
	// Check if this is a preview request (admins viewing draft content)
	preview := isPreviewRequest(c)

	// Skip cache check for preview mode - always fetch fresh data for admins
	if !preview {
		cacheKey := fmt.Sprintf("page:whitepapers:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Variables to hold whitepaper data regardless of preview vs published query
	var wpID, wpTopicID int64
	var wpTitle, wpMetaTitle string
	var wpSlug, wpOgImage string
	var wpMetaDesc sql.NullString
	var wpObj interface{}

	// Fetch whitepaper using different query based on preview mode
	if preview {
		// Preview mode: include draft whitepapers (not yet published)
		// This allows admins to see how content will look before publishing
		wp, err := h.queries.GetWhitepaperBySlugIncludeDrafts(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Whitepaper not found")
		}
		if err != nil {
			h.logger.Error("failed to load whitepaper", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Extract fields from query result
		wpID, wpTopicID = wp.ID, wp.TopicID
		wpTitle, wpSlug, wpMetaTitle = wp.Title, wp.Slug, wp.MetaTitle
		wpOgImage, wpMetaDesc = wp.OgImage, wp.MetaDescription
		wpObj = wp
	} else {
		// Normal mode: only published whitepapers visible to public
		wp, err := h.queries.GetWhitepaperBySlug(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Whitepaper not found")
		}
		if err != nil {
			h.logger.Error("failed to load whitepaper", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Extract fields from query result
		wpID, wpTopicID = wp.ID, wp.TopicID
		wpTitle, wpSlug, wpMetaTitle = wp.Title, wp.Slug, wp.MetaTitle
		wpOgImage, wpMetaDesc = wp.OgImage, wp.MetaDescription
		wpObj = wp
	}

	// Fetch learning points (key takeaways) for this whitepaper
	// Learning points are bullet-style highlights of what reader will learn
	learningPoints, err := h.queries.GetWhitepaperLearningPoints(ctx, wpID)
	if err != nil {
		h.logger.Error("failed to load whitepaper learning points", "error", err)
		learningPoints = []sqlc.GetWhitepaperLearningPointsRow{} // Default to empty slice
	}

	// Fetch related whitepapers from the same topic
	// Encourages visitors to explore more content, increasing engagement
	relatedPapers, err := h.queries.GetRelatedWhitepapers(ctx, sqlc.GetRelatedWhitepapersParams{
		ID:      wpID,      // Exclude current whitepaper
		TopicID: wpTopicID, // Find whitepapers in same topic
	})
	if err != nil {
		h.logger.Error("failed to load related whitepapers", "error", err)
		relatedPapers = []sqlc.GetRelatedWhitepapersRow{} // Default to empty slice
	}

	// Extract meta description for SEO, handling nullable field
	metaDesc := ""
	if wpMetaDesc.Valid {
		metaDesc = wpMetaDesc.String
	}

	// Build template data map with all fetched content
	data := map[string]interface{}{
		"Title":           wpTitle,                                    // Page title for <title> tag
		"MetaTitle":       wpMetaTitle,                                // Custom SEO title
		"MetaDescription": metaDesc,                                   // SEO description
		"MetaDesc":        metaDesc,                                   // Alias for template compatibility
		"OGImage":         wpOgImage,                                  // Open Graph image for social sharing
		"CanonicalURL":    fmt.Sprintf("/whitepapers/%s", wpSlug),     // Canonical URL for SEO
		"Whitepaper":      wpObj,                                      // Main whitepaper object
		"LearningPoints":  learningPoints,                             // Array of learning point objects
		"RelatedPapers":   relatedPapers,                              // Array of related whitepaper objects
		"CurrentPage":     "whitepapers",                              // Used by nav to highlight active link
	}

	// Handle preview mode differently - no caching and add admin edit link
	if preview {
		data["IsPreview"] = true // Shows preview banner in template
		data["EditURL"] = fmt.Sprintf("/admin/whitepapers/%d/edit", wpID) // Link to admin editor
		// No caching (ttl=0) for preview mode - always fresh content for admins
		return h.renderAndCache(c, "preview:whitepaper:"+slug, 0, http.StatusOK, "public/pages/whitepaper_detail.html", data)
	}

	// Normal mode: cache for 15 minutes since whitepaper content is relatively static
	return h.renderAndCache(c, fmt.Sprintf("page:whitepapers:%s", slug), 900, http.StatusOK, "public/pages/whitepaper_detail.html", data)
}

// WhitepaperDownload handles POST requests to /whitepapers/:slug/download
// Processes gated whitepaper download form submissions - validates input, stores lead information
// in database, increments download counter, and returns success page with download link.
//
// Route: POST /whitepapers/:slug/download
// Template: templates/public/pages/whitepaper_success.html (HTML fragment for HTMX swap)
//
// HTMX Behavior: This endpoint is designed for HTMX form submissions.
// The whitepaper detail page has a download form with hx-post="/whitepapers/{slug}/download".
// Upon successful submission, this returns an HTML fragment containing the download link
// and thank you message, which HTMX swaps into the page replacing the form.
//
// Form Fields:
//   - name (required): Visitor's full name
//   - email (required): Visitor's email address for lead capture
//   - company (required): Company/organization name
//   - designation (optional): Job title/role
//   - marketing_consent (optional): Checkbox for opt-in to marketing emails
//
// Business Logic: This implements a "gated content" strategy where visitors must provide
// contact information before accessing the whitepaper PDF. The information is stored for
// lead generation and marketing purposes. Download count is incremented for analytics.
//
// Returns: HTTP 200 with success HTML fragment containing download link, or HTTP 400/404/500 on error
func (h *WhitepapersHandler) WhitepaperDownload(c echo.Context) error {
	// Extract slug from URL path parameter
	slug := c.Param("slug")
	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Fetch whitepaper to verify it exists and get PDF file path
	whitepaper, err := h.queries.GetWhitepaperBySlug(ctx, slug)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Whitepaper not found")
	}
	if err != nil {
		h.logger.Error("failed to load whitepaper", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Parse and sanitize form values by trimming whitespace
	// This prevents issues with accidental spaces in input fields
	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	company := strings.TrimSpace(c.FormValue("company"))
	designation := strings.TrimSpace(c.FormValue("designation"))
	marketingConsent := c.FormValue("marketing_consent")

	// Validate required fields - reject submission if any are missing
	// Returns an error HTML fragment that HTMX will swap into the page
	if name == "" || email == "" || company == "" {
		return c.HTML(http.StatusBadRequest, `<div class="alert alert-error">Name, email, and company are required.</div>`)
	}

	// Parse marketing consent checkbox value
	// Checkboxes can send "on", "1", or "true" depending on implementation
	var consent int64
	if marketingConsent == "on" || marketingConsent == "1" || marketingConsent == "true" {
		consent = 1 // Store as 1 in database for TRUE
	} // Otherwise defaults to 0 for FALSE

	// Create whitepaper download record in database for lead tracking
	// This captures the visitor's information for marketing/sales follow-up
	_, err = h.queries.CreateWhitepaperDownload(ctx, sqlc.CreateWhitepaperDownloadParams{
		WhitepaperID: whitepaper.ID,
		Name:         name,
		Email:        email,
		Company:      company,
		// Designation is optional - use sql.NullString to handle empty value
		Designation: sql.NullString{
			String: designation,
			Valid:  designation != "", // Mark as valid only if non-empty
		},
		MarketingConsent: consent,
		// Capture visitor IP address for spam prevention and analytics
		IpAddress: sql.NullString{
			String: c.RealIP(),
			Valid:  c.RealIP() != "",
		},
		// Capture user agent for analytics and debugging
		UserAgent: sql.NullString{
			String: c.Request().UserAgent(),
			Valid:  c.Request().UserAgent() != "",
		},
	})
	if err != nil {
		h.logger.Error("failed to create whitepaper download", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Increment download count asynchronously in background goroutine
	// This updates the whitepaper's download_count field for analytics
	// We don't wait for this to complete - fire and forget for performance
	go func() {
		if err := h.queries.IncrementWhitepaperDownloadCount(ctx, whitepaper.ID); err != nil {
			h.logger.Error("failed to increment whitepaper download count", "error", err)
		}
	}()

	// Invalidate cache entries for this whitepaper since download count changed
	// This ensures next visitor sees updated download count
	h.cache.Delete(fmt.Sprintf("page:whitepapers:%s", slug))
	h.cache.Delete("page:whitepapers") // Also invalidate listing page

	// Build template data for success page fragment
	data := map[string]interface{}{
		"Whitepaper":    whitepaper,                      // Whitepaper object with title, description
		"Email":         email,                           // User's email for personalized message
		"WhitepaperURL": "/" + whitepaper.PdfFilePath,    // Path to PDF file for download link
	}
	// Inject settings from middleware for consistent branding
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}

	// Render success template fragment to buffer
	// This template shows thank you message and download link
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, "public/pages/whitepaper_success.html", data, c); err != nil {
		h.logger.Error("template render failed", "template", "public/pages/whitepaper_success.html", "error", err)
		return err
	}

	// Return success HTML fragment that HTMX will swap into the page
	// This replaces the download form with the success message and download link
	return c.HTML(http.StatusOK, buf.String())
}
