# Phase 02 - Sidebar Navigation Redesign

## Current State
- Flat list of 15+ links in dark blue sidebar (#004499)
- No grouping, no collapsing
- Settings and Page Sections are standalone top-level links
- No visual distinction between content types

## Goal
Grouped, collapsible navigation where each group contains its related pages AND its own settings link.

## New Sidebar Structure

```
[Logo: BlueJay Labs icon + text]

Dashboard                    (always visible, not in a group)

--- WEBSITE ---
v Homepage                   (collapsible group)
    Heroes
    Stats
    Testimonials
    CTAs
    Homepage Settings        <-- per-section settings

v About                     (collapsible group)
    Overview
    Mission/Vision/Values
    Core Values
    Milestones
    Certifications
    About Settings

v Header                    (single page, no children)

v Footer                    (single page, no children)

--- CONTENT ---
v Products                  (collapsible group)
    All Products
    Categories
    Product Settings         <-- visibility, display options

v Solutions                 (collapsible group)
    All Solutions
    Solution Settings

v Blog                      (collapsible group)
    All Posts
    Categories
    Authors
    Tags
    Blog Settings

v Case Studies              (single list+form, no sub-pages)

v Whitepapers               (collapsible group)
    All Whitepapers
    Topics
    Downloads

v Partners                  (collapsible group)
    All Partners
    Tiers
    Testimonials

--- ADMIN ---
Media Library               (new)
Contact Submissions
Navigation
Activity Log                (new)
Global Settings             (site name, SEO defaults, social media)

---
[View Site ->]              (pinned at bottom)
[User: Name / Role]         (pinned at bottom)
```

## Implementation Details

### Collapsible Groups
- Click group header to expand/collapse
- Chevron icon rotates on toggle (> when collapsed, v when expanded)
- Remember open/closed state in localStorage
- Active page's group auto-expands on page load
- Smooth height transition (CSS `max-height` animation)

### Section Labels
- "WEBSITE", "CONTENT", "ADMIN" as small gray uppercase labels
- Not clickable, just visual dividers

### Active State
- Active link: white bg, primary blue text, 3px left blue border
- Active group header: slightly lighter bg

### Per-Section Settings Links
- Appear as last item in each collapsible group
- Gear icon + "Settings" label
- These link to `/admin/{section}/settings` routes
- Each section's settings page only shows settings relevant to that section

### Visual Design
- Keep dark blue bg (#004499)
- White text, hover state lighter blue (#0066CC)
- Width: 260px (up from 256px to match mockup)
- Scrollable middle section if navigation overflows
- Logo and user footer pinned (not scrollable)

### Settings Distribution
Move settings from the monolithic settings page to per-section:

| Setting | Moves To |
|---------|----------|
| Homepage hero/stats visibility | Homepage Settings |
| Blog posts per page | Blog Settings |
| Product display options | Product Settings |
| About page section toggles | About Settings |
| Site name, tagline | Global Settings |
| Contact email, phone, address | Global Settings |
| SEO defaults | Global Settings |
| Social media URLs | Global Settings |
| GA tracking ID | Global Settings |

### Mobile Behavior
- Sidebar hidden by default on mobile
- Hamburger in header triggers slide-in overlay
- Semi-transparent black backdrop behind sidebar
- Tap backdrop or X button to close
- Same collapsible groups work on mobile

## Files to Create/Modify
| File | Action |
|------|--------|
| `templates/partials/admin-sidebar.html` | Rewrite |
| `public/css/admin-styles.css` | Add sidebar styles |
| `public/js/admin.js` | Create - sidebar toggle, collapse state |

## New Routes Needed
- `GET /admin/homepage/settings` - Homepage section settings
- `GET /admin/about/settings` - About section settings
- `GET /admin/products/settings` - Product section settings
- `GET /admin/solutions/settings` - Solution section settings
- `GET /admin/blog/settings` - Blog section settings
- `POST` variants for each above

## Database Changes
- May need to add columns to `settings` table for section-specific settings
- OR create a `section_settings` table with key-value pairs per section

## Dependencies
- Phase 01 (base layout) must be complete
