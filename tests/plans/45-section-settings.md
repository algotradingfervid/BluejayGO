# Test Plan: Section-Specific Settings

## Summary
Tests four separate section settings pages (About, Products, Solutions, Blog) each with GET+POST handlers and section-specific configuration options.

## Preconditions
- Server running on localhost:28090
- Admin user logged in (admin@bluejaylabs.com / password)
- Section settings initialized with default values
- Settings affect frontend display behavior for each section

## User Journey Steps
1. Navigate to /admin/about/settings
2. View About settings: checkboxes for show_mission, show_milestones, show_certifications, show_team
3. Toggle checkboxes and submit
4. Verify redirect to /admin/about/settings?saved=1
5. Navigate to /admin/products/settings
6. View Products settings: products_per_page (number), show_categories, show_search (checkboxes), default_sort (select)
7. Update values and submit
8. Verify redirect to /admin/products/settings?saved=1
9. Navigate to /admin/solutions/settings
10. View Solutions settings: solutions_per_page (number), show_industries, show_search (checkboxes)
11. Update values and submit
12. Verify redirect to /admin/solutions/settings?saved=1
13. Navigate to /admin/blog/settings
14. View Blog settings: blog_posts_per_page (number), show_author, show_date, show_categories, show_tags, show_search (checkboxes)
15. Update values and submit
16. Verify redirect to /admin/blog/settings?saved=1

## Test Cases

### Happy Path - About Settings
- **Load about settings**: Verifies GET /admin/about/settings loads with current values
- **Toggle show mission**: Checks/unchecks about_show_mission, verifies save
- **Toggle show milestones**: Checks/unchecks about_show_milestones, verifies save
- **Toggle show certifications**: Checks/unchecks about_show_certifications, verifies save
- **Toggle show team**: Checks/unchecks about_show_team, verifies save
- **All sections enabled**: Checks all checkboxes, verifies save
- **All sections disabled**: Unchecks all checkboxes, verifies save and frontend behavior

### Happy Path - Products Settings
- **Load products settings**: Verifies GET /admin/products/settings loads with current values
- **Update products per page**: Changes products_per_page to 12, verifies save
- **Toggle show categories**: Checks/unchecks products_show_categories, verifies save
- **Toggle show search**: Checks/unchecks products_show_search, verifies save
- **Update default sort**: Changes products_default_sort (e.g., "newest", "price_asc", "name"), verifies save
- **High per page value**: Sets products_per_page to 100, verifies save

### Happy Path - Solutions Settings
- **Load solutions settings**: Verifies GET /admin/solutions/settings loads with current values
- **Update solutions per page**: Changes solutions_per_page to 9, verifies save
- **Toggle show industries**: Checks/unchecks solutions_show_industries, verifies save
- **Toggle show search**: Checks/unchecks solutions_show_search, verifies save

### Happy Path - Blog Settings
- **Load blog settings**: Verifies GET /admin/blog/settings loads with current values
- **Update posts per page**: Changes blog_posts_per_page to 15, verifies save
- **Toggle show author**: Checks/unchecks blog_show_author, verifies save
- **Toggle show date**: Checks/unchecks blog_show_date, verifies save
- **Toggle show categories**: Checks/unchecks blog_show_categories, verifies save
- **Toggle show tags**: Checks/unchecks blog_show_tags, verifies save
- **Toggle show search**: Checks/unchecks blog_show_search, verifies save
- **All blog features enabled**: Checks all checkboxes, verifies save
- **Minimal blog display**: Unchecks author, date, categories, tags, verifies save

### Edge Cases / Error States
- **Zero per page**: Sets products/solutions/blog per page to 0, checks validation
- **Negative per page**: Sets per page to negative number, checks validation
- **Very large per page**: Sets per page to 9999, checks validation/limit
- **Invalid default sort**: Tests if invalid products_default_sort value is rejected
- **Empty default sort**: Tests if products_default_sort can be empty or requires selection
- **All about sections hidden**: Unchecks all about checkboxes, verifies about page behavior
- **Products without search or categories**: Disables both, verifies product listing behavior
- **Solutions without industries**: Disables solutions_show_industries, verifies display
- **Blog without metadata**: Disables all blog show checkboxes, verifies minimal post display
- **Checkbox to int64**: Verifies checkboxes correctly convert to 0/1 in database
- **Success redirect**: Confirms each page redirects to itself with ?saved=1
- **Success banner**: Verifies .alert-success appears on each page when ?saved=1

## Selectors & Elements

### About Settings
- Form: `form[action="/admin/about/settings"][method="POST"]`
- Show mission: `input[name="about_show_mission"][type="checkbox"]`
- Show milestones: `input[name="about_show_milestones"][type="checkbox"]`
- Show certifications: `input[name="about_show_certifications"][type="checkbox"]`
- Show team: `input[name="about_show_team"][type="checkbox"]`

### Products Settings
- Form: `form[action="/admin/products/settings"][method="POST"]`
- Per page: `input[name="products_per_page"][type="number"]`
- Show categories: `input[name="products_show_categories"][type="checkbox"]`
- Show search: `input[name="products_show_search"][type="checkbox"]`
- Default sort: `select[name="products_default_sort"]`
- Sort options: `option[value="newest"]`, `option[value="price_asc"]`, etc.

### Solutions Settings
- Form: `form[action="/admin/solutions/settings"][method="POST"]`
- Per page: `input[name="solutions_per_page"][type="number"]`
- Show industries: `input[name="solutions_show_industries"][type="checkbox"]`
- Show search: `input[name="solutions_show_search"][type="checkbox"]`

### Blog Settings
- Form: `form[action="/admin/blog/settings"][method="POST"]`
- Posts per page: `input[name="blog_posts_per_page"][type="number"]`
- Show author: `input[name="blog_show_author"][type="checkbox"]`
- Show date: `input[name="blog_show_date"][type="checkbox"]`
- Show categories: `input[name="blog_show_categories"][type="checkbox"]`
- Show tags: `input[name="blog_show_tags"][type="checkbox"]`
- Show search: `input[name="blog_show_search"][type="checkbox"]`

### Common Elements
- Submit button: `button[type="submit"]`
- Success banner: `.alert-success` (when ?saved=1)

## HTMX Interactions
- None - standard form POST with full page redirect for each settings page

## Dependencies
- Database settings table with section-specific configuration
- Templates: templates/admin/pages/about-settings.html, products-settings.html, solutions-settings.html, blog-settings.html
- Handlers: internal/handlers/about.go, products.go, solutions.go, blog.go (Get/Post settings methods)
- Checkbox values convert to int64 (0/1)
- Products default sort options defined (newest, oldest, price_asc, price_desc, name, etc.)
