-- ====================================================================
-- FOOTER CONFIGURATION QUERIES
-- ====================================================================
-- This file manages footer content and structure including:
-- - Global footer settings (columns layout, social links, copyright)
-- - Footer column items (text blocks, link groups per column)
-- - Footer links (navigation links within column items)
-- - Footer legal links (bottom legal/policy links)
--
-- Managed entities:
-- - settings: global footer configuration (single row, id=1)
-- - footer_column_items: content blocks in footer columns
-- - footer_links: individual links within column items
-- - footer_legal_links: bottom legal/policy links
--
-- Key concepts:
-- - column_index: which footer column (0, 1, 2, etc.)
-- - sort_order: display order within column
-- - type: 'text', 'links', 'newsletter', etc.
-- ====================================================================

-- ====================================================================
-- FOOTER SETTINGS
-- ====================================================================

-- name: UpdateFooterSettings :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Updates global footer configuration settings
-- Parameters (5 positional):
--   1. footer_columns (INTEGER): number of footer columns (2, 3, 4, etc.)
--   2. footer_bg_style (TEXT): background style ('dark', 'light', 'gradient', etc.)
--   3. footer_show_social (BOOLEAN): whether to display social media links
--   4. footer_social_style (TEXT): social links style ('icons', 'text', etc.)
--   5. footer_copyright (TEXT): copyright text
-- Return type: none
-- Note: Always updates row with id=1 (single settings row pattern)
UPDATE settings
SET footer_columns = ?,
    footer_bg_style = ?,
    footer_show_social = ?,
    footer_social_style = ?,
    footer_copyright = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- ====================================================================
-- FOOTER COLUMN ITEMS
-- ====================================================================

-- name: ListFooterColumnItems :many
-- Purpose: Lists all footer column content blocks
-- ORDER BY: column_index (left to right), then sort_order (top to bottom)
SELECT * FROM footer_column_items ORDER BY column_index, sort_order;

-- name: GetFooterColumnItem :one
-- Purpose: Retrieves specific footer column item for editing
SELECT * FROM footer_column_items WHERE id = ? LIMIT 1;

-- name: CreateFooterColumnItem :one
-- Purpose: Creates new content block in footer column
-- Parameters: column_index, type, heading, content, sort_order
-- type examples: 'text', 'links', 'newsletter'
INSERT INTO footer_column_items (column_index, type, heading, content, sort_order)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateFooterColumnItem :exec
-- Purpose: Updates existing footer column item
UPDATE footer_column_items
SET column_index = ?, type = ?, heading = ?, content = ?, sort_order = ?
WHERE id = ?;

-- name: DeleteFooterColumnItem :exec
-- Purpose: Removes specific footer column item
DELETE FROM footer_column_items WHERE id = ?;

-- name: DeleteFooterColumnItemsByIndex :exec
-- Purpose: Bulk delete all items in columns >= specified index
-- Use case: When reducing number of footer columns, remove excess column content
-- Parameters:
--   1. min_column_index (INTEGER): delete all items with column_index >= this
DELETE FROM footer_column_items WHERE column_index >= ?;

-- ====================================================================
-- FOOTER LINKS (within column items)
-- ====================================================================

-- name: ListFooterLinks :many
-- Purpose: Lists links belonging to a specific footer column item
-- Parameters:
--   1. column_item_id (INTEGER): parent column item
-- Use case: Getting links for a "Quick Links" or "Resources" column block
SELECT * FROM footer_links WHERE column_item_id = ? ORDER BY sort_order;

-- name: ListAllFooterLinks :many
-- Purpose: Lists all footer links across all column items
-- ORDER BY: groups by column_item_id, then sort_order within each group
SELECT * FROM footer_links ORDER BY column_item_id, sort_order;

-- name: CreateFooterLink :one
-- Purpose: Creates a new link within a footer column item
-- Parameters: column_item_id (parent), label (link text), url, sort_order
INSERT INTO footer_links (column_item_id, label, url, sort_order)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteFooterLinksByColumnItem :exec
-- Purpose: Bulk delete all links belonging to a column item
-- Use case: When deleting or rebuilding a footer column's links
-- Parameters:
--   1. column_item_id (INTEGER): delete all links for this column item
DELETE FROM footer_links WHERE column_item_id = ?;

-- ====================================================================
-- FOOTER LEGAL LINKS (bottom bar)
-- ====================================================================

-- name: ListFooterLegalLinks :many
-- Purpose: Lists legal/policy links shown in footer bottom bar
-- Examples: "Privacy Policy", "Terms of Service", "Accessibility"
-- ORDER BY sort_order: custom display sequence
SELECT * FROM footer_legal_links ORDER BY sort_order;

-- name: CreateFooterLegalLink :one
-- Purpose: Creates a new legal link in footer bottom bar
-- Parameters: label (link text), url, sort_order
INSERT INTO footer_legal_links (label, url, sort_order)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteAllFooterLegalLinks :exec
-- Purpose: Removes all legal links (used when rebuilding legal link set)
-- Note: No WHERE clause - deletes entire table contents
DELETE FROM footer_legal_links;
