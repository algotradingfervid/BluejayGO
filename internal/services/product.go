package services

import (
	"context"

	"github.com/narendhupati/bluejay-cms/db/sqlc"
)

type ProductService struct {
	queries *sqlc.Queries
}

func NewProductService(queries *sqlc.Queries) *ProductService {
	return &ProductService{queries: queries}
}

type ProductDetail struct {
	Product        sqlc.Product
	Category       sqlc.ProductCategory
	Specs          []sqlc.ProductSpec
	Images         []sqlc.ProductImage
	Features       []sqlc.ProductFeature
	Certifications []sqlc.ProductCertification
	Downloads      []sqlc.ProductDownload
}

func (s *ProductService) GetProductDetail(ctx context.Context, slug string) (*ProductDetail, error) {
	product, err := s.queries.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	category, err := s.queries.GetProductCategory(ctx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	specs, err := s.queries.ListProductSpecs(ctx, product.ID)
	if err != nil {
		specs = []sqlc.ProductSpec{}
	}

	images, err := s.queries.ListProductImages(ctx, product.ID)
	if err != nil {
		images = []sqlc.ProductImage{}
	}

	features, err := s.queries.ListProductFeatures(ctx, product.ID)
	if err != nil {
		features = []sqlc.ProductFeature{}
	}

	certifications, err := s.queries.ListProductCertifications(ctx, product.ID)
	if err != nil {
		certifications = []sqlc.ProductCertification{}
	}

	downloads, err := s.queries.ListProductDownloads(ctx, product.ID)
	if err != nil {
		downloads = []sqlc.ProductDownload{}
	}

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
