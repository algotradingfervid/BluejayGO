# Test Plan: Partners CRUD

## Summary
Testing partner creation, editing, listing with in-memory filters, deletion, and cache invalidation with UNIQUE name constraint.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 11 partners and partner_tiers

## User Journey Steps
1. Navigate to http://localhost:28090/admin/partners
2. Verify partner list displays with search, tier, status filters
3. Verify no pagination (in-memory filtering, all 11 partners shown)
4. Click "New Partner" button to navigate to /admin/partners/new
5. Fill required field: name (UNIQUE constraint)
6. Select tier_id from dropdown (partner_tiers)
7. Fill optional fields: logo_url, icon, website_url, description, display_order
8. Set is_active checkbox
9. Submit form via POST /admin/partners
10. Edit existing partner at /admin/partners/:id/edit
11. Delete partner using hx-delete
12. Verify cache invalidation for page:partners

## Test Cases

### Happy Path
- **List partners with filters**: Verify search by name, filter by tier/status work correctly
- **In-memory filtering**: All 11 partners loaded, filters applied client-side or server-side without pagination
- **Create new partner**: Required name filled, tier selected, partner saved
- **Unique name constraint**: Name "Acme Corp" saved, attempting duplicate shows error
- **Tier selection**: Select tier from dropdown, association saved
- **Logo and icon**: Fill logo_url and icon (Material icon name), saved successfully
- **Website URL**: Fill website_url "https://example.com", saved successfully
- **Description**: Fill description textarea with partner details
- **Display order**: Set numeric display_order for partner positioning
- **Active checkbox**: Toggle is_active checkbox, partner visibility changes
- **Edit partner**: Navigate to edit form, modify fields, save successfully
- **Delete partner**: hx-delete removes row from table without page reload
- **Cache invalidation**: After create/update/delete, page:partners cache cleared

### Edge Cases / Error States
- **Missing required name**: Empty name triggers validation error
- **Duplicate name**: Creating partner with existing name violates UNIQUE constraint, shows error
- **No tier selected**: tier_id null may show validation error or be accepted as optional
- **Invalid URL format**: website_url "not-a-url" may show validation warning or be accepted
- **Long description**: Description with 500+ characters accepted in textarea
- **Negative display_order**: Entering negative number validated or accepted
- **All filters applied**: Search + tier + status filters combined show subset of partners
- **Delete partner with testimonials**: Deleting partner with associated testimonials may cascade delete or show warning

## Selectors & Elements
- List route: GET /admin/partners
- Create form action: POST /admin/partners
- Edit route: GET /admin/partners/:id/edit
- Input names: name (required, UNIQUE), tier_id (select), logo_url, icon, website_url, description (textarea), display_order (number), is_active (checkbox)
- Delete button: hx-delete="/admin/partners/:id" hx-target="closest tr"
- Submit button: text "Create Partner" or "Update Partner"
- Filter inputs: search (text), tier (select), status (select active/inactive)

## HTMX Interactions
- **Delete partner**: hx-delete="/admin/partners/:id" hx-target="closest tr" hx-swap="outerHTML"
- Filters may use HTMX to update list without page reload, or standard form submission

## Dependencies
- 30-partner-testimonials.md (testimonials reference partners)
- Partner_tiers table must be seeded with data
