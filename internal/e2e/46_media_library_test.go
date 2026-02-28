package e2e_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestMediaLibraryList(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("load media library", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/media", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("search media files", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/media?search=test", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("sort media files", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/media?sort=name", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}

func TestMediaUpload(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("upload single file", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("files", "test.jpg")
		io.WriteString(part, "fake image data")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/admin/media/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK && rec.Code != http.StatusBadRequest {
			t.Errorf("expected 200 or 400, got %d", rec.Code)
		}
	})
}

func TestMediaGetFile(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	file, _ := queries.CreateMediaFile(ctx, sqlc.CreateMediaFileParams{
		Filename:         "test.jpg",
		OriginalFilename: "test.jpg",
		FilePath:         "/uploads/media/test.jpg",
		FileSize:         1024,
		MimeType:         "image/jpeg",
		Width:            sql.NullInt64{Int64: 800, Valid: true},
		Height:           sql.NullInt64{Int64: 600, Valid: true},
		AltText:          sql.NullString{String: "", Valid: true},
	})

	t.Run("get file by id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/media/%d", file.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		var response sqlc.MediaFile
		json.NewDecoder(rec.Body).Decode(&response)
		if response.ID != file.ID {
			t.Error("file ID mismatch")
		}
	})
}

func TestMediaUpdateAltText(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	file, _ := queries.CreateMediaFile(ctx, sqlc.CreateMediaFileParams{
		Filename:         "alt-test.jpg",
		OriginalFilename: "alt-test.jpg",
		FilePath:         "/uploads/media/alt-test.jpg",
		FileSize:         1024,
		MimeType:         "image/jpeg",
		AltText:          sql.NullString{String: "", Valid: true},
	})

	t.Run("update alt text", func(t *testing.T) {
		body := strings.NewReader(`{"alt_text":"Updated alt text"}`)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/admin/media/%d", file.ID), body)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		updated, _ := queries.GetMediaFile(ctx, file.ID)
		if updated.AltText.String != "Updated alt text" {
			t.Errorf("expected 'Updated alt text', got %q", updated.AltText.String)
		}
	})
}

func TestMediaDelete(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	ctx := context.Background()

	file, _ := queries.CreateMediaFile(ctx, sqlc.CreateMediaFileParams{
		Filename:         "delete-test.jpg",
		OriginalFilename: "delete-test.jpg",
		FilePath:         "/uploads/media/delete-test.jpg",
		FileSize:         1024,
		MimeType:         "image/jpeg",
		AltText:          sql.NullString{String: "", Valid: true},
	})

	t.Run("delete file", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/media/%d", file.ID), nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}

		_, err := queries.GetMediaFile(ctx, file.ID)
		if err != sql.ErrNoRows {
			t.Error("expected file to be deleted")
		}
	})
}

func TestMediaBrowser(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)

	t.Run("load media browser", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/media/browse", nil)
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}
