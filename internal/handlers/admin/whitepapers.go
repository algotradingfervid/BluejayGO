package admin

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

const whitepapersPerPage = 15

// List displays all whitepapers with filtering and pagination
func (h *WhitepapersHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	search := c.QueryParam("search")
	status := c.QueryParam("status")
	topicStr := c.QueryParam("topic")
	pageStr := c.QueryParam("page")

	var topicID int64
	if topicStr != "" {
		topicID, _ = strconv.ParseInt(topicStr, 10, 64)
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * whitepapersPerPage)

	filterParams := sqlc.ListWhitepapersAdminFilteredParams{
		FilterSearch: search,
		FilterTopic:  topicID,
		FilterStatus: status,
		PageLimit:    whitepapersPerPage,
		PageOffset:   offset,
	}

	whitepapers, err := h.queries.ListWhitepapersAdminFiltered(ctx, filterParams)
	if err != nil {
		h.logger.Error("Failed to list whitepapers", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load whitepapers")
	}

	total, err := h.queries.CountWhitepapersAdminFiltered(ctx, sqlc.CountWhitepapersAdminFilteredParams{
		FilterSearch: search,
		FilterTopic:  topicID,
		FilterStatus: status,
	})
	if err != nil {
		h.logger.Error("Failed to count whitepapers", "error", err)
		total = 0
	}

	topics, err := h.queries.ListWhitepaperTopics(ctx)
	if err != nil {
		h.logger.Error("Failed to list whitepaper topics", "error", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(whitepapersPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	showFrom := offset + 1
	showTo := offset + int64(len(whitepapers))
	if total == 0 {
		showFrom = 0
	}

	hasFilters := search != "" || status != "" || topicStr != ""

	return c.Render(http.StatusOK, "admin/pages/whitepapers_list.html", map[string]interface{}{
		"Title":       "Whitepapers",
		"Whitepapers": whitepapers,
		"Topics":      topics,
		"Search":      search,
		"Status":      status,
		"TopicID":     topicID,
		"HasFilters":  hasFilters,
		"Page":        page,
		"TotalPages":  totalPages,
		"Pages":       pages,
		"Total":       total,
		"ShowFrom":    showFrom,
		"ShowTo":      showTo,
	})
}

// New displays the form for creating a new whitepaper
func (h *WhitepapersHandler) New(c echo.Context) error {
	topics, err := h.queries.ListWhitepaperTopics(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list whitepaper topics", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load topics")
	}

	return c.Render(http.StatusOK, "admin/pages/whitepapers_form.html", map[string]interface{}{
		"Title":      "New Whitepaper",
		"FormAction": "/admin/whitepapers",
		"Item":       nil,
		"Topics":     topics,
		"IsNew":      true,
	})
}

// Create handles whitepaper creation
func (h *WhitepapersHandler) Create(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(50 << 20); err != nil {
		h.logger.Error("Failed to parse multipart form", "error", err)
		return c.String(http.StatusBadRequest, "Failed to parse form")
	}

	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	description := c.FormValue("description")
	publishedDate := c.FormValue("published_date")
	coverColorFrom := c.FormValue("cover_color_from")
	coverColorTo := c.FormValue("cover_color_to")
	metaDescription := c.FormValue("meta_description")

	topicID := int64(0)
	if v := c.FormValue("topic_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			topicID = parsed
		}
	}

	isPublished := int64(0)
	if c.FormValue("is_published") == "on" {
		isPublished = 1
	}

	pageCount := sql.NullInt64{}
	if v := c.FormValue("page_count"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			pageCount = sql.NullInt64{Int64: parsed, Valid: true}
		}
	}

	// Handle PDF upload
	pdfFilePath := ""
	var fileSizeBytes int64
	file, err := c.FormFile("pdf_file")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			h.logger.Error("Failed to open uploaded file", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to process upload")
		}
		defer src.Close()

		uploadDir := "public/uploads/whitepapers"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			h.logger.Error("Failed to create upload directory", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to create upload directory")
		}

		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), makeSlug(title), ext)
		dstPath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(dstPath)
		if err != nil {
			h.logger.Error("Failed to create destination file", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to save file")
		}
		defer dst.Close()

		written, err := io.Copy(dst, src)
		if err != nil {
			h.logger.Error("Failed to copy file", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to save file")
		}

		pdfFilePath = "/uploads/whitepapers/" + filename
		fileSizeBytes = written
	}

	params := sqlc.CreateWhitepaperParams{
		Title:           title,
		Slug:            slug,
		Description:     description,
		TopicID:         topicID,
		PdfFilePath:     pdfFilePath,
		FileSizeBytes:   fileSizeBytes,
		PageCount:       pageCount,
		PublishedDate:   publishedDate,
		IsPublished:     isPublished,
		CoverColorFrom:  coverColorFrom,
		CoverColorTo:    coverColorTo,
		MetaDescription: sql.NullString{String: metaDescription, Valid: metaDescription != ""},
	}

	whitepaper, err := h.queries.CreateWhitepaper(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create whitepaper", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to create whitepaper")
	}

	// Handle learning points
	learningPoints := c.Request().Form["learning_points[]"]
	for i, point := range learningPoints {
		if point == "" {
			continue
		}
		_, err := h.queries.CreateWhitepaperLearningPoint(c.Request().Context(), sqlc.CreateWhitepaperLearningPointParams{
			WhitepaperID: whitepaper.ID,
			PointText:    point,
			DisplayOrder: int64(i + 1),
		})
		if err != nil {
			h.logger.Error("Failed to create learning point", "error", err)
		}
	}

	h.cache.DeleteByPrefix("page:whitepapers")
	logActivity(c, "created", "whitepaper", 0, c.FormValue("title"), "Created Whitepaper '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/whitepapers")
}

