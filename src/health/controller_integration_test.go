//go:build integration

package health_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/health"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/stretchr/testify/assert"
)

const PORT = 3030

var HOST = fmt.Sprintf("localhost:%d", PORT)

func TestHealthControllerItTest(t *testing.T) {

	t.Run("GET /health", func(t *testing.T) {
		app := fiber.New()

		go utils.StartIntegrationApp(t, app)

		utils.WaitForConnection()

		req, err := http.NewRequest(fiber.MethodGet, utils.ConcatUrl("health"), nil)
		assert.NoError(t, err)
		resp, err := app.Test(req, 10000)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		var resp_body health.Health

		err = json.NewDecoder(resp.Body).Decode(&resp_body)

		fmt.Println(resp_body)

		assert.NoError(t, err)
		assert.Equal(t, "OK", string(resp_body.Status))

		err = app.Shutdown()
		assert.NoError(t, err)
	})
}
