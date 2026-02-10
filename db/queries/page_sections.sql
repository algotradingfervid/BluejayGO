-- ====================================================================
-- PAGE SECTIONS
-- ====================================================================

-- name: GetPageSection :one
SELECT * FROM page_sections
WHERE page_key = ? AND section_key = ? AND is_active = 1;

-- name: ListPageSections :many
SELECT * FROM page_sections
WHERE page_key = ? AND is_active = 1
ORDER BY display_order ASC;

-- name: GetPageSectionByID :one
SELECT * FROM page_sections WHERE id = ?;

-- name: ListAllPageSections :many
SELECT * FROM page_sections ORDER BY page_key, display_order ASC;

-- name: UpdatePageSection :exec
UPDATE page_sections SET
    heading = ?,
    subheading = ?,
    description = ?,
    label = ?,
    primary_button_text = ?,
    primary_button_url = ?,
    secondary_button_text = ?,
    secondary_button_url = ?,
    is_active = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
