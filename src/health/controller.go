package health

import (
	"github.com/labstack/echo/v4"
)

type Health struct {
	Status string `json:"status"`
}

func GetHealthHandler(c echo.Context) error {
	h := Health{
		Status: "OK",
	}
	return c.JSON(200, h)
}
