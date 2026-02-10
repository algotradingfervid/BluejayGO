# Architecture & Implementation Philosophy

---

## Architecture

Everything runs as a **single Go binary**. No containers, no separate database server, no frontend build pipeline in production.

```
┌─────────────────────────────────────────────────────┐
│                  Single Go Binary                    │
│                                                      │
│  ┌─────────┐   ┌──────────┐   ┌──────────────────┐  │
│  │  Caddy   │──▶│  Echo    │──▶│  Middleware       │  │
│  │ (reverse │   │  Router  │   │  - Recovery       │  │
│  │  proxy)  │   │          │   │  - Logging (slog) │  │
│  └─────────┘   └──────────┘   │  - Auth (sessions)│  │
│                                │  - Cache          │  │
│                                └────────┬─────────┘  │
│                          ┌──────────────┼──────┐     │
│                          ▼              ▼      │     │
│                   ┌────────────┐ ┌───────────┐ │     │
│                   │   Public   │ │   Admin   │ │     │
│                   │  Handlers  │ │  Handlers │ │     │
│                   └─────┬──────┘ └─────┬─────┘ │     │
│                         ▼              ▼       │     │
│                   ┌─────────────────────────┐  │     │
│                   │   Go html/template      │  │     │
│                   │   (server-rendered)      │  │     │
│                   └─────────────────────────┘  │     │
│                         │              │       │     │
│                         ▼              ▼       │     │
│                   ┌─────────────────────────┐  │     │
│                   │   sqlc-generated code   │  │     │
│                   │   (type-safe queries)   │  │     │
│                   └───────────┬─────────────┘  │     │
│                               ▼                │     │
│                   ┌─────────────────────────┐  │     │
│                   │   SQLite (modernc.org)  │  │     │
│                   │   + Local Filesystem    │  │     │
│                   └─────────────────────────┘  │     │
└─────────────────────────────────────────────────────┘
```

### Request Flow

**Public visitor** hits `/products/widget-pro`:

1. Caddy terminates TLS, proxies to Go binary on `:8090`
2. Echo router matches `handlers.ProductDetail`
3. Cache middleware checks in-memory `sync.Map` — if hit, return cached HTML (<1ms)
4. On miss: handler calls `queries.GetProductBySlug(ctx, "widget-pro")` — sqlc-generated, type-safe
5. Handler passes the product struct to `templates/public/product-detail.html`
6. Go renders HTML, stores in cache, responds — typically <5ms

**Content editor** saves a product:

1. Auth middleware checks gorilla/sessions cookie → redirect to `/admin/login` if missing
2. HTMX sends `POST /admin/products/save` with form data
3. Handler validates input, calls `queries.UpdateProduct(ctx, params)`
4. Handler invalidates related cache keys (`/products`, `/products/{slug}`, `/`)
5. Handler returns a toast HTML fragment — HTMX swaps it in, no full reload

### What This Means in Practice

Every layer is a standard Go construct. Handlers are functions that take `*sql.DB` and return `echo.HandlerFunc`. Templates are `html/template`. Queries are plain SQL that sqlc compiles to Go functions. When something breaks, you look at the SQL, the handler, or the template. There is no framework runtime to debug, no ORM query builder to reverse-engineer, no abstraction layer translating between your intent and the database.

### Dedicated Tables, Not JSON Blocks

Your website has known, fixed page structures. The homepage always has hero → features → testimonials → CTA. Product pages always have name, description, images, specs. This structure is designed once and rarely changes — only the content within it changes.

This means dedicated SQL tables, not a generic block/JSON system. You get SQL queryability (`WHERE category_id = 3`), type safety at the database level, straightforward forms (one field = one input), and easy migrations (`ALTER TABLE` instead of JSON walks). Editors can't accidentally delete structural sections because the structure lives in the schema, not in the content.

---

## Implementation Philosophy

### 1. Work Vertically

