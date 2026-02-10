# Test Plan: Public Product Detail

## Summary
Verify product detail page displays complete product information with interactive gallery and specifications accordion.

## Preconditions
- Server running on localhost:28090
- Database seeded with products including images, features, specifications, certifications, downloads, and related products
- No authentication required

## User Journey Steps
1. Navigate to GET /products/:category/:slug
2. View breadcrumb navigation and product header
3. Switch product images using gallery thumbnails
4. Watch product video if available
5. Expand/collapse specification sections using accordion
6. View certifications and downloads
7. Click related product links

## Test Cases

### Happy Path
- **Product detail page loads**: Verify GET /products/:category/:slug returns 200 status
- **Breadcrumb displays**: Verify "Home > Products > Category > Product" breadcrumb
- **Product header renders**: Verify SKU badge, tagline, and product name display
- **Gallery column displays**: Verify main product image and thumbnail images
- **Overview column displays**: Verify product overview content
- **Image switching**: Click thumbnail, verify switchImage() JS updates main image
- **Video section displays**: Verify video embed if product has video_url
- **Features section displays**: Verify checklist with icons for product features
- **Specifications accordion**: Verify specifications grouped by section
- **Accordion toggle**: Click specification section, verify toggleSpec() JS expands/collapses content
- **Certifications display**: Verify product certifications listed
- **Downloads section**: Verify downloadable files with links
- **Related products section**: Verify related products cards with links to other product detail pages

### Edge Cases / Error States
- **Product not found**: Navigate to invalid slug, verify 404 or error page
- **No video**: Verify video section hidden when video_url is null
- **No certifications**: Verify certifications section handles empty list
- **No downloads**: Verify downloads section handles empty list
- **No related products**: Verify related products section handles empty list
- **Single image**: Verify gallery works with only one image
- **Multiple accordion clicks**: Rapidly toggle specifications, verify state management

## Selectors & Elements
- Breadcrumb: text pattern "Home > Products > * > *"
- Product header: SKU badge element, tagline text, product name heading
- Two-column layout: gallery column, overview column
- Gallery: main image element, thumbnail container with clickable thumbnail images
- Video section: video embed or iframe (conditional)
- Features section: checklist container with icon elements
- Specifications: accordion container, section headers (clickable), section content (collapsible)
- Certifications: certification list container
- Downloads: download links container
- Related products: product cards grid with links to `/products/*/*`

## HTMX Interactions
- None (JavaScript-based interactions only)

## Dependencies
- JavaScript functions: switchImage(), toggleSpec()
- Template data: Category, Product, Images, Sections, Features, SpecSections
- Seeded product data with complete information
- Brutalist design: 2px solid black borders, manual box-shadows
- JetBrains Mono font
