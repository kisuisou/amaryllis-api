package controller

import (
	"amaryllis-api/model"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type create_user_book_req struct {
	BookID uint `json:"book_id"`
	IsRead bool `json:"is_read"`
}

type user_books_res struct {
	MetaData  model.Book
	ISBN      string
	IsRead    bool
	CreatedAt uint
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
	if req.BookID == 0 {
		return c.NoContent(http.StatusBadRequest)
	}
	book_data := new(model.Book)
	if err := model.DB.First(book_data, req.BookID).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	user_book_data := new(model.UserBooks)
	if err := model.DB.Where("book_id = ? AND user_id = ?", book_data.ID, user_id).First(user_book_data).Error; err == nil {
		return c.NoContent(http.StatusConflict)
	}
	user_book_data.BookID = book_data.ID
	user_book_data.UserID = user_id
	user_book_data.IsRead = req.IsRead
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
	user_books := new([]model.UserBooks)
	var response []user_books_res
	model.DB.Where("user_id = ?", user_id).Find(user_books)
	for _, user_book := range *user_books {
		var res user_books_res
		model.DB.First(&res.MetaData, user_book.BookID)
		res.ISBN = primaryISBN(user_book.BookID)
		res.CreatedAt = uint(user_book.CreatedAt.Unix())
		res.IsRead = user_book.IsRead
		response = append(response, res)
	}
	return c.JSON(http.StatusOK, response)
}
