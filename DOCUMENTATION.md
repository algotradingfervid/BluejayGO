# Bluejay CMS — Complete Documentation

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Tech Stack](#tech-stack)
4. [Directory Structure](#directory-structure)
5. [Database Design](#database-design)
6. [Backend Architecture](#backend-architecture)
7. [Frontend Architecture](#frontend-architecture)
8. [Authentication & Security](#authentication--security)
9. [Development Setup](#development-setup)
10. [Deployment Guide](#deployment-guide)
11. [Usage Guide](#usage-guide)
12. [API & Routes Reference](#api--routes-reference)
13. [Configuration Reference](#configuration-reference)

---

## Overview

Bluejay CMS is a full-featured content management system built in Go for managing a technology company's website. It provides both a **public-facing website** and a **brutalist-styled admin panel** for managing products, blog posts, solutions, case studies, whitepapers, partners, and more.

**Key Characteristics:**
- Server-side rendered with Go templates — no JavaScript framework
- HTMX for dynamic interactions without page reloads
- SQLite database with type-safe queries via sqlc
- Brutalist design system (JetBrains Mono, 2px black borders, 4px box shadows)
- 22 implementation phases, all completed

---

## Architecture

### System Design

```
┌─────────────────────────────────────────────────────────┐
│                        Caddy                            │
│              (Reverse Proxy + Auto TLS)                 │
│           yourdomain.com → localhost:28090               │
├─────────────────────────────────────────────────────────┤
│                     Echo v4 Server                       │
│                      Port 28090                          │
├──────────┬──────────┬───────────┬───────────────────────┤
│Middleware│ Handlers  │ Services  │ Template Renderer     │
│ Stack    │          │           │                        │
│ Recovery │ Admin    │ Product   │ Go html/template       │
│ Logging  │  Auth    │ Upload    │ + HTMX partials        │
│ Gzip     │  CRUD    │ Cache     │ + Tailwind CSS (CDN)   │
│ Security │ Public   │ Activity  │ + Trix Editor          │
│ Sessions │  Pages   │  Logger   │                        │
├──────────┴──────────┴───────────┴───────────────────────┤
│              sqlc Generated Queries                      │
├─────────────────────────────────────────────────────────┤
│           SQLite (WAL mode, FTS5 search)                │
├─────────────────────────────────────────────────────────┤
│        Litestream → S3 (continuous backup)              │
└─────────────────────────────────────────────────────────┘
```

### Request Flow

```
HTTP Request
  → Caddy (TLS termination, static files, gzip)
    → Echo Middleware Chain (recovery → logging → security → sessions)
      → Route Handler
        → Service Layer (business logic, caching)
          → sqlc Queries → SQLite
        → Template Renderer → HTML Response
```

### Data Flow for Admin Operations

```
POST /admin/products
  → RequireAuth middleware (session check)
  → ProductsHandler.Create()
    → Validate form input
    → UploadService.UploadProductImage() (if file attached)
    → sqlc.CreateProduct() → SQLite INSERT
    → ActivityLogService.Log() → audit trail
    → Cache.DeleteByPrefix("page:products") → invalidate
  → HTTP 302 Redirect → product list
```

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.25 | Backend runtime |
| Web Framework | Echo v4 | HTTP routing, middleware |
| Database | SQLite (modernc.org/sqlite) | Data storage, WAL mode |
| Query Gen | sqlc | Type-safe SQL → Go code |
| Migrations | golang-migrate | Schema versioning (34 migrations) |
| Sessions | gorilla/sessions | Cookie-based auth sessions |
| Crypto | golang.org/x/crypto | bcrypt password hashing |
| Templates | Go html/template | Server-side rendering |
| Interactivity | HTMX | Dynamic HTML updates |
| Rich Text | Trix Editor | Blog content editing |
| CSS | Tailwind CSS (CDN) | Utility-first styling |
| Reverse Proxy | Caddy | TLS, static files, headers |
| DB Backup | Litestream | Continuous SQLite → S3 replication |

---

## Directory Structure

```
bluejay-cms/
├── cmd/server/main.go           # Entry point — server setup, routing
├── internal/
│   ├── handlers/
│   │   ├── admin/               # 25 files — admin CRUD handlers
│   │   └── public/              # 13 files — public page handlers
│   ├── middleware/               # 10 files — auth, CSRF, logging, security, sessions, rate limit
│   ├── services/                # 4 services — product, upload, cache, activity log
│   ├── models/                  # Domain models
│   ├── database/                # SQLite init, migrations runner
│   └── templates/               # Template renderer (80+ template registrations)
├── db/
│   ├── migrations/              # 68 files (34 up + 34 down)
│   ├── queries/                 # 23 sqlc SQL files
│   ├── sqlc/                    # Auto-generated Go code (DO NOT EDIT)
│   └── seeds/                   # 24 seed data files
├── templates/
│   ├── admin/
│   │   ├── layouts/base.html    # Admin layout
│   │   ├── pages/               # 60 admin page templates
│   │   └── partials/            # 15 reusable HTMX fragments
│   ├── public/
│   │   ├── layouts/base.html    # Public layout
│   │   ├── pages/               # 16 public page templates
│   │   └── partials/            # 2 search fragments
│   └── partials/                # 5 shared components (header, footer, sidebar)
├── public/
│   ├── css/                     # styles.css, admin-styles.css, trix.css
│   ├── js/                      # admin.js, htmx.min.js, trix.js
│   └── uploads/                 # User-uploaded media (products, blog, etc.)
├── deploy/
│   ├── Caddyfile                # Reverse proxy config
│   ├── bluejay-cms.service      # systemd service unit
│   ├── build.sh                 # Cross-compile for Linux
│   └── litestream.yml           # SQLite backup to S3
├── plans/                       # 23 phase implementation docs
├── automation/                  # Phase tracker, orchestrator, logs
├── Makefile                     # Build, run, test, deploy commands
├── sqlc.yaml                    # sqlc code generation config
├── seed.sql                     # Database seed data
└── CLAUDE.md                    # Project instructions
```

---

## Database Design

### Schema Overview

SQLite with WAL journal mode, foreign keys enforced. 34 migrations define the schema.

### Core Tables

**Content Management:**
- `products` — Product catalog (name, slug, description, status, category, images)
- `product_specs`, `product_features`, `product_certifications`, `product_downloads`, `product_images` — Product details (1:M)
- `product_categories` — Product taxonomy
- `blog_posts` — Blog articles (title, body, excerpt, status, author, category)
- `blog_categories`, `blog_authors`, `blog_tags`, `blog_post_tags`, `blog_post_products` — Blog taxonomy and relations
- `solutions` — Solution pages with related stats, challenges, products, CTAs
- `case_studies` — Case studies with metrics and related products
- `whitepapers` — Downloadable content with download tracking
- `partners`, `partner_testimonials`, `partner_tiers` — Partner ecosystem

**Website Sections:**
- `homepage_heroes`, `homepage_stats`, `homepage_testimonials`, `homepage_ctas` — Homepage content
- `about_overview`, `about_values`, `about_milestones`, `about_certifications` — About page
- `contact_submissions`, `contact_offices` — Contact form data
- `page_sections` — Reusable content sections

**System:**
- `admin_users` — Authentication (email, bcrypt hash, role)
- `settings` — Global site settings (singleton)
- `header_settings`, `footer_settings` — Header/footer config
- `media_files` — Media library metadata
- `navigation_menus`, `navigation_items` — Menu structure
- `activity_log` — Audit trail
- `section_settings` — Per-section visibility toggles
- `fts5_search` — Full-text search index (products, blogs, solutions)

### Key Relationships

```
products ──→ product_categories (M:1)
products ──→ product_specs/images/features/certifications/downloads (1:M)
blog_posts ──→ blog_authors (M:1), blog_categories (M:1)
blog_posts ←──→ blog_tags (M:M via blog_post_tags)
blog_posts ←──→ products (M:M via blog_post_products)
solutions ←──→ products (M:M)
case_studies ←──→ products (M:M)
partners ──→ partner_tiers (M:1)
```

### SQLite Configuration

```
Journal Mode:  WAL (concurrent reads during writes)
Busy Timeout:  5000ms
Foreign Keys:  ON
Synchronous:   NORMAL
Cache Size:    2000 pages
Max Connections: 1 (SQLite design)
```

---

## Backend Architecture

### Middleware Stack

Executed in order for every request:

1. **Recovery** — Catches panics, returns 500
2. **Logging** — Structured JSON logs (method, path, status, duration, IP)
3. **Gzip** — Response compression
4. **SecurityHeaders** — CSP, X-Frame-Options, X-Content-Type-Options, XSS-Protection
5. **SessionMiddleware** — Gorilla session loading
6. **RequireAuth** — (admin routes only) Validates session, redirects to login
7. **SettingsLoader** — (public routes only) Loads global settings into context

### Handler Pattern

All handlers follow the same structure:

```go
type Handler struct {
    queries   *sqlc.Queries
    logger    *slog.Logger
    uploadSvc *services.UploadService  // optional
    cache     *services.Cache          // optional
}

func (h *Handler) ListProducts(c echo.Context) error {
    products, err := h.queries.ListProducts(c.Request().Context())
    if err != nil {
        return echo.NewHTTPError(500, "Failed to load products")
    }
    return c.Render(200, "admin-products-list", map[string]interface{}{
        "Items": products,
    })
}
```

### Services

| Service | Responsibilities |
|---------|-----------------|
| **ProductService** | Aggregates product with specs, images, features, certifications, downloads |
| **UploadService** | File validation (type, size), storage to `/public/uploads/`, cleanup |
| **Cache** | In-memory TTL cache with RWMutex, background cleanup every 5 min |
| **ActivityLogService** | Async audit logging — user, action, resource type, description |

### Template System

Custom renderer wrapping Go's `html/template`:
- 80+ templates registered at startup
- Template functions: `safeHTML`, `formatDate`, `truncate`, `slugify`, `formatFileSize`, `add`, `sub`, `seq`
- Inheritance via `{{template "admin-layout" .}}` with `{{block "content" .}}`
- HTMX endpoints return partials (no layout wrapper)

---

## Frontend Architecture

### Design System — Brutalist

- **Font:** JetBrains Mono everywhere (monospace)
- **Borders:** 2px solid black — no border-radius (0px on all elements)
- **Shadows:** `box-shadow: 4px 4px 0 0 #000` (manual, not Tailwind utilities)
- **Buttons:** Uppercase, thick borders, push-down effect on click (`transform: translate(2px, 2px)`)
- **Colors:** Black/white primary. Accent: `#0066CC` (blue), `#DC2626` (red), `#16A34A` (green)

### HTMX Patterns (221 occurrences across 41 files)

| Pattern | Example |
|---------|---------|
| Tab switching | `hx-get="/admin/products/1/features" hx-target="#detail-content"` |
| Delete with confirm | `hx-delete="/admin/products/1" hx-confirm="Delete?" hx-target="closest tr"` |
| Inline form submit | `hx-post="/admin/products/1/specs" hx-target="#specs-section" hx-swap="outerHTML"` |
| Auto-load on page | `hx-get="/admin/products/1/specs" hx-trigger="load"` |
| Search suggestions | `hx-get="/search/suggest?q=..." hx-trigger="keyup changed delay:300ms"` |

### Admin Layout

```
┌──────────────────────────────────────────────┐
│  Sidebar (fixed)  │  Header Bar              │
│  ─────────────    │──────────────────────────│
│  Dashboard        │                          │
│  WEBSITE          │  Page Content            │
│   Homepage        │  (list / form / detail)  │
│   About           │                          │
│   Header/Footer   │                          │
│  CONTENT          │                          │
│   Products        │                          │
│   Blog            │                          │
│   Solutions       │                          │
│  ADMIN            │                          │
│   Media Library   │                          │
│   Settings        │                          │
└──────────────────────────────────────────────┘
```

- Sidebar groups collapse/expand with localStorage persistence
- Mobile: drawer pattern with overlay toggle

### Public Layout

- Responsive grid layouts (1→2→3 columns)
- Grid-dotted hero backgrounds
- Material Symbols icons
- Inter font for display headings, JetBrains Mono for body

---

## Authentication & Security

### Authentication Flow

1. User submits email + password to `POST /admin/login`
2. Handler queries `admin_users` table by email
3. bcrypt comparison of password hash
4. On success: create Gorilla session with UserID, Email, DisplayName, Role
5. Session cookie: HttpOnly, SameSite=Lax, 7-day MaxAge
6. `RequireAuth` middleware checks `UserID > 0` on all `/admin/*` routes

### Security Measures

| Feature | Implementation |
|---------|---------------|
| Password Storage | bcrypt hash |
| Session Security | HttpOnly cookies, SameSite=Lax |
| SQL Injection | Prevented by sqlc parameterized queries |
| XSS Protection | `X-XSS-Protection` header, CSP policy |
| Clickjacking | `X-Frame-Options: DENY` |
| MIME Sniffing | `X-Content-Type-Options: nosniff` |
| Rate Limiting | Contact form: 5 requests/hour |
| CSRF | Middleware available (csrf.go) |

### Production Hardening Checklist

- [ ] Change session secret from default placeholder
- [ ] Set session `Secure: true` for HTTPS
- [ ] Tighten CSP to remove `unsafe-inline` and `unsafe-eval`
- [ ] Enable CSRF middleware on all POST routes
- [ ] Use environment variable for database path
- [ ] Set up Litestream S3 credentials

---

## Development Setup

### Prerequisites

- Go 1.25+
- sqlc (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)
- SQLite3 CLI (for seeding)
- air (optional, for hot-reload: `go install github.com/air-verse/air@latest`)

### Quick Start

```bash
# Clone the repository
git clone <repo-url>
cd bluejay-cms

# Generate sqlc code
make sqlc

# Build the application
make build

# Seed the database with sample data
make seed

# Run the server
make run
# → Server starts at http://localhost:28090

# Or, for development with hot-reload
make dev
```

### Default Admin Credentials

Check `seed.sql` for the default admin user created during seeding.

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Start server (`go run cmd/server/main.go`) |
| `make build` | Compile binary to `bin/bluejay-cms` |
| `make dev` | Hot-reload with air |
| `make sqlc` | Regenerate sqlc Go code from SQL |
| `make migrate-up` | Run pending migrations |
| `make migrate-down` | Rollback last migration |
| `make seed` | Load sample data into database |
| `make test` | Run all tests (`go test -v ./...`) |
| `make clean` | Remove binaries and database files |
| `make deploy` | Build, upload, restart on server |

### Adding a New SQL Query

1. Write the SQL in `db/queries/<entity>.sql`
2. Run `make sqlc` to regenerate Go code
3. Use the generated function in your handler via `h.queries.YourQuery(ctx, params)`

### Adding a New Migration

1. Create `db/migrations/NNN_description.up.sql` and `.down.sql`
2. Run `make migrate-up`
3. Run `make sqlc` (schema changed)

---

## Deployment Guide

### Architecture (Production)

```
Internet → Caddy (443/HTTPS) → Go Binary (28090) → SQLite
                                                      ↓
                                                  Litestream → S3
```

### Step-by-Step Deployment

#### 1. Build for Linux

```bash
# From development machine
make deploy-build
# Or manually:
GOOS=linux GOARCH=amd64 go build -o bluejay-cms cmd/server/main.go
```

#### 2. Prepare the Server

```bash
# SSH into your server
ssh user@yourserver

# Create application directory
sudo mkdir -p /var/www/bluejay-cms
sudo chown www-data:www-data /var/www/bluejay-cms

# Install Caddy (https://caddyserver.com/docs/install)
# Install Litestream (https://litestream.io/install/)
```

#### 3. Upload Files

```bash
# From development machine
scp bluejay-cms user@yourserver:/var/www/bluejay-cms/
scp -r public/ user@yourserver:/var/www/bluejay-cms/
scp -r templates/ user@yourserver:/var/www/bluejay-cms/
scp -r db/migrations/ user@yourserver:/var/www/bluejay-cms/db/migrations/
scp deploy/Caddyfile user@yourserver:/etc/caddy/Caddyfile
scp deploy/bluejay-cms.service user@yourserver:/etc/systemd/system/
scp deploy/litestream.yml user@yourserver:/etc/litestream.yml
```

#### 4. Configure

Edit `/etc/caddy/Caddyfile`:
- Replace `yourdomain.com` with your actual domain

Edit `/etc/systemd/system/bluejay-cms.service`:
- Verify `DATABASE_PATH` and `WorkingDirectory`

Edit `/etc/litestream.yml`:
- Set S3 bucket name and region
- Configure AWS credentials via environment variables

#### 5. Start Services

```bash
# Enable and start the CMS
sudo systemctl daemon-reload
sudo systemctl enable bluejay-cms
sudo systemctl start bluejay-cms

# Start Caddy
sudo systemctl enable caddy
sudo systemctl start caddy

# Start Litestream (for backups)
sudo systemctl enable litestream
sudo systemctl start litestream
```

#### 6. Verify

```bash
# Check service status
sudo systemctl status bluejay-cms

# Check logs
sudo journalctl -u bluejay-cms -f

# Test health endpoint
curl https://yourdomain.com/health
```

### Updating the Application

```bash
# From development machine — one command does it all
make deploy

# Or manually:
make deploy-build      # Cross-compile
make deploy-upload     # SCP binary + configs
make deploy-restart    # systemctl restart
```

### Backup & Recovery

**Continuous Backup:** Litestream streams every SQLite WAL change to S3.

**Restore from Backup:**
```bash
litestream restore -o /var/www/bluejay-cms/bluejay.db s3://your-backup-bucket/bluejay-cms
```

---

## Usage Guide

### Admin Panel

Access at `https://yourdomain.com/admin/login`

#### Dashboard
- Overview of content counts (products, blog posts, solutions, etc.)
- Recent activity log

#### Managing Products
1. Navigate to **Content → Products → All Products**
2. Click **+ New Product** to create
3. Fill in: name, slug, category, description, status (draft/published)
4. After creation, manage sub-items via tabs:
   - **Specs** — Technical specifications (key/value pairs)
   - **Features** — Product features with descriptions
   - **Certifications** — Compliance badges
   - **Downloads** — Datasheets, manuals (up to 50MB)
   - **Images** — Product photos (up to 5MB, jpg/png/webp)

#### Managing Blog Posts
1. Navigate to **Content → Blog → Posts**
2. Click **+ New Post**
3. Use the **Trix rich text editor** for the post body
4. Assign author, category, tags
5. Set status to **Published** when ready

#### Homepage Management
- **Website → Homepage** sections:
  - Heroes (rotating hero banners)
  - Stats (counter numbers)
  - Testimonials (customer quotes)
  - CTAs (call-to-action blocks)

#### Media Library
- Upload images and files
- Browse all uploads with search
- Files organized by content type in `/public/uploads/`

#### Navigation Editor
- Create and manage navigation menus
- Drag-and-drop menu item ordering
- Nested menu structures supported

#### Activity Log
- Audit trail of all admin actions
- Shows who did what, when, to which resource

#### Global Settings
- Site name, tagline, contact info
- Section visibility toggles
- SEO defaults

### Public Website

| Page | URL | Description |
|------|-----|-------------|
| Home | `/` | Hero, featured products, stats, testimonials |
| Products | `/products` | Filterable product catalog |
| Product Detail | `/products/:category/:slug` | Full product page with specs |
| Solutions | `/solutions` | Solution listing |
| Blog | `/blog` | Blog listing with category filter |
| Blog Post | `/blog/:slug` | Full article with related products |
| Case Studies | `/case-studies` | Customer success stories |
| Whitepapers | `/whitepapers` | Downloadable resources |
| About | `/about` | Company info, values, milestones |
| Partners | `/partners` | Partner directory |
| Contact | `/contact` | Contact form with rate limiting |
| Search | `/search` | Full-text search across all content |

---

## API & Routes Reference

### Public Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | HomeHandler | Homepage |
| GET | `/health` | inline | Health check (JSON) |
| GET | `/products` | ProductsHandler.List | Product catalog |
| GET | `/products/search` | ProductsHandler.Search | Product search |
| GET | `/products/:category` | ProductsHandler.Category | Products by category |
| GET | `/products/:category/:slug` | ProductsHandler.Detail | Product detail |
| GET | `/solutions` | SolutionsHandler.List | Solutions listing |
| GET | `/solutions/:slug` | SolutionsHandler.Detail | Solution detail |
| GET | `/blog` | BlogHandler.List | Blog listing |
| GET | `/blog/:slug` | BlogHandler.Detail | Blog post |
| GET | `/case-studies` | CaseStudiesHandler.List | Case studies |
| GET | `/case-studies/:slug` | CaseStudiesHandler.Detail | Case study detail |
| GET | `/whitepapers` | WhitepapersHandler.List | Whitepapers |
| GET | `/whitepapers/:slug` | WhitepapersHandler.Detail | Whitepaper detail |
| POST | `/whitepapers/:slug/download` | WhitepapersHandler.Download | Download tracking |
| GET | `/about` | AboutHandler.Show | About page |
| GET | `/partners` | PartnersHandler.Show | Partners page |
| GET | `/contact` | ContactHandler.Show | Contact form |
| POST | `/contact/submit` | ContactHandler.Submit | Submit contact (rate limited) |
| GET | `/search` | SearchHandler.Search | Full-text search |
| GET | `/search/suggest` | SearchHandler.Suggest | HTMX autocomplete |
| GET | `/sitemap.xml` | SitemapHandler | XML sitemap |
| GET | `/robots.txt` | RobotsHandler | Robots file |

### Admin Routes (all require authentication)

| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/admin/login` | Authentication |
| POST | `/admin/logout` | Logout |
| GET | `/admin/dashboard` | Dashboard |
| GET/POST | `/admin/settings` | Global settings |
| GET/POST | `/admin/header` | Header settings |
| GET/POST | `/admin/footer` | Footer settings |
| CRUD | `/admin/products/*` | Product management |
| CRUD | `/admin/product-categories/*` | Category management |
| CRUD | `/admin/blog-posts/*` | Blog post management |
| CRUD | `/admin/blog-categories/*` | Blog categories |
| CRUD | `/admin/blog-authors/*` | Blog authors |
| CRUD | `/admin/blog-tags/*` | Blog tags |
| CRUD | `/admin/solutions/*` | Solutions management |
| CRUD | `/admin/case-studies/*` | Case studies |
| CRUD | `/admin/whitepapers/*` | Whitepapers |
| CRUD | `/admin/whitepaper-topics/*` | Whitepaper topics |
| CRUD | `/admin/partners/*` | Partners |
| CRUD | `/admin/partner-tiers/*` | Partner tiers |
| CRUD | `/admin/industries/*` | Industries |
| CRUD | `/admin/about/*` | About page sections |
| CRUD | `/admin/homepage/*` | Homepage sections |
| CRUD | `/admin/media/*` | Media library |
| CRUD | `/admin/navigation/*` | Navigation menus |
| GET | `/admin/activity` | Activity log |
| CRUD | `/admin/contact/*` | Contact submissions |

*CRUD = GET list, GET new, POST create, GET :id/edit, POST :id, DELETE :id*

---

## Configuration Reference

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENVIRONMENT` | (none) | Set to `production` for production mode |
| `DATABASE_PATH` | `bluejay.db` | Path to SQLite database file |

### SQLite Pragmas

| Pragma | Value | Purpose |
|--------|-------|---------|
| `journal_mode` | WAL | Concurrent read/write |
| `busy_timeout` | 5000 | Wait 5s on lock |
| `foreign_keys` | ON | Enforce relationships |
| `synchronous` | NORMAL | Balance safety/speed |
| `cache_size` | 2000 | Pages in memory |

### Session Configuration

| Setting | Value |
|---------|-------|
| MaxAge | 7 days |
| HttpOnly | true |
| SameSite | Lax |
| Secure | false (set to true in production) |
| Secret | Must be changed from default |

### Upload Limits

| Type | Max Size | Allowed Formats |
|------|----------|----------------|
| Product Images | 5 MB | jpg, jpeg, png, webp |
| Product Downloads | 50 MB | Any |

### Rate Limits

| Endpoint | Limit |
|----------|-------|
| `POST /contact/submit` | 5 requests/hour per IP |

---

*Generated for Bluejay CMS — a Go-powered content management system with brutalist design.*
