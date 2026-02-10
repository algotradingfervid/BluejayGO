package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type HomepageHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewHomepageHandler(queries *sqlc.Queries, logger *slog.Logger) *HomepageHandler {
	return &HomepageHandler{queries: queries, logger: logger}
}

// ==================== HEROES ====================

func (h *HomepageHandler) HeroesList(c echo.Context) error {
	items, err := h.queries.ListAllHeroes(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list heroes", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/homepage_heroes_list.html", map[string]interface{}{
		"Title": "Homepage Heroes",
		"Items": items,
	})
}

func (h *HomepageHandler) HeroNew(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/homepage_hero_form.html", map[string]interface{}{
		"Title": "New Hero",
		"Item":  sqlc.HomepageHero{},
	})
}

func (h *HomepageHandler) HeroCreate(c echo.Context) error {
	displayOrder, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	_, err := h.queries.CreateHero(c.Request().Context(), sqlc.CreateHeroParams{
		Headline:         c.FormValue("headline"),
		Subheadline:      c.FormValue("subheadline"),
		BadgeText:        sql.NullString{String: c.FormValue("badge_text"), Valid: c.FormValue("badge_text") != ""},
		PrimaryCtaText:   c.FormValue("primary_cta_text"),
		PrimaryCtaUrl:    c.FormValue("primary_cta_url"),
		SecondaryCtaText: sql.NullString{String: c.FormValue("secondary_cta_text"), Valid: c.FormValue("secondary_cta_text") != ""},
		SecondaryCtaUrl:  sql.NullString{String: c.FormValue("secondary_cta_url"), Valid: c.FormValue("secondary_cta_url") != ""},
		BackgroundImage:  sql.NullString{String: c.FormValue("background_image"), Valid: c.FormValue("background_image") != ""},
		IsActive:         isActive,
		DisplayOrder:     displayOrder,
	})
	if err != nil {
		h.logger.Error("failed to create hero", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "hero", 0, c.FormValue("headline"), "Created Hero '%s'", c.FormValue("headline"))
	return c.Redirect(http.StatusSeeOther, "/admin/homepage/heroes")
}

func (h *HomepageHandler) HeroEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetHero(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get hero", "error", err)
		return echo.NewHTTPError(http.StatusNotFound)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/homepage_hero_form.html", map[string]interface{}{
		"Title": "Edit Hero",
		"Item":  item,
		"Saved": saved,
	})
}

func (h *HomepageHandler) HeroUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	displayOrder, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	err := h.queries.UpdateHero(c.Request().Context(), sqlc.UpdateHeroParams{
		ID:               id,
		Headline:         c.FormValue("headline"),
		Subheadline:      c.FormValue("subheadline"),
		BadgeText:        sql.NullString{String: c.FormValue("badge_text"), Valid: c.FormValue("badge_text") != ""},
		PrimaryCtaText:   c.FormValue("primary_cta_text"),
		PrimaryCtaUrl:    c.FormValue("primary_cta_url"),
		SecondaryCtaText: sql.NullString{String: c.FormValue("secondary_cta_text"), Valid: c.FormValue("secondary_cta_text") != ""},
		SecondaryCtaUrl:  sql.NullString{String: c.FormValue("secondary_cta_url"), Valid: c.FormValue("secondary_cta_url") != ""},
		BackgroundImage:  sql.NullString{String: c.FormValue("background_image"), Valid: c.FormValue("background_image") != ""},
		IsActive:         isActive,
		DisplayOrder:     displayOrder,
	})
	if err != nil {
		h.logger.Error("failed to update hero", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "hero", id, c.FormValue("headline"), "Updated Hero '%s'", c.FormValue("headline"))
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/homepage/heroes/%d/edit?saved=1", id))
}

func (h *HomepageHandler) HeroDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := h.queries.DeleteHero(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to delete hero", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "hero", id, "", "Deleted Hero #%d", id)
	return c.NoContent(http.StatusOK)
}

// ==================== STATS ====================

func (h *HomepageHandler) StatsList(c echo.Context) error {
	items, err := h.queries.ListAllStats(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list stats", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/homepage_stats_list.html", map[string]interface{}{
		"Title": "Homepage Stats",
		"Items": items,
	})
}

func (h *HomepageHandler) StatNew(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/homepage_stat_form.html", map[string]interface{}{
		"Title": "New Stat",
		"Item":  sqlc.HomepageStat{},
	})
}

func (h *HomepageHandler) StatCreate(c echo.Context) error {
	displayOrder, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	_, err := h.queries.CreateStat(c.Request().Context(), sqlc.CreateStatParams{
		StatValue:    c.FormValue("stat_value"),
		StatLabel:    c.FormValue("stat_label"),
		DisplayOrder: displayOrder,
		IsActive:     isActive,
	})
	if err != nil {
		h.logger.Error("failed to create stat", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "stat", 0, c.FormValue("stat_label"), "Created Stat '%s'", c.FormValue("stat_label"))
	return c.Redirect(http.StatusSeeOther, "/admin/homepage/stats")
}

func (h *HomepageHandler) StatEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetStat(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get stat", "error", err)
		return echo.NewHTTPError(http.StatusNotFound)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/homepage_stat_form.html", map[string]interface{}{
		"Title": "Edit Stat",
		"Item":  item,
		"Saved": saved,
	})
}

