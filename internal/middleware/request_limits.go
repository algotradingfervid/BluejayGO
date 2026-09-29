package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"mime"
	"net/http"
)

const (
	FormBodyLimit   = 4 << 20
	UploadBodyLimit = 64 << 20
	MaxFormValues   = 512
	MaxUploadFiles  = 32
)

// RequestLimits caps bodies before parsing, including chunked requests, and
// bounds multipart/text form complexity before handlers allocate work or write.
func RequestLimits() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			limit := int64(FormBodyLimit)
			if mediaType == "multipart/form-data" {
				limit = UploadBodyLimit
			}
			if r.ContentLength > limit {
				return echo.NewHTTPError(http.StatusRequestEntityTooLarge)
			}
			r.Body = http.MaxBytesReader(c.Response(), r.Body, limit)
			var err error
			if mediaType == "multipart/form-data" {
				err = r.ParseMultipartForm(8 << 20)
				if r.MultipartForm != nil {
					defer r.MultipartForm.RemoveAll()
				}
			} else if mediaType == "application/x-www-form-urlencoded" {
				err = r.ParseForm()
			}
			if err != nil {
				var tooLarge *http.MaxBytesError
				if errors.As(err, &tooLarge) {
					return echo.NewHTTPError(http.StatusRequestEntityTooLarge)
				}
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid form")
			}
			count := 0
			for _, values := range r.Form {
				count += len(values)
			}
			if count > MaxFormValues {
				return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "Too many form values")
			}
			if r.MultipartForm != nil {
				files := 0
				for _, entries := range r.MultipartForm.File {
					files += len(entries)
				}
				if files > MaxUploadFiles {
					return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "Too many upload files")
				}
			}
			return next(c)
		}
	}
}
