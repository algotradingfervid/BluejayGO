# Bluejay CMS Architecture

## System Overview

Bluejay CMS is a brutalist-designed content management system built with Go, providing both a public-facing website and an administrative backend. The system uses server-side rendering with HTMX for dynamic interactions.

```
                           ┌─────────────────┐
                           │   Caddy/Nginx   │
                           │   (Reverse      │
                           │    Proxy)       │
                           └────────┬────────┘
                                    │
                           ┌────────▼────────┐
                           │  Echo Web       │
                           │  Framework      │
                           │  :28090         │
                           └────────┬────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
          ┌─────────▼─────────┐    │    ┌─────────▼─────────┐
          │  Middleware Chain │    │    │  Static Files     │
          │  - Recovery       │    │    │  /public/*        │
          │  - Logging        │    │    │  /uploads/*       │
          │  - Gzip           │    │    └───────────────────┘
          │  - Security       │    │
          │  - Session        │    │
          │  - Auth (admin)   │    │
          │  - Settings       │    │
          │  - RateLimit      │    │
          └─────────┬─────────┘    │
                    │               │
          ┌─────────▼─────────┐    │
          │   Route Groups    │    │
          │  - Public         │    │
          │  - Admin (Auth)   │    │
          └─────────┬─────────┘    │
                    │               │
          ┌─────────▼─────────┐    │
          │    Handlers       │    │
          │  - Admin CRUD     │◄───┘
          │  - Public Pages   │
          └─────────┬─────────┘
                    │
          ┌─────────▼─────────┐
          │    Services       │
          │  - ProductService │
          │  - UploadService  │
          │  - Cache          │
          │  - ActivityLog    │
          └─────────┬─────────┘
                    │
          ┌─────────▼─────────┐
          │  sqlc Queries     │
          │  (Type-safe SQL)  │
          └─────────┬─────────┘
                    │
          ┌─────────▼─────────┐
          │  SQLite + WAL     │
          │  bluejay.db       │
          └─────────┬─────────┘
                    │
          ┌─────────▼─────────┐
          │   Litestream      │
          │   Replication     │
          │   → S3 Backup     │
          └───────────────────┘
```

## Request Lifecycle

### 1. Incoming HTTP Request
```
Client → Caddy/Nginx → Echo Server (port 28090)
```

### 2. Middleware Chain Execution (in order)
```go
Recovery()           // Catches panics, logs stack traces
  ↓
Logging()           // Logs request method, path, duration, status
  ↓
Gzip()              // Compresses responses
  ↓
SecurityHeaders()   // Sets X-Frame-Options, CSP, etc.
  ↓
SessionMiddleware() // Loads/creates session from cookie
  ↓
[Public Routes] → SettingsLoader()  // Loads site settings, footer data
  ↓
[Admin Routes] → RequireAuth()      // Checks session.UserID, redirects if not authenticated
  ↓
[Specific Routes] → RateLimiter()   // IP-based rate limiting (applied per-route, e.g., contact form: 5 requests/hour)
```

### 3. Route Handler Execution
```go
Handler receives echo.Context
  ↓
Extract parameters (path params, query strings, form values)
  ↓
Call Service Layer (if needed)
  ↓
Execute sqlc Queries
  ↓
Build template data map
  ↓
Return c.Render() or c.JSON()
```

### 4. Template Rendering
```go
Echo.Renderer.Render()
  ↓
Lookup template by name in pre-compiled map
  ↓
Execute template with data
  ↓
  ├─ Full Page: Renders layout + page content
  └─ HTMX Partial: Renders fragment only
  ↓
Write HTML to response
```

### 5. Response
```
HTML/JSON → Middleware (reverse order) → Client
```

## Directory Structure

