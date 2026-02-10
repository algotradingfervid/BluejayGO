// Package public provides HTTP handlers for the public-facing website.
// This file contains handlers for blog listing, category filtering,
// and individual blog post display with preview mode support.
package public

import (
	// Standard library imports
	"bytes"       // Buffer for template rendering to enable HTML caching
	"database/sql" // SQL error handling (sql.ErrNoRows for 404 detection)
	"fmt"         // String formatting for cache keys and template data
	"log/slog"    // Structured logging for debugging and error tracking
	"math"        // Math.Ceil for pagination calculations
	"net/http"    // HTTP status codes and request/response handling
	"strconv"     // String to integer conversion for page query parameters

	// Third-party imports
	"github.com/labstack/echo/v4" // Echo web framework - routing, context, rendering

	// Internal imports
	"github.com/narendhupati/bluejay-cms/db/sqlc"          // sqlc-generated database queries
	"github.com/narendhupati/bluejay-cms/internal/services" // Cache service for HTML caching
)

// BlogHandler handles all blog-related public routes including
// the main blog listing, category filtering, and individual post pages.
// It implements caching for improved performance on frequently accessed pages.
type BlogHandler struct {
	queries *sqlc.Queries    // Database query interface for blog posts and categories
	logger  *slog.Logger     // Structured logger for errors and debugging
	cache   *services.Cache  // In-memory cache for rendered HTML pages
}

// NewBlogHandler creates a new BlogHandler with the required dependencies.
// The cache is used to store rendered HTML to reduce database queries and
// template rendering overhead for frequently accessed blog pages.
func NewBlogHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *BlogHandler {
	return &BlogHandler{queries: queries, logger: logger, cache: cache}
}

