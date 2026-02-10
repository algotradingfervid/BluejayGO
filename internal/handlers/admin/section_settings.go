package admin

import (
	// Standard library imports
	"log/slog"  // Structured logging for error tracking and debugging
	"net/http"  // HTTP status codes and request/response handling
	"strconv"   // String to integer conversion for parsing numeric form fields
	"strings"   // String manipulation (trimming whitespace from form inputs)

	// Third-party framework
	"github.com/labstack/echo/v4" // Echo web framework for HTTP routing and context management

	// Internal dependencies
	"github.com/narendhupati/bluejay-cms/db/sqlc" // sqlc-generated database queries and models
)

// SectionSettingsHandler manages section-specific settings for different areas of the site.
// Handles configuration for About, Products, Solutions, and Blog sections.
// Each section has its own settings that control display options, pagination, and feature toggles.
type SectionSettingsHandler struct {
	queries *sqlc.Queries // Database query interface for section settings CRUD operations
	logger  *slog.Logger  // Structured logger for error tracking
}

// NewSectionSettingsHandler creates and initializes a new SectionSettingsHandler instance.
// Dependencies are injected to support database access and logging.
func NewSectionSettingsHandler(queries *sqlc.Queries, logger *slog.Logger) *SectionSettingsHandler {
	return &SectionSettingsHandler{queries: queries, logger: logger}
}

// ==================== ABOUT SETTINGS ====================
// About section settings control which components are displayed on the About page
// (mission statement, milestones timeline, certifications, team members)

// AboutSettings renders the About section settings form.
//
// HTTP Method: GET
// Route: /admin/about/settings
// Template: templates/admin/pages/about_settings.html (with admin-layout wrapper)
// HTMX: Returns full HTML page (not a fragment)
//
// Displays toggle controls for About page components:
// - about_show_mission: Display mission statement section
// - about_show_milestones: Display milestones timeline
// - about_show_certifications: Display certifications/accreditations
// - about_show_team: Display team members section
//
// Query Parameters:
// - saved: Set to "1" after successful update to show success message
//
// Template Data:
// - Title: Page title ("About Settings")
// - Settings: Current settings row from database (includes all section settings)
// - Saved: Boolean flag to display success banner
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) AboutSettings(c echo.Context) error {
	// Fetch current settings from database (single row table with all settings)
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Check for success flag from previous update operation
	saved := c.QueryParam("saved") == "1"

	// Render About settings form with current toggle states
	// Template path: templates/admin/pages/about_settings.html
	return c.Render(http.StatusOK, "admin/pages/about_settings.html", map[string]interface{}{
		"Title":    "About Settings",
		"Settings": settings,  // Contains about_show_* fields (int64: 0=off, 1=on)
		"Saved":    saved,     // Triggers success message banner
	})
}

