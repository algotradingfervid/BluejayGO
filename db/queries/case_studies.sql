-- ====================================================================
-- CASE STUDIES QUERIES
-- ====================================================================
-- This file manages case study content showcasing client success stories.
-- Case studies demonstrate real-world applications of products/services
-- and include challenge, solution, and outcome sections.
--
-- Managed entities:
-- - case_studies: main case study content with sections
-- - case_study_products: many-to-many link to featured products
-- - case_study_metrics: key performance indicators/results
--
-- Key concepts:
-- - is_published: controls public visibility (1 = published, 0 = draft)
-- - display_order: custom sorting for featured case studies
-- - industry_id: categorizes by client industry vertical
-- ====================================================================

-- ====================================================================
-- PUBLIC CASE STUDY QUERIES
-- ====================================================================

-- name: ListCaseStudies :many
-- sqlc annotation: :many returns slice of published case studies
-- Purpose: Lists all published case studies for public case studies index page
-- Parameters: none
-- Return type: slice of case studies with industry details
-- JOIN:
--   - INNER JOIN industries: gets industry name/slug for categorization
-- WHERE:
--   - is_published = 1: only show published case studies, hide drafts
-- ORDER BY:
--   - Primary: display_order ASC (manual featured ordering)
--   - Secondary: created_at DESC (newest first as fallback)
SELECT
    cs.id, cs.slug, cs.title, cs.client_name, cs.summary,
    cs.hero_image_url, cs.display_order,
    i.id as industry_id, i.name as industry_name, i.slug as industry_slug
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE cs.is_published = 1
ORDER BY cs.display_order ASC, cs.created_at DESC;

-- name: ListCaseStudiesByIndustry :many
-- sqlc annotation: :many returns case studies filtered by industry
-- Purpose: Lists published case studies for specific industry vertical
-- Parameters:
--   1. industry_id (INTEGER): industry to filter by
-- Return type: slice of case studies for that industry
-- WHERE:
--   - is_published = 1: published only
--   - industry_id = ?: filter to specific industry
SELECT
    cs.id, cs.slug, cs.title, cs.client_name, cs.summary,
    cs.hero_image_url, cs.display_order,
    i.id as industry_id, i.name as industry_name, i.slug as industry_slug
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE cs.is_published = 1 AND cs.industry_id = ?
ORDER BY cs.display_order ASC, cs.created_at DESC;

-- name: GetCaseStudyBySlug :one
-- sqlc annotation: :one returns single case study by slug
-- Purpose: Retrieves full published case study for public detail page
-- Parameters:
--   1. slug (TEXT): URL-safe case study identifier
-- Return type: complete case study with all content sections
-- SELECT fields:
--   - All case study sections (challenge, solution, outcome)
--   - SEO metadata (meta_title, meta_description, og_image)
--   - Industry details for breadcrumbs/categorization
-- WHERE:
--   - slug = ?: exact slug match
--   - is_published = 1: published only (hide drafts from public)
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
-- sqlc annotation: :one returns case study including drafts
-- Purpose: Retrieves case study for admin preview (allows viewing unpublished)
-- Parameters:
--   1. slug (TEXT): case study slug
-- Return type: same as GetCaseStudyBySlug but without publish filter
-- Note: Used for admin preview functionality; no is_published filter
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
-- sqlc annotation: :many returns products featured in a case study
-- Purpose: Retrieves products associated with/featured in a case study
-- Parameters:
--   1. case_study_id (INTEGER): case study to get products for
-- Return type: slice of products with category info
-- JOINs:
--   - case_study_products: junction table with display_order
--   - products: product details
--   - product_categories: for category name
-- WHERE:
--   - csp.case_study_id = ?: filter to specific case study
--   - p.status = 'published': only show published products
-- ORDER BY csp.display_order: custom sort order per case study
SELECT
    p.id, p.slug, p.name, p.tagline,
    p.primary_image, pc.name as category_name
FROM case_study_products csp
INNER JOIN products p ON csp.product_id = p.id
INNER JOIN product_categories pc ON p.category_id = pc.id
WHERE csp.case_study_id = ? AND p.status = 'published'
ORDER BY csp.display_order ASC;

-- name: GetCaseStudyMetrics :many
-- sqlc annotation: :many returns key metrics/results for a case study
-- Purpose: Retrieves performance metrics showcasing case study results
-- Parameters:
--   1. case_study_id (INTEGER): case study to get metrics for
-- Return type: slice of metrics (value + label pairs)
-- Example metrics: "50% reduction in costs", "2x faster deployment"
-- ORDER BY display_order: custom presentation order
SELECT id, metric_value, metric_label, display_order
FROM case_study_metrics
WHERE case_study_id = ?
ORDER BY display_order ASC;

-- name: CountCaseStudies :one
-- sqlc annotation: :one returns total published case studies count
-- Purpose: Counts published case studies for stats/pagination
-- Parameters: none
-- Return type: integer count
SELECT COUNT(*) FROM case_studies WHERE is_published = 1;

