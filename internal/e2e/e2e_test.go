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
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// Package-level variables

var (
	testLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
)

func TestMain(m *testing.M) {
	// Change to project root so template files can be found
	if err := os.Chdir("../../"); err != nil {
		fmt.Fprintf(os.Stderr, "chdir to project root: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

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
	t.Helper()

	db, queries, cleanup := testutil.SetupTestDB(t)
	customMiddleware.InitSessionStore("e2e-test-secret-at-least-32-characters-long")

	e := echo.New()
	e.HideBanner = true
	e.Renderer = &stubRenderer{}
	e.Use(customMiddleware.SecurityHeaders())
	e.Use(customMiddleware.SessionMiddleware())

	// Services
	productSvc := services.NewProductService(queries)
	uploadSvc := services.NewUploadService(t.TempDir())
	appCache := services.NewCache()
	activitySvc := services.NewActivityLogService(queries, testLogger)
	adminHandlers.SetActivityLogService(activitySvc)

	// Public routes
	homeHandler := publicHandlers.NewHomeHandler(queries, testLogger)
	e.GET("/", homeHandler.ShowHomePage)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	productsHandler := publicHandlers.NewProductsHandler(queries, testLogger, productSvc, appCache)
	e.GET("/products", productsHandler.ProductsList)
	e.GET("/products/search", productsHandler.ProductSearch)
	e.GET("/products/:category", productsHandler.ProductsByCategory)
	e.GET("/products/:category/:slug", productsHandler.ProductDetail)

	solutionsHandler := publicHandlers.NewSolutionsHandler(queries, testLogger, appCache)
	e.GET("/solutions", solutionsHandler.SolutionsList)
	e.GET("/solutions/:slug", solutionsHandler.SolutionDetail)

	blogHandler := publicHandlers.NewBlogHandler(queries, testLogger, appCache)
	e.GET("/blog", blogHandler.BlogListing)
	e.GET("/blog/:slug", blogHandler.BlogPost)

	whitepapersHandler := publicHandlers.NewWhitepapersHandler(queries, testLogger, appCache)
	e.GET("/whitepapers", whitepapersHandler.WhitepapersList)
	e.GET("/whitepapers/:slug", whitepapersHandler.WhitepaperDetail)
	e.POST("/whitepapers/:slug/download", whitepapersHandler.WhitepaperDownload)

	contactHandler := publicHandlers.NewContactHandler(queries, testLogger, appCache)
	contactLimiter := customMiddleware.NewRateLimiter(5, time.Hour)
	e.GET("/contact", contactHandler.ShowContactPage)
	e.POST("/contact/submit", contactHandler.SubmitContactForm, contactLimiter.Middleware())

	aboutHandler := publicHandlers.NewAboutHandler(queries, testLogger, appCache)
	e.GET("/about", aboutHandler.AboutPage)

	partnersPageHandler := publicHandlers.NewPartnersHandler(queries, testLogger, appCache)
	e.GET("/partners", partnersPageHandler.PartnersPage)

	searchHandler := publicHandlers.NewSearchHandler(db, testLogger)
	e.GET("/search", searchHandler.SearchPage)
	e.GET("/search/suggest", searchHandler.SearchSuggest)

	sitemapHandler := publicHandlers.NewSitemapHandler(queries, testLogger, "https://bluejaylabs.com")
	e.GET("/sitemap.xml", sitemapHandler.Sitemap)
	e.GET("/robots.txt", sitemapHandler.RobotsTxt)

	caseStudiesHandler := publicHandlers.NewCaseStudiesHandler(queries, testLogger, appCache)
	e.GET("/case-studies", caseStudiesHandler.CaseStudiesList)
	e.GET("/case-studies/:slug", caseStudiesHandler.CaseStudyDetail)

	// Admin auth routes
	authHandler := adminHandlers.NewAuthHandler(queries, testLogger)
	e.GET("/admin/login", authHandler.ShowLoginPage)
	e.POST("/admin/login", authHandler.LoginSubmit)
	e.POST("/admin/logout", authHandler.Logout)

	// Admin protected routes
	adminGroup := e.Group("/admin", customMiddleware.RequireAuth())

	dashHandler := adminHandlers.NewDashboardHandler(queries, testLogger)
	adminGroup.GET("/dashboard", dashHandler.ShowDashboard)

	// Product categories
	pcHandler := adminHandlers.NewProductCategoriesHandler(queries, testLogger)
	adminGroup.GET("/product-categories", pcHandler.List)
	adminGroup.GET("/product-categories/new", pcHandler.New)
	adminGroup.POST("/product-categories", pcHandler.Create)
	adminGroup.GET("/product-categories/:id/edit", pcHandler.Edit)
	adminGroup.POST("/product-categories/:id", pcHandler.Update)
	adminGroup.DELETE("/product-categories/:id", pcHandler.Delete)

	// Blog categories
	bcHandler := adminHandlers.NewBlogCategoriesHandler(queries, testLogger)
	adminGroup.GET("/blog-categories", bcHandler.List)
	adminGroup.GET("/blog-categories/new", bcHandler.New)
	adminGroup.POST("/blog-categories", bcHandler.Create)
	adminGroup.GET("/blog-categories/:id/edit", bcHandler.Edit)
	adminGroup.POST("/blog-categories/:id", bcHandler.Update)
	adminGroup.DELETE("/blog-categories/:id", bcHandler.Delete)

	// Blog authors
	baHandler := adminHandlers.NewBlogAuthorsHandler(queries, testLogger)
	adminGroup.GET("/blog-authors", baHandler.List)
	adminGroup.GET("/blog-authors/new", baHandler.New)
	adminGroup.POST("/blog-authors", baHandler.Create)
	adminGroup.GET("/blog-authors/:id/edit", baHandler.Edit)
	adminGroup.POST("/blog-authors/:id", baHandler.Update)
	adminGroup.DELETE("/blog-authors/:id", baHandler.Delete)

	// Industries
	indHandler := adminHandlers.NewIndustriesHandler(queries, testLogger)
	adminGroup.GET("/industries", indHandler.List)
	adminGroup.GET("/industries/new", indHandler.New)
	adminGroup.POST("/industries", indHandler.Create)
	adminGroup.GET("/industries/:id/edit", indHandler.Edit)
	adminGroup.POST("/industries/:id", indHandler.Update)
	adminGroup.DELETE("/industries/:id", indHandler.Delete)

	// Partner tiers
	ptHandler := adminHandlers.NewPartnerTiersHandler(queries, testLogger)
	adminGroup.GET("/partner-tiers", ptHandler.List)
	adminGroup.GET("/partner-tiers/new", ptHandler.New)
	adminGroup.POST("/partner-tiers", ptHandler.Create)
	adminGroup.GET("/partner-tiers/:id/edit", ptHandler.Edit)
	adminGroup.POST("/partner-tiers/:id", ptHandler.Update)
	adminGroup.DELETE("/partner-tiers/:id", ptHandler.Delete)

	// Whitepaper topics
	wtHandler := adminHandlers.NewWhitepaperTopicsHandler(queries, testLogger)
	adminGroup.GET("/whitepaper-topics", wtHandler.List)
	adminGroup.GET("/whitepaper-topics/new", wtHandler.New)
	adminGroup.POST("/whitepaper-topics", wtHandler.Create)
	adminGroup.GET("/whitepaper-topics/:id/edit", wtHandler.Edit)
	adminGroup.POST("/whitepaper-topics/:id", wtHandler.Update)
	adminGroup.DELETE("/whitepaper-topics/:id", wtHandler.Delete)

	// Header, footer, settings
	headerHandler := adminHandlers.NewHeaderHandler(queries, testLogger)
	adminGroup.GET("/header", headerHandler.Edit)
	adminGroup.POST("/header", headerHandler.Update)

	footerHandler := adminHandlers.NewFooterHandler(queries, testLogger)
	adminGroup.GET("/footer", footerHandler.Edit)
	adminGroup.POST("/footer", footerHandler.Update)

	settingsHandler := adminHandlers.NewSettingsHandler(queries, testLogger)
	adminGroup.GET("/settings", settingsHandler.Edit)
	adminGroup.POST("/settings", settingsHandler.Update)

	// Page sections
	psHandler := adminHandlers.NewPageSectionsHandler(queries, testLogger)
	adminGroup.GET("/page-sections", psHandler.List)
	adminGroup.GET("/page-sections/:id/edit", psHandler.Edit)
	adminGroup.POST("/page-sections/:id", psHandler.Update)

	// Products
	adminProductsHandler := adminHandlers.NewProductsHandler(queries, testLogger, uploadSvc, appCache)
	adminGroup.GET("/products", adminProductsHandler.List)
	adminGroup.GET("/products/new", adminProductsHandler.New)
	adminGroup.POST("/products", adminProductsHandler.Create)
	adminGroup.GET("/products/:id/edit", adminProductsHandler.Edit)
	adminGroup.POST("/products/:id", adminProductsHandler.Update)
	adminGroup.DELETE("/products/:id", adminProductsHandler.Delete)

	// Product details (specs, features, certs, downloads, images)
	pdHandler := adminHandlers.NewProductDetailsHandler(queries, testLogger, uploadSvc)
	adminGroup.GET("/products/:id/specs", pdHandler.ListSpecs)
	adminGroup.POST("/products/:id/specs", pdHandler.AddSpec)
	adminGroup.DELETE("/products/:id/specs", pdHandler.DeleteSpecs)
	adminGroup.DELETE("/products/:id/specs/:spec_id", pdHandler.DeleteSpec)
	adminGroup.GET("/products/:id/features", pdHandler.ListFeatures)
	adminGroup.POST("/products/:id/features", pdHandler.AddFeature)
	adminGroup.DELETE("/products/:id/features", pdHandler.DeleteFeatures)
	adminGroup.DELETE("/products/:id/features/:feature_id", pdHandler.DeleteFeature)
	adminGroup.GET("/products/:id/certifications", pdHandler.ListCertifications)
	adminGroup.POST("/products/:id/certifications", pdHandler.AddCertification)
	adminGroup.DELETE("/products/:id/certifications", pdHandler.DeleteCertifications)
	adminGroup.DELETE("/products/:id/certifications/:cert_id", pdHandler.DeleteCertification)
	adminGroup.GET("/products/:id/downloads", pdHandler.ListDownloads)
	adminGroup.POST("/products/:id/downloads", pdHandler.AddDownload)
	adminGroup.DELETE("/products/:id/downloads/:download_id", pdHandler.DeleteDownload)
	adminGroup.GET("/products/:id/images", pdHandler.ListImages)
	adminGroup.POST("/products/:id/images", pdHandler.AddImage)
	adminGroup.DELETE("/products/:id/images/:image_id", pdHandler.DeleteImage)

	// Blog posts
	adminBlogPostsHandler := adminHandlers.NewBlogPostsHandler(queries, testLogger, appCache)
	adminGroup.GET("/blog/posts", adminBlogPostsHandler.List)
	adminGroup.GET("/blog/posts/new", adminBlogPostsHandler.New)
	adminGroup.POST("/blog/posts", adminBlogPostsHandler.Create)
	adminGroup.GET("/blog/posts/:id/edit", adminBlogPostsHandler.Edit)
	adminGroup.POST("/blog/posts/:id", adminBlogPostsHandler.Update)
	adminGroup.DELETE("/blog/posts/:id", adminBlogPostsHandler.Delete)
	adminGroup.GET("/blog/products/search", adminBlogPostsHandler.SearchProducts)

	// Blog tags
	adminBlogTagsHandler := adminHandlers.NewBlogTagsHandler(queries, testLogger)
	adminGroup.GET("/blog/tags", adminBlogTagsHandler.List)
	adminGroup.POST("/blog/tags", adminBlogTagsHandler.Create)
	adminGroup.GET("/blog/tags/search", adminBlogTagsHandler.Search)
	adminGroup.POST("/blog/tags/quick-create", adminBlogTagsHandler.QuickCreate)
	adminGroup.DELETE("/blog/tags/:id", adminBlogTagsHandler.Delete)

	// Solutions
	adminSolutionsHandler := adminHandlers.NewSolutionsHandler(queries, testLogger, appCache)
	adminGroup.GET("/solutions", adminSolutionsHandler.List)
	adminGroup.GET("/solutions/new", adminSolutionsHandler.New)
	adminGroup.POST("/solutions", adminSolutionsHandler.Create)
	adminGroup.GET("/solutions/:id/edit", adminSolutionsHandler.Edit)
	adminGroup.POST("/solutions/:id", adminSolutionsHandler.Update)
	adminGroup.DELETE("/solutions/:id", adminSolutionsHandler.Delete)
	adminGroup.POST("/solutions/:id/stats", adminSolutionsHandler.AddStat)
	adminGroup.DELETE("/solutions/:id/stats/:statId", adminSolutionsHandler.DeleteStat)
	adminGroup.POST("/solutions/:id/challenges", adminSolutionsHandler.AddChallenge)
	adminGroup.DELETE("/solutions/:id/challenges/:challengeId", adminSolutionsHandler.DeleteChallenge)
	adminGroup.POST("/solutions/:id/products", adminSolutionsHandler.AddProduct)
	adminGroup.DELETE("/solutions/:id/products/:productId", adminSolutionsHandler.RemoveProduct)
	adminGroup.POST("/solutions/:id/ctas", adminSolutionsHandler.AddCTA)
	adminGroup.DELETE("/solutions/:id/ctas/:ctaId", adminSolutionsHandler.DeleteCTA)

	// Whitepapers admin
	adminWhitepapersHandler := adminHandlers.NewWhitepapersHandler(queries, testLogger, appCache)
	adminGroup.GET("/whitepapers", adminWhitepapersHandler.List)
	adminGroup.GET("/whitepapers/new", adminWhitepapersHandler.New)
	adminGroup.POST("/whitepapers", adminWhitepapersHandler.Create)
	adminGroup.GET("/whitepapers/:id/edit", adminWhitepapersHandler.Edit)
	adminGroup.POST("/whitepapers/:id", adminWhitepapersHandler.Update)
	adminGroup.DELETE("/whitepapers/:id", adminWhitepapersHandler.Delete)
	adminGroup.GET("/whitepapers/:id/downloads", adminWhitepapersHandler.Downloads)

	// Homepage admin
	homepageAdminHandler := adminHandlers.NewHomepageHandler(queries, testLogger)
	adminGroup.GET("/homepage/heroes", homepageAdminHandler.HeroesList)
	adminGroup.GET("/homepage/heroes/new", homepageAdminHandler.HeroNew)
	adminGroup.POST("/homepage/heroes", homepageAdminHandler.HeroCreate)
	adminGroup.GET("/homepage/heroes/:id/edit", homepageAdminHandler.HeroEdit)
	adminGroup.POST("/homepage/heroes/:id", homepageAdminHandler.HeroUpdate)
	adminGroup.DELETE("/homepage/heroes/:id", homepageAdminHandler.HeroDelete)
	adminGroup.GET("/homepage/stats", homepageAdminHandler.StatsList)
	adminGroup.GET("/homepage/stats/new", homepageAdminHandler.StatNew)
	adminGroup.POST("/homepage/stats", homepageAdminHandler.StatCreate)
	adminGroup.GET("/homepage/stats/:id/edit", homepageAdminHandler.StatEdit)
	adminGroup.POST("/homepage/stats/:id", homepageAdminHandler.StatUpdate)
	adminGroup.DELETE("/homepage/stats/:id", homepageAdminHandler.StatDelete)
	adminGroup.GET("/homepage/testimonials", homepageAdminHandler.TestimonialsList)
	adminGroup.GET("/homepage/testimonials/new", homepageAdminHandler.TestimonialNew)
	adminGroup.POST("/homepage/testimonials", homepageAdminHandler.TestimonialCreate)
	adminGroup.GET("/homepage/testimonials/:id/edit", homepageAdminHandler.TestimonialEdit)
	adminGroup.POST("/homepage/testimonials/:id", homepageAdminHandler.TestimonialUpdate)
	adminGroup.DELETE("/homepage/testimonials/:id", homepageAdminHandler.TestimonialDelete)
	adminGroup.GET("/homepage/cta", homepageAdminHandler.CTAList)
	adminGroup.GET("/homepage/cta/new", homepageAdminHandler.CTANew)
	adminGroup.POST("/homepage/cta", homepageAdminHandler.CTACreate)
	adminGroup.GET("/homepage/cta/:id/edit", homepageAdminHandler.CTAEdit)
	adminGroup.POST("/homepage/cta/:id", homepageAdminHandler.CTAUpdate)
	adminGroup.DELETE("/homepage/cta/:id", homepageAdminHandler.CTADelete)
	adminGroup.GET("/homepage/settings", homepageAdminHandler.Settings)
	adminGroup.POST("/homepage/settings", homepageAdminHandler.UpdateSettings)

	// Section settings
	sectionSettingsHandler := adminHandlers.NewSectionSettingsHandler(queries, testLogger)
	adminGroup.GET("/about/settings", sectionSettingsHandler.AboutSettings)
	adminGroup.POST("/about/settings", sectionSettingsHandler.UpdateAboutSettings)
	adminGroup.GET("/products/settings", sectionSettingsHandler.ProductsSettings)
	adminGroup.POST("/products/settings", sectionSettingsHandler.UpdateProductsSettings)
	adminGroup.GET("/solutions/settings", sectionSettingsHandler.SolutionsSettings)
	adminGroup.POST("/solutions/settings", sectionSettingsHandler.UpdateSolutionsSettings)
	adminGroup.GET("/blog/settings", sectionSettingsHandler.BlogSettings)
	adminGroup.POST("/blog/settings", sectionSettingsHandler.UpdateBlogSettings)

	// About admin
	adminAboutHandler := adminHandlers.NewAboutHandler(queries, testLogger, appCache)
	adminGroup.GET("/about/overview", adminAboutHandler.OverviewEdit)
	adminGroup.POST("/about/overview", adminAboutHandler.OverviewUpdate)
	adminGroup.GET("/about/mvv", adminAboutHandler.MVVEdit)
	adminGroup.POST("/about/mvv", adminAboutHandler.MVVUpdate)
	adminGroup.GET("/about/values", adminAboutHandler.CoreValuesList)
	adminGroup.GET("/about/values/new", adminAboutHandler.CoreValueNew)
	adminGroup.POST("/about/values", adminAboutHandler.CoreValueCreate)
	adminGroup.GET("/about/values/:id/edit", adminAboutHandler.CoreValueEdit)
	adminGroup.POST("/about/values/:id", adminAboutHandler.CoreValueUpdate)
	adminGroup.DELETE("/about/values/:id", adminAboutHandler.CoreValueDelete)
	adminGroup.GET("/about/milestones", adminAboutHandler.MilestonesList)
	adminGroup.GET("/about/milestones/new", adminAboutHandler.MilestoneNew)
	adminGroup.POST("/about/milestones", adminAboutHandler.MilestoneCreate)
	adminGroup.GET("/about/milestones/:id/edit", adminAboutHandler.MilestoneEdit)
	adminGroup.POST("/about/milestones/:id", adminAboutHandler.MilestoneUpdate)
	adminGroup.DELETE("/about/milestones/:id", adminAboutHandler.MilestoneDelete)
	adminGroup.GET("/about/certifications", adminAboutHandler.CertificationsList)
	adminGroup.GET("/about/certifications/new", adminAboutHandler.CertificationNew)
	adminGroup.POST("/about/certifications", adminAboutHandler.CertificationCreate)
	adminGroup.GET("/about/certifications/:id/edit", adminAboutHandler.CertificationEdit)
	adminGroup.POST("/about/certifications/:id", adminAboutHandler.CertificationUpdate)
	adminGroup.DELETE("/about/certifications/:id", adminAboutHandler.CertificationDelete)

	// Partners admin
	adminPartnersHandler := adminHandlers.NewPartnersHandler(queries, testLogger, appCache)
	adminGroup.GET("/partners", adminPartnersHandler.List)
	adminGroup.GET("/partners/new", adminPartnersHandler.New)
	adminGroup.POST("/partners", adminPartnersHandler.Create)
	adminGroup.GET("/partners/:id/edit", adminPartnersHandler.Edit)
	adminGroup.POST("/partners/:id", adminPartnersHandler.Update)
	adminGroup.DELETE("/partners/:id", adminPartnersHandler.Delete)
	adminGroup.GET("/partners/testimonials", adminPartnersHandler.TestimonialsList)
	adminGroup.GET("/partners/testimonials/new", adminPartnersHandler.TestimonialNew)
	adminGroup.POST("/partners/testimonials", adminPartnersHandler.TestimonialCreate)
	adminGroup.GET("/partners/testimonials/:id/edit", adminPartnersHandler.TestimonialEdit)
	adminGroup.POST("/partners/testimonials/:id", adminPartnersHandler.TestimonialUpdate)
	adminGroup.DELETE("/partners/testimonials/:id", adminPartnersHandler.TestimonialDelete)

	// Media library
	mediaHandler := adminHandlers.NewMediaHandler(queries, testLogger, t.TempDir())
	adminGroup.GET("/media", mediaHandler.List)
	adminGroup.POST("/media/upload", mediaHandler.Upload)
	adminGroup.GET("/media/browse", mediaHandler.Browse)
	adminGroup.GET("/media/:id", mediaHandler.GetFile)
	adminGroup.PUT("/media/:id", mediaHandler.UpdateAltText)
	adminGroup.DELETE("/media/:id", mediaHandler.Delete)

	// Navigation
	navHandler := adminHandlers.NewNavigationHandler(queries, testLogger)
	adminGroup.GET("/navigation", navHandler.List)
	adminGroup.POST("/navigation", navHandler.Create)
	adminGroup.GET("/navigation/:id", navHandler.Edit)
	adminGroup.POST("/navigation/:id/settings", navHandler.UpdateMenu)
	adminGroup.POST("/navigation/:id/items", navHandler.AddItem)
	adminGroup.POST("/navigation/items/:id", navHandler.UpdateItem)
	adminGroup.DELETE("/navigation/items/:id", navHandler.DeleteItem)
	adminGroup.DELETE("/navigation/:id", navHandler.DeleteMenu)
	adminGroup.POST("/navigation/:id/reorder", navHandler.Reorder)

	// Activity log
	activityHandler := adminHandlers.NewActivityHandler(queries, testLogger)
	adminGroup.GET("/activity", activityHandler.List)

	// Contact admin
	adminContactHandler := adminHandlers.NewAdminContactHandler(queries, testLogger, appCache)
	adminGroup.GET("/contact/submissions", adminContactHandler.ListSubmissions)
	adminGroup.GET("/contact/submissions/:id", adminContactHandler.ViewSubmission)
	adminGroup.POST("/contact/submissions/:id/status", adminContactHandler.UpdateSubmissionStatus)
	adminGroup.POST("/contact/submissions/bulk-mark-read", adminContactHandler.BulkMarkRead)
	adminGroup.DELETE("/contact/submissions/:id", adminContactHandler.DeleteSubmission)
	adminGroup.GET("/contact/offices", adminContactHandler.ListOffices)
	adminGroup.GET("/contact/offices/new", adminContactHandler.NewOffice)
	adminGroup.POST("/contact/offices", adminContactHandler.CreateOffice)
	adminGroup.GET("/contact/offices/:id/edit", adminContactHandler.EditOffice)
	adminGroup.POST("/contact/offices/:id", adminContactHandler.UpdateOffice)
	adminGroup.DELETE("/contact/offices/:id", adminContactHandler.DeleteOffice)

	// Case studies admin
	adminCaseStudiesHandler := adminHandlers.NewCaseStudiesHandler(queries, testLogger, appCache)
	adminGroup.GET("/case-studies", adminCaseStudiesHandler.List)
	adminGroup.GET("/case-studies/new", adminCaseStudiesHandler.New)
	adminGroup.POST("/case-studies", adminCaseStudiesHandler.Create)
	adminGroup.GET("/case-studies/:id/edit", adminCaseStudiesHandler.Edit)
	adminGroup.POST("/case-studies/:id", adminCaseStudiesHandler.Update)
	adminGroup.DELETE("/case-studies/:id", adminCaseStudiesHandler.Delete)
	adminGroup.POST("/case-studies/:id/products", adminCaseStudiesHandler.AddProduct)
	adminGroup.DELETE("/case-studies/:id/products/:productId", adminCaseStudiesHandler.RemoveProduct)
	adminGroup.POST("/case-studies/:id/metrics", adminCaseStudiesHandler.AddMetric)
	adminGroup.DELETE("/case-studies/:id/metrics/:metricId", adminCaseStudiesHandler.DeleteMetric)

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
