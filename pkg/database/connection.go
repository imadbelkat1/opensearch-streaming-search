package database

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database using the provided connection string.
func Connect(connStr string) error {
	var err error

	// sql.Open doesn't actually connect, just prepares
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database connection established successfully")
	return nil
}

// GetDB returns the established database connection.
// It panics if the connection has not been established.
func GetDB() *sql.DB {
	if DB == nil {
		panic("Database connection is not established")
	}
	return DB
}
