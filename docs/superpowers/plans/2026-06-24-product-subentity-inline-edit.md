# Product Sub-Entity Inline Edit Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let admins edit each product sub-item (features, specs, certifications, downloads, images) in place — text/value fields plus a numeric sort order — instead of delete-and-re-add.

**Architecture:** Server-rendered inline-edit toggle. Each list partial gains a per-row Edit (✎) button that does `hx-get=".../<entity>?edit=<id>"`. The existing `List<X>` handler reads the `edit` query param into an `EditingID` and the partial renders that one row as an inline `<form>`. SAVE posts to a new `POST /products/:id/<entity>/:<x>_id` route → `Update<X>` handler → re-renders the section (now re-sorted by `display_order`). CANCEL re-fetches the plain list. No custom JavaScript; one new route per entity.

**Tech Stack:** Go 1.25, Echo v4, SQLite (modernc.org/sqlite), sqlc, Go html/template, HTMX, Tailwind (brutalist design system).

## Global Constraints

- Module path: `github.com/narendhupati/bluejay-cms`. sqlc-generated package import: `github.com/narendhupati/bluejay-cms/db/sqlc` (aliased `sqlc`).
- After editing any `db/queries/*.sql`, run `sqlc generate` (binary at `/Users/narendhupati/go/bin/sqlc`; `make sqlc` also works). Never hand-edit `db/sqlc/*.go`.
- `go build ./cmd/...` must pass after every code change. Full suite: `go test ./...` (or `make test`).
- Routes exist in TWO places that must stay in sync: production `cmd/server/main.go` and the e2e harness `internal/e2e/e2e_test.go` (`setupApp`). Add every new route to BOTH.
- Brutalist design system: monospace (JetBrains Mono), 2px solid black borders, manual `box-shadow: Npx Npx 0px #000`, NO border-radius, uppercase button text. Edit forms must match the surrounding partial markup.
- Out of scope (do NOT implement): file/image replacement on downloads/images edit (metadata + order only); the `is_thumbnail` "primary" flag in the image edit form (preserve existing value, leave the pre-existing dead "Set Primary" button as-is); cache invalidation changes.
- Reuse the exact form field `name`s that the existing `Add*` forms use, so handler `c.FormValue(...)` reads stay consistent.
- `EditingID` is `int64`, defaults to `0`. Real row ids are ≥ 1, so `{{if eq .ID $.EditingID}}` is false for every row when not editing.

---

### Task 1: Add the five `UpdateProduct<X>` sqlc queries and regenerate

**Files:**
- Modify: `db/queries/products.sql` (append five queries)
- Regenerate (do not hand-edit): `db/sqlc/products.sql.go`, `db/sqlc/querier.go`

**Interfaces:**
- Consumes: nothing (first task).
- Produces (generated methods + param structs consumed by Tasks 2–6):
  - `UpdateProductFeature(ctx, UpdateProductFeatureParams{FeatureText string, DisplayOrder int64, ID int64}) error`
  - `UpdateProductSpec(ctx, UpdateProductSpecParams{SectionName string, SpecKey string, SpecValue string, DisplayOrder int64, ID int64}) error`
  - `UpdateProductCertification(ctx, UpdateProductCertificationParams{CertificationName string, CertificationCode sql.NullString, IconType sql.NullString, IconPath sql.NullString, DisplayOrder int64, ID int64}) error`
  - `UpdateProductDownload(ctx, UpdateProductDownloadParams{Title string, Description sql.NullString, FileType string, Version sql.NullString, DisplayOrder int64, ID int64}) error`
  - `UpdateProductImage(ctx, UpdateProductImageParams{AltText sql.NullString, Caption sql.NullString, DisplayOrder int64, ID int64}) error`

- [ ] **Step 1: Append the five Update queries to `db/queries/products.sql`**

Add this block at the END of `db/queries/products.sql`:

```sql
-- name: UpdateProductFeature :exec
UPDATE product_features
SET feature_text = ?, display_order = ?
WHERE id = ?;

-- name: UpdateProductSpec :exec
UPDATE product_specs
SET section_name = ?, spec_key = ?, spec_value = ?, display_order = ?
WHERE id = ?;

-- name: UpdateProductCertification :exec
UPDATE product_certifications
SET certification_name = ?, certification_code = ?, icon_type = ?, icon_path = ?, display_order = ?
WHERE id = ?;

-- name: UpdateProductDownload :exec
UPDATE product_downloads
SET title = ?, description = ?, file_type = ?, version = ?, display_order = ?
WHERE id = ?;

-- name: UpdateProductImage :exec
UPDATE product_images
SET alt_text = ?, caption = ?, display_order = ?
WHERE id = ?;
```

- [ ] **Step 2: Regenerate sqlc**

Run: `cd /Users/narendhupati/Documents/ClaudeWebsiteCreator && sqlc generate`
Expected: no output, exit 0. (If `sqlc` is not on PATH, use `/Users/narendhupati/go/bin/sqlc generate` or `make sqlc`.)

- [ ] **Step 3: Verify the generated methods exist**

Run: `grep -E "func \(q \*Queries\) UpdateProduct(Feature|Spec|Certification|Download|Image)\b" db/sqlc/products.sql.go`
Expected: exactly 5 matching lines (one per method).

- [ ] **Step 4: Verify the build compiles**

Run: `go build ./cmd/...`
Expected: no output, exit 0.

- [ ] **Step 5: Commit**

```bash
git add db/queries/products.sql db/sqlc/products.sql.go db/sqlc/querier.go
git commit -m "feat(products): add UpdateProduct* sqlc queries for sub-entities"
```

---

### Task 2: Features — inline edit + reorder

**Files:**
- Modify: `internal/handlers/admin/product_details.go` (edit `ListFeatures`, add `UpdateFeature`)
- Modify: `cmd/server/main.go` (add route after line 362)
- Modify: `internal/e2e/e2e_test.go` (add route after the `DeleteFeature` line ~296)
- Modify: `templates/admin/partials/product_features.html` (row loop, lines 30–41)
- Test: `internal/e2e/11_product_features_test.go` (append two tests)

**Interfaces:**
- Consumes: `UpdateProductFeature` / `UpdateProductFeatureParams` (Task 1).
- Produces: handler method `UpdateFeature(c echo.Context) error`; route `POST /admin/products/:id/features/:feature_id`.

- [ ] **Step 1: Write the failing tests** — append to `internal/e2e/11_product_features_test.go` (its import block already has `context, fmt, net/http, net/http/httptest, net/url, strings, testing, sqlc` — no import change needed):

```go
func TestProductFeaturesEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "FeatEditCat", Slug: "feat-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "FEAT-EDIT-1", Slug: "feature-edit-form", Name: "Edit Form Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	feat, _ := queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: product.ID, FeatureText: "Original feature text", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/features?edit=%d", product.ID, feat.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/features/%d"`, product.ID, feat.ID)) {
		t.Errorf("expected inline edit form for feature %d, body: %s", feat.ID, body)
	}
	if !strings.Contains(body, `value="Original feature text"`) {
		t.Errorf("expected pre-filled feature text, body: %s", body)
	}
}

func TestProductFeaturesUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "FeatUpdCat", Slug: "feat-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "FEAT-UPD-1", Slug: "feature-update", Name: "Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	f1, _ := queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: product.ID, FeatureText: "Alpha", DisplayOrder: 1,
	})
	queries.CreateProductFeature(ctx, sqlc.CreateProductFeatureParams{
		ProductID: product.ID, FeatureText: "Bravo", DisplayOrder: 2,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/features/%d", product.ID, f1.ID), strings.NewReader(url.Values{
		"feature_text":  {"Alpha Updated"},
		"display_order": {"5"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	features, _ := queries.ListProductFeatures(ctx, product.ID)
	if len(features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(features))
	}
	if features[0].FeatureText != "Bravo" {
		t.Errorf("expected 'Bravo' first after reorder, got %q", features[0].FeatureText)
	}
	if features[1].FeatureText != "Alpha Updated" {
		t.Errorf("expected 'Alpha Updated' second, got %q", features[1].FeatureText)
	}
	if features[1].DisplayOrder != 5 {
		t.Errorf("expected display_order 5, got %d", features[1].DisplayOrder)
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/e2e/ -run 'TestProductFeatures(EditForm|Update)_E2E' -v`
Expected: FAIL — `TestProductFeaturesUpdate_E2E` gets 405 (no POST route yet); `TestProductFeaturesEditForm_E2E` gets 200 but body lacks the edit form.

- [ ] **Step 3: Edit `ListFeatures` to read `?edit`** in `internal/handlers/admin/product_details.go`. Replace the body of `ListFeatures` (currently lines ~202–217) with:

```go
func (h *ProductDetailsHandler) ListFeatures(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	editingID, _ := strconv.ParseInt(c.QueryParam("edit"), 10, 64)

	// Fetch all features for this product
	features, err := h.queries.ListProductFeatures(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list features", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Render the features partial template
	return h.renderPartial(c, "product_features", map[string]interface{}{
		"ProductID": id,
		"Features":  features,
		"EditingID": editingID,
	})
}
```

- [ ] **Step 4: Add the `UpdateFeature` handler** in the same file, immediately after `DeleteFeature` (after line ~290):

```go
// UpdateFeature handles POST requests to /admin/products/:id/features/:feature_id
// Updates a single feature's text and display order, then returns the refreshed list.
func (h *ProductDetailsHandler) UpdateFeature(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	featureID, _ := strconv.ParseInt(c.Param("feature_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	if err := h.queries.UpdateProductFeature(ctx, sqlc.UpdateProductFeatureParams{
		FeatureText:  c.FormValue("feature_text"),
		DisplayOrder: order,
		ID:           featureID,
	}); err != nil {
		h.logger.Error("failed to update feature", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Updated feature for Product #%d", id)
	return h.ListFeatures(c)
}
```

- [ ] **Step 5: Register the route in `cmd/server/main.go`** — add directly after the `DeleteFeature` line (line 362):

```go
	adminGroup.POST("/products/:id/features/:feature_id", pdHandler.UpdateFeature)        // HTMX: update single feature
```

- [ ] **Step 6: Register the route in `internal/e2e/e2e_test.go`** — add directly after `adminGroup.DELETE("/products/:id/features/:feature_id", pdHandler.DeleteFeature)` (line ~296):

```go
		adminGroup.POST("/products/:id/features/:feature_id", pdHandler.UpdateFeature)
```

- [ ] **Step 7: Update the template** `templates/admin/partials/product_features.html` — replace the range body (lines 30–41, i.e. from `{{range .Features}}` through its matching `{{end}}`) with:

```html
        {{range .Features}}
        {{if eq .ID $.EditingID}}
        <!-- Inline edit form (only the row being edited) -->
        <form hx-post="/admin/products/{{$.ProductID}}/features/{{.ID}}"
              hx-target="#features-section"
              hx-swap="outerHTML"
              class="flex items-center gap-3 border-2 border-black px-4 py-3 bg-yellow-50" style="box-shadow: 2px 2px 0px #000;">
            <input type="text" name="feature_text" value="{{.FeatureText}}" required
                   class="flex-1 border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
            <input type="number" name="display_order" value="{{.DisplayOrder}}"
                   class="w-20 border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
            <button type="submit"
                    class="text-xs uppercase font-bold tracking-wider px-3 py-1 border-2 border-black bg-black text-white hover:bg-white hover:text-black transition-colors">Save</button>
            <a hx-get="/admin/products/{{$.ProductID}}/features"
               hx-target="#features-section"
               hx-swap="outerHTML"
               class="text-xs uppercase font-bold tracking-wider px-3 py-1 border-2 border-black cursor-pointer hover:bg-gray-100 transition-colors">Cancel</a>
        </form>
        {{else}}
        <div class="flex items-center gap-3 border-2 border-black px-4 py-3 bg-white group" style="box-shadow: 2px 2px 0px #000;">
            <span class="text-lg font-bold text-gray-400">&#x2022;</span>
            <span class="flex-1 text-sm">{{.FeatureText}}</span>
            <span class="text-xs text-gray-400 font-bold">#{{.DisplayOrder}}</span>
            <button hx-get="/admin/products/{{$.ProductID}}/features?edit={{.ID}}"
                    hx-target="#features-section"
                    hx-swap="outerHTML"
                    class="opacity-0 group-hover:opacity-100 transition-opacity text-gray-600 hover:text-black text-xs font-bold uppercase">&#9998;</button>
            <button hx-delete="/admin/products/{{$.ProductID}}/features/{{.ID}}"
                    hx-target="#features-section"
                    hx-swap="outerHTML"
                    hx-confirm="Delete this feature?"
                    class="opacity-0 group-hover:opacity-100 transition-opacity text-red-600 hover:text-red-800 text-xs font-bold uppercase">&#x2715;</button>
        </div>
        {{end}}
        {{end}}
```

- [ ] **Step 8: Run the tests to verify they pass**

Run: `go test ./internal/e2e/ -run 'TestProductFeatures(EditForm|Update)_E2E' -v`
Expected: PASS (both).

- [ ] **Step 9: Verify the build, then commit**

Run: `go build ./cmd/...`
Expected: exit 0.

```bash
git add internal/handlers/admin/product_details.go cmd/server/main.go internal/e2e/e2e_test.go templates/admin/partials/product_features.html internal/e2e/11_product_features_test.go
git commit -m "feat(products): inline edit + reorder for features"
```

---

### Task 3: Specs — inline edit + reorder

**Files:**
- Modify: `internal/handlers/admin/product_details.go` (edit `ListSpecs`, add `UpdateSpec`)
- Modify: `cmd/server/main.go` (add route after line 356)
- Modify: `internal/e2e/e2e_test.go` (add route after the `DeleteSpec` line ~292)
- Modify: `templates/admin/partials/product_specs.html` (per-spec row, lines 41–51)
- Test: `internal/e2e/10_product_specs_test.go` (append two tests)

**Interfaces:**
- Consumes: `UpdateProductSpec` / `UpdateProductSpecParams` (Task 1).
- Produces: handler `UpdateSpec(c echo.Context) error`; route `POST /admin/products/:id/specs/:spec_id`.

- [ ] **Step 1: Write the failing tests** — append to `internal/e2e/10_product_specs_test.go` (imports already include `net/url`):

```go
func TestProductSpecsEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "SpecEditCat", Slug: "spec-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "SPEC-EDIT-1", Slug: "spec-edit", Name: "Spec Edit Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	spec, _ := queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID: product.ID, SectionName: "Motor", SpecKey: "Voltage", SpecValue: "24V", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/specs?edit=%d", product.ID, spec.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/specs/%d"`, product.ID, spec.ID)) {
		t.Errorf("expected inline edit form for spec %d, body: %s", spec.ID, body)
	}
	if !strings.Contains(body, `value="Voltage"`) {
		t.Errorf("expected pre-filled spec key, body: %s", body)
	}
}

func TestProductSpecsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "SpecUpdCat", Slug: "spec-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "SPEC-UPD-1", Slug: "spec-upd", Name: "Spec Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	spec, _ := queries.CreateProductSpec(ctx, sqlc.CreateProductSpecParams{
		ProductID: product.ID, SectionName: "Motor", SpecKey: "Voltage", SpecValue: "24V", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/specs/%d", product.ID, spec.ID), strings.NewReader(url.Values{
		"section_name":  {"Power"},
		"spec_key":      {"Voltage"},
		"spec_value":    {"48V"},
		"display_order": {"3"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	specs, _ := queries.ListProductSpecs(ctx, product.ID)
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}
	if specs[0].SpecValue != "48V" {
		t.Errorf("expected spec value '48V', got %q", specs[0].SpecValue)
	}
	if specs[0].SectionName != "Power" {
		t.Errorf("expected section 'Power', got %q", specs[0].SectionName)
	}
	if specs[0].DisplayOrder != 3 {
		t.Errorf("expected display_order 3, got %d", specs[0].DisplayOrder)
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/e2e/ -run 'TestProductSpecs(EditForm|Update)_E2E' -v`
Expected: FAIL (Update gets 405; EditForm body lacks the form).

- [ ] **Step 3: Edit `ListSpecs` to read `?edit`** in `product_details.go`. Replace its body (lines ~96–111) with:

```go
func (h *ProductDetailsHandler) ListSpecs(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	editingID, _ := strconv.ParseInt(c.QueryParam("edit"), 10, 64)

	// Fetch all specifications for this product
	specs, err := h.queries.ListProductSpecs(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list specs", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Render the specs partial template
	return h.renderPartial(c, "product_specs", map[string]interface{}{
		"ProductID": id,
		"Specs":     specs,
		"EditingID": editingID,
	})
}
```

- [ ] **Step 4: Add the `UpdateSpec` handler** immediately after `DeleteSpec` (after line ~189):

```go
// UpdateSpec handles POST requests to /admin/products/:id/specs/:spec_id
// Updates a single specification's fields and display order, then returns the refreshed list.
func (h *ProductDetailsHandler) UpdateSpec(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	specID, _ := strconv.ParseInt(c.Param("spec_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	if err := h.queries.UpdateProductSpec(ctx, sqlc.UpdateProductSpecParams{
		SectionName:  c.FormValue("section_name"),
		SpecKey:      c.FormValue("spec_key"),
		SpecValue:    c.FormValue("spec_value"),
		DisplayOrder: order,
		ID:           specID,
	}); err != nil {
		h.logger.Error("failed to update spec", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Updated spec for Product #%d", id)
	return h.ListSpecs(c)
}
```

- [ ] **Step 5: Register the route in `cmd/server/main.go`** — add directly after the `DeleteSpec` line (line 356):

```go
	adminGroup.POST("/products/:id/specs/:spec_id", pdHandler.UpdateSpec)          // HTMX: update single spec
```

- [ ] **Step 6: Register the route in `internal/e2e/e2e_test.go`** — add after `adminGroup.DELETE("/products/:id/specs/:spec_id", pdHandler.DeleteSpec)` (line ~292):

```go
		adminGroup.POST("/products/:id/specs/:spec_id", pdHandler.UpdateSpec)
```

- [ ] **Step 7: Update the template** `templates/admin/partials/product_specs.html` — replace the single spec row block (lines 41–51, from `<div class="flex items-center gap-3 px-4 py-3 border-b ...">` through its closing `</div>`) with:

```html
            {{if eq .ID $.EditingID}}
            <form hx-post="/admin/products/{{$.ProductID}}/specs/{{.ID}}"
                  hx-target="#specs-section"
                  hx-swap="outerHTML"
                  class="px-4 py-3 border-b border-gray-300 last:border-b-0 bg-yellow-50 space-y-2">
                <div class="grid grid-cols-2 md:grid-cols-4 gap-2">
                    <input type="text" name="section_name" value="{{.SectionName}}" required placeholder="Section"
                           class="border-2 border-black px-2 py-1 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                    <input type="text" name="spec_key" value="{{.SpecKey}}" required placeholder="Label"
                           class="border-2 border-black px-2 py-1 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                    <input type="text" name="spec_value" value="{{.SpecValue}}" required placeholder="Value"
                           class="border-2 border-black px-2 py-1 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                    <input type="number" name="display_order" value="{{.DisplayOrder}}" placeholder="Order"
                           class="border-2 border-black px-2 py-1 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
                <div class="flex gap-2">
                    <button type="submit"
                            class="text-xs uppercase font-bold tracking-wider px-3 py-1 border-2 border-black bg-black text-white hover:bg-white hover:text-black transition-colors">Save</button>
                    <a hx-get="/admin/products/{{$.ProductID}}/specs"
                       hx-target="#specs-section"
                       hx-swap="outerHTML"
                       class="text-xs uppercase font-bold tracking-wider px-3 py-1 border-2 border-black cursor-pointer hover:bg-gray-100 transition-colors">Cancel</a>
                </div>
            </form>
            {{else}}
            <div class="flex items-center gap-3 px-4 py-3 border-b border-gray-300 last:border-b-0 group">
                <div class="flex-1 grid grid-cols-2 gap-3">
                    <div class="text-sm font-bold">{{.SpecKey}}</div>
                    <div class="text-sm text-gray-700">{{.SpecValue}}</div>
                </div>
                <button hx-get="/admin/products/{{$.ProductID}}/specs?edit={{.ID}}"
                        hx-target="#specs-section"
                        hx-swap="outerHTML"
                        class="opacity-0 group-hover:opacity-100 transition-opacity text-gray-600 hover:text-black text-xs font-bold uppercase">&#9998;</button>
                <button hx-delete="/admin/products/{{$.ProductID}}/specs/{{.ID}}"
                        hx-target="#specs-section"
                        hx-swap="outerHTML"
                        hx-confirm="Delete this spec?"
                        class="opacity-0 group-hover:opacity-100 transition-opacity text-red-600 hover:text-red-800 text-xs font-bold uppercase">&#x2715;</button>
            </div>
            {{end}}
```

- [ ] **Step 8: Run the tests to verify they pass**

Run: `go test ./internal/e2e/ -run 'TestProductSpecs(EditForm|Update)_E2E' -v`
Expected: PASS (both).

- [ ] **Step 9: Verify the build, then commit**

Run: `go build ./cmd/...`
Expected: exit 0.

```bash
git add internal/handlers/admin/product_details.go cmd/server/main.go internal/e2e/e2e_test.go templates/admin/partials/product_specs.html internal/e2e/10_product_specs_test.go
git commit -m "feat(products): inline edit + reorder for specs"
```

---

### Task 4: Certifications — inline edit + reorder

**Files:**
- Modify: `internal/handlers/admin/product_details.go` (edit `ListCertifications`, add `UpdateCertification`)
- Modify: `cmd/server/main.go` (add route after line 368)
- Modify: `internal/e2e/e2e_test.go` (add route after the `DeleteCertification` line ~300)
- Modify: `templates/admin/partials/product_certifications.html` (range body, lines 26–43)
- Test: `internal/e2e/12_product_certifications_test.go` (append two tests)

**Interfaces:**
- Consumes: `UpdateProductCertification` / `UpdateProductCertificationParams` (Task 1). `CertificationCode/IconType/IconPath` are `sql.NullString`.
- Produces: handler `UpdateCertification(c echo.Context) error`; route `POST /admin/products/:id/certifications/:cert_id`.

- [ ] **Step 1: Write the failing tests** — append to `internal/e2e/12_product_certifications_test.go` (imports already include `database/sql` and `net/url`):

```go
func TestProductCertificationsEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "CertEditCat", Slug: "cert-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "CERT-EDIT-1", Slug: "cert-edit", Name: "Cert Edit Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	cert, _ := queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID: product.ID, CertificationName: "CE",
		CertificationCode: sql.NullString{String: "EN60950", Valid: true}, DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/certifications?edit=%d", product.ID, cert.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/certifications/%d"`, product.ID, cert.ID)) {
		t.Errorf("expected inline edit form for cert %d, body: %s", cert.ID, body)
	}
	if !strings.Contains(body, `value="CE"`) {
		t.Errorf("expected pre-filled cert name, body: %s", body)
	}
}

func TestProductCertificationsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "CertUpdCat", Slug: "cert-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "CERT-UPD-1", Slug: "cert-upd", Name: "Cert Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	cert, _ := queries.CreateProductCertification(ctx, sqlc.CreateProductCertificationParams{
		ProductID: product.ID, CertificationName: "CE", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/certifications/%d", product.ID, cert.ID), strings.NewReader(url.Values{
		"certification_name": {"UL"},
		"certification_code": {"UL-123"},
		"icon_type":          {"shield"},
		"icon_path":          {"/icons/ul.svg"},
		"display_order":      {"2"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	certs, _ := queries.ListProductCertifications(ctx, product.ID)
	if len(certs) != 1 {
		t.Fatalf("expected 1 cert, got %d", len(certs))
	}
	if certs[0].CertificationName != "UL" {
		t.Errorf("expected cert name 'UL', got %q", certs[0].CertificationName)
	}
	if !certs[0].CertificationCode.Valid || certs[0].CertificationCode.String != "UL-123" {
		t.Errorf("expected cert code 'UL-123', got %+v", certs[0].CertificationCode)
	}
	if certs[0].DisplayOrder != 2 {
		t.Errorf("expected display_order 2, got %d", certs[0].DisplayOrder)
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

Run: `go test ./internal/e2e/ -run 'TestProductCertifications(EditForm|Update)_E2E' -v`
Expected: FAIL.

- [ ] **Step 3: Edit `ListCertifications` to read `?edit`** in `product_details.go`. Replace its body (lines ~303–318) with:

```go
func (h *ProductDetailsHandler) ListCertifications(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	editingID, _ := strconv.ParseInt(c.QueryParam("edit"), 10, 64)

	// Fetch all certifications for this product
	certs, err := h.queries.ListProductCertifications(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list certifications", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Render the certifications partial template
	return h.renderPartial(c, "product_certifications", map[string]interface{}{
		"ProductID":      id,
		"Certifications": certs,
		"EditingID":      editingID,
	})
}
```

- [ ] **Step 4: Add the `UpdateCertification` handler** immediately after `DeleteCertification` (after line ~402):

```go
// UpdateCertification handles POST requests to /admin/products/:id/certifications/:cert_id
// Updates a single certification's fields and display order, then returns the refreshed list.
func (h *ProductDetailsHandler) UpdateCertification(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	certID, _ := strconv.ParseInt(c.Param("cert_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	certCode := c.FormValue("certification_code")
	iconType := c.FormValue("icon_type")
	iconPath := c.FormValue("icon_path")

	if err := h.queries.UpdateProductCertification(ctx, sqlc.UpdateProductCertificationParams{
		CertificationName: c.FormValue("certification_name"),
		CertificationCode: sql.NullString{String: certCode, Valid: certCode != ""},
		IconType:          sql.NullString{String: iconType, Valid: iconType != ""},
		IconPath:          sql.NullString{String: iconPath, Valid: iconPath != ""},
		DisplayOrder:      order,
		ID:                certID,
	}); err != nil {
		h.logger.Error("failed to update certification", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Updated certification for Product #%d", id)
	return h.ListCertifications(c)
}
```

- [ ] **Step 5: Register the route in `cmd/server/main.go`** — add directly after the `DeleteCertification` line (line 368):

```go
	adminGroup.POST("/products/:id/certifications/:cert_id", pdHandler.UpdateCertification)   // HTMX: update single cert
```

- [ ] **Step 6: Register the route in `internal/e2e/e2e_test.go`** — add after `adminGroup.DELETE("/products/:id/certifications/:cert_id", pdHandler.DeleteCertification)` (line ~300):

```go
		adminGroup.POST("/products/:id/certifications/:cert_id", pdHandler.UpdateCertification)
```

- [ ] **Step 7: Update the template** `templates/admin/partials/product_certifications.html` — replace the range body (lines 26–43, the `<div class="flex items-center gap-4 ...">` card through its closing `</div>`) with:

```html
        {{if eq .ID $.EditingID}}
        <form hx-post="/admin/products/{{$.ProductID}}/certifications/{{.ID}}"
              hx-target="#certifications-section"
              hx-swap="outerHTML"
              class="border-2 border-black p-4 space-y-3 bg-yellow-50" style="box-shadow: 2px 2px 0px #000;">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">Name *</label>
                    <input type="text" name="certification_name" value="{{.CertificationName}}" required
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">Code</label>
                    <input type="text" name="certification_code" value="{{if .CertificationCode.Valid}}{{.CertificationCode.String}}{{end}}"
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">Icon Type</label>
                    <input type="text" name="icon_type" value="{{if .IconType.Valid}}{{.IconType.String}}{{end}}"
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">Order</label>
                    <input type="number" name="display_order" value="{{.DisplayOrder}}"
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
            </div>
            <div>
                <label class="block text-xs font-bold uppercase tracking-wider mb-1">Icon Path</label>
                <input type="text" name="icon_path" value="{{if .IconPath.Valid}}{{.IconPath.String}}{{end}}"
                       class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
            </div>
            <div class="flex gap-2">
                <button type="submit"
                        class="bg-black text-white px-6 py-2 text-sm font-bold uppercase tracking-wider border-2 border-black hover:bg-white hover:text-black transition-colors" style="box-shadow: 3px 3px 0px #000;">Save</button>
                <a hx-get="/admin/products/{{$.ProductID}}/certifications"
                   hx-target="#certifications-section"
                   hx-swap="outerHTML"
                   class="px-6 py-2 text-sm font-bold uppercase tracking-wider border-2 border-black cursor-pointer hover:bg-gray-100 transition-colors">Cancel</a>
            </div>
        </form>
        {{else}}
        <div class="flex items-center gap-4 border-2 border-black px-4 py-3 bg-white group" style="box-shadow: 2px 2px 0px #000;">
            <div class="w-10 h-10 border-2 border-black bg-green-100 flex items-center justify-center flex-shrink-0">
                <span class="text-green-700 font-bold text-sm">&#x2713;</span>
            </div>
            <div class="flex-1 min-w-0">
                <div class="text-sm font-bold">{{.CertificationName}}</div>
                <div class="text-xs text-gray-500 flex gap-3">
                    {{if .CertificationCode.Valid}}<span>Code: {{.CertificationCode.String}}</span>{{end}}
                    {{if .IconType.Valid}}<span>Type: {{.IconType.String}}</span>{{end}}
                </div>
            </div>
            <span class="text-xs text-gray-400 font-bold">#{{.DisplayOrder}}</span>
            <button hx-get="/admin/products/{{$.ProductID}}/certifications?edit={{.ID}}"
                    hx-target="#certifications-section"
                    hx-swap="outerHTML"
                    class="opacity-0 group-hover:opacity-100 transition-opacity text-gray-600 hover:text-black text-xs font-bold uppercase">&#9998;</button>
            <button hx-delete="/admin/products/{{$.ProductID}}/certifications/{{.ID}}"
                    hx-target="#certifications-section"
                    hx-swap="outerHTML"
                    hx-confirm="Delete this certification?"
                    class="opacity-0 group-hover:opacity-100 transition-opacity text-red-600 hover:text-red-800 text-xs font-bold uppercase">&#x2715;</button>
        </div>
        {{end}}
```

- [ ] **Step 8: Run the tests to verify they pass**

Run: `go test ./internal/e2e/ -run 'TestProductCertifications(EditForm|Update)_E2E' -v`
Expected: PASS (both).

- [ ] **Step 9: Verify the build, then commit**

Run: `go build ./cmd/...`
Expected: exit 0.

```bash
git add internal/handlers/admin/product_details.go cmd/server/main.go internal/e2e/e2e_test.go templates/admin/partials/product_certifications.html internal/e2e/12_product_certifications_test.go
git commit -m "feat(products): inline edit + reorder for certifications"
```

---

### Task 5: Downloads — inline edit + reorder (metadata only)

**Files:**
- Modify: `internal/handlers/admin/product_details.go` (edit `ListDownloads`, add `UpdateDownload`)
- Modify: `cmd/server/main.go` (add route after line 373)
- Modify: `internal/e2e/e2e_test.go` (add route after the `DeleteDownload` line ~303)
- Modify: `templates/admin/partials/product_downloads.html` (range body, lines 19–36)
- Modify + Test: `internal/e2e/13_product_downloads_test.go` (ADD `"net/url"` to imports, append two tests)

**Interfaces:**
- Consumes: `UpdateProductDownload` / `UpdateProductDownloadParams` (Task 1). `Description/Version` are `sql.NullString`. `FilePath`, `FileSize`, `DownloadCount` are NOT touched by the update (file stays put; `updated_at` auto-bumps via DB trigger).
- Produces: handler `UpdateDownload(c echo.Context) error`; route `POST /admin/products/:id/downloads/:download_id`.

- [ ] **Step 1: Add the `net/url` import to `internal/e2e/13_product_downloads_test.go`** — the import block currently is `context, fmt, mime/multipart, net/http, net/http/httptest, strings, testing, sqlc`. Add `"net/url"` (alphabetically after `net/http/httptest`):

```go
import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)
```

- [ ] **Step 2: Write the failing tests** — append to `internal/e2e/13_product_downloads_test.go`:

```go
func TestProductDownloadsEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "DlEditCat", Slug: "dl-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-EDIT-1", Slug: "dl-edit", Name: "Download Edit Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	dl, _ := queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: product.ID, Title: "Datasheet", FileType: ".pdf",
		FilePath: "/uploads/downloads/ds.pdf", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/downloads?edit=%d", product.ID, dl.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/downloads/%d"`, product.ID, dl.ID)) {
		t.Errorf("expected inline edit form for download %d, body: %s", dl.ID, body)
	}
	if !strings.Contains(body, `value="Datasheet"`) {
		t.Errorf("expected pre-filled download title, body: %s", body)
	}
}

func TestProductDownloadsUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "DlUpdCat", Slug: "dl-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "DL-UPD-1", Slug: "dl-upd", Name: "Download Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	dl, _ := queries.CreateProductDownload(ctx, sqlc.CreateProductDownloadParams{
		ProductID: product.ID, Title: "Datasheet", FileType: ".pdf",
		FilePath: "/uploads/downloads/ds.pdf", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/downloads/%d", product.ID, dl.ID), strings.NewReader(url.Values{
		"title":         {"User Manual"},
		"file_type":     {".pdf"},
		"version":       {"2.0"},
		"description":   {"Updated manual"},
		"display_order": {"4"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	downloads, _ := queries.ListProductDownloads(ctx, product.ID)
	if len(downloads) != 1 {
		t.Fatalf("expected 1 download, got %d", len(downloads))
	}
	if downloads[0].Title != "User Manual" {
		t.Errorf("expected title 'User Manual', got %q", downloads[0].Title)
	}
	if downloads[0].DisplayOrder != 4 {
		t.Errorf("expected display_order 4, got %d", downloads[0].DisplayOrder)
	}
	// Metadata-only edit must preserve the uploaded file path.
	if downloads[0].FilePath != "/uploads/downloads/ds.pdf" {
		t.Errorf("expected file path preserved, got %q", downloads[0].FilePath)
	}
}
```

- [ ] **Step 3: Run the tests to verify they fail**

Run: `go test ./internal/e2e/ -run 'TestProductDownloads(EditForm|Update)_E2E' -v`
Expected: FAIL.

- [ ] **Step 4: Edit `ListDownloads` to read `?edit`** in `product_details.go`. Replace its body (lines ~415–430) with:

```go
func (h *ProductDetailsHandler) ListDownloads(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	editingID, _ := strconv.ParseInt(c.QueryParam("edit"), 10, 64)

	// Fetch all downloadable files for this product
	downloads, err := h.queries.ListProductDownloads(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list downloads", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Render the downloads partial template
	return h.renderPartial(c, "product_downloads", map[string]interface{}{
		"ProductID": id,
		"Downloads": downloads,
		"EditingID": editingID,
	})
}
```

- [ ] **Step 5: Add the `UpdateDownload` handler** immediately after `DeleteDownload` (after line ~528):

```go
// UpdateDownload handles POST requests to /admin/products/:id/downloads/:download_id
// Updates a download's metadata and display order (NOT the file), then returns the refreshed list.
func (h *ProductDetailsHandler) UpdateDownload(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	downloadID, _ := strconv.ParseInt(c.Param("download_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	desc := c.FormValue("description")
	version := c.FormValue("version")
	fileType := c.FormValue("file_type")

	if err := h.queries.UpdateProductDownload(ctx, sqlc.UpdateProductDownloadParams{
		Title:        c.FormValue("title"),
		Description:  sql.NullString{String: desc, Valid: desc != ""},
		FileType:     fileType,
		Version:      sql.NullString{String: version, Valid: version != ""},
		DisplayOrder: order,
		ID:           downloadID,
	}); err != nil {
		h.logger.Error("failed to update download", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Updated download for Product #%d", id)
	return h.ListDownloads(c)
}
```

- [ ] **Step 6: Register the route in `cmd/server/main.go`** — add directly after the `DeleteDownload` line (line 373):

```go
	adminGroup.POST("/products/:id/downloads/:download_id", pdHandler.UpdateDownload)      // HTMX: update download metadata
```

- [ ] **Step 7: Register the route in `internal/e2e/e2e_test.go`** — add after `adminGroup.DELETE("/products/:id/downloads/:download_id", pdHandler.DeleteDownload)` (line ~303):

```go
		adminGroup.POST("/products/:id/downloads/:download_id", pdHandler.UpdateDownload)
```

- [ ] **Step 8: Update the template** `templates/admin/partials/product_downloads.html` — replace the range body (lines 19–36, the `<div class="flex items-center gap-4 ...">` card through its closing `</div>`) with:

```html
        {{if eq .ID $.EditingID}}
        <form hx-post="/admin/products/{{$.ProductID}}/downloads/{{.ID}}"
              hx-target="#downloads-section"
              hx-swap="outerHTML"
              class="border-2 border-black p-4 space-y-3 bg-yellow-50" style="box-shadow: 2px 2px 0px #000;">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">Title *</label>
                    <input type="text" name="title" value="{{.Title}}" required
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">File Type</label>
                    <input type="text" name="file_type" value="{{.FileType}}"
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">Version</label>
                    <input type="text" name="version" value="{{if .Version.Valid}}{{.Version.String}}{{end}}"
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
                <div>
                    <label class="block text-xs font-bold uppercase tracking-wider mb-1">Order</label>
                    <input type="number" name="display_order" value="{{.DisplayOrder}}"
                           class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
                </div>
            </div>
            <div>
                <label class="block text-xs font-bold uppercase tracking-wider mb-1">Description</label>
                <input type="text" name="description" value="{{if .Description.Valid}}{{.Description.String}}{{end}}"
                       class="w-full border-2 border-black px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
            </div>
            <p class="text-xs text-gray-500">Editing metadata only &mdash; to replace the file, delete and re-add the download.</p>
            <div class="flex gap-2">
                <button type="submit"
                        class="bg-black text-white px-6 py-2 text-sm font-bold uppercase tracking-wider border-2 border-black hover:bg-white hover:text-black transition-colors" style="box-shadow: 3px 3px 0px #000;">Save</button>
                <a hx-get="/admin/products/{{$.ProductID}}/downloads"
                   hx-target="#downloads-section"
                   hx-swap="outerHTML"
                   class="px-6 py-2 text-sm font-bold uppercase tracking-wider border-2 border-black cursor-pointer hover:bg-gray-100 transition-colors">Cancel</a>
            </div>
        </form>
        {{else}}
        <div class="flex items-center gap-4 border-2 border-black px-4 py-3 bg-white group" style="box-shadow: 2px 2px 0px #000;">
            <div class="w-10 h-10 border-2 border-black bg-blue-100 flex items-center justify-center flex-shrink-0">
                <span class="text-blue-700 font-bold text-xs uppercase">{{.FileType}}</span>
            </div>
            <div class="flex-1 min-w-0">
                <a href="{{.FilePath}}" target="_blank" class="text-sm font-bold hover:underline">{{.Title}}</a>
                <div class="text-xs text-gray-500 flex gap-3">
                    {{if .Version.Valid}}<span>v{{.Version.String}}</span>{{end}}
                    {{if .Description.Valid}}<span>{{.Description.String}}</span>{{end}}
                </div>
            </div>
            <span class="text-xs text-gray-400 font-bold">#{{.DisplayOrder}}</span>
            <button hx-get="/admin/products/{{$.ProductID}}/downloads?edit={{.ID}}"
                    hx-target="#downloads-section"
                    hx-swap="outerHTML"
                    class="opacity-0 group-hover:opacity-100 transition-opacity text-gray-600 hover:text-black text-xs font-bold uppercase">&#9998;</button>
            <button hx-delete="/admin/products/{{$.ProductID}}/downloads/{{.ID}}"
                    hx-target="#downloads-section"
                    hx-swap="outerHTML"
                    hx-confirm="Delete this download?"
                    class="opacity-0 group-hover:opacity-100 transition-opacity text-red-600 hover:text-red-800 text-xs font-bold uppercase">&#x2715;</button>
        </div>
        {{end}}
```

- [ ] **Step 9: Run the tests to verify they pass**

Run: `go test ./internal/e2e/ -run 'TestProductDownloads(EditForm|Update)_E2E' -v`
Expected: PASS (both).

- [ ] **Step 10: Verify the build, then commit**

Run: `go build ./cmd/...`
Expected: exit 0.

```bash
git add internal/handlers/admin/product_details.go cmd/server/main.go internal/e2e/e2e_test.go templates/admin/partials/product_downloads.html internal/e2e/13_product_downloads_test.go
git commit -m "feat(products): inline edit + reorder for downloads (metadata only)"
```

---

### Task 6: Images — inline edit + reorder (metadata only)

**Files:**
- Modify: `internal/handlers/admin/product_details.go` (edit `ListImages`, add `UpdateImage`)
- Modify: `cmd/server/main.go` (add route after line 378)
- Modify: `internal/e2e/e2e_test.go` (add route after the `DeleteImage` line ~306)
- Modify: `templates/admin/partials/product_images.html` (range body, lines 19–50)
- Modify + Test: `internal/e2e/14_product_images_test.go` (ADD `"net/url"` to imports, append two tests)

**Interfaces:**
- Consumes: `UpdateProductImage` / `UpdateProductImageParams` (Task 1). `AltText/Caption` are `sql.NullString`. `ImagePath` and `IsThumbnail` are NOT touched by the update (preserved).
- Produces: handler `UpdateImage(c echo.Context) error`; route `POST /admin/products/:id/images/:image_id`.

- [ ] **Step 1: Add the `net/url` import to `internal/e2e/14_product_images_test.go`** — the import block currently is `context, fmt, mime/multipart, net/http, net/http/httptest, strings, testing, sqlc`. Add `"net/url"`:

```go
import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)
```

- [ ] **Step 2: Write the failing tests** — append to `internal/e2e/14_product_images_test.go`:

```go
func TestProductImagesEditForm_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "ImgEditCat", Slug: "img-edit-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-EDIT-1", Slug: "img-edit", Name: "Image Edit Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	img, _ := queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: product.ID, ImagePath: "/uploads/products/img.jpg", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/products/%d/images?edit=%d", product.ID, img.ID), nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`hx-post="/admin/products/%d/images/%d"`, product.ID, img.ID)) {
		t.Errorf("expected inline edit form for image %d, body: %s", img.ID, body)
	}
	if !strings.Contains(body, `name="alt_text"`) {
		t.Errorf("expected alt_text input in edit form, body: %s", body)
	}
}

func TestProductImagesUpdate_E2E(t *testing.T) {
	e, queries, cleanup := setupApp(t)
	defer cleanup()
	createTestAdmin(t, queries)
	cookie := loginAndGetCookie(t, e)

	ctx := context.Background()
	cat, _ := queries.CreateProductCategory(ctx, sqlc.CreateProductCategoryParams{
		Name: "ImgUpdCat", Slug: "img-upd-cat", Description: "d", Icon: "i", SortOrder: 1,
	})
	product, _ := queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Sku: "IMG-UPD-1", Slug: "img-upd", Name: "Image Update Product",
		Description: "Test", CategoryID: cat.ID, Status: "draft",
	})
	img, _ := queries.CreateProductImage(ctx, sqlc.CreateProductImageParams{
		ProductID: product.ID, ImagePath: "/uploads/products/img.jpg", DisplayOrder: 1,
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/admin/products/%d/images/%d", product.ID, img.ID), strings.NewReader(url.Values{
		"alt_text":      {"Front view"},
		"caption":       {"Product front"},
		"display_order": {"7"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	images, _ := queries.ListProductImages(ctx, product.ID)
	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}
	if !images[0].AltText.Valid || images[0].AltText.String != "Front view" {
		t.Errorf("expected alt text 'Front view', got %+v", images[0].AltText)
	}
	if images[0].DisplayOrder != 7 {
		t.Errorf("expected display_order 7, got %d", images[0].DisplayOrder)
	}
	// Metadata-only edit must preserve the uploaded image path.
	if images[0].ImagePath != "/uploads/products/img.jpg" {
		t.Errorf("expected image path preserved, got %q", images[0].ImagePath)
	}
}
```

- [ ] **Step 3: Run the tests to verify they fail**

Run: `go test ./internal/e2e/ -run 'TestProductImages(EditForm|Update)_E2E' -v`
Expected: FAIL.

- [ ] **Step 4: Edit `ListImages` to read `?edit`** in `product_details.go`. Replace its body (lines ~541–556) with:

```go
func (h *ProductDetailsHandler) ListImages(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	editingID, _ := strconv.ParseInt(c.QueryParam("edit"), 10, 64)

	// Fetch all gallery images for this product
	images, err := h.queries.ListProductImages(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("failed to list images", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Render the images partial template
	return h.renderPartial(c, "product_images", map[string]interface{}{
		"ProductID": id,
		"Images":    images,
		"EditingID": editingID,
	})
}
```

- [ ] **Step 5: Add the `UpdateImage` handler** immediately after `DeleteImage` (after line ~646, at end of file):

```go
// UpdateImage handles POST requests to /admin/products/:id/images/:image_id
// Updates an image's alt text, caption, and display order (NOT the file or thumbnail flag),
// then returns the refreshed gallery.
func (h *ProductDetailsHandler) UpdateImage(c echo.Context) error {
	ctx := c.Request().Context()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	imageID, _ := strconv.ParseInt(c.Param("image_id"), 10, 64)
	order, _ := strconv.ParseInt(c.FormValue("display_order"), 10, 64)

	altText := c.FormValue("alt_text")
	caption := c.FormValue("caption")

	if err := h.queries.UpdateProductImage(ctx, sqlc.UpdateProductImageParams{
		AltText:      sql.NullString{String: altText, Valid: altText != ""},
		Caption:      sql.NullString{String: caption, Valid: caption != ""},
		DisplayOrder: order,
		ID:           imageID,
	}); err != nil {
		h.logger.Error("failed to update image", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	logActivity(c, "updated", "product", id, "", "Updated image for Product #%d", id)
	return h.ListImages(c)
}
```

- [ ] **Step 6: Register the route in `cmd/server/main.go`** — add directly after the `DeleteImage` line (line 378):

```go
	adminGroup.POST("/products/:id/images/:image_id", pdHandler.UpdateImage)       // HTMX: update image metadata
```

- [ ] **Step 7: Register the route in `internal/e2e/e2e_test.go`** — add after `adminGroup.DELETE("/products/:id/images/:image_id", pdHandler.DeleteImage)` (line ~306):

```go
		adminGroup.POST("/products/:id/images/:image_id", pdHandler.UpdateImage)
```

- [ ] **Step 8: Update the template** `templates/admin/partials/product_images.html` — replace the range body (lines 19–50, the `<div class="border-2 border-black bg-white relative group" ...>` card through its closing `</div>`) with the block below. Note: the existing dead "Set Primary" button is intentionally preserved (out of scope); only the Edit button and the edit form are added:

```html
        {{if eq $img.ID $.EditingID}}
        <form hx-post="/admin/products/{{$.ProductID}}/images/{{$img.ID}}"
              hx-target="#images-section"
              hx-swap="outerHTML"
              class="border-2 border-black bg-yellow-50 p-3 space-y-2" style="box-shadow: 3px 3px 0px #000;">
            <img src="{{$img.ImagePath}}" alt="" class="w-full h-24 object-cover border-2 border-black">
            <div>
                <label class="block text-xs font-bold uppercase tracking-wider mb-1">Alt Text</label>
                <input type="text" name="alt_text" value="{{if $img.AltText.Valid}}{{$img.AltText.String}}{{end}}"
                       class="w-full border-2 border-black px-2 py-1 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
            </div>
            <div>
                <label class="block text-xs font-bold uppercase tracking-wider mb-1">Caption</label>
                <input type="text" name="caption" value="{{if $img.Caption.Valid}}{{$img.Caption.String}}{{end}}"
                       class="w-full border-2 border-black px-2 py-1 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
            </div>
            <div>
                <label class="block text-xs font-bold uppercase tracking-wider mb-1">Order</label>
                <input type="number" name="display_order" value="{{$img.DisplayOrder}}"
                       class="w-full border-2 border-black px-2 py-1 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-yellow-300">
            </div>
            <div class="flex gap-2">
                <button type="submit"
                        class="flex-1 bg-black text-white px-3 py-1 text-xs font-bold uppercase tracking-wider border-2 border-black hover:bg-white hover:text-black transition-colors">Save</button>
                <a hx-get="/admin/products/{{$.ProductID}}/images"
                   hx-target="#images-section"
                   hx-swap="outerHTML"
                   class="flex-1 text-center px-3 py-1 text-xs font-bold uppercase tracking-wider border-2 border-black cursor-pointer hover:bg-gray-100 transition-colors">Cancel</a>
            </div>
        </form>
        {{else}}
        <div class="border-2 border-black bg-white relative group" style="box-shadow: 3px 3px 0px #000;">
            {{if $img.IsThumbnail}}
            <div class="absolute top-0 left-0 z-10 bg-yellow-300 border-b-2 border-r-2 border-black px-2 py-1">
                <span class="text-xs font-bold uppercase tracking-wider">Primary</span>
            </div>
            {{end}}
            <div class="relative overflow-hidden">
                <img src="{{$img.ImagePath}}" alt="{{if $img.AltText.Valid}}{{$img.AltText.String}}{{else}}Product image{{end}}" class="w-full h-40 object-cover">
                <!-- Hover overlay -->
                <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-50 transition-all flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100">
                    {{if not $img.IsThumbnail}}
                    <button hx-post="/admin/products/{{$.ProductID}}/images/{{$img.ID}}/primary"
                            hx-target="#images-section"
                            hx-swap="outerHTML"
                            class="bg-yellow-300 border-2 border-black px-3 py-1 text-xs font-bold uppercase hover:bg-yellow-400" style="box-shadow: 2px 2px 0px #000;">
                        Set Primary
                    </button>
                    {{end}}
                    <button hx-get="/admin/products/{{$.ProductID}}/images?edit={{$img.ID}}"
                            hx-target="#images-section"
                            hx-swap="outerHTML"
                            class="bg-white border-2 border-black px-3 py-1 text-xs font-bold uppercase hover:bg-gray-100" style="box-shadow: 2px 2px 0px #000;">
                        Edit
                    </button>
                    <button hx-delete="/admin/products/{{$.ProductID}}/images/{{$img.ID}}"
                            hx-target="#images-section"
                            hx-swap="outerHTML"
                            hx-confirm="Delete this image?"
                            class="bg-red-500 text-white border-2 border-black px-3 py-1 text-xs font-bold uppercase hover:bg-red-600" style="box-shadow: 2px 2px 0px #000;">
                        Delete
                    </button>
                </div>
            </div>
            <div class="px-3 py-2 border-t-2 border-black">
                {{if $img.AltText.Valid}}<div class="text-xs text-gray-600 truncate">{{$img.AltText.String}}</div>{{end}}
                <div class="text-xs text-gray-400 font-bold">#{{$img.DisplayOrder}}</div>
            </div>
        </div>
        {{end}}
```

- [ ] **Step 9: Run the tests to verify they pass**

Run: `go test ./internal/e2e/ -run 'TestProductImages(EditForm|Update)_E2E' -v`
Expected: PASS (both).

- [ ] **Step 10: Verify the build, then commit**

Run: `go build ./cmd/...`
Expected: exit 0.

```bash
git add internal/handlers/admin/product_details.go cmd/server/main.go internal/e2e/e2e_test.go templates/admin/partials/product_images.html internal/e2e/14_product_images_test.go
git commit -m "feat(products): inline edit + reorder for images (metadata only)"
```

---

### Task 7: Full-suite verification + manual smoke

**Files:** none (verification only)

- [ ] **Step 1: Run the entire test suite**

Run: `go test ./...`
Expected: PASS across all packages (in particular all `internal/e2e` product sub-entity tests, both pre-existing and the 10 new ones).

- [ ] **Step 2: Run go vet**

Run: `go vet ./...`
Expected: no findings.

- [ ] **Step 3: Confirm both route tables are in sync**

Run: `grep -nE "pdHandler.Update(Feature|Spec|Certification|Download|Image)" cmd/server/main.go internal/e2e/e2e_test.go`
Expected: 10 lines total — 5 in `cmd/server/main.go` and 5 in `internal/e2e/e2e_test.go`.

- [ ] **Step 4: Manual smoke test in the running app**

Run the server: `make run` (or `go run cmd/server/main.go`). Log in to the admin panel, open an existing product's edit page (`/admin/products/<id>/edit`), and for each tab (Features, Specs, Certifications, Downloads, Images):
  1. Click the ✎ (or "Edit" on images) on an existing row → the row becomes an inline form pre-filled with current values.
  2. Change a text field and the numeric Order, click Save → the list re-renders, the value is updated, and the row moves to its new sorted position.
  3. Click ✎ again then Cancel → the form reverts to the display row with no change.

Confirm no console errors and that downloads/images retain their file/image after a metadata edit.

- [ ] **Step 5: Final confirmation**

No commit needed (Tasks 1–6 each committed). Report the suite result and smoke-test outcome.

---

## Self-Review

**1. Spec coverage:**
- Edit feature/spec/cert/download/image details → Tasks 2–6 (handlers + templates + queries). ✓
- Change numeric sort order → every edit form includes a `display_order` number input; Update handlers persist it; List queries `ORDER BY display_order ASC` re-sort on save; reorder asserted in `TestProductFeaturesUpdate_E2E` and via display_order assertions in the others. ✓
- All five sub-entities → Tasks 2–6, one each. ✓
- Inline edit toggle UX → `?edit=<id>` + `{{if eq .ID $.EditingID}}` per partial. ✓
- File replacement out of scope (downloads/images) → Update handlers omit file columns; tests assert `FilePath`/`ImagePath` preserved. ✓
- `is_thumbnail` excluded / dead "Set Primary" left as-is → image UPDATE query omits `is_thumbnail`; template preserves the existing button. ✓
- Cache unchanged → no Cache dependency added to `ProductDetailsHandler`. ✓
- Routes in both main.go and e2e setupApp → each task registers in both; Task 7 Step 3 verifies. ✓

**2. Placeholder scan:** No TBD/TODO/"add error handling"/"similar to Task N". Every code step shows complete code. ✓

**3. Type consistency:** Handlers use the exact generated param structs from Task 1's Interfaces block (`UpdateProductFeatureParams{FeatureText, DisplayOrder, ID}`, etc.). Route param names match handler reads (`:feature_id`/`:spec_id`/`:cert_id`/`:download_id`/`:image_id`). Form field `name`s match `c.FormValue` reads and the existing Add forms. Container ids (`#features-section`, `#specs-section`, `#certifications-section`, `#downloads-section`, `#images-section`) match each partial. ✓
