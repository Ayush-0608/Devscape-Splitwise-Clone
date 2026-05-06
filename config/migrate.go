package config

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
		fullname TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		phone TEXT,
		password TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS friends(
		id1 INTEGER NOT NULL,
		id2 INTEGER NOT NULL,
		ACCEPTED BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT fk_requester FOREIGN KEY(id1) REFERENCES users(id) ON DELETE CASCADE,
		CONSTRAINT fk_receiver FOREIGN KEY(id2) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS groups(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_by INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT fk_creator FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE SET NULL
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS group_members(
		group_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT fk_group FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
		CONSTRAINT fk_member FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expenses(
		id SERIAL PRIMARY KEY,
		group_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		total_amount NUMERIC(10, 2)	NOT NULL,
		paid_by	INTEGER	 NOT NULL,
		split_type VARCHAR(10),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
		CONSTRAINT fk_payer FOREIGN KEY(paid_by) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS splits(
		id SERIAL PRIMARY KEY,
		expense_id INTEGER NOT NULL,
		user_id	INTEGER NOT NULL,
		amount_owed	NUMERIC(10, 2) NOT NULL,
		staus BOOLEAN DEFAULT FALSE,

		CONSTRAINT fk_expense FOREIGN KEY(expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
		CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}
