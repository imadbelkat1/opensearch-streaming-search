package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"internship-project/internal/cronjob"
	"internship-project/internal/services"
)

func main() {
	log.Println("Starting HackerNews Data Sync...")

	// Create HTTP client
	client := services.NewHackerNewsApiClient()

	// Create all services
	userService := services.NewUserApiService(client)
	storyService := services.NewStoryApiService(client)
	commentService := services.NewCommentApiService(client)
	jobService := services.NewJobApiService(client)
	askService := services.NewAskApiService(client)
	pollService := services.NewPollApiService(client)
	pollOptionService := services.NewPollOptionApiService(client)
	updateService := services.NewUpdateApiService(client)

	// Create and start data sync service
	dataSyncService, err := cronjob.NewDataSyncService(
		client,
		userService,
		storyService,
		commentService,
		jobService,
		askService,
		pollService,
		pollOptionService,
		updateService,
	)
	if err != nil {
		log.Fatal("Failed to create data sync service:", err)
	}

	//  Start all cron jobs
	if err := dataSyncService.Start(); err != nil {
		log.Fatal("Failed to start cron jobs:", err)
	}

	log.Println("All cron jobs started successfully!")
	log.Println("Data sync is now running automatically...")
	log.Println("Press Ctrl+C to stop")

	//  Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Stopping application...")
	if err := dataSyncService.Stop(); err != nil {
		log.Printf("Error stopping service: %v", err)
	} else {
		log.Println("Application stopped successfully")
	}
}
