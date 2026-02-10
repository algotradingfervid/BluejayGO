-- Education Case Study
INSERT INTO case_studies (
    slug, title, client_name, industry_id, hero_image_url, summary,
    challenge_title, challenge_content, challenge_bullets,
    solution_title, solution_content,
    outcome_title, outcome_content,
    is_published, display_order
) VALUES (
    'smart-classroom-transformation-abc-university',
    'Smart Classroom Transformation at ABC University',
    'ABC University',
    (SELECT id FROM industries WHERE slug = 'education'),
    '/uploads/case-studies/case-study-1.jpg',
    'How we deployed 500+ interactive panels across 50 campuses to revolutionize the learning experience.',
    'The Challenge',
    '<p>ABC University faced significant challenges in modernizing their teaching infrastructure across 50 campuses. Outdated projectors and whiteboards hindered interactive learning, while maintenance costs continued to rise.</p>',
    '["Inconsistent AV quality across campuses", "High maintenance costs for outdated equipment", "Limited collaboration capabilities", "Poor remote learning support infrastructure"]',
    'Our Solution',
    '<p>BlueJay worked closely with ABC University to design and implement a comprehensive smart classroom solution. Our team deployed state-of-the-art interactive flat panels integrated with collaborative software across all campuses.</p>',
    'The Outcome',
    '<p>The implementation was completed in three phases over 18 months with minimal disruption. The university reported significant improvements in student engagement and teaching effectiveness.</p>',
    1, 1
);

INSERT INTO case_study_metrics (case_study_id, metric_value, metric_label, display_order) VALUES
    ((SELECT id FROM case_studies WHERE slug = 'smart-classroom-transformation-abc-university'), '500+', 'Interactive Panels Deployed', 1),
    ((SELECT id FROM case_studies WHERE slug = 'smart-classroom-transformation-abc-university'), '50', 'Campuses Equipped', 2),
    ((SELECT id FROM case_studies WHERE slug = 'smart-classroom-transformation-abc-university'), '40%', 'Reduction in AV Maintenance', 3),
    ((SELECT id FROM case_studies WHERE slug = 'smart-classroom-transformation-abc-university'), '95%', 'Faculty Satisfaction Rate', 4);

-- Education Case Study - Products
INSERT INTO case_study_products (case_study_id, product_id, display_order) VALUES
    ((SELECT id FROM case_studies WHERE slug = 'smart-classroom-transformation-abc-university'), (SELECT id FROM products WHERE sku = 'BJ-IFP75'), 1),
    ((SELECT id FROM case_studies WHERE slug = 'smart-classroom-transformation-abc-university'), (SELECT id FROM products WHERE sku = 'BJ-OPS200'), 2),
    ((SELECT id FROM case_studies WHERE slug = 'smart-classroom-transformation-abc-university'), (SELECT id FROM products WHERE sku = 'BJ-WPS100'), 3);

-- Corporate Case Study
INSERT INTO case_studies (
    slug, title, client_name, industry_id, hero_image_url, summary,
    challenge_title, challenge_content, challenge_bullets,
    solution_title, solution_content,
    outcome_title, outcome_content,
    is_published, display_order
) VALUES (
    'digital-meeting-room-overhaul-techcorp',
    'Digital Meeting Room Overhaul at TechCorp',
    'TechCorp Industries',
    (SELECT id FROM industries WHERE slug = 'corporate'),
    '/uploads/case-studies/case-study-1.jpg',
    'Modernizing 200+ meeting rooms with unified collaboration technology.',
    'The Challenge',
    '<p>TechCorp struggled with fragmented video conferencing systems, complex IT support requirements, and inconsistent meeting experiences across their global offices.</p>',
    '["Incompatible legacy video conferencing systems", "Complex IT support requirements", "Low adoption rates for collaboration tools", "Inconsistent meeting experience across offices"]',
    'Our Solution',
    '<p>We deployed a unified meeting room solution with a standardized technology stack, including wireless presentation systems, video conferencing platforms, and centralized room management.</p>',
    'The Outcome',
    '<p>Achieved 98% system uptime, reduced IT support tickets by 60%, and significantly improved employee satisfaction with meeting room technology.</p>',
    1, 2
);

