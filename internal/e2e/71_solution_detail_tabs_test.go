package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// setupSolutionTabsApp builds a local Echo instance with the REAL template renderer
// (so partial bodies are actually rendered, not the shared stub) and registers the
// admin auth + solution detail-tab routes. The shared setupApp uses a stub renderer
// and does not register the *-tab routes, so it cannot exercise this bug.
func setupSolutionTabsApp(t *testing.T) (*echo.Echo, *sqlc.Queries, func()) {
	t.Helper()

	_, queries, cleanup := testutil.SetupTestDB(t)
	customMiddleware.InitSessionStore("e2e-test-secret-at-least-32-characters-long")

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")
	e.Use(customMiddleware.SecurityHeaders())
	e.Use(customMiddleware.SessionMiddleware())

	appCache := services.NewCache()
	uploadSvc := services.NewUploadService(t.TempDir())
	activitySvc := services.NewActivityLogService(queries, logger)
	adminHandlers.SetActivityLogService(activitySvc)

	authHandler := adminHandlers.NewAuthHandler(queries, logger)
	e.GET("/admin/login", authHandler.ShowLoginPage)
	e.POST("/admin/login", authHandler.LoginSubmit)

	adminGroup := e.Group("/admin", customMiddleware.RequireAuth())
	h := adminHandlers.NewSolutionsHandler(queries, logger, appCache, uploadSvc)
	adminGroup.GET("/solutions/:id/challenges-tab", h.ChallengesTab)
	adminGroup.GET("/solutions/:id/products-tab", h.ProductsTab)
	adminGroup.GET("/solutions/:id/stats-tab", h.StatsTab)
	adminGroup.GET("/solutions/:id/ctas-tab", h.CTAsTab)
	adminGroup.POST("/solutions/:id/challenges/:challengeId", h.UpdateChallenge)
	adminGroup.POST("/solutions/:id/stats/:statId", h.UpdateStat)
	adminGroup.POST("/solutions/:id/ctas/:ctaId", h.UpdateCTA)
	adminGroup.POST("/solutions/:id/products/:productId", h.UpdateProduct)

	return e, queries, cleanup
}

