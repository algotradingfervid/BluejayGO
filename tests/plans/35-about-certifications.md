# Test Plan: Certifications CRUD

## Summary
Tests the complete CRUD operations for company certifications including name, abbreviation, description, icon, and display order management.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with certifications data
- Material Icons available for icon field

## User Journey Steps
1. Navigate to /admin/about/certifications
2. View list of existing certifications (seeded data)
3. Click "New Certification" button
4. Fill in name, abbreviation, description, icon
5. Set display_order
6. Submit to create certification
7. Edit an existing certification
8. Update fields and save
9. Delete a certification using HTMX delete button
10. Verify certification removed without page reload

## Test Cases

### Happy Path
- **List certifications**: Verifies GET /admin/about/certifications shows all seeded certifications
- **Create new certification**: Adds certification with all fields, verifies creation
- **Edit existing certification**: Updates name and description, verifies save
- **Update abbreviation**: Changes abbreviation, verifies update
- **Update icon**: Changes Material icon name, verifies new icon displays
- **Reorder certifications**: Changes display_order values, verifies list reordering
- **View with abbreviation**: Verifies abbreviation displays correctly in list/cards

### Edge Cases / Error States
- **Empty name**: Tests required field validation on name
- **Empty abbreviation**: Tests if abbreviation is required or optional
- **Empty description**: Tests required field validation on description
- **Invalid icon**: Enters non-existent Material icon name, checks validation/fallback
- **Duplicate names**: Creates certifications with same name, verifies handling
- **Duplicate abbreviations**: Creates certifications with same abbreviation, checks handling
- **Very long name**: Tests character limit on name field
- **Very long abbreviation**: Tests abbreviation field length limit
- **Very long description**: Tests textarea character limits
- **Delete via HTMX**: Clicks delete button, verifies hx-delete removes item
- **Delete confirmation**: Verifies confirmation modal/prompt before deletion
- **Cache invalidation**: Confirms page:about cache cleared after changes

## Selectors & Elements
- Certifications list: `#certifications-list` or `.certifications-table`
- New certification button: `a[href="/admin/about/certifications/new"]` or `button#new-certification`
- Certification row: `.certification-row[data-id]` or `tr[data-certification-id]`
- Delete button: `button[hx-delete="/admin/about/certifications/{id}"]`
- Name input: `input[name="name"]`
- Abbreviation input: `input[name="abbreviation"]`
- Description textarea: `textarea[name="description"]`
- Icon input: `input[name="icon"]`
- Display order input: `input[name="display_order"][type="number"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`
- Icon preview: `.icon-preview` or Material icon display element

## HTMX Interactions
- **hx-delete**: Delete button uses `hx-delete="/admin/about/certifications/{id}"`
- **hx-target**: Targets parent row or container for removal
- **hx-swap**: Uses `outerHTML` or `delete` to remove element from DOM
- **hx-confirm**: Confirmation message before delete action

## Dependencies
- Database seeded with certifications
- Cache service for page:about
- Material Icons for icon display
- Template: templates/admin/pages/about-certifications-list.html, about-certifications-form.html
- Handler: internal/handlers/about.go (ListCertifications, NewCertification, CreateCertification, EditCertification, UpdateCertification, DeleteCertification)
- HTMX library loaded
