//go build:unit

package expenses_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPostExpenses(t *testing.T) {

	db, mock, err := sqlmock.New()
	sqlx_db := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatal(err)
	}
	app := fiber.New()

	controller := expenses.ExpensesController{
		Instance: app,
		Service: &expenses.ExpensesService{
			DB: sqlx_db,
		},
	}
	t.Run("Create Expenses", func(t *testing.T) {

		controller.Handle()
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
			dto := new(expenses.ExpensesDto)
			assert.NoError(t, json.Unmarshal([]byte(payload), dto))

			req := httptest.NewRequest(fiber.MethodPost, "/expenses", strings.NewReader(payload))
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			QUERY_STR := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id,title,amount,note,tags"
			mock.ExpectQuery(QUERY_STR)

			rec, err := app.Test(req, 100)

			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusCreated, rec.StatusCode)
				resp := new(expenses.ExpenseEntity)

				bodyResp, err := ioutil.ReadAll(rec.Body)
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(bodyResp, resp))

				assert.Equal(t, dto.Title, resp.Title)
				assert.Equal(t, dto.Amount, resp.Amount)
				assert.Equal(t, dto.Note, resp.Note)
				assert.Equal(t, dto.Tags, resp.Tags)
			}
		})

		t.Run("Should return error when property is not completed", func(t *testing.T) {
			payload := `
			{
				"title": "Test",
				"note": "Test Expense",
				"tags": [
					"Hello"
				]
			}
			`

			req := httptest.NewRequest(fiber.MethodPost, "/expenses", strings.NewReader(payload))

			rec, err := app.Test(req, 100)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Status)
		})
	})
}
