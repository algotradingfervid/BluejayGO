-- ====================================================================
-- WHITEPAPERS
-- ====================================================================

-- name: ListPublishedWhitepapers :many
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.file_size_bytes, w.page_count,
    w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.is_published = 1
ORDER BY w.published_date DESC, w.id DESC;

-- name: ListPublishedWhitepapersByTopic :many
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.file_size_bytes, w.page_count,
    w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.is_published = 1 AND w.topic_id = ?
ORDER BY w.published_date DESC, w.id DESC;

-- name: GetWhitepaperBySlug :one
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.pdf_file_path, w.file_size_bytes,
    w.page_count, w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    w.meta_description, w.meta_title, w.og_image,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.slug = ? AND w.is_published = 1;

-- name: GetWhitepaperBySlugIncludeDrafts :one
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.pdf_file_path, w.file_size_bytes,
    w.page_count, w.published_date, w.cover_color_from, w.cover_color_to, w.download_count,
    w.meta_description, w.meta_title, w.og_image,
    t.name as topic_name, t.color_hex as topic_color_hex
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
WHERE w.slug = ?;

-- name: GetWhitepaperLearningPoints :many
SELECT id, point_text, display_order
FROM whitepaper_learning_points
WHERE whitepaper_id = ?
ORDER BY display_order ASC, id ASC;

-- name: GetRelatedWhitepapers :many
SELECT
    w.id, w.title, w.slug, w.cover_color_from, w.cover_color_to
FROM whitepapers w
WHERE w.is_published = 1
  AND w.id != ?
  AND w.topic_id = ?
ORDER BY w.published_date DESC
LIMIT 3;

-- name: CreateWhitepaperDownload :one
INSERT INTO whitepaper_downloads (
    whitepaper_id, name, email, company, designation, marketing_consent, ip_address, user_agent
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, created_at;

-- name: IncrementWhitepaperDownloadCount :exec
UPDATE whitepapers
SET download_count = download_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: CountPublishedWhitepapers :one
SELECT COUNT(*) FROM whitepapers WHERE is_published = 1;

-- name: CountPublishedWhitepapersByTopic :one
SELECT COUNT(*) FROM whitepapers WHERE is_published = 1 AND topic_id = ?;

-- ====================================================================
-- WHITEPAPERS - ADMIN
-- ====================================================================

-- name: ListAllWhitepapers :many
SELECT
    w.id, w.title, w.slug, w.topic_id, w.is_published, w.published_date, w.download_count,
    w.created_at, w.updated_at,
    t.name as topic_name
FROM whitepapers w
INNER JOIN whitepaper_topics t ON w.topic_id = t.id
ORDER BY w.created_at DESC;

-- name: ListWhitepapersAdminFiltered :many
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
SELECT COUNT(*) FROM whitepapers w
WHERE
    (CASE WHEN @filter_search = '' THEN 1 ELSE w.title LIKE '%' || @filter_search || '%' END)
    AND (CASE WHEN @filter_topic = 0 THEN 1 ELSE w.topic_id = @filter_topic END)
    AND (CASE WHEN @filter_status = '' THEN 1
         WHEN @filter_status = 'published' THEN w.is_published = 1
         WHEN @filter_status = 'draft' THEN w.is_published = 0
         ELSE 1 END);

-- name: ListWhitepaperDownloadsFiltered :many
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
SELECT COUNT(*) FROM whitepaper_downloads wd
WHERE
    (CASE WHEN @filter_whitepaper = 0 THEN 1 ELSE wd.whitepaper_id = @filter_whitepaper END)
    AND (CASE WHEN @filter_date_from = '' THEN 1 ELSE wd.created_at >= @filter_date_from END)
    AND (CASE WHEN @filter_date_to = '' THEN 1 ELSE wd.created_at <= @filter_date_to || ' 23:59:59' END);

-- name: ListWhitepaperTopicsWithCount :many
SELECT t.id, t.name, t.slug, t.color_hex, t.icon, t.description, t.sort_order, t.created_at, t.updated_at,
    (SELECT COUNT(*) FROM whitepapers w WHERE w.topic_id = t.id) as whitepaper_count
FROM whitepaper_topics t
ORDER BY t.sort_order ASC, t.name ASC;

-- name: GetWhitepaperByID :one
SELECT
    w.id, w.title, w.slug, w.description, w.topic_id, w.pdf_file_path, w.file_size_bytes,
    w.page_count, w.published_date, w.is_published, w.cover_color_from, w.cover_color_to,
    w.meta_description, w.download_count, w.created_at, w.updated_at
FROM whitepapers w
WHERE w.id = ?;

-- name: CreateWhitepaper :one
INSERT INTO whitepapers (
    title, slug, description, topic_id, pdf_file_path, file_size_bytes, page_count,
    published_date, is_published, cover_color_from, cover_color_to, meta_description
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateWhitepaper :exec
UPDATE whitepapers
SET title = ?, slug = ?, description = ?, topic_id = ?, pdf_file_path = ?,
    file_size_bytes = ?, page_count = ?, published_date = ?, is_published = ?,
    cover_color_from = ?, cover_color_to = ?, meta_description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteWhitepaper :exec
DELETE FROM whitepapers WHERE id = ?;

-- name: CreateWhitepaperLearningPoint :one
INSERT INTO whitepaper_learning_points (whitepaper_id, point_text, display_order)
VALUES (?, ?, ?)
RETURNING id;

-- name: DeleteWhitepaperLearningPoints :exec
DELETE FROM whitepaper_learning_points WHERE whitepaper_id = ?;

-- name: ListWhitepaperDownloads :many
SELECT
    wd.id, wd.whitepaper_id, wd.name, wd.email, wd.company, wd.designation,
    wd.marketing_consent, wd.created_at,
    w.title as whitepaper_title
FROM whitepaper_downloads wd
INNER JOIN whitepapers w ON wd.whitepaper_id = w.id
ORDER BY wd.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountWhitepaperDownloads :one
SELECT COUNT(*) FROM whitepaper_downloads;

-- name: ListWhitepaperDownloadsByWhitepaperId :many
SELECT
    id, name, email, company, designation, marketing_consent, created_at
FROM whitepaper_downloads
WHERE whitepaper_id = ?
ORDER BY created_at DESC;
