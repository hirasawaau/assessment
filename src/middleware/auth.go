package middleware

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	auth_header := c.Get(fiber.HeaderAuthorization)
	isWrong, err := regexp.MatchString(`(.*)wrong_token$`, auth_header)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if isWrong {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}

	return c.Next()
}
