# Test Plan: Whitepapers CRUD

## Summary
Testing whitepaper creation with PDF upload, editing, listing with filters, deletion, and downloads analytics view with cache invalidation.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 12 whitepapers and whitepaper_topics

## User Journey Steps
1. Navigate to http://localhost:28090/admin/whitepapers
2. Verify whitepaper list displays with search, status, topic filters, pagination
3. Verify pagination shows 15 whitepapers per page
4. Click "New Whitepaper" button to navigate to /admin/whitepapers/new
5. Fill required field: title
6. Auto-generate slug on title blur
7. Fill description textarea
8. Select topic_id from dropdown (whitepaper_topics)
9. Fill published_date, cover_color_from, cover_color_to (gradient colors)
10. Fill meta_description
11. Set is_published checkbox
12. Enter page_count number
13. Upload PDF file via pdf_file input (multipart form)
14. Add learning_points (multiple text inputs)
15. Submit form via POST /admin/whitepapers (multipart/form-data)
16. Edit existing whitepaper at /admin/whitepapers/:id/edit
17. View downloads analytics at /admin/whitepapers/:id/downloads
18. Delete whitepaper using hx-delete
19. Verify cache invalidation for page:whitepapers

## Test Cases

### Happy Path
- **List whitepapers with filters**: Verify search by title, filter by status/topic work correctly
- **Pagination**: Verify 15 whitepapers per page and navigation works
- **Create new whitepaper**: Required title filled, slug auto-generated, PDF uploaded, whitepaper saved
- **Auto-slug generation**: Title "IoT Best Practices" generates slug "iot-best-practices" on blur
- **Topic selection**: Select topic from dropdown, association saved
- **PDF upload**: Upload test.pdf file, file stored, path saved to pdf_file column
- **Cover gradient**: Fill cover_color_from "#3B82F6", cover_color_to "#8B5CF6", gradient colors saved
- **Learning points**: Add 5 learning points via multiple text inputs, stored as array
- **Page count**: Enter page_count 24, saved successfully
- **Edit whitepaper**: Navigate to edit form, modify fields, re-upload PDF, save successfully
- **Downloads view**: Navigate to /admin/whitepapers/:id/downloads, see analytics data
- **Delete whitepaper**: hx-delete removes row from table without page reload
- **Cache invalidation**: After create/update/delete, page:whitepapers cache cleared

### Edge Cases / Error States
- **Missing required title**: Empty title triggers validation error
- **Duplicate slug**: Manual slug entry that conflicts shows error
- **No PDF uploaded**: Creating whitepaper without PDF may show warning or accept null
- **Invalid PDF file**: Uploading non-PDF file rejected or shows error
- **Large PDF file**: Uploading 50MB+ PDF may hit upload size limit
- **Empty learning_points**: No learning points entered, stored as empty array
- **Invalid color format**: cover_color_from "invalid" rejected, requires hex color
- **No topic selected**: topic_id null accepted, optional field
- **Negative page_count**: Entering negative number validated or accepted as 0

## Selectors & Elements
- List route: GET /admin/whitepapers
- Create form action: POST /admin/whitepapers (enctype="multipart/form-data")
- Edit route: GET /admin/whitepapers/:id/edit
- Downloads route: GET /admin/whitepapers/:id/downloads
- Input names: title, slug, description (textarea), topic_id (select), published_date, cover_color_from, cover_color_to, meta_description, is_published (checkbox), page_count (number), pdf_file (file), learning_points[] (multiple text inputs)
- Delete button: hx-delete="/admin/whitepapers/:id" hx-target="closest tr"
- Submit button: text "Create Whitepaper" or "Update Whitepaper"

## HTMX Interactions
- **Slug auto-generation**: hx-get="/admin/whitepapers/generate-slug" hx-trigger="blur from:#title" hx-target="#slug"
- **Delete whitepaper**: hx-delete="/admin/whitepapers/:id" hx-target="closest tr" hx-swap="outerHTML"
- Downloads view is separate page, no HTMX interactions

## Dependencies
- 28-whitepaper-downloads.md (downloads analytics view)
- Whitepaper_topics table must be seeded with data
- File upload storage configuration required
