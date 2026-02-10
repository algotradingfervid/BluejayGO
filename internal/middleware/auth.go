package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, ok := c.Get("session").(*Session)
			if !ok || sess.UserID == 0 {
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}
			return next(c)
		}
	}
}

func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, ok := c.Get("session").(*Session)
			if !ok || sess.UserID == 0 {
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}

			for _, role := range roles {
				if sess.Role == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
		}
	}
}
