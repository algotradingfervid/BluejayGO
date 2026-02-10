# Test Plan: Public Case Studies

## Summary
Verify case studies listing and detail pages display client stories with challenges, solutions, outcomes, and metrics.

## Preconditions
- Server running on localhost:28090
- Database seeded with 3+ case studies including metrics and linked products
- No authentication required

## User Journey Steps
1. Navigate to GET /case-studies for listing
2. View case study cards in grid
3. Click case study card to view detail
4. Navigate to GET /case-studies/:slug
5. View breadcrumb and hero section
6. Read challenge, solution, and outcome sections
7. Review metrics and linked products

## Test Cases

### Happy Path - Listing
- **Case studies listing loads**: Verify GET /case-studies returns 200 status
- **Grid layout displays**: Verify case study cards in grid layout
- **Card content**: Verify each card has title, client_name, industry, summary
- **Case study links**: Verify cards link to /case-studies/:slug
- **Card navigation**: Click card, verify navigation to detail page

### Happy Path - Detail
- **Case study detail loads**: Verify GET /case-studies/:slug returns 200 status
- **Breadcrumb displays**: Verify breadcrumb navigation
- **Hero section**: Verify case study hero with title and client information
- **Challenge section**: Verify title, content, and bullet points
- **Solution section**: Verify solution description and details
- **Outcome section**: Verify outcomes and results
- **Metrics display**: Verify metrics with value and label badges
- **Linked products**: Verify related products section with links to product detail pages
- **Product navigation**: Click linked product, verify navigation to /products/:category/:slug

### Edge Cases / Error States
- **Empty listing**: Verify listing page handles no case studies gracefully
- **Case study not found**: Navigate to invalid slug, verify 404 or error page
- **No metrics**: Verify metrics section handles empty metrics
- **No linked products**: Verify products section handles empty list
- **Long content**: Verify proper text wrapping for long challenge/solution/outcome content

## Selectors & Elements
- Listing page:
  - Grid container
  - Case study cards with title, client_name, industry, summary
  - Links to `/case-studies/*`
- Detail page:
  - Breadcrumb navigation
  - Hero section with title and client info
  - Challenge section: title, content text, bullet point list
  - Solution section: solution content
  - Outcome section: outcome content
  - Metrics section: metric badges with value and label
  - Linked products: products container with links to `/products/*/*`

## HTMX Interactions
- None (static content display)

## Dependencies
- Template data: case studies collection, metrics, linked products
- Seeded database with 3+ case studies
- Brutalist design system applied
- JetBrains Mono font
