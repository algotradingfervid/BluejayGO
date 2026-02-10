# Test Plan: Contact Form Submissions Management

## Summary
Tests the contact submissions manager with filtering, search, status updates, bulk operations, and detail view with prev/next navigation.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Database seeded with 15 contact submissions
- Submission types: contact (general inquiries), rfq (request for quote)
- Submission statuses: new, read, replied
- Pagination: 25 submissions per page
- New submissions highlighted in yellow

## User Journey Steps
1. Navigate to /admin/contact/submissions
2. View list with tab counts: All, Contact, RFQ, New
3. Click tabs to filter by type
4. Use status filter (new/read/replied)
5. Use search to find submissions
6. View status badges (yellow "New", gray "Read", green "Replied")
7. Click submission to view detail (GET /admin/contact/submissions/:id)
8. Navigate prev/next between submissions in detail view
9. Update status and add notes (POST /admin/contact/submissions/:id/status)
10. Use bulk mark-read (POST /admin/contact/submissions/bulk-mark-read)
11. Delete submission (DELETE /admin/contact/submissions/:id via hx-delete)

## Test Cases

### Happy Path - List View
- **Load all submissions**: Verifies GET /admin/contact/submissions shows all 15 seeded submissions
- **Tab counts**: Verifies tab badges show correct counts for All/Contact/RFQ/New
- **Filter by Contact type**: Clicks Contact tab, verifies only type=contact shown
- **Filter by RFQ type**: Clicks RFQ tab, verifies only type=rfq shown
- **Filter by New status**: Selects status=new, verifies only new submissions shown
- **Filter by Read status**: Selects status=read, verifies only reviewed submissions shown
- **Filter by Replied status**: Selects status=replied, verifies only replied submissions shown
- **Search submissions**: Enters email or name, verifies filtered results
- **Combined filters**: Applies type=rfq + status=new, verifies both filters applied
- **New submission highlight**: Verifies new submissions have yellow background/border
- **Status badges**: Verifies correct badge color/text for each status
- **Pagination**: Tests pagination with >25 submissions

### Happy Path - Detail View
- **View submission detail**: Clicks submission, verifies GET /admin/contact/submissions/:id loads detail
- **View all fields**: Confirms name, email, phone, company, message, type, status, created_at displayed
- **Prev/next navigation**: Clicks "Previous" and "Next" buttons, verifies navigation between submissions
- **Prev disabled**: On first submission, verifies "Previous" button disabled
- **Next disabled**: On last submission, verifies "Next" button disabled
- **Return to list**: Clicks back button, returns to /admin/contact/submissions

### Happy Path - Status Update
- **Update status to read**: Changes status dropdown to "read", adds notes, submits, verifies update
- **Update status to replied**: Changes status to "replied", adds notes, verifies update
- **Add notes**: Enters notes in textarea, verifies saved with status update
- **Update without notes**: Changes status without notes, verifies update
- **Status badge update**: After status change, verifies badge updates in list view

### Happy Path - Bulk Operations
- **Bulk mark as read**: Clicks "Mark All New as Read" button, verifies POST /admin/contact/submissions/bulk-mark-read
- **Verify bulk update**: After bulk operation, verifies all "new" submissions now "read"
- **Bulk with no new**: When no new submissions exist, verifies button disabled or hidden

### Happy Path - Delete
- **Delete submission**: Clicks delete button with hx-delete, verifies confirmation
- **Confirm delete**: Confirms deletion, verifies DELETE /admin/contact/submissions/:id
- **Submission removed**: After delete, verifies submission removed from list
- **Tab count update**: After delete, verifies tab counts update

