# AI Product-Creation & Editing Agent — Design Spec

**Date:** 2026-06-27
**Status:** Approved (design) — pending implementation plan
**Scope:** First AI feature in Bluejay CMS. A document-grounded chatbot that creates and edits Products, lives on the Products admin pages, uses tool-calling over the sqlc layer, and publishes only on explicit approval.

---

## 1. Goal & User Story

As an admin, I open a chatbot from the Products page, upload one or more documents (PDF, .docx, Excel/CSV, images) describing a product, and converse with it. The agent reads the documents, asks me about important missing details, and assembles a complete **draft** product (core fields, SEO, specifications, features, certifications, image gallery, downloads). It gives me a **preview rendered as the real public product page**, which I can **edit inline field-by-field**. The product becomes publicly visible **only when I approve**. The same chatbot can also **edit an existing product**, staging changes so the live page is untouched until I approve.

### Confirmed decisions (from interview)
| Topic | Decision |
|---|---|
| Input file types | PDF, .docx, Excel/CSV, Images |
| Model / provider | **Gemini Flash** via Gemini API (`GEMINI_API_KEY`) |
| Missing-info behavior | Ask for important fields; infer minor ones; propose SKU if absent |
| Draft privacy | Hide drafts from public (add status guard) |
| Content populated | Core+SEO, Specifications, Features+Certifications, Images+Downloads |
| Category | Must pick an **existing** category; agent never creates categories |
| Inline editing | Edit **every** field inline in the preview |
| Preview style | The **real public product page** (WYSIWYG) |
| Conversation persistence | Single session (chat transcript ephemeral); **draft row persists in DB** |
| Batch size | One product per conversation |
| SKU source | Extract from docs; **agent may generate** one for confirmation if absent |
| v1 scope | **Create + edit existing** products |
| Edit safety | Live page unchanged until approval (**shadow draft**) |
| Write guardrails | Write to draft freely; **only publish/apply is gated** |
| Entry points | "Create with AI" on Products list; "Edit with AI" on a product |
| Abandoned drafts | Keep as normal drafts in the products list |
| New deps / env | Approved: Gemini client + .docx + .xlsx/CSV parsers + `GEMINI_API_KEY` |

---

## 2. Architecture Overview

```
templates/admin (Products list/edit)
   │  "Create with AI" / "Edit with AI" buttons
   ▼
Chat panel (HTMX + small JS)  ── uploads files, streams messages ──►  AIProductAgentHandler (Echo, /admin/products/ai/*)
                                                                          │
                                                                          ▼
                                                            AIProductAgentService
                                                            ├─ Gemini client (function-calling loop)
                                                            ├─ Document parsing (PDF→Gemini, docx/xlsx/csv→text)
                                                            ├─ Tool implementations  ── over ──►  *sqlc.Queries (+ WithTx)
                                                            ├─ UploadService (+ new raw-bytes helper)
                                                            └─ ActivityLogService
                                                                          │
                                                                          ▼
                                                            Draft Product rows (status='draft')
                                                                          │  preview via ?preview=true (logged-in admin only)
                                                                          ▼
                                                            Real public product template  ──► Approve ──► status='published'
```

### New units (each single-purpose, independently testable)
- **`internal/config`** (new, small): centralizes env reads incl. `GEMINI_API_KEY`. Promotes the existing scattered `os.Getenv` calls only as needed — no broad refactor.
- **`internal/services/gemini/`**: thin Gemini API client — send messages + tool/function declarations, receive tool calls and text. Isolated so it can be mocked in tests. *Interface:* `Generate(ctx, history, tools) (reply, toolCalls, error)`.
- **`internal/services/docparse/`**: converts uploaded files to model-ready input. PDF/image → passthrough bytes + mime; `.docx`/`.xlsx`/`.csv` → extracted plain text / structured rows. *Interface:* `Parse(file) (ParsedDoc, error)`.
- **`internal/services/aiproduct/`**: the agent orchestrator. Holds `*sqlc.Queries`, `UploadService`, `ActivityLogService`, gemini client, docparse. Owns the tool registry, the function-calling loop, grounding rules, slug/SKU collision handling, transactions. *Interface:* `HandleMessage(ctx, session, userMsg, files) (AgentReply, error)`.
- **`internal/handlers/admin/ai_product_agent.go`**: Echo handler exposing the chat endpoints; manages the in-memory/session-scoped conversation; renders chat fragments (HTMX) and the preview link.
- **Templates:** `templates/admin/partials/ai_product_chat.html` (chat panel), plus entry-point buttons added to existing products list/edit pages, and an inline-edit overlay injected into the public product template in preview mode.

