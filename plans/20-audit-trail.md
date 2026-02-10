# Phase 20 - Activity Log / Audit Trail (NEW FEATURE)

## Current State
- No activity tracking exists
- No way to know who changed what

## Goal
Basic activity log: who did what and when. No diff view, no rollback (keep it simple).

## Activity Log Page

### Filter Bar
- User dropdown (if multiple admins in future)
- Action Type: All / Created / Updated / Deleted / Published / Login
  - Tooltip: "Filter by the type of action performed."
- Date Range picker
- Search input (searches description text)

### Tabs with Counts
- All Activity (total count)
- Create
- Update
- Delete
- Login/Logout

### Table
- Timestamp (e.g., "Feb 10, 2026 2:34 PM")
- User: avatar + name
- Action badge (color-coded):
  - Created: green
  - Updated: blue
  - Deleted: red
  - Published: purple
  - Login: teal
- Description: "Updated Product 'Industrial Cleaner'" (with link to the item)
- No IP address needed for solo admin (keep it simple)

### Pagination
- 50 per page
- "Showing 1-50 of X" + page numbers

### Empty State
- "No activity recorded yet. Actions will appear here as you use the admin panel."

## Backend Implementation

### Database Table
```sql
CREATE TABLE activity_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES admin_users(id),
    action TEXT NOT NULL,        -- 'created', 'updated', 'deleted', 'published', 'login', 'logout'
    resource_type TEXT NOT NULL,  -- 'product', 'blog_post', 'solution', etc.
    resource_id INTEGER,
    resource_title TEXT,          -- denormalized for display even if resource is deleted
    description TEXT NOT NULL,    -- human-readable: "Updated Product 'Industrial Cleaner'"
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activity_log_created_at ON activity_log(created_at DESC);
CREATE INDEX idx_activity_log_action ON activity_log(action);
```

### Logging Integration
Add activity logging calls to every handler's Create/Update/Delete methods:
```go
// In every handler after successful DB operation:
activityLog.Log(userID, "updated", "product", productID, product.Name,
    fmt.Sprintf("Updated Product '%s'", product.Name))
```

This means touching every existing handler file to add one line after each successful mutation. The logging function should be a simple service method that inserts a row.

### Activity Log Service
```go
type ActivityLogService struct {
    db *sql.DB
}

func (s *ActivityLogService) Log(userID int, action, resourceType string, resourceID int, resourceTitle, description string) error {
    // INSERT INTO activity_log ...
}

func (s *ActivityLogService) List(filters ActivityLogFilters) ([]ActivityLogEntry, int, error) {
    // SELECT with filters, pagination, count
}
```

### Routes
- `GET /admin/activity` - Activity log page

### Handler Files to Add Logging To
Every admin handler needs `activityLog.Log()` calls added:
- `products.go` (Create, Update, Delete)
- `product_details.go` (all sub-resource mutations)
- `solutions.go` (Create, Update, Delete + sub-resources)
- `case_studies.go`
- `blog_posts.go`
- `blog_categories.go`, `blog_authors.go`, `blog_tags.go`
- `whitepapers.go`, `whitepaper_topics.go`
- `partners.go`, `partner_tiers.go`
- `about.go` (all sub-forms)
- `homepage.go` (all sub-resources)
- `settings.go`
- `contact.go`
- `auth.go` (login/logout)

## Files to Create/Modify
| File | Action |
|------|--------|
| `templates/admin/pages/activity_log.html` | Create |
| `internal/handlers/admin/activity.go` | Create |
| `internal/services/activity_log.go` | Create |
| `db/migrations/032_activity_log.up.sql` | Create |
| `db/queries/activity_log.sql` | Create |
| `cmd/server/main.go` | Add route, inject service |
| All 20+ handler files | Add logging calls (1 line each) |

## Dependencies
- Phase 01, 02
- Ideally done early so all subsequent phase work automatically gets logged
- Dashboard (Phase 03) references activity data for recent activity feed
