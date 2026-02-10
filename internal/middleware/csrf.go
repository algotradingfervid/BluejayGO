package middleware

import (
	// crypto/rand provides cryptographically secure random number generation for
	// generating CSRF tokens. Unlike math/rand, crypto/rand is suitable for
	// security-sensitive applications where token predictability would be dangerous.
	"crypto/rand"

	// encoding/hex provides hexadecimal encoding/decoding. Used to convert the
	// random bytes from crypto/rand into a URL-safe, printable string format
	// suitable for inclusion in forms and HTTP headers.
	"encoding/hex"

	// net/http provides HTTP status codes and standard HTTP constants. Used for
	// the 403 Forbidden response when CSRF validation fails.
	"net/http"

	// github.com/labstack/echo/v4 is the Echo web framework, providing middleware
	// interfaces, context objects, and request/response utilities.
	"github.com/labstack/echo/v4"
)

// generateCSRFToken creates a cryptographically secure random token for CSRF protection.
// This function generates 32 random bytes using crypto/rand (a cryptographically secure
// random number generator) and encodes them as a 64-character hexadecimal string.
//
// The token is used to verify that state-changing requests (POST, PUT, DELETE) originated
// from the application's own forms and not from a malicious third-party site. This prevents
// Cross-Site Request Forgery (CSRF) attacks where an attacker tricks a user's browser into
// making unwanted requests to the application while the user is authenticated.
//
// Returns:
//   - string: A 64-character hexadecimal string (32 bytes encoded) that serves as the
//     CSRF token. This token should be included in forms and validated on submission.
//
// Security properties:
//   - Uses crypto/rand for cryptographic randomness (not predictable)
//   - 32 bytes provides 256 bits of entropy (2^256 possible values)
//   - Hexadecimal encoding makes the token URL-safe and HTML-safe
//   - Each token is unique and virtually impossible to guess or brute-force
//
// Note: This function does not handle read errors from crypto/rand. In practice, crypto/rand
// failing would indicate a serious system-level issue (e.g., no entropy available).
func generateCSRFToken() string {
	// Allocate a byte slice to hold 32 random bytes (256 bits of entropy)
	b := make([]byte, 32)

	// Fill the byte slice with cryptographically secure random data.
	// crypto/rand.Read() uses the operating system's CSPRNG (e.g., /dev/urandom
	// on Unix-like systems, CryptGenRandom on Windows) to generate unpredictable bytes.
	// Error handling is omitted because crypto/rand.Read() should never fail in practice.
	rand.Read(b)

	// Encode the random bytes as a hexadecimal string for safe transmission.
	// This converts 32 bytes into a 64-character string (each byte becomes two hex digits).
	// Hexadecimal encoding ensures the token can be safely included in HTML forms,
	// URLs, and HTTP headers without requiring additional escaping.
	return hex.EncodeToString(b)
}

// CSRF returns an Echo middleware that implements Cross-Site Request Forgery (CSRF) protection
// using the synchronizer token pattern. This middleware generates a unique token per session,
// stores it in the session, and validates it on all state-changing requests (POST, PUT, DELETE).
//
// The middleware follows this flow:
//  1. On each request, retrieve or generate a CSRF token stored in the user's session
//  2. Make the token available to templates so it can be embedded in forms
//  3. For state-changing requests (POST/PUT/DELETE), validate the submitted token
//  4. Reject requests with missing or mismatched tokens with a 403 Forbidden error
//
// This prevents CSRF attacks where a malicious site tricks a user's browser into making
// unwanted authenticated requests. The attacker cannot forge the token because it's:
//   - Generated with cryptographic randomness (unpredictable)
//   - Stored in the HTTP-only session cookie (inaccessible to JavaScript)
//   - Validated on the server side (cannot be bypassed)
//
// The token can be submitted in two ways:
//   - As a form field named "csrf_token" (for traditional form submissions)
//   - As an HTTP header "X-CSRF-Token" (for AJAX/fetch requests)
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that enforces CSRF protection.
//
// Example usage:
//
//	e.Use(middleware.SessionMiddleware())
//	e.Use(middleware.CSRF())
//
// Template usage:
//
//	<form method="POST">
//	    <input type="hidden" name="csrf_token" value="{{.csrf_token}}">
//	    <!-- other form fields -->
//	</form>
//
// AJAX usage:
//
//	fetch('/api/endpoint', {
//	    method: 'POST',
//	    headers: { 'X-CSRF-Token': csrfToken },
//	    body: JSON.stringify(data)
//	})
//
// Security considerations:
//   - Requires SessionMiddleware to be executed first in the middleware chain
//   - Tokens are bound to the session, so they expire when the session expires
//   - Safe methods (GET, HEAD, OPTIONS) are not validated (per HTTP semantics)
//   - Always use HTTPS in production to prevent token interception
func CSRF() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Retrieve the session from the Echo context. The session is populated
			// by SessionMiddleware earlier in the middleware chain.
			sess, ok := c.Get("session").(*Session)

			// If no session exists, skip CSRF protection. This can occur for:
			// - Unauthenticated endpoints that don't require session
			// - Errors in session initialization
			// In production, most protected routes should have RequireAuth() which
			// would catch missing sessions before reaching here.
			if !ok || sess == nil {
				return next(c)
			}

			// Attempt to retrieve an existing CSRF token from the session.
			// The token is stored in the session Values map with key "csrf_token".
			token, ok := sess.Values["csrf_token"].(string)

			// Generate a new token if none exists or if the stored value isn't a valid string.
			// This happens on the user's first request after session creation.
			if !ok || token == "" {
				// Generate a new cryptographically secure random token
				token = generateCSRFToken()

				// Store the token in the session so it persists across requests
				sess.Values["csrf_token"] = token

				// Save the session to ensure the token is stored in the session cookie.
				// This writes the updated session data to the HTTP response.
				sess.Session.Save(c.Request(), c.Response())
			}

			// Make the token available to templates via the Echo context.
			// Templates can access this using {{.csrf_token}} to embed in forms.
			// This allows the same token to be used across multiple forms on the same page.
			c.Set("csrf_token", token)

			// Validate the CSRF token for state-changing HTTP methods.
			// Per HTTP semantics, GET, HEAD, and OPTIONS should be safe (no side effects),
			// so we only validate POST, PUT, and DELETE requests.
			if c.Request().Method == "POST" || c.Request().Method == "PUT" || c.Request().Method == "DELETE" {
				// Extract the CSRF token from the form data (for traditional form submissions)
				formToken := c.FormValue("csrf_token")

				// Extract the CSRF token from the X-CSRF-Token header (for AJAX/API requests)
				headerToken := c.Request().Header.Get("X-CSRF-Token")

				// Prefer the form token if present, otherwise use the header token.
				// This allows the middleware to work with both traditional forms and modern
				// AJAX/fetch requests. Only one token needs to be present.
				submittedToken := formToken
				if submittedToken == "" {
					submittedToken = headerToken
				}

				// Compare the submitted token with the session-stored token.
				// If they don't match (or if no token was submitted), reject the request.
				// This prevents CSRF attacks because an attacker cannot know or guess the token.
				if submittedToken != token {
					// Return a 403 Forbidden error with a descriptive message.
					// Do not reveal the expected token in the error message (security through obscurity).
					return c.JSON(http.StatusForbidden, map[string]string{
						"error": "Invalid CSRF token",
					})
				}
			}

			// Token is valid (or request method doesn't require validation), proceed
			return next(c)
		}
	}
}
