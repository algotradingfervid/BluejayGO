-- Insert solutions
INSERT INTO solutions (title, slug, icon, short_description, hero_image_url, hero_title, hero_description, overview_content, meta_description, reference_code, is_published, display_order) VALUES
('Education', 'education', 'school', 'Transform classrooms with interactive technology that engages students and empowers educators.', '/uploads/solutions/education.jpg', 'Education Solutions', 'Transforming learning environments with innovative technology solutions', '<p>The education sector is undergoing a digital revolution. Modern classrooms require technology that engages students, facilitates interactive learning, and provides educators with powerful tools to deliver content effectively.</p><p>BlueJay''s education solutions are designed to transform traditional learning spaces into dynamic, interactive environments. From Interactive Flat Panels to robust computing solutions, we provide end-to-end technology that supports 21st-century education.</p>', 'Transform classrooms with interactive technology. BlueJay''s education solutions engage students and empower educators.', 'EDU-2026.01', 1, 1),
('Corporate', 'corporate', 'corporate_fare', 'Empower modern workplaces with collaborative solutions that boost productivity and communication.', '/uploads/solutions/corporate.jpg', 'Corporate Solutions', 'Enhancing workplace collaboration and productivity', '<p>In today''s rapidly evolving business landscape, Indian enterprises face unprecedented challenges in maintaining productivity across hybrid work environments. Studies show that 73% of organizations struggle with inconsistent meeting room experiences and technology fragmentation across office locations. BlueJay''s Corporate Solutions address these critical pain points with integrated technology infrastructure designed specifically for modern Indian businesses.</p><p>Our comprehensive portfolio combines enterprise-grade desktop computing, interactive collaboration displays, and intelligent IoT devices to create seamless hybrid work experiences. From boardrooms to individual workstations, BlueJay delivers reliable, scalable technology that empowers teams to collaborate effectively whether they''re in Mumbai or working remotely from tier-2 cities.</p><p>With over a decade of experience serving India''s corporate sector, we understand the unique requirements of domestic enterprises—from power efficiency and serviceability to integration with existing IT infrastructure. Our solutions are built to withstand India''s diverse operating conditions while delivering the performance and security that modern businesses demand.</p>', 'Corporate technology solutions for modern workplaces.', 'CORP-2026.01', 1, 2),
('Healthcare', 'healthcare', 'medical_services', 'Enhance patient care with reliable technology solutions designed for medical environments.', '/uploads/solutions/healthcare.jpg', 'Healthcare Solutions', 'Technology solutions for modern healthcare', '<p>India''s healthcare sector is undergoing rapid digital transformation, with hospitals and clinics embracing technology to improve patient outcomes and operational efficiency. From telemedicine platforms connecting rural patients with specialist doctors to digital patient engagement systems that streamline appointments and medical records, technology is bridging critical gaps in healthcare delivery across urban and rural India.</p><p>BlueJay''s healthcare technology solutions empower medical institutions with reliable, enterprise-grade hardware designed for demanding 24/7 clinical environments. Our interactive flat panels transform medical education and training, enabling collaborative learning for medical students and continuous professional development for healthcare practitioners. Digital signage solutions improve patient communication and wayfinding in hospital campuses, while IoT-enabled monitoring systems help facilities management teams maintain optimal environments for patient care.</p><p>With a focus on data security, HIPAA compliance readiness, and equipment reliability, BlueJay supports healthcare providers in delivering superior patient experiences while optimizing operational costs. Our solutions are deployed across hospitals, diagnostic centers, telemedicine hubs, and medical training institutions throughout India.</p>', 'Healthcare technology solutions for medical environments.', 'HCR-2026.01', 1, 3),
('Retail', 'retail', 'storefront', 'Create engaging customer experiences with digital signage and interactive displays.', '/uploads/solutions/retail.jpg', 'Retail Solutions', 'Engaging retail experiences through technology', '<p>Transform your retail experience with BlueJay''s comprehensive digital solutions designed specifically for the Indian market. Our integrated ecosystem of interactive displays, digital signage, point-of-sale systems, and IoT sensors enables retailers to create engaging customer experiences while optimizing operations and driving sales growth across physical and digital channels.</p><p>From flagship stores to multi-location chains, BlueJay empowers retailers with tools to deliver personalized shopping experiences, streamline inventory management, and gain real-time insights into customer behavior. Our solutions seamlessly integrate with existing retail systems, supporting omnichannel strategies that bridge online and offline touchpoints.</p><p>Leverage cutting-edge technology including interactive flat panels for customer engagement, robust digital signage for dynamic content delivery, and intelligent IoT sensors for foot traffic analysis and environmental monitoring. With proven deployments across 5000+ retail locations in India, BlueJay delivers the reliability, scalability, and innovation that modern retailers demand.</p>', 'Digital signage and interactive displays for retail.', 'RTL-2026.01', 1, 4),
('Government', 'government', 'account_balance', 'Support public sector digital transformation with GeM-registered, compliant solutions.', '/uploads/solutions/government.jpg', 'Government Solutions', 'Digital transformation for public sector', '<p>BlueJay empowers government departments and public sector organizations across India with cutting-edge technology solutions designed for Digital India initiatives. Our comprehensive portfolio supports e-governance, smart city projects, and citizen service delivery through GeM-compliant procurement processes.</p><p>From secure desktop workstations to interactive collaboration tools for command centers and municipal offices, we deliver Make in India solutions that ensure data sovereignty, meet stringent security standards, and modernize legacy infrastructure. Our technology enables efficient public service delivery while optimizing budget utilization.</p><p>Trusted by 100+ government departments and institutions, BlueJay provides end-to-end support for digital transformation projects, smart classroom initiatives, video conferencing solutions for administrative offices, and IoT-enabled infrastructure monitoring systems that drive transparent, accountable governance.</p>', 'GeM-registered technology solutions for government.', 'GOV-2026.01', 1, 5),
('Hospitality', 'hospitality', 'hotel', 'Deliver exceptional guest experiences with integrated technology and digital solutions.', '/uploads/solutions/hospitality.jpg', 'Hospitality Solutions', 'Elevating guest experiences with technology', '<p>India''s hospitality industry is experiencing unprecedented growth, with travelers demanding world-class experiences and seamless technology integration. From boutique hotels to large resorts and conference facilities, BlueJay''s hospitality solutions empower properties to deliver exceptional guest experiences while optimizing operational efficiency.</p><p>Our comprehensive technology suite transforms every touchpoint of the guest journey. Interactive displays in lobbies and conference rooms create engaging experiences, while smart room automation systems provide personalized comfort control. Digital signage solutions keep guests informed and enhance wayfinding throughout your property, reducing staff workload while improving satisfaction scores.</p><p>BlueJay understands the unique challenges of Indian hospitality operations—from managing diverse guest needs to maintaining consistent service quality across properties. Our solutions are designed for the demands of tropical climates, varying power conditions, and high-traffic environments, ensuring reliable performance that enhances your brand reputation and drives positive reviews.</p>', 'Technology solutions for hospitality industry.', 'HSP-2026.01', 1, 6);

