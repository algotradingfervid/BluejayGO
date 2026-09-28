-- Older installations have zero for every homepage content section because
-- ordering was not exposed in admin. Give only those unordered rows sensible
-- defaults; preserve every nonzero order an administrator has already chosen.
UPDATE page_sections
SET display_order = CASE section_key
    WHEN 'hero' THEN 1
    WHEN 'products_section' THEN 2
    WHEN 'solutions_section' THEN 3
    WHEN 'stats_section' THEN 4
    WHEN 'partners_section' THEN 5
    WHEN 'testimonials_section' THEN 6
    WHEN 'blog_section' THEN 7
    WHEN 'cta' THEN 8
END
WHERE page_key = 'home' AND display_order = 0
  AND section_key IN ('hero', 'products_section', 'solutions_section', 'stats_section',
                      'partners_section', 'testimonials_section', 'blog_section', 'cta');

-- Give the final homepage CTA the same visibility and ordering controls.
INSERT OR IGNORE INTO page_sections (page_key, section_key, heading, display_order)
VALUES ('home', 'cta', 'Call to Action', 8);
