-- ====================================================================
-- WHITEPAPERS QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing whitepapers and download tracking.
--
-- Main entities:
--   - whitepapers: PDF whitepapers with metadata
--   - whitepaper_learning_points: Key takeaways/bullets for each whitepaper
--   - whitepaper_downloads: Download tracking with lead capture (name, email, company)
--
-- Related: whitepaper_topics for categorization
--
-- Features:
--   - Lead generation via download forms
--   - Download analytics and tracking
--   - Gradient cover colors for visual appeal
--   - Topic-based filtering
--   - Published/draft workflow
--   - Admin filtering with date ranges
-- ====================================================================

-- ====================================================================
-- WHITEPAPERS - PUBLIC QUERIES
-- ====================================================================

-- name: ListPublishedWhitepapers :many
-- Retrieves all published whitepapers with topic metadata for public listing.
--
-- Parameters: none
-- Returns: []Whitepaper - Array of published whitepapers with topic information
--
-- JOIN logic:
--   - INNER JOIN whitepaper_topics t ON w.topic_id = t.id
--     Adds topic_name and topic_color_hex for display/filtering
--
-- Filtering: w.is_published = 1 - Only public whitepapers
--
-- Sorting logic:
--   1. w.published_date DESC - Newest whitepapers first
--   2. w.id DESC - ID fallback for same publication date
--
-- Use case: Whitepapers listing page, showing latest publications
-- Note: Selects subset of columns for performance (excludes PDF path, full metadata)
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.file_size_bytes, w.page_count,
    w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.is_published = 1
ORDER BY w.published_date DESC, w.id DESC;

-- name: ListPublishedWhitepapersByTopic :many
-- Retrieves published whitepapers filtered by a specific topic.
--
-- Parameters:
--   $1 (INTEGER) - topic_id: Filter by this whitepaper topic
-- Returns: []Whitepaper - Array of published whitepapers in the topic
--
-- JOIN logic: Same as ListPublishedWhitepapers
--
-- Filtering:
--   - w.is_published = 1 - Only public whitepapers
--   - w.topic_id = ? - Specific topic only
--
-- Sorting: Same as ListPublishedWhitepapers (newest first)
-- Use case: Topic-specific whitepaper listing pages
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.file_size_bytes, w.page_count,
    w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.is_published = 1 AND w.topic_id = ?
ORDER BY w.published_date DESC, w.id DESC;

-- name: GetWhitepaperBySlug :one
-- Retrieves a single published whitepaper by slug with full metadata.
--
-- Parameters:
--   $1 (TEXT) - slug: URL-safe identifier
-- Returns: Whitepaper - Single published whitepaper with complete data
--
-- JOIN logic: Adds topic_name and topic_color_hex
--
-- Filtering:
--   - w.slug = ? - Matches specific whitepaper
--   - w.is_published = 1 - Only published whitepapers
--
-- Use case: Public whitepaper detail/download page
-- Note: Includes pdf_file_path, SEO metadata (meta_title, meta_description, og_image)
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.pdf_file_path, w.file_size_bytes,
    w.page_count, w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    w.meta_description, w.meta_title, w.og_image,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.slug = ? AND w.is_published = 1;

-- name: GetWhitepaperBySlugIncludeDrafts :one
-- Retrieves a single whitepaper by slug regardless of published status.
--
-- Parameters:
--   $1 (TEXT) - slug: URL-safe identifier
-- Returns: Whitepaper - Single whitepaper (published or draft)
--
-- Use case: Admin preview mode, editing draft whitepapers
-- Note: Does NOT filter by is_published, returns draft whitepapers
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.pdf_file_path, w.file_size_bytes,
    w.page_count, w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    w.meta_description, w.meta_title, w.og_image,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.slug = ?;

-- name: GetWhitepaperLearningPoints :many
-- Retrieves all learning points (key takeaways) for a whitepaper.
--
-- Parameters:
--   $1 (INTEGER) - whitepaper_id: Whitepaper to fetch learning points for
-- Returns: []WhitepaperLearningPoint - Array of bullet points
--
-- Sorting:
--   1. display_order ASC - Admin-configured order
--   2. id ASC - ID fallback for same display_order
--
-- Use case: Displaying "What You'll Learn" section on whitepaper page
SELECT id, point_text, display_order
FROM whitepaper_learning_points
WHERE whitepaper_id = ?
ORDER BY display_order ASC, id ASC;

