package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type BlogPostsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewBlogPostsHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *BlogPostsHandler {
	return &BlogPostsHandler{queries: queries, logger: logger, cache: cache}
}

const blogPostsPerPage = 15

func (h *BlogPostsHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	search := c.QueryParam("search")
	status := c.QueryParam("status")
	categoryStr := c.QueryParam("category")
	authorStr := c.QueryParam("author")
	pageStr := c.QueryParam("page")

	var categoryID int64
	if categoryStr != "" {
		categoryID, _ = strconv.ParseInt(categoryStr, 10, 64)
	}
	var authorID int64
	if authorStr != "" {
		authorID, _ = strconv.ParseInt(authorStr, 10, 64)
	}

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * blogPostsPerPage)

	posts, err := h.queries.ListBlogPostsAdminFiltered(ctx, sqlc.ListBlogPostsAdminFilteredParams{
		FilterStatus:   status,
		FilterCategory: categoryID,
		FilterAuthor:   authorID,
		FilterSearch:   search,
		PageLimit:      blogPostsPerPage,
		PageOffset:     offset,
	})
	if err != nil {
		h.logger.Error("failed to list blog posts", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	total, err := h.queries.CountBlogPostsAdminFiltered(ctx, sqlc.CountBlogPostsAdminFilteredParams{
		FilterStatus:   status,
		FilterCategory: categoryID,
		FilterAuthor:   authorID,
		FilterSearch:   search,
	})
	if err != nil {
		h.logger.Error("failed to count blog posts", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	categories, _ := h.queries.ListBlogCategories(ctx)
	authors, _ := h.queries.ListBlogAuthors(ctx)

	totalPages := int(math.Ceil(float64(total) / float64(blogPostsPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	showFrom := offset + 1
	showTo := offset + int64(len(posts))
	if total == 0 {
		showFrom = 0
	}

	hasFilters := search != "" || status != "" || categoryStr != "" || authorStr != ""

	return c.Render(http.StatusOK, "admin/pages/blog_posts_list.html", map[string]interface{}{
		"Title":      "Manage Blog Posts",
		"Posts":      posts,
		"Categories": categories,
		"Authors":    authors,
		"Search":     search,
		"Status":     status,
		"CategoryID": categoryID,
		"AuthorID":   authorID,
		"HasFilters": hasFilters,
		"Page":       page,
		"TotalPages": totalPages,
		"Pages":      pages,
		"Total":      total,
		"ShowFrom":   showFrom,
		"ShowTo":     showTo,
	})
}

func (h *BlogPostsHandler) New(c echo.Context) error {
	ctx := c.Request().Context()
	categories, _ := h.queries.ListBlogCategories(ctx)
	authors, _ := h.queries.ListBlogAuthors(ctx)
	tags, _ := h.queries.ListAllBlogTags(ctx)
	return c.Render(http.StatusOK, "admin/pages/blog_post_form.html", map[string]interface{}{
		"Title":      "New Blog Post",
		"FormAction": "/admin/blog/posts",
		"Item":       nil,
		"Categories": categories,
		"Authors":    authors,
		"AllTags":    tags,
		"PostTags":     nil,
		"PostProducts": nil,
	})
}

func calculateReadingTime(body string) int64 {
	words := len(strings.Fields(strings.TrimSpace(body)))
	minutes := words / 200
	if minutes < 1 {
		minutes = 1
	}
	return int64(minutes)
}

func (h *BlogPostsHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	body := c.FormValue("body")
	readingTime, _ := strconv.ParseInt(c.FormValue("reading_time_minutes"), 10, 64)
	if readingTime == 0 {
		readingTime = calculateReadingTime(body)
	}

	categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	authorID, _ := strconv.ParseInt(c.FormValue("author_id"), 10, 64)

	featuredURL := c.FormValue("featured_image_url")
	featuredAlt := c.FormValue("featured_image_alt")
	metaDesc := c.FormValue("meta_description")
	excerpt := c.FormValue("excerpt")
	status := c.FormValue("status")

	var publishedAt sql.NullTime
	if status == "published" {
		pubStr := c.FormValue("published_at")
		if t, err := time.Parse("2006-01-02T15:04", pubStr); err == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		} else {
			publishedAt = sql.NullTime{Time: time.Now(), Valid: true}
		}
	}

	post, err := h.queries.CreateBlogPost(ctx, sqlc.CreateBlogPostParams{
		Title:              title,
		Slug:               slug,
		Excerpt:            excerpt,
		Body:               body,
		FeaturedImageUrl:   sql.NullString{String: featuredURL, Valid: featuredURL != ""},
		FeaturedImageAlt:   sql.NullString{String: featuredAlt, Valid: featuredAlt != ""},
		CategoryID:         categoryID,
		AuthorID:           authorID,
		MetaDescription:    sql.NullString{String: metaDesc, Valid: metaDesc != ""},
		ReadingTimeMinutes: sql.NullInt64{Int64: readingTime, Valid: readingTime > 0},
		Status:             status,
		PublishedAt:        publishedAt,
	})
	if err != nil {
		h.logger.Error("failed to create blog post", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Handle tags
	tagIDs := c.Request().Form["tag_ids"]
	for _, tagIDStr := range tagIDs {
		tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
		if tagID > 0 {
			h.queries.AddTagToPost(ctx, sqlc.AddTagToPostParams{
				BlogPostID: post.ID,
				BlogTagID:  tagID,
			})
		}
	}

	// Handle products
	productIDs := c.Request().Form["product_ids"]
	for i, pidStr := range productIDs {
		pid, _ := strconv.ParseInt(pidStr, 10, 64)
		if pid > 0 {
			h.queries.AddProductToPost(ctx, sqlc.AddProductToPostParams{
				BlogPostID:   post.ID,
				ProductID:    pid,
				DisplayOrder: sql.NullInt64{Int64: int64(i), Valid: true},
			})
		}
	}

	h.cache.DeleteByPrefix("page:blog")
	logActivity(c, "created", "blog_post", 0, title, "Created blog_post '%s'", title)
	return c.Redirect(http.StatusSeeOther, "/admin/blog/posts")
}

func (h *BlogPostsHandler) Edit(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	post, err := h.queries.GetBlogPost(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	}

	categories, _ := h.queries.ListBlogCategories(ctx)
	authors, _ := h.queries.ListBlogAuthors(ctx)
	allTags, _ := h.queries.ListAllBlogTags(ctx)
	postTags, _ := h.queries.GetPostTagsByPostID(ctx, id)
	postProducts, _ := h.queries.GetPostProductsByPostID(ctx, id)

	return c.Render(http.StatusOK, "admin/pages/blog_post_form.html", map[string]interface{}{
		"Title":        "Edit Blog Post",
		"FormAction":   fmt.Sprintf("/admin/blog/posts/%d", id),
		"Item":         post,
		"Categories":   categories,
		"Authors":      authors,
		"AllTags":      allTags,
		"PostTags":     postTags,
		"PostProducts": postProducts,
	})
}

func (h *BlogPostsHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	existing, err := h.queries.GetBlogPost(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Post not found")
	}

	title := c.FormValue("title")
	slug := c.FormValue("slug")
	if slug == "" {
		slug = makeSlug(title)
	}

	body := c.FormValue("body")
	readingTime, _ := strconv.ParseInt(c.FormValue("reading_time_minutes"), 10, 64)
	if readingTime == 0 {
		readingTime = calculateReadingTime(body)
	}

	categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)
	authorID, _ := strconv.ParseInt(c.FormValue("author_id"), 10, 64)

	featuredURL := c.FormValue("featured_image_url")
	featuredAlt := c.FormValue("featured_image_alt")
	metaDesc := c.FormValue("meta_description")
	excerpt := c.FormValue("excerpt")
	status := c.FormValue("status")

	publishedAt := existing.PublishedAt
	if status == "published" && !existing.PublishedAt.Valid {
		pubStr := c.FormValue("published_at")
		if t, err := time.Parse("2006-01-02T15:04", pubStr); err == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		} else {
			publishedAt = sql.NullTime{Time: time.Now(), Valid: true}
		}
	}

	_, err = h.queries.UpdateBlogPost(ctx, sqlc.UpdateBlogPostParams{
		ID:                 id,
		Title:              title,
		Slug:               slug,
		Excerpt:            excerpt,
		Body:               body,
		FeaturedImageUrl:   sql.NullString{String: featuredURL, Valid: featuredURL != ""},
		FeaturedImageAlt:   sql.NullString{String: featuredAlt, Valid: featuredAlt != ""},
		CategoryID:         categoryID,
		AuthorID:           authorID,
		MetaDescription:    sql.NullString{String: metaDesc, Valid: metaDesc != ""},
		ReadingTimeMinutes: sql.NullInt64{Int64: readingTime, Valid: readingTime > 0},
		Status:             status,
		PublishedAt:        publishedAt,
	})
	if err != nil {
		h.logger.Error("failed to update blog post", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Update tags: clear and re-add
	h.queries.ClearPostTags(ctx, id)
	tagIDs := c.Request().Form["tag_ids"]
	for _, tagIDStr := range tagIDs {
		tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
		if tagID > 0 {
			h.queries.AddTagToPost(ctx, sqlc.AddTagToPostParams{
				BlogPostID: id,
				BlogTagID:  tagID,
			})
		}
	}

	// Update products: clear and re-add
	h.queries.ClearPostProducts(ctx, id)
	productIDs := c.Request().Form["product_ids"]
	for i, pidStr := range productIDs {
		pid, _ := strconv.ParseInt(pidStr, 10, 64)
		if pid > 0 {
			h.queries.AddProductToPost(ctx, sqlc.AddProductToPostParams{
				BlogPostID:   id,
				ProductID:    pid,
				DisplayOrder: sql.NullInt64{Int64: int64(i), Valid: true},
			})
		}
	}

	h.cache.DeleteByPrefix("page:blog")
	logActivity(c, "updated", "blog_post", id, title, "Updated blog_post '%s'", title)
	return c.Redirect(http.StatusSeeOther, "/admin/blog/posts")
}

func (h *BlogPostsHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.queries.ClearPostProducts(c.Request().Context(), id)
	h.queries.ClearPostTags(c.Request().Context(), id)
	if err := h.queries.DeleteBlogPost(c.Request().Context(), id); err != nil {
		h.logger.Error("failed to delete blog post", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	h.cache.DeleteByPrefix("page:blog")
	logActivity(c, "deleted", "blog_post", id, "", "Deleted blog_post #%d", id)
	return c.NoContent(http.StatusOK)
}

func (h *BlogPostsHandler) SearchProducts(c echo.Context) error {
	q := strings.TrimSpace(c.QueryParam("_product_search"))
	if q == "" {
		return c.Render(http.StatusOK, "admin/partials/product_suggestions.html", map[string]interface{}{
			"Products": nil,
			"Query":    "",
		})
	}
	products, _ := h.queries.SearchPublishedProducts(c.Request().Context(), "%"+q+"%")
	return c.Render(http.StatusOK, "admin/partials/product_suggestions.html", map[string]interface{}{
		"Products": products,
		"Query":    q,
	})
}
