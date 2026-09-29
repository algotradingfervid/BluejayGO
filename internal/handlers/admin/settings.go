package admin

import (
	// Standard library imports
	"log/slog" // Structured logging for error tracking and debugging
	"net/http"
	"net/url"
	"strings"

	// Third-party framework
	"github.com/labstack/echo/v4" // Echo web framework for HTTP routing and context management

	// Internal dependencies
	"github.com/narendhupati/bluejay-cms/db/sqlc"           // sqlc-generated database queries and models
	"github.com/narendhupati/bluejay-cms/internal/services" // Cache service for invalidating page-level caches
)

// SettingsHandler handles global site settings management.
// Manages site-wide configuration including contact info, SEO metadata,
// analytics integration, and social media links.
type SettingsHandler struct {
	queries *sqlc.Queries   // Database query interface for settings CRUD operations
	logger  *slog.Logger    // Structured logger for error tracking
	cache   *services.Cache // Cache service for invalidating page-level caches
}

// NewSettingsHandler creates and initializes a new SettingsHandler instance.
// Dependencies are injected to support database access, logging, and cache invalidation.
func NewSettingsHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *SettingsHandler {
	return &SettingsHandler{queries: queries, logger: logger, cache: cache}
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
		return renderOperationError(c, "Settings unavailable", "We could not complete this request. Try again, or contact your site administrator if this continues.", "/admin/settings")
	}

	media, err := h.queries.ListMediaFiles(c.Request().Context(), sqlc.ListMediaFilesParams{Limit: -1})
	if err != nil {
		h.logger.Error("failed to load settings image choices", "error", err)
		return renderOperationError(c, "Settings unavailable", "We could not load the image library. Try again, or contact your site administrator if this continues.", "/admin/settings")
	}
	var images []sqlc.MediaFile
	for _, file := range media {
		if strings.HasPrefix(file.MimeType, "image/") {
			images = append(images, file)
		}
	}

	// Check for success flag from previous update operation
	saved := c.QueryParam("saved") == "1"

	// Determine which tab should be active (preserves tab state after form submission)
	activeTab := c.QueryParam("tab")
	if !validSettingsTab(activeTab) {
		activeTab = "general" // Default to general tab if not specified
	}

	// Render settings form with current data and UI state
	// Template path: templates/admin/pages/settings_form.html
	// Uses admin-layout wrapper for consistent navigation/header
	return c.Render(http.StatusOK, "admin/pages/settings_form.html", map[string]interface{}{
		"Title":        "Global Settings",
		"ImageChoices": images,
		"Settings":     settings, // Current settings data from database
		"Saved":        saved,    // Show success message if true
		"SocialFields": settingsSocialFields(settings),
		"ActiveTab":    activeTab, // Determines which tab is visible/active
	})
}