INSERT INTO case_study_metrics (case_study_id, metric_value, metric_label, display_order) VALUES
    ((SELECT id FROM case_studies WHERE slug = 'digital-meeting-room-overhaul-techcorp'), '200+', 'Meeting Rooms Upgraded', 1),
    ((SELECT id FROM case_studies WHERE slug = 'digital-meeting-room-overhaul-techcorp'), '60%', 'Reduction in IT Tickets', 2),
    ((SELECT id FROM case_studies WHERE slug = 'digital-meeting-room-overhaul-techcorp'), '98%', 'System Uptime', 3),
    ((SELECT id FROM case_studies WHERE slug = 'digital-meeting-room-overhaul-techcorp'), '4.8/5', 'User Satisfaction Score', 4);

-- Corporate Case Study - Products
INSERT INTO case_study_products (case_study_id, product_id, display_order) VALUES
    ((SELECT id FROM case_studies WHERE slug = 'digital-meeting-room-overhaul-techcorp'), (SELECT id FROM products WHERE sku = 'BJ-IFP65'), 1),
    ((SELECT id FROM case_studies WHERE slug = 'digital-meeting-room-overhaul-techcorp'), (SELECT id FROM products WHERE sku = 'BJ-CAM360'), 2),
    ((SELECT id FROM case_studies WHERE slug = 'digital-meeting-room-overhaul-techcorp'), (SELECT id FROM products WHERE sku = 'BJ-SPK200'), 3);

-- Healthcare Case Study
INSERT INTO case_studies (
    slug, title, client_name, industry_id, hero_image_url, summary,
    challenge_title, challenge_content, challenge_bullets,
    solution_title, solution_content,
    outcome_title, outcome_content,
    is_published, display_order
) VALUES (
    'patient-communication-system-metro-health',
    'Patient Communication System at Metro Health',
    'Metro Health Network',
    (SELECT id FROM industries WHERE slug = 'healthcare'),
    '/uploads/case-studies/case-study-1.jpg',
    'Implementing digital signage and communication across 12 hospital facilities.',
    'The Challenge',
    '<p>Metro Health Network relied on paper-based patient information systems with confusing hospital wayfinding and inefficient staff scheduling displays that impacted patient experience.</p>',
    '["Paper-based patient information systems", "Confusing hospital wayfinding", "Inefficient staff scheduling displays", "Limited real-time communication capability"]',
    'Our Solution',
    '<p>We implemented a comprehensive digital signage solution with real-time data integration, interactive wayfinding kiosks, and staff communication displays throughout all hospital facilities.</p>',
    'The Outcome',
    '<p>The system improved patient satisfaction scores, reduced perceived wait times by 35%, and enhanced staff communication efficiency across the entire network.</p>',
    1, 3
);

INSERT INTO case_study_metrics (case_study_id, metric_value, metric_label, display_order) VALUES
    ((SELECT id FROM case_studies WHERE slug = 'patient-communication-system-metro-health'), '12', 'Hospital Facilities', 1),
    ((SELECT id FROM case_studies WHERE slug = 'patient-communication-system-metro-health'), '35%', 'Reduced Wait Perception', 2),
    ((SELECT id FROM case_studies WHERE slug = 'patient-communication-system-metro-health'), '300+', 'Digital Displays', 3),
    ((SELECT id FROM case_studies WHERE slug = 'patient-communication-system-metro-health'), '92%', 'Patient Satisfaction', 4);

-- Healthcare Case Study - Products
INSERT INTO case_study_products (case_study_id, product_id, display_order) VALUES
    ((SELECT id FROM case_studies WHERE slug = 'patient-communication-system-metro-health'), (SELECT id FROM products WHERE sku = 'BJ-DS55'), 1),
    ((SELECT id FROM case_studies WHERE slug = 'patient-communication-system-metro-health'), (SELECT id FROM products WHERE sku = 'BJ-IOT100'), 2);
