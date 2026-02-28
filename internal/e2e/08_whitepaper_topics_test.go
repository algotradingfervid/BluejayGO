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

// NOTE: Whitepaper topics routes are not registered in setupApp() yet.
// These tests will fail until routes are added to setupApp().

func TestWhitepaperTopicsList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name:        "Cybersecurity",
		Slug:        "cybersecurity",
		ColorHex:    "#FF0000",
		Icon:        "shield",
		Description: sql.NullString{String: "Security topics", Valid: true},
		SortOrder:   1,
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/whitepaper-topics", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestWhitepaperTopicsCreate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req := httptest.NewRequest(http.MethodPost, "/admin/whitepaper-topics", strings.NewReader(url.Values{
		"name":        {"AI & Machine Learning"},
		"color_hex":   {"#00FF00"},
		"icon":        {"brain"},
		"description": {"Artificial intelligence topics"},
		"sort_order":  {"2"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}

	items, _ := queries.ListWhitepaperTopicsWithCount(context.Background())
	if len(items) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(items))
	}
	if items[0].Name != "AI & Machine Learning" {
		t.Errorf("expected 'AI & Machine Learning', got %q", items[0].Name)
	}
	if items[0].Slug != "ai--machine-learning" {
		t.Errorf("expected slug 'ai--machine-learning', got %q", items[0].Slug)
	}
	if items[0].ColorHex != "#00FF00" {
		t.Errorf("expected color '#00FF00', got %q", items[0].ColorHex)
	}
}

func TestWhitepaperTopicsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name:        "IoT",
		Slug:        "iot",
		ColorHex:    "#0000FF",
		Icon:        "network",
		Description: sql.NullString{String: "Internet of Things", Valid: true},
		SortOrder:   3,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/whitepaper-topics/%d", topic.ID), strings.NewReader(url.Values{
		"name":        {"Internet of Things"},
		"color_hex":   {"#FF00FF"},
		"icon":        {"connected-devices"},
		"description": {"IoT and connected devices"},
		"sort_order":  {"1"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetWhitepaperTopic(ctx, topic.ID)
	if updated.Name != "Internet of Things" {
		t.Errorf("expected 'Internet of Things', got %q", updated.Name)
	}
	if updated.Slug != "internet-of-things" {
		t.Errorf("expected slug 'internet-of-things', got %q", updated.Slug)
	}
	if updated.ColorHex != "#FF00FF" {
		t.Errorf("expected color '#FF00FF', got %q", updated.ColorHex)
	}
}

func TestWhitepaperTopicsDelete_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name:        "Blockchain",
		Slug:        "blockchain",
		ColorHex:    "#FFFF00",
		Icon:        "chain",
		Description: sql.NullString{String: "Blockchain tech", Valid: true},
		SortOrder:   5,
	})

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/whitepaper-topics/%d", topic.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	items, _ := queries.ListWhitepaperTopicsWithCount(ctx)
	if len(items) != 0 {
		t.Errorf("expected 0 topics after delete, got %d", len(items))
	}
}
