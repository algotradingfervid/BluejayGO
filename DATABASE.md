# Bluejay CMS Database Documentation

## Overview

Bluejay CMS uses SQLite as its database engine with Write-Ahead Logging (WAL) mode enabled for optimal concurrency and performance. SQLite was chosen for this project because:

- **Zero Configuration**: No separate database server required, simplifying deployment
- **Portable**: Single file database that can be easily backed up and moved
- **Reliable**: ACID-compliant with strong consistency guarantees
- **Fast**: Optimized for read-heavy workloads typical of CMS applications
- **Pure Go Driver**: Uses modernc.org/sqlite (no CGO required)
- **Cost-Effective**: Perfect for small to medium-sized deployments without licensing costs

The database is configured with production-ready settings optimized for CMS use cases with multiple concurrent readers and a single writer.

## Connection Configuration

### Database Initialization

The database connection is initialized in `internal/database/sqlite.go` with the following configuration:

```go
db, err := InitDB(Config{Path: "./data/cms.db"})
```

### Pragmas Explained

The connection DSN includes several SQLite pragmas that configure database behavior:

```
path/to/cms.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on&_synchronous=NORMAL&_cache_size=2000
```

**`_journal_mode=WAL`** (Write-Ahead Logging)
- Enables WAL mode for better concurrency
- Allows readers and writers to operate simultaneously
- Writers don't block readers and vice versa
- Recommended for web applications with moderate write activity
- WAL file automatically checkpoints to main database file

**`_busy_timeout=5000`** (5 seconds)
- Sets timeout when database is locked
- SQLite will retry operations for up to 5 seconds before returning SQLITE_BUSY error
- Handles concurrent access gracefully without immediate failures
- Appropriate for typical web application load

**`_foreign_keys=on`**
- Enables foreign key constraint enforcement
- SQLite has foreign keys disabled by default for backward compatibility
- Essential for maintaining referential integrity across tables
- Prevents orphaned records when parent records are deleted

**`_synchronous=NORMAL`**
- Balances durability and performance
- FULL mode would be safer but significantly slower
- NORMAL is sufficient for most applications when combined with WAL mode
- Ensures data safety while allowing better write performance

**`_cache_size=2000`** (approximately 8MB)
- Sets page cache to 2000 pages
- With SQLite's default 4KB page size, provides ~8MB cache
- Improves query performance for frequently accessed data
- Reduces disk I/O for hot data

### Connection Pool Configuration

SQLite supports only one concurrent writer, so the connection pool is configured accordingly:

```go
db.SetMaxOpenConns(1)      // Only one connection can be open at a time
db.SetMaxIdleConns(1)      // Keep one connection alive in the pool
db.SetConnMaxLifetime(0)   // Connections never expire due to age
db.SetConnMaxIdleTime(0)   // Idle connections are never closed
```

This configuration:
- Prevents "database is locked" errors from connection pool contention
- Ensures serialized writes as required by SQLite
- Keeps connection alive to avoid reconnection overhead
- Optimizes for long-lived application instances

## Schema Reference

### System Tables

#### `settings`
Global site configuration stored as a singleton record (id=1).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Always 1 (singleton pattern) |
| site_name | TEXT | NOT NULL | Site name display |
| site_tagline | TEXT | NOT NULL | Site tagline/slogan |
| contact_email | TEXT | NOT NULL | Primary contact email |
| contact_phone | TEXT | NOT NULL | Contact phone number |
| address | TEXT | NOT NULL | Physical address |
| footer_text | TEXT | NOT NULL | Footer copyright text |
| meta_description | TEXT | NOT NULL | Default meta description |
| meta_keywords | TEXT | NOT NULL | Default meta keywords |
| google_analytics_id | TEXT | NOT NULL | GA tracking ID |
| social_linkedin | TEXT | NOT NULL | LinkedIn profile URL |
| social_twitter | TEXT | NOT NULL | Twitter/X profile URL |
| social_github | TEXT | NOT NULL | GitHub profile URL |
| social_facebook | TEXT | NOT NULL | Facebook profile URL |
| social_youtube | TEXT | NOT NULL | YouTube channel URL |
| social_instagram | TEXT | NOT NULL | Instagram profile URL |
| business_hours | TEXT | NOT NULL | Business hours text |
| about_text | TEXT | NOT NULL | About section text |
| header_logo_path | TEXT | NOT NULL | Header logo file path |
| header_logo_alt | TEXT | NOT NULL | Header logo alt text |
| header_cta_enabled | BOOLEAN | NOT NULL | Enable header CTA button |
| header_cta_text | TEXT | NOT NULL | Header CTA button text |
| header_cta_url | TEXT | NOT NULL | Header CTA button URL |
| header_cta_style | TEXT | NOT NULL | Header CTA button style |
| header_show_phone | BOOLEAN | NOT NULL | Show phone in header |
| header_show_email | BOOLEAN | NOT NULL | Show email in header |
| header_show_social | BOOLEAN | NOT NULL | Show social links in header |
| header_social_style | TEXT | NOT NULL | Social links display style |
| footer_columns | INTEGER | NOT NULL | Number of footer columns |
| footer_bg_style | TEXT | NOT NULL | Footer background style |
| footer_show_social | INTEGER | NOT NULL | Show social links in footer |
| footer_social_style | TEXT | NOT NULL | Footer social links style |
| footer_copyright | TEXT | NOT NULL | Footer copyright text template |
| show_nav_home | BOOLEAN | NOT NULL | Show Home in navigation |
| show_nav_about | BOOLEAN | NOT NULL | Show About in navigation |
| show_nav_products | BOOLEAN | NOT NULL | Show Products in navigation |
| show_nav_solutions | BOOLEAN | NOT NULL | Show Solutions in navigation |
| show_nav_blog | BOOLEAN | NOT NULL | Show Blog in navigation |
| show_nav_partners | BOOLEAN | NOT NULL | Show Partners in navigation |
| show_nav_contact | BOOLEAN | NOT NULL | Show Contact in navigation |
| show_nav_case_studies | BOOLEAN | NOT NULL | Show Case Studies in navigation |
| show_nav_whitepapers | BOOLEAN | NOT NULL | Show Whitepapers in navigation |
| nav_label_home | TEXT | NOT NULL | Home navigation label |
| nav_label_about | TEXT | NOT NULL | About navigation label |
| nav_label_products | TEXT | NOT NULL | Products navigation label |
| nav_label_solutions | TEXT | NOT NULL | Solutions navigation label |
| nav_label_blog | TEXT | NOT NULL | Blog navigation label |
| nav_label_partners | TEXT | NOT NULL | Partners navigation label |
| nav_label_contact | TEXT | NOT NULL | Contact navigation label |
| nav_label_case_studies | TEXT | NOT NULL | Case Studies navigation label |
| nav_label_whitepapers | TEXT | NOT NULL | Whitepapers navigation label |
| show_footer_about | BOOLEAN | NOT NULL | Show About in footer |
| show_footer_socials | BOOLEAN | NOT NULL | Show Social links in footer |
| show_footer_products | BOOLEAN | NOT NULL | Show Products in footer |
| show_footer_solutions | BOOLEAN | NOT NULL | Show Solutions in footer |
| show_footer_resources | BOOLEAN | NOT NULL | Show Resources in footer |
| show_footer_contact | BOOLEAN | NOT NULL | Show Contact in footer |
| footer_heading_products | TEXT | NOT NULL | Products footer heading |
| footer_heading_solutions | TEXT | NOT NULL | Solutions footer heading |
| footer_heading_resources | TEXT | NOT NULL | Resources footer heading |
| footer_heading_contact | TEXT | NOT NULL | Contact footer heading |
| homepage_show_* | INTEGER | NOT NULL | Homepage section visibility |
| homepage_max_* | INTEGER | NOT NULL | Homepage section item limits |
| homepage_hero_autoplay | INTEGER | NOT NULL | Hero carousel autoplay |
| homepage_hero_interval | INTEGER | NOT NULL | Hero carousel interval (seconds) |
| about_show_* | INTEGER | NOT NULL | About page section visibility |
| products_per_page | INTEGER | NOT NULL | Products per page |
| products_show_* | INTEGER | NOT NULL | Products page feature toggles |
| products_default_sort | TEXT | NOT NULL | Default products sort order |
| solutions_per_page | INTEGER | NOT NULL | Solutions per page |
| solutions_show_* | INTEGER | NOT NULL | Solutions page feature toggles |
| blog_posts_per_page | INTEGER | NOT NULL | Blog posts per page |
| blog_show_* | INTEGER | NOT NULL | Blog page feature toggles |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Indexes:**
- `idx_settings_singleton` (UNIQUE on id) - Enforces singleton pattern

