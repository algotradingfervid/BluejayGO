# Phase 18 - Media Library (NEW FEATURE)

## Current State
- No media library exists
- Images uploaded per-resource (product images, blog featured images, etc.)
- Files stored in `public/uploads/` directory
- No way to browse or reuse previously uploaded images

## Goal
Central media library for browsing, uploading, and reusing images across all content types.

## Media Library Page

### Layout
- Full-width page (no sidebar folders needed initially - keep it simple)
- Toolbar at top
- Grid/List toggle view

### Toolbar
- "Upload" button (primary, opens upload modal)
- View toggle: Grid / List icons
- Search input (by filename)
  - Tooltip: "Search media files by name."
- Sort dropdown: Newest / Oldest / Name A-Z / Largest

### Grid View
- 4-column grid (3 on tablet, 2 on mobile)
- Each card:
  - Image preview (cover fit, square aspect)
  - Filename below (truncated)
  - File size
  - Hover overlay: shows "Select" / "Delete" buttons

### List View
- Table: Thumbnail (48x48), Filename, Type, Size, Dimensions, Upload Date, Used In (count), Actions

### Upload Modal
- Drag-and-drop zone (dashed border, "Drop files here or click to browse")
- Multi-file support
- Upload progress bar per file
- File size limit display: "Max 10MB per file"
- Supported formats: JPG, PNG, SVG, GIF, PDF
- After upload, files appear in the grid immediately

### File Detail Modal (click on any file)
- Large preview
- Metadata: filename, dimensions, size, upload date, format
- Alt text input (editable)
  - Tooltip: "Accessibility text for this image. Describe what the image shows."
- "Copy URL" button
- "Delete" button (with confirmation)
- "Used in" list: shows which content items reference this file

### Integration with Content Forms
- Anywhere there's an image upload field (product images, blog featured image, etc.), add a "Choose from Library" button alongside the regular upload
- Clicking it opens the media library in a modal/overlay
- User can select an existing image or upload new
- Selected image URL is inserted into the form field

## Backend Implementation

### New Database Table
```sql
CREATE TABLE media_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    mime_type TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    alt_text TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### New Routes
- `GET /admin/media` - Media library page
- `POST /admin/media/upload` - Upload file(s), return JSON
- `GET /admin/media/:id` - Get file details (JSON for modal)
- `PUT /admin/media/:id` - Update alt text
- `DELETE /admin/media/:id` - Delete file
- `GET /admin/media/browse` - HTMX partial for modal picker in forms

### New Handler
- `MediaHandler` with CRUD operations
- Upload handler saves to `public/uploads/media/` with UUID filenames
- Image dimension detection on upload (Go's `image` package)
- Scan existing `public/uploads/` directory on first load to populate media_files table

### Migration for Existing Files
- On first access, scan `public/uploads/` recursively
- Insert entries for all existing files into `media_files` table
- One-time migration, idempotent

## Files to Create
| File | Action |
|------|--------|
| `templates/admin/pages/media_library.html` | Create |
| `templates/admin/partials/media_picker.html` | Create (reusable modal) |
| `internal/handlers/admin/media.go` | Create |
| `db/migrations/030_media_library.up.sql` | Create |
| `db/migrations/030_media_library.down.sql` | Create |
| `db/queries/media.sql` | Create |
| `cmd/server/main.go` | Add routes |

## Dependencies
- Phase 01, 02
- Should be built before content form redesigns ideally, but can be added after