func (h *HomepageHandler) StatUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	displayOrder, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	err := h.queries.UpdateStat(c.Request().Context(), sqlc.UpdateStatParams{
		ID:           id,
		StatValue:    c.FormValue("stat_value"),
		StatLabel:    c.FormValue("stat_label"),
		DisplayOrder: displayOrder,
		IsActive:     isActive,
	})
	if err != nil {
		h.logger.Error("failed to update stat", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "stat", id, c.FormValue("stat_label"), "Updated Stat '%s'", c.FormValue("stat_label"))
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/homepage/stats/%d/edit?saved=1", id))
}

func (h *HomepageHandler) StatDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := h.queries.DeleteStat(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to delete stat", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "stat", id, "", "Deleted Stat #%d", id)
	return c.NoContent(http.StatusOK)
}

// ==================== TESTIMONIALS ====================

func (h *HomepageHandler) TestimonialsList(c echo.Context) error {
	items, err := h.queries.ListAllTestimonialsHomepage(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list testimonials", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/homepage_testimonials_list.html", map[string]interface{}{
		"Title": "Homepage Testimonials",
		"Items": items,
	})
}

func (h *HomepageHandler) TestimonialNew(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/homepage_testimonial_form.html", map[string]interface{}{
		"Title": "New Testimonial",
		"Item":  sqlc.HomepageTestimonial{},
	})
}

func (h *HomepageHandler) TestimonialCreate(c echo.Context) error {
	displayOrder, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	rating, _ := strconv.ParseInt(c.FormValue("rating"), 10, 64)
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	_, err := h.queries.CreateTestimonialHomepage(c.Request().Context(), sqlc.CreateTestimonialHomepageParams{
		Quote:         c.FormValue("quote"),
		AuthorName:    c.FormValue("author_name"),
		AuthorTitle:   sql.NullString{String: c.FormValue("author_title"), Valid: c.FormValue("author_title") != ""},
		AuthorCompany: sql.NullString{String: c.FormValue("author_company"), Valid: c.FormValue("author_company") != ""},
		AuthorImage:   sql.NullString{String: c.FormValue("author_image"), Valid: c.FormValue("author_image") != ""},
		Rating:        rating,
		DisplayOrder:  displayOrder,
		IsActive:      isActive,
	})
	if err != nil {
		h.logger.Error("failed to create testimonial", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "testimonial", 0, c.FormValue("author_name"), "Created Testimonial by '%s'", c.FormValue("author_name"))
	return c.Redirect(http.StatusSeeOther, "/admin/homepage/testimonials")
}

func (h *HomepageHandler) TestimonialEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetTestimonialHomepage(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get testimonial", "error", err)
		return echo.NewHTTPError(http.StatusNotFound)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/homepage_testimonial_form.html", map[string]interface{}{
		"Title": "Edit Testimonial",
		"Item":  item,
		"Saved": saved,
	})
}

func (h *HomepageHandler) TestimonialUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	displayOrder, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	rating, _ := strconv.ParseInt(c.FormValue("rating"), 10, 64)
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	err := h.queries.UpdateTestimonialHomepage(c.Request().Context(), sqlc.UpdateTestimonialHomepageParams{
		ID:            id,
		Quote:         c.FormValue("quote"),
		AuthorName:    c.FormValue("author_name"),
		AuthorTitle:   sql.NullString{String: c.FormValue("author_title"), Valid: c.FormValue("author_title") != ""},
		AuthorCompany: sql.NullString{String: c.FormValue("author_company"), Valid: c.FormValue("author_company") != ""},
		AuthorImage:   sql.NullString{String: c.FormValue("author_image"), Valid: c.FormValue("author_image") != ""},
		Rating:        rating,
		DisplayOrder:  displayOrder,
		IsActive:      isActive,
	})
	if err != nil {
		h.logger.Error("failed to update testimonial", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "testimonial", id, c.FormValue("author_name"), "Updated Testimonial by '%s'", c.FormValue("author_name"))
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/homepage/testimonials/%d/edit?saved=1", id))
}

func (h *HomepageHandler) TestimonialDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := h.queries.DeleteTestimonialHomepage(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to delete testimonial", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "testimonial", id, "", "Deleted Testimonial #%d", id)
	return c.NoContent(http.StatusOK)
}

// ==================== CTA ====================

func (h *HomepageHandler) CTAList(c echo.Context) error {
	items, err := h.queries.ListAllCTAs(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list CTAs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/homepage_cta_list.html", map[string]interface{}{
		"Title": "Homepage CTAs",
		"Items": items,
	})
}

