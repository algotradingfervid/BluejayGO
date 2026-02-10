package middleware

import (
	// log/slog is Go's structured logging package (introduced in Go 1.21). It provides
	// high-performance structured logging with support for key-value pairs, log levels,
	// and multiple output formats (JSON, text). Used here to log HTTP request metadata.
	"log/slog"

	// time provides time and duration measurement functionality. Used to calculate
	// the duration of each HTTP request by capturing timestamps before and after
	// handler execution.
	"time"

	// github.com/labstack/echo/v4 is the Echo web framework, providing middleware
	// interfaces, context objects, and request/response utilities.
	"github.com/labstack/echo/v4"
)

// Logging returns an Echo middleware that logs structured information about each HTTP request.
// This middleware captures key request metadata including HTTP method, path, response status,
// request duration, and client IP address. The logging occurs after the request is processed,
// allowing the middleware to record the actual response status and timing information.
//
// The middleware uses structured logging (slog) which outputs logs in a consistent, parseable
// format suitable for log aggregation systems like ELK, Splunk, or CloudWatch. Each log entry
// includes the following fields:
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - path: Request URL path (without query parameters)
//   - status: HTTP response status code (200, 404, 500, etc.)
//   - duration_ms: Request processing time in milliseconds
//   - ip: Client IP address (extracted from X-Real-IP, X-Forwarded-For, or remote addr)
//
// Parameters:
//   - logger: A configured slog.Logger instance that will receive the log entries.
//     This allows the caller to configure output format, log level filtering, and
//     output destination (stdout, file, remote logging service, etc.).
//
// Returns:
//   - echo.MiddlewareFunc: A middleware function that logs request information.
//
// Example usage:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	e.Use(middleware.Logging(logger))
//
// Performance considerations:
//   - time.Now() is called twice per request (negligible overhead)
//   - Structured logging is highly optimized in slog (minimal allocation)
//   - Logging happens after the response is sent, so it doesn't delay the client
//   - For high-traffic applications, consider sampling or filtering logs
//
// Example log output (JSON format):
//
//	{
//	  "time": "2024-01-15T10:30:45Z",
//	  "level": "INFO",
//	  "msg": "request",
//	  "method": "GET",
//	  "path": "/admin/products",
//	  "status": 200,
//	  "duration_ms": 45,
//	  "ip": "192.168.1.100"
//	}
func Logging(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Capture the request start time before executing the handler.
			// This allows us to calculate the total request processing duration,
			// including all subsequent middleware and the final handler execution.
			start := time.Now()

			// Execute the next handler in the middleware chain (or the final route handler).
			// We capture the error to return it after logging, ensuring that error responses
			// are still properly propagated even though we log first.
			err := next(c)

			// Log structured request information using slog's key-value format.
			// Info level is appropriate for normal request logging (not warnings or errors).
			//
			// Fields logged:
			// - method: HTTP verb (GET, POST, etc.) from the request
			// - path: URL path component (excludes query string and fragment)
			// - status: HTTP status code set by the handler (e.g., 200, 404, 500)
			// - duration_ms: Time elapsed from request start to completion, in milliseconds.
			//   This includes all middleware execution time and handler processing time.
			// - ip: Client's IP address. Echo's RealIP() intelligently extracts the IP from
			//   X-Real-IP, X-Forwarded-For headers (for proxied requests), or falls back
			//   to the direct connection's remote address.
			logger.Info("request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"duration_ms", time.Since(start).Milliseconds(),
				"ip", c.RealIP(),
			)

			// Return the error from the handler execution. This ensures that if the handler
			// or any subsequent middleware returned an error, it's properly propagated to
			// Echo's error handling mechanism, even though we've logged the request.
			return err
		}
	}
}
