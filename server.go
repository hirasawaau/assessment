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
	"github.com/hirasawaau/assessment/src/health"
)

type Health struct {
	Status string `json:"status"`
}

func main() {
	app := fiber.New()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "2565"
	}

	app.Use(logger.New())

	app.Get("/health", health.GetHealthHandler)

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
