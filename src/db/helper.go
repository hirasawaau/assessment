package db

import "github.com/jmoiron/sqlx"

func InitDB(d *sqlx.DB) error {

	_, err := d.Exec(`
	CREATE TABLE IF NOT EXISTS "expenses" (
		"id" SERIAL PRIMARY KEY,
		"title" VARCHAR(255) NOT NULL,
		"amount" INT NOT NULL,
		"note" TEXT NOT NULL,
		"category" TEXT NOT NULL,
		"tags" VARCHAR(255) [] NOT NULL,
	);
	`)

	return err
}
