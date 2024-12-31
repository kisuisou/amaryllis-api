package controller

import (
	"amaryllis-api/book"
	"amaryllis-api/model"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CreateUserBook(c echo.Context) error {
	isbn := c.Param("isbn")
	sess, _ := session.Get("session", c)
	user_id, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	book_data := new(model.Book)
	user_book_data := new(model.UserBooks)
	err := model.DB.Where("isbn = ?", isbn).First(book_data).Error
	err2 := model.DB.Where("isbn = ? AND user_id = ?", isbn, user_id).First(user_book_data).Error
	if err != nil {
		*book_data = book.GetMetaData(isbn)
		model.DB.Create(book_data)
	} else if err2 == nil {
		return c.NoContent(http.StatusConflict)
	}
	user_book_data.BookISBN = book_data.ISBN
	user_book_data.UserID = user_id
	model.DB.Omit("id").Create(user_book_data)
	return c.NoContent(http.StatusCreated)
}

func ReadUserBook(c echo.Context) error {
	isbn := c.Param("isbn")
	sess, _ := session.Get("session", c)
	_, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	book_data := new(model.Book)
	if err := model.DB.Where("isbn = ?", isbn).First(book_data).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, book_data)
}
