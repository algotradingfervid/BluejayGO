-- Seed products for Desktops category
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-D100', 'bj-d100-entry-desktop', 'BJ-D100', 'Entry Desktop', 'Affordable desktop solution for small businesses and home offices.', 'The BJ-D100 Entry Desktop is designed for everyday computing tasks with reliable performance.', (SELECT id FROM product_categories WHERE slug = 'desktops'), 'published', 0, NULL, '/uploads/products/bj-d100.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-D200', 'bj-d200-business-desktop', 'BJ-D200', 'Business Desktop', 'High-performance desktop solution for enterprise productivity.', 'The BJ-D200 Business Desktop is engineered for demanding enterprise environments.', (SELECT id FROM product_categories WHERE slug = 'desktops'), 'published', 1, 1, '/uploads/products/bj-d200.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-D300', 'bj-d300-compact-desktop', 'BJ-D300', 'Compact Desktop', 'Space-saving mini PC for modern workspaces.', 'The BJ-D300 Compact Desktop delivers powerful performance in an ultra-small form factor.', (SELECT id FROM product_categories WHERE slug = 'desktops'), 'published', 0, NULL, '/uploads/products/bj-d300.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-D500', 'bj-d500-premium-workstation', 'BJ-D500', 'Premium Workstation', 'High-end workstation for creative professionals and power users.', 'The BJ-D500 Premium Workstation delivers outstanding performance for demanding workloads.', (SELECT id FROM product_categories WHERE slug = 'desktops'), 'published', 1, 2, '/uploads/products/bj-d500.jpg', datetime('now'));

-- Seed products for OPS Modules
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-OPS100', 'bj-ops100-standard-module', 'BJ-OPS100', 'Standard OPS Module', 'Standard OPS computing module for interactive flat panels.', 'The BJ-OPS100 provides reliable computing for interactive displays.', (SELECT id FROM product_categories WHERE slug = 'ops-modules'), 'published', 0, NULL, '/uploads/products/bj-ops100.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-OPS200', 'bj-ops200-pro-module', 'BJ-OPS200', 'Pro OPS Module', 'High-performance OPS module for demanding applications.', 'The BJ-OPS200 Pro delivers enhanced performance for professional use.', (SELECT id FROM product_categories WHERE slug = 'ops-modules'), 'published', 1, 1, '/uploads/products/bj-ops200.jpg', datetime('now'));

-- Seed products for Interactive Flat Panels
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-IFP65', 'bj-ifp65-interactive-panel', 'BJ-IFP65', '65" Interactive Flat Panel', '65-inch 4K interactive display for classrooms and meeting rooms.', 'The BJ-IFP65 delivers an immersive 65-inch 4K interactive experience with 20-point touch.', (SELECT id FROM product_categories WHERE slug = 'interactive-flat-panels'), 'published', 1, 1, '/uploads/products/bj-ifp65.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-IFP75', 'bj-ifp75-interactive-panel', 'BJ-IFP75', '75" Interactive Flat Panel', '75-inch 4K interactive display for large classrooms and boardrooms.', 'The BJ-IFP75 provides a massive 75-inch 4K canvas for collaboration and teaching.', (SELECT id FROM product_categories WHERE slug = 'interactive-flat-panels'), 'published', 1, 2, '/uploads/products/bj-ifp75.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-IFP86', 'bj-ifp86-interactive-panel', 'BJ-IFP86', '86" Interactive Flat Panel', '86-inch 4K interactive display for auditoriums and large meeting spaces.', 'The BJ-IFP86 is our largest interactive panel, perfect for auditoriums and large venues.', (SELECT id FROM product_categories WHERE slug = 'interactive-flat-panels'), 'published', 0, NULL, '/uploads/products/bj-ifp86.jpg', datetime('now'));

-- Seed products for AV Accessories
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-WPS100', 'bj-wps100-wireless-presenter', 'BJ-WPS100', 'Wireless Presentation System', 'One-click wireless screen sharing for meeting rooms.', 'The BJ-WPS100 enables seamless wireless screen sharing from any device.', (SELECT id FROM product_categories WHERE slug = 'av-accessories'), 'published', 1, 1, '/uploads/products/bj-wps100.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-CAM360', 'bj-cam360-conference-camera', 'BJ-CAM360', '360° Conference Camera', 'AI-powered panoramic conference camera with auto-framing.', 'The BJ-CAM360 uses AI to automatically frame speakers and track conversation.', (SELECT id FROM product_categories WHERE slug = 'av-accessories'), 'published', 0, NULL, '/uploads/products/bj-cam360.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-SPK200', 'bj-spk200-speakerphone', 'BJ-SPK200', 'Conference Speakerphone', 'Professional USB/Bluetooth speakerphone for meeting rooms.', 'The BJ-SPK200 delivers crystal-clear audio for conference calls in any room size.', (SELECT id FROM product_categories WHERE slug = 'av-accessories'), 'published', 0, NULL, '/uploads/products/bj-spk200.jpg', datetime('now'));

-- Seed products for IoT Products
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-DS55', 'bj-ds55-digital-signage', 'BJ-DS55', '55" Digital Signage Display', 'Commercial-grade 55-inch digital signage display for 24/7 operation.', 'The BJ-DS55 is built for continuous commercial operation with high brightness and durability.', (SELECT id FROM product_categories WHERE slug = 'iot-products'), 'published', 1, 1, '/uploads/products/bj-ds55.jpg', datetime('now'));
INSERT INTO products (sku, slug, name, tagline, description, overview, category_id, status, is_featured, featured_order, primary_image, published_at) VALUES
('BJ-IOT100', 'bj-iot100-smart-sensor', 'BJ-IOT100', 'Smart Environment Sensor', 'Multi-sensor IoT device for monitoring temperature, humidity, and occupancy.', 'The BJ-IOT100 provides real-time environmental monitoring for smart buildings.', (SELECT id FROM product_categories WHERE slug = 'iot-products'), 'published', 0, NULL, '/uploads/products/bj-iot100.jpg', datetime('now'));

-- Seed specs for BJ-D200 (assuming product_id = 2)
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Processor', 'Processor Model', 'Intel Core i5-12400', 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Processor', 'Cores / Threads', '6 Cores / 12 Threads', 2),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Processor', 'Base Clock', '2.5 GHz', 3),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Processor', 'Turbo Clock', '4.4 GHz', 4),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Memory', 'RAM', '16GB DDR4', 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Memory', 'Memory Speed', '3200MHz', 2),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Memory', 'Max Memory', '64GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Storage', 'Primary Storage', '512GB NVMe SSD', 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Storage', 'Storage Interface', 'M.2 PCIe Gen4 x4', 2),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Connectivity', 'Ethernet', 'Intel I219-V Gigabit LAN', 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Connectivity', 'Wi-Fi', 'Intel Wi-Fi 6 AX201', 2),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.2', 3);

