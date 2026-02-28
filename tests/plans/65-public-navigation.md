# Test Plan: Public Navigation (Header and Footer)

## Summary
Verify header and footer navigation render correctly from settings with working links and configurable elements.

**IMPLEMENTATION NOTES**:
- Settings use PascalCase: `.Settings.ShowNavProducts`, `.Settings.ShowNavAbout`, etc. (not snake_case)
- Header CTA button (`header_cta_enabled`) exists in database schema but is NOT implemented in header template
- Header phone/email display (`header_show_phone`, `header_show_email`) exists in database but NOT rendered in header template
- Configurable navigation labels: `.Settings.NavLabelHome`, `.Settings.NavLabelAbout`, `.Settings.NavLabelProducts`, etc.

## Preconditions
- Server running on localhost:28090
- Database settings configured for navigation (ShowNavProducts, ShowNavSolutions, ShowNavBlog, ShowNavAbout, ShowNavContact, ShowNavPartners)
- Navigation labels: NavLabelHome, NavLabelAbout, NavLabelProducts, NavLabelSolutions, NavLabelBlog, NavLabelContact, NavLabelPartners
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
- **Nav links display**: Verify navigation links based on .Settings.ShowNav* flags
- **Products link**: If .Settings.ShowNavProducts=true, verify link to /products with label from .Settings.NavLabelProducts
- **Solutions link**: If .Settings.ShowNavSolutions=true, verify link to /solutions with label from .Settings.NavLabelSolutions
- **Blog link**: If .Settings.ShowNavBlog=true, verify link to /blog with label from .Settings.NavLabelBlog
- **About link**: If .Settings.ShowNavAbout=true, verify link to /about with label from .Settings.NavLabelAbout
- **Contact link**: If .Settings.ShowNavContact=true, verify link to /contact with label from .Settings.NavLabelContact
- **Partners link**: If .Settings.ShowNavPartners=true, verify link to /partners with label from .Settings.NavLabelPartners
- **CTA button displays**: NOT IMPLEMENTED — header_cta_enabled exists in DB but NOT in template
- **CTA button action**: NOT IMPLEMENTED
- **Contact phone displays**: NOT IMPLEMENTED — header_show_phone exists in DB but NOT rendered
- **Contact email displays**: NOT IMPLEMENTED — header_show_email exists in DB but NOT rendered
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
- **All nav settings disabled**: Set all .Settings.ShowNav* to false, verify minimal header
- **CTA disabled**: NOT APPLICABLE — CTA not implemented in template
- **No contact info**: NOT APPLICABLE — phone/email not implemented in template
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
- Settings (PascalCase): ShowNavProducts, ShowNavSolutions, ShowNavBlog, ShowNavAbout, ShowNavContact, ShowNavPartners
- Navigation labels: NavLabelHome, NavLabelAbout, NavLabelProducts, NavLabelSolutions, NavLabelBlog, NavLabelContact, NavLabelPartners
- NOT IMPLEMENTED: header_cta_enabled, header_show_phone, header_show_email (exist in DB but not in template)
- Footer configuration settings
- Brutalist design system applied
- JetBrains Mono font
