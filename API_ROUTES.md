# Bluejay CMS API Routes Documentation

## Overview

Bluejay CMS uses the Echo v4 web framework for HTTP routing. Routes are organized into two main groups:

1. **Public Routes** (`publicGroup`) - Publicly accessible content pages
2. **Admin Routes** (`adminGroup`) - Protected admin panel routes requiring authentication

### Middleware Stack

**Global Middleware** (applied to all routes):
- `customMiddleware.Recovery()` - Panic recovery
- `customMiddleware.Logging()` - Request logging
- `middleware.Gzip()` - Response compression
- `customMiddleware.SecurityHeaders()` - Security headers (CSP, X-Frame-Options, etc.)
- `customMiddleware.SessionMiddleware()` - Session management

**Public Group Middleware**:
- `customMiddleware.SettingsLoader()` - Loads site settings into context

**Admin Group Middleware**:
- `customMiddleware.RequireAuth()` - Requires authentication (checks session)

## Static File Routes

| Method | Path | Description |
|--------|------|-------------|
| Static | `/public` | Serves static assets from `public/` directory |
| Static | `/uploads` | Serves uploaded files from `public/uploads/` directory |

---

## Public Routes

All public routes use the `publicGroup` with `SettingsLoader` middleware.

### Health & Utility

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/health` | Inline handler | JSON | API | Health check endpoint, returns status and timestamp |

### Homepage

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/` | `homeHandler.ShowHomePage` | `public/pages/home.html` | Full Page | Homepage with heroes, stats, testimonials, CTA, featured products |

### Products

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/products` | `productsHandler.ProductsList` | `public/pages/products.html` | Full Page | Products landing page with categories |
| GET | `/products/search` | `productsHandler.ProductSearch` | `public/partials/product_search_results.html` | HTMX Fragment | HTMX search results for products |
| GET | `/products/:category` | `productsHandler.ProductsByCategory` | `public/pages/products_category.html` | Full Page | Products filtered by category |
| GET | `/products/:category/:slug` | `productsHandler.ProductDetail` | `public/pages/product_detail.html` | Full Page | Individual product detail page |

### Solutions

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/solutions` | `solutionsHandler.SolutionsList` | `public/pages/solutions.html` | Full Page | Solutions listing page |
| GET | `/solutions/:slug` | `solutionsHandler.SolutionDetail` | `public/pages/solution_detail.html` | Full Page | Individual solution detail page |

### Blog

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/blog` | `blogHandler.BlogListing` | `public/pages/blog.html` | Full Page | Blog post listing page with filters |
| GET | `/blog/:slug` | `blogHandler.BlogPost` | `public/pages/blog_post.html` | Full Page | Individual blog post detail page |

### Case Studies

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/case-studies` | `caseStudiesHandler.CaseStudiesList` | `public/pages/case_studies.html` | Full Page | Case studies listing page |
| GET | `/case-studies/:slug` | `caseStudiesHandler.CaseStudyDetail` | `public/pages/case_study_detail.html` | Full Page | Individual case study detail page |

### Whitepapers

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/whitepapers` | `whitepapersHandler.WhitepapersList` | `public/pages/whitepapers.html` | Full Page | Whitepapers listing page |
| GET | `/whitepapers/:slug` | `whitepapersHandler.WhitepaperDetail` | `public/pages/whitepaper_detail.html` | Full Page | Individual whitepaper detail page with download form |
| POST | `/whitepapers/:slug/download` | `whitepapersHandler.WhitepaperDownload` | N/A | Form Submit | Processes whitepaper download lead capture |

### About

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/about` | `aboutHandler.AboutPage` | `public/pages/about.html` | Full Page | About page with company overview, mission/vision, values, milestones |

### Partners

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/partners` | `partnersPageHandler.PartnersPage` | `public/pages/partners.html` | Full Page | Partners page with partner listing by tier |

### Contact

| Method | Path | Handler | Template | Type | Description | Rate Limited |
|--------|------|---------|----------|------|-------------|--------------|
| GET | `/contact` | `contactHandler.ShowContactPage` | `public/pages/contact.html` | Full Page | Contact page with office locations and form | No |
| POST | `/contact/submit` | `contactHandler.SubmitContactForm` | N/A | Form Submit | Processes contact form submission | **Yes** (5 per hour) |

### Search & SEO

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/search` | `searchHandler.SearchPage` | `public/pages/search.html` | Full Page | Global search results page |
| GET | `/search/suggest` | `searchHandler.SearchSuggest` | `public/partials/search_suggestions.html` | HTMX Fragment | HTMX autocomplete suggestions |
| GET | `/sitemap.xml` | `sitemapHandler.Sitemap` | XML | XML | XML sitemap for SEO |
| GET | `/robots.txt` | `sitemapHandler.RobotsTxt` | Text | Text | Robots.txt file |

---

## Admin Authentication Routes

