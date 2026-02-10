package admin

import (
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type ProductDetailsHandler struct {
	queries   *sqlc.Queries
	logger    *slog.Logger
	uploadSvc *services.UploadService
	partials  map[string]*template.Template
}

func NewProductDetailsHandler(queries *sqlc.Queries, logger *slog.Logger, uploadSvc *services.UploadService) *ProductDetailsHandler {
	h := &ProductDetailsHandler{
		queries:   queries,
		logger:    logger,
		uploadSvc: uploadSvc,
		partials:  make(map[string]*template.Template),
	}
	h.loadPartials()
	return h
}

func (h *ProductDetailsHandler) loadPartials() {
	basePath := "templates"
	names := []string{
		"product_specs",
		"product_features",
		"product_certifications",
		"product_downloads",
		"product_images",
	}
	for _, name := range names {
		h.partials[name] = template.Must(template.ParseFiles(
			filepath.Join(basePath, "admin/partials", name+".html"),
		))
	}
}

func (h *ProductDetailsHandler) renderPartial(c echo.Context, name string, data interface{}) error {
	tmpl, ok := h.partials[name]
	if !ok {
		return fmt.Errorf("partial not found: %s", name)
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return tmpl.ExecuteTemplate(c.Response(), name, data)
}

// --- Specs ---

func (h *ProductDetailsHandler) ListSpecs(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	specs, err := h.queries.ListProductSpecs(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list specs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return h.renderPartial(c, "product_specs", map[string]interface{}{
		"ProductID": id,
		"Specs":     specs,
	})
}

func (h *ProductDetailsHandler) AddSpec(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	_, err := h.queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID:    id,
		SectionName:  c.FormValue("section_name"),
		SpecKey:      c.FormValue("spec_key"),
		SpecValue:    c.FormValue("spec_value"),
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create spec", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Added spec to Product #%d", id)
	return h.ListSpecs(c)
}

func (h *ProductDetailsHandler) DeleteSpecs(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteProductSpecs(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete specs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "product", id, "", "Deleted specs from Product #%d", id)
	return h.ListSpecs(c)
}

// --- Features ---

func (h *ProductDetailsHandler) ListFeatures(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	features, err := h.queries.ListProductFeatures(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list features", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return h.renderPartial(c, "product_features", map[string]interface{}{
		"ProductID": id,
		"Features":  features,
	})
}

func (h *ProductDetailsHandler) AddFeature(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	_, err := h.queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    id,
		FeatureText:  c.FormValue("feature_text"),
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create feature", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Added feature to Product #%d", id)
	return h.ListFeatures(c)
}

func (h *ProductDetailsHandler) DeleteFeatures(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteProductFeatures(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete features", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "product", id, "", "Deleted features from Product #%d", id)
	return h.ListFeatures(c)
}

// --- Certifications ---

func (h *ProductDetailsHandler) ListCertifications(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	certs, err := h.queries.ListProductCertifications(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list certifications", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return h.renderPartial(c, "product_certifications", map[string]interface{}{
		"ProductID":      id,
		"Certifications": certs,
	})
}

func (h *ProductDetailsHandler) AddCertification(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	certCode := c.FormValue("certification_code")
	iconType := c.FormValue("icon_type")
	iconPath := c.FormValue("icon_path")

	_, err := h.queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID:         id,
		CertificationName: c.FormValue("certification_name"),
		CertificationCode: sql.NullString{String: certCode, Valid: certCode != ""},
		IconType:          sql.NullString{String: iconType, Valid: iconType != ""},
		IconPath:          sql.NullString{String: iconPath, Valid: iconPath != ""},
		DisplayOrder:      order,
	})
	if err != nil {
		h.logger.Error("failed to create certification", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Added certification to Product #%d", id)
	return h.ListCertifications(c)
}

func (h *ProductDetailsHandler) DeleteCertifications(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteProductCertifications(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete certifications", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "product", id, "", "Deleted certifications from Product #%d", id)
	return h.ListCertifications(c)
}

// --- Downloads ---

func (h *ProductDetailsHandler) ListDownloads(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	downloads, err := h.queries.ListProductDownloads(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list downloads", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return h.renderPartial(c, "product_downloads", map[string]interface{}{
		"ProductID": id,
		"Downloads": downloads,
	})
}

func (h *ProductDetailsHandler) AddDownload(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}

	path, err := h.uploadSvc.UploadProductDownload(fileHeader)
	if err != nil {
		h.logger.Error("failed to upload download", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to upload file: "+err.Error())
	}

	desc := c.FormValue("description")
	version := c.FormValue("version")
	fileType := c.FormValue("file_type")
	if fileType == "" {
		fileType = filepath.Ext(fileHeader.Filename)
	}

	_, err = h.queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID:    id,
		Title:        c.FormValue("title"),
		Description:  sql.NullString{String: desc, Valid: desc != ""},
		FileType:     fileType,
		FilePath:     path,
		FileSize:     sql.NullInt64{Int64: fileHeader.Size, Valid: true},
		Version:      sql.NullString{String: version, Valid: version != ""},
		DisplayOrder: order,
	})
	if err != nil {
		h.logger.Error("failed to create download", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Added download to Product #%d", id)
	return h.ListDownloads(c)
}

func (h *ProductDetailsHandler) DeleteDownload(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	downloadID, _ := strconv.ParseInt(c.Param("download_id"), 10, 64)
	if err := h.queries.DeleteProductDownload(c.Request().Context(), downloadID); err != nil {
		h.logger.Error("failed to delete download", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "product", id, "", "Deleted download from Product #%d", id)
	return h.ListDownloads(c)
}

// --- Images ---

func (h *ProductDetailsHandler) ListImages(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	images, err := h.queries.ListProductImages(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list images", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return h.renderPartial(c, "product_images", map[string]interface{}{
		"ProductID": id,
		"Images":    images,
	})
}

func (h *ProductDetailsHandler) AddImage(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
	isThumbnail := c.FormValue("is_thumbnail") == "1"

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Image file is required")
	}

	path, err := h.uploadSvc.UploadProductImage(fileHeader)
	if err != nil {
		h.logger.Error("failed to upload image", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to upload image: "+err.Error())
	}

	altText := c.FormValue("alt_text")
	caption := c.FormValue("caption")

	_, err = h.queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID:    id,
		ImagePath:    path,
		AltText:      sql.NullString{String: altText, Valid: altText != ""},
		Caption:      sql.NullString{String: caption, Valid: caption != ""},
		DisplayOrder: order,
		IsThumbnail:  isThumbnail,
	})
	if err != nil {
		h.logger.Error("failed to create image", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Added image to Product #%d", id)
	return h.ListImages(c)
}

func (h *ProductDetailsHandler) DeleteImage(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	imageID, _ := strconv.ParseInt(c.Param("image_id"), 10, 64)
	if err := h.queries.DeleteProductImage(c.Request().Context(), imageID); err != nil {
		h.logger.Error("failed to delete image", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	logActivity(c, "updated", "product", id, "", "Deleted image from Product #%d", id)
	return h.ListImages(c)
}
