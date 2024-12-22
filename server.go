package main

import (
	"amaryllis-api/book"
	"amaryllis-api/controller"
	"amaryllis-api/model"
	"flag"
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("loading .env failed")
	}
	flag.Parse()
	model.Connect()
	if len(flag.Args()) == 0 {
		e := echo.New()
		e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
		e.POST("/users", controller.CreateUser)
		e.GET("/signin", controller.ReadSession)
		e.POST("/signin", controller.CreateSession)
		e.DELETE("/signin", controller.DeleteSession)
		e.Logger.Fatal(e.Start(":1323"))
	} else if flag.Arg(0) == "migrate" {
		model.Migrate()
	} else if flag.Arg(0) == "add_user" {
		var user_id, password string
		user := new(model.User)
		fmt.Print("Enter UserID--->")
		fmt.Scan(&user_id)
		err := model.DB.Where("user_id = ?", user_id).First(user).Error
		if err == nil {
			fmt.Println("すでに使われているIDです")
		} else {
			fmt.Print("Enter Password--->")
			fmt.Scan(&password)
			hash, _ := argon2id.CreateHash(password, argon2id.DefaultParams)
			model.DB.Create(&model.User{UserID: user_id, PasswordHash: hash})
		}

	} else if flag.Arg(0) == "get_book_data" {
		fmt.Println(book.GetMetaData(flag.Arg(1)))
	}

}