Build one content type end-to-end before starting the next:

```
SQL migration → sqlc queries → handler → template → admin form → seed data
```

Don't build all migrations, then all handlers, then all templates. Finish products completely — public page, admin CRUD, image upload, cache invalidation — before touching team members. You'll have something demoable after day 2 and you'll discover real problems early instead of theoretical ones.

### 2. Don't Abstract Early

When you notice your admin handlers for products, team members, and testimonials look 80% the same, resist the urge to build a generic CRUD framework. Write the "boring" repetitive code for 3-4 content types first. Then — and only then — you'll know which parts actually benefit from abstraction and which just need to stay simple and explicit. Premature abstractions in a CMS tend to create worse problems than repetition does.

### 3. Cache Aggressively, Invalidate Precisely

A company website is read thousands of times for every one write. Public pages should almost always be served from an in-memory `sync.Map` — no database hit, no template rendering. Only touch the database when an editor saves something.

Cache key is the URL path. TTL is 1 hour as a safety net, but the real mechanism is **precise invalidation on save**: saving a product invalidates `/products`, `/products/{slug}`, and `/` (homepage). This gives you sub-millisecond response times for visitors and instant content updates for editors.

### 4. Own Your Data Layer

Write plain SQL. Let sqlc generate the Go code. Don't put an ORM or a framework DAO between you and your queries. When you need a complex join or a custom aggregation, you just write SQL — no fighting with a query builder. When you need to debug a slow query, you `EXPLAIN` it directly. The total SQL for this CMS will be maybe 40-50 queries. That's perfectly manageable by hand, and sqlc ensures they're type-checked at build time.

### 5. Boring Dependencies, Stable Foundations

Every dependency in this stack is either Go stdlib or a small, focused library that does one thing:

- **modernc.org/sqlite** — pure Go SQLite driver, nothing else
- **sqlc** — SQL → Go codegen, nothing else
- **golang-migrate** — run migration files, nothing else
- **gorilla/sessions** — session cookies, nothing else
- **Echo** — HTTP routing + middleware, nothing else

No dependency ships with features you don't use. No dependency has its own admin panel, its own plugin system, or its own opinions about how your CMS should work. When one of these has a breaking change, the blast radius is small and the fix is obvious.

### 6. Admin Panel: Functional Over Beautiful

Your public site is what visitors see — invest design effort there. The admin panel needs to be clear and usable, not pixel-perfect. Tailwind utility classes on standard HTML form elements are more than enough. A `<select>` dropdown, a text `<input>`, a Trix editor, and a Save button. If it's obvious what each field does and the save works reliably, the admin panel is done.

### 7. Vendor Your JS Dependencies

Copy `htmx.min.js` and `trix.js` into `/public/js/`. Don't use CDNs. Your site has zero external runtime dependencies and works even if every CDN goes down. Two JS files, both vendored, both under 50KB gzipped combined. That's the entire client-side footprint.

### 8. Seed Data Is Your Safety Net

Maintain a `seed.sql` that recreates all your content from scratch. You should be able to delete the database and get back to a known-good state in seconds (`make seed`). This makes development fearless — you can experiment with schema changes, test destructive operations, and always reset cleanly.

### 9. Deploy Like It's 2005 (With Modern TLS)

The deployment is: `scp` a binary to a VPS, restart a systemd service. Caddy handles HTTPS automatically. Litestream handles backups continuously. There's no Kubernetes, no Docker, no orchestration layer. A single €4/month Hetzner VPS will serve this site at thousands of requests per second. If the server dies, you provision a new one, restore from Litestream, and you're back online in minutes.

### 10. Complexity Budget

Every piece of complexity must justify itself against the alternative of "just write plain Go code." If a tool saves less time than it takes to learn, configure, debug, and maintain — skip it. This CMS serves a small team editing a company website. The architecture should reflect that scope, not the scope of a SaaS platform you might build someday.
