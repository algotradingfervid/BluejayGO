package admin_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	admin "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

func init() {
	middleware.InitSessionStore("test-secret-at-least-32-characters-long")
}

func postForm(e *echo.Echo, path string, values url.Values) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(values.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return rec, c
}

func TestProductCategoriesHandler_CRUD(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := admin.NewProductCategoriesHandler(queries, logger)

	ctx := context.Background()

	// Create a category directly
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test", Slug: "test", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Verify via queries
	items, err := queries.ListProductCategories(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("expected 1, got %d", len(items))
	}

	// Delete handler
	req := httptest.NewRequest(http.MethodDelete, "/admin/product-categories/"+strconv.FormatInt(cat.ID, 10), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(cat.ID, 10))

	if err := h.Delete(c); err != nil {
		t.Fatalf("Delete handler: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestAuthHandler_LoginSubmit_MissingFields(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := admin.NewAuthHandler(queries, logger)

	// Empty form
	rec, c := postForm(e, "/admin/login", url.Values{})

	// Set session context (required by handler)
	sess := &middleware.Session{}
	c.Set("session", sess)

	if err := h.LoginSubmit(c); err != nil {
		t.Fatalf("LoginSubmit error: %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "error=missing_fields") {
		t.Errorf("expected error=missing_fields redirect, got %q", loc)
	}
}

func TestAuthHandler_LoginSubmit_InvalidUser(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := admin.NewAuthHandler(queries, logger)

	rec, c := postForm(e, "/admin/login", url.Values{
		"email":    {"nonexistent@test.com"},
		"password": {"password"},
	})
	c.Set("session", &middleware.Session{})

	if err := h.LoginSubmit(c); err != nil {
		t.Fatalf("LoginSubmit error: %v", err)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "error=invalid_credentials") {
		t.Errorf("expected invalid_credentials redirect, got %q", loc)
	}
}

func TestAuthHandler_LoginSubmit_WrongPassword(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	_, err := queries.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{
		Email: "admin@test.com", PasswordHash: string(hash), DisplayName: "Admin", Role: "admin",
	})
	if err != nil {
		t.Fatalf("CreateAdminUser: %v", err)
	}

	e := echo.New()
	h := admin.NewAuthHandler(queries, logger)

	rec, c := postForm(e, "/admin/login", url.Values{
		"email":    {"admin@test.com"},
		"password": {"wrong-password"},
	})
	c.Set("session", &middleware.Session{})

	if err := h.LoginSubmit(c); err != nil {
		t.Fatalf("LoginSubmit error: %v", err)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "error=invalid_credentials") {
		t.Errorf("expected invalid_credentials, got %q", loc)
	}
}

func TestAuthHandler_LoginSubmit_Success(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	_, err := queries.CreateAdminUser(ctx, sqlc.CreateAdminUserParams{
		Email: "admin@test.com", PasswordHash: string(hash), DisplayName: "Admin", Role: "admin",
	})
	if err != nil {
		t.Fatalf("CreateAdminUser: %v", err)
	}

	e := echo.New()
	h := admin.NewAuthHandler(queries, logger)

	req := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(url.Values{
		"email":    {"admin@test.com"},
		"password": {"correct-password"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Need a real session from the store for Save to work
	gorillaSession, _ := middleware.SessionStore.Get(req, "bluejay_session")
	sess := &middleware.Session{Session: gorillaSession}
	c.Set("session", sess)

	if err := h.LoginSubmit(c); err != nil {
		t.Fatalf("LoginSubmit error: %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if loc != "/admin/dashboard" {
		t.Errorf("expected redirect to /admin/dashboard, got %q", loc)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := admin.NewAuthHandler(queries, logger)

	req := httptest.NewRequest(http.MethodPost, "/admin/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	gorillaSession, _ := middleware.SessionStore.Get(req, "bluejay_session")
	sess := &middleware.Session{Session: gorillaSession, UserID: 1, Email: "admin@test.com"}
	c.Set("session", sess)

	if err := h.Logout(c); err != nil {
		t.Fatalf("Logout error: %v", err)
	}
	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if loc != "/admin/login" {
		t.Errorf("expected redirect to /admin/login, got %q", loc)
	}
}

func TestBlogCategoriesHandler_Delete(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, err := queries.CreateBlogCategory(ctx, sqlc.CreateBlogCategoryParams{
		Name: "Tech", Slug: "tech", ColorHex: "#FF0000", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateBlogCategory: %v", err)
	}

	e := echo.New()
	h := admin.NewBlogCategoriesHandler(queries, logger)

	req := httptest.NewRequest(http.MethodDelete, "/admin/blog-categories/"+strconv.FormatInt(cat.ID, 10), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(cat.ID, 10))

	if err := h.Delete(c); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestIndustriesHandler_Delete(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	ind, err := queries.CreateIndustry(ctx, sqlc.CreateIndustryParams{
		Name: "Health", Slug: "health", Icon: "med", Description: "d", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateIndustry: %v", err)
	}

	e := echo.New()
	h := admin.NewIndustriesHandler(queries, logger)

	req := httptest.NewRequest(http.MethodDelete, "/admin/industries/"+strconv.FormatInt(ind.ID, 10), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(ind.ID, 10))

	if err := h.Delete(c); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestProductsHandler_Delete(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Cat", Slug: "cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateProductCategory: %v", err)
	}

	prod, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DEL-001", Slug: "del-prod", Name: "Del Prod", Description: "d", CategoryID: cat.ID, Status: "draft",
	})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	e := echo.New()
	uploadSvc := services.NewUploadService(t.TempDir())
	cache := services.NewCache()
	h := admin.NewProductsHandler(queries, logger, uploadSvc, cache)

	req := httptest.NewRequest(http.MethodDelete, "/admin/products/"+strconv.FormatInt(prod.ID, 10), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(prod.ID, 10))

	if err := h.Delete(c); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// Ensure sql import is used
var _ = sql.NullString{}
