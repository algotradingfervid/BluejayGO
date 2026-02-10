// Package public provides HTTP handlers for public-facing website features.
// This file implements full-text search functionality across products, blog posts, and case studies.
package public

import (
	"bytes"        // Used for buffering HTML fragments before sending to client
	"database/sql" // Provides database/sql driver interfaces for SQLite database access
	"log/slog"     // Structured logging for search query tracking and error reporting
	"net/http"     // HTTP status codes and constants
	"strings"      // String manipulation for query sanitization and FTS5 token processing

	"github.com/labstack/echo/v4" // Echo web framework for HTTP request/response handling
)

// SearchResult represents a single search result from any content type.
// Used to unify results from products, blog posts, and case studies into a consistent format.
type SearchResult struct {
	Type    string // Content type: "Product", "Article", or "Case Study"
	Title   string // Display title of the content
	URL     string // Relative URL path to the content detail page
	Excerpt string // Short preview text (tagline for products, excerpt for blog posts)
}

// SearchHandler processes full-text search requests across multiple content types.
// Uses SQLite FTS5 (Full-Text Search) indexes for fast query performance.
type SearchHandler struct {
	db     *sql.DB     // Database connection for executing FTS5 queries
	logger *slog.Logger // Structured logger for tracking search queries and debugging errors
}

// NewSearchHandler creates a new search handler with database and logger dependencies.
func NewSearchHandler(db *sql.DB, logger *slog.Logger) *SearchHandler {
	return &SearchHandler{db: db, logger: logger}
}

// sanitizeQuery removes FTS5 special characters and adds prefix matching.
// Prevents FTS5 syntax errors from user input while enabling partial word matching.
//
// Security: Neutralizes FTS5 query operators (quotes, wildcards, boolean operators)
// that could cause SQL injection or unexpected query behavior.
//
// Behavior:
//   - Removes special FTS5 characters: " * ( ) + - ^ : { } ~
//   - Splits input into individual words
//   - Wraps each word in quotes and appends * for prefix matching
//   - Example: "bluetooth speaker" becomes "bluetooth"* "speaker"*
//
// This allows users to find "speaker" when searching for "spea" but prevents
// malicious FTS5 syntax from being executed.
func sanitizeQuery(q string) string {
	q = strings.TrimSpace(q)
	if q == "" {
		return ""
	}
	// Remove FTS5 special characters to prevent query syntax errors
	// These characters have special meaning in FTS5 MATCH queries
	replacer := strings.NewReplacer(
		"\"", "", // FTS5 phrase delimiter
		"*", "",  // FTS5 prefix/wildcard operator
		"(", "",  // FTS5 grouping operator
		")", "",  // FTS5 grouping operator
		"+", "",  // FTS5 AND operator
		"-", "",  // FTS5 NOT operator
		"^", "",  // FTS5 initial token operator
		":", "",  // FTS5 column filter operator
		"{", "",  // FTS5 NEAR operator delimiter
		"}", "",  // FTS5 NEAR operator delimiter
		"~", "",  // FTS5 NOT operator (alternative syntax)
	)
	q = replacer.Replace(q)
	q = strings.TrimSpace(q)
	if q == "" {
		return ""
	}

	// Split into words and add prefix matching
	// Each word is quoted to prevent phrase parsing and suffixed with * for prefix matching
	// Example: "gaming headset" → "gaming"* "headset"*
	words := strings.Fields(q)
	for i, w := range words {
		words[i] = "\"" + w + "\"" + "*"
	}
	return strings.Join(words, " ")
}