// Update processes the global settings form submission and persists changes to database.
//
// HTTP Method: POST
// Route: /admin/settings
// Form Fields: All settings fields (site_name, contact_email, meta_description, etc.)
//
//	plus active_tab (hidden field to preserve tab state)
//
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
	if !validSettingsTab(activeTab) {
		activeTab = "general" // Default to general tab if not specified
	}

	// Reject unsafe destinations before any settings are changed.
	for _, field := range []string{"social_facebook", "social_twitter", "social_linkedin", "social_instagram", "social_youtube", "social_threads", "marketplace_gem_url", "marketplace_amazon_url"} {
		value := strings.TrimSpace(c.FormValue(field))
		if value != "" && !validSettingsURL(value) {
			return echo.NewHTTPError(http.StatusBadRequest, "Enter a full http:// or https:// URL for "+strings.ReplaceAll(field, "_", " "))
		}
	}
	logoPath := strings.TrimSpace(c.FormValue("footer_logo_path"))
	if logoPath != "" && !validSettingsURL(logoPath) && !(strings.HasPrefix(logoPath, "/") && !strings.HasPrefix(logoPath, "//") && !strings.ContainsAny(logoPath, "\\\r\n")) {
		return echo.NewHTTPError(http.StatusBadRequest, "Footer logo must be an image path starting with / or a full http:// or https:// URL")
	}

	ogImage := strings.TrimSpace(c.FormValue("default_og_image"))
	if ogImage != "" {
		image, err := h.queries.GetMediaFileByPath(c.Request().Context(), ogImage)
		if err != nil || !strings.HasPrefix(image.MimeType, "image/") {
			return echo.NewHTTPError(http.StatusBadRequest, "Choose an image from the media library, or remove the default image")
		}
	}

	// Update all global settings fields in database (single UPDATE query)
	// Settings table contains one row with all global configuration
	err := h.queries.UpdateGlobalSettings(c.Request().Context(), sqlc.UpdateGlobalSettingsParams{
		// General settings
		SiteName:      c.FormValue("site_name"),
		SiteTagline:   c.FormValue("site_tagline"),
		ContactEmail:  c.FormValue("contact_email"),
		ContactPhone:  c.FormValue("contact_phone"),
		Address:       c.FormValue("address"),
		BusinessHours: c.FormValue("business_hours"),

		// SEO settings
		MetaDescription:   c.FormValue("meta_description"),
		MetaKeywords:      c.FormValue("meta_keywords"),
		GoogleAnalyticsID: c.FormValue("google_analytics_id"),
		DefaultOgImage:    ogImage,

		// Social media links
		SocialFacebook:       strings.TrimSpace(c.FormValue("social_facebook")),
		SocialTwitter:        strings.TrimSpace(c.FormValue("social_twitter")),
		SocialLinkedin:       strings.TrimSpace(c.FormValue("social_linkedin")),
		SocialInstagram:      strings.TrimSpace(c.FormValue("social_instagram")),
		SocialYoutube:        strings.TrimSpace(c.FormValue("social_youtube")),
		SocialThreads:        strings.TrimSpace(c.FormValue("social_threads")),
		FooterLogoPath:       logoPath,
		MarketplaceGemUrl:    strings.TrimSpace(c.FormValue("marketplace_gem_url")),
		MarketplaceAmazonUrl: strings.TrimSpace(c.FormValue("marketplace_amazon_url")),
	})
	if err != nil {
		h.logger.Error("failed to update settings", "error", err)
		return renderOperationError(c, "Settings unavailable", "We could not complete this request. Try again, or contact your site administrator if this continues.", "/admin/settings")
	}

	// Invalidate ALL cached public pages. Global settings (contact info, site name,
	// etc.) render site-wide via the footer/header, so every "page:" cache entry is
	// potentially stale after a settings change.
	if h.cache != nil {
		h.cache.DeleteByPrefix("page:")
	}

	// Log settings update to activity_log for audit trail
	// entity_id=0 since settings is a singleton (no specific ID)
	logActivity(c, "updated", "settings", 0, "", "Updated Global Settings")

	// Redirect back to settings form with success flag and preserved tab state
	// saved=1 triggers success message banner in template
	// tab parameter ensures same tab is displayed after update
	return c.Redirect(http.StatusSeeOther, "/admin/settings?saved=1&tab="+activeTab)
}

// Social fields are data rather than a template slice call: Go's built-in slice
// extracts part of an existing slice and cannot construct a list of platforms.
type settingsSocialField struct {
	Label       string
	Name        string
	Value       string
	Placeholder string
}

func settingsSocialFields(settings sqlc.Setting) []settingsSocialField {
	return []settingsSocialField{
		{"Facebook", "social_facebook", settings.SocialFacebook, "https://facebook.com/yourpage"},
		{"Twitter / X", "social_twitter", settings.SocialTwitter, "https://x.com/yourhandle"},
		{"LinkedIn", "social_linkedin", settings.SocialLinkedin, "https://linkedin.com/company/yourco"},
		{"Instagram", "social_instagram", settings.SocialInstagram, "https://instagram.com/yourhandle"},
		{"Threads", "social_threads", settings.SocialThreads, "https://www.threads.net/@yourhandle"},
		{"YouTube", "social_youtube", settings.SocialYoutube, "https://youtube.com/@yourchannel"},
	}
}

func validSettingsURL(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && (parsed.Scheme == "https" || parsed.Scheme == "http") && parsed.Hostname() != "" && parsed.User == nil
}

func validSettingsTab(tab string) bool {
	switch tab {
	case "general", "contact", "social", "seo", "marketplaces":
		return true
	default:
		return false
	}
}
