# Inline Edit + Numeric Reordering for Product Sub-Details

**Date:** 2026-06-24
**Status:** Approved (design) — pending spec review
**Author:** Naren (with Claude)

## Problem

In the Bluejay CMS admin panel, a product's core fields (name, description, price,
category, image) can be edited via the existing `GET /admin/products/:id/edit` form.
But the five **sub-entity** tabs on the product form are **add/delete only**:

- Features
- Specs
- Certifications
- Downloads
- Images

There is no way to edit an existing sub-item's text/value or to change its sort
order. Every table already has a `display_order INTEGER NOT NULL DEFAULT 0` column
(rendered read-only as `#N`), and the list queries already
`ORDER BY display_order ASC` — but there is **no `Update*` query, handler, or route**
for any sub-entity. The only way to fix a typo or re-order today is to delete the
item and re-add it.

## Goal

Let admins edit each sub-item's details **and** change its numeric sort order, in
place, without deleting and re-adding. Reordering is done by editing a row's numeric
`display_order` (per the request: "the sorting order will be defined by numbers").

## Scope

In scope: inline edit (editable text/value fields + numeric `display_order`) for all
five sub-entities listed above.

Out of scope (explicit decisions):

1. **File replacement** for Downloads/Images. Editing changes metadata + order only;
   the uploaded file/image itself is unchanged. Replacing the actual file still uses
   delete + re-add. (Keeps the edit form simple — no multipart upload path.)
2. **`is_thumbnail` ("primary") flag** is not editable through the edit form. A plain
   `UPDATE` cannot set one image primary without un-setting the others, which would
   create two primaries. The update preserves the existing `is_thumbnail` value.
   The pre-existing dead "Set Primary" button (no route/handler exists) is a separate
   issue and remains out of scope.
3. **Cache invalidation is unchanged.** The existing Add/Delete sub-entity handlers do
   not invalidate the page cache (only the product-core handlers do), and
   `ProductDetailsHandler` has no `*services.Cache` dependency. The new Update handlers
   match existing behavior — no cache change.
4. No drag-and-drop reordering. Ordering is via the numeric field only.
5. No changes to the public-facing product pages.

## Approach

**Server-rendered inline-edit toggle via an `?edit=<id>` query param on the existing
List route.** Chosen over (a) a separate `GET .../edit` endpoint per row, and (b) a
client-side JS toggle rendering both views in the DOM. The chosen approach adds only
**one new route per entity** (the POST update), reuses the existing
`renderPartial` / `List*` plumbing, needs **zero custom JavaScript**, and keeps exactly
one row editable at a time.

### Interaction flow (per entity, Features shown as example)

1. Each display row gains an **Edit (✎)** button:
   `hx-get="/admin/products/{{$.ProductID}}/features?edit={{.ID}}"`,
   `hx-target="#features-section"`, `hx-swap="outerHTML"`.
2. `ListFeatures` reads `c.QueryParam("edit")` → `EditingID int64` (0 = none) and
   passes it into the partial data.
3. The partial renders that one row as an **inline edit form** when
   `{{if eq .ID $.EditingID}}`; all other rows render normally. The form contains the
   editable text fields + a numeric **Order** input pre-filled from the row, plus
   **SAVE** and **CANCEL**.
4. **SAVE** → `hx-post="/admin/products/{{$.ProductID}}/features/{{.ID}}"` → new
   `UpdateFeature` handler applies the update, then returns `h.ListFeatures(c)` (with
   no `edit` param), re-rendering the whole section — now sorted by the new
   `display_order`.
5. **CANCEL** → `hx-get="/admin/products/{{$.ProductID}}/features"` (plain list,
   `EditingID = 0`).

Because `EditingID` defaults to 0 and real ids are ≥ 1, the non-editing render is
unchanged.

### Why this is safe in Echo routing

Existing single-item route is `DELETE /products/:id/features/:feature_id`. Adding
`POST /products/:id/features/:feature_id` (different HTTP method, same path) does not
conflict. The List route `GET /products/:id/features` already exists; only its handler
changes to read `?edit`.

