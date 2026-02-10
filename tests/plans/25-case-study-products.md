# Test Plan: Case Study Products HTMX Management

## Summary
Testing HTMX-driven addition and removal of product associations to case studies with UNIQUE constraint enforcement and confirmation dialog.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with case studies and products
- On case study edit page at /admin/case-studies/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/case-studies/:id/edit
2. Locate #products-section container
3. Select product from product_id dropdown (select from all products)
4. Set display_order number
5. Click add button with hx-post="/admin/case-studies/:id/products"
6. Verify hx-target="#products-section" hx-swap="outerHTML" replaces entire section
7. Verify new product association appears in updated section
8. Click remove button on existing product with hx-delete="/admin/case-studies/:id/products/:productId"
9. Verify hx-confirm="Remove this product?" dialog appears
10. Confirm deletion, verify product removed from #products-section without page reload

## Test Cases

### Happy Path
- **Add product**: Select product from dropdown, set display_order 1, submit, product added
- **Multiple products**: Add 3 different products with different display_order values, verify all appear
- **Display order**: Products display in order based on display_order field
- **Remove product**: Click remove button, confirm dialog appears, confirm, product removed from section
- **Section swap**: After add/remove, entire #products-section swapped with updated HTML
- **No page reload**: All operations via HTMX, no full page refresh

### Edge Cases / Error States
- **Duplicate product**: Adding same product_id twice violates UNIQUE constraint (solution_id, product_id), shows error via HTMX response
- **Missing product_id**: Not selecting product from dropdown triggers validation error
- **All products already added**: Dropdown shows only unassociated products or disabled if all added
- **Display order conflict**: Multiple products with same display_order accepted, sorted arbitrarily
- **Cancel confirm dialog**: Clicking "Cancel" on hx-confirm dialog prevents deletion
- **HTMX error handling**: Server error on add/remove shows error message in section

## Selectors & Elements
- Section container: id="products-section"
- Add form action: hx-post="/admin/case-studies/:id/products" hx-target="#products-section" hx-swap="outerHTML"
- Input names: product_id (select required), display_order (number)
- Remove button: hx-delete="/admin/case-studies/:id/products/:productId" hx-target="#products-section" hx-swap="outerHTML" hx-confirm="Remove this product?"
- Add button: text "Add Product"

## HTMX Interactions
- **Add product**: hx-post="/admin/case-studies/:id/products" hx-target="#products-section" hx-swap="outerHTML" (returns updated products_section.html partial)
- **Remove product**: hx-delete="/admin/case-studies/:id/products/:productId" hx-target="#products-section" hx-swap="outerHTML" hx-confirm="Remove this product?" (returns updated products_section.html partial)
- Both operations swap entire #products-section to reflect current state

## Dependencies
- 24-case-studies-crud.md (parent case study edit page)
- Products table must be seeded with data
