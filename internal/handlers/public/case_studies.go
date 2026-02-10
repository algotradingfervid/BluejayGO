// Package public provides HTTP handlers for public-facing pages of the Bluejay CMS.
// These handlers serve the front-end website content to visitors (non-admin users).
package public

import (
	// bytes provides buffer operations for building HTML output before sending to client
	"bytes"
	// database/sql provides sql.NullString and sql.ErrNoRows for handling nullable fields and query results
	"database/sql"
	// encoding/json provides JSON parsing for challenge bullets stored as JSON array in database
	"encoding/json"
	// fmt provides string formatting for building cache keys and URLs
	"fmt"
	// log/slog is the structured logging library used for debug and error logging
	"log/slog"
	// net/http provides HTTP constants and status codes
	"net/http"
	// strconv provides string to integer conversion for parsing query parameters
	"strconv"

	// echo is the web framework used for routing and request/response handling
	"github.com/labstack/echo/v4"
	// sqlc provides type-safe database query interfaces generated from SQL files
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	// services provides business logic components like caching
	"github.com/narendhupati/bluejay-cms/internal/services"
)

// CaseStudiesHandler handles HTTP requests for case study pages.
// It manages both the case studies listing page (with optional industry filtering)
// and individual case study detail pages (with preview mode for admins).
type CaseStudiesHandler struct {
	queries *sqlc.Queries    // Database query interface for fetching case studies, industries, products, and metrics
	logger  *slog.Logger     // Structured logger for debugging and error tracking
	cache   *services.Cache  // In-memory cache for rendered HTML to improve response times
}

