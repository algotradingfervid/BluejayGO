package middleware

import (
	// net/http provides HTTP status codes and standard HTTP constants. Used here
	// for the 429 Too Many Requests status code when rate limits are exceeded.
	"net/http"

	// sync provides synchronization primitives. The sync.Map type is used for
	// thread-safe concurrent access to the visitor tracking map without requiring
	// explicit locking. This is crucial for high-concurrency web servers where
	// multiple goroutines handle requests simultaneously.
	"sync"

	// time provides time and duration measurement functionality. Used for tracking
	// when visitors last made requests and implementing sliding time windows for
	// rate limiting calculations.
	"time"

	// github.com/labstack/echo/v4 is the Echo web framework, providing middleware
	// interfaces, context objects, and utilities for extracting client IP addresses.
	"github.com/labstack/echo/v4"
)

// visitor represents a tracked client in the rate limiting system. It stores the request
// count and the timestamp of the most recent request for a specific IP address.
//
// Fields:
//   - count: The number of requests made by this visitor within the current time window.
//     This is reset to 0 when the time window expires.
//   - lastSeen: The timestamp of the visitor's most recent request. Used to calculate
//     whether the current request falls within the rate limit window and to identify
//     stale entries for cleanup.
//
// This struct is stored in the RateLimiter's visitors map, keyed by IP address.
type visitor struct {
	count    int       // Number of requests in current window
	lastSeen time.Time // Timestamp of most recent request
}

// RateLimiter implements a sliding window rate limiting algorithm to prevent abuse and
// ensure fair resource allocation. It tracks request counts per IP address and enforces
// a configurable maximum number of requests within a time window.
//
// The rate limiter uses an in-memory sync.Map for thread-safe concurrent access, making
// it suitable for high-traffic applications. A background goroutine periodically cleans
// up stale visitor entries to prevent memory leaks.
//
// Fields:
//   - visitors: A thread-safe map storing visitor data keyed by IP address. Uses sync.Map
//     for lock-free reads and writes under concurrent load.
//   - limit: The maximum number of requests allowed per visitor within the time window.
//     Once this limit is exceeded, subsequent requests receive a 429 error.
//   - window: The time duration for the rate limit window. If a visitor's last request
//     was more than this duration ago, their count is reset to 0.
//
// Algorithm: Sliding window with per-IP tracking
//   - Each IP address gets its own request counter
//   - When a request arrives, check if it's within the window (time since last request < window)
//   - If outside window, reset counter to 0
//   - Increment counter and check against limit
//   - Reject request if counter exceeds limit
//
// Memory management:
//   - Background goroutine runs every minute to clean up stale entries
//   - Entries are considered stale if lastSeen > window duration ago
//   - This prevents unbounded memory growth from tracking many one-time visitors
type RateLimiter struct {
	visitors sync.Map      // Map of IP -> visitor, thread-safe for concurrent access
	limit    int           // Maximum requests per window
	window   time.Duration // Time window for rate limiting
}

// NewRateLimiter creates and initializes a new RateLimiter with the specified rate limit
// parameters. It also starts a background goroutine that periodically cleans up stale
// visitor entries to prevent memory leaks.
//
// The rate limiter uses a sliding window algorithm where each IP address is allowed to make
// up to 'limit' requests within any 'window' duration. For example, with limit=100 and
// window=1 minute, an IP can make 100 requests per minute. If they make 100 requests in
// the first 30 seconds, they'll be rate limited until 30 seconds have passed from their
// first request.
//
// Parameters:
//   - limit: Maximum number of requests allowed per visitor within the time window.
//     Common values: 100 for general API endpoints, 10 for sensitive endpoints like login.
//   - window: Time duration for the rate limit window. Common values: 1 minute for API
//     endpoints, 15 minutes for login attempts to prevent brute force attacks.
//
// Returns:
//   - *RateLimiter: A fully initialized rate limiter with cleanup goroutine running.
//
// Example usage:
//
//	// Allow 100 requests per minute for general API
//	apiLimiter := middleware.NewRateLimiter(100, time.Minute)
//	e.Use(apiLimiter.Middleware())
//
//	// Stricter limit for login endpoint (10 attempts per 15 minutes)
//	loginLimiter := middleware.NewRateLimiter(10, 15*time.Minute)
//	e.POST("/admin/login", handler.Login, loginLimiter.Middleware())
//
// Memory and cleanup:
//   - Cleanup goroutine runs every 1 minute
//   - Removes visitor entries that haven't been seen in more than 'window' duration
//   - For a 1-minute window, entries are cleaned up 1-2 minutes after the visitor's last request
//   - This ensures memory usage stays bounded even with millions of unique visitors
//
// Concurrency safety:
//   - Safe to call from multiple goroutines simultaneously
//   - Uses sync.Map for lock-free concurrent access
//   - Cleanup goroutine safely iterates and deletes entries
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	// Initialize the RateLimiter struct with the provided parameters
	rl := &RateLimiter{limit: limit, window: window}

	// Start a background goroutine that runs for the lifetime of the application.
	// This goroutine periodically cleans up stale visitor entries to prevent memory leaks.
	// Without this cleanup, the visitors map would grow indefinitely as new IPs connect.
	go func() {
		for {
			// Sleep for 1 minute between cleanup cycles. This is a balance between:
			// - Frequent enough to prevent excessive memory usage
			// - Infrequent enough to avoid wasting CPU on cleanup
			time.Sleep(time.Minute)

			// Iterate over all visitor entries in the map. Range() is safe for
			// concurrent access and won't block other goroutines.
			rl.visitors.Range(func(key, value interface{}) bool {
				// Type assert the value to *visitor to access its fields
				v := value.(*visitor)

				// Check if this visitor's last request was more than 'window' ago.
				// If so, they're no longer actively using the service, so we can
				// reclaim the memory by deleting their entry.
				if time.Since(v.lastSeen) > rl.window {
					rl.visitors.Delete(key)
				}

				// Return true to continue iterating. Returning false would stop iteration.
				return true
			})
		}
	}()

	return rl
}