These routes are in the `adminAuthGroup` (no auth required for login page).

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/login` | `authHandler.ShowLoginPage` | `admin/pages/login.html` | Full Page | Admin login form |
| POST | `/admin/login` | `authHandler.LoginSubmit` | N/A | Form Submit | Processes login, creates session |
| POST | `/admin/logout` | `authHandler.Logout` | N/A | Form Submit | Destroys session, redirects to login |

---

## Admin Dashboard

All admin routes require authentication via `RequireAuth()` middleware.

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/dashboard` | `dashboardHandler.ShowDashboard` | `admin/pages/dashboard.html` | Full Page | Admin dashboard with statistics |

---

## Admin Master Tables CRUD

### Product Categories

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/product-categories` | `pcHandler.List` | `admin/pages/product_categories_list.html` | Full Page | List all product categories |
| GET | `/admin/product-categories/new` | `pcHandler.New` | `admin/pages/product_categories_form.html` | Full Page | New category form |
| POST | `/admin/product-categories` | `pcHandler.Create` | N/A | Form Submit | Create new category |
| GET | `/admin/product-categories/:id/edit` | `pcHandler.Edit` | `admin/pages/product_categories_form.html` | Full Page | Edit category form |
| POST | `/admin/product-categories/:id` | `pcHandler.Update` | N/A | Form Submit | Update category |
| DELETE | `/admin/product-categories/:id` | `pcHandler.Delete` | N/A | HTMX | Delete category |

### Blog Categories

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/blog-categories` | `bcHandler.List` | `admin/pages/blog_categories_list.html` | Full Page | List all blog categories |
| GET | `/admin/blog-categories/new` | `bcHandler.New` | `admin/pages/blog_categories_form.html` | Full Page | New category form |
| POST | `/admin/blog-categories` | `bcHandler.Create` | N/A | Form Submit | Create new category |
| GET | `/admin/blog-categories/:id/edit` | `bcHandler.Edit` | `admin/pages/blog_categories_form.html` | Full Page | Edit category form |
| POST | `/admin/blog-categories/:id` | `bcHandler.Update` | N/A | Form Submit | Update category |
| DELETE | `/admin/blog-categories/:id` | `bcHandler.Delete` | N/A | HTMX | Delete category |

### Blog Authors

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/blog-authors` | `baHandler.List` | `admin/pages/blog_authors_list.html` | Full Page | List all blog authors |
| GET | `/admin/blog-authors/new` | `baHandler.New` | `admin/pages/blog_authors_form.html` | Full Page | New author form |
| POST | `/admin/blog-authors` | `baHandler.Create` | N/A | Form Submit | Create new author |
| GET | `/admin/blog-authors/:id/edit` | `baHandler.Edit` | `admin/pages/blog_authors_form.html` | Full Page | Edit author form |
| POST | `/admin/blog-authors/:id` | `baHandler.Update` | N/A | Form Submit | Update author |
| DELETE | `/admin/blog-authors/:id` | `baHandler.Delete` | N/A | HTMX | Delete author |

### Industries

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/industries` | `indHandler.List` | `admin/pages/industries_list.html` | Full Page | List all industries |
| GET | `/admin/industries/new` | `indHandler.New` | `admin/pages/industries_form.html` | Full Page | New industry form |
| POST | `/admin/industries` | `indHandler.Create` | N/A | Form Submit | Create new industry |
| GET | `/admin/industries/:id/edit` | `indHandler.Edit` | `admin/pages/industries_form.html` | Full Page | Edit industry form |
| POST | `/admin/industries/:id` | `indHandler.Update` | N/A | Form Submit | Update industry |
| DELETE | `/admin/industries/:id` | `indHandler.Delete` | N/A | HTMX | Delete industry |

### Partner Tiers

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/partner-tiers` | `ptHandler.List` | `admin/pages/partner_tiers_list.html` | Full Page | List all partner tiers |
| GET | `/admin/partner-tiers/new` | `ptHandler.New` | `admin/pages/partner_tiers_form.html` | Full Page | New tier form |
| POST | `/admin/partner-tiers` | `ptHandler.Create` | N/A | Form Submit | Create new tier |
| GET | `/admin/partner-tiers/:id/edit` | `ptHandler.Edit` | `admin/pages/partner_tiers_form.html` | Full Page | Edit tier form |
| POST | `/admin/partner-tiers/:id` | `ptHandler.Update` | N/A | Form Submit | Update tier |
| DELETE | `/admin/partner-tiers/:id` | `ptHandler.Delete` | N/A | HTMX | Delete tier |

### Whitepaper Topics

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/whitepaper-topics` | `wtHandler.List` | `admin/pages/whitepaper_topics_list.html` | Full Page | List all whitepaper topics |
| GET | `/admin/whitepaper-topics/new` | `wtHandler.New` | `admin/pages/whitepaper_topics_form.html` | Full Page | New topic form |
| POST | `/admin/whitepaper-topics` | `wtHandler.Create` | N/A | Form Submit | Create new topic |
| GET | `/admin/whitepaper-topics/:id/edit` | `wtHandler.Edit` | `admin/pages/whitepaper_topics_form.html` | Full Page | Edit topic form |
| POST | `/admin/whitepaper-topics/:id` | `wtHandler.Update` | N/A | Form Submit | Update topic |
| DELETE | `/admin/whitepaper-topics/:id` | `wtHandler.Delete` | N/A | HTMX | Delete topic |

