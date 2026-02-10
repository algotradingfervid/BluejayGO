-- ====================================================================
-- BLOG POSTS QUERIES
-- ====================================================================
-- This file contains all queries for managing blog posts, including both
-- public-facing queries (listing published posts) and admin operations
-- (CRUD, filtering, searching). Includes relationship management for
-- tags and products associated with blog posts.
--
-- Managed entities:
-- - blog_posts: main blog content with status, metadata, SEO fields
-- - blog_post_tags: many-to-many relationship with blog_tags
-- - blog_post_products: many-to-many relationship with products
--
-- Key concepts:
-- - status: 'draft' or 'published' (only published shown on frontend)
-- - published_at: must be non-null for public display
-- - JOINs with categories/authors for denormalized display data
-- ====================================================================

-- ====================================================================
-- PUBLIC BLOG POST QUERIES
-- ====================================================================

-- name: ListPublishedPosts :many
-- sqlc annotation: :many returns slice of blog post rows
-- Purpose: Lists published blog posts for main blog index page
-- Parameters (positional):
--   1. LIMIT (INTEGER): number of posts per page
--   2. OFFSET (INTEGER): pagination offset
-- Return type: slice of denormalized blog post rows with category/author details
-- JOINs explained:
--   - INNER JOIN blog_categories: gets category name/slug/color for each post
--   - INNER JOIN blog_authors: gets author name/avatar for byline display
--   Note: INNER JOINs ensure posts without valid category/author are excluded
-- WHERE clause:
--   - status = 'published': only show published posts, hide drafts
--   - published_at IS NOT NULL: additional safety check for scheduled posts
-- ORDER BY published_at DESC: newest posts first (reverse chronological)
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.featured_image_url, bp.featured_image_alt,
    bp.category_id, bc.name AS category_name, bc.slug AS category_slug, bc.color_hex AS category_color,
    bp.author_id, ba.name AS author_name, ba.avatar_url AS author_avatar,
    bp.reading_time_minutes, bp.published_at, bp.created_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.status = 'published' AND bp.published_at IS NOT NULL
ORDER BY bp.published_at DESC
LIMIT ? OFFSET ?;

-- name: CountPublishedPosts :one
-- sqlc annotation: :one returns single integer count
-- Purpose: Counts total published posts for pagination calculations
-- Parameters: none
-- Return type: integer count
-- Note: WHERE clause must match ListPublishedPosts for accurate pagination
SELECT COUNT(*) FROM blog_posts
WHERE status = 'published' AND published_at IS NOT NULL;

-- name: ListPublishedPostsByCategory :many
-- sqlc annotation: :many returns slice of blog post rows
-- Purpose: Lists published posts filtered by category slug (category archive page)
-- Parameters (positional):
--   1. category slug (TEXT): URL-safe category identifier
--   2. LIMIT (INTEGER): posts per page
--   3. OFFSET (INTEGER): pagination offset
-- Return type: slice of denormalized blog post rows
-- WHERE clause:
--   - bp.status = 'published' AND bp.published_at IS NOT NULL: same as main list
--   - bc.slug = ?: filters by category slug (JOIN to blog_categories required)
-- Note: INNER JOIN ensures only valid category slugs return results
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.featured_image_url, bp.featured_image_alt,
    bp.category_id, bc.name AS category_name, bc.slug AS category_slug, bc.color_hex AS category_color,
    bp.author_id, ba.name AS author_name, ba.avatar_url AS author_avatar,
    bp.reading_time_minutes, bp.published_at, bp.created_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.status = 'published'
    AND bp.published_at IS NOT NULL
    AND bc.slug = ?
ORDER BY bp.published_at DESC
LIMIT ? OFFSET ?;

-- name: CountPublishedPostsByCategory :one
-- sqlc annotation: :one returns integer count
-- Purpose: Counts posts in specific category for pagination
-- Parameters:
--   1. category slug (TEXT)
-- Return type: integer count
-- Note: JOIN required to filter by category slug; WHERE must match ListPublishedPostsByCategory
SELECT COUNT(*) FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
WHERE bp.status = 'published'
    AND bp.published_at IS NOT NULL
    AND bc.slug = ?;

