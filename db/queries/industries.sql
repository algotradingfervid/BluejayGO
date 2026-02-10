-- name: ListIndustries :many
SELECT * FROM industries ORDER BY sort_order ASC, name ASC;

-- name: GetIndustry :one
SELECT * FROM industries WHERE id = ? LIMIT 1;

-- name: GetIndustryBySlug :one
SELECT * FROM industries WHERE slug = ? LIMIT 1;

-- name: CreateIndustry :one
INSERT INTO industries (name, slug, icon, description, sort_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateIndustry :one
UPDATE industries SET name = ?, slug = ?, icon = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteIndustry :exec
DELETE FROM industries WHERE id = ?;
