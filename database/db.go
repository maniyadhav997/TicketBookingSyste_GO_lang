package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB is a global variable that will hold our database connection.
// Using a global like this is simple and common for small beginner projects.
var DB *sql.DB

// InitDB opens the SQLite database file and creates the tables if they do not exist.
func InitDB() error {
	// Open will create the file if it does not exist.
	db, err := sql.Open("sqlite", "./database/app.db")
	if err != nil {
		return err
	}

	// We assign the opened db to our global variable so other packages can use it.
	DB = db

	// We create the users table first.
	err = createUsersTable()
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	// Then we create the tickets table.
	err = createTicketsTable()
	if err != nil {
		return fmt.Errorf("error creating tickets table: %w", err)
	}

	return nil
}

// createUsersTable runs a simple SQL statement to create the users table.
func createUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(query)
	return err
}

// createTicketsTable runs a simple SQL statement to create the tickets table.
func createTicketsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(query)
	return err
}

