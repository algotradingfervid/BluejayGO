# Test Plan: Partner Testimonials CRUD

## Summary
Testing partner testimonial creation, editing, listing, deletion, and cache invalidation at dedicated /admin/partners/testimonials routes.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 6 partner testimonials and 11 partners

## User Journey Steps
1. Navigate to http://localhost:28090/admin/partners/testimonials
2. Verify testimonial list displays all 6 seeded testimonials
3. Click "New Testimonial" button to navigate to /admin/partners/testimonials/new
4. Select partner_id from dropdown (select from partners)
5. Fill quote textarea (required)
6. Fill author_name (required)
7. Fill author_title (required)
8. Set display_order number
9. Set is_active checkbox
10. Submit form via POST /admin/partners/testimonials
11. Edit existing testimonial at /admin/partners/testimonials/:id/edit
12. Delete testimonial using hx-delete
13. Verify cache invalidation for page:partners

## Test Cases

### Happy Path
- **List testimonials**: Navigate to /admin/partners/testimonials, see all 6 seeded testimonials
- **Create new testimonial**: Select partner, fill quote, author_name, author_title, submit, testimonial saved
- **Partner selection**: Select partner from dropdown of all 11 partners, association saved
- **Quote textarea**: Fill quote with 200+ character testimonial text, saved successfully
- **Author details**: Fill author_name "John Doe", author_title "CTO", saved successfully
- **Display order**: Set numeric display_order for testimonial positioning
- **Active checkbox**: Toggle is_active checkbox, testimonial visibility changes
- **Multiple testimonials per partner**: Create 2 testimonials for same partner, both saved
- **Edit testimonial**: Navigate to edit form, modify fields, save successfully
- **Delete testimonial**: hx-delete removes row from table without page reload
- **Cache invalidation**: After create/update/delete, page:partners cache cleared

### Edge Cases / Error States
- **Missing required partner_id**: Not selecting partner triggers validation error
- **Missing required quote**: Empty quote textarea triggers validation error
- **Missing required author_name**: Empty author_name triggers validation error
- **Missing required author_title**: Empty author_title triggers validation error
- **Very long quote**: Quote with 1000+ characters accepted in textarea
- **Long author_title**: author_title "Senior Vice President of Technology" accepted, may wrap
- **Duplicate display_order**: Multiple testimonials with same display_order accepted, sorted arbitrarily
- **Inactive testimonial**: is_active false hides testimonial on frontend
- **Delete partner cascade**: Deleting partner with testimonials may cascade delete testimonials or show error

## Selectors & Elements
- List route: GET /admin/partners/testimonials
- Create form action: POST /admin/partners/testimonials
- Edit route: GET /admin/partners/testimonials/:id/edit
- Input names: partner_id (select required), quote (textarea required), author_name (required), author_title (required), display_order (number), is_active (checkbox)
- Delete button: hx-delete="/admin/partners/testimonials/:id" hx-target="closest tr"
- Submit button: text "Create Testimonial" or "Update Testimonial"

## HTMX Interactions
- **Delete testimonial**: hx-delete="/admin/partners/testimonials/:id" hx-target="closest tr" hx-swap="outerHTML"

## Dependencies
- 29-partners-crud.md (partners must exist to create testimonials)
- Partners table seeded with 11 partners