-- name: GetRelatedWhitepapers :many
-- Retrieves up to 3 related published whitepapers from the same topic.
--
-- Parameters:
--   $1 (INTEGER) - whitepaper_id: Current whitepaper ID (to exclude from results)
--   $2 (INTEGER) - topic_id: Topic to find related whitepapers in
-- Returns: []Whitepaper - Array of up to 3 related whitepapers (minimal fields)
--
-- Filtering:
--   - w.is_published = 1 - Only published whitepapers
--   - w.id != ? - Excludes current whitepaper
--   - w.topic_id = ? - Same topic only
--
-- Sorting: w.published_date DESC - Newest related whitepapers first
-- LIMIT: 3 - Maximum 3 recommendations
--
-- Use case: "Related Whitepapers" section on whitepaper detail page
SELECT
    w.id, w.title, w.slug, w.cover_color_from, w.cover_color_to
FROM whitepapers w
WHERE w.is_published = 1
  AND w.id != ?
  AND w.topic_id = ?
ORDER BY w.published_date DESC
LIMIT 3;

-- name: CreateWhitepaperDownload :one
-- Records a whitepaper download with lead capture information.
--
-- Parameters:
--   $1 (INTEGER) - whitepaper_id: Whitepaper being downloaded
--   $2 (TEXT) - name: Downloader's name
--   $3 (TEXT) - email: Downloader's email (lead capture)
--   $4 (TEXT) - company: Downloader's company (optional)
--   $5 (TEXT) - designation: Downloader's job title (optional)
--   $6 (BOOLEAN) - marketing_consent: Whether user opted into marketing
--   $7 (TEXT) - ip_address: Downloader's IP for analytics
--   $8 (TEXT) - user_agent: Browser user agent string
--
-- Returns: Partial WhitepaperDownload - Only id and created_at
--
-- Use case: Lead generation - capturing contact info when user downloads whitepaper
-- Note: This data feeds into CRM/marketing automation systems
INSERT INTO whitepaper_downloads (
    whitepaper_id, name, email, company, designation, marketing_consent, ip_address, user_agent
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, created_at;

-- name: IncrementWhitepaperDownloadCount :exec
-- Increments the download counter for analytics tracking.
--
-- Parameters:
--   $1 (INTEGER) - whitepaper_id: Whitepaper that was downloaded
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Tracking download popularity metrics
-- Note: Uses download_count + 1 for atomic increment without race conditions
UPDATE whitepapers
SET download_count = download_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: CountPublishedWhitepapers :one
-- Returns the total count of published whitepapers.
--
-- Parameters: none
-- Returns: INTEGER - Total number of published whitepapers
--
-- Use case: Site statistics, pagination calculations
SELECT COUNT(*) FROM whitepapers WHERE is_published = 1;

-- name: CountPublishedWhitepapersByTopic :one
-- Returns the count of published whitepapers in a specific topic.
--
-- Parameters:
--   $1 (INTEGER) - topic_id: Topic to count whitepapers in
-- Returns: INTEGER - Number of published whitepapers in the topic
--
-- Use case: Topic page statistics, showing "(12 whitepapers)" on topic badges
SELECT COUNT(*) FROM whitepapers WHERE is_published = 1 AND topic_id = ?;

-- ====================================================================
-- WHITEPAPERS - ADMIN QUERIES
-- ====================================================================

-- name: ListAllWhitepapers :many
-- Retrieves all whitepapers (published and draft) for admin dashboard.
--
-- Parameters: none
-- Returns: []Whitepaper - Array of all whitepapers with topic name
--
-- JOIN logic: Adds topic_name from whitepaper_topics
--
-- Sorting: w.created_at DESC - Newest whitepapers first (by creation date)
-- Use case: Admin whitepapers management listing
-- Note: Returns ALL whitepapers regardless of is_published status
-- Note: Selects subset of columns optimized for admin listing (not full content)
SELECT
    w.id, w.title, w.slug, w.topic_id, w.is_published, w.published_date, w.download_count,
    w.created_at, w.updated_at,
    t.name as topic_name
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
ORDER BY w.created_at DESC;

-- name: ListWhitepapersAdminFiltered :many
-- Retrieves paginated whitepapers with optional filters for admin interface.
--
-- Parameters (named parameters with @):
--   @filter_search (TEXT) - Search term for title (empty string for no search)
--   @filter_topic (INTEGER) - Filter by topic_id (0 for all topics)
--   @filter_status (TEXT) - Filter by status ("published", "draft", or "" for all)
--   @page_limit (INTEGER) - Results per page
--   @page_offset (INTEGER) - Pagination offset
-- Returns: []Whitepaper - Array of filtered whitepapers with topic name
--
-- Complex WHERE clause with CASE statements for optional filters:
--   1. Search filter: CASE WHEN @filter_search = '' THEN 1 (all) ELSE LIKE match on title
--   2. Topic filter: CASE WHEN @filter_topic = 0 THEN 1 (all) ELSE match topic_id
--   3. Status filter:
--      - "" (empty) -> Shows all whitepapers
--      - "published" -> Shows is_published = 1
--      - "draft" -> Shows is_published = 0
--
-- Sorting: w.created_at DESC - Newest whitepapers first
-- Use case: Admin whitepapers listing with search, topic dropdown, and status filter
SELECT
    w.id, w.title, w.slug, w.topic_id, w.is_published, w.published_date, w.download_count,
    w.created_at, w.updated_at,
    t.name as topic_name
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE
    (CASE WHEN @filter_search = '' THEN 1 ELSE w.title LIKE '%' || @filter_search || '%' END)
    AND (CASE WHEN @filter_topic = 0 THEN 1 ELSE w.topic_id = @filter_topic END)
    AND (CASE WHEN @filter_status = '' THEN 1
         WHEN @filter_status = 'published' THEN w.is_published = 1
         WHEN @filter_status = 'draft' THEN w.is_published = 0
         ELSE 1 END)
ORDER BY w.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountWhitepapersAdminFiltered :one
-- Returns count of whitepapers matching admin filters (for pagination).
--
-- Parameters: Same as ListWhitepapersAdminFiltered (@filter_search, @filter_topic, @filter_status)
-- Returns: INTEGER - Count of whitepapers matching filter criteria
--
-- Note: Uses identical WHERE clause as ListWhitepapersAdminFiltered for consistent counts
SELECT COUNT(*) FROM whitepapers w
WHERE
    (CASE WHEN @filter_search = '' THEN 1 ELSE w.title LIKE '%' || @filter_search || '%' END)
    AND (CASE WHEN @filter_topic = 0 THEN 1 ELSE w.topic_id = @filter_topic END)
    AND (CASE WHEN @filter_status = '' THEN 1
         WHEN @filter_status = 'published' THEN w.is_published = 1
         WHEN @filter_status = 'draft' THEN w.is_published = 0
         ELSE 1 END);

-- name: ListWhitepaperDownloadsFiltered :many
-- Retrieves paginated whitepaper download records with filters (lead management).
--
-- Parameters (named parameters with @):
--   @filter_whitepaper (INTEGER) - Filter by whitepaper_id (0 for all whitepapers)
--   @filter_date_from (TEXT) - Start date filter (YYYY-MM-DD format, empty for no filter)
--   @filter_date_to (TEXT) - End date filter (YYYY-MM-DD format, empty for no filter)
--   @page_limit (INTEGER) - Results per page
--   @page_offset (INTEGER) - Pagination offset
-- Returns: []WhitepaperDownload - Array of download records with whitepaper title
--
-- JOIN logic:
--   - INNER JOIN whitepapers w ON wd.whitepaper_id = w.id
--     Adds whitepaper_title for context in admin listing
--
-- Complex WHERE clause with CASE statements:
--   1. Whitepaper filter: CASE WHEN @filter_whitepaper = 0 THEN 1 (all) ELSE match whitepaper_id
--   2. Date range filters:
--      - filter_date_from: >= comparison for start date
--      - filter_date_to: <= comparison with ' 23:59:59' appended to include full day
--
-- Sorting: wd.created_at DESC - Newest downloads first
-- Use case: Admin lead management, download analytics, CRM export
SELECT
    wd.id, wd.whitepaper_id, wd.name, wd.email, wd.company, wd.designation,
    wd.marketing_consent, wd.created_at,
    w.title as whitepaper_title
FROM whitepaper_downloads wd
INNER JOIN whitepapers w ON wd.whitepaper_id = w.id
WHERE
    (CASE WHEN @filter_whitepaper = 0 THEN 1 ELSE wd.whitepaper_id = @filter_whitepaper END)
    AND (CASE WHEN @filter_date_from = '' THEN 1 ELSE wd.created_at >= @filter_date_from END)
    AND (CASE WHEN @filter_date_to = '' THEN 1 ELSE wd.created_at <= @filter_date_to || ' 23:59:59' END)
ORDER BY wd.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountWhitepaperDownloadsFiltered :one
-- Returns count of download records matching admin filters (for pagination).
--
-- Parameters: Same as ListWhitepaperDownloadsFiltered (@filter_whitepaper, @filter_date_from, @filter_date_to)
-- Returns: INTEGER - Count of download records matching filter criteria
--
-- Note: Uses identical WHERE clause as ListWhitepaperDownloadsFiltered for consistent counts
SELECT COUNT(*) FROM whitepaper_downloads wd
WHERE
    (CASE WHEN @filter_whitepaper = 0 THEN 1 ELSE wd.whitepaper_id = @filter_whitepaper END)
    AND (CASE WHEN @filter_date_from = '' THEN 1 ELSE wd.created_at >= @filter_date_from END)
    AND (CASE WHEN @filter_date_to = '' THEN 1 ELSE wd.created_at <= @filter_date_to || ' 23:59:59' END);

-- name: ListWhitepaperTopicsWithCount :many
-- Retrieves all topics with whitepaper count (aggregated via subquery).
--
-- Parameters: none
-- Returns: []WhitepaperTopic - Array of topics with whitepaper_count column
--
-- Subquery logic:
--   - (SELECT COUNT(*) FROM whitepapers w WHERE w.topic_id = t.id) as whitepaper_count
--     Counts total whitepapers (published + draft) for each topic
--
-- Sorting: t.sort_order ASC, t.name ASC
-- Use case: Admin topic management showing whitepaper count per topic
SELECT t.id, t.name, t.slug, t.color_hex, t.icon, t.description, t.sort_order, t.created_at, t.updated_at,
    (SELECT COUNT(*) FROM whitepapers w WHERE w.topic_id = t.id) as whitepaper_count
FROM whitepaper_topics t
ORDER BY t.sort_order ASC, t.name ASC;

-- name: GetWhitepaperByID :one
-- Retrieves a single whitepaper by its primary key ID (any status).
--
-- Parameters:
--   $1 (INTEGER) - whitepaper ID
-- Returns: Whitepaper - Single whitepaper record or error if not found
--
-- Use case: Admin editing, fetching whitepaper for update regardless of status
-- Note: Does NOT filter by is_published, returns draft whitepapers
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.pdf_file_path, w.file_size_bytes,
    w.page_count, w.published_date, w.is_published, w.cover_color_from, w.cover_color_to,
    w.meta_description, w.download_count, w.created_at, w.updated_at
FROM whitepapers w
WHERE w.id = ?;

-- name: CreateWhitepaper :one
-- Creates a new whitepaper record.
--
-- Parameters:
--   $1 (TEXT) - title: Whitepaper title
--   $2 (TEXT) - slug: URL-safe identifier
--   $3 (TEXT) - description: Brief description/summary
--   $4 (INTEGER) - topic_id: Foreign key to whitepaper_topics
--   $5 (TEXT) - pdf_file_path: Path to uploaded PDF file
--   $6 (INTEGER) - file_size_bytes: PDF file size in bytes
--   $7 (INTEGER) - page_count: Number of pages in PDF
--   $8 (TEXT) - published_date: Publication date (YYYY-MM-DD)
--   $9 (BOOLEAN) - is_published: Publication status (1=published, 0=draft)
--   $10 (TEXT) - cover_color_from: Gradient start color (hex)
--   $11 (TEXT) - cover_color_to: Gradient end color (hex)
--   $12 (TEXT) - meta_description: SEO meta description
--
-- Returns: Whitepaper - The newly created whitepaper with auto-generated ID and timestamps
--
-- Note: Learning points created separately via CreateWhitepaperLearningPoint
INSERT INTO whitepapers (
    title, slug, description, topic_id, pdf_file_path, file_size_bytes, page_count,
    published_date, is_published, cover_color_from, cover_color_to, meta_description
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateWhitepaper :exec
-- Updates an existing whitepaper's core fields.
--
-- Parameters: Same as CreateWhitepaper ($1-$12), plus:
--   $13 (INTEGER) - id: Whitepaper ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
UPDATE whitepapers
SET title = ?, slug = ?, description = ?, topic_id = ?, pdf_file_path = ?,
    file_size_bytes = ?, page_count = ?, published_date = ?, is_published = ?,
    cover_color_from = ?, cover_color_to = ?, meta_description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteWhitepaper :exec
-- Permanently deletes a whitepaper.
--
-- Parameters:
--   $1 (INTEGER) - whitepaper ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Should cascade delete related records (learning points, downloads)
-- Note: Physical PDF file should be deleted separately by application code
DELETE FROM whitepapers WHERE id = ?;

-- name: CreateWhitepaperLearningPoint :one
-- Adds a learning point (key takeaway) to a whitepaper.
--
-- Parameters:
--   $1 (INTEGER) - whitepaper_id: Parent whitepaper
--   $2 (TEXT) - point_text: Learning point description
--   $3 (INTEGER) - display_order: Position in learning points list
-- Returns: Partial WhitepaperLearningPoint - Only id
--
-- Use case: Building "What You'll Learn" bullet points during whitepaper creation/editing
INSERT INTO whitepaper_learning_points (whitepaper_id, point_text, display_order)
VALUES (?, ?, ?)
RETURNING id;

-- name: DeleteWhitepaperLearningPoints :exec
-- Deletes all learning points for a whitepaper (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - whitepaper_id: Whitepaper whose learning points to delete
-- Returns: (none)
--
-- Use case: Clearing learning points before rebuilding or deleting whitepaper
DELETE FROM whitepaper_learning_points WHERE whitepaper_id = ?;

-- name: ListWhitepaperDownloads :many
-- Retrieves paginated whitepaper download records (all whitepapers).
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Results per page
--   $2 (INTEGER) - OFFSET: Pagination offset
-- Returns: []WhitepaperDownload - Array of download records with whitepaper title
--
-- JOIN logic: Adds whitepaper_title from whitepapers table
-- Sorting: wd.created_at DESC - Newest downloads first
-- Use case: Admin lead management dashboard (all whitepapers combined)
SELECT
    wd.id, wd.whitepaper_id, wd.name, wd.email, wd.company, wd.designation,
    wd.marketing_consent, wd.created_at,
    w.title as whitepaper_title
FROM whitepaper_downloads wd
INNER JOIN whitepapers w ON wd.whitepaper_id = w.id
ORDER BY wd.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountWhitepaperDownloads :one
-- Returns the total count of all whitepaper downloads.
--
-- Parameters: none
-- Returns: INTEGER - Total number of download records across all whitepapers
--
-- Use case: Site statistics, lead generation metrics
SELECT COUNT(*) FROM whitepaper_downloads;

-- name: ListWhitepaperDownloadsByWhitepaperId :many
-- Retrieves all download records for a specific whitepaper (no pagination).
--
-- Parameters:
--   $1 (INTEGER) - whitepaper_id: Whitepaper to fetch downloads for
-- Returns: []WhitepaperDownload - Array of all download records for the whitepaper
--
-- Sorting: created_at DESC - Newest downloads first
-- Use case: Viewing all leads captured for a specific whitepaper
SELECT
    id, name, email, company, designation, marketing_consent, created_at
FROM whitepaper_downloads
WHERE whitepaper_id = ?
ORDER BY created_at DESC;
