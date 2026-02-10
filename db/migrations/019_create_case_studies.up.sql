CREATE TABLE IF NOT EXISTS case_studies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    client_name TEXT NOT NULL,
    industry_id INTEGER NOT NULL REFERENCES industries(id) ON DELETE RESTRICT,
    hero_image_url TEXT,
    summary TEXT NOT NULL,
    challenge_title TEXT NOT NULL DEFAULT 'The Challenge',
    challenge_content TEXT NOT NULL,
    challenge_bullets TEXT, -- JSON array
    solution_title TEXT NOT NULL DEFAULT 'Our Solution',
    solution_content TEXT NOT NULL,
    outcome_title TEXT NOT NULL DEFAULT 'The Outcome',
    outcome_content TEXT NOT NULL,
    meta_title TEXT,
    meta_description TEXT,
    is_published INTEGER NOT NULL DEFAULT 0,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_case_studies_slug ON case_studies(slug);
CREATE INDEX idx_case_studies_industry ON case_studies(industry_id);
CREATE INDEX idx_case_studies_published ON case_studies(is_published);
CREATE INDEX idx_case_studies_display_order ON case_studies(display_order);

CREATE TABLE IF NOT EXISTS case_study_products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    case_study_id INTEGER NOT NULL REFERENCES case_studies(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(case_study_id, product_id)
);
CREATE INDEX idx_case_study_products_case_study ON case_study_products(case_study_id);
CREATE INDEX idx_case_study_products_product ON case_study_products(product_id);

CREATE TABLE IF NOT EXISTS case_study_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    case_study_id INTEGER NOT NULL REFERENCES case_studies(id) ON DELETE CASCADE,
    metric_value TEXT NOT NULL,
    metric_label TEXT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_case_study_metrics_case_study ON case_study_metrics(case_study_id);
