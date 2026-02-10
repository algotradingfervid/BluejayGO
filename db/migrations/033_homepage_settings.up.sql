-- Homepage section visibility and display settings
ALTER TABLE settings ADD COLUMN homepage_show_heroes INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN homepage_show_stats INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN homepage_show_testimonials INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN homepage_show_cta INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN homepage_max_heroes INTEGER NOT NULL DEFAULT 5;
ALTER TABLE settings ADD COLUMN homepage_max_stats INTEGER NOT NULL DEFAULT 6;
ALTER TABLE settings ADD COLUMN homepage_max_testimonials INTEGER NOT NULL DEFAULT 3;
ALTER TABLE settings ADD COLUMN homepage_hero_autoplay INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN homepage_hero_interval INTEGER NOT NULL DEFAULT 5;