#### `admin_users`
Administrative users with authentication and role management.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | User ID |
| email | TEXT | NOT NULL, UNIQUE | Login email address |
| password_hash | TEXT | NOT NULL | Bcrypt password hash |
| display_name | TEXT | NOT NULL | User's display name |
| role | TEXT | NOT NULL, DEFAULT 'editor', CHECK IN ('admin', 'editor') | User role |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Account active status |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Account creation date |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last profile update |
| last_login_at | DATETIME | NULL | Last successful login |

**Indexes:**
- `idx_admin_users_email` (UNIQUE) - Fast email lookups
- `idx_admin_users_role` - Filter by role
- `idx_admin_users_is_active` - Filter active users

**Relationships:**
- Referenced by `activity_log.user_id` (optional)

#### `activity_log`
Audit trail of admin actions for compliance and debugging.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Log entry ID |
| user_id | INTEGER | REFERENCES admin_users(id) | Admin user (NULL for system) |
| action | TEXT | NOT NULL | Action performed (create/update/delete/etc.) |
| resource_type | TEXT | NOT NULL | Type of resource affected |
| resource_id | INTEGER | NULL | ID of affected resource |
| resource_title | TEXT | NULL | Title/name of resource |
| description | TEXT | NOT NULL | Human-readable description |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | When action occurred |

**Indexes:**
- `idx_activity_log_created_at` (DESC) - Recent activity queries
- `idx_activity_log_action` - Filter by action type

---

### Product Tables

#### `product_categories`
Product taxonomy for organizing products.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Category ID |
| name | TEXT | NOT NULL, UNIQUE | Category display name |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| description | TEXT | NOT NULL | Category description |
| icon | TEXT | NOT NULL | Icon identifier (Material Icons) |
| image_url | TEXT | NULL | Category image path |
| product_count | INTEGER | NOT NULL, DEFAULT 0 | Cached product count |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_product_categories_slug` - Slug lookups
- `idx_product_categories_sort` - Ordered listings

**Relationships:**
- Referenced by `products.category_id` (ON DELETE RESTRICT)

#### `products`
Main products table with metadata and status.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Product ID |
| sku | TEXT | NOT NULL, UNIQUE | Stock keeping unit |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| name | TEXT | NOT NULL | Product name |
| tagline | TEXT | NULL | Short tagline |
| description | TEXT | NOT NULL | Full description (HTML) |
| overview | TEXT | NULL | Overview section (HTML) |
| category_id | INTEGER | NOT NULL, FK to product_categories | Product category |
| status | TEXT | NOT NULL, DEFAULT 'draft', CHECK IN ('draft', 'published', 'archived') | Publication status |
| is_featured | BOOLEAN | NOT NULL, DEFAULT 0 | Featured product flag |
| featured_order | INTEGER | NULL | Order among featured products |
| meta_title | TEXT | NULL | SEO meta title |
| meta_description | TEXT | NULL | SEO meta description |
| og_image | TEXT | NOT NULL, DEFAULT '' | Open Graph image |
| primary_image | TEXT | NULL | Primary product image |
| video_url | TEXT | NULL | Product video URL |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update (auto-updated) |
| published_at | DATETIME | NULL | Publication timestamp |

**Indexes:**
- `idx_products_category` - Filter by category
- `idx_products_slug` - Slug lookups
- `idx_products_status` - Filter by status
- `idx_products_featured` - Featured products ordering
- `idx_products_published` - Published date ordering

**Triggers:**
- `update_products_timestamp` - Auto-updates updated_at on modification

**Relationships:**
- References `product_categories(id)` (ON DELETE RESTRICT)
- Referenced by `product_specs`, `product_images`, `product_features`, `product_certifications`, `product_downloads` (CASCADE)
- Referenced by `solution_products`, `case_study_products`, `blog_post_products` (CASCADE)

#### `product_specs`
Technical specifications organized by section.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Spec ID |
| product_id | INTEGER | NOT NULL, FK to products | Product reference |
| section_name | TEXT | NOT NULL | Specification section (e.g., "Processor", "Memory") |
| spec_key | TEXT | NOT NULL | Specification name |
| spec_value | TEXT | NOT NULL | Specification value |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order within section |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_product_specs_product` - Product specs lookup

**Relationships:**
- References `products(id)` (ON DELETE CASCADE)

#### `product_images`
Product gallery images with ordering and metadata.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Image ID |
| product_id | INTEGER | NOT NULL, FK to products | Product reference |
| image_path | TEXT | NOT NULL | Image file path |
| alt_text | TEXT | NULL | Image alt text for accessibility |
| caption | TEXT | NULL | Image caption |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Gallery display order |
| is_thumbnail | BOOLEAN | NOT NULL, DEFAULT 0 | Thumbnail flag |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Upload timestamp |

**Indexes:**
- `idx_product_images_product` - Product images lookup

**Relationships:**
- References `products(id)` (ON DELETE CASCADE)

#### `product_features`
Bullet-point features list for products.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Feature ID |
| product_id | INTEGER | NOT NULL, FK to products | Product reference |
| feature_text | TEXT | NOT NULL | Feature description |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_product_features_product` - Product features lookup

**Relationships:**
- References `products(id)` (ON DELETE CASCADE)

#### `product_certifications`
Product certifications and compliance badges.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Certification ID |
| product_id | INTEGER | NOT NULL, FK to products | Product reference |
| certification_name | TEXT | NOT NULL | Certification full name |
| certification_code | TEXT | NULL | Certification code/number |
| icon_type | TEXT | NULL | Icon type identifier |
| icon_path | TEXT | NULL | Custom icon path |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_product_certifications_product` - Product certifications lookup

**Relationships:**
- References `products(id)` (ON DELETE CASCADE)

#### `product_downloads`
Downloadable resources (datasheets, manuals, drivers).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Download ID |
| product_id | INTEGER | NOT NULL, FK to products | Product reference |
| title | TEXT | NOT NULL | Download title |
| description | TEXT | NULL | Download description |
| file_type | TEXT | NOT NULL | File type (PDF, ZIP, etc.) |
| file_path | TEXT | NOT NULL | File storage path |
| file_size | BIGINT | NULL | File size in bytes |
| version | TEXT | NULL | Resource version |
| download_count | INTEGER | NOT NULL, DEFAULT 0 | Download counter |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Upload timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update (auto-updated) |

**Indexes:**
- `idx_product_downloads_product` - Product downloads lookup

**Triggers:**
- `update_product_downloads_timestamp` - Auto-updates updated_at

**Relationships:**
- References `products(id)` (ON DELETE CASCADE)

---

### Blog Tables

#### `blog_categories`
Blog post categorization.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Category ID |
| name | TEXT | NOT NULL, UNIQUE | Category name |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| color_hex | TEXT | NOT NULL | Category color (hex code) |
| description | TEXT | NULL | Category description |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_blog_categories_slug` - Slug lookups

**Relationships:**
- Referenced by `blog_posts.category_id` (ON DELETE RESTRICT)

#### `blog_authors`
Blog post authors/contributors.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Author ID |
| name | TEXT | NOT NULL | Author full name |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| title | TEXT | NOT NULL | Job title/role |
| bio | TEXT | NULL | Author biography |
| avatar_url | TEXT | NULL | Profile image path |
| linkedin_url | TEXT | NULL | LinkedIn profile URL |
| email | TEXT | NULL | Author email |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_blog_authors_slug` - Slug lookups

**Relationships:**
- Referenced by `blog_posts.author_id` (ON DELETE RESTRICT)

#### `blog_posts`
Blog posts with rich metadata and SEO fields.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Post ID |
| title | TEXT | NOT NULL | Post title |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| excerpt | TEXT | NOT NULL | Post excerpt/summary |
| body | TEXT | NOT NULL | Post content (HTML) |
| featured_image_url | TEXT | NULL | Featured image path |
| featured_image_alt | TEXT | NULL | Featured image alt text |
| category_id | INTEGER | NOT NULL, FK to blog_categories | Post category |
| author_id | INTEGER | NOT NULL, FK to blog_authors | Post author |
| meta_description | TEXT | NULL | SEO meta description |
| meta_title | TEXT | NOT NULL, DEFAULT '' | SEO meta title |
| og_image | TEXT | NOT NULL, DEFAULT '' | Open Graph image |
| reading_time_minutes | INTEGER | NULL | Estimated reading time |
| status | TEXT | NOT NULL, DEFAULT 'draft', CHECK IN ('draft', 'published', 'archived') | Publication status |
| published_at | DATETIME | NULL | Publication timestamp |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_blog_posts_slug` - Slug lookups
- `idx_blog_posts_category` - Filter by category
- `idx_blog_posts_author` - Filter by author
- `idx_blog_posts_status_published` - Published posts queries
- `idx_blog_posts_created` - Creation date ordering

**Relationships:**
- References `blog_categories(id)` (ON DELETE RESTRICT)
- References `blog_authors(id)` (ON DELETE RESTRICT)
- Referenced by `blog_post_tags`, `blog_post_products` (CASCADE)

#### `blog_tags`
Tags for blog post classification.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Tag ID |
| name | TEXT | NOT NULL, UNIQUE | Tag name |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_blog_tags_slug` - Slug lookups

