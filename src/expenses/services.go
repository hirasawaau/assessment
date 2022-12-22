package expenses

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ExpensesService struct {
	DB *sqlx.DB
}

type ExpenseEntity struct {
	ID     int            `json:"id" db:"id"`
	Title  string         `json:"title" db:"title"`
	Amount int            `json:"amount" db:"amount"`
	Note   string         `json:"note" db:"note"`
	Tags   pq.StringArray `json:"tags" db:"tags"`
}

func (es *ExpensesService) CreateExpense(e ExpenseEntity) (*ExpenseEntity, error) {
	var insertedId int64
	INSERT_STRING := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id"
	if err := es.DB.QueryRowx(INSERT_STRING, e.Title, e.Amount, e.Note, pq.Array(e.Tags)).Scan(&insertedId); err != nil {
		return nil, err
	}

	res, err := es.GetExpenseById(insertedId)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (es *ExpensesService) GetExpenseById(id int64) (*ExpenseEntity, error) {
	res := new(ExpenseEntity)
	QUERY_STRING := "SELECT * FROM expenses WHERE id = $1"
	if err := es.DB.QueryRowx(QUERY_STRING, id).StructScan(res); err != nil {
		fmt.Println("THIS ERROR")
		return nil, err
	}

	return res, nil
}
