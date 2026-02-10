package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

var (
	testLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
)

// stubRenderer is a no-op renderer for e2e tests that don't need real templates.
type stubRenderer struct{}

func (r *stubRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	_, err := w.Write([]byte("<html>stub</html>"))
	return err
}

func setupApp(t *testing.T) (*echo.Echo, *sqlc.Queries, func()) {
	t.Helper()

	db, queries, cleanup := testutil.SetupTestDB(t)

	customMiddleware.InitSessionStore("e2e-test-secret-at-least-32-characters-long")

	e := echo.New()
	e.HideBanner = true
	e.Renderer = &stubRenderer{}

	e.Use(customMiddleware.SessionMiddleware())

	productSvc := services.NewProductService(queries)
	uploadSvc := services.NewUploadService(t.TempDir())

	// Public routes
	homeHandler := publicHandlers.NewHomeHandler(queries, testLogger)
	e.GET("/", homeHandler.ShowHomePage)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	cache := services.NewCache()
	productsHandler := publicHandlers.NewProductsHandler(queries, testLogger, productSvc, cache)
	e.GET("/products", productsHandler.ProductsList)
	e.GET("/products/:category", productsHandler.ProductsByCategory)
	e.GET("/products/:category/:slug", productsHandler.ProductDetail)

	// Admin auth routes (no auth middleware)
	authHandler := adminHandlers.NewAuthHandler(queries, testLogger)
	e.GET("/admin/login", authHandler.ShowLoginPage)
	e.POST("/admin/login", authHandler.LoginSubmit)
	e.POST("/admin/logout", authHandler.Logout)

	// Admin routes (with auth)
	adminGroup := e.Group("/admin", customMiddleware.RequireAuth())
	dashHandler := adminHandlers.NewDashboardHandler()
	adminGroup.GET("/dashboard", dashHandler.ShowDashboard)

	pcHandler := adminHandlers.NewProductCategoriesHandler(queries, testLogger)
	adminGroup.GET("/product-categories", pcHandler.List)
	adminGroup.POST("/product-categories", pcHandler.Create)
	adminGroup.GET("/product-categories/:id/edit", pcHandler.Edit)
	adminGroup.POST("/product-categories/:id", pcHandler.Update)
	adminGroup.DELETE("/product-categories/:id", pcHandler.Delete)

	adminProductsHandler := adminHandlers.NewProductsHandler(queries, testLogger, uploadSvc, cache)
	adminGroup.GET("/products", adminProductsHandler.List)
	adminGroup.POST("/products", adminProductsHandler.Create)
	adminGroup.DELETE("/products/:id", adminProductsHandler.Delete)

	_ = db
	return e, queries, cleanup
}

func createTestAdmin(t *testing.T, queries *sqlc.Queries) {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	_, err := queries.CreateAdminUser(context.Background(), sqlc.CreateAdminUserParams{
		Email: "admin@test.com", PasswordHash: string(hash), DisplayName: "Test Admin", Role: "admin",
	})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}
}

func loginAndGetCookie(t *testing.T, e *echo.Echo) *http.Cookie {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{
		"email": {"admin@test.com"}, "password": {"testpassword"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	for _, c := range rec.Result().Cookies() {
		if c.Name == "bluejay_session" {
			return c
		}
	}
	t.Fatal("no session cookie after login")
	return nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestHealthEndpoint(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", resp["status"])
	}
}

func TestLoginFlow(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)

	// Test missing fields
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Errorf("missing fields: expected 303, got %d", rec.Code)
	}

	// Test wrong password
	req = httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{
		"email": {"admin@test.com"}, "password": {"wrong"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if !strings.Contains(rec.Header().Get("Location"), "error=invalid_credentials") {
		t.Errorf("wrong password: expected invalid_credentials redirect, got Location=%q", rec.Header().Get("Location"))
	}

	// Test successful login
	cookie := loginAndGetCookie(t, e)

	// Test authenticated access to dashboard
	req = httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	// Should not redirect (200 if template renders, or 500 if no renderer, but not 303)
	if rec.Code == http.StatusSeeOther {
		t.Errorf("expected authenticated access, got redirect to %s", rec.Header().Get("Location"))
	}
}

func TestAdminProtectedRoutes_RequireAuth(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	protectedPaths := []string{
		"/admin/dashboard",
		"/admin/product-categories",
		"/admin/products",
	}

	for _, path := range protectedPaths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != http.StatusSeeOther {
				t.Errorf("expected 303 redirect for unauthenticated %s, got %d", path, rec.Code)
			}
			loc := rec.Header().Get("Location")
			if loc != "/admin/login" {
				t.Errorf("expected redirect to /admin/login, got %q", loc)
			}
		})
	}
}

func TestAdminProductCategoryCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// Create category
	req := httptest.NewRequest(http.MethodPost, "/admin/product-categories", strings.NewReader(url.Values{
		"name": {"E2E Category"}, "description": {"Test desc"}, "icon": {"test_icon"}, "sort_order": {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d; body: %s", rec.Code, rec.Body.String())
	}

	// Verify in DB
	cats, _ := queries.ListProductCategories(context.Background())
	if len(cats) != 1 {
		t.Fatalf("expected 1 category, got %d", len(cats))
	}
	if cats[0].Name != "E2E Category" {
		t.Errorf("expected 'E2E Category', got %q", cats[0].Name)
	}
	if cats[0].Slug != "e2e-category" {
		t.Errorf("expected slug 'e2e-category', got %q", cats[0].Slug)
	}

	// Delete category
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/product-categories/%d", cats[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	cats, _ = queries.ListProductCategories(context.Background())
	if len(cats) != 0 {
		t.Errorf("expected 0 categories after delete, got %d", len(cats))
	}
}

func TestPublicProducts_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	// Seed data
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Detectors", Slug: "detectors", Description: "Detection equipment", Icon: "radar", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DET-001", Slug: "alpha-detector", Name: "Alpha Detector",
		Description: "Detects alpha particles", CategoryID: cat.ID, Status: "published",
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	// Products list (no renderer so expect 500, but not 404)
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Error("products route not found")
	}

	// Products by category
	req = httptest.NewRequest(http.MethodGet, "/products/detectors", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Error("products category route not found")
	}

	// Product detail
	req = httptest.NewRequest(http.MethodGet, "/products/detectors/alpha-detector", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Error("product detail route not found")
	}

	// Nonexistent category should 404
	req = httptest.NewRequest(http.MethodGet, "/products/nonexistent", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent category, got %d", rec.Code)
	}

	// Wrong category for product should 404
	req = httptest.NewRequest(http.MethodGet, "/products/wrong-cat/alpha-detector", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for wrong category, got %d", rec.Code)
	}
}

func TestLogout_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// Logout
	req := httptest.NewRequest(http.MethodPost, "/admin/logout", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	// Try accessing protected route with cookies from logout response (should be logged out)
	req = httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	for _, c := range rec.Result().Cookies() {
		req.AddCookie(c)
	}
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect after logout, got %d", rec.Code)
	}
}

func TestAdminProductDelete_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DEL-E2E", Slug: "del-e2e", Name: "Delete Me",
		Description: "d", CategoryID: cat.ID, Status: "draft",
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/products/%d", prod.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	_, err := queries.GetProduct(ctx, prod.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}