**Relationships:**
- Referenced by `blog_post_tags.blog_tag_id` (CASCADE)

#### `blog_post_tags`
Many-to-many relationship between posts and tags.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| blog_post_id | INTEGER | NOT NULL, FK to blog_posts, PRIMARY KEY (composite) | Post reference |
| blog_tag_id | INTEGER | NOT NULL, FK to blog_tags, PRIMARY KEY (composite) | Tag reference |

**Indexes:**
- `idx_blog_post_tags_post` - Post tags lookup
- `idx_blog_post_tags_tag` - Tag posts lookup

**Relationships:**
- References `blog_posts(id)` (ON DELETE CASCADE)
- References `blog_tags(id)` (ON DELETE CASCADE)

#### `blog_post_products`
Many-to-many relationship between blog posts and products.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| blog_post_id | INTEGER | NOT NULL, FK to blog_posts, PRIMARY KEY (composite) | Post reference |
| product_id | INTEGER | NOT NULL, FK to products, PRIMARY KEY (composite) | Product reference |
| display_order | INTEGER | DEFAULT 0 | Display order |

**Relationships:**
- References `blog_posts(id)` (ON DELETE CASCADE)
- References `products(id)` (ON DELETE CASCADE)

---

### Solution Tables

#### `solutions`
Industry-specific solution pages.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Solution ID |
| title | TEXT | NOT NULL | Solution title |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| icon | TEXT | NOT NULL | Icon identifier |
| short_description | TEXT | NOT NULL | Brief description for listings |
| hero_image_url | TEXT | NULL | Hero section image |
| hero_title | TEXT | NULL | Hero section title |
| hero_description | TEXT | NULL | Hero section description |
| overview_content | TEXT | NULL | Overview section content (HTML) |
| meta_description | TEXT | NULL | SEO meta description |
| meta_title | TEXT | NOT NULL, DEFAULT '' | SEO meta title |
| og_image | TEXT | NOT NULL, DEFAULT '' | Open Graph image |
| reference_code | TEXT | NULL | Internal reference code |
| is_published | BOOLEAN | DEFAULT 0 | Publication status |
| display_order | INTEGER | DEFAULT 0 | Display order |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_solutions_slug` - Slug lookups
- `idx_solutions_published` - Published solutions queries

**Relationships:**
- Referenced by `solution_stats`, `solution_challenges`, `solution_products`, `solution_ctas` (CASCADE)

#### `solution_stats`
Industry statistics for solution pages.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Stat ID |
| solution_id | INTEGER | NOT NULL, FK to solutions | Solution reference |
| value | TEXT | NOT NULL | Statistic value (e.g., "85%") |
| label | TEXT | NOT NULL | Statistic label |
| display_order | INTEGER | DEFAULT 0 | Display order |

**Indexes:**
- `idx_solution_stats_solution` - Solution stats lookup

**Relationships:**
- References `solutions(id)` (ON DELETE CASCADE)

#### `solution_challenges`
Industry challenges addressed by solution.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Challenge ID |
| solution_id | INTEGER | NOT NULL, FK to solutions | Solution reference |
| title | TEXT | NOT NULL | Challenge title |
| description | TEXT | NOT NULL | Challenge description |
| icon | TEXT | NOT NULL | Icon identifier |
| display_order | INTEGER | DEFAULT 0 | Display order |

**Indexes:**
- `idx_solution_challenges_solution` - Solution challenges lookup

**Relationships:**
- References `solutions(id)` (ON DELETE CASCADE)

#### `solution_products`
Many-to-many relationship between solutions and products.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Relationship ID |
| solution_id | INTEGER | NOT NULL, FK to solutions, UNIQUE (composite) | Solution reference |
| product_id | INTEGER | NOT NULL, FK to products, UNIQUE (composite) | Product reference |
| display_order | INTEGER | DEFAULT 0 | Display order |
| is_featured | BOOLEAN | DEFAULT 0 | Featured product flag |

**Indexes:**
- `idx_solution_products_solution` - Solution products lookup
- `idx_solution_products_product` - Product solutions lookup

**Relationships:**
- References `solutions(id)` (ON DELETE CASCADE)
- References `products(id)` (ON DELETE CASCADE)

#### `solution_ctas`
Call-to-action sections for solution pages.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | CTA ID |
| solution_id | INTEGER | NOT NULL, FK to solutions | Solution reference |
| heading | TEXT | NOT NULL | CTA heading |
| subheading | TEXT | NULL | CTA subheading |
| primary_button_text | TEXT | NULL | Primary button label |
| primary_button_url | TEXT | NULL | Primary button URL |
| secondary_button_text | TEXT | NULL | Secondary button label |
| secondary_button_url | TEXT | NULL | Secondary button URL |
| phone_number | TEXT | NULL | Contact phone number |
| section_name | TEXT | NOT NULL | Section identifier |

**Indexes:**
- `idx_solution_ctas_solution` - Solution CTAs lookup

**Relationships:**
- References `solutions(id)` (ON DELETE CASCADE)

#### `solution_page_features`
"Why Choose BlueJay" features on solutions listing page.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Feature ID |
| title | TEXT | NOT NULL | Feature title |
| description | TEXT | NOT NULL | Feature description |
| icon | TEXT | NOT NULL | Icon identifier |
| display_order | INTEGER | DEFAULT 0 | Display order |
| is_active | BOOLEAN | DEFAULT 1 | Active status |

#### `solutions_listing_cta`
CTA section for solutions listing page.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | CTA ID |
| heading | TEXT | NOT NULL | CTA heading |
| subheading | TEXT | NULL | CTA subheading |
| primary_button_text | TEXT | NULL | Primary button label |
| primary_button_url | TEXT | NULL | Primary button URL |
| secondary_button_text | TEXT | NULL | Secondary button label |
| secondary_button_url | TEXT | NULL | Secondary button URL |
| is_active | BOOLEAN | DEFAULT 1 | Active status |

---

### Case Study Tables

#### `industries`
Industry taxonomy for case studies.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Industry ID |
| name | TEXT | NOT NULL, UNIQUE | Industry name |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| icon | TEXT | NOT NULL | Icon identifier |
| description | TEXT | NOT NULL | Industry description |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_industries_slug` - Slug lookups

**Relationships:**
- Referenced by `case_studies.industry_id` (ON DELETE RESTRICT)

#### `case_studies`
Customer success stories and case studies.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Case study ID |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| title | TEXT | NOT NULL | Case study title |
| client_name | TEXT | NOT NULL | Client/customer name |
| industry_id | INTEGER | NOT NULL, FK to industries | Industry reference |
| hero_image_url | TEXT | NULL | Hero section image |
| summary | TEXT | NOT NULL | Executive summary |
| challenge_title | TEXT | NOT NULL, DEFAULT 'The Challenge' | Challenge section title |
| challenge_content | TEXT | NOT NULL | Challenge description (HTML) |
| challenge_bullets | TEXT | NULL | Challenge bullet points (JSON array) |
| solution_title | TEXT | NOT NULL, DEFAULT 'Our Solution' | Solution section title |
| solution_content | TEXT | NOT NULL | Solution description (HTML) |
| outcome_title | TEXT | NOT NULL, DEFAULT 'The Outcome' | Outcome section title |
| outcome_content | TEXT | NOT NULL | Outcome description (HTML) |
| meta_title | TEXT | NULL | SEO meta title |
| meta_description | TEXT | NULL | SEO meta description |
| og_image | TEXT | NOT NULL, DEFAULT '' | Open Graph image |
| is_published | INTEGER | NOT NULL, DEFAULT 0 | Publication status |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_case_studies_slug` - Slug lookups
- `idx_case_studies_industry` - Filter by industry
- `idx_case_studies_published` - Published case studies
- `idx_case_studies_display_order` - Ordered listings

**Relationships:**
- References `industries(id)` (ON DELETE RESTRICT)
- Referenced by `case_study_products`, `case_study_metrics` (CASCADE)

#### `case_study_products`
Many-to-many relationship between case studies and products.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Relationship ID |
| case_study_id | INTEGER | NOT NULL, FK to case_studies, UNIQUE (composite) | Case study reference |
| product_id | INTEGER | NOT NULL, FK to products, UNIQUE (composite) | Product reference |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_case_study_products_case_study` - Case study products lookup
- `idx_case_study_products_product` - Product case studies lookup

