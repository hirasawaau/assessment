package health

import (
	"github.com/gofiber/fiber/v2"
)

type Health struct {
	Status string `json:"status"`
}

func GetHealthHandler(c *fiber.Ctx) error {
	h := Health{
		Status: "OK",
	}
	return c.Status(fiber.StatusOK).JSON(h)
}
