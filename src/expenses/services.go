package expenses

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ExpensesService struct {
	DB *sqlx.DB
}

type IExpensesService interface {
	CreateExpense(e ExpensesCreateDto) (*ExpenseEntity, error)
	GetExpenseById(id int64) (*ExpenseEntity, error)
	UpdateExpenseById(id int64, e ExpensesUpdateDto) (*ExpenseEntity, error)
}

type ExpenseEntity struct {
	ID     int64          `json:"id" db:"id"`
	Title  string         `json:"title" db:"title"`
	Amount int64          `json:"amount" db:"amount"`
	Note   string         `json:"note" db:"note"`
	Tags   pq.StringArray `json:"tags" db:"tags"`
}

func (es *ExpensesService) CreateExpense(e ExpensesCreateDto) (*ExpenseEntity, error) {
	res := new(ExpenseEntity)
	INSERT_STRING := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING *"
	if err := es.DB.QueryRowx(INSERT_STRING, e.Title, e.Amount, e.Note, pq.Array(e.Tags)).StructScan(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (es *ExpensesService) GetExpenseById(id int64) (*ExpenseEntity, error) {
	res := new(ExpenseEntity)
	QUERY_STRING := "SELECT * FROM expenses WHERE id = $1"
	if err := es.DB.QueryRowx(QUERY_STRING, id).StructScan(res); err != nil {
		fmt.Println("THIS ERROR", err)
		return nil, err
	}

	return res, nil
}

func GetUpdateExpenseQuery(e ExpensesUpdateDto) string {
	UPDATE_QUERY := `UPDATE expenses SET `
	if e.Title != "" {
		UPDATE_QUERY += `title = $1, `
	}
	if e.Amount != 0 {
		UPDATE_QUERY += `amount = $2, `
	}
	if e.Note != "" {
		UPDATE_QUERY += `note = $3, `
	}
	if len(e.Tags) != 0 {
		UPDATE_QUERY += `tags = $4 `
	}
	UPDATE_QUERY += `WHERE id = $5 RETURNING *`
	return UPDATE_QUERY
}

func (es *ExpensesService) UpdateExpenseById(id int64, e ExpensesUpdateDto) (*ExpenseEntity, error) {
	UPDATE_QUERY := `UPDATE expenses SET title = COALESCE($1, title), amount = COALESCE($2, amount), note = COALESCE($3, note), tags = COALESCE($4, tags) WHERE id = $5 RETURNING *`
	res := new(ExpenseEntity)
	if err := es.DB.QueryRowx(UPDATE_QUERY,
		e.Title, e.Amount, e.Note, pq.Array(e.Tags), id).StructScan(res); err != nil {
		return nil, err
	}

	return res, nil
}
