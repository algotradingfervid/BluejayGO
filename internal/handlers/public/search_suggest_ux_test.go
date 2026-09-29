package public_test

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/internal/database"
	"github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/templates"
)

func TestSearchSuggestEmptyAndRecovery(t *testing.T) {
	db, err := database.InitDB(database.Config{Path: filepath.Join(t.TempDir(), "search.db")})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := database.RunMigrations(db, "../../../db/migrations"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO product_categories(name,slug,description,icon) VALUES ('Search UX','search-ux','',''); INSERT INTO products(sku,slug,name,description,category_id,status) VALUES ('SEARCH-UX','search-ux','UniqueSearchFixture','',last_insert_rowid(),'published')`); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	h := public.NewSearchHandler(db, slog.New(slog.NewTextHandler(io.Discard, nil)))
	render := func(q string) string {
		t.Helper()
		rec := httptest.NewRecorder()
		c := e.NewContext(httptest.NewRequest("GET", "/search/suggest?q="+url.QueryEscape(q), nil), rec)
		if err := h.SearchSuggest(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != 200 {
			t.Fatal(rec.Code)
		}
		return rec.Body.String()
	}
	if body := render("  noSuchKeywordUX  "); !strings.Contains(body, "No results for “noSuchKeywordUX”") || !strings.Contains(body, `href="/products"`) || !strings.Contains(body, "shorter keyword") {
		t.Fatal(body)
	}
	if body := render("UniqueSearchFixture"); strings.Contains(body, "No results") || !strings.Contains(body, `href="/products/search-ux/search-ux"`) {
		t.Fatal(body)
	}
	if body := render(" \t "); strings.Contains(body, "No results") || !strings.Contains(body, "Enter a keyword") {
		t.Fatal(body)
	}
	if body := render(`<script>alert(1)</script>`); strings.Contains(body, "<script>") || !strings.Contains(body, "&lt;script&gt;") {
		t.Fatal(body)
	}
}

func TestHeroCTALabelMigrationPreservesCustomContent(t *testing.T) {
	db, err := database.InitDB(database.Config{Path: filepath.Join(t.TempDir(), "cta.db")})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := database.RunMigrations(db, "../../../db/migrations"); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`DELETE FROM homepage_hero; INSERT INTO homepage_hero(headline,subheadline,primary_cta_text,primary_cta_url) VALUES ('Seeded','','View Displays','/products'),('Custom destination','','View Displays','/products/displays'),('Custom copy','','Explore Products','/products')`)
	if err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../db/migrations/047_hero_cta_destination_label.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := db.Exec(string(migration)); err != nil {
			t.Fatal(err)
		}
	}
	for headline, want := range map[string]string{"Seeded": "View All Products", "Custom destination": "View Displays", "Custom copy": "Explore Products"} {
		var got string
		if err := db.QueryRow("SELECT primary_cta_text FROM homepage_hero WHERE headline=?", headline).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%s: got %q want %q", headline, got, want)
		}
	}
}
