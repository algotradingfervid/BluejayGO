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

// renderHomeBody renders the public home page using the REAL renderer (mirroring
// production) and returns the response body. A fresh Echo instance is built each
// call so no rendered-HTML cache can mask the template's actual output.
func renderHomeBody(t *testing.T, queries *sqlc.Queries) string {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

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
	return rec.Body.String()
}

// headerSection returns just the <header>...</header> portion of the rendered
// page so assertions about the header are not confounded by the footer (which
// also renders contact/social info via its own separate toggles).
func headerSection(t *testing.T, body string) string {
	t.Helper()
	start := strings.Index(body, "<header")
	if start < 0 {
		t.Fatalf("no <header> in rendered body")
	}
	end := strings.Index(body[start:], "</header>")
	if end < 0 {
		t.Fatalf("no </header> in rendered body")
	}
	return body[start : start+end]
}

// TestHeaderContactInfo_TogglesControlPublicHeader verifies that the admin
// "show phone / email / social" header toggles actually control whether the
// contact info renders in the public site header. The bug: the public header
// partial rendered no contact info at all, so the toggles controlled nothing.
func TestHeaderContactInfo_TogglesControlPublicHeader(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	const phone = "1-555-0100"
	const email = "hello@acme.test"
	const linkedin = "https://linkedin.com/company/acme"

	// Seed contact info + social link into global settings.
	if err := queries.UpdateGlobalSettings(ctx, sqlc.UpdateGlobalSettingsParams{
		SiteName:       "Acme",
		ContactEmail:   email,
		ContactPhone:   phone,
		SocialLinkedin: linkedin,
	}); err != nil {
		t.Fatalf("update global settings: %v", err)
	}

	// Helper to set the three header contact toggles to a given state.
	setToggles := func(on bool) {
		s, err := queries.GetSettings(ctx)
		if err != nil {
			t.Fatalf("get settings: %v", err)
		}
		if err := queries.UpdateHeaderSettings(ctx, sqlc.UpdateHeaderSettingsParams{
			HeaderLogoPath:      s.HeaderLogoPath,
			HeaderLogoAlt:       s.HeaderLogoAlt,
			HeaderCtaEnabled:    s.HeaderCtaEnabled,
			HeaderCtaText:       s.HeaderCtaText,
			HeaderCtaUrl:        s.HeaderCtaUrl,
			HeaderCtaStyle:      s.HeaderCtaStyle,
			HeaderShowPhone:     on,
			HeaderShowEmail:     on,
			HeaderShowSocial:    on,
			HeaderSocialStyle:   s.HeaderSocialStyle,
			ShowNavProducts:     s.ShowNavProducts,
			ShowNavSolutions:    s.ShowNavSolutions,
			ShowNavCaseStudies:  s.ShowNavCaseStudies,
			ShowNavAbout:        s.ShowNavAbout,
			ShowNavBlog:         s.ShowNavBlog,
			ShowNavWhitepapers:  s.ShowNavWhitepapers,
			ShowNavPartners:     s.ShowNavPartners,
			ShowNavContact:      s.ShowNavContact,
			NavLabelProducts:    s.NavLabelProducts,
			NavLabelSolutions:   s.NavLabelSolutions,
			NavLabelCaseStudies: s.NavLabelCaseStudies,
			NavLabelAbout:       s.NavLabelAbout,
			NavLabelBlog:        s.NavLabelBlog,
			NavLabelWhitepapers: s.NavLabelWhitepapers,
			NavLabelPartners:    s.NavLabelPartners,
			NavLabelContact:     s.NavLabelContact,
		}); err != nil {
			t.Fatalf("update header settings: %v", err)
		}
	}

	// --- Toggles ON: contact info must render in the header ---
	setToggles(true)
	header := headerSection(t, renderHomeBody(t, queries))

	if !strings.Contains(header, "tel:"+phone) {
		t.Errorf("toggles ON: expected phone link %q in header, but it was absent", "tel:"+phone)
	}
	if !strings.Contains(header, "mailto:"+email) {
		t.Errorf("toggles ON: expected email link %q in header, but it was absent", "mailto:"+email)
	}
	if !strings.Contains(header, linkedin) {
		t.Errorf("toggles ON: expected social link %q in header, but it was absent", linkedin)
	}

	// --- Toggles OFF: contact info must NOT render in the header ---
	setToggles(false)
	header = headerSection(t, renderHomeBody(t, queries))

	if strings.Contains(header, "tel:"+phone) {
		t.Errorf("toggles OFF: phone link %q should be hidden, but it rendered", "tel:"+phone)
	}
	if strings.Contains(header, "mailto:"+email) {
		t.Errorf("toggles OFF: email link %q should be hidden, but it rendered", "mailto:"+email)
	}
	if strings.Contains(header, linkedin) {
		t.Errorf("toggles OFF: social link %q should be hidden, but it rendered", linkedin)
	}
}