-- ============================================================
-- EDUCATION SOLUTION SUB-ENTITIES
-- ============================================================

-- Education solution stats
INSERT INTO solution_stats (solution_id, value, label, display_order) VALUES
((SELECT id FROM solutions WHERE slug='education'), '500+', 'Schools Served', 1),
((SELECT id FROM solutions WHERE slug='education'), '95%', 'Satisfaction Rate', 2),
((SELECT id FROM solutions WHERE slug='education'), '10,000+', 'Classrooms Transformed', 3);

-- Education solution challenges
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order) VALUES
((SELECT id FROM solutions WHERE slug='education'), 'Student Engagement', 'Keeping students focused and actively participating in lessons is increasingly challenging with traditional teaching methods.', 'target', 1),
((SELECT id FROM solutions WHERE slug='education'), 'Remote & Hybrid Learning', 'Enabling seamless learning experiences across physical and virtual classrooms with integrated technology.', 'video_call', 2),
((SELECT id FROM solutions WHERE slug='education'), 'Content Accessibility', 'Making educational content accessible to all students regardless of learning styles or abilities.', 'accessibility_new', 3),
((SELECT id FROM solutions WHERE slug='education'), 'Technology Integration', 'Ensuring new technology works seamlessly with existing infrastructure and is easy for educators to adopt.', 'extension', 4),
((SELECT id FROM solutions WHERE slug='education'), 'Budget Constraints', 'Delivering powerful technology solutions that meet institutional needs within limited education budgets.', 'savings', 5),
((SELECT id FROM solutions WHERE slug='education'), 'Long-term Reliability', 'Educational institutions need durable products that withstand daily use and provide years of reliable service.', 'verified_user', 6);

