-- Re-enable the "Case Studies" and "Whitepapers" footer RESOURCES links.
UPDATE page_sections SET is_active = 1 WHERE page_key = 'footer' AND section_key IN ('resource_1', 'resource_2');
