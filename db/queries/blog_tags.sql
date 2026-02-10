-- name: ListAllBlogTags :many
SELECT * FROM blog_tags ORDER BY name;

-- name: GetBlogTag :one
SELECT * FROM blog_tags WHERE id = ?;

-- name: GetBlogTagBySlug :one
SELECT * FROM blog_tags WHERE slug = ?;

-- name: CreateBlogTag :one
INSERT INTO blog_tags (name, slug)
VALUES (?, ?)
RETURNING *;

-- name: UpdateBlogTag :one
UPDATE blog_tags SET
    name = ?,
    slug = ?
WHERE id = ?
RETURNING *;

-- name: DeleteBlogTag :exec
DELETE FROM blog_tags WHERE id = ?;

-- name: SearchBlogTags :many
SELECT * FROM blog_tags WHERE name LIKE ? ORDER BY name LIMIT 10;
