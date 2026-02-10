# Phase 09 - Blog Posts List & Form

## Current State

### List
- Table with: Title (+ slug), Category, Author, Status, Published date, Actions
- Fixed column widths, no search/filter, no pagination

### Form
- Long form: title, slug, excerpt, Trix body editor, image URL, category, author, status, reading time, published date, meta description
- HTMX-powered tag and product multi-select with search
- Chip-based selection UI

## Changes to List Page

### Filter Bar
- Search input (title, slug)
  - Tooltip: "Search blog posts by title or URL slug."
- Category dropdown
- Author dropdown
- Status dropdown (All / Published / Draft)
- "Clear filters" link

### Table Improvements
- Add thumbnail column (featured image, 80x53 aspect ratio)
- Show excerpt preview below title (truncated, gray text)
- Author column: avatar circle with initials + name
- Better date formatting
- Pagination: 15 per page

### Empty/Filtered States
- No posts: "No blog posts yet. Write your first post!" + button
- No filter results: "No posts match your filters." + clear link

## Changes to Form Page

### Two-Column Layout (desktop only, stacks on mobile)

**Left Column (2/3 width) - Content:**

Section 1: Title & Slug (always open)
- Title (large input, text-2xl)
  - Tooltip: "The headline of your blog post. Make it engaging and descriptive."
- Slug (auto-generated, with "Edit" button to unlock)
  - Tooltip: "URL path for this post. Auto-generated from title. Only edit if you need a custom URL."

Section 2: Content (always open)
- Excerpt (textarea, 300-char counter)
  - Tooltip: "A brief summary shown on blog listing pages and in search results. Keep it compelling."
- Body (Trix rich text editor - full width)
  - Tooltip: "The main content of your post. Use the toolbar for formatting, images, and links."

**Right Column (1/3 width) - Sidebar Cards:**

Card 1: Publish
- Status dropdown (Draft / Published)
- Published Date (datetime picker)
  - Tooltip: "When this post goes live. Leave blank to publish immediately."
- Reading Time (number, minutes)
  - Tooltip: "Estimated reading time. Helps readers decide to read now or save for later."
- Action buttons: "Save Draft" / "Publish"

Card 2: Featured Image
- Upload area with drag-and-drop
- Preview of current image
- Alt text input
  - Tooltip: "Describe the image for accessibility. Also used when image can't load."
- Remove button

Card 3: Taxonomy
- Category dropdown
  - Tooltip: "The main category this post belongs to."
- Author dropdown
  - Tooltip: "Who wrote this post. Shown on the post page with their bio."
- Tags (HTMX search, chip selection)
  - Tooltip: "Tags help readers find related posts. Add up to 10 tags."

Card 4: SEO (collapsible, collapsed by default)
- Meta Title (70-char counter)
- Meta Description (160-char counter)
- Google search preview (live rendered)
  - Shows: title in blue, slug in green, description in gray

### Character Counter Behavior
- Green bar: under 80% of limit
- Yellow bar: 80-100% of limit
- Red bar: over limit
- Number shows "45/70" format

## Backend Changes
- Add pagination + filtering to blog posts list query
- Add count query for pagination

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/blog_posts_list.html` | Rewrite |
| `templates/admin/pages/blog_post_form.html` | Rewrite (two-column) |
| `internal/handlers/admin/blog_posts.go` | Add filtering, pagination |

## Dependencies
- Phase 01, 02
