-- Optional branding and destination links; blank values keep links hidden.
ALTER TABLE settings ADD COLUMN footer_logo_path TEXT NOT NULL DEFAULT '';
ALTER TABLE settings ADD COLUMN social_threads TEXT NOT NULL DEFAULT '';
ALTER TABLE settings ADD COLUMN marketplace_gem_url TEXT NOT NULL DEFAULT '';
ALTER TABLE settings ADD COLUMN marketplace_amazon_url TEXT NOT NULL DEFAULT '';