**Relationships:**
- References `case_studies(id)` (ON DELETE CASCADE)
- References `products(id)` (ON DELETE CASCADE)

#### `case_study_metrics`
Success metrics and results for case studies.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Metric ID |
| case_study_id | INTEGER | NOT NULL, FK to case_studies | Case study reference |
| metric_value | TEXT | NOT NULL | Metric value (e.g., "40% increase") |
| metric_label | TEXT | NOT NULL | Metric label/description |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_case_study_metrics_case_study` - Case study metrics lookup

**Relationships:**
- References `case_studies(id)` (ON DELETE CASCADE)

---

### Whitepaper Tables

#### `whitepaper_topics`
Topic taxonomy for whitepapers.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Topic ID |
| name | TEXT | NOT NULL, UNIQUE | Topic name |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| color_hex | TEXT | NOT NULL | Topic color (hex code) |
| icon | TEXT | NOT NULL | Icon identifier |
| description | TEXT | NULL | Topic description |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_whitepaper_topics_slug` - Slug lookups

**Relationships:**
- Referenced by `whitepapers.topic_id` (ON DELETE RESTRICT)

#### `whitepapers`
Downloadable whitepapers and reports.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Whitepaper ID |
| title | TEXT | NOT NULL | Whitepaper title |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| description | TEXT | NOT NULL | Whitepaper description |
| topic_id | INTEGER | NOT NULL, FK to whitepaper_topics | Topic reference |
| pdf_file_path | TEXT | NOT NULL | PDF file storage path |
| file_size_bytes | INTEGER | NOT NULL | File size in bytes |
| page_count | INTEGER | NULL | Number of pages |
| published_date | TEXT | NOT NULL | Publication date (ISO format) |
| is_published | INTEGER | NOT NULL, DEFAULT 0 | Publication status |
| cover_color_from | TEXT | NOT NULL, DEFAULT '#0066CC' | Gradient start color |
| cover_color_to | TEXT | NOT NULL, DEFAULT '#004499' | Gradient end color |
| meta_description | TEXT | NULL | SEO meta description |
| meta_title | TEXT | NOT NULL, DEFAULT '' | SEO meta title |
| og_image | TEXT | NOT NULL, DEFAULT '' | Open Graph image |
| download_count | INTEGER | NOT NULL, DEFAULT 0 | Download counter |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_whitepapers_slug` - Slug lookups
- `idx_whitepapers_topic_id` - Filter by topic
- `idx_whitepapers_published` - Published whitepapers

**Relationships:**
- References `whitepaper_topics(id)` (ON DELETE RESTRICT)
- Referenced by `whitepaper_learning_points`, `whitepaper_downloads` (CASCADE)

#### `whitepaper_learning_points`
Key learning points for whitepapers.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Learning point ID |
| whitepaper_id | INTEGER | NOT NULL, FK to whitepapers | Whitepaper reference |
| point_text | TEXT | NOT NULL | Learning point text |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_whitepaper_learning_points_whitepaper` - Whitepaper learning points lookup

**Relationships:**
- References `whitepapers(id)` (ON DELETE CASCADE)

#### `whitepaper_downloads`
Lead generation tracking for whitepaper downloads.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Download ID |
| whitepaper_id | INTEGER | NOT NULL, FK to whitepapers | Whitepaper reference |
| name | TEXT | NOT NULL | Downloader name |
| email | TEXT | NOT NULL | Downloader email |
| company | TEXT | NOT NULL | Company name |
| designation | TEXT | NULL | Job title |
| marketing_consent | INTEGER | NOT NULL, DEFAULT 0 | Marketing opt-in flag |
| ip_address | TEXT | NULL | IP address (for security) |
| user_agent | TEXT | NULL | Browser user agent |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Download timestamp |

**Indexes:**
- `idx_whitepaper_downloads_whitepaper` - Whitepaper downloads lookup
- `idx_whitepaper_downloads_email` - Email-based queries
- `idx_whitepaper_downloads_created` - Time-based analytics

**Relationships:**
- References `whitepapers(id)` (ON DELETE CASCADE)

---

### Partner Tables

#### `partner_tiers`
Partner tier/level taxonomy.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Tier ID |
| name | TEXT | NOT NULL, UNIQUE | Tier name (e.g., "Gold", "Silver") |
| slug | TEXT | NOT NULL, UNIQUE | URL-friendly identifier |
| description | TEXT | NOT NULL | Tier description |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_partner_tiers_slug` - Slug lookups

**Relationships:**
- Referenced by `partners.tier_id` (ON DELETE RESTRICT)

#### `partners`
Partner organizations and integrations.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Partner ID |
| name | TEXT | NOT NULL, UNIQUE | Partner name |
| tier_id | INTEGER | NOT NULL, FK to partner_tiers | Partner tier |
| logo_url | TEXT | NULL | Partner logo path |
| icon | TEXT | NULL | Icon identifier |
| website_url | TEXT | NULL | Partner website URL |
| description | TEXT | NULL | Partner description |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Active status |
| is_featured | INTEGER | NOT NULL, DEFAULT 0 | Featured on homepage |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_partners_tier` - Filter by tier
- `idx_partners_order` - Ordered listings
- `idx_partners_featured` - Featured partners

**Relationships:**
- References `partner_tiers(id)` (ON DELETE RESTRICT)
- Referenced by `partner_testimonials` (CASCADE)

#### `partner_testimonials`
Testimonials from partner organizations.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Testimonial ID |
| partner_id | INTEGER | NOT NULL, FK to partners | Partner reference |
| quote | TEXT | NOT NULL | Testimonial quote |
| author_name | TEXT | NOT NULL | Testimonial author name |
| author_title | TEXT | NOT NULL | Author job title |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Active status |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Indexes:**
- `idx_testimonials_active` - Active testimonials by display order

**Relationships:**
- References `partners(id)` (ON DELETE CASCADE)

---

### Website Content Tables

#### `homepage_hero`
Homepage hero carousel slides.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Hero ID |
| headline | TEXT | NOT NULL | Hero headline |
| subheadline | TEXT | NOT NULL | Hero subheadline |
| badge_text | TEXT | NULL | Badge/label text |
| primary_cta_text | TEXT | NOT NULL, DEFAULT 'Explore Products' | Primary button text |
| primary_cta_url | TEXT | NOT NULL, DEFAULT '/products' | Primary button URL |
| secondary_cta_text | TEXT | NULL | Secondary button text |
| secondary_cta_url | TEXT | NULL | Secondary button URL |
| background_image | TEXT | NULL | Background image path |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Active status |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Carousel slide order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

#### `homepage_stats`
Homepage statistics section.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Stat ID |
| stat_value | TEXT | NOT NULL | Statistic value (e.g., "500+") |
| stat_label | TEXT | NOT NULL | Statistic label |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Active status |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### `homepage_testimonials`
Homepage testimonials section.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Testimonial ID |
| quote | TEXT | NOT NULL | Testimonial quote |
| author_name | TEXT | NOT NULL | Author name |
| author_title | TEXT | NULL | Author job title |
| author_company | TEXT | NULL | Author company |
| author_image | TEXT | NULL | Author photo path |
| rating | INTEGER | NOT NULL, DEFAULT 5 | Star rating (1-5) |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Active status |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### `homepage_cta`
Homepage call-to-action section.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | CTA ID |
| headline | TEXT | NOT NULL | CTA headline |
| description | TEXT | NULL | CTA description |
| primary_cta_text | TEXT | NOT NULL, DEFAULT 'Schedule a Demo' | Primary button text |
| primary_cta_url | TEXT | NOT NULL, DEFAULT '/contact' | Primary button URL |
| secondary_cta_text | TEXT | NULL | Secondary button text |
| secondary_cta_url | TEXT | NULL | Secondary button URL |
| background_style | TEXT | DEFAULT 'primary' | Background style |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Active status |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

