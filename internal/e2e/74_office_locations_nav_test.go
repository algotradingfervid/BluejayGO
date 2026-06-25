package e2e_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestAdminSidebar_HasOfficeLocationsLink verifies the admin sidebar exposes a link
// to the Office Locations editor (Task 14 — the editor already worked but was not
// reachable from the nav). Uses the REAL renderer so the sidebar partial is rendered.
func TestAdminSidebar_HasOfficeLocationsLink(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	customMiddleware.InitSessionStore("e2e-test-secret-at-least-32-characters-long")

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")
	e.Use(customMiddleware.SecurityHeaders())
	e.Use(customMiddleware.SessionMiddleware())

	appCache := services.NewCache()
	activitySvc := services.NewActivityLogService(queries, logger)
	adminHandlers.SetActivityLogService(activitySvc)

	authHandler := adminHandlers.NewAuthHandler(queries, logger)
	e.GET("/admin/login", authHandler.ShowLoginPage)
	e.POST("/admin/login", authHandler.LoginSubmit)

	adminGroup := e.Group("/admin", customMiddleware.RequireAuth())
	contactHandler := adminHandlers.NewAdminContactHandler(queries, logger, appCache)
	adminGroup.GET("/contact/offices", contactHandler.ListOffices)

	cookie := loginTabsAdmin(t, e, queries)

	req := httptest.NewRequest(http.MethodGet, "/admin/contact/offices", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /admin/contact/offices, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	// Sidebar-specific marker (data-path is only used by sidebar links).
	if !strings.Contains(body, `data-path="/admin/contact/offices"`) {
		t.Errorf("admin sidebar should contain an Office Locations link (data-path=/admin/contact/offices), but it did not")
	}
	if !strings.Contains(body, "Office Locations") {
		t.Errorf("expected 'Office Locations' label in the admin sidebar")
	}
}