-- Education solution products
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured) VALUES
((SELECT id FROM solutions WHERE slug='education'), (SELECT id FROM products WHERE sku='BJ-IFP75'), 1, 1),
((SELECT id FROM solutions WHERE slug='education'), (SELECT id FROM products WHERE sku='BJ-IFP65'), 2, 1),
((SELECT id FROM solutions WHERE slug='education'), (SELECT id FROM products WHERE sku='BJ-OPS100'), 3, 0),
((SELECT id FROM solutions WHERE slug='education'), (SELECT id FROM products WHERE sku='BJ-D100'), 4, 0),
((SELECT id FROM solutions WHERE slug='education'), (SELECT id FROM products WHERE sku='BJ-WPS100'), 5, 0);

-- Education solution CTA
INSERT INTO solution_ctas (solution_id, heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name) VALUES
((SELECT id FROM solutions WHERE slug='education'), 'Ready to Transform Your Educational Institution?', 'Our education specialists can help you design the perfect technology solution for your school, college, or university.', 'Schedule a Demo', '/contact', 'Download Brochure', '/downloads/education-brochure', '+91-120-456-7890', 'main_cta');

-- ============================================================
-- CORPORATE SOLUTION SUB-ENTITIES
-- ============================================================

-- Corporate solution stats
INSERT INTO solution_stats (solution_id, value, label, display_order) VALUES
((SELECT id FROM solutions WHERE slug='corporate'), '1000+', 'Enterprises Served', 1),
((SELECT id FROM solutions WHERE slug='corporate'), '15000+', 'Meeting Rooms Equipped', 2),
((SELECT id FROM solutions WHERE slug='corporate'), '98.7%', 'Uptime Guarantee', 3);

-- Corporate solution challenges
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order) VALUES
((SELECT id FROM solutions WHERE slug='corporate'), 'Hybrid Work Complexity', 'Managing seamless collaboration between remote and in-office teams with inconsistent technology experiences across locations.', 'groups', 1),
((SELECT id FROM solutions WHERE slug='corporate'), 'Meeting Room Inefficiency', 'Outdated conference room technology leading to wasted time, poor video quality, and frustrated employees during critical client presentations.', 'video_call', 2),
((SELECT id FROM solutions WHERE slug='corporate'), 'IT Management Overhead', 'Fragmented vendor ecosystems and lack of centralized device management increasing IT support costs and complexity.', 'settings', 3),
((SELECT id FROM solutions WHERE slug='corporate'), 'Security & Compliance', 'Protecting sensitive corporate data across distributed workforces while meeting industry compliance standards and preventing breaches.', 'security', 4),
((SELECT id FROM solutions WHERE slug='corporate'), 'Collaboration Barriers', 'Siloed teams struggling with ineffective communication tools that hinder real-time collaboration and knowledge sharing.', 'forum', 5),
((SELECT id FROM solutions WHERE slug='corporate'), 'Digital Transformation', 'Legacy systems and resistance to change slowing down digital adoption and competitive positioning in the market.', 'trending_up', 6);

-- Corporate solution products
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured) VALUES
((SELECT id FROM solutions WHERE slug='corporate'), (SELECT id FROM products WHERE sku='BJ-D200'), 1, 1),
((SELECT id FROM solutions WHERE slug='corporate'), (SELECT id FROM products WHERE sku='BJ-IFP75'), 2, 1),
((SELECT id FROM solutions WHERE slug='corporate'), (SELECT id FROM products WHERE sku='BJ-OPS200'), 3, 1),
((SELECT id FROM solutions WHERE slug='corporate'), (SELECT id FROM products WHERE sku='BJ-CAM360'), 4, 0),
((SELECT id FROM solutions WHERE slug='corporate'), (SELECT id FROM products WHERE sku='BJ-WPS100'), 5, 0);

-- Corporate solution CTA
INSERT INTO solution_ctas (solution_id, heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name) VALUES
((SELECT id FROM solutions WHERE slug='corporate'), 'Transform Your Corporate Workspace', 'Discover how BlueJay''s integrated technology solutions can enhance productivity, streamline IT management, and future-proof your enterprise.', 'Schedule a Demo', '/contact', 'Download Brochure', '/downloads/corporate-brochure', '+91-120-456-7890', 'main_cta');

-- ============================================================
-- HEALTHCARE SOLUTION SUB-ENTITIES
-- ============================================================