## Per-entity change list

All sub-entities live in the same set of files:

- DB schema: `db/migrations/009_products.up.sql`
- sqlc SQL source: `db/queries/products.sql`
- sqlc generated: `db/sqlc/products.sql.go`, `db/sqlc/models.go`, `db/sqlc/querier.go`
- Handlers: `internal/handlers/admin/product_details.go`
- Routes: `cmd/server/main.go` (block ~350–378, under `adminGroup` = `/admin`)
- Templates: `templates/admin/partials/product_*.html`

For **each** sub-entity:

| Layer | Change |
|---|---|
| `db/queries/products.sql` | Add `UpdateProduct<X> :exec` setting all editable columns `WHERE id = ?`; run `sqlc generate` |
| `product_details.go` | `List<X>`: read `?edit` → `EditingID`, add to partial data. New `Update<X>` handler (POST): parse form, call `UpdateProduct<X>`, return `h.List<X>(c)` |
| `cmd/server/main.go` | One new route: `POST /products/:id/<entity>/:<x>_id` → `pdHandler.Update<X>` |
| `product_<x>.html` | Add ✎ Edit button per row + conditional inline edit `<form>` matching brutalist style; all HTMX swaps target the section container `outerHTML` |

### Editable fields, route params, container ids

| Entity | Editable fields | Single-item route param | Section container id |
|---|---|---|---|
| Features | `feature_text`, `display_order` | `:feature_id` | `#features-section` |
| Specs | `section_name`, `spec_key`, `spec_value`, `display_order` | `:spec_id` | `#specs-section` |
| Certifications | `certification_name`, `certification_code`, `icon_type`, `icon_path`, `display_order` | `:cert_id` | `#certifications-section` |
| Downloads | `title`, `description`, `file_type`, `version`, `display_order` | `:download_id` | `#downloads-section` |
| Images | `alt_text`, `caption`, `display_order` | `:image_id` | `#images-section` |

Nullable columns (`sql.NullString` / `sql.NullInt64`) are written exactly as the
existing Add handlers do: `sql.NullString{String: v, Valid: v != ""}`. Downloads/Images
update handlers are **metadata-only** — no `c.FormFile`, no multipart, file columns
(`file_path`, `file_size`, `image_path`, `download_count`, `is_thumbnail`) are not
touched by the `UPDATE`.

### New sqlc queries (exact shapes)

```sql
-- name: UpdateProductFeature :exec
UPDATE product_features SET feature_text = ?, display_order = ? WHERE id = ?;

-- name: UpdateProductSpec :exec
UPDATE product_specs SET section_name = ?, spec_key = ?, spec_value = ?, display_order = ? WHERE id = ?;

-- name: UpdateProductCertification :exec
UPDATE product_certifications
SET certification_name = ?, certification_code = ?, icon_type = ?, icon_path = ?, display_order = ?
WHERE id = ?;

-- name: UpdateProductDownload :exec
UPDATE product_downloads
SET title = ?, description = ?, file_type = ?, version = ?, display_order = ?
WHERE id = ?;

-- name: UpdateProductImage :exec
UPDATE product_images SET alt_text = ?, caption = ?, display_order = ? WHERE id = ?;
```

(`product_downloads` has an `update_product_downloads_timestamp` trigger that bumps
`updated_at` automatically on UPDATE — no extra handling needed.)

### Handler shape (Features example, others analogous)

