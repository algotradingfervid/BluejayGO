# Test Plan: Industries CRUD

## Summary
Verify complete CRUD operations for industries including icon fields and UNIQUE name validation.

## Preconditions
- User authenticated with valid session cookie
- Database seeded with 6 industries
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/industries
2. Verify list shows all 6 seeded industries
3. Click "New Industry" or navigate to http://localhost:28090/admin/industries/new
4. Fill form: name, icon, description, sort_order
5. Submit POST to /admin/industries
6. Verify redirect to list with new industry visible
7. Click "Edit" on an industry or navigate to http://localhost:28090/admin/industries/:id/edit
8. Modify fields and submit POST to /admin/industries/:id
9. Verify updated data appears in list
10. Click delete button with hx-delete attribute
11. Confirm deletion in browser confirmation dialog
12. Verify industry row removed from table via HTMX

## Test Cases

### Happy Path
- **List all industries**: All 6 seeded industries display with names and icons
- **Create new industry**: Form submission creates industry with auto-generated slug
- **Edit existing industry**: Form pre-fills with current data, updates successfully
- **Delete industry via HTMX**: Delete button removes row without page reload
- **Auto-slug generation**: Slug automatically generated from name field
- **Icon display**: Icon value displays as visual indicator in list view
- **Sort order respected**: Industries display in correct sort_order

### Edge Cases / Error States
- **Duplicate name validation**: Creating industry with existing name shows UNIQUE constraint error
- **Required name field**: Submitting without name shows validation error
- **Delete confirmation cancel**: Canceling hx-confirm dialog does not delete industry
- **Delete in-use industry**: Deleting industry assigned to products/partners prevents deletion or shows error
- **Empty icon field**: Industry can be created without icon (optional field)
- **Long description**: Very long description text saves and displays correctly

## Selectors & Elements
- List page: http://localhost:28090/admin/industries
- Create form: `action="/admin/industries" method="POST"`
- Edit form: `action="/admin/industries/:id" method="POST"`
- Input name: `name="name" type="text"` (required, UNIQUE)
- Input slug: `name="slug" type="text"` (auto-generated)
- Input icon: `name="icon" type="text"`
- Textarea description: `name="description"`
- Input sort_order: `name="sort_order" type="number"`
- Delete button: `hx-delete="/admin/industries/:id" hx-confirm="Delete this industry?"`
- Icon indicator: visual element showing icon in list view

## HTMX Interactions
- Delete action: `hx-delete="/admin/industries/:id"` with `hx-confirm` dialog
- Target: `hx-target="closest tr"` removes table row on successful delete
- Swap: `hx-swap="outerHTML"` replaces row element
- No HTMX on create/edit forms (standard POST with redirect)

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- Products CRUD test plan (industries may be used in product categorization)
- Partners CRUD test plan (industries used in partner profiles - not in current set)
