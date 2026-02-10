# Bluejay CMS — Brutalist Design System

**Version 2026.01** | Last Updated: February 10, 2026

---

## 1. Design Philosophy

### What is Brutalist Design?

Bluejay CMS embraces a **brutalist design aesthetic** characterized by:

- **Raw, unadorned elements** — No border-radius, no gradients, no soft shadows
- **High contrast** — Stark black borders on white backgrounds
- **Hard edges and sharp corners** — Everything is angular and precise
- **Functional over decorative** — Every visual element serves a purpose
- **Monospace typography** — Technical, utilitarian feel
- **Manual shadows** — Deliberate, offset box shadows instead of subtle blur
- **Bold, uppercase text** — Commands attention and authority
- **Grid patterns and technical references** — Dotted backgrounds, system labels

### Why Brutalism?

This design language communicates:
- **Technical precision** and engineering excellence
- **Industrial strength** and reliability
- **Transparency** — No hidden functionality or deceptive UI patterns
- **Authority** and confidence
- **Timeless aesthetic** — Won't feel dated in 5 years

The brutalist approach stands out in a sea of soft, rounded modern designs and creates a memorable, distinctive brand identity for technical/industrial products.

---

## 2. Typography

### Font Families

```css
/* Primary font — Used everywhere */
font-family: 'JetBrains Mono', monospace;

/* Display font — Used for large headings only */
font-family: 'Inter', sans-serif;
```

**JetBrains Mono** is a monospace font with excellent readability and technical character. It's used for:
- Body text
- UI labels
- Buttons
- Form inputs
- Navigation
- Code snippets
- All admin interface text

**Inter** is only used for:
- Large display headings (when you need extra impact)
- Marketing hero sections

### Font Loading

```html
<!-- Include in <head> -->
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600;700&display=swap" rel="stylesheet">
```

### Font Weight Hierarchy

- **400** (Regular) — Body text, descriptions
- **500** (Medium) — Subheadings, labels
- **600** (Semi-Bold) — Section headings
- **700** (Bold) — Buttons, primary headings, emphasis
- **800** (Extra-Bold) — Display headings (Inter only)

### Display Heading Class

```css
.font-display {
    font-family: 'Inter', sans-serif;
    text-transform: uppercase;
    font-weight: 800;
}
```

### Typography Scale

```css
/* Headings */
h1: 2rem (32px) - 3.75rem (60px)
h2: 1.5rem (24px) - 2rem (32px)
h3: 1.25rem (20px)

/* Body */
text-lg: 1.125rem (18px)
text-base: 1rem (16px)
text-sm: 0.875rem (14px)
text-xs: 0.75rem (12px)

/* Micro text */
text-[10px]: 0.625rem (10px) — Used for labels, metadata, captions
```

### Text Transform

- **Buttons** — ALWAYS uppercase
- **Labels** — ALWAYS uppercase
- **Section headings** — Usually uppercase
- **Body text** — Normal case
- **Navigation** — Usually uppercase

---

## 3. Color Palette

### CSS Custom Properties

**Admin Styles (admin-styles.css):**
```css
:root {
    --primary: #0066CC;
    --primary-dark: #004499;
    --navy-tech: #004499;
    --black: #000000;
    --white: #FFFFFF;
    --danger: #DC2626;
    --success: #16A34A;
    --warning: #EAB308;
    --gray-light: #F3F4F6;
}
```

**Public Styles (styles.css):**
```css
:root {
    --primary: #0066CC;
    --primary-dark: #004499;
    --navy-tech: #004499;
    --black: #000000;
    --white: #FFFFFF;
}
```

Note: Public styles use a subset of CSS variables. For danger, success, warning colors on public pages, use Tailwind utilities or inline hex values.

### Color Usage Guidelines

#### Primary Blue (`--primary: #0066CC`)
- Primary action buttons
- Links
- Active states
- Hover backgrounds
- Brand accent color
- Section dividers

#### Navy Tech (`--navy-tech: #004499`)
- Admin sidebar background
- Dark button hover states
- Technical emphasis

#### Black (`--black: #000000`)
- All borders (ALWAYS 2px solid)
- All text (default)
- Box shadows
- Header backgrounds
- High-contrast UI elements

#### White (`--white: #FFFFFF`)
- Page backgrounds
- Card backgrounds
- Button text on dark backgrounds
- Inverted UI elements

