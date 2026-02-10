# Test Plan: Solutions CRUD

## Summary
Testing solution creation, editing, listing with filters, deletion, and HTMX-driven sub-resource sections for stats, challenges, products, and CTAs.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 6 solutions

## User Journey Steps
1. Navigate to http://localhost:28090/admin/solutions
2. Verify solution list displays with search, status, page filters
3. Verify pagination shows 15 solutions per page
4. Click "New Solution" button to navigate to /admin/solutions/new
5. Fill required field: title
6. Auto-generate slug on title blur
7. Fill optional fields: icon (Material icon name), short_description, hero fields, overview_content, meta_description, reference_code
8. Set is_published checkbox and display_order number
9. Submit form via POST /admin/solutions
10. Navigate to edit form at /admin/solutions/:id/edit
11. Interact with HTMX sections: #stats-section, #challenges-section, #products-section, #ctas-section
12. Delete solution using hx-delete with hx-target="closest tr"
13. Verify cache invalidation for page:solutions

## Test Cases

### Happy Path
- **List solutions with filters**: Verify search by title, filter by status/page work correctly
- **Pagination**: Verify 15 solutions per page and navigation works
- **Create new solution**: Required title filled, slug auto-generated, solution saved
- **Auto-slug generation**: Title "IoT Platform" generates slug "iot-platform" on blur
- **Material icon field**: Enter "cloud_upload" as icon name, saved successfully
- **Hero section fields**: Fill hero_title, hero_description, hero_image_url
- **Overview content**: Enter rich text or textarea content for overview_content
- **Display order**: Set numeric display_order for solution positioning
- **Publish checkbox**: Toggle is_published checkbox, solution visibility changes
- **Edit solution**: Navigate to edit form, modify fields, save successfully
- **Delete solution**: hx-delete removes row from table without page reload
- **Cache invalidation**: After create/update/delete, page:solutions cache cleared

### Edge Cases / Error States
- **Missing required title**: Empty title triggers validation error
- **Duplicate slug**: Manual slug entry that conflicts with existing solution shows error
- **Invalid Material icon**: Non-existent icon name saved but may not render on frontend
- **Long short_description**: Very long text in short_description textarea accepted
- **Negative display_order**: Entering negative number in display_order field
- **Empty slug auto-generation**: Blank title doesn't generate slug until filled
- **Delete with sub-resources**: Deleting solution with stats/challenges/products may cascade delete or show warning

## Selectors & Elements
- List route: GET /admin/solutions
- Create form action: POST /admin/solutions
- Edit route: GET /admin/solutions/:id/edit
- Input names: title, slug, icon, short_description, hero_image_url, hero_title, hero_description, overview_content, meta_description, reference_code, is_published (checkbox), display_order (number)
- Delete button: hx-delete="/admin/solutions/:id" hx-target="closest tr"
- Submit button: text "Create Solution" or "Update Solution"
- HTMX sections: id="stats-section", id="challenges-section", id="products-section", id="ctas-section"

## HTMX Interactions
- **Slug auto-generation**: hx-get="/admin/solutions/generate-slug" hx-trigger="blur from:#title" hx-target="#slug"
- **Delete solution**: hx-delete="/admin/solutions/:id" hx-target="closest tr" hx-swap="outerHTML"
- **Sub-resource sections**: Each section loaded/updated via HTMX (see dependent test plans)

## Dependencies
- 20-solution-stats.md (stats HTMX section)
- 21-solution-challenges.md (challenges HTMX section)
- 22-solution-products.md (products HTMX section)
- 23-solution-ctas.md (CTAs HTMX section)
