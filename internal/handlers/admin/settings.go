package admin

import (
	// Standard library imports
	"log/slog"  // Structured logging for error tracking and debugging
	"net/http"  // HTTP status codes and request/response handling

	// Third-party framework
	"github.com/labstack/echo/v4" // Echo web framework for HTTP routing and context management

	// Internal dependencies
	"github.com/narendhupati/bluejay-cms/db/sqlc" // sqlc-generated database queries and models
)

// SettingsHandler handles global site settings management.
// Manages site-wide configuration including contact info, SEO metadata,
// analytics integration, and social media links.
type SettingsHandler struct {
	queries *sqlc.Queries // Database query interface for settings CRUD operations
	logger  *slog.Logger  // Structured logger for error tracking
}

// NewSettingsHandler creates and initializes a new SettingsHandler instance.
// Dependencies are injected to support database access and logging.
func NewSettingsHandler(queries *sqlc.Queries, logger *slog.Logger) *SettingsHandler {
	return &SettingsHandler{queries: queries, logger: logger}
}

// Edit renders the global settings form page with current settings data.
//
// HTTP Method: GET
// Route: /admin/settings
// Template: templates/admin/pages/settings_form.html (with admin-layout wrapper)
// HTMX: Returns full HTML page (not a fragment)
//
// Displays a tabbed interface for editing global site settings including:
// - General: Site name, tagline, contact info, business hours
// - SEO: Meta description, meta keywords, Google Analytics ID
// - Social: Facebook, Twitter, LinkedIn, Instagram, YouTube links
//
// Query Parameters:
// - saved: Set to "1" after successful update to show success message
// - tab: Active tab identifier (general, seo, social) - defaults to "general"
//
// Template Data:
// - Title: Page title ("Global Settings")
// - Settings: Current settings row from database (all fields)
// - Saved: Boolean flag to display success banner
// - ActiveTab: Which tab should be displayed/highlighted
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SettingsHandler) Edit(c echo.Context) error {
	// Fetch current settings from database (single row table)
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Check for success flag from previous update operation
	saved := c.QueryParam("saved") == "1"

	// Determine which tab should be active (preserves tab state after form submission)
	activeTab := c.QueryParam("tab")
	if activeTab == "" {
		activeTab = "general" // Default to general tab if not specified
	}

	// Render settings form with current data and UI state
	// Template path: templates/admin/pages/settings_form.html
	// Uses admin-layout wrapper for consistent navigation/header
	return c.Render(http.StatusOK, "admin/pages/settings_form.html", map[string]interface{}{
		"Title":     "Global Settings",
		"Settings":  settings,           // Current settings data from database
		"Saved":     saved,               // Show success message if true
		"ActiveTab": activeTab,           // Determines which tab is visible/active
	})
}

// Update processes the global settings form submission and persists changes to database.
//
// HTTP Method: POST
// Route: /admin/settings
// Form Fields: All settings fields (site_name, contact_email, meta_description, etc.)
//              plus active_tab (hidden field to preserve tab state)
// HTMX: Not used - standard form POST with redirect
//
// Updates all global settings fields in a single database operation.
// Settings table is a single-row table (no ID needed for update).
//
// Form Fields Processed:
// General Tab:
// - site_name: Site name/title
// - site_tagline: Site tagline/slogan
// - contact_email: Primary contact email
// - contact_phone: Primary contact phone number
// - address: Physical business address
// - business_hours: Operating hours text
//
// SEO Tab:
// - meta_description: Default meta description for SEO
// - meta_keywords: Default meta keywords for SEO
// - google_analytics_id: Google Analytics tracking ID
//
// Social Tab:
// - social_facebook: Facebook profile URL
// - social_twitter: Twitter profile URL
// - social_linkedin: LinkedIn profile URL
// - social_instagram: Instagram profile URL
// - social_youtube: YouTube channel URL
//
// Post-Update Behavior:
// - Logs activity to activity_log table for audit trail
// - Redirects back to settings form with saved=1 flag (shows success message)
// - Preserves active tab state in redirect URL
//
// Authentication: Requires valid session (enforced by middleware)
func (h *SettingsHandler) Update(c echo.Context) error {
	// Extract active tab from hidden form field to preserve UI state after redirect
	activeTab := c.FormValue("active_tab")
	if activeTab == "" {
		activeTab = "general" // Default to general tab if not specified
	}

	// Update all global settings fields in database (single UPDATE query)
	// Settings table contains one row with all global configuration
	err := h.queries.UpdateGlobalSettings(c.Request().Context(), sqlc.UpdateGlobalSettingsParams{
		// General settings
		SiteName:          c.FormValue("site_name"),
		SiteTagline:       c.FormValue("site_tagline"),
		ContactEmail:      c.FormValue("contact_email"),
		ContactPhone:      c.FormValue("contact_phone"),
		Address:           c.FormValue("address"),
		BusinessHours:     c.FormValue("business_hours"),

		// SEO settings
		MetaDescription:   c.FormValue("meta_description"),
		MetaKeywords:      c.FormValue("meta_keywords"),
		GoogleAnalyticsID: c.FormValue("google_analytics_id"),

		// Social media links
		SocialFacebook:    c.FormValue("social_facebook"),
		SocialTwitter:     c.FormValue("social_twitter"),
		SocialLinkedin:    c.FormValue("social_linkedin"),
		SocialInstagram:   c.FormValue("social_instagram"),
		SocialYoutube:     c.FormValue("social_youtube"),
	})
	if err != nil {
		h.logger.Error("failed to update settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Log settings update to activity_log for audit trail
	// entity_id=0 since settings is a singleton (no specific ID)
	logActivity(c, "updated", "settings", 0, "", "Updated Global Settings")

	// Redirect back to settings form with success flag and preserved tab state
	// saved=1 triggers success message banner in template
	// tab parameter ensures same tab is displayed after update
	return c.Redirect(http.StatusSeeOther, "/admin/settings?saved=1&tab="+activeTab)
}
