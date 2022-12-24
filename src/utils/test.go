package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const PORT = 3000

var BASE_URL = fmt.Sprintf("127.0.0.1:%d", PORT)

func ConcatUrl(path ...string) string {
	return fmt.Sprintf("http://%s/", BASE_URL) + strings.Join(path, "/")
}

func StartIntegrationApp(t *testing.T, app *fiber.App) error {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	assert.NoError(t, err)

	err = InitDB(db)

	assert.NoError(t, err)
	InjectApp(app, db)

	err = app.Listen(fmt.Sprintf(":%d", PORT))

	assert.NoError(t, err)

	return nil
}

func WaitForConnection() {

	for {
		conn, err := net.Dial("tcp", BASE_URL)
		if err != nil {
			log.Println("Retrying to Connect", err)
		}
		if conn != nil {
			conn.Close()
			return
		}
	}
}