-- name: GetPublishedPostBySlug :one
-- sqlc annotation: :one returns single blog post row or error if not found
-- Purpose: Retrieves full published blog post by slug for public post detail page
-- Parameters:
--   1. slug (TEXT): URL-safe post identifier
-- Return type: single denormalized blog post with full body content, SEO metadata
-- SELECT fields:
--   - Full post content including body (HTML)
--   - Category details (name, slug, color) for breadcrumbs/badges
--   - Extended author info (bio, linkedin) for author card display
--   - SEO fields (meta_title, meta_description, og_image)
-- WHERE clause:
--   - bp.slug = ?: exact slug match
--   - bp.status = 'published' AND bp.published_at IS NOT NULL: public posts only
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.body,
    bp.featured_image_url, bp.featured_image_alt,
    bp.category_id, bc.name AS category_name, bc.slug AS category_slug, bc.color_hex AS category_color,
    bp.author_id, ba.name AS author_name, ba.bio AS author_bio,
    ba.avatar_url AS author_avatar, ba.linkedin_url AS author_linkedin,
    bp.reading_time_minutes, bp.published_at, bp.meta_description,
    bp.meta_title, bp.og_image
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.slug = ? AND bp.status = 'published' AND bp.published_at IS NOT NULL;

-- name: GetPostBySlugIncludeDrafts :one
-- sqlc annotation: :one returns single blog post row including drafts
-- Purpose: Retrieves blog post for admin preview (allows viewing draft posts)
-- Parameters:
--   1. slug (TEXT): post slug
-- Return type: same as GetPublishedPostBySlug but without status filter
-- Note: Used for admin preview functionality; no status/published_at filter
--       allows editors to preview unpublished/draft content
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.body,
    bp.featured_image_url, bp.featured_image_alt,
    bp.category_id, bc.name AS category_name, bc.slug AS category_slug, bc.color_hex AS category_color,
    bp.author_id, ba.name AS author_name, ba.bio AS author_bio,
    ba.avatar_url AS author_avatar, ba.linkedin_url AS author_linkedin,
    bp.reading_time_minutes, bp.published_at, bp.meta_description,
    bp.meta_title, bp.og_image
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.slug = ?;

-- name: GetPostTagsByPostID :many
-- sqlc annotation: :many returns slice of blog_tags rows
-- Purpose: Retrieves all tags associated with a specific blog post
-- Parameters:
--   1. blog_post_id (INTEGER): post to get tags for
-- Return type: slice of blog_tags (id, name, slug)
-- JOIN explained:
--   - blog_post_tags is junction table for many-to-many relationship
--   - INNER JOIN ensures only tags actually linked to the post are returned
-- ORDER BY bt.name: alphabetical tag display
SELECT bt.id, bt.name, bt.slug
FROM blog_tags bt
INNER JOIN blog_post_tags bpt ON bt.id = bpt.blog_tag_id
WHERE bpt.blog_post_id = ?
ORDER BY bt.name;