-- Seed features for BJ-D200
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D200'), '12th Gen Intel Core i5 Processor', 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), '16GB DDR4 RAM (Expandable to 64GB)', 2),
((SELECT id FROM products WHERE sku = 'BJ-D200'), '512GB NVMe SSD Storage', 3),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Windows 11 Pro Pre-installed', 4),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'BIS & CE Certified', 5),
((SELECT id FROM products WHERE sku = 'BJ-D200'), '3-Year Warranty', 6),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Made in India', 7),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Enterprise-grade Reliability', 8);

-- Seed certifications for BJ-D200
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'BIS', 'R-41234567', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'RoHS', NULL, 'badge', 5);

-- Seed downloads for BJ-D200
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'BJ-D200 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-d200-datasheet.pdf', 2457600, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-d200-manual.pdf', 8498176, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'Driver Pack - Windows 11', 'All drivers for Windows 11 (64-bit)', 'zip', '/uploads/downloads/bj-d200-drivers.zip', 163577856, 'v2.1.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-D200'), 'BIOS Update', 'Latest BIOS firmware update utility', 'exe', '/uploads/downloads/bj-d200-bios.exe', 12582912, 'v1.08', 4);

-- Seed images for BJ-D200
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D200'), '/uploads/products/bj-d200-front.jpg', 'BJ-D200 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-D200'), '/uploads/products/bj-d200-back.jpg', 'BJ-D200 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-D200'), '/uploads/products/bj-d200-side.jpg', 'BJ-D200 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-D200'), '/uploads/products/bj-d200-ports.jpg', 'BJ-D200 Ports Detail', 'Ports Detail', 4, 0);

-- ============================================================================
-- BJ-D100 Entry Desktop (product_id = 1)
-- ============================================================================

-- Seed specs for BJ-D100
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Processor', 'Processor Model', 'Intel Core i3-12100', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Processor', 'Cores / Threads', '4 Cores / 8 Threads', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Processor', 'Base Clock', '3.3 GHz', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Processor', 'Turbo Clock', '4.3 GHz', 4),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Memory', 'RAM', '8GB DDR4', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Memory', 'Memory Speed', '3200MHz', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Memory', 'Max Memory', '32GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Storage', 'Primary Storage', '256GB NVMe SSD', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Storage', 'Storage Interface', 'M.2 PCIe Gen3', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Connectivity', 'Ethernet', 'Intel I219-LM Gigabit LAN', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 5 (AC7265)', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Connectivity', 'Bluetooth', 'Bluetooth 4.2', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Ports', 'USB Ports', '4x USB 3.2, 2x USB 2.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Ports', 'Display Outputs', 'HDMI, DisplayPort', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Ports', 'Audio', '3.5mm Headphone/Mic Combo', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Physical', 'Form Factor', 'SFF (Small Form Factor)', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Physical', 'Dimensions', '293 x 100 x 311 mm', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Physical', 'Weight', '4.2 kg', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Power', 'Power Supply', '180W 80+ Bronze', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Power', 'Operating System', 'Windows 11 Pro', 2);

-- Seed features for BJ-D100
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D100'), '12th Gen Intel Core i3 Processor', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), '8GB DDR4 RAM (Expandable to 32GB)', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), '256GB NVMe SSD Storage', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Windows 11 Pro Pre-installed', 4),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Small Form Factor Design', 5),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'BIS & CE Certified', 6),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Energy-Efficient 180W PSU', 7),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Ideal for Small Businesses and Home Offices', 8);

-- Seed certifications for BJ-D100
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'BIS', 'R-41234501', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'RoHS', NULL, 'badge', 5);

-- Seed downloads for BJ-D100
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'BJ-D100 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-d100-datasheet.pdf', 2125824, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-d100-manual.pdf', 7340032, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'Driver Pack - Windows 11', 'All drivers for Windows 11 (64-bit)', 'zip', '/uploads/downloads/bj-d100-drivers.zip', 142606336, 'v1.9.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-D100'), 'BIOS Update', 'Latest BIOS firmware update utility', 'exe', '/uploads/downloads/bj-d100-bios.exe', 10485760, 'v1.05', 4);

-- Seed images for BJ-D100
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D100'), '/uploads/products/bj-d100-front.jpg', 'BJ-D100 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-D100'), '/uploads/products/bj-d100-back.jpg', 'BJ-D100 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-D100'), '/uploads/products/bj-d100-side.jpg', 'BJ-D100 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-D100'), '/uploads/products/bj-d100-ports.jpg', 'BJ-D100 Ports Detail', 'Ports Detail', 4, 0);

-- ============================================================================
-- BJ-D300 Compact Desktop (product_id = 3)
-- ============================================================================

-- Seed specs for BJ-D300
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Processor', 'Processor Model', 'Intel Core i5-14500T', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Processor', 'Cores / Threads', '14 Cores / 20 Threads', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Processor', 'Base Clock', '1.7 GHz', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Processor', 'Turbo Clock', '4.8 GHz', 4),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Memory', 'RAM', '16GB DDR5', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Memory', 'Memory Speed', '5600MHz', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Memory', 'Max Memory', '64GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Storage', 'Primary Storage', '512GB NVMe SSD', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Storage', 'Storage Interface', 'M.2 PCIe Gen4', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Connectivity', 'Ethernet', 'Intel I226-V 2.5GbE', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 6E (AX211)', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.3', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Ports', 'USB Ports', '4x USB 3.2 Gen2, 2x USB-C', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Ports', 'Display Outputs', '2x HDMI 2.1, 1x DisplayPort 1.4', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Ports', 'Audio', '3.5mm Audio Jack', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Physical', 'Form Factor', 'Ultra-Compact Mini PC', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Physical', 'Dimensions', '182 x 179 x 34.5 mm', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Physical', 'Weight', '1.3 kg', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Physical', 'Mounting', 'VESA Mountable', 4),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Power', 'Power Supply', '65W External Adapter', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Power', 'Operating System', 'Windows 11 Pro', 2);

-- Seed features for BJ-D300
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D300'), '14th Gen Intel Core i5 Processor (14 Cores)', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), '16GB DDR5 5600MHz RAM (Expandable to 64GB)', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), '512GB PCIe Gen4 NVMe SSD Storage', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Ultra-Compact Form Factor (182x179x34.5mm)', 4),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Wi-Fi 6E and 2.5GbE Connectivity', 5),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'VESA Mountable Design', 6),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Windows 11 Pro Pre-installed', 7),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Perfect for Modern Workspaces', 8);

