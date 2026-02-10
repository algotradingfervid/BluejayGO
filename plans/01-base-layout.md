# Phase 01 - Base Layout & Design System

## Why First
Everything else depends on the base layout, header component, and shared CSS utilities.

## Current State
- `templates/admin/layouts/base.html`: Basic HTML shell with JetBrains Mono, Tailwind CDN, HTMX
- `public/css/styles.css`: ~1.5KB with `.manual-border`, `.manual-shadow`, `.btn-press`
- Inconsistent styling: dashboard uses brutalist style, forms use rounded Tailwind defaults

## Changes

### 1. Update `base.html`
- Add Material Symbols Outlined font (used in mockups for icons)
- Add Inter font for display headings alongside JetBrains Mono for body
- Add viewport meta for proper mobile behavior
- Add `<link>` to an expanded `admin-styles.css`

### 2. Create `public/css/admin-styles.css` (expand from styles.css)
Define reusable utility classes matching the brutalist design system:

**Border & Shadow:**
- `.manual-border` - 2px solid black (keep existing)
- `.manual-shadow` - 4px 4px 0 0 #000 (keep existing)
- `.manual-shadow-sm` - 2px 2px 0 0 #000 (new, for smaller elements)

**Buttons:**
- `.btn-primary` - Primary blue bg, white text, black border, manual shadow
- `.btn-secondary` - White bg, black border, primary text
- `.btn-danger` - Red bg, white text, black border
- `.btn-ghost` - No bg, no border, hover adds bg
- All buttons: no border-radius, active press effect

**Form Elements:**
- `.input-brutal` - 2px black border, no radius, white bg, focus ring in primary blue
- `.select-brutal` - Same as input with custom dropdown arrow
- `.textarea-brutal` - Same style, resizable vertically
- `.toggle-switch` - Custom toggle (already exists, standardize)

**Cards:**
- `.card-brutal` - White bg, manual-border, manual-shadow, p-6
- `.card-brutal-sm` - Same but p-4 and smaller shadow

**Badges:**
- `.badge-published` - Green bg, black text
- `.badge-draft` - Gray bg, black text
- `.badge-pending` - Yellow bg, black text
- `.badge-error` - Red bg, white text

**Tooltip System:**
- `.tooltip-trigger` - Relative position, cursor help
- `.tooltip-content` - Absolute positioned popup, hidden by default
- Show on hover with pure CSS (`:hover > .tooltip-content`)
- Simple black bg, white text, max-width 250px, small text
- Arrow pointer using CSS triangle

**Typography:**
- `.font-display` - Inter font, uppercase, extrabold (for page titles)
- Body text stays JetBrains Mono

### 3. Update Base Layout Structure
```
<body>
  <div class="flex h-screen overflow-hidden">
    <!-- Sidebar (included as partial) -->
    {{template "admin-sidebar" .}}

    <!-- Main area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Sticky header bar -->
      <header class="h-16 border-b-2 border-black bg-white flex items-center px-6 justify-between shrink-0">
        <!-- Mobile hamburger -->
        <!-- Page title (injected by each page) -->
        <!-- Right: help icon, notification bell, user avatar -->
      </header>

      <!-- Scrollable content area -->
      <main class="flex-1 overflow-y-auto bg-gray-50 p-6 lg:p-8">
        {{block "content" .}}{{end}}
      </main>
    </div>
  </div>
</body>
```

### 4. Sticky Header Component
- Left: Hamburger menu (mobile only) + breadcrumb/page title
- Right: Help link (?) + User avatar with name + Logout button
- Height: 64px, white bg, 2px black bottom border
- Mobile: hamburger toggles sidebar overlay

### 5. Mobile Sidebar Overlay
- On mobile (<1024px): sidebar is hidden off-screen (`-translate-x-full`)
- Hamburger click slides it in with backdrop overlay
- Close button (X) in top-right of sidebar on mobile

## Files to Create/Modify
| File | Action |
|------|--------|
| `templates/admin/layouts/base.html` | Modify |
| `public/css/admin-styles.css` | Create (replaces styles.css for admin) |
| `public/css/styles.css` | Keep for public site |
| `templates/partials/admin-header.html` | Create |

## Dependencies
- None (this is the foundation)

## Tooltip Specification
Every form field that needs explanation will use this pattern:
```html
<div class="tooltip-trigger">
  <label>Field Name</label>
  <span class="tooltip-icon">?</span>
  <div class="tooltip-content">
    Explanation of what this field does and expected format.
  </div>
</div>
```
The `?` icon appears next to the label. Hovering shows the tooltip popup.
