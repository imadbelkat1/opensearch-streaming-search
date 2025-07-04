package tests

import (
	"context"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestJobRepository(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewJobRepository()

	// Create a test job
	job := &models.Job{
		ID:         4001,
		Type:       "job",
		Title:      "Senior Go Developer at TechCorp (Remote)",
		Text:       "We're looking for experienced Go developers to join our team. Competitive salary and great benefits.",
		URL:        "https://techcorp.com/careers/senior-go-dev",
		Score:      30,
		Author:     "techcorp_hr",
		Created_At: time.Now().Unix(),
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, job)
		if err != nil {
			t.Errorf("Failed to create job: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		retrieved, err := repo.GetByID(ctx, job.ID)
		if err != nil {
			t.Errorf("Failed to get job: %v", err)
			return
		}

		if retrieved.ID != job.ID {
			t.Errorf("Expected ID %d, got %d", job.ID, retrieved.ID)
		}
		if retrieved.Title != job.Title {
			t.Errorf("Expected title %s, got %s", job.Title, retrieved.Title)
		}
		if retrieved.URL != job.URL {
			t.Errorf("Expected URL %s, got %s", job.URL, retrieved.URL)
		}
	})

	t.Run("Update", func(t *testing.T) {
		job.Title = "Updated: Lead Go Developer at TechCorp"
		job.Text = "Updated job description with new requirements"
		job.Score = 40

		err := repo.Update(ctx, job)
		if err != nil {
			t.Errorf("Failed to update job: %v", err)
			return
		}

		// Verify update
		retrieved, err := repo.GetByID(ctx, job.ID)
		if err != nil {
			t.Errorf("Failed to get updated job: %v", err)
			return
		}

		if retrieved.Title != job.Title {
			t.Errorf("Expected title %s, got %s", job.Title, retrieved.Title)
		}
		if retrieved.Text != job.Text {
			t.Errorf("Expected text %s, got %s", job.Text, retrieved.Text)
		}
		if retrieved.Score != 40 {
			t.Errorf("Expected score 40, got %d", retrieved.Score)
		}
	})

	t.Run("UpdateScore", func(t *testing.T) {
		err := repo.UpdateScore(ctx, job.ID, 50)
		if err != nil {
			t.Errorf("Failed to update score: %v", err)
			return
		}

		// Verify score update
		retrieved, err := repo.GetByID(ctx, job.ID)
		if err != nil {
			t.Errorf("Failed to get job after score update: %v", err)
			return
		}

		if retrieved.Score != 50 {
			t.Errorf("Expected score 50, got %d", retrieved.Score)
		}
	})

	t.Run("CreateBatch", func(t *testing.T) {
		jobs := []*models.Job{
			{
				ID:         4002,
				Type:       "job",
				Title:      "Go Backend Engineer at StartupXYZ",
				Text:       "Join our fast-growing startup as a backend engineer.",
				URL:        "https://startupxyz.com/jobs/backend",
				Score:      20,
				Author:     "startupxyz",
				Created_At: time.Now().Unix(),
			},
			{
				ID:         4003,
				Type:       "job",
				Title:      "Junior Go Developer at WebCo",
				Text:       "Great opportunity for junior developers.",
				URL:        "https://webco.com/jobs/junior-go",
				Score:      15,
				Author:     "webco_hr",
				Created_At: time.Now().Unix(),
			},
		}

		err := repo.CreateBatch(ctx, jobs)
		if err != nil {
			t.Errorf("Failed to create job batch: %v", err)
			return
		}

		// Verify batch creation
		for _, j := range jobs {
			exists, err := repo.Exists(ctx, j.ID)
			if err != nil {
				t.Errorf("Failed to check existence of batch job %d: %v", j.ID, err)
				continue
			}
			if !exists {
				t.Errorf("Batch job %d does not exist", j.ID)
			}
		}

		// Cleanup batch jobs
		for _, j := range jobs {
			_ = repo.Delete(ctx, j.ID)
		}
	})

	t.Run("GetByMinScore", func(t *testing.T) {
		jobs, err := repo.GetByMinScore(ctx, 40)
		if err != nil {
			t.Errorf("Failed to get jobs by score: %v", err)
			return
		}

		for _, j := range jobs {
			if j.Score < 40 {
				t.Errorf("Got job with score %d, expected >= 40", j.Score)
			}
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		jobs, err := repo.GetByAuthor(ctx, "techcorp_hr")
		if err != nil {
			t.Errorf("Failed to get jobs by author: %v", err)
			return
		}

		if len(jobs) == 0 {
			t.Error("Expected at least one job by 'techcorp_hr'")
		}
	})

	t.Run("GetRecent", func(t *testing.T) {
		jobs, err := repo.GetRecent(ctx, 5)
		if err != nil {
			t.Errorf("Failed to get recent jobs: %v", err)
			return
		}

		if len(jobs) == 0 {
			t.Error("Expected at least one recent job")
		}
	})

	t.Run("GetByDateRange", func(t *testing.T) {
		start := time.Now().Add(-24 * time.Hour).Unix()
		end := time.Now().Add(24 * time.Hour).Unix()

		jobs, err := repo.GetByDateRange(ctx, start, end)
		if err != nil {
			t.Errorf("Failed to get jobs by date range: %v", err)
			return
		}

		if len(jobs) == 0 {
			t.Error("Expected at least one job in date range")
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		jobs, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("Failed to get all jobs: %v", err)
			return
		}

		if len(jobs) == 0 {
			t.Error("Expected at least one job")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, job.ID)
		if err != nil {
			t.Errorf("Failed to check existence: %v", err)
			return
		}

		if !exists {
			t.Error("Expected job to exist")
		}
	})

	t.Run("GetCount", func(t *testing.T) {
		count, err := repo.GetCount(ctx)
		if err != nil {
			t.Errorf("Failed to get count: %v", err)
			return
		}

		if count < 1 {
			t.Error("Expected at least one job")
		}
	})

	t.Run("DeleteByAuthor", func(t *testing.T) {
		// Create a job to delete
		tempJob := &models.Job{
			ID:         4004,
			Type:       "job",
			Title:      "Job to Delete",
			Text:       "This will be deleted",
			URL:        "https://example.com/delete-job",
			Score:      10,
			Author:     "deletecompany",
			Created_At: time.Now().Unix(),
		}

		_ = repo.Create(ctx, tempJob)

		err := repo.DeleteByAuthor(ctx, "deletecompany")
		if err != nil {
			t.Errorf("Failed to delete by author: %v", err)
			return
		}

		exists, _ := repo.Exists(ctx, tempJob.ID)
		if exists {
			t.Error("Job should have been deleted")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, job.ID)
		if err != nil {
			t.Errorf("Failed to delete job: %v", err)
			return
		}

		// Verify deletion
		exists, _ := repo.Exists(ctx, job.ID)
		if exists {
			t.Error("Job should have been deleted")
		}
	})
}
