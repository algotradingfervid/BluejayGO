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

// TestFooterResources_CaseStudiesAndWhitepapersHidden verifies that migration 035
// deactivates the "Case Studies" and "Whitepapers" footer RESOURCES links so they
// no longer render on the public front end, while "Blog" and "Support" remain.
//
// The shared setupApp uses a stub renderer and does not wire the SettingsLoader
// middleware, so it cannot exercise the footer. This test builds a local Echo
// instance with the REAL renderer and the SettingsLoader middleware on the public
// group, then renders the /about page (which includes the base layout footer) and
// asserts the footer markup.
func TestFooterResources_CaseStudiesAndWhitepapersHidden(t *testing.T) {
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

	// The two hidden resource links must NOT be present anywhere in the footer.
	if strings.Contains(body, `href="/case-studies"`) {
		t.Errorf("footer should NOT contain a link to /case-studies, but it does")
	}
	if strings.Contains(body, `href="/whitepapers"`) {
		t.Errorf("footer should NOT contain a link to /whitepapers, but it does")
	}

	// The remaining resource links must STILL be present, proving only the two
	// were hidden and the RESOURCES section / mechanism is intact.
	if !strings.Contains(body, `href="/blog"`) {
		t.Errorf("footer should still contain a link to /blog, but it does not; body: %s", body)
	}
	// The Support resource link still renders; migration 036 repointed its URL
	// from the dead /support route to the existing /contact page.
	if !strings.Contains(body, `href="/contact"`) {
		t.Errorf("footer should still contain the Support link (now pointing to /contact), but it does not")
	}
}