-- Seed certifications for BJ-D300
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'BIS', 'R-41234503', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'RoHS', NULL, 'badge', 5);

-- Seed downloads for BJ-D300
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'BJ-D300 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-d300-datasheet.pdf', 2621440, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-d300-manual.pdf', 6291456, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'Driver Pack - Windows 11', 'All drivers for Windows 11 (64-bit)', 'zip', '/uploads/downloads/bj-d300-drivers.zip', 156237824, 'v2.3.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-D300'), 'BIOS Update', 'Latest BIOS firmware update utility', 'exe', '/uploads/downloads/bj-d300-bios.exe', 11534336, 'v1.12', 4);

-- Seed images for BJ-D300
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D300'), '/uploads/products/bj-d300-front.jpg', 'BJ-D300 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-D300'), '/uploads/products/bj-d300-back.jpg', 'BJ-D300 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-D300'), '/uploads/products/bj-d300-side.jpg', 'BJ-D300 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-D300'), '/uploads/products/bj-d300-ports.jpg', 'BJ-D300 Ports Detail', 'Ports Detail', 4, 0);

-- ============================================================================
-- BJ-D500 Premium Workstation (product_id = 4)
-- ============================================================================

-- Seed specs for BJ-D500
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Processor', 'Processor Model', 'Intel Core i7-13700K', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Processor', 'Cores / Threads', '16 Cores / 24 Threads', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Processor', 'Base Clock', '3.4 GHz', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Processor', 'Turbo Clock', '5.4 GHz', 4),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Memory', 'RAM', '32GB DDR5', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Memory', 'Memory Speed', '5600MHz', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Memory', 'Max Memory', '192GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Storage', 'Primary Storage', '1TB NVMe SSD', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Storage', 'Storage Interface', 'M.2 PCIe Gen4', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Connectivity', 'Graphics', 'NVIDIA RTX A2000 12GB', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Connectivity', 'Ethernet', 'Intel I226-V 2.5GbE', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 6E (AX211)', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.3', 4),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Ports', 'USB Ports', '6x USB 3.2 Gen2, 2x USB-C Thunderbolt 4', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Ports', 'Display Outputs', '4x DisplayPort (via GPU)', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Ports', 'Audio', 'Premium Audio with S/PDIF', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Physical', 'Form Factor', 'Tower Workstation', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Physical', 'Dimensions', '414 x 371 x 180 mm', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Physical', 'Weight', '12.5 kg', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Power', 'Power Supply', '550W 80+ Gold', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Power', 'Security', 'TPM 2.0, Intel vPro', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Power', 'Operating System', 'Windows 11 Pro for Workstations', 3);

-- Seed features for BJ-D500
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D500'), '13th Gen Intel Core i7-13700K (16 Cores / 24 Threads)', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), '32GB DDR5 5600MHz RAM (Expandable to 192GB)', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), '1TB PCIe Gen4 NVMe SSD Storage', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'NVIDIA RTX A2000 12GB Professional Graphics', 4),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Wi-Fi 6E and 2.5GbE Connectivity', 5),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'TPM 2.0 and Intel vPro Security', 6),
((SELECT id FROM products WHERE sku = 'BJ-D500'), '550W 80+ Gold PSU for Reliability', 7),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Windows 11 Pro for Workstations', 8);

-- Seed certifications for BJ-D500
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'BIS', 'R-41234504', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'RoHS', NULL, 'badge', 5);

-- Seed downloads for BJ-D500
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'BJ-D500 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-d500-datasheet.pdf', 3145728, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-d500-manual.pdf', 9437184, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'Driver Pack - Windows 11', 'All drivers for Windows 11 (64-bit)', 'zip', '/uploads/downloads/bj-d500-drivers.zip', 209715200, 'v3.0.1', 3),
((SELECT id FROM products WHERE sku = 'BJ-D500'), 'BIOS Update', 'Latest BIOS firmware update utility', 'exe', '/uploads/downloads/bj-d500-bios.exe', 14680064, 'v1.15', 4);

-- Seed images for BJ-D500
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-D500'), '/uploads/products/bj-d500-front.jpg', 'BJ-D500 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-D500'), '/uploads/products/bj-d500-back.jpg', 'BJ-D500 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-D500'), '/uploads/products/bj-d500-side.jpg', 'BJ-D500 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-D500'), '/uploads/products/bj-d500-ports.jpg', 'BJ-D500 Ports Detail', 'Ports Detail', 4, 0);
-- Additional seed data for OPS Modules and Interactive Flat Panels
-- Product IDs: 5=BJ-OPS100, 6=BJ-OPS200, 7=BJ-IFP65, 8=BJ-IFP75, 9=BJ-IFP86

-- ============================================================================
-- BJ-OPS100 Standard OPS Module (product_id = 5)
-- ============================================================================

-- Specs for BJ-OPS100
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Processor', 'Processor Model', 'Intel Core i5-10400', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Processor', 'Cores / Threads', '6 Cores / 12 Threads', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Processor', 'Base Clock', '2.9 GHz', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Processor', 'Turbo Clock', '4.3 GHz', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Memory', 'RAM', '8GB DDR4 SO-DIMM', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Memory', 'Memory Speed', '3200MHz', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Memory', 'Max Memory', '16GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Storage', 'Primary Storage', '256GB M.2 SATA SSD', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Storage', 'Storage Interface', 'M.2 SATA', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Connectivity', 'OPS Interface', '80-pin JAE connector', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Connectivity', 'Video Output', 'HDMI + DisplayPort', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Connectivity', 'USB Ports', '5x USB 2.0 + 2x USB 3.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Connectivity', 'Ethernet', 'RJ45 Gigabit LAN', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 6 AX201', 5),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.0', 6),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'System', 'Operating System', 'Windows 10 Pro', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'System', 'Form Factor', 'OPS Standard', 2);

-- Features for BJ-OPS100
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Intel Core i5-10400 Processor with 6 Cores', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), '8GB DDR4 RAM (Expandable to 16GB)', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), '256GB M.2 SATA SSD Storage', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'OPS Standard 80-pin JAE Connector', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Dual Video Output (HDMI + DisplayPort)', 5),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Wi-Fi 6 and Bluetooth 5.0', 6),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Windows 10 Pro Pre-installed', 7),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Plug-and-Play with OPS-Compatible Displays', 8);

-- Certifications for BJ-OPS100
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'BIS', 'R-41234501', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'RoHS', NULL, 'badge', 5);

