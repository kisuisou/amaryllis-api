package controller

import (
	"amaryllis-api/book"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func ReadBookImg(c echo.Context) error {
	isbn := c.Param("isbn")
	sess, _ := session.Get("session", c)
	_, is_ok := sess.Values["UserID"].(string)
	if !FindImg(isbn) {
		if is_ok {
			book.GetBookImg(isbn)
		} else {
			return c.NoContent(http.StatusForbidden)
		}
	}
	return c.File(fmt.Sprintf("./book_imgs/%s.jpg", isbn))
}

func FindImg(isbn string) bool {
	img_names, _ := os.ReadDir("./book_imgs")
	is_img_exist := false
	for _, img_name := range img_names {
		if (isbn + ".jpg") == img_name.Name() {
			is_img_exist = true
			break
		}
	}
	return is_img_exist
}