#### Danger Red (`--danger: #DC2626`)
- Delete buttons
- Error states
- Warning banners
- Critical actions

#### Success Green (`--success: #16A34A`)
- Published status badges
- Success messages
- Confirmation indicators

#### Warning Yellow (`--warning: #EAB308`)
- Draft status badges
- Preview mode banners
- Attention indicators

#### Gray Light (`--gray-light: #F3F4F6`)
- Page backgrounds (admin)
- Disabled states
- Subtle hover states

### Status Badge Colors

**Using CSS Classes (admin-styles.css):**
```css
/* Published */
.badge-published { background: #BBF7D0; } /* light green */

/* Draft */
.badge-draft { background: #E5E7EB; } /* light gray */

/* Pending */
.badge-pending { background: #FEF08A; } /* light yellow */

/* Error */
.badge-error { background: #DC2626; color: #FFFFFF; } /* danger red */
```

**Using Tailwind Utilities (common in templates):**
```html
<!-- Published -->
<span class="bg-green-400 text-black px-2 py-1 text-xs font-bold uppercase border-2 border-black">Published</span>

<!-- Draft -->
<span class="bg-yellow-300 text-black px-2 py-1 text-xs font-bold uppercase border-2 border-black">Draft</span>

<!-- Archived/Inactive -->
<span class="bg-gray-300 text-black px-2 py-1 text-xs font-bold uppercase border-2 border-black">Archived</span>
```

Note: Templates frequently use Tailwind color utilities (bg-green-400, bg-yellow-300) rather than the predefined CSS classes for flexibility.

---

## 4. Border Rules

### THE GOLDEN RULE: NO BORDER-RADIUS

```css
/* CORRECT */
border-radius: 0;

/* WRONG — Never use */
border-radius: 4px;
border-radius: 8px;
```

**Exception:** Full circles only (avatars, icons)
```css
border-radius: 9999px; /* or 50% for circles */
```

### Border Standard

```css
/* Default border style */
border: 2px solid var(--black);
border-radius: 0;
```

**ALWAYS 2px solid black** unless:
- White borders on dark backgrounds
- Colored borders for specific semantic meaning (danger, success)

### Border Utilities

**Public Styles (styles.css):**
```css
.manual-border {
    border: 2px solid var(--black);
    border-radius: 0;
}

.manual-border-white {
    border: 2px solid var(--white);
    border-radius: 0;
}
```

Note: These utilities exist in public styles only. In admin templates, use Tailwind utilities `border-2 border-black` or inline styles.

### Tailwind Border Radius Override

```javascript
// In tailwind.config
borderRadius: {
    "DEFAULT": "0px",
    "lg": "0px",
    "xl": "0px",
    "full": "9999px"  // Only for circles
}
```

This forces all Tailwind utilities to use sharp corners by default.

---

## 5. Shadow System

### Manual Box Shadows (No Blur)

The brutalist aesthetic uses **hard, offset shadows** with no blur radius.

```css
/* Standard shadow — Use for cards, buttons, large elements */
box-shadow: 4px 4px 0 0 var(--black);

/* Small shadow — Use for badges, small buttons, inline elements */
box-shadow: 2px 2px 0 0 var(--black);

/* Large shadow — Use for hero sections, major page sections */
box-shadow: 6px 6px 0 0 var(--black);
```

### Shadow Utilities

**Admin Styles (admin-styles.css):**
```css
.manual-shadow-sm {
    box-shadow: 2px 2px 0 0 var(--black);
}
```

**Public Styles (styles.css):**
```css
.manual-shadow {
    box-shadow: 4px 4px 0 0 var(--black);
}

.manual-shadow-white {
    box-shadow: 4px 4px 0 0 var(--white);
}
```

Note: For standard 4px shadow in admin templates, use inline style `style="box-shadow: 4px 4px 0px #000;"` or the Tailwind-compatible approach since `.manual-shadow` is only in public styles.

### Interactive Shadow States

```css
/* Hover state — No change to shadow */
.element:hover {
    background-color: var(--primary);
}

/* Active/pressed state — Push down effect */
.element:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px 0 0 var(--black);
}
```

The "push down" effect simulates a physical button press by moving the element and reducing the shadow.

### Push-Down Utility

**Public Styles (styles.css):**
```css
.btn-press:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px 0 0 var(--black);
}
```

