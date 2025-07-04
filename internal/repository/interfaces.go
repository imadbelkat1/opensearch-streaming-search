package repository

import (
	"context"
	"internship-project/internal/models"
)

// Common interface methods used across all repositories
type BaseRepositoryInterface interface {
	Exists(ctx context.Context, id int) (bool, error)
	GetCount(ctx context.Context) (int, error)
}

type StoryRepository interface {
	BaseRepositoryInterface

	// CRUD Operations
	Create(ctx context.Context, story *models.Story) error
	GetByID(ctx context.Context, id int) (*models.Story, error)
	Update(ctx context.Context, story *models.Story) error
	Delete(ctx context.Context, id int) error

	// Query Operations
	GetAll(ctx context.Context) ([]*models.Story, error)
	GetRecent(ctx context.Context, limit int) ([]*models.Story, error)
	GetByMinScore(ctx context.Context, minScore int) ([]*models.Story, error)
	GetByAuthor(ctx context.Context, author string) ([]*models.Story, error)
	GetByDateRange(ctx context.Context, start, end int64) ([]*models.Story, error)

	// Update specific fields
	UpdateScore(ctx context.Context, id int, score int) error
	UpdateCommentsCount(ctx context.Context, id int, count int) error

	// Batch operations
	CreateBatch(ctx context.Context, stories []*models.Story) error
	CreateBatchWithExistingIDs(ctx context.Context, stories []*models.Story) error
	DeleteByAuthor(ctx context.Context, author string) error
}

type CommentRepository interface {
	BaseRepositoryInterface

	// CRUD Operations
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id int) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id int) error

	// Query Operations
	GetAll(ctx context.Context) ([]*models.Comment, error)
	GetRecent(ctx context.Context, limit int) ([]*models.Comment, error)
	GetByAuthor(ctx context.Context, author string) ([]*models.Comment, error)
	GetByDateRange(ctx context.Context, start, end int64) ([]*models.Comment, error)

	// Batch operations
	DeleteByAuthor(ctx context.Context, author string) error
}

type AskRepository interface {
	BaseRepositoryInterface

	// CRUD Operations
	Create(ctx context.Context, ask *models.Ask) error
	GetByID(ctx context.Context, id int) (*models.Ask, error)
	Update(ctx context.Context, ask *models.Ask) error
	Delete(ctx context.Context, id int) error

	// Query Operations
	GetAll(ctx context.Context) ([]*models.Ask, error)
	GetRecent(ctx context.Context, limit int) ([]*models.Ask, error)
	GetByMinScore(ctx context.Context, minScore int) ([]*models.Ask, error)
	GetByAuthor(ctx context.Context, author string) ([]*models.Ask, error)
	GetByDateRange(ctx context.Context, start, end int64) ([]*models.Ask, error)

	// Update specific fields
	UpdateScore(ctx context.Context, id int, score int) error
	UpdateRepliesCount(ctx context.Context, id int, count int) error

	// Batch operations
	CreateBatch(ctx context.Context, asks []*models.Ask) error
	DeleteByAuthor(ctx context.Context, author string) error
}

type JobRepository interface {
	BaseRepositoryInterface

	// CRUD Operations
	Create(ctx context.Context, job *models.Job) error
	GetByID(ctx context.Context, id int) (*models.Job, error)
	Update(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, id int) error

	// Query Operations
	GetAll(ctx context.Context) ([]*models.Job, error)
	GetRecent(ctx context.Context, limit int) ([]*models.Job, error)
	GetByMinScore(ctx context.Context, minScore int) ([]*models.Job, error)
	GetByAuthor(ctx context.Context, author string) ([]*models.Job, error)
	GetByDateRange(ctx context.Context, start, end int64) ([]*models.Job, error)

	// Update specific fields
	UpdateScore(ctx context.Context, id int, score int) error

	// Batch operations
	CreateBatch(ctx context.Context, jobs []*models.Job) error
	DeleteByAuthor(ctx context.Context, author string) error
}

type PollRepository interface {
	BaseRepositoryInterface

	// CRUD Operations
	Create(ctx context.Context, poll *models.Poll) error
	GetByID(ctx context.Context, id int) (*models.Poll, error)
	Update(ctx context.Context, poll *models.Poll) error
	Delete(ctx context.Context, id int) error

	// Query Operations
	GetAll(ctx context.Context) ([]*models.Poll, error)
	GetRecent(ctx context.Context, limit int) ([]*models.Poll, error)
	GetByMinScore(ctx context.Context, minScore int) ([]*models.Poll, error)
	GetByAuthor(ctx context.Context, author string) ([]*models.Poll, error)
	GetByDateRange(ctx context.Context, start, end int64) ([]*models.Poll, error)

	// Update specific fields
	UpdateScore(ctx context.Context, id int, score int) error

	// Batch operations
	CreateBatch(ctx context.Context, polls []*models.Poll) error
	DeleteByAuthor(ctx context.Context, author string) error
}

type PollOptionRepository interface {
	BaseRepositoryInterface

	// CRUD Operations
	Create(ctx context.Context, pollOption *models.PollOption) error
	GetByID(ctx context.Context, id int) (*models.PollOption, error)
	Update(ctx context.Context, pollOption *models.PollOption) error
	Delete(ctx context.Context, id int) error

	// Query Operations
	GetAll(ctx context.Context) ([]*models.PollOption, error)
	GetByPollID(ctx context.Context, pollID int) ([]*models.PollOption, error)
	GetRecent(ctx context.Context, limit int) ([]*models.PollOption, error)
	GetByAuthor(ctx context.Context, author string) ([]*models.PollOption, error)
	GetByDateRange(ctx context.Context, start, end int64) ([]*models.PollOption, error)
	IncrementVotes(ctx context.Context, id int) error
	GetVoteCount(ctx context.Context, id int) (int, error)
	CountByPollID(ctx context.Context, pollID int) (int, error)
	GetTopVoted(ctx context.Context, pollID int, limit int) ([]*models.PollOption, error)

	// Update specific fields
	UpdateVotes(ctx context.Context, id int, votes int) error

	// Batch operations
	CreateBatch(ctx context.Context, pollOptions []*models.PollOption) error
	DeleteByAuthor(ctx context.Context, author string) error
	DeleteByPollID(ctx context.Context, pollID int) error
}