// NewCaseStudiesHandler constructs a new CaseStudiesHandler with required dependencies.
// This constructor is called during application initialization to wire up dependencies.
func NewCaseStudiesHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *CaseStudiesHandler {
	return &CaseStudiesHandler{
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
func (h *CaseStudiesHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// CaseStudiesList handles GET requests to /case-studies
// Renders the case studies listing page with optional industry filtering.
//
// Route: GET /case-studies (with optional ?industry=<id> query parameter)
// Template: templates/public/pages/case_studies.html (full page, not HTMX fragment)
// Cache: 600 seconds (10 minutes) - case studies content is relatively static
//
// HTMX Behavior: This endpoint returns a full HTML page, not an HTMX fragment.
// It is designed for direct browser navigation. However, the industry filter dropdown
// may use HTMX to reload the page content when selection changes.
//
// Query Parameters:
//   - industry (optional): Integer ID of industry to filter by (e.g., ?industry=3)
//
// Business Logic: If industry filter is applied, shows only case studies for that industry.
// Otherwise shows all published case studies. Each view has separate cache entries.
//
// Returns: HTTP 200 with rendered case_studies.html template, or HTTP 400 if industry param is invalid
func (h *CaseStudiesHandler) CaseStudiesList(c echo.Context) error {
	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Parse optional industry filter from query string
	industryParam := c.QueryParam("industry")
	var selectedIndustryID int64
	var caseStudies interface{}
	var totalCount int64
	var err error

	// Convert industry parameter to integer if provided
	if industryParam != "" {
		selectedIndustryID, err = strconv.ParseInt(industryParam, 10, 64)
		if err != nil {
			// Invalid industry ID format - return 400 error
			h.logger.Error("invalid industry parameter", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid industry parameter")
		}
	}

	// Determine cache key based on whether filter is applied
	// Different cache entries for filtered vs unfiltered views
	cacheKey := "page:case-studies"
	if selectedIndustryID > 0 {
		cacheKey = fmt.Sprintf("page:case-studies:industry:%d", selectedIndustryID)
	}

	// Check if cached version exists and return it immediately
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Fetch case studies based on whether industry filter is applied
	// Use different query depending on filter state
	if selectedIndustryID > 0 {
		// Filtered view: only case studies for selected industry
		caseStudies, err = h.queries.ListCaseStudiesByIndustry(ctx, selectedIndustryID)
		if err != nil {
			h.logger.Error("failed to list case studies by industry", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Get count for display (e.g., "Showing 5 case studies")
		totalCount, err = h.queries.CountCaseStudiesByIndustry(ctx, selectedIndustryID)
		if err != nil {
			h.logger.Error("failed to count case studies by industry", "error", err)
			totalCount = 0 // Graceful degradation - continue without count
		}
	} else {
		// Unfiltered view: all published case studies
		caseStudies, err = h.queries.ListCaseStudies(ctx)
		if err != nil {
			h.logger.Error("failed to list case studies", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Get total count of all case studies
		totalCount, err = h.queries.CountCaseStudies(ctx)
		if err != nil {
			h.logger.Error("failed to count case studies", "error", err)
			totalCount = 0 // Graceful degradation
		}
	}

	// Fetch all industries for filter dropdown
	// Dropdown allows users to filter case studies by industry
	industries, err := h.queries.ListIndustries(ctx)
	if err != nil {
		h.logger.Error("failed to list industries", "error", err)
		industries = []sqlc.Industry{} // Default to empty slice
	}

	// Build template data map
	data := map[string]interface{}{
		"Title":              "Case Studies",      // Page title for <title> tag and H1
		"CaseStudies":        caseStudies,         // Array of case study objects (filtered or all)
		"Industries":         industries,          // Array of industry objects for filter dropdown
		"SelectedIndustryID": selectedIndustryID,  // Currently selected industry ID (0 if none)
		"TotalCount":         totalCount,          // Number of case studies being displayed
		"CurrentPage":        "case-studies",      // Used by navigation to highlight active link
	}

	// Render template and cache for 10 minutes, return HTML to client
	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/case_studies.html", data)
}

// CaseStudyDetail handles GET requests to /case-studies/:slug
// Renders an individual case study detail page with full content, products, and metrics.
// Supports preview mode for admins to view draft case studies before publishing.
//
// Route: GET /case-studies/:slug (with optional ?preview=true query parameter for admins)
// Template: templates/public/pages/case_study_detail.html (full page, not HTMX fragment)
// Cache: 1800 seconds (30 minutes) for published, 0 seconds for preview mode
//
// HTMX Behavior: This endpoint returns a full HTML page, not an HTMX fragment.
// It is designed for direct browser navigation, not HTMX swaps.
//
// Preview Mode: When ?preview=true is in query string (checked by isPreviewRequest helper),
// the handler fetches draft case studies and skips caching. This allows admins to preview
// unpublished content before making it live. Preview mode adds edit link to the page.
//
// Related Content: Fetches associated products and success metrics to display alongside
// the case study narrative. Challenge bullets are stored as JSON array and parsed here.
//
// Returns: HTTP 200 with rendered case_study_detail.html, or HTTP 404 if slug not found
func (h *CaseStudiesHandler) CaseStudyDetail(c echo.Context) error {
	// Extract slug from URL path parameter (e.g., /case-studies/acme-corp-success)
	slug := c.Param("slug")
	// Check if this is a preview request (admins viewing draft content)
	preview := isPreviewRequest(c)

	// Skip cache check for preview mode - always fetch fresh data for admins
	if !preview {
		cacheKey := fmt.Sprintf("page:case-studies:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	// Extract request context for passing to database queries
	ctx := c.Request().Context()

	// Variables to hold case study data regardless of preview vs published query
	var csID int64
	var csTitle, csSlug, csOgImage string
	var csMetaTitle, csMetaDesc, csBullets sql.NullString
	var caseStudyObj interface{}

	// Fetch case study using different query based on preview mode
	if preview {
		// Preview mode: include draft case studies (not yet published)
		// This allows admins to see how content will look before publishing
		cs, err := h.queries.GetCaseStudyBySlugIncludeDrafts(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Case study not found")
		}
		if err != nil {
			h.logger.Error("failed to load case study", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Extract fields from query result
		csID, csTitle, csSlug, csOgImage = cs.ID, cs.Title, cs.Slug, cs.OgImage
		csMetaTitle, csMetaDesc, csBullets = cs.MetaTitle, cs.MetaDescription, cs.ChallengeBullets
		caseStudyObj = cs
	} else {
		// Normal mode: only published case studies visible to public
		cs, err := h.queries.GetCaseStudyBySlug(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Case study not found")
		}
		if err != nil {
			h.logger.Error("failed to load case study", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Extract fields from query result
		csID, csTitle, csSlug, csOgImage = cs.ID, cs.Title, cs.Slug, cs.OgImage
		csMetaTitle, csMetaDesc, csBullets = cs.MetaTitle, cs.MetaDescription, cs.ChallengeBullets
		caseStudyObj = cs
	}

	// Fetch products featured in this case study
	// Shows which company products/services were used to achieve the results
	products, err := h.queries.GetCaseStudyProducts(ctx, csID)
	if err != nil {
		h.logger.Error("failed to load case study products", "error", err)
		products = []sqlc.GetCaseStudyProductsRow{} // Default to empty slice
	}

	// Fetch success metrics for this case study
	// Metrics show quantifiable results (e.g., "50% increase in sales", "2x faster processing")
	metrics, err := h.queries.GetCaseStudyMetrics(ctx, csID)
	if err != nil {
		h.logger.Error("failed to load case study metrics", "error", err)
		metrics = []sqlc.GetCaseStudyMetricsRow{} // Default to empty slice
	}

	// Parse challenge bullets from JSON array stored in database
	// Challenge bullets are a list of key problems the client was facing
	var challengeBullets []string
	if csBullets.Valid {
		if err := json.Unmarshal([]byte(csBullets.String), &challengeBullets); err != nil {
			h.logger.Error("failed to parse challenge bullets", "error", err)
			challengeBullets = []string{} // Fallback to empty array on parse error
		}
	} else {
		challengeBullets = []string{} // No bullets configured
	}

	// Extract meta description for SEO, handling nullable field
	metaDesc := ""
	if csMetaDesc.Valid {
		metaDesc = csMetaDesc.String
	}
	// Extract meta title for SEO (custom title tag), handling nullable field
	metaTitle := ""
	if csMetaTitle.Valid && csMetaTitle.String != "" {
		metaTitle = csMetaTitle.String
	}

	// Build template data map with all fetched content
	data := map[string]interface{}{
		"Title":            csTitle,           // Page title (used as fallback if MetaTitle empty)
		"MetaTitle":        metaTitle,         // Custom SEO title for <title> tag
		"MetaDescription":  metaDesc,          // SEO description for <meta name="description">
		"MetaDesc":         metaDesc,          // Alias for template compatibility
		"OGImage":          csOgImage,         // Open Graph image for social media sharing
		"CanonicalURL":     fmt.Sprintf("/case-studies/%s", csSlug), // Canonical URL for SEO
		"CaseStudy":        caseStudyObj,      // Main case study object with narrative content
		"Products":         products,          // Array of product objects featured in case study
		"Metrics":          metrics,           // Array of success metric objects
		"ChallengeBullets": challengeBullets,  // Array of challenge bullet strings
		"CurrentPage":      "case-studies",    // Used by navigation to highlight active link
	}

	// Handle preview mode differently - no caching and add admin edit link
	if preview {
		data["IsPreview"] = true // Shows preview banner in template
		data["EditURL"] = fmt.Sprintf("/admin/case-studies/%d/edit", csID) // Link to admin editor
		// No caching (ttl=0) for preview mode - always fresh content for admins
		return h.renderAndCache(c, "preview:case-study:"+slug, 0, http.StatusOK, "public/pages/case_study_detail.html", data)
	}

	// Normal mode: cache for 30 minutes since case study content rarely changes
	return h.renderAndCache(c, fmt.Sprintf("page:case-studies:%s", slug), 1800, http.StatusOK, "public/pages/case_study_detail.html", data)
}