-- Healthcare solution stats
INSERT INTO solution_stats (solution_id, value, label, display_order) VALUES
((SELECT id FROM solutions WHERE slug='healthcare'), '200+', 'Hospitals Served', 1),
((SELECT id FROM solutions WHERE slug='healthcare'), '99.8%', 'System Uptime', 2),
((SELECT id FROM solutions WHERE slug='healthcare'), '50K+', 'Patient Touchpoints Daily', 3);

-- Healthcare solution challenges
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order) VALUES
((SELECT id FROM solutions WHERE slug='healthcare'), 'Telemedicine Infrastructure', 'Building reliable video consultation platforms that work across varying internet connectivity conditions in urban and rural areas.', 'video_camera_front', 1),
((SELECT id FROM solutions WHERE slug='healthcare'), 'Patient Engagement', 'Modern patients expect digital-first experiences including online appointment booking, digital medical records access, and interactive communication.', 'how_to_reg', 2),
((SELECT id FROM solutions WHERE slug='healthcare'), 'Medical Training & Education', 'Medical colleges and hospitals need interactive technology for anatomy visualization, surgical procedure training, and collaborative case studies.', 'school', 3),
((SELECT id FROM solutions WHERE slug='healthcare'), 'Data Security & Compliance', 'Protecting sensitive patient health information while ensuring compliance with data protection regulations and healthcare standards.', 'security', 4),
((SELECT id FROM solutions WHERE slug='healthcare'), 'Operational Efficiency', 'Hospital administrators must optimize bed management, OPD workflows, pharmacy operations, and facility management while controlling costs.', 'speed', 5),
((SELECT id FROM solutions WHERE slug='healthcare'), 'Equipment Reliability', 'Healthcare environments demand 24/7 equipment uptime with minimal failures, as technology downtime directly impacts patient care quality.', 'verified', 6);

-- Healthcare solution products
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured) VALUES
((SELECT id FROM solutions WHERE slug='healthcare'), (SELECT id FROM products WHERE sku='BJ-D200'), 1, 1),
((SELECT id FROM solutions WHERE slug='healthcare'), (SELECT id FROM products WHERE sku='BJ-IFP75'), 2, 1),
((SELECT id FROM solutions WHERE slug='healthcare'), (SELECT id FROM products WHERE sku='BJ-DS55'), 3, 1),
((SELECT id FROM solutions WHERE slug='healthcare'), (SELECT id FROM products WHERE sku='BJ-IOT100'), 4, 0),
((SELECT id FROM solutions WHERE slug='healthcare'), (SELECT id FROM products WHERE sku='BJ-CAM360'), 5, 0);

-- Healthcare solution CTA
INSERT INTO solution_ctas (solution_id, heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name) VALUES
((SELECT id FROM solutions WHERE slug='healthcare'), 'Transform Your Healthcare Facility', 'Schedule a consultation to discuss how BlueJay technology solutions can improve patient care and operational efficiency at your hospital or clinic.', 'Request Consultation', '/contact', 'Download Brochure', '/downloads/healthcare-brochure', '+91-120-456-7890', 'main_cta');

-- ============================================================
-- RETAIL SOLUTION SUB-ENTITIES
-- ============================================================

-- Retail solution stats
INSERT INTO solution_stats (solution_id, value, label, display_order) VALUES
((SELECT id FROM solutions WHERE slug='retail'), '5000+', 'Retail Locations', 1),
((SELECT id FROM solutions WHERE slug='retail'), '35%', 'Avg. Sales Increase', 2),
((SELECT id FROM solutions WHERE slug='retail'), '24/7', 'System Uptime', 3);

-- Retail solution challenges
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order) VALUES
((SELECT id FROM solutions WHERE slug='retail'), 'Customer Engagement', 'Capture attention and drive conversions with interactive experiences that differentiate your brand in crowded retail environments.', 'groups', 1),
((SELECT id FROM solutions WHERE slug='retail'), 'Digital Signage Management', 'Manage and update content across multiple displays and locations efficiently without the complexity of traditional signage systems.', 'display_settings', 2),
((SELECT id FROM solutions WHERE slug='retail'), 'Omnichannel Experience', 'Deliver seamless experiences across online and offline channels, ensuring consistency in messaging, pricing, and customer service.', 'swap_horiz', 3),
((SELECT id FROM solutions WHERE slug='retail'), 'Inventory & Operations', 'Optimize stock levels, reduce shrinkage, and streamline operations with real-time visibility into inventory and store performance.', 'inventory_2', 4),
((SELECT id FROM solutions WHERE slug='retail'), 'Store Analytics', 'Gain actionable insights into foot traffic patterns, dwell times, conversion rates, and customer behavior to optimize store layouts.', 'analytics', 5),
((SELECT id FROM solutions WHERE slug='retail'), 'Brand Consistency', 'Maintain uniform brand messaging and visual standards across all locations while allowing for localized promotions and customization.', 'verified', 6);