---

## Admin Products

### Products CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/products` | `adminProductsHandler.List` | `admin/pages/products_list.html` | Full Page | List products with filtering and pagination |
| GET | `/admin/products/new` | `adminProductsHandler.New` | `admin/pages/products_form.html` | Full Page | New product form |
| POST | `/admin/products` | `adminProductsHandler.Create` | N/A | Form Submit | Create new product |
| GET | `/admin/products/:id/edit` | `adminProductsHandler.Edit` | `admin/pages/products_form.html` | Full Page | Edit product form |
| POST | `/admin/products/:id` | `adminProductsHandler.Update` | N/A | Form Submit | Update product |
| DELETE | `/admin/products/:id` | `adminProductsHandler.Delete` | N/A | HTMX | Delete product |

### Product Sub-Entities (HTMX Fragments)

**Product Specifications:**

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/products/:id/specs` | `pdHandler.ListSpecs` | `admin/partials/product_specs.html` | HTMX Fragment | Get specs list |
| POST | `/admin/products/:id/specs` | `pdHandler.AddSpec` | `admin/partials/product_specs.html` | HTMX Fragment | Add spec, returns updated list |
| DELETE | `/admin/products/:id/specs` | `pdHandler.DeleteSpecs` | `admin/partials/product_specs.html` | HTMX Fragment | Delete all specs, returns empty list |

**Product Features:**

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/products/:id/features` | `pdHandler.ListFeatures` | `admin/partials/product_features.html` | HTMX Fragment | Get features list |
| POST | `/admin/products/:id/features` | `pdHandler.AddFeature` | `admin/partials/product_features.html` | HTMX Fragment | Add feature, returns updated list |
| DELETE | `/admin/products/:id/features` | `pdHandler.DeleteFeatures` | `admin/partials/product_features.html` | HTMX Fragment | Delete all features, returns empty list |

**Product Certifications:**

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/products/:id/certifications` | `pdHandler.ListCertifications` | `admin/partials/product_certifications.html` | HTMX Fragment | Get certifications list |
| POST | `/admin/products/:id/certifications` | `pdHandler.AddCertification` | `admin/partials/product_certifications.html` | HTMX Fragment | Add certification, returns updated list |
| DELETE | `/admin/products/:id/certifications` | `pdHandler.DeleteCertifications` | `admin/partials/product_certifications.html` | HTMX Fragment | Delete all certifications |

**Product Downloads:**

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/products/:id/downloads` | `pdHandler.ListDownloads` | `admin/partials/product_downloads.html` | HTMX Fragment | Get downloads list |
| POST | `/admin/products/:id/downloads` | `pdHandler.AddDownload` | `admin/partials/product_downloads.html` | HTMX Fragment | Upload file, returns updated list |
| DELETE | `/admin/products/:id/downloads/:download_id` | `pdHandler.DeleteDownload` | `admin/partials/product_downloads.html` | HTMX Fragment | Delete download, returns updated list |

**Product Images:**

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/products/:id/images` | `pdHandler.ListImages` | `admin/partials/product_images.html` | HTMX Fragment | Get images list |
| POST | `/admin/products/:id/images` | `pdHandler.AddImage` | `admin/partials/product_images.html` | HTMX Fragment | Upload image, returns updated list |
| DELETE | `/admin/products/:id/images/:image_id` | `pdHandler.DeleteImage` | `admin/partials/product_images.html` | HTMX Fragment | Delete image, returns updated list |

---

## Admin Blog

### Blog Posts CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/blog/posts` | `adminBlogPostsHandler.List` | `admin/pages/blog_posts_list.html` | Full Page | List blog posts with filtering and pagination |
| GET | `/admin/blog/posts/new` | `adminBlogPostsHandler.New` | `admin/pages/blog_post_form.html` | Full Page | New blog post form |
| POST | `/admin/blog/posts` | `adminBlogPostsHandler.Create` | N/A | Form Submit | Create new blog post |
| GET | `/admin/blog/posts/:id/edit` | `adminBlogPostsHandler.Edit` | `admin/pages/blog_post_form.html` | Full Page | Edit blog post form |
| POST | `/admin/blog/posts/:id` | `adminBlogPostsHandler.Update` | N/A | Form Submit | Update blog post |
| DELETE | `/admin/blog/posts/:id` | `adminBlogPostsHandler.Delete` | N/A | HTMX | Delete blog post |
| GET | `/admin/blog/products/search` | `adminBlogPostsHandler.SearchProducts` | `admin/partials/product_suggestions.html` | HTMX Fragment | HTMX product search for linking |

