package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"os/signal"
	"syscall"

	"github.com/hirasawaau/assessment/src/health"
	"github.com/labstack/echo/v4"
)

type Health struct {
	Status string `json:"status"`
}

func main() {
	e := echo.New()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "2565"
	}

	e.GET("/health", health.GetHealthHandler)

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