```go
// List<X> gains EditingID
func (h *ProductDetailsHandler) ListFeatures(c echo.Context) error {
    // ... existing fetch ...
    editingID, _ := strconv.ParseInt(c.QueryParam("edit"), 10, 64) // 0 when absent
    return h.renderPartial(c, "product_features", map[string]any{
        "ProductID": id,
        "Features":  features,
        "EditingID": editingID,
    })
}

// Update<X> is new — mirrors AddFeature, then re-lists
func (h *ProductDetailsHandler) UpdateFeature(c echo.Context) error {
    ctx := c.Request().Context()
    featureID, _ := strconv.ParseInt(c.Param("feature_id"), 10, 64)
    order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)
    if err := h.queries.UpdateProductFeature(ctx, sqlc.UpdateProductFeatureParams{
        FeatureText:  c.FormValue("feature_text"),
        DisplayOrder: order,
        ID:           featureID,
    }); err != nil {
        return err // match existing error handling in this file
    }
    return h.ListFeatures(c) // re-render section, EditingID resets to 0
}
```

(Validation parity with the existing Add handlers: minimal — required text fields are
marked `required` in the form; `display_order` parses to 0 on bad input.)

### Template shape (Features example, others analogous)

Within the existing `{{range .Features}}` loop, branch on the editing id:

```html
{{range .Features}}
  {{if eq .ID $.EditingID}}
    <form hx-post="/admin/products/{{$.ProductID}}/features/{{.ID}}"
          hx-target="#features-section" hx-swap="outerHTML"
          class="border-2 border-black px-4 py-3 bg-yellow-50 ..." style="box-shadow: 4px 4px 0px #000;">
      <input type="text" name="feature_text" value="{{.FeatureText}}" required class="...">
      <input type="number" name="display_order" value="{{.DisplayOrder}}" class="...">
      <button type="submit" class="... uppercase ...">Save</button>
      <a hx-get="/admin/products/{{$.ProductID}}/features"
         hx-target="#features-section" hx-swap="outerHTML" class="... uppercase ...">Cancel</a>
    </form>
  {{else}}
    <!-- existing display row, plus an Edit (✎) button: -->
    <button hx-get="/admin/products/{{$.ProductID}}/features?edit={{.ID}}"
            hx-target="#features-section" hx-swap="outerHTML"
            class="... text-xs font-bold uppercase">&#9998;</button>
  {{end}}
{{end}}
```

Styling follows the brutalist system already in these partials: 2px solid black
borders, `box-shadow: Npx Npx 0px #000`, JetBrains Mono, uppercase button text, no
border-radius. The edit form reuses the same input field names as the entity's Add
form so handler form-reads are consistent.

## Testing (TDD)

Follow the existing `internal/e2e/*_test.go` pattern (numbered files). For each
sub-entity, add a test that:

1. Seeds a product and one sub-item (with a known `display_order`).
2. `GET /admin/products/:id/<entity>?edit=<itemID>` → assert the response contains the
   inline edit form pre-filled with the item's current values.
3. `POST /admin/products/:id/<entity>/:<x>_id` with changed text + a different
   `display_order` → assert 200 and that the returned fragment / a follow-up list shows
   the updated value.
4. With two sub-items, set orders so they swap, then assert the list renders in the new
   `display_order ASC` sequence (reordering works).

Plus, after code changes:

- `sqlc generate` (regenerate after editing `db/queries/products.sql`).
- `go build ./cmd/...` must pass.
- Run the relevant e2e tests.

Write the failing test first (red), implement (green), then refactor — per entity.

## Risks / notes

- **Specs rows** don't currently display `#DisplayOrder`; the edit form will still
  expose the Order field. Acceptable (consistent with the other entities in edit mode).
- **Pre-existing dead "Set Primary" button** in `product_images.html` is left as-is
  (out of scope).
- **Cache staleness** on public product pages already exists for add/delete of
  sub-items (no invalidation); edit inherits the same behavior. If this becomes a
  problem, the follow-up is to inject `*services.Cache` into `ProductDetailsHandler`
  and invalidate `page:products` on all sub-entity mutations.
- Five near-identical implementations: per project convention
  (`CLAUDE.md`: "Follow existing code patterns — don't refactor unrelated code") we keep
  explicit per-entity handlers rather than introducing a generic abstraction, matching
  the existing `Add*`/`Delete*` style in `product_details.go`.
