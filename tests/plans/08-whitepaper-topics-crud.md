# Test Plan: Whitepaper Topics CRUD

## Summary
Verify complete CRUD operations for whitepaper topics including color codes, icons, and UNIQUE name validation.

## Preconditions
- User authenticated with valid session cookie
- Database seeded with 6 whitepaper topics
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/whitepaper-topics
2. Verify list shows all 6 seeded topics with color indicators and icons
3. Click "New Topic" or navigate to http://localhost:28090/admin/whitepaper-topics/new
4. Fill form: name, color_hex, icon, description, sort_order
5. Submit POST to /admin/whitepaper-topics
6. Verify redirect to list with new topic visible
7. Click "Edit" on a topic or navigate to http://localhost:28090/admin/whitepaper-topics/:id/edit
8. Modify fields and submit POST to /admin/whitepaper-topics/:id
9. Verify updated data appears in list
10. Click delete button with hx-delete attribute
11. Confirm deletion in browser confirmation dialog
12. Verify topic row removed from table via HTMX

## Test Cases

### Happy Path
- **List all topics**: All 6 seeded topics display with names, colors, and icons
- **Create new topic**: Form submission creates topic with auto-generated slug
- **Edit existing topic**: Form pre-fills with current data, updates successfully
- **Delete topic via HTMX**: Delete button removes row without page reload
- **Auto-slug generation**: Slug automatically generated from name field
- **Color preview**: Color hex value displays as visual indicator in list and form
- **Icon display**: Icon value displays in list view

### Edge Cases / Error States
- **Duplicate name validation**: Creating topic with existing name shows UNIQUE constraint error
- **Required name field**: Submitting without name shows validation error
- **Invalid color hex format**: Entering invalid hex code shows validation error
- **Delete confirmation cancel**: Canceling hx-confirm dialog does not delete topic
- **Delete in-use topic**: Deleting topic assigned to whitepapers prevents deletion or shows error
- **Empty icon field**: Topic can be created without icon (optional field)
- **Empty color hex**: System handles missing color or uses default value

## Selectors & Elements
- List page: http://localhost:28090/admin/whitepaper-topics
- Create form: `action="/admin/whitepaper-topics" method="POST"`
- Edit form: `action="/admin/whitepaper-topics/:id" method="POST"`
- Input name: `name="name" type="text"` (required, UNIQUE)
- Input slug: `name="slug" type="text"` (auto-generated)
- Input color_hex: `name="color_hex" type="text"` or `type="color"`
- Input icon: `name="icon" type="text"`
- Textarea description: `name="description"`
- Input sort_order: `name="sort_order" type="number"`
- Delete button: `hx-delete="/admin/whitepaper-topics/:id" hx-confirm="Delete this topic?"`
- Color indicator: visual element showing color_hex in list view
- Icon indicator: visual element showing icon in list view

## HTMX Interactions
- Delete action: `hx-delete="/admin/whitepaper-topics/:id"` with `hx-confirm` dialog
- Target: `hx-target="closest tr"` removes table row on successful delete
- Swap: `hx-swap="outerHTML"` replaces row element
- No HTMX on create/edit forms (standard POST with redirect)

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- Whitepapers CRUD test plan (topics used in whitepaper categorization - not in current set)
