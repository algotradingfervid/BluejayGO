# Test Plan: Admin Login & Logout

## Summary
Verify admin authentication flow including login, logout, session management, and access control redirects.

## Preconditions
- Server running on localhost:28090
- Admin user exists: admin@bluejaylabs.com / password
- Database initialized with user credentials

## User Journey Steps
1. Navigate to http://localhost:28090/admin/login
2. Enter email "admin@bluejaylabs.com" in `input[name="email"]`
3. Enter password "password" in `input[name="password"]`
4. Click button with text "Sign In"
5. Verify redirect to http://localhost:28090/admin/dashboard
6. Verify session cookie "bluejay_session" is set
7. Click logout button (POST /admin/logout)
8. Verify redirect to http://localhost:28090/admin/login
9. Verify session cookie is cleared

## Test Cases

### Happy Path
- **Valid login credentials**: User can log in with correct email and password, receives session cookie, and is redirected to dashboard
- **Successful logout**: Authenticated user can log out, session is cleared, and user is redirected to login page
- **Auth redirect for protected routes**: Unauthenticated user accessing /admin/dashboard is redirected to /admin/login
- **Already authenticated redirect**: User with valid session accessing /admin/login is redirected to /admin/dashboard

### Edge Cases / Error States
- **Invalid credentials error**: Login with wrong password shows error parameter ?error=invalid_credentials
- **Missing email field**: Submitting form without email shows ?error=missing_fields
- **Missing password field**: Submitting form without password shows ?error=missing_fields
- **Missing both fields**: Submitting empty form shows ?error=missing_fields
- **Session error**: Invalid or expired session shows ?error=session_error
- **Non-existent email**: Login with email not in database shows ?error=invalid_credentials

## Selectors & Elements
- Form: `action="/admin/login" method="POST"`
- Input email: `name="email" type="email"`
- Input password: `name="password" type="password"`
- Button: text "Sign In"
- Error message container: displays conditionally based on URL error parameter
- Logout form: `action="/admin/logout" method="POST"`

## HTMX Interactions
- No HTMX on login/logout forms (standard form POST)
- Full page redirects after authentication actions

## Dependencies
- 02-admin-dashboard.md (verifies dashboard landing page after login)
- All other admin test plans (require authenticated session)
