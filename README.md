# Bluejay CMS

A full-featured content management system built in Go for managing a technology company's website. Features a brutalist-styled admin panel and server-side rendering with HTMX for dynamic interactions.

## Screenshots

> Coming soon — add screenshots of admin panel and public pages here

## Features

- **Product Management** — Full product catalog with categories, specs, features, certifications, downloads, and image galleries
- **Blog System** — Rich text editing with Trix editor, categories, authors, tags, and related products
- **Solutions Pages** — Solution showcases with stats, challenges, related products, and CTAs
- **Case Studies** — Customer success stories with metrics and product references
- **Whitepapers** — Downloadable resources with download tracking and lead capture
- **Partners Directory** — Partner ecosystem with tiers and testimonials
- **Media Library** — Centralized asset management with search and metadata
- **Full-Text Search** — SQLite FTS5 search across products, blog posts, and solutions
- **Navigation Editor** — Drag-and-drop menu builder with nested structures
- **Activity Log** — Complete audit trail of admin actions
- **SEO Features** — XML sitemap generation, robots.txt, meta tags
- **Contact Forms** — Rate-limited contact submissions with office locations
- **Homepage Management** — Hero banners, stats counters, testimonials, and CTAs

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.25.5 | Backend runtime |
| Web Framework | Echo v4 | HTTP routing, middleware |
| Database | SQLite (modernc.org/sqlite) | Data storage, WAL mode |
| Query Gen | sqlc | Type-safe SQL → Go code |
| Migrations | golang-migrate | Schema versioning |
| Sessions | gorilla/sessions | Cookie-based authentication |
| Crypto | golang.org/x/crypto | bcrypt password hashing |
| Templates | Go html/template | Server-side rendering |
| Interactivity | HTMX | Dynamic HTML updates |
| Rich Text | Trix Editor | Blog content editing |
| CSS | Tailwind CSS (CDN) | Utility-first styling |
| Reverse Proxy | Caddy | TLS, static files, headers |
| DB Backup | Litestream | Continuous SQLite → S3 replication |

## Prerequisites

