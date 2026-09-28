package public

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
)

// NewPublicHTTPErrorHandler gives browser navigation a useful recovery page while
// retaining Echo's error responses for administrative, asset and API requests.
func NewPublicHTTPErrorHandler(e *echo.Echo, queries *sqlc.Queries) echo.HTTPErrorHandler {
	fallback := e.DefaultHTTPErrorHandler
	return func(err error, c echo.Context) {
		var httpErr *echo.HTTPError
		request := c.Request()
		path := request.URL.Path
		htmlRequest := strings.Contains(request.Header.Get("Accept"), "text/html")
		excluded := false
		for _, prefix := range []string{"/admin", "/api", "/public", "/uploads", "/health"} {
			if path == prefix || strings.HasPrefix(path, prefix+"/") {
				excluded = true
				break
			}
		}
		if c.Response().Committed || !errors.As(err, &httpErr) || httpErr.Code != http.StatusNotFound ||
			(request.Method != http.MethodGet && request.Method != http.MethodHead) ||
			request.Header.Get("HX-Request") == "true" || !htmlRequest || excluded {
			fallback(err, c)
			return
		}
		renderErr := customMiddleware.SettingsLoader(queries)(func(c echo.Context) error {
			data := map[string]interface{}{
				"Title": "Page not found", "NoIndex": true,
				"Settings": c.Get("settings"), "FooterCategories": c.Get("footer_categories"),
				"FooterSolutions": c.Get("footer_solutions"), "FooterResources": c.Get("footer_resources"),
			}
			var body bytes.Buffer
			if renderErr := e.Renderer.Render(&body, "public/pages/not_found.html", data, c); renderErr != nil {
				return renderErr
			}
			if request.Method == http.MethodHead {
				c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
				return c.NoContent(http.StatusNotFound)
			}
			return c.HTML(http.StatusNotFound, body.String())
		})(c)
		if renderErr != nil {
			e.Logger.Error(renderErr)
			fallback(err, c)
		}
	}
}
