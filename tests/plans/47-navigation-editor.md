# Test Plan: Navigation Menu Editor

## Summary
Tests the comprehensive navigation editor with menu creation, hierarchical item management, drag-drop reordering, and support for different link types (page/url/dropdown).

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Navigation system supports: header, footer, sidebar locations
- Link types: page (internal page reference), url (custom URL), dropdown (parent item)
- Hierarchical tree structure with parent_id relationships
- Drag-drop reordering with POST /admin/navigation/:id/reorder

## User Journey Steps
1. Navigate to /admin/navigation
2. View list of existing menus (may be empty initially)
3. Create new menu with name and location (header/footer/sidebar)
4. Navigate to menu editor GET /admin/navigation/:id
5. View hierarchical tree of menu items
6. Add new item: fill label, select link_type (page/url/dropdown)
7. If type=page: select page_identifier
8. If type=url: enter custom URL
9. If type=dropdown: item becomes parent for sub-items
10. Toggle open_new_tab checkbox
11. Toggle is_active checkbox
12. Submit to create item
13. Edit existing item: update label, link_type, url, page_identifier
14. Drag-drop items to reorder or change hierarchy
15. POST reorder data: [{id, parent_id, order}]
16. Delete menu item
17. Update menu settings (name, location)
18. Delete entire menu

## Test Cases

### Happy Path - Menu Management
- **List menus**: Verifies GET /admin/navigation shows all menus
- **Create new menu**: POSTs name="Main Menu" location="header", verifies creation
- **Create footer menu**: POSTs name="Footer Links" location="footer", verifies creation
- **Create sidebar menu**: POSTs name="Sidebar Nav" location="sidebar", verifies creation
- **Edit menu settings**: Updates menu name and location via POST /admin/navigation/:id/settings
- **Delete menu**: Deletes entire menu via DELETE /admin/navigation/:id, verifies removal

### Happy Path - Menu Items (Page Links)
- **Add page link**: Creates item with link_type=page, page_identifier="about", verifies creation
- **Edit page link**: Updates label and page_identifier, verifies save
- **Toggle open new tab**: Checks open_new_tab, verifies update
- **Toggle is_active**: Unchecks is_active, verifies item hidden

### Happy Path - Menu Items (Custom URLs)
- **Add URL link**: Creates item with link_type=url, url="https://example.com", verifies creation
- **Edit URL link**: Updates label and url, verifies save
- **External URL**: Adds external URL, checks open_new_tab by default
- **Internal URL**: Adds relative URL "/custom", verifies handling

### Happy Path - Menu Items (Dropdowns)
- **Add dropdown**: Creates item with link_type=dropdown, verifies parent item creation
- **Add child to dropdown**: Creates item with parent_id=dropdown_id, verifies hierarchy
- **Add nested dropdown**: Creates dropdown within dropdown, verifies multi-level nesting
- **Edit dropdown label**: Updates parent item label, verifies all children remain

### Happy Path - Reordering
- **Drag-drop reorder**: Changes item order, POSTs reorder JSON, verifies new order
- **Move item to dropdown**: Drags item into dropdown, POSTs new parent_id, verifies hierarchy
- **Move item out of dropdown**: Drags child to root level, POSTs parent_id=null, verifies
- **Reorder within dropdown**: Reorders children within parent, verifies order

### Happy Path - Deletion
- **Delete leaf item**: Deletes item with no children, verifies removal
- **Delete dropdown**: Deletes parent item, verifies handling of children (cascade or prevent)
- **Delete confirmation**: Verifies confirmation modal before deletion

