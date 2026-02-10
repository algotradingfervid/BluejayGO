CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    site_name TEXT NOT NULL DEFAULT 'BlueJay Innovative Labs',
    site_tagline TEXT NOT NULL DEFAULT 'Innovation Through Technology',
    contact_email TEXT NOT NULL DEFAULT 'info@bluejaylabs.com',
    contact_phone TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    footer_text TEXT NOT NULL DEFAULT '',
    meta_description TEXT NOT NULL DEFAULT '',
    meta_keywords TEXT NOT NULL DEFAULT '',
    google_analytics_id TEXT NOT NULL DEFAULT '',
    social_linkedin TEXT NOT NULL DEFAULT '',
    social_twitter TEXT NOT NULL DEFAULT '',
    social_github TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_settings_singleton ON settings(id);
INSERT INTO settings (id, site_name, site_tagline, contact_email) VALUES (1, 'BlueJay Innovative Labs', 'Innovation Through Technology', 'info@bluejaylabs.com');