### Blog Tags

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/blog/tags` | `adminBlogTagsHandler.List` | `admin/pages/blog_tags_list.html` | Full Page | List all blog tags |
| POST | `/admin/blog/tags` | `adminBlogTagsHandler.Create` | N/A | Form Submit | Create new tag |
| GET | `/admin/blog/tags/search` | `adminBlogTagsHandler.Search` | `admin/partials/tag_suggestions.html` | HTMX Fragment | HTMX tag search autocomplete |
| POST | `/admin/blog/tags/quick-create` | `adminBlogTagsHandler.QuickCreate` | `admin/partials/tag_chip.html` | HTMX Fragment | Quick create tag from form, returns chip |
| DELETE | `/admin/blog/tags/:id` | `adminBlogTagsHandler.Delete` | N/A | HTMX | Delete tag |

---

## Admin Solutions

### Solutions CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/solutions` | `adminSolutionsHandler.List` | `admin/pages/solutions_list.html` | Full Page | List solutions with filtering and pagination |
| GET | `/admin/solutions/new` | `adminSolutionsHandler.New` | `admin/pages/solutions_form.html` | Full Page | New solution form |
| POST | `/admin/solutions` | `adminSolutionsHandler.Create` | N/A | Form Submit | Create new solution |
| GET | `/admin/solutions/:id/edit` | `adminSolutionsHandler.Edit` | `admin/pages/solutions_form.html` | Full Page | Edit solution form |
| POST | `/admin/solutions/:id` | `adminSolutionsHandler.Update` | N/A | Form Submit | Update solution |
| DELETE | `/admin/solutions/:id` | `adminSolutionsHandler.Delete` | N/A | HTMX | Delete solution |

### Solution Sub-Entities (HTMX Fragments)

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| POST | `/admin/solutions/:id/stats` | `adminSolutionsHandler.AddStat` | `admin/partials/solution_stats.html` | HTMX Fragment | Add stat, returns updated list |
| DELETE | `/admin/solutions/:id/stats/:statId` | `adminSolutionsHandler.DeleteStat` | N/A | HTMX | Delete stat |
| POST | `/admin/solutions/:id/challenges` | `adminSolutionsHandler.AddChallenge` | `admin/partials/solution_challenges.html` | HTMX Fragment | Add challenge, returns updated list |
| DELETE | `/admin/solutions/:id/challenges/:challengeId` | `adminSolutionsHandler.DeleteChallenge` | N/A | HTMX | Delete challenge |
| POST | `/admin/solutions/:id/products` | `adminSolutionsHandler.AddProduct` | `admin/partials/solution_products.html` | HTMX Fragment | Link product, returns updated list |
| DELETE | `/admin/solutions/:id/products/:productId` | `adminSolutionsHandler.RemoveProduct` | N/A | HTMX | Unlink product |
| POST | `/admin/solutions/:id/ctas` | `adminSolutionsHandler.AddCTA` | `admin/partials/solution_ctas.html` | HTMX Fragment | Add CTA, returns updated list |
| DELETE | `/admin/solutions/:id/ctas/:ctaId` | `adminSolutionsHandler.DeleteCTA` | N/A | HTMX | Delete CTA |

---

## Admin Case Studies

### Case Studies CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/case-studies` | `adminCaseStudiesHandler.List` | `admin/pages/case_studies_list.html` | Full Page | List case studies with filtering and pagination |
| GET | `/admin/case-studies/new` | `adminCaseStudiesHandler.New` | `admin/pages/case_studies_form.html` | Full Page | New case study form |
| POST | `/admin/case-studies` | `adminCaseStudiesHandler.Create` | N/A | Form Submit | Create new case study |
| GET | `/admin/case-studies/:id/edit` | `adminCaseStudiesHandler.Edit` | `admin/pages/case_studies_form.html` | Full Page | Edit case study form |
| POST | `/admin/case-studies/:id` | `adminCaseStudiesHandler.Update` | N/A | Form Submit | Update case study |
| DELETE | `/admin/case-studies/:id` | `adminCaseStudiesHandler.Delete` | N/A | HTMX | Delete case study |

### Case Study Sub-Entities (HTMX Fragments)

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| POST | `/admin/case-studies/:id/products` | `adminCaseStudiesHandler.AddProduct` | `admin/partials/case_study_products.html` | HTMX Fragment | Link product, returns updated list |
| DELETE | `/admin/case-studies/:id/products/:productId` | `adminCaseStudiesHandler.RemoveProduct` | N/A | HTMX | Unlink product |
| POST | `/admin/case-studies/:id/metrics` | `adminCaseStudiesHandler.AddMetric` | `admin/partials/case_study_metrics.html` | HTMX Fragment | Add metric, returns updated list |
| DELETE | `/admin/case-studies/:id/metrics/:metricId` | `adminCaseStudiesHandler.DeleteMetric` | N/A | HTMX | Delete metric |

---

## Admin Whitepapers

### Whitepapers CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/whitepapers` | `adminWhitepapersHandler.List` | `admin/pages/whitepapers_list.html` | Full Page | List whitepapers with filtering and pagination |
| GET | `/admin/whitepapers/new` | `adminWhitepapersHandler.New` | `admin/pages/whitepapers_form.html` | Full Page | New whitepaper form |
| POST | `/admin/whitepapers` | `adminWhitepapersHandler.Create` | N/A | Form Submit | Create new whitepaper (with PDF upload) |
| GET | `/admin/whitepapers/:id/edit` | `adminWhitepapersHandler.Edit` | `admin/pages/whitepapers_form.html` | Full Page | Edit whitepaper form |
| POST | `/admin/whitepapers/:id` | `adminWhitepapersHandler.Update` | N/A | Form Submit | Update whitepaper (optional PDF replacement) |
| DELETE | `/admin/whitepapers/:id` | `adminWhitepapersHandler.Delete` | N/A | HTMX | Delete whitepaper |
| GET | `/admin/whitepapers/:id/downloads` | `adminWhitepapersHandler.Downloads` | `admin/pages/whitepapers_downloads.html` | Full Page | View whitepaper download leads |

