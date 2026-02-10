// Package e2e_test contains end-to-end integration tests for the Bluejay CMS application.
//
// These tests verify the complete request-response cycle including routing, middleware,
// handlers, services, and database operations. Unlike unit tests that isolate individual
// components, e2e tests validate the entire application stack working together.
//
// Test Coverage:
//   - Public routes (homepage, products, health checks)
//   - Admin authentication flow (login, logout, session management)
//   - Protected admin routes (requiring authentication)
//   - CRUD operations for product categories and products
//   - Authorization and access control
//
// Each test uses a fresh SQLite database with migrations applied, ensuring test isolation
// and repeatability. The tests use httptest to simulate HTTP requests without starting
// a real server.
package e2e_test

import (
	// Standard library imports
	"context"           // Context for database operations and request cancellation
	"database/sql"      // SQL error types like ErrNoRows for assertion checks
	"encoding/json"     // JSON decoding for validating API response payloads
	"fmt"               // String formatting for constructing dynamic URLs
	"io"                // I/O interfaces for the stub template renderer
	"log/slog"          // Structured logging for handlers (error-level in tests)
	"net/http"          // HTTP constants and types (StatusOK, StatusSeeOther, etc.)
	"net/http/httptest" // HTTP testing utilities for creating requests and recording responses
	"net/url"           // URL encoding for form data
	"os"                // File system access for logger output (stderr)
	"strings"           // String manipulation for request body creation
	"testing"           // Go testing framework

	// Third-party imports
	"github.com/labstack/echo/v4"    // Echo web framework - the HTTP router and context
	"golang.org/x/crypto/bcrypt"     // Password hashing for creating test admin users

	// Project imports
	"github.com/narendhupati/bluejay-cms/db/sqlc"                      // sqlc generated database queries
	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"   // Admin panel HTTP handlers
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public" // Public-facing HTTP handlers
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"    // Session and auth middleware
	"github.com/narendhupati/bluejay-cms/internal/services"            // Business logic layer (products, uploads, cache)
	"github.com/narendhupati/bluejay-cms/internal/testutil"            // Test database setup utilities
)

// Package-level variables

var (
	// testLogger is a structured logger configured for e2e tests.
	// It outputs to stderr and is set to ERROR level to reduce noise during test runs.
	// Only errors and critical issues will be logged, making test output cleaner
	// while still capturing important diagnostic information if tests fail.
	testLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
)

// stubRenderer is a minimal template renderer implementation for e2e tests.
//
// Most e2e tests focus on routing, middleware, and data layer behavior rather than
// HTML rendering. This stub allows handlers that call c.Render() to succeed without
// requiring real template files. It outputs a simple HTML stub instead.
//
// This approach:
//   - Speeds up tests by avoiding template parsing
//   - Eliminates template file dependencies in e2e tests
//   - Allows tests to focus on HTTP status codes, redirects, and data persistence
//
// Tests that need to verify actual template output should use a real renderer
// and check the response body content.
type stubRenderer struct{}

// Render implements the echo.Renderer interface by writing a simple HTML stub.
//
// Parameters:
//   - w: The io.Writer to write the rendered output to (typically the HTTP response)
//   - name: The template name (ignored by this stub implementation)
//   - data: The template data (ignored by this stub implementation)
//   - c: The Echo context (ignored by this stub implementation)
//
// Returns:
//   - error: Any error from writing to the writer, or nil on success
func (r *stubRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Write a minimal HTML response that satisfies handlers expecting template rendering
	_, err := w.Write([]byte("<html>stub</html>"))
	return err
}