// search executes full-text search queries across products, blog posts, and case studies.
// Returns a unified list of results from all content types, limited by the limit parameter.
//
// FTS5 Implementation:
//   - Uses SQLite FTS5 virtual tables (products_fts, blog_posts_fts, case_studies_fts)
//   - Joins FTS5 results back to main tables via rowid for full data access
//   - MATCH clause uses sanitized query to prevent syntax errors
//   - Only searches published content (status = 'published')
//
// Query Processing:
//   - Sanitizes query using sanitizeQuery to prevent FTS5 syntax errors
//   - Returns nil if sanitized query is empty (invalid input)
//   - Executes three separate queries sequentially (products → blog → case studies)
//   - Errors are logged but don't stop subsequent searches (graceful degradation)
//
// Performance:
//   - FTS5 indexes provide fast full-text matching across title, content, and metadata
//   - LIMIT parameter controls result count per content type
//   - Results are appended to single slice (not sorted by relevance)
func (h *SearchHandler) search(query string, limit int) []SearchResult {
	ftsQuery := sanitizeQuery(query)
	if ftsQuery == "" {
		return nil
	}

	var results []SearchResult

	// Search products across name, tagline, and description fields
	// FTS5 index: products_fts includes searchable product content
	// URL format: /products/{category-slug}/{product-slug}
	rows, err := h.db.Query(
		`SELECT p.name, pc.slug, p.slug, COALESCE(p.tagline, '') FROM products_fts f JOIN products p ON f.rowid = p.id JOIN product_categories pc ON p.category_id = pc.id WHERE products_fts MATCH ? AND p.status = 'published' LIMIT ?`,
		ftsQuery, limit,
	)
	if err != nil {
		h.logger.Error("products fts query failed", "error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var name, catSlug, slug, tagline string
			if err := rows.Scan(&name, &catSlug, &slug, &tagline); err == nil {
				// Build hierarchical URL with category slug for better SEO
				results = append(results, SearchResult{
					Type:    "Product",
					Title:   name,
					URL:     "/products/" + catSlug + "/" + slug,
					Excerpt: tagline, // Product tagline used as search result preview
				})
			}
		}
	}

	// Search blog posts across title, excerpt, and content fields
	// FTS5 index: blog_posts_fts includes title, excerpt, and full HTML content
	// URL format: /blog/{post-slug}
	rows2, err := h.db.Query(
		`SELECT bp.title, bp.slug, bp.excerpt FROM blog_posts_fts f JOIN blog_posts bp ON f.rowid = bp.id WHERE blog_posts_fts MATCH ? AND bp.status = 'published' LIMIT ?`,
		ftsQuery, limit,
	)
	if err != nil {
		h.logger.Error("blog_posts fts query failed", "error", err)
	} else {
		defer rows2.Close()
		for rows2.Next() {
			var title, slug, excerpt string
			if err := rows2.Scan(&title, &slug, &excerpt); err == nil {
				// Blog posts labeled as "Article" for user-facing display
				results = append(results, SearchResult{
					Type:    "Article",
					Title:   title,
					URL:     "/blog/" + slug,
					Excerpt: excerpt, // Blog excerpt used as search result preview
				})
			}
		}
	}

	// Search case studies across title and content fields
	// FTS5 index: case_studies_fts includes title and HTML content
	// URL format: /case-studies/{study-slug}
	rows3, err := h.db.Query(
		`SELECT cs.title, cs.slug FROM case_studies_fts f JOIN case_studies cs ON f.rowid = cs.id WHERE case_studies_fts MATCH ? AND cs.status = 'published' LIMIT ?`,
		ftsQuery, limit,
	)
	if err != nil {
		h.logger.Error("case_studies fts query failed", "error", err)
	} else {
		defer rows3.Close()
		for rows3.Next() {
			var title, slug string
			if err := rows3.Scan(&title, &slug); err == nil {
				// Case studies have no excerpt field, only title is returned
				results = append(results, SearchResult{
					Type:  "Case Study",
					Title: title,
					URL:   "/case-studies/" + slug,
				})
			}
		}
	}

	return results
}

// SearchPage handles the main search results page.
//
// HTTP Method: GET
// Route: /search
// Query Parameters:
//   - q: search query string (optional, if empty shows empty results)
//
// Template: public/pages/search.html (full page render)
// HTMX: Not an HTMX endpoint - returns full HTML page
//
// Template Data:
//   - Title: Page title for <title> tag and breadcrumbs
//   - Query: User's search query (echoed back for display in search box)
//   - Results: Array of SearchResult objects (empty if no query or no matches)
//   - Settings: Site settings from middleware (logo, site name, etc.)
//   - FooterCategories: Product categories for footer navigation
//   - FooterSolutions: Solutions for footer navigation
//   - FooterResources: Resource links for footer navigation
//
// SEO Behavior:
//   - Search results pages are not indexed (should include noindex meta tag)
//   - Query parameter preserved in URL for sharing search results
//   - Empty query shows search form without results
//
// Performance:
//   - Limits results to 10 per content type (30 max total)
//   - FTS5 queries are fast even on large content sets
func (h *SearchHandler) SearchPage(c echo.Context) error {
	query := c.QueryParam("q")

	var results []SearchResult
	if query != "" {
		// Execute search across all content types, max 10 results per type
		results = h.search(query, 10)
	}

	// Build template data with search results and shared layout data
	data := map[string]interface{}{
		"Title":   "Search", // Page title for <title> tag
		"Query":   query,    // Echo query back for search box value
		"Results": results,  // Search results or nil if no query
	}

	// Add shared layout data from middleware (footer nav, site settings)
	// These are populated by middleware.LoadFooterData and middleware.LoadSettings
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

	// Render full search results page with layout
	return c.Render(http.StatusOK, "public/pages/search.html", data)
}

// SearchSuggest provides live search suggestions as the user types.
//
// HTTP Method: GET
// Route: /search/suggest
// Query Parameters:
//   - q: partial search query string (required for results)
//
// Template: public/partials/search_suggestions.html (HTML fragment)
// HTMX: This is an HTMX endpoint - returns HTML fragment, not full page
//
// HTMX Integration:
//   - Triggered by hx-get on search input field
//   - Uses hx-trigger="keyup changed delay:300ms" for debounced suggestions
//   - HTML fragment replaces suggestion dropdown container
//   - No layout/header/footer - just search result list
//
// Template Data:
//   - Results: Array of SearchResult objects (max 5 items)
//
// UX Behavior:
//   - Shows fewer results (5) than full search page (10) for faster rendering
//   - Empty query returns empty results (no suggestions shown)
//   - Renders partial template directly into page via HTMX swap
//
// Performance:
//   - Limits to 5 results total for fast response time
//   - FTS5 prefix matching enables "type-ahead" behavior
//   - Buffer used to render template before returning (error handling)
func (h *SearchHandler) SearchSuggest(c echo.Context) error {
	query := c.QueryParam("q")

	var results []SearchResult
	if query != "" {
		// Execute search with lower limit (5) for faster suggestion response
		results = h.search(query, 5)
	}

	// Build minimal template data - no layout data needed for HTMX fragment
	data := map[string]interface{}{
		"Results": results,
	}

	// Render template to buffer to catch errors before sending response
	// This prevents partial HTML being sent on template errors
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, "public/partials/search_suggestions.html", data, c); err != nil {
		h.logger.Error("search suggestions render failed", "error", err)
		return err
	}

	// Return HTML fragment for HTMX to inject into page
	return c.HTML(http.StatusOK, buf.String())
}
