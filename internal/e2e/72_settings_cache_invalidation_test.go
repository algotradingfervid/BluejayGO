package e2e_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	publicHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	appmw "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestSettingsUpdate_InvalidatesPublicPageCache verifies that saving Global Settings
// in the admin panel busts the rendered-page cache so that cleared contact details
// disappear from public pages immediately (rather than after the 1h TTL).
//
// Bug: SettingsHandler.Update persisted the cleared value to the DB but never
// invalidated the "page:" rendered-HTML cache, so cached public pages (e.g. the
// footer on /contact) kept showing the old contact phone until the cache expired.
//
// This test shares ONE appCache between the public contact page and the admin
// settings route, uses the REAL renderer (the shared setupApp stub renderer cannot
// see template output), primes the cache by GETting /contact, clears contact_phone
// via POST /admin/settings, then re-GETs /contact and asserts the old phone is gone.
func TestSettingsUpdate_InvalidatesPublicPageCache(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	const oldPhone = "555-CACHE-BUST"

	// Seed the contact phone so it renders in the public footer.
	if err := queries.UpdateGlobalSettings(ctx, sqlc.UpdateGlobalSettingsParams{
		SiteName:     "Test Site",
		ContactPhone: oldPhone,
		ContactEmail: "info@example.com",
	}); err != nil {
		t.Fatalf("seed settings: %v", err)
	}

	// Shared cache instance used by BOTH the public handler and the admin handler.
	appCache := services.NewCache()

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")

	// Public contact page: SettingsLoader populates the footer settings, and the
	// contact handler renders+caches the page under "page:contact".
	contactHandler := publicHandlers.NewContactHandler(queries, logger, appCache)
	e.GET("/contact", contactHandler.ShowContactPage, appmw.SettingsLoader(queries))

	// Admin settings route sharing the SAME appCache.
	settingsHandler := adminHandlers.NewSettingsHandler(queries, logger, appCache)
	e.POST("/admin/settings", settingsHandler.Update)

	// 1) Prime the cache: GET /contact so its rendered HTML (with oldPhone) is cached.
	req := httptest.NewRequest(http.MethodGet, "/contact", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("prime GET /contact: expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), oldPhone) {
		t.Fatalf("precondition failed: expected old phone %q to render in /contact footer, body did not contain it", oldPhone)
	}

	// 2) Clear contact_phone via the admin settings form (keep other fields populated).
	form := url.Values{}
	form.Set("site_name", "Test Site")
	form.Set("contact_email", "info@example.com")
	form.Set("contact_phone", "") // <-- clearing the phone
	form.Set("active_tab", "general")

	postReq := httptest.NewRequest(http.MethodPost, "/admin/settings", strings.NewReader(form.Encode()))
	postReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	postRec := httptest.NewRecorder()
	e.ServeHTTP(postRec, postReq)
	if postRec.Code != http.StatusSeeOther {
		t.Fatalf("POST /admin/settings: expected 303, got %d; body: %s", postRec.Code, postRec.Body.String())
	}

	// (a) Data-layer clear must persist.
	settings, err := queries.GetSettings(ctx)
	if err != nil {
		t.Fatalf("GetSettings after update: %v", err)
	}
	if settings.ContactPhone != "" {
		t.Fatalf("expected ContactPhone cleared in DB, got %q", settings.ContactPhone)
	}

	// (b) Cache must have been invalidated: re-GET /contact and assert old phone is gone.
	req2 := httptest.NewRequest(http.MethodGet, "/contact", nil)
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("re-GET /contact: expected 200, got %d", rec2.Code)
	}
	if strings.Contains(rec2.Body.String(), oldPhone) {
		t.Fatalf("stale cache: old phone %q still present in /contact after settings update — cache was not invalidated", oldPhone)
	}
}
