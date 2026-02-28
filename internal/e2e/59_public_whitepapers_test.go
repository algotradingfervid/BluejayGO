package e2e_test

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
)

func TestPublicWhitepapersList(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "Security",
		Slug: "security",
	})

	_, err := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Security Best Practices",
		Slug:           "security-best-practices",
		Description:    "Comprehensive security guide",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/security.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    1,
		CoverColorFrom: "#667eea",
		CoverColorTo:   "#764ba2",
		MetaDescription: sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create whitepaper: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/whitepapers", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepapers listing route not found")
	}
}

func TestPublicWhitepapersList_TopicFilter(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "Performance",
		Slug: "performance",
	})

	_, _ = queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Performance Optimization",
		Slug:           "performance-optimization",
		Description:    "Speed up your systems",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/performance.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    1,
		CoverColorFrom: "#f093fb",
		CoverColorTo:   "#f5576c",
		MetaDescription: sql.NullString{},
	})

	req := httptest.NewRequest(http.MethodGet, "/whitepapers?topic="+strconv.FormatInt(topic.ID, 10), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepapers topic filter route not found")
	}
}

func TestPublicWhitepapersList_InvalidTopic(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/whitepapers?topic=invalid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid topic, got %d", rec.Code)
	}
}

func TestPublicWhitepaperDetail(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "Cloud",
		Slug: "cloud",
	})

	_, err := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Cloud Migration Guide",
		Slug:           "cloud-migration-guide",
		Description:    "Step-by-step cloud migration",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/cloud.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    1,
		CoverColorFrom: "#4facfe",
		CoverColorTo:   "#00f2fe",
		MetaDescription: sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create whitepaper: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/whitepapers/cloud-migration-guide", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepaper detail route not found")
	}
}

func TestPublicWhitepaperDetail_NotFound(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/whitepapers/nonexistent", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestPublicWhitepaperDetail_WithLearningPoints(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "DevOps",
		Slug: "devops",
	})

	wp, _ := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "DevOps Practices",
		Slug:           "devops-practices",
		Description:    "Modern DevOps strategies",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/devops.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    1,
		CoverColorFrom: "#fa709a",
		CoverColorTo:   "#fee140",
		MetaDescription: sql.NullString{},
	})

	_, err := queries.CreateWhitepaperLearningPoint(ctx, sqlc.CreateWhitepaperLearningPointParams{
		WhitepaperID: wp.ID,
		PointText:    "Continuous integration strategies",
		DisplayOrder: 1,
	})
	if err != nil {
		t.Fatalf("create learning point: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/whitepapers/devops-practices", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepaper with learning points route not found")
	}
}

func TestPublicWhitepaperDetail_PreviewMode(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "Draft Topic",
		Slug: "draft-topic",
	})

	_, err := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Draft Whitepaper",
		Slug:           "draft-whitepaper",
		Description:    "Draft description",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/draft.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    0,
		CoverColorFrom: "#ccc",
		CoverColorTo:   "#ddd",
		MetaDescription: sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create draft whitepaper: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/whitepapers/draft-whitepaper", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("draft should return 404, got %d", rec.Code)
	}

	// Preview mode requires authentication + ?preview=true
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	req = httptest.NewRequest(http.MethodGet, "/whitepapers/draft-whitepaper?preview=true", nil)
	req.AddCookie(cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("preview mode should show draft whitepapers")
	}
}

func TestPublicWhitepaperDownload_Success(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "AI",
		Slug: "ai",
	})

	_, err := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "AI Fundamentals",
		Slug:           "ai-fundamentals",
		Description:    "Introduction to AI",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/ai.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    1,
		CoverColorFrom: "#30cfd0",
		CoverColorTo:   "#330867",
		MetaDescription: sql.NullString{},
	})
	if err != nil {
		t.Fatalf("create whitepaper: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/whitepapers/ai-fundamentals/download", strings.NewReader(url.Values{
		"name":              {"John Doe"},
		"email":             {"john@example.com"},
		"company":           {"Acme Corp"},
		"designation":       {"Engineer"},
		"marketing_consent": {"on"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepaper download route not found")
	}
}

func TestPublicWhitepaperDownload_MissingFields(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "Testing",
		Slug: "testing",
	})

	_, _ = queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Testing Guide",
		Slug:           "testing-guide",
		Description:    "Testing strategies",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/testing.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    1,
		CoverColorFrom: "#a8edea",
		CoverColorTo:   "#fed6e3",
		MetaDescription: sql.NullString{},
	})

	req := httptest.NewRequest(http.MethodPost, "/whitepapers/testing-guide/download", strings.NewReader(url.Values{
		"name": {""},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", rec.Code)
	}
}

func TestPublicWhitepaperDownload_OptionalFields(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "Automation",
		Slug: "automation",
	})

	_, _ = queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:          "Automation Handbook",
		Slug:           "automation-handbook",
		Description:    "Automation strategies",
		TopicID:        topic.ID,
		PdfFilePath:    "whitepapers/automation.pdf",
		FileSizeBytes:  0,
		PageCount:      sql.NullInt64{},
		PublishedDate:  "2024-01-01",
		IsPublished:    1,
		CoverColorFrom: "#d299c2",
		CoverColorTo:   "#fef9d7",
		MetaDescription: sql.NullString{},
	})

	req := httptest.NewRequest(http.MethodPost, "/whitepapers/automation-handbook/download", strings.NewReader(url.Values{
		"name":    {"Jane Doe"},
		"email":   {"jane@example.com"},
		"company": {"Tech Inc"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepaper download with optional fields route not found")
	}
}

func TestPublicWhitepaperDownload_NotFound(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/whitepapers/nonexistent/download", strings.NewReader(url.Values{
		"name":    {"Test User"},
		"email":   {"test@example.com"},
		"company": {"Test Co"},
	}.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent whitepaper, got %d", rec.Code)
	}
}

func TestPublicWhitepaperDownload_MetaDescription(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	ctx := context.Background()

	topic, _ := queries.CreateWhitepaperTopic(ctx, sqlc.CreateWhitepaperTopicParams{
		Name: "SEO",
		Slug: "seo",
	})

	_, err := queries.CreateWhitepaper(ctx, sqlc.CreateWhitepaperParams{
		Title:           "SEO Whitepaper",
		Slug:            "seo-whitepaper",
		Description:     "SEO optimization guide",
		TopicID:         topic.ID,
		PdfFilePath:     "whitepapers/seo.pdf",
		FileSizeBytes:   0,
		PageCount:       sql.NullInt64{},
		PublishedDate:   "2024-01-01",
		IsPublished:     1,
		CoverColorFrom:  "#000",
		CoverColorTo:    "#fff",
		MetaDescription: sql.NullString{String: "Custom meta description", Valid: true},
	})
	if err != nil {
		t.Fatalf("create whitepaper with meta: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/whitepapers/seo-whitepaper", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Error("whitepaper with meta route not found")
	}
}