func loginTabsAdmin(t *testing.T, e *echo.Echo, queries *sqlc.Queries) *http.Cookie {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	_, err := queries.CreateAdminUser(context.Background(), sqlc.CreateAdminUserParams{
		Email:        "admin@test.com",
		PasswordHash: string(hash),
		DisplayName:  "Test Admin",
		Role:         "admin",
	})
	if err != nil {
		t.Fatalf("create admin: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{
		"email":    {"admin@test.com"},
		"password": {"testpassword"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	for _, c := range rec.Result().Cookies() {
		if c.Name == "bluejay_session" {
			return c
		}
	}
	t.Fatal("no session cookie after login")
	return nil
}

// seedSolutionWithChildren creates a solution plus one challenge, stat, and CTA so the
// tab partials render both the range items and the add form.
func seedSolutionWithChildren(t *testing.T, queries *sqlc.Queries) int64 {
	t.Helper()
	ctx := context.Background()
	sol, err := queries.CreateSolution(ctx, sqlc.CreateSolutionParams{
		Title:            "Test Solution",
		Slug:             "test-solution",
		Icon:             "icon",
		ShortDescription: "desc",
		IsPublished:      sql.NullBool{Bool: true, Valid: true},
		DisplayOrder:     sql.NullInt64{Int64: 0, Valid: true},
	})
	if err != nil {
		t.Fatalf("create solution: %v", err)
	}

	if _, err := queries.CreateSolutionChallenge(ctx, sqlc.CreateSolutionChallengeParams{
		SolutionID: sol.ID, Title: "Challenge A", Description: "d", Icon: "i",
		DisplayOrder: sql.NullInt64{Int64: 0, Valid: true},
	}); err != nil {
		t.Fatalf("create challenge: %v", err)
	}
	if _, err := queries.CreateSolutionStat(ctx, sqlc.CreateSolutionStatParams{
		SolutionID: sol.ID, Value: "99%", Label: "Uptime",
		DisplayOrder: sql.NullInt64{Int64: 0, Valid: true},
	}); err != nil {
		t.Fatalf("create stat: %v", err)
	}
	if _, err := queries.CreateSolutionCTA(ctx, sqlc.CreateSolutionCTAParams{
		SolutionID: sol.ID, Heading: "Get Started", SectionName: "footer",
	}); err != nil {
		t.Fatalf("create cta: %v", err)
	}
	return sol.ID
}

// TestSolutionDetailTabs_RenderWithCorrectSolutionID verifies the four detail-tab
// endpoints exist (return 200, not 404) AND that each rendered partial's add form
// posts to the correct solution id (proving SolutionID is passed at the top level
// of the render map — covering both the missing-route bug and the latent
// missing-SolutionID bug in the Add* handlers' re-renders).
func TestSolutionDetailTabs_RenderWithCorrectSolutionID(t *testing.T) {
	e, queries, cleanup := setupSolutionTabsApp(t)
	defer cleanup()
	cookie := loginTabsAdmin(t, e, queries)
	id := seedSolutionWithChildren(t, queries)

	cases := []struct {
		tab     string
		addPath string
	}{
		{"challenges-tab", "challenges"},
		{"products-tab", "products"},
		{"stats-tab", "stats"},
		{"ctas-tab", "ctas"},
	}

	for _, tc := range cases {
		t.Run(tc.tab, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/solutions/%d/%s", id, tc.tab), nil)
			req.AddCookie(cookie)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200 from /admin/solutions/%d/%s, got %d; body: %s", id, tc.tab, rec.Code, rec.Body.String())
			}

			body := rec.Body.String()
			wantPost := fmt.Sprintf(`hx-post="/admin/solutions/%d/%s"`, id, tc.addPath)
			if !strings.Contains(body, wantPost) {
				t.Fatalf("tab %s body should contain add-form %q (proving SolutionID passed at top level), got: %s", tc.tab, wantPost, body)
			}
			// Guard against the empty-id regression where SolutionID is missing.
			if strings.Contains(body, fmt.Sprintf(`hx-post="/admin/solutions//%s"`, tc.addPath)) {
				t.Fatalf("tab %s add-form posts to empty solution id (//%s) — SolutionID not passed at top level", tc.tab, tc.addPath)
			}
		})
	}
}

// TestSolutionDetailTabs_ChallengeInlineEdit verifies the challenges tab renders an
// inline edit form when ?edit={id} is passed, and that POSTing the edit form persists.
func TestSolutionDetailTabs_ChallengeInlineEdit(t *testing.T) {
	e, queries, cleanup := setupSolutionTabsApp(t)
	defer cleanup()
	cookie := loginTabsAdmin(t, e, queries)
	ctx := context.Background()
	id := seedSolutionWithChildren(t, queries)

	challenges, err := queries.GetSolutionChallenges(ctx, id)
	if err != nil || len(challenges) == 0 {
		t.Fatalf("seed challenge: %v (n=%d)", err, len(challenges))
	}
	cid := challenges[0].ID

	// GET with ?edit renders the inline edit form.
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/solutions/%d/challenges-tab?edit=%d", id, cid), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	wantPost := fmt.Sprintf(`hx-post="/admin/solutions/%d/challenges/%d"`, id, cid)
	if !strings.Contains(body, wantPost) {
		t.Fatalf("edit form missing %q; body: %s", wantPost, body)
	}
	if !strings.Contains(body, `value="Challenge A"`) {
		t.Fatalf("edit form not pre-filled with existing title; body: %s", body)
	}

	// POST the update persists.
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/challenges/%d", id, cid), strings.NewReader(url.Values{
		"title":         {"Updated Challenge"},
		"description":   {"new desc"},
		"icon":          {"bolt"},
		"display_order": {"3"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from update, got %d; body: %s", rec.Code, rec.Body.String())
	}
	got, _ := queries.GetSolutionChallenges(ctx, id)
	if len(got) == 0 || got[0].Title != "Updated Challenge" {
		t.Fatalf("challenge not updated; got %+v", got)
	}
}

// TestSolutionDetailTabs_StatInlineEdit verifies inline edit + persist for stats.
func TestSolutionDetailTabs_StatInlineEdit(t *testing.T) {
	e, queries, cleanup := setupSolutionTabsApp(t)
	defer cleanup()
	cookie := loginTabsAdmin(t, e, queries)
	ctx := context.Background()
	id := seedSolutionWithChildren(t, queries)

	stats, err := queries.GetSolutionStats(ctx, id)
	if err != nil || len(stats) == 0 {
		t.Fatalf("seed stat: %v (n=%d)", err, len(stats))
	}
	sid := stats[0].ID

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/solutions/%d/stats-tab?edit=%d", id, sid), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	wantPost := fmt.Sprintf(`hx-post="/admin/solutions/%d/stats/%d"`, id, sid)
	if !strings.Contains(body, wantPost) {
		t.Fatalf("edit form missing %q; body: %s", wantPost, body)
	}
	if !strings.Contains(body, `value="Uptime"`) {
		t.Fatalf("edit form not pre-filled with existing label; body: %s", body)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/stats/%d", id, sid), strings.NewReader(url.Values{
		"value":         {"42%"},
		"label":         {"Conversion"},
		"display_order": {"5"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from update, got %d; body: %s", rec.Code, rec.Body.String())
	}
	got, _ := queries.GetSolutionStats(ctx, id)
	if len(got) == 0 || got[0].Label != "Conversion" || got[0].Value != "42%" {
		t.Fatalf("stat not updated; got %+v", got)
	}
}

// TestSolutionDetailTabs_CTAInlineEdit verifies inline edit + persist for CTAs.
func TestSolutionDetailTabs_CTAInlineEdit(t *testing.T) {
	e, queries, cleanup := setupSolutionTabsApp(t)
	defer cleanup()
	cookie := loginTabsAdmin(t, e, queries)
	ctx := context.Background()
	id := seedSolutionWithChildren(t, queries)

	ctas, err := queries.GetSolutionCTAs(ctx, id)
	if err != nil || len(ctas) == 0 {
		t.Fatalf("seed cta: %v (n=%d)", err, len(ctas))
	}
	ctaID := ctas[0].ID

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/solutions/%d/ctas-tab?edit=%d", id, ctaID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	wantPost := fmt.Sprintf(`hx-post="/admin/solutions/%d/ctas/%d"`, id, ctaID)
	if !strings.Contains(body, wantPost) {
		t.Fatalf("edit form missing %q; body: %s", wantPost, body)
	}
	if !strings.Contains(body, `value="Get Started"`) {
		t.Fatalf("edit form not pre-filled with existing heading; body: %s", body)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/ctas/%d", id, ctaID), strings.NewReader(url.Values{
		"heading":      {"Contact Sales"},
		"section_name": {"footer"},
		"subheading":   {"Talk to us"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from update, got %d; body: %s", rec.Code, rec.Body.String())
	}
	got, _ := queries.GetSolutionCTAs(ctx, id)
	if len(got) == 0 || got[0].Heading != "Contact Sales" {
		t.Fatalf("cta not updated; got %+v", got)
	}
}

// TestSolutionDetailTabs_ProductInlineEdit verifies inline edit + persist for products
// (display_order + is_featured are editable; product_id is not).
func TestSolutionDetailTabs_ProductInlineEdit(t *testing.T) {
	e, queries, cleanup := setupSolutionTabsApp(t)
	defer cleanup()
	cookie := loginTabsAdmin(t, e, queries)
	ctx := context.Background()
	id := seedSolutionWithChildren(t, queries)

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	prod, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "PROD-001", Slug: "prod-001", Name: "Test Product", Description: "desc",
		CategoryID: cat.ID, Status: "published",
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}
	if err := queries.AddProductToSolution(ctx, sqlc.AddProductToSolutionParams{
		SolutionID: id, ProductID: prod.ID,
		DisplayOrder: sql.NullInt64{Int64: 0, Valid: true},
		IsFeatured:   sql.NullBool{Bool: false, Valid: true},
	}); err != nil {
		t.Fatalf("link product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/solutions/%d/products-tab?edit=%d", id, prod.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	wantPost := fmt.Sprintf(`hx-post="/admin/solutions/%d/products/%d"`, id, prod.ID)
	if !strings.Contains(body, wantPost) {
		t.Fatalf("edit form missing %q; body: %s", wantPost, body)
	}

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/solutions/%d/products/%d", id, prod.ID), strings.NewReader(url.Values{
		"display_order": {"7"},
		"is_featured":   {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from update, got %d; body: %s", rec.Code, rec.Body.String())
	}
	got, _ := queries.GetSolutionProducts(ctx, id)
	if len(got) == 0 || !got[0].IsFeatured.Bool || got[0].DisplayOrder.Int64 != 7 {
		t.Fatalf("product not updated; got %+v", got)
	}
}
