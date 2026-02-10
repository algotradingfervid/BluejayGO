package services

import (
	// Standard library imports for file handling, validation, and timestamp generation
	"fmt"            // Used for formatting error messages and generating timestamped filenames
	"io"             // Provides I/O primitives for copying uploaded files to disk
	"mime/multipart" // Handles multipart form data for file uploads from HTTP requests
	"os"             // Provides OS-level file operations (create, mkdir) for storing uploads
	"path/filepath"  // Handles cross-platform file path operations and extension extraction
	"strings"        // Used for filename sanitization and file extension validation
	"time"           // Generates Unix timestamps for unique filename prefixes
)

// UploadService handles file upload operations for product images and downloadable
// resources. It manages file validation, sanitization, storage organization, and
// provides public URLs for accessing uploaded files.
//
// The service enforces file type restrictions and size limits to prevent abuse and
// ensure storage efficiency. All uploaded files are stored with timestamped names
// to prevent filename collisions.
type UploadService struct {
	uploadDir string // Root directory for storing all uploaded files
}

// NewUploadService creates and initializes a new UploadService instance.
// The service organizes uploads into subdirectories (products, downloads) within
// the specified upload root directory.
//
// Parameters:
//   - uploadDir: Absolute path to the root directory for file uploads (e.g., "/var/www/uploads")
//
// Returns:
//   - *UploadService: Initialized service ready to handle file uploads
func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{uploadDir: uploadDir}
}

// UploadProductImage handles the upload of product image files with validation
// and storage. It enforces strict file type and size limits to ensure only valid
// image formats are accepted and storage is not abused.
//
// The function validates the file extension (jpg, jpeg, png, webp only) and size
// (max 5MB), then stores the file with a timestamped filename to prevent collisions.
// The uploaded file is saved to a "products" subdirectory within the upload root.
//
// File naming convention: {unix_timestamp}_{sanitized_original_filename}
// This ensures uniqueness while preserving some human-readable context.
//
// Parameters:
//   - file: Multipart file header from HTTP form upload
//
// Returns:
//   - string: Public URL path to the uploaded image (e.g., "/uploads/products/1234567890_image.jpg")
//   - error: Non-nil if validation fails or file cannot be saved
func (s *UploadService) UploadProductImage(file *multipart.FileHeader) (string, error) {
	// Extract and normalize the file extension for validation.
	// We convert to lowercase to handle case-insensitive comparisons.
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Validate file type: only common web image formats are allowed.
	// This prevents upload of executables, scripts, or other potentially harmful files.
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	// Enforce size limit to prevent storage abuse and ensure reasonable upload times.
	// 5MB is sufficient for high-quality product images while preventing abuse.
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("file too large (max 5MB)")
	}

	// Generate a unique filename using Unix timestamp prefix.
	// This prevents filename collisions and provides chronological ordering.
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, sanitizeFilename(file.Filename))

	// Ensure the products subdirectory exists.
	// Mode 0755 allows owner full access, group and others read/execute.
	productDir := filepath.Join(s.uploadDir, "products")
	if err := os.MkdirAll(productDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Open the uploaded file for reading
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close() // Ensure file handle is closed even if later operations fail

	// Create the destination file on disk
	dstPath := filepath.Join(productDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close() // Ensure file is properly closed and flushed

	// Copy the uploaded content to the destination file.
	// io.Copy efficiently handles the streaming transfer.
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return the public URL path for accessing the uploaded image.
	// This path should be served by a static file handler in the web server.
	return "/uploads/products/" + filename, nil
}

// UploadProductDownload handles the upload of downloadable product resources such
// as datasheets, manuals, CAD files, and other technical documents. Unlike image
// uploads, this function does not restrict file types, allowing various document
// formats to be uploaded.
//
// The function enforces a 50MB size limit to accommodate larger files like CAD
// models and detailed technical documentation while preventing excessive storage use.
// Files are stored with timestamped names in a "downloads" subdirectory.
//
// File naming convention: {unix_timestamp}_{sanitized_original_filename}
// This ensures uniqueness while preserving the original filename for user reference.
//
// Parameters:
//   - file: Multipart file header from HTTP form upload
//
// Returns:
//   - string: Public URL path to the uploaded file (e.g., "/uploads/downloads/1234567890_manual.pdf")
//   - error: Non-nil if validation fails or file cannot be saved
func (s *UploadService) UploadProductDownload(file *multipart.FileHeader) (string, error) {
	// Enforce size limit for downloadable files.
	// 50MB accommodates larger documents and CAD files while preventing abuse.
	// This is 10x larger than the image limit due to the nature of technical documents.
	if file.Size > 50*1024*1024 {
		return "", fmt.Errorf("file too large (max 50MB)")
	}

	// Generate a unique filename using Unix timestamp prefix.
	// This prevents filename collisions when multiple files with the same name are uploaded.
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, sanitizeFilename(file.Filename))

	// Ensure the downloads subdirectory exists.
	// Mode 0755 allows owner full access, group and others read/execute.
	downloadDir := filepath.Join(s.uploadDir, "downloads")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Open the uploaded file for reading
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close() // Ensure file handle is closed even if later operations fail

	// Create the destination file on disk
	dstPath := filepath.Join(downloadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close() // Ensure file is properly closed and flushed

	// Copy the uploaded content to the destination file.
	// io.Copy efficiently handles the streaming transfer without loading
	// the entire file into memory, which is important for larger files.
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Return the public URL path for accessing the uploaded file.
	// This path should be served by a static file handler in the web server.
	return "/uploads/downloads/" + filename, nil
}

// sanitizeFilename cleans uploaded filenames to ensure they are safe for filesystem
// storage and URL paths. This function prevents issues with special characters that
// could cause problems in file paths or URLs.
//
// Current implementation replaces spaces with underscores to create URL-friendly
// filenames. This is a basic sanitization approach that handles the most common
// problematic character.
//
// Note: This function could be extended to handle additional special characters
// (e.g., &, ?, #, /) or non-ASCII characters if needed for more comprehensive
// sanitization.
//
// Parameters:
//   - filename: Original filename from the uploaded file
//
// Returns:
//   - string: Sanitized filename safe for filesystem and URL use
func sanitizeFilename(filename string) string {
	// Replace spaces with underscores to create URL-friendly filenames.
	// Spaces in URLs must be encoded as %20, which is less readable and can
	// cause issues in some contexts. Underscores provide a clean alternative.
	filename = strings.ReplaceAll(filename, " ", "_")

	return filename
}
