package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func TestCertificationsCRUD_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	req := httptest.NewRequest(http.MethodPost, "/admin/about/certifications", strings.NewReader(url.Values{
		"name":          {"ISO 9001:2015"},
		"abbreviation":  {"ISO 9001"},
		"description":   {"Quality management systems certification"},
		"icon":          {"workspace_premium"},
		"display_order": {"1"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create: expected 303, got %d", rec.Code)
	}

	certifications, _ := queries.ListCertifications(ctx)
	if len(certifications) != 1 {
		t.Fatalf("expected 1 certification, got %d", len(certifications))
	}
	if certifications[0].Name != "ISO 9001:2015" {
		t.Errorf("expected 'ISO 9001:2015', got %q", certifications[0].Name)
	}
	if certifications[0].Abbreviation != "ISO 9001" {
		t.Errorf("expected 'ISO 9001', got %q", certifications[0].Abbreviation)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/about/certifications/%d", certifications[0].ID), nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("delete: expected 200, got %d", rec.Code)
	}

	_, err := queries.GetCertification(ctx, certifications[0].ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected ErrNoRows after delete, got %v", err)
	}
}

func TestCertificationsList_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	for i := 1; i <= 3; i++ {
		queries.CreateCertification(ctx, sqlc.CreateCertificationParams{
			Name:         fmt.Sprintf("Certification %d", i),
			Abbreviation: fmt.Sprintf("CERT-%d", i),
			Description:  sql.NullString{String: fmt.Sprintf("Description %d", i), Valid: true},
			Icon:         sql.NullString{String: "badge", Valid: true},
			DisplayOrder: int64(i),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/about/certifications", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("certifications list route not found")
	}
}

func TestCertificationEdit_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()

	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)
	ctx := context.Background()

	cert, _ := queries.CreateCertification(ctx, sqlc.CreateCertificationParams{
		Name:         "Original Name",
		Abbreviation: "ORIG",
		Description:  sql.NullString{String: "Original description", Valid: true},
		Icon:         sql.NullString{String: "badge", Valid: true},
		DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/about/certifications/%d", cert.ID), strings.NewReader(url.Values{
		"name":          {"Updated Name"},
		"abbreviation":  {"UPDT"},
		"description":   {"Updated description"},
		"icon":          {"verified"},
		"display_order": {"2"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("update: expected 303, got %d", rec.Code)
	}

	updated, _ := queries.GetCertification(ctx, cert.ID)
	if updated.Name != "Updated Name" {
		t.Errorf("expected 'Updated Name', got %q", updated.Name)
	}
	if updated.Abbreviation != "UPDT" {
		t.Errorf("expected 'UPDT', got %q", updated.Abbreviation)
	}
}
