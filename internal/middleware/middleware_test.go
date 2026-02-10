package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/internal/middleware"
)

func init() {
	middleware.InitSessionStore("test-secret-at-least-32-characters-long")
}

func TestSessionMiddleware_NewSession(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := middleware.SessionMiddleware()(func(c echo.Context) error {
		sess, ok := c.Get("session").(*middleware.Session)
		if !ok {
			t.Fatal("session not set in context")
		}
		if sess.UserID != 0 {
			t.Errorf("expected UserID 0 for new session, got %d", sess.UserID)
		}
		return c.String(http.StatusOK, "ok")
	})

	if err := handler(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
}

func TestRequireAuth_NoSession(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// No session set
	handler := middleware.RequireAuth()(func(c echo.Context) error {
		t.Fatal("should not reach handler")
		return nil
	})

	err := handler(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected redirect 303, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/admin/login" {
		t.Errorf("expected redirect to /admin/login, got %q", loc)
	}
}

func TestRequireAuth_WithSession(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set session with valid user
	c.Set("session", &middleware.Session{
		UserID:      1,
		Email:       "admin@test.com",
		DisplayName: "Admin",
		Role:        "admin",
	})

	reached := false
	handler := middleware.RequireAuth()(func(c echo.Context) error {
		reached = true
		return c.String(http.StatusOK, "ok")
	})

	if err := handler(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !reached {
		t.Error("handler was not reached")
	}
}

func TestRequireRole_Allowed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("session", &middleware.Session{
		UserID: 1, Role: "admin",
	})

	reached := false
	handler := middleware.RequireRole("admin", "superadmin")(func(c echo.Context) error {
		reached = true
		return c.String(http.StatusOK, "ok")
	})

	if err := handler(c); err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !reached {
		t.Error("handler was not reached for allowed role")
	}
}

func TestRequireRole_Denied(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("session", &middleware.Session{
		UserID: 1, Role: "editor",
	})

	handler := middleware.RequireRole("admin")(func(c echo.Context) error {
		t.Fatal("should not reach handler for denied role")
		return nil
	})

	err := handler(c)
	if err == nil {
		t.Fatal("expected error for denied role")
	}
	he, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected HTTPError, got %T", err)
	}
	if he.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", he.Code)
	}
}

func TestRequireAuth_EmptySession(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Session exists but UserID is 0
	c.Set("session", &middleware.Session{UserID: 0})

	handler := middleware.RequireAuth()(func(c echo.Context) error {
		t.Fatal("should not reach handler with UserID 0")
		return nil
	})

	if err := handler(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
}
