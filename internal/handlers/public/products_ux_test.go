package public_test

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestCatalogSearchWhitespaceAndRestorableFullPage(t *testing.T) {
	e, queries, cache, cleanup := setupProductsHandler(t)
	defer cleanup()
	e.Debug = true
	e.Renderer = templates.NewRenderer("../../../templates")
	h := public.NewProductsHandler(queries, logger, services.NewProductService(queries), cache)
	e.GET("/products/search", h.ProductSearch)
	cat := createTestCategory(t, queries, "Audio", "audio")
	createTestProduct(t, queries, "QA-1", "headphones", "QA Collaboration 2026", cat.ID, "published")
	createTestProduct(t, queries, "QA-2", "draft", "QA Collaboration 2026 Draft", cat.ID, "draft")
	get := func(path string, partial bool) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if partial {
			req.Header.Set("HX-Request", "true")
		}
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d: %s", path, rec.Code, rec.Body.String())
		}
		return rec
	}
	get("/products", false) // Populate generic-page cache first.
	for _, query := range []string{"Collaboration 2026", " Collaboration 2026 ", "\tCollaboration 2026\n"} {
		params := url.Values{"q": {query}}.Encode()
		fragment := get("/products/search?"+params, true)
		if fragment.Header().Get("HX-Replace-Url") != "/products?q=Collaboration+2026" {
			t.Fatalf("search history URL not canonical: %v", fragment.Header())
		}
		for _, path := range []string{"/products?" + params, "/products/search?" + params} {
			body := get(path, false).Body.String()
			for _, want := range []string{`value="Collaboration 2026"`, `href="/products/audio/headphones"`, "QA Collaboration 2026", `id="product-results"`} {
				if !strings.Contains(body, want) {
					t.Errorf("restored full page missing %q", want)
				}
			}
			if strings.Contains(body, "QA Collaboration 2026 Draft") {
				t.Error("draft leaked into search")
			}
		}
		if !strings.Contains(fragment.Body.String(), "QA Collaboration 2026") {
			t.Error("whitespace changed matches")
		}
	}
	if strings.Contains(get("/products", false).Body.String(), `href="/products/audio/headphones"`) {
		t.Error("search polluted generic page cache")
	}
	empty := get("/products/search?q=++", true)
	if strings.TrimSpace(empty.Body.String()) != "" || empty.Header().Get("HX-Replace-Url") != "/products" {
		t.Error("clearing query did not reset results and URL")
	}
	if !strings.Contains(get("/products?q=nomatch-fixture", false).Body.String(), "No products found") {
		t.Error("deep-linked empty result missing")
	}
}

func TestProductSpecificationsRenderUniqueLinkedStates(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()
	e.Debug = true
	e.Renderer = templates.NewRenderer("../../../templates")
	cat := createTestCategory(t, queries, "Audio", "audio")
	prod := createTestProduct(t, queries, "QA-SPEC", "qa-spec", "QA Spec", cat.ID, "published")
	var ids []int64
	for _, section := range []string{"Audio / USB", "Audio USB"} {
		spec, err := queries.CreateProductSpec(context.Background(), sqlc.CreateProductSpecParams{ProductID: prod.ID, SectionName: section, SpecKey: "Port", SpecValue: "USB"})
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, spec.ID)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/products/audio/qa-spec", nil))
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	for _, id := range ids {
		for _, want := range []string{fmt.Sprintf(`id="spec-toggle-%d" aria-expanded="false" aria-controls="spec-panel-%d"`, id, id), fmt.Sprintf(`id="spec-panel-%d" aria-labelledby="spec-toggle-%d"`, id, id)} {
			if !strings.Contains(rec.Body.String(), want) {
				t.Errorf("missing linked collapsed state %q", want)
			}
		}
	}
}

func TestProductQuoteLinkEscapesSKU(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()
	e.Debug = true
	e.Renderer = templates.NewRenderer("../../../templates")
	cat := createTestCategory(t, queries, "Audio", "audio")
	sku := "QA & USB+2026"
	createTestProduct(t, queries, sku, "quote-product", "Quote Product", cat.ID, "published")
	section, err := queries.GetPageSection(context.Background(), sqlc.GetPageSectionParams{PageKey: "product_detail", SectionKey: "cta"})
	if err != nil {
		t.Fatal(err)
	}
	err = queries.UpdatePageSection(context.Background(), sqlc.UpdatePageSectionParams{ID: section.ID, Heading: "Quote", PrimaryButtonText: "Request a Quote", PrimaryButtonUrl: "/contact?product={product_sku}", IsActive: true})
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/products/audio/quote-product", nil))
	if !strings.Contains(html.UnescapeString(rec.Body.String()), `href="/contact?product=`+url.QueryEscape(sku)+`"`) {
		t.Fatal("quote link did not preserve full SKU with reserved characters")
	}
}

