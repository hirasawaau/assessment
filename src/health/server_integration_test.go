// go:build integration
package health_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hirasawaau/assessment/src/health"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	ADDR := fmt.Sprintf("http://localhost:%s/health", os.Getenv("PORT"))
	t.Run("GET /health", func(t *testing.T) {
		resp, err := http.Get(ADDR)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		var resp_body health.Health

		err = json.NewDecoder(resp.Body).Decode(&resp_body)

		fmt.Println(resp_body)

		assert.NoError(t, err)
		assert.Equal(t, "OK", string(resp_body.Status))
	})
}
