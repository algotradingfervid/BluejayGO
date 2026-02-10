// Package public provides HTTP handlers for public-facing website features.
// This file implements preview functionality for viewing draft/unpublished content.
package public

import (
	"github.com/labstack/echo/v4"                                // Echo web framework for HTTP request/response handling
	"github.com/narendhupati/bluejay-cms/internal/middleware" // Session middleware for authentication checking
)

// isPreviewRequest checks if the request is a preview request with a valid admin session.
//
// Preview Mode:
//   - Allows authenticated admin users to view draft/unpublished content
//   - Triggered by ?preview=true query parameter
//   - Requires active admin session (user must be logged in)
//   - Used across all public content handlers (blog posts, products, case studies, etc.)
//
// Security:
//   - Only authenticated users can preview content (prevents public access to drafts)
//   - Session validation ensures preview parameter can't be exploited
//   - UserID > 0 check ensures session contains valid user data
//
// Usage Pattern:
//   In public content handlers:
//     if isPreviewRequest(c) {
//         // Fetch content regardless of status (include drafts)
//     } else {
//         // Only fetch published content
//     }
//
// Query Parameter:
//   - preview=true: Enable preview mode
//   - Any other value or omitted: Normal public mode
//
// Session Requirements:
//   - Must have valid session object in Echo context
//   - Session must contain UserID > 0 (authenticated admin user)
//   - Session created by middleware.SessionMiddleware
//
// Return Values:
//   - true: Valid preview request (authenticated admin with preview param)
//   - false: Normal request (no preview param OR not authenticated)
func isPreviewRequest(c echo.Context) bool {
	// Check for preview query parameter
	// Must be exactly "true" (case-sensitive)
	if c.QueryParam("preview") != "true" {
		return false
	}

	// Validate admin session exists and contains authenticated user
	// Session object is placed in context by middleware.SessionMiddleware
	sess, ok := c.Get("session").(*middleware.Session)
	if !ok || sess == nil || sess.UserID == 0 {
		// No valid session - deny preview access
		// This prevents unauthorized users from viewing drafts via ?preview=true
		return false
	}

	// Valid preview request: authenticated admin with preview parameter
	return true
}
