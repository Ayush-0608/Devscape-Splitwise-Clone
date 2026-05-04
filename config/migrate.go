package config

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT
	);`)

	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS friends(
	);`)

	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}
