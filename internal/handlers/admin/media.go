package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
)

type MediaHandler struct {
	queries   *sqlc.Queries
	logger    *slog.Logger
	uploadDir string
}

func NewMediaHandler(queries *sqlc.Queries, logger *slog.Logger, uploadDir string) *MediaHandler {
	return &MediaHandler{
		queries:   queries,
		logger:    logger,
		uploadDir: uploadDir,
	}
}

const mediaPerPage = 24

func (h *MediaHandler) List(c echo.Context) error {
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "newest"
	}
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * mediaPerPage)

	var files []sqlc.MediaFile
	var total int64
	var err error

	ctx := c.Request().Context()

	if search != "" {
		files, err = h.queries.SearchMediaFiles(ctx, sqlc.SearchMediaFilesParams{
			Search:     sql.NullString{String: search, Valid: true},
			PageLimit:  int64(mediaPerPage),
			PageOffset: offset,
		})
		if err != nil {
			h.logger.Error("failed to search media files", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to search media files")
		}
		total, err = h.queries.CountMediaFilesSearch(ctx, sql.NullString{String: search, Valid: true})
	} else {
		switch sort {
		case "oldest":
			files, err = h.queries.ListMediaFilesOldest(ctx, sqlc.ListMediaFilesOldestParams{
				Limit: int64(mediaPerPage), Offset: offset,
			})
		case "name":
			files, err = h.queries.ListMediaFilesByName(ctx, sqlc.ListMediaFilesByNameParams{
				Limit: int64(mediaPerPage), Offset: offset,
			})
		case "largest":
			files, err = h.queries.ListMediaFilesBySize(ctx, sqlc.ListMediaFilesBySizeParams{
				Limit: int64(mediaPerPage), Offset: offset,
			})
		default:
			files, err = h.queries.ListMediaFiles(ctx, sqlc.ListMediaFilesParams{
				Limit: int64(mediaPerPage), Offset: offset,
			})
		}
		if err != nil {
			h.logger.Error("failed to list media files", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to list media files")
		}
		total, err = h.queries.CountMediaFiles(ctx)
	}
	if err != nil {
		h.logger.Error("failed to count media files", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to count media files")
	}

	totalPages := (total + int64(mediaPerPage) - 1) / int64(mediaPerPage)

	data := map[string]interface{}{
		"Title":       "Media Library",
		"Files":       files,
		"Total":       total,
		"Page":        page,
		"TotalPages":  totalPages,
		"Search":      search,
		"Sort":        sort,
		"HasFilters":  search != "",
		"ActiveNav":   "media",
		"DisplayName": getSessionDisplayName(c),
		"Role":        getSessionRole(c),
	}

	return c.Render(http.StatusOK, "admin/pages/media_library.html", data)
}

func (h *MediaHandler) Upload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
	}

	formFiles := form.File["files"]
	if len(formFiles) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No files provided"})
	}

	allowedTypes := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".svg": true, ".pdf": true, ".webp": true,
	}

	mediaDir := filepath.Join(h.uploadDir, "media")
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
	}

	var uploaded []sqlc.MediaFile
	ctx := c.Request().Context()

	for _, file := range formFiles {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedTypes[ext] {
			continue
		}
		if file.Size > 10*1024*1024 {
			continue
		}

		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), sanitizeMediaFilename(file.Filename))
		webPath := "/uploads/media/" + filename
		dstPath := filepath.Join(mediaDir, filename)

		width, height := h.getImageDimensions(file)

		mimeType := file.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = getMimeTypeFromExt(ext)
		}

		src, err := file.Open()
		if err != nil {
			h.logger.Error("failed to open uploaded file", "error", err)
			continue
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			src.Close()
			h.logger.Error("failed to create destination file", "error", err)
			continue
		}

		if _, err = io.Copy(dst, src); err != nil {
			src.Close()
			dst.Close()
			h.logger.Error("failed to copy file", "error", err)
			continue
		}
		src.Close()
		dst.Close()

		mediaFile, err := h.queries.CreateMediaFile(ctx, sqlc.CreateMediaFileParams{
			Filename:         filename,
			OriginalFilename: file.Filename,
			FilePath:         webPath,
			FileSize:         file.Size,
			MimeType:         mimeType,
			Width:            sql.NullInt64{Int64: int64(width), Valid: width > 0},
			Height:           sql.NullInt64{Int64: int64(height), Valid: height > 0},
			AltText:          sql.NullString{String: "", Valid: true},
		})
		if err != nil {
			h.logger.Error("failed to save media file record", "error", err)
			continue
		}
		uploaded = append(uploaded, mediaFile)
	}

	if len(uploaded) > 0 {
		logActivity(c, "created", "media", 0, "", "Uploaded Media File")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"files":   uploaded,
		"count":   len(uploaded),
		"message": fmt.Sprintf("%d file(s) uploaded successfully", len(uploaded)),
	})
}

