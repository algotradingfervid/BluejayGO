# Phase 03 - Dashboard Homepage

## Current State
- 3 hardcoded stat cards showing "0" (Products, Blog Posts, Partners)
- Welcome message card
- No quick actions, no activity feed, no pending items

## Goal
An informative, actionable dashboard with: stats overview, quick action buttons, recent activity feed, and content status summary.

## Layout (Single Column, Stacked Sections)

### Section 1: Stats Cards Row
4-column grid (2 on mobile), each card shows:
- Icon (Material Symbols)
- Count number (large, bold)
- Label (e.g., "Published Products")
- "View All ->" link

**Cards:**
1. Published Products (count from DB)
2. Published Blog Posts (count from DB)
3. Contact Submissions (count, with "X new" badge if unread)
4. Total Partners (count from DB)

Tooltip on each card: "Shows the count of [resource] currently published on your site."

### Section 2: Quick Actions
2x2 grid of action buttons (large, primary styled):
- "New Product" -> `/admin/products/new`
- "New Blog Post" -> `/admin/blog/posts/new`
- "New Solution" -> `/admin/solutions/new`
- "New Case Study" -> `/admin/case-studies/new`

Each button: icon + label, manual-border, manual-shadow, full-width within grid cell.

### Section 3: Recent Activity (5 most recent)
Simple list showing:
- Colored avatar circle with user initials
- Action description: "Updated Product 'Industrial Cleaner'"
- Relative timestamp: "5 min ago"
- Action badge: Created (green), Updated (blue), Deleted (red), Published (purple)

"View All Activity ->" link at bottom.

**Note:** This requires the Audit Trail (Phase 20) to exist. If building dashboard before audit trail, show a placeholder: "Activity tracking coming soon."

### Section 4: Content Status Summary
Small cards showing draft/pending counts:
- "3 Draft Blog Posts" -> link to filtered list
- "2 Draft Products" -> link to filtered list
- "5 Unread Contact Submissions" -> link

Only show cards where count > 0. If everything is published and no pending items, show a success message: "All content is published. Nice work!"

## Backend Changes
- `DashboardHandler` needs to query counts from DB:
  - `SELECT COUNT(*) FROM products WHERE status = 'published'`
  - `SELECT COUNT(*) FROM blog_posts WHERE status = 'published'`
  - `SELECT COUNT(*) FROM contact_submissions`
  - `SELECT COUNT(*) FROM partners`
  - Draft counts for each content type
- Pass all counts to template via `DashboardData` struct

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/dashboard.html` | Rewrite |
| `internal/handlers/admin/dashboard.go` | Add DB queries for counts |

## Dependencies
- Phase 01 (base layout)
- Phase 02 (sidebar) - dashboard link must be active
- Phase 20 (audit trail) - for activity feed section (can stub initially)
