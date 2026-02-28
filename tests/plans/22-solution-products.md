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
2. Locate #products-list container
3. Enter product_id in number input field (not a dropdown select)
4. Set display_order and is_featured checkbox
5. Click add button with hx-post="/admin/solutions/:id/products"
6. Verify hx-target="#products-list" hx-swap="innerHTML" adds new product to list
7. Verify new product association appears in list
8. Click remove button on existing product with hx-delete="/admin/solutions/:id/products/:productId"
9. Verify product removed (hx-target="closest div" hx-swap="outerHTML") without page reload

## Test Cases

### Happy Path
- **Add product**: Enter product_id in number input, set display_order 1, is_featured true, submit, product added (innerHTML swap adds partial HTML)
- **Multiple products**: Add 3 different products with different display_order values, verify all appear
- **Display order**: Products display in order based on display_order field
- **Featured flag**: is_featured checkbox sets featured status for product on solution
- **Remove product**: Click remove button on product, hx-delete removes individual div (returns c.NoContent(http.StatusOK) with no HTML)
- **Add returns HTML**: Add handler returns partial HTML
- **Delete returns no content**: RemoveProduct returns c.NoContent(http.StatusOK) with no HTML
- **No page reload**: All operations via HTMX, no full page refresh

### Edge Cases / Error States
- **Duplicate product**: Adding same product_id twice violates UNIQUE constraint, shows error via HTMX response
- **Missing product_id**: Not entering product_id in number input triggers validation error
- **Invalid product_id**: Entering non-existent product_id may show error or foreign key constraint violation
- **Display order conflict**: Multiple products with same display_order accepted, sorted arbitrarily
- **Remove with confirmation**: hx-delete may include hx-confirm="Remove this product?" attribute
- **HTMX error handling**: Server error on add/remove shows error message in section

## Selectors & Elements
- Section container: id="products-list"
- Add form action: hx-post="/admin/solutions/:id/products" hx-target="#products-list" hx-swap="innerHTML"
- Input names: product_id (type="number" required), display_order (number), is_featured (checkbox)
- Remove button: hx-delete="/admin/solutions/:id/products/:productId" hx-target="closest div" hx-swap="outerHTML" hx-confirm="Remove this product?"
- Add button: text "Add Product"

## HTMX Interactions
- **Add product**: hx-post="/admin/solutions/:id/products" hx-target="#products-list" hx-swap="innerHTML" (returns partial HTML)
- **Remove product**: hx-delete="/admin/solutions/:id/products/:productId" hx-target="closest div" hx-swap="outerHTML" hx-confirm="Remove this product?" (returns c.NoContent(http.StatusOK) with no HTML content)
- Add operation inserts new product HTML into list, remove operation removes individual item div

## Dependencies
- 19-solutions-crud.md (parent solution edit page)
- Products table must be seeded with data
