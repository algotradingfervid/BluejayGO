package admin

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"golang.org/x/net/html"
)

// Empty editor HTML such as <div><br></div> is not an article. Drafts may
// remain incomplete, but publishing requires visible text or sourced media.
func hasPublishableBlogBody(body string) bool {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return false
	}
	meaningful := func(text string) bool {
		return strings.TrimFunc(text, func(r rune) bool {
			return unicode.IsSpace(r) || unicode.Is(unicode.Cf, r) || unicode.IsControl(r)
		}) != ""
	}
	var visit func(*html.Node) bool
	visit = func(node *html.Node) bool {
		if node.Type == html.TextNode {
			return meaningful(node.Data)
		}
		if node.Type == html.ElementNode {
			switch node.Data {
			case "head", "script", "style", "template", "noscript":
				return false
			}
			for _, attr := range node.Attr {
				if attr.Key == "hidden" {
					return false
				}
			}
			switch node.Data {
			case "img", "video", "audio", "source", "iframe":
				if node.Data == "source" && (node.Parent == nil || (node.Parent.Data != "video" && node.Parent.Data != "audio")) {
					return false
				}
				for _, attr := range node.Attr {
					if attr.Key == "src" && meaningful(attr.Val) {
						return true
					}
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if visit(child) {
				return true
			}
		}
		return false
	}
	return visit(doc)
}

// Rebuild the submitted form rather than redirecting to saved content. No post
// or association is written until publication validation has succeeded.
func (h *BlogPostsHandler) renderBlogBodyError(c echo.Context, item sqlc.BlogPost) error {
	ctx := c.Request().Context()
	savedSlug := item.Slug
	item.Title = c.FormValue("title")
	item.Slug = c.FormValue("slug")
	if item.Slug == "" {
		item.Slug = makeSlug(item.Title)
	}
	item.Body = c.FormValue("body")
	item.Excerpt = c.FormValue("excerpt")
	item.CategoryID, _ = strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	item.AuthorID, _ = strconv.ParseInt(c.FormValue("author_id"), 10, 64)
	item.FeaturedImageUrl = sql.NullString{String: c.FormValue("featured_image_url"), Valid: c.FormValue("featured_image_url") != ""}
	item.FeaturedImageAlt = sql.NullString{String: c.FormValue("featured_image_alt"), Valid: c.FormValue("featured_image_alt") != ""}
	item.MetaTitle = c.FormValue("meta_title")
	item.MetaDescription = sql.NullString{String: c.FormValue("meta_description"), Valid: c.FormValue("meta_description") != ""}
	item.OgImage = c.FormValue("og_image")
	categories, _ := h.queries.ListBlogCategories(ctx)
	authors, _ := h.queries.ListBlogAuthors(ctx)
	tags, _ := h.queries.ListAllBlogTags(ctx)
	var postTags []sqlc.BlogTag
	for _, value := range c.Request().Form["tag_ids"] {
		id, _ := strconv.ParseInt(value, 10, 64)
		if tag, err := h.queries.GetBlogTag(ctx, id); err == nil {
			postTags = append(postTags, tag)
		}
	}
	var postProducts []sqlc.Product
	for _, value := range c.Request().Form["product_ids"] {
		id, _ := strconv.ParseInt(value, 10, 64)
		if product, err := h.queries.GetProduct(ctx, id); err == nil {
			postProducts = append(postProducts, product)
		}
	}
	title, action := "New Blog Post", "/admin/blog/posts"
	if item.ID != 0 {
		title, action = "Edit Blog Post", fmt.Sprintf("/admin/blog/posts/%d", item.ID)
	}
	return c.Render(http.StatusUnprocessableEntity, "admin/pages/blog_post_form.html", map[string]interface{}{
		"Title": title, "FormAction": action, "Item": item,
		"SavedSlug":  savedSlug,
		"Categories": categories, "Authors": authors, "AllTags": tags,
		"PostTags": postTags, "PostProducts": postProducts,
		"BodyError":          "Add article text or media before publishing. You can save an unfinished draft.",
		"PublishedDateValue": c.FormValue("published_at"),
		"ReadingTimeValue":   c.FormValue("reading_time_minutes"),
	})
}
