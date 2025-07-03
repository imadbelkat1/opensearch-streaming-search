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
	insertJobQuery = `
		INSERT INTO jobs (id, type, title, url, score, author, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	selectJobFields = `
		SELECT id, type, title, url, score, author, created_at`

	selectJobByIDQuery = selectJobFields + ` FROM jobs WHERE id = $1`

	selectAllJobsQuery = selectJobFields + ` FROM jobs ORDER BY created_at DESC`

	selectRecentJobsQuery = selectJobFields + ` FROM jobs ORDER BY created_at DESC LIMIT $1`

	selectJobsByMinScoreQuery = selectJobFields + ` FROM jobs WHERE score >= $1 ORDER BY score DESC, created_at DESC`

	selectJobsByAuthorQuery = selectJobFields + ` FROM jobs WHERE author = $1 ORDER BY created_at DESC`

	selectJobsByDateRangeQuery = selectJobFields + ` FROM jobs WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`

	updateJobQuery = `
		UPDATE jobs SET type = $2, title = $3, url = $4, score = $5, author = $6, 
		created_at = $7 WHERE id = $1`

	updateJobScoreQuery = `UPDATE jobs SET score = $1 WHERE id = $2`

	deleteJobQuery = `DELETE FROM jobs WHERE id = $1`

	deleteJobsByAuthorQuery = `DELETE FROM jobs WHERE author = $1`

	jobExistsQuery = `SELECT EXISTS(SELECT 1 FROM jobs WHERE id = $1)`

	jobCountQuery = `SELECT COUNT(*) FROM jobs`
)

// JobRepository implements the JobRepository interface for PostgreSQL
type JobRepository struct {
	db *sql.DB
}

// NewJobRepository creates a new instance of JobRepository
func NewJobRepository() *JobRepository {
	return &JobRepository{
		db: database.GetDB(),
	}
}

// scanJob scans a single job from a row scanner
func (r *JobRepository) scanJob(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Job, error) {
	job := &models.Job{}
	err := scanner.Scan(
		&job.ID,
		&job.Type,
		&job.Title,
		&job.URL,
		&job.Score,
		&job.Author,
		&job.Created_At,
	)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// scanJobs scans multiple jobs from rows
func (r *JobRepository) scanJobs(rows *sql.Rows) ([]*models.Job, error) {
	var jobs []*models.Job

	for rows.Next() {
		job, err := r.scanJob(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return jobs, nil
}

// validateJob performs basic validation on job data
func (r *JobRepository) validateJob(job *models.Job) error {
	if job == nil {
		return fmt.Errorf("job cannot be nil")
	}
	if job.ID <= 0 {
		return fmt.Errorf("job ID must be positive")
	}
	if job.Author == "" {
		return fmt.Errorf("job author cannot be empty")
	}
	if job.Title == "" {
		return fmt.Errorf("job title cannot be empty")
	}
	return nil
}

// CREATE OPERATIONS

// CreateJob inserts a new job into the database
func (r *JobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	if err := r.validateJob(job); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := r.db.ExecContext(ctx, insertJobQuery,
		job.ID,
		job.Type,
		job.Title,
		job.URL,
		job.Score,
		job.Author,
		job.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}
	return nil
}

// CreateJobWithTransaction inserts a new job using the provided transaction
func (r *JobRepository) CreateJobWithTransaction(tx *sql.Tx, job *models.Job) error {
	if err := r.validateJob(job); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := tx.Exec(insertJobQuery,
		job.ID,
		job.Type,
		job.Title,
		job.URL,
		job.Score,
		job.Author,
		job.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to create job in transaction: %w", err)
	}
	return nil
}

// READ OPERATIONS

// GetJobByID retrieves a job by its ID
func (r *JobRepository) GetJobByID(ctx context.Context, id int) (*models.Job, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid job ID: %d", id)
	}

	job, err := r.scanJob(r.db.QueryRowContext(ctx, selectJobByIDQuery, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return job, nil
}

// GetAllJobs retrieves all jobs from the database
func (r *JobRepository) GetAllJobs(ctx context.Context) ([]*models.Job, error) {
	rows, err := r.db.QueryContext(ctx, selectAllJobsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get all jobs: %w", err)
	}
	defer rows.Close()

	return r.scanJobs(rows)
}

// GetRecentJobs retrieves the most recent jobs, limited by the specified count
func (r *JobRepository) GetRecentJobs(ctx context.Context, limit int) ([]*models.Job, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got: %d", limit)
	}
	if limit > 1000 {
		limit = 1000
	}

	rows, err := r.db.QueryContext(ctx, selectRecentJobsQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent jobs: %w", err)
	}
	defer rows.Close()

	return r.scanJobs(rows)
}

// GetJobsByScore retrieves jobs with score >= minScore
func (r *JobRepository) GetJobsByScore(ctx context.Context, minScore int) ([]*models.Job, error) {
	rows, err := r.db.QueryContext(ctx, selectJobsByMinScoreQuery, minScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by score: %w", err)
	}
	defer rows.Close()

	return r.scanJobs(rows)
}

// GetJobsByAuthor retrieves jobs by a specific author
func (r *JobRepository) GetJobsByAuthor(ctx context.Context, author string) ([]*models.Job, error) {
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}

	rows, err := r.db.QueryContext(ctx, selectJobsByAuthorQuery, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by author: %w", err)
	}
	defer rows.Close()

	return r.scanJobs(rows)
}

// GetJobsByDateRange retrieves jobs created within a specific date range
func (r *JobRepository) GetJobsByDateRange(ctx context.Context, start, end int64) ([]*models.Job, error) {
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("start and end timestamps must be non-negative")
	}
	if start > end {
		return nil, fmt.Errorf("start timestamp (%d) cannot be greater than end timestamp (%d)", start, end)
	}

	rows, err := r.db.QueryContext(ctx, selectJobsByDateRangeQuery, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by date range: %w", err)
	}
	defer rows.Close()

	return r.scanJobs(rows)
}

// UPDATE OPERATIONS

// UpdateJob updates an existing job in the database
func (r *JobRepository) UpdateJob(ctx context.Context, job *models.Job) error {
	if err := r.validateJob(job); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	result, err := r.db.ExecContext(ctx, updateJobQuery,
		job.ID,
		job.Type,
		job.Title,
		job.URL,
		job.Score,
		job.Author,
		job.Created_At,
	)

	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job with ID %d not found", job.ID)
	}

	return nil
}

// UpdateJobsScore updates the score of a specific job
func (r *JobRepository) UpdateJobsScore(ctx context.Context, id int, score int) error {
	if id <= 0 {
		return fmt.Errorf("invalid job ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, updateJobScoreQuery, score, id)
	if err != nil {
		return fmt.Errorf("failed to update score for job ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job with ID %d not found", id)
	}

	return nil
}

// DELETE OPERATIONS

// DeleteJob removes a job from the database by its ID
func (r *JobRepository) DeleteJob(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid job ID: %d", id)
	}

	result, err := r.db.ExecContext(ctx, deleteJobQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete job with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job with ID %d not found", id)
	}

	return nil
}

// DeleteJobsByAuthor removes all jobs made by a specific author
func (r *JobRepository) DeleteJobsByAuthor(ctx context.Context, author string) error {
	if author == "" {
		return fmt.Errorf("author cannot be empty")
	}

	result, err := r.db.ExecContext(ctx, deleteJobsByAuthorQuery, author)
	if err != nil {
		return fmt.Errorf("failed to delete jobs by author %s: %w", author, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no jobs found for author %s", author)
	}

	return nil
}

// UTILITY OPERATIONS

// JobExists checks if a job exists in the database by its ID
func (r *JobRepository) JobExists(ctx context.Context, id int) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid job ID: %d", id)
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, jobExistsQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if job exists: %w", err)
	}
	return exists, nil
}

// GetJobCount retrieves the total number of jobs in the database
func (r *JobRepository) GetJobCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, jobCountQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get job count: %w", err)
	}
	return count, nil
}

// BATCH OPERATIONS

// CreateJobsInBatch creates multiple jobs in a single transaction
func (r *JobRepository) CreateJobsInBatch(ctx context.Context, jobs []*models.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, job := range jobs {
		if err := r.CreateJobWithTransaction(tx, job); err != nil {
			return fmt.Errorf("failed to create job %d in batch: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch transaction: %w", err)
	}

	return nil
}
