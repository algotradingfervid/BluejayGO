package e2e_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestWhitepapersCRUD_E2E(t *testing.T) {
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

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "IoT Best Practices")
	writer.WriteField("slug", "iot-best-practices")
	writer.WriteField("description", "Guide to IoT security")
	writer.WriteField("topic_id", fmt.Sprintf("%d", topic.ID))
	writer.WriteField("published_date", "2024-01-15")
	writer.WriteField("cover_color_from", "#3B82F6")
	writer.WriteField("cover_color_to", "#8B5CF6")
	writer.WriteField("meta_description", "IoT security guide")
	writer.WriteField("is_published", "on")
	writer.WriteField("page_count", "24")
	part, _ := writer.CreateFormFile("pdf_file", "test.pdf")
	io.WriteString(part, "fake pdf data")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/whitepapers", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	whitepapers, _ := queries.ListAllWhitepapers(ctx)
	if len(whitepapers) != 1 {
		t.Fatalf("expected 1 whitepaper, got %d", len(whitepapers))
	}
	if whitepapers[0].Title != "IoT Best Practices" {
		t.Errorf("expected 'IoT Best Practices', got %q", whitepapers[0].Title)
	}
	if whitepapers[0].Slug != "iot-best-practices" {
		t.Errorf("expected slug 'iot-best-practices', got %q", whitepapers[0].Slug)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/whitepapers/%d", whitepapers[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	_, err := queries.GetWhitepaperByID(ctx, whitepapers[0].ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestWhitepapersList_E2E(t *testing.T) {
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

	for i := 1; i <= 3; i++ {
		queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
			Title:          fmt.Sprintf("Whitepaper %d", i),
			Slug:           fmt.Sprintf("whitepaper-%d", i),
			Description:    "Description",
			TopicID:        topic.ID,
			IsPublished:    int64(i % 2),
			PublishedDate:  "2024-01-15",
			CoverColorFrom: "#3B82F6",
			CoverColorTo:   "#8B5CF6",
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/whitepapers", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepapers list route not found")
	}
}

func TestWhitepaperEdit_E2E(t *testing.T) {
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
		Title:          "Original Title",
		Slug:           "original-title",
		Description:    "Original description",
		TopicID:        topic.ID,
		IsPublished:    1,
		PublishedDate:  "2024-01-15",
		CoverColorFrom: "#3B82F6",
		CoverColorTo:   "#8B5CF6",
	})

	updateBody := &bytes.Buffer{}
	updateWriter := multipart.NewWriter(updateBody)
	updateWriter.WriteField("title", "Updated Title")
	updateWriter.WriteField("slug", "updated-title")
	updateWriter.WriteField("description", "Updated description")
	updateWriter.WriteField("topic_id", fmt.Sprintf("%d", topic.ID))
	updateWriter.WriteField("published_date", "2024-02-20")
	updateWriter.WriteField("cover_color_from", "#FF0000")
	updateWriter.WriteField("cover_color_to", "#00FF00")
	updateWriter.WriteField("is_published", "on")
	updateWriter.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/whitepapers/%d", wp.ID), updateBody)
	req.Header.Set("Content-Type", updateWriter.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetWhitepaperByID(ctx, wp.ID)
	if updated.Title != "Updated Title" {
		t.Errorf("expected 'Updated Title', got %q", updated.Title)
	}
}
