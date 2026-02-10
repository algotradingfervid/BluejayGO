# Bluejay CMS — Playwright Test Plans Overview

## Test Environment
- **Base URL**: `http://localhost:28090`
- **Database**: SQLite (bluejay.db), seeded via `seed.sql`
- **Admin Credentials**:
  - Admin: `admin@bluejaylabs.com` / `password` (bcrypt hash in seed)
  - Editor: `editor@bluejaylabs.com` / `password` (bcrypt hash in seed)
- **Session Cookie**: `bluejay_session` (HttpOnly, SameSite=Lax, MaxAge=7 days)

## Test Plan Files

### Authentication & Core (01-02)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 01 | 01-admin-login-logout.md | Login, logout, session handling | Simple | None |
| 02 | 02-admin-dashboard.md | Dashboard stats & navigation | Simple | 01 |

### Master Table CRUD (03-08)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 03 | 03-product-categories-crud.md | Product categories | Simple | 01 |
| 04 | 04-blog-categories-crud.md | Blog categories | Simple | 01 |
| 05 | 05-blog-authors-crud.md | Blog authors | Simple | 01 |
| 06 | 06-industries-crud.md | Industries | Simple | 01 |
| 07 | 07-partner-tiers-crud.md | Partner tiers | Simple | 01 |
| 08 | 08-whitepaper-topics-crud.md | Whitepaper topics | Simple | 01 |

### Content CRUD — Main Entities (09-15)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 09 | 09-products-crud.md | Product listing, create, edit, delete | Complex | 01, 03 |
| 10 | 10-product-specs.md | Product specs (HTMX sub-resource) | Medium | 09 |
| 11 | 11-product-features.md | Product features (HTMX sub-resource) | Medium | 09 |
| 12 | 12-product-certifications.md | Product certifications (HTMX) | Medium | 09 |
| 13 | 13-product-downloads.md | Product file downloads (HTMX + upload) | Medium | 09 |
| 14 | 14-product-images.md | Product image gallery (HTMX + upload) | Medium | 09 |
| 15 | 15-blog-posts-crud.md | Blog posts listing, create, edit, delete | Complex | 01, 04, 05 |

### Content CRUD — Blog Sub-Features (16-18)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 16 | 16-blog-tags.md | Blog tag management & autocomplete | Medium | 01 |
| 17 | 17-blog-post-tags.md | Tag assignment in post editor (HTMX) | Medium | 15, 16 |
| 18 | 18-blog-post-products.md | Product linking in post editor (HTMX) | Medium | 15, 09 |

### Content CRUD — Solutions (19-23)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 19 | 19-solutions-crud.md | Solutions listing, create, edit, delete | Complex | 01 |
| 20 | 20-solution-stats.md | Solution stats (HTMX sub-resource) | Medium | 19 |
| 21 | 21-solution-challenges.md | Solution challenges (HTMX) | Medium | 19 |
| 22 | 22-solution-products.md | Solution product linking (HTMX) | Medium | 19, 09 |
| 23 | 23-solution-ctas.md | Solution CTAs (HTMX) | Medium | 19 |

### Content CRUD — Case Studies (24-26)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 24 | 24-case-studies-crud.md | Case studies listing, create, edit, delete | Complex | 01, 06 |
| 25 | 25-case-study-products.md | Case study product linking (HTMX) | Medium | 24, 09 |
| 26 | 26-case-study-metrics.md | Case study metrics (HTMX) | Medium | 24 |

### Content CRUD — Whitepapers (27-28)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 27 | 27-whitepapers-crud.md | Whitepapers listing, create, edit, delete | Complex | 01, 08 |
| 28 | 28-whitepaper-downloads.md | Whitepaper download analytics | Medium | 27 |

### Content CRUD — Partners (29-30)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 29 | 29-partners-crud.md | Partners listing, create, edit, delete | Medium | 01, 07 |
| 30 | 30-partner-testimonials.md | Partner testimonials CRUD | Medium | 29 |

