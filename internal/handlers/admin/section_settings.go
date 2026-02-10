package admin

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type SectionSettingsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewSectionSettingsHandler(queries *sqlc.Queries, logger *slog.Logger) *SectionSettingsHandler {
	return &SectionSettingsHandler{queries: queries, logger: logger}
}

// ==================== ABOUT SETTINGS ====================

func (h *SectionSettingsHandler) AboutSettings(c echo.Context) error {
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/about_settings.html", map[string]interface{}{
		"Title":    "About Settings",
		"Settings": settings,
		"Saved":    saved,
	})
}

func (h *SectionSettingsHandler) UpdateAboutSettings(c echo.Context) error {
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1
		}
		return 0
	}

	err := h.queries.UpdateAboutSettings(c.Request().Context(), sqlc.UpdateAboutSettingsParams{
		AboutShowMission:        boolToInt("about_show_mission"),
		AboutShowMilestones:     boolToInt("about_show_milestones"),
		AboutShowCertifications: boolToInt("about_show_certifications"),
		AboutShowTeam:           boolToInt("about_show_team"),
	})
	if err != nil {
		h.logger.Error("failed to update about settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "about_settings", 0, "", "Updated About Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/about/settings?saved=1")
}

// ==================== PRODUCTS SETTINGS ====================

func (h *SectionSettingsHandler) ProductsSettings(c echo.Context) error {
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/products_settings.html", map[string]interface{}{
		"Title":    "Products Settings",
		"Settings": settings,
		"Saved":    saved,
	})
}

func (h *SectionSettingsHandler) UpdateProductsSettings(c echo.Context) error {
	parseIntField := func(field string, defaultVal int64) int64 {
		v := strings.TrimSpace(c.FormValue(field))
		if v == "" {
			return defaultVal
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultVal
		}
		return n
	}
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1
		}
		return 0
	}

	err := h.queries.UpdateProductsSettings(c.Request().Context(), sqlc.UpdateProductsSettingsParams{
		ProductsPerPage:        parseIntField("products_per_page", 12),
		ProductsShowCategories: boolToInt("products_show_categories"),
		ProductsShowSearch:     boolToInt("products_show_search"),
		ProductsDefaultSort:    c.FormValue("products_default_sort"),
	})
	if err != nil {
		h.logger.Error("failed to update products settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "products_settings", 0, "", "Updated Products Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/products/settings?saved=1")
}

// ==================== SOLUTIONS SETTINGS ====================

func (h *SectionSettingsHandler) SolutionsSettings(c echo.Context) error {
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/solutions_settings.html", map[string]interface{}{
		"Title":    "Solutions Settings",
		"Settings": settings,
		"Saved":    saved,
	})
}

func (h *SectionSettingsHandler) UpdateSolutionsSettings(c echo.Context) error {
	parseIntField := func(field string, defaultVal int64) int64 {
		v := strings.TrimSpace(c.FormValue(field))
		if v == "" {
			return defaultVal
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultVal
		}
		return n
	}
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1
		}
		return 0
	}

	err := h.queries.UpdateSolutionsSettings(c.Request().Context(), sqlc.UpdateSolutionsSettingsParams{
		SolutionsPerPage:        parseIntField("solutions_per_page", 12),
		SolutionsShowIndustries: boolToInt("solutions_show_industries"),
		SolutionsShowSearch:     boolToInt("solutions_show_search"),
	})
	if err != nil {
		h.logger.Error("failed to update solutions settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "solutions_settings", 0, "", "Updated Solutions Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/solutions/settings?saved=1")
}

// ==================== BLOG SETTINGS ====================

func (h *SectionSettingsHandler) BlogSettings(c echo.Context) error {
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/blog_settings.html", map[string]interface{}{
		"Title":    "Blog Settings",
		"Settings": settings,
		"Saved":    saved,
	})
}

func (h *SectionSettingsHandler) UpdateBlogSettings(c echo.Context) error {
	parseIntField := func(field string, defaultVal int64) int64 {
		v := strings.TrimSpace(c.FormValue(field))
		if v == "" {
			return defaultVal
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultVal
		}
		return n
	}
	boolToInt := func(field string) int64 {
		if c.FormValue(field) == "on" {
			return 1
		}
		return 0
	}

	err := h.queries.UpdateBlogSettings(c.Request().Context(), sqlc.UpdateBlogSettingsParams{
		BlogPostsPerPage:   parseIntField("blog_posts_per_page", 10),
		BlogShowAuthor:     boolToInt("blog_show_author"),
		BlogShowDate:       boolToInt("blog_show_date"),
		BlogShowCategories: boolToInt("blog_show_categories"),
		BlogShowTags:       boolToInt("blog_show_tags"),
		BlogShowSearch:     boolToInt("blog_show_search"),
	})
	if err != nil {
		h.logger.Error("failed to update blog settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "blog_settings", 0, "", "Updated Blog Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/blog/settings?saved=1")
}
