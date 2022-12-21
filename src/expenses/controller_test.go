// go build:unit

package expenses_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostExpenses(t *testing.T) {

	db, mock, err := sqlmock.New()
	sqlx_db := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.Validator = &utils.Validator{
		Validator: validator.New(),
	}

	controller := expenses.ExpensesController{
		Instance: e,
		Service: &expenses.ExpensesService{
			DB: sqlx_db,
		},
	}
	t.Run("Create Expenses", func(t *testing.T) {
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

			req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ctx := controller.Instance.NewContext(req, rec)

			if assert.NoError(t, controller.PostExpensesHandler(ctx)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
				resp := new(expenses.ExpenseEntity)
				QUERY_STR := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id,title,amount,note,tags"
				mock.ExpectQuery(QUERY_STR).WithArgs(dto.Title, dto.Amount, dto.Note, dto.Tags)

				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), resp))
				assert.Equal(t, dto.Title, resp.Title)
				assert.Equal(t, dto.Amount, resp.Amount)
				assert.Equal(t, dto.Note, resp.Note)
				assert.Equal(t, dto.Tags, resp.Tags)
			}
		})

		t.Run("Should return error when property is not completed", func(t *testing.T) {
			var buf bytes.Buffer
			dto := expenses.ExpensesDto{
				Title:  "Test",
				Amount: 1,
				Note:   "Test Expense",
				Tags: []string{
					"Hello",
				},
			}
			assert.NoError(t, json.NewEncoder(&buf).Encode(dto))

			req := httptest.NewRequest("POST", "/expenses", &buf)
			rec := httptest.NewRecorder()

			ctx := controller.Instance.NewContext(req, rec)

			assert.Error(t, controller.PostExpensesHandler(ctx))
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	})
}
