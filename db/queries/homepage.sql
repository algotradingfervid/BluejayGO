-- ====================================================================
-- HOMEPAGE CONTENT QUERIES
-- ====================================================================
-- This file manages all homepage content sections including:
-- - Hero banner (headline, CTAs, background image)
-- - Statistics/metrics display
-- - Testimonials carousel
-- - Call-to-action (CTA) blocks
--
-- Managed entities:
-- - homepage_hero: hero banner variations (only one active at a time)
-- - homepage_stats: key metrics/statistics
-- - homepage_testimonials: customer testimonials
-- - homepage_cta: call-to-action sections
--
-- Key concepts:
-- - is_active: controls which variant is currently displayed
-- - display_order: custom sort sequence for multi-item sections
-- ====================================================================

-- ====================================================================
-- HOMEPAGE HERO
-- ====================================================================

-- name: GetActiveHero :one
-- sqlc annotation: :one returns the currently active hero banner
-- Purpose: Retrieves hero content displayed on public homepage
-- Parameters: none
-- Return type: single homepage_hero row
-- WHERE: is_active = 1 (only one hero should be active at a time)
-- ORDER BY display_order ASC LIMIT 1: safety fallback if multiple are accidentally active
SELECT * FROM homepage_hero
WHERE is_active = 1
ORDER BY display_order ASC
LIMIT 1;

-- name: ListAllHeroes :many
-- Purpose: Lists all hero variants (active + inactive) for admin management
-- ORDER BY display_order: custom sort for A/B testing or scheduling
SELECT * FROM homepage_hero ORDER BY display_order ASC;

-- name: GetHero :one
-- Purpose: Retrieves specific hero by ID for editing
SELECT * FROM homepage_hero WHERE id = ?;

