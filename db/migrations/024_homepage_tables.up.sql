-- Homepage Hero
CREATE TABLE IF NOT EXISTS homepage_hero (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    headline TEXT NOT NULL,
    subheadline TEXT NOT NULL,
    badge_text TEXT,
    primary_cta_text TEXT NOT NULL DEFAULT 'Explore Products',
    primary_cta_url TEXT NOT NULL DEFAULT '/products',
    secondary_cta_text TEXT,
    secondary_cta_url TEXT,
    background_image TEXT,
    is_active INTEGER NOT NULL DEFAULT 1,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Homepage Stats
CREATE TABLE IF NOT EXISTS homepage_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stat_value TEXT NOT NULL,
    stat_label TEXT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Homepage Testimonials
CREATE TABLE IF NOT EXISTS homepage_testimonials (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quote TEXT NOT NULL,
    author_name TEXT NOT NULL,
    author_title TEXT,
    author_company TEXT,
    author_image TEXT,
    rating INTEGER NOT NULL DEFAULT 5,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Homepage CTA
CREATE TABLE IF NOT EXISTS homepage_cta (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    headline TEXT NOT NULL,
    description TEXT,
    primary_cta_text TEXT NOT NULL DEFAULT 'Schedule a Demo',
    primary_cta_url TEXT NOT NULL DEFAULT '/contact',
    secondary_cta_text TEXT,
    secondary_cta_url TEXT,
    background_style TEXT DEFAULT 'primary',
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add is_featured to partners
ALTER TABLE partners ADD COLUMN is_featured INTEGER NOT NULL DEFAULT 0;
CREATE INDEX idx_partners_featured ON partners(is_featured);
