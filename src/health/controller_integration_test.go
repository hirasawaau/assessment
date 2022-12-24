//go:build integration

package health_test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthItTest(t *testing.T) {

	t.Run("Should return correct result", func(t *testing.T) {
		app := fiber.New()
		go utils.StartIntegrationApp(t, app)
		utils.WaitForConnection()

		req, err := http.NewRequest(fiber.MethodGet, utils.ConcatUrl("health"), nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)

		if assert.NoError(t, err) {
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		err = app.Shutdown()
		assert.NoError(t, err)
	})
}
