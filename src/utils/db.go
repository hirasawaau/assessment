package utils

import "github.com/jmoiron/sqlx"

func InitDB(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "expenses" (
		"id" SERIAL PRIMARY KEY,
		"title" VARCHAR(255) NOT NULL,
		"amount" INT NOT NULL,
		"note" TEXT NOT NULL,
		"tags" VARCHAR(255) [] NOT NULL
	);
	`)
	return err
}
