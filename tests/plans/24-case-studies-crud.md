# Test Plan: Case Studies CRUD

## Summary
Testing case study creation, editing, listing with filters, deletion, and HTMX-driven sub-resource sections for products and metrics.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 3+ case studies and industries

## User Journey Steps
1. Navigate to http://localhost:28090/admin/case-studies
2. Verify case study list displays with search, status, page filters
3. Click "New Case Study" button to navigate to /admin/case-studies/new
4. Fill required fields: title, client_name
5. Auto-generate slug on title blur
6. Select industry_id from dropdown
7. Fill summary textarea
8. Fill challenge section: challenge_title, challenge_content (textarea/rich), challenge_bullets (comma-separated)
9. Fill solution section: solution_title, solution_content
10. Fill outcome section: outcome_title, outcome_content
11. Fill SEO fields: meta_title, meta_description
12. Upload hero_image_url
13. Set is_published checkbox and display_order
14. Submit form via POST /admin/case-studies
15. Edit existing case study at /admin/case-studies/:id/edit with HTMX sections
16. Delete case study using hx-delete

## Test Cases

### Happy Path
- **List case studies with filters**: Verify search by title, filter by status/page work correctly
- **Create new case study**: Required fields (title, client_name) filled, slug auto-generated, case study saved
- **Auto-slug generation**: Title "Acme Corp Success" generates slug "acme-corp-success" on blur
- **Industry selection**: Select industry from dropdown, association saved
- **Challenge bullets**: Enter "Cost overruns, Delays, Quality issues" comma-separated, stored as JSON array
- **Rich text content**: challenge_content, solution_content, outcome_content support textarea or rich editor
- **SEO fields**: meta_title and meta_description filled, used for frontend SEO
- **Display order**: Set numeric display_order for case study positioning
- **Publish checkbox**: Toggle is_published checkbox, case study visibility changes
- **Edit case study**: Navigate to edit form, modify fields, interact with HTMX sections
- **Delete case study**: hx-delete removes row from table without page reload
- **Cache invalidation**: After create/update/delete, page:case-studies cache cleared

### Edge Cases / Error States
- **Missing required title**: Empty title triggers validation error
- **Missing required client_name**: Empty client_name triggers validation error
- **Duplicate slug**: Manual slug entry that conflicts shows error
- **Empty challenge_bullets**: No bullets entered, stored as empty JSON array
- **Invalid comma format**: Bullets "item1,item2,,item3" with double comma handled gracefully
- **Long summary**: Summary with 500+ characters accepted
- **No industry selected**: industry_id null accepted, optional field
- **Delete with sub-resources**: Deleting case study with products/metrics may cascade delete

## Selectors & Elements
- List route: GET /admin/case-studies
- Create form action: POST /admin/case-studies
- Edit route: GET /admin/case-studies/:id/edit
- Input names: title, slug, client_name, industry_id (select), summary (textarea), hero_image_url, challenge_title, challenge_content (textarea), challenge_bullets (textarea), solution_title, solution_content, outcome_title, outcome_content, meta_title, meta_description, is_published (checkbox), display_order (number)
- Delete button: hx-delete="/admin/case-studies/:id" hx-target="closest tr"
- Submit button: text "Create Case Study" or "Update Case Study"
- HTMX sections: id="products-section", id="metrics-section"

## HTMX Interactions
- **Slug auto-generation**: hx-get="/admin/case-studies/generate-slug" hx-trigger="blur from:#title" hx-target="#slug"
- **Delete case study**: hx-delete="/admin/case-studies/:id" hx-target="closest tr" hx-swap="outerHTML"
- **Sub-resource sections**: Products and metrics sections loaded/updated via HTMX (see dependent test plans)

## Dependencies
- 25-case-study-products.md (products HTMX section)
- 26-case-study-metrics.md (metrics HTMX section)
- Industries table must be seeded with data