// Edit displays the form for editing a whitepaper
func (h *WhitepapersHandler) Edit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid whitepaper ID")
	}

	whitepaper, err := h.queries.GetWhitepaperByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get whitepaper", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load whitepaper")
	}

	topics, err := h.queries.ListWhitepaperTopics(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list whitepaper topics", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load topics")
	}

	learningPoints, err := h.queries.GetWhitepaperLearningPoints(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get learning points", "error", err)
	}

	return c.Render(http.StatusOK, "admin/pages/whitepapers_form.html", map[string]interface{}{
		"Title":          "Edit Whitepaper",
		"FormAction":     "/admin/whitepapers/" + c.Param("id"),
		"Item":           whitepaper,
		"Topics":         topics,
		"LearningPoints": learningPoints,
		"IsNew":          false,
	})
}

// Update handles whitepaper updates
func (h *WhitepapersHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid whitepaper ID")
	}

	if err := c.Request().ParseMultipartForm(50 << 20); err != nil {
		h.logger.Error("Failed to parse multipart form", "error", err)
		return c.String(http.StatusBadRequest, "Failed to parse form")
	}

	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	description := c.FormValue("description")
	publishedDate := c.FormValue("published_date")
	coverColorFrom := c.FormValue("cover_color_from")
	coverColorTo := c.FormValue("cover_color_to")
	metaDescription := c.FormValue("meta_description")

	topicID := int64(0)
	if v := c.FormValue("topic_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			topicID = parsed
		}
	}

	isPublished := int64(0)
	if c.FormValue("is_published") == "on" {
		isPublished = 1
	}

	pageCount := sql.NullInt64{}
	if v := c.FormValue("page_count"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			pageCount = sql.NullInt64{Int64: parsed, Valid: true}
		}
	}

	// Get existing whitepaper for current PDF path
	existing, err := h.queries.GetWhitepaperByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get existing whitepaper", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load whitepaper")
	}

	pdfFilePath := existing.PdfFilePath
	fileSizeBytes := existing.FileSizeBytes

	// Handle optional PDF replacement
	file, err := c.FormFile("pdf_file")
	if err == nil && file != nil {
		src, err := file.Open()
		if err != nil {
			h.logger.Error("Failed to open uploaded file", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to process upload")
		}
		defer src.Close()

		uploadDir := "public/uploads/whitepapers"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			h.logger.Error("Failed to create upload directory", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to create upload directory")
		}

		// Remove old file if it exists
		if existing.PdfFilePath != "" {
			oldPath := filepath.Join("public", existing.PdfFilePath)
			os.Remove(oldPath)
		}

		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), makeSlug(title), ext)
		dstPath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(dstPath)
		if err != nil {
			h.logger.Error("Failed to create destination file", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to save file")
		}
		defer dst.Close()

		written, err := io.Copy(dst, src)
		if err != nil {
			h.logger.Error("Failed to copy file", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to save file")
		}

		pdfFilePath = "/uploads/whitepapers/" + filename
		fileSizeBytes = written
	}

	params := sqlc.UpdateWhitepaperParams{
		Title:           title,
		Slug:            slug,
		Description:     description,
		TopicID:         topicID,
		PdfFilePath:     pdfFilePath,
		FileSizeBytes:   fileSizeBytes,
		PageCount:       pageCount,
		PublishedDate:   publishedDate,
		IsPublished:     isPublished,
		CoverColorFrom:  coverColorFrom,
		CoverColorTo:    coverColorTo,
		MetaDescription: sql.NullString{String: metaDescription, Valid: metaDescription != ""},
		ID:              id,
	}

	err = h.queries.UpdateWhitepaper(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to update whitepaper", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update whitepaper")
	}

	// Replace learning points
	err = h.queries.DeleteWhitepaperLearningPoints(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete old learning points", "error", err)
	}

	learningPoints := c.Request().Form["learning_points[]"]
	for i, point := range learningPoints {
		if point == "" {
			continue
		}
		_, err := h.queries.CreateWhitepaperLearningPoint(c.Request().Context(), sqlc.CreateWhitepaperLearningPointParams{
			WhitepaperID: id,
			PointText:    point,
			DisplayOrder: int64(i + 1),
		})
		if err != nil {
			h.logger.Error("Failed to create learning point", "error", err)
		}
	}

	h.cache.DeleteByPrefix("page:whitepapers")
	logActivity(c, "updated", "whitepaper", id, c.FormValue("title"), "Updated Whitepaper '%s'", c.FormValue("title"))
	return c.Redirect(http.StatusSeeOther, "/admin/whitepapers")
}