-- name: CountCaseStudiesByIndustry :one
-- sqlc annotation: :one returns count for specific industry
-- Purpose: Counts published case studies in an industry
-- Parameters:
--   1. industry_id (INTEGER): industry to count
-- Return type: integer count
SELECT COUNT(*) FROM case_studies WHERE is_published = 1 AND industry_id = ?;

-- ====================================================================
-- ADMIN CASE STUDY QUERIES
-- ====================================================================

-- name: AdminListCaseStudies :many
-- sqlc annotation: :many returns all case studies for admin
-- Purpose: Simple list of all case studies (drafts + published) for admin table
-- Parameters: none
-- Return type: slice of case studies with minimal fields
-- Note: No filtering or pagination; includes all statuses
SELECT
    cs.id, cs.slug, cs.title, cs.client_name,
    cs.is_published, cs.display_order, cs.created_at,
    i.name as industry_name
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
ORDER BY cs.display_order ASC, cs.created_at DESC;

-- name: AdminListCaseStudiesFiltered :many
-- sqlc annotation: :many returns filtered/paginated case studies for admin
-- Purpose: Advanced admin list with search and status filtering
-- Parameters (named using @ prefix):
--   @filter_search (TEXT): search term for title/client_name ('' = no filter)
--   @filter_status (TEXT): 'published', 'draft', or '' for all
--   @page_limit (INTEGER): case studies per page
--   @page_offset (INTEGER): pagination offset
-- Return type: slice of case studies with display metadata
-- Complex WHERE using CASE statements:
--   - filter_search: LIKE match in title OR client_name fields
--   - filter_status: maps text values to is_published boolean (1/0)
--     'published' -> is_published = 1
--     'draft' -> is_published = 0
--     '' -> all statuses (condition always true)
SELECT
    cs.id, cs.slug, cs.title, cs.client_name,
    cs.hero_image_url, cs.is_published, cs.display_order,
    cs.updated_at,
    i.name as industry_name
FROM case_studies cs
INNER JOIN industries i ON cs.industry_id = i.id
WHERE
    -- Optional search filter: '' = all, otherwise LIKE match in title or client_name
    (CASE WHEN @filter_search = '' THEN 1 ELSE (cs.title LIKE '%' || @filter_search || '%' OR cs.client_name LIKE '%' || @filter_search || '%') END)
    -- Optional status filter: 'published' = 1, 'draft' = 0, '' = all
    AND (CASE WHEN @filter_status = '' THEN 1
         WHEN @filter_status = 'published' THEN cs.is_published = 1
         WHEN @filter_status = 'draft' THEN cs.is_published = 0
         ELSE 1 END)
ORDER BY cs.display_order ASC, cs.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountCaseStudiesAdminFiltered :one
-- sqlc annotation: :one returns count for filtered case studies
-- Purpose: Counts case studies matching admin filters (for pagination)
-- Parameters: same filters as AdminListCaseStudiesFiltered (except pagination)
-- Return type: integer count
-- Note: WHERE clause MUST match AdminListCaseStudiesFiltered exactly
SELECT COUNT(*) FROM case_studies cs
WHERE
    -- Exact same filter logic as AdminListCaseStudiesFiltered
    (CASE WHEN @filter_search = '' THEN 1 ELSE (cs.title LIKE '%' || @filter_search || '%' OR cs.client_name LIKE '%' || @filter_search || '%') END)
    AND (CASE WHEN @filter_status = '' THEN 1
         WHEN @filter_status = 'published' THEN cs.is_published = 1
         WHEN @filter_status = 'draft' THEN cs.is_published = 0
         ELSE 1 END);

-- name: AdminGetCaseStudy :one
-- sqlc annotation: :one returns single case study by ID for admin editing
-- Purpose: Retrieves complete case study for admin edit form
-- Parameters:
--   1. id (INTEGER): case study primary key
-- Return type: complete case_studies row with all fields
SELECT * FROM case_studies WHERE id = ?;

