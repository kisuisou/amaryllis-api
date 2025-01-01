package controller

import (
	"amaryllis-api/book"
	"amaryllis-api/model"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type create_user_book_req struct {
	ISBN string `json:"isbn"`
}

func CreateUserBook(c echo.Context) error {
	req := new(create_user_book_req)
	if err := c.Bind(req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sess, _ := session.Get("session", c)
	user_id, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	book_data := new(model.Book)
	user_book_data := new(model.UserBooks)
	err := model.DB.Where("isbn = ?", req.ISBN).First(book_data).Error
	err2 := model.DB.Where("book_isbn = ? AND user_id = ?", req.ISBN, user_id).First(user_book_data).Error
	if err != nil {
		*book_data = book.GetMetaData(req.ISBN)
		model.DB.Create(book_data)
	} else if err2 == nil {
		return c.NoContent(http.StatusConflict)
	}
	user_book_data.BookISBN = book_data.ISBN
	user_book_data.UserID = user_id
	model.DB.Omit("id").Create(user_book_data)
	return c.NoContent(http.StatusCreated)
}

func ReadUserBooks(c echo.Context) error {
	user_id := c.Param("user_id")
	sess, _ := session.Get("session", c)
	_, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	user := new(model.User)
	model.DB.Preload("Books").Where("id = ?", user_id).First(user)
	user_books := user.Books
	fmt.Println(user_books)
	return c.JSON(http.StatusOK, user_books)
}
