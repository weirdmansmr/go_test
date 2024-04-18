package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/weirdmansmr/go_test/cmd"
	"github.com/weirdmansmr/go_test/middleware"
)

func main() {
    e := echo.New()

	e.Use(middleware.CheckUserRole())
	
    e.GET("/bank", cmd.GetBankAccount)
    e.POST("/bankPost", cmd.UpdateBankAccount)

    fmt.Println("Server is running on port 8080")
    e.Logger.Fatal(e.Start(":8080"))
}