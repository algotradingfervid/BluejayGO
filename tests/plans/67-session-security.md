# Test Plan: Session Security

## Summary
Verify session cookie security attributes and session integrity for authenticated admin users.

## Preconditions
- Server running on localhost:28090
- Admin user: admin@bluejaylabs.com / password
- Session cookie name: bluejay_session

## User Journey Steps
1. Navigate to admin login page
2. Log in with admin credentials
3. Inspect session cookie attributes
4. Verify session persists across requests
5. Test tampered cookie rejection
6. Test expired session handling

## Test Cases

### Happy Path - Session Creation
- **Login succeeds**: POST login with admin@bluejaylabs.com / password, verify redirect to admin dashboard
- **Session cookie set**: Verify bluejay_session cookie is set
- **Cookie attributes - HttpOnly**: Verify HttpOnly=true (cookie not accessible via JavaScript)
- **Cookie attributes - SameSite**: Verify SameSite=Lax
- **Cookie attributes - Path**: Verify Path=/
- **Cookie attributes - MaxAge**: Verify MaxAge=604800 (7 days)
- **Cookie attributes - Secure**: Verify Secure=false (development environment)
- **Session persistence**: Make subsequent requests, verify session cookie sent and recognized
- **Session data**: Verify session contains UserID, Email, DisplayName, Role fields

### Session Validation
- **Valid session access**: With valid session cookie, access admin pages, verify successful
- **Session fields populated**: Verify UserID, Email, DisplayName, Role are correctly stored
- **Session across pages**: Navigate between admin pages, verify session persists

### Error Cases - Invalid Sessions
- **Tampered cookie rejected**: Modify session cookie value, make request, verify rejection
- **New session after tamper**: Verify new session created after tampered cookie rejected
- **Expired session redirects**: Set system time forward 8 days (or use expired cookie), verify redirect to login
- **Expired session message**: Verify appropriate message shown for expired session
- **Missing cookie redirects**: Delete session cookie, access admin page, verify redirect to login
- **Invalid cookie format**: Send malformed cookie, verify rejection and new session creation

### Security Verification
- **HttpOnly prevents JS access**: Attempt to access cookie via document.cookie, verify cookie not visible
- **SameSite Lax protection**: Verify cross-site POST requests don't include cookie (if testable)
- **Secure flag in production**: Verify Secure=false for localhost, but document expectation of Secure=true in production

## Selectors & Elements
- Login page: form with email and password fields
- Session cookie: bluejay_session in browser cookies
- Admin pages: any /admin/* route for testing session persistence

## HTMX Interactions
- None specific to session management

## Dependencies
- Session middleware
- Session cookie configuration: HttpOnly=true, SameSite=Lax, Path=/, MaxAge=604800, Secure=false (dev)
- Session storage mechanism
- Session validation logic
- Login handler
- Admin authentication middleware
- Session fields: UserID, Email, DisplayName, Role