func (h *HomepageHandler) CTANew(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/homepage_cta_form.html", map[string]interface{}{
		"Title": "New CTA",
		"Item":  sqlc.HomepageCtum{},
	})
}

func (h *HomepageHandler) CTACreate(c echo.Context) error {
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	_, err := h.queries.CreateCTA(c.Request().Context(), sqlc.CreateCTAParams{
		Headline:         c.FormValue("headline"),
		Description:      sql.NullString{String: c.FormValue("description"), Valid: c.FormValue("description") != ""},
		PrimaryCtaText:   c.FormValue("primary_cta_text"),
		PrimaryCtaUrl:    c.FormValue("primary_cta_url"),
		SecondaryCtaText: sql.NullString{String: c.FormValue("secondary_cta_text"), Valid: c.FormValue("secondary_cta_text") != ""},
		SecondaryCtaUrl:  sql.NullString{String: c.FormValue("secondary_cta_url"), Valid: c.FormValue("secondary_cta_url") != ""},
		BackgroundStyle:  sql.NullString{String: c.FormValue("background_style"), Valid: c.FormValue("background_style") != ""},
		IsActive:         isActive,
	})
	if err != nil {
		h.logger.Error("failed to create CTA", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "created", "cta", 0, c.FormValue("headline"), "Created CTA '%s'", c.FormValue("headline"))
	return c.Redirect(http.StatusSeeOther, "/admin/homepage/cta")
}

func (h *HomepageHandler) CTAEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	item, err := h.queries.GetCTA(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to get CTA", "error", err)
		return echo.NewHTTPError(http.StatusNotFound)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/homepage_cta_form.html", map[string]interface{}{
		"Title": "Edit CTA",
		"Item":  item,
		"Saved": saved,
	})
}

func (h *HomepageHandler) CTAUpdate(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var isActive int64
	if c.FormValue("is_active") == "on" {
		isActive = 1
	}
	err := h.queries.UpdateCTA(c.Request().Context(), sqlc.UpdateCTAParams{
		ID:               id,
		Headline:         c.FormValue("headline"),
		Description:      sql.NullString{String: c.FormValue("description"), Valid: c.FormValue("description") != ""},
		PrimaryCtaText:   c.FormValue("primary_cta_text"),
		PrimaryCtaUrl:    c.FormValue("primary_cta_url"),
		SecondaryCtaText: sql.NullString{String: c.FormValue("secondary_cta_text"), Valid: c.FormValue("secondary_cta_text") != ""},
		SecondaryCtaUrl:  sql.NullString{String: c.FormValue("secondary_cta_url"), Valid: c.FormValue("secondary_cta_url") != ""},
		BackgroundStyle:  sql.NullString{String: c.FormValue("background_style"), Valid: c.FormValue("background_style") != ""},
		IsActive:         isActive,
	})
	if err != nil {
		h.logger.Error("failed to update CTA", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "cta", id, c.FormValue("headline"), "Updated CTA '%s'", c.FormValue("headline"))
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/homepage/cta/%d/edit?saved=1", id))
}

func (h *HomepageHandler) CTADelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := h.queries.DeleteCTA(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to delete CTA", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "deleted", "cta", id, "", "Deleted CTA #%d", id)
	return c.NoContent(http.StatusOK)
}

// ==================== SETTINGS ====================

func (h *HomepageHandler) Settings(c echo.Context) error {
	settings, err := h.queries.GetSettings(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	saved := c.QueryParam("saved") == "1"
	return c.Render(http.StatusOK, "admin/pages/homepage_settings.html", map[string]interface{}{
		"Title":    "Homepage Settings",
		"Settings": settings,
		"Saved":    saved,
	})
}

func (h *HomepageHandler) UpdateSettings(c echo.Context) error {
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

	err := h.queries.UpdateHomepageSettings(c.Request().Context(), sqlc.UpdateHomepageSettingsParams{
		HomepageShowHeroes:       boolToInt("homepage_show_heroes"),
		HomepageShowStats:        boolToInt("homepage_show_stats"),
		HomepageShowTestimonials: boolToInt("homepage_show_testimonials"),
		HomepageShowCta:          boolToInt("homepage_show_cta"),
		HomepageMaxHeroes:        parseIntField("homepage_max_heroes", 5),
		HomepageMaxStats:         parseIntField("homepage_max_stats", 6),
		HomepageMaxTestimonials:  parseIntField("homepage_max_testimonials", 3),
		HomepageHeroAutoplay:     boolToInt("homepage_hero_autoplay"),
		HomepageHeroInterval:     parseIntField("homepage_hero_interval", 5),
	})
	if err != nil {
		h.logger.Error("failed to update homepage settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "homepage_settings", 0, "", "Updated Homepage Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/homepage/settings?saved=1")
}
