-- ====================================================================
-- SOLUTIONS QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing solutions and related entities.
--
-- Main entities:
--   - solutions: Industry solutions/use cases
--   - solution_stats: Quantifiable metrics (e.g., "50% faster", "99.9% uptime")
--   - solution_challenges: Problem statements solved by solution
--   - solution_products: Products associated with solution (many-to-many)
--   - solution_ctas: Call-to-action sections within solution page
--   - solution_page_features: "Why Choose BlueJay" features (shared across solutions)
--   - solutions_listing_cta: CTA for solutions listing page
--
-- Solution status: is_published (1=published, 0=draft)
-- Features:
--   - Rich content pages with heroes, stats, challenges, related products
--   - Draft/published workflow
--   - Admin filtering and search
--   - Featured products within solutions
-- ====================================================================

-- ====================================================================
-- SOLUTIONS - CORE CRUD OPERATIONS
-- ====================================================================

-- name: ListPublishedSolutions :many
-- Retrieves all published solutions in display order for public listing.
--
-- Parameters: none
-- Returns: []Solution - Array of published solutions (partial columns for performance)
--
-- Filtering: is_published = 1 - Only public solutions
--
-- Sorting logic:
--   1. display_order ASC - Custom admin-defined order
--   2. title ASC - Alphabetical fallback for same display_order
--
-- Use case: Solutions listing page, navigation menus
-- Note: Only selects necessary columns for listing (not full content fields)
SELECT id, title, slug, icon, short_description, display_order, created_at, updated_at
FROM solutions
WHERE is_published = 1
ORDER BY display_order ASC, title ASC;

-- name: GetSolutionBySlug :one
-- Retrieves a single published solution by its URL slug.
--
-- Parameters:
--   $1 (TEXT) - slug: URL-safe identifier (e.g., "industrial-automation")
-- Returns: Solution - Single published solution or error if not found/draft
--
-- Filtering:
--   - slug = ? - Matches specific solution
--   - is_published = 1 - Only published solutions (public view)
--
-- Use case: Public solution detail page
SELECT * FROM solutions
WHERE slug = ? AND is_published = 1
LIMIT 1;

-- name: GetSolutionBySlugIncludeDrafts :one
-- Retrieves a single solution by slug regardless of published status.
--
-- Parameters:
--   $1 (TEXT) - slug: URL-safe identifier
-- Returns: Solution - Single solution (published or draft) or error if not found
--
-- Use case: Admin preview mode, editing draft solutions
-- Note: Does NOT filter by is_published, returns draft solutions
SELECT * FROM solutions
WHERE slug = ?
LIMIT 1;

-- name: GetSolutionByID :one
-- Retrieves a single solution by its primary key ID (any status).
--
-- Parameters:
--   $1 (INTEGER) - solution ID
-- Returns: Solution - Single solution or error if not found
--
-- Use case: Admin editing, fetching solution for update
-- Note: No status filtering, returns drafts and published
SELECT * FROM solutions WHERE id = ?;