```
bluejay-cms/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── database/
│   │   ├── sqlite.go            # DB connection setup (WAL mode, pragmas)
│   │   └── migrate.go           # Migration runner (golang-migrate)
│   │
│   ├── testutil/                # Test utilities and helpers
│   ├── e2e/                     # End-to-end test suite
│   │
│   ├── handlers/
│   │   ├── admin/               # Admin panel handlers (CRUD operations)
│   │   │   ├── auth.go          # Login/logout
│   │   │   ├── dashboard.go     # Dashboard statistics
│   │   │   ├── products.go      # Product management
│   │   │   ├── blog_posts.go    # Blog post management
│   │   │   ├── solutions.go     # Solution management
│   │   │   ├── case_studies.go  # Case study management
│   │   │   ├── whitepapers.go   # Whitepaper management
│   │   │   ├── contact.go       # Contact submission management
│   │   │   ├── media.go         # Media library
│   │   │   ├── navigation.go    # Navigation menu editor
│   │   │   ├── settings.go      # Site settings
│   │   │   ├── header.go        # Header configuration
│   │   │   ├── footer.go        # Footer configuration
│   │   │   └── activity.go      # Activity log viewer
│   │   │
│   │   └── public/              # Public-facing handlers (read-only)
│   │       ├── home.go          # Homepage
│   │       ├── products.go      # Product listing, detail, search
│   │       ├── solutions.go     # Solution pages
│   │       ├── blog.go          # Blog listing, posts
│   │       ├── case_studies.go  # Case study pages
│   │       ├── whitepapers.go   # Whitepaper pages with download
│   │       ├── contact.go       # Contact form
│   │       ├── about.go         # About page
│   │       ├── partners.go      # Partners page
│   │       ├── search.go        # Global search
│   │       └── sitemap.go       # SEO sitemap/robots.txt
│   │
│   ├── middleware/
│   │   ├── recovery.go          # Panic recovery with stack traces
│   │   ├── logging.go           # Request/response logging
│   │   ├── security.go          # Security headers (CSP, X-Frame-Options)
│   │   ├── session.go           # Session management (gorilla/sessions)
│   │   ├── auth.go              # Authentication guard
│   │   ├── settings.go          # Settings loader for public pages
│   │   ├── ratelimit.go         # IP-based rate limiting
│   │   ├── csrf.go              # CSRF protection (available but not currently used)
│   │   ├── cache.go             # Cache control middleware
│   │   └── middleware_test.go   # Middleware unit tests
│   │
│   ├── services/
│   │   ├── product.go           # ProductService (aggregate product data)
│   │   ├── upload.go            # UploadService (file uploads)
│   │   ├── cache.go             # In-memory TTL cache
│   │   ├── activity_log.go      # ActivityLogService (audit logging)
│   │   ├── product_test.go      # ProductService unit tests
│   │   ├── upload_test.go       # UploadService unit tests
│   │   └── cache_test.go        # Cache unit tests
│   │
│   ├── templates/
│   │   └── template.go          # Template renderer with 80+ registrations
│   │
│   └── models/                  # (Not used - sqlc generates models)
│
├── templates/
│   ├── admin/
│   │   ├── layouts/
│   │   │   └── base.html        # Admin layout (sidebar, header)
│   │   ├── pages/               # Full admin pages (dashboard, forms, lists)
│   │   └── partials/            # HTMX fragments (solution_stats, tag_chip)
│   │
│   ├── public/
│   │   ├── layouts/
│   │   │   └── base.html        # Public layout (header, footer)
│   │   ├── pages/               # Full public pages (home, products, blog)
│   │   └── partials/            # HTMX fragments (search_suggestions)
│   │
│   └── partials/
│       ├── header.html          # Public site header
│       ├── footer.html          # Public site footer
│       └── admin-sidebar.html   # Admin navigation sidebar
│
├── db/
│   ├── migrations/              # SQL migration files (001_*.up.sql, *.down.sql)
│   ├── queries/                 # sqlc SQL query files (*.sql)
│   ├── sqlc/                    # Generated Go code from sqlc
│   └── seeds/                   # Seed data scripts
│
├── public/
│   ├── css/                     # Stylesheets (brutalist design tokens)
│   ├── js/                      # JavaScript (htmx, vendor libs)
│   │   └── vendor/              # Third-party JS libraries
│   └── uploads/                 # User-uploaded files
│       ├── products/            # Product images
│       ├── downloads/           # Product downloads
│       ├── whitepapers/         # Whitepaper PDFs
│       ├── authors/             # Blog author photos
│       ├── solutions/           # Solution images
│       ├── blog/                # Blog post images
│       ├── case-studies/        # Case study images
│       └── categories/          # Category images
│
├── mockups/                     # Design mockups and wireframes
├── automation/                  # Automation scripts and phase tracking
├── plans/                       # Project planning documents
├── deploy/                      # Deployment configurations
├── bin/                         # Compiled binaries
│
└── sqlc.yaml                    # sqlc configuration
```

