# Test Plan: Products CRUD

## Summary
Verify complete CRUD operations for products including search, filtering, pagination, file uploads, and HTMX tabbed sub-resources.

## Preconditions
- User authenticated with valid session cookie
- Database seeded with 13 products and 5 product categories
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/products
2. Verify list shows products with pagination (15 per page)
3. Test search by product name/SKU
4. Test filter by status (draft/published/archived)
5. Test filter by category
6. Click "New Product" or navigate to http://localhost:28090/admin/products/new
7. Fill form with multipart/form-data: sku, name, tagline, description, overview (Trix), category_id, status, is_featured, featured_order, primary_image (file), video_url, meta fields
8. Submit POST to /admin/products
9. Verify redirect to /admin/products with new product in list
10. Click "Edit" on product or navigate to http://localhost:28090/admin/products/:id/edit
11. Verify HTMX tabs load sub-resource sections (specs, features, certifications, downloads, images)
12. Modify product fields and submit POST to /admin/products/:id
13. Verify updates appear in list
14. Click delete button with hx-delete on product row
15. Confirm deletion dialog and verify row removed via HTMX

## Test Cases

### Happy Path
- **List products with pagination**: 13 products display across pages, 15 per page
- **Search products**: Searching by name or SKU filters list correctly
- **Filter by status**: Dropdown filters show only draft/published/archived products
- **Filter by category**: Category dropdown filters products correctly
- **Create product with file upload**: Multipart form creates product with primary_image file
- **Edit product**: Form pre-fills all fields including Trix editor content
- **HTMX tabs load sub-resources**: Clicking specs tab triggers hx-get with hx-trigger="load"
- **Delete product via HTMX**: Delete button removes row without page reload
- **Auto-slug generation**: Slug automatically generated from name
- **Featured checkbox**: is_featured checkbox toggles correctly
- **Trix editor**: Overview field uses Trix rich text editor

### Edge Cases / Error States
- **Duplicate SKU validation**: Creating product with existing SKU shows UNIQUE constraint error
- **Required name field**: Submitting without name shows validation error
- **Required SKU field**: Submitting without SKU shows validation error
- **Invalid category_id**: Submitting with non-existent category shows foreign key error
- **Invalid status value**: Submitting invalid status shows validation error
- **File upload size limit**: Uploading very large image shows size limit error
- **Invalid file type**: Uploading non-image file shows validation error
- **Delete confirmation cancel**: Canceling hx-confirm does not delete product
- **Empty search**: Searching with empty query returns all products
- **Pagination edge cases**: First page, last page, single product navigation
- **Missing primary_image**: Product can be created without image (optional)

## Selectors & Elements
- List page: http://localhost:28090/admin/products with query params for search/filter
- Create form: `action="/admin/products" method="POST" enctype="multipart/form-data"`
- Edit form: `action="/admin/products/:id" method="POST" enctype="multipart/form-data"`
- Input sku: `name="sku" type="text"` (required, UNIQUE)
- Input name: `name="name" type="text"` (required)
- Input slug: `name="slug" type="text"` (auto-generated)
- Input tagline: `name="tagline" type="text"`
- Textarea description: `name="description"`
- Trix editor: `name="overview"` with Trix editor component
- Select category_id: `name="category_id"` with options from product_categories
- Select status: `name="status"` with options (draft, published, archived)
- Checkbox is_featured: `name="is_featured" type="checkbox"`
- Input featured_order: `name="featured_order" type="number"`
- Input primary_image: `name="primary_image" type="file"`
- Input video_url: `name="video_url" type="url"`
- Input meta_title: `name="meta_title" type="text"`
- Textarea meta_description: `name="meta_description"`
- Search input: filter products by name/SKU
- Status filter: dropdown with draft/published/archived options
- Category filter: dropdown with category options
- Pagination controls: links for page navigation
- Delete button: `hx-delete="/admin/products/:id" hx-target="closest tr" hx-confirm="Delete this product?"`
- HTMX tabs container: sections for specs, features, certifications, downloads, images
- Specs tab: `hx-get="/admin/products/:id/specs" hx-trigger="load"` (or click)

## HTMX Interactions
- Delete product: `hx-delete="/admin/products/:id"` with `hx-target="closest tr"` and `hx-confirm`
- Specs section load: `hx-get="/admin/products/:id/specs"` with `hx-trigger="load"` in edit page
- Features section: `hx-get="/admin/products/:id/features"` (similar pattern)
- Certifications section: `hx-get="/admin/products/:id/certifications"`
- Downloads section: `hx-get="/admin/products/:id/downloads"`
- Images section: `hx-get="/admin/products/:id/images"`
- Target: `#specs-section`, `#features-section`, etc. for tab content swap
- Swap: `hx-swap="innerHTML"` or `hx-swap="outerHTML"` for tab content

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 03-product-categories-crud.md (categories used in product form)
- 10-product-specs.md (specs sub-resource in edit page)
- 11-product-features.md (features sub-resource in edit page)
- 12-product-certifications.md (certifications sub-resource in edit page)
- 13-product-downloads.md (downloads sub-resource in edit page)
- 14-product-images.md (images sub-resource in edit page)
