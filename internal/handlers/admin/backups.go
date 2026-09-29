package admin

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/internal/middleware"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

var backupName = regexp.MustCompile(`^bluejay-backup-[0-9]{8}T[0-9]{6}Z\.tar\.gz$`)

type BackupsHandler struct{ directory string }

func NewBackupsHandler(directory string) *BackupsHandler {
	return &BackupsHandler{directory: directory}
}

type backupFile struct{ Name, Created, Size string }
type backupStatus struct {
	Success  bool   `json:"success"`
	Finished string `json:"finished"`
}

func backupAccess(c echo.Context) error {
	sess, ok := c.Get("session").(*middleware.Session)
	if !ok || sess.UserID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	if sess.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	c.Response().Header().Set("Cache-Control", "no-store")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	return nil
}

func (h *BackupsHandler) List(c echo.Context) error {
	if err := backupAccess(c); err != nil {
		return err
	}
	files := []backupFile{}
	healthy := false
	var status backupStatus
	if h.directory != "" {
		entries, err := os.ReadDir(h.directory)
		if err != nil {
			return echo.NewHTTPError(http.StatusServiceUnavailable, "Backups are temporarily unavailable")
		}
		for _, entry := range entries {
			if !backupName.MatchString(entry.Name()) || entry.Type()&os.ModeSymlink != 0 {
				continue
			}
			info, err := entry.Info()
			if err != nil || !info.Mode().IsRegular() {
				continue
			}
			files = append(files, backupFile{entry.Name(), info.ModTime().UTC().Format("02 Jan 2006, 15:04 UTC"), fmt.Sprintf("%.1f MB", float64(info.Size())/(1024*1024))})
		}
		sort.Slice(files, func(i, j int) bool { return files[i].Name > files[j].Name })
		if b, err := os.ReadFile(filepath.Join(h.directory, "status.json")); err == nil && len(b) < 4096 {
			if json.Unmarshal(b, &status) == nil {
				finished, err := time.Parse(time.RFC3339, status.Finished)
				healthy = err == nil && status.Success && time.Since(finished) < 36*time.Hour && len(files) > 0
			}
		}
	}
	return c.Render(http.StatusOK, "admin/pages/backups.html", map[string]interface{}{"Title": "Backups", "Files": files, "Configured": h.directory != "", "Healthy": healthy, "LastRun": status.Finished})
}

func (h *BackupsHandler) Download(c echo.Context) error {
	if err := backupAccess(c); err != nil {
		return err
	}
	name := c.Param("name")
	if h.directory == "" || !backupName.MatchString(name) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	path := filepath.Join(h.directory, name)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	f, err := os.Open(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	c.Response().Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	c.Response().Header().Set("Content-Type", "application/gzip")
	// Backups can take longer than ordinary pages to transfer, but remain bounded.
	_ = http.NewResponseController(c.Response().Writer).SetWriteDeadline(time.Now().Add(30 * time.Minute))
	logActivity(c, "downloaded", "backup", 0, name, "Downloaded backup '%s'", name)
	http.ServeContent(c.Response(), c.Request(), name, opened.ModTime(), f)
	return nil
}
