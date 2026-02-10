# Test Plan: Office Locations CRUD

## Summary
Tests the complete CRUD operations for office locations with primary office enforcement, active status management, and display ordering.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with 6 offices: Bangalore HQ, Mumbai, Gurugram, Chennai, Hyderabad, Pune
- Only one office can be primary at a time (enforced by is_primary checkbox)
- Default country: India
- Cache: page:contact cleared on updates

## User Journey Steps
1. Navigate to /admin/contact/offices
2. View list of existing offices (6 seeded)
3. Identify which office is marked as primary (Bangalore HQ)
4. Click "New Office" button
5. Fill in name, address_line1 (required), city (required)
6. Optionally fill address_line2, state, postal_code
7. Country defaults to "India"
8. Fill phone and email
9. Check/uncheck is_primary checkbox
10. Check/uncheck is_active checkbox
11. Set display_order
12. Submit to create office
13. Edit existing office
14. Update fields including toggling is_primary
15. Delete office using HTMX delete button
16. Verify only one office can be primary at a time
17. Verify cache:page:contact is cleared

## Test Cases

### Happy Path
- **List offices**: Verifies GET /admin/contact/offices shows all 6 seeded offices
- **Create new office**: Adds office with required fields, verifies creation
- **Create with all fields**: Fills all fields including optional ones, verifies creation
- **Edit existing office**: Updates name and address, verifies save
- **Update address**: Changes address_line1 and address_line2, verifies update
- **Update city**: Changes city field, verifies update
- **Update state**: Changes state field, verifies update
- **Update postal code**: Changes postal_code, verifies update
- **Update phone**: Changes phone number, verifies update
- **Update email**: Changes email, verifies update
- **Default country**: Verifies country defaults to "India" on new office form
- **Change country**: Changes country from India to another country, verifies update
- **Toggle is_active**: Checks/unchecks is_active, verifies status change
- **Reorder offices**: Changes display_order values, verifies list reordering
- **Primary office indicator**: Verifies primary office has distinct badge/indicator

### Happy Path - Primary Office Enforcement
- **Set new primary**: Checks is_primary on Mumbai office, verifies it becomes primary
- **Previous primary reset**: After setting Mumbai as primary, verifies Bangalore HQ is_primary=false
- **Only one primary**: Verifies system enforces single primary office constraint
- **Create as primary**: Creates new office with is_primary=true, verifies old primary updated

### Happy Path - Deletion
- **Delete non-primary office**: Deletes Pune office via hx-delete, verifies removal
- **Delete confirmation**: Verifies confirmation prompt before deletion
- **Office removed**: After delete, verifies office removed from list without page reload

### Edge Cases / Error States
- **Empty name**: Tests required field validation on name
- **Empty address_line1**: Tests required field validation on address_line1
- **Empty city**: Tests required field validation on city
- **Empty state**: Tests if state is required or optional
- **Empty postal_code**: Tests if postal_code is required or optional
- **Empty country**: Tests if country can be empty or defaults to India
- **Empty phone**: Tests if phone is required or optional
- **Empty email**: Tests if email is required or optional
- **Invalid email format**: Enters malformed email, checks validation
- **Invalid phone format**: Enters invalid phone number, checks validation
- **Very long name**: Tests character limit on name field
- **Very long address**: Tests character limits on address_line1 and address_line2
- **Very long city**: Tests character limit on city field
- **Duplicate office names**: Creates offices with same name, verifies handling
- **All offices inactive**: Sets all is_active to false, verifies contact page behavior
- **Delete primary office**: Attempts to delete primary office, checks if prevented or handled
- **No primary office**: If all offices have is_primary=false, tests fallback behavior
- **Multiple primary attempt**: Tests if database constraint prevents multiple primary offices
- **Create with is_primary false**: Creates office without setting primary, verifies existing primary unchanged
- **Display order duplicates**: Creates offices with same display_order, verifies handling
- **Cache invalidation**: Confirms page:contact cache cleared after create/update/delete
- **Delete via HTMX**: Verifies hx-delete removes row from DOM
- **Delete non-existent office**: Attempts delete on non-existent ID, verifies error
- **Country dropdown**: If country is select dropdown, tests all options
- **Country text input**: If country is text input, tests validation

## Selectors & Elements
- Offices list: `#offices-list` or `.offices-table`
- New office button: `a[href="/admin/contact/offices/new"]` or `button#new-office`
- Office row: `.office-row[data-id]` or `tr[data-office-id]`
- Primary badge: `.badge-primary` or `.is-primary-indicator`
- Active badge: `.badge-active` or indicator for is_active=true
- Edit link: `a[href="/admin/contact/offices/{id}/edit"]`
- Delete button: `button[hx-delete="/admin/contact/offices/{id}"]`
- Name input: `input[name="name"]`
- Address line 1: `input[name="address_line1"]`
- Address line 2: `input[name="address_line2"]`
- City input: `input[name="city"]`
- State input: `input[name="state"]`
- Postal code input: `input[name="postal_code"]`
- Country input/select: `input[name="country"]` or `select[name="country"]` with default value="India"
- Phone input: `input[name="phone"]` with type="tel"
- Email input: `input[name="email"]` with type="email"
- Is primary checkbox: `input[name="is_primary"][type="checkbox"]`
- Is active checkbox: `input[name="is_active"][type="checkbox"]`
- Display order input: `input[name="display_order"][type="number"]`
- Submit button: `button[type="submit"]`
- Success message: `.alert-success`

## HTMX Interactions
- **hx-delete**: Delete button uses `hx-delete="/admin/contact/offices/{id}"`
- **hx-confirm**: Confirmation message before delete
- **hx-target**: Targets parent row for removal
- **hx-swap**: Uses `outerHTML` or `delete` to remove element

## Dependencies
- Database seeded with 6 offices (Bangalore HQ as primary, others in major Indian cities)
- office_locations table columns: id, name, address_line1, address_line2, city, state, postal_code, country, phone, email, is_primary, is_active, display_order
- Database constraint or application logic to enforce single is_primary=true
- Cache service for page:contact
- Template: templates/admin/pages/office-locations-list.html, office-locations-form.html
- Handler: internal/handlers/contact.go (ListOffices, NewOffice, CreateOffice, EditOffice, UpdateOffice, DeleteOffice)
- HTMX library loaded
- Country defaults to "India" for new offices
- Primary office enforcement: when setting new primary, old primary must be unset
