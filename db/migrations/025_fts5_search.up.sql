-- FTS5 virtual tables for full-text search
CREATE VIRTUAL TABLE IF NOT EXISTS products_fts USING fts5(
    name, tagline, description,
    content='products', content_rowid='id'
);

CREATE VIRTUAL TABLE IF NOT EXISTS blog_posts_fts USING fts5(
    title, excerpt, body,
    content='blog_posts', content_rowid='id'
);

CREATE VIRTUAL TABLE IF NOT EXISTS case_studies_fts USING fts5(
    title, client_name, challenge_content, solution_content,
    content='case_studies', content_rowid='id'
);

-- Triggers to keep FTS in sync with content tables
CREATE TRIGGER products_ai AFTER INSERT ON products BEGIN
    INSERT INTO products_fts(rowid, name, tagline, description) VALUES (new.id, new.name, COALESCE(new.tagline, ''), new.description);
END;

CREATE TRIGGER products_ad AFTER DELETE ON products BEGIN
    INSERT INTO products_fts(products_fts, rowid, name, tagline, description) VALUES('delete', old.id, old.name, COALESCE(old.tagline, ''), old.description);
END;

CREATE TRIGGER products_au AFTER UPDATE ON products BEGIN
    INSERT INTO products_fts(products_fts, rowid, name, tagline, description) VALUES('delete', old.id, old.name, COALESCE(old.tagline, ''), old.description);
    INSERT INTO products_fts(rowid, name, tagline, description) VALUES (new.id, new.name, COALESCE(new.tagline, ''), new.description);
END;

CREATE TRIGGER blog_posts_ai AFTER INSERT ON blog_posts BEGIN
    INSERT INTO blog_posts_fts(rowid, title, excerpt, body) VALUES (new.id, new.title, new.excerpt, COALESCE(new.body, ''));
END;

CREATE TRIGGER blog_posts_ad AFTER DELETE ON blog_posts BEGIN
    INSERT INTO blog_posts_fts(blog_posts_fts, rowid, title, excerpt, body) VALUES('delete', old.id, old.title, old.excerpt, COALESCE(old.body, ''));
END;

CREATE TRIGGER blog_posts_au AFTER UPDATE ON blog_posts BEGIN
    INSERT INTO blog_posts_fts(blog_posts_fts, rowid, title, excerpt, body) VALUES('delete', old.id, old.title, old.excerpt, COALESCE(old.body, ''));
    INSERT INTO blog_posts_fts(rowid, title, excerpt, body) VALUES (new.id, new.title, new.excerpt, COALESCE(new.body, ''));
END;

CREATE TRIGGER case_studies_ai AFTER INSERT ON case_studies BEGIN
    INSERT INTO case_studies_fts(rowid, title, client_name, challenge_content, solution_content) VALUES (new.id, new.title, COALESCE(new.client_name, ''), COALESCE(new.challenge_content, ''), COALESCE(new.solution_content, ''));
END;

CREATE TRIGGER case_studies_ad AFTER DELETE ON case_studies BEGIN
    INSERT INTO case_studies_fts(case_studies_fts, rowid, title, client_name, challenge_content, solution_content) VALUES('delete', old.id, COALESCE(old.title, ''), COALESCE(old.client_name, ''), COALESCE(old.challenge_content, ''), COALESCE(old.solution_content, ''));
END;

CREATE TRIGGER case_studies_au AFTER UPDATE ON case_studies BEGIN
    INSERT INTO case_studies_fts(case_studies_fts, rowid, title, client_name, challenge_content, solution_content) VALUES('delete', old.id, COALESCE(old.title, ''), COALESCE(old.client_name, ''), COALESCE(old.challenge_content, ''), COALESCE(old.solution_content, ''));
    INSERT INTO case_studies_fts(rowid, title, client_name, challenge_content, solution_content) VALUES (new.id, new.title, COALESCE(new.client_name, ''), COALESCE(new.challenge_content, ''), COALESCE(new.solution_content, ''));
END;

-- Populate FTS tables with existing data
INSERT INTO products_fts(rowid, name, tagline, description) SELECT id, name, COALESCE(tagline, ''), description FROM products;
INSERT INTO blog_posts_fts(rowid, title, excerpt, body) SELECT id, title, excerpt, COALESCE(body, '') FROM blog_posts;
INSERT INTO case_studies_fts(rowid, title, client_name, challenge_content, solution_content) SELECT id, title, COALESCE(client_name, ''), COALESCE(challenge_content, ''), COALESCE(solution_content, '') FROM case_studies;
