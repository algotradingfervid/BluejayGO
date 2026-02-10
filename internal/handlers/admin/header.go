package admin

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type HeaderHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewHeaderHandler(queries *sqlc.Queries, logger *slog.Logger) *HeaderHandler {
	return &HeaderHandler{queries: queries, logger: logger}
}

func (h *HeaderHandler) Edit(c echo.Context) error {
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/header_form.html", map[string]interface{}{
		"Title":    "Header Management",
		"Settings": settings,
		"Saved":    saved,
	})
}

func (h *HeaderHandler) Update(c echo.Context) error {
	err := h.queries.UpdateHeaderSettings(c.Request().Context(), sqlc.UpdateHeaderSettingsParams{
		HeaderLogoPath:     c.FormValue("header_logo_path"),
		HeaderLogoAlt:      c.FormValue("header_logo_alt"),
		HeaderCtaEnabled:   c.FormValue("header_cta_enabled") == "on",
		HeaderCtaText:      c.FormValue("header_cta_text"),
		HeaderCtaUrl:       c.FormValue("header_cta_url"),
		HeaderCtaStyle:     c.FormValue("header_cta_style"),
		HeaderShowPhone:    c.FormValue("header_show_phone") == "on",
		HeaderShowEmail:    c.FormValue("header_show_email") == "on",
		HeaderShowSocial:   c.FormValue("header_show_social") == "on",
		HeaderSocialStyle:  c.FormValue("header_social_style"),
		ShowNavProducts:    c.FormValue("show_nav_products") == "on",
		ShowNavSolutions:   c.FormValue("show_nav_solutions") == "on",
		ShowNavCaseStudies: c.FormValue("show_nav_case_studies") == "on",
		ShowNavAbout:       c.FormValue("show_nav_about") == "on",
		ShowNavBlog:        c.FormValue("show_nav_blog") == "on",
		ShowNavWhitepapers: c.FormValue("show_nav_whitepapers") == "on",
		ShowNavPartners:    c.FormValue("show_nav_partners") == "on",
		ShowNavContact:     c.FormValue("show_nav_contact") == "on",
		NavLabelProducts:    c.FormValue("nav_label_products"),
		NavLabelSolutions:   c.FormValue("nav_label_solutions"),
		NavLabelCaseStudies: c.FormValue("nav_label_case_studies"),
		NavLabelAbout:       c.FormValue("nav_label_about"),
		NavLabelBlog:        c.FormValue("nav_label_blog"),
		NavLabelWhitepapers: c.FormValue("nav_label_whitepapers"),
		NavLabelPartners:    c.FormValue("nav_label_partners"),
		NavLabelContact:     c.FormValue("nav_label_contact"),
	})
	if err != nil {
		h.logger.Error("failed to update header settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "header", 0, "", "Updated Header Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/header?saved=1")
}
