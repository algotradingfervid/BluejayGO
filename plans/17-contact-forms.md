# Phase 17 - Contact & Form Submissions

## Current State
- Contact submissions list + detail view
- Office locations list + form
- Basic table display

## Goal
Make submissions easy to review. Add "Request for Quote" form visibility alongside contact forms.

## Contact Submissions

### List Page
- Tabs: All | Contact Form | Request for Quote (if you add RFQ form later)
- Filter bar: Status (New / Read / Replied), Date range, Search
- Table: Name, Email, Subject, Status Badge, Submitted Date, Actions
- Status badges: New (yellow, bold), Read (gray), Replied (green)
- Unread count shown in tab label: "Contact Form (3)"
- Bulk "Mark as Read" action

### Detail Page
- Clean card layout showing full submission:
  - Contact info card: Name, Email, Phone, Company
  - Message card: Subject + full message body
  - Metadata card: Submitted date, IP address, referrer page
- Action buttons:
  - "Mark as Read" / "Mark as Unread"
  - "Reply via Email" (opens mailto: link)
  - "Delete"
- Previous / Next navigation between submissions

### Tooltips
- Status filter: "Filter submissions by their review status."
- Reply button: "Opens your email client with the sender's address pre-filled."

## Office Locations

### List
- Cards layout: Address + Phone + Email + Map link
- Sort order for display on public site

### Form
- Office Name
  - Tooltip: "Name for this location (e.g., 'Houston Headquarters', 'Dallas Office')."
- Address fields (street, city, state, zip, country)
- Phone
- Email
- Map URL
  - Tooltip: "Google Maps link for this location. Visitors can click to get directions."
- Is Primary (toggle)
  - Tooltip: "Primary office is shown first and used as the default contact address."
- Sort Order

## New Feature: Request for Quote Submissions
If a public RFQ form exists or will be added:
- Same list/detail pattern as contact submissions
- Additional fields: Product interested in, Quantity, Timeline
- Separate tab in submissions list

## Database Changes
- Add `status` column to `contact_submissions` if not exists (new/read/replied)
- Add `type` column to distinguish contact vs RFQ submissions

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/contact_submissions_list.html` | Rewrite |
| `templates/admin/pages/contact_submission_detail.html` | Rewrite |
| `templates/admin/pages/office_locations_list.html` | Redesign as cards |
| `templates/admin/pages/office_locations_form.html` | Add tooltips |
| `internal/handlers/admin/contact.go` | Add filtering, status updates |
| `db/migrations/029_contact_status.up.sql` | Add status column |

## Dependencies
- Phase 01, 02
