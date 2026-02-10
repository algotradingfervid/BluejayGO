-- ====================================================================
-- HOMEPAGE HERO
-- ====================================================================

-- name: GetActiveHero :one
SELECT * FROM homepage_hero
WHERE is_active = 1
ORDER BY display_order ASC
LIMIT 1;

-- name: ListAllHeroes :many
SELECT * FROM homepage_hero ORDER BY display_order ASC;

-- name: GetHero :one
SELECT * FROM homepage_hero WHERE id = ?;

-- name: CreateHero :one
INSERT INTO homepage_hero (
    headline, subheadline, badge_text,
    primary_cta_text, primary_cta_url,
    secondary_cta_text, secondary_cta_url,
    background_image, is_active, display_order
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateHero :exec
UPDATE homepage_hero
SET headline = ?, subheadline = ?, badge_text = ?,
    primary_cta_text = ?, primary_cta_url = ?,
    secondary_cta_text = ?, secondary_cta_url = ?,
    background_image = ?, is_active = ?, display_order = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteHero :exec
DELETE FROM homepage_hero WHERE id = ?;

-- ====================================================================
-- HOMEPAGE STATS
-- ====================================================================

-- name: ListActiveStats :many
SELECT * FROM homepage_stats
WHERE is_active = 1
ORDER BY display_order ASC;

-- name: ListAllStats :many
SELECT * FROM homepage_stats ORDER BY display_order ASC;

-- name: GetStat :one
SELECT * FROM homepage_stats WHERE id = ?;

-- name: CreateStat :one
INSERT INTO homepage_stats (stat_value, stat_label, display_order, is_active)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateStat :exec
UPDATE homepage_stats
SET stat_value = ?, stat_label = ?, display_order = ?, is_active = ?
WHERE id = ?;

-- name: DeleteStat :exec
DELETE FROM homepage_stats WHERE id = ?;

-- ====================================================================
-- HOMEPAGE TESTIMONIALS
-- ====================================================================

-- name: ListActiveTestimonialsHomepage :many
SELECT * FROM homepage_testimonials
WHERE is_active = 1
ORDER BY display_order ASC;

-- name: ListAllTestimonialsHomepage :many
SELECT * FROM homepage_testimonials ORDER BY display_order ASC;

-- name: GetTestimonialHomepage :one
SELECT * FROM homepage_testimonials WHERE id = ?;

-- name: CreateTestimonialHomepage :one
INSERT INTO homepage_testimonials (
    quote, author_name, author_title, author_company,
    author_image, rating, display_order, is_active
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTestimonialHomepage :exec
UPDATE homepage_testimonials
SET quote = ?, author_name = ?, author_title = ?, author_company = ?,
    author_image = ?, rating = ?, display_order = ?, is_active = ?
WHERE id = ?;

-- name: DeleteTestimonialHomepage :exec
DELETE FROM homepage_testimonials WHERE id = ?;

-- ====================================================================
-- HOMEPAGE CTA
-- ====================================================================

-- name: GetActiveCTA :one
SELECT * FROM homepage_cta
WHERE is_active = 1
LIMIT 1;

-- name: ListAllCTAs :many
SELECT * FROM homepage_cta ORDER BY id ASC;

-- name: GetCTA :one
SELECT * FROM homepage_cta WHERE id = ?;

-- name: CreateCTA :one
INSERT INTO homepage_cta (
    headline, description,
    primary_cta_text, primary_cta_url,
    secondary_cta_text, secondary_cta_url,
    background_style, is_active
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateCTA :exec
UPDATE homepage_cta
SET headline = ?, description = ?,
    primary_cta_text = ?, primary_cta_url = ?,
    secondary_cta_text = ?, secondary_cta_url = ?,
    background_style = ?, is_active = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteCTA :exec
DELETE FROM homepage_cta WHERE id = ?;
