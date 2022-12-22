//go:build unit
// +build unit

package expenses_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hirasawaau/assessment/src/expenses"
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