### About Page Management (31-35)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 31 | 31-about-overview.md | Company overview edit | Simple | 01 |
| 32 | 32-about-mvv.md | Mission/Vision/Values edit | Simple | 01 |
| 33 | 33-about-core-values.md | Core values CRUD | Simple | 01 |
| 34 | 34-about-milestones.md | Company milestones CRUD | Simple | 01 |
| 35 | 35-about-certifications.md | Company certifications CRUD | Simple | 01 |

### Homepage Management (36-41)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 36 | 36-homepage-heroes.md | Hero sections CRUD | Medium | 01 |
| 37 | 37-homepage-stats.md | Homepage stats CRUD | Simple | 01 |
| 38 | 38-homepage-testimonials.md | Homepage testimonials CRUD | Medium | 01 |
| 39 | 39-homepage-ctas.md | Homepage CTAs CRUD | Medium | 01 |
| 40 | 40-homepage-settings.md | Homepage feature toggles & config | Simple | 01 |

### Site Configuration (41-46)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 41 | 41-global-settings.md | Global settings (tabs: general/seo/social) | Medium | 01 |
| 42 | 42-header-management.md | Header logo, nav toggles, CTA | Medium | 01 |
| 43 | 43-footer-management.md | Footer columns, links, legal | Complex | 01 |
| 44 | 44-page-sections.md | Page sections editor | Simple | 01 |
| 45 | 45-section-settings.md | Section-specific settings (about/products/solutions/blog) | Medium | 01 |
| 46 | 46-media-library.md | Media upload, browse, delete, alt text | Complex | 01 |

### Navigation & Activity (47-48)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 47 | 47-navigation-editor.md | Menu CRUD, items, drag-drop reorder | Complex | 01 |
| 48 | 48-activity-log.md | Activity log viewing & filtering | Simple | 01 |

### Contact Management (49-50)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 49 | 49-contact-submissions.md | Submissions list, view, status, bulk, delete | Medium | 01 |
| 50 | 50-office-locations.md | Office locations CRUD | Simple | 01 |

### Public Pages (51-62)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 51 | 51-public-homepage.md | Homepage sections, hero carousel, links | Medium | Seed data |
| 52 | 52-public-products-listing.md | Product listing, category filter, HTMX search | Medium | Seed data |
| 53 | 53-public-product-detail.md | Product detail page, specs, gallery | Medium | Seed data |
| 54 | 54-public-solutions-listing.md | Solutions listing page | Simple | Seed data |
| 55 | 55-public-solution-detail.md | Solution detail with stats, challenges | Medium | Seed data |
| 56 | 56-public-blog-listing.md | Blog listing, category/tag filtering | Medium | Seed data |
| 57 | 57-public-blog-post.md | Blog post detail, author, tags | Medium | Seed data |
| 58 | 58-public-case-studies.md | Case studies listing & detail | Medium | Seed data |
| 59 | 59-public-whitepapers.md | Whitepapers listing, detail, lead capture | Complex | Seed data |
| 60 | 60-public-about-page.md | About page sections | Simple | Seed data |
| 61 | 61-public-partners-page.md | Partners directory & testimonials | Simple | Seed data |
| 62 | 62-public-contact-page.md | Contact form submission & offices | Medium | Seed data |

### Search & SEO (63-65)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 63 | 63-public-search.md | Full-text search & HTMX suggestions | Medium | Seed data |
| 64 | 64-sitemap-robots.md | XML sitemap & robots.txt | Simple | Seed data |
| 65 | 65-public-navigation.md | Header/footer nav rendering | Simple | Seed data |

### Security & Infrastructure (66-69)
| # | File | Flow | Complexity | Dependencies |
|---|------|------|------------|--------------|
| 66 | 66-rate-limiting.md | Contact form rate limiting (5/hr) | Medium | None |
| 67 | 67-session-security.md | Session handling, cookie security | Medium | 01 |
| 68 | 68-security-headers.md | CSP, X-Frame-Options, etc. | Simple | None |
| 69 | 69-health-check.md | /health endpoint | Simple | None |

## Recommended Execution Order

**Phase 1 — Foundation:**
01 → 02

**Phase 2 — Master Tables (parallel):**
03, 04, 05, 06, 07, 08

**Phase 3 — Products & Sub-resources:**
09 → 10, 11, 12, 13, 14 (parallel)

