-- ====================================================================
-- PAGE SECTIONS QUERY FILE
-- ====================================================================
-- This file contains SQL queries for managing configurable page sections
-- (heroes, CTAs, testimonials, stats, etc.) across different pages.
--
-- Entity: page_sections table
-- Purpose: Store content for modular page sections that can be edited via admin
-- Features:
--   - Multi-page support via page_key (e.g., "homepage", "about", "contact")
--   - Section-specific content via section_key (e.g., "hero", "cta", "stats")
--   - Flexible content fields (heading, subheading, description, buttons, etc.)
--   - Visibility control via is_active flag
--   - Display ordering for multiple sections on same page
--
-- Common page_key values: homepage, about, products, solutions, contact
-- Common section_key values: hero, cta, stats, testimonials, features
-- ====================================================================

-- name: GetPageSection :one
-- Retrieves a specific active section for a given page.
--
-- Parameters:
--   $1 (TEXT) - page_key: Page identifier (e.g., "homepage", "about")
--   $2 (TEXT) - section_key: Section identifier (e.g., "hero", "cta")
-- Returns: PageSection - Single section record or error if not found/inactive
--
-- Filtering logic:
--   - page_key = ? - Matches specific page
--   - section_key = ? - Matches specific section type
--   - is_active = 1 - Only returns visible sections
--
-- Use case: Fetching hero content for homepage, CTA for contact page, etc.
SELECT * FROM page_sections
WHERE page_key = ? AND section_key = ? AND is_active = 1;

-- name: ListPageSections :many
-- Retrieves all active sections for a specific page in display order.
--
-- Parameters:
--   $1 (TEXT) - page_key: Page identifier
-- Returns: []PageSection - Array of active sections for the page
--
-- Filtering: is_active = 1 - Only visible sections
-- Sorting: display_order ASC - Sections appear in configured order
--
-- Use case: Rendering all sections for a page, building page content dynamically
SELECT * FROM page_sections
WHERE page_key = ? AND is_active = 1
ORDER BY display_order ASC;

-- name: GetPageSectionByID :one
-- Retrieves a single page section by its primary key ID (ignores is_active).
--
-- Parameters:
--   $1 (INTEGER) - section ID
-- Returns: PageSection - Single section record or error if not found
--
-- Use case: Admin editing, fetching section for update regardless of active status
-- Note: Does NOT filter by is_active, so returns draft/inactive sections
SELECT * FROM page_sections WHERE id = ?;

-- name: ListAllPageSections :many
-- Retrieves all page sections across all pages, grouped and ordered.
--
-- Parameters: none
-- Returns: []PageSection - Array of all sections across entire site
--
-- Sorting logic:
--   1. page_key - Groups sections by page (alphabetical)
--   2. display_order ASC - Orders sections within each page
--
-- Use case: Admin section management dashboard, bulk operations, site-wide overview
-- Note: Returns ALL sections including inactive ones (is_active not filtered)
SELECT * FROM page_sections ORDER BY page_key, display_order ASC;

-- name: UpdatePageSection :exec
-- Updates the content fields of an existing page section.
--
-- Parameters:
--   $1 (TEXT) - heading: Main heading/title for the section
--   $2 (TEXT) - subheading: Secondary heading or subtitle
--   $3 (TEXT) - description: Longer description or body text
--   $4 (TEXT) - label: Short label or tag (e.g., "New", "Featured")
--   $5 (TEXT) - primary_button_text: Text for primary CTA button
--   $6 (TEXT) - primary_button_url: URL for primary button
--   $7 (TEXT) - secondary_button_text: Text for secondary button
--   $8 (TEXT) - secondary_button_url: URL for secondary button
--   $9 (BOOLEAN) - is_active: Visibility flag (1=visible, 0=hidden)
--   $10 (INTEGER) - id: Section ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP to track modifications
-- Use case: Admin content editing for heroes, CTAs, and other configurable sections
--
-- Content structure supports common section patterns:
--   - Hero: heading, subheading, description, primary_button (CTA)
--   - CTA: heading, description, primary_button, secondary_button
--   - Feature: label, heading, description
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
