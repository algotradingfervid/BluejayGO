CREATE TABLE IF NOT EXISTS whitepapers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    topic_id INTEGER NOT NULL REFERENCES whitepaper_topics(id) ON DELETE RESTRICT,
    pdf_file_path TEXT NOT NULL,
    file_size_bytes INTEGER NOT NULL,
    page_count INTEGER,
    published_date TEXT NOT NULL,
    is_published INTEGER NOT NULL DEFAULT 0,
    cover_color_from TEXT NOT NULL DEFAULT '#0066CC',
    cover_color_to TEXT NOT NULL DEFAULT '#004499',
    meta_description TEXT,
    download_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_whitepapers_slug ON whitepapers(slug);
CREATE INDEX idx_whitepapers_topic_id ON whitepapers(topic_id);
CREATE INDEX idx_whitepapers_published ON whitepapers(is_published);

CREATE TABLE IF NOT EXISTS whitepaper_learning_points (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    whitepaper_id INTEGER NOT NULL REFERENCES whitepapers(id) ON DELETE CASCADE,
    point_text TEXT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_whitepaper_learning_points_whitepaper ON whitepaper_learning_points(whitepaper_id);

CREATE TABLE IF NOT EXISTS whitepaper_downloads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    whitepaper_id INTEGER NOT NULL REFERENCES whitepapers(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    company TEXT NOT NULL,
    designation TEXT,
    marketing_consent INTEGER NOT NULL DEFAULT 0,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_whitepaper_downloads_whitepaper ON whitepaper_downloads(whitepaper_id);
CREATE INDEX idx_whitepaper_downloads_email ON whitepaper_downloads(email);
CREATE INDEX idx_whitepaper_downloads_created ON whitepaper_downloads(created_at);
