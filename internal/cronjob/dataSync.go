package cronjob

import (
	"context"
	"fmt"
	"log"
	"time"

	"internship-project/internal/repository/postgres"
	"internship-project/internal/services"
	"internship-project/pkg/database"

	"github.com/go-co-op/gocron/v2"
)

type DataSyncService struct {
	scheduler         gocron.Scheduler
	userService       *services.UserApiService
	storyService      *services.StoryApiService
	commentService    *services.CommentApiService
	jobService        *services.JobApiService
	askService        *services.AskApiService
	pollService       *services.PollApiService
	pollOptionService *services.PollOptionApiService
	updateService     *services.UpdateApiService
}

// NewDataSyncService creates a new data sync service
func NewDataSyncService(
	userService *services.UserApiService,
	storyService *services.StoryApiService,
	commentService *services.CommentApiService,
	jobService *services.JobApiService,
	askService *services.AskApiService,
	pollService *services.PollApiService,
	pollOptionService *services.PollOptionApiService,
	updateService *services.UpdateApiService,
) (*DataSyncService, error) {
	// Create a single scheduler for all jobs
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &DataSyncService{
		scheduler:         scheduler,
		userService:       userService,
		storyService:      storyService,
		askService:        askService,
		jobService:        jobService,
		commentService:    commentService,
		pollService:       pollService,
		pollOptionService: pollOptionService,
		updateService:     updateService,
	}, nil
}

// Start begins all scheduled jobs
func (d *DataSyncService) Start() error {
	// Connect to the database
	log.Println("Connecting to the database...")
	config := database.GetDefaultConfig()
	if err := database.Connect(config); err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	// Register all jobs
	if err := d.registerJobs(); err != nil {
		return fmt.Errorf("failed to register jobs: %w", err)
	}

	// Start the scheduler
	d.scheduler.Start()
	log.Println("DataSyncService started with all cron jobs and database connection established!")
	return nil
}

// Stop gracefully stops all jobs
func (d *DataSyncService) Stop() error {
	if err := d.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown scheduler: %w", err)
	}
	log.Println("DataSyncService stopped")

	// shutdown the database connection
	if err := database.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

// registerJobs sets up all the cron jobs
func (d *DataSyncService) registerJobs() error {
	jobs := []struct {
		name     string
		interval time.Duration
		task     func()
	}{
		{
			name:     "sync-stories",
			interval: 1 * time.Minute,
			task:     d.syncStories,
		},
		{
			name:     "sync-asks",
			interval: 60 * time.Second,
			task:     d.syncAsks,
		},
		{
			name:     "sync-jobs",
			interval: 60 * time.Second,
			task:     d.syncJobs,
		},
		{
			name:     "sync-comments",
			interval: 1 * time.Minute,
			task:     d.syncComments,
		},
	}

	for _, job := range jobs {
		_, err := d.scheduler.NewJob(
			gocron.DurationJob(job.interval),
			gocron.NewTask(job.task),
			gocron.WithName(job.name),
		)
		if err != nil {
			return fmt.Errorf("failed to create job %s: %w", job.name, err)
		}
		log.Printf("Registered job: %s (every %v)", job.name, job.interval)
	}

	return nil
}

// Job implementations
func (d *DataSyncService) syncStories() {
	log.Println("Starting story sync...")

	// Fetch top stories
	ctx := context.Background()
	ids, err := d.storyService.FetchTopStories(ctx)
	if err != nil {
		log.Printf("Error fetching top stories: %v", err)
		return
	}

	stories, err := d.storyService.FetchMultiple(ctx, ids)
	if err != nil {
		log.Printf("Error fetching story details: %v", err)
		return
	}

	log.Printf("Successfully synced %d stories", len(stories))

	log.Println("Saving stories to the database...")

	r := postgres.NewStoryRepository()
	r.CreateBatchWithExistingIDs(ctx, stories)

	log.Println("Story sync completed")
	log.Printf("Total stories synced: %d", len(stories))
}

func (d *DataSyncService) syncAsks() {
	log.Println("Starting ask sync...")

	ctx := context.Background()
	ids, err := d.askService.FetchAskStories(ctx)
	if err != nil {
		log.Printf("Error fetching ask stories: %v", err)
		return
	}

	if len(ids) > 10 {
		ids = ids[:10]
	}

	asks, err := d.askService.FetchMultiple(ctx, ids)
	if err != nil {
		log.Printf("Error fetching ask details: %v", err)
		return
	}

	log.Printf("Successfully synced %d asks", len(asks))

	log.Println("Saving asks to the database...")

	r := postgres.NewAskRepository()
	err = r.CreateBatchWithExistingIDs(ctx, asks)
	if err != nil {
		log.Printf("Error saving asks to the database: %v", err)
		return
	}

	log.Println("Ask sync completed")
	log.Printf("Total asks synced: %d", len(asks))
}

func (d *DataSyncService) syncJobs() {
	log.Println("Starting job sync...")

	ctx := context.Background()
	ids, err := d.jobService.FetchJobStories(ctx)
	if err != nil {
		log.Printf("Error fetching job stories: %v", err)
		return
	}

	jobs, err := d.jobService.FetchMultiple(ctx, ids)
	if err != nil {
		log.Printf("Error fetching job details: %v", err)
		return
	}

	log.Printf("Successfully synced %d jobs", len(jobs))

	log.Println("Saving jobs to the database...")

	r := postgres.NewJobRepository()
	err = r.CreateBatchWithExistingIDs(ctx, jobs)
	if err != nil {
		log.Printf("Error saving jobs to the database: %v", err)
		return
	}

	log.Println("Job sync completed")
	log.Printf("Total jobs synced: %d", len(jobs))
}

func (d *DataSyncService) syncComments() {
	log.Println("Starting comment sync...")

	// Get some story IDs first
	ctx := context.Background()
	storyIDs, err := d.storyService.FetchTopStories(ctx)
	if err != nil {
		log.Printf("Error fetching stories for comments: %v", err)
		return
	}

	// Fetch stories to get comment IDs
	stories, err := d.storyService.FetchMultiple(ctx, storyIDs)
	if err != nil {
		log.Printf("Error fetching story details: %v", err)
		return
	}

	// Collect comment IDs
	var commentIDs []int
	for _, story := range stories {
		if len(story.Comments_ids) > 0 {
			limit := 300
			if len(story.Comments_ids) < limit {
				limit = len(story.Comments_ids)
			}
			commentIDs = append(commentIDs, story.Comments_ids[:limit]...)
		}
	}

	if len(commentIDs) == 0 {
		log.Println("No comments to sync")
		return
	}

	comments, err := d.commentService.FetchMultiple(ctx, commentIDs)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		return
	}

	// Save comments to the database
	r := postgres.NewCommentRepository()
	err = r.CreateBatchWithExistingIDs(ctx, comments)
	if err != nil {
		log.Printf("Error saving comments to the database: %v", err)
		return
	}

	log.Printf("Successfully synced %d comments", len(comments))
}
