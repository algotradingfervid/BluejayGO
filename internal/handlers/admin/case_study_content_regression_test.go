package admin_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestCaseStudyLegacyEditAndBulletSaveRoundTrip(t *testing.T) {
	db, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	result, err := db.Exec(`INSERT INTO industries (name, slug, icon, description) VALUES ('Healthcare regression', 'healthcare-regression', 'health', 'Test')`)
	if err != nil {
		t.Fatal(err)
	}
	industryID, _ := result.LastInsertId()
	prose := strings.Repeat("Health teams need attendance in remote areas. ", 5)
	item, err := queries.AdminCreateCaseStudy(context.Background(), sqlc.AdminCreateCaseStudyParams{Slug: "legacy-case-regression", Title: "Legacy case", ClientName: "Test", IndustryID: industryID, Summary: "Summary", ChallengeTitle: prose, ChallengeContent: "<p>Existing story</p>", ChallengeBullets: sql.NullString{String: `["GPS, Wi-Fi and Bluetooth","Reliable attendance"]`, Valid: true}, SolutionTitle: "Solution", SolutionContent: "<p>Solution</p>", OutcomeTitle: "Outcome", OutcomeContent: "<p>Outcome</p>"})
	if err != nil {
		t.Fatal(err)
	}
	id := strconv.FormatInt(item.ID, 10)
	e := echo.New()
	e.Renderer = templates.NewRenderer("../../../templates")
	h := admin.NewCaseStudiesHandler(queries, logger, services.NewCache())
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/case-studies/"+id+"/edit", nil), rec)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if err := h.Edit(c); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`name="challenge_title" value="The Challenge"`, "GPS, Wi-Fi and Bluetooth\nReliable attendance", "Existing story"} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("edit form missing %q", want)
		}
	}
	unchanged, err := queries.AdminGetCaseStudy(context.Background(), item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if unchanged.ChallengeTitle != prose {
		t.Fatal("viewing edit form mutated source data")
	}
	form := url.Values{"title": {"Legacy case"}, "slug": {"legacy-case-regression"}, "client_name": {"Test"}, "industry_id": {strconv.FormatInt(industryID, 10)}, "summary": {"Summary"}, "challenge_title": {prose}, "challenge_content": {"<p>Existing story</p>"}, "challenge_bullets": {"GPS, Wi-Fi and Bluetooth\nReliable attendance"}, "solution_title": {"Solution"}, "solution_content": {"<p>Solution</p>"}, "outcome_title": {"Outcome"}, "outcome_content": {"<p>Outcome</p>"}}
	rec, c = postForm(e, "/admin/case-studies/"+id, form)
	c.SetParamNames("id")
	c.SetParamValues(id)
	if err := h.Update(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("save status %d", rec.Code)
	}
	saved, err := queries.AdminGetCaseStudy(context.Background(), item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if saved.ChallengeTitle != "The Challenge" || !strings.Contains(saved.ChallengeContent, strings.TrimSpace(prose)) || !strings.Contains(saved.ChallengeContent, "Existing story") {
		t.Fatal("normalized saved story lost content")
	}
	if saved.ChallengeBullets.String != `["GPS, Wi-Fi and Bluetooth","Reliable attendance"]` {
		t.Fatalf("bullets corrupted: %s", saved.ChallengeBullets.String)
	}
}
