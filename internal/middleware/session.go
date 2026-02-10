package middleware

import (
	// net/http provides HTTP status codes, cookie management, and standard HTTP constants.
	// Used here for SameSite cookie attributes and request/response handling in session operations.
	"net/http"

	// github.com/gorilla/sessions is a secure session management library that supports
	// cookie-based and filesystem-based session storage. It provides:
	// - Encrypted cookie storage (using HMAC for integrity and optionally AES for encryption)
	// - Session data serialization (gob encoding)
	// - Configurable security options (HttpOnly, Secure, SameSite, MaxAge)
	// Used here for secure cookie-based session storage of user authentication state.
	"github.com/gorilla/sessions"

	// github.com/labstack/echo/v4 is the Echo web framework, providing middleware
	// interfaces, context objects for storing session data, and request/response utilities.
	"github.com/labstack/echo/v4"
)

// Session is a wrapper around gorilla/sessions.Session that provides type-safe access to
// user authentication data stored in the session. It embeds the underlying Session object
// and adds strongly-typed fields for common session attributes.
//
// This struct serves two purposes:
//  1. Type safety: Provides structured access to session data instead of working with
//     map[interface{}]interface{} directly, reducing the risk of type assertion errors.
//  2. Convenience: Commonly accessed session fields are promoted to top-level struct fields
//     for easier access in handlers and middleware.
//
// Fields:
//   - Session: Embedded gorilla Session object providing raw access to session data, methods
//     like Save(), and access to the underlying Values map for storing arbitrary data.
//   - UserID: The authenticated user's database ID. A value of 0 indicates an unauthenticated
//     session (guest/anonymous user). This is the primary authentication indicator.
//   - Email: The authenticated user's email address. Used for display purposes and as a
//     unique identifier in the application.
//   - DisplayName: The user's display name (could be full name, username, etc.). Used in
//     the UI to personalize the experience (e.g., "Welcome, John").
//   - Role: The user's role for authorization purposes (e.g., "admin", "editor", "viewer").
//     Used by RequireRole() middleware and in templates to show/hide features.
//
// The struct is stored in the Echo context under the key "session" and can be accessed
// in handlers via: sess := c.Get("session").(*middleware.Session)
type Session struct {
	*sessions.Session            // Embedded gorilla Session for low-level access
	UserID            int64      // User database ID (0 = unauthenticated)
	Email             string     // User email address
	DisplayName       string     // User display name
	Role              string     // User role (admin, editor, etc.)
}

// SessionStore is a package-level variable holding the global CookieStore instance used
// for session management across the application. It's initialized once at application startup
// by InitSessionStore() and used by SessionMiddleware() to retrieve/create sessions.
//
// Gorilla's CookieStore encrypts session data using HMAC-SHA256 for integrity verification
// and optionally AES for encryption (if a second secret key is provided). The secret key
// should be:
//   - At least 32 bytes for HMAC-SHA256 security
//   - Randomly generated (not a hardcoded string)
//   - Stored securely (environment variable, secrets manager, etc.)
//   - Rotated periodically for security
//
// Note: This is a global variable for simplicity, but could be refactored to be injected
// as a dependency for better testability and to support multiple session stores.
var SessionStore *sessions.CookieStore

