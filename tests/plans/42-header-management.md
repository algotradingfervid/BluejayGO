# Test Plan: Header Configuration Management

## Summary
Tests the header management interface including logo settings, navigation visibility toggles, CTA configuration, contact display options, and social media settings.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Header settings initialized with default values
- Logo image uploaded and available
- Navigation sections: products, solutions, case_studies, about, blog, whitepapers, partners, contact

## User Journey Steps
1. Navigate to /admin/header
2. View logo section with header_logo_path, header_logo_alt, and logo preview
3. Update logo path and alt text
4. Toggle navigation visibility checkboxes (show_nav_products, show_nav_solutions, etc.)
5. Update navigation labels (nav_label_products, nav_label_solutions, etc.)
6. Configure CTA section: enable/disable, set text and URL, choose style (radio)
7. Toggle contact display (show_phone, show_email)
8. Configure social media: enable/disable, choose style (icons or icons_labels)
9. Submit form
10. Verify redirect to /admin/header?saved=1
11. Confirm success message and settings persist

## Test Cases

### Happy Path
- **Load header settings**: Verifies GET /admin/header loads with current configuration
- **Update logo path**: Changes header_logo_path, verifies save and preview update
- **Update logo alt text**: Changes header_logo_alt, verifies save
- **Toggle all nav items on**: Checks all show_nav_* checkboxes, verifies save
- **Toggle all nav items off**: Unchecks all show_nav_* checkboxes, verifies save
- **Update nav labels**: Changes nav_label_products to "Our Products", verifies save
- **Enable header CTA**: Checks header_cta_enabled, fills text and URL, verifies save
- **Disable header CTA**: Unchecks header_cta_enabled, verifies CTA hidden
- **Set CTA style primary**: Selects header_cta_style=primary radio, verifies save
- **Set CTA style secondary**: Selects header_cta_style=secondary radio, verifies save
- **Show contact phone**: Checks header_show_phone, verifies save
- **Show contact email**: Checks header_show_email, verifies save
- **Enable social with icons only**: Checks header_show_social, selects icons style, verifies save
- **Enable social with labels**: Checks header_show_social, selects icons_labels style, verifies save
- **Disable social**: Unchecks header_show_social, verifies social hidden

### Edge Cases / Error States
- **Invalid logo path**: Enters non-existent logo path, checks validation/preview error
- **Empty logo alt text**: Leaves header_logo_alt empty, checks if required or optional
- **All nav items hidden**: Unchecks all show_nav_* checkboxes, verifies header navigation behavior
- **Empty nav label**: Leaves nav_label_products empty, checks default or validation
- **CTA enabled without text**: Checks header_cta_enabled but leaves header_cta_text empty, checks validation
- **CTA enabled without URL**: Checks header_cta_enabled but leaves header_cta_url empty, checks validation
- **Invalid CTA URL**: Enters malformed header_cta_url, checks validation
- **CTA style not selected**: Enables CTA without selecting style radio, checks default
- **Both contact options off**: Unchecks both header_show_phone and header_show_email, verifies display
- **Social enabled without style**: Checks header_show_social without selecting style, checks default
- **Very long nav label**: Enters 100+ char nav label, checks truncation or limit
- **Logo preview**: Verifies logo preview image updates when header_logo_path changes

## Selectors & Elements
- Form: `form[action="/admin/header"][method="POST"]`
- Logo path input: `input[name="header_logo_path"]`
- Logo alt text input: `input[name="header_logo_alt"]`
- Logo preview image: `img#logo-preview` or `.logo-preview`
- Nav show checkboxes: `input[name="show_nav_products"][type="checkbox"]`, `input[name="show_nav_solutions"]`, etc.
- Nav label inputs: `input[name="nav_label_products"]`, `input[name="nav_label_solutions"]`, etc.
- CTA enabled checkbox: `input[name="header_cta_enabled"][type="checkbox"]`
- CTA text input: `input[name="header_cta_text"]`
- CTA URL input: `input[name="header_cta_url"]`
- CTA style radios: `input[name="header_cta_style"][value="primary"]`, `[value="secondary"]`
- Show phone checkbox: `input[name="header_show_phone"][type="checkbox"]`
- Show email checkbox: `input[name="header_show_email"][type="checkbox"]`
- Show social checkbox: `input[name="header_show_social"][type="checkbox"]`
- Social style radios: `input[name="header_social_style"][value="icons"]`, `[value="icons_labels"]`
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` (when ?saved=1)

## HTMX Interactions
- None - standard form POST with full page redirect
- Logo preview may update via JavaScript on path change

## Dependencies
- Database settings table with header configuration
- Logo image file accessible at header_logo_path
- Navigation sections defined (8 items: products, solutions, case_studies, about, blog, whitepapers, partners, contact)
- Template: templates/admin/pages/header-management.html
- Handler: internal/handlers/header.go (GetHeaderSettings, PostHeaderSettings)
