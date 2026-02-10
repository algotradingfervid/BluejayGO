# Phase 22 - Login & Auth Pages

## Current State
- Basic login form: email + password + submit button
- Minimal styling
- No forgot password flow

## Goal
Polished login page matching the brutalist design system. Clean, centered, intimidation-free.

## Login Page

### Layout
- Full-screen centered card (no sidebar, no header)
- Max-width 400px
- Site logo at top (from settings, or BlueJay Labs default)
- Login card below logo

### Card Content
- "Sign In" heading (font-display)
- Email input
  - Placeholder: "admin@example.com"
- Password input
  - Placeholder: "Enter your password"
  - Toggle visibility icon (eye)
- "Remember me" checkbox
- "Sign In" button (full-width, primary style)
- Error message area (red banner, hidden by default)

### Design
- White card with manual-border and manual-shadow
- Subtle background pattern or solid gray
- All brutalist styling (sharp edges, black borders)
- Focus states on inputs: thicker border or primary color border

### Error States
- Invalid credentials: red banner "Invalid email or password"
- Account locked (future): "Account locked. Contact administrator."

## Forgot Password Page (Stretch Goal)
Since this is a solo admin setup, this could be:
- Simple page with email input
- Since it's SQLite/local, password reset could be a CLI command instead
- For now: just add a help text below login: "Forgot your password? Reset it via the command line: `./server reset-password`"

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/login.html` | Rewrite |
| `internal/handlers/admin/auth.go` | Minor updates if needed |

## Dependencies
- Phase 01 (for design system CSS)
- No other dependencies (can be done anytime after Phase 01)
