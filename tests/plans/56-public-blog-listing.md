# Test Plan: Public Blog Listing

## Summary
Verify blog listing page displays posts with filtering by category/tag and pagination.

## Preconditions
- Server running on localhost:28090
- Database seeded with 7+ blog posts across multiple categories
- No authentication required

## User Journey Steps
1. Navigate to GET /blog
2. View featured post hero
3. Click category filter tabs
4. View filtered posts in 3-column grid
5. Navigate through pagination
6. Click post to view detail
7. Filter by tag using tag links

## Test Cases

### Happy Path
- **Blog listing page loads**: Verify GET /blog returns 200 status
- **Breadcrumb displays**: Verify breadcrumb navigation
- **Featured post hero**: Verify 2-column featured post card displays
- **Category filter tabs**: Verify "All Posts" tab plus per-category filter buttons
- **Posts grid displays**: Verify 3-column grid of blog post cards
- **Post card content**: Verify each card has image, category badge, title, excerpt, author avatar/name, publish date, reading time
- **Category filtering**: Click category button, verify navigation to /blog?category={slug}
- **Tag filtering**: Click tag link, verify navigation to /blog?tag={slug}
- **Post navigation**: Click post card, verify navigation to /blog/{slug}
- **Pagination displays**: Verify pagination controls when posts exceed page limit
- **Page navigation**: Click pagination link, verify navigation to /blog?page={number}

### Edge Cases / Error States
- **Invalid category filter**: Navigate to /blog?category=invalid, verify handling
- **Invalid tag filter**: Navigate to /blog?tag=invalid, verify handling
- **Invalid page number**: Navigate to /blog?page=999, verify handling
- **Empty category**: Filter by category with no posts, verify appropriate message
- **Last page**: Navigate to last page, verify "next" pagination disabled
- **First page**: Verify "previous" pagination disabled on page 1

## Selectors & Elements
- Breadcrumb: breadcrumb navigation
- Featured post: 2-column hero card with image, title, excerpt, link to `/blog/*`
- Category filters: tab/button container, "All Posts" button, category buttons
- Posts grid: 3-column grid container
- Post cards: card elements with:
  - Post image
  - Category badge with category-specific color
  - Post title
  - Excerpt text
  - Author avatar image
  - Author name
  - Publish date
  - Reading time indicator
  - Link to `/blog/*`
- Tag links: links with href pattern `/blog?tag=*`
- Pagination: pagination container with page links

## HTMX Interactions
- None (traditional page navigation with query parameters)

## Dependencies
- Template data: featured post, posts collection, categories, pagination info
- Seeded database with 7+ posts across categories and tags
- Category color mapping (CategoryColor)
- Brutalist design system applied
- JetBrains Mono font
