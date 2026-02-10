-- ====================================================================
-- INDUSTRIES QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing industry categories.
-- Industries are used to categorize solutions and content by target market
-- or business sector (e.g., Healthcare, Manufacturing, Finance).
--
-- Entity: industries table
-- Related tables: solutions (via solution_industries junction table)
-- ====================================================================

-- name: ListIndustries :many
-- Retrieves all industries ordered by custom sort order, then alphabetically.
--
-- Parameters: none
-- Returns: []Industry - Array of all industry records
--
-- Sorting logic:
--   1. sort_order ASC - Custom display order (lower numbers appear first)
--   2. name ASC - Alphabetical fallback for same sort_order values
--
-- Use case: Display industries in navigation, filters, or admin listing
SELECT * FROM industries ORDER BY sort_order ASC, name ASC;

-- name: GetIndustry :one
-- Retrieves a single industry by its primary key ID.
--
-- Parameters:
--   $1 (INTEGER) - industry ID
-- Returns: Industry - Single industry record or error if not found
--
-- Note: sqlc annotation :one expects exactly one row; will error if 0 or >1 rows returned
SELECT * FROM industries WHERE id = ? LIMIT 1;

-- name: GetIndustryBySlug :one
-- Retrieves a single industry by its URL-safe slug identifier.
--
-- Parameters:
--   $1 (TEXT) - industry slug (e.g., "healthcare", "manufacturing")
-- Returns: Industry - Single industry record or error if not found
--
-- Use case: Frontend routing, displaying industry-specific content
-- Note: Slugs should be unique (enforced by database constraint)
SELECT * FROM industries WHERE slug = ? LIMIT 1;

-- name: CreateIndustry :one
-- Inserts a new industry record and returns the created record.
--
-- Parameters:
--   $1 (TEXT) - name: Display name of the industry
--   $2 (TEXT) - slug: URL-safe identifier
--   $3 (TEXT) - icon: Icon identifier or CSS class (optional)
--   $4 (TEXT) - description: Industry description (optional)
--   $5 (INTEGER) - sort_order: Display position (lower = higher priority)
--
-- Returns: Industry - The newly created industry record with auto-generated ID and timestamps
--
-- Note: RETURNING * returns all columns including auto-generated created_at, updated_at
INSERT INTO industries (name, slug, icon, description, sort_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateIndustry :one
-- Updates an existing industry record and returns the updated record.
--
-- Parameters:
--   $1 (TEXT) - name: Updated display name
--   $2 (TEXT) - slug: Updated URL-safe identifier
--   $3 (TEXT) - icon: Updated icon identifier
--   $4 (TEXT) - description: Updated description
--   $5 (INTEGER) - sort_order: Updated display position
--   $6 (INTEGER) - id: Primary key of industry to update
--
-- Returns: Industry - The updated industry record
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP to track last modification
UPDATE industries SET name = ?, slug = ?, icon = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteIndustry :exec
-- Permanently deletes an industry record.
--
-- Parameters:
--   $1 (INTEGER) - industry ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count, no data
--
-- WARNING: This is a hard delete. Consider adding soft delete (is_active flag) for production.
-- Note: May fail if foreign key constraints exist (e.g., solutions referencing this industry)
DELETE FROM industries WHERE id = ?;