### Edge Cases / Error States
- **Empty submissions**: Tests display when no submissions exist, verifies empty state
- **Search no results**: Searches for non-existent email, verifies "no results" message
- **Filter no results**: Applies filters with no matches, verifies empty state
- **Very long message**: Tests submission with 5000+ char message, verifies truncation/display
- **Missing phone**: Tests submission without phone number, verifies optional handling
- **Missing company**: Tests submission without company, verifies optional handling
- **Invalid status filter**: Tests if invalid status param is rejected
- **Invalid type filter**: Tests if invalid type param is rejected
- **View non-existent ID**: Navigates to /admin/contact/submissions/99999, verifies 404
- **Update non-existent submission**: POSTs status update to non-existent ID, verifies error
- **Delete non-existent submission**: Attempts delete on non-existent ID, verifies error
- **Empty status update**: Submits status form without selecting status, checks validation
- **Very long notes**: Enters 2000+ char notes, checks validation/limit
- **Status update redirect**: After status update, verifies redirect back to detail or list
- **Bulk operation with zero new**: Tests bulk-mark-read when no new submissions, verifies handling
- **Delete with HTMX**: Verifies hx-delete removes row from DOM without page reload
- **Delete confirmation**: Verifies confirmation prompt before deletion
- **Pagination with filters**: Applies filters, navigates pages, verifies filter persistence
- **Tab persistence**: Clicks tab, performs action, returns, verifies tab remains active
- **New submission count**: After marking as read, verifies "New" tab count decrements
- **Status transition**: Tests all status transitions (new→read, read→replied, etc.)

## Selectors & Elements

### List View
- Submissions table: `#submissions-list` or `.submissions-table`
- Submission row: `.submission-row[data-id]` or `tr[data-submission-id]`
- New highlight: `.submission-row.new` or `tr.status-new` (yellow background)
- Status badge: `.status-badge[data-status="new"]`, `[data-status="read"]`, `[data-status="replied"]`
- Type column: `.submission-type`
- Name column: `.submission-name`
- Email column: `.submission-email`
- Date column: `.submission-date`
- View link: `a[href="/admin/contact/submissions/{id}"]`
- Delete button: `button[hx-delete="/admin/contact/submissions/{id}"]`
- Tabs: `.tab[data-type="all"]`, `[data-type="contact"]`, `[data-type="rfq"]`, `[data-type="new"]`
- Tab counts: `.tab-count` or badge within tab
- Status filter: `select[name="status"]` or `#status-filter`
- Type filter: `select[name="type"]` or `#type-filter` (may be implicit via tabs)
- Search input: `input[name="search"]` or `#submissions-search`
- Search button: `button[type="submit"]` or `#search-submit`
- Bulk mark read button: `button#bulk-mark-read` or `form[action="/admin/contact/submissions/bulk-mark-read"]`
- Pagination: `.pagination`
- Page links: `a[data-page]`
- Empty state: `.empty-state` or `#no-submissions`

### Detail View
- Detail container: `#submission-detail`
- Submission ID: `#submission-id` or `.detail-id`
- Name field: `.detail-name`
- Email field: `.detail-email`
- Phone field: `.detail-phone`
- Company field: `.detail-company`
- Message field: `.detail-message` or `pre.message-content`
- Type field: `.detail-type`
- Status field: `.detail-status`
- Created date: `.detail-created-at`
- Previous button: `a#prev-submission` or `.nav-prev`
- Next button: `a#next-submission` or `.nav-next`
- Back to list: `a[href="/admin/contact/submissions"]` or `.back-link`

### Status Update Form
- Form: `form[action="/admin/contact/submissions/{id}/status"][method="POST"]`
- Status select: `select[name="status"]`
- Status options: `option[value="new"]`, `option[value="read"]`, `option[value="replied"]`
- Notes textarea: `textarea[name="notes"]`
- Submit button: `button[type="submit"]` or `#update-status`
- Success message: `.alert-success`

## HTMX Interactions
- **hx-delete**: Delete button uses `hx-delete="/admin/contact/submissions/{id}"`
- **hx-confirm**: Confirmation message before delete
- **hx-target**: Targets parent row for removal
- **hx-swap**: Uses `outerHTML` or `delete` to remove from DOM

## Dependencies
- Database seeded with 15 submissions
- contact_submissions table columns: id, name, email, phone, company, message, type (contact/rfq), status (new/read/replied), created_at
- May have separate notes field or notes stored with status updates
- Template: templates/admin/pages/contact-submissions-list.html, contact-submission-detail.html
- Handler: internal/handlers/contact.go (ListSubmissions, GetSubmission, UpdateSubmissionStatus, BulkMarkRead, DeleteSubmission)
- HTMX library loaded
- Tab counts calculated from database queries
- Prev/next navigation requires ordered list context
- Bulk operation updates all status=new to status=read
- Pagination limit: 25 per page
