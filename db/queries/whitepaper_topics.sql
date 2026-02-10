-- ====================================================================
-- WHITEPAPER TOPICS QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing whitepaper topic categories.
--
-- Entity: whitepaper_topics table
-- Purpose: Categorize whitepapers by subject/topic (e.g., "IoT", "Cloud", "Security")
-- Related: whitepapers table references whitepaper_topics via topic_id foreign key
--
-- Features:
--   - Color-coded topics for visual organization (color_hex)
--   - Icon support for topic badges
--   - Custom display ordering
--   - Used for filtering whitepapers by topic
-- ====================================================================

-- name: ListWhitepaperTopics :many
-- Retrieves all whitepaper topics ordered by display priority, then alphabetically.
--
-- Parameters: none
-- Returns: []WhitepaperTopic - Array of all topic records
--
-- Sorting logic:
--   1. sort_order ASC - Custom display order (lower numbers appear first)
--   2. name ASC - Alphabetical fallback for same sort_order values
--
-- Use case: Topic navigation, whitepaper filters, admin topic listing
SELECT * FROM whitepaper_topics ORDER BY sort_order ASC, name ASC;

-- name: GetWhitepaperTopic :one
-- Retrieves a single whitepaper topic by its primary key ID.
--
-- Parameters:
--   $1 (INTEGER) - topic ID
-- Returns: WhitepaperTopic - Single topic record or error if not found
--
-- Use case: Editing a specific topic, fetching topic details
SELECT * FROM whitepaper_topics WHERE id = ? LIMIT 1;

-- name: GetWhitepaperTopicBySlug :one
-- Retrieves a single whitepaper topic by its URL-safe slug identifier.
--
-- Parameters:
--   $1 (TEXT) - topic slug (e.g., "iot", "cloud-computing", "cybersecurity")
-- Returns: WhitepaperTopic - Single topic record or error if not found
--
-- Use case: Frontend topic page routing, filtering whitepapers by topic URL
-- Note: Slugs should be unique (enforced by database constraint)
SELECT * FROM whitepaper_topics WHERE slug = ? LIMIT 1;

-- name: CreateWhitepaperTopic :one
-- Creates a new whitepaper topic category.
--
-- Parameters:
--   $1 (TEXT) - name: Display name of the topic
--   $2 (TEXT) - slug: URL-safe identifier
--   $3 (TEXT) - color_hex: Hex color code for topic badge/card (e.g., "#3B82F6")
--   $4 (TEXT) - icon: Icon identifier or CSS class (optional)
--   $5 (TEXT) - description: Topic description (optional)
--   $6 (INTEGER) - sort_order: Display position (lower = higher priority)
--
-- Returns: WhitepaperTopic - The newly created topic with auto-generated ID and timestamps
--
-- Note: color_hex is used for visual differentiation in topic badges and cards
INSERT INTO whitepaper_topics (name, slug, color_hex, icon, description, sort_order)
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateWhitepaperTopic :one
-- Updates an existing whitepaper topic.
--
-- Parameters:
--   $1 (TEXT) - name: Updated display name
--   $2 (TEXT) - slug: Updated URL-safe identifier
--   $3 (TEXT) - color_hex: Updated color code
--   $4 (TEXT) - icon: Updated icon identifier
--   $5 (TEXT) - description: Updated description
--   $6 (INTEGER) - sort_order: Updated display position
--   $7 (INTEGER) - id: Topic ID to update
--
-- Returns: WhitepaperTopic - The updated topic record
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
UPDATE whitepaper_topics SET name = ?, slug = ?, color_hex = ?, icon = ?, description = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteWhitepaperTopic :exec
-- Permanently deletes a whitepaper topic.
--
-- Parameters:
--   $1 (INTEGER) - topic ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Will fail if whitepapers reference this topic (foreign key constraint)
-- Note: Reassign or delete whitepapers in this topic before deletion
DELETE FROM whitepaper_topics WHERE id = ?;
