# Test Plan: Activity Log Viewer

## Summary
Tests the read-only activity log with filtering, search, and pagination showing all system actions (create/update/delete) across resources.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Activity log populated with various entries (creates, updates, deletes)
- Pagination: 50 entries per page
- No create/edit/delete functionality (read-only)

## User Journey Steps
1. Navigate to /admin/activity
2. View paginated list of activity entries (50 per page)
3. Use action filter dropdown to filter by create/update/delete
4. Use search to find entries by user_email, resource_title, or description
5. Navigate pagination to view older entries
6. View entry details: timestamp, user email, action badge, resource type, resource title, description
7. Click "Clear Filters" when HasFilters=true
8. Verify no create/edit/delete buttons exist

## Test Cases

### Happy Path - Viewing Activity
- **Load activity log**: Verifies GET /admin/activity shows 50 entries
- **View all fields**: Confirms each entry shows timestamp, user_email, action, resource_type, resource_title, description
- **Action badges**: Verifies create/update/delete actions have distinct badges/colors
- **Chronological order**: Verifies entries display newest first by default
- **Pagination**: Clicks page 2, verifies next 50 entries load
- **Page navigation**: Tests first page, middle page, last page

### Happy Path - Filtering by Action
- **Filter create actions**: Selects action=create, verifies only create entries shown
- **Filter update actions**: Selects action=update, verifies only update entries shown
- **Filter delete actions**: Selects action=delete, verifies only delete entries shown
- **Clear action filter**: Selects "All actions" or clears filter, verifies all entries shown

### Happy Path - Search
- **Search by user email**: Enters "admin@bluejaylabs.com", verifies filtered results
- **Search by resource title**: Searches for resource name, verifies matching entries
- **Search by description**: Searches description text, verifies matches
- **Search no results**: Searches for non-existent term, verifies "no results" message
- **Search with action filter**: Combines search + action filter, verifies both applied

### Happy Path - Clear Filters
- **HasFilters true**: Applies filter, verifies "Clear Filters" button appears
- **Clear filters**: Clicks clear, verifies redirect to /admin/activity without params
- **HasFilters false**: With no filters, verifies "Clear Filters" button hidden

### Edge Cases / Error States
- **Empty log**: Tests display when no activity entries exist, verifies empty state
- **Search special characters**: Searches with quotes, symbols, verifies handling
- **Search very long query**: Enters 200+ char search, verifies handling
- **Invalid page number**: Navigates to page=999 when only 5 pages exist, verifies handling
- **Page=0 or negative**: Tests invalid page param, verifies redirect or error
- **No create button**: Verifies absence of "New Entry" or similar button
- **No edit links**: Verifies entries are not editable
- **No delete buttons**: Verifies no delete option exists
- **Read-only enforcement**: Attempts POST/PUT/DELETE to activity endpoints, verifies rejection
- **Action filter edge**: Tests if action param accepts invalid values
- **Combined filters pagination**: Applies filters, navigates pages, verifies filter persistence
- **Very long resource title**: Tests display of 200+ char title, verifies truncation
- **Very long description**: Tests display of 500+ char description, verifies truncation
- **Missing user email**: Tests entry with null/empty user_email, verifies display
- **Missing resource title**: Tests entry with empty resource_title, verifies display
- **Timestamp formatting**: Verifies timestamps display in readable format with timezone
- **Pagination at end**: On last page, verifies "Next" button disabled or hidden
- **Pagination at start**: On first page, verifies "Previous" button disabled or hidden
- **50 entries per page**: Verifies exactly 50 entries on full pages
- **Partial last page**: On last page with <50 entries, verifies correct count

## Selectors & Elements
- Activity log table: `#activity-log` or `.activity-table`
- Activity entry row: `.activity-row` or `tr.activity-entry`
- Timestamp: `.activity-timestamp`
- User email: `.activity-user` or `.user-email`
- Action badge: `.action-badge[data-action="create"]`, `[data-action="update"]`, `[data-action="delete"]`
- Resource type: `.resource-type`
- Resource title: `.resource-title`
- Description: `.activity-description`
- Action filter select: `select[name="action"]` or `#action-filter`
- Action options: `option[value=""]` (All), `option[value="create"]`, `option[value="update"]`, `option[value="delete"]`
- Search input: `input[name="search"]` or `#activity-search`
- Search button: `button[type="submit"]` or `#search-submit`
- Clear filters button: `a[href="/admin/activity"]` or `button#clear-filters` (visible when HasFilters=true)
- Pagination: `.pagination`
- Page links: `a[data-page]` or `.page-link`
- Previous button: `.pagination .prev` or `a.page-prev`
- Next button: `.pagination .next` or `a.page-next`
- Current page indicator: `.page-item.active` or `span.current-page`
- Empty state: `.empty-state` or `#no-activity-message`
- No results message: `.no-results` or `#search-no-results`
- Verify absence: `button#new-activity`, `a[href*="/admin/activity/"][href*="/edit"]`, `.delete-activity`

## HTMX Interactions
- None - standard GET requests with query parameters
- No HTMX operations (read-only view)

## Dependencies
- Database activity_log table with columns: id, user_email, action (create/update/delete), resource_type, resource_title, description, created_at/timestamp
- Activity log automatically populated by system actions
- Template: templates/admin/pages/activity-log.html
- Handler: internal/handlers/activity.go (ListActivity)
- No create/update/delete handlers
- Query params: action, search, page
- Pagination limit: 50 entries per page
- HasFilters flag determined by presence of action or search params
