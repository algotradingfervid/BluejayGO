// Package admin provides HTTP handlers for the admin panel.
// This file contains utilities for logging user activity across all admin handlers.
package admin

import (
	// Standard library imports

	// fmt: String formatting for building activity descriptions
	"fmt"

	// Third-party and internal imports

	// echo/v4: Web framework providing context for request handling
	"github.com/labstack/echo/v4"

	// customMiddleware: Internal middleware providing session management
	// Aliased to avoid conflict with echo.Context methods
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"

	// services: Business logic layer containing ActivityLogService
	"github.com/narendhupati/bluejay-cms/internal/services"
)

// activityLog is a package-level reference to the activity log service.
// This singleton pattern allows all handlers in the admin package to log activities
// without explicitly passing the service through each handler's constructor.
//
// Lifecycle:
//   - Initialized via SetActivityLogService during application startup
//   - Remains available for the lifetime of the application
//   - Nil checks in logActivity ensure graceful handling if not initialized
var activityLog *services.ActivityLogService

// SetActivityLogService sets the package-level activity log service.
// This function must be called during application initialization, after the
// ActivityLogService has been instantiated with its database connection.
//
// Usage Pattern:
//   - Called once in main.go or server initialization code
//   - Must be called before any handlers process requests
//   - Enables all admin handlers to log activities via the logActivity helper
//
// Parameters:
//   - svc: Initialized ActivityLogService instance with database access
//
// Example:
//   activityLogService := services.NewActivityLogService(queries, logger)
//   admin.SetActivityLogService(activityLogService)
func SetActivityLogService(svc *services.ActivityLogService) {
	activityLog = svc
}

// getUserID extracts the user ID from the echo context session.
// This function is used by logActivity to associate activity log entries with
// the authenticated user who performed the action.
//
// Session Extraction:
//   - Retrieves the "session" key from the echo.Context (set by auth middleware)
//   - Type asserts to *customMiddleware.Session to access UserID field
//   - Returns 0 if session is missing or type assertion fails
//
// Middleware Dependency:
//   - Requires auth middleware to have run before this handler
//   - Auth middleware stores Session object in context on successful authentication
//   - Session contains UserID, UserEmail, and other authentication metadata
//
// Parameters:
//   - c: echo.Context containing the HTTP request and response
//
// Returns:
//   - int64: The authenticated user's ID, or 0 if session is unavailable
//
// Note: Returning 0 for missing sessions is safe because activity logging is
// non-critical and should not block request processing. The ActivityLogService
// handles userId=0 gracefully by storing it as-is (useful for debugging auth issues).
func getUserID(c echo.Context) int64 {
	// Attempt to retrieve the session object from the context
	// The "session" key is populated by the auth middleware for authenticated requests
	if sess, ok := c.Get("session").(*customMiddleware.Session); ok {
		// Successfully retrieved and type-asserted the session
		return sess.UserID
	}

	// Session not found or type assertion failed (e.g., unauthenticated request)
	// Return 0 to indicate "no user" rather than causing a panic
	return 0
}

// logActivity is a convenience helper used by all admin handlers to record user actions.
// This function provides a consistent interface for creating audit trail entries across
// the entire CMS, tracking all create, update, delete, and other significant operations.
//
// Audit Trail Purpose:
//   - Creates a permanent, searchable record of all administrative actions
//   - Enables compliance with audit requirements and security reviews
//   - Helps debug issues by showing what changed, when, and by whom
//   - Provides accountability for sensitive operations (deletions, publications, etc.)
//
// Activity Log Structure:
// Each log entry captures:
//   - Timestamp: When the action occurred (set by database)
//   - User ID: Who performed the action (extracted from session)
//   - Action: Type of operation (e.g., "create", "update", "delete", "publish")
//   - Resource Type: What kind of entity was affected (e.g., "product", "category", "page")
//   - Resource ID: Database ID of the specific entity
//   - Resource Title: Human-readable name/title of the entity (for quick reference)
//   - Description: Detailed explanation of what changed (formatted from descFmt and args)
//
// Usage Pattern in Handlers:
//   // After creating a new product:
//   logActivity(c, "create", "product", product.ID, product.Title,
//       "Created product with SKU %s and price $%.2f", product.SKU, product.Price)
//
//   // After updating a category:
//   logActivity(c, "update", "category", category.ID, category.Name,
//       "Updated category name from %q to %q", oldName, newName)
//
//   // After deleting a page:
//   logActivity(c, "delete", "page", page.ID, page.Title,
//       "Permanently deleted page and all associated content")
//
// Common Action Types:
//   - "create": New entity created
//   - "update": Existing entity modified
//   - "delete": Entity removed (soft or hard delete)
//   - "publish": Content made publicly visible
//   - "unpublish": Content hidden from public view
//   - "restore": Previously deleted entity recovered
//
// Common Resource Types:
//   - "product": E-commerce products
//   - "category": Product categories
//   - "page": Static content pages
//   - "media": Uploaded images/files
//   - "user": Admin users (for user management operations)
//
// Parameters:
//   - c: echo.Context for extracting user session and request context
//   - action: Type of operation performed (e.g., "create", "update", "delete")
//   - resourceType: Category of entity affected (e.g., "product", "category", "page")
//   - resourceID: Database ID of the affected entity
//   - resourceTitle: Human-readable name/title of the entity
//   - descFmt: Format string for the description (follows fmt.Sprintf conventions)
//   - args: Variadic arguments to be interpolated into descFmt
//
// Nil Safety:
//   - Checks if activityLog is initialized before attempting to log
//   - If activityLog is nil, function returns silently without error
//   - This prevents panics if SetActivityLogService was not called during initialization
//   - Activity logging is considered non-critical and should not block request processing
//
// Database Interaction:
//   - Logs are written asynchronously to avoid blocking the HTTP response
//   - Database errors are logged but do not affect handler return values
//   - Uses the request context for timeout/cancellation support
//
// HTMX Behavior Note:
//   - This helper is called by handlers regardless of whether they return full pages or fragments
//   - Activity logs track the underlying action, not the rendering format
//   - Both full page renders and HTMX fragment updates log the same activity data
func logActivity(c echo.Context, action, resourceType string, resourceID int64, resourceTitle, descFmt string, args ...interface{}) {
	// Nil check: ensure activityLog service was initialized during app startup
	// If not initialized, silently skip logging to avoid panics
	if activityLog != nil {
		// Format the description using fmt.Sprintf with the provided format string and arguments
		// This allows handlers to build descriptive audit messages with dynamic data
		// Example: "Updated product title from %q to %q" with args ["Old Title", "New Title"]
		//   becomes: "Updated product title from \"Old Title\" to \"New Title\""
		desc := fmt.Sprintf(descFmt, args...)

		// Call the ActivityLogService to persist the log entry to the database
		// Parameters passed through:
		//   - context: request context for timeout/cancellation
		//   - userID: extracted from session via getUserID helper
		//   - action: operation type ("create", "update", "delete", etc.)
		//   - resourceType: entity category ("product", "category", "page", etc.)
		//   - resourceID: database ID of affected entity
		//   - resourceTitle: human-readable entity name
		//   - desc: formatted description of what changed
		activityLog.Log(c.Request().Context(), getUserID(c), action, resourceType, resourceID, resourceTitle, desc)
	}
	// If activityLog is nil, function returns here without logging
	// This is intentional — activity logging is non-critical infrastructure
}
