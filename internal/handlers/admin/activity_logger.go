package admin

import (
	"fmt"

	"github.com/labstack/echo/v4"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

// activityLog is a package-level reference to the activity log service.
var activityLog *services.ActivityLogService

// SetActivityLogService sets the package-level activity log service.
func SetActivityLogService(svc *services.ActivityLogService) {
	activityLog = svc
}

// getUserID extracts the user ID from the echo context session.
func getUserID(c echo.Context) int64 {
	if sess, ok := c.Get("session").(*customMiddleware.Session); ok {
		return sess.UserID
	}
	return 0
}

// logActivity is a convenience helper used by all handlers.
func logActivity(c echo.Context, action, resourceType string, resourceID int64, resourceTitle, descFmt string, args ...interface{}) {
	if activityLog != nil {
		desc := fmt.Sprintf(descFmt, args...)
		activityLog.Log(c.Request().Context(), getUserID(c), action, resourceType, resourceID, resourceTitle, desc)
	}
}
