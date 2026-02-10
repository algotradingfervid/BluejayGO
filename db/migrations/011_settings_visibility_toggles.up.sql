-- Header nav link toggles
ALTER TABLE settings ADD COLUMN show_nav_home BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_nav_about BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_nav_products BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_nav_solutions BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_nav_blog BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_nav_partners BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_nav_contact BOOLEAN NOT NULL DEFAULT 1;
-- Footer section toggles
ALTER TABLE settings ADD COLUMN show_footer_about BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_footer_socials BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_footer_products BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_footer_solutions BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_footer_resources BOOLEAN NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN show_footer_contact BOOLEAN NOT NULL DEFAULT 1;