## Handler Architecture

### Handler Pattern
All handlers follow a consistent dependency injection pattern:

```go
type ProductsHandler struct {
    queries   *sqlc.Queries        // Database queries
    logger    *slog.Logger          // Structured logging
    uploadSvc *services.UploadService  // Optional: for file uploads
    cache     *services.Cache       // Optional: for caching
}

func NewProductsHandler(queries *sqlc.Queries, logger *slog.Logger, ...) *ProductsHandler {
    return &ProductsHandler{
        queries: queries,
        logger:  logger,
        // ...
    }
}

func (h *ProductsHandler) List(c echo.Context) error {
    ctx := c.Request().Context()

    // Extract parameters
    page, _ := strconv.Atoi(c.QueryParam("page"))

    // Query database
    products, err := h.queries.ListProducts(ctx, params)
    if err != nil {
        h.logger.Error("failed to list products", "error", err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    // Render template
    return c.Render(http.StatusOK, "admin/pages/products_list.html", map[string]interface{}{
        "Title": "Products",
        "Products": products,
    })
}
```

### Admin vs Public Handlers

**Admin Handlers** (`internal/handlers/admin/`):
- Protected by `RequireAuth()` middleware
- Full CRUD operations (Create, Read, Update, Delete)
- POST/DELETE methods for mutations
- Return 302 redirects after form submissions
- HTMX endpoints return HTML fragments
- Log activities via `ActivityLogService`
- Invalidate cache after mutations

**Public Handlers** (`internal/handlers/public/`):
- Read-only operations
- GET requests only
- Use caching extensively (TTL 600s)
- Load site settings via `SettingsLoader()` middleware
- SEO-optimized (meta tags, sitemaps)
- Rate-limited on sensitive endpoints (contact form)

## Middleware Chain

Middleware executes in the order registered in `main.go`:

### 1. Recovery Middleware
```go
func Recovery(logger *slog.Logger) echo.MiddlewareFunc
```
- Catches panics from downstream handlers
- Logs full stack trace with request context
- Returns 500 JSON error response
- Prevents server crashes

### 2. Logging Middleware
```go
func Logging(logger *slog.Logger) echo.MiddlewareFunc
```
- Logs every request with structured fields:
  - method, path, status, duration_ms, ip
- Uses `slog.Info()` for normal requests
- Executes after handler completes (records final status)

### 3. Gzip Middleware (Echo built-in)
```go
middleware.Gzip()
```
- Compresses response bodies
- Automatic content negotiation
- Reduces bandwidth usage

### 4. SecurityHeaders Middleware
```go
func SecurityHeaders() echo.MiddlewareFunc
```
Sets security headers:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy: default-src 'self'; ...`

### 5. SessionMiddleware
```go
func SessionMiddleware() echo.MiddlewareFunc
```
- Uses `gorilla/sessions` with cookie store
- Loads session from `bluejay_session` cookie
- Extracts user data (UserID, Email, DisplayName, Role)
- Stores in `c.Set("session", sess)`
- 7-day expiration, HttpOnly, SameSite=Lax

### 6. SettingsLoader Middleware (Public Routes Only)
```go
func SettingsLoader(queries *sqlc.Queries) echo.MiddlewareFunc
```
Loads site-wide data for public pages:
- Site settings (title, logo, social links)
- Footer categories (product categories)
- Footer solutions (published solutions)
- Footer resources (page sections)
- Stores in Echo context for template access

### 7. RequireAuth Middleware (Admin Routes Only)
```go
func RequireAuth() echo.MiddlewareFunc
```
- Checks `session.UserID > 0`
- Redirects to `/admin/login` if not authenticated
- Used as route group middleware: `e.Group("/admin", RequireAuth())`

### 8. RateLimiter Middleware (Specific Routes)
```go
func NewRateLimiter(limit int, window time.Duration) *RateLimiter
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc
```
- IP-based rate limiting with sliding window
- Applied per-route, not globally (e.g., contact form: 5 requests/hour)
- Example: `publicGroup.POST("/contact/submit", handler, contactLimiter.Middleware())`
- Returns 429 Too Many Requests when exceeded
- Background cleanup goroutine removes old entries

## Service Layer

Services encapsulate business logic and provide reusable operations across handlers.

### ProductService
```go
type ProductService struct {
    queries *sqlc.Queries
}

