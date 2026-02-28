# Test Plan: Case Study Metrics HTMX Management

## Summary
Testing HTMX-driven addition and deletion of metrics within case study edit page without full page reloads.

## KNOWN BUGS
- **Hardcoded case study ID in template**: Template has hardcoded case study ID `0` in delete URL (same bug as plan 25). Should be dynamic `{{.CaseStudyID}}`.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with case studies
- On case study edit page at /admin/case-studies/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/case-studies/:id/edit
2. Locate #metrics-list container
3. Fill add metrics form: metric_value (required, e.g. "40%"), metric_label (required, e.g. "Cost Reduction"), display_order
4. Click add button with hx-post="/admin/case-studies/:id/metrics"
5. Verify new metric appears in updated list (AddMetric returns partial HTML)
6. Click delete button on existing metric with hx-delete="/admin/case-studies/:id/metrics/:metricId"
7. Verify metric removed (DeleteMetric returns 204 No Content)

## Test Cases

### Happy Path
- **Add metric**: Fill metric_value "40%", metric_label "Cost Reduction", display_order 1, submit, metric added (AddMetric returns partial HTML)
- **Multiple metrics**: Add 3 metrics with different display_order values, verify all appear
- **Display order**: Metrics display in order based on display_order field
- **Various value formats**: metric_value "50%", "2x faster", "$1M saved", "99.9% uptime" all accepted
- **Delete metric**: Click delete button on metric, hx-delete removes it (returns 204 No Content, no HTML)
- **Individual item removal**: Delete uses hx-target="closest div" hx-swap="outerHTML" to remove individual metric item
- **No page reload**: All operations via HTMX, no full page refresh

### Edge Cases / Error States
- **Missing required metric_value**: Empty metric_value field triggers validation error via HTMX response
- **Missing required metric_label**: Empty metric_label field triggers validation error
- **Very long metric_value**: metric_value "999,999,999.99%" accepted and displayed
- **Very long metric_label**: metric_label with 100+ characters accepted, may wrap or truncate
- **Duplicate display_order**: Multiple metrics with same display_order accepted, sorted arbitrarily
- **Special characters**: metric_value with symbols "%", "$", "+", "x" accepted
- **Delete with confirmation**: hx-confirm attribute may prompt user before deletion
- **HTMX error handling**: Server error on add/delete shows error message in section

## Selectors & Elements
- Section container: id="metrics-list"
- Add form action: hx-post="/admin/case-studies/:id/metrics"
- Input names: metric_value (required), metric_label (required), display_order (number)
- Delete button: hx-delete="/admin/case-studies/:id/metrics/:metricId" hx-target="closest div" hx-swap="outerHTML"
- Add button: text "Add Metric"

## HTMX Interactions
- **Add metric**: hx-post="/admin/case-studies/:id/metrics" (returns partial HTML for the new metric item)
- **Delete metric**: hx-delete="/admin/case-studies/:id/metrics/:metricId" hx-target="closest div" hx-swap="outerHTML" (returns 204 No Content, no HTML)
- Delete removes individual metric items using closest div targeting

## Dependencies
- 24-case-studies-crud.md (parent case study edit page)
