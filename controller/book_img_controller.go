package controller

import (
	"amaryllis-api/book"
	"amaryllis-api/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func ReadBookImg(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sess, _ := session.Get("session", c)
	_, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	book_data := new(model.Book)
	err = model.DB.First(book_data, id).Error
	if err != nil || book_data.ImageStatus != "Success" {
		return c.NoContent(http.StatusNotFound)
	}
	isbn := primaryISBN(book_data.ID)
	if isbn == "" {
		return c.NoContent(http.StatusNotFound)
	}
	return c.File(fmt.Sprintf("./book_imgs/%s.jpg", isbn))
}

func ResolveBookImg(c echo.Context) error {
	req := new(resolveBookReq)
	if err := c.Bind(req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	sess, _ := session.Get("session", c)
	_, is_ok := sess.Values["UserID"].(string)
	if !is_ok {
		return c.NoContent(http.StatusForbidden)
	}
	book_data, err := resolveBook(req)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if book_data == nil {
		return c.NoContent(http.StatusAccepted)
	}
	isbn := primaryISBN(book_data.ID)
	if isbn == "" {
		return c.NoContent(http.StatusNotFound)
	}
	if book_data.ImageStatus == "Failed" {
		return c.NoContent(http.StatusNotFound)
	}
	if book_data.ImageStatus == "Waiting" {
		book_data.ImageStatus = "Processing"
		model.DB.Save(book_data)
		go func() {
			if !book.GetBookImg(isbn) {
				book_data.ImageStatus = "Failed"
				model.DB.Save(book_data)
			} else {
				book_data.ImageStatus = "Success"
				model.DB.Save(book_data)
			}
		}()
		return c.NoContent(http.StatusAccepted)
	} else if book_data.ImageStatus == "Processing" {
		return c.NoContent(http.StatusAccepted)
	} else if book_data.ImageStatus == "Success" {
		return c.File(fmt.Sprintf("./book_imgs/%s.jpg", isbn))
	} else {
		return c.NoContent(http.StatusNotFound)
	}
}
