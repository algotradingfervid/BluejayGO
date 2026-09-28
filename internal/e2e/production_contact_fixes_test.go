package e2e_test

import (
	"context"
	"fmt"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestProductionOfficeMapPersistenceAndRendering(t *testing.T) {
	app, queries, cleanup := setupApp(t)
	defer cleanup()
	app.Renderer = templates.NewRenderer("templates")
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, app)
	request := func(method, path string, form url.Values) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		return rec
	}
	form := url.Values{"name": {"Map Test Office"}, "address_line1": {"123 Main Street"}, "city": {"Hyderabad"}, "state": {"Telangana"}, "postal_code": {"500104"}, "country": {"India"}, "is_active": {"1"}, "map_url": {"https://www.google.com/maps?q=Hyderabad"}}
	if rec := request("POST", "/admin/contact/offices", form); rec.Code != http.StatusSeeOther {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	offices, err := queries.ListAllOfficeLocations(context.Background())
	if err != nil || len(offices) != 1 {
		t.Fatalf("offices: %v %v", offices, err)
	}
	id := offices[0].ID
	if offices[0].MapUrl != form.Get("map_url") {
		t.Fatal("map URL not persisted on create")
	}
	rec := request("GET", fmt.Sprintf("/admin/contact/offices/%d/edit", id), nil)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), form.Get("map_url")) {
		t.Fatalf("edit form lost stored URL: %d %s", rec.Code, rec.Body.String())
	}
	rec = request("GET", "/contact", nil)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `<iframe src="https://www.google.com/maps?q=Hyderabad`) || strings.Contains(rec.Body.String(), "Map Integration") {
		t.Fatalf("map not rendered: %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Security-Policy"), "frame-src https://www.google.com/maps https://www.google.com/maps/;") {
		t.Fatal("CSP must permit both search and embed map paths")
	}
	form.Set("map_url", "https://www.google.com/maps/embed?pb=!1m18!2m3")
	if rec := request("POST", fmt.Sprintf("/admin/contact/offices/%d", id), form); rec.Code != http.StatusSeeOther {
		t.Fatalf("embed update: %d", rec.Code)
	}
	rec = request("GET", "/contact", nil)
	if !strings.Contains(rec.Body.String(), "https://www.google.com/maps/embed?pb=") {
		t.Fatal("embed map not rendered")
	}
	form.Set("map_url", "https://maps.app.goo.gl/Example")
	if rec := request("POST", fmt.Sprintf("/admin/contact/offices/%d", id), form); rec.Code != http.StatusSeeOther {
		t.Fatalf("update: %d", rec.Code)
	}
	office, err := queries.GetOfficeLocationByID(context.Background(), id)
	if err != nil || office.MapUrl != form.Get("map_url") {
		t.Fatal("map URL not persisted on update")
	}
	rec = request("GET", "/contact", nil)
	if !strings.Contains(rec.Body.String(), "https://maps.app.goo.gl/Example") || !strings.Contains(rec.Body.String(), "maps?q=123") {
		t.Fatal("map update did not invalidate cached contact page")
	}
	form.Set("map_url", "https://evil.example/map")
	rec = request("POST", fmt.Sprintf("/admin/contact/offices/%d", id), form)
	if rec.Code != 400 || !strings.Contains(rec.Body.String(), "Google Maps HTTPS") || !strings.Contains(rec.Body.String(), "Map Test Office") {
		t.Fatal("unsafe URL must produce useful feedback and preserve form")
	}
	office, _ = queries.GetOfficeLocationByID(context.Background(), id)
	if office.MapUrl != "https://maps.app.goo.gl/Example" {
		t.Fatal("invalid URL changed stored map")
	}
}