### Why these boundaries
Each service answers "what does it do / how is it used / what it depends on" cleanly: gemini knows nothing about products; docparse knows nothing about Gemini; aiproduct composes them and is the only unit touching the DB. This keeps the model-specific and DB-specific concerns swappable and testable in isolation.

---

## 3. The Tool Set (function declarations given to Gemini)

All tools operate over the existing `Querier`. Reads ground the model; writes target the draft; commit is gated in the handler, not exposed as a free agent tool.

**Read / grounding**
- `list_categories()` → existing categories (id, name, slug). Agent must choose one of these.
- `check_sku_available(sku)` / `check_slug_available(slug)` → bool.
- `get_existing_product(slug|id)` → full product detail (edit mode only).

**Write to draft** (all within a tx where multi-row)
- `create_draft_product(core fields, category_id, status='draft')` → draft id.
- `add_spec(product_id, section_name, spec_key, spec_value, order?)`
- `add_feature(product_id, feature_text, order?)`
- `add_certification(product_id, name, code?, icon?)`
- `attach_image(product_id, uploaded_file_ref, alt_text?, caption?, is_thumbnail?)`
- `add_download(product_id, title, file_type, uploaded_file_ref, ...)`
- `update_draft_field(product_id, field, value)` — used by inline edits and chat corrections.

**Not agent tools (handler-controlled, user-gated):** `publish` (create flow) and `apply_to_live` (edit flow). The agent can *propose* publishing in chat, but the actual commit happens only when the user clicks Approve.

Every successful write calls `ActivityLogService.Log(...)` and invalidates the products cache prefix (`page:products`) exactly as the existing admin handlers do.

---

## 4. Grounding Strategy

System prompt instructs Gemini to:
1. Use **only** information present in the uploaded documents / user messages for factual fields (specs, descriptions, certifications).
2. **Ask** the user (in chat) when an *important* field is missing or ambiguous: category (offer the `list_categories` set to choose from), critical specs, product name, description.
3. **Infer** only low-risk derived fields: slug (from name), display orders, thumbnail choice. Mark inferred values in chat.
4. **Propose** a SKU if none found, derived from the name, and ask the user to confirm/replace before using it.
5. Never invent specifications, numbers, or certifications. If unknown, leave blank and note it.

---

## 5. Create Flow (detailed)

1. User clicks **Create with AI** on `/admin/products` → chat panel opens.
2. User uploads documents + describes the product. Handler stores uploaded files (via `UploadService` raw-bytes helper, see §8) and passes parsed content to the agent.
3. Agent extracts fields, calls `list_categories`, asks user to pick a category, asks about important gaps, proposes a SKU if needed.
4. Once it has the minimum required fields (`sku`, `slug`, `name`, `description`, `category_id`), it calls `create_draft_product` (status `draft`, `published_at` NULL) inside a tx, then loops the child `add_*` tools for specs/features/certs/images/downloads.
5. Slug/SKU **collision handling:** before insert, `check_*_available`; on collision, append `-2`, `-3`… and retry (bounded). DB UNIQUE constraint is the backstop.
6. Handler returns a **preview link**: `/products/{category}/{slug}?preview=true`. Because the viewer is a logged-in admin, the new status guard (§7) renders the draft; the public gets 404.
7. **Inline editing:** in preview mode the public template is augmented with an edit overlay (contenteditable fields + small controls) shown **only to logged-in admin**. Each edit posts to a small endpoint that calls `update_draft_field` / child update → **auto-saves to the draft**. Image reorder/remove hits child-update endpoints.
8. User clicks **Approve & Publish** → handler sets `status='published'`, `published_at=now()`, in a tx; logs activity; invalidates cache. Product now live and discoverable.
9. If the user leaves without approving, the draft remains in the products list (normal admin UI can resume/delete it).

---

## 6. Edit-Existing Flow (shadow draft)

