package admin_test

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// Use the actual renderer: template parsing alone does not detect the old
// built-in slice execution failure which made every Global Settings tab fail.
func TestGlobalSettingsRendersAllTabs(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	h := admin.NewSettingsHandler(queries, logger, nil)
	for _, tab := range []string{"general", "contact", "social", "marketplaces", "seo", "invalid"} {
		t.Run(tab, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/settings?tab="+tab, nil), rec)
			if err := h.Edit(c); err != nil {
				t.Fatalf("render settings: %v", err)
			}
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d", rec.Code)
			}
			for _, field := range []string{"contact_email", "footer_logo_path", "social_twitter", "social_linkedin", "social_instagram", "social_threads", "marketplace_gem_url", "marketplace_amazon_url"} {
				if !strings.Contains(rec.Body.String(), `name="`+field+`"`) {
					t.Errorf("missing field %s", field)
				}
			}
		})
	}
}

func settingsFormFixture() url.Values {
	return url.Values{
		"site_name": {"QA Company"}, "active_tab": {"social"},
		"contact_email": {"qa@example.com"}, "contact_phone": {"+91 9876543210"},
		"address": {"QA office\nSecond line"}, "business_hours": {"Monday–Friday"},
		"footer_logo_path": {"/uploads/qa-logo.png"},
		"social_twitter":   {"https://x.com/qa"}, "social_linkedin": {"https://linkedin.com/company/qa"},
		"social_instagram": {"https://instagram.com/qa"}, "social_threads": {"https://www.threads.net/@qa"},
		"marketplace_gem_url": {"https://gem.gov.in/qa-listing"}, "marketplace_amazon_url": {"https://amazon.in/qa-listing"},
	}
}

func renderSettingsFooter(t *testing.T, settings sqlc.Setting) string {
	t.Helper()
	tmpl, err := template.New("footer").Funcs(template.FuncMap{"now": time.Now}).ParseFiles("../../../templates/partials/footer.html")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := tmpl.ExecuteTemplate(&output, "footer", map[string]interface{}{"Settings": settings}); err != nil {
		t.Fatal(err)
	}
	return output.String()
}

func TestGlobalSettingsSaveAppearsInFooterAndClearsCache(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	cache := services.NewCache()
	cache.Set("page:home", "old footer", 60)
	e := echo.New()
	h := admin.NewSettingsHandler(queries, logger, cache)
	rec, c := postForm(e, "/admin/settings", settingsFormFixture())
	if err := h.Update(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status = %d", rec.Code)
	}
	if _, found := cache.Get("page:home"); found {
		t.Fatal("stale page was not invalidated")
	}
	settings, err := queries.GetSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	html := renderSettingsFooter(t, settings)
	for _, want := range []string{`src="/uploads/qa-logo.png"`, `mailto:qa@example.com`, `tel:&#43;91%209876543210`, "QA office", "Second line", "Monday–Friday", "https://x.com/qa", "https://linkedin.com/company/qa", "https://instagram.com/qa", "https://www.threads.net/@qa", "https://gem.gov.in/qa-listing", "https://amazon.in/qa-listing"} {
		if !strings.Contains(html, want) {
			t.Errorf("footer missing %q", want)
		}
	}
	// Clearing optional URLs must remove the links without affecting contact data.
	form := settingsFormFixture()
	for _, field := range []string{"social_threads", "marketplace_gem_url", "marketplace_amazon_url"} {
		form.Set(field, "")
	}
	_, c = postForm(e, "/admin/settings", form)
	if err := h.Update(c); err != nil {
		t.Fatal(err)
	}
	settings, err = queries.GetSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	html = renderSettingsFooter(t, settings)
	for _, absent := range []string{`aria-label="Threads"`, `aria-label="BlueJay on GeM"`, `aria-label="BlueJay on Amazon"`, "Find us on"} {
		if strings.Contains(html, absent) {
			t.Errorf("cleared footer link still renders: %s", absent)
		}
	}
}

func TestGlobalSettingsRejectsUnsafeURLsWithoutUpdating(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	original, err := queries.GetSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	h := admin.NewSettingsHandler(queries, logger, nil)
	for _, field := range []string{"social_threads", "social_twitter", "marketplace_gem_url", "marketplace_amazon_url", "footer_logo_path"} {
		form := settingsFormFixture()
		form.Set(field, "javascript:alert(1)")
		_, c := postForm(echo.New(), "/admin/settings", form)
		err := h.Update(c)
		httpErr, ok := err.(*echo.HTTPError)
		if !ok || httpErr.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected 400, got %v", field, err)
		}
		settings, err := queries.GetSettings(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if settings.SiteName != original.SiteName {
			t.Fatal("invalid settings partially persisted")
		}
	}
}

func TestFooterUsesHeaderLogoWhenFooterLogoUnset(t *testing.T) {
	html := renderSettingsFooter(t, sqlc.Setting{SiteName: "QA Company", ShowFooterAbout: true, HeaderLogoPath: "/uploads/header.png"})
	if !strings.Contains(html, `src="/uploads/header.png"`) {
		t.Fatal("existing header logo not used as footer fallback")
	}
}
