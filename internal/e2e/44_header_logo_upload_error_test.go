package e2e_test

import (
	"bytes"
	"context"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestHeaderLogoUpload_UnsupportedType_RerendersFormWithError verifies that when the
// admin uploads a logo file with an unsupported extension (e.g. .avif), the handler
// does NOT dump the user out to a raw JSON error page. Instead it re-renders the
// Header Management form (HTTP 400) with a friendly, styled inline error that names
// the offending type and lists the supported formats + size limit.
//
// The bug: HeaderHandler.Update returned echo.NewHTTPError(400, "Failed to upload
// logo: invalid file type: .avif"), which Echo renders as bare JSON
// {"message":"..."} — ugly and destroys the form.
//
// The shared setupApp uses a stub renderer that cannot show template output, so this
// test builds a local Echo with the REAL renderer (mirroring test 43) so the body
// assertions are meaningful.
func TestHeaderLogoUpload_UnsupportedType_RerendersFormWithError(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	_ = context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	uploadSvc := services.NewUploadService(t.TempDir())

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	h := adminHandlers.NewHeaderHandler(queries, logger, uploadSvc, services.NewCache())
	e.POST("/admin/header", h.Update)

	// Build a multipart form with an unsupported logo file extension.
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("header_logo_alt", "Bad Logo"); err != nil {
		t.Fatal(err)
	}
	part, err := w.CreateFormFile("header_logo_file", "bad-logo.avif")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("not-a-real-image")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/header", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 with re-rendered form, got %d; body: %s", rec.Code, rec.Body.String())
	}

	bodyStr := rec.Body.String()

	// Must NOT be the bare Echo JSON error page.
	if strings.Contains(bodyStr, `{"message"`) {
		t.Fatalf("expected rendered HTML form, got raw JSON error page: %s", bodyStr)
	}

	// Must be the rendered Header Management form.
	if !strings.Contains(bodyStr, "Header Management") {
		t.Errorf("expected re-rendered Header Management form, body did not contain it")
	}
	if !strings.Contains(bodyStr, `id="header-form"`) {
		t.Errorf("expected the header form to be present in the re-rendered page")
	}

	// Must contain a friendly error that names the bad extension and supported formats.
	if !strings.Contains(bodyStr, "not a supported image type") {
		t.Errorf("expected friendly 'not a supported image type' message, body: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, ".avif") {
		t.Errorf("expected the offending extension '.avif' to be named in the error")
	}
	if !strings.Contains(bodyStr, "SVG") {
		t.Errorf("expected the supported-formats hint (e.g. SVG) in the error")
	}
	if !strings.Contains(bodyStr, "5MB") {
		t.Errorf("expected the 5MB size limit mentioned in the error")
	}
}
