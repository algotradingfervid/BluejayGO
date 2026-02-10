// Package templates provides a custom template rendering engine for the Bluejay CMS.
// It integrates Go's html/template with Echo web framework and supports HTMX partial rendering.
//
// The package manages three types of templates:
// 1. Full page templates - Complete pages with layouts (admin/public)
// 2. HTMX partials - Standalone fragments for dynamic updates
// 3. Shared components - Headers, footers, sidebars used across pages
//
// Templates are pre-compiled at startup for performance and follow a strict
// directory structure under templates/ with admin/, public/, and partials/ subdirectories.
package templates

import (
	"fmt"        // For string formatting in error messages and file size display
	"html/template" // Go's HTML templating engine with auto-escaping for XSS protection
	"io"            // For writing rendered templates to HTTP response writers
	"path/filepath" // For cross-platform file path construction
	"strings"       // For string manipulation in template functions (ToUpper, ReplaceAll)
	"time"          // For date formatting in template functions

	"github.com/labstack/echo/v4" // Echo web framework - provides HTTP context for rendering
)

// Renderer implements Echo's echo.Renderer interface to integrate Go templates with Echo.
// It pre-compiles all templates at startup and caches them in memory for fast rendering.
//
// The templates map uses template paths as keys (e.g., "admin/pages/dashboard.html")
// and stores fully parsed template trees including all dependencies (layouts, partials).
//
// Integration with Echo:
// - Set as echo.Renderer via e.Renderer = NewRenderer(basePath)
// - Called automatically by c.Render(code, name, data) in handlers
// - Receives Echo context for potential context-aware rendering
//
// Integration with HTMX:
// - Full page templates render with layout (base.html wrapper)
// - HTMX partials render standalone (no layout) for hx-swap updates
// - Same Render method handles both by checking template structure
type Renderer struct {
	templates map[string]*template.Template // Map of template names to compiled template trees
	basePath  string                        // Root directory for template files (typically "templates/")
}

// NewRenderer creates and initializes a new template renderer.
// It immediately loads and compiles all templates from the basePath directory.
//
// Parameters:
//   - basePath: Root directory containing template files (e.g., "templates/")
//
// Returns:
//   - *Renderer: Fully initialized renderer with all templates pre-compiled
//
// The function will panic if template files are missing or contain syntax errors,
// ensuring that template problems are caught at startup rather than at request time.
//
// Example usage:
//   renderer := NewRenderer("templates/")
//   e.Renderer = renderer
func NewRenderer(basePath string) *Renderer {
	r := &Renderer{
		templates: make(map[string]*template.Template),
		basePath:  basePath,
	}
	r.loadTemplates()
	return r
}

// Render implements the echo.Renderer interface to render templates for HTTP responses.
// This method is called automatically by Echo when handlers use c.Render(code, name, data).
//
// Parameters:
//   - w: HTTP response writer where rendered HTML will be written
//   - name: Template identifier (e.g., "admin/pages/dashboard.html")
//   - data: Data passed to template for variable substitution (typically a struct or map)
//   - c: Echo context containing request/response info (currently unused but required by interface)
//
// Returns:
//   - error: Non-nil if template not found or execution fails
//
// Template execution behavior:
// - All templates execute their "base" block, which is defined in layout files
// - Full page templates: "base" includes complete HTML structure with <html>, <head>, <body>
// - HTMX partials: "base" is the fragment itself (no layout wrapper)
//
// Error handling:
// - Missing template returns error before rendering (fail-fast)
// - Template execution errors (missing data, type mismatches) return during rendering
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

