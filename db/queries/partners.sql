-- ====================================================================
-- PARTNERS QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing business partners and
-- their testimonials.
--
-- Entities:
--   - partners: Company partner records with tier/category
--   - partner_testimonials: Customer quotes and testimonials from partners
--
-- Related tables:
--   - partner_tiers: Defines partner levels (Platinum, Gold, Silver, etc.)
--
-- Features:
--   - Multi-tier partner management
--   - Active/inactive status for visibility control
--   - Featured partners for homepage highlights
--   - Partner testimonials with author details
--   - Display ordering within tiers
-- ====================================================================

-- ====================================================================
-- PARTNERS
-- ====================================================================

-- name: ListPartnersByTier :many
-- Retrieves all active partners organized by tier, with tier metadata.
--
-- Parameters: none
-- Returns: []Partner - Array of active partners with tier information
--
-- JOIN logic:
--   - JOIN partner_tiers pt ON p.tier_id = pt.id
--     Links partners to their tier category, adds tier_name and tier_level
--
-- Filtering: p.is_active = 1 - Only visible partners
--
-- Sorting logic:
--   1. pt.sort_order ASC - Higher tier partners first (Platinum before Gold)
--   2. p.display_order ASC - Custom ordering within each tier
--
-- Use case: Partners page display, showing partners grouped by tier
-- Note: Returns tier_name and tier_level (sort_order) as additional columns
SELECT p.*, pt.name AS tier_name, pt.sort_order AS tier_level
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
WHERE p.is_active = 1
ORDER BY pt.sort_order ASC, p.display_order ASC;

-- name: ListPartnersByTierID :many
-- Retrieves all active partners for a specific tier.
--
-- Parameters:
--   $1 (INTEGER) - tier_id: Partner tier ID to filter by
-- Returns: []Partner - Array of active partners in the specified tier
--
-- Filtering:
--   - tier_id = ? - Matches specific tier
--   - is_active = 1 - Only visible partners
--
-- Sorting: display_order ASC - Partners in custom-defined order
--
-- Use case: Displaying partners filtered by tier (e.g., "Show all Platinum Partners")
SELECT * FROM partners
WHERE tier_id = ? AND is_active = 1
ORDER BY display_order ASC;

-- name: GetPartner :one
-- Retrieves a single partner by ID with tier information.
--
-- Parameters:
--   $1 (INTEGER) - partner ID
-- Returns: Partner - Single partner record with tier metadata
--
-- JOIN logic:
--   - JOIN partner_tiers pt ON p.tier_id = pt.id
--     Adds tier_name and tier_level (sort_order) columns
--
-- Use case: Partner detail view, editing partner with tier context
-- Note: Does NOT filter by is_active, returns inactive partners for admin use
SELECT p.*, pt.name AS tier_name, pt.sort_order AS tier_level
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
WHERE p.id = ?;

