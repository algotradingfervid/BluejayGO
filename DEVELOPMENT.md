# Bluejay CMS — Development Guide

A practical developer guide for building features in the Bluejay CMS.

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Running the Application](#running-the-application)
3. [End-to-End Feature Walkthrough](#end-to-end-feature-walkthrough)
4. [Handler Conventions](#handler-conventions)
5. [Template Conventions](#template-conventions)
6. [HTMX Development Patterns](#htmx-development-patterns)
7. [Service Layer Patterns](#service-layer-patterns)
8. [Testing Approach](#testing-approach)
9. [Common Gotchas](#common-gotchas)
10. [Code Style and Conventions](#code-style-and-conventions)
11. [Debugging Tips](#debugging-tips)

---

## Development Environment Setup

### Prerequisites

1. **Go 1.25+**
   ```bash
   go version  # Should be 1.25 or higher
   ```

2. **sqlc** — Type-safe SQL code generator
   ```bash
   # Install sqlc
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

   # Verify installation
   sqlc version
   ```

3. **air** (optional) — Hot-reload for development
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

4. **SQLite3** — For manual database inspection
   ```bash
   # macOS
   brew install sqlite3

   # Ubuntu/Debian
   sudo apt-get install sqlite3
   ```

### Initial Setup

1. Clone the repository and navigate to the project directory:
   ```bash
   cd /path/to/bluejay-cms
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. The database will be created automatically on first run. Migrations run automatically at startup.

---

## Running the Application

### Standard Run

```bash
make run
# or
go run cmd/server/main.go
```

The server starts on `http://localhost:28090`

### Hot-Reload Development

For automatic reload on file changes:

```bash
make dev
# or
air
```

Air watches for changes in `.go`, `.html`, and `.sql` files and automatically rebuilds.

### Build Binary

```bash
make build
# Creates: bin/bluejay-cms
```

### Testing

```bash
# Run all tests
make test

# Run specific package tests
go test -v ./internal/handlers/admin/...
go test -v ./internal/services/...

# Run with coverage
go test -cover ./...
```

### Other Useful Commands

```bash
# Generate sqlc code (after SQL changes)
make sqlc

# Clean build artifacts and database
make clean
```

---

## End-to-End Feature Walkthrough

This section walks through adding a complete new feature: **Event Management** (a hypothetical feature to demonstrate the full pattern).

### Step 1: Create Migration Files

Create two migration files in `db/migrations/`:

**`020_events.up.sql`:**
```sql
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    location TEXT,
    event_date DATETIME NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_slug ON events(slug);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_date ON events(event_date);
```

**`020_events.down.sql`:**
```sql
DROP TABLE IF EXISTS events;
```

**Migration naming convention:**
- Format: `{number}_{description}.{up|down}.sql`
- Number must be sequential (001, 002, etc.)
- Always create both `.up.sql` and `.down.sql`

The migrations run automatically on server start. No manual migration command needed.

### Step 2: Write sqlc Queries

Create `db/queries/events.sql`:

```sql
-- name: CreateEvent :one
INSERT INTO events (title, slug, description, location, event_date, status)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetEvent :one
SELECT * FROM events WHERE id = ? LIMIT 1;

-- name: GetEventBySlug :one
SELECT * FROM events WHERE slug = ? LIMIT 1;

-- name: ListEvents :many
SELECT * FROM events
WHERE status = 'published'
ORDER BY event_date DESC
LIMIT ? OFFSET ?;

-- name: ListEventsAdmin :many
SELECT * FROM events
ORDER BY created_at DESC;

-- name: UpdateEvent :exec
UPDATE events
SET title = ?, slug = ?, description = ?, location = ?,
    event_date = ?, status = ?
WHERE id = ?;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = ?;
```

**Query naming conventions:**
- `:one` — Returns a single row (or error if not found)
- `:many` — Returns a slice of rows
- `:exec` — Executes without returning data (INSERT/UPDATE/DELETE)

### Step 3: Generate Go Code with sqlc

```bash
make sqlc
# or
sqlc generate
```

This generates type-safe Go code in `db/sqlc/events.sql.go` with methods like:
- `CreateEvent(ctx, params) (Event, error)`
- `GetEvent(ctx, id) (Event, error)`
- `ListEvents(ctx, params) ([]Event, error)`

**Important:** You must run this after ANY changes to SQL query files.

### Step 4: Create Handler Struct and Methods

Create `internal/handlers/admin/events.go`:

```go
package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

// EventsHandler manages all HTTP handlers for event CRUD operations
type EventsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

// NewEventsHandler creates a new EventsHandler with required dependencies
func NewEventsHandler(queries *sqlc.Queries, logger *slog.Logger) *EventsHandler {
	return &EventsHandler{
		queries: queries,
		logger:  logger,
	}
}

// List handles GET /admin/events
// Renders the events list page
func (h *EventsHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	events, err := h.queries.ListEventsAdmin(ctx)
	if err != nil {
		h.logger.Error("failed to list events", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "admin/pages/events_list.html", map[string]interface{}{
		"Title": "Manage Events",
		"Items": events,
	})
}

// New handles GET /admin/events/new
// Renders the event creation form
func (h *EventsHandler) New(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/events_form.html", map[string]interface{}{
		"Title":      "New Event",
		"FormAction": "/admin/events",
		"Item":       nil,
	})
}

// Create handles POST /admin/events
// Processes form submission and creates a new event
func (h *EventsHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	title := c.FormValue("title")
	slug := makeSlug(title)
	description := c.FormValue("description")
	location := c.FormValue("location")
	eventDateStr := c.FormValue("event_date")
	status := c.FormValue("status")

	// Parse event date
	eventDate, err := time.Parse("2006-01-02T15:04", eventDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event date")
	}

	_, err = h.queries.CreateEvent(ctx, sqlc.CreateEventParams{
		Title:       title,
		Slug:        slug,
		Description: description,
		Location:    sql.NullString{String: location, Valid: location != ""},
		EventDate:   eventDate,
		Status:      status,
	})
	if err != nil {
		h.logger.Error("failed to create event", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Log activity for audit trail
	logActivity(c, "created", "event", 0, title, "Created Event '%s'", title)

	// Redirect back to list page
	return c.Redirect(http.StatusSeeOther, "/admin/events")
}

// Edit handles GET /admin/events/:id/edit
// Renders the event edit form with existing data
func (h *EventsHandler) Edit(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	event, err := h.queries.GetEvent(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Event not found")
	}

	return c.Render(http.StatusOK, "admin/pages/events_form.html", map[string]interface{}{
		"Title":      "Edit Event",
		"FormAction": fmt.Sprintf("/admin/events/%d", id),
		"Item":       event,
	})
}

// Update handles POST /admin/events/:id
// Processes form submission and updates existing event
func (h *EventsHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	title := c.FormValue("title")
	slug := makeSlug(title)
	description := c.FormValue("description")
	location := c.FormValue("location")
	eventDateStr := c.FormValue("event_date")
	status := c.FormValue("status")

	eventDate, err := time.Parse("2006-01-02T15:04", eventDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event date")
	}

	err = h.queries.UpdateEvent(ctx, sqlc.UpdateEventParams{
		ID:          id,
		Title:       title,
		Slug:        slug,
		Description: description,
		Location:    sql.NullString{String: location, Valid: location != ""},
		EventDate:   eventDate,
		Status:      status,
	})
	if err != nil {
		h.logger.Error("failed to update event", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "event", id, title, "Updated Event '%s'", title)

	return c.Redirect(http.StatusSeeOther, "/admin/events")
}

// Delete handles DELETE /admin/events/:id
// Removes an event from the database
func (h *EventsHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := h.queries.DeleteEvent(ctx, id); err != nil {
		h.logger.Error("failed to delete event", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "deleted", "event", id, "", "Deleted Event #%d", id)

	// Return HTTP 200 with empty body (HTMX will remove the row)
	return c.NoContent(http.StatusOK)
}
```

### Step 5: Register Routes in main.go

Add to `cmd/server/main.go` in the admin routes section (after other CRUD handlers):

```go
// Event Management
eventsHandler := adminHandlers.NewEventsHandler(queries, logger)
adminGroup.GET("/events", eventsHandler.List)              // List all events
adminGroup.GET("/events/new", eventsHandler.New)           // Show creation form
adminGroup.POST("/events", eventsHandler.Create)           // Process new event
adminGroup.GET("/events/:id/edit", eventsHandler.Edit)     // Show edit form
adminGroup.POST("/events/:id", eventsHandler.Update)       // Process updates
adminGroup.DELETE("/events/:id", eventsHandler.Delete)     // Delete event (HTMX)
```

**Route pattern:**
- `GET /resource` — List page
- `GET /resource/new` — Create form
- `POST /resource` — Create action
- `GET /resource/:id/edit` — Edit form
- `POST /resource/:id` — Update action
- `DELETE /resource/:id` — Delete action

### Step 6: Register Templates in template.go

Add to `internal/templates/template.go` in the `loadTemplates()` function (after other master table pages):

```go
// Event admin pages
eventAdminPages := []string{
	"events_list", "events_form",
}
for _, page := range eventAdminPages {
	r.templates["admin/pages/"+page+".html"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		filepath.Join(r.basePath, "admin/layouts/base.html"),
		filepath.Join(r.basePath, "admin/pages/"+page+".html"),
		filepath.Join(r.basePath, "partials/admin-sidebar.html"),
	))
}
```

**Template registration patterns:**
- Full page templates include: `base.html`, the page itself, and `admin-sidebar.html`
- HTMX partials only include the partial file itself (no layout)

### Step 7: Create Template Files

**`templates/admin/pages/events_list.html`:**
```html
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Bluejay CMS</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body class="font-mono bg-white">
    <div class="flex">
        {{template "admin-sidebar" .}}

        <main class="flex-1 p-8">
            <div class="max-w-6xl mx-auto">
                <div class="flex justify-between items-center mb-8">
                    <h1 class="text-3xl uppercase border-b-4 border-black pb-2">{{.Title}}</h1>
                    <a href="/admin/events/new"
                       class="px-6 py-3 bg-black text-white uppercase border-2 border-black hover:bg-white hover:text-black transition-all">
                        Add New Event
                    </a>
                </div>

                {{if .Items}}
                <table class="w-full border-2 border-black">
                    <thead>
                        <tr class="bg-black text-white">
                            <th class="p-4 text-left uppercase">Title</th>
                            <th class="p-4 text-left uppercase">Date</th>
                            <th class="p-4 text-left uppercase">Status</th>
                            <th class="p-4 text-right uppercase">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Items}}
                        <tr class="border-b-2 border-black" id="event-row-{{.ID}}">
                            <td class="p-4">{{.Title}}</td>
                            <td class="p-4">{{formatDate .EventDate "Jan 2, 2006"}}</td>
                            <td class="p-4">
                                <span class="px-2 py-1 border border-black uppercase text-xs">{{.Status}}</span>
                            </td>
                            <td class="p-4 text-right">
                                <a href="/admin/events/{{.ID}}/edit"
                                   class="text-blue-600 hover:underline mr-4">Edit</a>
                                <button hx-delete="/admin/events/{{.ID}}"
                                        hx-confirm="Delete this event?"
                                        hx-target="#event-row-{{.ID}}"
                                        hx-swap="outerHTML swap:0.5s"
                                        class="text-red-600 hover:underline">
                                    Delete
                                </button>
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
                {{else}}
                <div class="border-4 border-black p-8 text-center">
                    <p class="text-xl uppercase mb-4">No events found</p>
                    <a href="/admin/events/new" class="underline">Create your first event</a>
                </div>
                {{end}}
            </div>
        </main>
    </div>
</body>
</html>
{{end}}
```

**`templates/admin/pages/events_form.html`:**
```html
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Bluejay CMS</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="font-mono bg-white">
    <div class="flex">
        {{template "admin-sidebar" .}}

        <main class="flex-1 p-8">
            <div class="max-w-4xl mx-auto">
                <h1 class="text-3xl uppercase border-b-4 border-black pb-2 mb-8">{{.Title}}</h1>

                <form method="POST" action="{{.FormAction}}" class="space-y-6">
                    <div>
                        <label class="block uppercase mb-2 font-bold">Title *</label>
                        <input type="text" name="title" required
                               value="{{if .Item}}{{.Item.Title}}{{end}}"
                               class="w-full p-3 border-2 border-black focus:outline-none focus:border-4">
                    </div>

                    <div>
                        <label class="block uppercase mb-2 font-bold">Description *</label>
                        <textarea name="description" rows="5" required
                                  class="w-full p-3 border-2 border-black focus:outline-none focus:border-4">{{if .Item}}{{.Item.Description}}{{end}}</textarea>
                    </div>

                    <div>
                        <label class="block uppercase mb-2 font-bold">Location</label>
                        <input type="text" name="location"
                               value="{{if .Item}}{{.Item.Location.String}}{{end}}"
                               class="w-full p-3 border-2 border-black focus:outline-none focus:border-4">
                    </div>

                    <div>
                        <label class="block uppercase mb-2 font-bold">Event Date *</label>
                        <input type="datetime-local" name="event_date" required
                               value="{{if .Item}}{{formatDate .Item.EventDate "2006-01-02T15:04"}}{{end}}"
                               class="w-full p-3 border-2 border-black focus:outline-none focus:border-4">
                    </div>

                    <div>
                        <label class="block uppercase mb-2 font-bold">Status *</label>
                        <select name="status" required
                                class="w-full p-3 border-2 border-black focus:outline-none focus:border-4">
                            <option value="draft" {{if .Item}}{{if eq .Item.Status "draft"}}selected{{end}}{{end}}>Draft</option>
                            <option value="published" {{if .Item}}{{if eq .Item.Status "published"}}selected{{end}}{{end}}>Published</option>
                        </select>
                    </div>

                    <div class="flex gap-4">
                        <button type="submit"
                                class="px-8 py-4 bg-black text-white uppercase border-2 border-black hover:bg-white hover:text-black transition-all">
                            Save Event
                        </button>
                        <a href="/admin/events"
                           class="px-8 py-4 border-2 border-black uppercase hover:bg-black hover:text-white transition-all inline-block">
                            Cancel
                        </a>
                    </div>
                </form>
            </div>
        </main>
    </div>
</body>
</html>
{{end}}
```

### Step 8: Add Sidebar Navigation Entry

Edit `templates/partials/admin-sidebar.html` and add:

```html
<a href="/admin/events"
   class="block px-6 py-3 hover:bg-black hover:text-white border-b border-gray-200 uppercase">
    Events
</a>
```

### Step 9: Test and Verify

1. **Restart the server** (migrations run automatically):
   ```bash
   make run
   ```

2. **Verify database schema**:
   ```bash
   sqlite3 bluejay.db ".schema events"
   ```

3. **Test the feature**:
   - Navigate to `/admin/events`
   - Create a new event
   - Edit the event
   - Delete the event

4. **Check logs** for any errors:
   ```bash
   # Server logs appear in stdout
   # Look for structured JSON logs with error level
   ```

---

## Handler Conventions

### Handler Structure Pattern

All handlers follow this consistent pattern:

```go
package admin

import (
	"log/slog"
	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

// Handler manages HTTP requests for a specific resource
type ResourceHandler struct {
	queries *sqlc.Queries   // Database queries (required)
	logger  *slog.Logger    // Structured logger (required)
	// Optional dependencies:
	uploadSvc *services.UploadService
	cache     *services.Cache
}

// Constructor initializes handler with dependencies
func NewResourceHandler(queries *sqlc.Queries, logger *slog.Logger) *ResourceHandler {
	return &ResourceHandler{
		queries: queries,
		logger:  logger,
	}
}

// Handler methods follow RESTful naming
func (h *ResourceHandler) List(c echo.Context) error { /* ... */ }
func (h *ResourceHandler) New(c echo.Context) error { /* ... */ }
func (h *ResourceHandler) Create(c echo.Context) error { /* ... */ }
func (h *ResourceHandler) Edit(c echo.Context) error { /* ... */ }
func (h *ResourceHandler) Update(c echo.Context) error { /* ... */ }
func (h *ResourceHandler) Delete(c echo.Context) error { /* ... */ }
```

### Error Handling

**Log and return HTTP errors:**
```go
if err != nil {
	h.logger.Error("failed to list products", "error", err)
	return echo.NewHTTPError(http.StatusInternalServerError)
}
```

**Return 404 for missing resources:**
```go
product, err := h.queries.GetProduct(ctx, id)
if err != nil {
	return echo.NewHTTPError(http.StatusNotFound, "Product not found")
}
```

**Return 400 for validation errors:**
```go
if c.FormValue("email") == "" {
	return echo.NewHTTPError(http.StatusBadRequest, "Email is required")
}
```

### Flash Messages and Redirects

**Redirect after successful create/update:**
```go
// Create action
logActivity(c, "created", "product", 0, name, "Created Product '%s'", name)
return c.Redirect(http.StatusSeeOther, "/admin/products")
```

**Redirect with error message:**
```go
return c.Redirect(http.StatusSeeOther, "/admin/login?error=invalid_credentials")
```

### Context and Database Access

Always use `c.Request().Context()` for database operations:

```go
ctx := c.Request().Context()
products, err := h.queries.ListProducts(ctx)
```

This ensures proper cancellation and timeout handling.

### Form Value Parsing

**Extract form values:**
```go
title := c.FormValue("title")
status := c.FormValue("status")

// Parse integers
categoryID, _ := strconv.ParseInt(c.FormValue("category_id"), 10, 64)

// Parse booleans
isFeatured := c.FormValue("is_featured") == "1"
```

**Handle nullable fields:**
```go
location := c.FormValue("location")
locationNull := sql.NullString{
	String: location,
	Valid:  location != "",
}
```

### File Uploads

```go
// Check if file was uploaded
file, err := c.FormFile("primary_image")
if err == nil {
	path, err := h.uploadSvc.UploadProductImage(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to upload image")
	}
	imagePath = sql.NullString{String: path, Valid: true}
}
```

---

## Template Conventions

### Template Structure

Templates follow a layout inheritance pattern:

**Base Layout** (`admin/layouts/base.html`):
```html
{{define "base"}}
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - Bluejay CMS</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    {{block "content" .}}{{end}}
</body>
</html>
{{end}}
```

**Page Template** (`admin/pages/products_list.html`):
```html
{{define "base"}}
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <!-- ... -->
</head>
<body>
    <div class="flex">
        {{template "admin-sidebar" .}}
        <main class="flex-1">
            <!-- Page content here -->
        </main>
    </div>
</body>
</html>
{{end}}
```

### Template Data Maps

Pass data to templates as `map[string]interface{}`:

```go
return c.Render(http.StatusOK, "admin/pages/products_list.html", map[string]interface{}{
	"Title":      "Manage Products",
	"Products":   products,
	"Categories": categories,
	"Page":       page,
	"TotalPages": totalPages,
	"HasFilters": hasFilters,
})
```

### Accessing Data in Templates

```html
<h1>{{.Title}}</h1>

{{if .Products}}
    {{range .Products}}
        <div>{{.Name}} - {{.SKU}}</div>
    {{end}}
{{else}}
    <p>No products found</p>
{{end}}

{{if .HasFilters}}
    <button>Clear Filters</button>
{{end}}
```

### Template Functions

Built-in custom functions available in all templates:

```html
<!-- Format dates -->
{{formatDate .CreatedAt "Jan 2, 2006"}}

<!-- Render HTML without escaping -->
{{safeHTML .Description}}

<!-- Truncate strings -->
{{truncate .LongText 100}}

<!-- Math operations -->
{{add .Page 1}}
{{sub .Total 5}}

<!-- String manipulation -->
{{upper .Name}}

<!-- File size formatting -->
{{formatFileSize .FileSize}}
```

### Conditional Rendering

```html
<!-- Check if value exists -->
{{if .Item}}
    <p>Editing: {{.Item.Name}}</p>
{{else}}
    <p>Creating new item</p>
{{end}}

<!-- Check equality -->
{{if eq .Status "published"}}
    <span class="badge-success">Published</span>
{{end}}

<!-- Check for nullable fields -->
{{if .Item.Location.Valid}}
    <p>Location: {{.Item.Location.String}}</p>
{{end}}
```

### Form Patterns

**Create/Edit form pattern:**
```html
<form method="POST" action="{{.FormAction}}">
    <input type="text" name="title"
           value="{{if .Item}}{{.Item.Title}}{{end}}" required>

    <select name="status">
        <option value="draft"
                {{if .Item}}{{if eq .Item.Status "draft"}}selected{{end}}{{end}}>
            Draft
        </option>
        <option value="published"
                {{if .Item}}{{if eq .Item.Status "published"}}selected{{end}}{{end}}>
            Published
        </option>
    </select>

    <button type="submit">Save</button>
</form>
```

---

## HTMX Development Patterns

HTMX enables dynamic page updates without writing JavaScript. The CMS uses HTMX extensively for inline editing, delete confirmations, and dynamic content loading.

### Basic Delete with Confirmation

```html
<button hx-delete="/admin/products/{{.ID}}"
        hx-confirm="Delete this product?"
        hx-target="#product-row-{{.ID}}"
        hx-swap="outerHTML swap:0.5s"
        class="text-red-600">
    Delete
</button>
```

**Handler returns HTTP 200 with empty body:**
```go
func (h *ProductsHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.queries.DeleteProduct(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	// Return 200 OK with no content - HTMX will remove the row
	return c.NoContent(http.StatusOK)
}
```

### Dynamic Tab Content

**Tabs in template:**
```html
<div class="tabs">
    <button hx-get="/admin/products/{{.ID}}/specs"
            hx-target="#tab-content"
            class="tab-active">Specs</button>
    <button hx-get="/admin/products/{{.ID}}/features"
            hx-target="#tab-content">Features</button>
</div>

<div id="tab-content">
    <!-- Content loads here -->
</div>
```

**Handler returns partial:**
```go
func (h *ProductDetailsHandler) ListSpecs(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	specs, _ := h.queries.ListProductSpecs(c.Request().Context(), id)

	// Returns partial template (no layout)
	return c.Render(http.StatusOK, "admin/partials/product_specs.html", map[string]interface{}{
		"Specs":     specs,
		"ProductID": id,
	})
}
```

### Inline Forms with HTMX

**Add feature inline:**
```html
<form hx-post="/admin/products/{{.ProductID}}/features"
      hx-target="#features-list"
      hx-swap="beforeend">
    <input type="text" name="feature_text" required>
    <button type="submit">Add Feature</button>
</form>

<div id="features-list">
    {{range .Features}}
        <div>{{.FeatureText}}</div>
    {{end}}
</div>
```

**Handler returns new item HTML:**
```go
func (h *ProductDetailsHandler) AddFeature(c echo.Context) error {
	ctx := c.Request().Context()
	productID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	feature, err := h.queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID:    productID,
		FeatureText:  c.FormValue("feature_text"),
		DisplayOrder: 0,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Return HTML fragment for the new feature
	return c.HTML(http.StatusOK, fmt.Sprintf(
		`<div class="feature-item">%s</div>`,
		feature.FeatureText,
	))
}
```

### Typeahead Search

**Search input:**
```html
<input type="text"
       name="search"
       hx-get="/admin/blog/products/search"
       hx-trigger="keyup changed delay:300ms"
       hx-target="#search-results"
       placeholder="Search products...">

<div id="search-results">
    <!-- Results appear here -->
</div>
```

**Handler returns results partial:**
```go
func (h *BlogPostsHandler) SearchProducts(c echo.Context) error {
	q := strings.TrimSpace(c.QueryParam("_product_search"))
	if q == "" {
		return c.Render(http.StatusOK, "admin/partials/product_suggestions.html", map[string]interface{}{
			"Products": nil,
		})
	}

	products, _ := h.queries.SearchPublishedProducts(c.Request().Context(), "%"+q+"%")
	return c.Render(http.StatusOK, "admin/partials/product_suggestions.html", map[string]interface{}{
		"Products": products,
		"Query":    q,
	})
}
```

### HTMX Response Headers

**Trigger client-side events:**
```go
// Trigger a success notification
c.Response().Header().Set("HX-Trigger", "showSuccess")

// Redirect after HTMX request
c.Response().Header().Set("HX-Redirect", "/admin/dashboard")

// Refresh page
c.Response().Header().Set("HX-Refresh", "true")
```

---

## Service Layer Patterns

### When to Create a Service

Create a service when you have:

1. **Complex business logic** spanning multiple queries
2. **Reusable operations** needed by multiple handlers
3. **Data aggregation** from multiple tables
4. **External integrations** (file uploads, email, APIs)

### Service Example: Product Detail Aggregation

**Service** (`internal/services/product.go`):
```go
package services

import (
	"context"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type ProductService struct {
	queries *sqlc.Queries
}

func NewProductService(queries *sqlc.Queries) *ProductService {
	return &ProductService{queries: queries}
}

// ProductDetail aggregates product data from multiple tables
type ProductDetail struct {
	Product        sqlc.Product
	Category       sqlc.ProductCategory
	Specs          []sqlc.ProductSpec
	Images         []sqlc.ProductImage
	Features       []sqlc.ProductFeature
	Certifications []sqlc.ProductCertification
	Downloads      []sqlc.ProductDownload
}

// GetProductDetail fetches a product with all related data
func (s *ProductService) GetProductDetail(ctx context.Context, slug string) (*ProductDetail, error) {
	product, err := s.queries.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	category, err := s.queries.GetProductCategory(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	// Fetch all related data (errors ignored, return empty slices)
	specs, _ := s.queries.ListProductSpecs(ctx, product.ID)
	images, _ := s.queries.ListProductImages(ctx, product.ID)
	features, _ := s.queries.ListProductFeatures(ctx, product.ID)
	certifications, _ := s.queries.ListProductCertifications(ctx, product.ID)
	downloads, _ := s.queries.ListProductDownloads(ctx, product.ID)

	return &ProductDetail{
		Product:        product,
		Category:       category,
		Specs:          specs,
		Images:         images,
		Features:       features,
		Certifications: certifications,
		Downloads:      downloads,
	}, nil
}
```

**Usage in handler:**
```go
detail, err := h.productSvc.GetProductDetail(c.Request().Context(), slug)
if err != nil {
	return echo.NewHTTPError(http.StatusNotFound)
}
```

### Service Pattern: File Upload

**Upload Service** (`internal/services/upload.go`):
```go
type UploadService struct {
	uploadDir string
}

func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{uploadDir: uploadDir}
}

func (s *UploadService) UploadProductImage(file *multipart.FileHeader) (string, error) {
	// Extract and validate file extension (jpg, jpeg, png, webp only)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("file too large (max 5MB)")
	}

	// Generate unique filename with timestamp
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)

	// Ensure products subdirectory exists
	productDir := filepath.Join(s.uploadDir, "products")
	os.MkdirAll(productDir, 0755)

	filePath := filepath.Join(productDir, filename)

	// Save file to disk
	src, _ := file.Open()
	defer src.Close()

	dst, _ := os.Create(filePath)
	defer dst.Close()

	io.Copy(dst, src)

	// Return public URL path
	return "/uploads/products/" + filename, nil
}
```

### Service Pattern: Caching

**Cache Service** (`internal/services/cache.go`):
```go
type Cache struct {
	store *sync.Map
}

func NewCache() *Cache {
	return &Cache{store: &sync.Map{}}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.store.Load(key)
}

func (c *Cache) Set(key string, value interface{}, ttlSeconds int) {
	c.store.Store(key, value)

	// Auto-expire after TTL
	if ttlSeconds > 0 {
		go func() {
			time.Sleep(time.Duration(ttlSeconds) * time.Second)
			c.store.Delete(key)
		}()
	}
}

func (c *Cache) DeleteByPrefix(prefix string) {
	c.store.Range(func(key, value interface{}) bool {
		if strings.HasPrefix(key.(string), prefix) {
			c.store.Delete(key)
		}
		return true
	})
}
```

**Usage:**
```go
// Check cache first
cacheKey := "page:products:detectors"
if cached, ok := h.cache.Get(cacheKey); ok {
	return c.HTML(http.StatusOK, cached.(string))
}

// Render and cache
html := renderTemplate(...)
h.cache.Set(cacheKey, html, 600) // 10 minutes
return c.HTML(http.StatusOK, html)

// Invalidate on mutations
h.cache.DeleteByPrefix("page:products")
```

---

## Testing Approach

The CMS uses three levels of testing:

### 1. Unit Tests for Services

Test business logic in isolation using a test database.

**Example** (`internal/services/product_test.go`):
```go
package services_test

import (
	"context"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestGetProductDetail(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create test data
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Detectors", Slug: "detectors", Description: "desc", Icon: "icon", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	prod, err := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DET-001", Slug: "detector-one", Name: "Detector One",
		Description: "A detector", CategoryID: cat.ID, Status: "published",
	})
	if err != nil {
		t.Fatalf("CreateProduct: %v", err)
	}

	// Test service method
	svc := services.NewProductService(queries)
	detail, err := svc.GetProductDetail(ctx, "detector-one")

	if err != nil {
		t.Fatalf("GetProductDetail: %v", err)
	}

	if detail.Product.Name != "Detector One" {
		t.Errorf("expected 'Detector One', got %q", detail.Product.Name)
	}

	if detail.Category.Name != "Detectors" {
		t.Errorf("expected 'Detectors', got %q", detail.Category.Name)
	}
}
```

**Run service tests:**
```bash
go test -v ./internal/services/...
```

### 2. Handler Tests

Test HTTP handlers in isolation with mock requests.

**Example** (`internal/handlers/admin/handlers_test.go`):
```go
package admin_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/handlers/admin"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

func TestProductCategoriesHandler_Delete(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	e := echo.New()
	h := admin.NewProductCategoriesHandler(queries, logger)

	ctx := context.Background()

	// Create a category directly using queries
	cat, err := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "Test", Slug: "test", Description: "d", Icon: "i", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Test delete handler
	req := httptest.NewRequest(http.MethodDelete, "/admin/product-categories/"+strconv.FormatInt(cat.ID, 10), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(cat.ID, 10))

	if err := h.Delete(c); err != nil {
		t.Fatalf("Delete handler: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
```

**Run handler tests:**
```bash
go test -v ./internal/handlers/...
```

### 3. End-to-End Tests

Test full request/response cycle through the entire application stack.

**Example** (`internal/e2e/e2e_test.go`):
```go
package e2e_test

import (
	"net/http/httptest"
	"testing"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestPublicHomePage_E2E(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	// Test homepage renders successfully
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHealthCheck_E2E(t *testing.T) {
	e, _, cleanup := setupApp(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	// Verify JSON response structure
	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", response["status"])
	}
}
```

**Run e2e tests:**
```bash
go test -v ./internal/e2e/...
```

### Test Database Setup

The `testutil` package provides helper functions:

```go
// SetupTestDB creates a temporary test database with migrations
db, queries, cleanup := testutil.SetupTestDB(t)
defer cleanup()  // Cleans up database file after test
```

---

## Common Gotchas

### 1. Must Run `sqlc generate` After SQL Changes

**Problem:** You modify a SQL query file but Go compilation fails with missing methods.

**Solution:** Always run after changing `.sql` files:
```bash
make sqlc
# or
sqlc generate
```

**Why:** sqlc generates Go code from SQL files. The generated code is what your handlers import and use.

### 2. Template Must Be Registered in template.go

**Problem:** Handler renders a template but you get "template not found" error.

**Solution:** Add template registration in `internal/templates/template.go`:
```go
r.templates["admin/pages/events_list.html"] = template.Must(...)
```

**Why:** Templates are pre-compiled at startup. Unregistered templates cannot be rendered.

### 3. SQLite Single-Writer Limitation

**Problem:** Database is locked during writes.

**Solution:**
- Keep transactions short
- Use `defer` to ensure connections are closed
- Consider retry logic for concurrent writes

```go
// Good: Quick transaction
_, err := h.queries.CreateProduct(ctx, params)

// Bad: Long-running transaction
tx, _ := db.Begin()
// ... lots of operations ...
tx.Commit()  // Blocks other writes
```

**Why:** SQLite only allows one writer at a time. Long transactions block the entire database.

### 4. Cache Invalidation After Mutations

**Problem:** Public pages show stale data after admin updates.

**Solution:** Invalidate cache in Create/Update/Delete handlers:
```go
h.cache.DeleteByPrefix("page:products")
```

**Why:** Public pages are cached for performance. Cache must be cleared when data changes.

### 5. HTMX Endpoints Must Return Partials

**Problem:** HTMX request returns full page HTML instead of fragment.

**Solution:** HTMX handlers should render partial templates (no layout):
```go
// Good: Partial template
return c.Render(http.StatusOK, "admin/partials/product_specs.html", data)

// Bad: Full page template
return c.Render(http.StatusOK, "admin/pages/products_list.html", data)
```

**Why:** HTMX swaps specific DOM elements. Returning full pages breaks the page structure.

### 6. Nullable Database Fields

**Problem:** Trying to access nullable field causes panic or wrong value.

**Solution:** Use `sql.Null*` types and check `.Valid`:
```go
// Writing
location := c.FormValue("location")
locationNull := sql.NullString{String: location, Valid: location != ""}

// Reading
if product.Location.Valid {
	fmt.Println(product.Location.String)
}
```

**Why:** SQLite allows NULL values. Go requires explicit handling of nullable types.

### 7. Form Values vs Query Params

**Problem:** Form data not found in handler.

**Solution:** Use correct method:
```go
// POST form data
title := c.FormValue("title")

// GET query parameters
search := c.QueryParam("search")

// URL path parameters
id := c.Param("id")
```

### 8. Redirect Status Codes

**Problem:** Browser doesn't follow redirect after POST.

**Solution:** Use `http.StatusSeeOther` (303) for POST-redirect-GET:
```go
// Good: POST → Redirect → GET
return c.Redirect(http.StatusSeeOther, "/admin/products")

// Bad: 301/302 may resubmit POST
return c.Redirect(http.StatusMovedPermanently, "/admin/products")
```

---

## Code Style and Conventions

### Handler Naming

- **Handler Struct:** `{Resource}Handler` (e.g., `ProductsHandler`, `EventsHandler`)
- **Constructor:** `New{Resource}Handler`
- **Methods:** RESTful verbs (`List`, `New`, `Create`, `Edit`, `Update`, `Delete`)

### Route Naming

Follow RESTful URL patterns:

```
GET    /admin/products             → List
GET    /admin/products/new         → New (form)
POST   /admin/products             → Create
GET    /admin/products/:id/edit    → Edit (form)
POST   /admin/products/:id         → Update
DELETE /admin/products/:id         → Delete
```

### Template Naming

- **Pages:** `{resource}_{action}.html` (e.g., `products_list.html`, `products_form.html`)
- **Partials:** `{resource}_{component}.html` (e.g., `product_specs.html`, `tag_suggestions.html`)
- **Location:**
  - Admin pages: `templates/admin/pages/`
  - Public pages: `templates/public/pages/`
  - Partials: `templates/admin/partials/` or `templates/public/partials/`

### File Organization

```
internal/
  handlers/
    admin/
      products.go          # Main CRUD
      product_details.go   # Sub-entities (specs, images, etc.)
    public/
      products.go          # Public-facing pages
  services/
    product.go             # Business logic
    product_test.go        # Service tests
  middleware/
    auth.go                # Authentication
    logging.go             # Request logging
```

### Import Organization

Group imports logically:

```go
import (
	// Standard library
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	// Third-party
	"github.com/labstack/echo/v4"

	// Internal
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)
```

### Variable Naming

- **Short names** for common types: `ctx`, `err`, `id`, `req`, `rec`
- **Descriptive names** for business logic: `categoryID`, `publishedAt`, `productDetail`
- **Acronyms uppercase:** `ID`, `URL`, `HTML` (not `Id`, `Url`, `Html`)

### Error Messages

```go
// Good: Specific, actionable
h.logger.Error("failed to create product", "error", err, "sku", sku)

// Bad: Generic, no context
h.logger.Error("error", "error", err)
```

### Comments

Document exported types and complex logic:

```go
// ProductsHandler manages all HTTP handlers for product CRUD operations.
// It handles listing with filters, creating, editing, updating, and deleting products.
type ProductsHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

// List handles GET /admin/products
// Renders the main product list page with filtering, searching, and pagination.
func (h *ProductsHandler) List(c echo.Context) error {
	// ...
}
```

---

## Debugging Tips

### Structured Logging

The application uses structured JSON logging:

```go
// Log with context
h.logger.Info("product created",
	"product_id", product.ID,
	"sku", product.SKU,
	"user_id", userID,
)

h.logger.Error("database query failed",
	"error", err,
	"query", "ListProducts",
	"filters", map[string]interface{}{
		"status": status,
		"category": categoryID,
	},
)
```

**Viewing logs:**
```bash
# All logs go to stdout as JSON
make run | jq .

# Filter errors only
make run | jq 'select(.level == "ERROR")'

# Filter by specific field
make run | jq 'select(.query == "ListProducts")'
```

### Health Endpoint

Check application health:

```bash
curl http://localhost:28090/health
# {"status":"ok","time":"2024-02-10T18:30:00Z"}
```

### Database Inspection

```bash
# Open database in SQLite CLI
sqlite3 bluejay.db

# List all tables
.tables

# Show schema for a table
.schema products

# Query data
SELECT * FROM products LIMIT 10;

# Check migrations
SELECT * FROM schema_migrations;

# Exit
.exit
```

### Template Debugging

Add debug output in templates:

```html
<!-- Show all data passed to template -->
<pre>{{printf "%+v" .}}</pre>

<!-- Check if variable exists -->
{{if .Products}}
	Products: {{len .Products}}
{{else}}
	No products
{{end}}
```

### HTMX Debugging

Add HTMX logging in browser console:

```html
<script>
document.body.addEventListener('htmx:configRequest', (e) => {
  console.log('HTMX Request:', e.detail);
});

document.body.addEventListener('htmx:afterSwap', (e) => {
  console.log('HTMX Swap:', e.detail);
});
</script>
```

### Request Debugging

Log incoming requests:

```go
func (h *ProductsHandler) Create(c echo.Context) error {
	// Log all form values
	h.logger.Info("form submitted",
		"form", c.Request().Form,
		"url", c.Request().URL.String(),
	)

	// ...
}
```

### Common Error Patterns

**"template not found"**
- Check template is registered in `template.go`
- Verify file exists at correct path
- Check template name matches registration

**"database is locked"**
- Close any open SQLite connections
- Check for long-running transactions
- Restart server to release locks

**"FOREIGN KEY constraint failed"**
- Check foreign key relationships
- Delete child records before parent
- Use CASCADE on foreign keys if appropriate

**"bind: address already in use"**
- Server already running on port 28090
- Kill existing process: `lsof -ti:28090 | xargs kill -9`

### VSCode Debugging

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server/main.go",
      "env": {},
      "args": []
    }
  ]
}
```

Set breakpoints and press F5 to debug.

---

## Quick Reference

### Common Commands

```bash
# Development
make run              # Start server
make dev              # Start with hot-reload
make build            # Build binary
make test             # Run all tests
make clean            # Clean artifacts

# Database
make sqlc             # Generate Go code from SQL
sqlite3 bluejay.db    # Open database CLI

# Testing
go test -v ./...                           # All tests
go test -v ./internal/handlers/admin/...   # Handler tests
go test -v ./internal/services/...         # Service tests
go test -run TestSpecificTest              # Single test
go test -cover ./...                       # With coverage
```

### File Patterns

```
db/migrations/XXX_name.up.sql      # Migration up
db/migrations/XXX_name.down.sql    # Migration down
db/queries/resource.sql            # SQL queries
internal/handlers/admin/resource.go # Admin handler
internal/handlers/public/resource.go # Public handler
internal/services/resource.go      # Business logic
templates/admin/pages/resource_list.html  # List page
templates/admin/pages/resource_form.html  # Form page
templates/admin/partials/resource_*.html  # HTMX partials
```

### Handler Template

```go
package admin

import (
	"log/slog"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type ResourceHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewResourceHandler(q *sqlc.Queries, l *slog.Logger) *ResourceHandler {
	return &ResourceHandler{queries: q, logger: l}
}

func (h *ResourceHandler) List(c echo.Context) error {
	items, err := h.queries.ListResources(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to list resources", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "admin/pages/resources_list.html", map[string]interface{}{
		"Title": "Resources",
		"Items": items,
	})
}
```

---

This guide covers the essential patterns and workflows for developing features in Bluejay CMS. For specific implementation details, refer to the existing codebase examples in `internal/handlers/admin/products.go`, `internal/handlers/admin/blog_posts.go`, and related files.