#### `company_overview`
About page company overview content.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Overview ID |
| headline | TEXT | NOT NULL | Main headline |
| tagline | TEXT | NOT NULL | Company tagline |
| description_main | TEXT | NOT NULL | Main description paragraph |
| description_secondary | TEXT | NULL | Secondary description |
| description_tertiary | TEXT | NULL | Tertiary description |
| hero_image_url | TEXT | NULL | Hero section image |
| company_image_url | TEXT | NULL | Company/team image |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

#### `mission_vision_values`
About page mission, vision, and values.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | MVV ID |
| mission | TEXT | NOT NULL | Mission statement |
| vision | TEXT | NOT NULL | Vision statement |
| values_summary | TEXT | NULL | Values summary |
| mission_icon | TEXT | NOT NULL, DEFAULT 'flag' | Mission icon |
| vision_icon | TEXT | NOT NULL, DEFAULT 'visibility' | Vision icon |
| values_icon | TEXT | NOT NULL, DEFAULT 'diamond' | Values icon |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

#### `core_values`
About page core values list.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Value ID |
| title | TEXT | NOT NULL | Value title |
| description | TEXT | NOT NULL | Value description |
| icon | TEXT | NOT NULL | Icon identifier |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### `milestones`
About page company timeline/milestones.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Milestone ID |
| year | INTEGER | NOT NULL | Milestone year |
| title | TEXT | NOT NULL | Milestone title |
| description | TEXT | NOT NULL | Milestone description |
| is_current | INTEGER | NOT NULL, DEFAULT 0 | Current milestone flag |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### `certifications`
Company certifications (About page).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Certification ID |
| name | TEXT | NOT NULL | Certification name |
| abbreviation | TEXT | NOT NULL | Certification abbreviation |
| description | TEXT | NULL | Certification description |
| icon | TEXT | NULL | Icon identifier |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

#### `contact_submissions`
Contact form submissions and inquiries.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Submission ID |
| name | TEXT | NOT NULL | Submitter name |
| email | TEXT | NOT NULL | Submitter email |
| phone | TEXT | NOT NULL | Contact phone |
| company | TEXT | NOT NULL | Company name |
| inquiry_type | TEXT | NULL | Type of inquiry |
| message | TEXT | NOT NULL | Inquiry message |
| ip_address | TEXT | NULL | IP address (security) |
| user_agent | TEXT | NULL | Browser user agent |
| status | TEXT | NOT NULL, DEFAULT 'new' | Submission status |
| notes | TEXT | NULL | Internal notes |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Submission timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_contact_submissions_status` - Filter by status
- `idx_contact_submissions_created` - Time-based queries
- `idx_contact_submissions_email` - Email-based searches

#### `office_locations`
Company office locations (Contact page).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Location ID |
| name | TEXT | NOT NULL | Office name |
| address_line1 | TEXT | NOT NULL | Address line 1 |
| address_line2 | TEXT | NULL | Address line 2 |
| city | TEXT | NOT NULL | City |
| state | TEXT | NOT NULL | State/province |
| postal_code | TEXT | NOT NULL | Postal/ZIP code |
| country | TEXT | NOT NULL, DEFAULT 'India' | Country |
| phone | TEXT | NULL | Office phone |
| email | TEXT | NULL | Office email |
| is_primary | INTEGER | NOT NULL, DEFAULT 0 | Primary office flag |
| is_active | INTEGER | NOT NULL, DEFAULT 1 | Active status |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_office_locations_active` - Active locations

#### `page_sections`
Reusable page section content (CTAs, headings, etc.).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Section ID |
| page_key | TEXT | NOT NULL, UNIQUE (composite) | Page identifier |
| section_key | TEXT | NOT NULL, UNIQUE (composite) | Section identifier |
| heading | TEXT | NOT NULL, DEFAULT '' | Section heading |
| subheading | TEXT | NOT NULL, DEFAULT '' | Section subheading |
| description | TEXT | NOT NULL, DEFAULT '' | Section description |
| label | TEXT | NOT NULL, DEFAULT '' | Section label |
| primary_button_text | TEXT | NOT NULL, DEFAULT '' | Primary button text |
| primary_button_url | TEXT | NOT NULL, DEFAULT '' | Primary button URL |
| secondary_button_text | TEXT | NOT NULL, DEFAULT '' | Secondary button text |
| secondary_button_url | TEXT | NOT NULL, DEFAULT '' | Secondary button URL |
| is_active | BOOLEAN | NOT NULL, DEFAULT 1 | Active status |
| display_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |
| created_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update |

**Indexes:**
- `idx_page_sections_key` (UNIQUE on page_key, section_key)

---

### Media and Navigation Tables

#### `media_files`
Centralized media library for uploads.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Media ID |
| filename | TEXT | NOT NULL | Stored filename |
| original_filename | TEXT | NOT NULL | Original uploaded filename |
| file_path | TEXT | NOT NULL | File storage path |
| file_size | INTEGER | NOT NULL | File size in bytes |
| mime_type | TEXT | NOT NULL | MIME type |
| width | INTEGER | NULL | Image width (if image) |
| height | INTEGER | NULL | Image height (if image) |
| alt_text | TEXT | DEFAULT '' | Default alt text |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | Upload timestamp |

**Indexes:**
- `idx_media_files_filename` - Filename lookups
- `idx_media_files_mime_type` - Filter by type
- `idx_media_files_created_at` - Recent uploads

#### `navigation_menus`
Navigation menu definitions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Menu ID |
| name | TEXT | NOT NULL | Menu name |
| location | TEXT | NOT NULL | Menu location (header, footer, etc.) |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updated_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | Last update |

**Relationships:**
- Referenced by `navigation_items.menu_id` (CASCADE)

#### `navigation_items`
Navigation menu items with hierarchy support.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Item ID |
| menu_id | INTEGER | NOT NULL, FK to navigation_menus | Menu reference |
| parent_id | INTEGER | NULL, FK to navigation_items | Parent item (for submenus) |
| label | TEXT | NOT NULL | Menu item label |
| link_type | TEXT | NOT NULL, DEFAULT 'page' | Link type (page, custom, etc.) |
| url | TEXT | NULL | Custom URL |
| page_identifier | TEXT | NULL | Internal page identifier |
| open_new_tab | INTEGER | DEFAULT 0 | Open in new tab flag |
| is_active | INTEGER | DEFAULT 1 | Active status |
| sort_order | INTEGER | DEFAULT 0 | Display order |
| created_at | DATETIME | DEFAULT CURRENT_TIMESTAMP | Creation timestamp |

**Relationships:**
- References `navigation_menus(id)` (ON DELETE CASCADE)
- Self-references `parent_id` for hierarchy (ON DELETE CASCADE)

#### `footer_column_items`
Footer column structure and content.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Column ID |
| column_index | INTEGER | NOT NULL | Column number (1-4) |
| type | TEXT | NOT NULL, DEFAULT 'links' | Column type |
| heading | TEXT | NOT NULL, DEFAULT '' | Column heading |
| content | TEXT | NOT NULL, DEFAULT '' | Column content |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |

**Relationships:**
- Referenced by `footer_links.column_item_id` (CASCADE)

#### `footer_links`
Links within footer columns.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Link ID |
| column_item_id | INTEGER | NOT NULL, FK to footer_column_items | Column reference |
| label | TEXT | NOT NULL, DEFAULT '' | Link label |
| url | TEXT | NOT NULL, DEFAULT '' | Link URL |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |

**Relationships:**
- References `footer_column_items(id)` (ON DELETE CASCADE)

#### `footer_legal_links`
Legal links in footer bottom bar.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Link ID |
| label | TEXT | NOT NULL, DEFAULT '' | Link label (e.g., "Privacy Policy") |
| url | TEXT | NOT NULL, DEFAULT '' | Link URL |
| sort_order | INTEGER | NOT NULL, DEFAULT 0 | Display order |

---

### FTS5 Search Index

SQLite's FTS5 (Full-Text Search version 5) provides fast text search across content tables.

#### `products_fts`
Virtual table for product full-text search.

**Indexed Columns:**
- `name` - Product name
- `tagline` - Product tagline
- `description` - Product description

**Content Source:** `products` table (content='products', content_rowid='id')

**Triggers:** Automatically synced with `products` table via INSERT/UPDATE/DELETE triggers

#### `blog_posts_fts`
Virtual table for blog post full-text search.

**Indexed Columns:**
- `title` - Post title
- `excerpt` - Post excerpt
- `body` - Post body content

**Content Source:** `blog_posts` table (content='blog_posts', content_rowid='id')

