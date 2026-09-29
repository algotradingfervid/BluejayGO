# BlueJay UI/UX review — implementation plan

Baseline: `dd1d30d`, branch `codex/ui-ux-review-20260929`. Local application at port 28123, disposable SQLite copy; no production writes. Review date: 29 September 2026.

## Evidence and decisions

132 rendered screen/viewport combinations (66 routes at 1440×1000 and 390×1000) are recorded in `output/playwright/ui-ux-20260929/before.json` and `before/`. No application JavaScript exceptions or HTTP 500 responses occurred. Document width alone missed clipped content inside the admin scroll container; screenshots and element bounds confirmed it.

| Priority | Finding | Evidence | Planned improvement |
|---|---|---|---|
| P1 | Mobile content tables hide status and edit/delete controls off-screen | `before/_admin_products-390.png`, solutions, blog, categories, contacts and activity | Labelled stacked records on small screens; retain semantic tables and desktop layout |
| P1 | Header navigation setting rows overflow on phones | `before/_admin_header-390.png` | Stack each visibility control and label field; fix checkbox visuals and keyboard focus |
| P1 | 174 controls without explicit accessible labels in desktop sample | `before.json`, solutions/header/partners forms | Add explicit names and associate existing visible labels across admin templates |
| P2 | Product editor puts the save action below lengthy content and optional sections | `before/_admin_products_new-1440.png` and 390px counterpart | Three freely navigable stages: Basics, Content & media, Visibility & SEO; show-all escape hatch; cross-stage error recovery and persistent field values |
| P2 | Contact enquiry form follows long office/map content on mobile | `before/_contact-390.png` | Place enquiry form first in DOM, preserve desktop two-column composition |
| P2 | Color inputs mark untouched whitepaper forms dirty because browser lowercases hex values | Whitepapers new form, uppercase defaults in template | Compare normalized default color value; preserve actual-change detection |
| P2 | Sidebar highlights both parent listing and nested settings page | Products settings navigation | Use most-specific matching destination and aria-current |
| P2 | Mislabelled seed image files render as broken images | Product JPEG paths contain SVG bytes; screenshot and file inspection | Graceful, accessible missing-image presentation; report seed-data repair separately |
| P3 | Interaction transitions lack consistent reduced-motion handling | Shared styles and sidebar | Subtle stage, menu and dialog transitions; first inspect slow timeline frames, then set final timings |
| P2 | Keyboard users must traverse site navigation to reach main content | Public/admin layouts | Visible-on-focus skip link and strong focus indicator |

Splitting every form would add unnecessary navigation. Limit the new staged experience to the product editor, where the observed length and five existing section groups support it. Keep all values in one form and one server submission; do not add draft storage or change publication semantics.

## Verification

Repeat route matrix, add 320px and tablet checks, keyboard and reduced-motion checks, product create/edit/save/reload/duplicate rejection, invalid fields on inactive stages, long unbroken content, mobile navigation and filters, contact validation, HTMX details, failed request recovery, and existing Go/JS/browser regression checks. Capture slow animation at multiple timeline positions before choosing final duration. Document exact checks and limitations; do not imply a complete accessibility certification or exhaustive test of every CRUD permutation.

Guidance: [WAI multi-page forms](https://www.w3.org/WAI/tutorials/forms/multi-page/), [WCAG reflow](https://www.w3.org/WAI/WCAG21/Understanding/reflow/), [reduced motion](https://www.w3.org/WAI/WCAG21/Techniques/css/C39). These inform staged orientation, reflow and optional motion; improvements are based on the local rendered evidence.

Merge and production deployment remain subject to user approval.
