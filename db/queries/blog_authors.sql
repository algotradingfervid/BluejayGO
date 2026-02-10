-- name: ListBlogAuthors :many
SELECT * FROM blog_authors ORDER BY sort_order ASC, name ASC;

-- name: GetBlogAuthor :one
SELECT * FROM blog_authors WHERE id = ? LIMIT 1;

-- name: GetBlogAuthorBySlug :one
SELECT * FROM blog_authors WHERE slug = ? LIMIT 1;

-- name: CreateBlogAuthor :one
INSERT INTO blog_authors (name, slug, title, bio, avatar_url, linkedin_url, email, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateBlogAuthor :one
UPDATE blog_authors SET name = ?, slug = ?, title = ?, bio = ?, avatar_url = ?, linkedin_url = ?, email = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteBlogAuthor :exec
DELETE FROM blog_authors WHERE id = ?;
