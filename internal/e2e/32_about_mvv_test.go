package e2e_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestAboutMVV_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	queries.UpsertMissionVisionValues(ctx, sqlc.UpsertMissionVisionValuesParams{
		Mission:       "To protect and secure",
		Vision:        "A safer world",
		ValuesSummary: sql.NullString{String: "Excellence, Innovation, Trust", Valid: true},
		MissionIcon:   "shield",
		VisionIcon:    "visibility",
		ValuesIcon:    "star",
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/about/mvv", strings.NewReader(url.Values{
		"mission":        {"To deliver comprehensive security solutions"},
		"vision":         {"Leading the industry in innovation"},
		"values_summary": {"Integrity, Quality, Customer Focus"},
		"mission_icon":   {"security"},
		"vision_icon":    {"lightbulb"},
		"values_icon":    {"favorite"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	mvv, _ := queries.GetMissionVisionValues(ctx)
	if mvv.Mission != "To deliver comprehensive security solutions" {
		t.Errorf("expected updated mission, got %q", mvv.Mission)
	}
	if mvv.Vision != "Leading the industry in innovation" {
		t.Errorf("expected updated vision, got %q", mvv.Vision)
	}
}

func TestAboutMVVLoad_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	queries.UpsertMissionVisionValues(ctx, sqlc.UpsertMissionVisionValuesParams{
		Mission:     "Test Mission",
		Vision:      "Test Vision",
		MissionIcon: "test_icon",
		VisionIcon:  "test_icon2",
		ValuesIcon:  "test_icon3",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/about/mvv", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("about mvv route not found")
	}
}
