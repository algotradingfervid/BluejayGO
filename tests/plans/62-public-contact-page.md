# Test Plan: Public Contact Page

## Summary
Verify contact page displays form with validation and rate limiting, plus office locations.

## Preconditions
- Server running on localhost:28090
- Database seeded with 6 office locations
- No authentication required
- Rate limiting: 5 submissions per hour per IP (contactLimiter middleware)

## User Journey Steps
1. Navigate to GET /contact
2. View contact form and office locations
3. Fill contact form with required fields
4. Submit form via POST /contact/submit
5. View success/error feedback
6. Test rate limiting by multiple submissions

## Test Cases

### Happy Path
- **Contact page loads**: Verify GET /contact returns 200 status
- **Contact form displays**: Verify form with all fields present
- **Office locations display**: Verify 6 office locations with addresses/details
- **Form field validation - name required**: Verify name field is required
- **Form field validation - email required**: Verify email field is required
- **Form field validation - message required**: Verify message textarea is required
- **Form field validation - optional fields**: Verify phone and company are optional
- **Inquiry type select**: Verify inquiry_type dropdown with options
- **Form submission success**: Fill all required fields, submit, verify success feedback
- **Success message displays**: Verify success feedback after POST /contact/submit

### Edge Cases / Error States
- **Form validation - missing name**: Submit without name, verify validation error
- **Form validation - missing email**: Submit without email, verify validation error
- **Form validation - invalid email**: Submit with invalid email format, verify error
- **Form validation - missing message**: Submit without message, verify validation error
- **Rate limiting - 5 submissions succeed**: Submit form 5 times within an hour, verify all succeed
- **Rate limiting - 6th submission blocked**: Submit 6th time, verify rate limit error response
- **Rate limiting error message**: Verify appropriate error message for rate-limited request
- **Rate limit reset**: Wait for rate limit window to expire, verify submissions work again
- **Long message text**: Submit with very long message, verify handling
- **Special characters in fields**: Submit with special characters, verify proper handling
- **XSS attempt in fields**: Submit with script tags, verify sanitization

## Selectors & Elements
- Contact form:
  - `form[action="/contact/submit"][method="POST"]`
  - `input[name="name"][required]`
  - `input[name="email"][required]`
  - `input[name="phone"]` (optional)
  - `input[name="company"]` (optional)
  - `select[name="inquiry_type"]` with options
  - `textarea[name="message"][required]`
  - Submit button
- Success/error feedback: feedback message container
- Office locations:
  - Locations container
  - 6 office location cards with address details

## HTMX Interactions
- Possible HTMX form submission for dynamic feedback (or traditional POST with redirect)

## Dependencies
- Form handler: POST /contact/submit
- contactLimiter middleware: 5 requests per hour per IP
- Seeded database with 6 office locations
- Form validation logic
- Rate limiting implementation
- Brutalist design system applied
- JetBrains Mono font
