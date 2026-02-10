# Test Plan: Blog Post Product Linking

## Summary
Testing the HTMX-driven product search, selection, and removal flow within blog post create/edit forms for linking related products.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with products and blog posts
- On blog post create/edit form at /admin/blog/posts/new or /admin/blog/posts/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/blog/posts/new
2. Locate product-search input in "Related Products" card
3. Type product name or SKU to trigger hx-get="/admin/blog/products/search"
4. Verify product suggestions appear in #product-suggestions after 200ms delay
5. Click a product suggestion to add product chip
6. Verify product chip with hidden product_ids[] input added to #selected-products
7. Click X on product chip to remove product
8. Verify hidden product_ids[] input removed from #selected-products
9. Submit post form with linked products

## Test Cases

### Happy Path
- **Search products**: Type "sensor", see matching products in dropdown within 200ms
- **Select product from suggestions**: Click "Temperature Sensor XL" suggestion, chip added to #selected-products with hidden input product_ids[]=:id
- **Multiple product selection**: Add 3 different products, verify 3 chips and 3 hidden inputs exist
- **Product chip display**: Chip shows product name and SKU/image if available
- **Remove product**: Click X on product chip, chip and hidden input removed from DOM
- **Focus trigger**: Focusing product-search input shows recent products or all products
- **Form submission**: Submit post form with 2 linked products, verify product_ids[] array sent to server

### Edge Cases / Error States
- **Input delay**: Typing rapidly waits 200ms after last keystroke before HTMX request
- **Empty search**: Focusing empty search input shows all products or placeholder
- **No results**: Search for "nonexistent123" shows "No products found" message
- **Duplicate product selection**: Selecting same product twice prevented by JS or shows warning
- **Remove last product**: Removing all products leaves #selected-products empty
- **Long product names**: Product chip with very long name truncates gracefully
- **JavaScript disabled**: Fallback behavior if addProduct() function unavailable

## Selectors & Elements
- Product search input: id="product-search", hx-get="/admin/blog/products/search", hx-trigger="input changed delay:200ms, focus", hx-target="#product-suggestions"
- Suggestions container: id="product-suggestions" (receives product_suggestions.html partial)
- Selected products container: id="selected-products" (contains product chips)
- Product chip structure: <div class="product-chip">Product Name <button onclick="removeProduct(id)">X</button><input type="hidden" name="product_ids[]" value=":id"></div>
- JavaScript functions: addProduct(id, name), removeProduct(id)

## HTMX Interactions
- **Product search**: hx-get="/admin/blog/products/search?q=keyword" hx-trigger="input changed delay:200ms, focus" hx-target="#product-suggestions" (returns product_suggestions.html partial with clickable suggestions)
- **Chip addition**: Clicking suggestion triggers addProduct() JS which manually creates chip HTML and appends to #selected-products
- **Chip removal**: Clicking X button calls removeProduct() JS which removes chip element from DOM

## Dependencies
- Products table must be seeded with data
- 15-blog-posts-crud.md (parent blog post form)