Note: This utility class exists in public styles. In admin templates, the push-down effect is typically implemented inline in component styles.

---

## 6. Button Components

### Primary Button

```css
.btn-primary {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background-color: var(--primary);
    color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    box-shadow: 4px 4px 0 0 var(--black);
    font-family: 'JetBrains Mono', monospace;
    font-weight: 700;
    font-size: 0.875rem;
    text-transform: uppercase;
    cursor: pointer;
    transition: background-color 0.15s, transform 0.1s, box-shadow 0.1s;
}

.btn-primary:hover {
    background-color: var(--primary-dark);
}

.btn-primary:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px 0 0 var(--black);
}
```

### Secondary Button

```css
.btn-secondary {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background-color: var(--white);
    color: var(--primary);
    border: 2px solid var(--black);
    border-radius: 0;
    box-shadow: 4px 4px 0 0 var(--black);
    font-family: 'JetBrains Mono', monospace;
    font-weight: 700;
    font-size: 0.875rem;
    text-transform: uppercase;
    cursor: pointer;
    transition: background-color 0.15s, transform 0.1s, box-shadow 0.1s;
}

.btn-secondary:hover {
    background-color: var(--gray-light);
}

.btn-secondary:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px 0 0 var(--black);
}
```

### Danger Button

```css
.btn-danger {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background-color: var(--danger);
    color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    box-shadow: 4px 4px 0 0 var(--black);
    font-family: 'JetBrains Mono', monospace;
    font-weight: 700;
    font-size: 0.875rem;
    text-transform: uppercase;
    cursor: pointer;
    transition: background-color 0.15s, transform 0.1s, box-shadow 0.1s;
}

.btn-danger:hover {
    background-color: #B91C1C;
}

.btn-danger:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px 0 0 var(--black);
}
```

### Ghost Button

```css
.btn-ghost {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background-color: transparent;
    color: var(--black);
    border: 2px solid transparent;
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-weight: 700;
    font-size: 0.875rem;
    text-transform: uppercase;
    cursor: pointer;
    transition: background-color 0.15s;
}

.btn-ghost:hover {
    background-color: var(--gray-light);
    border-color: var(--black);
}
```

### Button Usage

- **Primary** — Main CTA, submit forms, confirm actions
- **Secondary** — Cancel, back navigation, alternative actions
- **Danger** — Delete, remove, destructive actions
- **Ghost** — Subtle actions, tertiary options

---

## 7. Form Elements

### Input Field

```css
.input-brutal {
    width: 100%;
    padding: 0.625rem 0.75rem;
    background-color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.875rem;
    transition: border-color 0.15s, box-shadow 0.15s;
}

.input-brutal:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.2);
}

/* Note: Global focus styles in styles.css use 0.1 opacity:
   box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1); */
```

### Select Dropdown

```css
.select-brutal {
    width: 100%;
    padding: 0.625rem 2.5rem 0.625rem 0.75rem;
    background-color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.875rem;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23000' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    cursor: pointer;
    transition: border-color 0.15s, box-shadow 0.15s;
}

.select-brutal:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.2);
}

/* Note: Global focus styles in styles.css use 0.1 opacity:
   box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1); */
```

### Textarea

```css
.textarea-brutal {
    width: 100%;
    padding: 0.625rem 0.75rem;
    background-color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.875rem;
    resize: vertical;
    min-height: 6rem;
    transition: border-color 0.15s, box-shadow 0.15s;
}

.textarea-brutal:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.2);
}

/* Note: Global focus styles in styles.css use 0.1 opacity:
   box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1); */
```

### Form Labels

```html
<label class="block text-xs font-bold uppercase mb-1">
    Field Name
</label>
```

Always use:
- `text-xs` (12px)
- `font-bold`
- `uppercase`
- `mb-1` for consistent spacing

---

## 8. Card Components

### Standard Card

```css
.card-brutal {
    background-color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    box-shadow: 4px 4px 0 0 var(--black);
    padding: 1.5rem;
}
```

### Small Card

```css
.card-brutal-sm {
    background-color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    box-shadow: 2px 2px 0 0 var(--black);
    padding: 1rem;
}
```

### Interactive Card (Hover State)

```html
<a class="card-brutal group hover:bg-primary hover:text-white active:scale-[0.97] transition-all cursor-pointer">
    <!-- Card content -->
</a>
```

