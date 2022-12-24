package health

import (
	"github.com/gofiber/fiber/v2"
)

type Health struct {
	Status string `json:"status"`
}

type HealthController struct {
	Instance *fiber.App
}

func (hc *HealthController) getHealthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Health{Status: "OK"})
}

func (hc *HealthController) Handle() {
	g := hc.Instance.Group("/health")
	g.Get("", hc.getHealthHandler)
}
