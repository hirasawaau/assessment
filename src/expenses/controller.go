package expenses

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/middleware"
)

type ExpensesController struct {
	Instance *fiber.App
	Service  IExpensesService
}

type ExpensesCreateDto struct {
	Title  string   `json:"title" validate:"required"`
	Amount int      `json:"amount" validate:"required"`
	Note   string   `json:"note" validate:"required"`
	Tags   []string `json:"tags" validate:"required"`
}

type ExpensesUpdateDto struct {
	Title  string   `json:"title"`
	Amount int64    `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func (ec *ExpensesController) PostExpensesHandler(c *fiber.Ctx) error {
	validate := validator.New()
	body := new(ExpensesCreateDto)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	if err := validate.Struct(*body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	res, err := ec.Service.CreateExpense(*body)

	if err != nil {
		fmt.Println("SQL ERROR", err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (ec *ExpensesController) GetExpensesHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	resp, err := ec.Service.GetExpenseById(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ec *ExpensesController) PutExpensesHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	validate := validator.New()
	body := new(ExpensesUpdateDto)
	if err := c.BodyParser(&body); err != nil {
		fmt.Println("THIS ERROR")
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	if err := validate.Struct(*body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	resp, err := ec.Service.UpdateExpenseById(id, *body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)

}

func (ec *ExpensesController) GetAllExpensesHandler(c *fiber.Ctx) error {
	resp, err := ec.Service.GetExpenses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (ec *ExpensesController) Handle() {
	g := ec.Instance.Group("/expenses")
	g.Use(middleware.AuthMiddleware)
	g.Post("", ec.PostExpensesHandler)
	g.Get(":id", ec.GetExpensesHandler)
	g.Put(":id", ec.PutExpensesHandler)
	g.Get("", ec.GetAllExpensesHandler)
}