-- name: CreateHero :one
-- Purpose: Creates a new hero banner variant
-- Parameters (10 positional):
--   1. headline (TEXT): main hero headline
--   2. subheadline (TEXT): supporting text
--   3. badge_text (TEXT): optional badge/announcement text
--   4. primary_cta_text (TEXT): primary button text
--   5. primary_cta_url (TEXT): primary button link
--   6. secondary_cta_text (TEXT): optional secondary button text
--   7. secondary_cta_url (TEXT): secondary button link
--   8. background_image (TEXT): hero background/featured image URL
--   9. is_active (BOOLEAN): whether this hero is currently displayed
--   10. display_order (INTEGER): sort position
-- Note: Only one hero should have is_active = 1 at a time
INSERT INTO homepage_hero (
    headline, subheadline, badge_text,
    primary_cta_text, primary_cta_url,
    secondary_cta_text, secondary_cta_url,
    background_image, is_active, display_order
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateHero :exec
-- Purpose: Updates an existing hero banner
-- Parameters (11 positional): same as CreateHero + id (WHERE clause)
UPDATE homepage_hero
SET headline = ?, subheadline = ?, badge_text = ?,
    primary_cta_text = ?, primary_cta_url = ?,
    secondary_cta_text = ?, secondary_cta_url = ?,
    background_image = ?, is_active = ?, display_order = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteHero :exec
-- Purpose: Removes a hero banner variant
DELETE FROM homepage_hero WHERE id = ?;

-- ====================================================================
-- HOMEPAGE STATS / METRICS
-- ====================================================================

-- name: ListActiveStats :many
-- Purpose: Retrieves active statistics for public homepage display
-- WHERE: is_active = 1 (allows hiding stats without deleting them)
-- ORDER BY display_order: custom presentation sequence
-- Example stats: "500+ Clients", "10 Years Experience", "99% Satisfaction"
SELECT * FROM homepage_stats
WHERE is_active = 1
ORDER BY display_order ASC;

-- name: ListAllStats :many
-- Purpose: Lists all stats (active + inactive) for admin management
SELECT * FROM homepage_stats ORDER BY display_order ASC;

-- name: GetStat :one
-- Purpose: Retrieves specific stat by ID for editing
SELECT * FROM homepage_stats WHERE id = ?;

-- name: CreateStat :one
-- Purpose: Creates a new homepage statistic
-- Parameters (4 positional):
--   1. stat_value (TEXT): numeric value with unit (e.g., "500+", "10 Years")
--   2. stat_label (TEXT): description (e.g., "Happy Clients", "Experience")
--   3. display_order (INTEGER): sort position (left to right typically)
--   4. is_active (BOOLEAN): whether to display on homepage
INSERT INTO homepage_stats (stat_value, stat_label, display_order, is_active)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateStat :exec
-- Purpose: Updates an existing homepage statistic
-- Parameters (5 positional): same as CreateStat + id
UPDATE homepage_stats
SET stat_value = ?, stat_label = ?, display_order = ?, is_active = ?
WHERE id = ?;

-- name: DeleteStat :exec
-- Purpose: Removes a homepage statistic
DELETE FROM homepage_stats WHERE id = ?;

-- ====================================================================
-- HOMEPAGE TESTIMONIALS
-- ====================================================================

-- name: ListActiveTestimonialsHomepage :many
-- Purpose: Retrieves active testimonials for public homepage carousel/display
-- WHERE: is_active = 1 (allows managing testimonial rotation)
-- ORDER BY display_order: custom sequence for carousel slides
SELECT * FROM homepage_testimonials
WHERE is_active = 1
ORDER BY display_order ASC;

-- name: ListAllTestimonialsHomepage :many
-- Purpose: Lists all testimonials (active + inactive) for admin management
SELECT * FROM homepage_testimonials ORDER BY display_order ASC;

-- name: GetTestimonialHomepage :one
-- Purpose: Retrieves specific testimonial by ID for editing
SELECT * FROM homepage_testimonials WHERE id = ?;

-- name: CreateTestimonialHomepage :one
-- Purpose: Creates a new homepage testimonial
-- Parameters (8 positional):
--   1. quote (TEXT): testimonial text/content
--   2. author_name (TEXT): customer name
--   3. author_title (TEXT): job title
--   4. author_company (TEXT): company/organization name
--   5. author_image (TEXT): headshot/avatar URL
--   6. rating (INTEGER): star rating (typically 1-5)
--   7. display_order (INTEGER): carousel slide order
--   8. is_active (BOOLEAN): whether to display on homepage
INSERT INTO homepage_testimonials (
    quote, author_name, author_title, author_company,
    author_image, rating, display_order, is_active
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTestimonialHomepage :exec
-- Purpose: Updates an existing homepage testimonial
-- Parameters (9 positional): same as CreateTestimonialHomepage + id
UPDATE homepage_testimonials
SET quote = ?, author_name = ?, author_title = ?, author_company = ?,
    author_image = ?, rating = ?, display_order = ?, is_active = ?
WHERE id = ?;

-- name: DeleteTestimonialHomepage :exec
-- Purpose: Removes a homepage testimonial
DELETE FROM homepage_testimonials WHERE id = ?;

-- ====================================================================
-- HOMEPAGE CALL-TO-ACTION (CTA)
-- ====================================================================

-- name: GetActiveCTA :one
-- Purpose: Retrieves the currently active CTA block for public homepage
-- Parameters: none
-- Return type: single homepage_cta row
-- WHERE: is_active = 1 (only one CTA should be active at a time)
-- LIMIT 1: safety fallback if multiple CTAs accidentally marked active
SELECT * FROM homepage_cta
WHERE is_active = 1
LIMIT 1;

-- name: ListAllCTAs :many
-- Purpose: Lists all CTA variants (active + inactive) for admin management
-- ORDER BY id: chronological order (oldest first)
SELECT * FROM homepage_cta ORDER BY id ASC;

-- name: GetCTA :one
-- Purpose: Retrieves specific CTA by ID for editing
SELECT * FROM homepage_cta WHERE id = ?;

-- name: CreateCTA :one
-- Purpose: Creates a new CTA block variant
-- Parameters (8 positional):
--   1. headline (TEXT): CTA section headline
--   2. description (TEXT): supporting description/pitch
--   3. primary_cta_text (TEXT): primary button text
--   4. primary_cta_url (TEXT): primary button link
--   5. secondary_cta_text (TEXT): optional secondary button text
--   6. secondary_cta_url (TEXT): secondary button link
--   7. background_style (TEXT): visual style ('gradient', 'solid', 'image', etc.)
--   8. is_active (BOOLEAN): whether this CTA is currently displayed
-- Note: Only one CTA should have is_active = 1 at a time
INSERT INTO homepage_cta (
    headline, description,
    primary_cta_text, primary_cta_url,
    secondary_cta_text, secondary_cta_url,
    background_style, is_active
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateCTA :exec
-- Purpose: Updates an existing CTA block
-- Parameters (9 positional): same as CreateCTA + id (WHERE clause)
UPDATE homepage_cta
SET headline = ?, description = ?,
    primary_cta_text = ?, primary_cta_url = ?,
    secondary_cta_text = ?, secondary_cta_url = ?,
    background_style = ?, is_active = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteCTA :exec
-- Purpose: Removes a CTA block variant
DELETE FROM homepage_cta WHERE id = ?;
