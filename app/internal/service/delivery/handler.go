package delivery

import (
	"net/http"

	"github.com/labstack/echo"
)

func NewDelivery(e *echo.Group) {
	e.GET("/ping", ping)
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