-- Downloads for BJ-OPS100
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'BJ-OPS100 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-ops100-datasheet.pdf', 2048000, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-ops100-manual.pdf', 7340032, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), 'Driver Pack - Windows 10/11', 'All drivers for Windows 10/11 (64-bit)', 'zip', '/uploads/downloads/bj-ops100-drivers.zip', 145678336, 'v1.5.0', 3);

-- Images for BJ-OPS100
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), '/uploads/products/bj-ops100-front.jpg', 'BJ-OPS100 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), '/uploads/products/bj-ops100-back.jpg', 'BJ-OPS100 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), '/uploads/products/bj-ops100-side.jpg', 'BJ-OPS100 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-OPS100'), '/uploads/products/bj-ops100-ports.jpg', 'BJ-OPS100 Ports Detail', 'Ports Detail', 4, 0);

-- ============================================================================
-- BJ-OPS200 Pro OPS Module (product_id = 6)
-- ============================================================================

-- Specs for BJ-OPS200
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Processor', 'Processor Model', 'Intel Core i7-12650H', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Processor', 'Cores / Threads', '10 Cores / 16 Threads', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Processor', 'Base Clock', '2.3 GHz', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Processor', 'Turbo Clock', '4.7 GHz', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Memory', 'RAM', '16GB DDR4 SO-DIMM', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Memory', 'Memory Speed', '3200MHz', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Memory', 'Max Memory', '32GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Storage', 'Primary Storage', '512GB M.2 PCIe NVMe SSD', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Storage', 'Storage Interface', 'M.2 PCIe NVMe', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Connectivity', 'OPS Interface', '80-pin JAE connector', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Connectivity', 'Video Output', 'HDMI + DisplayPort + USB-C', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Connectivity', 'USB Ports', 'USB 3.0 + USB-C', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Connectivity', 'Ethernet', 'RJ45 Gigabit LAN', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 6E AX211', 5),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.2', 6),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Security', 'TPM', 'TPM 2.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'System', 'Operating System', 'Windows 11 Pro', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'System', 'Form Factor', 'OPS Standard', 2);

-- Features for BJ-OPS200
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), '12th Gen Intel Core i7-12650H Processor with 10 Cores', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), '16GB DDR4 RAM (Expandable to 32GB)', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), '512GB M.2 PCIe NVMe SSD Storage', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Triple Video Output (HDMI + DP + USB-C)', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'TPM 2.0 Hardware Security', 5),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Wi-Fi 6E and Bluetooth 5.2', 6),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Windows 11 Pro Pre-installed', 7),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Enterprise-grade Performance for Professional Use', 8);

-- Certifications for BJ-OPS200
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'BIS', 'R-41234502', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'RoHS', NULL, 'badge', 5);

-- Downloads for BJ-OPS200
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'BJ-OPS200 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-ops200-datasheet.pdf', 2304000, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-ops200-manual.pdf', 8192000, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Driver Pack - Windows 11', 'All drivers for Windows 11 (64-bit)', 'zip', '/uploads/downloads/bj-ops200-drivers.zip', 158654464, 'v2.0.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), 'Firmware Update', 'Latest firmware update utility', 'exe', '/uploads/downloads/bj-ops200-firmware.exe', 15728640, 'v1.12', 4);

-- Images for BJ-OPS200
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), '/uploads/products/bj-ops200-front.jpg', 'BJ-OPS200 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), '/uploads/products/bj-ops200-back.jpg', 'BJ-OPS200 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), '/uploads/products/bj-ops200-side.jpg', 'BJ-OPS200 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-OPS200'), '/uploads/products/bj-ops200-ports.jpg', 'BJ-OPS200 Ports Detail', 'Ports Detail', 4, 0);

-- ============================================================================
-- BJ-IFP65 65" Interactive Flat Panel (product_id = 7)
-- ============================================================================

-- Specs for BJ-IFP65
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Display', 'Screen Size', '65 inches', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Display', 'Resolution', '4K 3840x2160', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Display', 'Brightness', '350 nits', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Display', 'Contrast Ratio', '1200:1', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Display', 'Response Time', '8ms', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Display', 'Viewing Angle', '178° (H) / 178° (V)', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Touch', 'Touch Technology', '20-point Infrared Touch', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Audio', 'Built-in Speakers', '2x 20W Stereo Speakers', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Audio', 'Microphone Array', '8-microphone Array', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Physical', 'Dimensions (WxHxD)', '1488 x 938 x 118 mm', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Physical', 'Weight', '40.6 kg', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'HDMI', '2x HDMI 4K@60Hz', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'DisplayPort', 'DisplayPort', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'VGA', 'VGA Input', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'USB-C', 'USB-C with 65W Power Delivery', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'USB Ports', '4x USB 3.0 + 2x USB 2.0', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'Ethernet', 'RJ45 Gigabit LAN', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'RS232', 'RS232 Control Port', 7),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'OPS Slot', 'OPS Module Slot', 8),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 6', 9),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.2', 10),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'System', 'Operating System', 'Android 13', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'System', 'Processor', 'ARM Cortex-A73 Quad-core', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'System', 'RAM', '4GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'System', 'Storage', '32GB', 4);

-- Features for BJ-IFP65
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), '65-inch 4K Ultra HD Display (3840x2160)', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), '20-point Infrared Touch Technology', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Built-in Android 13 System with 4GB RAM', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'OPS Module Slot for Windows Integration', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'USB-C with 65W Power Delivery', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), '8-microphone Array with Noise Cancellation', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Dual 20W Stereo Speakers', 7),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Wi-Fi 6 and Bluetooth 5.2 Connectivity', 8);

-- Certifications for BJ-IFP65
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'BIS', 'R-41234565', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'RoHS', NULL, 'badge', 5);

-- Downloads for BJ-IFP65
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'BJ-IFP65 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-ifp65-datasheet.pdf', 3145728, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-ifp65-manual.pdf', 12582912, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), 'Android Firmware Update', 'Latest Android firmware update', 'zip', '/uploads/downloads/bj-ifp65-firmware.zip', 524288000, 'v3.2.1', 3);

-- Images for BJ-IFP65
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), '/uploads/products/bj-ifp65-front.jpg', 'BJ-IFP65 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), '/uploads/products/bj-ifp65-back.jpg', 'BJ-IFP65 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), '/uploads/products/bj-ifp65-side.jpg', 'BJ-IFP65 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-IFP65'), '/uploads/products/bj-ifp65-ports.jpg', 'BJ-IFP65 Ports Detail', 'Ports Detail', 4, 0);

-- ============================================================================
-- BJ-IFP75 75" Interactive Flat Panel (product_id = 8)
-- ============================================================================

