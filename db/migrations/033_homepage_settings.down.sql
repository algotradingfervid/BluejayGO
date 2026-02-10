-- SQLite doesn't support DROP COLUMN in older versions, but modernc/sqlite does
ALTER TABLE settings DROP COLUMN homepage_show_heroes;
ALTER TABLE settings DROP COLUMN homepage_show_stats;
ALTER TABLE settings DROP COLUMN homepage_show_testimonials;
ALTER TABLE settings DROP COLUMN homepage_show_cta;
ALTER TABLE settings DROP COLUMN homepage_max_heroes;
ALTER TABLE settings DROP COLUMN homepage_max_stats;
ALTER TABLE settings DROP COLUMN homepage_max_testimonials;
ALTER TABLE settings DROP COLUMN homepage_hero_autoplay;
ALTER TABLE settings DROP COLUMN homepage_hero_interval;
