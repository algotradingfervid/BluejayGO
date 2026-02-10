DROP INDEX IF EXISTS idx_partners_featured;
-- SQLite doesn't support DROP COLUMN, so we skip removing is_featured from partners
DROP TABLE IF EXISTS homepage_cta;
DROP TABLE IF EXISTS homepage_testimonials;
DROP TABLE IF EXISTS homepage_stats;
DROP TABLE IF EXISTS homepage_hero;
