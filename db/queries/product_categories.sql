-- name: ListProductCategories :many
SELECT * FROM product_categories ORDER BY sort_order ASC, name ASC;

-- name: GetProductCategory :one
SELECT * FROM product_categories WHERE id = ? LIMIT 1;

-- name: GetProductCategoryBySlug :one
SELECT * FROM product_categories WHERE slug = ? LIMIT 1;

-- name: CreateProductCategory :one
INSERT INTO product_categories (name, slug, description, icon, image_url, sort_order)
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateProductCategory :one
UPDATE product_categories SET name = ?, slug = ?, description = ?, icon = ?, image_url = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteProductCategory :exec
DELETE FROM product_categories WHERE id = ?;