**Triggers:** Automatically synced with `blog_posts` table

#### `case_studies_fts`
Virtual table for case study full-text search.

**Indexed Columns:**
- `title` - Case study title
- `client_name` - Client name
- `challenge_content` - Challenge description
- `solution_content` - Solution description

**Content Source:** `case_studies` table (content='case_studies', content_rowid='id')

**Triggers:** Automatically synced with `case_studies` table

**Usage Example:**
```sql
-- Search products
SELECT * FROM products WHERE id IN (
    SELECT rowid FROM products_fts WHERE products_fts MATCH 'desktop OR computer'
);

-- Search with ranking
SELECT p.*, rank FROM products p
JOIN (SELECT rowid, rank FROM products_fts WHERE products_fts MATCH 'touchscreen') fts
ON p.id = fts.rowid
ORDER BY rank;
```

## Entity Relationship Diagram

```
┌──────────────────┐
│    settings      │
│  (singleton=1)   │
└──────────────────┘

┌──────────────────┐         ┌──────────────────┐
│  admin_users     │────┐    │  activity_log    │
└──────────────────┘    └───>└──────────────────┘

┌──────────────────┐
│ product_categories│
└────────┬─────────┘
         │
         │ ON DELETE RESTRICT
         v
┌──────────────────┐         ┌───────────────────┐
│    products      │<───────>│  solution_products│
└────────┬─────────┘         └───────────────────┘
         │                            ^
         │ CASCADE                    │
         v                            │
    ┌────────┴──────────┬─────────────┼──────────┬─────────────┐
    v                   v             v          v             v
┌─────────────┐  ┌─────────────┐  ┌──────────┐ ┌─────────────┐ ┌──────────────┐
│product_specs│  │product_images│  │solutions │ │case_studies │ │ blog_posts   │
└─────────────┘  └─────────────┘  └────┬─────┘ └──────┬──────┘ └──────┬───────┘
                                        │              │               │
┌──────────────┐  ┌──────────────────┐ │              │               │
│product_      │  │product_          │ │              │               │
│features      │  │certifications    │ │              │               │
└──────────────┘  └──────────────────┘ │              │               │
                                        v              v               v
┌──────────────┐  ┌──────────────────┐ │       ┌──────────────┐ ┌──────────────┐
│product_      │  │product_downloads │ │       │case_study_   │ │blog_post_    │
│downloads     │  └──────────────────┘ │       │products      │ │products      │
└──────────────┘                        │       └──────────────┘ └──────────────┘
                                        │       ┌──────────────┐ ┌──────────────┐
                                        │       │case_study_   │ │blog_post_tags│
                                        │       │metrics       │ └──────┬───────┘
                                        │       └──────────────┘        │
                                        │                               v
                 ┌──────────────────────┼────────┬──────────┐   ┌─────────────┐
                 v                      v        v          v   │ blog_tags   │
         ┌───────────────┐     ┌─────────────┐ ┌────────────┐ └─────────────┘
         │solution_stats │     │solution_    │ │solution_   │
         └───────────────┘     │challenges   │ │ctas        │
                               └─────────────┘ └────────────┘

┌──────────────────┐         ┌──────────────────┐
│ blog_categories  │<────────│   blog_posts     │
└──────────────────┘         └───────┬──────────┘
                                     ^
┌──────────────────┐                │
│  blog_authors    │<───────────────┘
└──────────────────┘

┌──────────────────┐         ┌──────────────────┐
│   industries     │<────────│  case_studies    │
└──────────────────┘         └──────────────────┘

┌──────────────────┐         ┌──────────────────┐
│whitepaper_topics │<────────│   whitepapers    │
└──────────────────┘         └────────┬─────────┘
                                      │
                             ┌────────┴─────────┬───────────────────┐
                             v                  v                   v
                    ┌────────────────┐  ┌──────────────┐   ┌──────────────┐
                    │whitepaper_     │  │whitepaper_   │   │              │
                    │learning_points │  │downloads     │   │              │
                    └────────────────┘  └──────────────┘   │              │

┌──────────────────┐         ┌──────────────────┐
│  partner_tiers   │<────────│    partners      │
└──────────────────┘         └────────┬─────────┘
                                      │
                                      v
                             ┌────────────────────┐
                             │partner_testimonials│
                             └────────────────────┘

Homepage Tables:
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│ homepage_hero    │  │ homepage_stats   │  │homepage_         │  │ homepage_cta     │
└──────────────────┘  └──────────────────┘  │testimonials      │  └──────────────────┘
                                             └──────────────────┘

About Page Tables:
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│company_overview  │  │mission_vision_   │  │  core_values     │  │   milestones     │
└──────────────────┘  │values            │  └──────────────────┘  └──────────────────┘
                      └──────────────────┘
┌──────────────────┐
│ certifications   │
└──────────────────┘

Contact Tables:
┌──────────────────┐  ┌──────────────────┐
│contact_          │  │office_locations  │
│submissions       │  └──────────────────┘
└──────────────────┘

Navigation & Media:
┌──────────────────┐         ┌──────────────────┐
│navigation_menus  │<────────│navigation_items  │
└──────────────────┘         └──────────────────┘

┌──────────────────┐
│  media_files     │
└──────────────────┘

┌──────────────────┐         ┌──────────────────┐
│footer_column_    │<────────│  footer_links    │
│items             │         └──────────────────┘
└──────────────────┘
┌──────────────────┐
│footer_legal_links│
└──────────────────┘

┌──────────────────┐
│  page_sections   │
└──────────────────┘

Full-Text Search (Virtual Tables):
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│  products_fts    │  │ blog_posts_fts   │  │case_studies_fts  │
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

## Migrations Strategy

### Migration Tool

Bluejay CMS uses **golang-migrate** for database migrations. Migrations are SQL files stored in `db/migrations/`.

### Naming Convention

Migrations follow a strict naming pattern:

```
{version}_{description}.{direction}.sql
```

**Examples:**
- `001_settings.up.sql` - Initial settings table
- `001_settings.down.sql` - Rollback settings table
- `015_create_blog_posts.up.sql` - Create blog posts table
- `015_create_blog_posts.down.sql` - Drop blog posts table

**Rules:**
- Version numbers are sequential (001, 002, 003...)
- Use zero-padded 3-digit numbers for proper sorting
- Descriptions use snake_case
- Every `.up.sql` must have a corresponding `.down.sql`
- Down migrations should cleanly reverse up migrations

### Migration File Structure

**Up Migration** (applies changes):
```sql
-- Create table
CREATE TABLE IF NOT EXISTS example (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- Create indexes
CREATE INDEX idx_example_name ON example(name);

-- Seed initial data (if needed)
INSERT INTO example (name) VALUES ('default');
```

**Down Migration** (reverts changes):
```sql
-- Drop indexes first
DROP INDEX IF EXISTS idx_example_name;

-- Drop table
DROP TABLE IF EXISTS example;
```

### Creating New Migrations

1. Determine the next version number by checking existing migrations
2. Create both `.up.sql` and `.down.sql` files
3. Write SQL for up migration (CREATE, ALTER, INSERT, etc.)
4. Write SQL for down migration (DROP, reverse ALTERs, DELETE, etc.)
5. Test both directions locally

**Example workflow:**
```bash
# Check latest migration
ls db/migrations/ | tail -1
# Shows: 034_section_settings.up.sql

# Create new migration files
touch db/migrations/035_add_feature_flags.up.sql
touch db/migrations/035_add_feature_flags.down.sql

# Edit the files with your SQL
vim db/migrations/035_add_feature_flags.up.sql
vim db/migrations/035_add_feature_flags.down.sql

# Test migration up
migrate -path db/migrations -database "sqlite3://data/cms.db" up 1

# Test migration down
migrate -path db/migrations -database "sqlite3://data/cms.db" down 1
```

### Migration Best Practices

1. **Idempotent Migrations**: Use `IF NOT EXISTS` and `IF EXISTS` clauses
2. **Data Safety**: Never destructively modify production data without backup
3. **Foreign Keys**: Always specify ON DELETE behavior explicitly
4. **Indexes**: Create indexes for foreign keys and frequently queried columns
5. **Defaults**: Provide sensible defaults for new columns
6. **Triggers**: Drop triggers in down migrations before dropping tables
7. **Order Matters**: Drop dependent objects before parent objects

### Common Migration Patterns

**Adding a Column:**
```sql
-- up
ALTER TABLE products ADD COLUMN new_field TEXT NOT NULL DEFAULT '';

-- down
ALTER TABLE products DROP COLUMN new_field;
```

**Adding a Table with Foreign Key:**
```sql
-- up
CREATE TABLE child (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER NOT NULL REFERENCES parent(id) ON DELETE CASCADE
);

-- down
DROP TABLE IF EXISTS child;
```

**Creating an Index:**
```sql
-- up
CREATE INDEX idx_products_status ON products(status);

-- down
DROP INDEX IF EXISTS idx_products_status;
```

## sqlc Usage

### Configuration

sqlc is configured in `sqlc.yaml`:

```yaml
version: "2"
sql:
  - engine: "sqlite"
    queries: "db/queries"      # SQL query files location
    schema: "db/migrations"    # Schema source (migration files)
    gen:
      go:
        package: "sqlc"                    # Generated Go package name
        out: "db/sqlc"                     # Output directory
        emit_json_tags: true               # Add json tags to structs
        emit_prepared_queries: false       # Don't use prepared statements
        emit_interface: true               # Generate Queries interface
        emit_exact_table_names: false      # Use singular names (Product not products)
        emit_empty_slices: true            # Return empty slices, not nil
```

### Query File Format

Query files in `db/queries/` use sqlc's annotation format:

```sql
-- name: FunctionName :return_type
SELECT ... FROM ... WHERE ...;
```

**Return Type Annotations:**

- `:one` - Returns single row, error if not found or multiple found
- `:many` - Returns slice of rows (can be empty)
- `:exec` - Executes statement, returns error only
- `:execresult` - Returns sql.Result (for affected rows count)
- `:execrows` - Returns number of affected rows

### Generated Code Structure

sqlc generates:
- `db/sqlc/models.go` - Go structs for each table
- `db/sqlc/query.go` - Interface definition
- `db/sqlc/[table].sql.go` - Query functions per query file

**Example generated function:**
```go
// From: -- name: GetProduct :one
func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
    // ... generated code
}

