# Phase 06 - Footer Management

## Current State
- Footer toggles and column headings in global settings
- No dedicated management page
- Basic text fields only

## Goal
Dedicated "Footer" page in sidebar (under WEBSITE) for full footer customization.

## Page Layout

### Section 1: Footer Layout
- Column count selector: 2, 3, or 4 columns
  - Tooltip: "How many columns to display in the footer. 3 or 4 works best for most sites."
- Background style: Dark (navy) / Light (white) / Primary (blue)

### Section 2: Column Configuration
For each column (based on count selected):
- Column Heading (text input)
  - Tooltip: "The bold title shown above this column's links."
- Column Type: Links / Text / Contact Info
- If Links: manage list of label + URL pairs (add/remove)
- If Text: textarea for custom content
- If Contact Info: auto-pulls from Global Settings (display address, phone, email)

### Section 3: Social Media Row
- Toggle: Show social icons in footer
- Inherits URLs from Global Settings social media
- Display style: Icons only / Icons + platform name

### Section 4: Bottom Bar
- Copyright text (text input with year placeholder `{year}`)
  - Tooltip: "Footer copyright text. Use {year} for auto-updating year."
- Legal links: manage label + URL pairs (Privacy Policy, Terms, etc.)

## Backend Changes
- New route: `GET /admin/footer` and `POST /admin/footer`
- New columns in `settings` or new `footer_settings` table:
  - `footer_columns` (INTEGER - count)
  - `footer_bg_style` (TEXT)
  - `footer_show_social` (INTEGER 0/1)
  - `footer_social_style` (TEXT)
  - `footer_copyright` (TEXT)
- New table `footer_column_items` for column content:
  - `id`, `column_index`, `type`, `heading`, `content`, `sort_order`
- New table `footer_links` for link items:
  - `id`, `column_item_id`, `label`, `url`, `sort_order`
- New table `footer_legal_links`:
  - `id`, `label`, `url`, `sort_order`

## Files to Create/Modify
| File | Action |
|------|--------|
| `templates/admin/pages/footer_form.html` | Create |
| `internal/handlers/admin/footer.go` | Create |
| `db/migrations/028_footer_settings.up.sql` | Create |
| `db/migrations/028_footer_settings.down.sql` | Create |
| `db/queries/footer.sql` | Create |
| `cmd/server/main.go` | Add routes |

## Dependencies
- Phase 01, 02, 04
