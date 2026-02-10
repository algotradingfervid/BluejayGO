-- name: ListMediaFiles :many
SELECT * FROM media_files
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListMediaFilesByName :many
SELECT * FROM media_files
ORDER BY filename ASC
LIMIT ? OFFSET ?;

-- name: ListMediaFilesBySize :many
SELECT * FROM media_files
ORDER BY file_size DESC
LIMIT ? OFFSET ?;

-- name: ListMediaFilesOldest :many
SELECT * FROM media_files
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: SearchMediaFiles :many
SELECT * FROM media_files
WHERE original_filename LIKE '%' || @search || '%'
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountMediaFiles :one
SELECT COUNT(*) FROM media_files;

-- name: CountMediaFilesSearch :one
SELECT COUNT(*) FROM media_files
WHERE original_filename LIKE '%' || @search || '%';

-- name: GetMediaFile :one
SELECT * FROM media_files WHERE id = ? LIMIT 1;

-- name: GetMediaFileByPath :one
SELECT * FROM media_files WHERE file_path = ? LIMIT 1;

-- name: CreateMediaFile :one
INSERT INTO media_files (filename, original_filename, file_path, file_size, mime_type, width, height, alt_text)
VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateMediaFileAltText :exec
UPDATE media_files SET alt_text = ? WHERE id = ?;

-- name: DeleteMediaFile :exec
DELETE FROM media_files WHERE id = ?;
