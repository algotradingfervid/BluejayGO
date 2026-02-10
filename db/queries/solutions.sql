-- ====================================================================
-- SOLUTIONS
-- ====================================================================

-- name: ListPublishedSolutions :many
SELECT id, title, slug, icon, short_description, display_order, created_at, updated_at
FROM solutions
WHERE is_published = 1
ORDER BY display_order ASC, title ASC;

-- name: GetSolutionBySlug :one
SELECT * FROM solutions
WHERE slug = ? AND is_published = 1
LIMIT 1;

-- name: GetSolutionBySlugIncludeDrafts :one
SELECT * FROM solutions
WHERE slug = ?
LIMIT 1;

-- name: GetSolutionByID :one
SELECT * FROM solutions WHERE id = ?;

-- name: CreateSolution :one
INSERT INTO solutions (
    title, slug, icon, short_description, hero_image_url, hero_title,
    hero_description, overview_content, meta_description, reference_code,
    is_published, display_order
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolution :exec
UPDATE solutions
SET title = ?, slug = ?, icon = ?, short_description = ?,
    hero_image_url = ?, hero_title = ?, hero_description = ?,
    overview_content = ?, meta_description = ?, reference_code = ?,
    is_published = ?, display_order = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteSolution :exec
DELETE FROM solutions WHERE id = ?;

-- name: ListAllSolutions :many
SELECT * FROM solutions ORDER BY display_order ASC, title ASC;

-- name: ListSolutionsAdminFiltered :many
SELECT * FROM solutions
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE
        (CASE WHEN @filter_status = 'published' THEN is_published = 1
              WHEN @filter_status = 'draft' THEN (is_published = 0 OR is_published IS NULL)
              ELSE 1 END)
    END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (title LIKE '%' || @filter_search || '%') END)
ORDER BY display_order ASC, title ASC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountSolutionsAdminFiltered :one
SELECT COUNT(*) FROM solutions
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE
        (CASE WHEN @filter_status = 'published' THEN is_published = 1
              WHEN @filter_status = 'draft' THEN (is_published = 0 OR is_published IS NULL)
              ELSE 1 END)
    END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (title LIKE '%' || @filter_search || '%') END);

-- ====================================================================
-- SOLUTION STATS
-- ====================================================================

-- name: GetSolutionStats :many
SELECT * FROM solution_stats
WHERE solution_id = ?
ORDER BY display_order ASC;

-- name: CreateSolutionStat :one
INSERT INTO solution_stats (solution_id, value, label, display_order)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionStat :exec
UPDATE solution_stats
SET value = ?, label = ?, display_order = ?
WHERE id = ?;

-- name: DeleteSolutionStat :exec
DELETE FROM solution_stats WHERE id = ?;

-- name: DeleteSolutionStatsBySolutionID :exec
DELETE FROM solution_stats WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION CHALLENGES
-- ====================================================================

-- name: GetSolutionChallenges :many
SELECT * FROM solution_challenges
WHERE solution_id = ?
ORDER BY display_order ASC;

-- name: CreateSolutionChallenge :one
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionChallenge :exec
UPDATE solution_challenges
SET title = ?, description = ?, icon = ?, display_order = ?
WHERE id = ?;

-- name: DeleteSolutionChallenge :exec
DELETE FROM solution_challenges WHERE id = ?;

-- name: DeleteSolutionChallengesBySolutionID :exec
DELETE FROM solution_challenges WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION PRODUCTS
-- ====================================================================

-- name: GetSolutionProducts :many
SELECT sp.id, sp.solution_id, sp.product_id, sp.display_order, sp.is_featured,
       p.name AS product_name, p.slug AS product_slug, p.tagline AS product_tagline,
       p.primary_image AS product_image, p.status AS product_status
FROM solution_products sp
JOIN products p ON sp.product_id = p.id
WHERE sp.solution_id = ? AND p.status = 'published'
ORDER BY sp.display_order ASC;

-- name: AddProductToSolution :exec
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured)
VALUES (?, ?, ?, ?)
ON CONFLICT(solution_id, product_id) DO UPDATE SET display_order = excluded.display_order, is_featured = excluded.is_featured;

-- name: RemoveProductFromSolution :exec
DELETE FROM solution_products
WHERE solution_id = ? AND product_id = ?;

-- name: DeleteSolutionProductsBySolutionID :exec
DELETE FROM solution_products WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION CTAS
-- ====================================================================

-- name: GetSolutionCTAs :many
SELECT * FROM solution_ctas
WHERE solution_id = ?;

-- name: CreateSolutionCTA :one
INSERT INTO solution_ctas (
    solution_id, heading, subheading, primary_button_text,
    primary_button_url, secondary_button_text, secondary_button_url,
    phone_number, section_name
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionCTA :exec
UPDATE solution_ctas
SET heading = ?, subheading = ?, primary_button_text = ?,
    primary_button_url = ?, secondary_button_text = ?,
    secondary_button_url = ?, phone_number = ?, section_name = ?
WHERE id = ?;

-- name: DeleteSolutionCTA :exec
DELETE FROM solution_ctas WHERE id = ?;

-- name: DeleteSolutionCTAsBySolutionID :exec
DELETE FROM solution_ctas WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION PAGE FEATURES (Why Choose BlueJay)
-- ====================================================================

-- name: ListSolutionPageFeatures :many
SELECT * FROM solution_page_features
WHERE is_active = 1
ORDER BY display_order ASC;

-- name: ListAllSolutionPageFeatures :many
SELECT * FROM solution_page_features
ORDER BY display_order ASC;

-- name: CreateSolutionPageFeature :one
INSERT INTO solution_page_features (title, description, icon, display_order, is_active)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionPageFeature :exec
UPDATE solution_page_features
SET title = ?, description = ?, icon = ?, display_order = ?, is_active = ?
WHERE id = ?;

-- name: DeleteSolutionPageFeature :exec
DELETE FROM solution_page_features WHERE id = ?;

-- ====================================================================
-- SOLUTIONS LISTING CTA
-- ====================================================================

-- name: GetActiveSolutionsListingCTA :one
SELECT * FROM solutions_listing_cta
WHERE is_active = 1
LIMIT 1;

-- name: CreateSolutionsListingCTA :one
INSERT INTO solutions_listing_cta (
    heading, subheading, primary_button_text, primary_button_url,
    secondary_button_text, secondary_button_url, is_active
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionsListingCTA :exec
UPDATE solutions_listing_cta
SET heading = ?, subheading = ?, primary_button_text = ?,
    primary_button_url = ?, secondary_button_text = ?,
    secondary_button_url = ?, is_active = ?
WHERE id = ?;