func (s *ProductService) GetProductDetail(ctx context.Context, slug string) (*ProductDetail, error)
```
**Purpose**: Aggregates product data from multiple tables
- Fetches product, category, specs, images, features, certifications, downloads
- Returns single `ProductDetail` struct
- Handles missing optional data gracefully (empty slices)

**Used by**: Public product detail pages

### UploadService
```go
type UploadService struct {
    uploadDir string
}

func (s *UploadService) UploadProductImage(file *multipart.FileHeader) (string, error)
func (s *UploadService) UploadProductDownload(file *multipart.FileHeader) (string, error)
```
**Purpose**: Handles file uploads with validation
- Image validation: .jpg, .jpeg, .png, .webp (max 5MB)
- Download validation: any file (max 50MB)
- Generates timestamped filenames
- Creates subdirectories (`products/`, `downloads/`)
- Returns web-accessible path (`/uploads/products/123_image.jpg`)

**Used by**: Admin product handlers, media library

### Cache Service
```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]cacheItem
}

func (c *Cache) Get(key string) (interface{}, bool)
func (c *Cache) Set(key string, value interface{}, ttlSeconds int)
func (c *Cache) Delete(key string)
func (c *Cache) DeleteByPrefix(prefix string)
```
**Purpose**: In-memory TTL cache for rendered pages and data
- Thread-safe with RWMutex
- Automatic expiration with background cleanup goroutine
- Prefix-based invalidation (e.g., `cache.DeleteByPrefix("page:products")`)
- Used for public pages (600s TTL) and frequently accessed data

**Cache Keys**:
- `page:products` - Product listing page
- `page:products:category-slug` - Category pages
- `page:products:category-slug:product-slug` - Product detail pages
- `page:blog` - Blog listing
- Similar patterns for solutions, case studies, whitepapers

**Invalidation Strategy**:
- On create/update/delete in admin panel
- Delete by prefix to clear related pages
- Example: Updating a product clears `page:products:*`

### ActivityLogService
```go
type ActivityLogService struct {
    queries *sqlc.Queries
    logger  *slog.Logger
}

func (s *ActivityLogService) Log(ctx context.Context, userID int64, action, resourceType string, resourceID int64, resourceTitle, description string)
func (s *ActivityLogService) LogSimple(ctx context.Context, userID int64, action, description string)
```
**Purpose**: Audit logging for admin actions
- Fire-and-forget logging (errors logged but not returned)
- Tracks: user, action type, resource type/ID, timestamp
- Used throughout admin handlers for compliance and debugging
- Viewable in admin activity log page

**Example Actions**:
- "created", "updated", "deleted", "published", "login", "logout"

**Resource Types**:
- "product", "blog_post", "solution", "case_study", "whitepaper", "system"

## Template System

### Template Renderer
```go
type Renderer struct {
    templates map[string]*template.Template
    basePath  string
}
```

The renderer pre-compiles 80+ templates at startup using `template.Must()`. Templates are organized hierarchically:

```
Layouts (define "base")
  ↓
Pages (define "content", reference layouts)
  ↓
Partials (standalone or included)
```

### Template Inheritance Pattern

**Admin Layout** (`admin/layouts/base.html`):
```html
{{define "base"}}
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - Admin Panel</title>
    <script src="/public/js/vendor/htmx.min.js"></script>
</head>
<body>
    {{template "content" .}}
</body>
</html>
{{end}}
```

**Admin Page** (`admin/pages/dashboard.html`):
```html
{{define "content"}}
<div class="admin-container">
    {{template "admin-sidebar" .}}
    <main>
        <h1>{{.Title}}</h1>
        <!-- Page content -->
    </main>
</div>
{{end}}
```

**HTMX Partial** (`admin/partials/solution_stats.html`):
```html
{{define "base"}}
<div id="solution-stats">
    {{range .Stats}}
    <div class="stat-item">
        <span>{{.Label}}</span>
        <span>{{.Value}}</span>
    </div>
    {{end}}
