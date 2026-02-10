-- ====================================================================
-- BLOG TAGS QUERIES
-- ====================================================================
-- This file manages blog tags (keywords/topics) used for cross-referencing
-- and organizing blog posts. Tags enable many-to-many relationships with
-- posts (one post can have multiple tags, one tag can be on multiple posts).
--
-- Managed entity:
-- - blog_tags: tag definitions (name, slug)
--
-- Note: Tags are linked to posts via blog_post_tags junction table
-- ====================================================================

-- name: ListAllBlogTags :many
-- sqlc annotation: :many returns slice of all blog tags
-- Purpose: Lists all available tags for admin management or tag selection
-- Parameters: none
-- Return type: slice of blog_tags rows
-- ORDER BY name: alphabetical sorting for easier browsing
SELECT * FROM blog_tags ORDER BY name;

-- name: GetBlogTag :one
-- sqlc annotation: :one returns single blog tag by ID
-- Purpose: Retrieves specific tag for editing
-- Parameters:
--   1. id (INTEGER): tag primary key
-- Return type: single blog_tags row
SELECT * FROM blog_tags WHERE id = ?;

-- name: GetBlogTagBySlug :one
-- sqlc annotation: :one returns single blog tag by slug
-- Purpose: Retrieves tag by URL slug (for tag archive pages if implemented)
-- Parameters:
--   1. slug (TEXT): URL-safe tag identifier
-- Return type: single blog_tags row
-- Note: slug should be UNIQUE via database constraint
SELECT * FROM blog_tags WHERE slug = ?;

-- name: CreateBlogTag :one
-- sqlc annotation: :one returns the created tag
-- Purpose: Creates a new blog tag
-- Parameters (2 positional):
--   1. name (TEXT): display name of tag
--   2. slug (TEXT): URL-safe identifier (must be unique)
-- Return type: complete inserted row with ID
INSERT INTO blog_tags (name, slug)
VALUES (?, ?)
RETURNING *;

-- name: UpdateBlogTag :one
-- sqlc annotation: :one returns the updated tag
-- Purpose: Updates an existing blog tag
-- Parameters (3 positional):
--   1. name (TEXT): updated display name
--   2. slug (TEXT): updated slug
--   3. id (INTEGER): which tag to update (WHERE clause)
-- Return type: updated blog_tags row
UPDATE blog_tags SET
    name = ?,
    slug = ?
WHERE id = ?
RETURNING *;

-- name: DeleteBlogTag :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Permanently removes a blog tag
-- Parameters:
--   1. id (INTEGER): tag to delete
-- Return type: none
-- Note: May fail if tag is still associated with posts (foreign key constraint)
--       Consider removing tag associations first or using ON DELETE CASCADE
DELETE FROM blog_tags WHERE id = ?;

-- name: SearchBlogTags :many
-- sqlc annotation: :many returns filtered tags for autocomplete
-- Purpose: Searches tags by partial name match (for typeahead/autocomplete UI)
-- Parameters:
--   1. search_pattern (TEXT): LIKE pattern (e.g., "%java%")
-- Return type: slice of blog_tags (limited to 10 results)
-- LIMIT 10: restricts results for autocomplete dropdown performance
SELECT * FROM blog_tags WHERE name LIKE ? ORDER BY name LIMIT 10;
