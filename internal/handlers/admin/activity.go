package admin

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

const activityPerPage = 50

type ActivityHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewActivityHandler(queries *sqlc.Queries, logger *slog.Logger) *ActivityHandler {
	return &ActivityHandler{queries: queries, logger: logger}
}

func (h *ActivityHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	action := c.QueryParam("action")
	search := c.QueryParam("search")
	pageStr := c.QueryParam("page")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * activityPerPage)

	logs, err := h.queries.ListActivityLogs(ctx, sqlc.ListActivityLogsParams{
		FilterAction: action,
		FilterSearch: search,
		PageLimit:    activityPerPage,
		PageOffset:   offset,
	})
	if err != nil {
		h.logger.Error("failed to list activity logs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	total, err := h.queries.CountActivityLogs(ctx, sqlc.CountActivityLogsParams{
		FilterAction: action,
		FilterSearch: search,
	})
	if err != nil {
		h.logger.Error("failed to count activity logs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	totalPages := int(math.Ceil(float64(total) / float64(activityPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	showFrom := offset + 1
	showTo := offset + int64(len(logs))
	if total == 0 {
		showFrom = 0
	}

	hasFilters := action != "" || search != ""

	return c.Render(http.StatusOK, "admin/pages/activity_log.html", map[string]interface{}{
		"Title":      "Activity Log",
		"Logs":       logs,
		"Action":     action,
		"Search":     search,
		"HasFilters": hasFilters,
		"Page":       page,
		"TotalPages": totalPages,
		"Pages":      pages,
		"Total":      total,
		"ShowFrom":   showFrom,
		"ShowTo":     showTo,
	})
}
