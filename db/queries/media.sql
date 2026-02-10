-- ====================================================================
-- MEDIA FILES QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing uploaded media files
-- (images, PDFs, videos, etc.) in the media library.
--
-- Entity: media_files table
-- Purpose: Track uploaded files with metadata for reuse across the CMS
-- Features: Pagination, search, multiple sort orders, alt text for accessibility
-- ====================================================================

-- name: ListMediaFiles :many
-- Retrieves paginated media files sorted by upload date (newest first).
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Number of records to return per page
--   $2 (INTEGER) - OFFSET: Number of records to skip (for pagination)
-- Returns: []MediaFile - Array of media file records
--
-- Sorting: created_at DESC - Most recently uploaded files appear first
-- Use case: Default media library view showing latest uploads
SELECT * FROM media_files
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListMediaFilesByName :many
-- Retrieves paginated media files sorted alphabetically by filename.
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Number of records per page
--   $2 (INTEGER) - OFFSET: Pagination offset
-- Returns: []MediaFile - Array of media file records
--
-- Sorting: filename ASC - Alphabetical order (A-Z)
-- Use case: User prefers alphabetical file organization
SELECT * FROM media_files
ORDER BY filename ASC
LIMIT ? OFFSET ?;

-- name: ListMediaFilesBySize :many
-- Retrieves paginated media files sorted by file size (largest first).
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Number of records per page
--   $2 (INTEGER) - OFFSET: Pagination offset
-- Returns: []MediaFile - Array of media file records
--
-- Sorting: file_size DESC - Largest files appear first
-- Use case: Identifying large files for storage cleanup or optimization
SELECT * FROM media_files
ORDER BY file_size DESC
LIMIT ? OFFSET ?;

-- name: ListMediaFilesOldest :many
-- Retrieves paginated media files sorted by upload date (oldest first).
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Number of records per page
--   $2 (INTEGER) - OFFSET: Pagination offset
-- Returns: []MediaFile - Array of media file records
--
-- Sorting: created_at ASC - Oldest files appear first
-- Use case: Identifying old/unused files for archival or cleanup
SELECT * FROM media_files
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: SearchMediaFiles :many
-- Searches media files by original filename with pagination.
--
-- Parameters:
--   @search (TEXT) - Search term to match against original_filename
--   @page_limit (INTEGER) - Number of results per page
--   @page_offset (INTEGER) - Pagination offset
-- Returns: []MediaFile - Array of matching media file records
--
-- Search logic: LIKE '%' || @search || '%' - Case-insensitive partial match
-- Note: Searches original_filename (user's uploaded filename) not stored filename
-- Performance: May be slow on large tables without full-text search index
SELECT * FROM media_files
WHERE original_filename LIKE '%' || @search || '%'
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountMediaFiles :one
-- Returns the total count of all media files.
--
-- Parameters: none
-- Returns: INTEGER - Total number of media files in the library
--
-- Use case: Calculating total pages for pagination, displaying library statistics
SELECT COUNT(*) FROM media_files;

-- name: CountMediaFilesSearch :one
-- Returns the count of media files matching a search query.
--
-- Parameters:
--   @search (TEXT) - Search term to match against original_filename
-- Returns: INTEGER - Number of files matching the search
--
-- Use case: Pagination for search results
-- Note: Uses same LIKE pattern as SearchMediaFiles for consistent counts
SELECT COUNT(*) FROM media_files
WHERE original_filename LIKE '%' || @search || '%';

-- name: GetMediaFile :one
-- Retrieves a single media file by its primary key ID.
--
-- Parameters:
--   $1 (INTEGER) - media file ID
-- Returns: MediaFile - Single media file record or error if not found
--
-- Use case: Displaying file details, embedding in content
SELECT * FROM media_files WHERE id = ? LIMIT 1;

-- name: GetMediaFileByPath :one
-- Retrieves a single media file by its storage file path.
--
-- Parameters:
--   $1 (TEXT) - file_path: Relative or absolute path to file on disk
-- Returns: MediaFile - Single media file record or error if not found
--
-- Use case: Checking if a file already exists before upload, resolving paths to records
-- Note: file_path should be unique (enforced by database constraint)
SELECT * FROM media_files WHERE file_path = ? LIMIT 1;

-- name: CreateMediaFile :one
-- Inserts a new media file record after successful upload.
--
-- Parameters:
--   $1 (TEXT) - filename: System-generated unique filename (e.g., uuid.jpg)
--   $2 (TEXT) - original_filename: User's original upload filename
--   $3 (TEXT) - file_path: Relative/absolute path to stored file
--   $4 (INTEGER) - file_size: File size in bytes
--   $5 (TEXT) - mime_type: MIME type (e.g., image/jpeg, application/pdf)
--   $6 (INTEGER) - width: Image width in pixels (NULL for non-images)
--   $7 (INTEGER) - height: Image height in pixels (NULL for non-images)
--   $8 (TEXT) - alt_text: Accessibility alt text for images (optional)
--
-- Returns: MediaFile - The newly created media file record with auto-generated ID and timestamps
--
-- Note: width/height should be extracted from image files during upload processing
INSERT INTO media_files (filename, original_filename, file_path, file_size, mime_type, width, height, alt_text)
VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateMediaFileAltText :exec
-- Updates the alt text for an existing media file (accessibility).
--
-- Parameters:
--   $1 (TEXT) - alt_text: Updated accessibility description
--   $2 (INTEGER) - id: Media file ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Adding or editing alt text for screen readers and SEO
-- Note: Only updates alt_text, not other file metadata
UPDATE media_files SET alt_text = ? WHERE id = ?;

-- name: DeleteMediaFile :exec
-- Permanently deletes a media file record from the database.
--
-- Parameters:
--   $1 (INTEGER) - media file ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: This only deletes the database record, not the physical file on disk
-- Note: Application code should delete the physical file before calling this query
-- Caution: May orphan file references in products, solutions, or other content
DELETE FROM media_files WHERE id = ?;