// setupApp creates and configures a complete Echo application instance for e2e testing.
//
// This function mirrors the production application setup from cmd/server/main.go but
// with test-specific configuration:
//   - Uses a fresh temporary SQLite database for each test
//   - Initializes session middleware with a test-only secret key
//   - Configures all routes (public and admin)
//   - Uses a stub renderer instead of real templates
//   - Uses a temporary directory for file uploads
//
// The returned Echo instance is fully functional and can handle HTTP requests via
// httptest without starting a real server. This allows tests to verify the complete
// request-response cycle including routing, middleware, handlers, and database operations.
//
// Parameters:
//   - t: The testing.T instance for marking this as a helper and automatic cleanup
//
// Returns:
//   - *echo.Echo: A fully configured Echo application instance ready for testing
//   - *sqlc.Queries: A database query interface for test assertions and data seeding
//   - func(): A cleanup function that closes the database and removes temporary files.
//     Must be deferred immediately after calling setupApp.
//
// Example usage:
//
//	func TestSomeFeature(t *testing.T) {
//	    app, queries, cleanup := setupApp(t)
//	    defer cleanup()
//	    // Create test data with queries, make requests to app
//	}
func setupApp(t *testing.T) (*echo.Echo, *sqlc.Queries, func()) {
	// Mark this as a test helper for better error reporting
	t.Helper()

	// Create a fresh test database with all migrations applied
	db, queries, cleanup := testutil.SetupTestDB(t)

	// Initialize the session store with a test-specific secret key.
	// This key must be at least 32 characters for secure cookie encryption.
	// Using a fixed key in tests ensures consistent session behavior across test runs.
	customMiddleware.InitSessionStore("e2e-test-secret-at-least-32-characters-long")

	// Create a new Echo instance
	e := echo.New()
	e.HideBanner = true        // Suppress Echo startup banner in test output
	e.Renderer = &stubRenderer{} // Use stub renderer to avoid template file dependencies

	// Apply session middleware globally so all routes can access session data
	e.Use(customMiddleware.SessionMiddleware())

	// Initialize services that handlers depend on
	productSvc := services.NewProductService(queries)
	uploadSvc := services.NewUploadService(t.TempDir()) // Use temp directory for file uploads

	// Register public routes (no authentication required)
	homeHandler := publicHandlers.NewHomeHandler(queries, testLogger)
	e.GET("/", homeHandler.ShowHomePage)

	// Health check endpoint for monitoring
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Product catalog public routes
	cache := services.NewCache()
	productsHandler := publicHandlers.NewProductsHandler(queries, testLogger, productSvc, cache)
	e.GET("/products", productsHandler.ProductsList)
	e.GET("/products/:category", productsHandler.ProductsByCategory)
	e.GET("/products/:category/:slug", productsHandler.ProductDetail)

	// Admin authentication routes (no auth middleware - these handle login/logout)
	authHandler := adminHandlers.NewAuthHandler(queries, testLogger)
	e.GET("/admin/login", authHandler.ShowLoginPage)
	e.POST("/admin/login", authHandler.LoginSubmit)
	e.POST("/admin/logout", authHandler.Logout)

	// Admin routes group with authentication middleware
	// All routes in this group require a valid session with admin privileges
	adminGroup := e.Group("/admin", customMiddleware.RequireAuth())

	// Dashboard route
	dashHandler := adminHandlers.NewDashboardHandler(queries, testLogger)
	adminGroup.GET("/dashboard", dashHandler.ShowDashboard)

	// Product categories CRUD routes
	pcHandler := adminHandlers.NewProductCategoriesHandler(queries, testLogger)
	adminGroup.GET("/product-categories", pcHandler.List)
	adminGroup.POST("/product-categories", pcHandler.Create)
	adminGroup.GET("/product-categories/:id/edit", pcHandler.Edit)
	adminGroup.POST("/product-categories/:id", pcHandler.Update)
	adminGroup.DELETE("/product-categories/:id", pcHandler.Delete)

	// Products CRUD routes
	adminProductsHandler := adminHandlers.NewProductsHandler(queries, testLogger, uploadSvc, cache)
	adminGroup.GET("/products", adminProductsHandler.List)
	adminGroup.POST("/products", adminProductsHandler.Create)
	adminGroup.DELETE("/products/:id", adminProductsHandler.Delete)

	// Keep db reference to avoid unused variable error, though cleanup handles closing
	_ = db

	return e, queries, cleanup
}

