//go:build unit
// +build unit

package expenses_test

import (
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
		payload := expenses.ExpensesCreateDto{
			Title:  "Test",
			Amount: 1,
			Note:   "Test Expense",
			Tags:   pq.StringArray{"Hello"},
		}

		expectedRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, payload.Title, payload.Amount, payload.Note, pq.Array(payload.Tags))
		mock.ExpectQuery(INSERT_STR).WithArgs(payload.Title, payload.Amount, payload.Note, pq.Array(payload.Tags)).WillReturnRows(expectedRow)

		_, err := service.CreateExpense(payload)
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
		expectedRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(id, "title", 0, "note", pq.Array([]string{"tag1", "tag2"}))
		mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = (.+)").WithArgs(id).WillReturnRows(expectedRow)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		record, err := service.GetExpenseById(id)

		if assert.NoError(t, err) {
			assert.Equal(t, int64(id), record.ID)
		}
	})
}

func TestUpdateExpenseById(t *testing.T) {
	t.Run("should update expense with correct arguments", func(t *testing.T) {
		id := int64(1)
		dbx, mock, err := utils.GetMockDB()
		defer dbx.Close()
		service := &expenses.ExpensesService{
			DB: dbx,
		}

		dto := expenses.ExpensesUpdateDto{
			Title:  "title",
			Amount: 0,
			Note:   "_note",
			Tags:   []string{"tag1", "tag2"},
		}
		expectedRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(id, dto.Title, dto.Amount, dto.Note, pq.Array([]string{"tag1", "tag2"}))
		mock.ExpectQuery(`UPDATE expenses SET title = COALESCE\((.+), title\), amount = COALESCE\((.+), amount\), note = COALESCE\((.+), note\), tags = COALESCE\((.+), tags\) WHERE id = (.+) RETURNING (.+)`).WithArgs(dto.Title, dto.Amount, dto.Note, pq.Array([]string{"tag1", "tag2"}), id).WillReturnRows(expectedRow)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		record, err := service.UpdateExpenseById(id, dto)

		if assert.NoError(t, err) {
			assert.Equal(t, int64(id), record.ID)
			assert.Equal(t, "title", record.Title)
			assert.Equal(t, int64(0), record.Amount)
			assert.Equal(t, "_note", record.Note)
			assert.Equal(t, pq.StringArray{
				"tag1",
				"tag2",
			}, record.Tags)
		}
	})
}
