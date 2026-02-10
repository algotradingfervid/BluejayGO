CREATE TABLE IF NOT EXISTS page_sections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    page_key TEXT NOT NULL,
    section_key TEXT NOT NULL,
    heading TEXT NOT NULL DEFAULT '',
    subheading TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    label TEXT NOT NULL DEFAULT '',
    primary_button_text TEXT NOT NULL DEFAULT '',
    primary_button_url TEXT NOT NULL DEFAULT '',
    secondary_button_text TEXT NOT NULL DEFAULT '',
    secondary_button_url TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT 1,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_page_sections_key ON page_sections(page_key, section_key);

-- Home page
INSERT INTO page_sections (page_key, section_key, heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, display_order) VALUES
('home', 'hero', '', '', 'Get Started', '/contact', 'Learn More', '/about', 1);
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('home', 'solutions_section', 'Our Solutions', 2);

-- Products listing page
INSERT INTO page_sections (page_key, section_key, heading, description, label, display_order) VALUES
('products', 'hero', 'Our Products', 'Discover our comprehensive range of innovative computing and display solutions designed for modern enterprises.', 'Product Catalog', 1);
INSERT INTO page_sections (page_key, section_key, heading, display_order) VALUES
('products', 'categories_section', 'Browse by Product Type', 2);
INSERT INTO page_sections (page_key, section_key, heading, description, primary_button_text, primary_button_url, display_order) VALUES
('products', 'cta', 'Need Help Choosing?', 'Our product specialists can help you find the perfect solution for your organization''s needs.', 'Contact Sales', '/contact', 3);

-- Product detail page
INSERT INTO page_sections (page_key, section_key, heading, description, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, display_order) VALUES
('product_detail', 'cta', 'Interested in {product_name}?', 'Get in touch with our sales team for pricing, bulk orders, and custom configurations.', 'Request a Quote', '/contact?product={product_sku}', 'Contact Sales', '/contact', 1);

-- Solutions listing page
INSERT INTO page_sections (page_key, section_key, heading, description, label, display_order) VALUES
('solutions', 'hero', 'Industry Solutions', 'Tailored technology solutions designed to address the unique challenges of your industry. From education to healthcare, we deliver innovative products that drive transformation.', 'Industry Solutions', 1);
INSERT INTO page_sections (page_key, section_key, heading, label, display_order) VALUES
('solutions', 'grid_section', 'Explore Our Industry Solutions', 'Section 01', 2);
INSERT INTO page_sections (page_key, section_key, heading, label, display_order) VALUES
('solutions', 'features_section', 'Why Choose BlueJay', 'Section 02', 3);

-- Solution detail page
INSERT INTO page_sections (page_key, section_key, heading, label, display_order) VALUES
('solution_detail', 'overview_section', 'Industry Overview', 'Section 01', 1);
INSERT INTO page_sections (page_key, section_key, heading, label, display_order) VALUES
('solution_detail', 'challenges_section', 'Challenges We Address', 'Section 02', 2);
INSERT INTO page_sections (page_key, section_key, heading, label, display_order) VALUES
('solution_detail', 'products_section', 'Recommended Products for {solution_title}', 'Section 03', 3);
INSERT INTO page_sections (page_key, section_key, heading, label, display_order) VALUES
('solution_detail', 'other_solutions', 'Explore Other Industry Solutions', 'Section 05', 4);
INSERT INTO page_sections (page_key, section_key, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, display_order) VALUES
('solution_detail', 'hero_buttons', 'Get a Quote', '/contact', 'View Products', '/products', 5);
