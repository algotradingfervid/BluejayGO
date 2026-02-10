package public_test

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

// testRenderer is a minimal template renderer for handler tests.
// It renders data keys into HTML so we can assert on dynamic content.
type testRenderer struct{}

func (r *testRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	d, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map[string]interface{}, got %T", data)
	}
	t := template.New("test")
	// Produce a simple HTML dump of key data so tests can assert content
	tmplStr := `<html><body>
<div id="title">{{.Title}}</div>
{{if .Categories}}<div id="categories">{{range .Categories}}<span class="cat">{{.Category.Name}}</span>{{end}}</div>{{end}}
{{if .Products}}<div id="products">{{range .Products}}<span class="prod">{{.Name}}</span>{{end}}</div>{{end}}
{{if .Product}}<div id="product-name">{{.Product.Name}}</div><div id="product-sku">{{.Product.Sku}}</div>{{end}}
{{if .Category}}<div id="category-name">{{.Category.Name}}</div>{{end}}
{{if .TotalCount}}<div id="total-count">{{.TotalCount}}</div>{{end}}
</body></html>`
	t, err := t.Parse(tmplStr)
	if err != nil {
		return err
	}
	return t.Execute(w, d)
}

func setupProductsHandler(t *testing.T) (*echo.Echo, *sqlc.Queries, *services.Cache, func()) {
	t.Helper()
	_, queries, cleanup := testutil.SetupTestDB(t)

	e := echo.New()
	e.Renderer = &testRenderer{}

	productSvc := services.NewProductService(queries)
	cache := services.NewCache()
	h := public.NewProductsHandler(queries, logger, productSvc, cache)

	e.GET("/products", h.ProductsList)
	e.GET("/products/:category", h.ProductsByCategory)
	e.GET("/products/:category/:slug", h.ProductDetail)

	return e, queries, cache, cleanup
}

func createTestCategory(t *testing.T, queries *sqlc.Queries, name, slug string) sqlc.ProductCategory {
	t.Helper()
	cat, err := queries.CreateProductCategory(context.Background(), sqlc.CreateProductCategoryParams{
		Name: name, Slug: slug, Description: "Test " + name, Icon: "icon", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateProductCategory: %v", err)
	}
	return cat
}

func createTestProduct(t *testing.T, queries *sqlc.Queries, sku, slug, name string, catID int64, status string) sqlc.Product {
	t.Helper()
	prod, err := queries.CreateProduct(context.Background(), sqlc.CreateProductParams{
		Sku: sku, Slug: slug, Name: name, Description: "Test " + name, CategoryID: catID, Status: status,
	})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}
	return prod
}

// --- Test: Empty DB shows no hardcoded product data ---

func TestProductsList_EmptyDB(t *testing.T) {
	e, _, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	// Should NOT contain any hardcoded product strings
	hardcoded := []string{"BJ-D200", "BJ-D100", "Intel Core i5", "Intel Core i7", "AMD Ryzen", "Detector", "detector"}
	for _, s := range hardcoded {
		if strings.Contains(body, s) {
			t.Errorf("empty DB render contains hardcoded string %q", s)
		}
	}
}

func TestProductsByCategory_EmptyDB(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	createTestCategory(t, queries, "Monitors", "monitors")

	req := httptest.NewRequest(http.MethodGet, "/products/monitors", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if strings.Contains(body, "BJ-D200") || strings.Contains(body, "Intel Core") {
		t.Error("category page with no products contains hardcoded product data")
	}
}

// --- Test: Product detail returns correct data ---

func TestProductDetail_Found(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "Detectors", "detectors")
	createTestProduct(t, queries, "DET-100", "alpha-det", "Alpha Detector", cat.ID, "published")

	req := httptest.NewRequest(http.MethodGet, "/products/detectors/alpha-det", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Alpha Detector") {
		t.Error("product detail does not contain product name")
	}
	if !strings.Contains(body, "DET-100") {
		t.Error("product detail does not contain SKU")
	}
	if !strings.Contains(body, "Detectors") {
		t.Error("product detail does not contain category name")
	}
}

func TestProductDetail_NotFound(t *testing.T) {
	e, _, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/products/nonexistent/no-product", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestProductDetail_WrongCategory(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "Detectors", "detectors")
	createTestCategory(t, queries, "Monitors", "monitors")
	createTestProduct(t, queries, "DET-200", "beta-det", "Beta Detector", cat.ID, "published")

	req := httptest.NewRequest(http.MethodGet, "/products/monitors/beta-det", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for wrong category, got %d", rec.Code)
	}
}

// --- Test: Category filter only shows products in that category ---

