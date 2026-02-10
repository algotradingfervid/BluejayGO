# Test Plan: Blog Posts CRUD

## Summary
Comprehensive testing of blog post creation, editing, listing with filters, and deletion including HTMX-driven tag/product pickers and rich text editing.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 7+ blog posts, categories, authors, tags, products
- Trix editor CSS/JS loaded from unpkg CDN

## User Journey Steps
1. Navigate to http://localhost:28090/admin/blog/posts
2. Verify post list displays with search, status, category, author filters
3. Verify pagination shows 15 posts per page
4. Click "New Post" button to navigate to /admin/blog/posts/new
5. Fill required fields: title, category_id, author_id
6. Use Trix editor for body content (input="body-input")
7. Blur title field to trigger auto-slug generation
8. Use HTMX tag picker to add tags
9. Use HTMX product picker to link products
10. Fill SEO fields and verify live preview updates
11. Upload featured image
12. Submit form via POST /admin/blog/posts
13. Edit existing post at /admin/blog/posts/:id/edit
14. Delete post using hx-delete with hx-target="closest tr"

## Test Cases

### Happy Path
- **List posts with filters**: Verify search by title, filter by status/category/author work correctly
- **Pagination**: Verify 15 posts per page and navigation works
- **Create new post**: Required fields (title, category, author) filled, slug auto-generated, post saved
- **Auto-slug generation**: Title "My Blog Post" generates slug "my-blog-post" on blur
- **Trix editor integration**: Body content entered via Trix, stored in hidden input
- **Tag picker HTMX**: Search tags, select existing, quick-create new, chips added to #selected-tags
- **Product picker HTMX**: Search products, select, chips added to #selected-products
- **SEO preview**: Meta description updates live as typed
- **Reading time calculation**: Auto-calculated at 200 words per minute from body word count
- **Edit post**: Navigate to edit form, modify fields, save successfully
- **Delete post**: hx-delete removes row from table without page reload

### Edge Cases / Error States
- **Missing required fields**: Title, category_id, author_id empty triggers validation errors
- **Excerpt maxlength**: 300 character limit enforced on excerpt textarea
- **Duplicate slug**: Manual slug entry that conflicts with existing post shows error
- **Empty slug auto-generation**: Blank title doesn't generate slug until filled
- **Published_at datetime**: Invalid datetime-local format rejected
- **Tag removal**: Clicking X on tag chip removes hidden tag_ids[] input
- **Product removal**: Clicking X on product chip removes hidden product_ids[] input
- **Delete confirmation**: Verify HTMX confirm dialog appears before deletion

## Selectors & Elements
- Form action: POST /admin/blog/posts (create), POST /admin/blog/posts/:id (update)
- Input names: title, slug, excerpt, body (hidden), category_id, author_id, featured_image_url, featured_image_alt, meta_description, reading_time_minutes, status, published_at, tag_ids[], product_ids[]
- Trix editor: input="body-input", <trix-editor input="body-input">
- Tag search: id="tag-search"
- Tag suggestions container: id="tag-suggestions"
- Selected tags container: id="selected-tags"
- Product search: id="product-search"
- Product suggestions container: id="product-suggestions"
- Selected products container: id="selected-products"
- Delete button: hx-delete="/admin/blog/posts/:id" hx-target="closest tr"
- Submit button: text "Create Post" or "Update Post"
- Status select: options "draft", "published"

## HTMX Interactions
- **Slug auto-generation**: hx-get="/admin/blog/posts/generate-slug" hx-trigger="blur from:#title" hx-target="#slug"
- **Tag search**: hx-get="/admin/blog/tags/search" hx-trigger="input changed delay:200ms, focus" hx-target="#tag-suggestions" (returns tag_suggestions.html partial)
- **Tag quick-create**: hx-post="/admin/blog/tags/quick-create" hx-vals='{"name":"..."}' hx-target="#selected-tags" hx-swap="beforeend" (returns tag_chip.html partial)
- **Product search**: hx-get="/admin/blog/products/search" hx-trigger="input changed delay:200ms, focus" hx-target="#product-suggestions" (returns product_suggestions.html partial)
- **Delete post**: hx-delete="/admin/blog/posts/:id" hx-target="closest tr" hx-swap="outerHTML"
- **SEO preview update**: JS listener on meta_description input updates preview div live

## Dependencies
- 16-blog-tags.md (tag search and quick-create)
- 17-blog-post-tags.md (tag assignment flow)
- 18-blog-post-products.md (product linking flow)
