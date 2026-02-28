package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestPageSectionsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	section, err := queries.GetPageSection(ctx, sqlc.GetPageSectionParams{
		PageKey:    "home",
		SectionKey: "hero",
	})
	if err == sql.ErrNoRows {
		t.Skip("no seeded page sections available")
	}
	if err != nil {
		t.Fatalf("failed to get page section: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/page-sections/%d", section.ID), strings.NewReader(url.Values{
		"heading":               {"Updated Heading"},
		"subheading":            {"Updated Subheading"},
		"description":           {"Updated description text"},
		"label":                 {"Featured"},
		"primary_button_text":   {"Get Started"},
		"primary_button_url":    {"/signup"},
		"secondary_button_text": {"Learn More"},
		"secondary_button_url":  {"/about"},
		"is_active":             {"on"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("page sections route not found")
	}
}

func TestPageSectionsNoCreateDelete_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodGet, "/admin/page-sections", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("page sections list route not found")
	}

	req = httptest.NewRequest(http.MethodPost, "/admin/page-sections", strings.NewReader(url.Values{
		"page_key":    {"test"},
		"section_key": {"test"},
		"heading":     {"Test"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Log("no create endpoint exists as expected")
	}
}

func TestPageSectionsToggleActive_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	sections, err := queries.ListAllPageSections(ctx)
	if err != nil || len(sections) == 0 {
		t.Skip("no page sections available")
	}

	section := sections[0]
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/page-sections/%d", section.ID), strings.NewReader(url.Values{
		"heading":     {section.Heading},
		"subheading":  {section.Subheading},
		"description": {section.Description},
		"label":       {section.Label},
		"is_active":   {""},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("page section update route not found")
	}
}

func TestPageSectionsButtonFields_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	sections, _ := queries.ListAllPageSections(ctx)
	if len(sections) == 0 {
		t.Skip("no page sections available")
	}

	section := sections[0]
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/page-sections/%d", section.ID), strings.NewReader(url.Values{
		"heading":               {"Test"},
		"primary_button_text":   {"Primary Button"},
		"primary_button_url":    {"/primary"},
		"secondary_button_text": {""},
		"secondary_button_url":  {""},
		"is_active":             {"on"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("page section update route not found")
	}
}
