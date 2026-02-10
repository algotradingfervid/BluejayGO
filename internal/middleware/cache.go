package middleware

import (
	// github.com/labstack/echo/v4 is the Echo web framework, providing the middleware
	// interface and context needed for HTTP response header manipulation.
	"github.com/labstack/echo/v4"
)

// NoCache returns an Echo middleware that disables all HTTP caching for the response.
// This middleware sets multiple cache-control headers to ensure that browsers, proxies,
// and CDNs do not cache the response. This is critical for dynamic content that changes
// frequently or contains sensitive information that should not be stored in caches.
//
// The middleware sets three complementary headers:
//   - Cache-Control: Instructs HTTP/1.1 caches not to store or reuse the response
//   - Pragma: Provides backward compatibility with HTTP/1.0 caches
//   - Expires: Sets an explicit expiration date in the past (epoch time)
//
// This should be applied to:
//   - Admin panel pages that display real-time data
//   - Pages containing user-specific or sensitive information
//   - API endpoints that return dynamic data
//   - Authentication-related pages (login, logout, password reset)
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that disables caching.
//
// Example usage:
//
//	admin := e.Group("/admin", middleware.NoCache())
//
// Caching strategy explanation:
//   - "no-cache": Response can be stored but must be revalidated before reuse
//   - "no-store": Response must not be stored in any cache (stronger than no-cache)
//   - "must-revalidate": Once stale, cache must revalidate with origin server
//   - "Pragma: no-cache": HTTP/1.0 backward compatibility directive
//   - "Expires: 0": Tells cache the resource is already expired
func NoCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set Cache-Control header with multiple directives to prevent caching.
			// This is the primary HTTP/1.1 cache control mechanism.
			// - "no-cache": Allows storage but requires revalidation before use
			// - "no-store": Prevents any caching of the response (strongest directive)
			// - "must-revalidate": Requires stale cache entries to be revalidated
			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			// Set Pragma header for HTTP/1.0 backward compatibility. Older proxies
			// and browsers that don't understand Cache-Control will respect this.
			c.Response().Header().Set("Pragma", "no-cache")

			// Set Expires to 0 (epoch time), indicating the content is already expired.
			// This provides an additional layer of cache prevention for legacy systems.
			c.Response().Header().Set("Expires", "0")

			// Proceed to the next handler in the middleware chain
			return next(c)
		}
	}
}

// CacheStatic returns an Echo middleware that enables aggressive caching for static assets.
// This middleware sets cache headers to instruct browsers and CDNs to cache the response
// for one year (31536000 seconds). This is ideal for immutable static resources like fonts,
// images with content-based URLs, and versioned JavaScript/CSS bundles.
//
// The middleware sets the response as publicly cacheable with a maximum age of one year,
// which is effectively "cache forever" for web purposes. This dramatically reduces server
// load and improves page load times for returning visitors.
//
// This should be applied to:
//   - Static assets in /public/ with cache-busting URLs (e.g., app.v123.js)
//   - Immutable resources like fonts, icons, and versioned images
//   - Third-party libraries served from the application
//
// Do NOT use this for:
//   - HTML pages (content may change)
//   - API responses (data is dynamic)
//   - User-generated content (may be updated or deleted)
//   - Any resource without cache-busting in the URL
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that enables long-term caching.
//
// Example usage:
//
//	e.Static("/static", "public/static", middleware.CacheStatic())
//
// Caching strategy explanation:
//   - "public": Response can be cached by any cache (browser, proxy, CDN)
//   - "max-age=31536000": Cache can reuse response for 31536000 seconds (1 year)
//   - This relies on cache-busting strategies (versioned URLs) to force updates
//   - One year is the practical maximum; RFC 7234 recommends no more than 1 year
func CacheStatic() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set Cache-Control header to enable public caching for one year.
			// - "public": Allows any cache (shared or private) to store the response
			// - "max-age=31536000": Cache can serve this response for 365 days without
			//   revalidation. This is effectively "cache forever" since most cache
			//   implementations won't keep entries this long anyway.
			//
			// Note: This assumes the URL contains a version identifier or content hash.
			// When the file changes, the URL should change, forcing a cache miss.
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000")

			// Proceed to the next handler in the middleware chain
			return next(c)
		}
	}
}
