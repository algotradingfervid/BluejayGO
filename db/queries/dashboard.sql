-- ====================================================================
-- DASHBOARD STATISTICS QUERIES
-- ====================================================================
-- This file contains simple COUNT queries used to populate admin dashboard
-- widgets showing key metrics and pending items requiring attention.
--
-- All queries return single integer counts for dashboard cards/badges.
-- These queries help admins quickly see:
-- - Pending contact submissions needing response
-- - Total partner count
-- - Draft content needing review/publication
-- ====================================================================

-- name: CountNewContactSubmissions :one
-- sqlc annotation: :one returns integer count
-- Purpose: Counts unread contact form submissions (status = 'new')
-- Parameters: none
-- Return type: integer count
-- Used for: Dashboard alert badge showing pending inquiries
SELECT COUNT(*) FROM contact_submissions WHERE status = 'new';

-- name: CountPartners :one
-- sqlc annotation: :one returns integer count
-- Purpose: Counts total partners in the system
-- Parameters: none
-- Return type: integer count
-- Used for: Dashboard "Total Partners" statistic card
SELECT COUNT(*) FROM partners;

-- name: CountDraftProducts :one
-- sqlc annotation: :one returns integer count
-- Purpose: Counts unpublished products (status = 'draft')
-- Parameters: none
-- Return type: integer count
-- Used for: Dashboard alert showing draft products needing review
SELECT COUNT(*) FROM products WHERE status = 'draft';

-- name: CountDraftBlogPosts :one
-- sqlc annotation: :one returns integer count
-- Purpose: Counts unpublished blog posts (status = 'draft')
-- Parameters: none
-- Return type: integer count
-- Used for: Dashboard alert showing draft blog posts needing review
SELECT COUNT(*) FROM blog_posts WHERE status = 'draft';
