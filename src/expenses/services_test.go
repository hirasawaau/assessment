//go:build unit
// +build unit

package expenses_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "postgres")
	assert.NoError(t, err)

	service := &expenses.ExpensesService{
		DB: dbx,
	}

	t.Run("should create expense with correct arguments", func(t *testing.T) {
		INSERT_STR := "INSERT INTO expenses"
		payload := expenses.ExpenseEntity{
			Title:  "Test",
			Amount: 1,
			Note:   "Test Expense",
			Tags:   pq.StringArray{"Hello"},
		}

		mock.ExpectQuery(INSERT_STR).WithArgs(payload.Title, payload.Amount, payload.Note, payload.Tags).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		expectedRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, payload.Title, payload.Amount, payload.Note, payload.Tags)
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(expectedRow)

		record, err := service.CreateExpense(payload)
		fmt.Println(*record)
		assert.NoError(t, err)
	})
}

func TestGetExpenseById(t *testing.T) {

	t.Run("should get expense with correct arguments", func(t *testing.T) {
		id := int64(1)
		dbx, mock, err := utils.GetMockDB()
		defer dbx.Close()
		service := &expenses.ExpensesService{
			DB: dbx,
		}
		row := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(id, "title", 0, "note", pq.Array([]string{"tag1", "tag2"}))
		mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = (.+)").WithArgs(id).WillReturnRows(row)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		record, err := service.GetExpenseById(id)

		if assert.NoError(t, err) {
			assert.Equal(t, int64(id), record.ID)
		}
	})
}
