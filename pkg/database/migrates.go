package database

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Migrate runs the database migrations using the provided migration files.
func Migrate() error {

	// Read the SQL file
	content, err := ioutil.ReadFile("migrations/001_create_tables.sql")
	if err != nil {
		return fmt.Errorf("error reading migration file: %v", err)
	}

	// Execute the SQL
	_, err = DB.Exec(string(content))
	if err != nil {
		return fmt.Errorf("error executing migration: %v", err)
	}

	log.Println("Database migration completed successfully!")
	return nil
}
