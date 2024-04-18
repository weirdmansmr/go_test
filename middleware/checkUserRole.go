package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckUserRole() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            userRole := c.Request().Header.Get("User-Role")

            if userRole != "admin" && userRole != "client" {
                return echo.NewHTTPError(http.StatusForbidden, "access denied")
            }
            
            return next(c)
        }
    }
}
