# Test Plan: Blog Categories CRUD

## Summary
Verify complete CRUD operations for blog categories with color hex codes and standard category fields.

## Preconditions
- User authenticated with valid session cookie
- Database seeded with 4 blog categories: Industry News (#1E88E5), Product Updates (#43A047), How-To Guides (#FB8C00), Company Announcements (#7B1FA2)
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/blog-categories
2. Verify list shows all 4 seeded categories with color indicators
3. Click "New Category" or navigate to http://localhost:28090/admin/blog-categories/new
4. Fill form: name, color_hex, description, sort_order
5. Submit POST to /admin/blog-categories
6. Verify redirect to list with new category visible
7. Click "Edit" on a category or navigate to http://localhost:28090/admin/blog-categories/:id/edit
8. Modify fields including color_hex and submit POST to /admin/blog-categories/:id
9. Verify updated data and color appear in list
10. Click delete button with hx-delete attribute
11. Confirm deletion in browser confirmation dialog
12. Verify category row removed from table via HTMX

## Test Cases

### Happy Path
- **List all categories**: All 4 seeded categories display with correct names and color indicators
- **Create new category**: Form submission creates category with auto-generated slug
- **Edit existing category**: Form pre-fills data, updates successfully including color changes
- **Delete category via HTMX**: Delete button removes row without page reload
- **Auto-slug generation**: Slug automatically generated from name field
- **Color preview**: Color hex value displays as visual indicator in list and form

### Edge Cases / Error States
- **Duplicate name validation**: Creating category with existing name shows UNIQUE constraint error
- **Required name field**: Submitting without name shows validation error
- **Invalid color hex format**: Entering invalid hex code (e.g., "xyz123" or missing #) shows validation error
- **Delete confirmation cancel**: Canceling hx-confirm dialog does not delete category
- **Delete in-use category**: Deleting category assigned to blog posts prevents deletion or shows error
- **Empty color hex**: Submitting without color shows validation error or uses default color

## Selectors & Elements
- List page: http://localhost:28090/admin/blog-categories
- Create form: `action="/admin/blog-categories" method="POST"`
- Edit form: `action="/admin/blog-categories/:id" method="POST"`
- Input name: `name="name" type="text"` (required, UNIQUE)
- Input slug: `name="slug" type="text"` (auto-generated)
- Input color_hex: `name="color_hex" type="text"` or `type="color"`
- Textarea description: `name="description"`
- Input sort_order: `name="sort_order" type="number"`
- Delete button: `hx-delete="/admin/blog-categories/:id" hx-confirm="Delete this category?"`
- Color indicator: visual element showing color_hex value in list view

## HTMX Interactions
- Delete action: `hx-delete="/admin/blog-categories/:id"` with `hx-confirm` dialog
- Target: `hx-target="closest tr"` removes table row on successful delete
- Swap: `hx-swap="outerHTML"` replaces row element
- No HTMX on create/edit forms (standard POST with redirect)

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- Blog posts CRUD test plan (categories used in blog post creation - not in current set)
