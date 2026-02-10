# Test Plan: Product Certifications (HTMX Sub-Resource)

## Summary
Verify HTMX-based inline management of product certifications on product edit page with icon support.

## Preconditions
- User authenticated with valid session cookie
- Product exists with ID for editing
- Navigate to product edit page at http://localhost:28090/admin/products/:id/edit
- Server running on localhost:28090

## User Journey Steps
1. On product edit page, verify HTMX certifications section loads
2. Verify container `#certifications-section` displays existing certifications in display_order
3. Fill add certification form: certification_name, certification_code, icon_type, icon_path, display_order
4. Submit via `hx-post="/admin/products/:id/certifications"` with `hx-target="#certifications-section"` and `hx-swap="outerHTML"`
5. Verify new certification appears in list with icon without page reload
6. Click individual delete button on a certification
7. Verify certification removed from list via HTMX
8. Click "Delete All Certifications" button with `hx-delete="/admin/products/:id/certifications"` and `hx-confirm`
9. Confirm deletion dialog
10. Verify all certifications removed from section

## Test Cases

### Happy Path
- **Certifications section loads via HTMX**: On edit page load or tab click, certifications container populates
- **Add new certification**: Form submission adds certification with icon, updates UI inline
- **Display order respected**: Certifications display in correct display_order
- **Icon display**: Icon_type and icon_path render correct icon in list
- **Delete individual certification**: Individual delete removes single certification without reload
- **Delete all certifications**: Bulk delete removes all certifications for product after confirmation
- **Empty state**: When no certifications exist, section shows appropriate empty message

### Edge Cases / Error States
- **Required certification_name field**: Submitting without certification_name shows validation error
- **Required certification_code field**: Submitting without certification_code shows validation error
- **Delete all confirmation cancel**: Canceling hx-confirm does not delete certifications
- **Long certification name**: Very long names display correctly without breaking layout
- **Missing icon fields**: Certification can be created without icon_type or icon_path (optional)
- **Invalid icon_path**: System handles invalid or missing icon paths gracefully
- **Duplicate certification_code**: System allows or prevents duplicate codes based on business rules

## Selectors & Elements
- Container: `id="certifications-section"`
- Load trigger: `hx-get="/admin/products/:id/certifications" hx-trigger="load"` (or similar)
- Add form: `hx-post="/admin/products/:id/certifications" hx-target="#certifications-section" hx-swap="outerHTML"`
- Input certification_name: `name="certification_name" type="text"` (required)
- Input certification_code: `name="certification_code" type="text"` (required)
- Input icon_type: `name="icon_type" type="text"`
- Input icon_path: `name="icon_path" type="text"`
- Input display_order: `name="display_order" type="number"`
- Delete all button: `hx-delete="/admin/products/:id/certifications" hx-confirm="Delete all certifications for this product?"`
- Individual delete button: `hx-delete="/admin/products/:id/certifications/:cert_id"` targeting specific certification
- Empty state message: displayed when no certifications exist
- Icon element: displays based on icon_type and icon_path

## HTMX Interactions
- Initial load: `hx-get="/admin/products/:id/certifications"` populates `#certifications-section`
- Add certification: `hx-post="/admin/products/:id/certifications"` with `hx-target="#certifications-section"` and `hx-swap="outerHTML"`
- Delete all: `hx-delete="/admin/products/:id/certifications"` with `hx-confirm` dialog
- Delete individual: `hx-delete` on individual certification removes from list
- Target: `#certifications-section` for full section replacement
- Swap: `outerHTML` replaces entire section container with updated HTML

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 09-products-crud.md (requires product edit page context)
