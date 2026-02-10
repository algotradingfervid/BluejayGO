# Phase 07 - Products List & Form

## Current State

### List Page
- Simple table: SKU, Name, Status, Featured, Actions
- No search, no filtering, no pagination
- HTMX delete with confirmation
- Basic status badges

### Form Page
- Long single-column form with many fields
- HTMX tabs for sub-details (specs, features, certifications, downloads, images)
- Auto-slug generation
- File upload for primary image

## Changes to List Page

### Add Filter Bar
- Search input (searches name and SKU)
  - Tooltip: "Search products by name or SKU number."
- Status dropdown filter (All / Published / Draft / Archived)
- Category dropdown filter
- "Clear filters" button when any filter is active

### Improve Table
- Add product thumbnail column (small image, 48x48)
- Add category column
- Add "Updated" date column
- Status badges: use brutalist style (manual-border, color-coded)
- Pagination: show 15 per page, "Showing 1-15 of X" + page numbers

### Empty State
- When no products exist: illustration + "No products yet" + "Create your first product" button
- When filters return nothing: "No products match your filters" + "Clear filters" button

## Changes to Form Page

### Collapsible Sections (instead of one long form)
Replace the single scrolling form with collapsible sections that the user can open/close:

**Section 1: Basic Information** (open by default)
- Product Name
  - Tooltip: "The product name as it appears on your site."
- SKU
  - Tooltip: "Unique stock-keeping unit code for this product."
- Slug (auto-generated, editable)
  - Tooltip: "URL-friendly version of the name. Auto-generated but editable."
- Category dropdown
- Status dropdown

**Section 2: Description** (open by default)
- Tagline (short text)
  - Tooltip: "A one-line summary shown on product cards. Keep under 100 characters."
- Description (textarea)
- Overview (Trix rich text editor)
  - Tooltip: "Detailed product description with formatting. Shown on the product detail page."

**Section 3: Media** (collapsed by default)
- Primary Image upload with preview
- Video URL
  - Tooltip: "YouTube or Vimeo embed URL for the product video."

**Section 4: Display Options** (collapsed by default)
- Is Featured (toggle)
  - Tooltip: "Featured products appear in the homepage featured section."
- Featured Order (number, shown only if featured)

**Section 5: SEO** (collapsed by default)
- Meta Title (70-char counter)
- Meta Description (160-char counter)

### Collapsible Section Design
- Black header bar with section title + chevron icon
- Click header to expand/collapse
- Smooth animation
- Open sections show form fields below header
- Visual indicator if section has validation errors (red dot on header)

### Remove from this form
- Specs, Features, Certifications, Downloads, Images tabs stay as HTMX sub-pages (Phase 08)
- They appear BELOW the main form only when editing (not when creating new)

## Backend Changes
- Add pagination to `ListProducts` query (LIMIT/OFFSET)
- Add filter parameters to list handler (search, status, category)
- Count query for total products

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/products_list.html` | Rewrite |
| `templates/admin/pages/products_form.html` | Rewrite with collapsible sections |
| `internal/handlers/admin/products.go` | Add filtering, pagination |
| `db/queries/products.sql` | Add filtered list query |

## Dependencies
- Phase 01, 02
