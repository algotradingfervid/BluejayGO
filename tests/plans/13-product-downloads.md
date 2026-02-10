# Test Plan: Product Downloads (HTMX Sub-Resource)

## Summary
Verify HTMX-based inline management of product downloads on product edit page with file upload support.

## Preconditions
- User authenticated with valid session cookie
- Product exists with ID for editing
- Navigate to product edit page at http://localhost:28090/admin/products/:id/edit
- Server running on localhost:28090

## User Journey Steps
1. On product edit page, verify HTMX downloads section loads
2. Verify container `#downloads-section` displays existing downloads in display_order
3. Fill add download form: file (upload), title, description, file_type, version, display_order
4. Submit via `hx-post="/admin/products/:id/downloads"` with `hx-encoding="multipart/form-data"`, `hx-target="#downloads-section"`, `hx-swap="outerHTML"`
5. Verify new download appears in list with file link without page reload
6. Click individual delete button on a download
7. Verify download removed from list via HTMX using `hx-delete="/admin/products/:id/downloads/:download_id"`
8. Verify file metadata (title, size, type) displays correctly

## Test Cases

### Happy Path
- **Downloads section loads via HTMX**: On edit page load or tab click, downloads container populates
- **Add new download with file**: Multipart form submission uploads file and adds download, updates UI inline
- **Display order respected**: Downloads display in correct display_order
- **File metadata display**: Title, file_type, version, and file size display correctly
- **Delete individual download**: Individual delete removes download and file without reload
- **Download link functional**: Click download link retrieves uploaded file
- **Empty state**: When no downloads exist, section shows appropriate empty message

### Edge Cases / Error States
- **Required file field**: Submitting without file shows validation error
- **Required title field**: Submitting without title shows validation error
- **File size limit**: Uploading file exceeding size limit shows validation error
- **Invalid file type**: Uploading disallowed file type shows validation error
- **Delete confirmation**: Individual delete may have hx-confirm for user confirmation
- **Long description**: Very long description displays correctly without breaking layout
- **Missing optional fields**: Download can be created without description, file_type, version (if optional)
- **Special characters in filename**: Files with special characters in name upload and display correctly
- **Duplicate filenames**: System handles multiple downloads with same filename

## Selectors & Elements
- Container: `id="downloads-section"`
- Load trigger: `hx-get="/admin/products/:id/downloads" hx-trigger="load"` (or similar)
- Add form: `hx-post="/admin/products/:id/downloads" hx-encoding="multipart/form-data" hx-target="#downloads-section" hx-swap="outerHTML"`
- Input file: `name="file" type="file"` (required)
- Input title: `name="title" type="text"` (required)
- Textarea description: `name="description"`
- Input file_type: `name="file_type" type="text"`
- Input version: `name="version" type="text"`
- Input display_order: `name="display_order" type="number"`
- Individual delete button: `hx-delete="/admin/products/:id/downloads/:download_id"` targeting specific download
- Download link: anchor to download file
- Empty state message: displayed when no downloads exist
- File metadata display: shows title, size, type, version

## HTMX Interactions
- Initial load: `hx-get="/admin/products/:id/downloads"` populates `#downloads-section`
- Add download: `hx-post="/admin/products/:id/downloads"` with `hx-encoding="multipart/form-data"`, `hx-target="#downloads-section"`, `hx-swap="outerHTML"`
- Delete individual: `hx-delete="/admin/products/:id/downloads/:download_id"` removes specific download
- Target: `#downloads-section` for full section replacement
- Swap: `outerHTML` replaces entire section container with updated HTML
- Encoding: `multipart/form-data` for file upload handling
- Note: No "delete all" button mentioned, only individual deletes

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 09-products-crud.md (requires product edit page context)
