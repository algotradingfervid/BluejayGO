package admin

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type SettingsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewSettingsHandler(queries *sqlc.Queries, logger *slog.Logger) *SettingsHandler {
	return &SettingsHandler{queries: queries, logger: logger}
}

func (h *SettingsHandler) Edit(c echo.Context) error {
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	saved := c.QueryParam("saved") == "1"
	activeTab := c.QueryParam("tab")
	if activeTab == "" {
		activeTab = "general"
	}
	return c.Render(http.StatusOK, "admin/pages/settings_form.html", map[string]interface{}{
		"Title":     "Global Settings",
		"Settings":  settings,
		"Saved":     saved,
		"ActiveTab": activeTab,
	})
}

func (h *SettingsHandler) Update(c echo.Context) error {
	activeTab := c.FormValue("active_tab")
	if activeTab == "" {
		activeTab = "general"
	}

	err := h.queries.UpdateGlobalSettings(c.Request().Context(), sqlc.UpdateGlobalSettingsParams{
		SiteName:          c.FormValue("site_name"),
		SiteTagline:       c.FormValue("site_tagline"),
		ContactEmail:      c.FormValue("contact_email"),
		ContactPhone:      c.FormValue("contact_phone"),
		Address:           c.FormValue("address"),
		BusinessHours:     c.FormValue("business_hours"),
		MetaDescription:   c.FormValue("meta_description"),
		MetaKeywords:      c.FormValue("meta_keywords"),
		GoogleAnalyticsID: c.FormValue("google_analytics_id"),
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
	logActivity(c, "updated", "settings", 0, "", "Updated Global Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/settings?saved=1&tab="+activeTab)
}
