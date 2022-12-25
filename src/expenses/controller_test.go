//go:build unit
// +build unit

package expenses_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	UpdateCalled  int
	GetsCalled    int
}

func (m *MockService) CreateExpense(expense expenses.ExpensesCreateDto) (*expenses.ExpenseEntity, error) {
	m.CreatedCalled++
	return m.GetExpenseById(1)
}

func (m *MockService) UpdateExpenseById(id int64, expense expenses.ExpensesUpdateDto) (*expenses.ExpenseEntity, error) {
	m.UpdateCalled++
	return m.GetExpenseById(id)
}

func TestPostExpenses(t *testing.T) {
	mockService := &MockService{}
	controller := &expenses.ExpensesController{
		Service: mockService,
	}

	app := fiber.New()

	app.Post("/expenses", controller.PostExpensesHandler)

	t.Run("should create expense with correct arguments", func(t *testing.T) {
		dto := expenses.ExpensesCreateDto{
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
			assert.Equal(t, int64(dto.Amount), respEntity.Amount)
			assert.Equal(t, dto.Note, respEntity.Note)
			assert.ElementsMatch(t, dto.Tags, respEntity.Tags)
		}
	})

	t.Run("should return 400 when title is empty", func(t *testing.T) {
		dto := expenses.ExpensesCreateDto{
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

func (m *MockService) GetExpenseById(id int64) (*expenses.ExpenseEntity, error) {
	m.GetCalled++
	return &expenses.ExpenseEntity{
		ID:     id,
		Title:  "Test",
		Amount: 1,
		Note:   "Test Expense",
		Tags:   []string{"Hello"},
	}, nil
}

func TestGetExpensesById(t *testing.T) {
	mockService := &MockService{}
	controller := &expenses.ExpensesController{
		Service: mockService,
	}

	app := fiber.New()

	app.Get("/expenses/:id", controller.GetExpensesHandler)

	t.Run("should get expense with correct arguments and return correctly", func(t *testing.T) {
		id := int64(1)

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/expenses/%d", id), nil)
		resp, err := app.Test(req, 100)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, mockService.GetCalled)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			respEntity := new(expenses.ExpenseEntity)
			resp_bytes, err := io.ReadAll(resp.Body)
			fmt.Println(string(resp_bytes))
			assert.NoError(t, err)
			err = json.Unmarshal(resp_bytes, respEntity)
			assert.NoError(t, err)
			assert.Equal(t, id, respEntity.ID)
		}
	})

	t.Run("should return 400 when id is not int", func(t *testing.T) {
		id := "abc"

		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/expenses/%s", id), nil)
		resp, err := app.Test(req, 100)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, mockService.GetCalled)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		}
	})
}

func TestPutExpenseById(t *testing.T) {
	mockService := &MockService{}
	controller := &expenses.ExpensesController{
		Service: mockService,
	}

	app := fiber.New()
	app.Put("/expenses/:id", controller.PutExpensesHandler)
	t.Run("Should update expense with correct arguments", func(t *testing.T) {
		id := int64(1)
		dto := expenses.ExpensesCreateDto{
			Title:  "Test",
			Amount: 1,
			Note:   "Test Expense",
			Tags:   []string{"Hello"},
		}
		payload, err := json.Marshal(dto)
		assert.NoError(t, err)
		req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/expenses/%d", id), bytes.NewReader(payload))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req, 100)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, mockService.UpdateCalled)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			respEntity := new(expenses.ExpenseEntity)
			resp_bytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			err = json.Unmarshal(resp_bytes, respEntity)
			assert.NoError(t, err)
			assert.Equal(t, id, respEntity.ID)
		}
	})

}

func (m *MockService) GetExpenses() ([]*expenses.ExpenseEntity, error) {
	m.GetsCalled++
	return []*expenses.ExpenseEntity{
		{
			ID:     1,
			Title:  "Test",
			Amount: 1,
			Note:   "Test Expense",
			Tags:   []string{"Hello"},
		},
	}, nil
}

func TestGetsExpenses(t *testing.T) {
	mockService := &MockService{}
	controller := &expenses.ExpensesController{
		Service: mockService,
	}

	app := fiber.New()
	app.Get("/expenses", controller.GetExpensesHandler)
	t.Run("Should called get expenses service", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, "/expenses", nil)
		_, err := app.Test(req, 100)
		if assert.NoError(t, err) {
			assert.Equal(t, 1, mockService.GetsCalled)
		}
	})

}
