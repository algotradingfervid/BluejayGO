package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo/v4"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func CSRF() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, ok := c.Get("session").(*Session)
			if !ok || sess == nil {
				return next(c)
			}

			// Generate token if not exists
			token, ok := sess.Values["csrf_token"].(string)
			if !ok || token == "" {
				token = generateCSRFToken()
				sess.Values["csrf_token"] = token
				sess.Session.Save(c.Request(), c.Response())
			}

			// Make token available to templates
			c.Set("csrf_token", token)

			// Validate on POST/PUT/DELETE
			if c.Request().Method == "POST" || c.Request().Method == "PUT" || c.Request().Method == "DELETE" {
				formToken := c.FormValue("csrf_token")
				headerToken := c.Request().Header.Get("X-CSRF-Token")
				submittedToken := formToken
				if submittedToken == "" {
					submittedToken = headerToken
				}
				if submittedToken != token {
					return c.JSON(http.StatusForbidden, map[string]string{
						"error": "Invalid CSRF token",
					})
				}
			}

			return next(c)
		}
	}
}
