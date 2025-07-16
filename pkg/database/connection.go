package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"internship-project/internal/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetDefaultConfig returns default database configuration from environment variables
func GetDefaultConfig() *Config {
	return &Config{
		Host:     config.GetEnv("DB_HOST", "localhost"),
		Port:     config.GetEnv("DB_PORT", "5432"),
		User:     config.GetEnv("DB_USER", "postgres"),
		Password: config.GetEnv("DB_PASSWORD", "password"),
		DBName:   config.GetEnv("DB_NAME", "hackernews"),
		SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
	}
}

// DropAndRecreateDatabase drops the existing database and creates a new one
func DropAndRecreateDatabase(config *Config) error {
	if config == nil {
		config = GetDefaultConfig()
	}

	log.Printf("Dropping and recreating database: %s", config.DBName)

	// Connect to postgres database to drop/create target database
	tempConfig := *config
	tempConfig.DBName = "postgres"

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		tempConfig.Host, tempConfig.Port, tempConfig.User, tempConfig.Password, tempConfig.DBName, tempConfig.SSLMode,
	)

	tempDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer tempDB.Close()

	// Test connection
	if err := tempDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping postgres database: %w", err)
	}

	// Drop database if exists
	_, err = tempDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.DBName))
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	// Create database
	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	log.Printf("Successfully recreated database: %s", config.DBName)
	return nil
}

// Connect establishes database connection with provided config
func Connect(config *Config) error {
	if config == nil {
		config = GetDefaultConfig()
	}

	// Log connection attempt (without password)
	log.Printf("Attempting to connect to database: host=%s port=%s user=%s dbname=%s",
		config.Host, config.Port, config.User, config.DBName)

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database: %s", config.DBName)
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	if db == nil {
		panic("database connection not initialized - call Connect() first")
	}
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// Health checks database connectivity
func Health() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

// CleanDatabase drops all tables and recreates them
func CleanDatabase() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Cleaning database - dropping all tables...")

	// Drop all tables in reverse order to handle dependencies
	dropTables := []string{
		"DROP TABLE IF EXISTS poll_options CASCADE",
		"DROP TABLE IF EXISTS polls CASCADE",
		"DROP TABLE IF EXISTS comments CASCADE",
		"DROP TABLE IF EXISTS jobs CASCADE",
		"DROP TABLE IF EXISTS asks CASCADE",
		"DROP TABLE IF EXISTS stories CASCADE",
		"DROP TABLE IF EXISTS users CASCADE",
	}

	for _, dropSQL := range dropTables {
		_, err := db.Exec(dropSQL)
		if err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	log.Println("All tables dropped successfully")
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Running database migrations...")

	schema := `
-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    karma INTEGER NOT NULL DEFAULT 0 CHECK (karma >= 0),
    about TEXT NOT NULL DEFAULT '',
    created_at BIGINT NOT NULL,
    submitted_ids INTEGER[] DEFAULT '{}'
);
-- Stories table
CREATE TABLE IF NOT EXISTS stories (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Story' NOT NULL,
    title TEXT NOT NULL,
    url TEXT,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    comments_ids INTEGER[] DEFAULT '{}',     -- IDs of comments associated with the story
    comments_count INTEGER DEFAULT 0 CHECK (comments_count >= 0)
);

-- Asks table
CREATE TABLE IF NOT EXISTS asks (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Ask' NOT NULL,
    title TEXT NOT NULL,
    text TEXT,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    reply_ids INTEGER[] DEFAULT '{}',
    replies_count INTEGER DEFAULT 0 CHECK (replies_count >= 0),
    created_at BIGINT NOT NULL
);

-- Jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Job' NOT NULL,
    title TEXT NOT NULL,
    text TEXT,
    url TEXT,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL
);

-- Comments table 
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Comment' NOT NULL,
    text TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    parent_id INTEGER,
    reply_ids INTEGER[] DEFAULT '{}'
);

-- Polls table
CREATE TABLE IF NOT EXISTS polls (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Poll' NOT NULL,
    title TEXT NOT NULL,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    poll_options INTEGER[] DEFAULT '{}',
    reply_ids INTEGER[] DEFAULT '{}',
    created_at BIGINT NOT NULL
);

-- Poll Options table
CREATE TABLE IF NOT EXISTS poll_options (
    id INTEGER PRIMARY KEY NOT NULL,
    type VARCHAR(10) DEFAULT 'PollOption' NOT NULL,
    poll_id INTEGER NOT NULL,
    author VARCHAR(255) NOT NULL,
    option_text TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    votes INTEGER DEFAULT 0 CHECK (votes >= 0)
);
`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// FreshInit completely reinitializes the database
func FreshInit(config *Config) error {
	if config == nil {
		config = GetDefaultConfig()
	}

	log.Println("Starting fresh database initialization...")

	// Step 1: Drop and recreate database
	if err := DropAndRecreateDatabase(config); err != nil {
		return fmt.Errorf("failed to recreate database: %w", err)
	}

	// Step 2: Connect to new database
	if err := Connect(config); err != nil {
		return fmt.Errorf("failed to connect to new database: %w", err)
	}

	// Step 3: Run migrations
	if err := Migrate(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Step 4: Health check
	if err := Health(); err != nil {
		return fmt.Errorf("failed health check: %w", err)
	}

	log.Println("Fresh database initialization completed successfully!")
	return nil
}

// Transaction starts a new database transaction
func Transaction() (*sql.Tx, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return db.Begin()
}
