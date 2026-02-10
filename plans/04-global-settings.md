# Phase 04 - Global Settings Page

## Current State
- One massive `settings_form.html` with ALL settings in a single scrolling form
- Includes: site name, tagline, contact info, SEO, social media, nav toggles, footer toggles, nav labels, footer headings, footer text

## Goal
Slim down to ONLY global/site-wide settings. Section-specific settings move to their respective sidebar groups (handled in later phases).

## What Stays in Global Settings
Only settings that apply to the entire site:

### Tab 1: General
- Site Name (text input)
  - Tooltip: "Your company name. Appears in browser tab and site header."
- Site Tagline (text input)
  - Tooltip: "A short phrase describing your business. Used in SEO and header."
- Branding section:
  - Logo for light backgrounds (file upload with preview, recommended 240x60)
  - Logo for dark backgrounds (file upload with preview)
  - Favicon (file upload, recommended 32x32)

### Tab 2: Contact Info
- Email (email input)
  - Tooltip: "Primary contact email displayed on the site and used for form notifications."
- Phone (text input)
- Address (street, city, state, zip, country - grid layout)
- Business Hours (textarea)

### Tab 3: Social Media
- Facebook URL
- Twitter/X URL
- LinkedIn URL
- Instagram URL
- YouTube URL
- Each with URL validation icon (green check / red X)
- Tooltip on each: "Full URL to your [platform] profile page."

### Tab 4: SEO Defaults
- Default Meta Title (with 70-char counter)
  - Tooltip: "Fallback title for pages that don't have their own. Appears in search results."
- Default Meta Description (with 160-char counter)
  - Tooltip: "Fallback description for pages without their own. Appears in Google snippets."
- Google Analytics ID
  - Tooltip: "Your GA tracking ID (e.g., G-XXXXXXXXXX). Leave blank to disable tracking."
- Default OG Image (file upload)

## What Moves Out
| Setting | Moves To |
|---------|----------|
| Header nav toggle switches | Phase 05 (Header Management) |
| Navigation labels | Phase 05 (Header Management) |
| Footer toggle switches | Phase 06 (Footer Management) |
| Footer column headings | Phase 06 (Footer Management) |
| Footer text | Phase 06 (Footer Management) |

## UI Design
- Horizontal tab bar at top (desktop) / dropdown selector (mobile)
- Active tab: primary color bottom border (3px), bold text
- Each tab is a separate card section
- "Unsaved changes" yellow banner appears when form is dirty
- "Save" and "Discard" buttons in banner
- Character counters: green (ok) -> yellow (near limit) -> red (over)

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/settings_form.html` | Rewrite (slim down) |
| `internal/handlers/admin/settings.go` | Update to handle only global settings |

## Database Changes
- No schema changes needed; existing `settings` table columns still used
- Some columns will also be read by section-specific settings pages

## Dependencies
- Phase 01, 02
