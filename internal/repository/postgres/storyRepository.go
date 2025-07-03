package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"internship-project/internal/models"
	"internship-project/pkg/database"
)

// SQL query constants
const (
	insertStoryQuery = `
		INSERT INTO stories (id, type, title, url, score, author, created_at, comments_count) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	selectStoryFields = `
		SELECT id, type, title, url, score, author, created_at, comments_count`

	selectStoryByIDQuery = selectStoryFields + ` FROM stories WHERE id = $1`

	selectAllStoriesQuery = selectStoryFields + ` FROM stories ORDER BY created_at DESC`

	selectRecentStoriesQuery = selectStoryFields + ` FROM stories ORDER BY created_at DESC LIMIT $1`

	selectStoriesByMinScoreQuery = selectStoryFields + ` FROM stories WHERE score >= $1 ORDER BY score DESC, created_at DESC`

	selectStoriesByAuthorQuery = selectStoryFields + ` FROM stories WHERE author = $1 ORDER BY created_at DESC`

	selectStoriesByDateRangeQuery = selectStoryFields + ` FROM stories WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`

	updateStoryQuery = `
		UPDATE stories SET type = $2, title = $3, url = $4, score = $5, author = $6, 
		created_at = $7, comments_count = $8 WHERE id = $1`

	updateCommentsCountQuery = `UPDATE stories SET comments_count = $1 WHERE id = $2`

	updateStoryScoreQuery = `UPDATE stories SET score = $1 WHERE id = $2`

	deleteStoryQuery = `DELETE FROM stories WHERE id = $1`

	deleteStoriesByAuthorQuery = `DELETE FROM stories WHERE author = $1`

	storyExistsQuery = `SELECT EXISTS(SELECT 1 FROM stories WHERE id = $1)`

	storyCountQuery = `SELECT COUNT(*) FROM stories`
)

// StoryRepository implements the StoryRepository interface for PostgreSQL
type StoryRepository struct {
	db *sql.DB
}

// NewStoryRepository creates a new instance of StoryRepository
func NewStoryRepository() *StoryRepository {
	return &StoryRepository{
		db: database.GetDB(),
	}
}

// scanStory scans a single story from a row scanner
func (r *StoryRepository) scanStory(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Story, error) {
	story := &models.Story{}
	err := scanner.Scan(
		&story.ID,
		&story.Type,
		&story.Title,
		&story.URL,
		&story.Score,
		&story.Author,
		&story.Created_At,
		&story.Comments_count,
	)
	if err != nil {
		return nil, err
	}
	return story, nil
}

// scanStories scans multiple stories from rows
func (r *StoryRepository) scanStories(rows *sql.Rows) ([]*models.Story, error) {
	var stories []*models.Story

	for rows.Next() {
		story, err := r.scanStory(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan story: %w", err)
		}
		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return stories, nil
}

// validateStory performs basic validation on story data
func (r *StoryRepository) validateStory(story *models.Story) error {
	if story == nil {
		return fmt.Errorf("story cannot be nil")
	}
	if story.ID <= 0 {
		return fmt.Errorf("story ID must be positive")
	}
	if story.Author == "" {
		return fmt.Errorf("story author cannot be empty")
	}
	if story.Title == "" {
		return fmt.Errorf("story title cannot be empty")
	}
	return nil
}

// CREATE OPERATIONS

// CreateStory inserts a new story into the database
func (r *StoryRepository) CreateStory(ctx context.Context, story *models.Story) error {
	if err := r.validateStory(story); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := r.db.ExecContext(ctx, insertStoryQuery,
		story.ID,
		story.Type,
		story.Title,
		story.URL,
		story.Score,
		story.Author,
		story.Created_At,
		story.Comments_count,
	)

	if err != nil {
		return fmt.Errorf("failed to create story: %w", err)
	}
	return nil
}

// CreateStoryWithTransaction inserts a new story using the provided transaction
func (r *StoryRepository) CreateStoryWithTransaction(tx *sql.Tx, story *models.Story) error {
	if err := r.validateStory(story); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := tx.Exec(insertStoryQuery,
		story.ID,
		story.Type,
		story.Title,
		story.URL,
		story.Score,
		story.Author,
		story.Created_At,
		story.Comments_count,
	)

	if err != nil {
		return fmt.Errorf("failed to create story in transaction: %w", err)
	}
	return nil
}

// READ OPERATIONS

// GetStoryByID retrieves a story by its ID
func (r *StoryRepository) GetStoryByID(ctx context.Context, id int) (*models.Story, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid story ID: %d", id)
	}

	story, err := r.scanStory(r.db.QueryRowContext(ctx, selectStoryByIDQuery, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("story with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get story: %w", err)
	}
	return story, nil
}

// GetAllStories retrieves all stories from the database
func (r *StoryRepository) GetAllStories(ctx context.Context) ([]*models.Story, error) {
	rows, err := r.db.QueryContext(ctx, selectAllStoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stories: %w", err)
	}
	defer rows.Close()

	return r.scanStories(rows)
}

// GetRecentStories retrieves the most recent stories, limited by the specified count
func (r *StoryRepository) GetRecentStories(ctx context.Context, limit int) ([]*models.Story, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got: %d", limit)
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, err := r.db.QueryContext(ctx, selectRecentStoriesQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent stories: %w", err)
	}
	defer rows.Close()

	return r.scanStories(rows)
}

// GetStoriesByMinScore retrieves stories with score >= minScore
func (r *StoryRepository) GetStoriesByMinScore(ctx context.Context, minScore int) ([]*models.Story, error) {
	rows, err := r.db.QueryContext(ctx, selectStoriesByMinScoreQuery, minScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get stories by score: %w", err)
	}
	defer rows.Close()

	return r.scanStories(rows)
}

// GetStoriesByAuthor retrieves stories by a specific author
func (r *StoryRepository) GetStoriesByAuthor(ctx context.Context, author string) ([]*models.Story, error) {
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}

	rows, err := r.db.QueryContext(ctx, selectStoriesByAuthorQuery, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get stories by author: %w", err)
	}
	defer rows.Close()

	return r.scanStories(rows)
}

// GetStoriesByDateRange retrieves stories created within a specific date range
func (r *StoryRepository) GetStoriesByDateRange(ctx context.Context, start, end int64) ([]*models.Story, error) {
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("start and end timestamps must be non-negative")
	}
	if start > end {
		return nil, fmt.Errorf("start timestamp (%d) cannot be greater than end timestamp (%d)", start, end)
	}

	rows, err := r.db.QueryContext(ctx, selectStoriesByDateRangeQuery, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get stories by date range: %w", err)
	}
	defer rows.Close()

	return r.scanStories(rows)
}

// UPDATE OPERATIONS

// UpdateStory updates an existing story in the database
func (r *StoryRepository) UpdateStory(ctx context.Context, story *models.Story) error {
	if err := r.validateStory(story); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result, err := r.db.ExecContext(ctx, updateStoryQuery,
		story.ID,
		story.Type,
		story.Title,
		story.URL,
		story.Score,
		story.Author,
		story.Created_At,
		story.Comments_count,
	)

	if err != nil {
		return fmt.Errorf("failed to update story: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("story with ID %d not found", story.ID)
	}

	return nil
}

// UpdateStoriesCommentsCount updates the comments count for a specific story
func (r *StoryRepository) UpdateStoriesCommentsCount(ctx context.Context, id int, count int) error {
	if id <= 0 {
		return fmt.Errorf("invalid story ID: %d", id)
	}
	if count < 0 {
		return fmt.Errorf("comments count cannot be negative: %d", count)
	}

	result, err := r.db.ExecContext(ctx, updateCommentsCountQuery, count, id)
	if err != nil {
		return fmt.Errorf("failed to update comments count for story ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("story with ID %d not found", id)
	}

	return nil
}

// UpdateStoryScore updates the score of a specific story
func (r *StoryRepository) UpdateStoryScore(ctx context.Context, id int, score int) error {
	if id <= 0 {
		return fmt.Errorf("invalid story ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, updateStoryScoreQuery, score, id)
	if err != nil {
		return fmt.Errorf("failed to update score for story ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("story with ID %d not found", id)
	}

	return nil
}

// DELETE OPERATIONS

// DeleteStory removes a story from the database by its ID
func (r *StoryRepository) DeleteStory(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid story ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, deleteStoryQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete story with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("story with ID %d not found", id)
	}

	return nil
}

// DeleteStoriesByAuthor removes all stories made by a specific author
func (r *StoryRepository) DeleteStoriesByAuthor(ctx context.Context, author string) error {
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, deleteStoriesByAuthorQuery, author)
	if err != nil {
		return fmt.Errorf("failed to delete stories by author %s: %w", author, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no stories found for author %s", author)
	}

	return nil
}

// UTILITY OPERATIONS

// StoryExists checks if a story exists in the database by its ID
func (r *StoryRepository) StoryExists(ctx context.Context, id int) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid story ID: %d", id)
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, storyExistsQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if story exists: %w", err)
	}
	return exists, nil
}

// GetStoryCount retrieves the total number of stories in the database
func (r *StoryRepository) GetStoryCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, storyCountQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get story count: %w", err)
	}
	return count, nil
}

// BATCH OPERATIONS

// CreateStoriesInBatch creates multiple stories in a single transaction
func (r *StoryRepository) CreateStoriesInBatch(ctx context.Context, stories []*models.Story) error {
	if len(stories) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, story := range stories {
		if err := r.CreateStoryWithTransaction(tx, story); err != nil {
			return fmt.Errorf("failed to create story %d in batch: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch transaction: %w", err)
	}

	return nil
}
