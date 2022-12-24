//go:build unit
// +build unit

package expenses_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type MockService struct {
	CreatedCalled int
	GetCalled     int
}

func (m *MockService) CreateExpense(expense expenses.ExpenseEntity) (*expenses.ExpenseEntity, error) {
	m.CreatedCalled++
	return m.GetExpenseById(1)
}

func (m *MockService) GetExpenseById(id int64) (*expenses.ExpenseEntity, error) {
	m.GetCalled++
	return &expenses.ExpenseEntity{
		ID:     1,
		Title:  "Test",
		Amount: 1,
		Note:   "Test Expense",
		Tags:   []string{"Hello"},
	}, nil
}

func TestPostExpenses(t *testing.T) {
	mockService := &MockService{}
	controller := &expenses.ExpensesController{
		Service: mockService,
	}

	app := fiber.New()

	app.Post("/expenses", controller.PostExpensesHandler)

	t.Run("should create expense with correct arguments", func(t *testing.T) {
		dto := expenses.ExpensesDto{
			Title:  "Test",
			Amount: 1,
			Note:   "Test Expense",
			Tags:   pq.StringArray{"Hello"},
		}
		payload, err := json.Marshal(dto)
		assert.NoError(t, err)
		req := httptest.NewRequest(fiber.MethodPost, "/expenses", bytes.NewReader(payload))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req, 100)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, mockService.CreatedCalled)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
			respEntity := new(expenses.ExpenseEntity)
			resp_bytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			err = json.Unmarshal(resp_bytes, respEntity)
			assert.NoError(t, err)
			assert.Equal(t, dto.Title, respEntity.Title)
			assert.Equal(t, dto.Amount, respEntity.Amount)
			assert.Equal(t, dto.Note, respEntity.Note)
			assert.ElementsMatch(t, dto.Tags, respEntity.Tags)
		}
	})

	t.Run("should return 400 when title is empty", func(t *testing.T) {
		dto := expenses.ExpensesDto{
			Amount: 1,
			Note:   "Test Expense",
			Tags:   pq.StringArray{"Hello"},
		}
		payload, err := json.Marshal(dto)
		assert.NoError(t, err)
		req := httptest.NewRequest(fiber.MethodPost, "/expenses", bytes.NewReader(payload))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req, 100)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, mockService.CreatedCalled)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		}
	})
}
