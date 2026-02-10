# Phase 10 - Blog Categories, Authors & Tags

## Current State
- Separate list + form pages for categories, authors, tags
- Basic CRUD with minimal styling
- Tags have HTMX quick-create and search

## Goal
Polish these supporting pages. They're simpler than main content pages, so keep them lightweight.

## Blog Categories

### List
- Simple table: Name, Slug, Post Count, Actions
- No filters needed (usually < 20 categories)
- "New Category" button

### Form
- Name input
  - Tooltip: "Category name displayed on the site (e.g., 'Industry News', 'Product Updates')."
- Slug (auto-generated)
  - Tooltip: "URL-friendly version. Auto-generated from name."
- Description (textarea)
  - Tooltip: "Optional description shown on the category archive page."

## Blog Authors

### List
- Table: Avatar + Name, Email, Bio (truncated), Post Count, Actions
- Author avatar: colored circle with initials

### Form
- Name
  - Tooltip: "Author's display name as shown on blog posts."
- Email
  - Tooltip: "Author's email. Not displayed publicly, used for Gravatar."
- Bio (textarea)
  - Tooltip: "Short author biography shown on blog posts. Keep to 2-3 sentences."
- Profile Image URL
  - Tooltip: "URL to author's profile photo. Leave blank to auto-generate from initials."

## Blog Tags

### List
- Tag cloud view (not table): show all tags as chips with post count
- Larger chip = more posts using that tag
- Click tag to edit
- "New Tag" button
- Search input to filter tags

### Form (inline modal or slide-out, not separate page)
- Name input
  - Tooltip: "Tag name. Use lowercase, hyphen-separated (e.g., 'machine-learning')."
- Slug (auto-generated)
- Quick-create: type in tag search on blog post form, hit enter to create new

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/blog_categories_list.html` | Polish |
| `templates/admin/pages/blog_categories_form.html` | Add tooltips |
| `templates/admin/pages/blog_authors_list.html` | Polish, add avatar |
| `templates/admin/pages/blog_authors_form.html` | Add tooltips |
| `templates/admin/pages/blog_tags_list.html` | Redesign as tag cloud |
| `templates/admin/partials/tag_chip.html` | Polish |

## Dependencies
- Phase 01, 02
