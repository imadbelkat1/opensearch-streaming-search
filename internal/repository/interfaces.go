package repository

import (
	"context"
	"internship-project/internal/models"
)

type StoryRepository interface {
	// CREATE OPERATIONS
	CreateStory(ctx context.Context, story *models.Story) error

	// READ OPERATIONS
	GetStoryByID(ctx context.Context, id int) (*models.Story, error)
	GetAllStories(ctx context.Context) ([]*models.Story, error)
	GetRecentStories(ctx context.Context, limit int) ([]*models.Story, error)
	GetStoriesByMinScore(ctx context.Context, minScore int) ([]*models.Story, error)
	GetStoriesByAuthor(ctx context.Context, author string) ([]*models.Story, error)
	GetStoriesByDateRange(ctx context.Context, start, end int64) ([]*models.Story, error)

	// UPDATE OPERATIONS
	UpdateStory(ctx context.Context, story *models.Story) error
	UpdateStoriesCommentsCount(ctx context.Context, id int, count int) error
	UpdateStoryScore(ctx context.Context, id int, score int) error

	// DELETE OPERATIONS
	DeleteStory(ctx context.Context, id int) error
	DeleteStoriesByAuthor(ctx context.Context, author string) error

	// UTILITY OPERATIONS
	StoryExists(ctx context.Context, id int) (bool, error)
	GetStoryCount(ctx context.Context) (int, error)
}

type CommentRepository interface {
	// CREATE OPERATIONS
	CreateComment(ctx context.Context, comment *models.Comment) error

	// READ OPERATIONS
	GetCommentByID(ctx context.Context, id int) (*models.Comment, error)
	GetAllComments(ctx context.Context) ([]*models.Comment, error)
	GetCommentsByStoryID(ctx context.Context, storyID int) ([]*models.Comment, error)
	GetRecentComments(ctx context.Context, limit int) ([]*models.Comment, error)
	GetCommentsByAuthor(ctx context.Context, author string) ([]*models.Comment, error)
	GetCommentsByDateRange(ctx context.Context, start, end int64) ([]*models.Comment, error)

	// UPDATE OPERATIONS
	UpdateComment(ctx context.Context, comment *models.Comment) error
	UpdateCommentsScore(ctx context.Context, id int, score int) error

	// DELETE OPERATIONS
	DeleteComment(ctx context.Context, id int) error
	DeleteCommentsByAuthor(ctx context.Context, author string) error

	// UTILITY OPERATIONS
	CommentExists(ctx context.Context, id int) (bool, error)
	GetCommentCount(ctx context.Context) (int, error)
}

type AskRepository interface {
	// CREATE OPERATIONS
	CreateAsk(ctx context.Context, ask *models.Ask) error

	// READ OPERATIONS
	GetAskByID(ctx context.Context, id int) (*models.Ask, error)
	GetAllAsks(ctx context.Context) ([]*models.Ask, error)
	GetRecentAsks(ctx context.Context, limit int) ([]*models.Ask, error)
	GetAsksByScore(ctx context.Context, minScore int) ([]*models.Ask, error)
	GetAsksByAuthor(ctx context.Context, author string) ([]*models.Ask, error)
	GetAsksByDateRange(ctx context.Context, start, end int64) ([]*models.Ask, error)

	// UPDATE OPERATIONS
	UpdateAsk(ctx context.Context, ask *models.Ask) error
	UpdateAsksCommentsCount(ctx context.Context, id int, count int) error
	UpdateAskScore(ctx context.Context, id int, score int) error

	// DELETE OPERATIONS
	DeleteAsk(ctx context.Context, id int) error
	DeleteAsksByAuthor(ctx context.Context, author string) error

	// UTILITY OPERATIONS
	AskExists(ctx context.Context, id int) (bool, error)
	GetAskCount(ctx context.Context) (int, error)
}

type JobRepository interface {
	// CREATE OPERATIONS
	CreateJob(ctx context.Context, job *models.Job) error

	// READ OPERATIONS
	GetJobByID(ctx context.Context, id int) (*models.Job, error)
	GetAllJobs(ctx context.Context) ([]*models.Job, error)
	GetRecentJobs(ctx context.Context, limit int) ([]*models.Job, error)
	GetJobsByScore(ctx context.Context, minScore int) ([]*models.Job, error)
	GetJobsByAuthor(ctx context.Context, author string) ([]*models.Job, error)
	GetJobsByDateRange(ctx context.Context, start, end int64) ([]*models.Job, error)

	// UPDATE OPERATIONS
	UpdateJob(ctx context.Context, job *models.Job) error
	UpdateJobsScore(ctx context.Context, id int, score int) error

	// DELETE OPERATIONS
	DeleteJob(ctx context.Context, id int) error
	DeleteJobsByAuthor(ctx context.Context, author string) error

	// UTILITY OPERATIONS
	JobExists(ctx context.Context, id int) (bool, error)
	GetJobCount(ctx context.Context) (int, error)
}

type PollRepository interface {
	// CREATE OPERATIONS
	CreatePoll(ctx context.Context, poll *models.Poll) error

	// READ OPERATIONS
	GetPollByID(ctx context.Context, id int) (*models.Poll, error)
	GetAllPolls(ctx context.Context) ([]*models.Poll, error)
	GetRecentPolls(ctx context.Context, limit int) ([]*models.Poll, error)
	GetPollsByScore(ctx context.Context, minScore int) ([]*models.Poll, error)
	GetPollsByAuthor(ctx context.Context, author string) ([]*models.Poll, error)
	GetPollsByDateRange(ctx context.Context, start, end int64) ([]*models.Poll, error)

	// UPDATE OPERATIONS
	UpdatePoll(ctx context.Context, poll *models.Poll) error
	UpdatePollsScore(ctx context.Context, id int, score int) error

	// DELETE OPERATIONS
	DeletePoll(ctx context.Context, id int) error
	DeletePollsByAuthor(ctx context.Context, author string) error

	// UTILITY OPERATIONS
	PollExists(ctx context.Context, id int) (bool, error)
	GetPollCount(ctx context.Context) (int, error)
}
