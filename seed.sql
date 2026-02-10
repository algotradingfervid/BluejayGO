DELETE FROM admin_users;
DELETE FROM settings;

INSERT OR REPLACE INTO settings (
    id, site_name, site_tagline, contact_email, contact_phone,
    footer_text, meta_description, social_linkedin, social_twitter, social_github
) VALUES (
    1, 'BlueJay Innovative Labs', 'Innovation Through Technology',
    'info@bluejaylabs.com', '+1 (555) 123-4567',
    '© 2024 BlueJay Innovative Labs. All rights reserved.',
    'BlueJay Innovative Labs delivers cutting-edge technology solutions for modern businesses.',
    'https://linkedin.com/company/bluejaylabs',
    'https://twitter.com/bluejaylabs',
    'https://github.com/bluejaylabs'
);

INSERT INTO admin_users (email, password_hash, display_name, role, is_active)
VALUES (
    'admin@bluejaylabs.com',
    '$2a$10$8CjuoFtzSCcB7PjzJqXw5OVHev9Xtbas9y2KIoPYTYQTBuI.9dD8S',
    'Admin User', 'admin', 1
);

INSERT INTO admin_users (email, password_hash, display_name, role, is_active)
VALUES (
    'editor@bluejaylabs.com',
    '$2a$10$DwPGYt416wmeggliaWKZcOYsG/KXuxNg2EjECwCWpc1wFES1R8l/6',
    'Editor User', 'editor', 1
);

-- Phase 2: Master table seed data
DELETE FROM product_categories;
DELETE FROM blog_categories;
DELETE FROM blog_authors;
DELETE FROM industries;
DELETE FROM partner_tiers;
DELETE FROM whitepaper_topics;

.read db/seeds/003_product_categories.sql
.read db/seeds/004_blog_categories.sql
.read db/seeds/005_blog_authors.sql
.read db/seeds/006_industries.sql
.read db/seeds/007_partner_tiers.sql
.read db/seeds/008_whitepaper_topics.sql

-- Phase 3: Product seed data
DELETE FROM product_images;
DELETE FROM product_downloads;
DELETE FROM product_certifications;
DELETE FROM product_features;
DELETE FROM product_specs;
DELETE FROM case_study_products;
DELETE FROM products;

.read db/seeds/009_products.sql

-- Phase 4: Solutions seed data (after products for FK references)
DELETE FROM solution_ctas;
DELETE FROM solution_products;
DELETE FROM solution_challenges;
DELETE FROM solution_stats;
DELETE FROM solution_page_features;
DELETE FROM solutions_listing_cta;
DELETE FROM solutions;

.read db/seeds/004_solutions.sql

-- Phase 5: Blog seed data
DELETE FROM blog_post_tags;
DELETE FROM blog_tags;
DELETE FROM blog_posts;

-- Phase 6: Case study seed data
DELETE FROM case_study_metrics;
DELETE FROM case_study_products;
DELETE FROM case_studies;

.read db/seeds/012_case_studies.sql

.read db/seeds/009_blog_tags.sql
.read db/seeds/010_blog_posts.sql
.read db/seeds/011_blog_post_tags.sql

-- Phase 8: Whitepapers and Contact seed data
DELETE FROM whitepaper_downloads;
DELETE FROM whitepaper_learning_points;
DELETE FROM whitepapers;
DELETE FROM office_locations;
DELETE FROM contact_submissions;

.read db/seeds/020_whitepapers.sql
.read db/seeds/020b_whitepaper_downloads.sql
.read db/seeds/021_contact.sql

-- Phase 7: About Us and Partners seed data
DELETE FROM partner_testimonials;
DELETE FROM partners;
DELETE FROM certifications;
DELETE FROM milestones;
DELETE FROM core_values;
DELETE FROM mission_vision_values;
DELETE FROM company_overview;

.read db/seeds/022_about.sql
.read db/seeds/023_partners.sql

-- Phase 9: Homepage seed data
DELETE FROM homepage_cta;
DELETE FROM homepage_testimonials;
DELETE FROM homepage_stats;
DELETE FROM homepage_hero;

.read db/seeds/024_homepage.sql
