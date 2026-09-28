package public_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/database"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	public "github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
)

// Exercise the actual migration, queries, handlers and templates together: visibility
// was previously saved correctly in the DB but ignored by the homepage template.
func TestHomepageSectionAndPartnerRegressions(t *testing.T) {
	db, err := database.InitDB(database.Config{Path: filepath.Join(t.TempDir(), "home.db")})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := database.RunMigrations(db, "../../../db/migrations"); err != nil {
		t.Fatal(err)
	}
	q := sqlc.New(db)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	home := public.NewHomeHandler(q, logger)
	exec := func(stmt string) {
		t.Helper()
		if _, err := db.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	render := func() string {
		t.Helper()
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
		if err := home.ShowHomePage(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != 200 {
			t.Fatalf("home status %d", rec.Code)
		}
		return rec.Body.String()
	}
	exec("UPDATE page_sections SET is_active=0 WHERE page_key='home'")
	if body := render(); strings.Contains(body, "data-home-section=") {
		t.Fatal("disabled sections still rendered")
	}
	exec("UPDATE page_sections SET is_active=1,display_order=10 WHERE page_key='home' AND section_key='hero'")
	exec("UPDATE page_sections SET is_active=1,display_order=1 WHERE page_key='home' AND section_key='partners_section'")
	exec("UPDATE partners SET is_featured=0")
	exec("INSERT INTO partner_tiers(name,slug,description) VALUES ('Regression','regression','Test tier')")
	exec("INSERT INTO homepage_hero(headline,subheadline) VALUES ('Test hero','Test subheading')")
	exec("INSERT INTO homepage_cta(headline) VALUES ('Test CTA')")
	exec("INSERT INTO homepage_testimonials(quote,author_name,display_order) VALUES ('First regression quote','A',0),('Second regression quote','B',1)")
	var tierID int64
	if err := db.QueryRow("SELECT id FROM partner_tiers ORDER BY id LIMIT 1").Scan(&tierID); err != nil {
		t.Fatal(err)
	}
	partners := admin.NewPartnersHandler(q, logger, services.NewCache())
	for i, name := range []string{"Featured A", "Featured B"} {
		form := url.Values{"name": {name}, "tier_id": {strconv.FormatInt(tierID, 10)}, "is_featured": {"on"}, "display_order": {strconv.Itoa(i)}, "website_url": {"https://example.com/" + strconv.Itoa(i)}, "logo_url": {"/uploads/partner-" + strconv.Itoa(i) + ".png"}}
		req := httptest.NewRequest(http.MethodPost, "/admin/partners", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		c := e.NewContext(req, httptest.NewRecorder())
		if err := partners.Create(c); err != nil {
			t.Fatal(err)
		}
	}
	body := render()
	for _, name := range []string{"Featured A", "Featured B"} {
		if !strings.Contains(body, name) {
			t.Errorf("missing selected partner %q", name)
		}
	}
	if a, b := strings.Index(body, `data-home-section="partners_section"`), strings.Index(body, `data-home-section="hero"`); a < 0 || b < 0 || a >= b {
		t.Fatalf("configured section order ignored: partners=%d hero=%d", a, b)
	}
	var partnerID int64
	if err := db.QueryRow("SELECT id FROM partners WHERE name='Featured A'").Scan(&partnerID); err != nil {
		t.Fatal(err)
	}
	form := url.Values{"name": {"Featured A"}, "tier_id": {strconv.FormatInt(tierID, 10)}, "is_active": {"1"}, "display_order": {"0"}}
	req := httptest.NewRequest(http.MethodPost, "/admin/partners/"+strconv.FormatInt(partnerID, 10), strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c := e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(partnerID, 10))
	if err := partners.Update(c); err != nil {
		t.Fatal(err)
	}
	if body := render(); strings.Contains(body, "Featured A") {
		t.Fatal("unchecked partner still featured")
	}
	exec("UPDATE partners SET is_active=0 WHERE name='Featured B'")
	if body := render(); strings.Contains(body, "Featured B") {
		t.Fatal("inactive partner still featured")
	}

	// Both settings and page section controls must hide the complete block.
	exec("UPDATE page_sections SET is_active=1 WHERE page_key='home'")
	exec("UPDATE settings SET homepage_show_heroes=0,homepage_show_stats=0,homepage_show_testimonials=0,homepage_show_cta=0")
	body = render()
	for _, key := range []string{"hero", "stats_section", "testimonials_section", "cta"} {
		if strings.Contains(body, `data-home-section="`+key+`"`) {
			t.Errorf("homepage setting failed to hide %s", key)
		}
	}
	exec("UPDATE settings SET homepage_show_stats=1,homepage_show_testimonials=1,homepage_max_stats=1,homepage_max_testimonials=1")
	exec("UPDATE homepage_stats SET is_active=0")
	exec("INSERT INTO homepage_stats(stat_value,stat_label,display_order) VALUES ('1','First regression stat',0),('2','Second regression stat',1)")
	body = render()
	if !strings.Contains(body, "First regression quote") || strings.Contains(body, "Second regression quote") {
		t.Fatal("testimonial display limit ignored")
	}
	if !strings.Contains(body, "First regression stat") || strings.Contains(body, "Second regression stat") {
		t.Fatal("stats display limit ignored")
	}

	// Editing a section must persist order as well as visibility.
	section, err := q.GetPageSection(context.Background(), sqlc.GetPageSectionParams{PageKey: "home", SectionKey: "partners_section"})
	if err != nil {
		t.Fatal(err)
	}
	form = url.Values{"display_order": {"42"}, "heading": {"Partners updated"}}
	req = httptest.NewRequest(http.MethodPost, "/admin/page-sections", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c = e.NewContext(req, httptest.NewRecorder())
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(section.ID, 10))
	if err := admin.NewPageSectionsHandler(q, logger).Update(c); err != nil {
		t.Fatal(err)
	}
	updated, err := q.GetPageSectionByID(context.Background(), section.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.IsActive || updated.DisplayOrder != 42 {
		t.Fatalf("section admin save not respected: %+v", updated)
	}

	// Partner logos link to configured sites; missing URLs must not produce empty anchors.
	rec := httptest.NewRecorder()
	c = e.NewContext(httptest.NewRequest("GET", "/partners", nil), rec)
	data := map[string]interface{}{"Settings": sqlc.Setting{}, "Tiers": []sqlc.PartnerTier{{Name: "Test"}}, "PartnersByTier": map[string][]sqlc.ListPartnersByTierRow{"Test": {{Name: "Linked partner", WebsiteUrl: sql.NullString{String: "https://example.com/partner", Valid: true}, LogoUrl: sql.NullString{String: "/logo.png", Valid: true}}, {Name: "Unlinked partner"}}}}
	if err := e.Renderer.Render(rec, "public/pages/partners.html", data, c); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(rec.Body.String(), `href="https://example.com/partner"`) || strings.Contains(rec.Body.String(), `href=""`) {
		t.Fatal("partner website link missing or empty link rendered")
	}
}

func TestHomepageLegacySectionOrderMigration(t *testing.T) {
	db, err := database.InitDB(database.Config{Path: filepath.Join(t.TempDir(), "legacy.db")})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := database.RunMigrations(db, "../../../db/migrations"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("UPDATE page_sections SET display_order=0 WHERE page_key='home'; UPDATE page_sections SET display_order=1 WHERE page_key='home' AND section_key='hero'; UPDATE page_sections SET display_order=99 WHERE page_key='home' AND section_key='products_section'; UPDATE page_sections SET is_active=0 WHERE page_key='home' AND section_key='solutions_section'"); err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../db/migrations/041_homepage_cta_section.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ { // Also check idempotence; a rerun must preserve chosen values.
		if _, err := db.Exec(string(migration)); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := sqlc.New(db).ListAllPageSections(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]int64{"hero": 1, "products_section": 99, "solutions_section": 3, "stats_section": 4, "partners_section": 5, "testimonials_section": 6, "blog_section": 7, "cta": 8}
	for _, section := range rows {
		if section.PageKey != "home" {
			continue
		}
		if order, ok := expected[section.SectionKey]; ok && section.DisplayOrder != order {
			t.Errorf("%s order=%d want %d", section.SectionKey, section.DisplayOrder, order)
		}
		if section.SectionKey == "solutions_section" && section.IsActive {
			t.Error("migration changed saved visibility")
		}
	}
}
