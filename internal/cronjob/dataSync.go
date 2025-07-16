package cronjob

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"internship-project/internal/kafka"
	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
	"internship-project/internal/services"
	"internship-project/pkg/database"

	"github.com/go-co-op/gocron/v2"
)

type DataSyncService struct {
	scheduler         gocron.Scheduler
	apiClient         *services.HackerNewsApiClient
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
	apiClient *services.HackerNewsApiClient,
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
		apiClient:         apiClient,
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
		name      string
		interval  time.Duration
		task      func()
		immediate bool // Add this flag
	}{
		{
			name:     "sync-stories",
			interval: 50 * time.Minute,
			task:     d.syncStories,
		},
		{
			name:     "sync-asks",
			interval: 60 * time.Minute,
			task:     d.syncAsks,
		},
		{
			name:     "sync-jobs",
			interval: 60 * time.Minute,
			task:     d.syncJobs,
		},
		{
			name:     "sync-comments",
			interval: 60 * time.Minute,
			task:     d.syncComments,
		},
		{
			name:      "sync-updates",
			interval:  1 * time.Minute,
			task:      func() { d.syncUpdates() },
			immediate: true,
		},
	}

	for _, job := range jobs {
		// Run immediately
		if job.immediate {
			log.Printf("Running job %s immediately...", job.name)
			go job.task()
		}
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

func (d *DataSyncService) syncUpdates() {
	log.Println("Starting update sync...")

	ctx := context.Background()
	update, err := d.updateService.FetchUpdates(ctx)
	if err != nil {
		log.Printf("Error fetching updates: %v", err)
		return
	}

	if len(update.IDs) == 0 {
		log.Println("No items to sync in updates")
		return
	}

	log.Printf("Processing %d items from updates", len(update.IDs))

	// Initialize repositories
	storyRepo := postgres.NewStoryRepository()
	askRepo := postgres.NewAskRepository()
	commentRepo := postgres.NewCommentRepository()
	jobRepo := postgres.NewJobRepository()
	pollRepo := postgres.NewPollRepository()
	pollOptionRepo := postgres.NewPollOptionRepository()
	userRepo := postgres.NewUserRepository()

	// Collections for batch operations
	var stories []models.Story
	var asks []models.Ask
	var comments []models.Comment
	var jobs []models.Job
	var polls []models.Poll
	var pollOptions []models.PollOption
	var users []models.User

	// Collection of IDs for batch kafka producer
	var storiesIDs []int
	var asksIDs []int
	var commentsIDs []int
	var jobsIDs []int
	var pollsIDs []int
	var pollOptionsIDs []int

	// Process each item
	for _, itemID := range update.IDs {
		// Fetch raw item to determine type
		var rawItem map[string]interface{}
		err := d.apiClient.GetItem(ctx, itemID, &rawItem)
		if err != nil {
			log.Printf("Error fetching item %d: %v", itemID, err)
			continue
		}

		// Get item type
		itemType, ok := rawItem["type"].(string)
		if !ok {
			log.Printf("Item %d has no valid type", itemID)
			continue
		}

		log.Printf("Processing item %d of type: %s", itemID, itemType)

		// Route based on type
		switch itemType {
		case "story":
			var story models.Story
			err := d.apiClient.GetItem(ctx, itemID, &story)
			if err != nil {
				log.Printf("Error fetching story %d: %v", itemID, err)
				continue
			}
			if story.IsValid() {
				stories = append(stories, story)
				storiesIDs = append(storiesIDs, story.ID)
			}

		case "ask":
			var ask models.Ask
			err := d.apiClient.GetItem(ctx, itemID, &ask)
			if err != nil {
				log.Printf("Error fetching ask %d: %v", itemID, err)
				continue
			}
			if ask.IsValid() {
				asks = append(asks, ask)
				asksIDs = append(asksIDs, ask.ID)
			}

		case "comment":
			var comment models.Comment
			err := d.apiClient.GetItem(ctx, itemID, &comment)
			if err != nil {
				log.Printf("Error fetching comment %d: %v", itemID, err)
				continue
			}
			if comment.IsValid() {
				comments = append(comments, comment)
				commentsIDs = append(commentsIDs, comment.ID)

			}

		case "job":
			var job models.Job
			err := d.apiClient.GetItem(ctx, itemID, &job)
			if err != nil {
				log.Printf("Error fetching job %d: %v", itemID, err)
				continue
			}
			if job.IsValid() {
				jobs = append(jobs, job)
				jobsIDs = append(jobsIDs, job.ID)
			}

		case "poll":
			var poll models.Poll
			err := d.apiClient.GetItem(ctx, itemID, &poll)
			if err != nil {
				log.Printf("Error fetching poll %d: %v", itemID, err)
				continue
			}
			if poll.IsValid() {
				polls = append(polls, poll)
				pollsIDs = append(pollsIDs, poll.ID)
			}

		case "pollopt":
			var pollOption models.PollOption
			err := d.apiClient.GetItem(ctx, itemID, &pollOption)
			if err != nil {
				log.Printf("Error fetching poll option %d: %v", itemID, err)
				continue
			}
			if pollOption.IsValid() {
				pollOptions = append(pollOptions, pollOption)
				pollOptionsIDs = append(pollOptionsIDs, pollOption.ID)
			}

		default:
			log.Printf("Unknown item type '%s' for item %d", itemType, itemID)
		}
	}

	// Fetch users for all items
	for _, userID := range update.Profiles {
		var user models.User
		err := d.apiClient.Get(ctx, fmt.Sprintf("/user/%s.json", userID), &user)
		if err != nil {
			log.Printf("Error fetching user %s: %v", userID, err)
			continue
		}

		log.Printf("Processing user %s", userID)

		if user.IsValid() {
			users = append(users, user)
		}
	}

	// Batch save to database
	if len(stories) > 0 {
		storyPtrs := make([]*models.Story, len(stories))
		for i := range stories {
			storyPtrs[i] = &stories[i]
		}
		err = storyRepo.CreateBatchWithExistingIDs(ctx, storyPtrs)
		if err != nil {
			log.Printf("Error saving stories: %v", err)
		} else {
			err := kafka.NewProducer("StoriesTopic", storiesIDs)
			if err != nil {
				log.Printf("Error sending stories to Kafka: %v", err)
			} else {
				log.Printf("Sent %d stories to Kafka", len(stories))
			}
		}
	}

	if len(asks) > 0 {
		askPtrs := make([]*models.Ask, len(asks))
		for i := range asks {
			askPtrs[i] = &asks[i]
		}
		err = askRepo.CreateBatchWithExistingIDs(ctx, askPtrs)
		if err != nil {
			log.Printf("Error saving asks: %v", err)
		} else {
			err := kafka.NewProducer("AsksTopic", asksIDs)
			if err != nil {
				log.Printf("Error sending asks to Kafka: %v", err)
			} else {
				log.Printf("Sent %d asks to Kafka", len(asks))
			}
		}
	}

	if len(comments) > 0 {
		commentPtrs := make([]*models.Comment, len(comments))
		for i := range comments {
			commentPtrs[i] = &comments[i]
		}
		err = commentRepo.CreateBatchWithExistingIDs(ctx, commentPtrs)
		if err != nil {
			log.Printf("Error saving comments: %v", err)
		} else {
			err := kafka.NewProducer("CommentsTopic", commentsIDs)
			if err != nil {
				log.Printf("Error sending comments to Kafka: %v", err)
			} else {
				log.Printf("Sent %d comments to Kafka", len(comments))
			}
		}
	}

	if len(jobs) > 0 {
		jobPtrs := make([]*models.Job, len(jobs))
		for i := range jobs {
			jobPtrs[i] = &jobs[i]
		}
		err = jobRepo.CreateBatchWithExistingIDs(ctx, jobPtrs)
		if err != nil {
			log.Printf("Error saving jobs: %v", err)
		} else {
			err := kafka.NewProducer("JobsTopic", jobsIDs)
			if err != nil {
				log.Printf("Error sending jobs to Kafka: %v", err)
			} else {
				log.Printf("Sent %d jobs to Kafka", len(jobs))
			}
		}
	}

	if len(polls) > 0 {
		pollPtrs := make([]*models.Poll, len(polls))
		for i := range polls {
			pollPtrs[i] = &polls[i]
		}
		err = pollRepo.CreateBatchWithExistingIDs(ctx, pollPtrs)
		if err != nil {
			log.Printf("Error saving polls: %v", err)
		} else {
			err := kafka.NewProducer("PollsTopic", pollsIDs)
			if err != nil {
				log.Printf("Error sending polls to Kafka: %v", err)
			} else {
				log.Printf("Sent %d polls to Kafka", len(polls))
			}
		}
	}

	if len(pollOptions) > 0 {
		pollOptionPtrs := make([]*models.PollOption, len(pollOptions))
		for i := range pollOptions {
			pollOptionPtrs[i] = &pollOptions[i]
		}
		err = pollOptionRepo.CreateBatchWithExistingIDs(ctx, pollOptionPtrs)
		if err != nil {
			log.Printf("Error saving poll options: %v", err)
		} else {
			err := kafka.NewProducer("PollOptionsTopic", pollOptionsIDs)
			if err != nil {
				log.Printf("Error sending poll options to Kafka: %v", err)
			} else {
				log.Printf("Sent %d poll options to Kafka", len(pollOptions))
			}
		}
	}

	if len(users) > 0 {
		userPtrs := make([]*models.User, len(users))
		for i := range users {
			userPtrs[i] = &users[i]
		}
		err = userRepo.CreateBatchWithExistingIDs(ctx, userPtrs)
		if err != nil {
			log.Printf("Error saving users: %v", err)
		}
	}

	log.Printf("Update sync completed - Stories: %d, Asks: %d, Comments: %d, Jobs: %d, Polls: %d, Poll Options: %d, Users: %d",
		len(stories), len(asks), len(comments), len(jobs), len(polls), len(pollOptions), len(users))
}

func (d *DataSyncService) syncItemsFromMaxTo(items int, minusMaxItem int) {
	ctx := context.Background()

	// Initialize repositories
	storyRepo := postgres.NewStoryRepository()
	askRepo := postgres.NewAskRepository()
	commentRepo := postgres.NewCommentRepository()
	jobRepo := postgres.NewJobRepository()
	pollRepo := postgres.NewPollRepository()
	pollOptionRepo := postgres.NewPollOptionRepository()

	// Collections for batch operations
	var stories []models.Story
	var asks []models.Ask
	var comments []models.Comment
	var jobs []models.Job
	var polls []models.Poll
	var pollOptions []models.PollOption

	log.Printf("Starting sync for %d items...", items)

	maxItem, err := d.apiClient.GetMaxItemID()
	if err != nil {
		log.Printf("Error fetching max item ID: %v", err)
		return
	}

	maxItem -= minusMaxItem
	log.Printf("Max item ID is %d, starting sync from %d to %d", maxItem+minusMaxItem, maxItem-items+1, maxItem)

	// Process in batches of 100
	batchSize := 100
	for batch := 0; batch < items; batch += batchSize {
		end := batch + batchSize
		if end > items {
			end = items
		}

		var wg sync.WaitGroup
		var mu sync.Mutex

		// Process batch concurrently
		for i := batch; i < end; i++ {
			wg.Add(1)
			go func(itemID int) {
				defer wg.Done()

				var rawItem map[string]interface{}
				err := d.apiClient.GetItem(ctx, itemID, &rawItem)
				if err != nil {
					return
				}

				itemType, ok := rawItem["type"].(string)
				if !ok {
					return
				}

				switch itemType {
				case "story":
					var story models.Story
					if err := d.apiClient.GetItem(ctx, itemID, &story); err == nil && story.IsValid() {
						mu.Lock()
						stories = append(stories, story)
						mu.Unlock()
					}
				case "ask":
					var ask models.Ask
					if err := d.apiClient.GetItem(ctx, itemID, &ask); err == nil && ask.IsValid() {
						mu.Lock()
						asks = append(asks, ask)
						mu.Unlock()
					}
				case "comment":
					var comment models.Comment
					if err := d.apiClient.GetItem(ctx, itemID, &comment); err == nil && comment.IsValid() {
						mu.Lock()
						comments = append(comments, comment)
						mu.Unlock()
					}
				case "job":
					var job models.Job
					if err := d.apiClient.GetItem(ctx, itemID, &job); err == nil && job.IsValid() {
						mu.Lock()
						jobs = append(jobs, job)
						mu.Unlock()
					}
				case "poll":
					var poll models.Poll
					if err := d.apiClient.GetItem(ctx, itemID, &poll); err == nil && poll.IsValid() {
						mu.Lock()
						polls = append(polls, poll)
						mu.Unlock()
					}
				case "pollopt":
					var pollOption models.PollOption
					if err := d.apiClient.GetItem(ctx, itemID, &pollOption); err == nil && pollOption.IsValid() {
						mu.Lock()
						pollOptions = append(pollOptions, pollOption)
						mu.Unlock()
					}
				}
			}(maxItem - i)
		}

		wg.Wait()
		log.Printf("Processed batch %d-%d", batch, end)
	}

	// Save to database
	if len(stories) > 0 {
		storyPtrs := make([]*models.Story, len(stories))
		for i := range stories {
			storyPtrs[i] = &stories[i]
		}
		err = storyRepo.CreateBatchWithExistingIDs(ctx, storyPtrs)
		if err != nil {
			log.Printf("Error saving stories: %v", err)
		}
	}

	if len(asks) > 0 {
		askPtrs := make([]*models.Ask, len(asks))
		for i := range asks {
			askPtrs[i] = &asks[i]
		}
		err = askRepo.CreateBatchWithExistingIDs(ctx, askPtrs)
		if err != nil {
			log.Printf("Error saving asks: %v", err)
		}
	}

	if len(comments) > 0 {
		commentPtrs := make([]*models.Comment, len(comments))
		for i := range comments {
			commentPtrs[i] = &comments[i]
		}
		err = commentRepo.CreateBatchWithExistingIDs(ctx, commentPtrs)
		if err != nil {
			log.Printf("Error saving comments: %v", err)
		}
	}

	if len(jobs) > 0 {
		jobPtrs := make([]*models.Job, len(jobs))
		for i := range jobs {
			jobPtrs[i] = &jobs[i]
		}
		err = jobRepo.CreateBatchWithExistingIDs(ctx, jobPtrs)
		if err != nil {
			log.Printf("Error saving jobs: %v", err)
		}
	}

	if len(polls) > 0 {
		pollPtrs := make([]*models.Poll, len(polls))
		for i := range polls {
			pollPtrs[i] = &polls[i]
		}
		err = pollRepo.CreateBatchWithExistingIDs(ctx, pollPtrs)
		if err != nil {
			log.Printf("Error saving polls: %v", err)
		}
	}

	if len(pollOptions) > 0 {
		pollOptionPtrs := make([]*models.PollOption, len(pollOptions))
		for i := range pollOptions {
			pollOptionPtrs[i] = &pollOptions[i]
		}
		err = pollOptionRepo.CreateBatchWithExistingIDs(ctx, pollOptionPtrs)
		if err != nil {
			log.Printf("Error saving poll options: %v", err)
		}
	}

	log.Printf("Sync completed - Stories: %d, Asks: %d, Comments: %d, Jobs: %d, Polls: %d, Poll Options: %d",
		len(stories), len(asks), len(comments), len(jobs), len(polls), len(pollOptions))
}
