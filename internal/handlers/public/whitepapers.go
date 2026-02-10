package public

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type WhitepapersHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewWhitepapersHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *WhitepapersHandler {
	return &WhitepapersHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

func (h *WhitepapersHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
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

// GET /whitepapers
func (h *WhitepapersHandler) WhitepapersList(c echo.Context) error {
	ctx := c.Request().Context()

	// Get optional topic filter
	topicParam := c.QueryParam("topic")
	var selectedTopicID int64
	var whitepapers interface{}
	var totalCount int64
	var err error

	if topicParam != "" {
		selectedTopicID, err = strconv.ParseInt(topicParam, 10, 64)
		if err != nil {
			h.logger.Error("invalid topic parameter", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid topic parameter")
		}
	}

	// Determine cache key based on filter
	cacheKey := "page:whitepapers"
	if selectedTopicID > 0 {
		cacheKey = fmt.Sprintf("page:whitepapers:topic:%d", selectedTopicID)
	}

	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	// Fetch whitepapers based on filter
	if selectedTopicID > 0 {
		whitepapers, err = h.queries.ListPublishedWhitepapersByTopic(ctx, selectedTopicID)
		if err != nil {
			h.logger.Error("failed to list whitepapers by topic", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		totalCount, err = h.queries.CountPublishedWhitepapersByTopic(ctx, selectedTopicID)
		if err != nil {
			h.logger.Error("failed to count whitepapers by topic", "error", err)
			totalCount = 0
		}
	} else {
		whitepapers, err = h.queries.ListPublishedWhitepapers(ctx)
		if err != nil {
			h.logger.Error("failed to list whitepapers", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		totalCount, err = h.queries.CountPublishedWhitepapers(ctx)
		if err != nil {
			h.logger.Error("failed to count whitepapers", "error", err)
			totalCount = 0
		}
	}

	// Get topics for filter dropdown
	topics, err := h.queries.ListWhitepaperTopics(ctx)
	if err != nil {
		h.logger.Error("failed to list whitepaper topics", "error", err)
		topics = []sqlc.WhitepaperTopic{}
	}

	data := map[string]interface{}{
		"Title":           "Whitepapers",
		"Whitepapers":     whitepapers,
		"Topics":          topics,
		"SelectedTopicID": selectedTopicID,
		"TotalCount":      totalCount,
		"CurrentPage":     "whitepapers",
	}

	return h.renderAndCache(c, cacheKey, 600, http.StatusOK, "public/pages/whitepapers.html", data)
}

// GET /whitepapers/:slug
func (h *WhitepapersHandler) WhitepaperDetail(c echo.Context) error {
	slug := c.Param("slug")
	preview := isPreviewRequest(c)

	if !preview {
		cacheKey := fmt.Sprintf("page:whitepapers:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	ctx := c.Request().Context()

	var wpID, wpTopicID int64
	var wpTitle, wpMetaTitle string
	var wpSlug, wpOgImage string
	var wpMetaDesc sql.NullString
	var wpObj interface{}

	if preview {
		wp, err := h.queries.GetWhitepaperBySlugIncludeDrafts(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Whitepaper not found")
		}
		if err != nil {
			h.logger.Error("failed to load whitepaper", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		wpID, wpTopicID = wp.ID, wp.TopicID
		wpTitle, wpSlug, wpMetaTitle = wp.Title, wp.Slug, wp.MetaTitle
		wpOgImage, wpMetaDesc = wp.OgImage, wp.MetaDescription
		wpObj = wp
	} else {
		wp, err := h.queries.GetWhitepaperBySlug(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Whitepaper not found")
		}
		if err != nil {
			h.logger.Error("failed to load whitepaper", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		wpID, wpTopicID = wp.ID, wp.TopicID
		wpTitle, wpSlug, wpMetaTitle = wp.Title, wp.Slug, wp.MetaTitle
		wpOgImage, wpMetaDesc = wp.OgImage, wp.MetaDescription
		wpObj = wp
	}

	learningPoints, err := h.queries.GetWhitepaperLearningPoints(ctx, wpID)
	if err != nil {
		h.logger.Error("failed to load whitepaper learning points", "error", err)
		learningPoints = []sqlc.GetWhitepaperLearningPointsRow{}
	}

	relatedPapers, err := h.queries.GetRelatedWhitepapers(ctx, sqlc.GetRelatedWhitepapersParams{
		ID:      wpID,
		TopicID: wpTopicID,
	})
	if err != nil {
		h.logger.Error("failed to load related whitepapers", "error", err)
		relatedPapers = []sqlc.GetRelatedWhitepapersRow{}
	}

	metaDesc := ""
	if wpMetaDesc.Valid {
		metaDesc = wpMetaDesc.String
	}

	data := map[string]interface{}{
		"Title":           wpTitle,
		"MetaTitle":       wpMetaTitle,
		"MetaDescription": metaDesc,
		"MetaDesc":        metaDesc,
		"OGImage":         wpOgImage,
		"CanonicalURL":    fmt.Sprintf("/whitepapers/%s", wpSlug),
		"Whitepaper":      wpObj,
		"LearningPoints":  learningPoints,
		"RelatedPapers":   relatedPapers,
		"CurrentPage":     "whitepapers",
	}

	if preview {
		data["IsPreview"] = true
		data["EditURL"] = fmt.Sprintf("/admin/whitepapers/%d/edit", wpID)
		return h.renderAndCache(c, "preview:whitepaper:"+slug, 0, http.StatusOK, "public/pages/whitepaper_detail.html", data)
	}

	return h.renderAndCache(c, fmt.Sprintf("page:whitepapers:%s", slug), 900, http.StatusOK, "public/pages/whitepaper_detail.html", data)
}

// POST /whitepapers/:slug/download
func (h *WhitepapersHandler) WhitepaperDownload(c echo.Context) error {
	slug := c.Param("slug")
	ctx := c.Request().Context()

	whitepaper, err := h.queries.GetWhitepaperBySlug(ctx, slug)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "Whitepaper not found")
	}
	if err != nil {
		h.logger.Error("failed to load whitepaper", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Parse form values
	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	company := strings.TrimSpace(c.FormValue("company"))
	designation := strings.TrimSpace(c.FormValue("designation"))
	marketingConsent := c.FormValue("marketing_consent")

	// Validate required fields
	if name == "" || email == "" || company == "" {
		return c.HTML(http.StatusBadRequest, `<div class="alert alert-error">Name, email, and company are required.</div>`)
	}

	// Parse marketing consent
	var consent int64
	if marketingConsent == "on" || marketingConsent == "1" || marketingConsent == "true" {
		consent = 1
	}

	// Create download record
	_, err = h.queries.CreateWhitepaperDownload(ctx, sqlc.CreateWhitepaperDownloadParams{
		WhitepaperID: whitepaper.ID,
		Name:         name,
		Email:        email,
		Company:      company,
		Designation: sql.NullString{
			String: designation,
			Valid:  designation != "",
		},
		MarketingConsent: consent,
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
		h.logger.Error("failed to create whitepaper download", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Increment download count asynchronously
	go func() {
		if err := h.queries.IncrementWhitepaperDownloadCount(ctx, whitepaper.ID); err != nil {
			h.logger.Error("failed to increment whitepaper download count", "error", err)
		}
	}()

	// Invalidate cache for this whitepaper
	h.cache.Delete(fmt.Sprintf("page:whitepapers:%s", slug))
	h.cache.Delete("page:whitepapers")

	// Render success fragment
	data := map[string]interface{}{
		"Whitepaper":    whitepaper,
		"Email":         email,
		"WhitepaperURL": "/" + whitepaper.PdfFilePath,
	}
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}

	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, "public/pages/whitepaper_success.html", data, c); err != nil {
		h.logger.Error("template render failed", "template", "public/pages/whitepaper_success.html", "error", err)
		return err
	}

	return c.HTML(http.StatusOK, buf.String())
}
