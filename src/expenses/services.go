package expenses

import "github.com/jmoiron/sqlx"

type ExpensesService struct {
	DB *sqlx.DB
}

type ExpenseEntity struct {
	ID     int      `json:"id" db:"id"`
	Title  string   `json:"title" db:"title"`
	Amount int      `json:"amount" db:"amount"`
	Note   string   `json:"note" db:"note"`
	Tags   []string `json:"tags" db:"tags"`
}

func (es *ExpensesService) CreateExpense(e ExpenseEntity) (*ExpenseEntity, error) {
	res := new(ExpenseEntity)
	QUERY_STRING := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id,title,amount,note,tags"
	if err := es.DB.QueryRowx(QUERY_STRING, e.Title, e.Amount, e.Note, e.Tags).StructScan(&res); err != nil {
		return nil, err
	}

	return res, nil
}