---

## Admin About Page

### Company Overview & Mission/Vision/Values

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/about/overview` | `adminAboutHandler.OverviewEdit` | `admin/pages/about_overview_form.html` | Full Page | Edit company overview |
| POST | `/admin/about/overview` | `adminAboutHandler.OverviewUpdate` | N/A | Form Submit | Update company overview |
| GET | `/admin/about/mvv` | `adminAboutHandler.MVVEdit` | `admin/pages/about_mvv_form.html` | Full Page | Edit mission/vision/values |
| POST | `/admin/about/mvv` | `adminAboutHandler.MVVUpdate` | N/A | Form Submit | Update mission/vision/values |

### Core Values CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/about/values` | `adminAboutHandler.CoreValuesList` | `admin/pages/core_values_list.html` | Full Page | List core values |
| GET | `/admin/about/values/new` | `adminAboutHandler.CoreValueNew` | `admin/pages/core_values_form.html` | Full Page | New core value form |
| POST | `/admin/about/values` | `adminAboutHandler.CoreValueCreate` | N/A | Form Submit | Create core value |
| GET | `/admin/about/values/:id/edit` | `adminAboutHandler.CoreValueEdit` | `admin/pages/core_values_form.html` | Full Page | Edit core value form |
| POST | `/admin/about/values/:id` | `adminAboutHandler.CoreValueUpdate` | N/A | Form Submit | Update core value |
| DELETE | `/admin/about/values/:id` | `adminAboutHandler.CoreValueDelete` | N/A | HTMX | Delete core value |

### Milestones CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/about/milestones` | `adminAboutHandler.MilestonesList` | `admin/pages/milestones_list.html` | Full Page | List milestones |
| GET | `/admin/about/milestones/new` | `adminAboutHandler.MilestoneNew` | `admin/pages/milestones_form.html` | Full Page | New milestone form |
| POST | `/admin/about/milestones` | `adminAboutHandler.MilestoneCreate` | N/A | Form Submit | Create milestone |
| GET | `/admin/about/milestones/:id/edit` | `adminAboutHandler.MilestoneEdit` | `admin/pages/milestones_form.html` | Full Page | Edit milestone form |
| POST | `/admin/about/milestones/:id` | `adminAboutHandler.MilestoneUpdate` | N/A | Form Submit | Update milestone |
| DELETE | `/admin/about/milestones/:id` | `adminAboutHandler.MilestoneDelete` | N/A | HTMX | Delete milestone |

### About Certifications CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/about/certifications` | `adminAboutHandler.CertificationsList` | `admin/pages/certifications_list.html` | Full Page | List certifications |
| GET | `/admin/about/certifications/new` | `adminAboutHandler.CertificationNew` | `admin/pages/certifications_form.html` | Full Page | New certification form |
| POST | `/admin/about/certifications` | `adminAboutHandler.CertificationCreate` | N/A | Form Submit | Create certification |
| GET | `/admin/about/certifications/:id/edit` | `adminAboutHandler.CertificationEdit` | `admin/pages/certifications_form.html` | Full Page | Edit certification form |
| POST | `/admin/about/certifications/:id` | `adminAboutHandler.CertificationUpdate` | N/A | Form Submit | Update certification |
| DELETE | `/admin/about/certifications/:id` | `adminAboutHandler.CertificationDelete` | N/A | HTMX | Delete certification |

---

## Admin Partners

### Partners CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/partners` | `adminPartnersHandler.List` | `admin/pages/partners_list.html` | Full Page | List partners |
| GET | `/admin/partners/new` | `adminPartnersHandler.New` | `admin/pages/partners_form.html` | Full Page | New partner form |
| POST | `/admin/partners` | `adminPartnersHandler.Create` | N/A | Form Submit | Create partner |
| GET | `/admin/partners/:id/edit` | `adminPartnersHandler.Edit` | `admin/pages/partners_form.html` | Full Page | Edit partner form |
| POST | `/admin/partners/:id` | `adminPartnersHandler.Update` | N/A | Form Submit | Update partner |
| DELETE | `/admin/partners/:id` | `adminPartnersHandler.Delete` | N/A | HTMX | Delete partner |

### Partner Testimonials CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/partners/testimonials` | `adminPartnersHandler.TestimonialsList` | `admin/pages/partner_testimonials_list.html` | Full Page | List partner testimonials |
| GET | `/admin/partners/testimonials/new` | `adminPartnersHandler.TestimonialNew` | `admin/pages/partner_testimonials_form.html` | Full Page | New testimonial form |
| POST | `/admin/partners/testimonials` | `adminPartnersHandler.TestimonialCreate` | N/A | Form Submit | Create testimonial |
| GET | `/admin/partners/testimonials/:id/edit` | `adminPartnersHandler.TestimonialEdit` | `admin/pages/partner_testimonials_form.html` | Full Page | Edit testimonial form |
| POST | `/admin/partners/testimonials/:id` | `adminPartnersHandler.TestimonialUpdate` | N/A | Form Submit | Update testimonial |
| DELETE | `/admin/partners/testimonials/:id` | `adminPartnersHandler.TestimonialDelete` | N/A | HTMX | Delete testimonial |

