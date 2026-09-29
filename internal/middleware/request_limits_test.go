package middleware_test

import (
	"bytes"
	"github.com/labstack/echo/v4"
	m "github.com/narendhupati/bluejay-cms/internal/middleware"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLimits(t *testing.T) {
	for _, tc := range []struct {
		name, body, kind string
		chunked          bool
		want             int
	}{
		{"valid", "name=hello", "application/x-www-form-urlencoded", false, 204},
		{"many values", strings.Repeat("key=value&", 513), "application/x-www-form-urlencoded", false, 413},
		{"declared too large", strings.Repeat("x", m.FormBodyLimit+1), "application/x-www-form-urlencoded", false, 413},
		{"chunked too large", strings.Repeat("x", m.FormBodyLimit+1), "application/x-www-form-urlencoded", true, 413},
		{"malformed", "x=%GG", "application/x-www-form-urlencoded", false, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Use(m.RequestLimits())
			e.POST("/", func(c echo.Context) error { return c.NoContent(204) })
			r := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", tc.kind)
			if tc.chunked {
				r.ContentLength = -1
			}
			w := httptest.NewRecorder()
			e.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("got%d want%d", w.Code, tc.want)
			}
		})
	}
}
func TestMultipartLimitsPreserveValidFile(t *testing.T) {
	for _, count := range []int{1, 33} {
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		for i := 0; i < count; i++ {
			f, _ := writer.CreateFormFile("files", "test.txt")
			io.WriteString(f, "fixture")
		}
		writer.Close()
		e := echo.New()
		e.Use(m.RequestLimits())
		e.POST("/", func(c echo.Context) error {
			f, err := c.FormFile("files")
			if err != nil {
				return err
			}
			if f.Size != 7 {
				t.Fatal("upload changed")
			}
			return c.NoContent(http.StatusNoContent)
		})
		r := httptest.NewRequest("POST", "/", &b)
		r.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		e.ServeHTTP(w, r)
		want := 204
		if count == 33 {
			want = 413
		}
		if w.Code != want {
			t.Fatalf("files%d got%d want%d", count, w.Code, want)
		}
	}
}
