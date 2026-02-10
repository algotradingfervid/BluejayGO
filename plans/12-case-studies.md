# Phase 12 - Case Studies List & Form

## Current State
- Basic list + form pages
- HTMX sub-tabs for products and metrics

## Changes to List

### Filter Bar
- Search by title
- Status filter
- Pagination

### Table
- Columns: Image, Title, Client, Industry, Status, Updated, Actions
- Empty state

## Changes to Form

### Collapsible Sections

**Section 1: Basic Info** (open)
- Title
  - Tooltip: "Case study title. Include the client name or project for clarity."
- Slug (auto-generated)
- Client Name
  - Tooltip: "The company or organization featured in this case study."
- Industry dropdown
- Status dropdown

**Section 2: Story** (open)
- Challenge (textarea)
  - Tooltip: "What problem did the client face? Sets the stage for the solution."
- Solution (Trix editor)
  - Tooltip: "How your products/services solved the client's challenge."
- Results (textarea)
  - Tooltip: "Key outcomes and benefits the client achieved."

**Section 3: Media** (collapsed)
- Featured Image upload
  - Tooltip: "Hero image for the case study. Client facility or product in use recommended."
- Client Logo upload
  - Tooltip: "Client's company logo. Shown alongside the case study."

**Section 4: SEO** (collapsed)
- Meta Title + Meta Description

### Sub-pages (edit mode only)
**Products Tab:** linked products selector
**Metrics Tab:** number + label pairs (e.g., "45% increase in efficiency")
- Tooltip: "Quantifiable results to highlight. Use specific numbers for impact."

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/case_studies_list.html` | Rewrite |
| `templates/admin/pages/case_studies_form.html` | Rewrite |
| `templates/admin/partials/case_study_products.html` | Polish |
| `templates/admin/partials/case_study_metrics.html` | Polish |
| `internal/handlers/admin/case_studies.go` | Add filtering, pagination |

## Dependencies
- Phase 01, 02
