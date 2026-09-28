package public_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/internal/handlers/public"
	"github.com/narendhupati/bluejay-cms/internal/services"
	"github.com/narendhupati/bluejay-cms/internal/testutil"
)

func TestContactSubmissionValidation(t *testing.T) {
	_, queries, cleanup := testutil.SetupTestDB(t)
	defer cleanup()
	e := echo.New()
	h := public.NewContactHandler(queries, logger, services.NewCache())
	e.POST("/contact/submit", h.SubmitContactForm)
	cases := []struct {
		email, phone string
		valid        bool
	}{
		{"name@example.com", "+91 80194 44755", true},
		{"first.last+sales@example.co.in", "+1 (212) 555-1234", true},
		{"o'connor@example.com", "020 7946 0958", true},
		{"john@example.com", "555-1234", true},
		{"not-an-email", "+91 80194 44755", false},
		{"name@localhost", "+91 80194 44755", false},
		{"a..b@example.com", "+91 80194 44755", false},
		{".name@example.com", "+91 80194 44755", false},
		{"name@example..com", "+91 80194 44755", false},
		{"name@-example.com", "+91 80194 44755", false},
		{"Name <name@example.com>", "+91 80194 44755", false},
		{"name@example.com", "abc1234567", false},
		{"name@example.com", "123", false},
		{"name@example.com", "0000000000", false},
		{"name@example.com", "1234567890123456", false},
		{"name@example.com", "91+8019444755", false},
		{"name@example.com", "+91 (80194 44755", false},
	}
	for _, tc := range cases {
		t.Run(tc.email+"_"+tc.phone, func(t *testing.T) {
			before, err := queries.CountContactSubmissions(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			form := url.Values{"name": {"Jane O'Connor"}, "email": {tc.email}, "phone": {tc.phone}, "company": {"Example Ltd"}, "message": {"Please contact us."}}
			req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			after, err := queries.CountContactSubmissions(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if tc.valid {
				if rec.Code != 200 || after != before+1 {
					t.Fatalf("valid input: status=%d count %d -> %d body=%s", rec.Code, before, after, rec.Body.String())
				}
			} else {
				if rec.Code != 400 || after != before {
					t.Fatalf("invalid input stored: status=%d count %d -> %d", rec.Code, before, after)
				}
				if !strings.Contains(rec.Body.String(), `role="alert"`) {
					t.Fatal("missing accessible validation feedback")
				}
			}
		})
	}
}
