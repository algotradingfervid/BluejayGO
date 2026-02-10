# Test Plan: Media Library Management

## Summary
Tests the comprehensive media library with upload, search, filter, view modes, file details, alt text editing, URL copying, and deletion functionality using HTMX and JavaScript.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Media library initialized (may be empty or have sample files)
- File upload limit: 10MB per file
- Allowed extensions: .jpg, .jpeg, .png, .gif, .svg, .pdf, .webp
- Pagination: 24 files per page
- JavaScript functions: openUploadModal(), handleFiles(), handleDrop(), openDetailModal(id), updateAltText(), deleteMedia(id), copyUrl(), setView(grid/list), formatFileSize()

## User Journey Steps
1. Navigate to /admin/media
2. View media library in default grid view
3. Use search to filter files by name
4. Use sort dropdown (newest/oldest/name/largest)
5. Navigate pagination (24 files per page)
6. Click "Upload" to open upload modal
7. Drag-drop files to drop zone or select files
8. View upload progress via XHR
9. After upload, verify files appear in library
10. Switch between grid and list views
11. Click file to open detail modal
12. View file preview, metadata (size, dimensions, type, date)
13. Edit alt_text and save via PUT /admin/media/:id
14. Copy file URL to clipboard
15. Delete file via DELETE /admin/media/:id
16. Use media picker partial via GET /admin/media/browse (HTMX)

## Test Cases

### Happy Path - Library Browsing
- **Load media library**: Verifies GET /admin/media shows files with grid view
- **Search files**: Enters search term, verifies filtered results
- **Sort by newest**: Selects "newest" sort, verifies file order
- **Sort by oldest**: Selects "oldest" sort, verifies reverse chronological order
- **Sort by name**: Selects "name" sort, verifies alphabetical order
- **Sort by largest**: Selects "largest" sort, verifies files by size descending
- **Pagination**: Navigates to page 2, verifies next 24 files load
- **Switch to list view**: Clicks list view button, verifies table layout
- **Switch to grid view**: Clicks grid view button, verifies grid cards layout

### Happy Path - File Upload
- **Open upload modal**: Clicks upload button, verifies modal appears
- **Upload single file**: Selects 1 JPG file, verifies upload success and JSON response
- **Upload multiple files**: Selects 5 files at once, verifies all upload
- **Drag-drop upload**: Drags 3 PNG files to drop zone, verifies handleDrop() processes
- **Upload progress**: Monitors XHR progress during upload, verifies progress indicator
- **Upload different types**: Uploads JPG, PNG, GIF, SVG, PDF, WEBP, verifies all accepted
- **Close upload modal**: Closes modal after upload, verifies library refreshes

### Happy Path - File Details
- **Open detail modal**: Clicks file in grid, verifies openDetailModal(id) fetches JSON
- **View image preview**: Opens image file, verifies preview displays
- **View PDF**: Opens PDF file, verifies appropriate preview or icon
- **View metadata**: Verifies file size, dimensions, type, upload date display
- **Edit alt text**: Updates alt_text field, clicks save, verifies PUT /admin/media/:id
- **Copy URL**: Clicks "Copy URL" button, verifies copyUrl() copies to clipboard
- **Close detail modal**: Closes modal, returns to library

### Happy Path - File Deletion
- **Delete from detail modal**: Opens file, clicks delete, confirms, verifies DELETE /admin/media/:id
- **Delete confirmation**: Verifies confirmation prompt before deleteMedia(id)
- **File removed from library**: After delete, verifies file no longer in grid/list
- **File record deleted**: Verifies database record and physical file removed

### Happy Path - Media Picker (HTMX)
- **Load media picker**: Triggers GET /admin/media/browse, verifies HTMX loads media_picker.html partial
- **Select file from picker**: Clicks file in picker, verifies selection callback
- **Close picker**: Closes picker modal or panel

