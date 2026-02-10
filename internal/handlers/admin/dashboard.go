package admin

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
)

type DashboardHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewDashboardHandler(queries *sqlc.Queries, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{
		queries: queries,
		logger:  logger,
	}
}

type DashboardData struct {
	Title                  string
	ActiveNav              string
	DisplayName            string
	Email                  string
	Role                   string
	PublishedProducts      int64
	PublishedBlogPosts     int64
	ContactSubmissions     int64
	NewContactSubmissions  int64
	TotalPartners          int64
	DraftProducts          int64
	DraftBlogPosts         int64
}

func (h *DashboardHandler) ShowDashboard(c echo.Context) error {
	sess := c.Get("session").(*customMiddleware.Session)
	ctx := c.Request().Context()

	data := DashboardData{
		Title:       "Dashboard",
		ActiveNav:   "dashboard",
		DisplayName: sess.DisplayName,
		Email:       sess.Email,
		Role:        sess.Role,
	}

	// Fetch all counts, logging errors but not failing
	if count, err := h.queries.CountProducts(ctx); err == nil {
		data.PublishedProducts = count
	} else {
		h.logger.Error("dashboard: count products", "error", err)
	}

	if count, err := h.queries.CountPublishedPosts(ctx); err == nil {
		data.PublishedBlogPosts = count
	} else {
		h.logger.Error("dashboard: count blog posts", "error", err)
	}

	if count, err := h.queries.CountContactSubmissions(ctx); err == nil {
		data.ContactSubmissions = count
	} else {
		h.logger.Error("dashboard: count contact submissions", "error", err)
	}

	if count, err := h.queries.CountNewContactSubmissions(ctx); err == nil {
		data.NewContactSubmissions = count
	} else {
		h.logger.Error("dashboard: count new contact submissions", "error", err)
	}

	if count, err := h.queries.CountPartners(ctx); err == nil {
		data.TotalPartners = count
	} else {
		h.logger.Error("dashboard: count partners", "error", err)
	}

	if count, err := h.queries.CountDraftProducts(ctx); err == nil {
		data.DraftProducts = count
	} else {
		h.logger.Error("dashboard: count draft products", "error", err)
	}

	if count, err := h.queries.CountDraftBlogPosts(ctx); err == nil {
		data.DraftBlogPosts = count
	} else {
		h.logger.Error("dashboard: count draft blog posts", "error", err)
	}

	return c.Render(http.StatusOK, "admin/pages/dashboard.html", data)
}
