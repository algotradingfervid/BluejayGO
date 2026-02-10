# Test Plan: Homepage Testimonials CRUD

## Summary
Tests the complete CRUD operations for homepage testimonials including quotes, author information, ratings, images, active status, and display ordering.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with 3 testimonials
- Valid image URLs for author images
- Rating system supports 1-5 stars

## User Journey Steps
1. Navigate to /admin/homepage/testimonials
2. View list of existing testimonials (3 seeded)
3. Click "New Testimonial" button
4. Fill in quote (required textarea)
5. Fill in author_name (required)
6. Fill in author_title and author_company (optional)
7. Add author_image URL (optional)
8. Set rating (number 1-5)
9. Set display_order
10. Check/uncheck is_active checkbox
11. Submit to create testimonial
12. Edit existing testimonial
13. Update quote and author information
14. Delete a testimonial
15. Verify ordering and active status filtering

## Test Cases

### Happy Path
- **List testimonials**: Verifies GET /admin/homepage/testimonials shows all 3 seeded testimonials
- **Create new testimonial**: Adds testimonial with required fields, verifies creation
- **Create with all fields**: Adds testimonial with quote, author info, image, rating, verifies creation
- **Edit existing testimonial**: Updates quote text, verifies save
- **Update author info**: Changes author_name, author_title, author_company, verifies update
- **Update rating**: Changes rating from 4 to 5, verifies update
- **Add author image**: Adds author_image URL where none existed, verifies update
- **Remove author image**: Clears author_image URL, verifies removal
- **Toggle is_active**: Checks/unchecks is_active, verifies status change
- **Reorder testimonials**: Changes display_order, verifies list reordering

### Edge Cases / Error States
- **Empty quote**: Tests required field validation on quote textarea
- **Empty author_name**: Tests required field validation on author_name
- **Very long quote**: Tests textarea character limits
- **Very long author name**: Tests character limit on author_name
- **Invalid rating**: Enters rating outside 1-5 range (e.g., 0, 6, 10), checks validation
- **Negative rating**: Enters negative number, checks validation
- **Decimal rating**: Enters 3.5 or 4.7, checks if allowed or rounded
- **Invalid author image URL**: Enters malformed URL, checks validation
- **Missing author title**: Leaves author_title empty, verifies optional handling
- **Missing author company**: Leaves author_company empty, verifies optional handling
- **Duplicate display_order**: Creates testimonials with same order, verifies handling
- **All testimonials inactive**: Sets all is_active to false, checks empty state
- **Delete testimonial**: Deletes testimonial, verifies removal
- **Delete confirmation**: Verifies confirmation before deletion

## Selectors & Elements
- Testimonials list: `#testimonials-list` or `.testimonials-table`
- New testimonial button: `a[href="/admin/homepage/testimonials/new"]` or `button#new-testimonial`
- Testimonial row: `.testimonial-row[data-id]` or `tr[data-testimonial-id]`
- Edit link: `a[href="/admin/homepage/testimonials/{id}/edit"]`
- Delete button: `button[type="submit"]` in delete form or delete link
- Quote textarea: `textarea[name="quote"]`
- Author name input: `input[name="author_name"]`
- Author title input: `input[name="author_title"]`
- Author company input: `input[name="author_company"]`
- Author image URL: `input[name="author_image"]`
- Rating input: `input[name="rating"][type="number"]` with min="1" max="5"
- Display order input: `input[name="display_order"][type="number"]`
- Is active checkbox: `input[name="is_active"][type="checkbox"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`
- Star rating display: `.rating-stars` or visual rating indicator

## HTMX Interactions
- None specified - standard form submissions with redirects
- Delete may use HTMX hx-delete if implemented

## Dependencies
- Database seeded with 3 testimonials
- Template: templates/admin/pages/homepage-testimonials-list.html, homepage-testimonials-form.html
- Handler: internal/handlers/homepage.go (ListTestimonials, NewTestimonial, CreateTestimonial, EditTestimonial, UpdateTestimonial, DeleteTestimonial)
