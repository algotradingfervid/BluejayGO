# Test Plan: Public Navigation (Header and Footer)

## Summary
Verify header and footer navigation render correctly from settings with working links and configurable elements.

## Preconditions
- Server running on localhost:28090
- Database settings configured for navigation (show_nav_* settings, header_cta_enabled, show_phone, show_email)
- No authentication required

## User Journey Steps
1. Load any public page
2. View header navigation
3. Verify logo, nav links, CTA button, contact info
4. Click navigation links
5. View footer navigation
6. Verify footer columns, social links, legal links
7. Click footer links

## Test Cases

### Happy Path - Header
- **Header renders on all pages**: Verify header present on /, /products, /solutions, /blog, /about, /contact, etc.
- **Logo displays**: Verify company logo image/text in header
- **Logo links to home**: Click logo, verify navigation to /
- **Nav links display**: Verify navigation links based on show_nav_* settings
- **Products link**: If show_nav_products=true, verify link to /products
- **Solutions link**: If show_nav_solutions=true, verify link to /solutions
- **Blog link**: If show_nav_blog=true, verify link to /blog
- **About link**: If show_nav_about=true, verify link to /about
- **Contact link**: If show_nav_contact=true, verify link to /contact
- **Partners link**: If show_nav_partners=true, verify link to /partners
- **CTA button displays**: If header_cta_enabled=true, verify CTA button present
- **CTA button action**: Click CTA button, verify configured action (link or modal)
- **Contact phone displays**: If show_phone=true, verify phone number in header
- **Contact email displays**: If show_email=true, verify email address in header
- **Navigation link click**: Click each nav link, verify navigation to correct page

### Happy Path - Footer
- **Footer renders on all pages**: Verify footer present on all public pages
- **Footer columns display**: Verify configurable footer columns with links
- **Footer links work**: Click footer links, verify navigation
- **Social links display**: Verify social media icons/links
- **Social links open**: Click social link, verify opens to correct social profile
- **Legal links display**: Verify privacy policy, terms of service links
- **Copyright notice**: Verify copyright text with current year
- **Footer layout**: Verify multi-column footer layout matches settings

### Edge Cases / Error States
- **All nav settings disabled**: Set all show_nav_* to false, verify minimal header
- **CTA disabled**: Set header_cta_enabled=false, verify CTA button hidden
- **No contact info**: Set show_phone=false and show_email=false, verify contact info hidden
- **Long nav menu**: Enable all nav links, verify layout handles many items
- **Missing social links**: Verify footer handles missing social media URLs
- **Empty footer columns**: Verify footer handles columns with no links

## Selectors & Elements
- Header:
  - Header container
  - Logo image/text with link to `/`
  - Navigation menu container
  - Nav links: `/products`, `/solutions`, `/blog`, `/about`, `/contact`, `/partners` (conditional)
  - CTA button (conditional)
  - Phone number (conditional)
  - Email address (conditional)
- Footer:
  - Footer container
  - Footer columns
  - Footer links
  - Social media icons/links
  - Legal links (privacy policy, terms of service)
  - Copyright text

## HTMX Interactions
- None (static navigation rendering)

## Dependencies
- Template: header partial, footer partial
- Settings: show_nav_products, show_nav_solutions, show_nav_blog, show_nav_about, show_nav_contact, show_nav_partners, header_cta_enabled, show_phone, show_email
- Footer configuration settings
- Brutalist design system applied
- JetBrains Mono font
