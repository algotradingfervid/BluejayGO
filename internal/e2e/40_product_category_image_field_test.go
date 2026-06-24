package e2e_test

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	adminHandlers "github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/templates"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

// TestProductCategoryForm_ImageFieldAcceptsRelativePath verifies that the admin
// product category form's "Image URL" field can hold the relative upload paths
// the CMS actually uses (e.g. "/uploads/categories/desktops.jpg" — the format in
// the seed data and rendered directly into <img src> on the public pages).
//
// The bug: the field was rendered as <input type="url">, whose HTML5 validation
// rejects relative paths (it requires an absolute URL with a scheme). That makes
// the browser silently block form submission, so an admin's display image never
// gets saved and the category image "link" never works.
//
// setupApp uses a stub renderer that can't see template output, so this test
// builds a local Echo with the REAL renderer (TestMain has chdir'd to project
// root, so "templates" resolves) and asserts on the rendered markup.
func TestProductCategoryForm_ImageFieldAcceptsRelativePath(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	// Seed a category whose display image is a relative upload path, the format
	// used everywhere else in the CMS (seed data, public <img src>).
	const relPath = "/uploads/categories/desktops.jpg"
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name:        "Desktops",
		Slug:        "desktops",
		Description: "d",
		Icon:        "computer",
		ImageUrl:    sql.NullString{String: relPath, Valid: true},
		SortOrder:   1,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Renderer = templates.NewRenderer("templates")
	h := adminHandlers.NewProductCategoriesHandler(queries, logger)
	e.GET("/admin/product-categories/:id/edit", h.Edit)

	req := httptest.NewRequest(http.MethodGet, "/admin/product-categories/"+strconv.FormatInt(cat.ID, 10)+"/edit", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from edit form, got %d; body: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()

	// Sanity: the stored relative path should be pre-filled into the field.
	if !strings.Contains(body, relPath) {
		t.Fatalf("expected stored image path %q to be pre-filled in the form; body:\n%s", relPath, body)
	}

	// The image field must NOT be type="url": HTML5 url validation rejects the
	// relative "/uploads/..." paths the CMS uses, blocking the admin from saving
	// a working display image. It should be a plain text input.
	if strings.Contains(body, `name="image_url"`) && strings.Contains(body, `type="url" name="image_url"`) {
		t.Errorf("image_url field is type=\"url\", which rejects relative upload paths like %q; "+
			"expected a text input so admins can save the CMS's relative image paths.\nbody:\n%s", relPath, body)
	}
}