- **Go 1.25.5+** — [Download Go](https://go.dev/dl/)
- **sqlc** — Install: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- **SQLite3 CLI** — For seeding ([Download SQLite](https://sqlite.org/download.html))
- **air** (optional) — For hot-reload: `go install github.com/air-verse/air@latest`

## Quick Start

```bash
# Clone the repository
git clone <repo-url>
cd bluejay-cms

# Generate sqlc code from SQL queries
make sqlc

# Seed the database with sample data
make seed

# Run the server
make run
```

The server will start at **http://localhost:28090**

### Default Admin Credentials

After running `make seed`, you can login at `/admin/login` with:

- **Email:** `admin@bluejaylabs.com`
- **Password:** `admin123`

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Start the server (go run) |
| `make build` | Compile binary to `bin/bluejay-cms` |
| `make dev` | Start with hot-reload using air |
| `make sqlc` | Regenerate sqlc Go code from SQL queries |
| `make migrate-up` | Run all pending database migrations |
| `make migrate-down` | Rollback all migrations |
| `make seed` | Load sample data into database |
| `make test` | Run all tests (`go test -v ./...`) |
| `make clean` | Remove binaries and database files |
| `make deploy` | Full deploy: build, upload, restart |
| `make deploy-build` | Cross-compile for Linux (GOOS=linux GOARCH=amd64) |
| `make deploy-upload` | Upload binary and configs to server via SCP |
| `make deploy-restart` | Restart services on remote server |

## Project Structure

```
bluejay-cms/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── admin/               # Admin panel CRUD handlers (25 files)
│   │   └── public/              # Public-facing page handlers (13 files)
│   ├── middleware/              # Auth, logging, security, sessions, rate limiting
│   ├── services/                # Business logic (product, upload, cache, activity)
│   ├── models/                  # Domain models
│   ├── database/                # DB initialization, migrations runner
│   └── templates/               # Template rendering engine
├── db/
│   ├── migrations/              # 68 migration files (34 up + 34 down)
│   ├── queries/                 # 24 sqlc SQL query files
│   ├── sqlc/                    # Auto-generated Go code (DO NOT EDIT)
│   └── seeds/                   # 18 seed data files
├── templates/
│   ├── admin/                   # Admin panel templates
│   │   ├── layouts/             # Base admin layout
│   │   ├── pages/               # 60+ admin page templates
│   │   └── partials/            # Reusable HTMX fragments
│   └── public/                  # Public website templates
│       ├── layouts/             # Base public layout
│       ├── pages/               # 16 public page templates
│       └── partials/            # Search and shared fragments
├── public/
│   ├── css/                     # Stylesheets
│   ├── js/                      # htmx.min.js, trix.js, admin.js
│   └── uploads/                 # User-uploaded media
├── deploy/
│   ├── Caddyfile                # Reverse proxy configuration
│   ├── bluejay-cms.service      # systemd service file
│   ├── build.sh                 # Cross-compile script
│   └── litestream.yml           # SQLite backup to S3 config
├── Makefile                     # Build and deployment commands
├── sqlc.yaml                    # sqlc code generation config
└── seed.sql                     # Database seed data loader
```

## Development Workflow

### Adding a New SQL Query

1. Write your SQL in `db/queries/<entity>.sql`
2. Run `make sqlc` to regenerate Go code
3. Use the generated function in your handler: `h.queries.YourQuery(ctx, params)`

### Adding a New Migration

1. Create `db/migrations/NNN_description.up.sql` and `NNN_description.down.sql`
2. Run `make migrate-up` to apply the migration
3. Run `make sqlc` if the schema changed (to regenerate query code)

### Resetting the Database

```bash
make clean  # Removes database files
make seed   # Recreates database with sample data
```

## Architecture

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

### Brutalist Design System

The admin panel follows a brutalist design philosophy:

- **Font:** JetBrains Mono (monospace everywhere)
- **Borders:** 2px solid black, no border-radius (sharp corners)
- **Shadows:** Manual `box-shadow: 4px 4px 0px #000` (no drop-shadow utilities)
- **Buttons:** Uppercase text, thick borders, hover shifts shadow
- **Colors:** Black/white primary with accent colors per section

### Key Components

- **Handlers** — Echo route handlers that process requests
- **Services** — Business logic layer (ProductService, UploadService, Cache, ActivityLog)
- **Middleware** — Request processing chain (auth, logging, security headers, sessions)
- **Templates** — Server-rendered HTML with HTMX for dynamic updates
- **sqlc Queries** — Type-safe, compile-time checked SQL queries

## Deployment

### Production Architecture

```
Internet → Caddy (443/HTTPS) → Go Binary (28090) → SQLite
                                                      ↓
                                                  Litestream → S3
```

### Basic Deployment Steps

1. **Build for Linux:**
   ```bash
   make deploy-build
   ```

2. **Prepare Server:**
   ```bash
   ssh user@yourserver
   sudo mkdir -p /var/www/bluejay-cms
   sudo chown www-data:www-data /var/www/bluejay-cms
   ```

3. **Upload Files:**
   ```bash
   make deploy-upload
   ```

4. **Start Services:**
   ```bash
   ssh user@yourserver
   sudo systemctl daemon-reload
   sudo systemctl enable bluejay-cms
   sudo systemctl start bluejay-cms
   sudo systemctl enable caddy
   sudo systemctl start caddy
   ```

5. **Verify:**
   ```bash
   curl https://yourdomain.com/health
   ```

### Continuous Deployment

```bash
make deploy  # One command: build → upload → restart
```

### Backup & Recovery

Litestream continuously streams SQLite WAL changes to S3.

**Restore from backup:**
```bash
litestream restore -o /var/www/bluejay-cms/bluejay.db s3://your-backup-bucket/bluejay-cms
```

## Documentation

For detailed information, see:

- **[DOCUMENTATION.md](DOCUMENTATION.md)** — Complete technical documentation
- **[architecture-and-philosophy.md](architecture-and-philosophy.md)** — Design philosophy and implementation principles
- **[CLAUDE.md](CLAUDE.md)** — Project instructions and conventions

Additional documentation:

- **Database Schema** — See `db/migrations/` for all tables and relationships
- **API Routes** — Full route reference in DOCUMENTATION.md
- **Deployment Guide** — Production deployment steps in DOCUMENTATION.md
- **Development Setup** — Local development workflow in DOCUMENTATION.md

## License

> Add your license information here

---

**Bluejay CMS** — A Go-powered content management system with brutalist design, built for simplicity and performance.
