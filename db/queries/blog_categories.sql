-- name: ListBlogCategories :many
SELECT * FROM blog_categories ORDER BY sort_order ASC, name ASC;

-- name: GetBlogCategory :one
SELECT * FROM blog_categories WHERE id = ? LIMIT 1;

-- name: GetBlogCategoryBySlug :one
SELECT * FROM blog_categories WHERE slug = ? LIMIT 1;

-- name: CreateBlogCategory :one
INSERT INTO blog_categories (name, slug, color_hex, description, sort_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateBlogCategory :one
UPDATE blog_categories SET name = ?, slug = ?, color_hex = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteBlogCategory :exec
DELETE FROM blog_categories WHERE id = ?;
