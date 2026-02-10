ALTER TABLE contact_submissions ADD COLUMN submission_type TEXT NOT NULL DEFAULT 'contact';
CREATE INDEX idx_contact_submissions_type ON contact_submissions(submission_type);
