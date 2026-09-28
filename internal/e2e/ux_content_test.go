package e2e_test

import (
	"context"
	"database/sql"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
)

func contentRequest(app *echo.Echo, path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Accept", "text/html")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	return response
}

func TestUXContactQuoteContextAndPrivacy(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()
	app.Renderer = templates.NewRenderer("templates")
	app.Use(customMiddleware.SettingsLoader(queries))
	contact := publicHandlers.NewContactHandler(queries, testLogger, services.NewCache())
	app.GET("/privacy", contact.ShowPrivacyNotice)
	ctx := context.Background()
	category, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{Name: "Quote products", Slug: "quote-products", Description: "Test", Icon: "screen"})
	if err != nil {
		t.Fatal(err)
	}
	for _, product := range []sqlc.CreateProductParams{
		{Sku: "TEST A&B", Slug: "quote-a", Name: "Display <A>", Description: "Test", CategoryID: category.ID, Status: "published"},
		{Sku: "TEST-B", Slug: "quote-b", Name: "Display B", Description: "Test", CategoryID: category.ID, Status: "published"},
		{Sku: "DRAFT", Slug: "draft", Name: "Hidden product", Description: "Test", CategoryID: category.ID, Status: "draft"},
	} {
		if _, err := queries.CreateProduct(ctx, product); err != nil {
			t.Fatal(err)
		}
	}
	// Warm the generic cache before exercising distinct quote contexts.
	generic := contentRequest(app, "/contact")
	if generic.Code != 200 || strings.Contains(generic.Body.String(), `id="quote-context"`) {
		t.Fatal("generic form should not have quote context")
	}
	expected := "Please send me a quote for Display <A> (SKU: TEST A&B)."
	quote := contentRequest(app, "/contact?product="+url.QueryEscape("TEST A&B"))
	if quote.Code != 200 || !strings.Contains(quote.Body.String(), html.EscapeString(expected)) || !strings.Contains(quote.Body.String(), `value="sales" selected`) {
		t.Fatalf("quote missing context: %d %s", quote.Code, quote.Body.String())
	}
	if strings.Contains(quote.Body.String(), "Display <A>") {
		t.Fatal("product name must be escaped")
	}
	second := contentRequest(app, "/contact?product=TEST-B")
	if !strings.Contains(second.Body.String(), "Display B (SKU: TEST-B)") || strings.Contains(second.Body.String(), "TEST A&amp;B") {
		t.Fatal("product context leaked through cache")
	}
	for _, path := range []string{"/contact?product=DRAFT", "/contact?product=missing"} {
		response := contentRequest(app, path)
		if !strings.Contains(response.Body.String(), "This product is no longer available here") || strings.Contains(response.Body.String(), "Hidden product") {
			t.Fatal("unpublished/unknown SKU exposed or not explained")
		}
	}
	if strings.Contains(contentRequest(app, "/contact").Body.String(), `id="quote-context"`) {
		t.Fatal("generic cache contaminated by quote")
	}
	form := url.Values{"name": {"UX quote test"}, "email": {"quote@example.test"}, "phone": {"+1 202 555 0147"}, "company": {"Example Test"}, "inquiry_type": {"sales"}, "message": {expected + " Please include installation."}}
	request := httptest.NewRequest("POST", "/contact/submit", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("submit: %d %s", response.Code, response.Body.String())
	}
	rows, err := queries.ListContactSubmissions(ctx, sqlc.ListContactSubmissionsParams{Limit: 10})
	if err != nil || len(rows) != 1 || rows[0].Message != form.Get("message") {
		t.Fatalf("edited quote not saved: %v %v", rows, err)
	}
	privacy := contentRequest(app, "/privacy")
	for _, required := range []string{"Privacy Notice", "IP address and browser information", "review and respond", `href="/contact"`} {
		if privacy.Code != 200 || !strings.Contains(privacy.Body.String(), required) {
			t.Fatalf("privacy page missing %q", required)
		}
	}
	if !strings.Contains(generic.Body.String(), `href="/privacy"`) || !strings.Contains(generic.Body.String(), "We use the details you submit") {
		t.Fatal("form privacy explanation missing")
	}
}

func TestUXPublic404RecoveryPreservesOtherResponses(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()
	app.Renderer = templates.NewRenderer("templates")
	app.HTTPErrorHandler = publicHandlers.NewPublicHTTPErrorHandler(app, queries)
	for _, path := range []string{"/missing-page", "/blog/unknown-article", "/products/missing/missing"} {
		response := contentRequest(app, path)
		if response.Code != 404 || !strings.Contains(response.Header().Get("Content-Type"), "text/html") {
			t.Fatalf("%s: expected HTML404, got %d %s", path, response.Code, response.Body.String())
		}
		for _, text := range []string{"Page not found", `aria-label="Page recovery"`, `action="/search"`, `name="robots" content="noindex"`, `href="/privacy"`} {
			if !strings.Contains(response.Body.String(), text) {
				t.Fatalf("%s missing %q", path, text)
			}
		}
	}
	for _, test := range []struct{ path, method, accept, hx string }{
		{"/api/missing", "GET", "text/html", ""}, {"/admin-missing", "GET", "application/json", ""},
		{"/admin/missing", "GET", "text/html", ""}, {"/public/missing.css", "GET", "text/html", ""},
		{"/uploads/missing.png", "GET", "text/html", ""}, {"/missing", "GET", "application/json", ""},
		{"/missing", "GET", "text/html", "true"}, {"/missing", "POST", "text/html", ""},
	} {
		request := httptest.NewRequest(test.method, test.path, nil)
		request.Header.Set("Accept", test.accept)
		request.Header.Set("HX-Request", test.hx)
		response := httptest.NewRecorder()
		app.ServeHTTP(response, request)
		if strings.Contains(response.Body.String(), "Page not found") {
			t.Fatalf("unexpected full page for %+v", test)
		}
	}
	request := httptest.NewRequest("HEAD", "/missing", nil)
	request.Header.Set("Accept", "text/html")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != 404 || response.Body.Len() != 0 {
		t.Fatal("HEAD must preserve404 with no body")
	}
}

