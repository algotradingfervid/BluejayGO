# Test Plan: Public Products Listing

## Summary
Verify products listing page displays categories and supports HTMX search functionality.

## Preconditions
- Server running on localhost:28090
- Database seeded with product categories and products
- No authentication required

## User Journey Steps
1. Navigate to GET /products
2. View breadcrumb navigation
3. Enter search query in search input
4. View HTMX search results update
5. Click category card to navigate to category page

## Test Cases

### Happy Path
- **Products page loads**: Verify GET /products returns 200 status
- **Breadcrumb displays**: Verify "Home > Products" breadcrumb navigation
- **Page header renders**: Verify page header with search input
- **Search input configured**: Verify input has hx-get="/products/search", hx-trigger="keyup changed delay:300ms", hx-target="#product-results"
- **Category cards display**: Verify 2-column grid of category cards with links to /products/{category-slug}
- **CTA section displays**: Verify page CTA section present
- **HTMX search triggers**: Type in search input, verify HTMX request to /products/search after 300ms delay
- **Search results update**: Verify #product-results div updates with search results HTML
- **Category navigation**: Click category card, verify navigation to /products/{category-slug}

### Edge Cases / Error States
- **Empty search query**: Verify behavior when search input is empty
- **No search results**: Enter query with no matches, verify appropriate message in #product-results
- **Fast typing debounce**: Type rapidly, verify only final query triggers search after delay
- **Special characters in search**: Test search with special characters, verify proper handling

## Selectors & Elements
- Breadcrumb: text "Home > Products"
- Page header: container with search input
- Search input: `input[hx-get="/products/search"][hx-trigger="keyup changed delay:300ms"][hx-target="#product-results"]`
- Search results container: `#product-results`
- Category grid: 2-column grid container
- Category cards: links with href pattern `/products/*`
- CTA section: page CTA container

## HTMX Interactions
- **Search input**: hx-get="/products/search", hx-trigger="keyup changed delay:300ms", hx-target="#product-results"
- **Expected behavior**: Typing triggers debounced HTMX request, response replaces #product-results content
- **Response type**: HTML fragment (partial)

## Dependencies
- HTMX library loaded
- Template data: PageHero, Query, Categories, PageCTA
- Seeded product categories and products
- Brutalist design system applied
