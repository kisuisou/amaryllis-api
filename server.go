package main

import (
	"amaryllis-api/controller"
	"amaryllis-api/model"
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("loading .env failed")
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		e := echo.New()
		e.POST("/users", controller.CreateUser)
		e.Logger.Fatal(e.Start(":1323"))
	} else if flag.Arg(0) == "migrate" {
		model.Connect()
		model.Migrate()
	}

}