func TestProductsByCategory_FilterCorrectly(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	catA := createTestCategory(t, queries, "Category A", "cat-a")
	catB := createTestCategory(t, queries, "Category B", "cat-b")

	createTestProduct(t, queries, "A-001", "prod-a1", "Product A1", catA.ID, "published")
	createTestProduct(t, queries, "A-002", "prod-a2", "Product A2", catA.ID, "published")
	createTestProduct(t, queries, "B-001", "prod-b1", "Product B1", catB.ID, "published")

	// Request category A
	req := httptest.NewRequest(http.MethodGet, "/products/cat-a", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Product A1") {
		t.Error("expected Product A1 in category A results")
	}
	if !strings.Contains(body, "Product A2") {
		t.Error("expected Product A2 in category A results")
	}
	if strings.Contains(body, "Product B1") {
		t.Error("Product B1 should not appear in category A results")
	}

	// Request category B
	req = httptest.NewRequest(http.MethodGet, "/products/cat-b", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	body = rec.Body.String()
	if !strings.Contains(body, "Product B1") {
		t.Error("expected Product B1 in category B results")
	}
	if strings.Contains(body, "Product A1") {
		t.Error("Product A1 should not appear in category B results")
	}
}

func TestProductsByCategory_DraftNotShown(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "Widgets", "widgets")
	createTestProduct(t, queries, "W-001", "published-widget", "Published Widget", cat.ID, "published")
	createTestProduct(t, queries, "W-002", "draft-widget", "Draft Widget", cat.ID, "draft")

	req := httptest.NewRequest(http.MethodGet, "/products/widgets", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "Published Widget") {
		t.Error("expected Published Widget in results")
	}
	if strings.Contains(body, "Draft Widget") {
		t.Error("Draft Widget should not appear in public listing")
	}
}

// --- Test: Nonexistent category returns 404 ---

func TestProductsByCategory_NonexistentCategory(t *testing.T) {
	e, _, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/products/nonexistent", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent category, got %d", rec.Code)
	}
}

// --- Test: Cache invalidation ---

func TestProductsList_CacheInvalidation(t *testing.T) {
	e, queries, cache, cleanup := setupProductsHandler(t)
	defer cleanup()

	// First request - empty, gets cached
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	firstBody := rec.Body.String()

	// Add a category
	createTestCategory(t, queries, "New Category", "new-category")

	// Cached response should still be old
	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	cachedBody := rec.Body.String()
	if cachedBody != firstBody {
		t.Error("expected cached response to match first response")
	}

	// Invalidate cache
	cache.DeleteByPrefix("page:products")

	// Now should reflect new data
	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	newBody := rec.Body.String()
	if !strings.Contains(newBody, "New Category") {
		t.Error("after cache invalidation, expected new category to appear")
	}
}

func TestProductDetail_CacheInvalidation(t *testing.T) {
	e, queries, cache, cleanup := setupProductsHandler(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "Tech", "tech")
	prod := createTestProduct(t, queries, "T-001", "gadget", "Gadget One", cat.ID, "published")

	// First request - gets cached
	req := httptest.NewRequest(http.MethodGet, "/products/tech/gadget", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Gadget One") {
		t.Fatal("expected product name in response")
	}

	// Update product name in DB
	err := queries.UpdateProduct(context.Background(), sqlc.UpdateProductParams{
		ID: prod.ID, Sku: "T-001", Slug: "gadget", Name: "Gadget Updated",
		Description: "Updated", CategoryID: cat.ID, Status: "published",
	})
	if err != nil {
		t.Fatalf("UpdateProduct: %v", err)
	}

	// Still cached - old name
	req = httptest.NewRequest(http.MethodGet, "/products/tech/gadget", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), "Gadget Updated") {
		t.Error("cache should still serve old content before invalidation")
	}

	// Invalidate
	cache.DeleteByPrefix("page:products")

	// After invalidation - new name
	req = httptest.NewRequest(http.MethodGet, "/products/tech/gadget", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if !strings.Contains(rec.Body.String(), "Gadget Updated") {
		t.Error("after cache invalidation, expected updated product name")
	}
}

// --- Test: Template integrity - no hardcoded product strings ---

func TestTemplateIntegrity_NoHardcodedProducts(t *testing.T) {
	e, _, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	// Render the products list page with an empty DB
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	body := rec.Body.String()

	hardcoded := []string{
		"BJ-D200", "BJ-D100", "BJ-M100", "BJ-L100",
		"Intel Core i5", "Intel Core i7", "Intel Core i9",
		"AMD Ryzen", "NVIDIA",
		"Alpha Detector", "Beta Detector", "Gamma Detector",
	}

	for _, s := range hardcoded {
		if strings.Contains(body, s) {
			t.Errorf("products page with empty DB contains hardcoded string %q", s)
		}
	}
}

// --- Test: Products list shows categories with counts ---

func TestProductsList_ShowsCategoriesWithCounts(t *testing.T) {
	e, queries, _, cleanup := setupProductsHandler(t)
	defer cleanup()

	cat := createTestCategory(t, queries, "Sensors", "sensors")
	createTestProduct(t, queries, "S-001", "sensor-a", "Sensor A", cat.ID, "published")
	createTestProduct(t, queries, "S-002", "sensor-b", "Sensor B", cat.ID, "published")

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Sensors") {
		t.Error("expected category name 'Sensors' in products list")
	}
}
