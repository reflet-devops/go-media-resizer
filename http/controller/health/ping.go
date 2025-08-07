package health

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetPing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