// loadTemplates discovers and compiles all templates from the filesystem.
// This method is called once during initialization and will panic on any template errors,
// ensuring all templates are valid before the application starts serving requests.
//
// Template organization:
// - Full pages reference layouts (admin/layouts/base.html or public/layouts/base.html)
// - Layouts define {{block "content" .}} where page content is injected
// - Partials (header, footer, sidebar) are included in layouts or pages via {{template "name"}}
// - HTMX fragments are standalone files with no layout dependencies
//
// Template function map:
// All templates have access to custom functions for data formatting and manipulation.
// These functions are registered before template parsing and available in all templates.
func (r *Renderer) loadTemplates() {
	// funcMap registers custom functions available to all templates.
	// Functions provide data formatting, math operations, and string manipulation.
	// These extend Go's built-in template functions (len, printf, etc.)
	funcMap := template.FuncMap{
		"safeHTML":   safeHTML,   // Renders HTML without escaping (use carefully!)
		"formatDate": formatDate, // Formats time.Time to human-readable string
		"truncate":   truncate,   // Shortens strings with ellipsis
		"slugify":    slugify,    // Converts strings to URL-safe slugs
		"now":        time.Now,   // Returns current timestamp
		"add":        func(a, b int) int { return a + b }, // Integer addition for templates
		"sub":        func(a, b int) int { return a - b }, // Integer subtraction for templates
		"upper":      strings.ToUpper,                     // Converts string to uppercase
		// seq generates integer sequence for range loops ({{range seq 5}} generates 0,1,2,3,4)
		"seq": func(n int64) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"int64": func(i int) int64 { return int64(i) },     // Type conversion for int to int64
		// formatFileSize converts bytes to human-readable format (B, KB, MB)
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

	// Public homepage template
	// Uses: public/layouts/base.html (defines <html>, <head>, <body> structure)
	// Includes: partials/header.html (site navigation), partials/footer.html (site footer)
	// Content: public/pages/home.html defines {{block "content"}} for hero, stats, testimonials
	r.templates["public/pages/home.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/layouts/base.html"),
		filepath.Join(r.basePath, "public/pages/home.html"),
		filepath.Join(r.basePath, "partials/header.html"),
		filepath.Join(r.basePath, "partials/footer.html"),
	))

	// Admin login page template
	// Uses: admin/layouts/base.html (minimal layout, no sidebar)
	// Content: admin/pages/login.html defines login form with username/password fields
	// Note: No sidebar included - login page is accessed before authentication
	r.templates["admin/pages/login.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/layouts/base.html"),
		filepath.Join(r.basePath, "admin/pages/login.html"),
	))

	// Admin dashboard template
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation menu)
	// Content: admin/pages/dashboard.html shows statistics and recent activity
	r.templates["admin/pages/dashboard.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/layouts/base.html"),
		filepath.Join(r.basePath, "admin/pages/dashboard.html"),
		filepath.Join(r.basePath, "partials/admin-sidebar.html"),
	))

	// Phase 3: Public product pages
	// Uses: public/layouts/base.html for consistent site structure
	// Includes: partials/header.html (main navigation), partials/footer.html (site footer)
	// Templates:
	//   - products.html: Product catalog grid with filters and search
	//   - products_category.html: Category-filtered product listing
	//   - product_detail.html: Individual product page with specs, media, related products
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
	// Uses: admin/layouts/base.html (admin panel structure with nav)
	// Includes: partials/admin-sidebar.html (admin menu with active state highlighting)
	// Pattern: Each entity has _list (table view) and _form (create/edit) templates
	// Templates manage:
	//   - Categories: Product categories, blog categories, whitepaper topics
	//   - Authors: Blog authors with profiles
	//   - Industries: Industry classifications for products/solutions
	//   - Partner tiers: Partnership level definitions
	//   - Products: Full product CRUD with media, specs, categories
	//   - Settings: Global site settings (name, SEO, social links)
	//   - Page sections: Reusable content blocks for pages
	//   - Header/Footer: Site-wide navigation and footer content management
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
	// Uses: public/layouts/base.html for consistent public site structure
	// Includes: partials/header.html (navigation), partials/footer.html (footer)
	// Templates:
	//   - solutions_list.html: Industry solutions grid with filtering
	//   - solution_detail.html: Solution overview with stats, challenges, products, CTAs
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
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates:
	//   - solutions_list.html: Table of all solutions with edit/delete actions
	//   - solutions_form.html: Create/edit form with Trix editor and related content management
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
	// These are HTMX target fragments that render without any layout wrapper.
	// Used for dynamic updates within solutions_form.html via hx-get/hx-post.
	// Partials:
	//   - solution_stats.html: Editable statistics section (updates via HTMX)
	//   - solution_challenges.html: Industry challenges list (updates via HTMX)
	//   - solution_products.html: Related products selector (updates via HTMX)
	//   - solution_ctas.html: Call-to-action buttons editor (updates via HTMX)
	solutionPartials := []string{
		"solution_stats", "solution_challenges", "solution_products", "solution_ctas",
	}
	for _, partial := range solutionPartials {
		r.templates["admin/partials/"+partial+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/partials/"+partial+".html"),
		))
	}

	// Phase 6: Public case study pages
	// Uses: public/layouts/base.html for public site structure
	// Includes: partials/header.html (navigation), partials/footer.html (footer)
	// Templates:
	//   - case_studies.html: Grid of case studies with industry/product filtering
	//   - case_study_detail.html: Full case study with challenge, solution, results, metrics
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
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates:
	//   - case_studies_list.html: Table of case studies with status, industry, featured flag
	//   - case_studies_form.html: Create/edit form with client info, challenge, solution, results
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
	// HTMX target fragments for dynamic updates within case_studies_form.html.
	// Partials:
	//   - case_study_products.html: Related products multi-select (updates via HTMX)
	//   - case_study_metrics.html: Success metrics editor with labels/values (updates via HTMX)
	caseStudyPartials := []string{
		"case_study_products", "case_study_metrics",
	}
	for _, partial := range caseStudyPartials {
		r.templates["admin/partials/"+partial+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/partials/"+partial+".html"),
		))
	}

	// Phase 5: Blog tag partials (HTMX fragments - standalone, no layout)
	// HTMX fragments for dynamic tag and product selection in blog post editor.
	// Partials:
	//   - tag_suggestions.html: Autocomplete tag suggestions (hx-get on input)
	//   - tag_chip.html: Individual tag chip with remove button (hx-delete on click)
	//   - product_suggestions.html: Related product autocomplete (hx-get on input)
	blogPartials := []string{"tag_suggestions", "tag_chip", "product_suggestions"}
	for _, partial := range blogPartials {
		r.templates["admin/partials/"+partial+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
			filepath.Join(r.basePath, "admin/partials/"+partial+".html"),
		))
	}

	// Product search partial (HTMX fragment - standalone, no layout)
	// Used on public products page for live search results (hx-get on search input).
	// Returns filtered product grid without page reload.
	r.templates["public/partials/product_search_results.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/partials/product_search_results.html"),
	))

	// Phase 5: Public blog pages
	// Uses: public/layouts/base.html for public site structure
	// Includes: partials/header.html (navigation), partials/footer.html (footer)
	// Templates:
	//   - blog_listing.html: Blog archive with category/tag/author filtering, pagination
	//   - blog_post.html: Full blog post with author bio, tags, related posts, comments
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
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates:
	//   - blog_posts_list.html: Table of posts with status, category, author, publish date
	//   - blog_post_form.html: Create/edit form with Trix editor, tags, SEO, scheduling
	//   - blog_tags_list.html: Tag management with usage count, merge/delete actions
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
	// Uses: public/layouts/base.html for public site structure
	// Includes: partials/header.html (navigation), partials/footer.html (footer)
	// Templates:
	//   - whitepapers.html: Grid of whitepapers with topic filtering
	//   - whitepaper_detail.html: Whitepaper overview with gated download form
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
	// Shown after whitepaper download form submission via HTMX.
	// Displays success message and download link without page reload.
	r.templates["public/pages/whitepaper_success.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/pages/whitepaper_success.html"),
	))

	// Phase 8: Public contact page
	// Uses: public/layouts/base.html for public site structure
	// Includes: partials/header.html (navigation), partials/footer.html (footer)
	// Content: Contact form with office locations map
	r.templates["public/pages/contact.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/layouts/base.html"),
		filepath.Join(r.basePath, "public/pages/contact.html"),
		filepath.Join(r.basePath, "partials/header.html"),
		filepath.Join(r.basePath, "partials/footer.html"),
	))

	// Phase 8: Admin whitepaper pages
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates:
	//   - whitepapers_list.html: Table of whitepapers with status, topic, downloads count
	//   - whitepapers_form.html: Create/edit form with file upload, gating options
	//   - whitepapers_downloads.html: Analytics page showing download history, user info
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
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates:
	//   - contact_submissions_list.html: Table of contact form submissions with status, date
	//   - contact_submission_detail.html: Full submission view with user info, message, actions
	//   - office_locations_list.html: Table of office locations with address, contact info
	//   - office_locations_form.html: Create/edit form for office location details
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
	// Uses: public/layouts/base.html for public site structure
	// Includes: partials/header.html (navigation), partials/footer.html (footer)
	// Templates:
	//   - about.html: Company overview, mission/vision/values, milestones, certifications
	//   - partners.html: Partner ecosystem with tier-based filtering, logos, testimonials
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
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates manage homepage components and page-level settings:
	//   - homepage_heroes_list/form: Hero banners with background images, headlines, CTAs
	//   - homepage_stats_list/form: Key statistics displayed on homepage
	//   - homepage_testimonials_list/form: Customer testimonials with ratings
	//   - homepage_cta_list/form: Call-to-action sections
	//   - homepage_settings: Homepage metadata (title, description, featured content)
	//   - products_settings: Products page configuration (layout, filters)
	//   - solutions_settings: Solutions page configuration (display options)
	//   - blog_settings: Blog page configuration (posts per page, sidebar content)
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
	// Uses: public/layouts/base.html for public site structure
	// Includes: partials/header.html (navigation), partials/footer.html (footer)
	// Content: Global site search with results across products, blog, solutions, case studies
	r.templates["public/pages/search.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/layouts/base.html"),
		filepath.Join(r.basePath, "public/pages/search.html"),
		filepath.Join(r.basePath, "partials/header.html"),
		filepath.Join(r.basePath, "partials/footer.html"),
	))

	// Phase 9: Search suggestions partial (HTMX fragment - standalone, no layout)
	// Autocomplete suggestions shown while user types in search box (hx-get on input).
	// Returns filtered results without page reload for instant search experience.
	r.templates["public/partials/search_suggestions.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "public/partials/search_suggestions.html"),
	))

	// Phase 7: Admin about and partners pages
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates manage about page content components:
	//   - about_overview_form: Company description and background text
	//   - about_mvv_form: Mission, vision, and values content editor
	//   - core_values_list/form: Individual core values with icons and descriptions
	//   - milestones_list/form: Company history timeline with dates and achievements
	//   - certifications_list/form: Industry certifications with logos and details
	//   - partners_list/form: Partner organizations with tier, logo, description
	//   - testimonials_list/form: Customer testimonials (different from homepage testimonials)
	//   - about_settings: About page configuration (layout, sections visibility)
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
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Content: Centralized media management with upload, organize, search, embed
	r.templates["admin/pages/media_library.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/layouts/base.html"),
		filepath.Join(r.basePath, "admin/pages/media_library.html"),
		filepath.Join(r.basePath, "partials/admin-sidebar.html"),
	))

	// Phase 18: Media picker partial (HTMX fragment - standalone, no layout)
	// Modal overlay for selecting media from library in forms (hx-get on media button click).
	// Allows browsing, searching, and selecting images/files without page navigation.
	r.templates["admin/partials/media_picker.html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/partials/media_picker.html"),
	))

	// Phase 19: Navigation editor pages
	// Uses: admin/layouts/base.html (admin panel structure)
	// Includes: partials/admin-sidebar.html (admin navigation)
	// Templates:
	//   - navigation_list.html: List of navigation menus (header, footer, etc.)
	//   - navigation_editor.html: Drag-and-drop menu editor with nested items, links, ordering
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

