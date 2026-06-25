package e2e_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestHeaderLogo_RendersConfiguredLogoOnPublicHeader verifies that when an admin
// configures a header logo path (via the header/branding settings), the public site
// header actually renders that logo as an <img>. The bug: the public header partial
// ignored Settings.HeaderLogoPath entirely, so uploading a logo had no visible effect.
//
// The shared setupApp uses a stub renderer, so this test builds a local Echo with the
// REAL renderer (mirroring production) and asserts on the rendered home page markup.
func TestHeaderLogo_RendersConfiguredLogoOnPublicHeader(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	const logoPath = "/uploads/branding/test-logo.png"
	const logoAlt = "Acme Logo"

	// Simulate the admin saving a logo through the header settings form.
	if err := queries.UpdateHeaderSettings(ctx, sqlc.UpdateHeaderSettingsParams{
		HeaderLogoPath: logoPath,
		HeaderLogoAlt:  logoAlt,
	}); err != nil {
		t.Fatalf("update header settings: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	h := publicHandlers.NewHomeHandler(queries, logger)
	e.GET("/", h.ShowHomePage)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	// The configured logo image path must appear in the rendered header.
	if !strings.Contains(body, logoPath) {
		t.Errorf("expected configured logo path %q to render in the public header, but it was absent", logoPath)
	}
	// It should be rendered as an <img> with the configured alt text.
	if !strings.Contains(body, `src="`+logoPath+`"`) {
		t.Errorf("expected an <img src=%q> in the header, body did not contain it", logoPath)
	}
	if !strings.Contains(body, logoAlt) {
		t.Errorf("expected logo alt text %q to render, but it was absent", logoAlt)
	}
}
