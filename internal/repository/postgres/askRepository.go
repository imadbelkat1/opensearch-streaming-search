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
	insertAskQuery = `
		INSERT INTO asks (id, type, title, text, score, author, reply_ids, replies_count, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	selectAskFields = `
		SELECT id, type, title, text, score, author, reply_ids, replies_count, created_at`

	selectAskByIDQuery = selectAskFields + ` FROM asks WHERE id = $1`

	selectAllAsksQuery = selectAskFields + ` FROM asks ORDER BY created_at DESC`

	selectRecentAsksQuery = selectAskFields + ` FROM asks ORDER BY created_at DESC LIMIT $1`

	selectAsksByMinScoreQuery = selectAskFields + ` FROM asks WHERE score >= $1 ORDER BY score DESC, created_at DESC`

	selectAsksByAuthorQuery = selectAskFields + ` FROM asks WHERE author = $1 ORDER BY created_at DESC`

	selectAsksByDateRangeQuery = selectAskFields + ` FROM asks WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`

	updateAskQuery = `
		UPDATE asks SET type = $2, title = $3, text = $4, score = $5, author = $6, 
		reply_ids = $7, replies_count = $8, created_at = $9 WHERE id = $1`

	updateRepliesCountQuery = `UPDATE asks SET replies_count = $1 WHERE id = $2`

	updateScoreQuery = `UPDATE asks SET score = $1 WHERE id = $2`

	deleteAskQuery = `DELETE FROM asks WHERE id = $1`

	deleteAsksByAuthorQuery = `DELETE FROM asks WHERE author = $1`

	askExistsQuery = `SELECT EXISTS(SELECT 1 FROM asks WHERE id = $1)`

	askCountQuery = `SELECT COUNT(*) FROM asks`
)

// AskRepository implements the AskRepository interface for PostgreSQL
type AskRepository struct {
	db *sql.DB
}

// NewAskRepository creates a new instance of AskRepository
func NewAskRepository() *AskRepository {
	return &AskRepository{
		db: database.GetDB(),
	}
}

// scanAsk scans a single ask from a row scanner
func (r *AskRepository) scanAsk(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Ask, error) {
	ask := &models.Ask{}
	err := scanner.Scan(
		&ask.ID,
		&ask.Type,
		&ask.Title,
		&ask.Text,
		&ask.Score,
		&ask.Author,
		&ask.Reply_ids,
		&ask.Replies_count,
		&ask.Created_At,
	)
	if err != nil {
		return nil, err
	}
	return ask, nil
}

// scanAsks scans multiple asks from rows
func (r *AskRepository) scanAsks(rows *sql.Rows) ([]*models.Ask, error) {
	var asks []*models.Ask

	for rows.Next() {
		ask, err := r.scanAsk(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ask: %w", err)
		}

		if ask.IsValid() {
			asks = append(asks, ask)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return asks, nil
}

// validateAsk performs basic validation on ask data
func (r *AskRepository) validateAsk(ask *models.Ask) error {
	if ask == nil {
		return fmt.Errorf("ask cannot be nil")
	}
	if ask.ID <= 0 {
		return fmt.Errorf("ask ID must be positive")
	}
	if ask.Author == "" {
		return fmt.Errorf("ask author cannot be empty")
	}
	if ask.Title == "" {
		return fmt.Errorf("ask title cannot be empty")
	}
	return nil
}

// CREATE OPERATIONS

// CreateAsk inserts a new ask into the database
func (r *AskRepository) CreateAsk(ctx context.Context, ask *models.Ask) error {
	if err := r.validateAsk(ask); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := r.db.ExecContext(ctx, insertAskQuery,
		ask.ID,
		ask.Type,
		ask.Title,
		ask.Text,
		ask.Score,
		ask.Author,
		ask.Reply_ids,
		ask.Replies_count,
		ask.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create ask: %w", err)
	}
	return nil
}

// CreateAskWithTransaction inserts a new ask using the provided transaction
func (r *AskRepository) CreateAskWithTransaction(tx *sql.Tx, ask *models.Ask) error {
	if err := r.validateAsk(ask); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := tx.Exec(insertAskQuery,
		ask.ID,
		ask.Type,
		ask.Title,
		ask.Text,
		ask.Score,
		ask.Author,
		ask.Reply_ids,
		ask.Replies_count,
		ask.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create ask in transaction: %w", err)
	}
	return nil
}

// READ OPERATIONS

// GetAskByID retrieves an ask by its ID
func (r *AskRepository) GetAskByID(ctx context.Context, id int) (*models.Ask, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ask ID: %d", id)
	}

	ask, err := r.scanAsk(r.db.QueryRowContext(ctx, selectAskByIDQuery, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ask with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get ask: %w", err)
	}
	return ask, nil
}

// GetAllAsks retrieves all asks from the database
func (r *AskRepository) GetAllAsks(ctx context.Context) ([]*models.Ask, error) {
	rows, err := r.db.QueryContext(ctx, selectAllAsksQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get all asks: %w", err)
	}
	defer rows.Close()

	return r.scanAsks(rows)
}

// GetRecentAsks retrieves the most recent asks, limited by the specified count
func (r *AskRepository) GetRecentAsks(ctx context.Context, limit int) ([]*models.Ask, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got: %d", limit)
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, err := r.db.QueryContext(ctx, selectRecentAsksQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent asks: %w", err)
	}
	defer rows.Close()

	return r.scanAsks(rows)
}

// GetAsksByScore retrieves asks with score >= minScore
func (r *AskRepository) GetAsksByScore(ctx context.Context, minScore int) ([]*models.Ask, error) {
	rows, err := r.db.QueryContext(ctx, selectAsksByMinScoreQuery, minScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get asks by score: %w", err)
	}
	defer rows.Close()

	return r.scanAsks(rows)
}

// GetAsksByAuthor retrieves asks by a specific author
func (r *AskRepository) GetAsksByAuthor(ctx context.Context, author string) ([]*models.Ask, error) {
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}

	rows, err := r.db.QueryContext(ctx, selectAsksByAuthorQuery, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get asks by author: %w", err)
	}
	defer rows.Close()

	return r.scanAsks(rows)
}

// GetAsksByDateRange retrieves asks created within a specific date range
func (r *AskRepository) GetAsksByDateRange(ctx context.Context, start, end int64) ([]*models.Ask, error) {
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("start and end timestamps must be non-negative")
	}
	if start > end {
		return nil, fmt.Errorf("start timestamp (%d) cannot be greater than end timestamp (%d)", start, end)
	}

	rows, err := r.db.QueryContext(ctx, selectAsksByDateRangeQuery, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get asks by date range: %w", err)
	}
	defer rows.Close()

	return r.scanAsks(rows)
}

// UPDATE OPERATIONS

// UpdateAsk updates an existing ask in the database
func (r *AskRepository) UpdateAsk(ctx context.Context, ask *models.Ask) error {
	if err := r.validateAsk(ask); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result, err := r.db.ExecContext(ctx, updateAskQuery,
		ask.ID,
		ask.Type,
		ask.Title,
		ask.Text,
		ask.Score,
		ask.Author,
		ask.Reply_ids,
		ask.Replies_count,
		ask.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to update ask: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ask with ID %d not found", ask.ID)
	}

	return nil
}

// UpdateAsksCommentsCount updates the replies count for a specific ask
func (r *AskRepository) UpdateAsksCommentsCount(ctx context.Context, id int, count int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ask ID: %d", id)
	}
	if count < 0 {
		return fmt.Errorf("replies count cannot be negative: %d", count)
	}

	result, err := r.db.ExecContext(ctx, updateRepliesCountQuery, count, id)
	if err != nil {
		return fmt.Errorf("failed to update replies count for ask ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ask with ID %d not found", id)
	}

	return nil
}

// UpdateAskScore updates the score of a specific ask
func (r *AskRepository) UpdateAskScore(ctx context.Context, id int, score int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ask ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, updateScoreQuery, score, id)
	if err != nil {
		return fmt.Errorf("failed to update score for ask ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ask with ID %d not found", id)
	}

	return nil
}

// DELETE OPERATIONS

// DeleteAsk removes an ask from the database by its ID
func (r *AskRepository) DeleteAsk(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ask ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, deleteAskQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete ask with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ask with ID %d not found", id)
	}

	return nil
}

// DeleteAsksByAuthor removes all asks made by a specific author
func (r *AskRepository) DeleteAsksByAuthor(ctx context.Context, author string) error {
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, deleteAsksByAuthorQuery, author)
	if err != nil {
		return fmt.Errorf("failed to delete asks by author %s: %w", author, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no asks found for author %s", author)
	}

	return nil
}

// UTILITY OPERATIONS

// AskExists checks if an ask exists in the database by its ID
func (r *AskRepository) AskExists(ctx context.Context, id int) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid ask ID: %d", id)
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, askExistsQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if ask exists: %w", err)
	}
	return exists, nil
}

// GetAskCount retrieves the total number of asks in the database
func (r *AskRepository) GetAskCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, askCountQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get ask count: %w", err)
	}
	return count, nil
}

// BATCH OPERATIONS

// CreateAsksInBatch creates multiple asks in a single transaction
func (r *AskRepository) CreateAsksInBatch(ctx context.Context, asks []*models.Ask) error {
	if len(asks) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, ask := range asks {
		if err := r.CreateAskWithTransaction(tx, ask); err != nil {
			return fmt.Errorf("failed to create ask %d in batch: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch transaction: %w", err)
	}

	return nil
}
