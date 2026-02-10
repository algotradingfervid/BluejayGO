package models

import (
	"database/sql"
	"time"
)

type ProductDetail struct {
	Product        ProductView
	Category       CategoryView
	Specs          []SpecView
	Images         []ImageView
	Features       []FeatureView
	Certifications []CertificationView
	Downloads      []DownloadView
}

type ProductView struct {
	ID              int64
	SKU             string
	Slug            string
	Name            string
	Tagline         sql.NullString
	Description     string
	Overview        sql.NullString
	CategoryID      int64
	Status          string
	IsFeatured      int64
	FeaturedOrder   sql.NullInt64
	MetaTitle       sql.NullString
	MetaDescription sql.NullString
	PrimaryImage    sql.NullString
	VideoURL        sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PublishedAt     sql.NullTime
}

type CategoryView struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	Icon        string
	ImageUrl    sql.NullString
}

type SpecView struct {
	ID           int64
	ProductID    int64
	SectionName  string
	SpecKey      string
	SpecValue    string
	DisplayOrder int64
}

type ImageView struct {
	ID           int64
	ProductID    int64
	ImagePath    string
	AltText      sql.NullString
	Caption      sql.NullString
	DisplayOrder int64
	IsThumbnail  int64
}

type FeatureView struct {
	ID           int64
	ProductID    int64
	FeatureText  string
	DisplayOrder int64
}

type CertificationView struct {
	ID                int64
	ProductID         int64
	CertificationName string
	CertificationCode sql.NullString
	IconType          sql.NullString
	IconPath          sql.NullString
	DisplayOrder      int64
}

type DownloadView struct {
	ID            int64
	ProductID     int64
	Title         string
	Description   sql.NullString
	FileType      string
	FilePath      string
	FileSize      sql.NullInt64
	Version       sql.NullString
	DownloadCount int64
	DisplayOrder  int64
}

// SpecSection groups specs by section name for template rendering
type SpecSection struct {
	Name  string
	Specs []SpecView
}
