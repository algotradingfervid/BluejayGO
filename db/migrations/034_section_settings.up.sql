-- About section settings
ALTER TABLE settings ADD COLUMN about_show_mission INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN about_show_milestones INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN about_show_certifications INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN about_show_team INTEGER NOT NULL DEFAULT 1;

-- Products section settings
ALTER TABLE settings ADD COLUMN products_per_page INTEGER NOT NULL DEFAULT 12;
ALTER TABLE settings ADD COLUMN products_show_categories INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN products_show_search INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN products_default_sort TEXT NOT NULL DEFAULT 'name_asc';

-- Solutions section settings
ALTER TABLE settings ADD COLUMN solutions_per_page INTEGER NOT NULL DEFAULT 12;
ALTER TABLE settings ADD COLUMN solutions_show_industries INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN solutions_show_search INTEGER NOT NULL DEFAULT 1;

-- Blog section settings
ALTER TABLE settings ADD COLUMN blog_posts_per_page INTEGER NOT NULL DEFAULT 10;
ALTER TABLE settings ADD COLUMN blog_show_author INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN blog_show_date INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN blog_show_categories INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN blog_show_tags INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN blog_show_search INTEGER NOT NULL DEFAULT 1;