1. User clicks **Edit with AI** on a product → chat opens pre-loaded via `get_existing_product`.
2. Handler creates a **shadow draft**: a cloned product row (`status='draft'`) + cloned child rows, linked to the original via a new nullable column `products.shadow_of_id` (or a small `product_edit_drafts` mapping table — chosen at plan time; `shadow_of_id` preferred for simplicity). Unique constraints handled by suffixing the shadow's slug/sku (shadow is never published as-is).
3. Agent and inline editing operate on the **shadow**, previewed at the shadow's `?preview=true` URL. The live product and its public page are **untouched**.
4. **Approve & Apply** → in a single tx: copy the shadow's editable fields onto the original product, replace the original's child rows with the shadow's (delete + reinsert, or diff), delete the shadow. Log activity; invalidate cache. The live page updates atomically.
5. Abandoned shadow drafts: remain as drafts (filtered out of normal product lists by their `shadow_of_id IS NOT NULL` marker so they don't clutter), resumable/deletable. *(No auto-cleanup job in v1.)*

> Note: the shadow is the only way to satisfy "live page unchanged until approval" without a full versioning system. Its child-table copy/replace logic is the main added complexity in v1 and must be covered by tests.

---

## 7. Required Fix to Existing Code (draft privacy)

`internal/handlers/public/products.go` `ProductDetail` (and/or `GetProductDetail`/`GetProductBySlug` path) currently returns **any** product regardless of status. Change:

- If the product's `status != 'published'`: render it **only** when the request is an authenticated-admin preview (`isPreviewRequest` is already true for logged-in admin + `?preview=true`). Otherwise return **404**.
- Discovery surfaces (`ListProducts`, `ListProductsByCategory`, search, sitemap, featured) already filter `status='published'` — no change needed there.

This is a small, contained change but is **required** for the approval gate to mean anything. It must ship with this feature and have a regression test (public gets 404 on a draft; admin-preview gets 200).

---

## 8. Image & Download Upload

- `UploadService` currently takes `*multipart.FileHeader`. Add `UploadProductImageBytes(data []byte, filename string) (string, error)` (and a downloads equivalent) so the agent can store files it received through the chat endpoint. Reuse existing validation (ext allowlist, 5MB cap) and storage path conventions (`/uploads/products/...`).
- Agent references uploaded files by an opaque `uploaded_file_ref` the handler maps to the stored path, then `attach_image`/`add_download` insert the child row.

---

## 9. Data Model Changes

- **No new column needed for drafts** — `products.status` already exists.
- **New:** `products.shadow_of_id INTEGER NULL REFERENCES products(id) ON DELETE CASCADE` (migration), marking a row as a staged edit of another product. Lists/admin views exclude `shadow_of_id IS NOT NULL`. (Alternative `product_edit_drafts` table decided at plan time.)
- `og_image` has a column default; `CreateProduct` doesn't set it — confirm default is acceptable or set explicitly.

---

## 10. Error Handling

- **Transactions:** all multi-row create/apply operations use sqlc `WithTx`; partial failure rolls back so no half-built product/edit.
- **Slug/SKU collisions:** pre-check + bounded suffix retry; UNIQUE constraint is the backstop, surfaced to the user as a chat message (not a 500).
- **Gemini/API failures:** surfaced as a friendly chat error; no DB writes on a failed turn. Timeouts and rate limits retried with backoff in the gemini client.
- **Document parse failures:** the offending file is reported in chat ("couldn't read X"); the agent continues with what it has.
- **FK violations (bad category):** prevented by forcing category choice from `list_categories`.

---

## 11. Testing Strategy

- **gemini client:** unit tests with a stubbed HTTP transport (no live API). Verify request shape (tools/function declarations) and parsing of tool-call responses.
- **docparse:** table tests with sample .docx/.xlsx/.csv/PDF fixtures → expected extracted text/rows.
- **aiproduct service:** tests with a fake gemini client scripting tool-call sequences; assert correct sqlc calls, tx rollback on injected failure, slug/SKU collision retry, grounding (no writes for missing required fields).
- **Draft privacy fix:** handler test — public request to a draft → 404; admin-preview request → 200.
- **Shadow-draft apply:** create product, shadow-edit, apply → original updated, children replaced, shadow deleted, live page reflects changes; verify live page unchanged *before* apply.
- **Activity log:** assert entries written for create/publish/apply.
- Follow the project's existing test patterns (`internal/e2e/*` for end-to-end where applicable).

---

## 12. Out of Scope (v1 / YAGNI)

- Persisting chat transcripts / resume-across-sessions (explicitly single-session).
- Multiple products from one conversation/batch.
- Auto-cleanup job for abandoned drafts.
- Creating new categories via the agent.
- Translation, SEO-audit, bulk-ops, or any other agent (separate future features).
- Token-based shareable preview links (preview stays admin-session gated).
- Streaming token-by-token UI (acceptable to return full message per turn in v1 unless trivial).

---

## 13. Open Implementation Choices (resolve in the plan, not blocking design)

- Exact Gemini model id (latest Flash) — configurable via env, default chosen at plan time.
- Go libraries for `.docx` (e.g. a maintained docx text extractor) and `.xlsx` (e.g. excelize) and CSV (stdlib).
- `shadow_of_id` column vs. `product_edit_drafts` mapping table.
- Gemini Go SDK vs. thin hand-rolled HTTP client (lean toward thin client to avoid heavy deps).
- Whether to set `og_image` explicitly in the agent's create path.
