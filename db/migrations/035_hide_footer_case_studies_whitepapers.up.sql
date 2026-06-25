-- Hide the "Case Studies" and "Whitepapers" footer RESOURCES links from the public
-- front end by deactivating their page_sections rows. This is reversible: the admin
-- Page Sections "Active" toggle (or the matching .down.sql) can re-enable them later.
UPDATE page_sections SET is_active = 0 WHERE page_key = 'footer' AND section_key IN ('resource_1', 'resource_2');
