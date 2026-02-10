-- name: GetSettings :one
SELECT * FROM settings WHERE id = 1 LIMIT 1;

-- name: UpdateSettings :exec
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

-- name: UpdateHeaderSettings :exec
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
UPDATE settings
SET about_show_mission = ?,
    about_show_milestones = ?,
    about_show_certifications = ?,
    about_show_team = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateProductsSettings :exec
UPDATE settings
SET products_per_page = ?,
    products_show_categories = ?,
    products_show_search = ?,
    products_default_sort = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateSolutionsSettings :exec
UPDATE settings
SET solutions_per_page = ?,
    solutions_show_industries = ?,
    solutions_show_search = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = 1;

-- name: UpdateBlogSettings :exec
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
