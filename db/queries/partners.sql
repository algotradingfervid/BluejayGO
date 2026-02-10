-- name: ListPartnersByTier :many
SELECT p.*, pt.name AS tier_name, pt.sort_order AS tier_level
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
WHERE p.is_active = 1
ORDER BY pt.sort_order ASC, p.display_order ASC;

-- name: ListPartnersByTierID :many
SELECT * FROM partners
WHERE tier_id = ? AND is_active = 1
ORDER BY display_order ASC;

-- name: GetPartner :one
SELECT p.*, pt.name AS tier_name, pt.sort_order AS tier_level
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
WHERE p.id = ?;

-- name: CreatePartner :one
INSERT INTO partners (
    name, tier_id, logo_url, icon, website_url, description, display_order
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdatePartner :one
UPDATE partners
SET name = ?, tier_id = ?, logo_url = ?, icon = ?,
    website_url = ?, description = ?, display_order = ?,
    is_active = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeletePartner :exec
DELETE FROM partners WHERE id = ?;

-- name: ListActiveTestimonials :many
SELECT pt.*, p.name AS partner_name, p.logo_url AS partner_logo_url, p.icon AS partner_icon
FROM partner_testimonials pt
JOIN partners p ON pt.partner_id = p.id
WHERE pt.is_active = 1
ORDER BY pt.display_order ASC;

-- name: GetTestimonial :one
SELECT pt.*, p.name AS partner_name
FROM partner_testimonials pt
JOIN partners p ON pt.partner_id = p.id
WHERE pt.id = ?;

-- name: CreateTestimonial :one
INSERT INTO partner_testimonials (
    partner_id, quote, author_name, author_title, display_order
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTestimonial :one
UPDATE partner_testimonials
SET partner_id = ?, quote = ?, author_name = ?,
    author_title = ?, display_order = ?, is_active = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTestimonial :exec
DELETE FROM partner_testimonials WHERE id = ?;

-- name: ListAllPartners :many
SELECT p.*, pt.name AS tier_name
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
ORDER BY pt.sort_order ASC, p.display_order ASC;

-- name: ListFeaturedPartners :many
SELECT p.*, pt.name AS tier_name
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
WHERE p.is_featured = 1 AND p.is_active = 1
ORDER BY p.display_order ASC
LIMIT ?;