// From: -- name: ListProducts :many
func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
    // ... generated code
}

// From: -- name: DeleteProduct :exec
func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
    // ... generated code
}
```

### Parameter Binding

**Positional Parameters** (simple queries):
```sql
-- name: GetProductByID :one
SELECT * FROM products WHERE id = ?;
```

**Named Parameters** (complex queries):
```sql
-- name: ListProductsAdminFiltered :many
SELECT p.* FROM products p
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE p.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE p.category_id = @filter_category END)
LIMIT @page_limit OFFSET @page_offset;
```

Generates:
```go
type ListProductsAdminFilteredParams struct {
    FilterStatus   string
    FilterCategory int64
    PageLimit      int64
    PageOffset     int64
}
```

### Query Naming Conventions

**Naming Pattern:** `[Action][Entity][Qualifier]`

**Examples:**
- `GetProduct` - Get single product by ID
- `GetProductBySlug` - Get single product by slug
- `ListProducts` - List all products
- `ListPublishedProducts` - List published products only
- `ListProductsByCategory` - List products filtered by category
- `CreateProduct` - Insert new product
- `UpdateProduct` - Update existing product
- `DeleteProduct` - Delete product
- `CountProducts` - Count products
- `SearchProducts` - Full-text search products

### Common Query Patterns

**Simple Read:**
```sql
-- name: GetProduct :one
SELECT * FROM products WHERE id = ? LIMIT 1;
```

**List with Pagination:**
```sql
-- name: ListProducts :many
SELECT * FROM products
WHERE status = 'published'
ORDER BY published_at DESC
LIMIT ? OFFSET ?;
```

**Complex Join:**
```sql
-- name: ListPublishedPosts :many
SELECT
    bp.id, bp.title, bp.slug,
    bc.name AS category_name,
    ba.name AS author_name
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.status = 'published'
ORDER BY bp.published_at DESC
LIMIT ? OFFSET ?;
```

**Conditional Filtering:**
```sql
-- name: ListProductsFiltered :many
SELECT * FROM products
WHERE
    (CASE WHEN @status = '' THEN 1 ELSE status = @status END)
    AND (CASE WHEN @category_id = 0 THEN 1 ELSE category_id = @category_id END)