// InitSessionStore initializes the global session store with the provided secret key and
// configures cookie security options. This function MUST be called once at application
// startup before any requests are processed.
//
// The session store uses cookie-based storage where session data is serialized, signed
// (and optionally encrypted), and stored directly in the client's browser cookie. This
// eliminates the need for server-side session storage (Redis, database, etc.) but means:
//   - Session data size is limited (typically 4KB per cookie)
//   - Session data is sent with every request (bandwidth overhead)
//   - Session data cannot be invalidated server-side (revocation requires token-based auth)
//
// Parameters:
//   - secret: The secret key used for HMAC signing of session cookies. Should be at least
//     32 bytes of random data (64 hex characters). If you need encryption in addition to
//     signing, pass two keys: []byte(secret), []byte(encryptionKey).
//
// Cookie options configured:
//   - Path: "/" (session cookie is sent for all paths on the domain)
//   - MaxAge: 604800 seconds (7 days) - session expires after 1 week of inactivity
//   - HttpOnly: true (cookie cannot be accessed via JavaScript, mitigating XSS attacks)
//   - Secure: false (allows HTTP during development; MUST be true in production for HTTPS-only)
//   - SameSite: Lax (cookie is sent for same-site requests and top-level navigation, but not
//     cross-site subrequests like images/iframes, providing CSRF protection)
//
// Example usage:
//
//	secret := os.Getenv("SESSION_SECRET") // Load from environment
//	if secret == "" {
//	    log.Fatal("SESSION_SECRET environment variable required")
//	}
//	middleware.InitSessionStore(secret)
//
// Production security checklist:
//   - Set Secure: true (require HTTPS)
//   - Use a cryptographically random 32+ byte secret
//   - Store secret in environment variable or secrets manager (not in code)
//   - Consider setting SameSite: Strict for maximum CSRF protection (may break some workflows)
//   - Consider shortening MaxAge for sensitive applications (e.g., 1 hour instead of 7 days)
//   - Implement session rotation on privilege escalation (e.g., after login)
func InitSessionStore(secret string) {
	// Create a new CookieStore with the provided secret key for HMAC signing.
	// gorilla/sessions uses gob encoding to serialize session data, then signs it
	// with HMAC-SHA256 using the secret key. This prevents tampering but does not
	// encrypt the data (clients can decode and read it, but cannot modify it without
	// detection). To add encryption, pass a second key: []byte(encryptionSecret).
	SessionStore = sessions.NewCookieStore([]byte(secret))

	// Configure cookie options for security and usability
	SessionStore.Options = &sessions.Options{
		// Path: "/" means the session cookie is sent for all paths under the domain.
		// This allows the session to be accessed across the entire application.
		// Could be restricted to "/admin" if sessions are only needed for admin routes.
		Path: "/",

		// MaxAge: 604800 seconds = 7 days. The session cookie expires after 7 days
		// of inactivity. After expiration, the user must log in again. This balances
		// security (shorter = more secure) with UX (longer = fewer logins).
		// Note: Changing this requires existing users to log in again.
		MaxAge: 86400 * 7, // 7 days in seconds

		// HttpOnly: true means the cookie cannot be accessed via JavaScript (document.cookie).
		// This is a critical XSS mitigation: even if an attacker injects JavaScript into
		// the page, they cannot steal the session cookie and impersonate the user.
		// ALWAYS set this to true for authentication cookies.
		HttpOnly: true,

		// Secure: false allows the cookie to be sent over HTTP (not just HTTPS).
		// WARNING: This should be false ONLY during local development. In production,
		// this MUST be true to prevent session hijacking via man-in-the-middle attacks.
		// When true, the cookie is only sent over HTTPS connections.
		// TODO: Set this to true in production or when HTTPS is enabled.
		Secure: false,

		// SameSite: Lax provides CSRF protection by controlling when the cookie is sent
		// in cross-site contexts. With Lax mode:
		// - Cookie IS sent for same-site requests (normal navigation within the app)
		// - Cookie IS sent for top-level navigation from external sites (clicking a link)
		// - Cookie is NOT sent for cross-site subrequests (images, iframes, AJAX from other sites)
		// This prevents most CSRF attacks while maintaining good UX.
		// Alternatives: Strict (strongest CSRF protection, may break workflows),
		//               None (no protection, requires Secure=true)
		SameSite: http.SameSiteLaxMode,
	}
}