// UpdateAboutSettings processes the About settings form submission and persists changes.
//
// HTTP Method: POST
// Route: /admin/about/settings
// Form Fields: about_show_mission, about_show_milestones, about_show_certifications, about_show_team
// HTMX: Not used - standard form POST with redirect
//
// Form Field Behavior:
// - Checkboxes send "on" when checked, no value when unchecked
// - boolToInt helper converts: "on" -> 1 (true), missing -> 0 (false)
// - Database stores boolean toggles as int64 (0 or 1)
//
// Post-Update Behavior:
// - Logs activity to activity_log table for audit trail
// - Redirects back to settings form with saved=1 flag (shows success message)
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) UpdateAboutSettings(c echo.Context) error {
	// Helper function to convert HTML checkbox values to database int64 boolean
	// HTML checkboxes send "on" when checked, no value when unchecked
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1 // Checkbox is checked (enabled)
		}
		return 0 // Checkbox is unchecked (disabled)
	}

	// Update About section settings in database
	// Only updates about_show_* fields, leaves other settings unchanged
	err := h.queries.UpdateAboutSettings(c.Request().Context(), sqlc.UpdateAboutSettingsParams{
		AboutShowMission:        boolToInt("about_show_mission"),        // Toggle mission statement section
		AboutShowMilestones:     boolToInt("about_show_milestones"),     // Toggle milestones timeline
		AboutShowCertifications: boolToInt("about_show_certifications"), // Toggle certifications section
		AboutShowTeam:           boolToInt("about_show_team"),           // Toggle team members section
	})
	if err != nil {
		h.logger.Error("failed to update about settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Log settings update to activity_log for audit trail
	logActivity(c, "updated", "about_settings", 0, "", "Updated About Settings")

	// Redirect back to settings form with success flag
	return c.Redirect(http.StatusSeeOther, "/admin/about/settings?saved=1")
}

// ==================== PRODUCTS SETTINGS ====================
// Products section settings control pagination, filtering, and sorting options
// for the products listing page

// ProductsSettings renders the Products section settings form.
//
// HTTP Method: GET
// Route: /admin/products/settings
// Template: templates/admin/pages/products_settings.html (with admin-layout wrapper)
// HTMX: Returns full HTML page (not a fragment)
//
// Displays controls for Products page behavior:
// - products_per_page: Number of products to display per page (integer)
// - products_show_categories: Display category filter sidebar (toggle)
// - products_show_search: Display search box (toggle)
// - products_default_sort: Default sort order (dropdown: name, date, price, etc.)
//
// Query Parameters:
// - saved: Set to "1" after successful update to show success message
//
// Template Data:
// - Title: Page title ("Products Settings")
// - Settings: Current settings row from database
// - Saved: Boolean flag to display success banner
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) ProductsSettings(c echo.Context) error {
	// Fetch current settings from database
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Check for success flag from previous update operation
	saved := c.QueryParam("saved") == "1"

	// Render Products settings form
	// Template path: templates/admin/pages/products_settings.html
	return c.Render(http.StatusOK, "admin/pages/products_settings.html", map[string]interface{}{
		"Title":    "Products Settings",
		"Settings": settings,  // Contains products_* configuration fields
		"Saved":    saved,     // Triggers success message banner
	})
}

// UpdateProductsSettings processes the Products settings form submission and persists changes.
//
// HTTP Method: POST
// Route: /admin/products/settings
// Form Fields: products_per_page, products_show_categories, products_show_search, products_default_sort
// HTMX: Not used - standard form POST with redirect
//
// Form Field Processing:
// - products_per_page: Numeric input, parsed as int64 (defaults to 12 if invalid/empty)
// - products_show_categories: Checkbox ("on" -> 1, missing -> 0)
// - products_show_search: Checkbox ("on" -> 1, missing -> 0)
// - products_default_sort: Text input (e.g., "name_asc", "date_desc", "price_asc")
//
// Helper Functions:
// - parseIntField: Safely converts string to int64, returns default value on error/empty
// - boolToInt: Converts checkbox "on" value to 1, missing to 0
//
// Post-Update Behavior:
// - Logs activity to activity_log table for audit trail
// - Redirects back to settings form with saved=1 flag (shows success message)
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) UpdateProductsSettings(c echo.Context) error {
	// Helper function to safely parse integer form fields with fallback default
	parseIntField := func(field string, defaultVal int64) int64 {
		v := strings.TrimSpace(c.FormValue(field)) // Remove whitespace
		if v == "" {
			return defaultVal // Return default if field is empty
		}
		n, err := strconv.ParseInt(v, 10, 64) // Parse as base-10 int64
		if err != nil {
			return defaultVal // Return default if parsing fails
		}
		return n
	}

	// Helper function to convert HTML checkbox values to database int64 boolean
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1 // Checkbox is checked
		}
		return 0 // Checkbox is unchecked
	}

	// Update Products section settings in database
	err := h.queries.UpdateProductsSettings(c.Request().Context(), sqlc.UpdateProductsSettingsParams{
		ProductsPerPage:        parseIntField("products_per_page", 12), // Default: 12 products per page
		ProductsShowCategories: boolToInt("products_show_categories"),  // Toggle category filter
		ProductsShowSearch:     boolToInt("products_show_search"),      // Toggle search box
		ProductsDefaultSort:    c.FormValue("products_default_sort"),   // Default sort order string
	})
	if err != nil {
		h.logger.Error("failed to update products settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Log settings update to activity_log for audit trail
	logActivity(c, "updated", "products_settings", 0, "", "Updated Products Settings")

	// Redirect back to settings form with success flag
	return c.Redirect(http.StatusSeeOther, "/admin/products/settings?saved=1")
}

// ==================== SOLUTIONS SETTINGS ====================
// Solutions section settings control pagination and filtering options
// for the solutions listing page

// SolutionsSettings renders the Solutions section settings form.
//
// HTTP Method: GET
// Route: /admin/solutions/settings
// Template: templates/admin/pages/solutions_settings.html (with admin-layout wrapper)
// HTMX: Returns full HTML page (not a fragment)
//
// Displays controls for Solutions page behavior:
// - solutions_per_page: Number of solutions to display per page (integer)
// - solutions_show_industries: Display industry filter sidebar (toggle)
// - solutions_show_search: Display search box (toggle)
//
// Query Parameters:
// - saved: Set to "1" after successful update to show success message
//
// Template Data:
// - Title: Page title ("Solutions Settings")
// - Settings: Current settings row from database
// - Saved: Boolean flag to display success banner
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) SolutionsSettings(c echo.Context) error {
	// Fetch current settings from database
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Check for success flag from previous update operation
	saved := c.QueryParam("saved") == "1"

	// Render Solutions settings form
	// Template path: templates/admin/pages/solutions_settings.html
	return c.Render(http.StatusOK, "admin/pages/solutions_settings.html", map[string]interface{}{
		"Title":    "Solutions Settings",
		"Settings": settings,  // Contains solutions_* configuration fields
		"Saved":    saved,     // Triggers success message banner
	})
}

