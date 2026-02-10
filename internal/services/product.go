package services

import (
	// Standard library imports for context handling
	"context" // Provides context for request cancellation and timeout handling

	// Internal application imports
	"github.com/narendhupati/bluejay-cms/db/sqlc" // Generated database query code from sqlc
)

// ProductService provides business logic for product-related operations.
// It aggregates product data from multiple related tables to provide complete
// product information including specifications, images, features, certifications,
// and downloadable resources.
type ProductService struct {
	queries *sqlc.Queries // Database query interface for executing product-related queries
}

// NewProductService creates and initializes a new ProductService instance.
// This service handles complex product data retrieval operations that require
// joining multiple related tables.
//
// Parameters:
//   - queries: Database query interface from sqlc for executing product operations
//
// Returns:
//   - *ProductService: Initialized service ready for product operations
func NewProductService(queries *sqlc.Queries) *ProductService {
	return &ProductService{queries: queries}
}

// ProductDetail is an aggregate data structure that combines a product with all
// of its related information from multiple tables. This provides a complete view
// of a product for display on product detail pages or in administrative interfaces.
//
// The structure includes the core product data along with optional collections of
// related resources. Each collection (specs, images, features, etc.) is represented
// as a slice that may be empty if no related records exist.
type ProductDetail struct {
	Product        sqlc.Product                 // Core product information (name, description, pricing, etc.)
	Category       sqlc.ProductCategory         // Product category details for navigation and organization
	Specs          []sqlc.ProductSpec           // Technical specifications (e.g., dimensions, weight, materials)
	Images         []sqlc.ProductImage          // Product images for gallery display
	Features       []sqlc.ProductFeature        // Key product features and selling points
	Certifications []sqlc.ProductCertification  // Industry certifications and compliance information
	Downloads      []sqlc.ProductDownload       // Downloadable resources (datasheets, manuals, CAD files)
}

// GetProductDetail retrieves complete product information by slug, aggregating data
// from multiple related tables into a single ProductDetail structure. This method
// performs multiple database queries to gather all product-related information.
//
// The function uses a graceful degradation strategy for optional collections: if
// fetching specs, images, features, certifications, or downloads fails, it defaults
// to an empty slice rather than failing the entire operation. This ensures that
// products can still be displayed even if some related data is missing or
// inaccessible.
//
// However, if the core product or its category cannot be retrieved, the entire
// operation fails and returns an error. This is because these are considered
// essential data required for any product display.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout handling
//   - slug: URL-friendly identifier for the product (e.g., "hydraulic-pump-model-x")
//
// Returns:
//   - *ProductDetail: Complete product information with all related data
//   - error: Non-nil if the product or category cannot be found, nil on success
func (s *ProductService) GetProductDetail(ctx context.Context, slug string) (*ProductDetail, error) {
	// Retrieve the core product record by its URL slug.
	// This is a required operation - if it fails, we cannot proceed.
	product, err := s.queries.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Retrieve the product's category information.
	// This is also required as it provides essential navigation context.
	category, err := s.queries.GetProductCategory(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	// Retrieve product specifications (e.g., dimensions, weight, materials).
	// If this fails, we default to an empty slice to allow graceful degradation.
	// The product page can still render without specs, though it's less informative.
	specs, err := s.queries.ListProductSpecs(ctx, product.ID)
	if err != nil {
		specs = []sqlc.ProductSpec{}
	}

	// Retrieve product images for the gallery.
	// Defaulting to empty slice allows the product to display even without images,
	// though a placeholder image should ideally be shown in the UI.
	images, err := s.queries.ListProductImages(ctx, product.ID)
	if err != nil {
		images = []sqlc.ProductImage{}
	}

	// Retrieve product features and selling points.
	// Empty slice default allows product display without feature highlights.
	features, err := s.queries.ListProductFeatures(ctx, product.ID)
	if err != nil {
		features = []sqlc.ProductFeature{}
	}

	// Retrieve certifications and compliance information.
	// Empty slice default is acceptable as not all products have certifications.
	certifications, err := s.queries.ListProductCertifications(ctx, product.ID)
	if err != nil {
		certifications = []sqlc.ProductCertification{}
	}

	// Retrieve downloadable resources (datasheets, manuals, CAD files).
	// Empty slice default is acceptable as downloads are optional resources.
	downloads, err := s.queries.ListProductDownloads(ctx, product.ID)
	if err != nil {
		downloads = []sqlc.ProductDownload{}
	}

	// Assemble all retrieved data into a comprehensive ProductDetail structure
	return &ProductDetail{
		Product:        product,
		Category:       category,
		Specs:          specs,
		Images:         images,
		Features:       features,
		Certifications: certifications,
		Downloads:      downloads,
	}, nil
}