// Middleware returns an Echo middleware function that enforces rate limiting based on client
// IP address. This middleware should be applied to routes or route groups that need protection
// from abuse, brute force attacks, or excessive resource consumption.
//
// The middleware implements a sliding window algorithm:
//  1. Extract the client's IP address from the request
//  2. Look up or create a visitor entry for this IP
//  3. Check if the request falls within the rate limit window
//  4. If outside the window, reset the request count (sliding window behavior)
//  5. Increment the request count
//  6. If count exceeds the limit, reject with 429 Too Many Requests
//  7. Otherwise, allow the request to proceed
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that can be used with Echo's Use() method
//     or applied to specific routes/groups.
//
// Example usage:
//
//	limiter := middleware.NewRateLimiter(100, time.Minute)
//	e.Use(limiter.Middleware())  // Apply globally
//
//	// Or apply to specific routes
//	e.POST("/api/resource", handler.Create, limiter.Middleware())
//
// IP address extraction:
//   - Uses Echo's RealIP() which intelligently handles proxied requests
//   - Checks X-Real-IP and X-Forwarded-For headers
//   - Falls back to direct connection remote address
//   - Important: Ensure your reverse proxy (nginx, etc.) sets these headers correctly
//
// Response on rate limit exceeded:
//   - HTTP 429 Too Many Requests status code (standard rate limiting response)
//   - JSON error message instructing the client to try again later
//   - No Retry-After header is set (could be added for better client behavior)
//
// Limitations and considerations:
//   - In-memory storage: Does not persist across server restarts
//   - Single-server only: Does not share state across multiple application instances
//   - For distributed systems, consider Redis-based rate limiting
//   - IP-based: Can be bypassed by attackers with multiple IPs or VPNs
//   - Consider combining with authentication-based rate limiting for logged-in users
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the client's IP address. Echo's RealIP() method intelligently
			// handles proxied requests by checking X-Real-IP and X-Forwarded-For headers
			// before falling back to the direct connection's remote address.
			ip := c.RealIP()

			// Load the existing visitor entry for this IP, or create a new one if this is
			// the first request from this IP. LoadOrStore is atomic and thread-safe.
			// The underscore ignores the boolean "loaded" return value (true if entry existed).
			val, _ := rl.visitors.LoadOrStore(ip, &visitor{})

			// Type assert the interface{} value to *visitor so we can access its fields.
			// This is safe because we only ever store *visitor values in the map.
			v := val.(*visitor)

			// Capture the current time for time window calculations
			now := time.Now()

			// Check if the time elapsed since the visitor's last request exceeds the
			// rate limit window. If so, this is a new time window, so reset the counter.
			// This implements the "sliding window" behavior: if you made 100 requests
			// 2 minutes ago (and window is 1 minute), your counter resets to 0.
			if now.Sub(v.lastSeen) > rl.window {
				v.count = 0
			}

			// Increment the request counter for this visitor. This counts the current request.
			v.count++

			// Update the lastSeen timestamp to now. This is used for:
			// 1. Calculating whether future requests are in the same window
			// 2. Determining when to clean up this entry (if inactive for > window duration)
			v.lastSeen = now

			// Check if the visitor has exceeded the rate limit for this time window.
			// Note: We use > (not >=) so a limit of 100 allows exactly 100 requests.
			if v.count > rl.limit {
				// Visitor has exceeded the rate limit. Return a 429 Too Many Requests error
				// with a JSON body explaining the issue. The client should implement
				// exponential backoff or wait until the time window expires.
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests. Please try again later.",
				})
			}

			// Visitor is within the rate limit, proceed to the next handler
			return next(c)
		}
	}
}
