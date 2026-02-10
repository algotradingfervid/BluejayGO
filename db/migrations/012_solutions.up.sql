-- Solutions main table
CREATE TABLE IF NOT EXISTS solutions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    icon TEXT NOT NULL,
    short_description TEXT NOT NULL,
    hero_image_url TEXT,
    hero_title TEXT,
    hero_description TEXT,
    overview_content TEXT,
    meta_description TEXT,
    reference_code TEXT,
    is_published BOOLEAN DEFAULT 0,
    display_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_solutions_slug ON solutions(slug);
CREATE INDEX IF NOT EXISTS idx_solutions_published ON solutions(is_published, display_order);

-- Solution statistics (Industry Overview section)
CREATE TABLE IF NOT EXISTS solution_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    solution_id INTEGER NOT NULL,
    value TEXT NOT NULL,
    label TEXT NOT NULL,
    display_order INTEGER DEFAULT 0,
    FOREIGN KEY (solution_id) REFERENCES solutions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_solution_stats_solution ON solution_stats(solution_id);

-- Solution challenges
CREATE TABLE IF NOT EXISTS solution_challenges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    solution_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    icon TEXT NOT NULL,
    display_order INTEGER DEFAULT 0,
    FOREIGN KEY (solution_id) REFERENCES solutions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_solution_challenges_solution ON solution_challenges(solution_id);

-- Solution-Product relationship
CREATE TABLE IF NOT EXISTS solution_products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    solution_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    display_order INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT 0,
    FOREIGN KEY (solution_id) REFERENCES solutions(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE(solution_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_solution_products_solution ON solution_products(solution_id);
CREATE INDEX IF NOT EXISTS idx_solution_products_product ON solution_products(product_id);

-- Solution CTAs
CREATE TABLE IF NOT EXISTS solution_ctas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    solution_id INTEGER NOT NULL,
    heading TEXT NOT NULL,
    subheading TEXT,
    primary_button_text TEXT,
    primary_button_url TEXT,
    secondary_button_text TEXT,
    secondary_button_url TEXT,
    phone_number TEXT,
    section_name TEXT NOT NULL,
    FOREIGN KEY (solution_id) REFERENCES solutions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_solution_ctas_solution ON solution_ctas(solution_id);

-- Listing page features (Why Choose BlueJay)
CREATE TABLE IF NOT EXISTS solution_page_features (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    icon TEXT NOT NULL,
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT 1
);

-- Listing page CTA
CREATE TABLE IF NOT EXISTS solutions_listing_cta (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    heading TEXT NOT NULL,
    subheading TEXT,
    primary_button_text TEXT,
    primary_button_url TEXT,
    secondary_button_text TEXT,
    secondary_button_url TEXT,
    is_active BOOLEAN DEFAULT 1
);