-- name: GetRelatedPosts :many
-- sqlc annotation: :many returns slice of related blog posts
-- Purpose: Retrieves posts from same category for "Related Posts" widget
-- Parameters (positional):
--   1. category_id (INTEGER): category to match
--   2. current_post_id (INTEGER): post to exclude (don't show current post as related)
--   3. LIMIT (INTEGER): max number of related posts (typically 3-5)
-- Return type: slice of minimal blog post data (no body content)
-- WHERE clause:
--   - bp.status = 'published' AND bp.published_at IS NOT NULL: public posts only
--   - bp.category_id = ?: same category as current post
--   - bp.id != ?: exclude current post from results
-- ORDER BY bp.published_at DESC: newest related posts first
SELECT
    bp.id, bp.title, bp.slug, bp.featured_image_url, bp.featured_image_alt,
    bc.name AS category_name, bc.slug AS category_slug, bc.color_hex AS category_color,
    bp.reading_time_minutes, bp.published_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
WHERE bp.status = 'published'
    AND bp.published_at IS NOT NULL
    AND bp.category_id = ?
    AND bp.id != ?
ORDER BY bp.published_at DESC
LIMIT ?;

-- name: GetFeaturedPost :one
-- sqlc annotation: :one returns single featured blog post
-- Purpose: Retrieves most recent published post for homepage/featured display
-- Parameters: none
-- Return type: single blog post with category/author details (no body)
-- Logic: Simply gets the newest published post (highest published_at date)
-- ORDER BY bp.published_at DESC LIMIT 1: most recent post only
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.featured_image_url, bp.featured_image_alt,
    bp.category_id, bc.name AS category_name, bc.slug AS category_slug, bc.color_hex AS category_color,
    bp.author_id, ba.name AS author_name, ba.avatar_url AS author_avatar,
    bp.reading_time_minutes, bp.published_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.status = 'published' AND bp.published_at IS NOT NULL
ORDER BY bp.published_at DESC
LIMIT 1;

-- ====================================================================
-- ADMIN BLOG POST QUERIES
-- ====================================================================

-- name: ListAllBlogPosts :many
-- sqlc annotation: :many returns slice of all blog posts (drafts + published)
-- Purpose: Simple admin list of all posts without filtering/pagination
-- Parameters: none
-- Return type: slice of blog posts with basic metadata (no body/excerpt)
-- Note: Includes all statuses (draft, published); used for simple admin views
--       ORDER BY created_at DESC shows newest posts first
SELECT
    bp.id, bp.title, bp.slug, bp.status, bp.category_id, bc.name AS category_name,
    bp.author_id, ba.name AS author_name, bp.published_at, bp.created_at, bp.updated_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
ORDER BY bp.created_at DESC;

-- name: ListBlogPostsAdminFiltered :many
-- sqlc annotation: :many returns filtered/paginated blog posts for admin table
-- Purpose: Advanced admin list with multiple filter options and pagination
-- Parameters (named using @ prefix):
--   @filter_status (TEXT): filter by status ('draft', 'published', or '' for all)
--   @filter_category (INTEGER): filter by category ID (0 = all categories)
--   @filter_author (INTEGER): filter by author ID (0 = all authors)
--   @filter_search (TEXT): search term for title/slug ('' = no search)
--   @page_limit (INTEGER): posts per page
--   @page_offset (INTEGER): pagination offset
-- Return type: slice of blog posts with full display metadata
-- Complex WHERE clause using CASE statements:
--   - CASE WHEN @filter_status = '' THEN 1: when filter is empty, condition is always true (1)
--   - ELSE bp.status = @filter_status: otherwise, apply the filter
--   - This pattern allows optional filters without complex NULL handling
--   - filter_category/filter_author use 0 as "no filter" sentinel value
--   - filter_search uses LIKE for partial matching in title OR slug
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.featured_image_url, bp.featured_image_alt,
    bp.status, bp.category_id, bc.name AS category_name,
    bp.author_id, ba.name AS author_name, ba.avatar_url AS author_avatar,
    bp.reading_time_minutes, bp.published_at, bp.created_at, bp.updated_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE
    -- Optional status filter: '' = all statuses, otherwise exact match
    (CASE WHEN @filter_status = '' THEN 1 ELSE bp.status = @filter_status END)
    -- Optional category filter: 0 = all categories, otherwise exact ID match
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE bp.category_id = @filter_category END)
    -- Optional author filter: 0 = all authors, otherwise exact ID match
    AND (CASE WHEN @filter_author = 0 THEN 1 ELSE bp.author_id = @filter_author END)
    -- Optional search filter: '' = no search, otherwise LIKE match in title or slug
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (bp.title LIKE '%' || @filter_search || '%' OR bp.slug LIKE '%' || @filter_search || '%') END)
ORDER BY bp.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountBlogPostsAdminFiltered :one
-- sqlc annotation: :one returns integer count for pagination
-- Purpose: Counts total posts matching admin filters
-- Parameters: same filters as ListBlogPostsAdminFiltered (except pagination)
-- Return type: integer count
-- Note: WHERE clause MUST exactly match ListBlogPostsAdminFiltered for accuracy
--       Does NOT include JOINs since we only need count
SELECT COUNT(*) FROM blog_posts bp
WHERE
    -- Exact same filter logic as ListBlogPostsAdminFiltered
    (CASE WHEN @filter_status = '' THEN 1 ELSE bp.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE bp.category_id = @filter_category END)
    AND (CASE WHEN @filter_author = 0 THEN 1 ELSE bp.author_id = @filter_author END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (bp.title LIKE '%' || @filter_search || '%' OR bp.slug LIKE '%' || @filter_search || '%') END);

