package controller

import (
	"amaryllis-api/book"
	"amaryllis-api/model"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func ReadBookImg(c echo.Context) error {
	isbn := c.Param("isbn")
	sess, _ := session.Get("session", c)
	_, is_ok := sess.Values["UserID"].(string)
	book_data := new(model.Book)
	err := model.DB.Where("isbn = ?", isbn).First(book_data).Error
	if err != nil || book_data.Image == "Failed" {
		return c.NoContent(http.StatusNotFound)
	}
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	} else if book_data.Image == "" {
		if !book.GetBookImg(isbn) {
			book_data.Image = "Failed"
			model.DB.Save(book_data)
			return c.NoContent(http.StatusNotFound)
		} else {
			book_data.Image = "Success"
			model.DB.Save(book_data)
		}
	}
	return c.File(fmt.Sprintf("./book_imgs/%s.jpg", isbn))
}
