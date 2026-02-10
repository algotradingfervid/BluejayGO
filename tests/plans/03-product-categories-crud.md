# Test Plan: Product Categories CRUD

## Summary
Verify complete CRUD operations for product categories including list, create, update, and HTMX delete.

## Preconditions
- User authenticated with valid session cookie
- Database seeded with 5 product categories: Desktops, OPS Modules, Interactive Flat Panels, AV Accessories, IoT Products
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/product-categories
2. Verify list shows all 5 seeded categories
3. Click "New Category" or navigate to http://localhost:28090/admin/product-categories/new
4. Fill form: name, description, icon, image_url, sort_order
5. Submit POST to /admin/product-categories
6. Verify redirect to list with new category visible
7. Click "Edit" on a category or navigate to http://localhost:28090/admin/product-categories/:id/edit
8. Modify fields and submit POST to /admin/product-categories/:id
9. Verify updated data appears in list
10. Click delete button with hx-delete attribute
11. Confirm deletion in browser confirmation dialog
12. Verify category row removed from table via HTMX

## Test Cases

### Happy Path
- **List all categories**: All seeded categories display with correct data
- **Create new category**: Form submission creates category, auto-generates slug, redirects to list
- **Edit existing category**: Form pre-fills with current data, updates on submission
- **Delete category via HTMX**: Delete button triggers HTMX request, row removed without page reload
- **Auto-slug generation**: Slug automatically generated from name field
- **Sort order respected**: Categories display in correct sort_order

### Edge Cases / Error States
- **Duplicate name validation**: Creating category with existing name shows UNIQUE constraint error
- **Required name field**: Submitting without name shows validation error
- **Delete confirmation cancel**: Canceling hx-confirm dialog does not delete category
- **Delete in-use category**: Deleting category assigned to products shows foreign key error or prevents deletion
- **Invalid sort_order**: Non-numeric sort_order shows validation error
- **Long description**: Very long description text truncates or wraps correctly in list view

## Selectors & Elements
- List page: http://localhost:28090/admin/product-categories
- Create form: `action="/admin/product-categories" method="POST"`
- Edit form: `action="/admin/product-categories/:id" method="POST"`
- Input name: `name="name" type="text"` (required)
- Input slug: `name="slug" type="text"` (auto-generated)
- Textarea description: `name="description"`
- Input icon: `name="icon" type="text"`
- Input image_url: `name="image_url" type="url"`
- Input sort_order: `name="sort_order" type="number"`
- Delete button: `hx-delete="/admin/product-categories/:id" hx-confirm="Delete this category?"`
- Table rows: one per category in list view

## HTMX Interactions
- Delete action: `hx-delete="/admin/product-categories/:id"` with `hx-confirm` dialog
- Target: `hx-target="closest tr"` removes table row on successful delete
- Swap: `hx-swap="outerHTML"` (or similar) replaces row element
- No HTMX on create/edit forms (standard POST with redirect)

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 09-products-crud.md (categories used in product creation)
