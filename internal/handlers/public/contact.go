package public

import (
	"bytes"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type ContactHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewContactHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *ContactHandler {
	return &ContactHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

func (h *ContactHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}
	if cats := c.Get("footer_categories"); cats != nil {
		data["FooterCategories"] = cats
	}
	if sols := c.Get("footer_solutions"); sols != nil {
		data["FooterSolutions"] = sols
	}
	if res := c.Get("footer_resources"); res != nil {
		data["FooterResources"] = res
	}
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		h.logger.Error("template render failed", "template", templateName, "error", err)
		return err
	}
	html := buf.String()
	h.cache.Set(cacheKey, html, ttlSeconds)
	return c.HTML(statusCode, html)
}

// GET /contact
func (h *ContactHandler) ShowContactPage(c echo.Context) error {
	cacheKey := "page:contact"
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()

	offices, err := h.queries.GetActiveOfficeLocations(ctx)
	if err != nil {
		h.logger.Error("failed to load office locations", "error", err)
		offices = []sqlc.GetActiveOfficeLocationsRow{}
	}

	data := map[string]interface{}{
		"Title":       "Contact Us",
		"Offices":     offices,
		"CurrentPage": "contact",
	}

	return h.renderAndCache(c, cacheKey, 3600, http.StatusOK, "public/pages/contact.html", data)
}

// POST /contact/submit
func (h *ContactHandler) SubmitContactForm(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse form values
	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	phone := strings.TrimSpace(c.FormValue("phone"))
	company := strings.TrimSpace(c.FormValue("company"))
	message := strings.TrimSpace(c.FormValue("message"))
	inquiryType := strings.TrimSpace(c.FormValue("inquiry_type"))

	// Validate required fields
	if name == "" || email == "" || phone == "" || company == "" || message == "" {
		return c.HTML(http.StatusBadRequest, `<div class="alert alert-error">Name, email, phone, company, and message are required.</div>`)
	}

	// Create contact submission
	_, err := h.queries.CreateContactSubmission(ctx, sqlc.CreateContactSubmissionParams{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Company: company,
		Message: message,
		InquiryType: sql.NullString{
			String: inquiryType,
			Valid:  inquiryType != "",
		},
		IpAddress: sql.NullString{
			String: c.RealIP(),
			Valid:  c.RealIP() != "",
		},
		UserAgent: sql.NullString{
			String: c.Request().UserAgent(),
			Valid:  c.Request().UserAgent() != "",
		},
	})
	if err != nil {
		h.logger.Error("failed to create contact submission", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.HTML(http.StatusOK, `<div class="alert alert-success">Thank you for your message. We will get back to you shortly.</div>`)
}
