package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Session struct {
	*sessions.Session
	UserID      int64
	Email       string
	DisplayName string
	Role        string
}

var SessionStore *sessions.CookieStore

func InitSessionStore(secret string) {
	SessionStore = sessions.NewCookieStore([]byte(secret))
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := SessionStore.Get(c.Request(), "bluejay_session")
			if err != nil {
				session, _ = SessionStore.New(c.Request(), "bluejay_session")
			}

			sess := &Session{Session: session}
			if userID, ok := session.Values["user_id"].(int64); ok {
				sess.UserID = userID
			}
			if email, ok := session.Values["email"].(string); ok {
				sess.Email = email
			}
			if displayName, ok := session.Values["display_name"].(string); ok {
				sess.DisplayName = displayName
			}
			if role, ok := session.Values["role"].(string); ok {
				sess.Role = role
			}

			c.Set("session", sess)
			return next(c)
		}
	}
}

func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	s.Values["user_id"] = s.UserID
	s.Values["email"] = s.Email
	s.Values["display_name"] = s.DisplayName
	s.Values["role"] = s.Role
	return s.Session.Save(r, w)
}
