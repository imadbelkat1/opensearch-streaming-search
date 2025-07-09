package services

import (
	"context"

	"internship-project/internal/models"
)

// ApiDataFetcher defines the base interface for all API fetchers
type ApiDataFetcher[T any] interface {
	FetchByID(ctx context.Context, id int) (*T, error)
	FetchMultiple(ctx context.Context, ids []int) ([]*T, error)
	FetchTopItems(ctx context.Context) ([]int, error)
}

// UserApiFetcher defines the interface for user API operations
type UserApiFetcher interface {
	FetchByID(ctx context.Context, id int) (*models.User, error)
	FetchMultiple(ctx context.Context, ids []int) ([]*models.User, error)
	FetchByUsername(ctx context.Context, username string) (*models.User, error)
}

// StoryApiFetcher defines the interface for story API operations
type StoryApiFetcher interface {
	ApiDataFetcher[models.Story]
	FetchTopStories(ctx context.Context) ([]int, error)
	FetchNewStories(ctx context.Context) ([]int, error)
	FetchBestStories(ctx context.Context) ([]int, error)
}

// CommentApiFetcher defines the interface for comment API operations
type CommentApiFetcher interface {
	ApiDataFetcher[models.Comment]
}

// AskApiFetcher defines the interface for ask story API operations
type AskApiFetcher interface {
	ApiDataFetcher[models.Ask]
	FetchAskStories(ctx context.Context) ([]int, error)
}

// JobApiFetcher defines the interface for job API operations
type JobApiFetcher interface {
	ApiDataFetcher[models.Job]
	FetchJobStories(ctx context.Context) ([]int, error)
}

// PollApiFetcher defines the interface for poll API operations
type PollApiFetcher interface {
	ApiDataFetcher[models.Poll]
}

// PollOptionApiFetcher defines the interface for poll option API operations
type PollOptionApiFetcher interface {
	ApiDataFetcher[models.PollOption]
}

type UpdateSApiFetcher interface {
	FetchUpdates(ctx context.Context) ([]*models.Update, error)
}

// Example usage:
/*
func main() {
	ctx := context.Background()
	factory := NewHackerNewsApiServiceFactory()

	// Create services
	storyService := factory.CreateStoryService()
	userService := factory.CreateUserService()

	// Fetch top stories
	topStoryIds, err := storyService.FetchTopStories(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch first 10 stories
	stories, err := storyService.FetchMultiple(ctx, topStoryIds[:10])
	if err != nil {
		log.Fatal(err)
	}

	// Fetch a specific user
	user, err := userService.FetchByUsername(ctx, "pg")
	if err != nil {
		log.Fatal(err)
	}

	// Use the data to save to your repos
	for _, story := range stories {
		// Transform and save to your story repo
		// storyRepo.Save(transformStoryData(story))
	}
}
*/
