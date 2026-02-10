-- ====================================================================
-- BLOG CATEGORIES QUERIES
-- ====================================================================
-- This file manages blog category taxonomy for organizing posts into
-- distinct content verticals (e.g., "Technology", "Industry News", "Tutorials").
--
-- Managed entity:
-- - blog_categories: post classification with color coding for UI
--
-- Note: slug field enables category archive pages (e.g., /blog/category/{slug})
--       color_hex allows color-coded category badges in UI
-- ====================================================================

-- name: ListBlogCategories :many
-- sqlc annotation: :many returns slice of blog_categories rows
-- Purpose: Lists all blog categories for admin management or post categorization
-- Parameters: none
-- Return type: slice of all blog_categories rows
-- ORDER BY logic:
--   - Primary sort: sort_order ASC (custom display order)
--   - Secondary sort: name ASC (alphabetical fallback)
SELECT * FROM blog_categories ORDER BY sort_order ASC, name ASC;

-- name: GetBlogCategory :one
-- sqlc annotation: :one returns single blog_categories row or error
-- Purpose: Retrieves specific category by ID for editing
-- Parameters:
--   1. id (INTEGER): category primary key
-- Return type: single complete blog_categories row
SELECT * FROM blog_categories WHERE id = ? LIMIT 1;

-- name: GetBlogCategoryBySlug :one
-- sqlc annotation: :one returns single blog_categories row or error
-- Purpose: Retrieves category by URL slug for public category archive pages
-- Parameters:
--   1. slug (TEXT): URL-safe category identifier (e.g., "technology")
-- Return type: single blog_categories row
-- Note: slug should be UNIQUE via database constraint
SELECT * FROM blog_categories WHERE slug = ? LIMIT 1;

-- name: CreateBlogCategory :one
-- sqlc annotation: :one returns the created category row
-- Purpose: Creates a new blog category
-- Parameters (5 positional):
--   1. name (TEXT): category display name
--   2. slug (TEXT): URL-safe identifier (must be unique)
--   3. color_hex (TEXT): hex color for UI badges (e.g., "#3B82F6")
--   4. description (TEXT): category description for archive page
--   5. sort_order (INTEGER): display order in category lists
-- Return type: complete inserted row with generated ID and timestamps
INSERT INTO blog_categories (name, slug, color_hex, description, sort_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateBlogCategory :one
-- sqlc annotation: :one returns the updated category row
-- Purpose: Updates an existing blog category
-- Parameters (6 positional):
--   1-5. updated field values (name, slug, color_hex, description, sort_order)
--   6. id (INTEGER): which category to update (WHERE clause)
-- Return type: updated blog_categories row
-- Note: updated_at explicitly set to CURRENT_TIMESTAMP
UPDATE blog_categories SET name = ?, slug = ?, color_hex = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteBlogCategory :exec
-- sqlc annotation: :exec returns no data, only error or success
-- Purpose: Permanently removes a blog category
-- Parameters:
--   1. id (INTEGER): category to delete
-- Return type: none
-- Warning: This will fail if category has associated blog posts (foreign key constraint)
--          Consider reassigning posts before deletion
DELETE FROM blog_categories WHERE id = ?;
