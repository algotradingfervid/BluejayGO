# Dependency refresh and merged release

The production fixes, both UX review rounds, and deployment-security changes are integrated into `main`.

Go 1.27.1 is the latest stable Go release verified against the official release feed on 29 September 2026. All modules declared in `go.mod` are updated to their latest available versions within their existing module paths. Notable updates include golang-migrate 4.20.1, modernc SQLite 1.60.0, Echo 4.16.0, crypto 0.57.0, net 0.59.0 and their used transitive dependencies.

HTMX is vendored at 2.0.11 (the upstream `latest` channel), and Trix is current at 2.1.19. Vendor metadata records registry sources and integrity hashes. Tailwind is pinned at 3.4.19, the latest v3 LTS release, with forms 0.5.11 and container queries 0.1.1. Styles are compiled locally and checked in; no Tailwind compiler runs from a third-party CDN in visitors' browsers. Run `npm ci --ignore-scripts && npm run build:css` after changing template classes. Deployment rebuilds styles before uploading.

This compatibility-focused refresh retains the supported Echo 4 and Tailwind 3 major versions. Echo 5 and Tailwind 4 require separate API/style and browser-support migrations; HTMX 4 is on upstream's `next` channel. They are not represented as upgraded here.

Browser testing found and fixed a Trix/form-state integration issue: the newly form-associated editor duplicated its backing textarea in dirty-state snapshots. Saved articles now start clean, actual edits trigger protection, and reverting text restores the clean state. Combining the branches also created a redundant mobile menu on the backup page; it now uses the shared accessible navigation.

Validation before GitHub push/deployment included the full application Go test suite and integration tests, vet, module verification, JavaScript tests, npm audit, and Linux govulncheck. Browser checks covered case-study publishing, blog empty-body validation and publishing, contact submission through the admin queue, media upload, backup download, mobile routes/navigation, and dirty-form recovery. The Linux release was tested under the hardened systemd profile against an isolated production backup, including migrations 44→48, database integrity, ten public and seven admin routes, secure login and backup range downloads.

The npm audit found no advisories. Go scanning found no reachable vulnerabilities and none in imported packages, while retaining one advisory in an unused part of a required module. These checks are bounded regression/security checks, not a guarantee that every possible defect or the remaining application-security audit findings have been resolved.
