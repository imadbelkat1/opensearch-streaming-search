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
	insertCommentQuery = `
		INSERT INTO comments (id, type, by, text, parent, story_id, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	selectCommentFields = `SELECT id, type, by, text, parent, story_id, created_at`

	selectCommentByIDQuery = selectCommentFields + ` FROM comments WHERE id = $1`

	selectAllCommentsQuery = selectCommentFields + ` FROM comments ORDER BY created_at DESC`

	selectCommentsByStoryIDQuery = selectCommentFields + ` FROM comments WHERE story_id = $1 ORDER BY created_at ASC`

	selectRecentCommentsQuery = selectCommentFields + ` FROM comments ORDER BY created_at DESC LIMIT $1`

	selectCommentsByAuthorQuery = selectCommentFields + ` FROM comments WHERE by = $1 ORDER BY created_at DESC`

	selectCommentsByDateRangeQuery = selectCommentFields + ` FROM comments WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`

	updateCommentQuery = `
		UPDATE comments SET type = $2, by = $3, text = $4, parent = $5, story_id = $6, created_at = $7 
		WHERE id = $1`

	updateCommentScoreQuery = `UPDATE comments SET score = $1 WHERE id = $2`

	deleteCommentQuery = `DELETE FROM comments WHERE id = $1`

	deleteCommentsByAuthorQuery = `DELETE FROM comments WHERE by = $1`

	commentExistsQuery = `SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1)`

	commentCountQuery = `SELECT COUNT(*) FROM comments`
)

// CommentRepository implements the CommentRepository interface for PostgreSQL
type CommentRepository struct {
	db *sql.DB
}

// NewCommentRepository creates a new instance of CommentRepository
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		db: database.GetDB(),
	}
}

// scanComment scans a single comment from a row scanner
func (r *CommentRepository) scanComment(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Comment, error) {
	comment := &models.Comment{}
	err := scanner.Scan(
		&comment.ID,
		&comment.Type,
		&comment.Author,
		&comment.Text,
		&comment.Parent,
		&comment.Replies,
		&comment.Created_At,
	)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// scanComments scans multiple comments from rows
func (r *CommentRepository) scanComments(rows *sql.Rows) ([]*models.Comment, error) {
	var comments []*models.Comment

	for rows.Next() {
		comment, err := r.scanComment(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return comments, nil
}

// validateComment performs basic validation on comment data
func (r *CommentRepository) validateComment(comment *models.Comment) error {
	if comment == nil {
		return fmt.Errorf("comment cannot be null")
	}
	if comment.ID <= 0 {
		return fmt.Errorf("comment ID must be positive")
	}
	if comment.Author == "" {
		return fmt.Errorf("comment author cannot be empty")
	}
	if comment.Text == "" {
		return fmt.Errorf("comment text cannot be empty")
	}
	return nil
}

// CREATE OPERATIONS

// CreateComment inserts a new comment into the database
func (r *CommentRepository) CreateComment(comment *models.Comment) error {
	return r.CreateCommentWithContext(context.Background(), comment)
}

// CreateCommentWithContext inserts a new comment into the database with context
func (r *CommentRepository) CreateCommentWithContext(ctx context.Context, comment *models.Comment) error {
	if err := r.validateComment(comment); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := r.db.ExecContext(ctx, insertCommentQuery,
		comment.ID,
		comment.Type,
		comment.Author,
		comment.Text,
		comment.Parent,
		comment.Replies,
		comment.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

// CreateCommentWithTransaction inserts a new comment using the provided transaction
func (r *CommentRepository) CreateCommentWithTransaction(tx *sql.Tx, comment *models.Comment) error {
	if err := r.validateComment(comment); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := tx.Exec(insertCommentQuery,
		comment.ID,
		comment.Type,
		comment.Author,
		comment.Text,
		comment.Parent,
		comment.Replies,
		comment.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create comment in transaction: %w", err)
	}
	return nil
}

// READ OPERATIONS

// GetCommentByID retrieves a comment by its ID
func (r *CommentRepository) GetCommentByID(id int) (*models.Comment, error) {
	return r.GetCommentByIDWithContext(context.Background(), id)
}

// GetCommentByIDWithContext retrieves a comment by its ID with context
func (r *CommentRepository) GetCommentByIDWithContext(ctx context.Context, id int) (*models.Comment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid comment ID: %d", id)
	}

	comment, err := r.scanComment(r.db.QueryRowContext(ctx, selectCommentByIDQuery, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return comment, nil
}

// GetAllComments retrieves all comments from the database
func (r *CommentRepository) GetAllComments() ([]*models.Comment, error) {
	return r.GetAllCommentsWithContext(context.Background())
}

// GetAllCommentsWithContext retrieves all comments from the database with context
func (r *CommentRepository) GetAllCommentsWithContext(ctx context.Context) ([]*models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, selectAllCommentsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// GetCommentsByStoryID retrieves comments for a specific story by its ID
func (r *CommentRepository) GetCommentsByStoryID(storyID int) ([]*models.Comment, error) {
	return r.GetCommentsByStoryIDWithContext(context.Background(), storyID)
}

// GetCommentsByStoryIDWithContext retrieves comments for a specific story with context
func (r *CommentRepository) GetCommentsByStoryIDWithContext(ctx context.Context, storyID int) ([]*models.Comment, error) {
	if storyID <= 0 {
		return nil, fmt.Errorf("invalid story ID: %d", storyID)
	}

	rows, err := r.db.QueryContext(ctx, selectCommentsByStoryIDQuery, storyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by story ID: %w", err)
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// GetRecentComments retrieves the most recent comments, limited by the specified count
func (r *CommentRepository) GetRecentComments(limit int) ([]*models.Comment, error) {
	return r.GetRecentCommentsWithContext(context.Background(), limit)
}

// GetRecentCommentsWithContext retrieves the most recent comments with context
func (r *CommentRepository) GetRecentCommentsWithContext(ctx context.Context, limit int) ([]*models.Comment, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got: %d", limit)
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, err := r.db.QueryContext(ctx, selectRecentCommentsQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent comments: %w", err)
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// GetCommentsByAuthor retrieves comments made by a specific author
func (r *CommentRepository) GetCommentsByAuthor(author string) ([]*models.Comment, error) {
	return r.GetCommentsByAuthorWithContext(context.Background(), author)
}

// GetCommentsByAuthorWithContext retrieves comments made by a specific author with context
func (r *CommentRepository) GetCommentsByAuthorWithContext(ctx context.Context, author string) ([]*models.Comment, error) {
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}

	rows, err := r.db.QueryContext(ctx, selectCommentsByAuthorQuery, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by author: %w", err)
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// GetCommentsByDateRange retrieves comments created within a specific date range
func (r *CommentRepository) GetCommentsByDateRange(start, end int64) ([]*models.Comment, error) {
	return r.GetCommentsByDateRangeWithContext(context.Background(), start, end)
}

// GetCommentsByDateRangeWithContext retrieves comments created within a specific date range with context
func (r *CommentRepository) GetCommentsByDateRangeWithContext(ctx context.Context, start, end int64) ([]*models.Comment, error) {
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("start and end timestamps must be non-negative")
	}
	if start > end {
		return nil, fmt.Errorf("start timestamp (%d) cannot be greater than end timestamp (%d)", start, end)
	}

	rows, err := r.db.QueryContext(ctx, selectCommentsByDateRangeQuery, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by date range: %w", err)
	}
	defer rows.Close()

	return r.scanComments(rows)
}

// UPDATE OPERATIONS

// UpdateComment updates an existing comment in the database
func (r *CommentRepository) UpdateComment(comment *models.Comment) error {
	return r.UpdateCommentWithContext(context.Background(), comment)
}

// UpdateCommentWithContext updates an existing comment in the database with context
func (r *CommentRepository) UpdateCommentWithContext(ctx context.Context, comment *models.Comment) error {
	if err := r.validateComment(comment); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result, err := r.db.ExecContext(ctx, updateCommentQuery,
		comment.ID,
		comment.Type,
		comment.Author,
		comment.Text,
		comment.Parent,
		comment.Replies,
		comment.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %d not found", comment.ID)
	}

	return nil
}

// UpdateCommentScore updates the score of a specific comment
func (r *CommentRepository) UpdateCommentScore(id int, score int) error {
	return r.UpdateCommentScoreWithContext(context.Background(), id, score)
}

// UpdateCommentScoreWithContext updates the score of a specific comment with context
func (r *CommentRepository) UpdateCommentScoreWithContext(ctx context.Context, id int, score int) error {
	if id <= 0 {
		return fmt.Errorf("invalid comment ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, updateCommentScoreQuery, score, id)
	if err != nil {
		return fmt.Errorf("failed to update comment score for ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %d not found", id)
	}

	return nil
}

// DELETE OPERATIONS

// DeleteComment removes a comment from the database by its ID
func (r *CommentRepository) DeleteComment(id int) error {
	return r.DeleteCommentWithContext(context.Background(), id)
}

// DeleteCommentWithContext removes a comment from the database by its ID with context
func (r *CommentRepository) DeleteCommentWithContext(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid comment ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, deleteCommentQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %d not found", id)
	}

	return nil
}

// DeleteCommentsByAuthor removes all comments made by a specific author
func (r *CommentRepository) DeleteCommentsByAuthor(author string) error {
	return r.DeleteCommentsByAuthorWithContext(context.Background(), author)
}

// DeleteCommentsByAuthorWithContext removes all comments made by a specific author with context
func (r *CommentRepository) DeleteCommentsByAuthorWithContext(ctx context.Context, author string) error {
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, deleteCommentsByAuthorQuery, author)
	if err != nil {
		return fmt.Errorf("failed to delete comments by author %s: %w", author, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no comments found for author %s", author)
	}

	return nil
}

// UTILITY OPERATIONS

// CommentExists checks if a comment exists in the database by its ID
func (r *CommentRepository) CommentExists(id int) (bool, error) {
	return r.CommentExistsWithContext(context.Background(), id)
}

// CommentExistsWithContext checks if a comment exists in the database by its ID with context
func (r *CommentRepository) CommentExistsWithContext(ctx context.Context, id int) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid comment ID: %d", id)
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, commentExistsQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if comment exists: %w", err)
	}
	return exists, nil
}

// GetCommentCount retrieves the total number of comments in the database
func (r *CommentRepository) GetCommentCount() (int, error) {
	return r.GetCommentCountWithContext(context.Background())
}

// GetCommentCountWithContext retrieves the total number of comments in the database with context
func (r *CommentRepository) GetCommentCountWithContext(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, commentCountQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get comment count: %w", err)
	}
	return count, nil
}

// BATCH OPERATIONS

// CreateCommentsInBatch creates multiple comments in a single transaction
func (r *CommentRepository) CreateCommentsInBatch(comments []*models.Comment) error {
	return r.CreateCommentsInBatchWithContext(context.Background(), comments)
}

// CreateCommentsInBatchWithContext creates multiple comments in a single transaction with context
func (r *CommentRepository) CreateCommentsInBatchWithContext(ctx context.Context, comments []*models.Comment) error {
	if len(comments) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, comment := range comments {
		if err := r.CreateCommentWithTransaction(tx, comment); err != nil {
			return fmt.Errorf("failed to create comment %d in batch: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch transaction: %w", err)
	}

	return nil
}

// GetCommentsWithPagination retrieves comments with offset-based pagination
func (r *CommentRepository) GetCommentsWithPagination(offset, limit int) ([]*models.Comment, error) {
	return r.GetCommentsWithPaginationWithContext(context.Background(), offset, limit)
}

// GetCommentsWithPaginationWithContext retrieves comments with offset-based pagination and context
func (r *CommentRepository) GetCommentsWithPaginationWithContext(ctx context.Context, offset, limit int) ([]*models.Comment, error) {
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative: %d", offset)
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive: %d", limit)
	}
	if limit > 1000 {
		limit = 1000
	}

	query := selectCommentFields + ` FROM comments ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments with pagination: %w", err)
	}
	defer rows.Close()

	return r.scanComments(rows)
}
