package e2e_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestFooterSupportLink_PointsToContact verifies that migration 036 repoints the
// footer RESOURCES "Support" link from the dead /support route (which 404s) to the
// existing /contact page, which already handles support inquiries.
//
// The shared setupApp uses a stub renderer and does not wire the SettingsLoader
// middleware, so it cannot exercise the footer. This test builds a local Echo
// instance with the REAL renderer and the SettingsLoader middleware on the public
// group, then renders the /about page (which includes the base layout footer) and
// asserts the footer markup.
func TestFooterSupportLink_PointsToContact(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	appCache := services.NewCache()
	aboutHandler := publicHandlers.NewAboutHandler(queries, logger, appCache)

	// Public group with the real settings/footer loader, mirroring production wiring.
	pub := e.Group("", customMiddleware.SettingsLoader(queries))
	pub.GET("/about", aboutHandler.AboutPage)

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /about, got %d; body: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()

	// The Support link must point to the real /contact page.
	if !strings.Contains(body, `href="/contact"`) {
		t.Errorf("footer should contain a Support link to /contact, but it does not; body: %s", body)
	}

	// The dead /support link must NOT be present anywhere in the footer.
	if strings.Contains(body, `href="/support"`) {
		t.Errorf("footer should NOT contain a link to /support (it 404s), but it does")
	}
}
