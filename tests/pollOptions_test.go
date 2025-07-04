package tests

import (
	"context"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestPollOptionRepository(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	// Create poll options for testing
	options := []*models.PollOption{
		{
			ID:         1,
			Type:       "PollOption",
			PollID:     5001,
			Author:     "poll_creator",
			OptionText: "Gin",
			CreatedAt:  time.Now().Unix(),
			Votes:      0,
		},
		{
			ID:         2,
			Type:       "PollOption",
			PollID:     5001,
			Author:     "poll_creator",
			OptionText: "Echo",
			CreatedAt:  time.Now().Unix(),
			Votes:      0,
		},
		{
			ID:         3,
			Type:       "PollOption",
			PollID:     5001,
			Author:     "poll_creator",
			OptionText: "Fiber",
			CreatedAt:  time.Now().Unix(),
			Votes:      0,
		},
	}

	t.Run("Create", func(t *testing.T) {
		for _, opt := range options {
			err := repo.Create(ctx, opt)
			if err != nil {
				t.Errorf("Failed to create option %s: %v", opt.OptionText, err)
			}
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		option, err := repo.GetByID(ctx, options[0].ID)
		if err != nil {
			t.Errorf("Failed to get option: %v", err)
			return
		}

		if option.ID != options[0].ID {
			t.Errorf("Expected ID %d, got %d", options[0].ID, option.ID)
		}
		if option.OptionText != options[0].OptionText {
			t.Errorf("Expected text %s, got %s", options[0].OptionText, option.OptionText)
		}
		if option.PollID != options[0].PollID {
			t.Errorf("Expected poll ID %d, got %d", options[0].PollID, option.PollID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		options[0].OptionText = "Gin Framework"
		options[0].Votes = 5

		err := repo.Update(ctx, options[0])
		if err != nil {
			t.Errorf("Failed to update option: %v", err)
			return
		}

		// Verify update
		retrieved, err := repo.GetByID(ctx, options[0].ID)
		if err != nil {
			t.Errorf("Failed to get updated option: %v", err)
			return
		}

		if retrieved.OptionText != "Gin Framework" {
			t.Errorf("Expected text 'Gin Framework', got %s", retrieved.OptionText)
		}
		if retrieved.Votes != 5 {
			t.Errorf("Expected votes 5, got %d", retrieved.Votes)
		}
	})

	t.Run("GetByPollID", func(t *testing.T) {
		pollOptions, err := repo.GetByPollID(ctx, 5001)
		if err != nil {
			t.Errorf("Failed to get options by poll ID: %v", err)
			return
		}

		if len(pollOptions) != 3 {
			t.Errorf("Expected 3 options for poll 5001, got %d", len(pollOptions))
		}
	})

	t.Run("UpdateVotes", func(t *testing.T) {
		err := repo.UpdateVotes(ctx, options[0].ID, 10)
		if err != nil {
			t.Errorf("Failed to update votes: %v", err)
			return
		}

		// Verify vote update
		count, err := repo.GetVoteCount(ctx, options[0].ID)
		if err != nil {
			t.Errorf("Failed to get vote count: %v", err)
			return
		}

		if count != 10 {
			t.Errorf("Expected vote count 10, got %d", count)
		}
	})

	t.Run("IncrementVotes", func(t *testing.T) {
		err := repo.IncrementVotes(ctx, options[0].ID)
		if err != nil {
			t.Errorf("Failed to increment votes: %v", err)
			return
		}

		// Verify increment
		count, err := repo.GetVoteCount(ctx, options[0].ID)
		if err != nil {
			t.Errorf("Failed to get vote count after increment: %v", err)
			return
		}

		if count != 11 {
			t.Errorf("Expected vote count 11 after increment, got %d", count)
		}
	})

	t.Run("GetVoteCount", func(t *testing.T) {
		count, err := repo.GetVoteCount(ctx, options[1].ID)
		if err != nil {
			t.Errorf("Failed to get vote count: %v", err)
			return
		}

		if count != 0 {
			t.Errorf("Expected vote count 0 for untouched option, got %d", count)
		}
	})

	t.Run("CountByPollID", func(t *testing.T) {
		count, err := repo.CountByPollID(ctx, 5001)
		if err != nil {
			t.Errorf("Failed to count options: %v", err)
			return
		}

		if count != 3 {
			t.Errorf("Expected 3 options for poll 5001, got %d", count)
		}
	})

	t.Run("GetTopVoted", func(t *testing.T) {
		// Set different vote counts
		_ = repo.UpdateVotes(ctx, options[1].ID, 20)
		_ = repo.UpdateVotes(ctx, options[2].ID, 15)

		topOptions, err := repo.GetTopVoted(ctx, 5001, 2)
		if err != nil {
			t.Errorf("Failed to get top voted: %v", err)
			return
		}

		if len(topOptions) != 2 {
			t.Errorf("Expected 2 top options, got %d", len(topOptions))
			return
		}

		// First should be Echo with 20 votes
		if topOptions[0].OptionText != "Echo" || topOptions[0].Votes != 20 {
			t.Errorf("Expected top option to be Echo with 20 votes, got %s with %d votes",
				topOptions[0].OptionText, topOptions[0].Votes)
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		authorOptions, err := repo.GetByAuthor(ctx, "poll_creator")
		if err != nil {
			t.Errorf("Failed to get options by author: %v", err)
			return
		}

		if len(authorOptions) < 3 {
			t.Errorf("Expected at least 3 options by 'poll_creator', got %d", len(authorOptions))
		}
	})

	t.Run("GetRecent", func(t *testing.T) {
		recentOptions, err := repo.GetRecent(ctx, 5)
		if err != nil {
			t.Errorf("Failed to get recent options: %v", err)
			return
		}

		if len(recentOptions) == 0 {
			t.Error("Expected at least one recent option")
		}
	})

	t.Run("GetByDateRange", func(t *testing.T) {
		start := time.Now().Add(-24 * time.Hour).Unix()
		end := time.Now().Add(24 * time.Hour).Unix()

		rangeOptions, err := repo.GetByDateRange(ctx, start, end)
		if err != nil {
			t.Errorf("Failed to get options by date range: %v", err)
			return
		}

		if len(rangeOptions) == 0 {
			t.Error("Expected at least one option in date range")
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		allOptions, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("Failed to get all options: %v", err)
			return
		}

		if len(allOptions) == 0 {
			t.Error("Expected at least one option")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, options[0].ID)
		if err != nil {
			t.Errorf("Failed to check existence: %v", err)
			return
		}

		if !exists {
			t.Error("Expected option to exist")
		}
	})

	t.Run("GetCount", func(t *testing.T) {
		count, err := repo.GetCount(ctx)
		if err != nil {
			t.Errorf("Failed to get total count: %v", err)
			return
		}

		if count < 3 {
			t.Errorf("Expected at least 3 options total, got %d", count)
		}
	})

	t.Run("CreateBatch", func(t *testing.T) {
		batchOptions := []*models.PollOption{
			{
				ID:         10,
				Type:       "PollOption",
				PollID:     5002,
				Author:     "batch_creator",
				OptionText: "Option A",
				CreatedAt:  time.Now().Unix(),
				Votes:      0,
			},
			{
				ID:         11,
				Type:       "PollOption",
				PollID:     5002,
				Author:     "batch_creator",
				OptionText: "Option B",
				CreatedAt:  time.Now().Unix(),
				Votes:      0,
			},
		}

		err := repo.CreateBatch(ctx, batchOptions)
		if err != nil {
			t.Errorf("Failed to create batch: %v", err)
			return
		}

		// Verify batch creation
		count, err := repo.CountByPollID(ctx, 5002)
		if err != nil {
			t.Errorf("Failed to count batch options: %v", err)
			return
		}

		if count != 2 {
			t.Errorf("Expected 2 options for poll 5002, got %d", count)
		}

		// Cleanup
		_ = repo.DeleteByPollID(ctx, 5002)
	})

	t.Run("DeleteByAuthor", func(t *testing.T) {
		// Create an option to delete
		tempOption := &models.PollOption{
			ID:         20,
			Type:       "PollOption",
			PollID:     5003,
			Author:     "delete_author",
			OptionText: "Delete Me",
			CreatedAt:  time.Now().Unix(),
			Votes:      0,
		}

		_ = repo.Create(ctx, tempOption)

		err := repo.DeleteByAuthor(ctx, "delete_author")
		if err != nil {
			t.Errorf("Failed to delete by author: %v", err)
			return
		}

		exists, _ := repo.Exists(ctx, tempOption.ID)
		if exists {
			t.Error("Option should have been deleted")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		for _, opt := range options {
			err := repo.Delete(ctx, opt.ID)
			if err != nil {
				t.Errorf("Failed to delete option %d: %v", opt.ID, err)
			}
		}

		// Verify deletion
		count, _ := repo.CountByPollID(ctx, 5001)
		if count != 0 {
			t.Errorf("Expected 0 options for poll 5001 after deletion, got %d", count)
		}
	})

	t.Run("DeleteByPollID", func(t *testing.T) {
		// Create options for a new poll
		testOptions := []*models.PollOption{
			{
				ID:         30,
				Type:       "PollOption",
				PollID:     6000,
				Author:     "test_creator",
				OptionText: "Test 1",
				CreatedAt:  time.Now().Unix(),
				Votes:      0,
			},
			{
				ID:         31,
				Type:       "PollOption",
				PollID:     6000,
				Author:     "test_creator",
				OptionText: "Test 2",
				CreatedAt:  time.Now().Unix(),
				Votes:      0,
			},
			{
				ID:         32,
				Type:       "PollOption",
				PollID:     6000,
				Author:     "test_creator",
				OptionText: "Test 3",
				CreatedAt:  time.Now().Unix(),
				Votes:      0,
			},
		}
		for _, opt := range testOptions {
			err := repo.Create(ctx, opt)
			if err != nil {
				t.Errorf("Failed to create test option %s: %v", opt.OptionText, err)
			}
		}
		err := repo.DeleteByPollID(ctx, 6000)
		if err != nil {
			t.Errorf("Failed to delete options by poll ID: %v", err)
			return
		}
		// Verify deletion
		count, err := repo.CountByPollID(ctx, 6000)
		if err != nil {
			t.Errorf("Failed to count options after deletion: %v", err)
			return
		}
		if count != 0 {
			t.Errorf("Expected 0 options for poll 6000 after deletion, got %d", count)
		}
	})
}