---

## Admin Homepage Management

### Homepage Heroes CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/homepage/heroes` | `homepageAdminHandler.HeroesList` | `admin/pages/homepage_heroes_list.html` | Full Page | List heroes |
| GET | `/admin/homepage/heroes/new` | `homepageAdminHandler.HeroNew` | `admin/pages/homepage_hero_form.html` | Full Page | New hero form |
| POST | `/admin/homepage/heroes` | `homepageAdminHandler.HeroCreate` | N/A | Form Submit | Create hero |
| GET | `/admin/homepage/heroes/:id/edit` | `homepageAdminHandler.HeroEdit` | `admin/pages/homepage_hero_form.html` | Full Page | Edit hero form |
| POST | `/admin/homepage/heroes/:id` | `homepageAdminHandler.HeroUpdate` | N/A | Form Submit | Update hero |
| DELETE | `/admin/homepage/heroes/:id` | `homepageAdminHandler.HeroDelete` | N/A | HTMX | Delete hero |

### Homepage Stats CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/homepage/stats` | `homepageAdminHandler.StatsList` | `admin/pages/homepage_stats_list.html` | Full Page | List stats |
| GET | `/admin/homepage/stats/new` | `homepageAdminHandler.StatNew` | `admin/pages/homepage_stat_form.html` | Full Page | New stat form |
| POST | `/admin/homepage/stats` | `homepageAdminHandler.StatCreate` | N/A | Form Submit | Create stat |
| GET | `/admin/homepage/stats/:id/edit` | `homepageAdminHandler.StatEdit` | `admin/pages/homepage_stat_form.html` | Full Page | Edit stat form |
| POST | `/admin/homepage/stats/:id` | `homepageAdminHandler.StatUpdate` | N/A | Form Submit | Update stat |
| DELETE | `/admin/homepage/stats/:id` | `homepageAdminHandler.StatDelete` | N/A | HTMX | Delete stat |

### Homepage Testimonials CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/homepage/testimonials` | `homepageAdminHandler.TestimonialsList` | `admin/pages/homepage_testimonials_list.html` | Full Page | List testimonials |
| GET | `/admin/homepage/testimonials/new` | `homepageAdminHandler.TestimonialNew` | `admin/pages/homepage_testimonial_form.html` | Full Page | New testimonial form |
| POST | `/admin/homepage/testimonials` | `homepageAdminHandler.TestimonialCreate` | N/A | Form Submit | Create testimonial |
| GET | `/admin/homepage/testimonials/:id/edit` | `homepageAdminHandler.TestimonialEdit` | `admin/pages/homepage_testimonial_form.html` | Full Page | Edit testimonial form |
| POST | `/admin/homepage/testimonials/:id` | `homepageAdminHandler.TestimonialUpdate` | N/A | Form Submit | Update testimonial |
| DELETE | `/admin/homepage/testimonials/:id` | `homepageAdminHandler.TestimonialDelete` | N/A | HTMX | Delete testimonial |

### Homepage CTAs CRUD

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/homepage/cta` | `homepageAdminHandler.CTAList` | `admin/pages/homepage_cta_list.html` | Full Page | List CTAs |
| GET | `/admin/homepage/cta/new` | `homepageAdminHandler.CTANew` | `admin/pages/homepage_cta_form.html` | Full Page | New CTA form |
| POST | `/admin/homepage/cta` | `homepageAdminHandler.CTACreate` | N/A | Form Submit | Create CTA |
| GET | `/admin/homepage/cta/:id/edit` | `homepageAdminHandler.CTAEdit` | `admin/pages/homepage_cta_form.html` | Full Page | Edit CTA form |
| POST | `/admin/homepage/cta/:id` | `homepageAdminHandler.CTAUpdate` | N/A | Form Submit | Update CTA |
| DELETE | `/admin/homepage/cta/:id` | `homepageAdminHandler.CTADelete` | N/A | HTMX | Delete CTA |

### Homepage Settings

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/homepage/settings` | `homepageAdminHandler.Settings` | `admin/pages/homepage_settings.html` | Full Page | Edit homepage settings (show/hide sections, limits) |
| POST | `/admin/homepage/settings` | `homepageAdminHandler.UpdateSettings` | N/A | Form Submit | Update homepage settings |

---

