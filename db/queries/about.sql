-- ====================================================================
-- ABOUT PAGE QUERIES
-- ====================================================================
-- This file contains all SQL queries for managing "About Us" page content
-- including company overview, mission/vision/values, core values list,
-- company milestones timeline, and certifications/credentials display.
--
-- Managed entities:
-- - company_overview: main about page hero content
-- - mission_vision_values: strategic statements section
-- - core_values: individual value items with icons
-- - milestones: company history timeline
-- - certifications: professional credentials and certifications
-- ====================================================================

-- ====================================================================
-- COMPANY OVERVIEW
-- ====================================================================

-- name: GetCompanyOverview :one
-- sqlc annotation: :one returns a single row or error if not found
-- Purpose: Retrieves the most recent company overview content for the About page
-- Parameters: none
-- Return type: single company_overview row
-- Note: Uses ORDER BY id DESC to get the latest entry (highest ID)
SELECT * FROM company_overview ORDER BY id DESC LIMIT 1;

-- name: UpsertCompanyOverview :one
-- sqlc annotation: :one returns the newly inserted row
-- Purpose: Creates a new company overview entry (upsert pattern via insert-only)
-- Parameters (7 positional):
--   1. headline (TEXT): main hero headline
--   2. tagline (TEXT): supporting tagline
--   3. description_main (TEXT): primary description paragraph
--   4. description_secondary (TEXT): secondary description
--   5. description_tertiary (TEXT): tertiary description
--   6. hero_image_url (TEXT): hero section image
--   7. company_image_url (TEXT): additional company image
-- Return type: complete inserted row with generated ID and timestamps
-- Note: Called "Upsert" but actually inserts new row each time
INSERT INTO company_overview (
    headline, tagline, description_main, description_secondary,
    description_tertiary, hero_image_url, company_image_url
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- ====================================================================
-- MISSION, VISION & VALUES
-- ====================================================================

-- name: GetMissionVisionValues :one
-- sqlc annotation: :one returns single row or error
-- Purpose: Retrieves the latest mission/vision/values content block
-- Parameters: none
-- Return type: single mission_vision_values row
-- Note: ORDER BY id DESC ensures we get the most recent version
SELECT * FROM mission_vision_values ORDER BY id DESC LIMIT 1;

-- name: UpsertMissionVisionValues :one
-- sqlc annotation: :one returns inserted row
-- Purpose: Creates new mission/vision/values entry
-- Parameters (6 positional):
--   1. mission (TEXT): company mission statement
--   2. vision (TEXT): company vision statement
--   3. values_summary (TEXT): overview of company values
--   4. mission_icon (TEXT): icon identifier for mission
--   5. vision_icon (TEXT): icon identifier for vision
--   6. values_icon (TEXT): icon identifier for values
-- Return type: complete inserted row
INSERT INTO mission_vision_values (
    mission, vision, values_summary,
    mission_icon, vision_icon, values_icon
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- ====================================================================
-- CORE VALUES
-- ====================================================================

-- name: ListCoreValues :many
-- sqlc annotation: :many returns slice of rows (0 or more)
-- Purpose: Lists all core values for display on About page
-- Parameters: none
-- Return type: slice of core_values rows
-- Note: ORDER BY display_order ensures consistent presentation order
SELECT * FROM core_values ORDER BY display_order ASC;

-- name: GetCoreValue :one
-- sqlc annotation: :one returns single row or error
-- Purpose: Retrieves a specific core value by ID for editing
-- Parameters:
--   1. id (INTEGER): core value primary key
-- Return type: single core_values row
SELECT * FROM core_values WHERE id = ? LIMIT 1;

-- name: CreateCoreValue :one
-- sqlc annotation: :one returns the created row
-- Purpose: Creates a new core value entry
-- Parameters (4 positional):
--   1. title (TEXT): value name/title
--   2. description (TEXT): detailed explanation
--   3. icon (TEXT): icon identifier (e.g., "shield", "star")
--   4. display_order (INTEGER): sort position in list
-- Return type: complete inserted row with generated ID
INSERT INTO core_values (title, description, icon, display_order)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdateCoreValue :one
-- sqlc annotation: :one returns the updated row
-- Purpose: Updates an existing core value
-- Parameters (5 positional):
--   1. title (TEXT): updated title
--   2. description (TEXT): updated description
--   3. icon (TEXT): updated icon
--   4. display_order (INTEGER): updated sort position
--   5. id (INTEGER): which core value to update (WHERE clause)
-- Return type: updated row with new values
UPDATE core_values SET title = ?, description = ?, icon = ?, display_order = ?
WHERE id = ? RETURNING *;

-- name: DeleteCoreValue :exec
-- sqlc annotation: :exec returns no data, only error or nil
-- Purpose: Permanently removes a core value entry
-- Parameters:
--   1. id (INTEGER): core value to delete
-- Return type: none (exec returns only error status)
DELETE FROM core_values WHERE id = ?;

-- ====================================================================
-- MILESTONES / COMPANY TIMELINE
-- ====================================================================

-- name: ListMilestones :many
-- sqlc annotation: :many returns slice of milestone rows
-- Purpose: Lists all company milestones for timeline display
-- Parameters: none
-- Return type: slice of milestones rows
-- Note: ORDER BY display_order for chronological or custom ordering
SELECT * FROM milestones ORDER BY display_order ASC;

-- name: GetMilestone :one
-- sqlc annotation: :one returns single milestone or error
-- Purpose: Retrieves specific milestone for editing
-- Parameters:
--   1. id (INTEGER): milestone primary key
-- Return type: single milestones row
SELECT * FROM milestones WHERE id = ? LIMIT 1;

-- name: CreateMilestone :one
-- sqlc annotation: :one returns created milestone
-- Purpose: Creates a new milestone entry for company timeline
-- Parameters (5 positional):
--   1. year (TEXT): year of milestone (stored as text for flexibility)
--   2. title (TEXT): milestone heading
--   3. description (TEXT): milestone details
--   4. is_current (BOOLEAN): marks if this is the current/ongoing milestone
--   5. display_order (INTEGER): sort position in timeline
-- Return type: complete inserted row
INSERT INTO milestones (year, title, description, is_current, display_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateMilestone :one
-- sqlc annotation: :one returns updated row
-- Purpose: Updates an existing milestone entry
-- Parameters (6 positional):
--   1-5. updated field values (year, title, description, is_current, display_order)
--   6. id (INTEGER): which milestone to update
-- Return type: updated milestone row
UPDATE milestones SET year = ?, title = ?, description = ?, is_current = ?, display_order = ?
WHERE id = ? RETURNING *;

-- name: DeleteMilestone :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Removes a milestone from the timeline
-- Parameters:
--   1. id (INTEGER): milestone to delete
-- Return type: none
DELETE FROM milestones WHERE id = ?;

-- ====================================================================
-- CERTIFICATIONS & CREDENTIALS
-- ====================================================================

-- name: ListCertifications :many
-- sqlc annotation: :many returns slice of certification rows
-- Purpose: Lists all certifications/credentials for About page display
-- Parameters: none
-- Return type: slice of certifications rows
-- Note: ORDER BY display_order for custom presentation sequence
SELECT * FROM certifications ORDER BY display_order ASC;

-- name: GetCertification :one
-- sqlc annotation: :one returns single certification or error
-- Purpose: Retrieves specific certification for editing
-- Parameters:
--   1. id (INTEGER): certification primary key
-- Return type: single certifications row
SELECT * FROM certifications WHERE id = ? LIMIT 1;

-- name: CreateCertification :one
-- sqlc annotation: :one returns created certification
-- Purpose: Adds a new certification/credential entry
-- Parameters (5 positional):
--   1. name (TEXT): full certification name
--   2. abbreviation (TEXT): short form (e.g., "ISO 9001")
--   3. description (TEXT): certification details
--   4. icon (TEXT): icon identifier for display
--   5. display_order (INTEGER): sort position
-- Return type: complete inserted row with ID
INSERT INTO certifications (name, abbreviation, description, icon, display_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateCertification :one
-- sqlc annotation: :one returns updated row
-- Purpose: Updates an existing certification entry
-- Parameters (6 positional):
--   1-5. updated field values
--   6. id (INTEGER): which certification to update
-- Return type: updated certification row
UPDATE certifications SET name = ?, abbreviation = ?, description = ?, icon = ?, display_order = ?
WHERE id = ? RETURNING *;

-- name: DeleteCertification :exec
-- sqlc annotation: :exec returns no data
-- Purpose: Removes a certification from the list
-- Parameters:
--   1. id (INTEGER): certification to delete
-- Return type: none
DELETE FROM certifications WHERE id = ?;
