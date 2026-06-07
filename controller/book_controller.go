package controller

import (
	"amaryllis-api/book"
	"amaryllis-api/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type resolveBookReq struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	ISBN  string `json:"isbn"`
}

type bookRes struct {
	MetaData    model.Book
	Identifiers []model.BookIdentifier
}

func ReadBook(c echo.Context) error {
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
	if err := model.DB.First(book_data, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, buildBookRes(*book_data))
}

func ResolveBook(c echo.Context) error {
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
	return c.JSON(http.StatusOK, buildBookRes(*book_data))
}

func resolveBook(req *resolveBookReq) (*model.Book, error) {
	identifier_type := req.Type
	identifier_value := req.Value
	if req.ISBN != "" {
		identifier_type = "isbn"
		identifier_value = req.ISBN
	}
	identifier_type = strings.ToLower(strings.TrimSpace(identifier_type))
	identifier_value = normalizeIdentifier(identifier_value)
	if identifier_type == "" {
		identifier_type = "isbn"
	}
	if identifier_type != "isbn" || identifier_value == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest)
	}

	book_identifier := new(model.BookIdentifier)
	if err := model.DB.Where("type = ? AND value = ?", identifier_type, identifier_value).First(book_identifier).Error; err == nil {
		book_data := new(model.Book)
		if err := model.DB.First(book_data, book_identifier.BookID).Error; err == nil {
			return book_data, nil
		}
	}

	book_data := book.GetMetaData(identifier_value)
	if err := model.DB.Create(&book_data).Error; err != nil {
		return nil, err
	}
	book_identifier = &model.BookIdentifier{
		BookID: book_data.ID,
		Type:   identifier_type,
		Value:  identifier_value,
	}
	if err := model.DB.Create(book_identifier).Error; err != nil {
		return nil, err
	}
	return &book_data, nil
}

func primaryISBN(book_id uint) string {
	book_identifier := new(model.BookIdentifier)
	if err := model.DB.Where("book_id = ? AND type = ?", book_id, "isbn").First(book_identifier).Error; err != nil {
		return ""
	}
	return book_identifier.Value
}

func buildBookRes(book_data model.Book) bookRes {
	var identifiers []model.BookIdentifier
	model.DB.Where("book_id = ?", book_data.ID).Find(&identifiers)
	return bookRes{
		MetaData:    book_data,
		Identifiers: identifiers,
	}
}

func normalizeIdentifier(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, " ", "")
	return value
}
