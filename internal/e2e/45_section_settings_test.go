package e2e_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAboutSettings(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("load settings page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/about/settings", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("update about settings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/about/settings", strings.NewReader(url.Values{
			"about_show_mission":        {"on"},
			"about_show_milestones":     {"on"},
			"about_show_certifications": {},
			"about_show_team":           {"on"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		settings, _ := queries.GetSettings(context.Background())
		if settings.AboutShowMission != 1 {
			t.Error("expected about_show_mission to be 1")
		}
		if settings.AboutShowCertifications != 0 {
			t.Error("expected about_show_certifications to be 0")
		}
	})
}

func TestProductsSettings(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("update products settings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/products/settings", strings.NewReader(url.Values{
			"products_per_page":        {"24"},
			"products_show_categories": {"on"},
			"products_show_search":     {},
			"products_default_sort":    {"newest"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		settings, _ := queries.GetSettings(context.Background())
		if settings.ProductsPerPage != 24 {
			t.Errorf("expected products_per_page 24, got %d", settings.ProductsPerPage)
		}
		if settings.ProductsShowSearch != 0 {
			t.Error("expected products_show_search to be 0")
		}
	})
}

func TestSolutionsSettings(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("update solutions settings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/solutions/settings", strings.NewReader(url.Values{
			"solutions_per_page":        {"18"},
			"solutions_show_industries": {"on"},
			"solutions_show_search":     {"on"},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		settings, _ := queries.GetSettings(context.Background())
		if settings.SolutionsPerPage != 18 {
			t.Errorf("expected solutions_per_page 18, got %d", settings.SolutionsPerPage)
		}
	})
}

func TestBlogSettings(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("update blog settings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/settings", strings.NewReader(url.Values{
			"blog_posts_per_page":   {"15"},
			"blog_show_author":      {"on"},
			"blog_show_date":        {"on"},
			"blog_show_categories":  {},
			"blog_show_tags":        {"on"},
			"blog_show_search":      {},
		}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("expected 303, got %d", rec.Code)
		}

		settings, _ := queries.GetSettings(context.Background())
		if settings.BlogPostsPerPage != 15 {
			t.Errorf("expected blog_posts_per_page 15, got %d", settings.BlogPostsPerPage)
		}
		if settings.BlogShowCategories != 0 {
			t.Error("expected blog_show_categories to be 0")
		}
	})
}
