-- name: ListPartnerTiers :many
SELECT * FROM partner_tiers ORDER BY sort_order ASC, name ASC;

-- name: GetPartnerTier :one
SELECT * FROM partner_tiers WHERE id = ? LIMIT 1;

-- name: GetPartnerTierBySlug :one
SELECT * FROM partner_tiers WHERE slug = ? LIMIT 1;

-- name: CreatePartnerTier :one
INSERT INTO partner_tiers (name, slug, description, sort_order)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdatePartnerTier :one
UPDATE partner_tiers SET name = ?, slug = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeletePartnerTier :exec
DELETE FROM partner_tiers WHERE id = ?;
