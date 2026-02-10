# Test Plan: Solution Products HTMX Management

## Summary
Testing HTMX-driven addition and removal of product associations to solutions with UNIQUE constraint enforcement.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with solutions and products
- On solution edit page at /admin/solutions/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/solutions/:id/edit
2. Locate #products-section container
3. Select product from product_id dropdown (select from all products)
4. Set display_order and is_featured checkbox
5. Click add button with hx-post="/admin/solutions/:id/products"
6. Verify hx-target="#products-section" hx-swap="outerHTML" replaces entire section
7. Verify new product association appears in updated section
8. Click remove button on existing product with hx-delete="/admin/solutions/:id/products/:productId"
9. Verify product removed from #products-section without page reload

## Test Cases

### Happy Path
- **Add product**: Select product from dropdown, set display_order 1, is_featured true, submit, product added
- **Multiple products**: Add 3 different products with different display_order values, verify all appear
- **Display order**: Products display in order based on display_order field
- **Featured flag**: is_featured checkbox sets featured status for product on solution
- **Remove product**: Click remove button on product, hx-delete removes it from section
- **Section swap**: After add/remove, entire #products-section swapped with updated HTML
- **No page reload**: All operations via HTMX, no full page refresh

### Edge Cases / Error States
- **Duplicate product**: Adding same product_id twice violates UNIQUE constraint, shows error via HTMX response
- **Missing product_id**: Not selecting product from dropdown triggers validation error
- **All products already added**: Dropdown shows only unassociated products or disabled if all added
- **Display order conflict**: Multiple products with same display_order accepted, sorted arbitrarily
- **Remove with confirmation**: hx-delete may include hx-confirm="Remove this product?" attribute
- **HTMX error handling**: Server error on add/remove shows error message in section

## Selectors & Elements
- Section container: id="products-section"
- Add form action: hx-post="/admin/solutions/:id/products" hx-target="#products-section" hx-swap="outerHTML"
- Input names: product_id (select required), display_order (number), is_featured (checkbox)
- Remove button: hx-delete="/admin/solutions/:id/products/:productId" hx-target="#products-section" hx-swap="outerHTML" hx-confirm="Remove this product?"
- Add button: text "Add Product"

## HTMX Interactions
- **Add product**: hx-post="/admin/solutions/:id/products" hx-target="#products-section" hx-swap="outerHTML" (returns updated products_section.html partial)
- **Remove product**: hx-delete="/admin/solutions/:id/products/:productId" hx-target="#products-section" hx-swap="outerHTML" hx-confirm="Remove this product?" (returns updated products_section.html partial)
- Both operations swap entire #products-section to reflect current state

## Dependencies
- 19-solutions-crud.md (parent solution edit page)
- Products table must be seeded with data
