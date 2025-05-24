package controller

import (
	"amaryllis-api/book"
	"amaryllis-api/model"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func ReadBook(c echo.Context) error {
	isbn := c.Param("isbn")
	sess, _ := session.Get("session", c)
	_, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	book_data := new(model.Book)
	if err := model.DB.Where("isbn = ?", isbn).First(book_data).Error; err != nil {
		*book_data = book.GetMetaData(isbn)
		model.DB.Create(book_data)
	}
	return c.JSON(http.StatusOK, book_data)
}
