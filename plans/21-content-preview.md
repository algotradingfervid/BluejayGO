# Phase 21 - Content Preview System (NEW FEATURE)

## Current State
- No preview capability
- Must publish content to see how it looks on the public site

## Goal
"Preview" button on content forms that opens the public page in a new tab showing draft/unsaved content.

## How It Works

### User Flow
1. User is editing a product (draft status)
2. Clicks "Preview" button in the form header
3. New browser tab opens showing the public product page with the current form data
4. User can review, close the tab, and continue editing

### Implementation Approach: Token-Based Preview

1. When "Preview" is clicked, the form data is saved as a temporary preview (not published)
2. A unique preview token is generated
3. New tab opens: `/preview/{resource_type}/{id}?token={preview_token}`
4. The public handler checks for the preview token
5. If valid token, renders the page using draft/preview data instead of published data
6. Preview tokens expire after 30 minutes

### Alternative (Simpler): Save as Draft + View

Since this is simpler:
1. "Preview" button first saves the content as draft (HTMX POST)
2. Then opens `/products/{slug}?preview=true&token={session_token}` in new tab
3. Public handler: if `preview=true` and user has valid admin session, show draft content
4. Preview banner at top of page: "This is a preview. This content is not yet published."

**Go with the simpler approach** - leverages existing session auth.

### Preview Banner
When viewing a preview, inject a sticky banner at the top:
- Yellow background, black text
- "PREVIEW MODE - This content is not published yet"
- "Edit" button links back to admin form
- "Close Preview" button

### Which Content Types Get Preview
- Products -> `/products/{slug}`
- Blog Posts -> `/blog/{slug}`
- Solutions -> `/solutions/{slug}`
- Case Studies -> `/case-studies/{slug}`
- Whitepapers -> `/whitepapers/{slug}`

### Preview Button Placement
- In the form header area, next to "Save" button
- Secondary style button with eye icon
- Only shown when editing existing content (not when creating new)
- Tooltip: "Preview how this content will look on the public site."

## Backend Changes

### Modify Public Handlers
For each previewable content type's public handler:
```go
func (h *ProductsPublicHandler) Show(c echo.Context) error {
    slug := c.Param("slug")
    preview := c.QueryParam("preview") == "true"

    if preview {
        // Check admin session is valid
        if !isAuthenticated(c) {
            return c.Redirect(302, "/admin/login")
        }
        // Load product regardless of status (include drafts)
        product, err := h.queries.GetProductBySlugIncludeDrafts(ctx, slug)
        // Render with preview banner
    } else {
        // Normal: only published content
        product, err := h.queries.GetProductBySlug(ctx, slug)
    }
}
```

### New Queries
- `GetProductBySlugIncludeDrafts` - same as GetProductBySlug but without `WHERE status = 'published'`
- Same pattern for blog posts, solutions, case studies, whitepapers

### Preview Banner Partial
- `templates/partials/preview-banner.html` - injected into public templates when preview mode

## Files to Create/Modify
| File | Action |
|------|--------|
| `templates/partials/preview-banner.html` | Create |
| `templates/admin/pages/products_form.html` | Add preview button |
| `templates/admin/pages/blog_post_form.html` | Add preview button |
| `templates/admin/pages/solutions_form.html` | Add preview button |
| `templates/admin/pages/case_studies_form.html` | Add preview button |
| `templates/admin/pages/whitepapers_form.html` | Add preview button |
| `internal/handlers/public/*.go` | Add preview mode logic (5 files) |
| `db/queries/*.sql` | Add "include drafts" variants (5 queries) |

## Dependencies
- Phase 01, 02
- Respective content form phases should be done first (07, 09, 11, 12, 13)
- Public site templates must exist (not modifying public site design, just adding preview logic)
