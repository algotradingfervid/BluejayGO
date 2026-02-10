# Test Plan: Sitemap and Robots.txt

## Summary
Verify dynamic XML sitemap lists all public pages and robots.txt provides crawl directives.

## Preconditions
- Server running on localhost:28090
- Database seeded with all public content types
- Base URL: https://bluejaylabs.com (configured in settings)
- No authentication required

## User Journey Steps
1. Navigate to GET /sitemap.xml
2. Verify XML structure and content URLs
3. Navigate to GET /robots.txt
4. Verify crawl directives

## Test Cases

### Happy Path - Sitemap
- **Sitemap loads**: Verify GET /sitemap.xml returns 200 status
- **XML content type**: Verify Content-Type header is application/xml or text/xml
- **Valid XML structure**: Verify XML is well-formed with urlset root element
- **Sitemap namespace**: Verify xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
- **Base URL correct**: Verify all URLs use base https://bluejaylabs.com
- **Static pages included**: Verify URLs for /, /products, /solutions, /blog, /case-studies, /whitepapers, /about, /contact, /partners
- **Product URLs included**: Verify all published products have URLs in format /products/{category}/{slug}
- **Solution URLs included**: Verify all published solutions have URLs in format /solutions/{slug}
- **Blog post URLs included**: Verify all published blog posts have URLs in format /blog/{slug}
- **Case study URLs included**: Verify all published case studies have URLs in format /case-studies/{slug}
- **Whitepaper URLs included**: Verify all published whitepapers have URLs in format /whitepapers/{slug}
- **URL elements**: Verify each URL entry has <loc> element
- **Optional elements**: Check for optional <lastmod>, <changefreq>, <priority> elements

### Happy Path - Robots.txt
- **Robots.txt loads**: Verify GET /robots.txt returns 200 status
- **Text content type**: Verify Content-Type header is text/plain
- **User-agent directive**: Verify User-agent: * or specific bot directives
- **Disallow directives**: Verify appropriate disallow rules (e.g., Disallow: /admin/)
- **Sitemap reference**: Verify Sitemap: https://bluejaylabs.com/sitemap.xml directive

### Edge Cases / Error States
- **Unpublished content excluded**: Verify draft/unpublished items not in sitemap
- **Deleted content excluded**: Verify deleted items not in sitemap
- **Empty content type**: If no blog posts exist, verify section handled gracefully
- **URL encoding**: Verify special characters in slugs are properly URL-encoded
- **Sitemap size**: Verify sitemap doesn't exceed 50MB or 50,000 URLs (standard limits)

## Selectors & Elements
- Sitemap XML:
  - Root element: `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
  - URL entries: `<url>` elements
  - Location: `<loc>` elements with full URLs
  - Optional: `<lastmod>`, `<changefreq>`, `<priority>` elements
- Robots.txt:
  - User-agent lines
  - Disallow/Allow lines
  - Sitemap reference line

## HTMX Interactions
- None (static file responses)

## Dependencies
- Sitemap handler: GET /sitemap.xml
- Robots.txt handler: GET /robots.txt
- Base URL configuration setting
- Database queries for all published content
- XML generation logic
- Proper filtering of published vs. unpublished content