Use Tailwind's `group` utility to style child elements on hover:

```html
<h3 class="text-primary group-hover:text-white">Heading</h3>
<p class="opacity-60 group-hover:opacity-80">Description</p>
```

---

## 9. Badge/Status Components

### Published Badge

```css
.badge-published {
    display: inline-block;
    padding: 0.125rem 0.5rem;
    background-color: #BBF7D0;
    color: var(--black);
    border: 2px solid var(--black);
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
}
```

### Draft Badge

```css
.badge-draft {
    display: inline-block;
    padding: 0.125rem 0.5rem;
    background-color: #E5E7EB;
    color: var(--black);
    border: 2px solid var(--black);
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
}
```

### Pending Badge

```css
.badge-pending {
    display: inline-block;
    padding: 0.125rem 0.5rem;
    background-color: #FEF08A;
    color: var(--black);
    border: 2px solid var(--black);
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
}
```

### Error Badge

```css
.badge-error {
    display: inline-block;
    padding: 0.125rem 0.5rem;
    background-color: var(--danger);
    color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    font-family: 'JetBrains Mono', monospace;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
}
```

---

## 10. Table Styling

### Admin List Tables

```html
<div class="bg-white border-2 border-black mb-6" style="box-shadow: 4px 4px 0px #000;">
    <table class="w-full">
        <thead>
            <tr class="border-b-2 border-black bg-gray-100">
                <th class="px-4 py-3 text-left text-xs font-bold uppercase">Column</th>
            </tr>
        </thead>
        <tbody>
            <tr class="border-b border-gray-200 hover:bg-gray-50">
                <td class="px-4 py-3 text-sm">Cell content</td>
            </tr>
        </tbody>
    </table>
</div>
```

### Table Design Rules

- **Container** — White background, 2px black border, 4px shadow
- **Header row** — `bg-gray-100`, 2px bottom border, uppercase text
- **Body rows** — 1px gray border between rows
- **Hover state** — Light gray background (`bg-gray-50`)
- **Cell padding** — `px-4 py-3`
- **Text size** — `text-xs` for headers, `text-sm` for cells

---

## 11. Sidebar Navigation

### Structure

- **Fixed width** — 260px
- **Background** — Navy (`#004499`)
- **Border** — 4px solid black on the right
- **Sections** — Labeled with uppercase micro text

### Sidebar Link

```css
.sidebar-link {
    display: flex;
    align-items: center;
    gap: 0.625rem;
    padding: 0.5rem 0.75rem;
    margin-bottom: 1px;
    font-size: 0.8125rem;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.85);
    text-decoration: none;
    transition: background-color 0.15s;
    font-family: 'JetBrains Mono', monospace;
    border-left: 3px solid transparent;
}

.sidebar-link:hover {
    background-color: #0066CC;
    color: #fff;
}

.sidebar-link.active {
    background-color: #fff;
    color: #004499;
    font-weight: 700;
    border-left-color: #0066CC;
}
```

### Collapsible Group

```css
.sidebar-group-header {
    display: flex;
    align-items: center;
    gap: 0.625rem;
    width: 100%;
    padding: 0.5rem 0.75rem;
    font-size: 0.8125rem;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.85);
    background: none;
    border: none;
    cursor: pointer;
    transition: background-color 0.15s;
    font-family: 'JetBrains Mono', monospace;
    text-align: left;
    border-left: 3px solid transparent;
}

.sidebar-group-header:hover {
    background-color: #0066CC;
    color: #fff;
}

.sidebar-chevron {
    transition: transform 0.2s ease;
}

.sidebar-group.open .sidebar-chevron {
    transform: rotate(90deg);
}
```

### Sub-links

```css
.sidebar-sublink {
    display: block;
    padding: 0.375rem 0.75rem 0.375rem 2.75rem;
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.7);
    text-decoration: none;
    transition: background-color 0.15s, color 0.15s;
    font-family: 'JetBrains Mono', monospace;
    border-left: 3px solid transparent;
}

.sidebar-sublink:hover {
    background-color: #0066CC;
    color: #fff;
}

.sidebar-sublink.active {
    background-color: #fff;
    color: #004499;
    font-weight: 700;
    border-left-color: #0066CC;
}
```

### Mobile Sidebar

