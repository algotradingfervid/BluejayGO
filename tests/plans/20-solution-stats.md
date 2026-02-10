# Test Plan: Solution Stats HTMX Management

## Summary
Testing HTMX-driven addition and deletion of stats within solution edit page without full page reloads.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with solutions
- On solution edit page at /admin/solutions/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/solutions/:id/edit
2. Locate #stats-section container
3. Fill add stats form: value (required), label (required), display_order
4. Click add button with hx-post="/admin/solutions/:id/stats"
5. Verify hx-target="#stats-section" hx-swap="outerHTML" replaces entire section
6. Verify new stat appears in updated section
7. Click delete button on existing stat with hx-delete="/admin/solutions/:id/stats/:statId"
8. Verify stat removed from #stats-section without page reload

## Test Cases

### Happy Path
- **Add stat**: Fill value "99%", label "Uptime", display_order 1, submit, stat added to #stats-section
- **Multiple stats**: Add 3 stats with different display_order values, verify all appear
- **Display order**: Stats display in order based on display_order field
- **Delete stat**: Click delete button on stat, hx-delete removes it from section
- **Section swap**: After add/delete, entire #stats-section swapped with updated HTML
- **No page reload**: All operations via HTMX, no full page refresh

### Edge Cases / Error States
- **Missing required value**: Empty value field triggers validation error via HTMX response
- **Missing required label**: Empty label field triggers validation error
- **Duplicate display_order**: Multiple stats with same display_order accepted, sorted arbitrarily
- **Very long value**: Value "999,999,999+" accepted and displayed
- **Very long label**: Label with 100+ characters accepted, may wrap or truncate
- **Delete with confirmation**: hx-confirm attribute prompts user before deletion
- **HTMX error handling**: Server error on add/delete shows error message in section

## Selectors & Elements
- Section container: id="stats-section"
- Add form action: hx-post="/admin/solutions/:id/stats" hx-target="#stats-section" hx-swap="outerHTML"
- Input names: value (required), label (required), display_order (number)
- Delete button: hx-delete="/admin/solutions/:id/stats/:statId" hx-target="#stats-section" hx-swap="outerHTML"
- Add button: text "Add Stat"

## HTMX Interactions
- **Add stat**: hx-post="/admin/solutions/:id/stats" hx-target="#stats-section" hx-swap="outerHTML" (returns updated stats_section.html partial)
- **Delete stat**: hx-delete="/admin/solutions/:id/stats/:statId" hx-target="#stats-section" hx-swap="outerHTML" (returns updated stats_section.html partial)
- Both operations swap entire #stats-section to reflect current state

## Dependencies
- 19-solutions-crud.md (parent solution edit page)
