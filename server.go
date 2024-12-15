package main

import (
	"amaryllis-api/controller"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/users", controller.CreateUser)
	e.Logger.Fatal(e.Start(":1323"))
}
