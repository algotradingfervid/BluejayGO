-- ====================================================================
-- PRODUCTS QUERY FILE
-- ====================================================================
-- This file contains all SQL queries for managing products and related entities.
--
-- Main entities:
--   - products: Core product records
--   - product_specs: Technical specifications (grouped by section)
--   - product_images: Gallery images for product detail pages
--   - product_features: Bullet-point feature lists
--   - product_certifications: Compliance badges (UL, CE, ISO, etc.)
--   - product_downloads: Datasheets, manuals, CAD files, etc.
--
-- Product status values: "published", "draft", "archived"
-- Features:
--   - Multi-status workflow (draft/published/archived)
--   - Featured products for homepage highlights
--   - Rich metadata for SEO (meta_title, meta_description)
--   - Category-based organization
--   - Search and filtering capabilities
-- ====================================================================

-- ====================================================================
-- PRODUCTS - CORE CRUD OPERATIONS
-- ====================================================================

-- name: CreateProduct :one
-- Creates a new product record with all core fields.
--
-- Parameters:
--   $1 (TEXT) - sku: Stock Keeping Unit / product code
--   $2 (TEXT) - slug: URL-safe identifier for product page
--   $3 (TEXT) - name: Product display name
--   $4 (TEXT) - tagline: Short marketing tagline
--   $5 (TEXT) - description: Full product description (HTML/Markdown)
--   $6 (TEXT) - overview: Brief product overview
--   $7 (INTEGER) - category_id: Foreign key to product_categories
--   $8 (TEXT) - status: "published", "draft", or "archived"
--   $9 (BOOLEAN) - is_featured: Whether product appears in featured listings
--   $10 (INTEGER) - featured_order: Display position for featured products
--   $11 (TEXT) - meta_title: SEO page title
--   $12 (TEXT) - meta_description: SEO meta description
--   $13 (TEXT) - primary_image: Main product image path
--   $14 (TEXT) - video_url: Product demo/promo video URL (optional)
--   $15 (TIMESTAMP) - published_at: Publication date/time
--
-- Returns: Product - The newly created product with auto-generated ID and timestamps
--
-- Note: Related entities (specs, images, features) are created separately after product creation
INSERT INTO products (
    sku, slug, name, tagline, description, overview,
    category_id, status, is_featured, featured_order,
    meta_title, meta_description, primary_image, video_url, published_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetProduct :one
-- Retrieves a single product by its primary key ID (all statuses).
--
-- Parameters:
--   $1 (INTEGER) - product ID
-- Returns: Product - Single product record or error if not found
--
-- Use case: Admin editing, fetching product for update regardless of status
-- Note: Does NOT filter by status, returns draft/archived products
SELECT * FROM products WHERE id = ? LIMIT 1;

-- name: GetProductBySlug :one
-- Retrieves a single product by its URL-safe slug (all statuses).
--
-- Parameters:
--   $1 (TEXT) - product slug (e.g., "industrial-sensor-x200")
-- Returns: Product - Single product record or error if not found
--
-- Use case: Frontend product detail page, preview mode
-- Note: Does NOT filter by status; application should check status before displaying
SELECT * FROM products WHERE slug = ? LIMIT 1;

-- name: GetProductBySKU :one
-- Retrieves a single product by its SKU/product code.
--
-- Parameters:
--   $1 (TEXT) - sku: Stock Keeping Unit (e.g., "PROD-12345")
-- Returns: Product - Single product record or error if not found
--
-- Use case: SKU lookup, validation during data import, order processing
-- Note: SKUs should be unique (enforced by database constraint)
SELECT * FROM products WHERE sku = ? LIMIT 1;

-- ====================================================================
-- PRODUCTS - PUBLIC LISTING QUERIES
-- ====================================================================

-- name: ListProducts :many
-- Retrieves paginated published products with featured products first.
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Number of products per page
--   $2 (INTEGER) - OFFSET: Pagination offset
-- Returns: []Product - Array of published products
--
-- Filtering: status = 'published' - Only live, public products
--
-- Sorting logic (complex CASE statement):
--   1. Featured products (is_featured = 1) ordered by featured_order (1, 2, 3...)
--   2. Non-featured products ordered by published_at DESC (newest first)
--   - CASE WHEN is_featured = 1 THEN featured_order ELSE 999999 END
--     This ensures featured products (orders 1-100) appear before non-featured (999999)
--
-- Use case: Main products catalog page, products listing
SELECT * FROM products
WHERE status = 'published'
ORDER BY
    CASE WHEN is_featured = 1 THEN featured_order ELSE 999999 END ASC,
    published_at DESC
LIMIT ? OFFSET ?;

-- name: ListProductsByCategory :many
-- Retrieves paginated published products for a specific category.
--
-- Parameters:
--   $1 (INTEGER) - category_id: Filter by this product category
--   $2 (INTEGER) - LIMIT: Number of products per page
--   $3 (INTEGER) - OFFSET: Pagination offset
-- Returns: []Product - Array of published products in the category
--
-- Filtering:
--   - category_id = ? - Specific category only
--   - status = 'published' - Only public products
--
-- Sorting: Same complex ordering as ListProducts (featured first, then by date)
--
-- Use case: Category-specific product listing pages
SELECT * FROM products
WHERE category_id = ? AND status = 'published'
ORDER BY
    CASE WHEN is_featured = 1 THEN featured_order ELSE 999999 END ASC,
    published_at DESC
LIMIT ? OFFSET ?;

-- name: ListFeaturedProducts :many
-- Retrieves a limited number of featured products with category slug.
--
-- Parameters:
--   $1 (INTEGER) - LIMIT: Maximum number of featured products to return
-- Returns: []Product - Array of featured products with category_slug
--
-- JOIN logic:
--   - INNER JOIN product_categories pc ON p.category_id = pc.id
--     Adds category_slug for building product URLs
--
-- Filtering:
--   - p.is_featured = 1 - Only featured products
--   - p.status = 'published' - Only public products
--
-- Sorting: p.featured_order ASC - Featured products in custom order (1, 2, 3...)
-- LIMIT: Controls homepage featured products count (e.g., 4 or 6)
--
-- Use case: Homepage featured products section, product highlights carousel
SELECT p.*, pc.slug AS category_slug
FROM products p
INNER JOIN product_categories pc ON p.category_id = pc.id
WHERE p.is_featured = 1 AND p.status = 'published'
ORDER BY p.featured_order ASC
LIMIT ?;

-- name: CountProducts :one
-- Returns the total count of published products.
--
-- Parameters: none
-- Returns: INTEGER - Total number of published products
--
-- Filtering: status = 'published' - Only counts public products
-- Use case: Pagination calculations, site statistics
SELECT COUNT(*) FROM products WHERE status = 'published';

-- name: CountProductsByCategory :one
-- Returns the count of published products in a specific category.
--
-- Parameters:
--   $1 (INTEGER) - category_id: Category to count products in
-- Returns: INTEGER - Number of published products in the category
--
-- Use case: Category page pagination, category statistics
SELECT COUNT(*) FROM products WHERE category_id = ? AND status = 'published';

-- name: SearchProducts :many
-- Searches published products by name, description, or tagline.
--
-- Parameters:
--   $1 (TEXT) - Search term for name (LIKE pattern, e.g., "%sensor%")
--   $2 (TEXT) - Search term for description (same pattern as $1)
--   $3 (TEXT) - Search term for tagline (same pattern as $1)
--   $4 (INTEGER) - LIMIT: Results per page
--   $5 (INTEGER) - OFFSET: Pagination offset
-- Returns: []Product - Array of matching published products
--
-- Search logic: OR condition across three fields (name, description, tagline)
-- Note: Caller must add wildcards (%) to search term for partial matching
-- Example: SearchProducts("%industrial%", "%industrial%", "%industrial%", 20, 0)
--
-- Performance: May be slow without full-text search index on large datasets
-- Sorting: published_at DESC - Newest matching products first
SELECT * FROM products
WHERE status = 'published'
    AND (name LIKE ? OR description LIKE ? OR tagline LIKE ?)
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateProduct :exec
-- Updates all core fields of an existing product.
--
-- Parameters:
--   $1-$15 - Same as CreateProduct parameters (sku through published_at)
--   $16 (INTEGER) - id: Product ID to update
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Note: This updates the main product record only, not related entities
-- Separate queries handle specs, images, features, certifications, downloads
UPDATE products
SET sku = ?, slug = ?, name = ?, tagline = ?, description = ?, overview = ?,
    category_id = ?, status = ?, is_featured = ?, featured_order = ?,
    meta_title = ?, meta_description = ?, primary_image = ?, video_url = ?, published_at = ?
WHERE id = ?;

-- name: DeleteProduct :exec
-- Permanently deletes a product record.
--
-- Parameters:
--   $1 (INTEGER) - product ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- WARNING: Should cascade delete related records (specs, images, features, etc.)
-- Note: Ensure foreign key constraints are configured for CASCADE DELETE
-- Alternative: Set status='archived' for soft delete instead
DELETE FROM products WHERE id = ?;

-- ====================================================================
-- PRODUCTS - ADMIN QUERIES
-- ====================================================================

-- name: ListAllProductsAdmin :many
-- Retrieves all products (any status) for admin dashboard.
--
-- Parameters: none
-- Returns: []Product - Array of all products including drafts and archived
--
-- Sorting: created_at DESC - Newest products first
-- Use case: Admin product management dashboard
-- Note: No status filtering - shows published, draft, and archived products
SELECT * FROM products
ORDER BY created_at DESC;

-- name: ListProductsAdminFiltered :many
-- Retrieves paginated products with optional filters for admin interface.
--
-- Parameters (named parameters with @):
--   @filter_status (TEXT) - Filter by status ("published", "draft", "archived", or "" for all)
--   @filter_category (INTEGER) - Filter by category_id (0 for all categories)
--   @filter_search (TEXT) - Search term for name/SKU (empty string for no search)
--   @page_limit (INTEGER) - Results per page
--   @page_offset (INTEGER) - Pagination offset
-- Returns: []Product - Array of filtered products
--
-- Complex WHERE clause with CASE statements for optional filters:
--   1. Status filter: CASE WHEN @filter_status = '' THEN 1 (all) ELSE match status
--   2. Category filter: CASE WHEN @filter_category = 0 THEN 1 (all) ELSE match category_id
--   3. Search filter: CASE WHEN @filter_search = '' THEN 1 (all) ELSE LIKE match on name/SKU
--
-- The CASE WHEN pattern allows filters to be "turned off" by passing default values
-- ("" for strings, 0 for category ID)
--
-- Sorting: p.created_at DESC - Newest products first
-- Use case: Admin product listing with filter dropdowns and search bar
SELECT p.* FROM products p
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE p.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE p.category_id = @filter_category END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (p.name LIKE '%' || @filter_search || '%' OR p.sku LIKE '%' || @filter_search || '%') END)
ORDER BY p.created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountProductsAdminFiltered :one
-- Returns count of products matching admin filters (for pagination).
--
-- Parameters: Same as ListProductsAdminFiltered (@filter_status, @filter_category, @filter_search)
-- Returns: INTEGER - Count of products matching filter criteria
--
-- Note: Uses identical WHERE clause as ListProductsAdminFiltered for consistent counts
-- Use case: Calculating total pages for filtered admin product listing
SELECT COUNT(*) FROM products p
WHERE
    (CASE WHEN @filter_status = '' THEN 1 ELSE p.status = @filter_status END)
    AND (CASE WHEN @filter_category = 0 THEN 1 ELSE p.category_id = @filter_category END)
    AND (CASE WHEN @filter_search = '' THEN 1 ELSE (p.name LIKE '%' || @filter_search || '%' OR p.sku LIKE '%' || @filter_search || '%') END);

