package expenses

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ExpensesController struct {
	Instance *echo.Echo
	Service  *ExpensesService
}

type ExpensesDto struct {
	Title  string   `json:"title" validate:"required"`
	Amount int      `json:"amount" validate:"required"`
	Note   string   `json:"note" validate:"required"`
	Tags   []string `json:"tags" validate:"required"`
}

func (ec *ExpensesController) PostExpensesHandler(c echo.Context) error {

	body := new(ExpensesDto)
	if err := c.Bind(body); err != nil {
		fmt.Println("THIS CASE")
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	res, err := ec.Service.CreateExpense(ExpenseEntity{
		Title:  body.Title,
		Amount: body.Amount,
		Note:   body.Note,
		Tags:   body.Tags,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, res)
}

func (ec *ExpensesController) Handle() {
	g := ec.Instance.Group("/expenses")
	g.POST("", ec.PostExpensesHandler)
}
