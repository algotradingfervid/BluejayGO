// Package middleware provides HTTP middleware components for the Bluejay CMS application.
// This package includes authentication, authorization, session management, security headers,
// CSRF protection, rate limiting, caching strategies, request logging, panic recovery, and
// application settings injection. All middleware functions follow the Echo v4 middleware
// pattern and can be chained together in the routing configuration.
package middleware

import (
	// net/http provides HTTP status codes and standard HTTP constants used for
	// redirects, authorization failures, and other HTTP response scenarios.
	"net/http"

	// github.com/labstack/echo/v4 is the Echo web framework. It provides the core
	// middleware interface (echo.MiddlewareFunc) and context object (echo.Context)
	// used throughout all middleware implementations.
	"github.com/labstack/echo/v4"
)

// RequireAuth returns an Echo middleware that enforces authentication for protected routes.
// This middleware checks if a valid user session exists by retrieving the session from the
// Echo context (populated by SessionMiddleware). If no valid session is found, or if the
// session exists but UserID is 0 (indicating an unauthenticated user), the middleware
// redirects the user to the admin login page.
//
// This is the primary authentication gate for admin panel routes. It should be applied
// to all routes that require a logged-in user, regardless of role.
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that can be used in Echo route groups
//     or individual routes to enforce authentication requirements.
//
// Example usage:
//
//	admin := e.Group("/admin", middleware.RequireAuth())
//
// Security considerations:
//   - Relies on SessionMiddleware being executed first in the middleware chain
//   - Uses 303 See Other redirect to ensure POST requests are converted to GET
//   - Does not perform role-based authorization (see RequireRole for that)
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve the session from the Echo context. The session is populated
			// by SessionMiddleware earlier in the middleware chain.
			sess, ok := c.Get("session").(*Session)

			// Check if session retrieval was successful and if the user is authenticated.
			// A UserID of 0 indicates an unauthenticated session (no user logged in).
			if !ok || sess.UserID == 0 {
				// Redirect unauthenticated users to the login page using HTTP 303 See Other.
				// This status code ensures that the redirect is always performed with GET,
				// even if the original request was POST, preventing form resubmission issues.
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}

			// User is authenticated, proceed to the next handler in the chain
			return next(c)
		}
	}
}

// RequireRole returns an Echo middleware that enforces role-based access control (RBAC)
// for protected routes. This middleware first checks authentication (similar to RequireAuth),
// then verifies that the authenticated user has one of the specified roles.
//
// The middleware checks the user's role against a whitelist of allowed roles. If the user's
// role matches any of the provided roles, access is granted. Otherwise, a 403 Forbidden
// error is returned.
//
// Parameters:
//   - roles: Variable number of role strings that are allowed to access the route.
//     Common roles include "admin", "editor", "viewer", etc. At least one role
//     should match for the user to gain access.
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that enforces role-based authorization.
//
// Example usage:
//
//	// Only admins can access this route
//	e.DELETE("/admin/users/:id", handler.DeleteUser, middleware.RequireRole("admin"))
//
//	// Admins and editors can access this route group
//	content := admin.Group("/content", middleware.RequireRole("admin", "editor"))
//
// Security considerations:
//   - Performs authentication check first (redirects to login if not authenticated)
//   - Role comparison is case-sensitive
//   - Returns 403 Forbidden (not 401 Unauthorized) for authenticated but unauthorized users
//   - Relies on SessionMiddleware to populate the session with accurate role information
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve the session from the Echo context
			sess, ok := c.Get("session").(*Session)

			// First check if the user is authenticated at all. If not, redirect to login.
			// This ensures that unauthenticated users don't receive a "Forbidden" error,
			// which would be confusing (they should see the login page instead).
			if !ok || sess.UserID == 0 {
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}

			// Check if the user's role matches any of the allowed roles.
			// This is a whitelist approach: if any match is found, access is granted.
			for _, role := range roles {
				if sess.Role == role {
					// User has an authorized role, proceed to the next handler
					return next(c)
				}
			}

			// User is authenticated but does not have any of the required roles.
			// Return a 403 Forbidden error with a clear message. Note that we use
			// 403 (not 401) because the user IS authenticated, just not authorized.
			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
		}
	}
}
