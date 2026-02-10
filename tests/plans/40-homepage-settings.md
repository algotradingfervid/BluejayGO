# Test Plan: Homepage Settings Management

## Summary
Tests the homepage settings configuration including visibility toggles, maximum display counts, and hero carousel autoplay settings.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Homepage settings initialized with default values
- Settings stored as int64 (checkboxes convert to 0/1)

## User Journey Steps
1. Navigate to /admin/homepage/settings
2. View current settings with all visibility toggles
3. Toggle visibility checkboxes (show_heroes, show_stats, show_testimonials, show_cta)
4. Update maximum counts (max_heroes, max_stats, max_testimonials)
5. Configure hero carousel settings (autoplay checkbox, interval number)
6. Submit form
7. Verify redirect to /admin/homepage/settings?saved=1
8. Confirm success message appears
9. Verify settings persist on page reload

## Test Cases

### Happy Path
- **Load settings form**: Verifies GET /admin/homepage/settings loads with current values
- **Toggle all visibility on**: Checks all show checkboxes, verifies save
- **Toggle all visibility off**: Unchecks all show checkboxes, verifies save
- **Update max heroes**: Changes homepage_max_heroes value, verifies save
- **Update max stats**: Changes homepage_max_stats value, verifies save
- **Update max testimonials**: Changes homepage_max_testimonials value, verifies save
- **Enable hero autoplay**: Checks homepage_hero_autoplay, sets interval, verifies save
- **Disable hero autoplay**: Unchecks homepage_hero_autoplay, verifies save
- **Change autoplay interval**: Updates homepage_hero_interval value, verifies save
- **Success redirect**: Confirms ?saved=1 parameter and success banner appear

### Edge Cases / Error States
- **Zero max values**: Sets max_heroes/stats/testimonials to 0, checks validation
- **Negative max values**: Enters negative numbers in max fields, checks validation
- **Very large max values**: Enters 999 or larger, checks validation/limits
- **Zero autoplay interval**: Sets homepage_hero_interval to 0, checks validation
- **Negative autoplay interval**: Enters negative interval, checks validation
- **Very large interval**: Enters 99999ms interval, checks validation
- **Autoplay enabled without interval**: Checks autoplay but leaves interval empty, checks default
- **All sections hidden**: Unchecks all show checkboxes, verifies homepage behavior
- **Checkbox to int64 conversion**: Verifies checkboxes correctly convert to 0/1 in database
- **Form validation**: Tests required field validation if any fields are required

## Selectors & Elements
- Form: `form[action="/admin/homepage/settings"][method="POST"]`
- Show heroes checkbox: `input[name="homepage_show_heroes"][type="checkbox"]`
- Show stats checkbox: `input[name="homepage_show_stats"][type="checkbox"]`
- Show testimonials checkbox: `input[name="homepage_show_testimonials"][type="checkbox"]`
- Show CTA checkbox: `input[name="homepage_show_cta"][type="checkbox"]`
- Max heroes input: `input[name="homepage_max_heroes"][type="number"]`
- Max stats input: `input[name="homepage_max_stats"][type="number"]`
- Max testimonials input: `input[name="homepage_max_testimonials"][type="number"]`
- Hero autoplay checkbox: `input[name="homepage_hero_autoplay"][type="checkbox"]`
- Hero interval input: `input[name="homepage_hero_interval"][type="number"]`
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` (when ?saved=1 present)

## HTMX Interactions
- None - standard form POST with full page redirect

## Dependencies
- Database settings table with homepage configuration
- Settings service/handler for get and update operations
- Template: templates/admin/pages/homepage-settings.html
- Handler: internal/handlers/homepage.go (GetHomepageSettings, PostHomepageSettings)
- Checkbox values correctly convert to int64 (1 for checked, 0 for unchecked)
