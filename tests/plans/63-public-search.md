# Test Plan: Public Search

## Summary
Verify full-text search across products, blog posts, and case studies with HTMX autocomplete suggestions.

## Preconditions
- Server running on localhost:28090
- Database with FTS5 virtual tables for products, blog posts, and case studies
- No authentication required

## User Journey Steps
1. Navigate to GET /search (or access search from navigation)
2. Enter search query in search input
3. View HTMX autocomplete suggestions
4. Submit search with query parameter
5. Navigate to GET /search?q={query}
6. View grouped search results by type
7. Click result to navigate to detail page

## Test Cases

### Happy Path - Search Page
- **Search page loads**: Verify GET /search returns 200 status
- **Search input displays**: Verify search input field present
- **Search with query**: Submit search, verify navigation to GET /search?q={query}
- **Results display**: Verify search results grouped by type (Products, Blog Posts, Case Studies)
- **Product results**: Verify matching products with name, tagline, description snippets
- **Blog post results**: Verify matching posts with title, excerpt, body snippets
- **Case study results**: Verify matching case studies with title, client_name, content snippets
- **Result links**: Verify each result links to appropriate detail page
- **Product link navigation**: Click product result, verify navigation to /products/:category/:slug
- **Blog link navigation**: Click blog result, verify navigation to /blog/:slug
- **Case study link navigation**: Click case study result, verify navigation to /case-studies/:slug

### Happy Path - Autocomplete
- **Autocomplete input configured**: Verify input has hx-get="/search/suggest", hx-trigger with debounce
- **Autocomplete triggers**: Type in search input, verify HTMX request to GET /search/suggest?q={query}
- **Suggestions display**: Verify search_suggestions.html partial renders with suggestions
- **Suggestion selection**: Click suggestion, verify search executes with selected term

### Edge Cases / Error States
- **Empty query**: Submit empty search, verify appropriate handling
- **No results**: Search for non-existent term, verify "no results" message
- **Very short query**: Search with 1-2 characters, verify behavior
- **Special characters**: Search with special characters, verify proper handling
- **FTS5 operators**: Search with FTS5 syntax (quotes, AND/OR), verify proper parsing
- **Long query**: Search with very long query string, verify handling
- **Autocomplete debounce**: Type rapidly, verify only final query triggers suggestion request

## Selectors & Elements
- Search page:
  - Search input field
  - Search submit button
  - Results container
  - Results grouped by type: Products section, Blog Posts section, Case Studies section
  - Result items with title, snippet, link
- Autocomplete:
  - Search input with HTMX attributes: `hx-get="/search/suggest"`, debounced trigger, target for suggestions
  - Suggestions container (rendered from search_suggestions.html partial)
  - Suggestion items (clickable)

## HTMX Interactions
- **Autocomplete input**: hx-get="/search/suggest?q={query}", debounced trigger (e.g., "keyup changed delay:300ms"), target for suggestions dropdown
- **Response**: search_suggestions.html partial with suggestion items

## Dependencies
- FTS5 virtual tables for full-text search
- Search handler: GET /search?q={query}
- Autocomplete handler: GET /search/suggest?q={query}
- Template: search_suggestions.html partial
- Seeded database with searchable products, blog posts, case studies
- HTMX library loaded
- Brutalist design system applied
- JetBrains Mono font
