# Test Plan: Product Specifications (HTMX Sub-Resource)

## Summary
Verify HTMX-based inline management of product specifications on product edit page with grouped sections.

## Preconditions
- User authenticated with valid session cookie
- Product exists with ID for editing
- Navigate to product edit page at http://localhost:28090/admin/products/:id/edit
- Server running on localhost:28090

## User Journey Steps
1. On product edit page, verify HTMX specs section loads via `hx-get="/admin/products/:id/specs"` with `hx-trigger="load"`
2. Verify container `#specs-section` displays existing specs grouped by section_name
3. Verify collapsible panels for each section_name group
4. Fill add spec form: section_name, spec_key, spec_value, display_order
5. Submit via `hx-post="/admin/products/:id/specs"` with `hx-target="#specs-section"` and `hx-swap="outerHTML"`
6. Verify new spec appears in correct section group without page reload
7. Click individual delete button on a spec
8. Verify spec removed from list via HTMX
9. Click "Delete All Specs" button with `hx-delete="/admin/products/:id/specs"` and `hx-confirm`
10. Confirm deletion dialog
11. Verify all specs removed from section

## Test Cases

### Happy Path
- **Specs section loads via HTMX**: On edit page load, specs container populates with existing specs
- **Add new spec**: Form submission adds spec to correct section group, updates UI inline
- **Specs grouped by section**: Specs with same section_name display together in collapsible panel
- **Display order respected**: Specs within section display in correct display_order
- **Delete individual spec**: Individual delete removes single spec without reload
- **Delete all specs**: Bulk delete removes all specs for product after confirmation
- **Empty state**: When no specs exist, section shows appropriate empty message

### Edge Cases / Error States
- **Required fields validation**: Submitting without section_name, spec_key, or spec_value shows validation error
- **Duplicate spec_key in section**: System allows or prevents duplicate keys based on business rules
- **Delete all confirmation cancel**: Canceling hx-confirm does not delete specs
- **Long spec values**: Very long spec_value text displays correctly without breaking layout
- **Special characters in keys**: Spec keys with special characters save and display correctly
- **Negative display_order**: System handles negative or zero display_order values
- **Collapsible panel state**: Expanding/collapsing section panels maintains state during inline operations

## Selectors & Elements
- Container: `id="specs-section"`
- Load trigger: `hx-get="/admin/products/:id/specs" hx-trigger="load"`
- Add form: `hx-post="/admin/products/:id/specs" hx-target="#specs-section" hx-swap="outerHTML"`
- Input section_name: `name="section_name" type="text"` (required)
- Input spec_key: `name="spec_key" type="text"` (required)
- Input spec_value: `name="spec_value" type="text"` (required)
- Input display_order: `name="display_order" type="number"`
- Delete all button: `hx-delete="/admin/products/:id/specs" hx-confirm="Delete all specs for this product?"`
- Individual delete button: `hx-delete="/admin/products/:id/specs/:spec_id"` (or similar endpoint)
- Section groups: collapsible panels grouped by section_name with toggle UI
- Empty state message: displayed when no specs exist

## HTMX Interactions
- Initial load: `hx-get="/admin/products/:id/specs"` with `hx-trigger="load"` populates `#specs-section`
- Add spec: `hx-post="/admin/products/:id/specs"` with `hx-target="#specs-section"` and `hx-swap="outerHTML"`
- Delete all: `hx-delete="/admin/products/:id/specs"` with `hx-confirm` dialog
- Delete individual: `hx-delete` on individual spec removes from list
- Target: `#specs-section` for full section replacement
- Swap: `outerHTML` replaces entire section container with updated HTML

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 09-products-crud.md (requires product edit page context)
