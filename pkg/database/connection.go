package database

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

// init loads .env file if it exists
func init() {
	loadEnvFile()
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		// .env file is optional
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Only set if not already set
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

// GetDefaultConfig returns default database configuration from environment variables
func GetDefaultConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "hackernews"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
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

// Migrate runs database migrations
func Migrate() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	schema := `
-- Stories table
CREATE TABLE IF NOT EXISTS stories (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Story' NOT NULL,
    title TEXT NOT NULL,
    url TEXT,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
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
    story_id INTEGER NOT NULL,
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Comment' NOT NULL,
    text TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    parent_id INTEGER,
    reply_ids INTEGER[] DEFAULT '{}',
    FOREIGN KEY (story_id) REFERENCES stories(id) ON DELETE CASCADE
);

-- Polls table
CREATE TABLE IF NOT EXISTS polls (
    id INTEGER PRIMARY KEY,
    type VARCHAR(10) DEFAULT 'Poll' NOT NULL,
    title TEXT NOT NULL,
    score INTEGER DEFAULT 0 CHECK (score >= 0),
    author VARCHAR(255) NOT NULL,
    poll_options TEXT[] DEFAULT '{}',
    reply_ids INTEGER[] DEFAULT '{}',
    created_at BIGINT NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_stories_author ON stories(author);
CREATE INDEX IF NOT EXISTS idx_stories_created_at ON stories(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_stories_score ON stories(score DESC);

CREATE INDEX IF NOT EXISTS idx_asks_author ON asks(author);
CREATE INDEX IF NOT EXISTS idx_asks_created_at ON asks(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_asks_score ON asks(score DESC);

CREATE INDEX IF NOT EXISTS idx_jobs_author ON jobs(author);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_comments_story_id ON comments(story_id);
CREATE INDEX IF NOT EXISTS idx_comments_author ON comments(author);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);

CREATE INDEX IF NOT EXISTS idx_polls_author ON polls(author);
CREATE INDEX IF NOT EXISTS idx_polls_created_at ON polls(created_at DESC);`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Transaction starts a new database transaction
func Transaction() (*sql.Tx, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return db.Begin()
}