// safeHTML marks a string as safe HTML content, bypassing Go's auto-escaping.
// This is required for rendering HTML content from WYSIWYG editors (Trix) and
// database-stored HTML that should be displayed as-is rather than escaped.
//
// Parameters:
//   - s: Raw HTML string to render without escaping
//
// Returns:
//   - template.HTML: Marked-safe HTML that won't be escaped in templates
//
// SECURITY WARNING: Only use with trusted content. User input passed through
// this function without sanitization creates XSS vulnerabilities.
//
// Usage in templates: {{.Content | safeHTML}}
func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

// formatDate converts a time.Time value to a human-readable date string.
// Supports custom format strings using Go's time formatting syntax.
//
// Parameters:
//   - t: Time value to format
//   - format: Go time format string (e.g., "2006-01-02", "Jan 2, 2006")
//            If empty, defaults to "January 2, 2006"
//
// Returns:
//   - string: Formatted date string
//
// Go time format reference:
//   - 2006: four-digit year
//   - 01 or Jan: month
//   - 02: day
//   - 15: 24-hour, 03: 12-hour
//   - 04: minute, 05: second
//
// Usage in templates: {{.PublishedAt | formatDate "Jan 2, 2006"}}
func formatDate(t time.Time, format string) string {
	if format == "" {
		format = "January 2, 2006"
	}
	return t.Format(format)
}

// truncate shortens a string to a maximum length, appending "..." if truncated.
// Used for excerpt generation, preview text, and card descriptions.
//
// Parameters:
//   - s: String to potentially truncate
//   - length: Maximum length in characters (excluding ellipsis)
//
// Returns:
//   - string: Original string if <= length, otherwise truncated with "..." suffix
//
// Note: Truncates at exact character count, not word boundaries.
// For word-boundary truncation, implement custom logic in handlers.
//
// Usage in templates: {{.Description | truncate 100}}
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// slugify converts a string to a URL-safe slug format.
// Converts to lowercase and replaces spaces with hyphens.
//
// Parameters:
//   - s: String to convert (e.g., "Product Name")
//
// Returns:
//   - string: URL-safe slug (e.g., "product-name")
//
// Current limitations:
//   - Only handles spaces, not other special characters
//   - No Unicode normalization (accented characters not simplified)
//   - No duplicate hyphen removal
//
// For production, consider using a dedicated slugify library that handles:
//   - Special characters (punctuation, symbols)
//   - Unicode transliteration (é -> e, ñ -> n)
//   - Duplicate separator removal
//
// Usage in templates: {{.Title | slugify}}
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