</div>
{{end}}
```

### Template Registration
All templates are registered in `internal/templates/template.go`:

```go
func (r *Renderer) loadTemplates() {
    funcMap := template.FuncMap{
        "safeHTML":       safeHTML,       // Renders unescaped HTML
        "formatDate":     formatDate,     // Formats time.Time
        "truncate":       truncate,       // Truncates strings
        "slugify":        slugify,        // Creates URL slugs
        "formatFileSize": formatFileSize, // Bytes to KB/MB
        "now":            time.Now,       // Current timestamp
        "add":            func(a, b int) int { return a + b },
        "sub":            func(a, b int) int { return a - b },
        "upper":          strings.ToUpper,
        "int64":          func(i int) int64 { return int64(i) }, // Type conversion
        "seq":            func(n int64) []int { ... }, // Range helper
    }

    // Full page templates with layouts
    r.templates["admin/pages/dashboard.html"] = template.Must(
        template.New("base").Funcs(funcMap).ParseFiles(
            "templates/admin/layouts/base.html",
            "templates/admin/pages/dashboard.html",
            "templates/partials/admin-sidebar.html",
        )
    )

    // HTMX partials (standalone, no layout)
    r.templates["admin/partials/solution_stats.html"] = template.Must(
        template.New("base").Funcs(funcMap).ParseFiles(
            "templates/admin/partials/solution_stats.html",
        )
    )
}
```

### Template Functions
Custom functions available in all templates:

| Function | Purpose | Example |
|----------|---------|---------|
| `safeHTML` | Renders HTML without escaping | `{{.Content \| safeHTML}}` |
| `formatDate` | Formats dates | `{{formatDate .PublishedAt "Jan 2, 2006"}}` |
| `truncate` | Truncates strings | `{{truncate .Description 100}}` |
| `slugify` | Creates URL slugs | `{{slugify .Name}}` |
| `formatFileSize` | Formats bytes | `{{formatFileSize .Size}}` |
| `now` | Returns current time | `{{now}}` |
| `add` | Integer addition | `{{add .Page 1}}` |
| `sub` | Integer subtraction | `{{sub .Total .Used}}` |
| `upper` | Uppercase string | `{{upper .Status}}` |
| `int64` | Convert int to int64 | `{{int64 .Count}}` |
| `seq` | Generate sequence | `{{range seq 5}}` (0,1,2,3,4) |

## HTMX Integration

HTMX enables dynamic page updates without full page reloads. The backend returns HTML fragments instead of JSON.

### HTMX Patterns in Bluejay CMS

#### 1. Dynamic List Updates
```html
<!-- Admin tag management -->
<div id="tag-list">
    <button hx-get="/admin/blog/tags/search?q=golang"
            hx-target="#tag-suggestions"
            hx-trigger="click">
        Search Tags
    </button>
    <div id="tag-suggestions"></div>
</div>
```

**Handler**:
```go
func (h *BlogTagsHandler) Search(c echo.Context) error {
    tags, _ := h.queries.SearchTags(ctx, query)
    return c.Render(http.StatusOK, "admin/partials/tag_suggestions.html", map[string]interface{}{
        "Tags": tags,
    })
}
```

#### 2. Form Submission with Inline Updates
```html
<form hx-post="/admin/products/123/features"
      hx-target="#feature-list"
      hx-swap="beforeend">
    <input name="name" placeholder="Feature name">
    <button type="submit">Add Feature</button>
</form>
<div id="feature-list">
    <!-- Features rendered here -->
</div>
```

**Handler**:
```go
func (h *ProductDetailsHandler) AddFeature(c echo.Context) error {
    // Create feature in database
    feature, _ := h.queries.CreateProductFeature(ctx, params)

    // Return single feature HTML fragment
    return c.Render(http.StatusOK, "admin/partials/feature_item.html", map[string]interface{}{
        "Feature": feature,
    })
}
```

#### 3. Delete with Removal
```html
<div id="feature-{{.ID}}">
    <span>{{.Name}}</span>
    <button hx-delete="/admin/products/{{.ProductID}}/features/{{.ID}}"
            hx-target="#feature-{{.ID}}"
            hx-swap="outerHTML"
            hx-confirm="Delete this feature?">
        Delete
    </button>
