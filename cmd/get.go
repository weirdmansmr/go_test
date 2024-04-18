package cmd

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type BankAccount struct {
    ID       int    `json:"id"`
    Money    int    `json:"money"`
}

func GetBankAccount(c echo.Context) error {
	db, err := ConnectDB()
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