// UpdateSolutionsSettings processes the Solutions settings form submission and persists changes.
//
// HTTP Method: POST
// Route: /admin/solutions/settings
// Form Fields: solutions_per_page, solutions_show_industries, solutions_show_search
// HTMX: Not used - standard form POST with redirect
//
// Form Field Processing:
// - solutions_per_page: Numeric input, parsed as int64 (defaults to 12 if invalid/empty)
// - solutions_show_industries: Checkbox ("on" -> 1, missing -> 0)
// - solutions_show_search: Checkbox ("on" -> 1, missing -> 0)
//
// Helper Functions:
// - parseIntField: Safely converts string to int64, returns default value on error/empty
// - boolToInt: Converts checkbox "on" value to 1, missing to 0
//
// Post-Update Behavior:
// - Logs activity to activity_log table for audit trail
// - Redirects back to settings form with saved=1 flag (shows success message)
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) UpdateSolutionsSettings(c echo.Context) error {
	// Helper function to safely parse integer form fields with fallback default
	parseIntField := func(field string, defaultVal int64) int64 {
		v := strings.TrimSpace(c.FormValue(field)) // Remove whitespace
		if v == "" {
			return defaultVal // Return default if field is empty
		}
		n, err := strconv.ParseInt(v, 10, 64) // Parse as base-10 int64
		if err != nil {
			return defaultVal // Return default if parsing fails
		}
		return n
	}

	// Helper function to convert HTML checkbox values to database int64 boolean
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1 // Checkbox is checked
		}
		return 0 // Checkbox is unchecked
	}

	// Update Solutions section settings in database
	err := h.queries.UpdateSolutionsSettings(c.Request().Context(), sqlc.UpdateSolutionsSettingsParams{
		SolutionsPerPage:        parseIntField("solutions_per_page", 12), // Default: 12 solutions per page
		SolutionsShowIndustries: boolToInt("solutions_show_industries"),  // Toggle industry filter
		SolutionsShowSearch:     boolToInt("solutions_show_search"),      // Toggle search box
	})
	if err != nil {
		h.logger.Error("failed to update solutions settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Log settings update to activity_log for audit trail
	logActivity(c, "updated", "solutions_settings", 0, "", "Updated Solutions Settings")

	// Redirect back to settings form with success flag
	return c.Redirect(http.StatusSeeOther, "/admin/solutions/settings?saved=1")
}

// ==================== BLOG SETTINGS ====================
// Blog section settings control pagination, metadata display, and filtering options
// for the blog listing and individual blog post pages

// BlogSettings renders the Blog section settings form.
//
// HTTP Method: GET
// Route: /admin/blog/settings
// Template: templates/admin/pages/blog_settings.html (with admin-layout wrapper)
// HTMX: Returns full HTML page (not a fragment)
//
// Displays controls for Blog page behavior:
// - blog_posts_per_page: Number of blog posts to display per page (integer)
// - blog_show_author: Display author name on posts (toggle)
// - blog_show_date: Display publication date on posts (toggle)
// - blog_show_categories: Display category filter/tags (toggle)
// - blog_show_tags: Display article tags (toggle)
// - blog_show_search: Display search box (toggle)
//
// Query Parameters:
// - saved: Set to "1" after successful update to show success message
//
// Template Data:
// - Title: Page title ("Blog Settings")
// - Settings: Current settings row from database
// - Saved: Boolean flag to display success banner
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) BlogSettings(c echo.Context) error {
	// Fetch current settings from database
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Check for success flag from previous update operation
	saved := c.QueryParam("saved") == "1"

	// Render Blog settings form
	// Template path: templates/admin/pages/blog_settings.html
	return c.Render(http.StatusOK, "admin/pages/blog_settings.html", map[string]interface{}{
		"Title":    "Blog Settings",
		"Settings": settings,  // Contains blog_* configuration fields
		"Saved":    saved,     // Triggers success message banner
	})
}

