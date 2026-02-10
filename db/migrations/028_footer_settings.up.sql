-- Add footer settings columns to settings table
ALTER TABLE settings ADD COLUMN footer_columns INTEGER NOT NULL DEFAULT 4;
ALTER TABLE settings ADD COLUMN footer_bg_style TEXT NOT NULL DEFAULT 'dark';
ALTER TABLE settings ADD COLUMN footer_show_social INTEGER NOT NULL DEFAULT 1;
ALTER TABLE settings ADD COLUMN footer_social_style TEXT NOT NULL DEFAULT 'icons';
ALTER TABLE settings ADD COLUMN footer_copyright TEXT NOT NULL DEFAULT '© {year} All rights reserved.';

-- Footer column items
CREATE TABLE IF NOT EXISTS footer_column_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    column_index INTEGER NOT NULL,
    type TEXT NOT NULL DEFAULT 'links',
    heading TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0
);

-- Footer links (for columns of type "links")
CREATE TABLE IF NOT EXISTS footer_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    column_item_id INTEGER NOT NULL,
    label TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (column_item_id) REFERENCES footer_column_items(id) ON DELETE CASCADE
);

-- Footer legal links (bottom bar)
CREATE TABLE IF NOT EXISTS footer_legal_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    label TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0
);