-- name: AdminCreateCaseStudy :one
-- sqlc annotation: :one returns the created case study
-- Purpose: Creates a new case study (draft or published)
-- Parameters (17 positional):
--   1. slug (TEXT): URL-safe identifier (must be unique)
--   2. title (TEXT): case study headline
--   3. client_name (TEXT): client/company name
--   4. industry_id (INTEGER): foreign key to industries table
--   5. hero_image_url (TEXT): hero/header image
--   6. summary (TEXT): brief overview for listing pages
--   7. challenge_title (TEXT): "Challenge" section heading
--   8. challenge_content (TEXT): challenge description (HTML)
--   9. challenge_bullets (TEXT): bullet points for challenge (JSON/text)
--   10. solution_title (TEXT): "Solution" section heading
--   11. solution_content (TEXT): solution description (HTML)
--   12. outcome_title (TEXT): "Outcome" section heading
--   13. outcome_content (TEXT): results/outcome description (HTML)
--   14. meta_title (TEXT): SEO page title
--   15. meta_description (TEXT): SEO meta description
--   16. is_published (BOOLEAN): 1 = published, 0 = draft
--   17. display_order (INTEGER): featured sort position
-- Return type: complete inserted case study with ID and timestamps
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
-- sqlc annotation: :one returns the updated case study
-- Purpose: Updates an existing case study (all fields except ID/created_at)
-- Parameters (18 positional):
--   1-17. updated field values (same order as AdminCreateCaseStudy)
--   18. id (INTEGER): which case study to update (WHERE clause)
-- Return type: updated case_studies row
-- Note: updated_at explicitly set to CURRENT_TIMESTAMP
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
-- sqlc annotation: :exec returns no data
-- Purpose: Permanently removes a case study
-- Parameters:
--   1. id (INTEGER): case study to delete
-- Return type: none
-- Note: May cascade delete related case_study_products and case_study_metrics
DELETE FROM case_studies WHERE id = ?;

-- ====================================================================
-- CASE STUDY PRODUCTS MANAGEMENT
-- ====================================================================

-- name: AdminAddCaseStudyProduct :one
-- sqlc annotation: :one returns the association row (new or updated)
-- Purpose: Associates a product with a case study (many-to-many relationship)
-- Parameters:
--   1. case_study_id (INTEGER)
--   2. product_id (INTEGER)
--   3. display_order (INTEGER): sort position for this product in this case study
-- Return type: case_study_products row
-- ON CONFLICT: If association already exists, update the display_order
-- Note: Assumes UNIQUE constraint on (case_study_id, product_id)
INSERT INTO case_study_products (case_study_id, product_id, display_order)
VALUES (?, ?, ?)
ON CONFLICT (case_study_id, product_id) DO UPDATE
SET display_order = EXCLUDED.display_order
RETURNING *;

-- name: AdminRemoveCaseStudyProduct :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Removes a specific product association from a case study
-- Parameters:
--   1. case_study_id (INTEGER)
--   2. product_id (INTEGER)
-- Return type: none
DELETE FROM case_study_products
WHERE case_study_id = ? AND product_id = ?;

-- name: AdminListCaseStudyProducts :many
-- sqlc annotation: :many returns products linked to a case study (admin view)
-- Purpose: Lists products associated with a case study for admin management
-- Parameters:
--   1. case_study_id (INTEGER): case study to list products for
-- Return type: slice of products with association metadata
-- JOIN: Gets product name/slug for display in admin UI
-- ORDER BY display_order: shows custom sort order for this case study
SELECT
    csp.id, csp.product_id, csp.display_order,
    p.name as product_name, p.slug as product_slug
FROM case_study_products csp
INNER JOIN products p ON csp.product_id = p.id
WHERE csp.case_study_id = ?
ORDER BY csp.display_order ASC;

-- ====================================================================
-- CASE STUDY METRICS MANAGEMENT
-- ====================================================================

-- name: AdminCreateMetric :one
-- sqlc annotation: :one returns the created metric
-- Purpose: Creates a new performance metric for a case study
-- Parameters (4 positional):
--   1. case_study_id (INTEGER): foreign key to case_studies
--   2. metric_value (TEXT): numeric value with unit (e.g., "50%", "2x")
--   3. metric_label (TEXT): description (e.g., "reduction in costs")
--   4. display_order (INTEGER): sort position in metrics list
-- Return type: complete inserted case_study_metrics row
-- Example: metric_value="75%", metric_label="faster deployment time"
INSERT INTO case_study_metrics (case_study_id, metric_value, metric_label, display_order)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: AdminUpdateMetric :one
-- sqlc annotation: :one returns the updated metric
-- Purpose: Updates an existing case study metric
-- Parameters (4 positional):
--   1-3. updated field values (metric_value, metric_label, display_order)
--   4. id (INTEGER): which metric to update (WHERE clause)
-- Return type: updated case_study_metrics row
UPDATE case_study_metrics SET
    metric_value = ?,
    metric_label = ?,
    display_order = ?
WHERE id = ?
RETURNING *;

-- name: AdminDeleteMetric :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Removes a metric from a case study
-- Parameters:
--   1. id (INTEGER): metric to delete
-- Return type: none
DELETE FROM case_study_metrics WHERE id = ?;

-- name: AdminListMetrics :many
-- sqlc annotation: :many returns all metrics for a case study
-- Purpose: Lists metrics for admin editing interface
-- Parameters:
--   1. case_study_id (INTEGER): case study to list metrics for
-- Return type: slice of case_study_metrics rows
-- ORDER BY display_order: shows custom presentation sequence
SELECT id, case_study_id, metric_value, metric_label, display_order
FROM case_study_metrics
WHERE case_study_id = ?
ORDER BY display_order ASC;