### Edge Cases / Error States
- **Empty library**: Tests display when no files exist, verifies empty state message
- **Search no results**: Searches for non-existent file, verifies "no results" message
- **Upload oversized file**: Uploads file >10MB, verifies error message
- **Upload invalid extension**: Uploads .exe or .zip file, verifies rejection
- **Upload with spaces in name**: Uploads "my file.jpg", verifies filename handling
- **Upload duplicate filename**: Uploads file with existing name, verifies handling (overwrite/rename)
- **Upload failure**: Simulates network error during upload, verifies error handling
- **XHR progress error**: Tests upload progress bar edge cases
- **Invalid file ID**: Calls openDetailModal() with non-existent ID, verifies error
- **GET /admin/media/:id not found**: Requests JSON for deleted file, verifies 404
- **Update alt text empty**: Saves alt_text as empty string, verifies handling
- **Update alt text very long**: Enters 500+ char alt text, checks validation/limit
- **PUT /admin/media/:id error**: Simulates update failure, verifies error message
- **Delete non-existent file**: Attempts DELETE /admin/media/:id for missing file, verifies error
- **Copy URL to clipboard failure**: Tests copyUrl() when clipboard API unavailable
- **Grid view overflow**: Tests display with 100+ files
- **List view performance**: Tests list view with many files
- **File size formatting**: Verifies formatFileSize() correctly displays B, KB, MB
- **Drag-drop invalid files**: Drops non-image files, verifies handleDrop() validation
- **Multiple concurrent uploads**: Uploads 10 files simultaneously, verifies handling
- **Pagination edge**: Tests page=999 when only 3 pages exist, verifies handling

## Selectors & Elements
- Media grid: `#media-grid`
- Media list: `#media-list`
- Grid view button: `button[data-view="grid"]` or `#view-grid`
- List view button: `button[data-view="list"]` or `#view-list`
- Search input: `input[name="search"]` or `#media-search`
- Sort select: `select[name="sort"]` or `#media-sort`
- Sort options: `option[value="newest"]`, `option[value="oldest"]`, `option[value="name"]`, `option[value="largest"]`
- Pagination links: `.pagination a[data-page]`
- Upload button: `button#upload-media` or `#open-upload`
- Upload modal: `#upload-modal`
- Drop zone: `#drop-zone` or `.upload-drop-zone`
- File input: `input[type="file"][name="files[]"][multiple]`
- Upload progress: `.upload-progress` or `#upload-progress-bar`
- Detail modal: `#detail-modal`
- File preview: `#file-preview` or `.file-preview img`
- File metadata: `.file-metadata` (size, dimensions, type, date)
- Alt text input: `input[name="alt_text"]` or `#alt-text-input`
- Save alt text button: `button#save-alt-text`
- Copy URL button: `button#copy-url` or `.copy-url-btn`
- Delete button (detail): `button#delete-file` or `.delete-media-btn`
- Media file card: `.media-file-card[data-id]`
- Media file row: `tr.media-file-row[data-id]`
- Empty state: `.empty-state` or `#no-media-message`
- Error message: `.error-message` or `.alert-error`
- Success message: `.success-message` or `.alert-success`

## HTMX Interactions
- **Media picker**: GET /admin/media/browse returns HTMX partial (media_picker.html)
- **hx-get**: May use HTMX for lazy loading pagination or search results
- **hx-trigger**: Upload modal may use HTMX for dynamic content

## Dependencies
- Database media/files table with columns: id, filename, filepath, file_size, mime_type, alt_text, uploaded_at
- File storage directory for uploaded media
- File upload handler: POST /admin/media/upload (multipart form)
- File detail endpoint: GET /admin/media/:id (returns JSON)
- Update alt text: PUT /admin/media/:id (JSON body: {alt_text: "..."})
- Delete file: DELETE /admin/media/:id
- Media picker partial: GET /admin/media/browse
- JavaScript file with all listed functions
- Templates: templates/admin/pages/media-library.html, partials/media_picker.html
- Handler: internal/handlers/media.go (ListMedia, UploadMedia, GetMediaFile, UpdateMediaFile, DeleteMediaFile, MediaBrowser)
- HTMX library loaded
- Clipboard API for copyUrl()
- XHR upload with progress events
