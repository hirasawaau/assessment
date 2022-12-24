package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hirasawaau/assessment/src/expenses"
	"github.com/hirasawaau/assessment/src/health"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InjectApp(app *fiber.App, db *sqlx.DB) {
	healthController := health.HealthController{Instance: app}
	healthController.Handle()

	expensesService := expenses.ExpensesService{DB: db}
	expensesController := expenses.ExpensesController{Instance: app, Service: &expensesService}
	expensesController.Handle()
}

func main() {
	app := fiber.New()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "2565"
	}

	app.Use(logger.New())

	db := sqlx.MustOpen("postgres", os.Getenv("DATABASE_URL"))

	defer db.Close()

	if err := utils.InitDB(db); err != nil {
		log.Fatal(err)
	}

	utils.InjectApp(app, db)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", PORT)); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
