-- Add missing SEO columns to content tables

-- blog_posts: add meta_title and og_image
ALTER TABLE blog_posts ADD COLUMN meta_title TEXT NOT NULL DEFAULT '';
ALTER TABLE blog_posts ADD COLUMN og_image TEXT NOT NULL DEFAULT '';

-- whitepapers: add meta_title and og_image
ALTER TABLE whitepapers ADD COLUMN meta_title TEXT NOT NULL DEFAULT '';
ALTER TABLE whitepapers ADD COLUMN og_image TEXT NOT NULL DEFAULT '';

-- solutions: add meta_title and og_image
ALTER TABLE solutions ADD COLUMN meta_title TEXT NOT NULL DEFAULT '';
ALTER TABLE solutions ADD COLUMN og_image TEXT NOT NULL DEFAULT '';

-- products: add og_image
ALTER TABLE products ADD COLUMN og_image TEXT NOT NULL DEFAULT '';

-- case_studies: add og_image
ALTER TABLE case_studies ADD COLUMN og_image TEXT NOT NULL DEFAULT '';
