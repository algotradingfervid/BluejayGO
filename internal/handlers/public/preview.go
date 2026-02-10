package public

import (
	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/internal/middleware"
)

// isPreviewRequest checks if the request is a preview request with a valid admin session.
func isPreviewRequest(c echo.Context) bool {
	if c.QueryParam("preview") != "true" {
		return false
	}
	sess, ok := c.Get("session").(*middleware.Session)
	if !ok || sess == nil || sess.UserID == 0 {
		return false
	}
	return true
}
