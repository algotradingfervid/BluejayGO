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
2. Locate #stats-list container
3. Fill add stats form: value (required), label (required), display_order
4. Click add button with hx-post="/admin/solutions/:id/stats"
5. Verify hx-target="#stats-list" hx-swap="innerHTML" adds new stat to list
6. Verify new stat appears in list
7. Click delete button on existing stat with hx-delete="/admin/solutions/:id/stats/:statId"
8. Verify stat removed (hx-target="closest div" hx-swap="outerHTML swap:0.3s") without page reload

## Test Cases

### Happy Path
- **Add stat**: Fill value "99%", label "Uptime", display_order 1, submit, stat added to #stats-list (innerHTML swap adds partial)
- **Multiple stats**: Add 3 stats with different display_order values, verify all appear
- **Display order**: Stats display in order based on display_order field
- **Delete stat**: Click delete button on stat, hx-delete removes individual div (returns 200 NoContent with no HTML)
- **Add returns HTML**: Add handler returns solution_stats.html partial
- **Delete returns no content**: Delete handler returns c.NoContent(http.StatusOK) with no HTML
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
- Section container: id="stats-list"
- Add form action: hx-post="/admin/solutions/:id/stats" hx-target="#stats-list" hx-swap="innerHTML"
- Input names: value (required), label (required), display_order (number)
- Delete button: hx-delete="/admin/solutions/:id/stats/:statId" hx-target="closest div" hx-swap="outerHTML swap:0.3s"
- Add button: text "Add Stat"

## HTMX Interactions
- **Add stat**: hx-post="/admin/solutions/:id/stats" hx-target="#stats-list" hx-swap="innerHTML" (returns solution_stats.html partial)
- **Delete stat**: hx-delete="/admin/solutions/:id/stats/:statId" hx-target="closest div" hx-swap="outerHTML swap:0.3s" (returns c.NoContent(http.StatusOK) with no HTML content)
- Add operation inserts new stat HTML into list, delete operation removes individual item div

## Dependencies
- 19-solutions-crud.md (parent solution edit page)
