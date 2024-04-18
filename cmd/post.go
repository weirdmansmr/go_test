package cmd

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type SumValue struct {
	Sum int `json:"sum"`
}

func UpdateBankAccount(c echo.Context) error {
    db, err := ConnectDB()
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