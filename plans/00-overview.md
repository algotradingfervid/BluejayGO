# Admin Panel Redesign - Master Plan

## Philosophy
- Keep the retro brutalist aesthetic (JetBrains Mono, thick borders, manual shadows)
- Reduce cognitive overload: fewer elements per page, clear visual hierarchy
- Every form field gets a hover tooltip explaining what it does
- Settings live where they belong (per-section in sidebar, not one giant page)
- Mockups are directional reference, not pixel-perfect targets

## Plan Hierarchy (Execute in Order)

### Foundation (Must do first)
1. **Phase 01** - Base Layout & Design System (`01-base-layout.md`)
2. **Phase 02** - Sidebar Navigation Redesign (`02-sidebar-navigation.md`)

### Core Pages
3. **Phase 03** - Dashboard Homepage (`03-dashboard.md`)
4. **Phase 04** - Global Settings Page (`04-global-settings.md`)
5. **Phase 05** - Header Management (`05-header.md`)
6. **Phase 06** - Footer Management (`06-footer.md`)

### Content Management
7. **Phase 07** - Products List & Form (`07-products.md`)
8. **Phase 08** - Product Detail Sub-pages (`08-product-details.md`)
9. **Phase 09** - Blog Posts List & Form (`09-blog.md`)
10. **Phase 10** - Blog Categories, Authors, Tags (`10-blog-taxonomy.md`)
11. **Phase 11** - Solutions List & Form (`11-solutions.md`)
12. **Phase 12** - Case Studies List & Form (`12-case-studies.md`)
13. **Phase 13** - Whitepapers List & Form (`13-whitepapers.md`)
14. **Phase 14** - Partners & Testimonials (`14-partners-testimonials.md`)
15. **Phase 15** - About Page Sections (`15-about.md`)
16. **Phase 16** - Homepage Customization (`16-homepage.md`)

### Admin Features
17. **Phase 17** - Contact & Form Submissions (`17-contact-forms.md`)
18. **Phase 18** - Media Library (NEW) (`18-media-library.md`)
19. **Phase 19** - Navigation Editor (`19-navigation.md`)
20. **Phase 20** - Activity Log / Audit Trail (NEW) (`20-audit-trail.md`)
21. **Phase 21** - Content Preview System (NEW) (`21-content-preview.md`)

### Auth & Users
22. **Phase 22** - Login & Auth Pages (`22-auth.md`)

## New Features Being Added
- **Media Library**: Browse/reuse uploaded images across all content
- **Audit Trail**: Basic log of who did what and when
- **Content Preview**: "Preview" button opens public page in new tab with draft content
- **Header Management**: Logo upload, nav links, CTA button, contact display, social icons
- **Footer Management**: Column layout, links, social icons, copyright
- **Form Submissions Viewer**: See contact/quote request submissions from public site
- **Per-section Settings**: Each sidebar group has its own settings instead of one giant page

## Tech Stack (No Changes)
- Go/Echo backend, HTML templates, Tailwind CSS, HTMX, Trix editor
- SQLite database with sqlc
- No new frontend frameworks
