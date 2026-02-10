-- name: ListWhitepaperTopics :many
SELECT * FROM whitepaper_topics ORDER BY sort_order ASC, name ASC;

-- name: GetWhitepaperTopic :one
SELECT * FROM whitepaper_topics WHERE id = ? LIMIT 1;

-- name: GetWhitepaperTopicBySlug :one
SELECT * FROM whitepaper_topics WHERE slug = ? LIMIT 1;

-- name: CreateWhitepaperTopic :one
INSERT INTO whitepaper_topics (name, slug, color_hex, icon, description, sort_order)
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateWhitepaperTopic :one
UPDATE whitepaper_topics SET name = ?, slug = ?, color_hex = ?, icon = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteWhitepaperTopic :exec
DELETE FROM whitepaper_topics WHERE id = ?;
