package main

import (
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/hirasawaau/assessment/src/health"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func InjectApp(e *echo.Echo, db *sqlx.DB) {
	healthController := health.HealthController{Instance: e}
	healthController.Handle()

	expensesService := expenses.ExpensesService{DB: db}
	expensesController := expenses.ExpensesController{Instance: e, Service: &expensesService}
	expensesController.Handle()
}
