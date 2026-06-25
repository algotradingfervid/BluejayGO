-- Repoint the footer RESOURCES "Support" link from the dead /support route (which
-- has no handler and 404s) to the existing /contact page, which already handles
-- support inquiries. This is reversible: the matching .down.sql (or the admin
-- Page Sections editor) can point it back to /support later.
UPDATE page_sections SET primary_button_url = '/contact' WHERE page_key = 'footer' AND section_key = 'resource_4';
