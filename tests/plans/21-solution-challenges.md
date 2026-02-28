# Test Plan: Solution Challenges HTMX Management

## Summary
Testing HTMX-driven addition and deletion of challenges within solution edit page without full page reloads.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with solutions
- On solution edit page at /admin/solutions/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/solutions/:id/edit
2. Locate #challenges-list container
3. Fill add challenge form: title (required), description (textarea required), icon, display_order
4. Click add button with hx-post="/admin/solutions/:id/challenges"
5. Verify hx-target="#challenges-list" hx-swap="innerHTML" adds new challenge to list
6. Verify new challenge appears in list
7. Click delete button on existing challenge with hx-delete="/admin/solutions/:id/challenges/:challengeId"
8. Verify challenge removed (hx-target="closest div" hx-swap="outerHTML swap:0.3s") without page reload

## Test Cases

### Happy Path
- **Add challenge**: Fill title "Data Silos", description "Legacy systems...", icon "storage", display_order 1, submit, challenge added (innerHTML swap adds partial HTML)
- **Multiple challenges**: Add 3 challenges with different display_order values, verify all appear
- **Display order**: Challenges display in order based on display_order field
- **Delete challenge**: Click delete button on challenge, hx-delete removes individual div (returns c.NoContent(http.StatusOK) with no HTML)
- **Add returns HTML**: Add handler returns partial HTML
- **Delete returns no content**: DeleteChallenge returns c.NoContent(http.StatusOK) with no HTML
- **No page reload**: All operations via HTMX, no full page refresh
- **Material icon**: Icon field accepts Material icon names like "warning", "error_outline"

### Edge Cases / Error States
- **Missing required title**: Empty title field triggers validation error via HTMX response
- **Missing required description**: Empty description textarea triggers validation error
- **Long description**: Description with 500+ characters accepted in textarea
- **Optional icon**: Leaving icon field empty accepted, no icon displayed
- **Duplicate display_order**: Multiple challenges with same display_order accepted, sorted arbitrarily
- **Delete with confirmation**: hx-confirm attribute prompts user before deletion
- **HTMX error handling**: Server error on add/delete shows error message in section

## Selectors & Elements
- Section container: id="challenges-list"
- Add form action: hx-post="/admin/solutions/:id/challenges" hx-target="#challenges-list" hx-swap="innerHTML"
- Input names: title (required), description (textarea required), icon, display_order (number)
- Delete button: hx-delete="/admin/solutions/:id/challenges/:challengeId" hx-target="closest div" hx-swap="outerHTML swap:0.3s"
- Add button: text "Add Challenge"

## HTMX Interactions
- **Add challenge**: hx-post="/admin/solutions/:id/challenges" hx-target="#challenges-list" hx-swap="innerHTML" (returns partial HTML)
- **Delete challenge**: hx-delete="/admin/solutions/:id/challenges/:challengeId" hx-target="closest div" hx-swap="outerHTML swap:0.3s" (returns c.NoContent(http.StatusOK) with no HTML content)
- Add operation inserts new challenge HTML into list, delete operation removes individual item div

## Dependencies
- 19-solutions-crud.md (parent solution edit page)