-- name: CreateSolution :one
-- Creates a new solution record.
--
-- Parameters:
--   $1 (TEXT) - title: Solution display name
--   $2 (TEXT) - slug: URL-safe identifier
--   $3 (TEXT) - icon: Icon identifier or CSS class
--   $4 (TEXT) - short_description: Brief description for cards/listings
--   $5 (TEXT) - hero_image_url: Hero section background image
--   $6 (TEXT) - hero_title: Hero section heading
--   $7 (TEXT) - hero_description: Hero section subheading/description
--   $8 (TEXT) - overview_content: Full solution overview (HTML/Markdown)
--   $9 (TEXT) - meta_description: SEO meta description
--   $10 (TEXT) - reference_code: Internal reference/code (optional)
--   $11 (BOOLEAN) - is_published: Publication status (1=published, 0=draft)
--   $12 (INTEGER) - display_order: Position in solutions listing
--
-- Returns: Solution - The newly created solution with auto-generated ID and timestamps
--
-- Note: Related entities (stats, challenges, products, CTAs) are created separately
INSERT INTO solutions (
    title, slug, icon, short_description, hero_image_url, hero_title,
    hero_description, overview_content, meta_description, reference_code,
    is_published, display_order
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolution :exec
-- Updates an existing solution's core fields.
--
-- Parameters: Same as CreateSolution ($1-$12), plus:
--   $13 (INTEGER) - id: Solution ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
-- Related entities updated via separate queries
UPDATE solutions
SET title = ?, slug = ?, icon = ?, short_description = ?,
    hero_image_url = ?, hero_title = ?, hero_description = ?,
    overview_content = ?, meta_description = ?, reference_code = ?,
    is_published = ?, display_order = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteSolution :exec
-- Permanently deletes a solution.
--
-- Parameters:
--   $1 (INTEGER) - solution ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Should cascade delete related records (stats, challenges, products, CTAs)
-- Note: Ensure foreign key constraints are configured for CASCADE DELETE
DELETE FROM solutions WHERE id = ?;

-- name: ListAllSolutions :many
-- Retrieves all solutions (published and draft) for admin dashboard.
--
-- Parameters: none
-- Returns: []Solution - Array of all solutions
--
-- Sorting:
--   1. display_order ASC - Custom order
--   2. title ASC - Alphabetical fallback
--
-- Use case: Admin solutions management listing
-- Note: Returns ALL solutions regardless of is_published status
SELECT * FROM solutions ORDER BY display_order ASC, title ASC;

-- ====================================================================
-- SOLUTIONS - ADMIN QUERIES
-- ====================================================================

-- name: ListSolutionsAdminFiltered :many
-- Retrieves paginated solutions with optional status and search filters.
--
-- Parameters (named parameters with @):
--   @filter_status (TEXT) - Filter by status ("published", "draft", or "" for all)
--   @filter_search (TEXT) - Search term for title (empty string for no search)
--   @page_limit (INTEGER) - Results per page
--   @page_offset (INTEGER) - Pagination offset
-- Returns: []Solution - Array of filtered solutions
--
-- Complex WHERE clause with nested CASE statements:
--   1. Status filter:
--      - "" (empty) -> Shows all solutions
--      - "published" -> Shows is_published = 1
--      - "draft" -> Shows is_published = 0 OR NULL
--   2. Search filter:
--      - "" (empty) -> No filtering
--      - Non-empty -> LIKE search on title
--
-- Sorting: display_order ASC, title ASC - Custom order then alphabetical
-- Use case: Admin solutions management with status filter and search bar
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
-- Returns count of solutions matching admin filters (for pagination).
--
-- Parameters: Same as ListSolutionsAdminFiltered (@filter_status, @filter_search)
-- Returns: INTEGER - Count of solutions matching filter criteria
--
-- Note: Uses identical WHERE clause as ListSolutionsAdminFiltered for consistent counts
SELECT COUNT(*) FROM solutions
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE
        (CASE WHEN @filter_status = 'published' THEN is_published = 1
              WHEN @filter_status = 'draft' THEN (is_published = 0 OR is_published IS NULL)
              ELSE 1 END)
    END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (title LIKE '%' || @filter_search || '%') END);

-- ====================================================================
-- SOLUTION STATS (Quantifiable Metrics)
-- ====================================================================
-- Stats display key metrics like "50% Faster", "99.9% Uptime", "24/7 Support"

-- name: GetSolutionStats :many
-- Retrieves all stats for a solution in display order.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution to fetch stats for
-- Returns: []SolutionStat - Array of stats
--
-- Sorting: display_order ASC - Stats appear in admin-configured order
-- Use case: Displaying stats section on solution detail page
SELECT * FROM solution_stats
WHERE solution_id = ?
ORDER BY display_order ASC;

