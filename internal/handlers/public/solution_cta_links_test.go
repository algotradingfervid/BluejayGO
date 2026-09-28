package public

import (
	"bytes"
	"database/sql"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestSolutionBrochureAvailabilityAndVisibility(t *testing.T) {
	publicDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(publicDir, "uploads"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(publicDir, "uploads", "brochure.pdf"), []byte("%PDF-1.4 test"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		url  string
		want bool
	}{
		{"", false}, {"#", false}, {"/downloads/corporate.pdf", false}, {"/uploads/missing.pdf", false},
		{"/uploads/brochure.pdf", true}, {"/public/uploads/brochure.pdf", true},
		{"https://example.com/brochure.pdf", true}, {"javascript:alert(1)", false}, {"//example.com/brochure.pdf", false},
		{"/public/../../secret", false},
	} {
		if got := solutionButtonAvailable("Download Brochure", tc.url, publicDir); got != tc.want {
			t.Errorf("%q: got %v want %v", tc.url, got, tc.want)
		}
	}
	str := func(v string) sql.NullString { return sql.NullString{String: v, Valid: true} }
	cta := sqlc.SolutionCta{Heading: "Corporate CTA", SectionName: "main_cta", IsActive: true, SecondaryButtonEnabled: true, PrimaryButtonText: str("Contact Us"), PrimaryButtonUrl: str("/contact"), SecondaryButtonText: str("Download Brochure"), SecondaryButtonUrl: str("/uploads/brochure.pdf")}
	render := func(ctas []sqlc.SolutionCta) string {
		t.Helper()
		tmpl, err := template.New("solution").Funcs(template.FuncMap{"safeHTML": func(v string) template.HTML { return template.HTML(v) }}).ParseFiles("../../../templates/public/pages/solution_detail.html")
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		data := map[string]interface{}{"Solution": sqlc.Solution{}, "Sections": map[string]interface{}{}, "CTAs": visibleSolutionCTAs(ctas, publicDir)}
		if err := tmpl.ExecuteTemplate(&buf, "content", data); err != nil {
			t.Fatal(err)
		}
		return buf.String()
	}
	if html := render([]sqlc.SolutionCta{cta}); !strings.Contains(html, `href="/uploads/brochure.pdf"`) {
		t.Fatal("configured brochure missing")
	}
	cta.SecondaryButtonEnabled = false
	if html := render([]sqlc.SolutionCta{cta}); strings.Contains(html, "Download Brochure") || !strings.Contains(html, `href="/contact"`) {
		t.Fatal("button toggle did not preserve main action")
	}
	cta.SecondaryButtonEnabled = true
	cta.SecondaryButtonUrl = str("#")
	if html := render([]sqlc.SolutionCta{cta}); strings.Contains(html, "Download Brochure") {
		t.Fatal("placeholder brochure rendered")
	}
	cta.IsActive = false
	if html := render([]sqlc.SolutionCta{cta}); strings.Contains(html, "Corporate CTA") {
		t.Fatal("disabled section rendered")
	}
}
