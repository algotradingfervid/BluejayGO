# Test Plan: Blog Authors CRUD

## Summary
Verify complete CRUD operations for blog authors including profile fields, social links, and auto-slug generation.

## Preconditions
- User authenticated with valid session cookie
- Database seeded with 6 blog authors
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/blog-authors
2. Verify list shows all 6 seeded authors
3. Click "New Author" or navigate to http://localhost:28090/admin/blog-authors/new
4. Fill form: name, title, bio, avatar_url, linkedin_url, email, sort_order
5. Submit POST to /admin/blog-authors
6. Verify redirect to list with new author visible
7. Click "Edit" on an author or navigate to http://localhost:28090/admin/blog-authors/:id/edit
8. Modify fields and submit POST to /admin/blog-authors/:id
9. Verify updated data appears in list
10. Click delete button with hx-delete attribute
11. Confirm deletion in browser confirmation dialog
12. Verify author row removed from table via HTMX

## Test Cases

### Happy Path
- **List all authors**: All 6 seeded authors display with names, titles, and avatars
- **Create new author**: Form submission creates author with auto-generated slug
- **Edit existing author**: Form pre-fills with current data, updates successfully
- **Delete author via HTMX**: Delete button removes row without page reload
- **Auto-slug generation**: Slug automatically generated from name field
- **Avatar display**: Avatar URL shows preview image in list and form
- **LinkedIn link validation**: Valid LinkedIn URL accepted and stored

### Edge Cases / Error States
- **Required name field**: Submitting without name shows validation error
- **Invalid email format**: Entering invalid email shows validation error
- **Invalid URL format**: Entering invalid avatar_url or linkedin_url shows validation error
- **Delete confirmation cancel**: Canceling hx-confirm dialog does not delete author
- **Delete in-use author**: Deleting author assigned to blog posts prevents deletion or shows error
- **Long bio text**: Very long bio content saves correctly and displays properly
- **Empty optional fields**: Author can be created with only required name field

## Selectors & Elements
- List page: http://localhost:28090/admin/blog-authors
- Create form: `action="/admin/blog-authors" method="POST"`
- Edit form: `action="/admin/blog-authors/:id" method="POST"`
- Input name: `name="name" type="text"` (required)
- Input slug: `name="slug" type="text"` (auto-generated)
- Input title: `name="title" type="text"`
- Textarea bio: `name="bio"`
- Input avatar_url: `name="avatar_url" type="url"`
- Input linkedin_url: `name="linkedin_url" type="url"`
- Input email: `name="email" type="email"`
- Input sort_order: `name="sort_order" type="number"`
- Delete button: `hx-delete="/admin/blog-authors/:id" hx-confirm="Delete this author?"`
- Avatar preview: img element showing avatar_url in list/form

## HTMX Interactions
- Delete action: `hx-delete="/admin/blog-authors/:id"` with `hx-confirm` dialog
- Target: `hx-target="closest tr"` removes table row on successful delete
- Swap: `hx-swap="outerHTML"` replaces row element
- No HTMX on create/edit forms (standard POST with redirect)

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- Blog posts CRUD test plan (authors used in blog post creation - not in current set)
