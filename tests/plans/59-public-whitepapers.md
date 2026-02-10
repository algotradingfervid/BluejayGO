# Test Plan: Public Whitepapers

## Summary
Verify whitepapers listing, detail pages, and lead capture download flow with form validation.

## Preconditions
- Server running on localhost:28090
- Database seeded with 12 whitepapers
- No authentication required for viewing
- Lead capture form required for downloads

## User Journey Steps
1. Navigate to GET /whitepapers for listing
2. View whitepaper cards organized by topic
3. Click whitepaper to view detail
4. Navigate to GET /whitepapers/:slug
5. Read description and learning points
6. Fill download form with lead information
7. Submit POST /whitepapers/:slug/download
8. View success page with download link
9. Download whitepaper file

## Test Cases

### Happy Path - Listing
- **Whitepapers listing loads**: Verify GET /whitepapers returns 200 status
- **Cards by topic**: Verify whitepapers organized by topic with gradient covers
- **12 seeded whitepapers**: Verify all 12 whitepapers display
- **Whitepaper links**: Verify cards link to /whitepapers/:slug
- **Card navigation**: Click card, verify navigation to detail page

### Happy Path - Detail & Download
- **Whitepaper detail loads**: Verify GET /whitepapers/:slug returns 200 status
- **Description displays**: Verify whitepaper description text
- **Learning points display**: Verify learning points/benefits listed
- **Download form displays**: Verify form with name, email, company, designation, marketing_consent fields
- **Form validation - required fields**: Verify name and email are required
- **Form submission**: Fill all required fields, submit, verify POST /whitepapers/:slug/download
- **Success page**: Verify redirect to /whitepapers/:slug/success
- **Download link available**: Verify success page (whitepaper_success.html) contains download link
- **Download file**: Click download link, verify file downloads

### Edge Cases / Error States
- **Whitepaper not found**: Navigate to invalid slug, verify 404 or error page
- **Form validation - missing name**: Submit without name, verify validation error
- **Form validation - missing email**: Submit without email, verify validation error
- **Form validation - invalid email**: Submit with invalid email format, verify error
- **Marketing consent checkbox**: Verify checkbox is optional (can submit unchecked)
- **Company field optional**: Verify form submits without company
- **Designation field optional**: Verify form submits without designation
- **Form resubmission**: Submit form twice for same whitepaper, verify handling
- **Success page direct access**: Navigate directly to /whitepapers/:slug/success, verify behavior

## Selectors & Elements
- Listing page:
  - Cards container organized by topic
  - Whitepaper cards with gradient covers, title, topic
  - Links to `/whitepapers/*`
- Detail page:
  - Whitepaper description text
  - Learning points list
  - Download form with fields:
    - `input[name="name"]` (required)
    - `input[name="email"]` (required)
    - `input[name="company"]` (optional)
    - `input[name="designation"]` (optional)
    - `input[type="checkbox"][name="marketing_consent"]` (optional)
  - Submit button
- Success page:
  - Success message
  - Download link to whitepaper file

## HTMX Interactions
- None (traditional form submission with POST and redirect)

## Dependencies
- Template: whitepaper_success.html
- Form handler: POST /whitepapers/:slug/download
- Lead capture data storage
- File download mechanism
- Seeded database with 12 whitepapers
- Brutalist design system applied
- JetBrains Mono font