ORDER BY created_at DESC;
```

**Insert with RETURNING:**
```sql
-- name: CreateProduct :one
INSERT INTO products (
    sku, slug, name, description, category_id, status
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;
```

**Update:**
```sql
-- name: UpdateProduct :exec
UPDATE products
SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
```

**Delete:**
```sql
-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;
```

**Count:**
```sql
-- name: CountProducts :one
SELECT COUNT(*) FROM products WHERE status = 'published';
```

## How to Add New Queries

### Step-by-Step Process

1. **Identify the Table and Operation**
   - Determine which table(s) you need to query
   - Decide operation type (SELECT, INSERT, UPDATE, DELETE)

2. **Create or Edit Query File**
   - Navigate to `db/queries/`
   - Edit existing file or create new one (e.g., `db/queries/products.sql`)

3. **Write SQL Query with sqlc Annotation**
   ```sql
   -- name: GetProductBySKU :one
   SELECT * FROM products WHERE sku = ? LIMIT 1;
   ```

4. **Choose Correct Return Type**
   - `:one` for single row (errors if not found or multiple)
   - `:many` for multiple rows (empty slice if none found)
   - `:exec` for statements that don't return data

5. **Generate Go Code**
   ```bash
   sqlc generate
   ```

6. **Verify Generated Code**
   - Check `db/sqlc/` for new function
   - Review parameter types and return values

7. **Use in Handler**
   ```go
   product, err := queries.GetProductBySKU(ctx, "SKU-12345")
   if err != nil {
       // handle error
   }
   ```

### Example: Adding a New Query

**Requirement:** Get featured products for homepage

**Step 1:** Add query to `db/queries/products.sql`
```sql
-- name: ListFeaturedProducts :many
SELECT p.*, pc.slug AS category_slug
FROM products p
INNER JOIN product_categories pc ON p.category_id = pc.id
WHERE p.is_featured = 1 AND p.status = 'published'
ORDER BY p.featured_order ASC
LIMIT ?;
```

**Step 2:** Generate code
```bash
sqlc generate
```

**Step 3:** Use in handler
```go
func (h *Handler) GetHomepage(c echo.Context) error {
    queries := sqlc.New(h.db)

    featuredProducts, err := queries.ListFeaturedProducts(c.Request().Context(), 6)
    if err != nil {
        return err
    }

    return c.Render(http.StatusOK, "homepage", map[string]interface{}{
        "FeaturedProducts": featuredProducts,
    })
}
```

### Query Development Tips

1. **Test SQL First**: Test raw SQL in sqlite3 CLI before adding to query file
2. **Use Descriptive Names**: Function names should clearly indicate purpose
3. **Add Comments**: Document complex queries with SQL comments
4. **Consider Performance**: Add indexes for frequently queried columns
5. **Handle NULL Values**: Use `sql.Null*` types for nullable columns
6. **Validate Joins**: Ensure foreign key relationships are correct

## Seed Data

### Available Seed Files

Seed files populate the database with initial/sample data. Located in `db/seeds/`:

- `003_product_categories.sql` - Product category taxonomy
- `004_blog_categories.sql` - Blog category taxonomy
- `004_solutions.sql` - Industry solutions
- `005_blog_authors.sql` - Sample blog authors
- `006_industries.sql` - Industry taxonomy for case studies
- `007_partner_tiers.sql` - Partner tier levels
- `008_whitepaper_topics.sql` - Whitepaper topic taxonomy
- `009_blog_tags.sql` - Blog tags
- `009_products.sql` - Sample products with full product data
- `010_blog_posts.sql` - Sample blog posts
- `011_blog_post_tags.sql` - Blog post tag relationships
- `012_case_studies.sql` - Sample case studies with metrics
- `020_whitepapers.sql` - Sample whitepapers with learning points
- `020b_whitepaper_downloads.sql` - Sample whitepaper download records
- `021_contact.sql` - Contact page content and office locations
- `022_about.sql` - About page content (company overview, mission, values, milestones, certifications)
- `023_partners.sql` - Partner organizations and testimonials
- `024_homepage.sql` - Homepage content (hero, stats, testimonials, CTA)

### Seed File Format

Standard INSERT statements:

```sql
INSERT INTO product_categories (name, slug, description, icon, image_url, product_count, sort_order) VALUES
('Desktops', 'desktops', 'High-performance desktop computers...', 'computer', '/uploads/categories/desktops.jpg', 12, 1),
('OPS Modules', 'ops-modules', 'Open Pluggable Specification modules...', 'memory', '/uploads/categories/ops-modules.jpg', 8, 2);
```

### Using Seed Data

**Manual Seeding:**
```bash
# Seed a single file
sqlite3 data/cms.db < db/seeds/003_product_categories.sql

# Seed all files in order
for file in db/seeds/*.sql; do
    sqlite3 data/cms.db < "$file"
done
```

**Programmatic Seeding:**
```go
import (
    "database/sql"
    "io/ioutil"
)

func SeedDatabase(db *sql.DB, seedFile string) error {
    content, err := ioutil.ReadFile(seedFile)
    if err != nil {
        return err
    }

    _, err = db.Exec(string(content))
    return err
}
```

### Development vs Production

**Development:**
- Seed comprehensive sample data for testing UI/UX
- Include realistic content for all features
- Use seed data for local development

**Production:**
- Seed only essential taxonomies (categories, tiers, topics)
- Skip sample content (blog posts, products, case studies)
- Let admins create actual content via CMS

## Backup Strategy

### Litestream to S3

Bluejay CMS uses **Litestream** for continuous replication to S3-compatible storage.

### Why Litestream?

- **Continuous Backup**: Streams changes to S3 in real-time
- **Point-in-Time Recovery**: Restore to any point in history
- **Disaster Recovery**: Automatic failover to S3 backup
- **Cost-Effective**: Minimal storage costs with S3
- **Zero Downtime**: Backups don't block database operations
- **Simple Setup**: Single binary, simple configuration

### Litestream Configuration

Typical configuration (litestream.yml):

```yaml
dbs:
  - path: /var/www/bluejay-cms/bluejay.db
    replicas:
      - type: s3
        bucket: your-backup-bucket
        path: bluejay-cms
        region: us-east-1
        # access-key-id: $AWS_ACCESS_KEY_ID
        # secret-access-key: $AWS_SECRET_ACCESS_KEY
```

### Backup Operations

**Start Replication:**
```bash
litestream replicate
```

**List Snapshots:**
```bash
litestream snapshots /path/to/cms.db
```

**Restore from Backup:**
```bash
litestream restore -o /path/to/cms.db s3://bucket-name/db-backups
```

**Restore to Specific Time:**
```bash
litestream restore -timestamp 2024-01-15T10:30:00Z -o /path/to/cms.db s3://bucket-name/db-backups
```

### Manual Backup

For ad-hoc backups without Litestream:

**Simple File Copy:**
```bash
# Backup
sqlite3 data/cms.db ".backup 'data/cms-backup-$(date +%Y%m%d).db'"

# Restore
cp data/cms-backup-20240115.db data/cms.db
```

**SQL Dump:**
```bash
# Backup
sqlite3 data/cms.db .dump > backup.sql

# Restore
sqlite3 data/cms.db < backup.sql
```

### Backup Best Practices

1. **Automate Backups**: Use Litestream for continuous replication
2. **Test Restores**: Regularly test restore procedures
3. **Monitor Replication**: Check Litestream logs for errors
4. **Retention Policy**: Keep 30+ days of backups
5. **Offsite Storage**: Use S3 for geographic redundancy
6. **Pre-Migration Backups**: Always backup before running migrations

## Common Query Patterns

### Pattern 1: List with Pagination

**Use Case:** Paginated listings (products, blog posts, etc.)

```sql
-- name: ListProducts :many
SELECT * FROM products
WHERE status = 'published'
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: CountProducts :one
SELECT COUNT(*) FROM products WHERE status = 'published';
```

**Usage:**
```go
perPage := 20
page := 1
offset := (page - 1) * perPage

products, _ := queries.ListProducts(ctx, perPage, offset)
total, _ := queries.CountProducts(ctx)
totalPages := (total + perPage - 1) / perPage
```

### Pattern 2: Filtered Listings

**Use Case:** Admin lists with multiple filter criteria

```sql
-- name: ListProductsFiltered :many
SELECT * FROM products
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE category_id = @filter_category END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE name LIKE '%' || @filter_search || '%' END)
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;
```

**Usage:**
```go
params := ListProductsFilteredParams{
    FilterStatus:   "published",  // or "" for all
    FilterCategory: 5,            // or 0 for all
    FilterSearch:   "desktop",    // or "" for all
    PageLimit:      20,
    PageOffset:     0,
}
products, _ := queries.ListProductsFiltered(ctx, params)
```

### Pattern 3: Single Record with Details

**Use Case:** Product/blog post detail page with joined data

```sql
-- name: GetPublishedPostBySlug :one
SELECT
    bp.id, bp.title, bp.slug, bp.excerpt, bp.body,
    bp.featured_image_url, bp.featured_image_alt,
    bc.name AS category_name, bc.slug AS category_slug,
    ba.name AS author_name, ba.bio AS author_bio,
    bp.reading_time_minutes, bp.published_at
FROM blog_posts bp
INNER JOIN blog_categories bc ON bp.category_id = bc.id
INNER JOIN blog_authors ba ON bp.author_id = ba.id
WHERE bp.slug = ? AND bp.status = 'published';
```

### Pattern 4: Many-to-Many Relationships

**Use Case:** Fetching related records (tags, products, etc.)

```sql
-- name: GetPostTags :many
SELECT bt.id, bt.name, bt.slug
FROM blog_tags bt
INNER JOIN blog_post_tags bpt ON bt.id = bpt.blog_tag_id
WHERE bpt.blog_post_id = ?
ORDER BY bt.name;

-- name: AddTagToPost :exec
INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
VALUES (?, ?)
ON CONFLICT DO NOTHING;

-- name: ClearPostTags :exec
DELETE FROM blog_post_tags WHERE blog_post_id = ?;
```

**Usage:**
```go
// Get tags for post
tags, _ := queries.GetPostTags(ctx, postID)

// Update tags (clear then add)
queries.ClearPostTags(ctx, postID)
for _, tagID := range newTagIDs {
    queries.AddTagToPost(ctx, postID, tagID)
}
```

### Pattern 5: Full-Text Search

**Use Case:** Search across products, blog posts, case studies

```sql
-- name: SearchProducts :many
SELECT p.* FROM products p
WHERE p.id IN (
    SELECT rowid FROM products_fts WHERE products_fts MATCH ?
)
AND p.status = 'published'
ORDER BY p.published_at DESC
LIMIT ? OFFSET ?;
```

**Usage:**
```go
searchQuery := "desktop computer touchscreen"
products, _ := queries.SearchProducts(ctx, searchQuery, 20, 0)
```

### Pattern 6: Increment Counter

**Use Case:** Download counts, view counts, etc.

```sql
-- name: IncrementDownloadCount :exec
UPDATE product_downloads
SET download_count = download_count + 1
WHERE id = ?;

-- name: IncrementWhitepaperDownloads :exec
UPDATE whitepapers
SET download_count = download_count + 1
WHERE id = ?;
```

### Pattern 7: Conditional INSERT (Upsert)

**Use Case:** Insert if not exists, ignore duplicates

```sql
-- name: AddTagToPost :exec
INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
VALUES (?, ?)
ON CONFLICT DO NOTHING;
```

### Pattern 8: Bulk Delete with CASCADE

**Use Case:** Deleting parent records with child dependencies

```sql
-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;
-- Child records (specs, images, features) automatically deleted via CASCADE
```

### Pattern 9: Ordered List with Display Order

**Use Case:** Draggable admin lists, featured items

```sql
-- name: ListFeaturedProducts :many
SELECT * FROM products
WHERE is_featured = 1 AND status = 'published'
ORDER BY featured_order ASC
LIMIT ?;

-- name: UpdateProductOrder :exec
UPDATE products
SET featured_order = ?
WHERE id = ?;
```

### Pattern 10: Singleton Settings Pattern

**Use Case:** Global site settings (always id=1)

```sql
-- name: GetSettings :one
SELECT * FROM settings WHERE id = 1 LIMIT 1;

-- name: UpdateSettings :exec
UPDATE settings
SET site_name = ?, site_tagline = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = 1;
```

---

## Summary

This document provides a complete reference for the Bluejay CMS database architecture, covering:

- SQLite configuration with WAL mode, pragmas, and connection pooling
- Complete schema with 50+ tables organized by domain
- Entity relationships and foreign key constraints
- Migration strategy using golang-migrate
- sqlc code generation workflow and patterns
- Seed data management
- Backup strategy with Litestream
- Common query patterns for typical CMS operations

For questions or schema changes, refer to:
- Migration files: `db/migrations/`
- Query files: `db/queries/`
- Generated code: `db/sqlc/`
- Database config: `internal/database/sqlite.go`
