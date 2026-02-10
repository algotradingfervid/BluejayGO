-- Homepage Hero
INSERT INTO homepage_hero (headline, subheadline, badge_text, primary_cta_text, primary_cta_url, secondary_cta_text, secondary_cta_url, is_active, display_order) VALUES
('Innovative Solutions With Unmatched Support', 'Transforming spaces with cutting-edge interactive technology and dedicated customer excellence.', 'System Status: Operational', 'Explore Products', '/products', 'Request a Demo', '/contact', 1, 1);

-- Homepage Stats
INSERT INTO homepage_stats (stat_value, stat_label, display_order, is_active) VALUES
('15+', 'Years of Experience', 1, 1),
('5,000+', 'Products Deployed', 2, 1),
('50+', 'Countries Served', 3, 1),
('98%', 'Customer Satisfaction', 4, 1);

-- Homepage Testimonials
INSERT INTO homepage_testimonials (quote, author_name, author_title, author_company, rating, display_order, is_active) VALUES
('BlueJay''s interactive displays have completely transformed how we conduct training sessions. The support team is incredibly responsive and the product quality is outstanding.', 'Sarah Johnson', 'IT Director', 'Global Education Inc.', 5, 1, 1),
('The implementation was seamless and the results exceeded our expectations. Our conference rooms are now state-of-the-art collaboration spaces.', 'Michael Chen', 'CTO', 'TechVentures Corp', 5, 2, 1),
('Outstanding product quality and unmatched customer service. BlueJay has been our trusted technology partner for over 5 years.', 'Emily Rodriguez', 'Operations Manager', 'Metro Healthcare Systems', 5, 3, 1);

-- Homepage CTA
INSERT INTO homepage_cta (headline, description, primary_cta_text, primary_cta_url, secondary_cta_text, secondary_cta_url, background_style, is_active) VALUES
('Ready to Transform Your Space?', 'Let our experts help you find the perfect solution for your organization''s unique needs.', 'Schedule a Demo', '/contact', 'Contact Sales', '/contact', 'primary', 1);

-- Mark some partners as featured for homepage
UPDATE partners SET is_featured = 1 WHERE display_order <= 7;

-- Mark some products as featured if not already
UPDATE products SET is_featured = 1, featured_order = id WHERE is_featured = 0 AND status = 'published' LIMIT 6;
