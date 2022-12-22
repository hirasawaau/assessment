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
	res := new(ExpenseEntity)
	INSERT_STRING := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id,title,amount,note,tags"
	result, err := es.DB.Exec(INSERT_STRING, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	err = es.DB.Get(res, "SELECT id,title,amount,note,tags FROM expenses WHERE id = $1", lastId)

	if err != nil {
		fmt.Println("This Error")
		return nil, err
	}

	return res, nil
}