-- Specs for BJ-IFP75
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Display', 'Screen Size', '75 inches', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Display', 'Resolution', '4K 3840x2160', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Display', 'Brightness', '450 nits', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Display', 'Contrast Ratio', '1200:1', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Display', 'Response Time', '8ms', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Display', 'Viewing Angle', '178° (H) / 178° (V)', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Touch', 'Touch Technology', '40-point Infrared Touch', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Audio', 'Built-in Speakers', '2.1ch (2x 20W + Subwoofer)', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Audio', 'Microphone Array', '8-microphone Array', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Physical', 'Dimensions (WxHxD)', '1709 x 1061 x 118 mm', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Physical', 'Weight', '52.2 kg', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'HDMI', '4x HDMI 4K@60Hz', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'DisplayPort', 'DisplayPort', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'USB-C', '2x USB-C with 65W Power Delivery', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'USB Ports', '4x USB 3.0', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'Ethernet', 'RJ45 Gigabit LAN', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'RS232', 'RS232 Control Port', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'OPS Slot', 'OPS Module Slot', 7),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 6', 8),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.2', 9),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'System', 'Operating System', 'Android 13', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'System', 'Processor', 'ARM Cortex-A73 (4x) + A53 (4x) Octa-core', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'System', 'RAM', '8GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'System', 'Storage', '128GB', 4);

-- Features for BJ-IFP75
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '75-inch 4K Ultra HD Display (3840x2160)', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '40-point Infrared Touch Technology', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Built-in Android 13 System with 8GB RAM', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'OPS Module Slot for Windows Integration', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Dual USB-C with 65W Power Delivery', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '2.1 Channel Audio System with Subwoofer', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '8-microphone Array with AI Noise Reduction', 7),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '4x HDMI Inputs for Multiple Sources', 8);

-- Certifications for BJ-IFP75
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'BIS', 'R-41234575', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'RoHS', NULL, 'badge', 5);

-- Downloads for BJ-IFP75
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'BJ-IFP75 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-ifp75-datasheet.pdf', 3407872, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-ifp75-manual.pdf', 13631488, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), 'Android Firmware Update', 'Latest Android firmware update', 'zip', '/uploads/downloads/bj-ifp75-firmware.zip', 536870912, 'v3.2.1', 3);

-- Images for BJ-IFP75
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '/uploads/products/bj-ifp75-front.jpg', 'BJ-IFP75 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '/uploads/products/bj-ifp75-back.jpg', 'BJ-IFP75 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '/uploads/products/bj-ifp75-side.jpg', 'BJ-IFP75 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-IFP75'), '/uploads/products/bj-ifp75-ports.jpg', 'BJ-IFP75 Ports Detail', 'Ports Detail', 4, 0);

-- ============================================================================
-- BJ-IFP86 86" Interactive Flat Panel (product_id = 9)
-- ============================================================================

-- Specs for BJ-IFP86
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Display', 'Screen Size', '86 inches', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Display', 'Resolution', '4K 3840x2160', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Display', 'Brightness', '450 nits', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Display', 'Contrast Ratio', '1200:1', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Display', 'Response Time', '8ms', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Display', 'Viewing Angle', '178° (H) / 178° (V)', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Touch', 'Touch Technology', '50-point Infrared Touch', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Audio', 'Built-in Speakers', '2.1ch (2x 20W + Subwoofer)', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Audio', 'Microphone Array', '8-microphone Array', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Physical', 'Dimensions (WxHxD)', '1957 x 1201 x 118 mm', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Physical', 'Weight', '68.7 kg', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'HDMI', '4x HDMI 4K@60Hz', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'DisplayPort', 'DisplayPort', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'USB-C', '2x USB-C with 65W Power Delivery', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'USB Ports', '4x USB 3.0', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'Ethernet', 'RJ45 Gigabit LAN', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'RS232', 'RS232 Control Port', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'OPS Slot', 'OPS Module Slot', 7),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'Wi-Fi', 'Wi-Fi 6', 8),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Connectivity', 'Bluetooth', 'Bluetooth 5.2', 9),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'System', 'Operating System', 'Android 13', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'System', 'Processor', 'ARM Cortex-A73 (4x) + A53 (4x) Octa-core', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'System', 'RAM', '8GB', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'System', 'Storage', '128GB', 4);

-- Features for BJ-IFP86
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), '86-inch 4K Ultra HD Display (3840x2160)', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), '50-point Infrared Touch Technology', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Built-in Android 13 System with 8GB RAM', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'OPS Module Slot for Windows Integration', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Dual USB-C with 65W Power Delivery', 5),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Premium 2.1 Channel Audio System', 6),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Advanced 8-microphone Array with AI Processing', 7),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Ideal for Large Venues and Auditoriums', 8);

-- Certifications for BJ-IFP86
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'BIS', 'R-41234586', 'badge', 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'CE', NULL, 'badge', 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'FCC', NULL, 'badge', 3),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Energy Star', NULL, 'icon', 4),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'RoHS', NULL, 'badge', 5);

-- Downloads for BJ-IFP86
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'BJ-IFP86 Datasheet', 'Product specifications and features overview', 'pdf', '/uploads/downloads/bj-ifp86-datasheet.pdf', 3670016, NULL, 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'User Manual', 'Installation and operation guide', 'pdf', '/uploads/downloads/bj-ifp86-manual.pdf', 14680064, NULL, 2),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), 'Android Firmware Update', 'Latest Android firmware update', 'zip', '/uploads/downloads/bj-ifp86-firmware.zip', 536870912, 'v3.2.1', 3);

-- Images for BJ-IFP86
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), '/uploads/products/bj-ifp86-front.jpg', 'BJ-IFP86 Front View', 'Front View', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), '/uploads/products/bj-ifp86-back.jpg', 'BJ-IFP86 Back View', 'Back View', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), '/uploads/products/bj-ifp86-side.jpg', 'BJ-IFP86 Side View', 'Side View', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-IFP86'), '/uploads/products/bj-ifp86-ports.jpg', 'BJ-IFP86 Ports Detail', 'Ports Detail', 4, 0);
-- Product Sub-Entities for AV & IoT Products
-- Products: BJ-WPS100, BJ-CAM360, BJ-SPK200, BJ-DS55, BJ-IOT100

-- ============================================================================
-- BJ-WPS100 Wireless Presenter (Product ID: 10)
-- ============================================================================