-- Retail solution products
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured) VALUES
((SELECT id FROM solutions WHERE slug='retail'), (SELECT id FROM products WHERE sku='BJ-DS55'), 1, 1),
((SELECT id FROM solutions WHERE slug='retail'), (SELECT id FROM products WHERE sku='BJ-IFP65'), 2, 1),
((SELECT id FROM solutions WHERE slug='retail'), (SELECT id FROM products WHERE sku='BJ-D100'), 3, 0),
((SELECT id FROM solutions WHERE slug='retail'), (SELECT id FROM products WHERE sku='BJ-IOT100'), 4, 0),
((SELECT id FROM solutions WHERE slug='retail'), (SELECT id FROM products WHERE sku='BJ-OPS100'), 5, 0);

-- Retail solution CTA
INSERT INTO solution_ctas (solution_id, heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name) VALUES
((SELECT id FROM solutions WHERE slug='retail'), 'Ready to Transform Your Retail Experience?', 'Discover how BlueJay''s retail solutions can increase customer engagement, boost sales, and streamline operations across your stores.', 'Schedule a Consultation', '/contact', 'Download Brochure', '/downloads/retail-brochure', '+91-120-456-7890', 'main_cta');

-- ============================================================
-- GOVERNMENT SOLUTION SUB-ENTITIES
-- ============================================================

-- Government solution stats
INSERT INTO solution_stats (solution_id, value, label, display_order) VALUES
((SELECT id FROM solutions WHERE slug='government'), '100+', 'Government Departments', 1),
((SELECT id FROM solutions WHERE slug='government'), '500+', 'Smart City Installations', 2),
((SELECT id FROM solutions WHERE slug='government'), '99.8%', 'GeM Compliance Rate', 3);

-- Government solution challenges
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order) VALUES
((SELECT id FROM solutions WHERE slug='government'), 'GeM Compliance', 'Navigate Government e-Marketplace procurement with certified, compliant technology solutions that meet all regulatory requirements.', 'verified_user', 1),
((SELECT id FROM solutions WHERE slug='government'), 'Data Sovereignty', 'Ensure sensitive government data remains secure within national boundaries with Make in India solutions that meet data localization mandates.', 'shield', 2),
((SELECT id FROM solutions WHERE slug='government'), 'Legacy Modernization', 'Transform outdated infrastructure into modern, efficient digital systems without disrupting critical government operations.', 'system_update_alt', 3),
((SELECT id FROM solutions WHERE slug='government'), 'Citizen Service Delivery', 'Enable seamless, transparent public service delivery through digital platforms that improve accessibility and accountability.', 'groups', 4),
((SELECT id FROM solutions WHERE slug='government'), 'Security & Compliance', 'Meet stringent government security standards and audit requirements with enterprise-grade, certified technology solutions.', 'security', 5),
((SELECT id FROM solutions WHERE slug='government'), 'Budget Optimization', 'Maximize value from limited public funds with cost-effective, reliable solutions that deliver long-term operational efficiency.', 'account_balance_wallet', 6);

-- Government solution products
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured) VALUES
((SELECT id FROM solutions WHERE slug='government'), (SELECT id FROM products WHERE sku='BJ-D200'), 1, 1),
((SELECT id FROM solutions WHERE slug='government'), (SELECT id FROM products WHERE sku='BJ-D100'), 2, 0),
((SELECT id FROM solutions WHERE slug='government'), (SELECT id FROM products WHERE sku='BJ-IFP75'), 3, 1),
((SELECT id FROM solutions WHERE slug='government'), (SELECT id FROM products WHERE sku='BJ-DS55'), 4, 0),
((SELECT id FROM solutions WHERE slug='government'), (SELECT id FROM products WHERE sku='BJ-IOT100'), 5, 0);

-- Government solution CTA
INSERT INTO solution_ctas (solution_id, heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name) VALUES
((SELECT id FROM solutions WHERE slug='government'), 'Ready to Transform Government Services?', 'Partner with BlueJay for GeM-compliant, Make in India technology solutions that drive Digital India forward.', 'Request Consultation', '/contact', 'Download Brochure', '/downloads/government-brochure', '+91-120-456-7890', 'main_cta');

