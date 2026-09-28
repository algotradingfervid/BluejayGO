---
name: tdd-bug-fix
description: Use when fixing a bug, regression, or "X is broken/wrong/returns nothing/500s" report in the Bluejay CMS (Go + Echo + HTMX + sqlc) — before editing any handler, service, query, or template.
---

# TDD Bug Fix (Bluejay CMS)

## Overview

A bug is not fixed until a test that **failed because of the bug now passes because of your fix**. Encode the bug as a failing test FIRST, against the real in-process e2e harness, then make it pass with the smallest change.

**The Iron Law:** No fix without a failing test first. If you wrote the fix before the test, `git stash` (or delete) the fix, write the test, watch it fail for the right reason, then restore the fix.

**Violating the letter of this rule is violating its spirit.** A test written after the fix that passes immediately proves nothing — it never demonstrated it can catch the bug.

## The Loop (RED → GREEN → VERIFY)

1. **Reproduce** the bug as a user/HTTP request. Confirm the broken behavior with your own eyes before theorizing about cause. See superpowers:systematic-debugging if the cause is unclear.
2. **RED** — Write a test that asserts the *correct* behavior. Run it. It MUST fail, and the failure message must describe the actual bug (wrong status, wrong count, wrong body), not a compile error or typo.
3. **GREEN** — Make the smallest change that turns the test green. Don't refactor unrelated code (per CLAUDE.md scope rules).
4. **VERIFY** — Run the gates below. All must pass before you claim done.

## Where the test goes

| Bug surface | Test location | Pattern |
|---|---|---|
| Route / status / redirect / auth / DB persistence | `internal/e2e/NN_<area>_test.go` | Use `setupApp(t)`, `createTestAdmin(t, queries)`, `loginAndGetCookie(t, e)` |
| Handler logic in isolation | `internal/handlers/{admin,public}/<name>_test.go` | Follow the existing `_test.go` in that package |
| Service / business logic (slugs, cache, uploads) | `internal/services/<name>_test.go` | Pure unit test |
| Template **output** (not just routing) | e2e test with a **real** renderer | The shared `setupApp` uses a stub renderer that emits `<html>stub</html>` — it cannot see template bugs. Build a local Echo instance with `templates.NewRenderer("templates")` and assert on `rec.Body.String()`. |

The e2e harness (`internal/e2e/e2e_test.go`) spins up every route against a **fresh temp SQLite DB** with migrations applied — no server needed. Seed data via `queries.Create...`, drive via `httptest`, assert on `rec.Code` / `rec.Body` / DB state. Numbered files mirror `tests/plans/NN-*.md` — read the matching plan for expected behavior.

## Verification gates (ALL required before "done")

```bash
go build ./cmd/...          # phase-runner rule: must compile (use ./... to catch test files)
go test ./internal/e2e/...  # the package you added the test to
go test ./...               # nothing else regressed
```

If you touched `db/queries/*.sql` or `db/schema`, run `sqlc generate` BEFORE building (CLAUDE.md rule), and commit the regenerated `db/sqlc/*.go`.

For a template/visual fix, also run the app and look at the page — a 200 with broken markup is still a bug:
```bash
./restart.sh    # kills :28090, rebuilds, runs; or: go run cmd/server/main.go
curl -s localhost:28090/<path> | grep <expected>
```
Template fixes must keep the brutalist system (CLAUDE.md): 2px solid black borders, `box-shadow: 4px 4px 0px #000`, NO border-radius, JetBrains Mono, uppercase buttons.

## Example: "DELETE on a missing product 500s instead of 404"

```go
// internal/e2e/09_products_extended_test.go
func TestAdminProductDelete_NotFound_Returns404(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	// No product with id 99999 exists in the fresh DB.
	req := httptest.NewRequest(http.MethodDelete, "/admin/products/99999", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// RED: bug returns 500; correct behavior is 404.
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing product, got %d; body: %s", rec.Code, rec.Body.String())
	}
}
```
Run `go test ./internal/e2e/ -run TestAdminProductDelete_NotFound_Returns404` → watch it fail with `got 500` → fix the handler to map `sql.ErrNoRows` to 404 → watch it pass → run the full gates.

## Red Flags — STOP, you're rationalizing

- "I can see the fix, the test is a formality" → write it first anyway; it's the only proof the bug is real and fixed.
- "I'll add the test after the fix compiles" → tests-after never demonstrate they catch the bug. RED first.
- "The stub renderer test passes, ship it" → stub renderer can't see template bugs. Use a real renderer for output assertions.
- "Build passes, done" → build is not test. Run `go test ./...`.
- "I'll also clean up this nearby code" → out of scope (CLAUDE.md). One bug, one minimal fix.
- "Changed a .sql query, build is green" → did you run `sqlc generate` and rebuild? Stale generated code lies.

## Rationalization table

| Excuse | Reality |
|---|---|
| "Too simple to need a test" | Simple handlers return the wrong status all the time. The test is 10 lines. |
| "Test passes immediately after I write it" | Then it never saw the bug. Make it fail first (revert the fix), confirm the message, restore. |
| "e2e is slow, I'll skip it" | `go test ./internal/e2e/ -run TestName` runs one test in seconds. |
| "It's a template, hard to test" | Real-renderer e2e + `rec.Body.String()` assertion. The example shows how. |
| "Manual curl confirmed it" | Manual checks vanish. A committed failing→passing test is the regression guard. |