func TestUXBlogFilterResultsAndCertificationVisibility(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()
	app.Renderer = templates.NewRenderer("templates")
	ctx := context.Background()
	category, err := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{Name: "Technology", Slug: "technology", ColorHex: "#0066CC"})
	if err != nil {
		t.Fatal(err)
	}
	author, err := queries.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{Name: "Test Author", Slug: "test-author", Title: "Editor"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{Title: "Example article", Slug: "example-article", Excerpt: "Example excerpt", Body: "<h2>Meaningful section</h2><p>Example body</p>", PublishedAt: sql.NullTime{Time: time.Now(), Valid: true}, CategoryID: category.ID, AuthorID: author.ID, Status: "published"})
	if err != nil {
		t.Fatal(err)
	}
	response := contentRequest(app, "/blog?category=technology")
	for _, text := range []string{`id="blog-results"`, `role="status"`, "Technology: 1 post", `/blog?category=technology#blog-results`, `aria-current="page"`, `/blog#blog-results`} {
		if response.Code != 200 || !strings.Contains(response.Body.String(), text) {
			t.Fatalf("filter missing %q: %d", text, response.Code)
		}
	}
	empty := contentRequest(app, "/blog?category=does-not-exist")
	if !strings.Contains(empty.Body.String(), "Unknown category: 0 posts") {
		t.Fatal("empty/unknown results not explained")
	}
	article := contentRequest(app, "/blog/example-article")
	if article.Code != 200 || !strings.Contains(article.Body.String(), `class="blog-article-body`) || !strings.Contains(article.Body.String(), "<h2>Meaningful section</h2>") {
		t.Fatal("article headings not semantically rendered/styled")
	}
	about := contentRequest(app, "/about")
	if about.Code != 200 || strings.Contains(about.Body.String(), "Certifications &amp; Compliance") || strings.Contains(about.Body.String(), "No certifications configured") {
		t.Fatal("empty certification section visible")
	}
	_, err = queries.CreateCertification(ctx, sqlc.CreateCertificationParams{Name: "Synthetic credential", Abbreviation: "TEST", Description: sql.NullString{String: "Local test only", Valid: true}})
	if err != nil {
		t.Fatal(err)
	}
	// New handler with an empty cache simulates the admin invalidation path.
	aboutHandler := publicHandlers.NewAboutHandler(queries, testLogger, services.NewCache())
	app.GET("/about-with-cert", aboutHandler.AboutPage)
	populated := contentRequest(app, "/about-with-cert")
	if populated.Code != 200 || !strings.Contains(populated.Body.String(), "Certifications &amp; Compliance") || !strings.Contains(populated.Body.String(), "Synthetic credential") {
		t.Fatal("configured certification hidden")
	}
}

func TestUXArticleHeadingMigrationPreservesOtherContent(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err = db.Exec(`CREATE TABLE blog_posts(slug TEXT PRIMARY KEY,body TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	slug := "the-invisible-engine-why-ops-modules-are-crucial-for-enterprise-interactive-displays"
	original := `<div>What is Open Pluggable Specification (OPS)?<br><br></div><div>Keep <a href="/products">this link</a>.</div><div>1. Seamless Lifecycle Management &amp; Future-Proofing</div><div>Architecture Built for Performance<br><br></div><div>Edited section</div>`
	for _, key := range []string{slug, "another-post"} {
		if _, err = db.Exec("INSERT INTO blog_posts VALUES (?,?)", key, original); err != nil {
			t.Fatal(err)
		}
	}
	migration, err := os.ReadFile("db/migrations/045_ops_article_heading_structure.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	var changed, untouched string
	if err = db.QueryRow("SELECT body FROM blog_posts WHERE slug=?", slug).Scan(&changed); err != nil {
		t.Fatal(err)
	}
	if err = db.QueryRow("SELECT body FROM blog_posts WHERE slug='another-post'").Scan(&untouched); err != nil {
		t.Fatal(err)
	}
	if untouched != original {
		t.Fatal("unrelated article modified")
	}
	for _, part := range []string{"<h2>What is Open Pluggable Specification (OPS)?</h2>", "<h3>1. Seamless Lifecycle Management &amp; Future-Proofing</h3>", `<div>Keep <a href="/products">this link</a>.</div>`, "<div>Edited section</div>"} {
		if !strings.Contains(changed, part) {
			t.Fatalf("lost expected content %q", part)
		}
	}
	if _, err = db.Exec(string(migration)); err != nil {
		t.Fatal(err)
	}
	var repeated string
	db.QueryRow("SELECT body FROM blog_posts WHERE slug=?", slug).Scan(&repeated)
	if repeated != changed {
		t.Fatal("migration not idempotent")
	}
}
