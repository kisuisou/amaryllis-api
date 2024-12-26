package controller

import (
	"amaryllis-api/book"
	"amaryllis-api/model"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type create_user_book_req struct {
	ISBN string `json:"isbn"`
}

func CreateUserBook(c echo.Context) error {
	sess, _ := session.Get("session", c)
	user_id, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	r := new(create_user_book_req)
	if err := c.Bind(r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	book_data := new(model.Book)
	err := model.DB.Where("isbn = ?", r.ISBN).First(book_data).Error
	if err != nil {
		*book_data = book.GetMetaData(r.ISBN)
		model.DB.Create(book_data)
	}
	user_book_data := new(model.UserBooks)
	user_book_data.BookISBN = book_data.ISBN
	user_book_data.UserID = user_id
	model.DB.Omit("id").Create(user_book_data)
	return c.NoContent(http.StatusCreated)
}
