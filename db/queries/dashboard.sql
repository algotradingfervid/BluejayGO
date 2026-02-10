-- name: CountNewContactSubmissions :one
SELECT COUNT(*) FROM contact_submissions WHERE status = 'new';

-- name: CountPartners :one
SELECT COUNT(*) FROM partners;

-- name: CountDraftProducts :one
SELECT COUNT(*) FROM products WHERE status = 'draft';

-- name: CountDraftBlogPosts :one
SELECT COUNT(*) FROM blog_posts WHERE status = 'draft';