-- Product Specs
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Video Output', 'Output Interface', 'HDMI 2.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Video Output', 'Resolution', '4K @ 60fps', 2),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Video Output', 'Split Screen', 'Up to 4 presenters simultaneously', 3),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Connectivity', 'Wireless', 'Wi-Fi 6 (802.11ax)', 4),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Connectivity', 'Range', '30 meters', 5),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Connectivity', 'Encryption', 'AES 128-bit', 6),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Hardware', 'Buttons', 'HDMI + USB-C connection buttons', 7),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Hardware', 'Package Contents', 'Host unit, 2 presenter buttons, charging cradle', 8),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Compatibility', 'Operating Systems', 'Windows, macOS, iOS, Android', 9),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Compatibility', 'Software Required', 'None - plug and play', 10);

-- Product Features
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'HDMI 2.0 output supports crystal-clear 4K resolution at 60fps for stunning presentations', 1),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Wi-Fi 6 technology ensures fast, stable wireless connection with minimal latency', 2),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Split-screen capability allows up to 4 presenters to share content simultaneously', 3),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'AES 128-bit encryption protects your presentation content from unauthorized access', 4),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), '30-meter wireless range provides freedom to move around large conference rooms', 5),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Simple HDMI and USB-C buttons make connection quick and intuitive', 6),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Zero software installation required - true plug-and-play operation', 7),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Cross-platform compatibility with Windows, Mac, iOS, and Android devices', 8);

-- Product Certifications
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'CE Marking', 'CE', 'ce', 1),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'FCC Certified', 'FCC ID: 2A3BJ-WPS100', 'fcc', 2),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'RoHS Compliant', 'RoHS 2011/65/EU', 'rohs', 3),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'HDMI 2.0 Certified', 'HDMI 2.0b', 'hdmi', 4),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Wi-Fi Alliance', 'Wi-Fi 6 Certified', 'wifi', 5);

-- Product Downloads
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Product Datasheet', 'Technical specifications and features overview', 'pdf', '/uploads/downloads/bj-wps100-datasheet.pdf', 2202009, '1.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'User Manual', 'Complete installation and operation guide', 'pdf', '/uploads/downloads/bj-wps100-manual.pdf', 4718592, '1.0', 2),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Quick Start Guide', 'Fast setup instructions', 'pdf', '/uploads/downloads/bj-wps100-quickstart.pdf', 1258291, '1.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), 'Firmware Update', 'Latest firmware version for host unit', 'zip', '/uploads/downloads/bj-wps100-firmware.zip', 16042189, '1.2.0', 4);

-- Product Images
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), '/uploads/products/bj-wps100-front.jpg', 'BJ-WPS100 Wireless Presenter Front View', 'Front view showing host unit and presenter buttons', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), '/uploads/products/bj-wps100-setup.jpg', 'BJ-WPS100 Complete Setup', 'Complete package with host, buttons, and charging cradle', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), '/uploads/products/bj-wps100-button.jpg', 'BJ-WPS100 Presenter Button Detail', 'Close-up of presenter button with HDMI and USB-C options', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), '/uploads/products/bj-wps100-split.jpg', 'BJ-WPS100 Split Screen Demo', '4-way split screen presentation in action', 4, 0),
((SELECT id FROM products WHERE sku = 'BJ-WPS100'), '/uploads/products/bj-wps100-lifestyle.jpg', 'BJ-WPS100 In Use', 'Presenter using wireless button in conference room', 5, 0);

-- ============================================================================
-- BJ-CAM360 360° Conference Camera (Product ID: 11)
-- ============================================================================

-- Product Specs
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Camera', 'Field of View', '360 degrees', 1),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Camera', 'Resolution', '1080p Full HD', 2),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Camera', 'Sensors', '3x 13MP cameras with AI stitching', 3),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Camera', 'Digital Zoom', '6x zoom', 4),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Audio', 'Microphones', '8 beamforming microphones', 5),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Audio', 'Pickup Range', '18 feet (5.5 meters)', 6),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Audio', 'Speaker', '360° tri-speaker system, 76dB output', 7),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Connectivity', 'Primary Interface', 'USB-C', 8),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Connectivity', 'Wireless', 'Wi-Fi for firmware updates', 9),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Hardware', 'Processor', 'Qualcomm Snapdragon 605', 10),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Hardware', 'Dimensions', '111 x 111 x 272 mm', 11),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Hardware', 'Weight', '1.2 kg', 12),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Compatibility', 'Video Platforms', 'Microsoft Teams, Zoom, Google Meet certified', 13);

-- Product Features
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '360-degree field of view captures everyone in the room without blind spots', 1),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Triple 13MP camera array with AI-powered stitching for seamless panoramic video', 2),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '8 beamforming microphones with 18-foot pickup range ensure crystal-clear audio from all directions', 3),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '360-degree tri-speaker system delivers room-filling 76dB audio output', 4),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '6x digital zoom allows focus on active speakers with smooth transitions', 5),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Qualcomm Snapdragon 605 processor powers intelligent AI features and real-time processing', 6),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'USB-C connectivity provides simple plug-and-play setup with any conferencing system', 7),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Certified compatibility with Microsoft Teams, Zoom, and Google Meet', 8);

-- Product Certifications
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'CE Marking', 'CE', 'ce', 1),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'FCC Certified', 'FCC ID: 2A3BJ-CAM360', 'fcc', 2),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'RoHS Compliant', 'RoHS 2011/65/EU', 'rohs', 3),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Microsoft Teams Certified', 'Teams Rooms Certified', 'teams', 4),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Zoom Certified', 'Zoom Rooms Certified', 'zoom', 5),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Google Meet Certified', 'Meet Hardware Certified', 'google', 6);

-- Product Downloads
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Product Datasheet', 'Technical specifications and features overview', 'pdf', '/uploads/downloads/bj-cam360-datasheet.pdf', 3355443, '1.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'User Manual', 'Complete installation, setup, and operation guide', 'pdf', '/uploads/downloads/bj-cam360-manual.pdf', 7130317, '1.0', 2),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Quick Installation Guide', 'Fast setup instructions for conference rooms', 'pdf', '/uploads/downloads/bj-cam360-quickstart.pdf', 1572864, '1.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Firmware Update', 'Latest firmware with enhanced AI features', 'zip', '/uploads/downloads/bj-cam360-firmware.zip', 134742016, '2.1.3', 4),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), 'Integration Guide', 'Platform-specific setup for Teams, Zoom, and Meet', 'pdf', '/uploads/downloads/bj-cam360-integration.pdf', 2831155, '1.0', 5);

