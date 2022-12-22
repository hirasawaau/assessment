//go:build integration

package expenses_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func ConcatUrl(url ...string) string {
	BASE_URL := os.Getenv("BASE_URL") + ":" + os.Getenv("PORT")
	return strings.Join(append([]string{BASE_URL}, url...), "/")
}

func TestIntegrationPostExpenses(t *testing.T) {
	t.Run("Should return correct result", func(t *testing.T) {
		payload := `
		{
			"title": "Test",
			"amount": 1,
			"note": "Test Expense",
			"tags": [
				"Hello"
			]
		}
		`
		uri := ConcatUrl("expenses")
		fmt.Println(uri)
		resp, err := http.Post(ConcatUrl("expenses"), fiber.MIMEApplicationJSON, strings.NewReader(payload))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})
}
