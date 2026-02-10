-- Seed data for About page
-- BlueJay Innovative Labs company information

-- Company Overview
INSERT INTO company_overview (
    headline,
    tagline,
    description_main,
    description_secondary,
    description_tertiary
) VALUES (
    'Transforming Collaboration Through Innovative Display Technology',
    'Where Innovation Meets Intelligence',
    'BlueJay Innovative Labs is a leading B2B technology company specializing in cutting-edge interactive flat panels, high-performance desktops, OPS modules, professional AV accessories, document cameras, and custom kiosks. We empower education and enterprise sectors with solutions that enhance collaboration, productivity, and engagement.',
    'Our approach combines advanced engineering with user-centric design, ensuring every product delivers exceptional performance and reliability. We work closely with our partners to understand their unique challenges and develop technology solutions that drive measurable results and transform the way teams work together.',
    'With a global footprint spanning over 45 countries, BlueJay Innovative Labs has become a trusted partner for thousands of educational institutions and enterprises worldwide. Our commitment to innovation and quality has made us a preferred choice for organizations seeking to modernize their technology infrastructure and create dynamic, connected environments.'
);

-- Mission, Vision, and Values
INSERT INTO mission_vision_values (
    mission,
    vision,
    values_summary,
    mission_icon,
    vision_icon,
    values_icon
) VALUES (
    'To empower organizations with innovative technology solutions that transform collaboration, enhance learning experiences, and drive business success through cutting-edge interactive displays and intelligent systems.',
    'To be the global leader in collaborative technology, setting the standard for innovation, quality, and customer success while creating a connected world where every interaction is meaningful and productive.',
    'Our core values guide every decision we make, from product development to customer partnerships. We are committed to innovation, quality excellence, customer success, integrity, collaboration, and environmental sustainability.',
    'flag',
    'visibility',
    'diamond'
);

-- Core Values
INSERT INTO core_values (title, description, icon, display_order) VALUES
(
    'Innovation',
    'We continuously push the boundaries of what''s possible, investing in research and development to create breakthrough technologies that anticipate and exceed market needs. Our culture of innovation drives us to challenge conventions and deliver solutions that define the future of collaboration.',
    'lightbulb',
    1
),
(
    'Quality Excellence',
    'Every product bearing the BlueJay name undergoes rigorous testing and quality assurance processes. We are committed to delivering solutions that meet the highest standards of performance, reliability, and durability, ensuring our customers receive exceptional value.',
    'verified',
    2
),
(
    'Customer Success',
    'Our customers'' success is our success. We go beyond selling products to become trusted partners, providing comprehensive support, training, and consulting services that help organizations maximize their technology investments and achieve their strategic objectives.',
    'support_agent',
    3
),
(
    'Integrity',
    'We operate with unwavering honesty, transparency, and ethical standards in all our business practices. Building trust with our customers, partners, and employees is fundamental to who we are, and we honor our commitments with consistency and accountability.',
    'handshake',
    4
),
(
    'Collaboration',
    'Just as our products enable collaboration, we believe in the power of working together. We foster partnerships across industries, engage with our communities, and create an inclusive workplace where diverse perspectives drive better solutions and stronger outcomes.',
    'group',
    5
),
(
    'Sustainability',
    'We are dedicated to environmental responsibility and sustainable business practices. From energy-efficient product designs to eco-friendly manufacturing processes, we strive to minimize our environmental impact and contribute to a healthier planet for future generations.',
    'eco',
    6
);

-- Company Milestones
INSERT INTO milestones (year, title, description, is_current, display_order) VALUES
(
    2008,
    'Company Founded',
    'BlueJay Innovative Labs was established with a vision to revolutionize collaborative technology for education and enterprise sectors.',
    0,
    1
),
(
    2010,
    'First Interactive Flat Panel Launch',
    'Introduced our flagship interactive flat panel series, setting new standards for touch responsiveness and display quality in educational environments.',
    0,
    2
),
(
    2013,
    'Global Expansion Begins',
    'Opened international offices and distribution channels across Asia-Pacific and Europe, expanding our reach to serve customers worldwide.',
    0,
    3
),
(
    2016,
    'ISO 9001 & ISO 14001 Certification',
    'Achieved ISO 9001 quality management and ISO 14001 environmental management certifications, demonstrating our commitment to excellence and sustainability.',
    0,
    4
),
(
    2018,
    'OPS Module Innovation',
    'Launched our advanced OPS module lineup, offering seamless computing integration with interactive displays for enhanced performance and flexibility.',
    0,
    5
),
(
    2021,
    'Strategic Enterprise Partnerships',
    'Formed key partnerships with Fortune 500 companies and leading educational institutions, solidifying our position as a trusted B2B technology provider.',
    0,
    6
),
(
    2023,
    'AI-Powered Interactive Solutions',
    'Introduced AI-enhanced features across our product portfolio, including smart gesture recognition, intelligent collaboration tools, and predictive maintenance capabilities.',
    0,
    7
),
(
    2024,
    'Industry Leadership & Global Recognition',
    'Serving over 10,000 organizations across 45+ countries with comprehensive technology solutions, recognized as a leader in interactive display and collaborative technology innovation.',
    1,
    8
);

-- Certifications and Standards
INSERT INTO certifications (name, abbreviation, description, icon, display_order) VALUES
(
    'ISO 9001:2015 Quality Management',
    'ISO 9001',
    'International standard for quality management systems, ensuring consistent product quality and customer satisfaction through rigorous process controls and continuous improvement.',
    'verified_user',
    1
),
(
    'ISO 14001:2015 Environmental Management',
    'ISO 14001',
    'Environmental management system certification demonstrating our commitment to sustainable practices, waste reduction, and minimizing environmental impact across all operations.',
    'eco',
    2
),
(
    'ISO 27001:2013 Information Security',
    'ISO 27001',
    'Information security management standard ensuring robust protection of customer data, intellectual property, and sensitive business information through comprehensive security controls.',
    'security',
    3
),
(
    'Federal Communications Commission',
    'FCC',
    'FCC compliance certification for electromagnetic interference standards, ensuring all products meet regulatory requirements for safe operation in commercial and educational environments.',
    'cell_tower',
    4
),
(
    'Conformité Européenne',
    'CE',
    'European conformity marking indicating compliance with EU safety, health, and environmental protection standards for product distribution across European markets.',
    'public',
    5
),
(
    'ENERGY STAR Certified',
    'ENERGY STAR',
    'EPA recognition for superior energy efficiency, demonstrating our commitment to reducing power consumption and operating costs while minimizing environmental impact.',
    'bolt',
    6
);