// SessionMiddleware returns an Echo middleware that loads session data for each request and
// makes it available to handlers via the Echo context. This middleware should be one of the
// first in the middleware chain (after logging/recovery) so that session data is available
// to all subsequent middleware and handlers.
//
// The middleware performs the following steps:
//  1. Retrieves the session cookie from the request (if it exists)
//  2. Decrypts and validates the session data
//  3. Extracts user authentication data (UserID, Email, DisplayName, Role)
//  4. Stores a Session object in the Echo context for handler access
//  5. If the session cookie is invalid or missing, creates a new empty session
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that loads session data for each request.
//
// Example usage:
//
//	e.Use(middleware.SessionMiddleware())
//
// Accessing session in handlers:
//
//	func (h *Handler) SomeHandler(c echo.Context) error {
//	    sess := c.Get("session").(*middleware.Session)
//	    if sess.UserID != 0 {
//	        // User is authenticated
//	        fmt.Printf("Logged in as: %s\n", sess.Email)
//	    }
//	    return c.String(200, "OK")
//	}
//
// Modifying session in handlers:
//
//	sess.UserID = user.ID
//	sess.Email = user.Email
//	sess.Role = user.Role
//	sess.Save(c.Request(), c.Response())
//
// Session data persistence:
//   - Session data is stored in a signed (and optionally encrypted) cookie
//   - Changes to the Session object are NOT persisted until Save() is called
//   - The middleware loads session data on each request; it doesn't automatically save
//   - Handlers must explicitly call sess.Save() after modifying session fields
//
// Error handling:
//   - If session decoding fails (e.g., cookie tampered with, secret key changed), creates new session
//   - This ensures the application doesn't crash on invalid session data
//   - Users will be logged out if their session cookie becomes invalid
func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Attempt to retrieve an existing session from the request cookie.
			// SessionStore.Get() will:
			// 1. Look for a cookie named "bluejay_session" in the request
			// 2. Decode the cookie value (base64 -> gob)
			// 3. Verify the HMAC signature to ensure it hasn't been tampered with
			// 4. Return the session object with the Values map populated
			session, err := SessionStore.Get(c.Request(), "bluejay_session")

			// If session retrieval fails (cookie invalid, tampered, or secret key changed),
			// create a new empty session instead of returning an error. This ensures:
			// - The application doesn't crash on invalid session data
			// - Users with corrupted cookies can still access the site (as guests)
			// - The middleware always provides a Session object to handlers
			if err != nil {
				// Create a new empty session with default values.
				// The underscore ignores the error (New() rarely fails).
				session, _ = SessionStore.New(c.Request(), "bluejay_session")
			}

			// Create our custom Session wrapper around the gorilla Session.
			// Initially, all user fields (UserID, Email, etc.) are zero values.
			sess := &Session{Session: session}

			// Attempt to extract user_id from the session Values map.
			// Type assertion is necessary because Values is map[interface{}]interface{}.
			// If the value exists and is an int64, populate the UserID field.
			// A UserID of 0 (zero value) indicates an unauthenticated session.
			if userID, ok := session.Values["user_id"].(int64); ok {
				sess.UserID = userID
			}

			// Extract email from session Values. If present and is a string, populate Email field.
			if email, ok := session.Values["email"].(string); ok {
				sess.Email = email
			}

			// Extract display name from session Values. Used for UI personalization.
			if displayName, ok := session.Values["display_name"].(string); ok {
				sess.DisplayName = displayName
			}

			// Extract role from session Values. Used for authorization checks.
			if role, ok := session.Values["role"].(string); ok {
				sess.Role = role
			}

			// Store the Session object in the Echo context so handlers and other middleware
			// can access it via c.Get("session"). This is the primary way session data is
			// accessed throughout the application.
			c.Set("session", sess)

			// Proceed to the next handler in the middleware chain.
			// The session is now available in the Echo context for all subsequent handlers.
			return next(c)
		}
	}
}

// Save is a convenience method on the Session struct that persists the current session
// state (UserID, Email, DisplayName, Role) back to the session cookie. This method must
// be called explicitly after modifying session fields to ensure changes are persisted.
//
// The method:
//  1. Updates the underlying session.Values map with current field values
//  2. Calls the gorilla Session.Save() method to serialize, sign, and set the cookie
//  3. Writes the Set-Cookie header to the HTTP response
//
// Parameters:
//   - r: The HTTP request object (needed for session encoding)
//   - w: The HTTP response writer (where the Set-Cookie header is written)
//
// Returns:
//   - error: Any error that occurred during session encoding or cookie writing. In practice,
//     errors are rare and usually indicate serious issues (e.g., response already sent).
//
// Example usage in a login handler:
//
//	sess := c.Get("session").(*middleware.Session)
//	sess.UserID = user.ID
//	sess.Email = user.Email
//	sess.DisplayName = user.Name
//	sess.Role = user.Role
//	if err := sess.Save(c.Request(), c.Response()); err != nil {
//	    return err
//	}
//	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
//
// Example usage in a logout handler:
//
//	sess := c.Get("session").(*middleware.Session)
//	sess.UserID = 0  // Clear authentication
//	sess.Email = ""
//	sess.DisplayName = ""
//	sess.Role = ""
//	sess.Save(c.Request(), c.Response())
//	return c.Redirect(http.StatusSeeOther, "/admin/login")
//
// Important notes:
//   - Must be called BEFORE writing the response body (headers cannot be modified after)
//   - Creates a new Set-Cookie header that overwrites any existing session cookie
//   - The cookie includes all data in session.Values, not just the typed fields
//   - Session size is limited by cookie size constraints (typically 4KB)
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	// Update the underlying session.Values map with the current values of our
	// type-safe fields. This ensures that the struct fields and the Values map
	// stay in sync. Without this, modifications to UserID, Email, etc. would
	// not be persisted to the cookie.
	s.Values["user_id"] = s.UserID
	s.Values["email"] = s.Email
	s.Values["display_name"] = s.DisplayName
	s.Values["role"] = s.Role

	// Call the gorilla Session.Save() method to persist the session data.
	// This method:
	// 1. Serializes session.Values using gob encoding
	// 2. Signs the serialized data with HMAC using the secret key
	// 3. Base64-encodes the signed data
	// 4. Writes a Set-Cookie header to the response with the encoded session data
	// 5. Applies the configured cookie options (Path, MaxAge, HttpOnly, Secure, SameSite)
	//
	// If this returns an error, the session data was not saved and the user's
	// authentication state will be lost on the next request.
	return s.Session.Save(r, w)
}
