package public

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type BlogHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewBlogHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *BlogHandler {
	return &BlogHandler{queries: queries, logger: logger, cache: cache}
}

func (h *BlogHandler) renderAndCache(c echo.Context, cacheKey string, ttlSeconds int, statusCode int, templateName string, data map[string]interface{}) error {
	if settings := c.Get("settings"); settings != nil {
		data["Settings"] = settings
	}
	if cats := c.Get("footer_categories"); cats != nil {
		data["FooterCategories"] = cats
	}
	if sols := c.Get("footer_solutions"); sols != nil {
		data["FooterSolutions"] = sols
	}
	if res := c.Get("footer_resources"); res != nil {
		data["FooterResources"] = res
	}
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		h.logger.Error("template render failed", "template", templateName, "error", err)
		return err
	}
	html := buf.String()
	h.cache.Set(cacheKey, html, ttlSeconds)
	return c.HTML(statusCode, html)
}

// GET /blog
func (h *BlogHandler) BlogListing(c echo.Context) error {
	pageStr := c.QueryParam("page")
	categorySlug := c.QueryParam("category")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("page:blog:page:%d:category:%s", page, categorySlug)
	if cached, ok := h.cache.Get(cacheKey); ok {
		return c.HTML(http.StatusOK, cached.(string))
	}

	ctx := c.Request().Context()
	limit := int64(9)
	offset := int64((page - 1)) * limit

	var posts []sqlc.ListPublishedPostsRow
	var catPosts []sqlc.ListPublishedPostsByCategoryRow
	var totalCount int64
	var err error

	if categorySlug != "" {
		catPosts, err = h.queries.ListPublishedPostsByCategory(ctx, sqlc.ListPublishedPostsByCategoryParams{
			Slug:   categorySlug,
			Limit:  limit,
			Offset: offset,
		})
		totalCount, _ = h.queries.CountPublishedPostsByCategory(ctx, categorySlug)
	} else {
		posts, err = h.queries.ListPublishedPosts(ctx, sqlc.ListPublishedPostsParams{
			Limit:  limit,
			Offset: offset,
		})
		totalCount, _ = h.queries.CountPublishedPosts(ctx)
	}

	if err != nil {
		h.logger.Error("failed to list blog posts", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	featuredPost, _ := h.queries.GetFeaturedPost(ctx)
	categories, _ := h.queries.ListBlogCategories(ctx)

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	data := map[string]interface{}{
		"Title":           "Blog",
		"Posts":           posts,
		"CatPosts":        catPosts,
		"FeaturedPost":    featuredPost,
		"Categories":      categories,
		"CurrentCategory": categorySlug,
		"CurrentPage":     "blog",
		"Page":            page,
		"TotalPages":      totalPages,
		"TotalCount":      totalCount,
	}

	return h.renderAndCache(c, cacheKey, 300, http.StatusOK, "public/pages/blog_listing.html", data)
}

// GET /blog/:slug
func (h *BlogHandler) BlogPost(c echo.Context) error {
	slug := c.Param("slug")
	preview := isPreviewRequest(c)

	if !preview {
		cacheKey := fmt.Sprintf("page:blog:post:%s", slug)
		if cached, ok := h.cache.Get(cacheKey); ok {
			return c.HTML(http.StatusOK, cached.(string))
		}
	}

	ctx := c.Request().Context()

	var post interface{}
	var postID int64
	var postTitle, postSlug, postMetaTitle, metaDesc string
	var postOgImage string
	var postMetaDesc sql.NullString

	if preview {
		p, err := h.queries.GetPostBySlugIncludeDrafts(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		}
		if err != nil {
			h.logger.Error("failed to load blog post", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		post = p
		postID = p.ID
		postTitle = p.Title
		postSlug = p.Slug
		postMetaTitle = p.MetaTitle
		postOgImage = p.OgImage
		postMetaDesc = p.MetaDescription
	} else {
		p, err := h.queries.GetPublishedPostBySlug(ctx, slug)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Post not found")
		}
		if err != nil {
			h.logger.Error("failed to load blog post", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		post = p
		postID = p.ID
		postTitle = p.Title
		postSlug = p.Slug
		postMetaTitle = p.MetaTitle
		postOgImage = p.OgImage
		postMetaDesc = p.MetaDescription
	}

	tags, _ := h.queries.GetPostTagsByPostID(ctx, postID)
	relatedProducts, _ := h.queries.GetPostProductsByPostID(ctx, postID)

	metaDesc = ""
	if postMetaDesc.Valid && postMetaDesc.String != "" {
		metaDesc = postMetaDesc.String
	}

	data := map[string]interface{}{
		"Title":           postTitle,
		"MetaTitle":       postMetaTitle,
		"MetaDescription": metaDesc,
		"OGImage":         postOgImage,
		"CanonicalURL":    fmt.Sprintf("/blog/%s", postSlug),
		"Post":            post,
		"Tags":            tags,
		"RelatedProducts": relatedProducts,
		"CurrentPage":     "blog",
	}

	if preview {
		data["IsPreview"] = true
		data["EditURL"] = fmt.Sprintf("/admin/blog/posts/%d/edit", postID)
		return h.renderAndCache(c, "preview:blog:"+slug, 0, http.StatusOK, "public/pages/blog_post.html", data)
	}

	return h.renderAndCache(c, fmt.Sprintf("page:blog:post:%s", slug), 600, http.StatusOK, "public/pages/blog_post.html", data)
}
