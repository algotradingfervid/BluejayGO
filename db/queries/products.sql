-- ====================================================================
-- PRODUCTS
-- ====================================================================

-- name: CreateProduct :one
INSERT INTO products (
    sku, slug, name, tagline, description, overview,
    category_id, status, is_featured, featured_order,
    meta_title, meta_description, primary_image, video_url, published_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: GetProductBySlug :one
SELECT * FROM products WHERE slug = ? LIMIT 1;

-- name: GetProductBySKU :one
SELECT * FROM products WHERE sku = ? LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE status = 'published'
ORDER BY
    CASE WHEN is_featured = 1 THEN featured_order ELSE 999999 END ASC,
    published_at DESC
LIMIT ? OFFSET ?;

-- name: ListProductsByCategory :many
SELECT * FROM products
WHERE category_id = ? AND status = 'published'
ORDER BY
    CASE WHEN is_featured = 1 THEN featured_order ELSE 999999 END ASC,
    published_at DESC
LIMIT ? OFFSET ?;

-- name: ListFeaturedProducts :many
SELECT p.*, pc.slug AS category_slug
FROM products p
INNER JOIN product_categories pc ON p.category_id = pc.id
WHERE p.is_featured = 1 AND p.status = 'published'
ORDER BY p.featured_order ASC
LIMIT ?;

-- name: CountProducts :one
SELECT COUNT(*) FROM products WHERE status = 'published';

-- name: CountProductsByCategory :one
SELECT COUNT(*) FROM products WHERE category_id = ? AND status = 'published';

-- name: SearchProducts :many
SELECT * FROM products
WHERE status = 'published'
    AND (name LIKE ? OR description LIKE ? OR tagline LIKE ?)
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateProduct :exec
UPDATE products
SET sku = ?, slug = ?, name = ?, tagline = ?, description = ?, overview = ?,
    category_id = ?, status = ?, is_featured = ?, featured_order = ?,
    meta_title = ?, meta_description = ?, primary_image = ?, video_url = ?, published_at = ?
WHERE id = ?;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = ?;

-- name: ListAllProductsAdmin :many
SELECT * FROM products
ORDER BY created_at DESC;

-- name: ListProductsAdminFiltered :many
SELECT p.* FROM products p
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE p.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE p.category_id = @filter_category END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (p.name LIKE '%' || @filter_search || '%' OR p.sku LIKE '%' || @filter_search || '%') END)
ORDER BY p.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountProductsAdminFiltered :one
SELECT COUNT(*) FROM products p
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE p.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE p.category_id = @filter_category END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (p.name LIKE '%' || @filter_search || '%' OR p.sku LIKE '%' || @filter_search || '%') END);

-- ====================================================================
-- PRODUCT SPECS
-- ====================================================================

-- name: CreateProductSpec :one
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: ListProductSpecs :many
SELECT * FROM product_specs
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductSpecs :exec
DELETE FROM product_specs WHERE product_id = ?;

-- ====================================================================
-- PRODUCT IMAGES
-- ====================================================================

-- name: CreateProductImage :one
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListProductImages :many
SELECT * FROM product_images
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductImage :exec
DELETE FROM product_images WHERE id = ?;

-- name: DeleteProductImages :exec
DELETE FROM product_images WHERE product_id = ?;

-- ====================================================================
-- PRODUCT FEATURES
-- ====================================================================

-- name: CreateProductFeature :one
INSERT INTO product_features (product_id, feature_text, display_order)
VALUES (?, ?, ?)
RETURNING *;

-- name: ListProductFeatures :many
SELECT * FROM product_features
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductFeatures :exec
DELETE FROM product_features WHERE product_id = ?;

-- ====================================================================
-- PRODUCT CERTIFICATIONS
-- ====================================================================

-- name: CreateProductCertification :one
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, icon_path, display_order)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListProductCertifications :many
SELECT * FROM product_certifications
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductCertifications :exec
DELETE FROM product_certifications WHERE product_id = ?;

-- ====================================================================
-- PRODUCT DOWNLOADS
-- ====================================================================

-- name: CreateProductDownload :one
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetProductDownload :one
SELECT * FROM product_downloads WHERE id = ? LIMIT 1;

-- name: ListProductDownloads :many
SELECT * FROM product_downloads
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: IncrementDownloadCount :exec
UPDATE product_downloads
SET download_count = download_count + 1
WHERE id = ?;

-- name: DeleteProductDownload :exec
DELETE FROM product_downloads WHERE id = ?;

-- name: DeleteProductDownloads :exec
DELETE FROM product_downloads WHERE product_id = ?;
