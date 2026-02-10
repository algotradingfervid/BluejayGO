package admin

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type FooterHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

// FooterColumnData is a template-friendly struct for each column
type FooterColumnData struct {
	Heading string
	Type    string
	Content string
	Links   []sqlc.FooterLink
}

func NewFooterHandler(queries *sqlc.Queries, logger *slog.Logger) *FooterHandler {
	return &FooterHandler{queries: queries, logger: logger}
}

func (h *FooterHandler) Edit(c echo.Context) error {
	ctx := c.Request().Context()

	settings, err := h.queries.GetSettings(ctx)
	if err != nil {
		h.logger.Error("failed to load settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	columnItems, err := h.queries.ListFooterColumnItems(ctx)
	if err != nil {
		h.logger.Error("failed to load footer column items", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	allLinks, err := h.queries.ListAllFooterLinks(ctx)
	if err != nil {
		h.logger.Error("failed to load footer links", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	legalLinks, err := h.queries.ListFooterLegalLinks(ctx)
	if err != nil {
		h.logger.Error("failed to load footer legal links", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Build links map keyed by column_item_id
	linksMap := make(map[int64][]sqlc.FooterLink)
	for _, link := range allLinks {
		linksMap[link.ColumnItemID] = append(linksMap[link.ColumnItemID], link)
	}

	// Build columns map keyed by column_index
	columnsMap := make(map[int64]sqlc.FooterColumnItem)
	for _, item := range columnItems {
		columnsMap[item.ColumnIndex] = item
	}

	// Build template-friendly ColumnData slice (always 4 entries)
	columnData := make([]FooterColumnData, 4)
	for i := 0; i < 4; i++ {
		col, ok := columnsMap[int64(i)]
		if ok {
			columnData[i] = FooterColumnData{
				Heading: col.Heading,
				Type:    col.Type,
				Content: col.Content,
				Links:   linksMap[col.ID],
			}
		} else {
			columnData[i] = FooterColumnData{
				Type: "links",
			}
		}
	}

	saved := c.QueryParam("saved") == "1"

	return c.Render(http.StatusOK, "admin/pages/footer_form.html", map[string]interface{}{
		"Title":      "Footer Management",
		"Settings":   settings,
		"Saved":      saved,
		"ColumnData": columnData,
		"LegalLinks": legalLinks,
	})
}

func (h *FooterHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse footer columns count
	footerColumns, _ := strconv.ParseInt(c.FormValue("footer_columns"), 10, 64)
	if footerColumns < 2 || footerColumns > 4 {
		footerColumns = 4
	}

	// Parse social toggle
	var footerShowSocial int64
	if c.FormValue("footer_show_social") == "on" {
		footerShowSocial = 1
	}

	// Update footer settings
	err := h.queries.UpdateFooterSettings(ctx, sqlc.UpdateFooterSettingsParams{
		FooterColumns:     footerColumns,
		FooterBgStyle:     c.FormValue("footer_bg_style"),
		FooterShowSocial:  footerShowSocial,
		FooterSocialStyle: c.FormValue("footer_social_style"),
		FooterCopyright:   c.FormValue("footer_copyright"),
	})
	if err != nil {
		h.logger.Error("failed to update footer settings", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Update column items - delete and recreate
	existingItems, err := h.queries.ListFooterColumnItems(ctx)
	if err != nil {
		h.logger.Error("failed to list existing column items", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	for _, item := range existingItems {
		_ = h.queries.DeleteFooterLinksByColumnItem(ctx, item.ID)
		_ = h.queries.DeleteFooterColumnItem(ctx, item.ID)
	}

	// Create column items based on form data
	for i := int64(0); i < footerColumns; i++ {
		prefix := fmt.Sprintf("col_%d_", i)
		colType := c.FormValue(prefix + "type")
		if colType == "" {
			colType = "links"
		}

		heading := c.FormValue(prefix + "heading")
		content := c.FormValue(prefix + "content")

		colItem, err := h.queries.CreateFooterColumnItem(ctx, sqlc.CreateFooterColumnItemParams{
			ColumnIndex: i,
			Type:        colType,
			Heading:     heading,
			Content:     content,
			SortOrder:   i,
		})
		if err != nil {
			h.logger.Error("failed to create footer column item", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// If type is links, create the link items
		if colType == "links" {
			labels := c.Request().Form[prefix+"link_label[]"]
			urls := c.Request().Form[prefix+"link_url[]"]
			for j := 0; j < len(labels) && j < len(urls); j++ {
				if labels[j] == "" && urls[j] == "" {
					continue
				}
				_, err := h.queries.CreateFooterLink(ctx, sqlc.CreateFooterLinkParams{
					ColumnItemID: colItem.ID,
					Label:        labels[j],
					Url:          urls[j],
					SortOrder:    int64(j),
				})
				if err != nil {
					h.logger.Error("failed to create footer link", "error", err)
				}
			}
		}
	}

	// Update legal links - delete and recreate
	_ = h.queries.DeleteAllFooterLegalLinks(ctx)
	legalLabels := c.Request().Form["legal_link_label[]"]
	legalUrls := c.Request().Form["legal_link_url[]"]
	for i := 0; i < len(legalLabels) && i < len(legalUrls); i++ {
		if legalLabels[i] == "" && legalUrls[i] == "" {
			continue
		}
		_, err := h.queries.CreateFooterLegalLink(ctx, sqlc.CreateFooterLegalLinkParams{
			Label:     legalLabels[i],
			Url:       legalUrls[i],
			SortOrder: int64(i),
		})
		if err != nil {
			h.logger.Error("failed to create legal link", "error", err)
		}
	}

	logActivity(c, "updated", "footer", 0, "", "Updated Footer Settings")
	return c.Redirect(http.StatusSeeOther, "/admin/footer?saved=1")
}
