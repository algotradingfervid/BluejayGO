package admin_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
	"golang.org/x/net/html"
)

func TestEditorialBlogActionsPersistPromisedState(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	cat, err := q.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{Name: "Test", Slug: "editorial", ColorHex: "#000000"})
	if err != nil {
		t.Fatal(err)
	}
	author, err := q.CreateBlogAuthor(ctx, sqlc.CreateBlogAuthorParams{Name: "Test", Slug: "editorial", Title: "Editor"})
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	h := admin.NewBlogPostsHandler(q, logger, services.NewCache())
	f := url.Values{"title": {"Editorial"}, "slug": {"editorial"}, "body": {"<p>Saved <strong>story</strong></p>"}, "category_id": {strconv.FormatInt(cat.ID, 10)}, "author_id": {strconv.FormatInt(author.ID, 10)}, "status": {"draft"}, "submit_action": {"publish"}}
	rec, c := postForm(e, "/admin/blog/posts", f)
	if err := h.Create(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create status %d", rec.Code)
	}
	published, err := q.GetPublishedPostBySlug(ctx, "editorial")
	if err != nil {
		t.Fatalf("Publish did not make article public: %v", err)
	}
	post, err := q.GetBlogPost(ctx, published.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !post.PublishedAt.Valid || post.Body != f.Get("body") {
		t.Fatal("publish lost body/date")
	}
	id := strconv.FormatInt(post.ID, 10)
	for _, action := range []string{"draft", "publish"} {
		// Deliberately conflicting stale status proves the button wins in both directions.
		f.Set("status", post.Status)
		f.Set("submit_action", action)
		f.Set("body", "<p>Updated "+action+"</p>")
		f.Set("published_at", "2030-01-02T12:30")
		rec, c = postForm(e, "/admin/blog/posts/"+id, f)
		c.SetParamNames("id")
		c.SetParamValues(id)
		if err := h.Update(c); err != nil {
			t.Fatal(err)
		}
		if rec.Header().Get("Location") != "/admin/blog/posts/"+id+"/edit?saved=1" {
			t.Fatal("missing state confirmation redirect")
		}
		post, err = q.GetBlogPost(ctx, post.ID)
		if err != nil {
			t.Fatal(err)
		}
		want := "draft"
		if action == "publish" {
			want = "published"
		}
		if post.PublishedAt.Time.Format("2006-01-02T15:04") != "2030-01-02T12:30" {
			t.Fatal("article date edit did not persist")
		}
		if post.Status != want || post.Body != f.Get("body") {
			t.Fatalf("%s yielded %s or lost body", action, post.Status)
		}
		_, publicErr := q.GetPublishedPostBySlug(ctx, "editorial")
		if action == "draft" && publicErr != sql.ErrNoRows {
			t.Fatal("draft remained public")
		}
		if action == "publish" && publicErr != nil {
			t.Fatal(publicErr)
		}
		rec = httptest.NewRecorder()
		c = e.NewContext(httptest.NewRequest("GET", "/admin/blog/posts/"+id+"/edit?saved=1", nil), rec)
		c.SetParamNames("id")
		c.SetParamValues(id)
		if err := h.Edit(c); err != nil {
			t.Fatal(err)
		}
		wantCopy := "Draft saved."
		if action == "publish" {
			wantCopy = "Post published."
		}
		if !strings.Contains(rec.Body.String(), wantCopy) {
			t.Fatal("missing truthful feedback")
		}
	}
}

