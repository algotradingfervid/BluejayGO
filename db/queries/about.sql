-- name: GetCompanyOverview :one
SELECT * FROM company_overview ORDER BY id DESC LIMIT 1;

-- name: UpsertCompanyOverview :one
INSERT INTO company_overview (
    headline, tagline, description_main, description_secondary,
    description_tertiary, hero_image_url, company_image_url
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetMissionVisionValues :one
SELECT * FROM mission_vision_values ORDER BY id DESC LIMIT 1;

-- name: UpsertMissionVisionValues :one
INSERT INTO mission_vision_values (
    mission, vision, values_summary,
    mission_icon, vision_icon, values_icon
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListCoreValues :many
SELECT * FROM core_values ORDER BY display_order ASC;

-- name: GetCoreValue :one
SELECT * FROM core_values WHERE id = ? LIMIT 1;

-- name: CreateCoreValue :one
INSERT INTO core_values (title, description, icon, display_order)
VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdateCoreValue :one
UPDATE core_values SET title = ?, description = ?, icon = ?, display_order = ?
WHERE id = ? RETURNING *;

-- name: DeleteCoreValue :exec
DELETE FROM core_values WHERE id = ?;

-- name: ListMilestones :many
SELECT * FROM milestones ORDER BY display_order ASC;

-- name: GetMilestone :one
SELECT * FROM milestones WHERE id = ? LIMIT 1;

-- name: CreateMilestone :one
INSERT INTO milestones (year, title, description, is_current, display_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateMilestone :one
UPDATE milestones SET year = ?, title = ?, description = ?, is_current = ?, display_order = ?
WHERE id = ? RETURNING *;

-- name: DeleteMilestone :exec
DELETE FROM milestones WHERE id = ?;

-- name: ListCertifications :many
SELECT * FROM certifications ORDER BY display_order ASC;

-- name: GetCertification :one
SELECT * FROM certifications WHERE id = ? LIMIT 1;

-- name: CreateCertification :one
INSERT INTO certifications (name, abbreviation, description, icon, display_order)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateCertification :one
UPDATE certifications SET name = ?, abbreviation = ?, description = ?, icon = ?, display_order = ?
WHERE id = ? RETURNING *;

-- name: DeleteCertification :exec
DELETE FROM certifications WHERE id = ?;
