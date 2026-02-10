# Phase 16 - Homepage Customization

## Current State
- Separate list + form pages for: heroes, stats, testimonials, CTAs
- Basic CRUD for each

## Goal
All homepage sections managed under the "Homepage" sidebar group. Each section gets a clean, visual editor.

## Heroes Section
### List
- Visual cards (not table): show hero image + title + subtitle + CTA button text
- Sort order badges
- Status indicator (active/inactive)

### Form
- Title
  - Tooltip: "Hero banner headline. Keep it bold and concise (5-8 words)."
- Subtitle (textarea)
  - Tooltip: "Supporting text below the headline. 1-2 sentences max."
- Background Image upload
  - Tooltip: "Full-width banner image. Recommended: 1920x600. Text overlays this image."
- CTA Button Text
  - Tooltip: "Text on the hero's call-to-action button (e.g., 'Get Started')."
- CTA Button Link
  - Tooltip: "Where the CTA button links to."
- Sort Order
- Active toggle
  - Tooltip: "Only active heroes are displayed. Use sort order to control which shows first."

## Stats Section
### List
- Row of stat cards preview: number + label
- Each card editable inline or via form

### Form
- Number/Value (text - allows "500+" format)
  - Tooltip: "The stat number displayed prominently (e.g., '500+', '98%', '24/7')."
- Label
  - Tooltip: "What the number represents (e.g., 'Products Sold', 'Customer Satisfaction')."
- Icon (Material Symbol name)
  - Tooltip: "Optional icon displayed with the stat."
- Sort Order

## Testimonials Section
### List
- Select from existing testimonials (created in Partners > Testimonials)
- Checkbox list to pick which testimonials appear on homepage
- Sort order for selected ones

### Form
- No separate form - uses testimonial selector
- Tooltip: "Choose which testimonials to feature on the homepage. Create new testimonials under Partners > Testimonials."

## CTAs Section
### List
- Visual cards showing CTA blocks

### Form
- Heading
  - Tooltip: "CTA section heading (e.g., 'Ready to get started?')."
- Description (textarea)
  - Tooltip: "Supporting text encouraging action."
- Button Text
  - Tooltip: "Primary button label."
- Button Link
- Background Style: Light / Dark / Primary
- Sort Order

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/homepage_heroes_list.html` | Redesign as visual cards |
| `templates/admin/pages/homepage_hero_form.html` | Add tooltips |
| `templates/admin/pages/homepage_stats_list.html` | Redesign |
| `templates/admin/pages/homepage_stat_form.html` | Add tooltips |
| `templates/admin/pages/homepage_testimonials_list.html` | Redesign as selector |
| `templates/admin/pages/homepage_testimonial_form.html` | Replace with selector |
| `templates/admin/pages/homepage_cta_list.html` | Redesign |
| `templates/admin/pages/homepage_cta_form.html` | Add tooltips |

## Dependencies
- Phase 01, 02
- Phase 14 (for testimonial data to select from)
