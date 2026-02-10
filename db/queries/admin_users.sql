-- name: GetAdminUserByEmail :one
SELECT id, email, password_hash, display_name, role, is_active, last_login_at
FROM admin_users
WHERE email = ? AND is_active = 1
LIMIT 1;

-- name: UpdateLastLogin :exec
UPDATE admin_users
SET last_login_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: CreateAdminUser :one
INSERT INTO admin_users (email, password_hash, display_name, role)
VALUES (?, ?, ?, ?)
RETURNING id, email, display_name, role, created_at;

-- name: ListAdminUsers :many
SELECT id, email, display_name, role, is_active, created_at, last_login_at
FROM admin_users
ORDER BY created_at DESC;
