package services_test

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/narendhupati/bluejay-cms/internal/services"
)

func createMultipartFileHeader(t *testing.T, filename string, content []byte, contentType string) *multipart.FileHeader {
	t.Helper()
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", contentType)

	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("CreatePart: %v", err)
	}
	part.Write(content)
	w.Close()

	r := multipart.NewReader(&b, w.Boundary())
	form, err := r.ReadForm(32 << 20)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("no files in form")
	}
	return files[0]
}

func TestUploadProductImage_ValidJPG(t *testing.T) {
	tmpDir := t.TempDir()
	svc := services.NewUploadService(tmpDir)

	fh := createMultipartFileHeader(t, "test.jpg", []byte("fake-jpg-data"), "image/jpeg")
	path, err := svc.UploadProductImage(fh)
	if err != nil {
		t.Fatalf("UploadProductImage: %v", err)
	}

	if !strings.HasPrefix(path, "/uploads/products/") {
		t.Errorf("expected path prefix /uploads/products/, got %q", path)
	}
	if !strings.HasSuffix(path, "_test.jpg") {
		t.Errorf("expected path to end with _test.jpg, got %q", path)
	}

	// Verify file exists on disk
	localPath := filepath.Join(tmpDir, "products", filepath.Base(path))
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		t.Errorf("uploaded file does not exist: %s", localPath)
	}
}

func TestUploadProductImage_ValidPNG(t *testing.T) {
	tmpDir := t.TempDir()
	svc := services.NewUploadService(tmpDir)

	fh := createMultipartFileHeader(t, "image.png", []byte("fake-png"), "image/png")
	path, err := svc.UploadProductImage(fh)
	if err != nil {
		t.Fatalf("UploadProductImage: %v", err)
	}
	if !strings.HasSuffix(path, "_image.png") {
		t.Errorf("unexpected path: %q", path)
	}
}

func TestUploadProductImage_ValidWebP(t *testing.T) {
	tmpDir := t.TempDir()
	svc := services.NewUploadService(tmpDir)

	fh := createMultipartFileHeader(t, "photo.webp", []byte("fake-webp"), "image/webp")
	_, err := svc.UploadProductImage(fh)
	if err != nil {
		t.Fatalf("UploadProductImage: %v", err)
	}
}

func TestUploadProductImage_InvalidType(t *testing.T) {
	tmpDir := t.TempDir()
	svc := services.NewUploadService(tmpDir)

	fh := createMultipartFileHeader(t, "doc.pdf", []byte("fake-pdf"), "application/pdf")
	_, err := svc.UploadProductImage(fh)
	if err == nil {
		t.Fatal("expected error for invalid file type")
	}
	if !strings.Contains(err.Error(), "invalid file type") {
		t.Errorf("expected 'invalid file type' error, got: %v", err)
	}
}

func TestUploadProductImage_TooLarge(t *testing.T) {
	tmpDir := t.TempDir()
	svc := services.NewUploadService(tmpDir)

	// Create a file header that reports > 5MB
	fh := createMultipartFileHeader(t, "big.jpg", make([]byte, 6*1024*1024), "image/jpeg")
	_, err := svc.UploadProductImage(fh)
	if err == nil {
		t.Fatal("expected error for large file")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Errorf("expected 'too large' error, got: %v", err)
	}
}
