-- Add nav label fields to settings
ALTER TABLE settings ADD COLUMN nav_label_home TEXT NOT NULL DEFAULT 'Home';
ALTER TABLE settings ADD COLUMN nav_label_about TEXT NOT NULL DEFAULT 'About Us';
ALTER TABLE settings ADD COLUMN nav_label_products TEXT NOT NULL DEFAULT 'Products';
ALTER TABLE settings ADD COLUMN nav_label_solutions TEXT NOT NULL DEFAULT 'Solutions';
ALTER TABLE settings ADD COLUMN nav_label_blog TEXT NOT NULL DEFAULT 'Blog';
ALTER TABLE settings ADD COLUMN nav_label_partners TEXT NOT NULL DEFAULT 'Partners';
ALTER TABLE settings ADD COLUMN nav_label_contact TEXT NOT NULL DEFAULT 'Contact Us';

-- Add footer column heading fields to settings
ALTER TABLE settings ADD COLUMN footer_heading_products TEXT NOT NULL DEFAULT 'Products';
ALTER TABLE settings ADD COLUMN footer_heading_solutions TEXT NOT NULL DEFAULT 'Solutions';
ALTER TABLE settings ADD COLUMN footer_heading_resources TEXT NOT NULL DEFAULT 'Resources';
ALTER TABLE settings ADD COLUMN footer_heading_contact TEXT NOT NULL DEFAULT 'Contact Us';

-- Product detail sub-section headings
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('product_detail', 'overview_section', 'Overview', 1);
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('product_detail', 'video_section', 'Product Video', 2);
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('product_detail', 'features_section', 'Key Features', 3);
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('product_detail', 'specs_section', 'Technical Specifications', 4);
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('product_detail', 'certifications_section', 'Certifications', 5);
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('product_detail', 'downloads_section', 'Downloads & Resources', 6);

-- Products category page sections
INSERT INTO page_sections (page_key, section_key, label, heading, display_order) VALUES
('products_category', 'hero', 'Product Category', '', 1);
INSERT INTO page_sections (page_key, section_key, heading, primary_button_text, primary_button_url, display_order) VALUES
('products_category', 'empty_state', 'No products found', 'Back to All Products', '/products', 2);

-- Footer resource links as page sections
INSERT INTO page_sections (page_key, section_key, heading, primary_button_url, display_order) VALUES
('footer', 'resource_1', 'Case Studies', '/case-studies', 1);
INSERT INTO page_sections (page_key, section_key, heading, primary_button_url, display_order) VALUES
('footer', 'resource_2', 'Whitepapers', '/whitepapers', 2);
INSERT INTO page_sections (page_key, section_key, heading, primary_button_url, display_order) VALUES
('footer', 'resource_3', 'Blog', '/blog', 3);
INSERT INTO page_sections (page_key, section_key, heading, primary_button_url, display_order) VALUES
('footer', 'resource_4', 'Support', '/support', 4);

-- Watch button label for product video
UPDATE page_sections SET label = 'Watch Product Video' WHERE page_key = 'product_detail' AND section_key = 'video_section';