**Phase 4 — Blog & Sub-features:**
15 → 16 → 17, 18 (parallel)

**Phase 5 — Solutions & Sub-resources:**
19 → 20, 21, 22, 23 (parallel)

**Phase 6 — Case Studies & Sub-resources:**
24 → 25, 26 (parallel)

**Phase 7 — Whitepapers:**
27 → 28

**Phase 8 — Partners:**
29 → 30

**Phase 9 — About Page (parallel):**
31, 32, 33, 34, 35

**Phase 10 — Homepage (parallel):**
36, 37, 38, 39, 40

**Phase 11 — Site Config (parallel):**
41, 42, 43, 44, 45, 46

**Phase 12 — Navigation, Activity, Contact:**
47, 48, 49, 50

**Phase 13 — Public Pages (parallel):**
51-65

**Phase 14 — Security & Infrastructure (parallel):**
66, 67, 68, 69

## Dependency Graph
```
01-login
├── 02-dashboard
├── 03-product-categories ──┐
├── 04-blog-categories ────┐│
├── 05-blog-authors ──────┐││
├── 06-industries ────────┐│││
├── 07-partner-tiers ────┐││││
├── 08-whitepaper-topics ┐│││││
│                        │││││└──> 09-products ──> 10,11,12,13,14
│                        ││││└───> 15-blog-posts ──> 17,18
│                        │││└────> 16-blog-tags ──> 17
│                        ││└─────> 19-solutions ──> 20,21,22,23
│                        │└──────> 24-case-studies ──> 25,26
│                        └───────> 27-whitepapers ──> 28
├── 29-partners ──> 30-partner-testimonials
├── 31-35 (about, parallel)
├── 36-40 (homepage, parallel)
├── 41-50 (config/contact, parallel)
├── 47-navigation
└── 48-activity-log

Public pages (51-65) depend on seed data, not admin flows.
Security tests (66-69) are independent.
```

## Seed Data Summary
| Entity | Count | Key Items |
|--------|-------|-----------|
| Admin Users | 2 | admin@bluejaylabs.com, editor@bluejaylabs.com |
| Product Categories | 5 | Desktops, OPS Modules, Interactive Flat Panels, AV Accessories, IoT Products |
| Blog Categories | 4 | Industry News, Product Updates, How-To Guides, Company Announcements |
| Blog Authors | 6 | Vikram Patel, Priya Mehta, Rahul Sharma, Anjali Desai, Arjun Singh, Tech Team |
| Industries | 6 | Education, Corporate, Healthcare, Retail, Government, Hospitality |
| Partner Tiers | 2 | Technology Partners, Channel Partners |
| Whitepaper Topics | 6 | Interactive Displays, Digital Classroom, Enterprise Computing, IoT Solutions, AV Technology, Industry Trends |
| Blog Tags | 13 | IFP, Education, EdTech, Classroom Technology, etc. |
| Products | 13 | BJ-D100 through BJ-IOT100 (with specs, features, certs, downloads, images) |
| Solutions | 6 | Education, Corporate, Healthcare, Retail, Government, Hospitality (with stats, challenges, products, CTAs) |
| Blog Posts | 7+ | Various topics with tags and product links |
| Case Studies | 3+ | ABC University, TechCorp, Metro Health (with metrics and products) |
| Whitepapers | 12 | 2 per topic, with learning points |
| Whitepaper Downloads | 29 | Sample lead capture data |
| Contact Submissions | 15 | Mix of sales, demo, partnership, support inquiries |
| Office Locations | 6 | Bangalore (HQ), Mumbai, Gurugram, Chennai, Hyderabad, Pune |
| Partners | 11 | Intel, Microsoft, Google, Samsung, Qualcomm, NVIDIA + 5 channel partners |
| Partner Testimonials | 6 | From various partners |
| Homepage Heroes | 1 | Main hero with CTAs |
| Homepage Stats | 4 | Company achievement numbers |
| Homepage Testimonials | 3 | Customer quotes with ratings |
| Homepage CTAs | 1 | Bottom CTA section |
| About | Full | Overview, MVV, 6 core values, milestones (2008-present), certifications |
