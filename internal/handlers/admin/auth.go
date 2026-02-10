package admin

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	customMiddleware "github.com/narendhupati/bluejay-cms/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewAuthHandler(queries *sqlc.Queries, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *AuthHandler) ShowLoginPage(c echo.Context) error {
	sess, _ := c.Get("session").(*customMiddleware.Session)
	if sess != nil && sess.UserID > 0 {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}

	return c.Render(http.StatusOK, "admin/pages/login.html", map[string]interface{}{
		"Title": "Admin Login",
		"Error": c.QueryParam("error"),
	})
}

func (h *AuthHandler) LoginSubmit(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=missing_fields")
	}

	user, err := h.queries.GetAdminUserByEmail(c.Request().Context(), email)
	if err != nil {
		h.logger.Warn("login attempt failed", "email", email, "error", "user_not_found")
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=invalid_credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		h.logger.Warn("login attempt failed", "email", email, "error", "invalid_password")
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=invalid_credentials")
	}

	if err := h.queries.UpdateLastLogin(c.Request().Context(), user.ID); err != nil {
		h.logger.Error("failed to update last login", "user_id", user.ID, "error", err)
	}

	sess := c.Get("session").(*customMiddleware.Session)
	sess.UserID = user.ID
	sess.Email = user.Email
	sess.DisplayName = user.DisplayName
	sess.Role = user.Role

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.logger.Error("failed to save session", "error", err)
		return c.Redirect(http.StatusSeeOther, "/admin/login?error=session_error")
	}

	h.logger.Info("user logged in", "user_id", user.ID, "email", user.Email)
	logActivity(c, "login", "system", user.ID, user.DisplayName, "User '%s' logged in", user.DisplayName)
	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	sess := c.Get("session").(*customMiddleware.Session)

	sess.UserID = 0
	sess.Email = ""
	sess.DisplayName = ""
	sess.Role = ""
	sess.Options.MaxAge = -1

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		h.logger.Error("failed to destroy session", "error", err)
	}

	h.logger.Info("user logged out")
	logActivity(c, "logout", "system", 0, "", "User logged out")
	return c.Redirect(http.StatusSeeOther, "/admin/login")
}
