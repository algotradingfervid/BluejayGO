-- ====================================================================
-- SETTINGS QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing global site settings.
--
-- Entity: settings table (singleton pattern - only one row with id=1)
-- Purpose: Store site-wide configuration, metadata, feature toggles
--
-- Settings categories:
--   - Global: Site name, contact info, social media links
--   - Header: Logo, CTA buttons, navigation toggles
--   - Homepage: Section visibility, carousel settings
--   - Page-specific: About, Products, Solutions, Blog configurations
--
-- Note: All queries target id=1 (singleton row)
-- ====================================================================

-- name: GetSettings :one
-- Retrieves the global settings record (singleton).
--
-- Parameters: none
-- Returns: Settings - The single settings record
--
-- Note: Always targets id=1; settings table should only contain one row
-- Use case: Loading site configuration on application startup, template rendering
SELECT * FROM settings WHERE id = 1 LIMIT 1;

-- name: UpdateSettings :exec
-- Updates the global settings record with comprehensive site configuration.
--
-- Parameters (47 total):
--   $1-$8: Site identity (name, tagline, contact info, footer text, SEO)
--   $9: google_analytics_id - GA tracking ID
--   $10-$15: Social media URLs (LinkedIn, Twitter, GitHub, Facebook, YouTube, Instagram)
--   $16-$17: Business info (hours, about text)
--   $18-$24: Navigation visibility toggles (show_nav_*)
--   $25-$30: Footer section visibility toggles (show_footer_*)
--   $31-$37: Navigation labels (nav_label_*)
--   $38-$41: Footer section headings (footer_heading_*)
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Note: updated_at is automatically set to CURRENT_TIMESTAMP
-- Use case: Comprehensive settings update from admin settings page (legacy query)
-- Recommendation: Use specific Update*Settings queries for better maintainability
UPDATE settings
SET site_name = ?,
    site_tagline = ?,
    contact_email = ?,
    contact_phone = ?,
    address = ?,
    footer_text = ?,
    meta_description = ?,
    meta_keywords = ?,
    google_analytics_id = ?,
    social_linkedin = ?,
    social_twitter = ?,
    social_github = ?,
    social_facebook = ?,
    social_youtube = ?,
    social_instagram = ?,
    business_hours = ?,
    about_text = ?,
    show_nav_home = ?,
    show_nav_about = ?,
    show_nav_products = ?,
    show_nav_solutions = ?,
    show_nav_blog = ?,
    show_nav_partners = ?,
    show_nav_contact = ?,
    show_footer_about = ?,
    show_footer_socials = ?,
    show_footer_products = ?,
    show_footer_solutions = ?,
    show_footer_resources = ?,
    show_footer_contact = ?,
    nav_label_home = ?,
    nav_label_about = ?,
    nav_label_products = ?,
    nav_label_solutions = ?,
    nav_label_blog = ?,
    nav_label_partners = ?,
    nav_label_contact = ?,
    footer_heading_products = ?,
    footer_heading_solutions = ?,
    footer_heading_resources = ?,
    footer_heading_contact = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- ====================================================================
-- SETTINGS - SECTION-SPECIFIC UPDATES
-- ====================================================================

-- name: UpdateHeaderSettings :exec
-- Updates header and navigation-specific settings.
--
-- Parameters:
--   $1-$2: Logo settings (path, alt text)
--   $3-$6: Header CTA button (enabled, text, URL, style)
--   $7-$10: Header contact display toggles (phone, email, social, social style)
--   $11-$18: Navigation item visibility toggles (show_nav_*)
--   $19-$26: Navigation item custom labels (nav_label_*)
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Admin header/navigation settings page
-- Note: Scoped update - only affects header-related fields
UPDATE settings
SET header_logo_path = ?,
    header_logo_alt = ?,
    header_cta_enabled = ?,
    header_cta_text = ?,
    header_cta_url = ?,
    header_cta_style = ?,
    header_show_phone = ?,
    header_show_email = ?,
    header_show_social = ?,
    header_social_style = ?,
    show_nav_products = ?,
    show_nav_solutions = ?,
    show_nav_case_studies = ?,
    show_nav_about = ?,
    show_nav_blog = ?,
    show_nav_whitepapers = ?,
    show_nav_partners = ?,
    show_nav_contact = ?,
    nav_label_products = ?,
    nav_label_solutions = ?,
    nav_label_case_studies = ?,
    nav_label_about = ?,
    nav_label_blog = ?,
    nav_label_whitepapers = ?,
    nav_label_partners = ?,
    nav_label_contact = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateHomepageSettings :exec
