# Test Plan: About Overview Management

## Summary
Tests the company overview page editing functionality including headline, taglines, descriptions, and image URLs.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- About overview section seeded with company data
- Valid image URLs for hero and company images

## User Journey Steps
1. Navigate to /admin/about/overview
2. View pre-populated form with existing overview data
3. Edit headline, tagline, and description fields
4. Update hero and company image URLs
5. Submit form
6. Verify redirect to /admin/about/overview?saved=1
7. Confirm success message appears
8. Verify updated data persists on page reload

## Test Cases

### Happy Path
- **Load overview form**: Verifies GET /admin/about/overview loads with all fields populated from seeded data
- **Update all fields**: Updates headline, tagline, all three description fields, and both image URLs, verifies successful save
- **Update single field**: Changes only headline, verifies other fields remain unchanged
- **Success indicator**: Verifies ?saved=1 query param shows success banner
- **Cache invalidation**: Confirms page:about cache is cleared after update

### Edge Cases / Error States
- **Empty required fields**: Tests validation when headline or tagline is empty
- **Long text handling**: Enters very long text in textarea fields, verifies truncation or validation
- **Invalid image URLs**: Enters malformed URLs in image fields, checks validation
- **Concurrent edits**: Two users editing simultaneously, last save wins
- **XSS prevention**: Enters script tags in text fields, verifies sanitization

## Selectors & Elements
- Form: `form[action="/admin/about/overview"][method="POST"]`
- Headline input: `input[name="headline"]`
- Tagline input: `input[name="tagline"]`
- Description main: `textarea[name="description_main"]`
- Description secondary: `textarea[name="description_secondary"]`
- Description tertiary: `textarea[name="description_tertiary"]`
- Hero image URL: `input[name="hero_image_url"]`
- Company image URL: `input[name="company_image_url"]`
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` or `[data-saved="true"]` (check for ?saved=1)

## HTMX Interactions
- None - standard form POST with full page redirect

## Dependencies
- Database seeded with about overview data
- Cache service configured for page:about
- Template: templates/admin/pages/about-overview.html
- Handler: internal/handlers/about.go (GetAboutOverview, PostAboutOverview)
