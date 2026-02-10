# Test Plan: Whitepaper Downloads Analytics

## Summary
Testing the read-only analytics view for whitepaper download lead capture data with filtering and pagination.

## Preconditions
- Server running on localhost:28090
- Logged in as admin@bluejaylabs.com / password
- Database seeded with 29 whitepaper download records

## User Journey Steps
1. Navigate to http://localhost:28090/admin/whitepapers/downloads
2. Verify download records list displays with whitepaper filter, date_from, date_to, page params
3. Verify each record shows: name, email, company, designation, marketing_consent, download date
4. Apply whitepaper filter to show downloads for specific whitepaper
5. Apply date range filters: date_from and date_to
6. Navigate through pagination
7. Verify this is read-only view with no edit/delete actions
8. Navigate to specific whitepaper downloads via /admin/whitepapers/:id/downloads

## Test Cases

### Happy Path
- **List all downloads**: Navigate to /admin/whitepapers/downloads, see all 29 seeded records
- **Whitepaper filter**: Select specific whitepaper from dropdown, list filtered to show only that whitepaper's downloads
- **Date range filter**: Set date_from "2024-01-01", date_to "2024-12-31", see downloads in that range
- **Combined filters**: Apply whitepaper + date range filters together, see filtered results
- **Pagination**: Navigate through pages if more than 15 records per page
- **Download details**: Each record shows name, email, company, designation, marketing_consent boolean, download timestamp
- **Whitepaper-specific view**: Navigate to /admin/whitepapers/:id/downloads, see downloads for that whitepaper only
- **Read-only view**: No edit or delete buttons present, only view capability

### Edge Cases / Error States
- **No filters applied**: All downloads shown, may be paginated if 29+ records
- **Filter no results**: Apply filters that match no downloads, see "No downloads found" message
- **Invalid date range**: date_from after date_to may show validation error or no results
- **Empty date filter**: Leaving date_from or date_to empty accepted, filters on populated field only
- **Marketing consent true/false**: marketing_consent column shows boolean value or checkbox icon
- **Long company name**: Company name with 100+ characters displayed, may wrap or truncate
- **Missing designation**: designation field null shows empty or "N/A"
- **Email validation**: Email format displayed as-is, no validation on view page

## Selectors & Elements
- List route: GET /admin/whitepapers/downloads
- Whitepaper-specific route: GET /admin/whitepapers/:id/downloads
- Query params: whitepaper (filter by whitepaper ID), date_from (date filter), date_to (date filter), page (pagination)
- Display columns: name, email, company, designation, marketing_consent, download_date
- Filter form: whitepaper (select dropdown), date_from (date input), date_to (date input)
- Submit filter button: text "Filter" or "Apply"
- Pagination controls: page numbers or next/prev links

## HTMX Interactions
- No HTMX interactions, this is a standard server-rendered paginated list view
- Filters applied via form submission GET request with query params

## Dependencies
- 27-whitepapers-crud.md (whitepapers must exist to have downloads)
- Whitepaper_downloads table seeded with 29 records including lead capture data
