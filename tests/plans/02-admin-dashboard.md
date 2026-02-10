# Test Plan: Admin Dashboard

## Summary
Verify admin dashboard displays correct statistics, navigation, and quick action links.

## Preconditions
- User authenticated with valid session cookie
- Database contains seeded content (products, blog posts, partners, contact submissions)
- Server running on localhost:28090

## User Journey Steps
1. Navigate to http://localhost:28090/admin/dashboard (or be redirected after login)
2. Verify page title shows "Dashboard"
3. Verify statistics cards display counts for: PublishedProducts, PublishedBlogPosts, ContactSubmissions, NewContactSubmissions, TotalPartners, DraftProducts, DraftBlogPosts
4. Verify sidebar navigation is present with collapsible groups
5. Click on navigation group to test toggleGroup() JavaScript
6. Verify quick action links to create content are present and functional

## Test Cases

### Happy Path
- **Dashboard loads successfully**: Authenticated user can access dashboard with all stat cards visible
- **Statistics display correct counts**: Each stat card shows accurate count from database
- **Sidebar navigation functional**: All nav links are clickable and collapsible groups toggle correctly
- **Quick action links work**: "Create Product", "Create Blog Post", etc. navigate to correct creation forms
- **Page title correct**: H1 or page title element shows "Dashboard"

### Edge Cases / Error States
- **Empty database stats**: Dashboard shows "0" for all stats when database is empty
- **Large numbers formatting**: Stats display correctly formatted when counts exceed 1000
- **Navigation group persistence**: Collapsed/expanded state persists during session (if implemented)

## Selectors & Elements
- Page title: text "Dashboard"
- Stat cards: containers showing count and label for each metric
  - PublishedProducts
  - PublishedBlogPosts
  - ContactSubmissions
  - NewContactSubmissions
  - TotalPartners
  - DraftProducts
  - DraftBlogPosts
- Sidebar navigation: nav element with groups having `data-group` attribute
- Collapsible group toggle: JavaScript function `toggleGroup()`
- Quick action links: links to /admin/products/new, /admin/blog-posts/new, etc.

## HTMX Interactions
- No HTMX on dashboard page (static content load)
- Stats loaded server-side on GET request

## Dependencies
- 01-admin-login-logout.md (requires authenticated session)
- 03-product-categories-crud.md, 09-products-crud.md (products data for stats)
- Blog posts test plan (blog stats - not in current set)
- Contact submissions test plan (contact stats - not in current set)
