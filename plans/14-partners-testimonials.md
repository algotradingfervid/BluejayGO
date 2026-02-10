# Phase 14 - Partners & Testimonials

## Current State
- Partners: list + form
- Partner tiers: list + form
- Testimonials: list + form
- All basic CRUD

## Partners List Changes
- Search by name
- Tier dropdown filter
- Status filter
- Table: Logo (small), Name, Tier Badge, Website, Status, Actions
- Empty state

## Partners Form

### Collapsible Sections

**Section 1: Company Info** (open)
- Company Name
  - Tooltip: "Partner company's official name."
- Slug (auto-generated)
- Tier dropdown (references Partner Tiers)
  - Tooltip: "Partnership level. Manage tiers in the sidebar under Partners > Tiers."
- Website URL
  - Tooltip: "Link to partner's website. Shown as a clickable link on the partners page."
- Status dropdown

**Section 2: Branding** (open)
- Company Logo upload
  - Tooltip: "Partner's logo. Displayed on the partners page. Recommended: 200x100 PNG with transparency."
- Description (textarea)
  - Tooltip: "Brief description of the partnership and what they offer."

**Section 3: SEO** (collapsed)
- Meta Title + Meta Description

## Partner Tiers Page
- Simple list: Name, Color Badge, Partner Count, Sort Order, Actions
- Form: Name, Description, Sort Order, Badge Color
- Tooltip: "Tiers define partnership levels (e.g., Gold, Silver, Bronze). Partners are grouped by tier on the public page."

## Testimonials List Changes
- Table: Quote (truncated), Author, Company, Rating Stars, Featured, Actions
- Filter: Featured only toggle
- Empty state

## Testimonials Form
- Author Name
  - Tooltip: "Name of the person giving the testimonial."
- Author Title
  - Tooltip: "Job title (e.g., 'VP of Operations')."
- Company
  - Tooltip: "Company the testimonial author works for."
- Quote (textarea)
  - Tooltip: "The testimonial text. Keep it concise and specific about results."
- Rating (1-5 star selector)
  - Tooltip: "Star rating. Only 4-5 star testimonials are recommended for display."
- Featured (toggle)
  - Tooltip: "Featured testimonials appear on the homepage."
- Author Photo upload (optional)

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/partners_list.html` | Rewrite |
| `templates/admin/pages/partners_form.html` | Rewrite |
| `templates/admin/pages/partner_tiers_list.html` | Polish |
| `templates/admin/pages/partner_tiers_form.html` | Add tooltips |
| `templates/admin/pages/testimonials_list.html` | Rewrite |
| `templates/admin/pages/testimonials_form.html` | Rewrite |
| `internal/handlers/admin/partners.go` | Add filtering |

## Dependencies
- Phase 01, 02