func (h *MediaHandler) GetFile(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	file, err := h.queries.GetMediaFile(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}

	return c.JSON(http.StatusOK, file)
}

func (h *MediaHandler) UpdateAltText(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var body struct {
		AltText string `json:"alt_text"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		body.AltText = c.FormValue("alt_text")
	}

	err = h.queries.UpdateMediaFileAltText(c.Request().Context(), sqlc.UpdateMediaFileAltTextParams{
		AltText: sql.NullString{String: body.AltText, Valid: true},
		ID:      id,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update alt text"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Alt text updated"})
}

func (h *MediaHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	file, err := h.queries.GetMediaFile(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}

	// Delete physical file
	fsPath := filepath.Join(h.uploadDir, strings.TrimPrefix(file.FilePath, "/uploads/"))
	os.Remove(fsPath)

	if err := h.queries.DeleteMediaFile(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete file"})
	}

	logActivity(c, "deleted", "media", id, "", "Deleted Media File #%d", id)
	return c.JSON(http.StatusOK, map[string]string{"message": "File deleted"})
}

func (h *MediaHandler) Browse(c echo.Context) error {
	search := c.QueryParam("search")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * mediaPerPage)

	var files []sqlc.MediaFile
	var total int64
	var err error
	ctx := c.Request().Context()

	if search != "" {
		files, err = h.queries.SearchMediaFiles(ctx, sqlc.SearchMediaFilesParams{
			Search:     sql.NullString{String: search, Valid: true},
			PageLimit:  int64(mediaPerPage),
			PageOffset: offset,
		})
		if err == nil {
			total, err = h.queries.CountMediaFilesSearch(ctx, sql.NullString{String: search, Valid: true})
		}
	} else {
		files, err = h.queries.ListMediaFiles(ctx, sqlc.ListMediaFilesParams{
			Limit: int64(mediaPerPage), Offset: offset,
		})
		if err == nil {
			total, err = h.queries.CountMediaFiles(ctx)
		}
	}
	if err != nil {
		h.logger.Error("failed to browse media files", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load media files")
	}

	totalPages := (total + int64(mediaPerPage) - 1) / int64(mediaPerPage)

	data := map[string]interface{}{
		"Files":      files,
		"Total":      total,
		"Page":       page,
		"TotalPages": totalPages,
		"Search":     search,
	}

	return c.Render(http.StatusOK, "admin/partials/media_picker.html", data)
}

func (h *MediaHandler) getImageDimensions(file *multipart.FileHeader) (int, int) {
	src, err := file.Open()
	if err != nil {
		return 0, 0
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".svg" || ext == ".pdf" {
		return 0, 0
	}

	config, _, err := image.DecodeConfig(src)
	if err != nil {
		return 0, 0
	}
	return config.Width, config.Height
}

func sanitizeMediaFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	return filename
}

func getMimeTypeFromExt(ext string) string {
	types := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".webp": "image/webp",
	}
	if t, ok := types[ext]; ok {
		return t
	}
	return "application/octet-stream"
}

func getSessionDisplayName(c echo.Context) string {
	if sess, ok := c.Get("session").(*customMiddleware.Session); ok {
		return sess.DisplayName
	}
	return "Admin"
}

func getSessionRole(c echo.Context) string {
	if sess, ok := c.Get("session").(*customMiddleware.Session); ok {
		return sess.Role
	}
	return ""
}