-- name: CreatePartner :one
-- Creates a new partner record.
--
-- Parameters:
--   $1 (TEXT) - name: Partner company name
--   $2 (INTEGER) - tier_id: Partner tier category ID
--   $3 (TEXT) - logo_url: Path/URL to partner logo image
--   $4 (TEXT) - icon: Icon identifier or CSS class (optional)
--   $5 (TEXT) - website_url: Partner website URL (optional)
--   $6 (TEXT) - description: Partner description or relationship details (optional)
--   $7 (INTEGER) - display_order: Position within tier for ordering
--
-- Returns: Partner - The newly created partner with auto-generated ID and timestamps
--
-- Note: is_active defaults to 1 (true) via database schema, is_featured defaults to 0
INSERT INTO partners (
    name, tier_id, logo_url, icon, website_url, description, display_order
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdatePartner :one
-- Updates an existing partner record.
--
-- Parameters:
--   $1 (TEXT) - name: Updated company name
--   $2 (INTEGER) - tier_id: Updated tier assignment
--   $3 (TEXT) - logo_url: Updated logo path
--   $4 (TEXT) - icon: Updated icon identifier
--   $5 (TEXT) - website_url: Updated website URL
--   $6 (TEXT) - description: Updated description
--   $7 (INTEGER) - display_order: Updated display position
--   $8 (BOOLEAN) - is_active: Updated visibility status
--   $9 (INTEGER) - id: Partner ID to update
--
-- Returns: Partner - The updated partner record
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
UPDATE partners
SET name = ?, tier_id = ?, logo_url = ?, icon = ?,
    website_url = ?, description = ?, display_order = ?,
    is_active = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeletePartner :exec
-- Permanently deletes a partner record.
--
-- Parameters:
--   $1 (INTEGER) - partner ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Will cascade delete related testimonials if foreign key configured
-- Note: Consider setting is_active=0 for soft delete instead
DELETE FROM partners WHERE id = ?;

-- ====================================================================
-- PARTNER TESTIMONIALS
-- ====================================================================

-- name: ListActiveTestimonials :many
-- Retrieves all active partner testimonials with partner details for display.
--
-- Parameters: none
-- Returns: []PartnerTestimonial - Array of active testimonials with partner metadata
--
-- JOIN logic:
--   - JOIN partners p ON pt.partner_id = p.id
--     Adds partner_name, partner_logo_url, partner_icon for display context
--
-- Filtering: pt.is_active = 1 - Only visible testimonials
-- Sorting: pt.display_order ASC - Custom-defined testimonial sequence
--
-- Use case: Homepage testimonials section, partners page testimonials carousel
SELECT pt.*, p.name AS partner_name, p.logo_url AS partner_logo_url, p.icon AS partner_icon
FROM partner_testimonials pt
JOIN partners p ON pt.partner_id = p.id
WHERE pt.is_active = 1
ORDER BY pt.display_order ASC;

-- name: GetTestimonial :one
-- Retrieves a single testimonial by ID with partner name.
--
-- Parameters:
--   $1 (INTEGER) - testimonial ID
-- Returns: PartnerTestimonial - Single testimonial with partner_name
--
-- JOIN logic: Adds partner_name for display/admin context
-- Use case: Editing testimonial, displaying single testimonial detail
SELECT pt.*, p.name AS partner_name
FROM partner_testimonials pt
JOIN partners p ON pt.partner_id = p.id
WHERE pt.id = ?;

-- name: CreateTestimonial :one
-- Creates a new partner testimonial.
--
-- Parameters:
--   $1 (INTEGER) - partner_id: ID of partner providing the testimonial
--   $2 (TEXT) - quote: Testimonial text/quote
--   $3 (TEXT) - author_name: Name of person giving testimonial
--   $4 (TEXT) - author_title: Job title of testimonial author
--   $5 (INTEGER) - display_order: Position in testimonial sequence
--
-- Returns: PartnerTestimonial - The newly created testimonial with auto-generated ID
--
-- Note: is_active defaults to 1 (true) via database schema
INSERT INTO partner_testimonials (
    partner_id, quote, author_name, author_title, display_order
) VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateTestimonial :one
-- Updates an existing partner testimonial.
--
-- Parameters:
--   $1 (INTEGER) - partner_id: Updated partner assignment
--   $2 (TEXT) - quote: Updated testimonial text
--   $3 (TEXT) - author_name: Updated author name
--   $4 (TEXT) - author_title: Updated author job title
--   $5 (INTEGER) - display_order: Updated display position
--   $6 (BOOLEAN) - is_active: Updated visibility status
--   $7 (INTEGER) - id: Testimonial ID to update
--
-- Returns: PartnerTestimonial - The updated testimonial record
UPDATE partner_testimonials
SET partner_id = ?, quote = ?, author_name = ?,
    author_title = ?, display_order = ?, is_active = ?
WHERE id = ?
RETURNING *;

-- name: DeleteTestimonial :exec
-- Permanently deletes a partner testimonial.
--
-- Parameters:
--   $1 (INTEGER) - testimonial ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Hard delete, consider setting is_active=0 for soft delete instead
DELETE FROM partner_testimonials WHERE id = ?;

-- ====================================================================
-- PARTNER LISTING QUERIES
-- ====================================================================

-- name: ListAllPartners :many
-- Retrieves all partners (active and inactive) with tier information.
--
-- Parameters: none
-- Returns: []Partner - Array of all partners with tier_name
--
-- JOIN logic: Adds tier_name from partner_tiers table
--
-- Sorting:
--   1. pt.sort_order ASC - Higher tiers first
--   2. p.display_order ASC - Custom order within tier
--
-- Use case: Admin partner management, bulk operations
-- Note: Returns ALL partners regardless of is_active status
SELECT p.*, pt.name AS tier_name
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
ORDER BY pt.sort_order ASC, p.display_order ASC;

-- name: ListFeaturedPartners :many
-- Retrieves a limited number of featured active partners.
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Maximum number of featured partners to return
-- Returns: []Partner - Array of featured partners with tier_name
--
-- JOIN logic: Adds tier_name from partner_tiers table
--
-- Filtering:
--   - p.is_featured = 1 - Only partners marked as featured
--   - p.is_active = 1 - Only visible partners
--
-- Sorting: p.display_order ASC - Featured partners in custom order
-- LIMIT: Controls how many featured partners to display (e.g., 6 for homepage)
--
-- Use case: Homepage featured partners section, highlighting key partnerships
SELECT p.*, pt.name AS tier_name
FROM partners p
JOIN partner_tiers pt ON p.tier_id = pt.id
WHERE p.is_featured = 1 AND p.is_active = 1
ORDER BY p.display_order ASC
LIMIT ?;
