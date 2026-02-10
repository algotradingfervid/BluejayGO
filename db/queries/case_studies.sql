-- ====================================================================
-- CASE STUDIES
-- ====================================================================

-- name: ListCaseStudies :many
SELECT
    cs.id, cs.slug, cs.title, cs.client_name, cs.summary,
    cs.hero_image_url, cs.display_order,
    i.id as industry_id, i.name as industry_name, i.slug as industry_slug
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE cs.is_published = 1
ORDER BY cs.display_order ASC, cs.created_at DESC;

-- name: ListCaseStudiesByIndustry :many
SELECT
    cs.id, cs.slug, cs.title, cs.client_name, cs.summary,
    cs.hero_image_url, cs.display_order,
    i.id as industry_id, i.name as industry_name, i.slug as industry_slug
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE cs.is_published = 1 AND cs.industry_id = ?
ORDER BY cs.display_order ASC, cs.created_at DESC;

-- name: GetCaseStudyBySlug :one
SELECT
    cs.id, cs.slug, cs.title, cs.client_name, cs.summary,
    cs.industry_id, cs.hero_image_url,
    cs.challenge_title, cs.challenge_content, cs.challenge_bullets,
    cs.solution_title, cs.solution_content,
    cs.outcome_title, cs.outcome_content,
    cs.meta_title, cs.meta_description, cs.og_image, cs.created_at,
    i.name as industry_name, i.slug as industry_slug
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE cs.slug = ? AND cs.is_published = 1;

-- name: GetCaseStudyBySlugIncludeDrafts :one
SELECT
    cs.id, cs.slug, cs.title, cs.client_name, cs.summary,
    cs.industry_id, cs.hero_image_url,
    cs.challenge_title, cs.challenge_content, cs.challenge_bullets,
    cs.solution_title, cs.solution_content,
    cs.outcome_title, cs.outcome_content,
    cs.meta_title, cs.meta_description, cs.og_image, cs.created_at,
    i.name as industry_name, i.slug as industry_slug
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE cs.slug = ?;

-- name: GetCaseStudyProducts :many
SELECT
    p.id, p.slug, p.name, p.tagline,
    p.primary_image, pc.name as category_name
FROM case_study_products csp
INNER JOIN products p ON csp.product_id = p.id
INNER JOIN product_categories pc ON p.category_id = pc.id
WHERE csp.case_study_id = ? AND p.status = 'published'
ORDER BY csp.display_order ASC;

-- name: GetCaseStudyMetrics :many
SELECT id, metric_value, metric_label, display_order
FROM case_study_metrics
WHERE case_study_id = ?
ORDER BY display_order ASC;

-- name: CountCaseStudies :one
SELECT COUNT(*) FROM case_studies WHERE is_published = 1;

-- name: CountCaseStudiesByIndustry :one
SELECT COUNT(*) FROM case_studies WHERE is_published = 1 AND industry_id = ?;

-- Admin queries

-- name: AdminListCaseStudies :many
SELECT
    cs.id, cs.slug, cs.title, cs.client_name,
    cs.is_published, cs.display_order, cs.created_at,
    i.name as industry_name
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
ORDER BY cs.display_order ASC, cs.created_at DESC;

-- name: AdminListCaseStudiesFiltered :many
SELECT
    cs.id, cs.slug, cs.title, cs.client_name,
    cs.hero_image_url, cs.is_published, cs.display_order,
    cs.updated_at,
    i.name as industry_name
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE
    (CASE WHEN @filter_search = '' THEN 1 ELSE (cs.title LIKE '%' || @filter_search || '%' OR cs.client_name LIKE '%' || @filter_search || '%') END)
    AND (CASE WHEN @filter_status = '' THEN 1
         WHEN @filter_status = 'published' THEN cs.is_published = 1
         WHEN @filter_status = 'draft' THEN cs.is_published = 0
         ELSE 1 END)
ORDER BY cs.display_order ASC, cs.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountCaseStudiesAdminFiltered :one
SELECT COUNT(*) FROM case_studies cs
WHERE
    (CASE WHEN @filter_search = '' THEN 1 ELSE (cs.title LIKE '%' || @filter_search || '%' OR cs.client_name LIKE '%' || @filter_search || '%') END)
    AND (CASE WHEN @filter_status = '' THEN 1
         WHEN @filter_status = 'published' THEN cs.is_published = 1
         WHEN @filter_status = 'draft' THEN cs.is_published = 0
         ELSE 1 END);

-- name: AdminGetCaseStudy :one
SELECT * FROM case_studies WHERE id = ?;

-- name: AdminCreateCaseStudy :one
INSERT INTO case_studies (
    slug, title, client_name, industry_id, hero_image_url, summary,
    challenge_title, challenge_content, challenge_bullets,
    solution_title, solution_content,
    outcome_title, outcome_content,
    meta_title, meta_description, is_published, display_order
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?,
    ?, ?,
    ?, ?,
    ?, ?, ?, ?
) RETURNING *;

-- name: AdminUpdateCaseStudy :one
UPDATE case_studies SET
    slug = ?,
    title = ?,
    client_name = ?,
    industry_id = ?,
    hero_image_url = ?,
    summary = ?,
    challenge_title = ?,
    challenge_content = ?,
    challenge_bullets = ?,
    solution_title = ?,
    solution_content = ?,
    outcome_title = ?,
    outcome_content = ?,
    meta_title = ?,
    meta_description = ?,
    is_published = ?,
    display_order = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: AdminDeleteCaseStudy :exec
DELETE FROM case_studies WHERE id = ?;

-- Case study products management

-- name: AdminAddCaseStudyProduct :one
INSERT INTO case_study_products (case_study_id, product_id, display_order)
VALUES (?, ?, ?)
ON CONFLICT (case_study_id, product_id) DO UPDATE
SET display_order = EXCLUDED.display_order
RETURNING *;

-- name: AdminRemoveCaseStudyProduct :exec
DELETE FROM case_study_products
WHERE case_study_id = ? AND product_id = ?;

-- name: AdminListCaseStudyProducts :many
SELECT
    csp.id, csp.product_id, csp.display_order,
    p.name as product_name, p.slug as product_slug
FROM case_study_products csp
INNER JOIN products p ON csp.product_id = p.id
WHERE csp.case_study_id = ?
ORDER BY csp.display_order ASC;

-- Metrics management

-- name: AdminCreateMetric :one
INSERT INTO case_study_metrics (case_study_id, metric_value, metric_label, display_order)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: AdminUpdateMetric :one
UPDATE case_study_metrics SET
    metric_value = ?,
    metric_label = ?,
    display_order = ?
WHERE id = ?
RETURNING *;

-- name: AdminDeleteMetric :exec
DELETE FROM case_study_metrics WHERE id = ?;

-- name: AdminListMetrics :many
SELECT id, case_study_id, metric_value, metric_label, display_order
FROM case_study_metrics
WHERE case_study_id = ?
ORDER BY display_order ASC;
