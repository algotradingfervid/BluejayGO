# Test Plan: Page Sections Management

## Summary
Tests the editing of pre-seeded page sections with no create/delete functionality, only updating existing sections for specific pages.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with page sections: home, products, product_detail, solutions, solution_detail, footer
- No create or delete functionality available (edit only)

## User Journey Steps
1. Navigate to /admin/page-sections
2. View list of all pre-seeded page sections (6 items)
3. Click edit link for a section (e.g., "home" section)
4. Navigate to /admin/page-sections/:id/edit
5. Update heading, subheading, description
6. Update label field
7. Update primary button text and URL
8. Update secondary button text and URL
9. Toggle is_active checkbox
10. Update display_order
11. Submit form via POST /admin/page-sections/:id
12. Verify redirect back to edit page or list
13. Confirm updated data persists
14. Verify no "New Section" or "Delete" buttons exist

## Test Cases

### Happy Path
- **List page sections**: Verifies GET /admin/page-sections shows all 6 seeded sections
- **Edit home section**: Updates home section fields, verifies save
- **Edit products section**: Updates products section fields, verifies save
- **Edit product_detail section**: Updates product detail fields, verifies save
- **Edit solutions section**: Updates solutions section fields, verifies save
- **Edit solution_detail section**: Updates solution detail fields, verifies save
- **Edit footer section**: Updates footer section fields, verifies save
- **Update heading**: Changes heading text, verifies save
- **Update subheading**: Changes subheading text, verifies save
- **Update description**: Changes description textarea, verifies save
- **Update label**: Changes label field, verifies save
- **Update primary button**: Changes primary_button_text and primary_button_url, verifies save
- **Remove primary button**: Clears primary button fields, verifies removal
- **Update secondary button**: Changes secondary_button_text and secondary_button_url, verifies save
- **Remove secondary button**: Clears secondary button fields, verifies removal
- **Toggle is_active**: Checks/unchecks is_active checkbox, verifies status change
- **Update display_order**: Changes display_order value, verifies list reordering

### Edge Cases / Error States
- **No create button**: Verifies "New Section" button does not exist on list page
- **No delete button**: Verifies delete button/option does not exist on edit page
- **Empty heading**: Tests validation when heading is cleared
- **Empty subheading**: Tests if subheading is required or optional
- **Empty description**: Tests if description is required or optional
- **Empty label**: Tests if label is required or optional
- **Primary button partial**: Enters text without URL or URL without text, checks validation
- **Secondary button partial**: Enters text without URL or URL without text, checks validation
- **Invalid primary URL**: Enters malformed primary_button_url, checks validation
- **Invalid secondary URL**: Enters malformed secondary_button_url, checks validation
- **Very long heading**: Tests character limits on heading field
- **Very long description**: Tests textarea character limits
- **All sections inactive**: Sets all sections is_active=false, checks site behavior
- **Duplicate display_order**: Sets multiple sections to same order, verifies handling
- **Section identifier immutable**: Verifies section identifier (home, products, etc.) cannot be changed

## Selectors & Elements
- Sections list: `#page-sections-list` or `.sections-table`
- Section row: `.section-row[data-id]` or `tr[data-section-id]`
- Edit link: `a[href="/admin/page-sections/{id}/edit"]`
- No create button: Verify absence of `a[href="/admin/page-sections/new"]` or `button#new-section`
- No delete button: Verify absence of delete button/form on edit page
- Form: `form[action="/admin/page-sections/{id}"][method="POST"]`
- Heading input: `input[name="heading"]`
- Subheading input: `input[name="subheading"]`
- Description textarea: `textarea[name="description"]`
- Label input: `input[name="label"]`
- Primary button text: `input[name="primary_button_text"]`
- Primary button URL: `input[name="primary_button_url"]`
- Secondary button text: `input[name="secondary_button_text"]`
- Secondary button URL: `input[name="secondary_button_url"]`
- Is active checkbox: `input[name="is_active"][type="checkbox"]`
- Display order input: `input[name="display_order"][type="number"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`
- Section identifier display: `.section-identifier` or read-only field showing page type

## HTMX Interactions
- None - standard form POST with redirect to list or edit page

## Dependencies
- Database seeded with exactly 6 page sections (home, products, product_detail, solutions, solution_detail, footer)
- No create/delete handlers or routes
- Template: templates/admin/pages/page-sections-list.html, page-sections-edit.html
- Handler: internal/handlers/page_sections.go (ListPageSections, EditPageSection, UpdatePageSection)
- Section identifiers are immutable and tied to specific pages
