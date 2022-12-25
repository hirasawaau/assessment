//go:build integration

package expenses_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPostExpenses(t *testing.T) {

	t.Run("Should return correct result", func(t *testing.T) {
		app := fiber.New()
		go utils.StartIntegrationApp(t, app)

		utils.WaitForConnection()
		payload := `{"title":"Test","amount":1,"note":"Test Expense","tags":["Hello"]}`
		expected := regexp.MustCompile(`{"id":(.+),"title":"Test","amount":1,"note":"Test Expense","tags":\["Hello"\]}`)
		req, err := http.NewRequest(fiber.MethodPost, utils.ConcatUrl("expenses"), strings.NewReader(payload))
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		client := http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		bodyRes, err := io.ReadAll(resp.Body)

		if assert.NoError(t, err) {
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			assert.Regexp(t, expected, strings.TrimSpace(string(bodyRes)))
		}

		assert.NoError(t, err)
		err = app.Shutdown()
		assert.NoError(t, err)
	})

}

func TestIntegrationGetExpense(t *testing.T) {

	t.Run("Should return correct result", func(t *testing.T) {
		app := fiber.New()
		go utils.StartIntegrationApp(t, app)

		utils.WaitForConnection()
		id := "1"
		req, err := http.NewRequest(fiber.MethodGet, utils.ConcatUrl("expenses", id), nil)
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		client := http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		bodyRes, err := io.ReadAll(resp.Body)

		if assert.NoError(t, err) {
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			fmt.Println(string(bodyRes))
			resp_body := new(expenses.ExpenseEntity)
			err = json.Unmarshal(bodyRes, resp_body)
			assert.NoError(t, err)
			idx, err := strconv.ParseInt(id, 10, 64)
			assert.NoError(t, err)
			assert.Equal(t, resp_body.ID, idx)
		}

		assert.NoError(t, err)
		err = app.Shutdown()
		assert.NoError(t, err)
	})

}

func TestIntegrationPutExpenses(t *testing.T) {
	t.Run("Should return correct result", func(t *testing.T) {
		id := int64(1)
		app := fiber.New()

		go utils.StartIntegrationApp(t, app)
		utils.WaitForConnection()

		payload := `{"title":"Test","amount":1,"note":"Test Expense","tags":["Hello"]}`
		expected := fmt.Sprintf(`{"id":%d,"title":"Test","amount":1,"note":"Test Expense","tags":["Hello"]}`, id)

		req, err := http.NewRequest(fiber.MethodPut, utils.ConcatUrl("expenses", fmt.Sprintf("%d", id)), strings.NewReader(payload))
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		client := http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		bodyRes, err := io.ReadAll(resp.Body)

		if assert.NoError(t, err) {
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			assert.Equal(t, expected, strings.TrimSpace(string(bodyRes)))
		}

		assert.NoError(t, err)
		err = app.Shutdown()
		assert.NoError(t, err)
	})

}

func TestIntegrationGetExpenses(t *testing.T) {
	t.Run("Should return array that len >= 1", func(t *testing.T) {
		app := fiber.New()
		go utils.StartIntegrationApp(t, app)

		utils.WaitForConnection()
		req, err := http.NewRequest(fiber.MethodGet, utils.ConcatUrl("expenses"), nil)
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		client := http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		bodyRes, err := io.ReadAll(resp.Body)

		if assert.NoError(t, err) {
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			resp_body := new([]expenses.ExpenseEntity)
			err = json.Unmarshal(bodyRes, resp_body)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(*resp_body), 1)
		}

		assert.NoError(t, err)
		err = app.Shutdown()
		assert.NoError(t, err)
	})
}
