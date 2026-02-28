# Test Plan: Product Features (HTMX Sub-Resource)

## Summary
Verify HTMX-based inline management of product features on product edit page with display ordering.

## Preconditions
- User authenticated with valid session cookie
- Product exists with ID for editing
- Navigate to product edit page at http://localhost:28090/admin/products/:id/edit
- Server running on localhost:28090

## User Journey Steps
1. On product edit page, verify HTMX features section loads (similar pattern to specs)
2. Verify container `#features-section` displays existing features in display_order
3. Fill add feature form: feature_text, display_order
4. Submit via `hx-post="/admin/products/:id/features"` with `hx-target="#features-section"` and `hx-swap="outerHTML"`
5. Verify new feature appears in list without page reload
6. Click "Delete All Features" button with `hx-delete="/admin/products/:id/features"` and `hx-confirm`
7. Confirm deletion dialog
8. Verify all features removed from section

## Test Cases

### Happy Path
- **Features section loads via HTMX**: On edit page load or tab click, features container populates
- **Add new feature**: Form submission adds feature, updates UI inline
- **Display order respected**: Features display in correct display_order
- **Delete all features**: Bulk delete removes all features for product after confirmation
- **Empty state**: When no features exist, section shows appropriate empty message

### Edge Cases / Error States
- **Required feature_text field**: Submitting without feature_text shows validation error
- **Delete all confirmation cancel**: Canceling hx-confirm does not delete features
- **Long feature text**: Very long feature_text displays correctly without breaking layout
- **HTML in feature text**: Feature text with HTML entities/tags is properly escaped or rendered
- **Duplicate display_order**: Multiple features with same display_order handled gracefully
- **Negative display_order**: System handles negative or zero display_order values

## Selectors & Elements
- Container: `id="features-section"`
- Load trigger: `hx-get="/admin/products/:id/features" hx-trigger="load"` (or similar)
- Add form: `hx-post="/admin/products/:id/features" hx-target="#features-section" hx-swap="outerHTML"`
- Textarea feature_text: `name="feature_text"` (required)
- Input display_order: `name="display_order" type="number"`
- Delete all button: `hx-delete="/admin/products/:id/features" hx-confirm="Delete all features for this product?"`
- Empty state message: displayed when no features exist

**Note**: Individual delete buttons appear in the template UI but the backend route (`/admin/products/:id/features/:feature_id`) is NOT implemented. This is a known bug. Only bulk deletion via `DELETE /admin/products/:id/features` is currently supported.
- Feature list items: ordered by display_order

## HTMX Interactions
- Initial load: `hx-get="/admin/products/:id/features"` populates `#features-section` (or loads with page)
- Add feature: `hx-post="/admin/products/:id/features"` with `hx-target="#features-section"` and `hx-swap="outerHTML"`
- Delete all: `hx-delete="/admin/products/:id/features"` with `hx-confirm` dialog
- Target: `#features-section` for full section replacement
- Swap: `outerHTML` replaces entire section container with updated HTML

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 09-products-crud.md (requires product edit page context)
