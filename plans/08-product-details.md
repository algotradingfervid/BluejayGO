# Phase 08 - Product Detail Sub-pages

## Current State
- HTMX tabs below product form for: Specs, Features, Certifications, Downloads, Images
- Each tab loads content via HTMX GET
- Basic forms within each tab

## Goal
Polish these sub-pages with better UX, clearer instructions, and consistent brutalist styling.

## Sub-page Improvements

### Specs Tab
- Group specs by category (Motor, Power, Physical, Environmental)
- Each group is collapsible
- "Add Spec" button within each group
- Each spec row: Label input + Value input + Delete button
- Tooltip on "Add Spec": "Add a technical specification. Use Label for the spec name and Value for the measurement."
- Drag handles for reordering within groups (stretch goal)

### Features Tab
- Simple list of feature text items
- Each row: Feature text input + Delete button
- "Add Feature" button at bottom
- Tooltip: "Product features shown as bullet points on the product page."
- Max 10 features recommended (show count: "3 of 10 recommended")

### Certifications Tab
- Each row: Certification Name + Issuing Body + Certificate Number (optional) + Delete
- "Add Certification" button
- Tooltip: "Industry certifications this product holds. Displayed with certification badges."

### Downloads Tab
- Each row: File name + File upload + Description + Delete
- "Add Download" button
- Show file size after upload
- Tooltip: "Downloadable files like datasheets, manuals, or CAD drawings. PDF format recommended."

### Images Tab
- Grid of uploaded images (3 columns)
- First image marked as "Primary" with badge
- Click to set as primary
- Delete button overlay on hover
- "Add Images" button (multi-file upload)
- Tooltip: "Product gallery images. First image is used as the main product photo."
- Drag to reorder (stretch goal)

## Shared Sub-page Pattern
All sub-pages follow this pattern:
1. Section title + tooltip with brief explanation
2. List of existing items
3. "Add" button at bottom
4. Each item has inline edit + delete
5. Changes save via HTMX (no full page reload)
6. Empty state: "No [items] added yet. Click 'Add [Item]' to get started."

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/partials/product_specs.html` | Redesign |
| `templates/admin/partials/product_features.html` | Redesign |
| `templates/admin/partials/product_certifications.html` | Redesign |
| `templates/admin/partials/product_downloads.html` | Redesign |
| `templates/admin/partials/product_images.html` | Redesign |

## Dependencies
- Phase 01, 07
