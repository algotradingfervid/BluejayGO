package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession_CookieSetOnLogin(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	if cookie.Name != "bluejay_session" {
		t.Errorf("expected cookie name 'bluejay_session', got %s", cookie.Name)
	}

	if !cookie.HttpOnly {
		t.Error("expected HttpOnly=true")
	}

	if cookie.Path != "/" {
		t.Errorf("expected Path=/, got %s", cookie.Path)
	}

	if cookie.MaxAge != 86400*7 {
		t.Errorf("expected MaxAge=604800, got %d", cookie.MaxAge)
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax, got %d", cookie.SameSite)
	}
}

func TestSession_ValidSessionAccessesProtectedRoute(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusSeeOther {
		t.Error("valid session should not redirect to login")
	}
}

func TestSession_TamperedCookieRejected(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	cookie.Value = "tampered_value"

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("tampered cookie should redirect, got %d", rec.Code)
	}
}

func TestSession_MissingCookieRedirects(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("missing cookie should redirect, got %d", rec.Code)
	}

	if rec.Header().Get("Location") != "/admin/login" {
		t.Errorf("expected redirect to /admin/login, got %s", rec.Header().Get("Location"))
	}
}

func TestSession_LogoutInvalidatesSession(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/logout", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("logout should redirect, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	for _, c := range rec.Result().Cookies() {
		req.AddCookie(c)
	}
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("invalidated session should redirect, got %d", rec.Code)
	}
}
