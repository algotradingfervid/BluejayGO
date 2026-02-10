# Test Plan: Homepage Call-to-Action CRUD

## Summary
Tests the complete CRUD operations for homepage CTA sections including headlines, descriptions, primary/secondary CTAs, background styles, and active status.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with 1 CTA section
- Background style options available (e.g., light, dark, primary, gradient)

## User Journey Steps
1. Navigate to /admin/homepage/cta
2. View list of existing CTA sections (1 seeded)
3. Click "New CTA" button
4. Fill in headline
5. Fill in description (textarea)
6. Set primary CTA text and URL
7. Optionally set secondary CTA text and URL
8. Select background_style from dropdown
9. Check/uncheck is_active checkbox
10. Submit to create CTA
11. Edit existing CTA
12. Update headline, description, and CTAs
13. Change background_style
14. Delete a CTA
15. Verify active status filtering

## Test Cases

### Happy Path
- **List CTAs**: Verifies GET /admin/homepage/cta shows seeded CTA section
- **Create new CTA**: Adds CTA with all fields, verifies creation
- **Edit existing CTA**: Updates headline and description, verifies save
- **Update primary CTA**: Changes primary_cta_text and primary_cta_url, verifies update
- **Add secondary CTA**: Adds secondary CTA where none existed, verifies save
- **Remove secondary CTA**: Clears secondary CTA fields, verifies removal
- **Change background style**: Selects different background_style option, verifies update
- **Toggle is_active**: Checks/unchecks is_active, verifies status change
- **Multiple CTAs**: Creates multiple CTA sections, verifies all listed

### Edge Cases / Error States
- **Empty headline**: Tests validation when headline is empty
- **Empty description**: Tests validation when description is empty
- **Empty primary CTA text**: Tests validation when primary CTA text missing
- **Empty primary CTA URL**: Tests validation when primary CTA URL missing
- **Secondary CTA partial**: Enters only text or only URL for secondary CTA, checks validation
- **Invalid primary URL**: Enters malformed primary_cta_url, checks validation
- **Invalid secondary URL**: Enters malformed secondary_cta_url, checks validation
- **Very long headline**: Tests character limits on headline
- **Very long description**: Tests textarea character limits
- **Invalid background style**: Tests if invalid style value is rejected
- **Multiple active CTAs**: Creates multiple CTAs with is_active=true, verifies handling
- **All CTAs inactive**: Sets all is_active to false, checks empty state
- **Delete CTA**: Deletes CTA, verifies removal from list
- **Delete confirmation**: Verifies confirmation before deletion

## Selectors & Elements
- CTAs list: `#cta-list` or `.cta-table`
- New CTA button: `a[href="/admin/homepage/cta/new"]` or `button#new-cta`
- CTA row: `.cta-row[data-id]` or `tr[data-cta-id]`
- Edit link: `a[href="/admin/homepage/cta/{id}/edit"]`
- Delete button: `button[type="submit"]` in delete form or delete link
- Headline input: `input[name="headline"]`
- Description textarea: `textarea[name="description"]`
- Primary CTA text: `input[name="primary_cta_text"]`
- Primary CTA URL: `input[name="primary_cta_url"]`
- Secondary CTA text: `input[name="secondary_cta_text"]`
- Secondary CTA URL: `input[name="secondary_cta_url"]`
- Background style select: `select[name="background_style"]`
- Background style options: `option[value="light"]`, `option[value="dark"]`, etc.
- Is active checkbox: `input[name="is_active"][type="checkbox"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`
- Background preview: Element showing selected background_style

## HTMX Interactions
- None specified - standard form submissions with redirects
- Delete may use HTMX hx-delete if implemented

## Dependencies
- Database seeded with 1 CTA section
- Background style options defined (light/dark/primary/gradient/etc.)
- Template: templates/admin/pages/homepage-cta-list.html, homepage-cta-form.html
- Handler: internal/handlers/homepage.go (ListCTAs, NewCTA, CreateCTA, EditCTA, UpdateCTA, DeleteCTA)