-- name: GetBlogPost :one
-- sqlc annotation: :one returns single blog post by ID
-- Purpose: Retrieves blog post for admin editing (all statuses)
-- Parameters:
--   1. id (INTEGER): post primary key
-- Return type: complete blog_posts row with all fields
SELECT * FROM blog_posts WHERE id = ?;

-- name: CreateBlogPost :one
-- sqlc annotation: :one returns the created blog post
-- Purpose: Creates a new blog post (draft or published)
-- Parameters (12 positional):
--   1. title (TEXT): post headline
--   2. slug (TEXT): URL-safe identifier (must be unique)
--   3. excerpt (TEXT): short summary for listings
--   4. body (TEXT): full HTML content
--   5. featured_image_url (TEXT): hero image URL
--   6. featured_image_alt (TEXT): image alt text for accessibility
--   7. category_id (INTEGER): foreign key to blog_categories
--   8. author_id (INTEGER): foreign key to blog_authors
--   9. meta_description (TEXT): SEO description
--   10. reading_time_minutes (INTEGER): estimated read time
--   11. status (TEXT): 'draft' or 'published'
--   12. published_at (TIMESTAMP): publication datetime (NULL for drafts)
-- Return type: complete inserted row with ID and timestamps
INSERT INTO blog_posts (
    title, slug, excerpt, body, featured_image_url, featured_image_alt,
    category_id, author_id, meta_description, reading_time_minutes,
    status, published_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateBlogPost :one
-- sqlc annotation: :one returns the updated blog post
-- Purpose: Updates existing blog post (all fields except ID/created_at)
-- Parameters (13 positional):
--   1-12. updated field values (same order as CreateBlogPost)
--   13. id (INTEGER): which post to update (WHERE clause)
-- Return type: updated blog_posts row
-- Note: updated_at explicitly set to CURRENT_TIMESTAMP
UPDATE blog_posts SET
    title = ?,
    slug = ?,
    excerpt = ?,
    body = ?,
    featured_image_url = ?,
    featured_image_alt = ?,
    category_id = ?,
    author_id = ?,
    meta_description = ?,
    reading_time_minutes = ?,
    status = ?,
    published_at = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteBlogPost :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Permanently removes a blog post
-- Parameters:
--   1. id (INTEGER): post to delete
-- Return type: none
-- Note: This will cascade delete related blog_post_tags and blog_post_products
--       if foreign keys are configured with ON DELETE CASCADE
DELETE FROM blog_posts WHERE id = ?;

-- ====================================================================
-- BLOG POST TAGS RELATIONSHIP
-- ====================================================================

-- name: AddTagToPost :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Associates a tag with a blog post (many-to-many relationship)
-- Parameters:
--   1. blog_post_id (INTEGER): post ID
--   2. blog_tag_id (INTEGER): tag ID
-- Return type: none
-- Note: ON CONFLICT DO NOTHING prevents duplicate tag associations
--       (blog_post_tags should have UNIQUE constraint on (blog_post_id, blog_tag_id))
INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
VALUES (?, ?)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromPost :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Removes a specific tag association from a post
-- Parameters:
--   1. blog_post_id (INTEGER)
--   2. blog_tag_id (INTEGER)
-- Return type: none
DELETE FROM blog_post_tags
WHERE blog_post_id = ? AND blog_tag_id = ?;

-- name: ClearPostTags :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Removes all tag associations from a post (used before re-adding tags)
-- Parameters:
--   1. blog_post_id (INTEGER): post to clear tags from
-- Return type: none
-- Note: Typically used when updating post tags (clear then re-add)
DELETE FROM blog_post_tags WHERE blog_post_id = ?;

-- ====================================================================
-- BLOG POST PRODUCTS RELATIONSHIP
-- ====================================================================

-- name: GetPostProductsByPostID :many
-- sqlc annotation: :many returns slice of products associated with a post
-- Purpose: Retrieves products featured/mentioned in a blog post
-- Parameters:
--   1. blog_post_id (INTEGER): post to get products for
-- Return type: slice of products with category info
-- JOINs:
--   - blog_post_products: junction table with display_order
--   - products: product details
--   - product_categories: for category name/slug
-- ORDER BY: display_order (custom order), then name (alphabetical fallback)
SELECT p.id, p.name, p.slug, p.primary_image, p.tagline, pc.slug AS category_slug, pc.name AS category_name
FROM products p
INNER JOIN blog_post_products bpp ON p.id = bpp.product_id
INNER JOIN product_categories pc ON p.category_id = pc.id
WHERE bpp.blog_post_id = ?
ORDER BY bpp.display_order, p.name;

-- name: AddProductToPost :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Associates a product with a blog post with specific display order
-- Parameters:
--   1. blog_post_id (INTEGER)
--   2. product_id (INTEGER)
--   3. display_order (INTEGER): position in product list for this post
-- Return type: none
-- Note: ON CONFLICT DO NOTHING prevents duplicate associations
INSERT INTO blog_post_products (blog_post_id, product_id, display_order)
VALUES (?, ?, ?)
ON CONFLICT DO NOTHING;

-- name: ClearPostProducts :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Removes all product associations from a post
-- Parameters:
--   1. blog_post_id (INTEGER)
-- Return type: none
DELETE FROM blog_post_products WHERE blog_post_id = ?;

-- name: SearchPublishedProducts :many
-- sqlc annotation: :many returns slice of products for search autocomplete
-- Purpose: Searches published products by name for adding to blog posts
-- Parameters:
--   1. search_pattern (TEXT): LIKE pattern (e.g., "%laptop%")
-- Return type: slice of minimal product data (id, name, slug, image)
-- WHERE: status = 'published' ensures only published products can be linked
-- LIMIT 10: restricts results for autocomplete/typeahead UI
SELECT id, name, slug, primary_image FROM products
WHERE status = 'published' AND name LIKE ?
ORDER BY name LIMIT 10;

-- ====================================================================
-- UTILITY QUERIES
-- ====================================================================

-- name: ListLatestPublishedPosts :many
-- sqlc annotation: :many returns recent published posts
-- Purpose: Lists N most recent published posts (for homepage, sidebar widgets)
-- Parameters:
--   1. LIMIT (INTEGER): number of posts to return
-- Return type: slice of blog posts with category/author details
-- Note: Similar to ListPublishedPosts but without pagination offset
--       Used for "Recent Posts" widgets with fixed count
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.featured_image_url, bp.featured_image_alt,
    bp.category_id, bc.name AS category_name, bc.slug AS category_slug, bc.color_hex AS category_color,
    bp.author_id, ba.name AS author_name, ba.avatar_url AS author_avatar,
    bp.reading_time_minutes, bp.published_at, bp.created_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.status = 'published' AND bp.published_at IS NOT NULL
ORDER BY bp.published_at DESC
LIMIT ?;