func TestAudioCopyMigrationPreservesCustomDescriptions(t *testing.T) {
	db, _, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	old := "Deliver sharp video quality and smooth performance with our high-definition web cameras. Designed for online meetings, virtual classrooms, and live sessions, they provide crystal-clear visuals, auto light correction, and reliable plug-and-play connectivity."
	up, err := os.ReadFile("../../../db/migrations/046_audio_category_copy.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	down, err := os.ReadFile("../../../db/migrations/046_audio_category_copy.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO product_categories (name,slug,description,icon) VALUES ('Audio Test','audio-solutions',?,'test'),('Camera Test','camera-test',?,'test')`, old, old)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(up)); err != nil {
		t.Fatal(err)
	}
	var audio, camera string
	db.QueryRow(`SELECT description FROM product_categories WHERE slug='audio-solutions'`).Scan(&audio)
	db.QueryRow(`SELECT description FROM product_categories WHERE slug='camera-test'`).Scan(&camera)
	if audio == old || !strings.Contains(audio, "headphones") || camera != old {
		t.Fatal("migration scope/correction wrong")
	}
	if _, err = db.Exec(string(down)); err != nil {
		t.Fatal(err)
	}
	db.QueryRow(`SELECT description FROM product_categories WHERE slug='audio-solutions'`).Scan(&audio)
	if audio != old {
		t.Fatal("known text did not roll back")
	}
	db.Exec(`UPDATE product_categories SET description='Owner customized audio copy' WHERE slug='audio-solutions'`)
	if _, err = db.Exec(string(up)); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(down)); err != nil {
		t.Fatal(err)
	}
	db.QueryRow(`SELECT description FROM product_categories WHERE slug='audio-solutions'`).Scan(&audio)
	if audio != "Owner customized audio copy" {
		t.Fatal("migration overwrote customized owner copy")
	}
}

func TestProductDraftPrivacyPreviewAndUnpublish(t *testing.T) {
	e, queries, cache, cleanup := setupProductsHandler(t)
	defer cleanup()
	cat := createTestCategory(t, queries, "Audio", "audio")
	product := createTestProduct(t, queries, "QA-PREVIEW", "preview-check", "Preview Check", cat.ID, "draft")
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Test-only authenticated context; application session auth is unchanged.
			if c.Request().Header.Get("X-Test-Admin") == "yes" {
				c.Set("session", &middleware.Session{UserID: 1})
			}
			return next(c)
		}
	})
	e.POST("/admin/products/:id", admin.NewProductsHandler(queries, logger, nil, cache).Update)
	get := func(path string, authenticated bool, want int) {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if authenticated {
			req.Header.Set("X-Test-Admin", "yes")
		}
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Fatalf("GET %s (admin=%v): got %d want %d: %s", path, authenticated, rec.Code, want, rec.Body.String())
		}
	}
	get("/products/audio/preview-check", false, 404)
	get("/products/audio/preview-check?preview=true", false, 404)
	get("/products/audio/preview-check?preview=true", true, 200)
	get("/products/audio/preview-check", false, 404) // Preview cannot populate public cache.
	update := func(status string) {
		t.Helper()
		values := url.Values{"name": {"Preview Check"}, "sku": {"QA-PREVIEW"}, "slug": {"preview-check"}, "category_id": {fmt.Sprint(cat.ID)}, "status": {status}, "description": {"Synthetic fixture"}}
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d", product.ID), strings.NewReader(values.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != 303 {
			t.Fatalf("status update: %d: %s", rec.Code, rec.Body.String())
		}
	}
	update("published")
	get("/products/audio/preview-check", false, 200) // Populate published-page cache.
	update("draft")
	get("/products/audio/preview-check", false, 404) // Real admin update must invalidate it.
	get("/products/audio/preview-check?preview=true", true, 200)
}

func TestProductOverviewRichMarkupAndPlainDescriptionFallback(t *testing.T) {
	e, queries, cache, cleanup := setupProductsHandler(t)
	defer cleanup()
	e.Renderer = templates.NewRenderer("../../../templates")
	cat := createTestCategory(t, queries, "Audio", "audio")
	product := createTestProduct(t, queries, "QA-RICH", "rich-overview", "Rich Overview", cat.ID, "published")
	rich := `<p>Compare <strong>clear audio</strong> and <em>comfort</em>.</p><ul><li>USB connection</li><li>Microphone</li></ul><p><a href="/contact">Ask a question</a></p>`
	update := func(overview sql.NullString, description string) {
		t.Helper()
		if err := queries.UpdateProduct(context.Background(), sqlc.UpdateProductParams{ID: product.ID, Sku: product.Sku, Slug: product.Slug, Name: product.Name, CategoryID: cat.ID, Status: "published", Overview: overview, Description: description}); err != nil {
			t.Fatal(err)
		}
		cache.DeleteByPrefix("page:products")
	}
	render := func() string {
		t.Helper()
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/products/audio/rich-overview", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("render: %d: %s", rec.Code, rec.Body.String())
		}
		return rec.Body.String()
	}
	update(sql.NullString{String: rich, Valid: true}, "Plain summary")
	body := render()
	// A block container preserves valid paragraph/list structure from the rich editor.
	if !strings.Contains(body, `<div class="product-overview text-sm leading-relaxed opacity-80">`+rich+`</div>`) {
		t.Fatal("rich overview lost its formatting or block container")
	}
	if strings.Contains(body, `&lt;strong&gt;clear audio`) || strings.Contains(body, `&lt;ul&gt;`) {
		t.Fatal("overview displayed literal HTML tags")
	}
	update(sql.NullString{}, `A plain <strong>summary</strong>`)
	body = render()
	if !strings.Contains(body, `<p class="text-sm leading-relaxed opacity-80">A plain &lt;strong&gt;summary&lt;/strong&gt;</p>`) {
		t.Fatal("description fallback must remain escaped plain text")
	}
}
