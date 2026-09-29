package public_test

import (
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/internal/database"
	"github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/templates"
)

func TestSearchSuggestionsContinueThroughEveryPublishedResult(t *testing.T) {
	db, err := database.InitDB(database.Config{Path: filepath.Join(t.TempDir(), "search.db")})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := database.RunMigrations(db, "../../../db/migrations"); err != nil {
		t.Fatal(err)
	}
	exec := func(query string, args ...interface{}) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`INSERT INTO product_categories(name,slug,description,icon) VALUES ('Paged search','paged-search','','')`)
	exec(`INSERT INTO blog_categories(name,slug,color_hex) VALUES ('Paged search','paged-search','#000000')`)
	exec(`INSERT INTO blog_authors(name,slug,title) VALUES ('Paged author','paged-author','Editor')`)
	for i := 0; i < 32; i++ {
		status := "published"
		if i == 31 {
			status = "draft"
		}
		exec(`INSERT INTO products(sku,slug,name,description,category_id,status) VALUES (?,?,?,'Café & audio', (SELECT id FROM product_categories WHERE slug='paged-search'),?)`, fmt.Sprintf("PAGED-%02d", i), fmt.Sprintf("paged-%02d", i), fmt.Sprintf("RoundTwo Product %02d", i), status)
	}
	for i := 0; i < 13; i++ {
		status := "published"
		if i == 12 {
			status = "draft"
		}
		exec(`INSERT INTO blog_posts(title,slug,excerpt,body,category_id,author_id,status) VALUES (?,?,'Excerpt','Body',(SELECT id FROM blog_categories WHERE slug='paged-search'),(SELECT id FROM blog_authors WHERE slug='paged-author'),?)`, fmt.Sprintf("RoundTwo Article %02d", i), fmt.Sprintf("paged-%02d", i), status)
	}
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	h := public.NewSearchHandler(db, slog.New(slog.NewTextHandler(io.Discard, nil)))
	e.GET("/search", h.SearchPage)
	e.GET("/search/suggest", h.SearchSuggest)
	get := func(path string) string {
		t.Helper()
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d: %s", path, rec.Code, rec.Body.String())
		}
		return rec.Body.String()
	}
	resultLinks := regexp.MustCompile(`href="(/(?:products/paged-search|blog)/paged-[0-9]+)"`)
	suggestions := get("/search/suggest?q=RoundTwo")
	if len(resultLinks.FindAllStringSubmatch(suggestions, -1)) != 10 || !strings.Contains(suggestions, "10 search suggestions") || strings.Contains(suggestions, "matching results") {
		t.Fatal("suggestions must identify their capped count while retaining both content types")
	}
	if !strings.Contains(suggestions, `href="/search?q=RoundTwo"`) || !strings.Contains(suggestions, "View all results") {
		t.Fatal("missing continuation")
	}
	nextLink := regexp.MustCompile(`href="([^"]+)" rel="next"`)
	current := "/search?q=RoundTwo"
	seen := map[string]bool{}
	for page, wantCount := range []int{20, 20, 3} {
		body := get(current)
		matches := resultLinks.FindAllStringSubmatch(body, -1)
		if len(matches) != wantCount {
			t.Fatalf("page %d: got %d results, want %d", page+1, len(matches), wantCount)
		}
		for _, match := range matches {
			if seen[match[1]] {
				t.Fatalf("duplicate across pages: %s", match[1])
			}
			seen[match[1]] = true
		}
		if page > 0 && !strings.Contains(body, `rel="prev"`) {
			t.Fatal("missing previous page")
		}
		next := nextLink.FindStringSubmatch(body)
		if page < 2 {
			if len(next) != 2 {
				t.Fatal("missing next page")
			}
			current = html.UnescapeString(next[1])
		} else if len(next) != 0 {
			t.Fatal("last page should not offer an empty next page")
		}
	}
	for i := 0; i < 31; i++ {
		if !seen[fmt.Sprintf("/products/paged-search/paged-%02d", i)] {
			t.Fatalf("missing product %d", i)
		}
	}
	for i := 0; i < 12; i++ {
		if !seen[fmt.Sprintf("/blog/paged-%02d", i)] {
			t.Fatalf("missing article %d", i)
		}
	}
	if seen["/products/paged-search/paged-31"] || seen["/blog/paged-12"] {
		t.Fatal("draft appeared in public search")
	}
	for _, value := range []string{"0", "-1", "invalid", "999999999999999999999999999"} {
		if got := len(resultLinks.FindAllStringSubmatch(get("/search?q=RoundTwo&page="+url.QueryEscape(value)), -1)); got != 20 {
			t.Fatalf("invalid page %q: got %d results", value, got)
		}
	}
	if body := get("/search?q=RoundTwo&page=999"); !strings.Contains(body, "No results on this page") || !strings.Contains(body, `href="/search?q=RoundTwo"`) {
		t.Fatal("out-of-range page has no recovery")
	}
	if body := get("/search?q=++&page=2"); strings.Contains(body, "No results") || strings.Contains(body, "Next page") {
		t.Fatal("empty query should show only the search form")
	}
	// Prefix matching is identical in suggestions and full results; punctuation
	// is preserved in continuation URLs rather than becoming URL syntax.
	for _, query := range []string{"RoundTwo+", `"RoundTwo"`, "Café", "RoundTwo &"} {
		body := get("/search/suggest?q=" + url.QueryEscape(query))
		want := "/search?" + url.Values{"q": {query}}.Encode()
		if !strings.Contains(html.UnescapeString(body), `href="`+want+`"`) {
			t.Fatalf("continuation lost query %q: %s", query, body)
		}
	}
	unsafe := `<script>alert(1)</script> & "quoted" + 東京`
	for _, path := range []string{"/search", "/search/suggest"} {
		body := get(path + "?q=" + url.QueryEscape(unsafe))
		if strings.Contains(body, "<script>alert(1)</script>") || !strings.Contains(body, "&lt;script&gt;") {
			t.Fatal("query must remain escaped text")
		}
	}
}