## Admin Section Settings

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/about/settings` | `sectionSettingsHandler.AboutSettings` | `admin/pages/about_settings.html` | Full Page | About section settings |
| POST | `/admin/about/settings` | `sectionSettingsHandler.UpdateAboutSettings` | N/A | Form Submit | Update about settings |
| GET | `/admin/products/settings` | `sectionSettingsHandler.ProductsSettings` | `admin/pages/products_settings.html` | Full Page | Products section settings |
| POST | `/admin/products/settings` | `sectionSettingsHandler.UpdateProductsSettings` | N/A | Form Submit | Update products settings |
| GET | `/admin/solutions/settings` | `sectionSettingsHandler.SolutionsSettings` | `admin/pages/solutions_settings.html` | Full Page | Solutions section settings |
| POST | `/admin/solutions/settings` | `sectionSettingsHandler.UpdateSolutionsSettings` | N/A | Form Submit | Update solutions settings |
| GET | `/admin/blog/settings` | `sectionSettingsHandler.BlogSettings` | `admin/pages/blog_settings.html` | Full Page | Blog section settings |
| POST | `/admin/blog/settings` | `sectionSettingsHandler.UpdateBlogSettings` | N/A | Form Submit | Update blog settings |

---

## Admin Site-Wide Settings

### Global Settings

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/settings` | `settingsHandler.Edit` | `admin/pages/settings_form.html` | Full Page | Edit global settings (site name, contact info, social links, analytics) |
| POST | `/admin/settings` | `settingsHandler.Update` | N/A | Form Submit | Update global settings |

### Header Settings

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/header` | `headerHandler.Edit` | `admin/pages/header_form.html` | Full Page | Edit header settings (logo, CTA, navigation labels) |
| POST | `/admin/header` | `headerHandler.Update` | N/A | Form Submit | Update header settings |

### Footer Settings

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/footer` | `footerHandler.Edit` | `admin/pages/footer_form.html` | Full Page | Edit footer settings (columns, legal links, social) |
| POST | `/admin/footer` | `footerHandler.Update` | N/A | Form Submit | Update footer settings |

### Page Sections

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/page-sections` | `psHandler.List` | `admin/pages/page_sections_list.html` | Full Page | List all editable page sections |
| GET | `/admin/page-sections/:id/edit` | `psHandler.Edit` | `admin/pages/page_sections_form.html` | Full Page | Edit page section (headings, buttons, etc.) |
| POST | `/admin/page-sections/:id` | `psHandler.Update` | N/A | Form Submit | Update page section |

---

## Admin Media Library

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/media` | `mediaHandler.List` | `admin/pages/media_library.html` | Full Page | Media library with search, filtering, pagination |
| POST | `/admin/media/upload` | `mediaHandler.Upload` | JSON | API | Upload multiple files, returns JSON array of media files |
| GET | `/admin/media/browse` | `mediaHandler.Browse` | `admin/partials/media_picker.html` | HTMX Fragment | Media picker modal for selecting files |
| GET | `/admin/media/:id` | `mediaHandler.GetFile` | JSON | API | Get single media file metadata as JSON |
| PUT | `/admin/media/:id` | `mediaHandler.UpdateAltText` | JSON | API | Update file alt text, returns JSON |
| DELETE | `/admin/media/:id` | `mediaHandler.Delete` | N/A | HTMX | Delete media file |

---

## Admin Navigation Management

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/navigation` | `navHandler.List` | `admin/pages/navigation_list.html` | Full Page | List all navigation menus |
| POST | `/admin/navigation` | `navHandler.Create` | N/A | Form Submit | Create new navigation menu |
| GET | `/admin/navigation/:id` | `navHandler.Edit` | `admin/pages/navigation_editor.html` | Full Page | Edit navigation menu (drag-drop items) |
| POST | `/admin/navigation/:id/settings` | `navHandler.UpdateMenu` | N/A | Form Submit | Update menu settings |
| POST | `/admin/navigation/:id/items` | `navHandler.AddItem` | N/A | HTMX | Add navigation item to menu |
| POST | `/admin/navigation/items/:id` | `navHandler.UpdateItem` | N/A | HTMX | Update navigation item |
| DELETE | `/admin/navigation/items/:id` | `navHandler.DeleteItem` | N/A | HTMX | Delete navigation item |
| DELETE | `/admin/navigation/:id` | `navHandler.DeleteMenu` | N/A | HTMX | Delete entire navigation menu |
| POST | `/admin/navigation/:id/reorder` | `navHandler.Reorder` | N/A | HTMX | Reorder navigation items (drag-drop) |

---

## Admin Contact Management

### Contact Submissions

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/contact/submissions` | `adminContactHandler.ListSubmissions` | `admin/pages/contact_submissions_list.html` | Full Page | List contact form submissions with filtering |
| GET | `/admin/contact/submissions/:id` | `adminContactHandler.ViewSubmission` | `admin/pages/contact_submission_detail.html` | Full Page | View single submission details |
| POST | `/admin/contact/submissions/:id/status` | `adminContactHandler.UpdateSubmissionStatus` | N/A | HTMX | Update submission status (new/read/replied) |
| POST | `/admin/contact/submissions/bulk-mark-read` | `adminContactHandler.BulkMarkRead` | N/A | HTMX | Mark multiple submissions as read |
| DELETE | `/admin/contact/submissions/:id` | `adminContactHandler.DeleteSubmission` | N/A | HTMX | Delete submission |

