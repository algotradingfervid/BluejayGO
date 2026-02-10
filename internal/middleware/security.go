package middleware

import (
	// github.com/labstack/echo/v4 is the Echo web framework, providing middleware
	// interfaces, context objects, and response header manipulation utilities.
	"github.com/labstack/echo/v4"
)

// SecurityHeaders returns an Echo middleware that sets security-related HTTP response headers
// to protect against common web vulnerabilities. This middleware implements defense-in-depth
// by adding multiple layers of browser-based security controls that help prevent attacks like
// XSS (Cross-Site Scripting), clickjacking, MIME-sniffing attacks, and information leakage.
//
// The middleware sets the following security headers on every response:
//
//  1. X-Content-Type-Options: nosniff
//     Prevents MIME-sniffing attacks where browsers try to detect content types and might
//     execute malicious content disguised as a safe file type (e.g., JavaScript disguised as an image).
//
//  2. X-Frame-Options: DENY
//     Prevents the application from being embedded in iframes, protecting against clickjacking
//     attacks where attackers overlay invisible iframes to trick users into clicking malicious elements.
//
//  3. X-XSS-Protection: 1; mode=block
//     Enables the browser's built-in XSS filter and instructs it to block the page rather than
//     sanitize when XSS is detected. Note: This is a legacy header; modern browsers rely on CSP instead.
//
//  4. Referrer-Policy: strict-origin-when-cross-origin
//     Controls how much referrer information is included in requests. Sends full URL for same-origin
//     requests, but only sends the origin (no path) for cross-origin requests, balancing analytics
//     needs with privacy protection.
//
//  5. Content-Security-Policy (CSP)
//     Defines which resources (scripts, styles, fonts, images) can be loaded and from where. This is
//     the most powerful header for preventing XSS attacks by whitelisting trusted content sources.
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that should be applied globally to add
//     security headers to all responses.
//
// Example usage:
//
//	e.Use(middleware.SecurityHeaders())
//
// Security considerations:
//   - Should be one of the first middleware in the chain to ensure headers are set early
//   - CSP policy may need adjustment when adding new third-party services
//   - 'unsafe-inline' and 'unsafe-eval' in CSP reduce security; consider removing in production
//   - Modern browsers support CSP level 3; older browsers may not respect all directives
//   - These headers are defense-in-depth; they complement (not replace) server-side security
//
// Content Security Policy breakdown:
//   - default-src 'self': By default, only load resources from the same origin
//   - script-src: JavaScript sources (includes 'unsafe-inline' and 'unsafe-eval' for compatibility
//     with Tailwind CDN, HTMX, and inline scripts; consider removing for better security)
//   - style-src: CSS sources (allows inline styles for Tailwind CDN and utility classes)
//   - font-src: Web font sources (Google Fonts)
//   - img-src: Image sources (self-hosted, data URIs, and any HTTPS source for flexibility)
//
// Trade-offs in current CSP:
//   - 'unsafe-inline' allows inline scripts/styles but reduces XSS protection
//   - 'unsafe-eval' allows eval() and Function() but enables certain attack vectors
//   - These are needed for Tailwind CDN's JIT mode and some third-party libraries
//   - For maximum security, use a build process for Tailwind and remove 'unsafe-*' directives
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// X-Content-Type-Options: nosniff
			// Instructs browsers to strictly follow the Content-Type header and not attempt
			// to MIME-sniff the response. This prevents attacks where malicious JavaScript
			// is uploaded as an image but executed because the browser detects it as script.
			// Example attack prevented: Uploading a file named "image.jpg" containing JavaScript,
			// then linking to it with <script src="image.jpg"> hoping the browser sniffs and executes it.
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// X-Frame-Options: DENY
			// Prevents the application from being displayed in any iframe, frame, or embed element.
			// This protects against clickjacking attacks where an attacker overlays the application
			// in an invisible iframe and tricks users into clicking on sensitive actions.
			// Alternative values: SAMEORIGIN (allow same-origin framing), ALLOW-FROM uri (allow specific origin).
			// DENY is the most secure option unless you specifically need iframe embedding.
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// X-XSS-Protection: 1; mode=block
			// Enables the browser's built-in XSS filter (if available) and instructs it to block
			// the entire page when XSS is detected, rather than trying to sanitize the attack.
			// Note: This is a legacy header. Modern browsers have deprecated it in favor of CSP.
			// However, it still provides defense-in-depth for older browsers.
			// "1; mode=block" means: enable filter (1) and block the page (mode=block) rather than sanitize.
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			// Referrer-Policy: strict-origin-when-cross-origin
			// Controls how much referrer information (the URL the user came from) is included
			// in requests to other sites. This balances analytics/tracking needs with privacy:
			// - Same-origin requests: Send full URL (path and query string included)
			// - Cross-origin requests to HTTPS: Send only the origin (scheme + host, no path)
			// - Cross-origin requests to HTTP from HTTPS: Send nothing (downgrade protection)
			// This prevents leaking sensitive information in URL parameters to third-party sites.
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Content-Security-Policy (CSP)
			// The most powerful security header, defining a whitelist of sources from which
			// various resource types can be loaded. This is the primary defense against XSS attacks.
			//
			// Directive breakdown:
			//
			// - default-src 'self':
			//   By default, only allow resources from the same origin (scheme + host + port).
			//   This is the fallback for any directive not explicitly specified.
			//
			// - script-src 'self' 'unsafe-inline' 'unsafe-eval' cdn.tailwindcss.com cdn.jsdelivr.net fonts.googleapis.com:
			//   Allow JavaScript from:
			//   * 'self': Same origin (our own scripts)
			//   * 'unsafe-inline': Inline <script> tags and event handlers (onclick, etc.)
			//     WARNING: This reduces XSS protection significantly. Consider using nonces in production.
			//   * 'unsafe-eval': eval(), Function(), and similar dynamic code execution
			//     WARNING: This enables certain XSS attack vectors. Required for some libraries.
			//   * cdn.tailwindcss.com: Tailwind CDN for JIT compilation
			//   * cdn.jsdelivr.net: Common CDN for third-party libraries
			//   * fonts.googleapis.com: Google Fonts API (may load scripts)
			//
			// - style-src 'self' 'unsafe-inline' fonts.googleapis.com cdn.tailwindcss.com:
			//   Allow CSS from:
			//   * 'self': Same origin (our own stylesheets)
			//   * 'unsafe-inline': Inline <style> tags and style attributes
			//     Required for Tailwind utility classes and inline styles.
			//   * fonts.googleapis.com: Google Fonts CSS
			//   * cdn.tailwindcss.com: Tailwind CDN styles
			//
			// - font-src 'self' fonts.gstatic.com:
			//   Allow web fonts from:
			//   * 'self': Same origin (our own font files)
			//   * fonts.gstatic.com: Google Fonts font files (WOFF, WOFF2, TTF, etc.)
			//
			// - img-src 'self' data: https::
			//   Allow images from:
			//   * 'self': Same origin (our own images)
			//   * data:: Data URIs (inline images encoded in base64)
			//   * https:: Any HTTPS source (allows loading images from external HTTPS URLs)
			//     This is permissive but necessary for user-generated content and external images.
			//
			// Production hardening recommendations:
			// 1. Remove 'unsafe-inline' and 'unsafe-eval' by:
			//    - Using a build process for Tailwind (no CDN)
			//    - Moving inline scripts to external files
			//    - Using CSP nonces or hashes for necessary inline scripts
			// 2. Replace 'https:' in img-src with specific whitelisted domains
			// 3. Add report-uri or report-to directives to monitor CSP violations
			// 4. Consider adding frame-ancestors directive (redundant with X-Frame-Options but more flexible)
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' cdn.tailwindcss.com cdn.jsdelivr.net fonts.googleapis.com; style-src 'self' 'unsafe-inline' fonts.googleapis.com cdn.tailwindcss.com; font-src 'self' fonts.gstatic.com; img-src 'self' data: https:;")

			// Proceed to the next handler in the middleware chain.
			// The security headers are already set on the response and will be sent
			// when the response is finalized.
			return next(c)
		}
	}
}
