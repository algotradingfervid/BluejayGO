// Package public provides HTTP handlers for public-facing website features.
// This file implements SEO-critical sitemap.xml and robots.txt generation.
package public

import (
	"encoding/xml" // XML marshaling for sitemap.xml standard format
	"fmt"          // String formatting for constructing URLs
	"log/slog"     // Structured logging for tracking sitemap generation errors
	"net/http"     // HTTP status codes and constants
	"time"         // Date formatting for sitemap lastmod timestamps

	"github.com/labstack/echo/v4"                      // Echo web framework for HTTP request/response handling
	"github.com/narendhupati/bluejay-cms/db/sqlc" // Generated sqlc database queries for fetching published content
)

// URLSet represents the root element of an XML sitemap following the sitemaps.org protocol.
// Conforms to https://www.sitemaps.org/protocol.html specification.
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`         // Root element name in XML
	XMLNS   string   `xml:"xmlns,attr"`     // XML namespace attribute (required by sitemap spec)
	URLs    []URL    `xml:"url"`            // Collection of URL entries in the sitemap
}

// URL represents a single URL entry in the sitemap with SEO metadata.
// All fields except Loc are optional but recommended for SEO best practices.
type URL struct {
	Loc        string `xml:"loc"`                    // Absolute URL of the page (required)
	LastMod    string `xml:"lastmod,omitempty"`      // Last modification date (ISO 8601 format: YYYY-MM-DD)
	ChangeFreq string `xml:"changefreq,omitempty"`   // Expected change frequency (always|hourly|daily|weekly|monthly|yearly|never)
	Priority   string `xml:"priority,omitempty"`     // Priority relative to other URLs (0.0 to 1.0, default 0.5)
}

// SitemapHandler generates XML sitemaps and robots.txt for search engine optimization.
// Sitemaps help search engines discover and index all website content efficiently.
type SitemapHandler struct {
	queries *sqlc.Queries // Database queries for fetching published content
	logger  *slog.Logger  // Structured logger for error tracking
	baseURL string        // Base URL for the website (e.g., "https://example.com")
}

// NewSitemapHandler creates a new sitemap handler with database queries, logger, and base URL.
// The baseURL should be the production domain without trailing slash.
func NewSitemapHandler(queries *sqlc.Queries, logger *slog.Logger, baseURL string) *SitemapHandler {
	return &SitemapHandler{queries: queries, logger: logger, baseURL: baseURL}
}

// Sitemap generates a complete XML sitemap for the website.
//
// HTTP Method: GET
// Route: /sitemap.xml
// Content-Type: application/xml
//
// Template: None - generates XML directly
// HTMX: Not an HTMX endpoint - returns XML document
//
// SEO Purpose:
//   - Helps search engines discover all website pages
//   - Provides crawl priority hints via priority field
//   - Indicates update frequency via changefreq field
//   - Includes lastmod timestamps for content freshness
//
// Sitemap Structure:
//   1. Static pages (homepage, category indexes, about, contact)
//   2. Dynamic content pages (solutions, blog posts, case studies, whitepapers)
//   3. All URLs are absolute (include baseURL)
//   4. Only published content is included
//
// Priority Guidelines:
//   - 1.0: Homepage (highest priority)
//   - 0.9: Main category indexes (products, solutions)
//   - 0.8: Blog index, case studies index
//   - 0.7: Individual content pages, about page
//   - 0.6: Contact page (lowest priority)
//
// Change Frequency Guidelines:
//   - daily: Blog index (new posts frequently)
//   - weekly: Category indexes, individual solutions
//   - monthly: Individual blog posts, case studies, static pages
//
// Error Handling:
//   - Database query errors are logged but don't fail sitemap generation
//   - Graceful degradation: partial sitemap still generated if one query fails
//   - XML marshaling errors return 500 error
func (h *SitemapHandler) Sitemap(c echo.Context) error {
	// Use current date as lastmod for static pages (updated during deploys)
	now := time.Now().Format("2006-01-02")

	// Initialize URLSet with required XML namespace
	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9", // Sitemap protocol 0.9 namespace
	}

	// Static pages with SEO metadata
	// These pages don't come from database - hardcoded navigation structure
	staticPages := []struct {
		path       string
		changeFreq string
		priority   string
	}{
		{"/", "weekly", "1.0"},              // Homepage: highest priority, updated weekly
		{"/products", "weekly", "0.9"},      // Product catalog index
		{"/solutions", "weekly", "0.9"},     // Solutions index
		{"/blog", "daily", "0.8"},           // Blog index: updated daily with new posts
		{"/case-studies", "weekly", "0.8"},  // Case studies index
		{"/whitepapers", "weekly", "0.8"},   // Whitepapers index
		{"/about", "monthly", "0.7"},        // About page: rarely changes
		{"/contact", "monthly", "0.6"},      // Contact page: lowest priority
		{"/partners", "monthly", "0.7"},     // Partners page
	}

	// Add static pages to sitemap with current date as lastmod
	for _, page := range staticPages {
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        h.baseURL + page.path, // Construct absolute URL
			LastMod:    now,                   // Use current date since static pages update with deploys
			ChangeFreq: page.changeFreq,
			Priority:   page.priority,
		})
	}

	// Solutions: individual solution detail pages
	// URL format: /solutions/{slug}
	// Only includes published solutions (status filtering in query)
	solutions, err := h.queries.ListPublishedSolutions(c.Request().Context())
	if err != nil {
		// Log error but continue generating sitemap with remaining content
		h.logger.Error("sitemap: failed to list solutions", "error", err)
	} else {
		for _, s := range solutions {
			u := URL{
				Loc:        fmt.Sprintf("%s/solutions/%s", h.baseURL, s.Slug),
				ChangeFreq: "weekly",   // Solutions updated weekly with new features/pricing
				Priority:   "0.8",      // High priority - key conversion pages
			}
			// Use actual update timestamp if available for accurate lastmod
			if s.UpdatedAt.Valid {
				u.LastMod = s.UpdatedAt.Time.Format("2006-01-02")
			}
			urlset.URLs = append(urlset.URLs, u)
		}
	}

	// Blog posts: individual article pages
	// URL format: /blog/{slug}
	// Limit set to 1000 - should cover most blogs, adjust if needed
	// Note: Very large blogs (>1000 posts) should implement sitemap index
	posts, err := h.queries.ListPublishedPosts(c.Request().Context(), sqlc.ListPublishedPostsParams{
		Limit:  1000, // Maximum posts to include in sitemap
		Offset: 0,    // Start from newest posts
	})
	if err != nil {
		// Log error but continue generating sitemap with remaining content
		h.logger.Error("sitemap: failed to list blog posts", "error", err)
	} else {
		for _, p := range posts {
			urlset.URLs = append(urlset.URLs, URL{
				Loc:        fmt.Sprintf("%s/blog/%s", h.baseURL, p.Slug),
				ChangeFreq: "monthly",  // Blog posts rarely updated after publication
				Priority:   "0.7",      // Medium priority - good for SEO but not conversion pages
			})
		}
	}

	// Case studies: customer success stories
	// URL format: /case-studies/{slug}
	// Important for B2B SEO - shows proof of value
	caseStudies, err := h.queries.ListCaseStudies(c.Request().Context())
	if err != nil {
		// Log error but continue generating sitemap with remaining content
		h.logger.Error("sitemap: failed to list case studies", "error", err)
	} else {
		for _, cs := range caseStudies {
			urlset.URLs = append(urlset.URLs, URL{
				Loc:        fmt.Sprintf("%s/case-studies/%s", h.baseURL, cs.Slug),
				ChangeFreq: "monthly",  // Case studies updated monthly with new metrics
				Priority:   "0.7",      // Medium-high priority - good for conversion
			})
		}
	}

	// Whitepapers: downloadable content/lead magnets
	// URL format: /whitepapers/{slug}
	// Often gated content for lead generation
	whitepapers, err := h.queries.ListPublishedWhitepapers(c.Request().Context())
	if err != nil {
		// Log error but continue generating sitemap with remaining content
		h.logger.Error("sitemap: failed to list whitepapers", "error", err)
	} else {
		for _, w := range whitepapers {
			urlset.URLs = append(urlset.URLs, URL{
				Loc:        fmt.Sprintf("%s/whitepapers/%s", h.baseURL, w.Slug),
				ChangeFreq: "monthly",  // Whitepapers rarely change after publication
				Priority:   "0.7",      // Medium priority - valuable for lead gen
			})
		}
	}

	// Marshal URLSet to formatted XML with 2-space indentation
	// Pretty-printed XML is easier for humans to read when debugging
	xmlData, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		h.logger.Error("sitemap: failed to marshal XML", "error", err)
		return c.String(http.StatusInternalServerError, "failed to generate sitemap")
	}

	// Prepend XML declaration header (required by XML spec)
	// Results in: <?xml version="1.0" encoding="UTF-8"?>
	xmlData = append([]byte(xml.Header), xmlData...)

	// Return XML with correct Content-Type header
	// application/xml tells search engines this is a sitemap
	return c.Blob(http.StatusOK, "application/xml", xmlData)
}

// RobotsTxt generates the robots.txt file for search engine crawlers.
//
// HTTP Method: GET
// Route: /robots.txt
// Content-Type: text/plain
//
// Template: None - generates plain text directly
// HTMX: Not an HTMX endpoint - returns plain text document
//
// SEO Purpose:
//   - Controls which paths search engines can crawl
//   - Points crawlers to sitemap.xml location
//   - Prevents indexing of admin panel and private areas
//
// Robots.txt Directives:
//   - User-agent: * (applies to all search engines)
//   - Allow: / (allow crawling of all public pages)
//   - Disallow: /admin/ (block admin panel from search engines)
//   - Sitemap: {baseURL}/sitemap.xml (tell crawlers where sitemap is)
//
// Security:
//   - Blocking /admin/ prevents admin URLs from appearing in search results
//   - robots.txt is NOT a security mechanism (add proper auth to admin routes)
//   - Well-behaved crawlers respect these rules; malicious bots may ignore them
//
// Best Practices:
//   - Keep simple - complex rules can confuse crawlers
//   - Always include sitemap location
//   - Use Disallow sparingly (only for truly private areas)
func (h *SitemapHandler) RobotsTxt(c echo.Context) error {
	// Generate robots.txt with sitemap reference
	// Format follows robots exclusion protocol standard
	robots := `User-agent: *
Allow: /
Disallow: /admin/

Sitemap: ` + h.baseURL + `/sitemap.xml`

	// Return plain text with text/plain Content-Type
	return c.String(http.StatusOK, robots)
}
