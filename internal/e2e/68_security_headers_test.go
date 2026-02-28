package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSecurityHeaders_PublicPages(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	pages := []string{"/", "/products", "/contact", "/health"}

	for _, page := range pages {
		t.Run(page, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, page, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			headers := rec.Header()

			if headers.Get("X-Content-Type-Options") != "nosniff" {
				t.Errorf("expected X-Content-Type-Options=nosniff, got %s", headers.Get("X-Content-Type-Options"))
			}

			if headers.Get("X-Frame-Options") != "DENY" {
				t.Errorf("expected X-Frame-Options=DENY, got %s", headers.Get("X-Frame-Options"))
			}

			if headers.Get("X-XSS-Protection") != "1; mode=block" {
				t.Errorf("expected X-XSS-Protection=1; mode=block, got %s", headers.Get("X-XSS-Protection"))
			}

			if headers.Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
				t.Errorf("expected Referrer-Policy=strict-origin-when-cross-origin, got %s", headers.Get("Referrer-Policy"))
			}

			csp := headers.Get("Content-Security-Policy")
			if !strings.Contains(csp, "default-src 'self'") {
				t.Errorf("CSP missing default-src 'self'")
			}
		})
	}
}

func TestSecurityHeaders_AdminPages(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	headers := rec.Header()

	if headers.Get("X-Content-Type-Options") != "nosniff" {
		t.Error("X-Content-Type-Options not set on admin page")
	}

	if headers.Get("X-Frame-Options") != "DENY" {
		t.Error("X-Frame-Options not set on admin page")
	}
}

func TestSecurityHeaders_CSP_Details(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	csp := rec.Header().Get("Content-Security-Policy")

	expectedDirectives := []string{
		"default-src 'self'",
		"script-src",
		"style-src",
		"font-src",
		"img-src",
	}

	for _, directive := range expectedDirectives {
		if !strings.Contains(csp, directive) {
			t.Errorf("CSP missing directive: %s", directive)
		}
	}
}

func TestSecurityHeaders_ErrorResponses(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("security headers missing on 404 response")
	}
}

func TestSecurityHeaders_StaticFiles(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/public/test.css", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("security headers missing on static file response")
	}
}
