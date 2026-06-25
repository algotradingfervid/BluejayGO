-- Remove the home page_sections rows seeded by 037.
DELETE FROM page_sections
WHERE page_key = 'home'
  AND section_key IN ('products_section', 'testimonials_section', 'blog_section', 'partners_section', 'stats_section');

-- Restore the solutions_section heading to the 013 default, but only if it still holds
-- the value 037 set (so we don't clobber an admin-customized heading on rollback).
UPDATE page_sections
SET heading = 'Our Solutions',
    subheading = ''
WHERE page_key = 'home' AND section_key = 'solutions_section' AND heading = 'Solutions By Industry';
