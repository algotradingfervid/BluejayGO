# Phase 19 - Navigation Editor

## Current State
- No navigation editor exists in the current admin panel
- Nav links are controlled by toggle switches in settings

## Goal
Visual navigation editor for managing site menu structure.

## Navigation Editor Page

### Two-Column Layout

**Left Column (1/3): Add Items**

Menu Settings:
- Menu Name
  - Tooltip: "Internal name for this menu (e.g., 'Main Header Nav')."
- Location dropdown: Header / Footer / Footer Legal / Mobile
  - Tooltip: "Where this menu appears on the site."

Add Item Options:
- Page Link: dropdown of existing pages (Products, Solutions, About, Blog, etc.)
  - Tooltip: "Link to an existing page on your site."
- Custom Link: Label + URL inputs
  - Tooltip: "Link to any URL, internal or external."
- Dropdown: Name input (creates a parent with child items)
  - Tooltip: "Creates a dropdown menu with child links."
- "Add to Menu" button for each type

**Right Column (2/3): Menu Structure**

- Draggable list of menu items
- Each item shows:
  - Drag handle (6 dots icon)
  - Type badge: Page (blue), Custom (green), Dropdown (purple)
  - Label text
  - URL (truncated)
  - Edit / Delete buttons
- Nested items (children of dropdowns) indented with left border
- Drag to reorder or nest under dropdowns

### Edit Item Modal
- Label
  - Tooltip: "Display text for this menu link."
- Link Type: Page / Custom / Dropdown (visual selector cards)
- If Page: page dropdown
- If Custom: URL input
- Open in: Same Window / New Tab
  - Tooltip: "Choose 'New Tab' for external links."
- Active toggle
  - Tooltip: "Inactive items are hidden from the navigation."

### Preview Section
- Below the editor: live preview showing how the nav looks
- Updates as you make changes (HTMX partial refresh)

## Backend Implementation

### Database Tables
```sql
CREATE TABLE navigation_menus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    location TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE navigation_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    menu_id INTEGER NOT NULL REFERENCES navigation_menus(id),
    parent_id INTEGER REFERENCES navigation_items(id),
    label TEXT NOT NULL,
    link_type TEXT NOT NULL, -- 'page', 'custom', 'dropdown'
    url TEXT,
    page_identifier TEXT,
    open_new_tab INTEGER DEFAULT 0,
    is_active INTEGER DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Routes
- `GET /admin/navigation` - List all menus
- `GET /admin/navigation/:id` - Edit specific menu
- `POST /admin/navigation` - Create menu
- `POST /admin/navigation/:id/items` - Add item
- `PUT /admin/navigation/items/:id` - Update item
- `DELETE /admin/navigation/items/:id` - Delete item
- `POST /admin/navigation/:id/reorder` - Update sort order (HTMX)

## Files to Create/Modify
| File | Action |
|------|--------|
| `templates/admin/pages/navigation_list.html` | Create |
| `templates/admin/pages/navigation_editor.html` | Create |
| `internal/handlers/admin/navigation.go` | Create |
| `db/migrations/031_navigation.up.sql` | Create |
| `db/queries/navigation.sql` | Create |
| `cmd/server/main.go` | Add routes |

## Dependencies
- Phase 01, 02
- Phase 05 (header) references navigation structure
