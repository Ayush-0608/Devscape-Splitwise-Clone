package config

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() *sql.DB {
	connectionStr := "postgres://postgres:1234@localhost:5432/splitwise?sslmode=disable"

	db, err := sql.Open("pgx", connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the postgrSQL server ")

	return db
}
