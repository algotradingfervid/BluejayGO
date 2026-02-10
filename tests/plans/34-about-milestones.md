# Test Plan: Company Milestones CRUD

## Summary
Tests the complete CRUD operations for company milestones timeline including year, title, description, current status, and ordering.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with milestones from 2008 to present
- Milestones displayed in chronological or reverse chronological order

## User Journey Steps
1. Navigate to /admin/about/milestones
2. View list of existing milestones (seeded data 2008-present)
3. Click "New Milestone" button
4. Fill in year (number), title, description
5. Check/uncheck "is_current" checkbox
6. Set display_order
7. Submit to create milestone
8. Edit an existing milestone
9. Update fields including toggling is_current
10. Delete a milestone using HTMX delete button
11. Verify milestone removed without page reload

## Test Cases

### Happy Path
- **List milestones**: Verifies GET /admin/about/milestones shows all seeded milestones
- **Create new milestone**: Adds milestone for current year with valid data
- **Edit existing milestone**: Updates title and description, verifies save
- **Toggle is_current**: Checks/unchecks is_current checkbox, verifies update
- **Set display order**: Changes display_order, verifies list reordering
- **Year validation**: Enters valid 4-digit year, verifies acceptance
- **Multiple current milestones**: Verifies system allows or prevents multiple is_current=true

### Edge Cases / Error States
- **Invalid year format**: Enters 3-digit or 5-digit year, tests validation
- **Future year**: Enters year in future, checks if allowed
- **Historical year**: Enters very old year (e.g., 1900), checks validation
- **Empty title**: Tests required field validation on title
- **Empty description**: Tests required field validation on description
- **Delete via HTMX**: Clicks delete button, verifies hx-delete removes item
- **Delete confirmation**: Verifies confirmation before deletion
- **Duplicate years**: Creates multiple milestones for same year, verifies handling
- **Very long description**: Tests textarea character limits
- **Cache invalidation**: Confirms page:about cache cleared after changes

## Selectors & Elements
- Milestones list: `#milestones-list` or `.milestones-table`
- New milestone button: `a[href="/admin/about/milestones/new"]` or `button#new-milestone`
- Milestone row: `.milestone-row[data-id]` or `tr[data-milestone-id]`
- Delete button: `button[hx-delete="/admin/about/milestones/{id}"]`
- Year input: `input[name="year"][type="number"]`
- Title input: `input[name="title"]`
- Description textarea: `textarea[name="description"]`
- Is current checkbox: `input[name="is_current"][type="checkbox"]`
- Display order input: `input[name="display_order"][type="number"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`
- Current badge: `.badge-current` or indicator for is_current=true

## HTMX Interactions
- **hx-delete**: Delete button uses `hx-delete="/admin/about/milestones/{id}"`
- **hx-target**: Targets parent row or container for removal
- **hx-swap**: Uses `outerHTML` or `delete` to remove element
- **hx-confirm**: Confirmation message before delete action

## Dependencies
- Database seeded with milestones (2008-present)
- Cache service for page:about
- Template: templates/admin/pages/about-milestones-list.html, about-milestones-form.html
- Handler: internal/handlers/about.go (ListMilestones, NewMilestone, CreateMilestone, EditMilestone, UpdateMilestone, DeleteMilestone)
- HTMX library loaded
