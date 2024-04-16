package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CheckUserRole(allowedRole string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            userRole := c.Request().Header.Get("User-Role")

            if userRole != allowedRole {
                return echo.NewHTTPError(http.StatusUnauthorized, "Only admins are allowed")
            }
            
            return next(c)
        }
    }
}
