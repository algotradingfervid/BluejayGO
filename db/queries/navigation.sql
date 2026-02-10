-- ====================================================================
-- NAVIGATION QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing navigation menus and menu items.
--
-- Entities:
--   - navigation_menus: Container for menu items (e.g., "Main Menu", "Footer Menu")
--   - navigation_items: Individual links/items within a menu (supports nesting)
--
-- Features:
--   - Multiple menus per site (header, footer, sidebar, etc.)
--   - Hierarchical navigation (parent/child items for dropdowns)
--   - Flexible link types (URL, page identifier, custom)
--   - Drag-and-drop reordering via sort_order
-- ====================================================================

-- ====================================================================
-- NAVIGATION MENUS
-- ====================================================================

-- name: ListNavigationMenus :many
-- Retrieves all navigation menus sorted by creation date (newest first).
--
-- Parameters: none
-- Returns: []NavigationMenu - Array of all menu containers
--
-- Use case: Admin panel menu management, displaying available menus
SELECT * FROM navigation_menus ORDER BY created_at DESC;

-- name: GetNavigationMenu :one
-- Retrieves a single navigation menu by its primary key ID.
--
-- Parameters:
--   $1 (INTEGER) - menu ID
-- Returns: NavigationMenu - Single menu record or error if not found
--
-- Use case: Editing a specific menu, fetching menu details
SELECT * FROM navigation_menus WHERE id = ? LIMIT 1;

-- name: CreateNavigationMenu :one
-- Creates a new navigation menu container.
--
-- Parameters:
--   $1 (TEXT) - name: Display name for admin (e.g., "Main Navigation")
--   $2 (TEXT) - location: Identifier for template placement (e.g., "header", "footer")
-- Returns: NavigationMenu - The newly created menu with auto-generated ID and timestamps
--
-- Use case: Setting up a new menu location (header, footer, sidebar, etc.)
INSERT INTO navigation_menus (name, location) VALUES (?, ?) RETURNING *;

-- name: UpdateNavigationMenu :exec
-- Updates an existing navigation menu's metadata.
--
-- Parameters:
--   $1 (TEXT) - name: Updated menu display name
--   $2 (TEXT) - location: Updated location identifier
--   $3 (INTEGER) - id: Menu ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
UPDATE navigation_menus SET name = ?, location = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: DeleteNavigationMenu :exec
-- Permanently deletes a navigation menu container.
--
-- Parameters:
--   $1 (INTEGER) - menu ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: This is a hard delete. Should cascade delete all navigation_items in this menu
-- Note: Ensure foreign key constraints are configured for cascading deletes
DELETE FROM navigation_menus WHERE id = ?;

-- ====================================================================
-- NAVIGATION ITEMS
-- ====================================================================

-- name: ListNavigationItems :many
-- Retrieves all navigation items for a specific menu, sorted by display order.
--
-- Parameters:
--   $1 (INTEGER) - menu_id: Parent menu ID
-- Returns: []NavigationItem - Array of menu items in display order
--
-- Sorting: sort_order ASC - Items appear in manually configured order
-- Use case: Rendering menu items in templates, admin item listing
-- Note: Returns both parent and child items; application must build hierarchy
SELECT * FROM navigation_items WHERE menu_id = ? ORDER BY sort_order ASC;

-- name: GetNavigationItem :one
-- Retrieves a single navigation item by its primary key ID.
--
-- Parameters:
--   $1 (INTEGER) - navigation item ID
-- Returns: NavigationItem - Single item record or error if not found
--
-- Use case: Editing a specific menu item
SELECT * FROM navigation_items WHERE id = ? LIMIT 1;

-- name: CreateNavigationItem :one
-- Creates a new navigation item (link) within a menu.
--
-- Parameters:
--   $1 (INTEGER) - menu_id: Parent menu container ID
--   $2 (INTEGER) - parent_id: Parent item ID for nested items (NULL for top-level)
--   $3 (TEXT) - label: Display text for the link
--   $4 (TEXT) - link_type: Type of link ("url", "page", "custom")
--   $5 (TEXT) - url: External or custom URL (NULL if using page_identifier)
--   $6 (TEXT) - page_identifier: Internal page slug/ID (NULL if using url)
--   $7 (BOOLEAN) - open_new_tab: Whether link opens in new window
--   $8 (BOOLEAN) - is_active: Whether item is visible (soft delete alternative)
--   $9 (INTEGER) - sort_order: Display position within menu
--
-- Returns: NavigationItem - The newly created item with auto-generated ID and timestamps
--
-- Link type logic:
--   - "url": Uses the url field for external/custom links
--   - "page": Uses page_identifier to reference internal content
--   - Application code resolves page_identifier to actual URL
INSERT INTO navigation_items (menu_id, parent_id, label, link_type, url, page_identifier, open_new_tab, is_active, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateNavigationItem :exec
-- Updates an existing navigation item's properties.
--
-- Parameters:
--   $1 (TEXT) - label: Updated link text
--   $2 (TEXT) - link_type: Updated link type
--   $3 (TEXT) - url: Updated URL
--   $4 (TEXT) - page_identifier: Updated page reference
--   $5 (BOOLEAN) - open_new_tab: Updated target behavior
--   $6 (BOOLEAN) - is_active: Updated visibility status
--   $7 (INTEGER) - parent_id: Updated parent (for moving in hierarchy)
--   $8 (INTEGER) - sort_order: Updated display position
--   $9 (INTEGER) - id: Item ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Editing menu items, changing hierarchy, reordering
UPDATE navigation_items
SET label = ?, link_type = ?, url = ?, page_identifier = ?, open_new_tab = ?, is_active = ?, parent_id = ?, sort_order = ?
WHERE id = ?;

-- name: DeleteNavigationItem :exec
-- Permanently deletes a single navigation item.
--
-- Parameters:
--   $1 (INTEGER) - navigation item ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: If this is a parent item, child items may become orphaned or cascade delete
-- Note: Consider setting is_active=0 for soft delete instead
DELETE FROM navigation_items WHERE id = ?;

-- name: DeleteNavigationItemsByMenu :exec
-- Permanently deletes all navigation items belonging to a specific menu.
--
-- Parameters:
--   $1 (INTEGER) - menu_id: Parent menu ID
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Clearing all items from a menu before rebuilding, cleanup after menu deletion
-- WARNING: This is a bulk hard delete operation
DELETE FROM navigation_items WHERE menu_id = ?;

-- name: UpdateNavigationItemOrder :exec
-- Updates a navigation item's sort order and parent (for drag-and-drop reordering).
--
-- Parameters:
--   $1 (INTEGER) - sort_order: New display position
--   $2 (INTEGER) - parent_id: New parent ID (NULL for top-level, or parent item ID)
--   $3 (INTEGER) - id: Item ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Drag-and-drop reordering in admin interface
-- Note: Application should handle recalculating sort_order for all affected items
UPDATE navigation_items SET sort_order = ?, parent_id = ? WHERE id = ?;