### Office Locations

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/contact/offices` | `adminContactHandler.ListOffices` | `admin/pages/office_locations_list.html` | Full Page | List office locations |
| GET | `/admin/contact/offices/new` | `adminContactHandler.NewOffice` | `admin/pages/office_locations_form.html` | Full Page | New office form |
| POST | `/admin/contact/offices` | `adminContactHandler.CreateOffice` | N/A | Form Submit | Create office location |
| GET | `/admin/contact/offices/:id/edit` | `adminContactHandler.EditOffice` | `admin/pages/office_locations_form.html` | Full Page | Edit office form |
| POST | `/admin/contact/offices/:id` | `adminContactHandler.UpdateOffice` | N/A | Form Submit | Update office location |
| DELETE | `/admin/contact/offices/:id` | `adminContactHandler.DeleteOffice` | N/A | HTMX | Delete office location |

---

## Admin Activity Log

| Method | Path | Handler | Template | Type | Description |
|--------|------|---------|----------|------|-------------|
| GET | `/admin/activity` | `activityHandler.List` | `admin/pages/activity_log.html` | Full Page | Activity log with filtering and pagination |

---

## HTMX-Specific Endpoints Summary

The following endpoints return HTML fragments (partials) instead of full pages, designed for use with HTMX's `hx-get`, `hx-post`, or `hx-delete` attributes:

### Product Management Fragments
- `GET /admin/products/:id/specs` - Product specifications list
- `POST /admin/products/:id/specs` - Add specification
- `DELETE /admin/products/:id/specs` - Delete all specifications
- `GET /admin/products/:id/features` - Product features list
- `POST /admin/products/:id/features` - Add feature
- `DELETE /admin/products/:id/features` - Delete all features
- `GET /admin/products/:id/certifications` - Product certifications list
- `POST /admin/products/:id/certifications` - Add certification
- `DELETE /admin/products/:id/certifications` - Delete all certifications
- `GET /admin/products/:id/downloads` - Product downloads list
- `POST /admin/products/:id/downloads` - Add download
- `DELETE /admin/products/:id/downloads/:download_id` - Delete download
- `GET /admin/products/:id/images` - Product images list
- `POST /admin/products/:id/images` - Add image
- `DELETE /admin/products/:id/images/:image_id` - Delete image

### Blog Management Fragments
- `GET /admin/blog/products/search` - Product search autocomplete
- `GET /admin/blog/tags/search` - Tag search autocomplete
- `POST /admin/blog/tags/quick-create` - Quick create tag chip

### Solution Management Fragments
- `POST /admin/solutions/:id/stats` - Add stat
- `POST /admin/solutions/:id/challenges` - Add challenge
- `POST /admin/solutions/:id/products` - Link product
- `POST /admin/solutions/:id/ctas` - Add CTA

### Case Study Management Fragments
- `POST /admin/case-studies/:id/products` - Link product
- `POST /admin/case-studies/:id/metrics` - Add metric

### Media Library Fragments
- `GET /admin/media/browse` - Media picker modal

### Public Search Fragments
- `GET /search/suggest` - Search autocomplete suggestions
- `GET /products/search` - Product search results

### All Delete Operations
Most entity delete operations return no content (204) and are triggered via HTMX `hx-delete` to remove elements from the DOM.

---

## Rate-Limited Endpoints

Only one endpoint is rate-limited:

| Endpoint | Rate Limit | Middleware |
|----------|-----------|------------|
| `POST /contact/submit` | 5 requests per hour per IP | `contactLimiter.Middleware()` |

---

## API Response Types

### Full Pages
Return complete HTML documents with layout wrapper (`admin-layout.html` or public layouts).

### HTMX Fragments
Return HTML snippets (partials) without layout, intended to replace specific DOM elements.

### Form Submissions
Typically redirect with `http.StatusSeeOther (303)` after successful POST.

### JSON Responses
Media library API endpoints return JSON for programmatic access.

### XML/Text
SEO endpoints (`sitemap.xml`, `robots.txt`) return XML/text content.

---

## Notes

1. **Authentication**: All `/admin/*` routes except `/admin/login` require valid session authentication via `RequireAuth()` middleware.

2. **HTMX Pattern**: Admin panel heavily uses HTMX for interactive UI. Fragment endpoints return partial HTML that replaces specific page sections without full page reload.

3. **Caching**: Public pages use cache service with TTL (e.g., 600 seconds for product pages, 3600 for contact page).

4. **File Uploads**: Product images, downloads, whitepaper PDFs, and media library uploads are handled via multipart form data.

5. **Pagination**: Most list endpoints support `?page=N` query parameter with per-page limits defined as constants (e.g., `productsPerPage = 15`).

6. **Filtering**: List endpoints support various filters via query parameters (e.g., `?status=published&category=5&search=keyword`).

7. **Slug-based URLs**: Public content uses slug-based URLs for SEO (`/products/:category/:slug`, `/blog/:slug`, etc.).

8. **Activity Logging**: Most admin mutations log activity via `logActivity()` helper function.

9. **Cache Invalidation**: Admin handlers call `cache.DeleteByPrefix()` after mutations to ensure public pages reflect changes.

10. **Session Store**: Initialized in main with secret key: `change-this-secret-in-production-minimum-32-chars`.
