# Test Plan: Security Headers

## Summary
Verify SecurityHeaders middleware adds required security headers to all HTTP responses.

## Preconditions
- Server running on localhost:28090
- SecurityHeaders middleware configured and active
- No authentication required for public pages

## User Journey Steps
1. Make request to any public page
2. Inspect HTTP response headers
3. Verify all security headers present
4. Verify header values are correct

## Test Cases

### Happy Path - Public Pages
- **Homepage headers**: Request GET /, verify all security headers present
- **Products page headers**: Request GET /products, verify headers
- **Blog page headers**: Request GET /blog, verify headers
- **Contact page headers**: Request GET /contact, verify headers
- **Static asset headers**: Request static file, verify headers (if middleware applies)

### Happy Path - Admin Pages
- **Admin login headers**: Request GET /admin/login, verify headers
- **Admin dashboard headers**: Request GET /admin/dashboard (authenticated), verify headers

### Security Headers Validation
- **Content-Security-Policy present**: Verify CSP header exists
- **CSP directives**: Verify CSP contains appropriate directives (default-src, script-src, style-src, etc.)
- **X-Frame-Options present**: Verify header exists
- **X-Frame-Options value**: Verify value is DENY or SAMEORIGIN
- **X-Content-Type-Options present**: Verify header exists
- **X-Content-Type-Options value**: Verify value is nosniff
- **Referrer-Policy present**: Verify header exists
- **Referrer-Policy value**: Verify appropriate policy (e.g., no-referrer, strict-origin-when-cross-origin)
- **X-XSS-Protection present**: Verify header exists
- **X-XSS-Protection value**: Verify value (e.g., 1; mode=block)

### Header Value Details
- **CSP default-src**: Verify default-src directive (e.g., 'self')
- **CSP script-src**: Verify script-src allows necessary sources (e.g., 'self', CDN for HTMX/Tailwind)
- **CSP style-src**: Verify style-src allows necessary sources (e.g., 'self', CDN for Tailwind, 'unsafe-inline' if needed)
- **CSP img-src**: Verify img-src allows necessary sources (e.g., 'self', data:)
- **Frame options**: Verify DENY or SAMEORIGIN prevents clickjacking

### Edge Cases
- **Multiple requests**: Make multiple requests, verify headers consistent
- **Different content types**: Request HTML, JSON, verify headers on all response types
- **Error responses**: Trigger 404 or 500, verify security headers still present
- **AJAX requests**: Make HTMX or fetch request, verify headers on partial responses

## Selectors & Elements
- HTTP Response Headers:
  - Content-Security-Policy
  - X-Frame-Options
  - X-Content-Type-Options
  - Referrer-Policy
  - X-XSS-Protection

## HTMX Interactions
- Verify security headers present on HTMX partial responses (e.g., /products/search)

## Dependencies
- SecurityHeaders middleware
- Middleware applied to all routes
- CSP configuration matching application needs (HTMX, Tailwind CDN, etc.)
- Header configuration constants or settings