</div>
```

**Handler**:
```go
func (h *ProductDetailsHandler) DeleteFeature(c echo.Context) error {
    h.queries.DeleteProductFeature(ctx, id)
    // Return empty response - hx-swap="outerHTML" removes the element
    return c.NoContent(http.StatusOK)
}
```

#### 4. Autocomplete Search
```html
<input type="text"
       hx-get="/search/suggest"
       hx-trigger="keyup changed delay:300ms"
       hx-target="#suggestions"
       name="q">
<div id="suggestions"></div>
```

**Handler**:
```go
func (h *SearchHandler) SearchSuggest(c echo.Context) error {
    results := performSearch(query)
    return c.Render(http.StatusOK, "public/partials/search_suggestions.html", map[string]interface{}{
        "Results": results,
    })
}
```

### HTMX Attributes Used

- `hx-get`, `hx-post`, `hx-delete` - HTTP methods
- `hx-target` - CSS selector for update target
- `hx-swap` - Swap strategy (innerHTML, outerHTML, beforeend, etc.)
- `hx-trigger` - Event trigger (click, keyup, load, etc.)
- `hx-confirm` - Confirmation dialog before action
- `hx-include` - Include additional form fields
- `hx-indicator` - Show loading indicator

## Database Layer

### SQLite Configuration
```go
DSN: "bluejay.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on&_synchronous=NORMAL&_cache_size=2000"

Connection Pool:
- SetMaxOpenConns(1)     // Single writer (SQLite limitation)
- SetMaxIdleConns(1)     // Keep connection alive
- SetConnMaxLifetime(0)  // No expiration
- SetConnMaxIdleTime(0)  // Never close idle
```

**WAL Mode Benefits**:
- Readers don't block writers
- Writers don't block readers
- Better concurrency for web apps
- Automatic checkpointing

### sqlc Code Generation

**Configuration** (`sqlc.yaml`):
```yaml
version: "2"
sql:
  - engine: "sqlite"
    queries: "db/queries"
    schema: "db/migrations"
    gen:
      go:
        package: "sqlc"
        out: "db/sqlc"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
```

**SQL Query File** (`db/queries/products.sql`):
```sql
-- name: GetProduct :one
SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE status = 'published'
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: CreateProduct :one
INSERT INTO products (name, slug, description, ...)
VALUES (?, ?, ?, ...)
RETURNING *;
```

**Generated Go Code** (`db/sqlc/products.sql.go`):
```go
type Product struct {
    ID          int64
    Name        string
    Slug        string
    Description string
    // ...
}

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error)
func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
```

### Migrations

**Migration File** (`db/migrations/001_create_products.up.sql`):
```sql
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sku TEXT UNIQUE NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    category_id INTEGER NOT NULL,
    status TEXT DEFAULT 'draft',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE RESTRICT
);

CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_category ON products(category_id);
```

**Rollback File** (`db/migrations/001_create_products.down.sql`):
```sql
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_products_status;
DROP TABLE IF EXISTS products;
```

**Migration Runner**:
```go
func RunMigrations(db *sql.DB, migrationsPath string) error {
    driver, _ := sqlite.WithInstance(db, &sqlite.Config{})
    m, _ := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "sqlite", driver)
    return m.Up() // Runs all pending migrations
}
```

## Caching Strategy

### Cache Implementation
```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]cacheItem
}

