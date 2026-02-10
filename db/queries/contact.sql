-- ====================================================================
-- CONTACT SUBMISSIONS QUERIES
-- ====================================================================
-- This file manages contact form submissions and office location data.
-- Includes both public-facing queries (form submission, office listings)
-- and admin operations (submission management, filtering, status updates).
--
-- Managed entities:
-- - contact_submissions: form submissions with status tracking
-- - office_locations: physical office addresses and contact info
--
-- Key concepts:
-- - status: 'new', 'read', 'replied', 'archived' (tracks submission lifecycle)
-- - submission_type: categorizes inquiry type
-- - is_active: controls office location visibility on public site
-- ====================================================================

-- ====================================================================
-- PUBLIC CONTACT QUERIES
-- ====================================================================

-- name: CreateContactSubmission :one
-- sqlc annotation: :one returns minimal info after creating submission
-- Purpose: Records a new contact form submission from public website
-- Parameters (8 positional):
--   1. name (TEXT): submitter's full name
--   2. email (TEXT): submitter's email address
--   3. phone (TEXT): optional phone number
--   4. company (TEXT): optional company/organization name
--   5. inquiry_type (TEXT): type of inquiry (e.g., "sales", "support", "general")
--   6. message (TEXT): inquiry message content
--   7. ip_address (TEXT): submitter's IP for spam prevention
--   8. user_agent (TEXT): browser user agent for tracking
-- Return type: id and created_at only (minimal response)
-- Note: status defaults to 'new' via schema default
INSERT INTO contact_submissions (
    name, email, phone, company, inquiry_type, message, ip_address, user_agent
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, created_at;

-- name: GetActiveOfficeLocations :many
-- sqlc annotation: :many returns active office locations for public display
-- Purpose: Lists office locations shown on contact page
-- Parameters: none
-- Return type: slice of office_locations with contact details
-- WHERE: is_active = 1 (only show enabled locations)
-- ORDER BY:
--   - Primary: is_primary DESC (primary office first)
--   - Secondary: display_order ASC (custom sort order)
--   - Tertiary: id ASC (stable fallback)
SELECT id, name, address_line1, address_line2, city, state, postal_code, country, phone, email, is_primary
FROM office_locations
WHERE is_active = 1
ORDER BY is_primary DESC, display_order ASC, id ASC;

-- ====================================================================
-- CONTACT SUBMISSIONS - ADMIN QUERIES
-- ====================================================================
-- Admin queries support multiple filtering patterns:
-- - ListContactSubmissions: all submissions (basic pagination)
-- - ListContactSubmissionsByStatus: filter by status (new/read/replied/archived)
-- - ListContactSubmissionsByType: filter by submission_type
-- - ListContactSubmissionsByStatusAndType: combined filters
-- - SearchContactSubmissions: text search in name/email/company/message
--
-- Each List* query has a corresponding Count* query for pagination.
-- Navigation queries (GetPrevious/NextSubmissionID) enable prev/next buttons.
-- ====================================================================

-- name: ListContactSubmissions :many
-- sqlc annotation: :many returns paginated contact submissions
-- Purpose: Basic paginated list of all contact submissions
-- Parameters:
--   1. LIMIT (INTEGER): submissions per page
--   2. OFFSET (INTEGER): pagination offset
-- Return type: slice of contact_submissions (excludes ip_address/user_agent for security)
-- Note: Ordered by created_at DESC (newest first)
SELECT id, name, email, phone, company, inquiry_type, message, status, submission_type, created_at, updated_at
FROM contact_submissions
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
SELECT id, name, email, phone, company, inquiry_type, message, status, submission_type, created_at, updated_at
FROM contact_submissions
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListContactSubmissionsByStatus :many
SELECT id, name, email, phone, company, inquiry_type, message, status, submission_type, created_at, updated_at
FROM contact_submissions
WHERE status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListContactSubmissionsByType :many
SELECT id, name, email, phone, company, inquiry_type, message, status, submission_type, created_at, updated_at
FROM contact_submissions
WHERE submission_type = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: ListContactSubmissionsByStatusAndType :many
SELECT id, name, email, phone, company, inquiry_type, message, status, submission_type, created_at, updated_at
FROM contact_submissions
WHERE status = ? AND submission_type = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: SearchContactSubmissions :many
-- Purpose: Full-text search across multiple fields for admin search functionality
-- Parameters:
--   1-4. search_term (TEXT): same search term repeated 4 times for each field
--   5. LIMIT (INTEGER)
--   6. OFFSET (INTEGER)
-- WHERE clause: LIKE search in name OR email OR company OR message
-- Note: SQLite doesn't support full-text search here, so LIKE is used for simplicity
SELECT id, name, email, phone, company, inquiry_type, message, status, submission_type, created_at, updated_at
FROM contact_submissions
WHERE (name LIKE '%' || ? || '%' OR email LIKE '%' || ? || '%' OR company LIKE '%' || ? || '%' OR message LIKE '%' || ? || '%')
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountContactSubmissions :one
SELECT COUNT(*) FROM contact_submissions;

-- name: CountContactSubmissionsByStatus :one
SELECT COUNT(*) FROM contact_submissions WHERE status = ?;

-- name: CountContactSubmissionsByType :one
SELECT COUNT(*) FROM contact_submissions WHERE submission_type = ?;

-- name: CountContactSubmissionsByStatusAndType :one
SELECT COUNT(*) FROM contact_submissions WHERE status = ? AND submission_type = ?;

-- name: CountContactSubmissionsSearch :one
SELECT COUNT(*) FROM contact_submissions
WHERE (name LIKE '%' || ? || '%' OR email LIKE '%' || ? || '%' OR company LIKE '%' || ? || '%' OR message LIKE '%' || ? || '%');

-- name: GetContactSubmissionByID :one
SELECT id, name, email, phone, company, inquiry_type, message, ip_address, user_agent,
       status, notes, submission_type, created_at, updated_at
FROM contact_submissions
WHERE id = ?;

-- name: GetPreviousSubmissionID :one
-- Purpose: Gets ID of submission created AFTER current one (for "previous" navigation button)
-- Parameters:
--   1. current_id (INTEGER): current submission ID
-- Logic: Finds submission with created_at > current submission's created_at
--        Orders ASC and takes first = chronologically next newer submission
-- Subquery: Gets created_at of current submission for comparison
SELECT cs.id FROM contact_submissions cs WHERE cs.created_at > (SELECT cs2.created_at FROM contact_submissions cs2 WHERE cs2.id = ?) ORDER BY cs.created_at ASC LIMIT 1;

-- name: GetNextSubmissionID :one
-- Purpose: Gets ID of submission created BEFORE current one (for "next" navigation button)
-- Parameters:
--   1. current_id (INTEGER): current submission ID
-- Logic: Finds submission with created_at < current submission's created_at
--        Orders DESC and takes first = chronologically next older submission
SELECT cs.id FROM contact_submissions cs WHERE cs.created_at < (SELECT cs2.created_at FROM contact_submissions cs2 WHERE cs2.id = ?) ORDER BY cs.created_at DESC LIMIT 1;

-- name: UpdateContactSubmissionStatus :exec
UPDATE contact_submissions
SET status = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: BulkMarkContactSubmissionsRead :exec
UPDATE contact_submissions
SET status = 'read', updated_at = CURRENT_TIMESTAMP
WHERE status = 'new';

-- name: DeleteContactSubmission :exec
DELETE FROM contact_submissions WHERE id = ?;

-- ====================================================================
-- OFFICE LOCATIONS - ADMIN QUERIES
-- ====================================================================
-- Manages physical office location data displayed on contact page.
-- - is_primary: marks the main/headquarters office (only one should be primary)
-- - is_active: controls public visibility
-- - display_order: custom sort order for multiple locations
-- ====================================================================

-- name: ListAllOfficeLocations :many
SELECT id, name, address_line1, address_line2, city, state, postal_code, country,
       phone, email, is_primary, is_active, display_order, created_at, updated_at
FROM office_locations
ORDER BY display_order ASC, id ASC;

-- name: GetOfficeLocationByID :one
SELECT id, name, address_line1, address_line2, city, state, postal_code, country,
       phone, email, is_primary, is_active, display_order
FROM office_locations
WHERE id = ?;

-- name: CreateOfficeLocation :one
INSERT INTO office_locations (
    name, address_line1, address_line2, city, state, postal_code, country,
    phone, email, is_primary, is_active, display_order
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, created_at, updated_at;

-- name: UpdateOfficeLocation :exec
UPDATE office_locations
SET name = ?, address_line1 = ?, address_line2 = ?, city = ?, state = ?,
    postal_code = ?, country = ?, phone = ?, email = ?, is_primary = ?,
    is_active = ?, display_order = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteOfficeLocation :exec
DELETE FROM office_locations WHERE id = ?;

-- name: UnsetPrimaryOfficeLocations :exec
UPDATE office_locations SET is_primary = 0 WHERE is_primary = 1;
