package health

import (
	"github.com/labstack/echo/v4"
)

type Health struct {
	Status string `json:"status"`
}

type HealthController struct {
	Instance *echo.Echo
}

func (hc *HealthController) getHealthHandler(c echo.Context) error {
	return c.JSON(200, Health{Status: "OK"})
}

func (hc *HealthController) Handle() {
	g := hc.Instance.Group("/health")
	g.GET("", hc.getHealthHandler)
}
