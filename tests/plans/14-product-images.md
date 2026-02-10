# Test Plan: Product Images (HTMX Sub-Resource)

## Summary
Verify HTMX-based inline management of product images on product edit page with file upload and thumbnail support.

## Preconditions
- User authenticated with valid session cookie
- Product exists with ID for editing
- Navigate to product edit page at http://localhost:28090/admin/products/:id/edit
- Server running on localhost:28090

## User Journey Steps
1. On product edit page, verify HTMX images section loads
2. Verify container `#images-section` displays existing images in display_order with thumbnails
3. Fill add image form: image (upload), alt_text, caption, is_thumbnail (checkbox), display_order
4. Submit via `hx-post="/admin/products/:id/images"` with `hx-encoding="multipart/form-data"`, `hx-target="#images-section"`, `hx-swap="outerHTML"`
5. Verify new image appears in gallery without page reload
6. Verify image thumbnail renders correctly
7. Click individual delete button on an image
8. Verify image removed from gallery via HTMX using `hx-delete="/admin/products/:id/images/:image_id"`
9. Test is_thumbnail checkbox functionality

## Test Cases

### Happy Path
- **Images section loads via HTMX**: On edit page load or tab click, images container populates with gallery
- **Add new image with file**: Multipart form submission uploads image and adds to gallery, updates UI inline
- **Display order respected**: Images display in correct display_order
- **Thumbnail display**: Uploaded images render as thumbnails in gallery view
- **Alt text and caption**: Alt text and caption metadata save and display correctly
- **Is_thumbnail checkbox**: Checkbox toggles thumbnail designation for product
- **Delete individual image**: Individual delete removes image and file without reload
- **Empty state**: When no images exist, section shows appropriate empty message or upload prompt

### Edge Cases / Error States
- **Required image field**: Submitting without image file shows validation error
- **Invalid file type**: Uploading non-image file (e.g., PDF, TXT) shows validation error
- **File size limit**: Uploading image exceeding size limit shows validation error
- **Delete confirmation**: Individual delete may have hx-confirm for user confirmation
- **Long alt text**: Very long alt_text saves correctly
- **Long caption**: Very long caption displays correctly without breaking layout
- **Missing optional fields**: Image can be uploaded without alt_text, caption (if optional)
- **Multiple thumbnails**: System handles or restricts multiple images marked as is_thumbnail
- **Special characters in filename**: Image files with special characters upload correctly
- **Large image dimensions**: Very large dimension images are handled (resized or accepted as-is)

## Selectors & Elements
- Container: `id="images-section"`
- Load trigger: `hx-get="/admin/products/:id/images" hx-trigger="load"` (or similar)
- Add form: `hx-post="/admin/products/:id/images" hx-encoding="multipart/form-data" hx-target="#images-section" hx-swap="outerHTML"`
- Input image: `name="image" type="file"` (required)
- Input alt_text: `name="alt_text" type="text"`
- Textarea caption: `name="caption"`
- Checkbox is_thumbnail: `name="is_thumbnail" type="checkbox"`
- Input display_order: `name="display_order" type="number"`
- Individual delete button: `hx-delete="/admin/products/:id/images/:image_id"` targeting specific image
- Image gallery: displays thumbnails of all product images
- Empty state message: displayed when no images exist
- Thumbnail indicator: visual indicator showing which image is marked as thumbnail

## HTMX Interactions
- Initial load: `hx-get="/admin/products/:id/images"` populates `#images-section`
- Add image: `hx-post="/admin/products/:id/images"` with `hx-encoding="multipart/form-data"`, `hx-target="#images-section"`, `hx-swap="outerHTML"`
- Delete individual: `hx-delete="/admin/products/:id/images/:image_id"` removes specific image
- Target: `#images-section` for full section replacement
- Swap: `outerHTML` replaces entire section container with updated HTML
- Encoding: `multipart/form-data` for file upload handling
- Note: No "delete all" button mentioned, only individual deletes

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 09-products-crud.md (requires product edit page context)
