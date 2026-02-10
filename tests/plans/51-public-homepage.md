# Test Plan: Public Homepage

## Summary
Verify homepage displays all sections correctly with seeded content and interactive carousel navigation.

## Preconditions
- Server running on localhost:28090
- Database seeded with 1 hero, 4 stats, 3 testimonials, 1 CTA, featured products, and featured partners
- No authentication required

## User Journey Steps
1. Navigate to GET /
2. Scroll through homepage sections
3. Interact with testimonial carousel using JS navigation
4. Click links to product details, solutions, blog posts

## Test Cases

### Happy Path
- **Homepage loads successfully**: Verify GET / returns 200 status
- **Hero section displays**: Verify headline, subheadline, badge, CTAs, background image, slide indicators present
- **Featured Products section**: Verify 3-column grid with product cards linking to /products/{category}/{slug}
- **Solutions By Industry section**: Verify 6-column icon grid with links to /solutions/{slug}
- **Company Statistics section**: Verify 4-stat bar with seeded statistics
- **Client Testimonials carousel**: Verify 3 testimonials with rating stars, JS navigation controls, and indicators
- **Trusted Partners section**: Verify 5-column logo grid with partner logos
- **Latest News section**: Verify 3-column blog card grid linking to /blog/{slug}
- **CTA section displays**: Verify final CTA section with call-to-action content
- **Product detail link navigation**: Click product card, verify navigation to /products/{category}/{slug}
- **Solution link navigation**: Click solution icon, verify navigation to /solutions/{slug}
- **Blog post link navigation**: Click blog card, verify navigation to /blog/{slug}

### Edge Cases / Error States
- **Empty seeded content**: Verify sections handle missing optional content gracefully
- **Carousel at boundaries**: Click previous on first testimonial, click next on last testimonial
- **Long content overflow**: Verify text truncation or wrapping in cards

## Selectors & Elements
- Hero section: headline text, subheadline text, badge element, CTA buttons, background image, slide indicators
- Featured Products: `.featured-products` container, 3-column grid, product card links with href pattern `/products/*/`
- Solutions grid: 6-column layout, icon elements, solution links with href pattern `/solutions/*`
- Statistics bar: 4 stat elements with value and label
- Testimonials carousel: testimonial cards, rating stars (5-star display), navigation buttons (prev/next), carousel indicators
- Partners grid: 5-column layout, partner logo images
- News section: 3-column grid, blog cards with image, title, excerpt, link to `/blog/*`
- CTA section: CTA content container

## HTMX Interactions
- None on homepage (static content display)

## Dependencies
- JetBrains Mono font loaded
- Brutalist design: 2px solid black borders, manual box-shadows (4px 4px 0px #000)
- JavaScript for testimonial carousel navigation
- Seeded database content for all sections
