-- name: ListNavigationMenus :many
SELECT * FROM navigation_menus ORDER BY created_at DESC;

-- name: GetNavigationMenu :one
SELECT * FROM navigation_menus WHERE id = ? LIMIT 1;

-- name: CreateNavigationMenu :one
INSERT INTO navigation_menus (name, location) VALUES (?, ?) RETURNING *;

-- name: UpdateNavigationMenu :exec
UPDATE navigation_menus SET name = ?, location = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: DeleteNavigationMenu :exec
DELETE FROM navigation_menus WHERE id = ?;

-- name: ListNavigationItems :many
SELECT * FROM navigation_items WHERE menu_id = ? ORDER BY sort_order ASC;

-- name: GetNavigationItem :one
SELECT * FROM navigation_items WHERE id = ? LIMIT 1;

-- name: CreateNavigationItem :one
INSERT INTO navigation_items (menu_id, parent_id, label, link_type, url, page_identifier, open_new_tab, is_active, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateNavigationItem :exec
UPDATE navigation_items
SET label = ?, link_type = ?, url = ?, page_identifier = ?, open_new_tab = ?, is_active = ?, parent_id = ?, sort_order = ?
WHERE id = ?;

-- name: DeleteNavigationItem :exec
DELETE FROM navigation_items WHERE id = ?;

-- name: DeleteNavigationItemsByMenu :exec
DELETE FROM navigation_items WHERE menu_id = ?;

-- name: UpdateNavigationItemOrder :exec
UPDATE navigation_items SET sort_order = ?, parent_id = ? WHERE id = ?;