-- Product Images
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '/uploads/products/bj-cam360-front.jpg', 'BJ-CAM360 360° Camera Front View', 'Front view showing camera array and speaker grille', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '/uploads/products/bj-cam360-top.jpg', 'BJ-CAM360 Top View', 'Top view showing 360° camera configuration', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '/uploads/products/bj-cam360-room.jpg', 'BJ-CAM360 In Conference Room', 'Camera positioned in center of conference table', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '/uploads/products/bj-cam360-detail.jpg', 'BJ-CAM360 Camera Detail', 'Close-up of triple camera array and microphones', 4, 0),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '/uploads/products/bj-cam360-interface.jpg', 'BJ-CAM360 Connection Panel', 'USB-C port and control interface', 5, 0),
((SELECT id FROM products WHERE sku = 'BJ-CAM360'), '/uploads/products/bj-cam360-action.jpg', 'BJ-CAM360 During Video Call', '360° view during active video conference', 6, 0);

-- ============================================================================
-- BJ-SPK200 Smart Speakerphone (Product ID: 12)
-- ============================================================================

-- Product Specs
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Audio Output', 'Driver', '65mm full-range driver', 1),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Audio Output', 'Sound Pressure Level', '76dB SPL', 2),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Audio Output', 'Room Coverage', '16 x 16 feet', 3),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Audio Input', 'Microphones', '4 noise-cancelling beamforming microphones', 4),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Audio Input', 'Technology', 'Full duplex, super wideband', 5),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Connectivity', 'Wired', 'USB-C + USB-A', 6),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Connectivity', 'Wireless', 'Bluetooth 5.2, 30m range', 7),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Power', 'Battery Life', '32 hours continuous use', 8),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Durability', 'IP Rating', 'IP64 (dust and water resistant)', 9),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Physical', 'Weight', '0.61 kg', 10),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Compatibility', 'Certifications', 'Microsoft Teams, Zoom, Google Meet certified', 11);

-- Product Features
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '65mm full-range driver delivers powerful 76dB output for clear audio across 16x16 foot rooms', 1),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '4 noise-cancelling beamforming microphones capture every voice with precision', 2),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Full duplex technology allows natural, simultaneous two-way conversations', 3),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Dual connectivity with USB-C, USB-A wired and Bluetooth 5.2 wireless options', 4),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Extended 32-hour battery life ensures all-day meetings without recharging', 5),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'IP64 rating provides dust and water resistance for demanding environments', 6),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '30-meter Bluetooth range offers flexibility for various room configurations', 7),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Certified for Microsoft Teams, Zoom, and Google Meet with optimized performance', 8);

-- Product Certifications
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'CE Marking', 'CE', 'ce', 1),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'FCC Certified', 'FCC ID: 2A3BJ-SPK200', 'fcc', 2),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'RoHS Compliant', 'RoHS 2011/65/EU', 'rohs', 3),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'IP64 Rated', 'IP64', 'ip', 4),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Microsoft Teams Certified', 'Teams Rooms Certified', 'teams', 5),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Zoom Certified', 'Zoom Certified Peripheral', 'zoom', 6),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Google Meet Certified', 'Meet Hardware Certified', 'google', 7);

-- Product Downloads
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Product Datasheet', 'Technical specifications and features overview', 'pdf', '/uploads/downloads/bj-spk200-datasheet.pdf', 1887436, '1.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'User Manual', 'Complete operation and maintenance guide', 'pdf', '/uploads/downloads/bj-spk200-manual.pdf', 4089446, '1.0', 2),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Quick Start Guide', 'Fast setup for immediate use', 'pdf', '/uploads/downloads/bj-spk200-quickstart.pdf', 943718, '1.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), 'Firmware Update', 'Latest firmware for enhanced audio performance', 'zip', '/uploads/downloads/bj-spk200-firmware.zip', 8598323, '1.4.2', 4);

-- Product Images
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '/uploads/products/bj-spk200-front.jpg', 'BJ-SPK200 Speakerphone Front View', 'Front view showing speaker grille and control buttons', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '/uploads/products/bj-spk200-top.jpg', 'BJ-SPK200 Top View', 'Top view showing microphone array pattern', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '/uploads/products/bj-spk200-ports.jpg', 'BJ-SPK200 Connection Ports', 'USB-C and USB-A ports detail', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '/uploads/products/bj-spk200-desk.jpg', 'BJ-SPK200 On Office Desk', 'Speakerphone in typical office environment', 4, 0),
((SELECT id FROM products WHERE sku = 'BJ-SPK200'), '/uploads/products/bj-spk200-lifestyle.jpg', 'BJ-SPK200 In Use', 'Small team using speakerphone during conference call', 5, 0);

-- ============================================================================
-- BJ-DS55 Digital Signage Display (Product ID: 13)
-- ============================================================================

-- Product Specs
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Display', 'Size', '55 inches', 1),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Display', 'Resolution', '4K UHD (3840 x 2160)', 2),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Display', 'Panel Type', 'IPS', 3),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Display', 'Brightness', '500 nits', 4),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Display', 'Contrast Ratio', '4000:1', 5),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Display', 'Color Depth', '10-bit color', 6),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Display', 'Operation', '24/7 continuous operation', 7),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Connectivity', 'HDMI Input', '3x HDMI 2.0', 8),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Connectivity', 'HDMI Output', '1x HDMI (loop-through)', 9),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Connectivity', 'Other Inputs', 'DisplayPort, USB, RJ45, RS232', 10),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Connectivity', 'Wireless', 'Wi-Fi 802.11ac', 11),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'System', 'Operating System', 'webOS 6.0', 12),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Physical', 'Dimensions', '1231 x 707 x 30 mm', 13),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Physical', 'Weight', '16.1 kg', 14),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Physical', 'VESA Mount', '300 x 300 mm', 15),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Physical', 'Orientation', 'Landscape or Portrait', 16),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Durability', 'IP Rating', 'IP5X (dust protection)', 17);

-- Product Features
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '55-inch 4K UHD IPS display delivers stunning 3840x2160 resolution with vibrant colors', 1),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '500-nit brightness and 4000:1 contrast ratio ensure excellent visibility in any lighting', 2),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '10-bit color depth displays over 1 billion colors for accurate, lifelike images', 3),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Commercial-grade panel designed for reliable 24/7 continuous operation', 4),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'webOS 6.0 operating system provides powerful content management and scheduling', 5),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Extensive connectivity with 3x HDMI inputs, DisplayPort, USB, RJ45, and RS232', 6),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Flexible installation supports both landscape and portrait orientations', 7),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'VESA 300x300mm mount compatible with standard commercial mounting solutions', 8);

