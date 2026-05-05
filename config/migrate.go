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
		phone TEXT NULLABLE,
		password TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS friends(
		id1 INTEGER FOREIGN KEY(users.id) NOT NULL,
		id2 INTEGER FOREIGN KEY(users.id) NOT NULL,
		ACCEPTED BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS groups(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_by INTEGER FOREIGN KEY(users.id) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS group_members(
		group_id INTEGER FOREIGN KEY(groups.id) NOT NULL,
		user_id INTEGER FOREIGN KEY(users.id) NOT NULL,
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expenses(
		id SERIAL PRIMARY KEY,
		group_id INTEGER FOREIGN KEY(groups.id) NULLABLE,
		name TEXT NOT NULL,
		description TEXT NULLABLE,
		total_amount NUMERIC(10, 2)	NOT NULL,
		paid_by	INTEGER	 NOT NULL,
		expense_date TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

		CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL,
		CONSTRAINT fk_paid FOREIGN KEY(paid_by) REFERENCES users(id) ON DELETE CASCADE,
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS splits(
		id SERIAL PRIMARY KEY,
		expense_id INTEGER FOREIGN KEY(expenses.id) NOT NULL,
		user_id	INTEGER	FOREIGN KEY(users.id) NOT NULL,
		amount_owed	NUMERIC(10, 2) NOT NULL

		CONSTRAINT fk_expense FOREIGN KEY(expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
		CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS payments (
    	id SERIAL PRIMARY KEY,
    	payer_id INTEGER NOT NULL,
    	payee_id INTEGER NOT NULL,
    	group_id INTEGER,
    	amount NUMERIC(10, 2) NOT NULL,
    	status VARCHAR(10) NOT NULL DEFAULT 'PENDING',
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	completed_at TIMESTAMP,
    
    	CONSTRAINT fk_payer FOREIGN KEY(payer_id) REFERENCES users(id) ON DELETE CASCADE,
    	CONSTRAINT fk_payee FOREIGN KEY(payee_id) REFERENCES users(id) ON DELETE CASCADE,
    	CONSTRAINT fk_group FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE SET NULL
	);`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}
