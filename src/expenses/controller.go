package expenses

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type ExpensesController struct {
	Instance *fiber.App
	Service  *ExpensesService
}

type ExpensesDto struct {
	Title  string   `json:"title" validate:"required"`
	Amount int      `json:"amount" validate:"required"`
	Note   string   `json:"note" validate:"required"`
	Tags   []string `json:"tags" validate:"required"`
}

func (ec *ExpensesController) PostExpensesHandler(c *fiber.Ctx) error {
	validate := validator.New()
	body := new(ExpensesDto)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	if err := validate.Struct(*body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	res, err := ec.Service.CreateExpense(ExpenseEntity{
		Title:  body.Title,
		Amount: body.Amount,
		Note:   body.Note,
		Tags:   body.Tags,
	})

	if err != nil {
		fmt.Println("SQL ERROR", err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (ec *ExpensesController) Handle() {
	g := ec.Instance.Group("/expenses")
	g.Post("", ec.PostExpensesHandler)
}
