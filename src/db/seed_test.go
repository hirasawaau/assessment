//go:build seed

package db_test

import (
	"os"
	"testing"

	"github.com/hirasawaau/assessment/src/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSeedDB(t *testing.T) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	utils.InitDB(db)
	if assert.NoError(t, err) {
		QUERY := `INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id`
		_, err := db.Exec(QUERY, "Test", 100, "Test", pq.Array([]string{"test"}))
		assert.NoError(t, err)
	}
}
