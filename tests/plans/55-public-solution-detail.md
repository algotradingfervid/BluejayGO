# Test Plan: Public Solution Detail

## Summary
Verify solution detail page displays hero, stats, challenges, linked products, and CTAs.

## Preconditions
- Server running on localhost:28090
- Database seeded with solutions including stats, challenges, and linked products
- No authentication required

## User Journey Steps
1. Navigate to GET /solutions/:slug
2. View breadcrumb and hero section
3. Review stats bar
4. Read challenges section
5. View linked products
6. Interact with CTAs
7. Navigate to other solutions

## Test Cases

### Happy Path
- **Solution detail page loads**: Verify GET /solutions/:slug returns 200 status
- **Breadcrumb displays**: Verify breadcrumb navigation present
- **Hero section displays**: Verify title, description, and hero image
- **Stats bar displays**: Verify statistics with value and label
- **Challenges section**: Verify title, description, and icon for each challenge
- **Products section**: Verify linked products display with links to product detail pages
- **CTAs section**: Verify call-to-action section with appropriate content
- **Other solutions links**: Verify links to other related solutions
- **Product navigation**: Click linked product, verify navigation to /products/:category/:slug

### Edge Cases / Error States
- **Solution not found**: Navigate to invalid slug, verify 404 or error page
- **No stats**: Verify stats bar handles empty statistics
- **No challenges**: Verify challenges section handles empty list
- **No linked products**: Verify products section handles empty list
- **No related solutions**: Verify other solutions section handles empty list

## Selectors & Elements
- Breadcrumb: breadcrumb navigation container
- Hero section: title heading, description text, hero image element
- Stats bar: statistics container with value and label elements
- Challenges section: challenges container, challenge title, challenge description, challenge icon
- Products section: products container, product cards with links to `/products/*/*`
- CTAs section: CTA container with action buttons/links
- Other solutions: related solutions links container with links to `/solutions/*`

## HTMX Interactions
- None (static content display)

## Dependencies
- Template data: solution details, stats, challenges, linked products
- Seeded database with complete solution data
- Brutalist design system applied
- JetBrains Mono font