-- ============================================================
-- HOSPITALITY SOLUTION SUB-ENTITIES
-- ============================================================

-- Hospitality solution stats
INSERT INTO solution_stats (solution_id, value, label, display_order) VALUES
((SELECT id FROM solutions WHERE slug='hospitality'), '300+', 'Hotels & Resorts', 1),
((SELECT id FROM solutions WHERE slug='hospitality'), '98%', 'Guest Satisfaction', 2),
((SELECT id FROM solutions WHERE slug='hospitality'), '40%', 'Energy Savings', 3);

-- Hospitality solution challenges
INSERT INTO solution_challenges (solution_id, title, description, icon, display_order) VALUES
((SELECT id FROM solutions WHERE slug='hospitality'), 'Guest Experience Enhancement', 'Delivering personalized, technology-enabled experiences that exceed modern traveler expectations and drive positive reviews.', 'hotel', 1),
((SELECT id FROM solutions WHERE slug='hospitality'), 'Conference & Event Facilities', 'Providing state-of-the-art presentation technology for seamless meetings, weddings, and corporate events.', 'meeting_room', 2),
((SELECT id FROM solutions WHERE slug='hospitality'), 'Digital Wayfinding & Signage', 'Guiding guests through complex properties with dynamic digital displays that reduce confusion and enhance navigation.', 'signpost', 3),
((SELECT id FROM solutions WHERE slug='hospitality'), 'Room Automation & Comfort', 'Implementing smart controls for lighting, temperature, and entertainment that provide personalized comfort while optimizing energy usage.', 'nest_thermostat', 4),
((SELECT id FROM solutions WHERE slug='hospitality'), 'Energy Management', 'Reducing operational costs through intelligent IoT systems that monitor and control energy consumption across the property.', 'energy_savings_leaf', 5),
((SELECT id FROM solutions WHERE slug='hospitality'), 'Operational Efficiency', 'Streamlining front desk operations, housekeeping coordination, and guest services through integrated technology solutions.', 'speed', 6);

-- Hospitality solution products
INSERT INTO solution_products (solution_id, product_id, display_order, is_featured) VALUES
((SELECT id FROM solutions WHERE slug='hospitality'), (SELECT id FROM products WHERE sku='BJ-DS55'), 1, 1),
((SELECT id FROM solutions WHERE slug='hospitality'), (SELECT id FROM products WHERE sku='BJ-IFP65'), 2, 1),
((SELECT id FROM solutions WHERE slug='hospitality'), (SELECT id FROM products WHERE sku='BJ-WPS100'), 3, 0),
((SELECT id FROM solutions WHERE slug='hospitality'), (SELECT id FROM products WHERE sku='BJ-IOT100'), 4, 0),
((SELECT id FROM solutions WHERE slug='hospitality'), (SELECT id FROM products WHERE sku='BJ-D300'), 5, 0);

-- Hospitality solution CTA
INSERT INTO solution_ctas (solution_id, heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, phone_number, section_name) VALUES
((SELECT id FROM solutions WHERE slug='hospitality'), 'Transform Your Guest Experience', 'Discover how BlueJay''s hospitality solutions can elevate service quality, improve operational efficiency, and drive higher guest satisfaction scores.', 'Schedule a Property Consultation', '/contact', 'Download Brochure', '/downloads/hospitality-brochure', '+91-120-456-7890', 'main_cta');

-- ============================================================
-- SHARED PAGE ELEMENTS
-- ============================================================

-- Solution page features (Why Choose BlueJay)
INSERT INTO solution_page_features (title, description, icon, display_order, is_active) VALUES
('Industry Expertise', 'Deep understanding of vertical needs and requirements', 'insights', 1, 1),
('End-to-End Solutions', 'From hardware to after-sales support', 'sync_alt', 2, 1),
('Certified Quality', 'ISO, BIS, CE, FCC certified products', 'verified', 3, 1),
('Local Support', 'Dedicated support teams across regions', 'support_agent', 4, 1);

-- Solutions listing CTA
INSERT INTO solutions_listing_cta (heading, subheading, primary_button_text, primary_button_url, secondary_button_text, secondary_button_url, is_active) VALUES
('Not Sure Which Solution Fits Your Needs?', 'Our team can help you identify the right products and solutions for your specific industry requirements.', 'Contact Our Team', '/contact', 'View Products', '/products', 1);
