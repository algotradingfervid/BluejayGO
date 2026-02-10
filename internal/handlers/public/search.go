package public

import (
	"bytes"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

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

// sanitizeQuery removes FTS5 special characters and adds prefix matching
func sanitizeQuery(q string) string {
	q = strings.TrimSpace(q)
	if q == "" {
		return ""
	}
	// Remove FTS5 special characters
	replacer := strings.NewReplacer(
		"\"", "",
		"*", "",
		"(", "",
		")", "",
		"+", "",
		"-", "",
		"^", "",
		":", "",
		"{", "",
		"}", "",
		"~", "",
	)
	q = replacer.Replace(q)
	q = strings.TrimSpace(q)
	if q == "" {
		return ""
	}

	// Split into words and add prefix matching
	words := strings.Fields(q)
	for i, w := range words {
		words[i] = "\"" + w + "\"" + "*"
	}
	return strings.Join(words, " ")
}

func (h *SearchHandler) search(query string, limit int) []SearchResult {
	ftsQuery := sanitizeQuery(query)
	if ftsQuery == "" {
		return nil
	}

	var results []SearchResult

	// Search products
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
				results = append(results, SearchResult{
					Type:    "Product",
					Title:   name,
					URL:     "/products/" + catSlug + "/" + slug,
					Excerpt: tagline,
				})
			}
		}
	}

	// Search blog posts
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
				results = append(results, SearchResult{
					Type:    "Article",
					Title:   title,
					URL:     "/blog/" + slug,
					Excerpt: excerpt,
				})
			}
		}
	}

	// Search case studies
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

// GET /search
func (h *SearchHandler) SearchPage(c echo.Context) error {
	query := c.QueryParam("q")

	var results []SearchResult
	if query != "" {
		results = h.search(query, 10)
	}

	data := map[string]interface{}{
		"Title":   "Search",
		"Query":   query,
		"Results": results,
	}

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

	return c.Render(http.StatusOK, "public/pages/search.html", data)
}

// GET /search/suggest
func (h *SearchHandler) SearchSuggest(c echo.Context) error {
	query := c.QueryParam("q")

	var results []SearchResult
	if query != "" {
		results = h.search(query, 5)
	}

	data := map[string]interface{}{
		"Results": results,
	}

	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, "public/partials/search_suggestions.html", data, c); err != nil {
		h.logger.Error("search suggestions render failed", "error", err)
		return err
	}
	return c.HTML(http.StatusOK, buf.String())
}
