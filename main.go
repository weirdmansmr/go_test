package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/weirdmansmr/go_test/middleware"
)

type BankAccount struct {
    ID       int    `json:"id"`
    Money    int    `json:"money"`
}

type SumValue struct {
	Sum int `json:"sum"`
}

func connectDB() (*sql.DB, error) {
    connStr := "postgres://postgres:Saykhanov01@localhost/bank_test?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %v", err)
    }
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error pinging the database: %v", err)
    }
    return db, nil
}

func getBankAccount(c echo.Context) error {
    db, err := connectDB()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, money FROM public.bank_account")
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error querying the database: %v", err))
    }
    defer rows.Close()

    var bank []BankAccount

    for rows.Next() {
        var b BankAccount
        if err := rows.Scan(&b.ID, &b.Money); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error scanning rows: %v", err))
        }
        bank = append(bank, b)
    }

    if err := rows.Err(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error iterating over rows: %v", err))
    }

    return c.JSON(http.StatusOK, bank)
}

func updateBankAccount(c echo.Context) error {
    db, err := connectDB()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    defer db.Close()

    var currentMoney int
    err = db.QueryRow("SELECT money FROM public.bank_account WHERE id = 1").Scan(&currentMoney)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error getting current money value: %v", err))
    }

    var sumValue SumValue
	if err := c.Bind(&sumValue); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

    if userRole := c.Request().Header.Get("User-Role"); userRole == "admin" {
        _, err = db.Exec("UPDATE public.bank_account SET money = $1 WHERE id = 1", currentMoney + sumValue.Sum)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error updating employee money: %v", err))
	}
    } else if userRole == "client" {
        if sumValue.Sum > currentMoney {
            return echo.NewHTTPError(http.StatusForbidden, "Insufficient funds")
        } else {
            _, err = db.Exec("UPDATE public.bank_account SET money = $1 WHERE id = 1", currentMoney - sumValue.Sum)
        }
    }

    return c.JSON(http.StatusOK, "s")
}

func main() {
    e := echo.New()

	e.Use(middleware.CheckUserRole())
	
    e.GET("/bank", getBankAccount)
    e.POST("/bankPost", updateBankAccount)

    fmt.Println("Server is running on port 8080")
    e.Logger.Fatal(e.Start(":8080"))
}