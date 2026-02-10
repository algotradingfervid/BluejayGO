package middleware

import (
	// fmt provides string formatting utilities. Used here to convert non-error panic
	// values into error types for consistent logging and error handling.
	"fmt"

	// log/slog is Go's structured logging package, providing high-performance structured
	// logging with support for key-value pairs. Used to log panic details including
	// error messages and stack traces for debugging.
	"log/slog"

	// net/http provides HTTP status codes and standard HTTP constants. Used here for
	// the 500 Internal Server Error response when a panic is recovered.
	"net/http"

	// runtime/debug provides access to runtime debugging information. The Stack() function
	// captures the full stack trace at the point of a panic, which is essential for
	// debugging unexpected errors in production environments.
	"runtime/debug"

	// github.com/labstack/echo/v4 is the Echo web framework, providing middleware
	// interfaces, context objects, and response utilities.
	"github.com/labstack/echo/v4"
)

// Recovery returns an Echo middleware that recovers from panics in HTTP handlers and
// prevents the entire application from crashing. When a panic occurs, this middleware:
//  1. Catches the panic using Go's recover() mechanism
//  2. Converts the panic value to an error for structured logging
//  3. Logs the error with full stack trace for debugging
//  4. Returns a generic 500 Internal Server Error response to the client
//  5. Allows the application to continue serving other requests
//
// This is a critical safety net for production applications. Without panic recovery,
// a single unhandled panic in a request handler would crash the entire server process,
// taking down all active connections and requiring a restart.
//
// Parameters:
//   - logger: A configured slog.Logger instance for logging panic details. The logger
//     receives structured log entries with the error message, stack trace, and request path.
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that should be applied globally at the
//     application level to catch panics in any handler.
//
// Example usage:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	e.Use(middleware.Recovery(logger))
//
// What gets logged on panic:
//   - error: The panic value converted to an error message
//   - stack: Full stack trace showing the panic location and call chain
//   - path: The request URL path that triggered the panic
//
// Security and UX considerations:
//   - Returns a generic error message to clients (doesn't leak implementation details)
//   - Logs detailed information server-side for debugging
//   - Prevents cascading failures from taking down the entire application
//   - Should be one of the first middleware in the chain to catch panics in other middleware
//
// Common causes of panics this catches:
//   - Nil pointer dereferences
//   - Array/slice index out of bounds
//   - Type assertion failures
//   - Explicit panic() calls in code
//   - Integer divide by zero
//
// Example log output (JSON format):
//
//	{
//	  "time": "2024-01-15T10:30:45Z",
//	  "level": "ERROR",
//	  "msg": "panic recovered",
//	  "error": "runtime error: invalid memory address or nil pointer dereference",
//	  "stack": "goroutine 45 [running]:\n...",
//	  "path": "/admin/products/123"
//	}
func Recovery(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Use defer to ensure the panic recovery logic runs even if a panic occurs.
			// Defer statements are executed when the surrounding function returns, even
			// if that return is caused by a panic. This is the key mechanism that allows
			// us to catch panics.
			defer func() {
				// recover() returns nil if no panic occurred, or the panic value if one did.
				// This is the only way to catch a panic in Go; it must be called within
				// a deferred function.
				if r := recover(); r != nil {
					// Try to convert the panic value to an error type. Panics can be of any type
					// (string, int, custom struct, etc.), but we want to log them consistently
					// as errors. If the panic value is already an error, use it directly.
					err, ok := r.(error)

					// If the panic value isn't an error type, convert it to one using fmt.Errorf.
					// This handles cases like panic("string message") or panic(42).
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					// Log the panic details using structured logging. This creates a permanent
					// record for debugging, including:
					// - error: The error message describing what went wrong
					// - stack: Full stack trace showing where the panic occurred and the call chain
					// - path: The request URL path that triggered the panic (helps identify patterns)
					//
					// We use Error level (not Info or Warn) because panics are always serious issues
					// that require investigation and fixing.
					logger.Error("panic recovered",
						"error", err,
						"stack", string(debug.Stack()),
						"path", c.Request().URL.Path,
					)

					// Return a generic 500 Internal Server Error to the client. We intentionally
					// do not include the panic details in the response for security reasons:
					// - Stack traces can reveal implementation details and file paths
					// - Error messages might contain sensitive information
					// - Clients don't need technical details; they just need to know something went wrong
					//
					// The detailed information is logged server-side where developers can access it.
					// Note: We call c.JSON() directly rather than returning an error because we're
					// in a deferred function that doesn't return a value to the normal control flow.
					c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "Internal server error",
					})
				}
			}()

			// Execute the next handler in the middleware chain. If this or any subsequent
			// handler panics, the deferred function above will catch it. If no panic occurs,
			// the request completes normally and the deferred function does nothing.
			return next(c)
		}
	}
}