### Edge Cases / Error States
- **Create menu empty name**: Tests validation when menu name is empty
- **Create menu no location**: Tests validation when location not selected
- **Duplicate menu names**: Creates menus with same name, verifies handling
- **Add item empty label**: Tests validation when label is empty
- **Add item no link type**: Tests validation when link_type not selected
- **Page link without identifier**: Selects link_type=page but no page_identifier, checks validation
- **URL link without URL**: Selects link_type=url but no url, checks validation
- **Invalid URL format**: Enters malformed URL, checks validation
- **Dropdown with external URL**: Tests if dropdown can have both children and URL
- **Very deep nesting**: Creates 5+ levels of nested dropdowns, verifies handling
- **Circular reference**: Attempts to make item its own parent, verifies prevention
- **Reorder invalid data**: POSTs malformed reorder JSON, verifies error handling
- **Reorder non-existent item**: Includes invalid item ID in reorder, verifies error
- **Delete item with children**: Deletes parent with children, verifies cascade or error
- **Update non-existent item**: POSTs to /admin/navigation/items/99999, verifies 404
- **Delete entire menu confirmation**: Verifies strong confirmation before menu deletion
- **Menu with saved=1**: Updates menu settings, verifies redirect with ?saved=1
- **Item position constraints**: Tests if items can be reordered across different menus
- **Maximum items per menu**: Tests performance/limits with 100+ items
- **Long label text**: Enters 200+ char label, checks truncation or limit
- **Special characters in label**: Enters HTML/emojis in label, verifies sanitization

## Selectors & Elements

### Menu List
- Menus list: `#navigation-menus` or `.menus-table`
- New menu button: `button#new-menu` or `a[href="/admin/navigation/new"]`
- Menu row: `.menu-row[data-id]` or `tr[data-menu-id]`
- Edit menu link: `a[href="/admin/navigation/{id}"]`
- Delete menu button: `button[data-delete-menu="{id}"]` or delete form

### Menu Editor
- Menu tree: `#menu-tree` or `.menu-items-tree`
- Add item button: `button#add-menu-item`
- Menu item node: `.menu-item[data-id]` (draggable)
- Drag handle: `.drag-handle`
- Item label: `.item-label`
- Edit item link: `button[data-edit-item="{id}"]`
- Delete item button: `button[data-delete-item="{id}"]`

### Menu Settings Form
- Form: `form[action="/admin/navigation/{id}/settings"][method="POST"]`
- Menu name: `input[name="name"]`
- Location select: `select[name="location"]`
- Location options: `option[value="header"]`, `option[value="footer"]`, `option[value="sidebar"]`
- Submit button: `button[type="submit"]`

### Add/Edit Item Form
- Form: `form[action="/admin/navigation/{menu_id}/items"][method="POST"]` (create)
- Form: `form[action="/admin/navigation/items/{id}"][method="POST"]` (update)
- Label input: `input[name="label"]`
- Link type radios: `input[name="link_type"][value="page"]`, `[value="url"]`, `[value="dropdown"]`
- Page identifier select: `select[name="page_identifier"]`
- Page options: `option[value="home"]`, `option[value="about"]`, `option[value="products"]`, etc.
- URL input: `input[name="url"][type="url"]`
- Open new tab checkbox: `input[name="open_new_tab"][type="checkbox"]`
- Is active checkbox: `input[name="is_active"][type="checkbox"]`
- Submit button: `button[type="submit"]`

### Reorder
- Reorder endpoint: POST /admin/navigation/:id/reorder
- Request body: JSON array `[{id: 1, parent_id: null, order: 0}, {id: 2, parent_id: 1, order: 1}, ...]`

### Common Elements
- Success message: `.alert-success`
- Error message: `.alert-error`
- Confirmation modal: `#confirm-modal` or `.confirm-dialog`

## HTMX Interactions
- May use HTMX for adding/deleting items without page reload
- Item forms may use hx-post for inline updates
- Delete may use hx-delete="/admin/navigation/items/:id"

## Dependencies
- Database tables: navigation_menus, navigation_items
- navigation_menus columns: id, name, location
- navigation_items columns: id, menu_id, parent_id, label, link_type, url, page_identifier, open_new_tab, is_active, order
- Drag-drop JavaScript library (e.g., Sortable.js) for reordering
- Template: templates/admin/pages/navigation-list.html, navigation-editor.html
- Handler: internal/handlers/navigation.go (ListMenus, CreateMenu, GetMenu, UpdateMenuSettings, DeleteMenu, CreateMenuItem, UpdateMenuItem, DeleteMenuItem, ReorderItems)
- Page identifier options defined (home, about, products, solutions, blog, contact, etc.)
- Hierarchical rendering for nested items
- Reorder handler validates parent_id relationships
