package controller

import (
	"amaryllis-api/model"
	"net/http"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type create_session_req struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

func createSession(c echo.Context) error {
	env := os.Getenv("ENV")
	is_cookie_secure := env == "production"
	r := new(create_session_req)
	if err := c.Bind(r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	user := new(model.User)
	err := model.DB.Where("UserID = ?", r.UserID).First(user).Error
	if err != nil {
		return c.NoContent(http.StatusForbidden)
	}
	is_pw_collect, err := argon2id.ComparePasswordAndHash(r.Password, user.PasswordHash)
	if is_pw_collect {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   604800,
			Secure:   is_cookie_secure,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		sess.Values["UserId"] = user.UserID
		sess.Save(c.Request(), c.Response())
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusForbidden)
}
