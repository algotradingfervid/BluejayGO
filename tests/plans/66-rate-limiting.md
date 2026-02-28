# Test Plan: Rate Limiting

## Summary
Verify contact form rate limiting restricts to 5 submissions per hour per IP address.

**IMPLEMENTATION NOTES**:
- Rate limit error returns JSON (`{"error": "Too many requests. Please try again later."}`) with Content-Type `application/json`, but contact form expects HTML fragments via HTMX — potential UX issue
- Rate limiting middleware runs BEFORE handler validation, so invalid form submissions also count toward the limit
- GET requests to `/contact` do NOT count toward rate limit (limiter only applied to POST route)

## Preconditions
- Server running on localhost:28090
- contactLimiter middleware configured: 5 requests per hour per IP
- No authentication required

## User Journey Steps
1. Navigate to GET /contact
2. Submit contact form 5 times successfully
3. Attempt 6th submission within same hour
4. Verify rate limit error response
5. Wait for rate limit window to reset (or test reset mechanism)
6. Verify submissions work again after reset

## Test Cases

### Happy Path - Within Limit
- **First submission succeeds**: Fill form, submit POST /contact/submit, verify success (200 status)
- **Second submission succeeds**: Submit again, verify success
- **Third submission succeeds**: Submit again, verify success
- **Fourth submission succeeds**: Submit again, verify success
- **Fifth submission succeeds**: Submit again, verify success (5th request within limit)

### Error Cases - Rate Limited
- **Sixth submission blocked**: Submit 6th time within same hour, verify rate limit error
- **Rate limit status code**: Verify response status is 429 (Too Many Requests)
- **Rate limit response format**: Response is JSON: `{"error": "Too many requests. Please try again later."}` with Content-Type `application/json`
- **HTMX compatibility issue**: HTMX form expects HTML fragment, but gets JSON — may not display error properly
- **User feedback**: Verify error message handling (JSON response may not integrate with HTMX UI)
- **Seventh submission also blocked**: Attempt 7th submission, verify still blocked

### Rate Limit Reset
- **Rate limit window expiry**: Wait for 1 hour (or use clock manipulation if available), verify submissions work again
- **Counter reset**: After reset, verify 5 new submissions allowed
- **Per-IP isolation**: From different IP (if testable), verify separate rate limit counter

### Edge Cases
- **Rapid submissions**: Submit 5 times rapidly in succession, verify all 5 succeed before rate limit kicks in
- **Partial form data**: Verify incomplete form submissions still count toward rate limit
- **Failed validations**: Middleware runs BEFORE handler validation, so invalid submissions DO count toward limit
- **GET requests not limited**: GET /contact does NOT count toward rate limit (only POST /contact/submit is limited)
- **Rate limit persistence**: Restart server (if applicable), verify rate limit counters persist or reset as designed

## Selectors & Elements
- Contact form: `form[action="/contact/submit"]`
- Submit button
- Success message container
- Error message container (for rate limit error)

## HTMX Interactions
- None specific to rate limiting (applies to POST /contact/submit regardless of submission method)

## Dependencies
- contactLimiter middleware
- Rate limiting implementation (5 requests per hour per IP)
- IP address detection mechanism
- Rate limit counter storage (in-memory or persistent)
- Error response formatting
- Contact form handler: POST /contact/submit
