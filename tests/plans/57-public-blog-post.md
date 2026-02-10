# Test Plan: Public Blog Post Detail

## Summary
Verify blog post detail page displays complete post content with author bio, tags, and related products.

## Preconditions
- Server running on localhost:28090
- Database seeded with blog posts including tags, author information, and related products
- No authentication required

## User Journey Steps
1. Navigate to GET /blog/:slug
2. View breadcrumb navigation
3. Read post header with metadata
4. View featured image
5. Read article body content
6. Click tag links for filtering
7. View author bio card
8. View related products section

## Test Cases

### Happy Path
- **Blog post page loads**: Verify GET /blog/:slug returns 200 status
- **Breadcrumb displays**: Verify "Home > Blog > Category > Title" breadcrumb
- **Post header renders**: Verify category badge (with CategoryColor), reading time, title, excerpt
- **Author metadata**: Verify author avatar, author name, publish date display
- **Featured image displays**: Verify post featured image renders
- **Article body renders**: Verify post content displays with safeHTML rendering
- **Tags section displays**: Verify post tags with links to /blog?tag={slug}
- **Author bio card**: Verify author avatar, name, title, bio, LinkedIn link
- **Related products section**: Verify related products display with links
- **Tag navigation**: Click tag link, verify navigation to /blog?tag={slug}
- **LinkedIn link**: Verify author LinkedIn link opens correctly

### Edge Cases / Error States
- **Post not found**: Navigate to invalid slug, verify 404 or error page
- **No tags**: Verify tags section handles posts without tags
- **No related products**: Verify related products section handles empty list
- **No author LinkedIn**: Verify author bio handles missing LinkedIn URL
- **HTML content safety**: Verify safeHTML prevents XSS in article body
- **Long article body**: Verify proper rendering of very long content

## Selectors & Elements
- Breadcrumb: text pattern "Home > Blog > * > *"
- Post header:
  - Category badge with background color from CategoryColor
  - Reading time indicator
  - Post title heading
  - Excerpt text
  - Author avatar image
  - Author name
  - Publish date
- Featured image: post image element
- Article body: content container with HTML content
- Tags section: tags container with tag links to `/blog?tag=*`
- Author bio card:
  - Author avatar
  - Author name
  - Author title/role
  - Author bio text
  - LinkedIn link (if available)
- Related products: products container with product cards/links

## HTMX Interactions
- None (static content display)

## Dependencies
- Template data: post details, author info, tags, related products
- CategoryColor mapping for badge styling
- safeHTML function for secure content rendering
- Seeded posts with tags and related data
- Brutalist design system applied
- JetBrains Mono font
