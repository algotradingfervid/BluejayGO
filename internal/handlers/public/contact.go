// Package public provides HTTP handlers for public-facing pages of the Bluejay CMS.
// These handlers serve the front-end website content to visitors (non-admin users).
package public

import (
	// bytes provides buffer operations for building HTML output before sending to client
	"bytes"
	// database/sql provides sql.NullString and other SQL nullable types
	"database/sql"
	// log/slog is the structured logging library used for debug and error logging
	"log/slog"
	// net/http provides HTTP constants and status codes
	"net/http"
	// strings provides string manipulation utilities like TrimSpace for input sanitization
	"strings"

	// echo is the web framework used for routing and request/response handling
	"github.com/labstack/echo/v4"
	// sqlc provides type-safe database query interfaces generated from SQL files
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	// services provides business logic components like caching
	"github.com/narendhupati/bluejay-cms/internal/services"
)

// ContactHandler handles HTTP requests for the Contact Us page and form submissions.
// It manages displaying office locations and processing contact form submissions with validation.
type ContactHandler struct {
	queries *sqlc.Queries    // Database query interface for fetching office locations and storing submissions
	logger  *slog.Logger     // Structured logger for debugging and error tracking
	cache   *services.Cache  // In-memory cache for rendered HTML to improve response times
}

// NewContactHandler constructs a new ContactHandler with required dependencies.
// This constructor is called during application initialization to wire up dependencies.
func NewContactHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *ContactHandler {
	return &ContactHandler{
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
func (h *ContactHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// ShowContactPage handles GET requests to /contact
// Renders the Contact Us page with contact form and office location information.
//
// Route: GET /contact
// Template: templates/public/pages/contact.html (full page, not HTMX fragment)
// Cache: 3600 seconds (1 hour) - office locations rarely change
//
// HTMX Behavior: This endpoint returns a full HTML page, not an HTMX fragment.
// It is designed for direct browser navigation, not HTMX swaps.
//
// Returns: HTTP 200 with rendered contact.html template
func (h *ContactHandler) ShowContactPage(c echo.Context) error {
	// Define cache key for this page
	cacheKey := "page:contact"

	// Check if cached version exists and return it immediately to improve performance
	// Contact page can be cached longer (1 hour) since office locations rarely change
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Fetch active office locations (only offices marked as active in admin)
	// Each office includes address, phone, email, hours, and optional map coordinates
	offices, err := h.queries.GetActiveOfficeLocations(ctx)
	if err != nil {
		h.logger.Error("failed to load office locations", "error", err)
		offices = []sqlc.GetActiveOfficeLocationsRow{} // Default to empty slice to prevent template errors
	}

	// Build template data map
	data := map[string]interface{}{
		"Title":       "Contact Us",  // Page title for <title> tag and H1
		"Offices":     offices,       // Array of office location objects for display
		"CurrentPage": "contact",     // Used by navigation to highlight active link
	}

	// Render template and cache for 1 hour, return HTML to client
	return h.renderAndCache(c, cacheKey, 3600, http.StatusOK, "public/pages/contact.html", data)
}

// SubmitContactForm handles POST requests to /contact/submit
// Processes contact form submissions, validates input, stores in database, and returns success/error message.
//
// Route: POST /contact/submit
// Template: None (returns inline HTML fragment for HTMX swap)
//
// HTMX Behavior: This endpoint is designed for HTMX form submissions.
// It returns a small HTML fragment (success or error message) that HTMX swaps into the page.
// The contact form uses hx-post="/contact/submit" to submit asynchronously without page reload.
//
// Form Fields:
//   - name (required): Contact person's full name
//   - email (required): Contact email address
//   - phone (required): Contact phone number
//   - company (required): Company/organization name
//   - message (required): The inquiry message
//   - inquiry_type (optional): Category like "Sales", "Support", "Partnership"
//
// Returns: HTTP 200 with success message, or HTTP 400 with error message (both as HTML fragments)
func (h *ContactHandler) SubmitContactForm(c echo.Context) error {
	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Parse and sanitize form values by trimming whitespace
	// This prevents issues with accidental spaces in input fields
	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	phone := strings.TrimSpace(c.FormValue("phone"))
	company := strings.TrimSpace(c.FormValue("company"))
	message := strings.TrimSpace(c.FormValue("message"))
	inquiryType := strings.TrimSpace(c.FormValue("inquiry_type"))

	// Validate required fields - reject submission if any are missing
	// Returns an error HTML fragment that HTMX will swap into the page
	if name == "" || email == "" || phone == "" || company == "" || message == "" {
		return c.HTML(http.StatusBadRequest, `<div class="alert alert-error">Name, email, phone, company, and message are required.</div>`)
	}

	// Create contact submission record in database
	// This stores the inquiry for admin review in the admin panel
	_, err := h.queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Company: company,
		Message: message,
		// InquiryType is optional - use sql.NullString to handle empty value
		InquiryType: sql.NullString{
			String: inquiryType,
			Valid:  inquiryType != "", // Mark as valid only if non-empty
		},
		// Capture visitor IP address for spam prevention and analytics
		IpAddress: sql.NullString{
			String: c.RealIP(),
			Valid:  c.RealIP() != "",
		},
		// Capture user agent for debugging and analytics
		UserAgent: sql.NullString{
			String: c.Request().UserAgent(),
			Valid:  c.Request().UserAgent() != "",
		},
	})
	if err != nil {
		h.logger.Error("failed to create contact submission", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Return success message HTML fragment that HTMX will swap into the page
	// This replaces the form or displays below it depending on hx-target configuration
	return c.HTML(http.StatusOK, `<div class="alert alert-success">Thank you for your message. We will get back to you shortly.</div>`)
}
