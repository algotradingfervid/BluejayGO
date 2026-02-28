package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestWhitepaperDownloadsList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name:        "Security",
		Slug:        "security",
		Description: sql.NullString{String: "Security topics", Valid: true},
		SortOrder:   1,
	})

	wp, _ := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Test Whitepaper",
		Slug:           "test-whitepaper",
		Description:    "Test description",
		TopicID:        topic.ID,
		IsPublished:    1,
		PublishedDate:  "2024-01-15",
		CoverColorFrom: "#3B82F6",
		CoverColorTo:   "#8B5CF6",
	})

	for i := 1; i <= 5; i++ {
		queries.CreateWhitepaperDownload(ctx, sqlc.CreateWhitepaperDownloadParams{
			WhitepaperID:     wp.ID,
			Name:             "Test User",
			Email:            "test@example.com",
			Company:          "Test Corp",
			Designation:      sql.NullString{String: "Engineer", Valid: true},
			MarketingConsent: int64(i % 2),
			IpAddress:        sql.NullString{},
			UserAgent:        sql.NullString{},
		})
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/whitepapers/%d/downloads", wp.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepaper downloads route not found")
	}
}

func TestWhitepaperDownloadsFilter_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name:        "Security",
		Slug:        "security",
		Description: sql.NullString{String: "Security topics", Valid: true},
		SortOrder:   1,
	})

	wp, _ := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Test Whitepaper",
		Slug:           "test-whitepaper",
		Description:    "Test description",
		TopicID:        topic.ID,
		IsPublished:    1,
		PublishedDate:  "2024-01-15",
		CoverColorFrom: "#3B82F6",
		CoverColorTo:   "#8B5CF6",
	})

	queries.CreateWhitepaperDownload(ctx, sqlc.CreateWhitepaperDownloadParams{
		WhitepaperID:     wp.ID,
		Name:             "Test User",
		Email:            "test@example.com",
		Company:          "Test Corp",
		Designation:      sql.NullString{String: "Engineer", Valid: true},
		MarketingConsent: 1,
		IpAddress:        sql.NullString{},
		UserAgent:        sql.NullString{},
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/whitepapers/%d/downloads?date_from=2024-01-01&date_to=2024-12-31", wp.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepaper downloads with date filter route not found")
	}
}