-- ====================================================================
-- PRODUCT SPECS (Technical Specifications)
-- ====================================================================

-- name: CreateProductSpec :one
-- Creates a single technical specification for a product.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Foreign key to parent product
--   $2 (TEXT) - section_name: Grouping label (e.g., "Electrical", "Physical", "Environmental")
--   $3 (TEXT) - spec_key: Specification label (e.g., "Voltage", "Weight", "Operating Temp")
--   $4 (TEXT) - spec_value: Specification value (e.g., "24V DC", "2.5 kg", "-40°C to 85°C")
--   $5 (INTEGER) - display_order: Position within product specs list
--
-- Returns: ProductSpec - The newly created spec with auto-generated ID
--
-- Use case: Adding technical specifications during product creation/editing
-- Note: Specs can be grouped by section_name for tabbed or sectioned display
INSERT INTO product_specs (product_id, section_name, spec_key, spec_value, display_order)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: ListProductSpecs :many
-- Retrieves all technical specifications for a product in display order.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product to fetch specs for
-- Returns: []ProductSpec - Array of specifications ordered by display_order
--
-- Sorting: display_order ASC - Specs appear in admin-configured order
-- Use case: Displaying specs table on product detail page
-- Note: Application code should group by section_name for organized display
SELECT * FROM product_specs
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductSpecs :exec
-- Deletes all technical specifications for a product (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product whose specs to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Clearing all specs before re-importing or rebuilding spec list
-- WARNING: Deletes ALL specs for the product in one operation
DELETE FROM product_specs WHERE product_id = ?;

-- ====================================================================
-- PRODUCT IMAGES (Gallery Images)
-- ====================================================================

-- name: CreateProductImage :one
-- Adds a gallery image to a product.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Foreign key to parent product
--   $2 (TEXT) - image_path: Path to image file
--   $3 (TEXT) - alt_text: Accessibility alt text for screen readers
--   $4 (TEXT) - caption: Optional image caption for display
--   $5 (INTEGER) - display_order: Position in image gallery
--   $6 (BOOLEAN) - is_thumbnail: Whether this image is the thumbnail (usually first image)
--
-- Returns: ProductImage - The newly created image record with auto-generated ID
--
-- Use case: Building product image gallery during product creation/editing
-- Note: Typically only one image should have is_thumbnail=1 per product
INSERT INTO product_images (product_id, image_path, alt_text, caption, display_order, is_thumbnail)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListProductImages :many
-- Retrieves all gallery images for a product in display order.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product to fetch images for
-- Returns: []ProductImage - Array of images ordered by display_order
--
-- Sorting: display_order ASC - Images appear in admin-configured order
-- Use case: Rendering product image gallery, lightbox, thumbnails
SELECT * FROM product_images
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductImage :exec
-- Deletes a single product gallery image.
--
-- Parameters:
--   $1 (INTEGER) - image ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Removing individual images from product gallery
-- WARNING: Only deletes database record; application should delete physical file
DELETE FROM product_images WHERE id = ?;

-- name: DeleteProductImages :exec
-- Deletes all gallery images for a product (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product whose images to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Clearing all images before re-importing or deleting product
-- WARNING: Deletes ALL images for the product; physical files should also be removed
DELETE FROM product_images WHERE product_id = ?;

-- ====================================================================
-- PRODUCT FEATURES (Bullet-Point Feature Lists)
-- ====================================================================

-- name: CreateProductFeature :one
-- Adds a single feature bullet point to a product.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Foreign key to parent product
--   $2 (TEXT) - feature_text: Feature description (e.g., "IP67 waterproof rating")
--   $3 (INTEGER) - display_order: Position in features list
--
-- Returns: ProductFeature - The newly created feature with auto-generated ID
--
-- Use case: Building key features list during product creation/editing
-- Note: Features are typically short, marketing-focused bullet points
INSERT INTO product_features (product_id, feature_text, display_order)
VALUES (?, ?, ?)
RETURNING *;

-- name: ListProductFeatures :many
-- Retrieves all features for a product in display order.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product to fetch features for
-- Returns: []ProductFeature - Array of features ordered by display_order
--
-- Sorting: display_order ASC - Features appear in admin-configured order
-- Use case: Displaying key features list on product detail page
SELECT * FROM product_features
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductFeatures :exec
-- Deletes all features for a product (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product whose features to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Clearing all features before re-importing or rebuilding features list
-- WARNING: Deletes ALL features for the product in one operation
DELETE FROM product_features WHERE product_id = ?;

-- ====================================================================
-- PRODUCT CERTIFICATIONS (Compliance Badges & Standards)
-- ====================================================================

-- name: CreateProductCertification :one
-- Adds a certification/compliance badge to a product.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Foreign key to parent product
--   $2 (TEXT) - certification_name: Display name (e.g., "UL Listed", "CE Certified")
--   $3 (TEXT) - certification_code: Official code/number (e.g., "UL 508", "CE 2014/35/EU")
--   $4 (TEXT) - icon_type: Icon source type ("upload", "library", "font-icon")
--   $5 (TEXT) - icon_path: Path to icon file or icon class name
--   $6 (INTEGER) - display_order: Position in certifications list
--
-- Returns: ProductCertification - The newly created certification with auto-generated ID
--
-- Use case: Adding compliance badges during product creation/editing
-- Note: Common certifications include UL, CE, FCC, RoHS, ISO, CSA
INSERT INTO product_certifications (product_id, certification_name, certification_code, icon_type, icon_path, display_order)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListProductCertifications :many
-- Retrieves all certifications for a product in display order.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product to fetch certifications for
-- Returns: []ProductCertification - Array of certifications ordered by display_order
--
-- Sorting: display_order ASC - Certifications appear in admin-configured order
-- Use case: Displaying compliance badges on product detail page
SELECT * FROM product_certifications
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: DeleteProductCertifications :exec
-- Deletes all certifications for a product (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product whose certifications to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Clearing all certifications before re-importing or deleting product
-- WARNING: Deletes ALL certifications for the product in one operation
DELETE FROM product_certifications WHERE product_id = ?;

-- ====================================================================
-- PRODUCT DOWNLOADS (Datasheets, Manuals, CAD Files, etc.)
-- ====================================================================

-- name: CreateProductDownload :one
-- Adds a downloadable file (datasheet, manual, CAD, etc.) to a product.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Foreign key to parent product
--   $2 (TEXT) - title: Download display name (e.g., "Product Datasheet", "User Manual")
--   $3 (TEXT) - description: Download description or notes (optional)
--   $4 (TEXT) - file_type: File type label (e.g., "PDF", "DWG", "STEP", "ZIP")
--   $5 (TEXT) - file_path: Path to downloadable file
--   $6 (INTEGER) - file_size: File size in bytes
--   $7 (TEXT) - version: Document/file version (e.g., "v2.1", "Rev A")
--   $8 (INTEGER) - display_order: Position in downloads list
--
-- Returns: ProductDownload - The newly created download with auto-generated ID
--
-- Use case: Adding technical documents during product creation/editing
-- Note: download_count initializes to 0 via database schema default
INSERT INTO product_downloads (product_id, title, description, file_type, file_path, file_size, version, display_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetProductDownload :one
-- Retrieves a single product download by its ID.
--
-- Parameters:
--   $1 (INTEGER) - download ID
-- Returns: ProductDownload - Single download record or error if not found
--
-- Use case: Fetching download metadata before serving file, tracking analytics
SELECT * FROM product_downloads WHERE id = ? LIMIT 1;

-- name: ListProductDownloads :many
-- Retrieves all downloadable files for a product in display order.
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product to fetch downloads for
-- Returns: []ProductDownload - Array of downloads ordered by display_order
--
-- Sorting: display_order ASC - Downloads appear in admin-configured order
-- Use case: Displaying downloads list on product detail page
SELECT * FROM product_downloads
WHERE product_id = ?
ORDER BY display_order ASC;

-- name: IncrementDownloadCount :exec
-- Increments the download counter for analytics tracking.
--
-- Parameters:
--   $1 (INTEGER) - download ID
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Tracking download analytics when user downloads a file
-- Note: Uses download_count + 1 to atomically increment without race conditions
UPDATE product_downloads
SET download_count = download_count + 1
WHERE id = ?;

-- name: DeleteProductDownload :exec
-- Deletes a single product download.
--
-- Parameters:
--   $1 (INTEGER) - download ID to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Removing individual downloads from product
-- WARNING: Only deletes database record; application should delete physical file
DELETE FROM product_downloads WHERE id = ?;

-- name: DeleteProductDownloads :exec
-- Deletes all downloads for a product (bulk delete).
--
-- Parameters:
--   $1 (INTEGER) - product_id: Product whose downloads to delete
-- Returns: (none) - sqlc annotation :exec returns only row count
--
-- Use case: Clearing all downloads before re-importing or deleting product
-- WARNING: Deletes ALL downloads for the product; physical files should also be removed
DELETE FROM product_downloads WHERE product_id = ?;
