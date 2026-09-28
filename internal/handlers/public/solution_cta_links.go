package public

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

// solutionButtonAvailable rejects placeholder links. Local downloads must point
// to an existing file under a directory actually served by the application.
func solutionButtonAvailable(label, rawURL, publicDir string) bool {
	label, rawURL = strings.TrimSpace(label), strings.TrimSpace(rawURL)
	if label == "" || rawURL == "" || strings.HasPrefix(rawURL, "#") {
		return false
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.User != nil {
		return false
	}
	if u.Scheme != "" {
		return (u.Scheme == "https" || u.Scheme == "http") && u.Host != ""
	}
	if u.Host != "" || !strings.HasPrefix(u.Path, "/") || strings.Contains(u.Path, "\\") {
		return false
	}
	isDownload := strings.Contains(strings.ToLower(label), "brochure") || strings.Contains(strings.ToLower(label), "download")
	if !isDownload {
		return true
	}
	var relative string
	switch {
	case strings.HasPrefix(u.Path, "/uploads/"):
		relative = strings.TrimPrefix(u.Path, "/")
	case strings.HasPrefix(u.Path, "/public/"):
		relative = strings.TrimPrefix(u.Path, "/public/")
	default:
		return false
	}
	if !filepath.IsLocal(relative) {
		return false
	}
	info, err := os.Stat(filepath.Join(publicDir, filepath.FromSlash(relative)))
	return err == nil && info.Mode().IsRegular()
}

func visibleSolutionCTAs(ctas []sqlc.SolutionCta, publicDir string) []sqlc.SolutionCta {
	visible := make([]sqlc.SolutionCta, 0, len(ctas))
	for _, cta := range ctas {
		if !cta.IsActive {
			continue
		}
		cta.PrimaryButtonText.Valid = solutionButtonAvailable(cta.PrimaryButtonText.String, cta.PrimaryButtonUrl.String, publicDir)
		cta.SecondaryButtonText.Valid = cta.SecondaryButtonEnabled && solutionButtonAvailable(cta.SecondaryButtonText.String, cta.SecondaryButtonUrl.String, publicDir)
		visible = append(visible, cta)
	}
	return visible
}
