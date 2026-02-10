# Phase 05 - Header Management

## Current State
- No dedicated header management page
- Nav visibility toggles and labels buried in global settings form
- No logo upload capability
- No CTA button management

## Goal
A dedicated "Header" page in the sidebar (under WEBSITE group) for full header customization.

## Page Layout

### Section 1: Logo
- Upload area for site logo (drag-and-drop or click to upload)
  - Tooltip: "Your site logo displayed in the header. Recommended size: 240x60px. PNG or SVG."
- Preview of current logo
- Option to remove/replace
- Alt text input
  - Tooltip: "Accessibility text for the logo. Typically your company name."

### Section 2: Navigation Links
- Toggle switches to show/hide each nav section:
  - Products, Solutions, Case Studies, About, Blog, Whitepapers, Partners, Contact
  - Tooltip on each: "Toggle whether this page appears in the main navigation menu."
- Editable label for each link (so "Products" could become "Our Products")
  - Tooltip: "Customize the display text for this navigation link."
- Drag handle for reordering (stretch goal - can be Phase 2)

### Section 3: CTA Button
- Toggle: Show/Hide CTA button in header
- CTA Text (e.g., "Request a Quote")
  - Tooltip: "Button text shown in the header. Keep it short and action-oriented."
- CTA Link URL
  - Tooltip: "Where the button links to. Use /contact for the contact page."
- CTA Style: Primary (filled) or Secondary (outlined)

### Section 4: Contact Info Display
- Toggle: Show phone number in header
- Toggle: Show email in header
- These pull from Global Settings contact info (display only, not editable here)
- Note text: "Edit contact info in Global Settings"

### Section 5: Social Icons
- Toggle: Show social media icons in header
- Which platforms to show (checkboxes for each social platform that has a URL in Global Settings)
- Display style: Icons only / Icons + labels

## Backend Changes
- New route: `GET /admin/header` -> render header management page
- New route: `POST /admin/header` -> save header settings
- Header settings can use existing `settings` table columns for toggles/labels
- New columns needed in `settings`:
  - `header_logo_path` (TEXT)
  - `header_logo_alt` (TEXT)
  - `header_cta_enabled` (INTEGER 0/1)
  - `header_cta_text` (TEXT)
  - `header_cta_url` (TEXT)
  - `header_cta_style` (TEXT)
  - `header_show_phone` (INTEGER 0/1)
  - `header_show_email` (INTEGER 0/1)
  - `header_show_social` (INTEGER 0/1)
  - `header_social_style` (TEXT)

## Files to Create/Modify
| File | Action |
|------|--------|
| `templates/admin/pages/header_form.html` | Create |
| `internal/handlers/admin/header.go` | Create |
| `db/migrations/027_header_settings.up.sql` | Create |
| `db/migrations/027_header_settings.down.sql` | Create |
| `cmd/server/main.go` | Add routes |

## Dependencies
- Phase 01, 02, 04 (Global Settings has contact/social data this page references)
