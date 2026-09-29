package admin_test

import (
	"github.com/labstack/echo/v4"
	a "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	m "github.com/narendhupati/bluejay-cms/internal/middleware"
	templates "github.com/narendhupati/bluejay-cms/internal/templates"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBackupAccessAndDownload(t *testing.T) {
	dir := t.TempDir()
	name := "bluejay-backup-20260929T023000Z.tar.gz"
	os.WriteFile(filepath.Join(dir, name), []byte("synthetic-archive"), 0600)
	h := a.NewBackupsHandler(dir)
	for _, tc := range []struct {
		role string
		id   int64
		want int
	}{{"", 0, 401}, {"editor", 1, 403}, {"admin", 1, 200}} {
		e := echo.New()
		e.GET("/:name", func(c echo.Context) error {
			c.Set("session", &m.Session{UserID: tc.id, Role: tc.role})
			return h.Download(c)
		})
		w := httptest.NewRecorder()
		e.ServeHTTP(w, httptest.NewRequest("GET", "/"+name, nil))
		if w.Code != tc.want {
			t.Fatalf("role%q got%d", tc.role, w.Code)
		}
		if w.Code == 200 {
			if w.Body.String() != "synthetic-archive" || w.Header().Get("Cache-Control") != "no-store" || !strings.HasPrefix(w.Header().Get("Content-Disposition"), "attachment;") {
				t.Fatal("unsafe download response")
			}
		}
	}
	// A correctly-shaped symlink is also forbidden.
	linked := "bluejay-backup-20260928T023000Z.tar.gz"
	os.Symlink(filepath.Join(dir, name), filepath.Join(dir, linked))
	for _, file := range []string{linked, "../secret", "status.json"} {
		e := echo.New()
		r := httptest.NewRequest("GET", "/", nil)
		c := e.NewContext(r, httptest.NewRecorder())
		c.SetParamNames("name")
		c.SetParamValues(file)
		c.Set("session", &m.Session{UserID: 1, Role: "admin"})
		if err := h.Download(c); err == nil {
			t.Fatal("unsafe file served:", file)
		}
	}
}
func TestBackupPageRenders(t *testing.T) {
	dir := t.TempDir()
	h := a.NewBackupsHandler(dir)
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	e.GET("/admin/backups", func(c echo.Context) error { c.Set("session", &m.Session{UserID: 1, Role: "admin"}); return h.List(c) })
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest("GET", "/admin/backups", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "No completed backups") {
		t.Fatalf("page failed: %d %s", w.Code, w.Body.String())
	}
}
