-- ====================================================================
-- ADMIN USERS QUERIES
-- ====================================================================
-- This file manages authentication and admin user account operations.
-- These queries handle login verification, user creation, and admin
-- user management within the CMS.
--
-- Managed entity:
-- - admin_users: CMS admin accounts with roles and permissions
--
-- Security notes:
-- - password_hash is stored, never plain text passwords
-- - is_active flag allows soft-disable of accounts
-- - last_login_at tracks account activity
-- ====================================================================

-- name: GetAdminUserByEmail :one
-- sqlc annotation: :one returns single admin_users row or error
-- Purpose: Retrieves admin user by email for authentication during login
-- Parameters:
--   1. email (TEXT): email address to look up
-- Return type: partial admin_users row (excludes sensitive created_at/updated_at)
-- WHERE clause notes:
--   - email = ? ensures exact match (case-sensitive in SQLite)
--   - is_active = 1 prevents disabled accounts from logging in
-- Security: Only returns active accounts; password_hash included for verification
SELECT id, email, password_hash, display_name, role, is_active, last_login_at
FROM admin_users
WHERE email = ? AND is_active = 1
LIMIT 1;

-- name: UpdateLastLogin :exec
-- sqlc annotation: :exec returns no data, only error or success
-- Purpose: Updates last_login_at timestamp after successful authentication
-- Parameters:
--   1. id (INTEGER): admin user ID to update
-- Return type: none
-- Note: CURRENT_TIMESTAMP uses database server time (UTC in SQLite)
--       updated_at also refreshed to track any account modifications
UPDATE admin_users
SET last_login_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: CreateAdminUser :one
-- sqlc annotation: :one returns the newly created user (partial)
-- Purpose: Creates a new admin user account
-- Parameters (4 positional):
--   1. email (TEXT): unique email address for login
--   2. password_hash (TEXT): bcrypt/argon2 hashed password
--   3. display_name (TEXT): human-friendly name for UI
--   4. role (TEXT): permission level (e.g., "admin", "editor", "viewer")
-- Return type: partial row (id, email, display_name, role, created_at)
-- Note: Does NOT return password_hash for security
--       is_active defaults to 1 (true) via schema default
--       UNIQUE constraint on email enforced at database level
INSERT INTO admin_users (email, password_hash, display_name, role)
VALUES (?, ?, ?, ?)
RETURNING id, email, display_name, role, created_at;

-- name: ListAdminUsers :many
-- sqlc annotation: :many returns slice of admin_users rows
-- Purpose: Lists all admin users for management dashboard
-- Parameters: none
-- Return type: slice of partial admin_users rows (excludes password_hash)
-- Note: Ordered by created_at DESC to show newest accounts first
--       Does NOT return password_hash for security
SELECT id, email, display_name, role, is_active, created_at, last_login_at
FROM admin_users
ORDER BY created_at DESC;
