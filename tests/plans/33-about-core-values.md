# Test Plan: Core Values CRUD

## Summary
Tests the complete CRUD operations for company core values including creation, editing, deletion, and display order management.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with 6 core values
- Material Icons available for icon field

## User Journey Steps
1. Navigate to /admin/about/values
2. View list of existing core values (6 seeded items)
3. Click "New Value" button
4. Fill in title, description, icon, and display_order
5. Submit to create new value
6. Edit an existing value
7. Update fields and save
8. Delete a value using HTMX delete button
9. Verify value removed from list without full page reload
10. Verify display order affects list sorting

## Test Cases

### Happy Path
- **List core values**: Verifies GET /admin/about/values shows all 6 seeded values
- **Create new value**: Fills form with valid data, verifies successful creation
- **Edit existing value**: Updates title and description, verifies save
- **Update icon**: Changes icon to different Material icon name, verifies update
- **Reorder values**: Changes display_order values, verifies list reordering
- **View empty state**: Deletes all values, verifies empty state message

### Edge Cases / Error States
- **Create with empty title**: Tests required field validation
- **Create with empty description**: Tests required field validation
- **Duplicate display_order**: Creates values with same order, verifies handling
- **Delete value via HTMX**: Clicks delete button, verifies hx-delete removes item from DOM
- **Invalid icon name**: Enters non-existent Material icon, checks validation
- **Very long title**: Tests character limit on title field
- **Very long description**: Tests textarea limit
- **Delete confirmation**: Verifies confirmation modal/prompt before deletion
- **Cache invalidation**: Confirms page:about cache cleared after create/update/delete

## Selectors & Elements
- Values list: `#values-list` or `.values-table`
- New value button: `a[href="/admin/about/values/new"]` or `button#new-value`
- Value row: `.value-row[data-id]` or `tr[data-value-id]`
- Delete button: `button[hx-delete="/admin/about/values/{id}"]`
- Title input: `input[name="title"]`
- Description textarea: `textarea[name="description"]`
- Icon input: `input[name="icon"]`
- Display order input: `input[name="display_order"][type="number"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`

## HTMX Interactions
- **hx-delete**: Delete button uses `hx-delete="/admin/about/values/{id}"` to remove value
- **hx-target**: Delete targets parent row for removal
- **hx-swap**: Uses `outerHTML` or similar to remove element
- **hx-confirm**: May include confirmation message before delete

## Dependencies
- Database seeded with 6 core values
- Cache service for page:about
- Material Icons for icon display
- Template: templates/admin/pages/about-values-list.html, about-values-form.html
- Handler: internal/handlers/about.go (ListValues, NewValue, CreateValue, EditValue, UpdateValue, DeleteValue)
- HTMX library loaded
