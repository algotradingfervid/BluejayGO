# Bluejay CMS — Admin Panel Redesign

## Tech Stack
- **Backend**: Go 1.25, Echo v4, SQLite (modernc.org/sqlite)
- **Frontend**: Go html/template, HTMX, Tailwind CSS (CDN), Trix editor
- **DB Queries**: sqlc (config: sqlc.yaml)
- **Module**: github.com/narendhupati/bluejay-cms

## Brutalist Design System
- Font: JetBrains Mono (monospace everywhere)
- Borders: 2px solid black
- Shadows: manual box-shadow (4px 4px 0px #000), NO drop-shadow utilities
- NO border-radius anywhere — everything is sharp corners
- Colors: black/white primary, accent colors per section
- Buttons: uppercase text, thick borders, hover shifts shadow

## Directory Structure
```
cmd/server/          — main.go entry point
internal/
  handlers/          — Echo route handlers
  middleware/        — Auth, logging middleware
  models/            — Data structures
  services/          — Business logic
  database/          — DB connection, migrations
  templates/         — Template renderer
templates/admin/
  layouts/           — Base layouts (admin-layout.html)
  pages/             — Full page templates
  partials/          — Reusable components (sidebar, header, etc.)
db/queries/          — sqlc SQL query files
public/              — Static assets (CSS, JS, images)
```

## Conventions
- Handlers: `func (h *Handler) GetProducts(c echo.Context) error`
- Templates: `{{template "admin-layout" .}}` with `{{block "content" .}}...{{end}}`
- HTMX endpoints return HTML fragments (partials), not full pages
- Routes: `/admin/...` for all admin pages
- Form handlers: GET renders form, POST processes submission

## Current Phase
Check `automation/phase-tracker.json` for the current phase being implemented.

## Rules
- Only modify files within scope of the current phase plan
- Run `go build ./cmd/...` to verify compilation after changes
- Run `sqlc generate` after any SQL query file changes
- Follow existing code patterns — don't refactor unrelated code
