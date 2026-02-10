// Package admin provides HTTP handlers for the admin panel.
// This file contains handlers for the activity log viewing functionality.
package admin

import (
	// Standard library imports

	// log/slog: Structured logging for error and debug messages
	"log/slog"

	// math: Mathematical operations, used for pagination calculations (Ceil)
	"math"

	// net/http: HTTP status codes and standard HTTP types
	"net/http"

	// strconv: String conversion utilities, used for parsing query parameters
	"strconv"

	// Third-party imports

	// echo/v4: Web framework providing routing, context, and HTTP handling
	"github.com/labstack/echo/v4"

	// sqlc: Generated database query client for type-safe SQL operations
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

// activityPerPage defines the number of activity log entries displayed per page.
// This constant controls pagination throughout the activity log interface.
const activityPerPage = 50

// ActivityHandler handles HTTP requests for viewing and filtering activity logs.
// It provides read-only access to the audit trail, allowing administrators to
// review all actions taken within the CMS.
type ActivityHandler struct {
	// queries provides type-safe database access via sqlc-generated methods
	queries *sqlc.Queries

	// logger outputs structured error and debug messages
	logger *slog.Logger
}

// NewActivityHandler creates a new ActivityHandler with the provided dependencies.
// This constructor is typically called during application initialization to wire
// up the database queries and logger instances.
//
// Parameters:
//   - queries: sqlc.Queries instance for database operations
//   - logger: slog.Logger for structured logging
//
// Returns:
//   - *ActivityHandler: Initialized handler ready to process HTTP requests
func NewActivityHandler(queries *sqlc.Queries, logger *slog.Logger) *ActivityHandler {
	return &ActivityHandler{queries: queries, logger: logger}
}

// List displays a paginated, filterable view of all activity logs in the system.
//
// HTTP Method: GET
// Route: /admin/activity
// Returns: Full HTML page (not an HTMX fragment)
// Template: admin/pages/activity_log.html (uses admin-layout.html as base)
//
// Query Parameters:
//   - action: Filter by specific action type (e.g., "create", "update", "delete")
//   - search: Free-text search across user_email, resource_title, and description
//   - page: Current page number (defaults to 1 if missing or invalid)
//
// Behavior:
//   - This endpoint renders the complete activity log page, including filters and pagination
//   - It is NOT an HTMX fragment — it returns a full page wrapped in the admin layout
//   - Filtering is performed at the database level via sqlc-generated queries
//   - Pagination displays 50 entries per page (controlled by activityPerPage constant)
//   - When no results exist, displays "showing 0 to 0" instead of "1 to 0"
//
// Audit Trail Context:
//   - This handler provides READ-ONLY access to the activity log
//   - Activity logs are created by other handlers via the logActivity helper
//   - Logs track all CUD (Create, Update, Delete) operations across the CMS
//   - Each log entry contains: timestamp, user, action, resource type/ID/title, description
//
// Template Data:
//   - Title: Page title ("Activity Log")
//   - Logs: Array of activity log entries for current page
//   - Action: Current action filter value (for preserving filter state)
//   - Search: Current search term (for preserving filter state)
//   - HasFilters: Boolean indicating if any filters are active (for UI state)
//   - Page: Current page number (1-indexed)
//   - TotalPages: Total number of pages based on filtered results
//   - Pages: Array of page numbers for pagination UI
//   - Total: Total number of filtered log entries
//   - ShowFrom: First entry number being displayed (e.g., "51" on page 2)
//   - ShowTo: Last entry number being displayed (e.g., "100" on page 2)
//
// Error Handling:
//   - Database query failures return HTTP 500 and log error details
//   - Invalid page numbers are silently corrected to page 1
func (h *ActivityHandler) List(c echo.Context) error {
	// Extract request context for database operations with timeout/cancellation support
	ctx := c.Request().Context()

	// Parse query parameters for filtering and pagination
	// action: filters by specific action types (create/update/delete/publish/unpublish)
	action := c.QueryParam("action")

	// search: performs full-text search across user_email, resource_title, and description
	search := c.QueryParam("search")

	// page: current page number, defaults to 1 if missing or invalid
	pageStr := c.QueryParam("page")

	// Convert page parameter to integer, defaulting to 1 for invalid values
	// Ignoring error because we handle invalid values with the fallback
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	// Calculate database offset for LIMIT/OFFSET pagination
	// Example: page 1 = offset 0, page 2 = offset 50, page 3 = offset 100
	offset := int64((page - 1) * activityPerPage)

	// Query the database for activity logs matching the current filters
	// ListActivityLogs performs filtering, ordering, and pagination in a single query
	logs, err := h.queries.ListActivityLogs(ctx, sqlc.ListActivityLogsParams{
		FilterAction: action,        // Filter by action type (empty string = no filter)
		FilterSearch: search,        // Filter by search term (empty string = no filter)
		PageLimit:    activityPerPage, // Always fetch exactly 50 rows
		PageOffset:   offset,         // Skip rows from previous pages
	})
	if err != nil {
		// Log database errors with structured context for debugging
		h.logger.Error("failed to list activity logs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Get the total count of matching logs for pagination calculation
	// This query respects the same filters but doesn't apply LIMIT/OFFSET
	total, err := h.queries.CountActivityLogs(ctx, sqlc.CountActivityLogsParams{
		FilterAction: action,
		FilterSearch: search,
	})
	if err != nil {
		// Log database errors with structured context for debugging
		h.logger.Error("failed to count activity logs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Calculate total number of pages needed to display all results
	// math.Ceil ensures we round up (e.g., 51 entries = 2 pages, not 1.02 pages)
	totalPages := int(math.Ceil(float64(total) / float64(activityPerPage)))
	if totalPages < 1 {
		// Ensure at least 1 page exists even with 0 results (avoids divide-by-zero in UI)
		totalPages = 1
	}

	// Build an array of page numbers for the pagination UI
	// Example: totalPages=5 produces [1, 2, 3, 4, 5]
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	// Calculate the entry range being displayed (e.g., "Showing 51-100 of 237")
	// offset + 1 gives the first entry number (1-indexed, not 0-indexed)
	showFrom := offset + 1

	// showTo is the offset plus the number of logs actually returned
	// On the last page, this might be less than a full page
	showTo := offset + int64(len(logs))

	// Special case: when no results exist, show "0 to 0" instead of "1 to 0"
	if total == 0 {
		showFrom = 0
	}

	// Determine if any filters are active (used to show/hide "Clear Filters" button)
	hasFilters := action != "" || search != ""

	// Render the full activity log page template
	// This is NOT an HTMX fragment — it includes the full admin layout with sidebar
	// Template path: templates/admin/pages/activity_log.html
	// Base layout: templates/admin/layouts/admin-layout.html
	return c.Render(http.StatusOK, "admin/pages/activity_log.html", map[string]interface{}{
		"Title":      "Activity Log",     // Browser title and page heading
		"Logs":       logs,                // Array of activity log entries for current page
		"Action":     action,              // Current action filter (preserves form state)
		"Search":     search,              // Current search term (preserves form state)
		"HasFilters": hasFilters,          // Whether any filters are active (UI visibility)
		"Page":       page,                // Current page number (for pagination state)
		"TotalPages": totalPages,          // Total pages (for pagination limits)
		"Pages":      pages,               // Page number array (for pagination UI)
		"Total":      total,               // Total filtered results (for "X total entries" display)
		"ShowFrom":   showFrom,            // First entry number on page (for "Showing X-Y" display)
		"ShowTo":     showTo,              // Last entry number on page (for "Showing X-Y" display)
	})
}
