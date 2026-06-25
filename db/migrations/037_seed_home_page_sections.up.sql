-- Seed the remaining "home" page_sections so every homepage content section has an
-- EDITABLE row in the /admin/page-sections editor (which only edits existing rows).
-- Migration 013 only seeded 'hero' and 'solutions_section'.
--
-- INSERT OR IGNORE keeps this idempotent against the unique (page_key, section_key)
-- index, so it won't clobber rows that already exist in a live DB.
INSERT OR IGNORE INTO page_sections (page_key, section_key, heading, subheading, display_order) VALUES
('home', 'products_section', 'Featured Products', 'Discover our most popular interactive solutions', 2),
('home', 'stats_section', '', '', 3),
('home', 'testimonials_section', 'What Our Clients Say', 'Trusted by organizations worldwide', 4),
('home', 'partners_section', 'Trusted Partners', '', 5),
('home', 'blog_section', 'Latest News & Insights', 'Stay updated with industry trends', 6);

-- Align the solutions_section heading with the homepage's prior visible default, but
-- only when it still holds the untouched 013 seed value ('Our Solutions') so we never
-- overwrite an admin-customized heading.
UPDATE page_sections
SET heading = 'Solutions By Industry',
    subheading = 'Tailored technology solutions for every sector'
WHERE page_key = 'home' AND section_key = 'solutions_section' AND heading = 'Our Solutions';
