# Test Plan: Blog Tags CRUD

## Summary
Testing tag creation, listing, HTMX autocomplete search, and inline quick-create functionality used in blog post forms.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 13 tags

## User Journey Steps
1. Navigate to http://localhost:28090/admin/blog/tags
2. Verify list of all tags displays
3. Create new tag via POST /admin/blog/tags with name field
4. Test HTMX autocomplete search via GET /admin/blog/tags/search?q=keyword
5. Verify search returns tag_suggestions.html partial
6. Test quick-create via POST /admin/blog/tags/quick-create from blog post form
7. Verify quick-create returns tag_chip.html appended to #selected-tags
8. Delete tag using DELETE /admin/blog/tags/:id

## Test Cases

### Happy Path
- **List all tags**: Navigate to /admin/blog/tags, verify 13 seeded tags display
- **Create new tag**: Submit form with name "Technology", tag created successfully
- **Search tags autocomplete**: Type "tech" in search input, HTMX returns matching tag_suggestions.html partial
- **Quick-create from post form**: Type new tag name, click quick-create, tag_chip.html appended to #selected-tags
- **Delete tag**: Click delete button, tag removed via hx-delete

### Edge Cases / Error States
- **Empty tag name**: Creating tag with blank name shows validation error
- **Duplicate tag name**: Creating tag with existing name shows uniqueness error
- **Search with no results**: Query "xyz123nonexistent" returns empty suggestions partial
- **Search delay trigger**: Typing rapidly waits 200ms before firing HTMX request
- **Delete tag in use**: Deleting tag assigned to posts may show constraint error or warning
- **Quick-create duplicate**: Quick-creating existing tag name handles gracefully

## Selectors & Elements
- List route: GET /admin/blog/tags
- Create form action: POST /admin/blog/tags
- Input name: name (required)
- Search endpoint: GET /admin/blog/tags/search with query param q
- Quick-create endpoint: POST /admin/blog/tags/quick-create
- Delete button: hx-delete="/admin/blog/tags/:id"
- Tag search input: id="tag-search"
- Suggestions container: id="tag-suggestions"
- Selected tags container: id="selected-tags"

## HTMX Interactions
- **Autocomplete search**: hx-get="/admin/blog/tags/search" hx-trigger="input changed delay:200ms, focus" hx-target="#tag-suggestions" (returns tag_suggestions.html partial)
- **Quick-create**: hx-post="/admin/blog/tags/quick-create" hx-vals='{"name":"New Tag"}' hx-target="#selected-tags" hx-swap="beforeend" (returns tag_chip.html partial)
- **Delete**: hx-delete="/admin/blog/tags/:id" hx-target="closest li" or similar

## Dependencies
- Used by 17-blog-post-tags.md for tag assignment in blog posts
