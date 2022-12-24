package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/hirasawaau/assessment/src/health"
	"github.com/jmoiron/sqlx"
)

func InjectApp(app *fiber.App, db *sqlx.DB) {
	healthController := health.HealthController{Instance: app}
	healthController.Handle()

	expensesService := expenses.ExpensesService{DB: db}
	expensesController := expenses.ExpensesController{Instance: app, Service: &expensesService}
	expensesController.Handle()
}
