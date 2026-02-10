package admin

import (
	// Standard library imports
	"log/slog" // Structured logging for error tracking and debugging
	"net/http" // HTTP status codes and request/response handling

	// Third-party framework
	"github.com/labstack/echo/v4" // Echo web framework for HTTP routing and context management

	// Internal dependencies
	"github.com/narendhupati/bluejay-cms/db/sqlc"           // sqlc-generated database queries and models
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware" // Session management and authentication middleware
)

// DashboardHandler handles all dashboard-related HTTP requests.
// Responsible for aggregating statistics from multiple content sections
// and rendering the admin dashboard overview page.
type DashboardHandler struct {
	queries *sqlc.Queries // Database query interface for fetching dashboard statistics
	logger  *slog.Logger  // Structured logger for error and activity logging
}

// NewDashboardHandler creates and initializes a new DashboardHandler instance.
// Dependencies are injected to support database access and logging.
func NewDashboardHandler(queries *sqlc.Queries, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{
		queries: queries,
		logger:  logger,
	}
}

// DashboardData contains all data needed to render the admin dashboard page.
// Includes user session information and aggregated statistics from all content sections.
// This struct is passed to the template renderer for display.
type DashboardData struct {
	Title                  string // Page title ("Dashboard")
	ActiveNav              string // Active navigation item identifier ("dashboard")
	DisplayName            string // Current user's display name from session
	Email                  string // Current user's email from session
	Role                   string // Current user's role (admin, editor, etc.) from session
	PublishedProducts      int64  // Count of published products
	PublishedBlogPosts     int64  // Count of published blog posts
	ContactSubmissions     int64  // Total count of contact form submissions
	NewContactSubmissions  int64  // Count of unread/new contact submissions
	TotalPartners          int64  // Count of partner organizations
	DraftProducts          int64  // Count of unpublished/draft products
	DraftBlogPosts         int64  // Count of unpublished/draft blog posts
}

// ShowDashboard renders the admin dashboard overview page.
//
// HTTP Method: GET
// Route: /admin/dashboard
// Template: templates/admin/pages/dashboard.html (with admin-layout wrapper)
// HTMX: Returns full HTML page (not a fragment)
//
// Aggregates statistics from multiple database tables (products, blog posts,
// contact submissions, partners) and displays them in a single dashboard view.
// Uses graceful degradation: if any count query fails, it logs the error but
// continues rendering the page with zero values for failed queries.
//
// Authentication: Requires valid session (enforced by middleware)
// Session data: Retrieves user DisplayName, Email, and Role for display
func (h *DashboardHandler) ShowDashboard(c echo.Context) error {
	// Extract authenticated session from Echo context (set by auth middleware)
	sess := c.Get("session").(*customMiddleware.Session)
	ctx := c.Request().Context()

	// Initialize template data with session information
	data := DashboardData{
		Title:       "Dashboard",
		ActiveNav:   "dashboard", // Highlights "Dashboard" in sidebar navigation
		DisplayName: sess.DisplayName,
		Email:       sess.Email,
		Role:        sess.Role,
	}

	// Fetch all statistics counts from database
	// Strategy: Log errors but don't fail the entire page if one query fails
	// This ensures the dashboard remains accessible even with partial data

	// Count published products (status = 'published')
	if count, err := h.queries.CountProducts(ctx); err == nil {
		data.PublishedProducts = count
	} else {
		h.logger.Error("dashboard: count products", "error", err)
	}

	// Count published blog posts (status = 'published')
	if count, err := h.queries.CountPublishedPosts(ctx); err == nil {
		data.PublishedBlogPosts = count
	} else {
		h.logger.Error("dashboard: count blog posts", "error", err)
	}

	// Count all contact form submissions (total submissions ever received)
	if count, err := h.queries.CountContactSubmissions(ctx); err == nil {
		data.ContactSubmissions = count
	} else {
		h.logger.Error("dashboard: count contact submissions", "error", err)
	}

	// Count new/unread contact submissions (status = 'new')
	if count, err := h.queries.CountNewContactSubmissions(ctx); err == nil {
		data.NewContactSubmissions = count
	} else {
		h.logger.Error("dashboard: count new contact submissions", "error", err)
	}

	// Count total partner organizations
	if count, err := h.queries.CountPartners(ctx); err == nil {
		data.TotalPartners = count
	} else {
		h.logger.Error("dashboard: count partners", "error", err)
	}

	// Count draft products (status = 'draft')
	if count, err := h.queries.CountDraftProducts(ctx); err == nil {
		data.DraftProducts = count
	} else {
		h.logger.Error("dashboard: count draft products", "error", err)
	}

	// Count draft blog posts (status = 'draft')
	if count, err := h.queries.CountDraftBlogPosts(ctx); err == nil {
		data.DraftBlogPosts = count
	} else {
		h.logger.Error("dashboard: count draft blog posts", "error", err)
	}

	// Render the dashboard template with admin layout wrapper
	// Template path: templates/admin/pages/dashboard.html
	// Layout: Uses {{template "admin-layout" .}} for consistent admin UI
	return c.Render(http.StatusOK, "admin/pages/dashboard.html", data)
}