-- Updates homepage-specific feature toggles and limits.
--
-- Parameters:
--   $1-$4: Section visibility toggles (heroes, stats, testimonials, CTA)
--   $5-$7: Maximum items to display (heroes, stats, testimonials)
--   $8-$9: Hero carousel settings (autoplay enabled, interval in milliseconds)
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Admin homepage configuration page
-- Note: Controls which homepage sections are visible and their behavior
UPDATE settings
SET homepage_show_heroes = ?,
    homepage_show_stats = ?,
    homepage_show_testimonials = ?,
    homepage_show_cta = ?,
    homepage_max_heroes = ?,
    homepage_max_stats = ?,
    homepage_max_testimonials = ?,
    homepage_hero_autoplay = ?,
    homepage_hero_interval = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateAboutSettings :exec
-- Updates About page section visibility toggles.
--
-- Parameters:
--   $1-$4: Section visibility toggles (mission, milestones, certifications, team)
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Admin About page configuration
-- Note: Controls which About page sections are displayed
UPDATE settings
SET about_show_mission = ?,
    about_show_milestones = ?,
    about_show_certifications = ?,
    about_show_team = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateProductsSettings :exec
-- Updates Products page display and filter settings.
--
-- Parameters:
--   $1: products_per_page - Number of products per page (pagination)
--   $2: products_show_categories - Show category filter
--   $3: products_show_search - Show search bar
--   $4: products_default_sort - Default sort order ("newest", "name", "price", etc.)
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Admin Products page configuration
UPDATE settings
SET products_per_page = ?,
    products_show_categories = ?,
    products_show_search = ?,
    products_default_sort = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateSolutionsSettings :exec
-- Updates Solutions page display and filter settings.
--
-- Parameters:
--   $1: solutions_per_page - Number of solutions per page (pagination)
--   $2: solutions_show_industries - Show industry filter
--   $3: solutions_show_search - Show search bar
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Admin Solutions page configuration
UPDATE settings
SET solutions_per_page = ?,
    solutions_show_industries = ?,
    solutions_show_search = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateBlogSettings :exec
-- Updates Blog page display and metadata settings.
--
-- Parameters:
--   $1: blog_posts_per_page - Number of posts per page (pagination)
--   $2-$5: Metadata visibility toggles (author, date, categories, tags)
--   $6: blog_show_search - Show search bar
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Admin Blog page configuration
-- Note: Controls blog post listing metadata and features
UPDATE settings
SET blog_posts_per_page = ?,
    blog_show_author = ?,
    blog_show_date = ?,
    blog_show_categories = ?,
    blog_show_tags = ?,
    blog_show_search = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateGlobalSettings :exec
-- Updates site-wide global settings (identity, contact, SEO, social).
--
-- Parameters:
--   $1-$2: Site identity (name, tagline)
--   $3-$6: Contact information (email, phone, address, hours)
--   $7-$8: SEO metadata (meta_description, meta_keywords)
--   $9: google_analytics_id - GA tracking ID
--   $10-$14: Social media URLs (Facebook, Twitter, LinkedIn, Instagram, YouTube)
--
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Admin global settings page (site identity and contact info)
-- Note: Most commonly updated settings for basic site configuration
UPDATE settings
SET site_name = ?,
    site_tagline = ?,
    contact_email = ?,
    contact_phone = ?,
    address = ?,
    business_hours = ?,
    meta_description = ?,
    meta_keywords = ?,
    google_analytics_id = ?,
    social_facebook = ?,
    social_twitter = ?,
    social_linkedin = ?,
    social_instagram = ?,
    social_youtube = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;
