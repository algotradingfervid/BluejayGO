# Test Plan: Case Study Metrics HTMX Management

## Summary
Testing HTMX-driven addition and deletion of metrics within case study edit page without full page reloads.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with case studies
- On case study edit page at /admin/case-studies/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/case-studies/:id/edit
2. Locate #metrics-section container
3. Fill add metrics form: metric_value (required, e.g. "40%"), metric_label (required, e.g. "Cost Reduction"), display_order
4. Click add button with hx-post="/admin/case-studies/:id/metrics"
5. Verify hx-target="#metrics-section" hx-swap="outerHTML" replaces entire section
6. Verify new metric appears in updated section
7. Click delete button on existing metric with hx-delete="/admin/case-studies/:id/metrics/:metricId"
8. Verify metric removed from #metrics-section without page reload

## Test Cases

### Happy Path
- **Add metric**: Fill metric_value "40%", metric_label "Cost Reduction", display_order 1, submit, metric added to #metrics-section
- **Multiple metrics**: Add 3 metrics with different display_order values, verify all appear
- **Display order**: Metrics display in order based on display_order field
- **Various value formats**: metric_value "50%", "2x faster", "$1M saved", "99.9% uptime" all accepted
- **Delete metric**: Click delete button on metric, hx-delete removes it from section
- **Section swap**: After add/delete, entire #metrics-section swapped with updated HTML
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
- Section container: id="metrics-section"
- Add form action: hx-post="/admin/case-studies/:id/metrics" hx-target="#metrics-section" hx-swap="outerHTML"
- Input names: metric_value (required), metric_label (required), display_order (number)
- Delete button: hx-delete="/admin/case-studies/:id/metrics/:metricId" hx-target="#metrics-section" hx-swap="outerHTML"
- Add button: text "Add Metric"

## HTMX Interactions
- **Add metric**: hx-post="/admin/case-studies/:id/metrics" hx-target="#metrics-section" hx-swap="outerHTML" (returns updated metrics_section.html partial)
- **Delete metric**: hx-delete="/admin/case-studies/:id/metrics/:metricId" hx-target="#metrics-section" hx-swap="outerHTML" (returns updated metrics_section.html partial)
- Both operations swap entire #metrics-section to reflect current state

## Dependencies
- 24-case-studies-crud.md (parent case study edit page)
