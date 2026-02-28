# Test Plan: Solution CTAs HTMX Management

## Summary
Testing HTMX-driven addition and deletion of call-to-action sections within solution edit page without full page reloads.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with solutions
- On solution edit page at /admin/solutions/:id/edit

## User Journey Steps
1. Navigate to http://localhost:28090/admin/solutions/:id/edit
2. Locate #ctas-list container
3. Fill add CTA form: heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name
4. Click add button with hx-post="/admin/solutions/:id/ctas"
5. Verify hx-target="#ctas-list" hx-swap="innerHTML" adds new CTA to list
6. Verify new CTA appears in list
7. Click delete button on existing CTA with hx-delete="/admin/solutions/:id/ctas/:ctaId"
8. Verify CTA removed (hx-target="closest div" hx-swap="outerHTML") without page reload

## Test Cases

### Happy Path
- **Add CTA with primary button**: Fill heading "Get Started Today" (required), subheading "Transform...", primary_button_text "Contact Sales", primary_button_url "/contact", section_name "hero-cta" (required), submit, CTA added (innerHTML swap adds partial HTML)
- **Add CTA with secondary button**: Fill both primary and secondary button fields, verify both buttons appear
- **Add CTA with phone**: Fill phone_number "1-800-123-4567", verify phone displays
- **Section name**: Fill section_name "hero-cta" (required) for internal reference
- **Multiple CTAs**: Add 2 CTAs with different section_name values, verify all appear
- **Delete CTA**: Click delete button on CTA, hx-delete removes individual div (returns c.NoContent(http.StatusOK) with no HTML)
- **Add returns HTML**: Add handler returns partial HTML
- **Delete returns no content**: DeleteCTA returns c.NoContent(http.StatusOK) with no HTML
- **No page reload**: All operations via HTMX, no full page refresh

### Edge Cases / Error States
- **Required fields**: heading and section_name are required (have required attribute in template)
- **Missing required heading**: Empty heading field triggers validation error
- **Missing required section_name**: Empty section_name field triggers validation error
- **Only secondary button**: Fill only secondary button fields without primary, accepted
- **Invalid URL**: primary_button_url with invalid format may show validation warning
- **Long subheading**: Subheading with 200+ characters accepted, may wrap in display
- **Phone number format**: Phone without specific format accepted, formatted on frontend
- **Delete with confirmation**: hx-confirm attribute prompts user before deletion
- **HTMX error handling**: Server error on add/delete shows error message in section

## Selectors & Elements
- Section container: id="ctas-list"
- Add form action: hx-post="/admin/solutions/:id/ctas" hx-target="#ctas-list" hx-swap="innerHTML"
- Input names: heading (required), subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name (required)
- Delete button: hx-delete="/admin/solutions/:id/ctas/:ctaId" hx-target="closest div" hx-swap="outerHTML"
- Add button: text "Add CTA"

## HTMX Interactions
- **Add CTA**: hx-post="/admin/solutions/:id/ctas" hx-target="#ctas-list" hx-swap="innerHTML" (returns partial HTML)
- **Delete CTA**: hx-delete="/admin/solutions/:id/ctas/:ctaId" hx-target="closest div" hx-swap="outerHTML" (returns c.NoContent(http.StatusOK) with no HTML content)
- Add operation inserts new CTA HTML into list, delete operation removes individual item div

## Dependencies
- 19-solutions-crud.md (parent solution edit page)
