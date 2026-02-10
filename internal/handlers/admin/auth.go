package admin

import (
	// Standard library imports
	"log/slog"  // Structured logging for authentication events and security auditing
	"net/http"  // HTTP status codes and request/response handling

	// Third-party framework
	"github.com/labstack/echo/v4" // Echo web framework for HTTP routing and context management

	// Internal dependencies
	"github.com/narendhupati/bluejay-cms/db/sqlc"           // sqlc-generated database queries and models
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware" // Session management and authentication middleware

	// Security
	"golang.org/x/crypto/bcrypt" // Password hashing and verification using bcrypt algorithm
)

// AuthHandler handles authentication-related HTTP requests.
// Manages user login, logout, and session lifecycle for admin panel access.
// All authentication events are logged for security auditing.
type AuthHandler struct {
	queries *sqlc.Queries // Database query interface for user credential verification
	logger  *slog.Logger  // Structured logger for security event tracking
}

// NewAuthHandler creates and initializes a new AuthHandler instance.
// Dependencies are injected to support database access and security logging.
func NewAuthHandler(queries *sqlc.Queries, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		queries: queries,
		logger:  logger,
	}
}

// ShowLoginPage renders the admin login form page.
//
// HTTP Method: GET
// Route: /admin/login
// Template: templates/admin/pages/login.html (standalone, no admin layout wrapper)
// HTMX: Returns full HTML page (not a fragment)
//
// Behavior:
// - If user is already authenticated (valid session), redirects to /admin/dashboard
// - Otherwise, displays login form with optional error message from query parameter
//
// Query Parameters:
// - error: Optional error code (e.g., "invalid_credentials", "missing_fields", "session_error")
//   Used to display error messages after failed login attempts
//
// Authentication: None required (public access)
func (h *AuthHandler) ShowLoginPage(c echo.Context) error {
	// Check if user is already logged in by inspecting session
	sess, _ := c.Get("session").(*customMiddleware.Session)
	if sess != nil && sess.UserID > 0 {
		// User already authenticated, redirect to dashboard to prevent re-login
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}

	// Render login form with potential error message from previous failed attempt
	return c.Render(http.StatusOK, "admin/pages/login.html", map[string]interface{}{
		"Title": "Admin Login",
		"Error": c.QueryParam("error"), // Error codes: invalid_credentials, missing_fields, session_error
	})
}

// LoginSubmit processes the login form submission and authenticates the user.
//
// HTTP Method: POST
// Route: /admin/login
// Form Fields: email, password
// HTMX: Not used - standard form POST with redirects
//
// Authentication Flow:
// 1. Validate that email and password fields are not empty
// 2. Query database for user with matching email
// 3. Verify password using bcrypt hash comparison
// 4. Update user's last_login timestamp in database
// 5. Create authenticated session with user data
// 6. Log successful authentication event
// 7. Redirect to admin dashboard
//
// Security Features:
// - Passwords are verified using bcrypt (constant-time comparison)
// - Failed attempts are logged with email (for security monitoring)
// - Generic "invalid_credentials" error prevents user enumeration
// - Session is saved server-side with secure cookie
//
// Error Handling:
// - All authentication failures redirect back to login with error query param
// - Errors are logged for security auditing but user sees generic message
//
// Authentication: None required (this is the login endpoint)
func (h *AuthHandler) LoginSubmit(c echo.Context) error {
	// Extract form values from POST request
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Validate required fields are present
	if email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=missing_fields")
	}

	// Query database for user with matching email
	user, err := h.queries.GetAdminUserByEmail(c.Request().Context(), email)
	if err != nil {
		// User not found - log for security monitoring but show generic error
		h.logger.Warn("login attempt failed", "email", email, "error", "user_not_found")
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=invalid_credentials")
	}

	// Verify password using bcrypt constant-time comparison
	// CompareHashAndPassword prevents timing attacks
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// Password mismatch - log for security monitoring but show generic error
		h.logger.Warn("login attempt failed", "email", email, "error", "invalid_password")
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=invalid_credentials")
	}

	// Update user's last_login timestamp for activity tracking
	// Non-critical operation: log error but don't fail authentication
	if err := h.queries.UpdateLastLogin(c.Request().Context(), user.ID); err != nil {
		h.logger.Error("failed to update last login", "user_id", user.ID, "error", err)
	}

	// Retrieve session and populate with authenticated user data
	sess := c.Get("session").(*customMiddleware.Session)
	sess.UserID = user.ID
	sess.Email = user.Email
	sess.DisplayName = user.DisplayName
	sess.Role = user.Role

	// Persist session to cookie (server-side session storage)
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.logger.Error("failed to save session", "error", err)
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=session_error")
	}

	// Log successful authentication for security audit trail
	h.logger.Info("user logged in", "user_id", user.ID, "email", user.Email)
	// Log activity to activity_log table for admin dashboard
	logActivity(c, "login", "system", user.ID, user.DisplayName, "User '%s' logged in", user.DisplayName)

	// Redirect to admin dashboard on successful authentication
	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

// Logout terminates the user's session and logs them out of the admin panel.
//
// HTTP Method: POST (or GET if configured)
// Route: /admin/logout
// HTMX: Not used - standard redirect flow
//
// Logout Process:
// 1. Retrieve current session from context
// 2. Clear all session data (UserID, Email, DisplayName, Role)
// 3. Set session MaxAge to -1 (instructs browser to delete cookie)
// 4. Save session changes to persist deletion
// 5. Log logout event for audit trail
// 6. Redirect to login page
//
// Security Notes:
// - Session is completely destroyed, not just cleared
// - MaxAge = -1 ensures browser deletes the session cookie
// - Logout events are logged for security monitoring
// - Even if session save fails, user is redirected to login (fail-safe)
//
// Authentication: Requires valid session (enforced by middleware)
func (h *AuthHandler) Logout(c echo.Context) error {
	// Retrieve existing session from context
	sess := c.Get("session").(*customMiddleware.Session)

	// Clear all user data from session
	sess.UserID = 0
	sess.Email = ""
	sess.DisplayName = ""
	sess.Role = ""
	// Set MaxAge to -1 to instruct browser to delete the session cookie
	sess.Options.MaxAge = -1

	// Persist session destruction (saves empty session and triggers cookie deletion)
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		// Log error but continue logout process (fail-safe approach)
		h.logger.Error("failed to destroy session", "error", err)
	}

	// Log logout event for security audit trail
	h.logger.Info("user logged out")
	// Log activity to activity_log table (user_id=0 since session is cleared)
	logActivity(c, "logout", "system", 0, "", "User logged out")

	// Redirect to login page (user must re-authenticate to access admin panel)
	return c.Redirect(http.StatusSeeOther, "/admin/login")
}
