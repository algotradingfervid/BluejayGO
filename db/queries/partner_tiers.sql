-- ====================================================================
-- PARTNER TIERS QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing partner tier categories.
--
-- Entity: partner_tiers table
-- Purpose: Categorize partners by level/tier (e.g., Platinum, Gold, Silver, Bronze)
-- Related: partners table references partner_tiers via tier_id foreign key
--
-- Business logic:
--   - Higher tier partners (lower sort_order) appear first in listings
--   - Used to organize partners page, display partner logos by importance
--   - Tier descriptions can explain benefits or requirements
-- ====================================================================

-- name: ListPartnerTiers :many
-- Retrieves all partner tiers ordered by display priority, then alphabetically.
--
-- Parameters: none
-- Returns: []PartnerTier - Array of all tier records
--
-- Sorting logic:
--   1. sort_order ASC - Custom priority (lower = higher tier, e.g., Platinum=1, Gold=2)
--   2. name ASC - Alphabetical fallback for tiers with same sort_order
--
-- Use case: Displaying tier dropdown in partner form, organizing partners page by tier
SELECT * FROM partner_tiers ORDER BY sort_order ASC, name ASC;

-- name: GetPartnerTier :one
-- Retrieves a single partner tier by its primary key ID.
--
-- Parameters:
--   $1 (INTEGER) - tier ID
-- Returns: PartnerTier - Single tier record or error if not found
--
-- Use case: Editing a specific tier, fetching tier details for validation
SELECT * FROM partner_tiers WHERE id = ? LIMIT 1;

-- name: GetPartnerTierBySlug :one
-- Retrieves a single partner tier by its URL-safe slug identifier.
--
-- Parameters:
--   $1 (TEXT) - tier slug (e.g., "platinum", "gold", "silver")
-- Returns: PartnerTier - Single tier record or error if not found
--
-- Use case: Frontend routing, filtering partners by tier via URL parameter
-- Note: Slugs should be unique (enforced by database constraint)
SELECT * FROM partner_tiers WHERE slug = ? LIMIT 1;

-- name: CreatePartnerTier :one
-- Creates a new partner tier category.
--
-- Parameters:
--   $1 (TEXT) - name: Display name (e.g., "Platinum Partner", "Gold Partner")
--   $2 (TEXT) - slug: URL-safe identifier (e.g., "platinum", "gold")
--   $3 (TEXT) - description: Tier description or benefits (optional)
--   $4 (INTEGER) - sort_order: Display priority (lower = higher tier)
--
-- Returns: PartnerTier - The newly created tier with auto-generated ID and timestamps
--
-- Note: RETURNING * includes auto-generated created_at, updated_at timestamps
INSERT INTO partner_tiers (name, slug, description, sort_order)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdatePartnerTier :one
-- Updates an existing partner tier.
--
-- Parameters:
--   $1 (TEXT) - name: Updated display name
--   $2 (TEXT) - slug: Updated URL-safe identifier
--   $3 (TEXT) - description: Updated description
--   $4 (INTEGER) - sort_order: Updated display priority
--   $5 (INTEGER) - id: Tier ID to update
--
-- Returns: PartnerTier - The updated tier record
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
UPDATE partner_tiers SET name = ?, slug = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeletePartnerTier :exec
-- Permanently deletes a partner tier.
--
-- Parameters:
--   $1 (INTEGER) - tier ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Will fail if partners reference this tier (foreign key constraint)
-- Note: Consider archiving or reassigning partners before deletion
DELETE FROM partner_tiers WHERE id = ?;
