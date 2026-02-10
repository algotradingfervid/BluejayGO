DELETE FROM page_sections WHERE page_key = 'product_detail' AND section_key IN ('overview_section', 'video_section', 'features_section', 'specs_section', 'certifications_section', 'downloads_section');
DELETE FROM page_sections WHERE page_key = 'products_category';
DELETE FROM page_sections WHERE page_key = 'footer';