// createTestAdmin creates a test admin user in the database for authentication tests.
//
// This helper function creates a user with known credentials that can be used to test
// the login flow and access protected admin routes. The password is hashed using bcrypt
// to match the production authentication behavior.
//
// Test credentials:
//   - Email: admin@test.com
//   - Password: testpassword (plaintext, hashed before storage)
//   - Role: admin
//
// Parameters:
//   - t: The testing.T instance for error reporting
//   - queries: The database query interface for creating the user
//
// Note: This function ignores the bcrypt hashing error because bcrypt.GenerateFromPassword
// only returns an error if the cost parameter is invalid, which won't happen with
// bcrypt.DefaultCost. In production code, this error should be checked.
func createTestAdmin(t *testing.T, queries *sqlc.Queries) {
	// Mark as helper for better error stack traces
	t.Helper()

	// Hash the test password using bcrypt with default cost (currently 10).
	// We ignore the error here because DefaultCost is guaranteed to be valid.
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)

	// Insert the admin user into the database
	_, err := queries.CreateAdminUser(context.Background(), sqlc.CreateAdminUserParams{
		Email:        "admin@test.com",
		PasswordHash: string(hash),
		DisplayName:  "Test Admin",
		Role:         "admin",
	})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}
}

// loginAndGetCookie performs a login request and extracts the session cookie.
//
// This helper automates the common test pattern of logging in as the test admin user
// and retrieving the session cookie for subsequent authenticated requests. The function
// submits a POST request to /admin/login with the test admin credentials, then searches
// the response cookies for the session cookie.
//
// This is used in tests that need to access protected admin routes. The returned cookie
// should be added to subsequent requests using req.AddCookie(cookie).
//
// Parameters:
//   - t: The testing.T instance for error reporting
//   - e: The Echo application instance to send the login request to
//
// Returns:
//   - *http.Cookie: The session cookie returned after successful login
//
// Panics:
//   - Calls t.Fatal if no session cookie is found in the response, which indicates
//     the login failed or the session middleware isn't working correctly.
//
// Example usage:
//
//	cookie := loginAndGetCookie(t, app)
//	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
//	req.AddCookie(cookie)
func loginAndGetCookie(t *testing.T, e *echo.Echo) *http.Cookie {
	// Mark as helper for better error stack traces
	t.Helper()

	// Create a POST request to the login endpoint with form-encoded credentials
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{
		"email":    {"admin@test.com"},
		"password": {"testpassword"},
	}.Encode()))

	// Set the content type to application/x-www-form-urlencoded so Echo's form parsing works
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	// Record the response
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Search for the session cookie in the response
	for _, c := range rec.Result().Cookies() {
		if c.Name == "bluejay_session" {
			return c
		}
	}

	// If we reach here, login failed or session middleware didn't set a cookie
	t.Fatal("no session cookie after login")
	return nil
}

// ---------------------------------------------------------------------------
// Test Cases
// ---------------------------------------------------------------------------

// TestHealthEndpoint verifies the health check endpoint returns the expected response.
//
// This test ensures:
//   - The /health endpoint is accessible without authentication
//   - The response status code is 200 OK
//   - The response body is valid JSON
//   - The JSON contains {"status": "ok"}
//
// The health endpoint is typically used by load balancers, monitoring systems, and
// orchestration platforms (like Kubernetes) to verify the application is running.
func TestHealthEndpoint(t *testing.T) {
	// Set up the application with a fresh database
	e, _, cleanup := setupApp(t)
	defer cleanup()

	// Create a GET request to the health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify the response status code is 200 OK
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// Parse the JSON response body
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify the response contains the expected status field
	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", resp["status"])
	}
}

