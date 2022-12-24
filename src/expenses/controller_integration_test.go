//go:build integration

package expenses_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPostExpenses(t *testing.T) {

	t.Run("Should return correct result", func(t *testing.T) {
		app := fiber.New()
		go utils.StartIntegrationApp(t, app)

		utils.WaitForConnection()
		payload := `{"id":1,"title":"Test","amount":1,"note":"Test Expense","tags":["Hello"]}`
		req, err := http.NewRequest(fiber.MethodPost, utils.ConcatUrl("expenses"), strings.NewReader(payload))
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		client := http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		bodyRes, err := io.ReadAll(resp.Body)

		if assert.NoError(t, err) {
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			assert.Equal(t, strings.TrimSpace(payload), strings.TrimSpace(string(bodyRes)))
		}

		assert.NoError(t, err)
		err = app.Shutdown()
		assert.NoError(t, err)
	})

}