On screens smaller than 1024px (`max-width: 1023px`), the sidebar becomes a drawer that slides in from the left:

```css
@media (max-width: 1023px) {
    .admin-sidebar {
        position: fixed;
        top: 0;
        left: 0;
        bottom: 0;
        z-index: 50;
        transform: translateX(-100%);
        transition: transform 0.3s ease;
    }

    .admin-sidebar.open {
        transform: translateX(0);
    }
}
```

Breakpoint: `1023px` (corresponds to Tailwind's `lg:` breakpoint at `1024px`)

---

## 12. Tooltip System

### Structure

```html
<div class="tooltip-trigger">
    Label
    <span class="tooltip-icon">ⓘ</span>
    <div class="tooltip-content">
        Helpful explanation text
    </div>
</div>
```

### Styles

```css
.tooltip-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.125rem;
    height: 1.125rem;
    background-color: var(--black);
    color: var(--white);
    font-size: 0.75rem;
    font-weight: 700;
    cursor: help;
    border: none;
    line-height: 1;
}

.tooltip-content {
    display: none;
    position: absolute;
    bottom: calc(100% + 8px);
    left: 50%;
    transform: translateX(-50%);
    background-color: var(--black);
    color: var(--white);
    padding: 0.5rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 400;
    max-width: 250px;
    width: max-content;
    z-index: 50;
    line-height: 1.4;
}

.tooltip-trigger:hover > .tooltip-content {
    display: block;
}
```

Use tooltips for:
- Form field help text
- Icon explanations
- Feature descriptions
- Technical terms

---

## 13. Toggle Switches

### Structure

```html
<div class="toggle-switch" onclick="this.classList.toggle('active')"></div>
```

### Styles

```css
.toggle-switch {
    position: relative;
    width: 3rem;
    height: 1.5rem;
    background-color: #D1D5DB;
    border: 2px solid var(--black);
    cursor: pointer;
    transition: background-color 0.2s;
}

.toggle-switch::after {
    content: '';
    position: absolute;
    top: 1px;
    left: 1px;
    width: 1.125rem;
    height: 1.125rem;
    background-color: var(--white);
    border: 2px solid var(--black);
    transition: transform 0.2s;
}

.toggle-switch.active {
    background-color: var(--primary);
}

.toggle-switch.active::after {
    transform: translateX(1.5rem);
}
```

The toggle has **no border-radius** — it's a sliding square inside a rectangle.

---

## 14. HTMX Interaction Patterns

### Tab Switching

```html
<nav class="flex" id="detail-tabs">
    <button class="px-4 py-2 text-sm font-bold uppercase bg-black text-white border-2 border-black border-b-0"
            hx-get="/admin/products/123/specs"
            hx-target="#detail-content"
            hx-swap="innerHTML"
            onclick="setActiveTab(this)">
        Specs
    </button>
    <button class="px-4 py-2 text-sm font-bold uppercase bg-white text-black border-2 border-black border-b-0 border-l-0 hover:bg-gray-100"
            hx-get="/admin/products/123/features"
            hx-target="#detail-content"
            hx-swap="innerHTML"
            onclick="setActiveTab(this)">
        Features
    </button>
</nav>
<div id="detail-content"></div>
```

Active tab styling:
```javascript
function setActiveTab(el) {
    document.querySelectorAll('#detail-tabs button').forEach(function(btn) {
        btn.className = 'px-4 py-2 text-sm font-bold uppercase bg-white text-black border-2 border-black border-b-0 border-l-0 hover:bg-gray-100';
    });
    el.className = 'px-4 py-2 text-sm font-bold uppercase bg-black text-white border-2 border-black border-b-0';
}
```

### Delete with Confirmation

```html
<button hx-delete="/admin/products/123"
        hx-confirm="Delete this product?"
        hx-target="closest tr"
        hx-swap="outerHTML swap:0.3s"
        class="btn-danger">
    Delete
</button>
```

### Loading State

**Public Styles (styles.css):**
```css
.htmx-request {
    opacity: 0.7;
    pointer-events: none;
}
```

This class is automatically applied by HTMX to elements making requests.

---

## 15. Tailwind Config Overrides

### Location

Tailwind configuration is defined inline in the public site base template at `/templates/public/layouts/base.html` within a `<script>` tag. There is no separate `tailwind.config.js` file.

Admin templates use Tailwind CDN without custom configuration overrides.

### Border Radius Override (Public Site Only)

```javascript
tailwind.config = {
    theme: {
        extend: {
            borderRadius: {
                "DEFAULT": "0px",
                "lg": "0px",
                "xl": "0px",
                "full": "9999px"
            }
        }
    }
}
```

This ensures that ALL Tailwind utilities generate sharp corners on the public site.

### Custom Colors (Public Site Only)

```javascript
colors: {
    "primary": "#0066CC",
    "primary-dark": "#004499",
    "navy-tech": "#004499",
    "background-light": "#F8F9FA",
    "background-dark": "#1A1A2E",
    "text-primary": "#333333",
    "text-secondary": "#666666",
}
```

### Font Families (Public Site Only)

```javascript
fontFamily: {
    "display": ["Inter", "sans-serif"],
    "mono": ["JetBrains Mono", "monospace"]
}
```

---

## 16. Icons

### Material Symbols Outlined

**Admin Templates:**
```html
<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200" rel="stylesheet">
```

**Public Templates:**
```html
<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined" rel="stylesheet">
```

The admin template uses additional font parameters for more customization options (optical size, weight, fill, and grade variations).

Usage:
```html
<span class="material-symbols-outlined">dashboard</span>
<span class="material-symbols-outlined text-2xl">settings</span>
```

### Common Icons

- `dashboard` — Dashboard
- `inventory_2` — Products
- `article` — Blog posts
- `lightbulb` — Solutions
- `handshake` — Partners
- `settings` — Settings
- `add_box` — Create new
- `edit` — Edit
- `delete` — Delete
- `chevron_right` — Navigation arrow
- `expand_more` — Dropdown arrow
- `open_in_new` — External link

---

## 17. Responsive Patterns

### Mobile Sidebar Toggle

```html
<!-- Hamburger menu button (mobile only) -->
<button onclick="toggleSidebar()" class="lg:hidden">
    <span class="material-symbols-outlined">menu</span>
</button>

<!-- Sidebar overlay -->
<div class="sidebar-overlay" onclick="toggleSidebar()"></div>
```

### Grid Breakpoints

Standard Tailwind breakpoints are used:
- `sm:` — 640px and up
- `md:` — 768px and up
- `lg:` — 1024px and up
- `xl:` — 1280px and up

```html
<!-- 1 column on mobile, 2 on tablet, 3 on desktop -->
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
```

### Font Size Adjustments

**Public Styles (styles.css):**
```css
@media (max-width: 768px) {
    html { font-size: 14px; }
}
```

On mobile devices, the base font size reduces from 16px to 14px for better readability on smaller screens.

---

## 18. File Structure and Component Usage

### CSS File Organization

**public/css/admin-styles.css**
- Used by admin panel templates (`/templates/admin/**`)
- Contains: button components, form elements, badges, cards, toggle switches, tooltips, sidebar navigation
- Includes full CSS variable set (primary, danger, success, warning, etc.)
- Loaded in `/templates/admin/layouts/base.html`

**public/css/styles.css**
- Used by public site templates (`/templates/public/**`)
- Contains: basic utilities (manual-border, manual-shadow, btn-press, grid-dotted)
- Includes minimal CSS variable set (primary, navy-tech, black, white only)
- Loaded in `/templates/public/layouts/base.html`

### Component Usage Patterns

**Admin Templates:**
- Use Tailwind utilities for most styling: `border-2 border-black`, `bg-white`, etc.
- Use inline `style="box-shadow: 4px 4px 0px #000;"` for shadows (since `.manual-shadow` is not in admin-styles.css)
- Use CSS classes from admin-styles.css: `.btn-primary`, `.btn-danger`, `.card-brutal`, `.sidebar-link`, etc.
- Interactive elements use HTMX for dynamic content loading

**Public Templates:**
- Use Tailwind with custom config (inline in base.html)
- Use utility classes from styles.css: `.manual-border`, `.manual-shadow`, `.btn-press`
- Custom components should follow the brutalist design principles

### Interactive Behavior

**Sidebar Navigation** (`/public/js/admin.js`):
- `toggleGroup(groupName)` — Expands/collapses sidebar groups
- `toggleSidebar()` — Shows/hides mobile sidebar drawer
- State persistence using localStorage (`bluejay_sidebar_groups` key)
- Auto-expands group containing active page

**Form Interactions:**
- Focus states show blue ring: `box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1)`
- HTMX loading states reduce opacity to 0.7

## 19. How to Create a New Component

### Checklist

- [ ] **NO border-radius** (except circles)
- [ ] **2px solid black borders**
- [ ] **4px 4px 0 0 shadow** for standard elements
- [ ] **Uppercase text** for buttons and labels
- [ ] **JetBrains Mono** font
- [ ] **Bold weight (700)** for interactive elements
- [ ] **Push-down effect** on active state
- [ ] **High contrast** — avoid subtle grays
- [ ] **Sharp transitions** — 0.15s or less
- [ ] **Hard edges** — everything is angular

### Example: New Card Component

```css
.card-feature {
    background-color: var(--white);
    border: 2px solid var(--black);
    border-radius: 0;
    box-shadow: 4px 4px 0 0 var(--black);
    padding: 2rem;
    transition: background-color 0.15s, transform 0.1s;
}

.card-feature:hover {
    background-color: var(--primary);
    color: var(--white);
}

.card-feature:active {
    transform: translate(2px, 2px);
    box-shadow: 2px 2px 0 0 var(--black);
}
```

---

## 20. Do's and Don'ts

### ✅ DO

- Use sharp corners everywhere
- Use 2px solid borders
- Use hard, offset shadows (no blur)
- Use uppercase text for UI elements
- Use JetBrains Mono for consistency
- Use high contrast (black/white primary)
- Use bold weights for emphasis
- Use the push-down effect for buttons
- Use monospace for technical feel
- Use grid patterns for backgrounds
- Use technical labels (FIG 01, REF. MANUAL, etc.)

### ❌ DON'T

- Use border-radius (except full circles)
- Use soft, blurred shadows
- Use gradients
- Use rounded buttons
- Use subtle grays for borders
- Use thin font weights for UI
- Use mixed font families
- Use decorative elements
- Use animations longer than 0.3s
- Use pastel colors
- Use lowercase for buttons
- Use serif fonts

### Common Mistakes

**WRONG:**
```css
button {
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    font-family: 'Helvetica', sans-serif;
}
```

**CORRECT:**
```css
button {
    border-radius: 0;
    box-shadow: 4px 4px 0 0 var(--black);
    font-family: 'JetBrains Mono', monospace;
}
```

---

## 21. Accessibility Notes

While brutalist design is bold, it must still be accessible:

- **Focus states** — Always visible (blue ring or border change)
- **Contrast ratios** — Black/white provides AAA contrast
- **Font size** — Never smaller than 12px (0.75rem)
- **Touch targets** — Minimum 44x44px for mobile
- **Semantic HTML** — Use proper heading hierarchy
- **Alt text** — Always provide for images
- **ARIA labels** — Add for icon-only buttons

---

## Version History

- **v2026.01** — Initial design system documentation (Feb 10, 2026)

---

---

## 22. Quick Reference: Where to Find CSS Classes

| Class | Location | Used In |
|-------|----------|---------|
| `.btn-primary`, `.btn-secondary`, `.btn-danger`, `.btn-ghost` | admin-styles.css | Admin templates |
| `.input-brutal`, `.select-brutal`, `.textarea-brutal` | admin-styles.css | Admin forms |
| `.card-brutal`, `.card-brutal-sm` | admin-styles.css | Admin templates |
| `.badge-published`, `.badge-draft`, `.badge-pending`, `.badge-error` | admin-styles.css | Admin templates |
| `.toggle-switch` | admin-styles.css | Admin forms |
| `.tooltip-trigger`, `.tooltip-icon`, `.tooltip-content` | admin-styles.css | Admin templates |
| `.sidebar-link`, `.sidebar-group-header`, `.sidebar-sublink` | admin-styles.css | Admin sidebar |
| `.manual-border`, `.manual-border-white` | styles.css | Public templates |
| `.manual-shadow`, `.manual-shadow-white` | styles.css | Public templates |
| `.btn-press` | styles.css | Public templates |
| `.grid-dotted` | styles.css | Public templates |
| `.fade-in` | styles.css | Public templates (animation) |
| `.htmx-request` | styles.css | Automatic HTMX class |
| `.font-display` | admin-styles.css | Both (Inter font) |

---

**Questions?** Check existing components in `/templates` and `/public/css` for reference implementations.
