# Test Plan: Homepage Hero Sections CRUD

## Summary
Tests the complete CRUD operations for homepage hero sections including headlines, CTAs, background images, active status, and display ordering.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with 1 hero section
- Valid image URLs for background images

## User Journey Steps
1. Navigate to /admin/homepage/heroes
2. View list of existing hero sections (1 seeded)
3. Click "New Hero" button
4. Fill in headline (required), subheadline, badge_text
5. Set primary CTA text and URL
6. Optionally set secondary CTA text and URL
7. Add background_image URL
8. Check/uncheck is_active checkbox
9. Set display_order
10. Submit to create hero
11. Edit existing hero from list
12. Update fields and save (redirect to edit page with ?saved=1)
13. Delete a hero
14. Verify proper ordering and active status filtering

## Test Cases

### Happy Path
- **List heroes**: Verifies GET /admin/homepage/heroes shows seeded hero
- **Create new hero**: Adds hero with all required fields, verifies creation
- **Edit existing hero**: Updates headline and subheadline, verifies save with ?saved=1 redirect
- **Update primary CTA**: Changes primary_cta_text and primary_cta_url, verifies update
- **Add secondary CTA**: Adds secondary CTA where none existed, verifies save
- **Remove secondary CTA**: Clears secondary CTA fields, verifies removal
- **Toggle is_active**: Checks/unchecks is_active, verifies active status change
- **Update background image**: Changes background_image URL, verifies update
- **Reorder heroes**: Changes display_order, verifies list reordering
- **Badge text**: Adds/updates badge_text, verifies display

### Edge Cases / Error States
- **Empty headline**: Tests required field validation on headline
- **Empty primary CTA text**: Tests validation when primary CTA text missing but URL present
- **Empty primary CTA URL**: Tests validation when primary CTA URL missing but text present
- **Invalid image URL**: Enters malformed background_image URL, checks validation
- **Secondary CTA partial**: Enters only text or only URL for secondary CTA, checks validation
- **Multiple active heroes**: Creates multiple heroes with is_active=true, verifies handling
- **Very long headline**: Tests character limits on headline field
- **Very long subheadline**: Tests character limits on subheadline
- **Delete hero**: Deletes hero, verifies removal from list
- **Delete confirmation**: Verifies confirmation before deletion
- **Edit redirect**: Confirms redirect to /admin/homepage/heroes/:id/edit?saved=1 after update

## Selectors & Elements
- Heroes list: `#heroes-list` or `.heroes-table`
- New hero button: `a[href="/admin/homepage/heroes/new"]` or `button#new-hero`
- Hero row: `.hero-row[data-id]` or `tr[data-hero-id]`
- Edit link: `a[href="/admin/homepage/heroes/{id}/edit"]`
- Delete button: `button[type="submit"]` or link in delete form
- Headline input: `input[name="headline"]`
- Subheadline input: `input[name="subheadline"]`
- Badge text input: `input[name="badge_text"]`
- Primary CTA text: `input[name="primary_cta_text"]`
- Primary CTA URL: `input[name="primary_cta_url"]`
- Secondary CTA text: `input[name="secondary_cta_text"]`
- Secondary CTA URL: `input[name="secondary_cta_url"]`
- Background image URL: `input[name="background_image"]`
- Is active checkbox: `input[name="is_active"][type="checkbox"]`
- Display order input: `input[name="display_order"][type="number"]`
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` (when ?saved=1 present)

## HTMX Interactions
- None - standard form submissions with redirects
- Delete may use standard form POST or DELETE method

## Dependencies
- Database seeded with 1 hero section
- Template: templates/admin/pages/homepage-heroes-list.html, homepage-heroes-form.html
- Handler: internal/handlers/homepage.go (ListHeroes, NewHero, CreateHero, EditHero, UpdateHero, DeleteHero)
