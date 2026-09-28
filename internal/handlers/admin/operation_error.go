package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Keep administrative load failures inside the authenticated navigation, without
// exposing database details or replacing a failed request with a success page.
func renderOperationError(c echo.Context, title, message, retryURL string) error {
	return c.Render(http.StatusInternalServerError, "admin/pages/operation_error.html", map[string]interface{}{
		"Title": title, "Message": message, "RetryURL": retryURL,
	})
}
