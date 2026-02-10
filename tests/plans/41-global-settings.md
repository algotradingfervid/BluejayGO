# Test Plan: Global Settings Management

## Summary
Tests the tabbed global settings interface including general site settings, SEO configuration, and social media links with character counters and validation.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Global settings initialized with default values
- JavaScript functions available: switchTab(), validateUrl(), showUnsavedBanner()
- Character counters functional for meta fields

## User Journey Steps
1. Navigate to /admin/settings?tab=general (default tab)
2. View general settings tab with site info and contact details
3. Fill in site_name, site_tagline, contact_email, contact_phone
4. Fill in address (textarea) and business_hours (textarea)
5. Switch to SEO tab using switchTab() JavaScript
6. Fill in meta_keywords (max 100 chars), meta_description (max 250 chars)
7. Fill in google_analytics_id
8. Observe character counters update
9. Switch to social media tab
10. Fill in social media URLs (Facebook, Twitter, LinkedIn, Instagram, YouTube)
11. Validate URLs using validateUrl() function
12. Submit form with hidden field active_tab
13. Verify redirect to /admin/settings?saved=1&tab=X
14. Confirm success banner appears
15. Test unsaved changes warning when navigating away

## Test Cases

### Happy Path
- **Load general tab**: Verifies GET /admin/settings?tab=general shows general settings
- **Load SEO tab**: Verifies GET /admin/settings?tab=seo shows SEO settings
- **Load social tab**: Verifies GET /admin/settings?tab=social shows social media settings
- **Default tab**: Navigates to /admin/settings without tab param, verifies general tab active
- **Update general settings**: Changes site_name and contact_email, verifies save
- **Update address**: Fills multi-line address textarea, verifies save
- **Update business hours**: Fills business_hours textarea, verifies save
- **Update SEO settings**: Changes meta_keywords and meta_description, verifies save
- **Update analytics ID**: Enters Google Analytics ID, verifies save
- **Character counter**: Types in meta_description, verifies character count updates
- **Max length enforcement**: Reaches 250 chars in meta_description, verifies limit
- **Update social links**: Fills all social media URLs, verifies save
- **Tab persistence**: Saves from SEO tab, verifies redirect to ?saved=1&tab=seo
- **Success banner**: Confirms .alert-success appears when ?saved=1 present

### Edge Cases / Error States
- **Empty site_name**: Tests required field validation
- **Invalid email format**: Enters malformed contact_email, checks validation
- **Invalid phone format**: Enters invalid contact_phone, checks validation
- **Meta keywords over 100**: Enters 150 chars, verifies truncation or validation
- **Meta description over 250**: Enters 300 chars, verifies truncation or validation
- **Invalid Google Analytics ID**: Enters malformed GA ID, checks validation
- **Invalid social URLs**: Enters non-URL text in social fields, verifies validateUrl() error
- **Partial social URLs**: Enters "facebook.com" without https://, checks validation
- **Mixed valid/invalid URLs**: Enters valid Facebook URL and invalid Twitter URL, checks handling
- **Tab switching with unsaved changes**: Edits field, switches tabs, verifies showUnsavedBanner()
- **Hidden active_tab field**: Verifies hidden input correctly sets active tab on submit
- **Character counter accuracy**: Tests counter at 0, 50%, 90%, 100%, 110% capacity
- **Empty optional fields**: Leaves social links empty, verifies optional handling

## Selectors & Elements
- Form: `form[action="/admin/settings"][method="POST"]`
- Hidden active tab: `input[name="active_tab"][type="hidden"]`
- Tab buttons: `.tab-button[data-tab="general"]`, `[data-tab="seo"]`, `[data-tab="social"]`
- Active tab indicator: `.tab-button.active`
- Tab content: `#general-tab`, `#seo-tab`, `#social-tab`
- Site name: `input[name="site_name"]`
- Site tagline: `input[name="site_tagline"]`
- Contact email: `input[name="contact_email"][type="email"]`
- Contact phone: `input[name="contact_phone"][type="tel"]`
- Address: `textarea[name="address"]`
- Business hours: `textarea[name="business_hours"]`
- Meta keywords: `input[name="meta_keywords"][maxlength="100"]`
- Meta description: `textarea[name="meta_description"][maxlength="250"]`
- Character counter (keywords): `#keywords-counter` or `.char-counter`
- Character counter (description): `#description-counter` or `.char-counter`
- Google Analytics ID: `input[name="google_analytics_id"]`
- Social Facebook: `input[name="social_facebook"][type="url"]`
- Social Twitter: `input[name="social_twitter"][type="url"]`
- Social LinkedIn: `input[name="social_linkedin"][type="url"]`
- Social Instagram: `input[name="social_instagram"][type="url"]`
- Social YouTube: `input[name="social_youtube"][type="url"]`
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` (when ?saved=1)
- Unsaved changes banner: `#unsaved-banner` or `.unsaved-warning`
- URL validation error: `.url-error` or validation message element

## HTMX Interactions
- None - standard form POST with full page redirect
- JavaScript handles tab switching client-side

## Dependencies
- Database settings table with global configuration
- JavaScript file with switchTab(), validateUrl(), showUnsavedBanner() functions
- Character counter JavaScript for meta fields
- Template: templates/admin/pages/global-settings.html
- Handler: internal/handlers/settings.go (GetGlobalSettings, PostGlobalSettings)
- URL validation for social media fields
