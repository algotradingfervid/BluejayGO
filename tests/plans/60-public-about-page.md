# Test Plan: Public About Page

## Summary
Verify about page displays company overview, mission/vision/values, core values grid, timeline, and certifications.

## Preconditions
- Server running on localhost:28090
- Database seeded with company information, values, milestones, and certifications
- No authentication required

## User Journey Steps
1. Navigate to GET /about
2. View company overview section
3. Read mission, vision, and values sections
4. View core values grid with icons
5. Review milestones timeline
6. View certifications section

## Test Cases

### Happy Path
- **About page loads**: Verify GET /about returns 200 status
- **Company overview section**: Verify headline, tagline, description texts, and images display
- **Mission section**: Verify mission statement with icon
- **Vision section**: Verify vision statement with icon
- **Values section**: Verify company values with icon
- **Core values grid**: Verify grid layout of core values with icons and descriptions
- **Milestones timeline**: Verify chronological timeline of company milestones
- **Certifications section**: Verify company certifications displayed

### Edge Cases / Error States
- **Missing images**: Verify graceful handling if overview images missing
- **Empty milestones**: Verify timeline section handles empty milestones
- **Empty certifications**: Verify certifications section handles empty list
- **Long text content**: Verify proper text wrapping and layout with extensive content

## Selectors & Elements
- Company overview section:
  - Headline text
  - Tagline text
  - Description paragraphs
  - Company images
- Mission section: mission text with icon
- Vision section: vision text with icon
- Values section: values text with icon
- Core values grid:
  - Grid container
  - Value items with icons and descriptions
- Milestones timeline:
  - Timeline container
  - Milestone entries with dates and descriptions
- Certifications section:
  - Certifications container
  - Certification items/badges

## HTMX Interactions
- None (static content display)

## Dependencies
- Template data: company info, mission/vision/values, core values, milestones, certifications
- Seeded database with complete about page content
- Brutalist design system applied
- JetBrains Mono font
- Icons for mission/vision/values and core values
