package tests

import (
	"context"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestPollRepository(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	// Create a test poll
	poll := &models.Poll{
		ID:          5001,
		Type:        "poll",
		Title:       "What's your favorite Go web framework?",
		Score:       35,
		Author:      "poll_creator",
		PollOptions: []int{1, 2, 3, 4, 5},
		Reply_Ids:   []int{},
		Created_At:  time.Now().Unix(),
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, poll)
		if err != nil {
			t.Errorf("Failed to create poll: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		retrieved, err := repo.GetByID(ctx, poll.ID)
		if err != nil {
			t.Errorf("Failed to get poll: %v", err)
			return
		}

		if retrieved.ID != poll.ID {
			t.Errorf("Expected ID %d, got %d", poll.ID, retrieved.ID)
		}
		if retrieved.Title != poll.Title {
			t.Errorf("Expected title %s, got %s", poll.Title, retrieved.Title)
		}
		if len(retrieved.PollOptions) != len(poll.PollOptions) {
			t.Errorf("Expected %d options, got %d", len(poll.PollOptions), len(retrieved.PollOptions))
		}
	})

	t.Run("Update", func(t *testing.T) {
		poll.Score = 50
		poll.Reply_Ids = []int{101, 102, 103}
		poll.Title = "Updated: Best Go Web Framework Survey"

		err := repo.Update(ctx, poll)
		if err != nil {
			t.Errorf("Failed to update poll: %v", err)
			return
		}

		// Verify update
		retrieved, err := repo.GetByID(ctx, poll.ID)
		if err != nil {
			t.Errorf("Failed to get updated poll: %v", err)
			return
		}

		if retrieved.Score != 50 {
			t.Errorf("Expected score 50, got %d", retrieved.Score)
		}
		if retrieved.Title != poll.Title {
			t.Errorf("Expected title %s, got %s", poll.Title, retrieved.Title)
		}
		if len(retrieved.Reply_Ids) != 3 {
			t.Errorf("Expected 3 reply IDs, got %d", len(retrieved.Reply_Ids))
		}
	})

	t.Run("UpdateScore", func(t *testing.T) {
		err := repo.UpdateScore(ctx, poll.ID, 60)
		if err != nil {
			t.Errorf("Failed to update score: %v", err)
			return
		}

		// Verify score update
		retrieved, err := repo.GetByID(ctx, poll.ID)
		if err != nil {
			t.Errorf("Failed to get poll after score update: %v", err)
			return
		}

		if retrieved.Score != 60 {
			t.Errorf("Expected score 60, got %d", retrieved.Score)
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		polls, err := repo.GetByAuthor(ctx, "poll_creator")
		if err != nil {
			t.Errorf("Failed to get polls by author: %v", err)
			return
		}

		if len(polls) == 0 {
			t.Error("Expected at least one poll by 'poll_creator'")
		}
	})

	t.Run("GetByMinScore", func(t *testing.T) {
		polls, err := repo.GetByMinScore(ctx, 50)
		if err != nil {
			t.Errorf("Failed to get polls by score: %v", err)
			return
		}

		for _, p := range polls {
			if p.Score < 50 {
				t.Errorf("Got poll with score %d, expected >= 50", p.Score)
			}
		}
	})

	t.Run("GetRecent", func(t *testing.T) {
		polls, err := repo.GetRecent(ctx, 5)
		if err != nil {
			t.Errorf("Failed to get recent polls: %v", err)
			return
		}

		if len(polls) == 0 {
			t.Error("Expected at least one recent poll")
		}
	})

	t.Run("GetByDateRange", func(t *testing.T) {
		start := time.Now().Add(-24 * time.Hour).Unix()
		end := time.Now().Add(24 * time.Hour).Unix()

		polls, err := repo.GetByDateRange(ctx, start, end)
		if err != nil {
			t.Errorf("Failed to get polls by date range: %v", err)
			return
		}

		if len(polls) == 0 {
			t.Error("Expected at least one poll in date range")
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		polls, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("Failed to get all polls: %v", err)
			return
		}

		if len(polls) == 0 {
			t.Error("Expected at least one poll")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, poll.ID)
		if err != nil {
			t.Errorf("Failed to check existence: %v", err)
			return
		}

		if !exists {
			t.Error("Expected poll to exist")
		}
	})

	t.Run("GetCount", func(t *testing.T) {
		count, err := repo.GetCount(ctx)
		if err != nil {
			t.Errorf("Failed to get count: %v", err)
			return
		}

		if count < 1 {
			t.Error("Expected at least one poll")
		}
	})

	t.Run("CreateBatch", func(t *testing.T) {
		polls := []*models.Poll{
			{
				ID:          5002,
				Type:        "poll",
				Title:       "Favorite Go testing framework?",
				Score:       10,
				Author:      "test_creator",
				PollOptions: []int{10, 11, 12},
				Reply_Ids:   []int{},
				Created_At:  time.Now().Unix(),
			},
			{
				ID:          5003,
				Type:        "poll",
				Title:       "Best Go IDE?",
				Score:       15,
				Author:      "test_creator",
				PollOptions: []int{20, 21, 22, 23},
				Reply_Ids:   []int{},
				Created_At:  time.Now().Unix(),
			},
		}

		err := repo.CreateBatch(ctx, polls)
		if err != nil {
			t.Errorf("Failed to create batch: %v", err)
			return
		}

		// Verify batch creation
		for _, p := range polls {
			exists, err := repo.Exists(ctx, p.ID)
			if err != nil {
				t.Errorf("Failed to check existence of batch poll %d: %v", p.ID, err)
				continue
			}
			if !exists {
				t.Errorf("Batch poll %d does not exist", p.ID)
			}
		}

		// Cleanup batch polls
		for _, p := range polls {
			_ = repo.Delete(ctx, p.ID)
		}
	})

	t.Run("DeleteByAuthor", func(t *testing.T) {
		// Create a poll to delete
		tempPoll := &models.Poll{
			ID:          5004,
			Type:        "poll",
			Title:       "Poll to Delete",
			Score:       5,
			Author:      "deleteuser",
			PollOptions: []int{30, 31},
			Reply_Ids:   []int{},
			Created_At:  time.Now().Unix(),
		}

		_ = repo.Create(ctx, tempPoll)

		err := repo.DeleteByAuthor(ctx, "deleteuser")
		if err != nil {
			t.Errorf("Failed to delete by author: %v", err)
			return
		}

		exists, _ := repo.Exists(ctx, tempPoll.ID)
		if exists {
			t.Error("Poll should have been deleted")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, poll.ID)
		if err != nil {
			t.Errorf("Failed to delete poll: %v", err)
			return
		}

		// Verify deletion
		exists, _ := repo.Exists(ctx, poll.ID)
		if exists {
			t.Error("Poll should have been deleted")
		}
	})
}
