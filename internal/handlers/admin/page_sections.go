package admin

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type PageSectionsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewPageSectionsHandler(queries *sqlc.Queries, logger *slog.Logger) *PageSectionsHandler {
	return &PageSectionsHandler{queries: queries, logger: logger}
}

func (h *PageSectionsHandler) List(c echo.Context) error {
	sections, err := h.queries.ListAllPageSections(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to load page sections", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Group sections by PageKey
	grouped := make(map[string][]sqlc.PageSection)
	for _, section := range sections {
		grouped[section.PageKey] = append(grouped[section.PageKey], section)
	}

	return c.Render(http.StatusOK, "admin/pages/page_sections_list.html", map[string]interface{}{
		"Title":           "Page Sections",
		"GroupedSections": grouped,
	})
}

func (h *PageSectionsHandler) Edit(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	section, err := h.queries.GetPageSectionByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to load page section", "error", err, "id", id)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	saved := c.QueryParam("saved") == "1"

	return c.Render(http.StatusOK, "admin/pages/page_sections_form.html", map[string]interface{}{
		"Title":   "Edit Page Section",
		"Section": section,
		"Saved":   saved,
	})
}

func (h *PageSectionsHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	err = h.queries.UpdatePageSection(c.Request().Context(), sqlc.UpdatePageSectionParams{
		Heading:             c.FormValue("heading"),
		Subheading:          c.FormValue("subheading"),
		Description:         c.FormValue("description"),
		Label:               c.FormValue("label"),
		PrimaryButtonText:   c.FormValue("primary_button_text"),
		PrimaryButtonUrl:    c.FormValue("primary_button_url"),
		SecondaryButtonText: c.FormValue("secondary_button_text"),
		SecondaryButtonUrl:  c.FormValue("secondary_button_url"),
		IsActive:            c.FormValue("is_active") == "on",
		ID:                  id,
	})

	if err != nil {
		h.logger.Error("failed to update page section", "error", err, "id", id)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/page-sections/"+c.Param("id")+"/edit?saved=1")
}
