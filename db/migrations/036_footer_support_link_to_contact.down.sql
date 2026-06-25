-- Reverse 036: point the footer RESOURCES "Support" link back to its original
-- /support target.
UPDATE page_sections SET primary_button_url = '/support' WHERE page_key = 'footer' AND section_key = 'resource_4';
