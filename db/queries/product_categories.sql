-- ====================================================================
-- PRODUCT CATEGORIES QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing product category taxonomy.
--
-- Entity: product_categories table
-- Purpose: Organize products into logical groups (e.g., "Hardware", "Software", "Services")
-- Related: products table references product_categories via category_id foreign key
--
-- Features:
--   - Hierarchical organization of product catalog
--   - Category pages with custom imagery and descriptions
--   - Icon support for visual navigation
--   - Custom display ordering
-- ====================================================================

-- name: ListProductCategories :many
-- Retrieves all product categories ordered by display priority, then alphabetically.
--
-- Parameters: none
-- Returns: []ProductCategory - Array of all category records
--
-- Sorting logic:
--   1. sort_order ASC - Custom display order (lower numbers appear first)
--   2. name ASC - Alphabetical fallback for same sort_order values
--
-- Use case: Category navigation menu, product filters, admin category listing
SELECT * FROM product_categories ORDER BY sort_order ASC, name ASC;

-- name: GetProductCategory :one
-- Retrieves a single product category by its primary key ID.
--
-- Parameters:
--   $1 (INTEGER) - category ID
-- Returns: ProductCategory - Single category record or error if not found
--
-- Use case: Editing a specific category, fetching category details
SELECT * FROM product_categories WHERE id = ? LIMIT 1;

-- name: GetProductCategoryBySlug :one
-- Retrieves a single product category by its URL-safe slug identifier.
--
-- Parameters:
--   $1 (TEXT) - category slug (e.g., "hardware", "software", "services")
-- Returns: ProductCategory - Single category record or error if not found
--
-- Use case: Frontend category page routing, filtering products by category URL
-- Note: Slugs should be unique (enforced by database constraint)
SELECT * FROM product_categories WHERE slug = ? LIMIT 1;

-- name: CreateProductCategory :one
-- Creates a new product category.
--
-- Parameters:
--   $1 (TEXT) - name: Display name of the category
--   $2 (TEXT) - slug: URL-safe identifier
--   $3 (TEXT) - description: Category description for category page (optional)
--   $4 (TEXT) - icon: Icon identifier or CSS class (optional)
--   $5 (TEXT) - image_url: Header/banner image for category page (optional)
--   $6 (INTEGER) - sort_order: Display position (lower = higher priority)
--
-- Returns: ProductCategory - The newly created category with auto-generated ID and timestamps
--
-- Note: RETURNING * includes auto-generated created_at, updated_at
INSERT INTO product_categories (name, slug, description, icon, image_url, sort_order)
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateProductCategory :one
-- Updates an existing product category.
--
-- Parameters:
--   $1 (TEXT) - name: Updated display name
--   $2 (TEXT) - slug: Updated URL-safe identifier
--   $3 (TEXT) - description: Updated description
--   $4 (TEXT) - icon: Updated icon identifier
--   $5 (TEXT) - image_url: Updated header image path
--   $6 (INTEGER) - sort_order: Updated display position
--   $7 (INTEGER) - id: Category ID to update
--
-- Returns: ProductCategory - The updated category record
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
UPDATE product_categories SET name = ?, slug = ?, description = ?, icon = ?, image_url = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteProductCategory :exec
-- Permanently deletes a product category.
--
-- Parameters:
--   $1 (INTEGER) - category ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Will fail if products reference this category (foreign key constraint)
-- Note: Reassign or delete products in this category before deletion
DELETE FROM product_categories WHERE id = ?;
