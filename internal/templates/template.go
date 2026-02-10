package templates

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Renderer struct {
	templates map[string]*template.Template
	basePath  string
}

func NewRenderer(basePath string) *Renderer {
	r := &Renderer{
		templates: make(map[string]*template.Template),
		basePath:  basePath,
	}
	r.loadTemplates()
	return r
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

func (r *Renderer) loadTemplates() {
	funcMap := template.FuncMap{
		"safeHTML":   safeHTML,
		"formatDate": formatDate,
		"truncate":   truncate,
		"slugify":    slugify,
		"now":        time.Now,
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"upper":      strings.ToUpper,
		"seq": func(n int64) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"int64": func(i int) int64 { return int64(i) },
		"formatFileSize": func(size int64) string {
			if size < 1024 {
				return fmt.Sprintf("%d B", size)
			}
			if size < 1024*1024 {
				return fmt.Sprintf("%.1f KB", float64(size)/1024)
			}
			return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
		},
	}

	r.templates["public/pages/home.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/layouts/base.html"),
		filepath.Join(r.basePath, "public/pages/home.html"),
		filepath.Join(r.basePath, "partials/header.html"),
		filepath.Join(r.basePath, "partials/footer.html"),
	))

	r.templates["admin/pages/login.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/layouts/base.html"),
		filepath.Join(r.basePath, "admin/pages/login.html"),
	))

	r.templates["admin/pages/dashboard.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/layouts/base.html"),
		filepath.Join(r.basePath, "admin/pages/dashboard.html"),
		filepath.Join(r.basePath, "partials/admin-sidebar.html"),
	))

	// Phase 3: Public product pages
	publicProductPages := []string{
		"products", "products_category", "product_detail",
	}
	for _, page := range publicProductPages {
		r.templates["public/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "public/layouts/base.html"),
			filepath.Join(r.basePath, "public/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/header.html"),
			filepath.Join(r.basePath, "partials/footer.html"),
		))
	}

	// Phase 2: Master table admin pages
	masterPages := []string{
		"product_categories_list", "product_categories_form",
		"blog_categories_list", "blog_categories_form",
		"blog_authors_list", "blog_authors_form",
		"industries_list", "industries_form",
		"partner_tiers_list", "partner_tiers_form",
		"whitepaper_topics_list", "whitepaper_topics_form",
		"products_list", "products_form",
		"settings_form",
		"page_sections_list", "page_sections_form",
		"header_form",
		"footer_form",
	}
	for _, page := range masterPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}

	// Phase 4: Public solution pages
	publicSolutionPages := []string{
		"solutions_list", "solution_detail",
	}
	for _, page := range publicSolutionPages {
		r.templates["public/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "public/layouts/base.html"),
			filepath.Join(r.basePath, "public/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/header.html"),
			filepath.Join(r.basePath, "partials/footer.html"),
		))
	}

	// Phase 4: Admin solution pages
	solutionAdminPages := []string{
		"solutions_list", "solutions_form",
	}
	for _, page := range solutionAdminPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}

	// Phase 4: Admin solution partials (HTMX fragments - standalone, no layout)
	solutionPartials := []string{
		"solution_stats", "solution_challenges", "solution_products", "solution_ctas",
	}
	for _, partial := range solutionPartials {
		r.templates["admin/partials/"+partial+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/partials/"+partial+".html"),
		))
	}

	// Phase 6: Public case study pages
	publicCaseStudyPages := []string{
		"case_studies", "case_study_detail",
	}
	for _, page := range publicCaseStudyPages {
		r.templates["public/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "public/layouts/base.html"),
			filepath.Join(r.basePath, "public/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/header.html"),
			filepath.Join(r.basePath, "partials/footer.html"),
		))
	}

	// Phase 6: Admin case study pages
	caseStudyAdminPages := []string{
		"case_studies_list", "case_studies_form",
	}
	for _, page := range caseStudyAdminPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}

	// Phase 6: Admin case study partials (HTMX fragments - standalone, no layout)
	caseStudyPartials := []string{
		"case_study_products", "case_study_metrics",
	}
	for _, partial := range caseStudyPartials {
		r.templates["admin/partials/"+partial+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/partials/"+partial+".html"),
		))
	}

	// Phase 5: Blog tag partials (HTMX fragments - standalone, no layout)
	blogPartials := []string{"tag_suggestions", "tag_chip", "product_suggestions"}
	for _, partial := range blogPartials {
		r.templates["admin/partials/"+partial+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/partials/"+partial+".html"),
		))
	}

	// Product search partial (standalone, no layout)
	r.templates["public/partials/product_search_results.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/partials/product_search_results.html"),
	))

	// Phase 5: Public blog pages
	publicBlogPages := []string{
		"blog_listing", "blog_post",
	}
	for _, page := range publicBlogPages {
		r.templates["public/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "public/layouts/base.html"),
			filepath.Join(r.basePath, "public/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/header.html"),
			filepath.Join(r.basePath, "partials/footer.html"),
		))
	}

	// Phase 5: Admin blog pages
	blogAdminPages := []string{
		"blog_posts_list", "blog_post_form", "blog_tags_list",
	}
	for _, page := range blogAdminPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}
	// Phase 8: Public whitepaper pages
	publicWhitepaperPages := []string{
		"whitepapers", "whitepaper_detail",
	}
	for _, page := range publicWhitepaperPages {
		r.templates["public/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "public/layouts/base.html"),
			filepath.Join(r.basePath, "public/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/header.html"),
			filepath.Join(r.basePath, "partials/footer.html"),
		))
	}

	// Phase 8: Whitepaper success partial (HTMX fragment - standalone, no layout)
	r.templates["public/pages/whitepaper_success.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/pages/whitepaper_success.html"),
	))

	// Phase 8: Public contact page
	r.templates["public/pages/contact.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/layouts/base.html"),
		filepath.Join(r.basePath, "public/pages/contact.html"),
		filepath.Join(r.basePath, "partials/header.html"),
		filepath.Join(r.basePath, "partials/footer.html"),
	))

	// Phase 8: Admin whitepaper pages
	whitepaperAdminPages := []string{
		"whitepapers_list", "whitepapers_form", "whitepapers_downloads",
	}
	for _, page := range whitepaperAdminPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}

	// Phase 8: Admin contact pages
	contactAdminPages := []string{
		"contact_submissions_list", "contact_submission_detail",
		"office_locations_list", "office_locations_form",
	}
	for _, page := range contactAdminPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}

	// Phase 7: Public about and partners pages
	publicAboutPages := []string{"about", "partners"}
	for _, page := range publicAboutPages {
		r.templates["public/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "public/layouts/base.html"),
			filepath.Join(r.basePath, "public/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/header.html"),
			filepath.Join(r.basePath, "partials/footer.html"),
		))
	}

	// Phase 9: Admin homepage pages
	homepageAdminPages := []string{
		"homepage_heroes_list", "homepage_hero_form",
		"homepage_stats_list", "homepage_stat_form",
		"homepage_testimonials_list", "homepage_testimonial_form",
		"homepage_cta_list", "homepage_cta_form",
		"homepage_settings",
		"products_settings",
		"solutions_settings",
		"blog_settings",
	}
	for _, page := range homepageAdminPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}

	// Phase 9: Search page
	r.templates["public/pages/search.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/layouts/base.html"),
		filepath.Join(r.basePath, "public/pages/search.html"),
		filepath.Join(r.basePath, "partials/header.html"),
		filepath.Join(r.basePath, "partials/footer.html"),
	))

	// Phase 9: Search suggestions partial (HTMX fragment)
	r.templates["public/partials/search_suggestions.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/partials/search_suggestions.html"),
	))

	// Phase 7: Admin about and partners pages
	aboutAdminPages := []string{
		"about_overview_form", "about_mvv_form",
		"core_values_list", "core_values_form",
		"milestones_list", "milestones_form",
		"certifications_list", "certifications_form",
		"partners_list", "partners_form",
		"testimonials_list", "testimonials_form",
		"about_settings",
	}
	for _, page := range aboutAdminPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}

	// Phase 18: Media library page
	r.templates["admin/pages/media_library.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/layouts/base.html"),
		filepath.Join(r.basePath, "admin/pages/media_library.html"),
		filepath.Join(r.basePath, "partials/admin-sidebar.html"),
	))

	// Phase 18: Media picker partial (HTMX fragment - standalone, no layout)
	r.templates["admin/partials/media_picker.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/partials/media_picker.html"),
	))

	// Phase 19: Navigation editor pages
	navigationPages := []string{
		"navigation_list", "navigation_editor",
	}
	for _, page := range navigationPages {
		r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/layouts/base.html"),
			filepath.Join(r.basePath, "admin/pages/"+page+".html"),
			filepath.Join(r.basePath, "partials/admin-sidebar.html"),
		))
	}
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func formatDate(t time.Time, format string) string {
	if format == "" {
		format = "January 2, 2006"
	}
	return t.Format(format)
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
