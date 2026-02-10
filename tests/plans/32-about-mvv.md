# Test Plan: Mission, Vision, Values Management

## Summary
Tests the editing of company mission, vision, and values summary with associated icon selections.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- MVV section pre-seeded with initial data
- Material Icons available for icon selection

## User Journey Steps
1. Navigate to /admin/about/mvv
2. View pre-populated mission, vision, and values summary
3. Edit mission statement textarea
4. Edit vision statement textarea
5. Edit values summary textarea
6. Update icon selections for each section
7. Submit form
8. Verify redirect to /admin/about/mvv?saved=1
9. Confirm success message and data persistence

## Test Cases

### Happy Path
- **Load MVV form**: Verifies GET /admin/about/mvv loads with seeded mission, vision, values data
- **Update mission only**: Changes mission text and icon, verifies vision and values unchanged
- **Update all sections**: Edits all three textareas and all three icons, verifies save
- **Icon selection**: Updates icon fields with valid Material icon names, verifies display
- **Success redirect**: Confirms ?saved=1 parameter and success banner appear

### Edge Cases / Error States
- **Empty mission**: Tests validation when mission textarea is cleared
- **Empty vision**: Tests validation when vision textarea is cleared
- **Invalid icon names**: Enters non-existent Material icon names, checks validation/fallback
- **Maximum textarea length**: Tests very long text in mission/vision/values fields
- **Special characters**: Enters markdown/HTML in textareas, verifies sanitization
- **Cache verification**: Confirms page:about cache invalidation after update

## Selectors & Elements
- Form: `form[action="/admin/about/mvv"][method="POST"]`
- Mission textarea: `textarea[name="mission"]`
- Vision textarea: `textarea[name="vision"]`
- Values summary textarea: `textarea[name="values_summary"]`
- Mission icon: `input[name="mission_icon"]`
- Vision icon: `input[name="vision_icon"]`
- Values icon: `input[name="values_icon"]`
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` (when ?saved=1 present)

## HTMX Interactions
- None - standard form POST with full page redirect

## Dependencies
- Database seeded with MVV data
- Cache service for page:about
- Material Icons CDN or icon validation
- Template: templates/admin/pages/about-mvv.html
- Handler: internal/handlers/about.go (GetAboutMVV, PostAboutMVV)
