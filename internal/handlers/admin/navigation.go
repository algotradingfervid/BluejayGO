package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type NavigationHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewNavigationHandler(queries *sqlc.Queries, logger *slog.Logger) *NavigationHandler {
	return &NavigationHandler{queries: queries, logger: logger}
}

// NavigationItemView is a template-friendly struct with children
type NavigationItemView struct {
	sqlc.NavigationItem
	Children []sqlc.NavigationItem
}

func (h *NavigationHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	menus, err := h.queries.ListNavigationMenus(ctx)
	if err != nil {
		h.logger.Error("failed to list navigation menus", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Count items per menu
	menuItemCounts := make(map[int64]int)
	for _, menu := range menus {
		items, err := h.queries.ListNavigationItems(ctx, menu.ID)
		if err == nil {
			menuItemCounts[menu.ID] = len(items)
		}
	}

	return c.Render(http.StatusOK, "admin/pages/navigation_list.html", map[string]interface{}{
		"Title":          "Navigation Menus",
		"Menus":          menus,
		"MenuItemCounts": menuItemCounts,
	})
}

func (h *NavigationHandler) Edit(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	menu, err := h.queries.GetNavigationMenu(ctx, id)
	if err != nil {
		h.logger.Error("failed to get navigation menu", "error", err)
		return echo.NewHTTPError(http.StatusNotFound)
	}

	items, err := h.queries.ListNavigationItems(ctx, id)
	if err != nil {
		h.logger.Error("failed to list navigation items", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Build tree: separate top-level items and children
	topLevel := []NavigationItemView{}
	childrenMap := make(map[int64][]sqlc.NavigationItem)
	for _, item := range items {
		if item.ParentID.Valid {
			childrenMap[item.ParentID.Int64] = append(childrenMap[item.ParentID.Int64], item)
		}
	}
	for _, item := range items {
		if !item.ParentID.Valid {
			topLevel = append(topLevel, NavigationItemView{
				NavigationItem: item,
				Children:       childrenMap[item.ID],
			})
		}
	}

	saved := c.QueryParam("saved") == "1"

	pageOptions := []string{
		"Products", "Solutions", "About", "Blog", "Case Studies",
		"Whitepapers", "Partners", "Contact",
	}

	return c.Render(http.StatusOK, "admin/pages/navigation_editor.html", map[string]interface{}{
		"Title":       fmt.Sprintf("Edit Menu: %s", menu.Name),
		"Menu":        menu,
		"Items":       topLevel,
		"AllItems":    items,
		"Saved":       saved,
		"PageOptions": pageOptions,
	})
}

func (h *NavigationHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	name := c.FormValue("name")
	location := c.FormValue("location")

	if name == "" {
		name = "New Menu"
	}
	if location == "" {
		location = "header"
	}

	menu, err := h.queries.CreateNavigationMenu(ctx, sqlc.CreateNavigationMenuParams{
		Name:     name,
		Location: location,
	})
	if err != nil {
		h.logger.Error("failed to create navigation menu", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "created", "navigation", 0, c.FormValue("name"), "Created Navigation Menu '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/navigation/%d?saved=1", menu.ID))
}

func (h *NavigationHandler) DeleteMenu(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// Delete items first, then menu
	if err := h.queries.DeleteNavigationItemsByMenu(ctx, id); err != nil {
		h.logger.Error("failed to delete navigation items", "error", err)
	}
	if err := h.queries.DeleteNavigationMenu(ctx, id); err != nil {
		h.logger.Error("failed to delete navigation menu", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "deleted", "navigation", id, "", "Deleted Navigation Menu #%d", id)
	return c.NoContent(http.StatusOK)
}

func (h *NavigationHandler) AddItem(c echo.Context) error {
	ctx := c.Request().Context()

	menuID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	label := c.FormValue("label")
	linkType := c.FormValue("link_type")
	url := c.FormValue("url")
	pageIdentifier := c.FormValue("page_identifier")

	if label == "" {
		label = "New Item"
	}
	if linkType == "" {
		linkType = "page"
	}

	// Get current max sort order
	items, _ := h.queries.ListNavigationItems(ctx, menuID)
	nextSort := int64(len(items))

	// For page type, generate URL from page identifier
	if linkType == "page" && pageIdentifier != "" && url == "" {
		url = "/" + slugifyNav(pageIdentifier)
		label = pageIdentifier
	}

	_, err = h.queries.CreateNavigationItem(ctx, sqlc.CreateNavigationItemParams{
		MenuID:         menuID,
		Label:          label,
		LinkType:       linkType,
		Url:            sql.NullString{String: url, Valid: url != ""},
		PageIdentifier: sql.NullString{String: pageIdentifier, Valid: pageIdentifier != ""},
		OpenNewTab:     sql.NullInt64{Int64: 0, Valid: true},
		IsActive:       sql.NullInt64{Int64: 1, Valid: true},
		SortOrder:      sql.NullInt64{Int64: nextSort, Valid: true},
	})
	if err != nil {
		h.logger.Error("failed to create navigation item", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "navigation", 0, "", "Updated Navigation Items")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/navigation/%d?saved=1", menuID))
}

func (h *NavigationHandler) UpdateItem(c echo.Context) error {
	ctx := c.Request().Context()

	itemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	item, err := h.queries.GetNavigationItem(ctx, itemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	label := c.FormValue("label")
	linkType := c.FormValue("link_type")
	url := c.FormValue("url")
	pageIdentifier := c.FormValue("page_identifier")
	openNewTab := c.FormValue("open_new_tab") == "on"
	isActive := c.FormValue("is_active") == "on"

	var openNewTabInt int64
	if openNewTab {
		openNewTabInt = 1
	}
	var isActiveInt int64
	if isActive {
		isActiveInt = 1
	}

	err = h.queries.UpdateNavigationItem(ctx, sqlc.UpdateNavigationItemParams{
		ID:             itemID,
		Label:          label,
		LinkType:       linkType,
		Url:            sql.NullString{String: url, Valid: url != ""},
		PageIdentifier: sql.NullString{String: pageIdentifier, Valid: pageIdentifier != ""},
		OpenNewTab:     sql.NullInt64{Int64: openNewTabInt, Valid: true},
		IsActive:       sql.NullInt64{Int64: isActiveInt, Valid: true},
		ParentID:       item.ParentID,
		SortOrder:      item.SortOrder,
	})
	if err != nil {
		h.logger.Error("failed to update navigation item", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "navigation", 0, "", "Updated Navigation Items")
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/navigation/%d?saved=1", item.MenuID))
}

func (h *NavigationHandler) DeleteItem(c echo.Context) error {
	ctx := c.Request().Context()

	itemID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := h.queries.DeleteNavigationItem(ctx, itemID); err != nil {
		h.logger.Error("failed to delete navigation item", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "navigation", 0, "", "Updated Navigation Items")
	return c.NoContent(http.StatusOK)
}

func (h *NavigationHandler) Reorder(c echo.Context) error {
	ctx := c.Request().Context()

	type ReorderItem struct {
		ID       int64  `json:"id"`
		ParentID *int64 `json:"parent_id"`
		Order    int64  `json:"order"`
	}

	var items []ReorderItem
	if err := json.NewDecoder(c.Request().Body).Decode(&items); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}

	for _, item := range items {
		parentID := sql.NullInt64{}
		if item.ParentID != nil {
			parentID = sql.NullInt64{Int64: *item.ParentID, Valid: true}
		}
		err := h.queries.UpdateNavigationItemOrder(ctx, sqlc.UpdateNavigationItemOrderParams{
			ID:        item.ID,
			SortOrder: sql.NullInt64{Int64: item.Order, Valid: true},
			ParentID:  parentID,
		})
		if err != nil {
			h.logger.Error("failed to reorder navigation item", "error", err, "id", item.ID)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *NavigationHandler) UpdateMenu(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	name := c.FormValue("name")
	location := c.FormValue("location")

	err = h.queries.UpdateNavigationMenu(ctx, sqlc.UpdateNavigationMenuParams{
		ID:       id,
		Name:     name,
		Location: location,
	})
	if err != nil {
		h.logger.Error("failed to update navigation menu", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "navigation", id, c.FormValue("name"), "Updated Navigation Menu '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/navigation/%d?saved=1", id))
}

func slugifyNav(s string) string {
	result := ""
	for _, c := range s {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' {
			result += string(c)
		} else if c >= 'A' && c <= 'Z' {
			result += string(c + 32)
		} else if c == ' ' {
			result += "-"
		}
	}
	return result
}
