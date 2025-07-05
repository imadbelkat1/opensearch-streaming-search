package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestCreateJob(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()
	randomNum := rand.Intn(2000)

	job := &models.Job{
		ID:         randomNum,
		Type:       "job",
		Title:      "Senior Go Developer at TechCorp (Remote/Hybrid)",
		Text:       "We're looking for experienced Go developers to join our growing team. Competitive salary, great benefits, and flexible work arrangements.",
		URL:        "https://techcorp.com/careers/senior-go-dev-2025",
		Score:      rand.Intn(50),
		Author:     "enhanced_techcorp_hr",
		Created_At: time.Now().Unix(),
	}

	err := repo.Create(ctx, job)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	} else {
		t.Logf("Job created successfully: %v", job)
	}
}

func TestJobGetByID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()
	id := 35

	job, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get job by ID: %v", err)
	}

	if job == nil {
		t.Fatalf("Expected job to be found, but got nil")
	}

	t.Logf("Successfully retrieved job: ID=%d, Title=%s, URL=%s", job.ID, job.Title, job.URL)
}

func TestJobUpdate(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	job := &models.Job{
		ID:         35,
		Type:       "job",
		Title:      "Updated: Lead Go Developer at TechCorp (Remote/On-site)",
		Text:       "Updated job description with new requirements and enhanced benefits package. Looking for experienced Go developers with leadership skills.",
		URL:        "https://techcorp.com/careers/lead-go-dev-updated",
		Score:      75,
		Author:     "enhanced_techcorp_hr",
		Created_At: time.Now().Unix(),
	}

	err := repo.Update(ctx, job)
	if err != nil {
		t.Fatalf("Failed to update job ID: %d", job.ID)
	} else {
		t.Logf("Job with ID %d updated successfully", job.ID)
	}

	retrieved, err := repo.GetByID(ctx, job.ID)
	if err != nil {
		t.Fatalf("Failed to get job by ID %d", job.ID)
	}

	if retrieved.Title != job.Title {
		t.Fatalf("Failed to update title")
	}

	if retrieved.Text != job.Text {
		t.Fatalf("Failed to update text")
	}

	if retrieved.Score != job.Score {
		t.Fatalf("Failed to update score")
	}

	if retrieved.URL != job.URL {
		t.Fatalf("Failed to update URL")
	}
}

func TestJobUpdateScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	err := repo.UpdateScore(ctx, 35, 90)
	if err != nil {
		t.Errorf("Failed to update job score: %v", err)
		return
	}

	// Verify score update
	retrieved, err := repo.GetByID(ctx, 35)
	if err != nil {
		t.Errorf("Failed to get job after score update: %v", err)
		return
	}

	if retrieved.Score != 90 {
		t.Errorf("Expected score 90, got %d", retrieved.Score)
	}

	t.Logf("Successfully updated job score to %d", retrieved.Score)
}

func TestJobGetByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	jobs, err := repo.GetByAuthor(ctx, "enhanced_techcorp_hr")
	if err != nil {
		t.Errorf("Failed to get jobs by author: %v", err)
		return
	}

	t.Logf("Found %d jobs by author 'enhanced_techcorp_hr'", len(jobs))
	for _, job := range jobs {
		t.Logf("Job: ID=%d, Title=%s, Score=%d", job.ID, job.Title, job.Score)
	}
}

func TestJobGetByMinScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	jobs, err := repo.GetByMinScore(ctx, 60)
	if err != nil {
		t.Errorf("Failed to get jobs by min score: %v", err)
		return
	}

	for _, job := range jobs {
		if job.Score < 60 {
			t.Errorf("Got job with score %d, expected >= 60", job.Score)
		}
	}

	t.Logf("Found %d jobs with score >= 60", len(jobs))
}