type cacheItem struct {
    value     interface{}
    expiresAt time.Time
}
```

### Caching Pattern in Handlers
```go
func (h *ProductsHandler) ProductsList(c echo.Context) error {
    cacheKey := "page:products"

    // Check cache first
    if cached, ok := h.cache.Get(cacheKey); ok {
        return c.HTML(http.StatusOK, cached.(string))
    }

    // Fetch data
    products, _ := h.queries.ListProducts(ctx)

    // Render template
    var buf bytes.Buffer
    c.Echo().Renderer.Render(&buf, "public/pages/products.html", data, c)
    html := buf.String()

    // Store in cache (600 second TTL)
    h.cache.Set(cacheKey, html, 600)

    return c.HTML(http.StatusOK, html)
}
```

### Cache Invalidation
```go
// When admin updates a product
func (h *AdminProductsHandler) Update(c echo.Context) error {
    // Update database
    h.queries.UpdateProduct(ctx, params)

    // Invalidate related cache entries
    h.cache.DeleteByPrefix("page:products")  // Clears all product pages

    return c.Redirect(http.StatusSeeOther, "/admin/products")
}
```

### Cache TTLs
- Public pages: 600 seconds (10 minutes)
- Frequently changing data: Not cached
- Admin pages: Not cached (always fresh)

### Cache Cleanup
Background goroutine runs every 5 minutes to remove expired entries:
```go
func (c *Cache) cleanupLoop() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        c.mu.Lock()
        for k, item := range c.items {
            if time.Now().After(item.expiresAt) {
                delete(c.items, k)
            }
        }
        c.mu.Unlock()
    }
}
```

## Error Handling Patterns

### 1. Database Query Errors
```go
product, err := h.queries.GetProduct(ctx, id)
if err == sql.ErrNoRows {
    return echo.NewHTTPError(http.StatusNotFound, "Product not found")
}
if err != nil {
    h.logger.Error("failed to get product", "id", id, "error", err)
    return echo.NewHTTPError(http.StatusInternalServerError)
}
```

### 2. Form Validation Errors
```go
if name == "" || slug == "" {
    return c.Render(http.StatusBadRequest, "admin/pages/product_form.html", map[string]interface{}{
        "Error": "Name and slug are required",
        "Item":  formData,
    })
}
```

### 3. Middleware Panic Recovery
```go
func Recovery(logger *slog.Logger) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            defer func() {
                if r := recover(); r != nil {
                    logger.Error("panic recovered",
                        "error", r,
                        "stack", string(debug.Stack()))
                    c.JSON(http.StatusInternalServerError, map[string]string{
                        "error": "Internal server error",
                    })
                }
            }()
            return next(c)
        }
    }
}
```

### 4. Graceful Degradation
```go
// Dashboard loads statistics, but doesn't fail if one query fails
if count, err := h.queries.CountProducts(ctx); err == nil {
    data.PublishedProducts = count
} else {
    h.logger.Error("dashboard: count products", "error", err)
    // Continue with zero count
}
```

### 5. Activity Logging Errors (Fire-and-Forget)
```go
func (s *ActivityLogService) Log(ctx context.Context, ...) {
    err := s.queries.CreateActivityLog(ctx, params)
    if err != nil {
        s.logger.Error("failed to log activity", "error", err)
        // Don't return error - activity logging should never fail the main operation
    }
}
```

## Data Flow Diagrams

### Creating a Product (Admin)

```
User fills form
     ↓
POST /admin/products
     ↓
AdminProductsHandler.Create()
     ↓
├─ Validate form inputs
├─ Upload product image (UploadService)
├─ Generate slug from name
├─ Insert product (sqlc: CreateProduct)
├─ Log activity (ActivityLogService)
└─ Invalidate cache (cache.DeleteByPrefix("page:products"))
     ↓
Redirect to /admin/products (list page)
```

### Viewing a Public Product Page

```
User visits /products/category/product-slug
     ↓
GET /products/category/product-slug
     ↓
PublicProductsHandler.ProductDetail()
     ↓
Check cache: cache.Get("page:products:category:slug")
     ↓
├─ Cache HIT → Return cached HTML
└─ Cache MISS ↓
       ├─ ProductService.GetProductDetail(slug)
       │   ├─ GetProductBySlug()
       │   ├─ GetProductCategory()
       │   ├─ ListProductSpecs()
       │   ├─ ListProductImages()
       │   ├─ ListProductFeatures()
       │   ├─ ListProductCertifications()
       │   └─ ListProductDownloads()
       ├─ Render template with data
       ├─ Store in cache (TTL: 600s)
       └─ Return HTML response
```

### Contact Form Submission

```
User fills contact form
     ↓
POST /contact/submit (Rate Limited: 5/hour per IP)
     ↓
RateLimiter checks IP
     ↓
├─ Exceeded → 429 Too Many Requests
└─ Allowed ↓
       ContactHandler.SubmitContactForm()
       ↓
       ├─ Validate form (email format, required fields)
       ├─ Insert submission (CreateContactSubmission)
       ├─ Optional: Send email notification
       └─ Return success HTML fragment (HTMX swap)
```

### Admin Login Flow

```
User enters credentials
     ↓
POST /admin/login
     ↓
AuthHandler.LoginSubmit()
     ↓