// renderAndCache is a helper method that renders a template to HTML,
// caches the result, and returns it to the client.
//
// Purpose:
// This method centralizes the render-and-cache pattern used across all
// blog handlers. It ensures consistent injection of middleware data
// (settings, footer navigation) and caches the rendered HTML to reduce
// database load and template rendering overhead.
//
// Parameters:
//   - cacheKey: Unique identifier for the cached HTML (e.g., "page:blog:page:1:category:")
//   - ttlSeconds: Time-to-live for the cached HTML (0 means no expiration)
//   - statusCode: HTTP status code to return (typically http.StatusOK)
//   - templateName: Name of the template to render (e.g., "public/pages/blog_listing.html")
//   - data: Template data map to populate the template
//
// Caching Strategy:
//   - Blog listing: 300 seconds (5 minutes) - moderate TTL for fresh content
//   - Blog posts: 600 seconds (10 minutes) - longer TTL for stable content
//   - Preview mode: No cache (TTL=0) - to show live changes
//
// Middleware Data Injection:
//   - Settings: Global site configuration
//   - FooterCategories: Product categories for footer navigation
//   - FooterSolutions: Solutions for footer navigation
//   - FooterResources: Resources for footer navigation
func (h *BlogHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// BlogListing handles GET requests to the main blog listing page.
//
// HTTP Method: GET
// Route: /blog
// Query Parameters: ?page=N (optional), ?category=slug (optional)
// Template: public/pages/blog_listing.html (full page, not HTMX fragment)
// HTMX: Returns complete HTML page
// Cache TTL: 300 seconds (5 minutes)
//
// Purpose:
// Displays a paginated list of published blog posts, optionally filtered
// by category. Shows a featured post at the top (if one is set) and
// provides category navigation for filtering.
//
// Query Parameters:
//   - page: Page number for pagination (default: 1, minimum: 1)
//   - category: Category slug to filter posts (optional, shows all if empty)
//
// Template Data:
//   - Title: "Blog" - Browser tab title
//   - Posts: []sqlc.ListPublishedPostsRow - All published posts (when no category filter)
//   - CatPosts: []sqlc.ListPublishedPostsByCategoryRow - Category-filtered posts
//   - FeaturedPost: sqlc.BlogPost - Featured post to highlight at top
//   - Categories: []sqlc.BlogCategory - All blog categories for navigation
//   - CurrentCategory: string - Currently selected category slug (or empty)
//   - CurrentPage: "blog" - For navigation highlighting
//   - Page: int - Current page number
//   - TotalPages: int - Total number of pages for pagination
//   - TotalCount: int64 - Total number of posts (for "X results" display)
//
// Pagination:
//   - 9 posts per page (hard-coded limit)
//   - Invalid/negative page numbers default to page 1
//   - Total pages calculated using ceiling division
//
// Category Filtering:
//   - When category query param is present, filters posts by that category
//   - Uses different query results (Posts vs CatPosts) for all/filtered views
//   - Template can check which to display based on presence
//
// SEO Considerations:
//   - Blog is a key content marketing page for organic traffic
//   - Pagination helps manage large post archives
//   - Featured post highlights important/recent content
//   - Category filtering helps users find relevant content
//
// Caching:
//   - Cached for 5 minutes to balance freshness with performance
//   - Cache key includes page number and category for granular invalidation
//   - Cache is invalidated when new posts are published or categories change
func (h *BlogHandler) BlogListing(c echo.Context) error {
	// Extract query parameters
	pageStr := c.QueryParam("page")
	categorySlug := c.QueryParam("category")

	// Parse and validate page number
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1 // Default to first page for invalid values
	}

	// Check cache for this specific page/category combination
	cacheKey := fmt.Sprintf("page:blog:page:%d:category:%s", page, categorySlug)
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	// Set up pagination parameters
	limit := int64(9) // Show 9 posts per page (3x3 grid works well)
	offset := int64((page - 1)) * limit

	// Variables to hold query results (one will be populated based on category filter)
	var posts []sqlc.ListPublishedPostsRow          // All posts (no category filter)
	var catPosts []sqlc.ListPublishedPostsByCategoryRow // Category-filtered posts
	var totalCount int64
	var err error

	// Execute different query based on category filter
	if categorySlug != "" {
		// Fetch posts filtered by category
		catPosts, err = h.queries.ListPublishedPostsByCategory(ctx, sqlc.ListPublishedPostsByCategoryParams{
			Slug:   categorySlug,
			Limit:  limit,
			Offset: offset,
		})
		totalCount, _ = h.queries.CountPublishedPostsByCategory(ctx, categorySlug)
	} else {
		// Fetch all published posts
		posts, err = h.queries.ListPublishedPosts(ctx, sqlc.ListPublishedPostsParams{
			Limit:  limit,
			Offset: offset,
		})
		totalCount, _ = h.queries.CountPublishedPosts(ctx)
	}

	if err != nil {
		h.logger.Error("failed to list blog posts", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Fetch featured post to display prominently (errors ignored for graceful degradation)
	featuredPost, _ := h.queries.GetFeaturedPost(ctx)

	// Fetch all categories for navigation/filtering UI
	categories, _ := h.queries.ListBlogCategories(ctx)

	// Calculate total pages for pagination controls
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Assemble template data
	data := map[string]interface{}{
		"Title":           "Blog",          // Browser tab title
		"Posts":           posts,           // All posts (when no filter)
		"CatPosts":        catPosts,        // Category-filtered posts
		"FeaturedPost":    featuredPost,    // Featured post to highlight
		"Categories":      categories,      // All categories for nav
		"CurrentCategory": categorySlug,    // Selected category
		"CurrentPage":     "blog",          // For nav highlighting
		"Page":            page,            // Current page number
		"TotalPages":      totalPages,      // Total pages
		"TotalCount":      totalCount,      // Total posts count
	}

	// Render template and cache for 5 minutes
	// Template: templates/public/pages/blog_listing.html
	return h.renderAndCache(c, cacheKey, 300, http.StatusOK, "public/pages/blog_listing.html", data)
}

// BlogPost handles GET requests to view individual blog posts.
//
// HTTP Method: GET
// Route: /blog/:slug (e.g., /blog/introducing-our-new-sensor-lineup)
// Query Parameters: ?preview=1 (optional, for admin preview)
// Template: public/pages/blog_post.html (full page, not HTMX fragment)
// HTMX: Returns complete HTML page
// Cache TTL: 600 seconds (10 minutes) for published posts, no cache for previews
//
// Purpose:
// Displays a complete blog post with:
//   - Post content (HTML from Trix editor)
//   - Featured image
//   - Author information
//   - Publication date
//   - Tags for topical organization
//   - Related products (if applicable)
//   - SEO metadata (title, description, OG image)
//
// URL Parameters:
//   - slug: Unique slug identifier for the post (e.g., "new-product-announcement")
//
// Query Parameters:
//   - preview: If "1", shows draft/unpublished posts for admin preview
//
// Template Data:
//   - Title: Post title - Browser tab title
//   - MetaTitle: Custom SEO title (if set)
//   - MetaDescription: SEO meta description
//   - OGImage: Open Graph image for social sharing
//   - CanonicalURL: Canonical URL for SEO
//   - Post: sqlc.BlogPost - Core post data (title, content, author, date, etc.)
//   - Tags: []string - Topic tags for the post
//   - RelatedProducts: []sqlc.Product - Products mentioned/related to this post
//   - CurrentPage: "blog" - For navigation highlighting
//   - IsPreview: bool - True if viewing in preview mode (only in preview)
//   - EditURL: string - Admin edit URL (only in preview mode)
//
// SEO Considerations:
//   - Blog posts are primary content marketing assets for organic traffic
//   - Custom meta titles and descriptions improve search visibility
//   - OG images enhance social media sharing
//   - Canonical URLs prevent duplicate content issues
//   - Tags provide topical signals for search engines
//
// Preview Mode:
//   - When ?preview=1 is set, shows unpublished/draft posts
//   - Bypasses published status check for admin review
//   - Adds edit link for admin convenience
//   - Disables caching to show live changes immediately
//
// Related Content:
//   - Tags help users discover related content
//   - Related products link blog content to product pages
//   - This creates cross-linking opportunities for SEO
//
// Error Handling:
//   - Returns 404 if post slug doesn't exist
//   - Returns 500 on database errors
//   - Gracefully handles missing tags/products
//
// Caching:
//   - Published posts cached for 10 minutes
//   - Preview mode has no cache (TTL=0) for immediate updates
//   - Cache is invalidated when post is updated in admin
func (h *BlogHandler) BlogPost(c echo.Context) error {
	slug := c.Param("slug")
	preview := isPreviewRequest(c) // Check if this is an admin preview request

	// Skip cache lookup for preview mode to show live changes
	if !preview {
		cacheKey := fmt.Sprintf("page:blog:post:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	ctx := c.Request().Context()

	// Variables to store post data (extracted from different query result types)
	var post interface{}
	var postID int64
	var postTitle, postSlug, postMetaTitle, metaDesc string
	var postOgImage string
	var postMetaDesc sql.NullString

	// Execute different query based on preview mode
	if preview {
		// Preview mode: include draft/unpublished posts
		p, err := h.queries.GetPostBySlugIncludeDrafts(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		}
		if err != nil {
			h.logger.Error("failed to load blog post", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Extract fields from preview query result
		post = p
		postID = p.ID
		postTitle = p.Title
		postSlug = p.Slug
		postMetaTitle = p.MetaTitle
		postOgImage = p.OgImage
		postMetaDesc = p.MetaDescription
	} else {
		// Normal mode: only show published posts
		p, err := h.queries.GetPublishedPostBySlug(ctx, slug)
		if err == sql.ErrNoRows {
			// Post doesn't exist or is not published
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		}
		if err != nil {
			h.logger.Error("failed to load blog post", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		// Extract fields from published query result
		post = p
		postID = p.ID
		postTitle = p.Title
		postSlug = p.Slug
		postMetaTitle = p.MetaTitle
		postOgImage = p.OgImage
		postMetaDesc = p.MetaDescription
	}

	// Fetch associated data (tags and related products)
	// Errors are ignored for graceful degradation
	tags, _ := h.queries.GetPostTagsByPostID(ctx, postID)                   // Topic tags for this post
	relatedProducts, _ := h.queries.GetPostProductsByPostID(ctx, postID)    // Products mentioned in this post

	// Extract meta description with null-safety
	metaDesc = ""
	if postMetaDesc.Valid && postMetaDesc.String != "" {
		metaDesc = postMetaDesc.String
	}

	// Assemble template data
	data := map[string]interface{}{
		"Title":           postTitle,                              // Browser tab title
		"MetaTitle":       postMetaTitle,                          // SEO title
		"MetaDescription": metaDesc,                               // SEO description
		"OGImage":         postOgImage,                            // Social sharing image
		"CanonicalURL":    fmt.Sprintf("/blog/%s", postSlug),      // SEO canonical URL
		"Post":            post,                                   // Full post data
		"Tags":            tags,                                   // Topic tags
		"RelatedProducts": relatedProducts,                        // Related products
		"CurrentPage":     "blog",                                 // For nav highlighting
	}

	// Handle preview mode
	if preview {
		data["IsPreview"] = true // Show preview banner in template
		data["EditURL"] = fmt.Sprintf("/admin/blog/posts/%d/edit", postID) // Link to admin editor
		// Don't cache preview pages (TTL=0)
		return h.renderAndCache(c, "preview:blog:"+slug, 0, http.StatusOK, "public/pages/blog_post.html", data)
	}

	// Render and cache for 10 minutes (600 seconds)
	// Template: templates/public/pages/blog_post.html
	return h.renderAndCache(c, fmt.Sprintf("page:blog:post:%s", slug), 600, http.StatusOK, "public/pages/blog_post.html", data)
}
