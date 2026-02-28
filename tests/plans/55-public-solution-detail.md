# Test Plan: Public Solution Detail

## Summary
Verify solution detail page displays hero, stats grid, challenges, linked products, reference code, and CTAs.

**KNOWN BUG**: Product links in template use `/products/{{.ProductSlug}}` but should be `/products/{{.ProductCategorySlug}}/{{.ProductSlug}}` — this will cause 404 errors when clicking product links.

## Preconditions
- Server running on localhost:28090
- Database seeded with solutions including stats, challenges, and linked products
- No authentication required

## User Journey Steps
1. Navigate to GET /solutions/:slug
2. View breadcrumb and hero section
3. Review stats grid (3-column layout)
4. View reference code display
5. Read challenges section
6. View linked products
7. Interact with CTAs
8. Navigate to other solutions

## Test Cases

### Happy Path
- **Solution detail page loads**: Verify GET /solutions/:slug returns 200 status
- **Breadcrumb displays**: Verify breadcrumb navigation present
- **Hero section displays**: Verify title, description, and hero image
- **Stats grid displays**: Verify statistics in 3-column grid layout (`sm:grid-cols-3`) with value and label
- **Reference code display**: Verify reference code shown if available
- **Challenges section**: Verify title, description, and icon for each challenge
- **Products section**: Verify linked products display with links to product detail pages
- **CTAs section**: Verify call-to-action section with appropriate content
- **Other solutions links**: Verify links to other related solutions
- **Product navigation**: Click linked product, verify navigation to /products/:category/:slug (NOTE: Known bug - template uses incorrect URL pattern)

### Edge Cases / Error States
- **Solution not found**: Navigate to invalid slug, verify 404 or error page
- **No stats**: Verify stats grid handles empty statistics
- **No reference code**: Verify reference code section handles missing data
- **No challenges**: Verify challenges section handles empty list
- **No linked products**: Verify products section handles empty list
- **No related solutions**: Verify other solutions section handles empty list
- **Product link 404**: Due to known bug, verify product links using `/products/{{.ProductSlug}}` result in 404 errors

## Selectors & Elements
- Breadcrumb: breadcrumb navigation container
- Hero section: title heading, description text, hero image element
- Stats grid: statistics container with 3-column grid layout (`sm:grid-cols-3`), value and label elements
- Reference code: reference code display element
- Challenges section: challenges container, challenge title, challenge description, challenge icon
- Products section: products container, product cards with links to `/products/*/*` (NOTE: Bug - template uses `/products/{{.ProductSlug}}` instead)
- CTAs section: CTA container with action buttons/links
- Other solutions: related solutions links container with links to `/solutions/*`

## HTMX Interactions
- None (static content display)

## Dependencies
- Template data: solution details, stats, challenges, linked products
- Seeded database with complete solution data
- Brutalist design system applied
- JetBrains Mono font