├─ GetAdminUserByEmail(email)
│   └─ Not found → Redirect with error
├─ bcrypt.CompareHashAndPassword(stored, input)
│   └─ Mismatch → Redirect with error
├─ UpdateLastLogin(userID)
├─ Save session (UserID, Email, DisplayName, Role)
├─ Log activity ("login", "system")
└─ Redirect to /admin/dashboard
     ↓
RequireAuth() middleware checks session.UserID
     ↓
DashboardHandler.ShowDashboard()
     ↓
├─ Aggregate statistics (products, posts, contacts)
└─ Render dashboard template
```

## Performance Considerations

### 1. Template Pre-compilation
- All templates compiled at startup (80+ templates)
- Zero runtime compilation overhead
- Fast template lookup from map

### 2. Connection Pooling
- Single connection pool (SQLite limitation)
- Long-lived connection (no reconnection overhead)
- WAL mode for concurrent reads

### 3. In-Memory Caching
- Rendered HTML cached (no re-rendering)
- 600s TTL for public pages
- Prefix-based invalidation (efficient cache clearing)

### 4. HTMX Partial Updates
- Only updated sections re-rendered
- Reduced bandwidth (HTML fragments vs full pages)
- No client-side JS rendering

### 5. Database Indexing
- Indexes on frequently queried columns (status, category_id, slug)
- Foreign keys for referential integrity

### 6. Static Asset Serving
- Static files served directly by Echo
- Gzip compression enabled
- CDN-ready (Caddy/Nginx can add caching headers)

## Security Measures

### 1. Authentication
- Bcrypt password hashing (cost 10)
- Session-based auth (HttpOnly cookies)
- SameSite=Lax (CSRF protection)

### 2. Authorization
- `RequireAuth()` middleware on admin routes
- Session validation on every request

### 3. Input Validation
- Form validation in handlers
- File upload restrictions (type, size)
- SQL injection prevented by sqlc parameterized queries

### 4. Security Headers
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- Content-Security-Policy: restrictive
- X-XSS-Protection: 1; mode=block

### 5. Rate Limiting
- IP-based rate limiting on sensitive endpoints
- Example: Contact form (5 requests/hour)

### 6. Panic Recovery
- Graceful error handling
- Stack traces logged (not exposed to users)

## Deployment Architecture

```
Internet
   ↓
Caddy (Reverse Proxy)
   ├─ TLS termination
   ├─ Static file caching
   ├─ Gzip compression
   └─ Proxy to :28090
      ↓
Bluejay CMS (Go binary)
   ├─ Port 28090
   ├─ Echo web server
   └─ SQLite database
      ↓
Litestream (Continuous backup)
   └─ Replicates to S3 every 10s
```

### Production Checklist
- [ ] Change session secret (32+ chars)
- [ ] Enable TLS in Caddy/Nginx
- [ ] Configure Litestream for S3 backups
- [ ] Set up log rotation
- [ ] Configure systemd service
- [ ] Set environment variables (PORT, DB_PATH)
- [ ] Run database migrations
- [ ] Seed admin user
- [ ] Configure firewall (allow 80, 443, SSH)
- [ ] Set up monitoring (Prometheus, Grafana)

## Testing Strategy

### Unit Tests
- Service layer tests (ProductService, Cache)
- Middleware tests (auth, rate limiting)
- Template function tests

### Integration Tests
- Handler tests with mock database
- Database query tests (sqlc-generated code)

### E2E Tests
- Admin workflow tests (create product, publish blog post)
- Public page rendering tests
- HTMX interaction tests

## Future Enhancements

### Planned Features
- [ ] Multi-user admin roles (editor, viewer, admin)
- [ ] Content versioning and rollback
- [ ] Image optimization pipeline
- [ ] Full-text search (SQLite FTS5)
- [ ] Email notifications (SMTP integration)
- [ ] Webhooks for external integrations
- [ ] API endpoints (JSON REST API)
- [ ] Import/export functionality
- [ ] Advanced SEO tools (schema.org markup)
- [ ] A/B testing framework

### Scalability Considerations
- Current architecture supports ~1000 req/s on single server
- For higher traffic: Add read replicas (Litestream restore)
- Consider PostgreSQL migration for high-write workloads
- Add Redis for distributed caching
- Implement CDN for static assets

---

**Last Updated**: 2026-02-10
**Version**: 1.0
**Maintainer**: Bluejay CMS Team
