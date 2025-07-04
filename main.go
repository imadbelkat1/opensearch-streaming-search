package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
	"internship-project/pkg/database"
)

func main() {
	// Initialize database connection
	if err := initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Example usage of repositories
	if err := exampleUsage(); err != nil {
		log.Printf("Example usage error: %v", err)
	}

	// Setup graceful shutdown
	gracefulShutdown()
}

func initDatabase() error {
	// Connect to database
	config := database.GetDefaultConfig()
	if err := database.Connect(config); err != nil {
		return err
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		return err
	}

	// Check health
	if err := database.Health(); err != nil {
		return err
	}

	return nil
}

func exampleUsage() error {
	ctx := context.Background()

	// Initialize repositories
	storyRepo := postgres.NewStoryRepository()
	commentRepo := postgres.NewCommentRepository()

	// Example: Create a story
	story := &models.Story{
		ID:             12348,
		Type:           "Story",
		Title:          "Example Story Title",
		URL:            "https://example.com/story",
		Score:          42,
		Author:         "testuser",
		Created_At:     time.Now().Unix(),
		Comments_count: 0,
	}

	if err := storyRepo.Create(ctx, story); err != nil {
		log.Printf("Failed to create story: %v", err)
	} else {
		log.Println("Story created successfully")
	}

	// Example: Get recent stories
	stories, err := storyRepo.GetRecent(ctx, 10)
	if err != nil {
		log.Printf("Failed to get recent stories: %v", err)
	} else {
		log.Printf("Found %d recent stories", len(stories))
	}

	// Example: Create a comment
	comment := &models.Comment{
		ID:         54319,
		StoryID:    12345,
		Type:       "Comment",
		Text:       "This is an example comment",
		Author:     "commenter",
		Parent:     0,
		Replies:    []int{},
		Created_At: time.Now().Unix(),
	}

	if err := commentRepo.Create(ctx, comment); err != nil {
		log.Printf("Failed to create comment: %v", err)
	} else {
		log.Println("Comment created successfully")
	}

	// Example: Check if story exists
	exists, err := storyRepo.Exists(ctx, 12345)
	if err != nil {
		log.Printf("Failed to check story existence: %v", err)
	} else {
		log.Printf("Story exists: %v", exists)
	}

	// Example: Get count of stories
	count, err := storyRepo.GetCount(ctx)
	if err != nil {
		log.Printf("Failed to get story count: %v", err)
	} else {
		log.Printf("Total stories: %d", count)
	}

	return nil
}

func gracefulShutdown() {
	// Setup signal catching
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until signal received
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	// Cleanup
	log.Println("Shutting down gracefully...")

	// Give ongoing operations time to complete
	time.Sleep(2 * time.Second)

	log.Println("Shutdown complete")
}