-- Product Certifications
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'CE Marking', 'CE', 'ce', 1),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'FCC Certified', 'FCC ID: 2A3BJ-DS55', 'fcc', 2),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'RoHS Compliant', 'RoHS 2011/65/EU', 'rohs', 3),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'IP5X Rated', 'IP5X Dust Protection', 'ip', 4),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Energy Star', 'Energy Star 8.0', 'energy', 5),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'UL Listed', 'UL 62368-1', 'ul', 6);

-- Product Downloads
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Product Datasheet', 'Technical specifications and features overview', 'pdf', '/uploads/downloads/bj-ds55-datasheet.pdf', 4299161, '1.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'User Manual', 'Complete installation and operation guide', 'pdf', '/uploads/downloads/bj-ds55-manual.pdf', 9122611, '1.0', 2),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Installation Guide', 'Mounting and setup instructions', 'pdf', '/uploads/downloads/bj-ds55-installation.pdf', 3355443, '1.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'webOS Content Manager', 'Software for content management and scheduling', 'zip', '/uploads/downloads/bj-ds55-cms.zip', 164403814, '6.0.2', 4),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'Firmware Update', 'Latest display firmware', 'zip', '/uploads/downloads/bj-ds55-firmware.zip', 257214054, '3.2.1', 5),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), 'CAD Drawings', 'Technical drawings for installation planning', 'DWG', '/uploads/downloads/bj-ds55-cad.dwg', 2936012, '1.0', 6);

-- Product Images
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '/uploads/products/bj-ds55-front.jpg', 'BJ-DS55 Digital Signage Display Front View', 'Front view showing ultra-slim bezel design', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '/uploads/products/bj-ds55-back.jpg', 'BJ-DS55 Rear Panel', 'Rear view showing connectivity ports and VESA mount', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '/uploads/products/bj-ds55-retail.jpg', 'BJ-DS55 In Retail Environment', 'Display showing promotional content in store', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '/uploads/products/bj-ds55-portrait.jpg', 'BJ-DS55 Portrait Orientation', 'Display mounted vertically in portrait mode', 4, 0),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '/uploads/products/bj-ds55-corporate.jpg', 'BJ-DS55 In Corporate Lobby', 'Display showing corporate information and wayfinding', 5, 0),
((SELECT id FROM products WHERE sku = 'BJ-DS55'), '/uploads/products/bj-ds55-menu.jpg', 'BJ-DS55 As Menu Board', 'Display used as digital menu board in restaurant', 6, 0);

-- ============================================================================
-- BJ-IOT100 Smart Environmental Sensor (Product ID: 14)
-- ============================================================================

-- Product Specs
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Sensors', 'Temperature', '±1°C accuracy', 1),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Sensors', 'Humidity', '±3% RH accuracy', 2),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Sensors', 'Air Pressure', '300-1100 hPa range', 3),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Sensors', 'Air Quality', 'VOC gas sensor', 4),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Physical', 'Dimensions', '3 x 3 x 1 mm', 5),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Connectivity', 'Interfaces', 'I2C + SPI', 6),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Connectivity', 'Wireless', 'Bluetooth Low Energy + Wi-Fi', 7),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Connectivity', 'Protocols', 'MQTT, HTTP', 8),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Power', 'Operating Voltage', '1.2 - 3.6V', 9),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Power', 'Battery Life', '2 years (typical use)', 10),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Environmental', 'Operating Temperature', '-40°C to +85°C', 11),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Installation', 'Mounting', 'Wall or ceiling mount', 12);

-- Product Features
INSERT INTO product_features (product_id, feature_text, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'High-precision temperature sensor with ±1°C accuracy for reliable environmental monitoring', 1),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Humidity sensor with ±3% RH accuracy tracks moisture levels effectively', 2),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Air pressure sensor covers wide 300-1100 hPa range for comprehensive atmospheric monitoring', 3),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'VOC gas sensor detects air quality issues and volatile organic compounds', 4),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Ultra-compact 3x3x1mm form factor fits virtually any IoT application', 5),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Dual connectivity with BLE and Wi-Fi supports flexible deployment scenarios', 6),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'MQTT and HTTP protocol support enables seamless cloud integration', 7),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Extended 2-year battery life reduces maintenance requirements', 8);

-- Product Certifications
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'CE Marking', 'CE', 'ce', 1),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'FCC Certified', 'FCC ID: 2A3BJ-IOT100', 'fcc', 2),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'RoHS Compliant', 'RoHS 2011/65/EU', 'rohs', 3),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Bluetooth SIG', 'Bluetooth 5.0 Qualified', 'bluetooth', 4),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Wi-Fi Alliance', 'Wi-Fi Certified', 'wifi', 5);

-- Product Downloads
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Product Datasheet', 'Technical specifications and features overview', 'pdf', '/uploads/downloads/bj-iot100-datasheet.pdf', 1572864, '1.0', 1),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Integration Guide', 'API documentation and integration instructions', 'pdf', '/uploads/downloads/bj-iot100-integration.pdf', 4404019, '1.0', 2),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Quick Start Guide', 'Fast deployment instructions', 'pdf', '/uploads/downloads/bj-iot100-quickstart.pdf', 838860, '1.0', 3),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'Firmware Update', 'Latest sensor firmware', 'zip', '/uploads/downloads/bj-iot100-firmware.zip', 3250585, '2.0.5', 4),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'SDK Package', 'Software development kit for custom applications', 'zip', '/uploads/downloads/bj-iot100-sdk.zip', 30094745, '1.3.0', 5),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), 'MQTT Configuration', 'MQTT broker setup and configuration guide', 'pdf', '/uploads/downloads/bj-iot100-mqtt.pdf', 1992294, '1.0', 6);

-- Product Images
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail) VALUES
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), '/uploads/products/bj-iot100-front.jpg', 'BJ-IOT100 Smart Sensor Front View', 'Compact sensor module with visible components', 1, 1),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), '/uploads/products/bj-iot100-size.jpg', 'BJ-IOT100 Size Comparison', 'Sensor shown next to coin for scale reference', 2, 0),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), '/uploads/products/bj-iot100-mounted.jpg', 'BJ-IOT100 Wall Mounted', 'Sensor installed on wall in office environment', 3, 0),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), '/uploads/products/bj-iot100-pcb.jpg', 'BJ-IOT100 PCB Layout', 'Detailed view of sensor circuit board', 4, 0),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), '/uploads/products/bj-iot100-dashboard.jpg', 'BJ-IOT100 Monitoring Dashboard', 'Web dashboard showing real-time sensor data', 5, 0),
((SELECT id FROM products WHERE sku = 'BJ-IOT100'), '/uploads/products/bj-iot100-install.jpg', 'BJ-IOT100 Installation Options', 'Various mounting configurations', 6, 0);