-- name: CreateSolutionStat :one
-- Adds a stat metric to a solution.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Parent solution
--   $2 (TEXT) - value: Metric value (e.g., "50%", "99.9%", "24/7")
--   $3 (TEXT) - label: Metric description (e.g., "Faster Processing", "Uptime")
--   $4 (INTEGER) - display_order: Position in stats list
-- Returns: SolutionStat - Newly created stat
INSERT INTO solution_stats (solution_id, value, label, display_order)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionStat :exec
-- Updates an existing solution stat.
--
-- Parameters:
--   $1 (TEXT) - value: Updated metric value
--   $2 (TEXT) - label: Updated label
--   $3 (INTEGER) - display_order: Updated position
--   $4 (INTEGER) - id: Stat ID to update
-- Returns: (none)
UPDATE solution_stats
SET value = ?, label = ?, display_order = ?
WHERE id = ?;

-- name: DeleteSolutionStat :exec
-- Deletes a single solution stat.
--
-- Parameters:
--   $1 (INTEGER) - stat ID to delete
-- Returns: (none)
DELETE FROM solution_stats WHERE id = ?;

-- name: DeleteSolutionStatsBySolutionID :exec
-- Deletes all stats for a solution (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution whose stats to delete
-- Returns: (none)
--
-- Use case: Clearing stats before rebuilding or deleting solution
DELETE FROM solution_stats WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION CHALLENGES (Problem Statements)
-- ====================================================================
-- Challenges describe customer pain points that the solution addresses

-- name: GetSolutionChallenges :many
-- Retrieves all challenges for a solution in display order.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution to fetch challenges for
-- Returns: []SolutionChallenge - Array of challenge problem statements
--
-- Sorting: display_order ASC - Challenges appear in configured order
-- Use case: Displaying "Challenges We Solve" section on solution page
SELECT * FROM solution_challenges
WHERE solution_id = ?
ORDER BY display_order ASC;

-- name: CreateSolutionChallenge :one
-- Adds a challenge/problem statement to a solution.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Parent solution
--   $2 (TEXT) - title: Challenge heading (e.g., "Slow Production Times")
--   $3 (TEXT) - description: Detailed problem description
--   $4 (TEXT) - icon: Icon identifier for visual representation
--   $5 (INTEGER) - display_order: Position in challenges list
-- Returns: SolutionChallenge - Newly created challenge
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionChallenge :exec
-- Updates an existing solution challenge.
--
-- Parameters:
--   $1 (TEXT) - title: Updated heading
--   $2 (TEXT) - description: Updated description
--   $3 (TEXT) - icon: Updated icon identifier
--   $4 (INTEGER) - display_order: Updated position
--   $5 (INTEGER) - id: Challenge ID to update
-- Returns: (none)
UPDATE solution_challenges
SET title = ?, description = ?, icon = ?, display_order = ?
WHERE id = ?;

-- name: DeleteSolutionChallenge :exec
-- Deletes a single solution challenge.
--
-- Parameters:
--   $1 (INTEGER) - challenge ID to delete
-- Returns: (none)
DELETE FROM solution_challenges WHERE id = ?;

-- name: DeleteSolutionChallengesBySolutionID :exec
-- Deletes all challenges for a solution (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution whose challenges to delete
-- Returns: (none)
DELETE FROM solution_challenges WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION PRODUCTS (Many-to-Many Relationship)
-- ====================================================================
-- Links products to solutions showing which products support each solution

-- name: GetSolutionProducts :many
-- Retrieves all published products associated with a solution.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution to fetch products for
-- Returns: []SolutionProduct - Array of products with solution-specific metadata
--
-- JOIN logic:
--   - JOIN products p ON sp.product_id = p.id
--     Brings in product details (name, slug, tagline, image, status)
--
-- Filtering:
--   - sp.solution_id = ? - Products for this solution
--   - p.status = 'published' - Only public products
--
-- Sorting: sp.display_order ASC - Products in admin-configured order
-- Use case: Displaying "Related Products" section on solution page
SELECT sp.id, sp.solution_id, sp.product_id, sp.display_order, sp.is_featured,
       p.name AS product_name, p.slug AS product_slug, p.tagline AS product_tagline,
       p.primary_image AS product_image, p.status AS product_status
FROM solution_products sp
JOIN products p ON sp.product_id = p.id
WHERE sp.solution_id = ? AND p.status = 'published'
ORDER BY sp.display_order ASC;

-- name: AddProductToSolution :exec
-- Adds or updates a product association with a solution (upsert).
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution to link product to
--   $2 (INTEGER) - product_id: Product to link
--   $3 (INTEGER) - display_order: Position in product list
--   $4 (BOOLEAN) - is_featured: Whether product is highlighted
-- Returns: (none)
--
-- Note: ON CONFLICT clause performs upsert - updates if association exists
-- Use case: Adding products to solution or reordering existing associations
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured)
VALUES (?, ?, ?, ?)
ON CONFLICT(solution_id, product_id) DO UPDATE SET display_order = excluded.display_order, is_featured = excluded.is_featured;

-- name: RemoveProductFromSolution :exec
-- Removes a product association from a solution.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution to remove product from
--   $2 (INTEGER) - product_id: Product to unlink
-- Returns: (none)
DELETE FROM solution_products
WHERE solution_id = ? AND product_id = ?;

-- name: DeleteSolutionProductsBySolutionID :exec
-- Removes all product associations for a solution (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution to clear products from
-- Returns: (none)
DELETE FROM solution_products WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION CTAS (Call-to-Action Sections)
-- ====================================================================
-- CTAs are action-oriented sections within solution pages (e.g., "Get Started", "Contact Sales")

-- name: GetSolutionCTAs :many
-- Retrieves all CTA sections for a solution.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution to fetch CTAs for
-- Returns: []SolutionCTA - Array of CTA sections
--
-- Use case: Displaying CTAs within solution page (typically 1-2 per solution)
-- Note: Multiple CTAs can exist per solution (e.g., mid-page and bottom CTAs)
SELECT * FROM solution_ctas
WHERE solution_id = ?;

-- name: CreateSolutionCTA :one
-- Creates a CTA section for a solution.
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Parent solution
--   $2 (TEXT) - heading: CTA main heading
--   $3 (TEXT) - subheading: CTA description/subheading
--   $4 (TEXT) - primary_button_text: Primary button label
--   $5 (TEXT) - primary_button_url: Primary button link
--   $6 (TEXT) - secondary_button_text: Secondary button label (optional)
--   $7 (TEXT) - secondary_button_url: Secondary button link (optional)
--   $8 (TEXT) - phone_number: Contact phone for "Call Now" CTAs (optional)
--   $9 (TEXT) - section_name: Identifier for positioning (e.g., "mid-page", "footer")
-- Returns: SolutionCTA - Newly created CTA
INSERT INTO solution_ctas (
    solution_id, heading, subheading, primary_button_text,
    primary_button_url, secondary_button_text, secondary_button_url,
    phone_number, section_name
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionCTA :exec
-- Updates an existing solution CTA.
--
-- Parameters: $1-$8 same as CreateSolutionCTA (minus solution_id), plus:
--   $9 (INTEGER) - id: CTA ID to update
-- Returns: (none)
UPDATE solution_ctas
SET heading = ?, subheading = ?, primary_button_text = ?,
    primary_button_url = ?, secondary_button_text = ?,
    secondary_button_url = ?, phone_number = ?, section_name = ?
WHERE id = ?;

-- name: DeleteSolutionCTA :exec
-- Deletes a single solution CTA.
--
-- Parameters:
--   $1 (INTEGER) - CTA ID to delete
-- Returns: (none)
DELETE FROM solution_ctas WHERE id = ?;

-- name: DeleteSolutionCTAsBySolutionID :exec
-- Deletes all CTAs for a solution (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - solution_id: Solution whose CTAs to delete
-- Returns: (none)
DELETE FROM solution_ctas WHERE solution_id = ?;

-- ====================================================================
-- SOLUTION PAGE FEATURES ("Why Choose BlueJay" Section)
-- ====================================================================
-- Shared features that appear on all solution pages (company differentiators)

-- name: ListSolutionPageFeatures :many
-- Retrieves all active solution page features.
--
-- Parameters: none
-- Returns: []SolutionPageFeature - Array of active features
--
-- Filtering: is_active = 1 - Only visible features
-- Sorting: display_order ASC - Features in configured order
-- Use case: Displaying "Why Choose BlueJay" section on solution pages
SELECT * FROM solution_page_features
WHERE is_active = 1
ORDER BY display_order ASC;

-- name: ListAllSolutionPageFeatures :many
-- Retrieves all solution page features (active and inactive) for admin.
--
-- Parameters: none
-- Returns: []SolutionPageFeature - Array of all features
--
-- Sorting: display_order ASC
-- Use case: Admin management of solution page features
SELECT * FROM solution_page_features
ORDER BY display_order ASC;

-- name: CreateSolutionPageFeature :one
-- Creates a new solution page feature.
--
-- Parameters:
--   $1 (TEXT) - title: Feature heading
--   $2 (TEXT) - description: Feature description
--   $3 (TEXT) - icon: Icon identifier
--   $4 (INTEGER) - display_order: Position in features list
--   $5 (BOOLEAN) - is_active: Visibility status
-- Returns: SolutionPageFeature - Newly created feature
INSERT INTO solution_page_features (title, description, icon, display_order, is_active)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionPageFeature :exec
-- Updates an existing solution page feature.
--
-- Parameters:
--   $1-$5: Same as CreateSolutionPageFeature
--   $6 (INTEGER) - id: Feature ID to update
-- Returns: (none)
UPDATE solution_page_features
SET title = ?, description = ?, icon = ?, display_order = ?, is_active = ?
WHERE id = ?;

-- name: DeleteSolutionPageFeature :exec
-- Deletes a solution page feature.
--
-- Parameters:
--   $1 (INTEGER) - feature ID to delete
-- Returns: (none)
DELETE FROM solution_page_features WHERE id = ?;

-- ====================================================================
-- SOLUTIONS LISTING CTA (Singleton for Solutions Index Page)
-- ====================================================================
-- CTA displayed on the main solutions listing page

-- name: GetActiveSolutionsListingCTA :one
-- Retrieves the active CTA for solutions listing page (singleton).
--
-- Parameters: none
-- Returns: SolutionsListingCTA - The active CTA or error if none active
--
-- Filtering: is_active = 1 - Only the currently active CTA
-- Note: Typically only one CTA should be active at a time
SELECT * FROM solutions_listing_cta
WHERE is_active = 1
LIMIT 1;

-- name: CreateSolutionsListingCTA :one
-- Creates a new solutions listing CTA.
--
-- Parameters:
--   $1 (TEXT) - heading: CTA main heading
--   $2 (TEXT) - subheading: CTA description
--   $3 (TEXT) - primary_button_text: Primary button label
--   $4 (TEXT) - primary_button_url: Primary button URL
--   $5 (TEXT) - secondary_button_text: Secondary button label (optional)
--   $6 (TEXT) - secondary_button_url: Secondary button URL (optional)
--   $7 (BOOLEAN) - is_active: Whether this CTA is active
-- Returns: SolutionsListingCTA - Newly created CTA
INSERT INTO solutions_listing_cta (
    heading, subheading, primary_button_text, primary_button_url,
    secondary_button_text, secondary_button_url, is_active
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateSolutionsListingCTA :exec
-- Updates an existing solutions listing CTA.
--
-- Parameters: $1-$7 same as CreateSolutionsListingCTA, plus:
--   $8 (INTEGER) - id: CTA ID to update
-- Returns: (none)
UPDATE solutions_listing_cta
SET heading = ?, subheading = ?, primary_button_text = ?,
    primary_button_url = ?, secondary_button_text = ?,
    secondary_button_url = ?, is_active = ?
WHERE id = ?;
