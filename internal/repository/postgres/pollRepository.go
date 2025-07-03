package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"internship-project/internal/models"
	"internship-project/pkg/database"

	"github.com/lib/pq"
)

// SQL query constants
const (
	insertPollQuery = `
		INSERT INTO polls (id, type, title, score, author, poll_options, reply_ids, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	selectPollFields = `
		SELECT id, type, title, score, author, poll_options, reply_ids, created_at`

	selectPollByIDQuery = selectPollFields + ` FROM polls WHERE id = $1`

	selectAllPollsQuery = selectPollFields + ` FROM polls ORDER BY created_at DESC`

	selectRecentPollsQuery = selectPollFields + ` FROM polls ORDER BY created_at DESC LIMIT $1`

	selectPollsByMinScoreQuery = selectPollFields + ` FROM polls WHERE score >= $1 ORDER BY score DESC, created_at DESC`

	selectPollsByAuthorQuery = selectPollFields + ` FROM polls WHERE author = $1 ORDER BY created_at DESC`

	selectPollsByDateRangeQuery = selectPollFields + ` FROM polls WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`

	updatePollQuery = `
		UPDATE polls SET type = $2, title = $3, score = $4, author = $5, 
		poll_options = $6, reply_ids = $7, created_at = $8 WHERE id = $1`

	updatePollScoreQuery = `UPDATE polls SET score = $1 WHERE id = $2`

	deletePollQuery = `DELETE FROM polls WHERE id = $1`

	deletePollsByAuthorQuery = `DELETE FROM polls WHERE author = $1`

	pollExistsQuery = `SELECT EXISTS(SELECT 1 FROM polls WHERE id = $1)`

	pollCountQuery = `SELECT COUNT(*) FROM polls`
)

// PollRepository implements the PollRepository interface for PostgreSQL
type PollRepository struct {
	db *sql.DB
}

// NewPollRepository creates a new instance of PollRepository
func NewPollRepository() *PollRepository {
	return &PollRepository{
		db: database.GetDB(),
	}
}

// scanPoll scans a single poll from a row scanner
func (r *PollRepository) scanPoll(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Poll, error) {
	poll := &models.Poll{}
	var pollOptions pq.Int64Array

	err := scanner.Scan(
		&poll.ID,
		&poll.Type,
		&poll.Title,
		&poll.Score,
		&poll.Author,
		&poll.Parts,
		&poll.Reply_Ids,
		&poll.Created_At,
	)
	if err != nil {
		return nil, err
	}

	// Convert pq.Int64Array to []int for Parts field
	poll.Parts = make([]int, len(pollOptions))
	for i, v := range pollOptions {
		poll.Parts[i] = int(v)
	}

	return poll, nil
}

// scanPolls scans multiple polls from rows
func (r *PollRepository) scanPolls(rows *sql.Rows) ([]*models.Poll, error) {
	var polls []*models.Poll

	for rows.Next() {
		poll, err := r.scanPoll(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan poll: %w", err)
		}

		if poll.IsValid() {
			polls = append(polls, poll)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return polls, nil
}

// validatePoll performs basic validation on poll data
func (r *PollRepository) validatePoll(poll *models.Poll) error {
	if poll == nil {
		return fmt.Errorf("poll cannot be nil")
	}
	if poll.ID <= 0 {
		return fmt.Errorf("poll ID must be positive")
	}
	if poll.Author == "" {
		return fmt.Errorf("poll author cannot be empty")
	}
	if poll.Title == "" {
		return fmt.Errorf("poll title cannot be empty")
	}
	return nil
}

// CREATE OPERATIONS
// CreatePoll inserts a new poll into the database
func (r *PollRepository) CreatePoll(ctx context.Context, poll *models.Poll) error {
	if err := r.validatePoll(poll); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Convert []int to pq.Int64Array for database storage
	pollOptions := make(pq.Int64Array, len(poll.Parts))
	for i, v := range poll.Parts {
		pollOptions[i] = int64(v)
	}

	_, err := r.db.ExecContext(ctx, insertPollQuery,
		poll.ID,
		poll.Type,
		poll.Title,
		poll.Score,
		poll.Author,
		poll.Parts,
		poll.Reply_Ids,
		poll.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create poll: %w", err)
	}
	return nil
}

// CreatePollWithTransaction inserts a new poll using the provided transaction
func (r *PollRepository) CreatePollWithTransaction(tx *sql.Tx, poll *models.Poll) error {
	if err := r.validatePoll(poll); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Convert []int to pq.Int64Array for database storage
	pollOptions := make(pq.Int64Array, len(poll.Parts))
	for i, v := range poll.Parts {
		pollOptions[i] = int64(v)
	}

	_, err := tx.Exec(insertPollQuery,
		poll.ID,
		poll.Type,
		poll.Title,
		poll.Score,
		poll.Author,
		poll.Parts,
		poll.Reply_Ids,
		poll.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create poll in transaction: %w", err)
	}
	return nil
}

// READ OPERATIONS
// GetPollByID retrieves a poll by its ID
func (r *PollRepository) GetPollByID(ctx context.Context, id int) (*models.Poll, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid poll ID: %d", id)
	}

	poll, err := r.scanPoll(r.db.QueryRowContext(ctx, selectPollByIDQuery, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("poll with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}
	return poll, nil
}

// GetAllPolls retrieves all polls from the database
func (r *PollRepository) GetAllPolls(ctx context.Context) ([]*models.Poll, error) {
	rows, err := r.db.QueryContext(ctx, selectAllPollsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get all polls: %w", err)
	}
	defer rows.Close()

	return r.scanPolls(rows)
}

// GetRecentPolls retrieves the most recent polls, limited by the specified count
func (r *PollRepository) GetRecentPolls(ctx context.Context, limit int) ([]*models.Poll, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got: %d", limit)
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, err := r.db.QueryContext(ctx, selectRecentPollsQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent polls: %w", err)
	}
	defer rows.Close()

	return r.scanPolls(rows)
}

// GetPollsByScore retrieves polls with score >= minScore
func (r *PollRepository) GetPollsByScore(ctx context.Context, minScore int) ([]*models.Poll, error) {
	rows, err := r.db.QueryContext(ctx, selectPollsByMinScoreQuery, minScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get polls by score: %w", err)
	}
	defer rows.Close()

	return r.scanPolls(rows)
}

// GetPollsByAuthor retrieves polls by a specific author
func (r *PollRepository) GetPollsByAuthor(ctx context.Context, author string) ([]*models.Poll, error) {
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}

	rows, err := r.db.QueryContext(ctx, selectPollsByAuthorQuery, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get polls by author: %w", err)
	}
	defer rows.Close()

	return r.scanPolls(rows)
}

// GetPollsByDateRange retrieves polls created within a specific date range
func (r *PollRepository) GetPollsByDateRange(ctx context.Context, start, end int64) ([]*models.Poll, error) {
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("start and end timestamps must be non-negative")
	}
	if start > end {
		return nil, fmt.Errorf("start timestamp (%d) cannot be greater than end timestamp (%d)", start, end)
	}

	rows, err := r.db.QueryContext(ctx, selectPollsByDateRangeQuery, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get polls by date range: %w", err)
	}
	defer rows.Close()

	return r.scanPolls(rows)
}

// UPDATE OPERATIONS
// UpdatePoll updates an existing poll in the database
func (r *PollRepository) UpdatePoll(ctx context.Context, poll *models.Poll) error {
	if err := r.validatePoll(poll); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Convert []int to pq.Int64Array for database storage
	pollOptions := make(pq.Int64Array, len(poll.Parts))
	for i, v := range poll.Parts {
		pollOptions[i] = int64(v)
	}

	result, err := r.db.ExecContext(ctx, updatePollQuery,
		poll.ID,
		poll.Type,
		poll.Title,
		poll.Score,
		poll.Author,
		poll.Parts,
		poll.Reply_Ids,
		poll.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to update poll: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("poll with ID %d not found", poll.ID)
	}

	return nil
}

// UpdatePollsScore updates the score of a specific poll
func (r *PollRepository) UpdatePollsScore(ctx context.Context, id int, score int) error {
	if id <= 0 {
		return fmt.Errorf("invalid poll ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, updatePollScoreQuery, score, id)
	if err != nil {
		return fmt.Errorf("failed to update score for poll ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("poll with ID %d not found", id)
	}

	return nil
}

// DELETE OPERATIONS
// DeletePoll removes a poll from the database by its ID
func (r *PollRepository) DeletePoll(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid poll ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, deletePollQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete poll with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("poll with ID %d not found", id)
	}

	return nil
}

// DeletePollsByAuthor removes all polls made by a specific author
func (r *PollRepository) DeletePollsByAuthor(ctx context.Context, author string) error {
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, deletePollsByAuthorQuery, author)
	if err != nil {
		return fmt.Errorf("failed to delete polls by author %s: %w", author, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no polls found for author %s", author)
	}

	return nil
}

// UTILITY OPERATIONS
// PollExists checks if a poll exists in the database by its ID
func (r *PollRepository) PollExists(ctx context.Context, id int) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid poll ID: %d", id)
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, pollExistsQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if poll exists: %w", err)
	}
	return exists, nil
}

// GetPollCount retrieves the total number of polls in the database
func (r *PollRepository) GetPollCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, pollCountQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get poll count: %w", err)
	}
	return count, nil
}

// BATCH OPERATIONS
// CreatePollsInBatch creates multiple polls in a single transaction
func (r *PollRepository) CreatePollsInBatch(ctx context.Context, polls []*models.Poll) error {
	if len(polls) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, poll := range polls {
		if err := r.CreatePollWithTransaction(tx, poll); err != nil {
			return fmt.Errorf("failed to create poll %d in batch: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch transaction: %w", err)
	}

	return nil
}
