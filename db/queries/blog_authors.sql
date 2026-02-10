-- ====================================================================
-- BLOG AUTHORS QUERIES
-- ====================================================================
-- This file manages blog author profiles used for article attribution.
-- Authors can be team members, guest writers, or external contributors.
--
-- Managed entity:
-- - blog_authors: author profiles with bio, avatar, social links
--
-- Note: slug field is used for author archive pages (e.g., /blog/author/{slug})
-- ====================================================================

-- name: ListBlogAuthors :many
-- sqlc annotation: :many returns slice of blog_authors rows
-- Purpose: Lists all blog authors for admin management or author selection
-- Parameters: none
-- Return type: slice of all blog_authors rows
-- ORDER BY logic:
--   - Primary sort: sort_order ASC (allows manual ordering)
--   - Secondary sort: name ASC (alphabetical fallback if sort_order is same)
SELECT * FROM blog_authors ORDER BY sort_order ASC, name ASC;

-- name: GetBlogAuthor :one
-- sqlc annotation: :one returns single blog_authors row or error
-- Purpose: Retrieves specific author by ID for editing
-- Parameters:
--   1. id (INTEGER): author primary key
-- Return type: single complete blog_authors row
SELECT * FROM blog_authors WHERE id = ? LIMIT 1;

-- name: GetBlogAuthorBySlug :one
-- sqlc annotation: :one returns single blog_authors row or error
-- Purpose: Retrieves author by URL slug for public author archive pages
-- Parameters:
--   1. slug (TEXT): URL-safe author identifier (e.g., "john-doe")
-- Return type: single blog_authors row
-- Note: slug should be UNIQUE via database constraint to prevent duplicates
SELECT * FROM blog_authors WHERE slug = ? LIMIT 1;

-- name: CreateBlogAuthor :one
-- sqlc annotation: :one returns the created author row
-- Purpose: Creates a new blog author profile
-- Parameters (8 positional):
--   1. name (TEXT): author's full name
--   2. slug (TEXT): URL-safe identifier (must be unique)
--   3. title (TEXT): job title or role (e.g., "Senior Editor")
--   4. bio (TEXT): author biography for byline display
--   5. avatar_url (TEXT): profile image URL
--   6. linkedin_url (TEXT): LinkedIn profile link (optional)
--   7. email (TEXT): contact email (optional, may not be publicly displayed)
--   8. sort_order (INTEGER): display order in author lists
-- Return type: complete inserted row with generated ID and timestamps
INSERT INTO blog_authors (name, slug, title, bio, avatar_url, linkedin_url, email, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateBlogAuthor :one
-- sqlc annotation: :one returns the updated author row
-- Purpose: Updates an existing blog author profile
-- Parameters (9 positional):
--   1-8. updated field values (name, slug, title, bio, avatar_url, linkedin_url, email, sort_order)
--   9. id (INTEGER): which author to update (WHERE clause)
-- Return type: updated blog_authors row
-- Note: updated_at explicitly set to CURRENT_TIMESTAMP to track modifications
UPDATE blog_authors SET name = ?, slug = ?, title = ?, bio = ?, avatar_url = ?, linkedin_url = ?, email = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteBlogAuthor :exec
-- sqlc annotation: :exec returns no data, only error or success
-- Purpose: Permanently removes a blog author
-- Parameters:
--   1. id (INTEGER): author to delete
-- Return type: none
-- Warning: This will fail if author has associated blog posts (foreign key constraint)
--          Consider soft-delete or reassigning posts before deletion
DELETE FROM blog_authors WHERE id = ?;
