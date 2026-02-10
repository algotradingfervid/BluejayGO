-- name: ListPublishedPosts :many
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
SELECT COUNT(*) FROM blog_posts
WHERE status = 'published' AND published_at IS NOT NULL;

-- name: ListPublishedPostsByCategory :many
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
SELECT COUNT(*) FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
WHERE bp.status = 'published'
    AND bp.published_at IS NOT NULL
    AND bc.slug = ?;

-- name: GetPublishedPostBySlug :one
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
SELECT bt.id, bt.name, bt.slug
FROM blog_tags bt
INNER JOIN blog_post_tags bpt ON bt.id = bpt.blog_tag_id
WHERE bpt.blog_post_id = ?
ORDER BY bt.name;

-- name: GetRelatedPosts :many
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

-- name: ListAllBlogPosts :many
SELECT
    bp.id, bp.title, bp.slug, bp.status, bp.category_id, bc.name AS category_name,
    bp.author_id, ba.name AS author_name, bp.published_at, bp.created_at, bp.updated_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
ORDER BY bp.created_at DESC;

-- name: ListBlogPostsAdminFiltered :many
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.featured_image_url, bp.featured_image_alt,
    bp.status, bp.category_id, bc.name AS category_name,
    bp.author_id, ba.name AS author_name, ba.avatar_url AS author_avatar,
    bp.reading_time_minutes, bp.published_at, bp.created_at, bp.updated_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE bp.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE bp.category_id = @filter_category END)
    AND (CASE WHEN @filter_author = 0 THEN 1 ELSE bp.author_id = @filter_author END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (bp.title LIKE '%' || @filter_search || '%' OR bp.slug LIKE '%' || @filter_search || '%') END)
ORDER BY bp.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountBlogPostsAdminFiltered :one
SELECT COUNT(*) FROM blog_posts bp
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE bp.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE bp.category_id = @filter_category END)
    AND (CASE WHEN @filter_author = 0 THEN 1 ELSE bp.author_id = @filter_author END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (bp.title LIKE '%' || @filter_search || '%' OR bp.slug LIKE '%' || @filter_search || '%') END);

-- name: GetBlogPost :one
SELECT * FROM blog_posts WHERE id = ?;

-- name: CreateBlogPost :one
INSERT INTO blog_posts (
    title, slug, excerpt, body, featured_image_url, featured_image_alt,
    category_id, author_id, meta_description, reading_time_minutes,
    status, published_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: UpdateBlogPost :one
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
DELETE FROM blog_posts WHERE id = ?;

-- name: AddTagToPost :exec
INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
VALUES (?, ?)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromPost :exec
DELETE FROM blog_post_tags
WHERE blog_post_id = ? AND blog_tag_id = ?;

-- name: ClearPostTags :exec
DELETE FROM blog_post_tags WHERE blog_post_id = ?;

-- name: GetPostProductsByPostID :many
SELECT p.id, p.name, p.slug, p.primary_image, p.tagline, pc.slug AS category_slug, pc.name AS category_name
FROM products p
INNER JOIN blog_post_products bpp ON p.id = bpp.product_id
INNER JOIN product_categories pc ON p.category_id = pc.id
WHERE bpp.blog_post_id = ?
ORDER BY bpp.display_order, p.name;

-- name: AddProductToPost :exec
INSERT INTO blog_post_products (blog_post_id, product_id, display_order)
VALUES (?, ?, ?)
ON CONFLICT DO NOTHING;

-- name: ClearPostProducts :exec
DELETE FROM blog_post_products WHERE blog_post_id = ?;

-- name: SearchPublishedProducts :many
SELECT id, name, slug, primary_image FROM products
WHERE status = 'published' AND name LIKE ?
ORDER BY name LIMIT 10;

-- name: ListLatestPublishedPosts :many
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
