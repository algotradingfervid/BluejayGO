package public

import (
	"bytes"
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// SearchResult is a published product or article returned by site search.
type SearchResult struct {
	Type    string
	Title   string
	URL     string
	Excerpt string
}

type SearchHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewSearchHandler(db *sql.DB, logger *slog.Logger) *SearchHandler {
	return &SearchHandler{db: db, logger: logger}
}

// sanitizeQuery turns user words into literal FTS prefix terms.
func sanitizeQuery(q string) string {
	q = strings.TrimSpace(q)
	if q == "" {
		return ""
	}
	// Remove FTS5 special characters to prevent query syntax errors
	// These characters have special meaning in FTS5 MATCH queries
	replacer := strings.NewReplacer(
		"\"", "", // FTS5 phrase delimiter
		"*", "", // FTS5 prefix/wildcard operator
		"(", "", // FTS5 grouping operator
		")", "", // FTS5 grouping operator
		"+", "", // FTS5 AND operator
		"-", "", // FTS5 NOT operator
		"^", "", // FTS5 initial token operator
		":", "", // FTS5 column filter operator
		"{", "", // FTS5 NEAR operator delimiter
		"}", "", // FTS5 NEAR operator delimiter
		"~", "", // FTS5 NOT operator (alternative syntax)
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

// search returns a stable page of published products and articles. The same
// matching rules serve suggestions and the complete results page. A negative
// perTypeLimit leaves each type uncapped for complete paginated results.
func (h *SearchHandler) search(ctx context.Context, query string, limit, offset, perTypeLimit int64) ([]SearchResult, error) {
	ftsQuery := sanitizeQuery(query)
	if ftsQuery == "" {
		return nil, nil
	}
	rows, err := h.db.QueryContext(ctx, `
        WITH product_matches AS (
            SELECT 0 AS type_order, 'Product' AS result_type, p.name AS title,
                '/products/' || pc.slug || '/' || p.slug AS url,
                COALESCE(p.tagline, '') AS excerpt
            FROM products_fts f
            JOIN products p ON f.rowid = p.id
            JOIN product_categories pc ON p.category_id = pc.id
            WHERE products_fts MATCH ? AND p.status = 'published'
            ORDER BY p.name COLLATE NOCASE, p.slug LIMIT ?
        ), article_matches AS (
            SELECT 1 AS type_order, 'Article' AS result_type, bp.title AS title,
                '/blog/' || bp.slug AS url, COALESCE(bp.excerpt, '') AS excerpt
            FROM blog_posts_fts f
            JOIN blog_posts bp ON f.rowid = bp.id
            WHERE blog_posts_fts MATCH ? AND bp.status = 'published'
            ORDER BY bp.title COLLATE NOCASE, bp.slug LIMIT ?
        )
        SELECT result_type, title, url, excerpt FROM (
            SELECT * FROM product_matches UNION ALL SELECT * FROM article_matches
        ) ORDER BY type_order, title COLLATE NOCASE, url
        LIMIT ? OFFSET ?`, ftsQuery, perTypeLimit, ftsQuery, perTypeLimit, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []SearchResult
	for rows.Next() {
		var result SearchResult
		if err := rows.Scan(&result.Type, &result.Title, &result.URL, &result.Excerpt); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

const searchPageSize int64 = 20

func searchURL(query string, page int64) string {
	values := url.Values{"q": {query}}
	if page > 1 {
		values.Set("page", strconv.FormatInt(page, 10))
	}
	return "/search?" + values.Encode()
}

// SearchPage uses one extra result to detect another page without presenting a
// capped result count as a total. Invalid page input returns the first page.
func (h *SearchHandler) SearchPage(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	page := int64(1)
	if value, err := strconv.ParseUint(c.QueryParam("page"), 10, 31); err == nil && value > 0 {
		page = int64(value)
	}
	if query == "" {
		page = 1
	}
	results, err := h.search(c.Request().Context(), query, searchPageSize+1, (page-1)*searchPageSize, -1)
	if err != nil {
		h.logger.Error("site search failed", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Search is temporarily unavailable. Please try again.")
	}
	hasNext := int64(len(results)) > searchPageSize
	if hasNext {
		results = results[:searchPageSize]
	}
	data := map[string]interface{}{
		"Title": "Search", "Query": query, "Results": results, "Page": page,
		"FirstURL": searchURL(query, 1),
	}
	if page > 1 {
		data["PreviousURL"] = searchURL(query, page-1)
	}
	if hasNext {
		data["NextURL"] = searchURL(query, page+1)
	}
	for key, contextKey := range map[string]string{
		"Settings": "settings", "FooterCategories": "footer_categories",
		"FooterSolutions": "footer_solutions", "FooterResources": "footer_resources",
	} {
		if value := c.Get(contextKey); value != nil {
			data[key] = value
		}
	}
	return c.Render(http.StatusOK, "public/pages/search.html", data)
}

// SearchSuggest deliberately shows a short list with a route to every result.
func (h *SearchHandler) SearchSuggest(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	results, err := h.search(c.Request().Context(), query, 10, 0, 5)
	if err != nil {
		h.logger.Error("search suggestions failed", "error", err)
		// The modal can recover through the normal results page after a retry.
		return echo.NewHTTPError(http.StatusInternalServerError, "Search is temporarily unavailable. Please try again.")
	}
	data := map[string]interface{}{
		"Results": results, "Query": query, "SearchURL": searchURL(query, 1),
	}
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, "public/partials/search_suggestions.html", data, c); err != nil {
		h.logger.Error("search suggestions render failed", "error", err)
		return err
	}
	return c.HTML(http.StatusOK, buf.String())
}