// TestLoginFlow verifies the complete admin authentication flow including error cases.
//
// This test validates:
//   - Login with missing credentials fails with redirect
//   - Login with incorrect password fails with error parameter
//   - Login with correct credentials succeeds and returns a session cookie
//   - The session cookie grants access to protected admin routes
//
// The test covers both the authentication logic and session management, ensuring
// that the auth middleware correctly validates sessions created by the login handler.
func TestLoginFlow(t *testing.T) {
	// Set up the application and create a test admin user
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)

	// Test Case 1: Login with missing fields should fail
	// Submit an empty form to trigger validation errors
	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Expect a 303 redirect back to the login page with error message
	if rec.Code != http.StatusSeeOther {
		t.Errorf("missing fields: expected 303, got %d", rec.Code)
	}

	// Test Case 2: Login with wrong password should fail
	req = httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{
		"email":    {"admin@test.com"},
		"password": {"wrong"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Expect a redirect with invalid_credentials error parameter in the URL
	if !strings.Contains(rec.Header().Get("Location"), "error=invalid_credentials") {
		t.Errorf("wrong password: expected invalid_credentials redirect, got Location=%q", rec.Header().Get("Location"))
	}

	// Test Case 3: Login with correct credentials should succeed
	// This helper submits valid credentials and extracts the session cookie
	cookie := loginAndGetCookie(t, e)

	// Test Case 4: Verify the session cookie grants access to protected routes
	req = httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Should NOT redirect to login page (which would be a 303).
	// We expect either 200 (if template renders) or 500 (if renderer fails),
	// but definitely not a redirect which would indicate auth failure.
	if rec.Code == http.StatusSeeOther {
		t.Errorf("expected authenticated access, got redirect to %s", rec.Header().Get("Location"))
	}
}

// TestAdminProtectedRoutes_RequireAuth verifies that admin routes require authentication.
//
// This test ensures that the RequireAuth() middleware is properly applied to all admin
// routes and correctly redirects unauthenticated requests to the login page.
//
// The test verifies:
//   - Unauthenticated requests to admin routes return 303 See Other (redirect)
//   - The redirect location is /admin/login
//   - This behavior is consistent across all protected admin routes
//
// This is a critical security test - if it fails, unauthorized users could access
// admin functionality.
func TestAdminProtectedRoutes_RequireAuth(t *testing.T) {
	// Set up the application (no admin user needed - we're testing unauthorized access)
	e, _, cleanup := setupApp(t)
	defer cleanup()

	// List of admin routes that should require authentication
	protectedPaths := []string{
		"/admin/dashboard",
		"/admin/product-categories",
		"/admin/products",
	}

	// Test each protected route in a subtest for clear failure reporting
	for _, path := range protectedPaths {
		t.Run(path, func(t *testing.T) {
			// Make a request WITHOUT a session cookie (unauthenticated)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			// Verify the response is a redirect (303 See Other)
			if rec.Code != http.StatusSeeOther {
				t.Errorf("expected 303 redirect for unauthenticated %s, got %d", path, rec.Code)
			}

			// Verify the redirect goes to the login page
			loc := rec.Header().Get("Location")
			if loc != "/admin/login" {
				t.Errorf("expected redirect to /admin/login, got %q", loc)
			}
		})
	}
}

// TestAdminProductCategoryCRUD_E2E tests the complete product category lifecycle.
//
// This end-to-end test validates:
//   - Creating a product category via POST request
//   - Automatic slug generation from the category name
//   - Database persistence of the created category
//   - Deleting a category via DELETE request
//   - Database cleanup after deletion
//
// The test simulates the full user workflow: authenticate, create a category through
// the admin panel, verify it exists in the database, delete it, and confirm removal.
func TestAdminProductCategoryCRUD_E2E(t *testing.T) {
	// Set up application, test admin, and get session cookie
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// Step 1: Create a product category via POST request
	req := httptest.NewRequest(http.MethodPost, "/admin/product-categories", strings.NewReader(url.Values{
		"name":        {"E2E Category"},
		"description": {"Test desc"},
		"icon":        {"test_icon"},
		"sort_order":  {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie) // Include auth cookie
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify the handler redirects after successful creation
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d; body: %s", rec.Code, rec.Body.String())
	}

	// Step 2: Verify the category was persisted to the database
	cats, _ := queries.ListProductCategories(context.Background())
	if len(cats) != 1 {
		t.Fatalf("expected 1 category, got %d", len(cats))
	}

	// Verify the category data matches what we submitted
	if cats[0].Name != "E2E Category" {
		t.Errorf("expected 'E2E Category', got %q", cats[0].Name)
	}

	// Verify the slug was auto-generated correctly (spaces -> hyphens, lowercase)
	if cats[0].Slug != "e2e-category" {
		t.Errorf("expected slug 'e2e-category', got %q", cats[0].Slug)
	}

	// Step 3: Delete the category via DELETE request
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/product-categories/%d", cats[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// DELETE endpoints typically return 200 OK for HTMX responses
	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	// Step 4: Verify the category was removed from the database
	cats, _ = queries.ListProductCategories(context.Background())
	if len(cats) != 0 {
		t.Errorf("expected 0 categories after delete, got %d", len(cats))
	}
}

// TestPublicProducts_E2E verifies the public product catalog routes and error handling.
//
// This test validates:
//   - Public access to product listing pages (no authentication required)
//   - Correct routing for category-based product views
//   - Product detail page routing with category and slug
//   - Proper 404 responses for nonexistent categories and products
//   - Category validation in product detail URLs
//
// The test seeds test data (a category and product) then verifies all public routes
// return appropriate status codes. Since we use a stub renderer, we verify route
// existence rather than template output.
func TestPublicProducts_E2E(t *testing.T) {
	// Set up the application with a fresh database
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	// Seed test data: create a product category
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Detectors",
		Slug:        "detectors",
		Description: "Detection equipment",
		Icon:        "radar",
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	// Seed test data: create a product in that category
	_, err = queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "DET-001",
		Slug:        "alpha-detector",
		Name:        "Alpha Detector",
		Description: "Detects alpha particles",
		CategoryID:  cat.ID,
		Status:      "published",
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	// Test Case 1: Products list route should exist
	// Note: With stub renderer, we may get 500 but NOT 404
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Error("products route not found")
	}

	// Test Case 2: Products by category route should exist
	req = httptest.NewRequest(http.MethodGet, "/products/detectors", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Error("products category route not found")
	}

	// Test Case 3: Product detail route should exist
	// URL pattern: /products/:category/:slug
	req = httptest.NewRequest(http.MethodGet, "/products/detectors/alpha-detector", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusNotFound {
		t.Error("product detail route not found")
	}

	// Test Case 4: Nonexistent category should return 404
	req = httptest.NewRequest(http.MethodGet, "/products/nonexistent", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent category, got %d", rec.Code)
	}

	// Test Case 5: Product with wrong category should return 404
	// This verifies the handler validates the category matches the product's actual category
	req = httptest.NewRequest(http.MethodGet, "/products/wrong-cat/alpha-detector", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for wrong category, got %d", rec.Code)
	}
}

// TestLogout_E2E verifies the logout functionality and session invalidation.
//
// This test ensures:
//   - The logout endpoint processes POST requests correctly
//   - Logout invalidates the session
//   - After logout, protected routes redirect to login
//   - Session cookies are properly cleared or invalidated
//
// The test simulates logging in, logging out, then attempting to access a protected
// route with the cookies from the logout response to verify the session is truly
// invalidated.
func TestLogout_E2E(t *testing.T) {
	// Set up application and authenticate
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// Step 1: Perform logout with valid session
	req := httptest.NewRequest(http.MethodPost, "/admin/logout", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify logout redirects (typically to login page or home)
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}

	// Step 2: Try accessing a protected route with cookies from logout response
	// The logout handler should have invalidated the session, so even if we send
	// the cookies back, we should be redirected to login.
	req = httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	for _, c := range rec.Result().Cookies() {
		req.AddCookie(c)
	}
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Should redirect to login because session is invalidated
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect after logout, got %d", rec.Code)
	}
}

// TestAdminProductDelete_E2E verifies product deletion functionality.
//
// This test ensures:
//   - Authenticated admin users can delete products via DELETE request
//   - Successful deletion returns HTTP 200 OK
//   - The product is removed from the database
//   - Attempting to fetch the deleted product returns sql.ErrNoRows
//
// The test creates a test product, deletes it via the admin endpoint, then verifies
// the product no longer exists in the database.
func TestAdminProductDelete_E2E(t *testing.T) {
	// Set up application and authenticate
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// Seed test data: create a product category
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Cat",
		Slug:        "cat",
		Description: "d",
		Icon:        "i",
		SortOrder:   1,
	})

	// Seed test data: create a product to delete
	prod, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku:         "DEL-E2E",
		Slug:        "del-e2e",
		Name:        "Delete Me",
		Description: "d",
		CategoryID:  cat.ID,
		Status:      "draft",
	})

	// Step 1: Send DELETE request to remove the product
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/products/%d", prod.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify the deletion succeeded (HTMX DELETE endpoints return 200)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	// Step 2: Verify the product was removed from the database
	// Attempting to fetch it should return sql.ErrNoRows
	_, err := queries.GetProduct(ctx, prod.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}
