# Test Plan: Rate Limiting

## Summary
Verify contact form rate limiting restricts to 5 submissions per hour per IP address.

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
- **Rate limit status code**: Verify response status is 429 (Too Many Requests) or similar error code
- **Rate limit error message**: Verify error message indicates rate limit exceeded
- **User feedback**: Verify user-friendly error message displayed on page
- **Seventh submission also blocked**: Attempt 7th submission, verify still blocked

### Rate Limit Reset
- **Rate limit window expiry**: Wait for 1 hour (or use clock manipulation if available), verify submissions work again
- **Counter reset**: After reset, verify 5 new submissions allowed
- **Per-IP isolation**: From different IP (if testable), verify separate rate limit counter

### Edge Cases
- **Rapid submissions**: Submit 5 times rapidly in succession, verify all 5 succeed before rate limit kicks in
- **Partial form data**: Verify incomplete form submissions still count toward rate limit
- **Failed validations**: Verify failed validation attempts count (or don't count) toward limit appropriately
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
