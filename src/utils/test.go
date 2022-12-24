package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func ConcatUrl(base string, path ...string) string {
	return fmt.Sprintf("http://%s/", base) + strings.Join(path, "/")
}

func IntegrationApp(t *testing.T, app *fiber.App, port int) error {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	assert.NoError(t, err)

	InitDB(db)

	InjectApp(app, db)

	app.Listen(fmt.Sprintf(":%d", port))

	return nil
}