// Delete handles whitepaper deletion
func (h *WhitepapersHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid whitepaper ID")
	}

	// Get whitepaper to find PDF path for cleanup
	whitepaper, err := h.queries.GetWhitepaperByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get whitepaper for deletion", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load whitepaper")
	}

	err = h.queries.DeleteWhitepaper(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete whitepaper", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete whitepaper")
	}

	// Remove PDF file
	if whitepaper.PdfFilePath != "" {
		oldPath := filepath.Join("public", whitepaper.PdfFilePath)
		os.Remove(oldPath)
	}

	h.cache.DeleteByPrefix("page:whitepapers")
	logActivity(c, "deleted", "whitepaper", id, "", "Deleted Whitepaper #%d", id)
	return c.NoContent(http.StatusOK)
}

// Downloads lists all whitepaper download leads with filtering and pagination
func (h *WhitepapersHandler) Downloads(c echo.Context) error {
	ctx := c.Request().Context()

	whitepaperStr := c.QueryParam("whitepaper")
	dateFrom := c.QueryParam("date_from")
	dateTo := c.QueryParam("date_to")
	pageStr := c.QueryParam("page")

	var whitepaperID int64
	if whitepaperStr != "" {
		whitepaperID, _ = strconv.ParseInt(whitepaperStr, 10, 64)
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPage := int64(25)
	offset := int64((page - 1)) * perPage

	downloads, err := h.queries.ListWhitepaperDownloadsFiltered(ctx, sqlc.ListWhitepaperDownloadsFilteredParams{
		FilterWhitepaper: whitepaperID,
		FilterDateFrom:   dateFrom,
		FilterDateTo:     dateTo,
		PageLimit:        perPage,
		PageOffset:       offset,
	})
	if err != nil {
		h.logger.Error("Failed to list whitepaper downloads", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load downloads")
	}

	totalCount, err := h.queries.CountWhitepaperDownloadsFiltered(ctx, sqlc.CountWhitepaperDownloadsFilteredParams{
		FilterWhitepaper: whitepaperID,
		FilterDateFrom:   dateFrom,
		FilterDateTo:     dateTo,
	})
	if err != nil {
		h.logger.Error("Failed to count whitepaper downloads", "error", err)
		totalCount = 0
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	// Get whitepapers list for filter dropdown
	whitepapers, err := h.queries.ListAllWhitepapers(ctx)
	if err != nil {
		h.logger.Error("Failed to list whitepapers for filter", "error", err)
	}

	hasFilters := whitepaperStr != "" || dateFrom != "" || dateTo != ""

	return c.Render(http.StatusOK, "admin/pages/whitepapers_downloads.html", map[string]interface{}{
		"Title":        "Whitepaper Downloads",
		"Downloads":    downloads,
		"Whitepapers":  whitepapers,
		"WhitepaperID": whitepaperID,
		"DateFrom":     dateFrom,
		"DateTo":       dateTo,
		"HasFilters":   hasFilters,
		"Page":         page,
		"TotalPages":   totalPages,
		"Pages":        pages,
		"TotalCount":   totalCount,
	})
}