// UpdateBlogSettings processes the Blog settings form submission and persists changes.
//
// HTTP Method: POST
// Route: /admin/blog/settings
// Form Fields: blog_posts_per_page, blog_show_author, blog_show_date,
//              blog_show_categories, blog_show_tags, blog_show_search
// HTMX: Not used - standard form POST with redirect
//
// Form Field Processing:
// - blog_posts_per_page: Numeric input, parsed as int64 (defaults to 10 if invalid/empty)
// - blog_show_author: Checkbox ("on" -> 1, missing -> 0)
// - blog_show_date: Checkbox ("on" -> 1, missing -> 0)
// - blog_show_categories: Checkbox ("on" -> 1, missing -> 0)
// - blog_show_tags: Checkbox ("on" -> 1, missing -> 0)
// - blog_show_search: Checkbox ("on" -> 1, missing -> 0)
//
// Helper Functions:
// - parseIntField: Safely converts string to int64, returns default value on error/empty
// - boolToInt: Converts checkbox "on" value to 1, missing to 0
//
// Post-Update Behavior:
// - Logs activity to activity_log table for audit trail
// - Redirects back to settings form with saved=1 flag (shows success message)
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SectionSettingsHandler) UpdateBlogSettings(c echo.Context) error {
	// Helper function to safely parse integer form fields with fallback default
	parseIntField := func(field string, defaultVal int64) int64 {
		v := strings.TrimSpace(c.FormValue(field)) // Remove whitespace
		if v == "" {
			return defaultVal // Return default if field is empty
		}
		n, err := strconv.ParseInt(v, 10, 64) // Parse as base-10 int64
		if err != nil {
			return defaultVal // Return default if parsing fails
		}
		return n
	}

	// Helper function to convert HTML checkbox values to database int64 boolean
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1 // Checkbox is checked
		}
		return 0 // Checkbox is unchecked
	}

	// Update Blog section settings in database
	err := h.queries.UpdateBlogSettings(c.Request().Context(), sqlc.UpdateBlogSettingsParams{
		BlogPostsPerPage:   parseIntField("blog_posts_per_page", 10), // Default: 10 posts per page
		BlogShowAuthor:     boolToInt("blog_show_author"),            // Toggle author name display
		BlogShowDate:       boolToInt("blog_show_date"),              // Toggle publication date display
		BlogShowCategories: boolToInt("blog_show_categories"),        // Toggle category filter/tags
		BlogShowTags:       boolToInt("blog_show_tags"),              // Toggle article tags display
		BlogShowSearch:     boolToInt("blog_show_search"),            // Toggle search box
	})
	if err != nil {
		h.logger.Error("failed to update blog settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Log settings update to activity_log for audit trail
	logActivity(c, "updated", "blog_settings", 0, "", "Updated Blog Settings")

	// Redirect back to settings form with success flag
	return c.Redirect(http.StatusSeeOther, "/admin/blog/settings?saved=1")
}
