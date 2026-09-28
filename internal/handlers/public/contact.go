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
	"net/mail"
	"regexp"
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
	queries *sqlc.Queries   // Database query interface for fetching office locations and storing submissions
	logger  *slog.Logger    // Structured logger for debugging and error tracking
	cache   *services.Cache // In-memory cache for rendered HTML to improve response times
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
	if cacheKey != "" {
		h.cache.Set(cacheKey, html, ttlSeconds)
	}

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
	productSKU := strings.TrimSpace(c.QueryParam("product"))
	// Product-specific forms bypass shared cache: context must never leak between visitors.
	if productSKU != "" {
		cacheKey = ""
	}

	// Check if cached version exists and return it immediately to improve performance
	// Contact page can be cached longer (1 hour) since office locations rarely change
	if cached, ok := h.cache.Get(cacheKey); cacheKey != "" && ok {
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

	// Build map URLs separately from stored data so arbitrary URLs cannot become frames.
	type contactOffice struct {
		sqlc.GetActiveOfficeLocationsRow
		MapEmbedURL      string
		MapDirectionsURL string
	}
	mappedOffices := make([]contactOffice, 0, len(offices))
	for _, office := range offices {
		address := strings.Join([]string{office.AddressLine1, office.AddressLine2.String, office.City, office.State, office.PostalCode, office.Country}, ", ")
		embed, directions := services.OfficeMapURLs(office.MapUrl, address)
		mappedOffices = append(mappedOffices, contactOffice{office, embed, directions})
	}

	// Build template data map
	data := map[string]interface{}{
		"Title":       "Contact Us",  // Page title for <title> tag and H1
		"Offices":     mappedOffices, // Array of office location objects for display
		"CurrentPage": "contact",     // Used by navigation to highlight active link
	}

	if productSKU != "" {
		product, err := h.queries.GetProductBySKU(ctx, productSKU)
		if err == nil && product.Status == "published" {
			data["QuoteProduct"] = product
			data["InitialMessage"] = "Please send me a quote for " + product.Name + " (SKU: " + product.Sku + ")."
		} else {
			data["QuoteUnavailable"] = true
		}
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
		return contactFormError(c, "Name, email, phone, company, and message are required.")
	}

	if !validContactEmail(email) {
		return contactFormError(c, "Enter a valid email address, such as name@example.com.")
	}
	if !validContactPhone(phone) {
		return contactFormError(c, "Enter a valid phone number with 7–15 digits. Spaces, parentheses, hyphens and a leading + are allowed.")
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
	return c.HTML(http.StatusOK, `<div role="status" class="alert alert-success border-2 border-green-700 bg-green-50 p-4 font-mono">Thank you for your message. We will get back to you shortly.</div>`)
}

// Errors keep the form and its entered values visible. The page explicitly swaps
// validation responses, since HTMX does not swap HTTP 400 responses by default.
func contactFormError(c echo.Context, message string) error {
	return c.HTML(http.StatusBadRequest, `<div role="alert" class="border-2 border-red-700 bg-red-50 p-4 font-mono text-red-800">`+message+`</div>`)
}

var contactDomainLabel = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
var contactTLD = regexp.MustCompile(`^[a-zA-Z]{2,63}$`)

func validContactEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) > 64 || strings.ContainsAny(parts[0], "\" ()<>\\") {
		return false
	}
	labels := strings.Split(parts[1], ".")
	if len(labels) < 2 || !contactTLD.MatchString(labels[len(labels)-1]) {
		return false
	}
	for _, label := range labels {
		if !contactDomainLabel.MatchString(label) {
			return false
		}
	}
	return true
}

func validContactPhone(phone string) bool {
	if len(phone) > 40 {
		return false
	}
	digits := ""
	open := false
	groupDigits := 0
	for i, ch := range phone {
		switch {
		case ch >= '0' && ch <= '9':
			digits += string(ch)
			if open {
				groupDigits++
			}
		case ch == '+' && i == 0:
		case ch == '(' && !open:
			open = true
			groupDigits = 0
		case ch == ')' && open && groupDigits > 0:
			open = false
		case ch == ' ' || ch == '-':
		default:
			return false
		}
	}
	if open || len(digits) < 7 || len(digits) > 15 {
		return false
	}
	return strings.Trim(digits, digits[:1]) != ""
}

// ShowPrivacyNotice explains the actual enquiry fields and their intended use.
func (h *ContactHandler) ShowPrivacyNotice(c echo.Context) error {
	return h.renderAndCache(c, "page:privacy", 3600, http.StatusOK, "public/pages/privacy.html", map[string]interface{}{
		"Title": "Privacy Notice", "CanonicalURL": "/privacy", "CurrentPage": "privacy",
	})
}
