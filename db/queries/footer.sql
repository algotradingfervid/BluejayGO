-- name: UpdateFooterSettings :exec
UPDATE settings
SET footer_columns = ?,
    footer_bg_style = ?,
    footer_show_social = ?,
    footer_social_style = ?,
    footer_copyright = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: ListFooterColumnItems :many
SELECT * FROM footer_column_items ORDER BY column_index, sort_order;

-- name: GetFooterColumnItem :one
SELECT * FROM footer_column_items WHERE id = ? LIMIT 1;

-- name: CreateFooterColumnItem :one
INSERT INTO footer_column_items (column_index, type, heading, content, sort_order)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateFooterColumnItem :exec
UPDATE footer_column_items
SET column_index = ?, type = ?, heading = ?, content = ?, sort_order = ?
WHERE id = ?;

-- name: DeleteFooterColumnItem :exec
DELETE FROM footer_column_items WHERE id = ?;

-- name: DeleteFooterColumnItemsByIndex :exec
DELETE FROM footer_column_items WHERE column_index >= ?;

-- name: ListFooterLinks :many
SELECT * FROM footer_links WHERE column_item_id = ? ORDER BY sort_order;

-- name: ListAllFooterLinks :many
SELECT * FROM footer_links ORDER BY column_item_id, sort_order;

-- name: CreateFooterLink :one
INSERT INTO footer_links (column_item_id, label, url, sort_order)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteFooterLinksByColumnItem :exec
DELETE FROM footer_links WHERE column_item_id = ?;

-- name: ListFooterLegalLinks :many
SELECT * FROM footer_legal_links ORDER BY sort_order;

-- name: CreateFooterLegalLink :one
INSERT INTO footer_legal_links (label, url, sort_order)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteAllFooterLegalLinks :exec
DELETE FROM footer_legal_links;
