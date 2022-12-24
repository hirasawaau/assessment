//go:build integration

package expenses_test

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const SERVER_PORT = 3000

func TestIntegrationPostExpenses(t *testing.T) {
	app := fiber.New()

	HOST := fmt.Sprintf("localhost:%d", SERVER_PORT)

	go func(app *fiber.App) {
		db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatal(err)
		}

		utils.InjectApp(app, db)
		app.Listen(fmt.Sprintf(":%d", SERVER_PORT))
	}(app)
	start := time.Now()
	for {
		conn, err := net.Dial("tcp", HOST)
		if err != nil {
			continue
		}
		if conn != nil || time.Since(start) > 30*time.Second {
			conn.Close()
			break
		}
	}

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
		req, err := http.NewRequest(fiber.MethodPost, fmt.Sprintf("http://%s", utils.ConcatUrl(HOST, "expenses")), strings.NewReader(payload))
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		client := http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		bodyRes, err := io.ReadAll(resp.Body)

		if assert.NoError(t, err) {
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			assert.Equal(t, payload, strings.TrimSpace(string(bodyRes)))
		}

		err = app.Shutdown()
		assert.NoError(t, err)
	})
}