func TestJobGetRecent(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	jobs, err := repo.GetRecent(ctx, 5)
	if err != nil {
		t.Errorf("Failed to get recent jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		t.Error("Expected at least one recent job")
	}

	t.Logf("Retrieved %d recent jobs", len(jobs))
}

func TestJobGetByDateRange(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	start := time.Now().Add(-24 * time.Hour).Unix()
	end := time.Now().Add(24 * time.Hour).Unix()

	jobs, err := repo.GetByDateRange(ctx, start, end)
	if err != nil {
		t.Errorf("Failed to get jobs by date range: %v", err)
		return
	}

	t.Logf("Found %d jobs in date range", len(jobs))
}

func TestJobGetAll(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	jobs, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get all jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		t.Error("Expected at least one job")
	}

	t.Logf("Retrieved %d total jobs", len(jobs))
}

func TestJobExists(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	exists, err := repo.Exists(ctx, 35)
	if err != nil {
		t.Errorf("Failed to check job existence: %v", err)
		return
	}

	if !exists {
		t.Error("Expected job to exist")
	}

	t.Logf("Job exists: %v", exists)
}

func TestJobGetCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Errorf("Failed to get job count: %v", err)
		return
	}

	t.Logf("Total job count: %d", count)
}

func TestJobCreateBatch(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	jobs := []*models.Job{
		{
			ID:         4201,
			Type:       "job",
			Title:      "DevOps Engineer - Go/Kubernetes",
			Text:       "Looking for DevOps engineer with Go experience and Kubernetes expertise.",
			URL:        "https://devops-company.com/jobs/devops-go-k8s",
			Score:      40,
			Author:     "batch_job_creator",
			Created_At: time.Now().Unix(),
		},
		{
			ID:         4202,
			Type:       "job",
			Title:      "Full Stack Developer (Go + React)",
			Text:       "Full stack position with Go backend and React frontend development.",
			URL:        "https://fullstack-inc.com/jobs/fullstack-go-react",
			Score:      35,
			Author:     "batch_job_creator",
			Created_At: time.Now().Unix(),
		},
	}

	err := repo.CreateBatch(ctx, jobs)
	if err != nil {
		t.Errorf("Failed to create batch jobs: %v", err)
		return
	}

	// Verify batch creation
	for _, job := range jobs {
		exists, err := repo.Exists(ctx, job.ID)
		if err != nil {
			t.Errorf("Failed to check existence of batch job %d: %v", job.ID, err)
			continue
		}
		if !exists {
			t.Errorf("Batch job %d does not exist", job.ID)
		}
	}

	t.Logf("Successfully created %d jobs in batch", len(jobs))

	// Cleanup
	for _, job := range jobs {
		_ = repo.Delete(ctx, job.ID)
	}
}

func TestJobDeleteByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	// Create a job to delete
	tempJob := &models.Job{
		ID:         4301,
		Type:       "job",
		Title:      "Job to Delete",
		Text:       "This job will be deleted",
		URL:        "https://example.com/delete-job",
		Score:      15,
		Author:     "deletejobuser",
		Created_At: time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempJob)

	err := repo.DeleteByAuthor(ctx, "deletejobuser")
	if err != nil {
		t.Errorf("Failed to delete by author: %v", err)
		return
	}

	exists, _ := repo.Exists(ctx, tempJob.ID)
	if exists {
		t.Error("Job should have been deleted")
	}

	t.Logf("Successfully deleted jobs by author 'deletejobuser'")
}

func TestJobDelete(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	// Create a job to delete
	tempJob := &models.Job{
		ID:         4401,
		Type:       "job",
		Title:      "Temporary job for deletion test",
		Text:       "This job is created for testing deletion functionality",
		URL:        "https://example.com/temp-job",
		Score:      20,
		Author:     "tempuser",
		Created_At: time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempJob)

	err := repo.Delete(ctx, tempJob.ID)
	if err != nil {
		t.Errorf("Failed to delete job: %v", err)
		return
	}

	// Verify deletion
	exists, _ := repo.Exists(ctx, tempJob.ID)
	if exists {
		t.Error("Job should have been deleted")
	}

	t.Logf("Successfully deleted job ID: %d", tempJob.ID)
}
