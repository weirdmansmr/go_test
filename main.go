package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/weirdmansmr/go_test/middleware"
)

type Employee struct {
    ID       int    `json:"id"`
    Money    int    `json:"money"`
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

func getEmployees(c echo.Context) error {
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

    var employees []Employee

    for rows.Next() {
        var e Employee
        if err := rows.Scan(&e.ID, &e.Money); err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error scanning rows: %v", err))
        }
        employees = append(employees, e)
    }

    if err := rows.Err(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error iterating over rows: %v", err))
    }

    return c.JSON(http.StatusOK, employees)
}

func main() {
    e := echo.New()

		e.Use(middleware.CheckUserRole("admin"))
		
    e.GET("/bank", getEmployees)

    fmt.Println("Server is running on port 8080")
    e.Logger.Fatal(e.Start(":8080"))
}