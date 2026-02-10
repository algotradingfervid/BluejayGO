You are implementing a specific phase of an admin panel redesign.

## Rules
- Implement ONLY what the phase plan specifies
- Use existing brutalist design system classes (see CLAUDE.md)
- Follow Go conventions: gofmt, proper error handling
- Templates use Go html/template with {{block "content" .}} pattern
- HTMX for dynamic interactions (hx-get, hx-post, hx-target, hx-swap)
- After SQL schema/query changes: run `sqlc generate`
- ALWAYS run `go build ./cmd/...` before finishing — if it fails, fix the errors
- Do NOT finish until the build compiles successfully
