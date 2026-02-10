# Test Plan: Homepage Stats CRUD

## Summary
Tests the complete CRUD operations for homepage statistics display including stat values, labels, active status, and display ordering.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with 4 homepage stats
- Stats displayed in order by display_order field

## User Journey Steps
1. Navigate to /admin/homepage/stats
2. View list of existing stats (4 seeded items)
3. Click "New Stat" button
4. Fill in stat_value (required) - e.g., "500+", "99%"
5. Fill in stat_label (required) - e.g., "Projects Completed"
6. Set display_order
7. Check/uncheck is_active checkbox
8. Submit to create stat
9. Edit an existing stat
10. Update stat_value and stat_label
11. Delete a stat
12. Verify ordering and active status filtering

## Test Cases

### Happy Path
- **List stats**: Verifies GET /admin/homepage/stats shows all 4 seeded stats
- **Create new stat**: Adds stat with valid stat_value and stat_label, verifies creation
- **Edit existing stat**: Updates stat_value, verifies save
- **Update label**: Changes stat_label, verifies update
- **Toggle is_active**: Checks/unchecks is_active checkbox, verifies status change
- **Reorder stats**: Changes display_order values, verifies list reordering
- **Various value formats**: Tests different stat_value formats (numbers, percentages, ranges)

### Edge Cases / Error States
- **Empty stat_value**: Tests required field validation
- **Empty stat_label**: Tests required field validation
- **Very long stat_value**: Tests character limit on stat_value
- **Very long stat_label**: Tests character limit on stat_label
- **Special characters in value**: Enters "+", "%", "-" in stat_value, verifies handling
- **Duplicate display_order**: Creates stats with same order, verifies handling
- **Multiple active stats**: Ensures all active stats display correctly
- **All stats inactive**: Sets all is_active to false, checks empty state or warning
- **Delete stat**: Deletes stat, verifies removal from list
- **Delete confirmation**: Verifies confirmation before deletion
- **Numeric validation**: Tests if stat_value accepts/validates numeric patterns

## Selectors & Elements
- Stats list: `#stats-list` or `.stats-table`
- New stat button: `a[href="/admin/homepage/stats/new"]` or `button#new-stat`
- Stat row: `.stat-row[data-id]` or `tr[data-stat-id]`
- Edit link: `a[href="/admin/homepage/stats/{id}/edit"]`
- Delete button: `button[type="submit"]` in delete form or delete link
- Stat value input: `input[name="stat_value"]`
- Stat label input: `input[name="stat_label"]`
- Display order input: `input[name="display_order"][type="number"]`
- Is active checkbox: `input[name="is_active"][type="checkbox"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`
- Active badge: `.badge-active` or indicator for is_active=true

## HTMX Interactions
- None specified - standard form submissions with redirects
- Delete may use HTMX hx-delete if implemented

## Dependencies
- Database seeded with 4 homepage stats
- Template: templates/admin/pages/homepage-stats-list.html, homepage-stats-form.html
- Handler: internal/handlers/homepage.go (ListStats, NewStat, CreateStat, EditStat, UpdateStat, DeleteStat)