func TestEditorialProductConflictPreservesValuesAndSupportsCorrection(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	cat, err := q.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{Name: "Test", Slug: "test", Description: "Test", Icon: "test"})
	if err != nil {
		t.Fatal(err)
	}
	original, err := q.CreateProduct(ctx, sqlc.CreateProductParams{Name: "Existing", Slug: "existing", Sku: "USED", Description: "Original", Status: "draft", CategoryID: cat.ID})
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	h := admin.NewProductsHandler(q, logger, services.NewUploadService(t.TempDir()), services.NewCache())
	f := url.Values{"name": {"Preserve my product"}, "sku": {"USED"}, "slug": {"custom-url"}, "category_id": {strconv.FormatInt(cat.ID, 10)}, "status": {"draft"}, "description": {"Preserve description"}, "overview": {"<p>Preserve <strong>formatting</strong></p>"}, "tagline": {"Keep tagline"}, "meta_title": {"Keep SEO"}, "meta_description": {"Keep description"}, "is_featured": {"1"}, "featured_order": {"4"}, "video_url": {"https://example.com/video"}}
	rec, c := postForm(e, "/admin/products", f)
	if err := h.Create(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 422 {
		t.Fatalf("conflict status %d", rec.Code)
	}
	for _, want := range []string{"SKU already exists", `aria-invalid="true"`, "Preserve my product", "Preserve description", "Keep tagline", "Keep SEO", "Keep description", "custom-url", "Create Product", "/admin/products/" + strconv.FormatInt(original.ID, 10) + "/edit"} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("error form lost %q", want)
		}
	}
	if strings.Contains(rec.Body.String(), "/admin/products/0/") {
		t.Fatal("new record exposed edit subpages")
	}
	f.Set("sku", "UNIQUE")
	for _, invalid := range []string{"demo/model", "demo?variant", "demo#section"} {
		f.Set("slug", invalid)
		rec, c = postForm(e, "/admin/products", f)
		if err := h.Create(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != 422 || !strings.Contains(rec.Body.String(), "Slug must start") || !strings.Contains(rec.Body.String(), invalid) {
			t.Fatal("invalid slug not rejected with values retained")
		}
	}
	f.Set("slug", "custom-url")
	rec, c = postForm(e, "/admin/products", f)
	if err := h.Create(c); err != nil {
		t.Fatal(err)
	}
	saved, err := q.GetProductBySKU(ctx, "UNIQUE")
	if err != nil {
		t.Fatal(err)
	}
	if saved.Slug != "custom-url" || saved.Overview.String != f.Get("overview") {
		t.Fatal("correction lost entered URL/content")
	}
	id := strconv.FormatInt(saved.ID, 10)
	f.Set("sku", "USED")
	f.Set("name", "Pending update")
	rec, c = postForm(e, "/admin/products/"+id, f)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if err := h.Update(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 422 || !strings.Contains(rec.Body.String(), "Pending update") {
		t.Fatal("update conflict did not retain form")
	}
	unchanged, _ := q.GetProduct(ctx, saved.ID)
	if unchanged.Name != saved.Name {
		t.Fatal("conflict partially changed record")
	}
	f.Set("sku", "UNIQUE")
	rec, c = postForm(e, "/admin/products/"+id, f)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if err := h.Update(c); err != nil {
		t.Fatal(err)
	}
	updated, _ := q.GetProduct(ctx, saved.ID)
	if updated.Name != "Pending update" || updated.Slug != "custom-url" {
		t.Fatal("valid update failed or changed URL")
	}
}

func TestEditorialProductFallbackAndAccessibleFields(t *testing.T) {
	_, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	h := admin.NewProductsHandler(q, logger, nil, services.NewCache())
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest("GET", "/admin/products/new", nil), rec)
	if err := h.New(c); err != nil {
		t.Fatal(err)
	}
	doc, err := html.Parse(strings.NewReader(rec.Body.String()))
	if err != nil {
		t.Fatal(err)
	}
	labels := map[string]bool{}
	var controls []*html.Node
	var visit func(*html.Node)
	attr := func(n *html.Node, k string) string {
		for _, a := range n.Attr {
			if a.Key == k {
				return a.Val
			}
		}
		return ""
	}
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "label" {
				labels[attr(n, "for")] = true
			}
			if (n.Data == "input" || n.Data == "textarea" || n.Data == "select") && attr(n, "name") != "" && attr(n, "type") != "hidden" {
				controls = append(controls, n)
			}
		}
		for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
			visit(ch)
		}
	}
	visit(doc)
	fallback := false
	for _, n := range controls {
		if !labels[attr(n, "id")] && attr(n, "aria-label") == "" {
			t.Errorf("unnamed control %s", attr(n, "name"))
		}
		if attr(n, "name") == "overview" {
			fallback = n.Data == "textarea"
		}
	}
	if !fallback {
		t.Fatal("Overview lacks native named fallback")
	}
	for _, asset := range []string{"/public/css/trix.css", "/public/js/vendor/trix.js", "/public/js/admin-editors.js"} {
		if !strings.Contains(rec.Body.String(), asset) {
			t.Fatalf("missing local asset %s", asset)
		}
	}
}
