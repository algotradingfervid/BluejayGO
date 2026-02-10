# Phase 13 - Whitepapers List & Form

## Current State
- List + form + topics management + downloads tracking page
- Basic CRUD

## Changes to List

### Filter Bar
- Search by title
- Topic dropdown filter
- Status filter
- Pagination

### Table
- Columns: Title, Topic, Status, Downloads (count), Updated, Actions
- Download count as a badge number

## Changes to Form

### Collapsible Sections

**Section 1: Basic Info** (open)
- Title
  - Tooltip: "Whitepaper title. Use a clear, benefit-driven title."
- Slug (auto-generated)
- Topic dropdown
  - Tooltip: "The subject category for this whitepaper."
- Status dropdown

**Section 2: Content** (open)
- Summary (textarea)
  - Tooltip: "Executive summary shown on the listing page. Convince readers to download."
- PDF Upload
  - Tooltip: "Upload the whitepaper PDF file. This is what users download."
- Preview of current PDF (filename + size)

**Section 3: Landing Page** (collapsed)
- Description (Trix editor)
  - Tooltip: "Full description for the whitepaper's landing page."
- Featured Image upload
  - Tooltip: "Cover image or thumbnail for the whitepaper."
- Key Takeaways (list of text items)
  - Tooltip: "Bullet-point benefits. Shown above the download button."

**Section 4: Gating** (collapsed)
- Require form submission to download (toggle)
  - Tooltip: "If enabled, visitors must submit their email before downloading."
- Form fields to collect (checkboxes: Name, Email, Company, Phone, Job Title)

**Section 5: SEO** (collapsed)
- Meta Title + Meta Description

## Whitepaper Topics Page
- Simple list: Name, Whitepaper Count, Actions
- Inline form (no separate page needed)
- Tooltip: "Topics help organize whitepapers by subject area."

## Downloads Tracking Page
- Read-only table: Whitepaper Title, Downloader Email, Company, Download Date
- Filter by whitepaper, date range
- Export to CSV button
- Tooltip: "Track who has downloaded your gated whitepapers."

## Files to Modify
| File | Action |
|------|--------|
| `templates/admin/pages/whitepapers_list.html` | Rewrite |
| `templates/admin/pages/whitepapers_form.html` | Rewrite |
| `templates/admin/pages/whitepaper_topics_list.html` | Polish |
| `templates/admin/pages/whitepaper_topics_form.html` | Polish |
| `templates/admin/pages/whitepapers_downloads.html` | Rewrite |
| `internal/handlers/admin/whitepapers.go` | Add filtering, pagination |

## Dependencies
- Phase 01, 02
