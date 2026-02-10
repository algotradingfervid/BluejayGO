# Test Plan: Partner Tiers CRUD

## Summary
Verify complete CRUD operations for partner tiers with UNIQUE name validation and sort order management.

## Preconditions
- User authenticated with valid session cookie
- Database seeded with 2 partner tiers
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/partner-tiers
2. Verify list shows both seeded tiers
3. Click "New Tier" or navigate to http://localhost:28090/admin/partner-tiers/new
4. Fill form: name, description, sort_order
5. Submit POST to /admin/partner-tiers
6. Verify redirect to list with new tier visible
7. Click "Edit" on a tier or navigate to http://localhost:28090/admin/partner-tiers/:id/edit
8. Modify fields and submit POST to /admin/partner-tiers/:id
9. Verify updated data appears in list
10. Click delete button with hx-delete attribute
11. Confirm deletion in browser confirmation dialog
12. Verify tier row removed from table via HTMX

## Test Cases

### Happy Path
- **List all tiers**: Both seeded tiers display with names and descriptions
- **Create new tier**: Form submission creates tier with auto-generated slug
- **Edit existing tier**: Form pre-fills with current data, updates successfully
- **Delete tier via HTMX**: Delete button removes row without page reload
- **Auto-slug generation**: Slug automatically generated from name field
- **Sort order respected**: Tiers display in correct sort_order

### Edge Cases / Error States
- **Duplicate name validation**: Creating tier with existing name shows UNIQUE constraint error
- **Required name field**: Submitting without name shows validation error
- **Delete confirmation cancel**: Canceling hx-confirm dialog does not delete tier
- **Delete in-use tier**: Deleting tier assigned to partners prevents deletion or shows error
- **Empty description**: Tier can be created without description (optional field)
- **Negative sort_order**: System handles negative sort_order values correctly

## Selectors & Elements
- List page: http://localhost:28090/admin/partner-tiers
- Create form: `action="/admin/partner-tiers" method="POST"`
- Edit form: `action="/admin/partner-tiers/:id" method="POST"`
- Input name: `name="name" type="text"` (required, UNIQUE)
- Input slug: `name="slug" type="text"` (auto-generated)
- Textarea description: `name="description"`
- Input sort_order: `name="sort_order" type="number"`
- Delete button: `hx-delete="/admin/partner-tiers/:id" hx-confirm="Delete this tier?"`
- Table rows: one per tier in list view

## HTMX Interactions
- Delete action: `hx-delete="/admin/partner-tiers/:id"` with `hx-confirm` dialog
- Target: `hx-target="closest tr"` removes table row on successful delete
- Swap: `hx-swap="outerHTML"` replaces row element
- No HTMX on create/edit forms (standard POST with redirect)

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- Partners CRUD test plan (tiers used in partner profiles - not in current set)
