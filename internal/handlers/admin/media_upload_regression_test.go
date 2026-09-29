package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

type uploadFixture struct {
	name string
	size int64
}
type uploadOutcome struct {
	Count int `json:"count"`
	Files []struct {
		ID       int64  `json:"id"`
		FilePath string `json:"file_path"`
	} `json:"files"`
	Results []struct {
		Index    int    `json:"index"`
		Filename string `json:"filename"`
		Status   string `json:"status"`
		Error    string `json:"error"`
	} `json:"results"`
}

func mediaUploadRequest(t *testing.T, h *admin.MediaHandler, files []uploadFixture) (*httptest.ResponseRecorder, uploadOutcome) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for _, file := range files {
		part, err := writer.CreateFormFile("files", file.name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.CopyN(part, bytes.NewReader(make([]byte, file.size)), file.size); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/admin/media/upload", &body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	if err := h.Upload(echo.New().NewContext(req, rec)); err != nil {
		t.Fatal(err)
	}
	var outcome uploadOutcome
	if err := json.Unmarshal(rec.Body.Bytes(), &outcome); err != nil {
		t.Fatal(err)
	}
	return rec, outcome
}

func TestMediaUploadReportsPerFileOutcomes(t *testing.T) {
	for _, tc := range []struct {
		name      string
		files     []uploadFixture
		wantSaved int
		wantCode  int
	}{
		{"unsupported", []uploadFixture{{"unsupported.txt", 30}}, 0, http.StatusBadRequest},
		{"oversized", []uploadFixture{{"oversized.png", 10*1024*1024 + 1}}, 0, http.StatusBadRequest},
		{"mixed", []uploadFixture{{"same.pdf", 30}, {"unsupported.txt", 30}, {"same.pdf", 40}, {"oversized.png", 10*1024*1024 + 1}}, 2, http.StatusOK},
		{"valid", []uploadFixture{{"valid.pdf", 30}}, 1, http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, q, cleanup := testutil.SetupTestDB(t)
			defer cleanup()
			dir := t.TempDir()
			h := admin.NewMediaHandler(q, logger, dir)
			rec, out := mediaUploadRequest(t, h, tc.files)
			if rec.Code != tc.wantCode || out.Count != tc.wantSaved || len(out.Files) != tc.wantSaved || len(out.Results) != len(tc.files) {
				t.Fatalf("status=%d outcome=%+v", rec.Code, out)
			}
			saved := 0
			for i, result := range out.Results {
				if result.Index != i || result.Filename != tc.files[i].name {
					t.Fatalf("lost file identity at%d: %+v", i, result)
				}
				if result.Status == "uploaded" {
					saved++
					if result.Error != "" {
						t.Fatal("saved item has error")
					}
				} else if result.Status != "error" || result.Error == "" {
					t.Fatalf("missing failure detail: %+v", result)
				}
			}
			count, err := q.CountMediaFiles(context.Background())
			if err != nil || int(count) != tc.wantSaved || saved != tc.wantSaved {
				t.Fatalf("stored count=%d err=%v", count, err)
			}
			entries, err := os.ReadDir(filepath.Join(dir, "media"))
			if err != nil || len(entries) != tc.wantSaved {
				t.Fatalf("disk files=%d err=%v", len(entries), err)
			}
			for _, file := range out.Files {
				if file.ID == 0 || file.FilePath == "" {
					t.Fatal("successful response missing persisted file identity")
				}
			}
		})
	}
}

func TestMediaUploadDatabaseFailureRemovesOrphanFile(t *testing.T) {
	db, q, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	dir := t.TempDir()
	h := admin.NewMediaHandler(q, logger, dir)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	rec, out := mediaUploadRequest(t, h, []uploadFixture{{"orphan.pdf", 30}})
	if rec.Code != http.StatusBadRequest || out.Count != 0 || len(out.Results) != 1 || out.Results[0].Error == "" {
		t.Fatalf("unexpected failed upload: %d %+v", rec.Code, out)
	}
	entries, err := os.ReadDir(filepath.Join(dir, "media"))
	if err != nil || len(entries) != 0 {
		t.Fatalf("orphan remains: entries=%v err=%v", entries, err)
	}
}
