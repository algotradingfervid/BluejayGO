# Phase 11 - Solutions List & Form

## Current State
- List page with basic table
- Form with HTMX sub-tabs for: challenges, products, stats, CTAs

## Changes to List

### Filter Bar
- Search by name
- Status filter (Published / Draft)
- Pagination (15 per page)

### Table
- Columns: Thumbnail, Name, Status, Industry, Updated, Actions
- Status badges (brutalist style)
- Empty state

## Changes to Form

### Collapsible Sections

**Section 1: Basic Info** (open)
- Title
  - Tooltip: "Solution name as displayed on the site."
- Slug (auto-generated)
- Industry dropdown
  - Tooltip: "Which industry this solution serves."
- Status dropdown
- Short Description (textarea)
  - Tooltip: "1-2 sentence summary shown on solution cards."

**Section 2: Content** (open)
- Full Description (Trix editor)
  - Tooltip: "Detailed description of the solution, including how it works and benefits."
- Hero Image upload
  - Tooltip: "Main banner image for the solution page. Recommended 1200x600."

**Section 3: SEO** (collapsed)
- Meta Title (70-char counter)
- Meta Description (160-char counter)

### Sub-pages (below form, edit mode only)
Keep existing HTMX tabs but polish:

**Challenges Tab:**
- List of challenge items (icon + title + description)
- "Add Challenge" button
- Tooltip: "Problems this solution addresses. Shown as pain points on the page."

**Products Tab:**
- Searchable product selector (checkbox list)
- Shows selected products as cards
- Tooltip: "Products used in this solution. Links to product pages."

**Stats Tab:**
- Each stat: Number + Label + Description
- "Add Stat" button
- Tooltip: "Key metrics that demonstrate the solution's impact (e.g., '40% cost reduction')."

**CTAs Tab:**
- Each CTA: Text + URL + Style (primary/secondary)
- "Add CTA" button
- Tooltip: "Call-to-action buttons shown at the bottom of the solution page."

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/solutions_list.html` | Rewrite |
| `templates/admin/pages/solutions_form.html` | Rewrite with collapsible sections |
| `templates/admin/partials/solution_challenges.html` | Polish |
| `templates/admin/partials/solution_products.html` | Polish |
| `templates/admin/partials/solution_stats.html` | Polish |
| `templates/admin/partials/solution_ctas.html` | Polish |
| `internal/handlers/admin/solutions.go` | Add filtering, pagination |

## Dependencies
- Phase 01, 02
