# Test Plan: Public Partners Page

## Summary
Verify partners page displays partner directory with tier filtering and testimonials.

## Preconditions
- Server running on localhost:28090
- Database seeded with 11 partners and 6 testimonials
- No authentication required

## User Journey Steps
1. Navigate to GET /partners
2. View partner directory
3. Filter partners by tier
4. View partner cards with details
5. Click partner website links
6. Read partner testimonials section

## Test Cases

### Happy Path
- **Partners page loads**: Verify GET /partners returns 200 status
- **Partner directory displays**: Verify partner cards layout
- **11 seeded partners**: Verify all 11 partners display
- **Tier filtering**: Verify tier filter controls present
- **Partner card content**: Verify each card has logo, name, tier, description, website link
- **Filter by tier**: Click tier filter, verify only matching partners display
- **Website links**: Verify partner website links are clickable and properly formatted
- **Testimonials section**: Verify 6 partner testimonials display
- **Testimonial content**: Verify each testimonial has quote, author name, author info

### Edge Cases / Error States
- **Filter to empty tier**: Filter by tier with no partners, verify appropriate handling
- **Missing partner logo**: Verify fallback display when logo missing
- **Missing website**: Verify partner card handles missing website URL
- **Long descriptions**: Verify text truncation or wrapping for long partner descriptions
- **External link behavior**: Verify partner website links open appropriately (target="_blank" or same window)

## Selectors & Elements
- Partner directory:
  - Tier filter controls
  - Partner cards container
  - Partner cards with:
    - Partner logo image
    - Partner name
    - Tier badge/indicator
    - Description text
    - Website link
- Testimonials section:
  - Testimonials container
  - Testimonial cards with:
    - Quote text
    - Author name
    - Author company/role info

## HTMX Interactions
- Possible tier filtering via HTMX (or JavaScript/traditional navigation)

## Dependencies
- Template data: partners collection, testimonials
- Seeded database with 11 partners and 6 testimonials
- Partner tier classifications
- Brutalist design system applied
- JetBrains Mono font
