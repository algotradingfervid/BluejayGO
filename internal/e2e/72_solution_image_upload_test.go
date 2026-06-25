package e2e_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestSolutionImageUpload_E2E verifies that uploading a hero image file when
// creating a solution persists the file to disk and stores its public
// /uploads/solutions/ path as the solution's hero image, overriding any pasted URL.
func TestSolutionImageUpload_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// Build a multipart form with an uploaded hero image plus a (different) pasted
	// URL; the uploaded file must win.
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fields := map[string]string{
		"title":          "Aviation Solutions",
		"slug":           "aviation-solutions",
		"hero_image_url": "/should/be/overridden.png",
		"is_published":   "1",
	}
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	part, err := w.CreateFormFile("hero_image_file", "hero.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("\x89PNG\r\n\x1a\nfake-png-bytes")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/solutions", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d (body: %s)", rec.Code, rec.Body.String())
	}

	sol, err := queries.GetSolutionBySlugIncludeDrafts(ctx, "aviation-solutions")
	if err != nil {
		t.Fatalf("failed to fetch created solution: %v", err)
	}
	if !sol.HeroImageUrl.Valid || !strings.HasPrefix(sol.HeroImageUrl.String, "/uploads/solutions/") {
		t.Errorf("expected uploaded hero image under /uploads/solutions/, got %q", sol.HeroImageUrl.String)
	}
	if !strings.HasSuffix(sol.HeroImageUrl.String, "_hero.png") {
		t.Errorf("expected saved path to keep original filename, got %q", sol.HeroImageUrl.String)
	}
}

// TestSolutionImageUpload_RejectsBadType_E2E verifies an unsupported upload type
// is rejected with a 400 rather than being silently saved.
func TestSolutionImageUpload_RejectsBadType_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	_ = w.WriteField("title", "Bad Upload")
	_ = w.WriteField("slug", "bad-upload")
	part, _ := w.CreateFormFile("hero_image_file", "evil.avif")
	_, _ = part.Write([]byte("not-an-image"))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/solutions", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unsupported file type, got %d", rec.Code)
	}
}
