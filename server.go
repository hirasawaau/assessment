package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/hirasawaau/assessment/src/db"
	"github.com/hirasawaau/assessment/src/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "2565"
	}

	database, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer database.Close()

	if err = db.InitDB(database); err != nil {
		e.Logger.Fatal(err)
	}
	e.Validator = &utils.Validator{Validator: validator.New()}

	InjectApp(e, database)

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", PORT)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
