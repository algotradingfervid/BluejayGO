# Phase 15 - About Page Sections

## Current State
- About overview form (text fields)
- Mission/Vision/Values form
- Core values list + form
- Milestones list + form
- Certifications list + form

## Goal
All About sub-pages live under the "About" sidebar group. Polish forms with tooltips and consistent styling.

## About Overview Form
- Page Title
  - Tooltip: "The main heading for the About page."
- Subtitle
  - Tooltip: "Subheading text below the main title."
- Content (Trix editor)
  - Tooltip: "Full description of your company. Shown as the main About page content."
- Hero Image upload
  - Tooltip: "Banner image for the About page hero section."

## Mission, Vision & Values Form
- Mission Statement (textarea)
  - Tooltip: "Your company's mission - what you do and why. Keep to 1-2 sentences."
- Vision Statement (textarea)
  - Tooltip: "Where your company is headed. An aspirational future state."
- Values Introduction (textarea)
  - Tooltip: "Brief intro text before the list of core values."

## Core Values
### List
- Cards layout instead of table (each value as a card with icon + title + description)
- Sort order badges
- Reorder by changing sort number

### Form
- Title
  - Tooltip: "The value name (e.g., 'Innovation', 'Integrity')."
- Icon (dropdown or text input for Material Symbol name)
  - Tooltip: "Material Symbols icon name. Browse icons at fonts.google.com/icons."
- Description (textarea)
  - Tooltip: "1-2 sentences explaining what this value means to your company."
- Sort Order (number)

## Milestones (Timeline)
### List
- Vertical timeline layout instead of table
- Each milestone: Year + Title + Description
- Chronological order

### Form
- Year
  - Tooltip: "The year this milestone occurred (e.g., 2015)."
- Title
  - Tooltip: "Brief milestone headline (e.g., 'Founded in Houston, TX')."
- Description (textarea)
  - Tooltip: "Details about this milestone event."
- Sort Order

## Certifications
### List
- Grid of certification cards (logo + name + issuer)

### Form
- Name
  - Tooltip: "Certification name (e.g., 'ISO 9001:2015')."
- Issuing Body
  - Tooltip: "Organization that granted the certification."
- Certificate Number (optional)
- Logo/Image upload
  - Tooltip: "Certification badge or logo image."
- Description (textarea)
- Valid Until (date, optional)
  - Tooltip: "Expiration date if applicable. Expired certs are flagged."

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/about_overview_form.html` | Add tooltips, polish |
| `templates/admin/pages/about_mvv_form.html` | Add tooltips, polish |
| `templates/admin/pages/core_values_list.html` | Redesign as cards |
| `templates/admin/pages/core_values_form.html` | Add tooltips |
| `templates/admin/pages/milestones_list.html` | Redesign as timeline |
| `templates/admin/pages/milestones_form.html` | Add tooltips |
| `templates/admin/pages/certifications_list.html` | Redesign as grid |
| `templates/admin/pages/certifications_form.html` | Add tooltips |

## Dependencies
- Phase 01, 02
