-- ====================================================================
-- CONTACT SUBMISSIONS
-- ====================================================================

-- name: CreateContactSubmission :one
INSERT INTO contact_submissions (
    name, email, phone, company, inquiry_type, message, ip_address, user_agent
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, created_at;

-- name: GetActiveOfficeLocations :many
SELECT id, name, address_line1, address_line2, city, state, postal_code, country, phone, email, is_primary
FROM office_locations
WHERE is_active = 1
ORDER BY is_primary DESC, display_order ASC, id ASC;

-- ====================================================================
-- CONTACT - ADMIN
-- ====================================================================

-- name: ListContactSubmissions :many
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
SELECT cs.id FROM contact_submissions cs WHERE cs.created_at > (SELECT cs2.created_at FROM contact_submissions cs2 WHERE cs2.id = ?) ORDER BY cs.created_at ASC LIMIT 1;

-- name: GetNextSubmissionID :one
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
-- OFFICE LOCATIONS - ADMIN
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
