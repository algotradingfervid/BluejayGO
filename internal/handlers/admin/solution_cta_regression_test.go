package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestSolutionCTAControlsPersistAndInvalidateCache(t *testing.T) {
	db, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	result, err := db.Exec(`INSERT INTO solutions (title,slug,icon,short_description) VALUES ('Corporate regression','corporate-regression','business','Test')`)
	if err != nil {
		t.Fatal(err)
	}
	solutionID, _ := result.LastInsertId()
	id := strconv.FormatInt(solutionID, 10)
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	cache := services.NewCache()
	h := admin.NewSolutionsHandler(queries, logger, cache, nil)
	form := url.Values{"heading": {"Corporate CTA"}, "section_name": {"main_cta"}, "primary_button_text": {"Contact Us"}, "primary_button_url": {"/contact"}, "secondary_button_text": {"Download Brochure"}, "secondary_button_url": {"https://example.com/brochure.pdf"}, "is_active": {"1"}, "secondary_button_enabled": {"1"}}
	rec, c := postForm(e, "/admin/solutions/"+id+"/ctas", form)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if err := h.AddCTA(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("add status %d: %s", rec.Code, rec.Body.String())
	}
	ctAs, err := queries.GetSolutionCTAs(context.Background(), solutionID)
	if err != nil || len(ctAs) != 1 {
		t.Fatalf("read CTAs: %v %v", ctAs, err)
	}
	cta := ctAs[0]
	if !cta.IsActive || !cta.SecondaryButtonEnabled {
		t.Fatal("initial controls not enabled")
	}
	ctaID := strconv.FormatInt(cta.ID, 10)
	// The edit form renders both visibility controls and preserves checked state.
	rec = httptest.NewRecorder()
	c = e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/solutions/"+id+"/ctas-tab?edit="+ctaID, nil), rec)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if err := h.CTAsTab(c); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`name="is_active" value="1" checked`, `name="secondary_button_enabled" value="1" checked`} {
		if !strings.Contains(rec.Body.String(), field) {
			t.Errorf("missing checked control %s", field)
		}
	}
	form.Del("is_active")
	form.Del("secondary_button_enabled")
	cache.Set("page:solutions:corporate-regression", "stale", 600)
	rec, c = postForm(e, "/admin/solutions/"+id+"/ctas/"+ctaID, form)
	c.SetParamNames("id", "ctaId")
	c.SetParamValues(id, ctaID)
	if err := h.UpdateCTA(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("update status %d", rec.Code)
	}
	ctAs, err = queries.GetSolutionCTAs(context.Background(), solutionID)
	if err != nil {
		t.Fatal(err)
	}
	if ctAs[0].IsActive || ctAs[0].SecondaryButtonEnabled {
		t.Fatal("disabled controls not saved")
	}
	if ctAs[0].SecondaryButtonUrl.String != "https://example.com/brochure.pdf" {
		t.Fatal("disabling erased configured brochure")
	}
	if _, ok := cache.Get("page:solutions:corporate-regression"); ok {
		t.Fatal("stale public page remains cached")
	}
}
