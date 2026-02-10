package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

func SettingsLoader(queries *sqlc.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			settings, err := queries.GetSettings(ctx)
			if err != nil {
				slog.Warn("settings middleware: failed to load settings", "error", err)
			} else {
				c.Set("settings", settings)
			}

			// Load footer categories
			categories, err := queries.ListProductCategories(ctx)
			if err != nil {
				slog.Warn("settings middleware: failed to load footer categories", "error", err)
			} else {
				c.Set("footer_categories", categories)
			}

			// Load footer solutions
			solutions, err := queries.ListPublishedSolutions(ctx)
			if err != nil {
				slog.Warn("settings middleware: failed to load footer solutions", "error", err)
			} else {
				c.Set("footer_solutions", solutions)
			}

			// Load footer resource links
			resources, err := queries.ListPageSections(ctx, "footer")
			if err != nil {
				slog.Warn("settings middleware: failed to load footer resources", "error", err)
			} else {
				c.Set("footer_resources", resources)
			}

			return next(c)
		}
	}
}
