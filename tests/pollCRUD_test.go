package tests

import (
	"context"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestCreatePoll(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	poll := &models.Poll{
		ID:          5101,
		Type:        "poll",
		Title:       "What's your favorite Go web framework for 2025?",
		Score:       45,
		Author:      "enhanced_poll_creator",
		PollOptions: []int{1, 2, 3, 4, 5, 6},
		Reply_Ids:   []int{},
		Created_At:  time.Now().Unix(),
	}

	err := repo.Create(ctx, poll)
	if err != nil {
		t.Fatalf("Failed to create poll: %v", err)
	} else {
		t.Logf("Poll created successfully: %v", poll)
	}
}

func TestPollGetByID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()
	id := 5101

	poll, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get poll by ID: %v", err)
	}

	if poll == nil {
		t.Fatalf("Expected poll to be found, but got nil")
	}

	t.Logf("Successfully retrieved poll: ID=%d, Title=%s, Options=%v", poll.ID, poll.Title, poll.PollOptions)
}

func TestPollUpdate(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	poll := &models.Poll{
		ID:          5101,
		Type:        "poll",
		Title:       "Updated: Best Go Web Framework Survey 2025",
		Score:       75,
		Author:      "enhanced_poll_creator",
		PollOptions: []int{1, 2, 3, 4, 5, 6, 7},
		Reply_Ids:   []int{101, 102, 103, 104},
		Created_At:  time.Now().Unix(),
	}

	err := repo.Update(ctx, poll)
	if err != nil {
		t.Fatalf("Failed to update poll ID: %d", poll.ID)
	} else {
		t.Logf("Poll with ID %d updated successfully", poll.ID)
	}

	retrieved, err := repo.GetByID(ctx, poll.ID)
	if err != nil {
		t.Fatalf("Failed to get poll by ID %d", poll.ID)
	}

	if retrieved.Title != poll.Title {
		t.Fatalf("Failed to update title")
	}

	if retrieved.Score != poll.Score {
		t.Fatalf("Failed to update score")
	}

	if len(retrieved.PollOptions) != len(poll.PollOptions) {
		t.Fatalf("Expected poll options length %d, got %d", len(poll.PollOptions), len(retrieved.PollOptions))
	}

	if len(retrieved.Reply_Ids) != len(poll.Reply_Ids) {
		t.Fatalf("Expected reply IDs length %d, got %d", len(poll.Reply_Ids), len(retrieved.Reply_Ids))
	}
}

func TestPollUpdateScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	err := repo.UpdateScore(ctx, 5101, 100)
	if err != nil {
		t.Errorf("Failed to update poll score: %v", err)
		return
	}

	// Verify score update
	retrieved, err := repo.GetByID(ctx, 5101)
	if err != nil {
		t.Errorf("Failed to get poll after score update: %v", err)
		return
	}

	if retrieved.Score != 100 {
		t.Errorf("Expected score 100, got %d", retrieved.Score)
	}

	t.Logf("Successfully updated poll score to %d", retrieved.Score)
}

func TestPollGetByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	polls, err := repo.GetByAuthor(ctx, "enhanced_poll_creator")
	if err != nil {
		t.Errorf("Failed to get polls by author: %v", err)
		return
	}

	t.Logf("Found %d polls by author 'enhanced_poll_creator'", len(polls))
	for _, poll := range polls {
		t.Logf("Poll: ID=%d, Title=%s, Score=%d", poll.ID, poll.Title, poll.Score)
	}
}

func TestPollGetByMinScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	polls, err := repo.GetByMinScore(ctx, 50)
	if err != nil {
		t.Errorf("Failed to get polls by min score: %v", err)
		return
	}

	for _, poll := range polls {
		if poll.Score < 50 {
			t.Errorf("Got poll with score %d, expected >= 50", poll.Score)
		}
	}

	t.Logf("Found %d polls with score >= 50", len(polls))
}

func TestPollGetRecent(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	polls, err := repo.GetRecent(ctx, 5)
	if err != nil {
		t.Errorf("Failed to get recent polls: %v", err)
		return
	}

	if len(polls) == 0 {
		t.Error("Expected at least one recent poll")
	}

	t.Logf("Retrieved %d recent polls", len(polls))
}

func TestPollGetByDateRange(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	start := time.Now().Add(-24 * time.Hour).Unix()
	end := time.Now().Add(24 * time.Hour).Unix()

	polls, err := repo.GetByDateRange(ctx, start, end)
	if err != nil {
		t.Errorf("Failed to get polls by date range: %v", err)
		return
	}

	t.Logf("Found %d polls in date range", len(polls))
}

func TestPollGetAll(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	polls, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get all polls: %v", err)
		return
	}

	if len(polls) == 0 {
		t.Error("Expected at least one poll")
	}

	t.Logf("Retrieved %d total polls", len(polls))
}

func TestPollExists(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	exists, err := repo.Exists(ctx, 5101)
	if err != nil {
		t.Errorf("Failed to check poll existence: %v", err)
		return
	}

	if !exists {
		t.Error("Expected poll to exist")
	}

	t.Logf("Poll exists: %v", exists)
}

func TestPollGetCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Errorf("Failed to get poll count: %v", err)
		return
	}

	t.Logf("Total poll count: %d", count)
}

func TestPollCreateBatch(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	polls := []*models.Poll{
		{
			ID:          5201,
			Type:        "poll",
			Title:       "Favorite Go concurrency pattern?",
			Score:       25,
			Author:      "batch_poll_creator",
			PollOptions: []int{10, 11, 12},
			Reply_Ids:   []int{},
			Created_At:  time.Now().Unix(),
		},
		{
			ID:          5202,
			Type:        "poll",
			Title:       "Best Go deployment strategy?",
			Score:       30,
			Author:      "batch_poll_creator",
			PollOptions: []int{13, 14, 15, 16},
			Reply_Ids:   []int{},
			Created_At:  time.Now().Unix(),
		},
	}

	err := repo.CreateBatch(ctx, polls)
	if err != nil {
		t.Errorf("Failed to create batch polls: %v", err)
		return
	}

	// Verify batch creation
	for _, poll := range polls {
		exists, err := repo.Exists(ctx, poll.ID)
		if err != nil {
			t.Errorf("Failed to check existence of batch poll %d: %v", poll.ID, err)
			continue
		}
		if !exists {
			t.Errorf("Batch poll %d does not exist", poll.ID)
		}
	}

	t.Logf("Successfully created %d polls in batch", len(polls))

	// Cleanup
	for _, poll := range polls {
		_ = repo.Delete(ctx, poll.ID)
	}
}

func TestPollDeleteByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	// Create a poll to delete
	tempPoll := &models.Poll{
		ID:          5301,
		Type:        "poll",
		Title:       "Poll to Delete",
		Score:       5,
		Author:      "deletepolluser",
		PollOptions: []int{30, 31},
		Reply_Ids:   []int{},
		Created_At:  time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempPoll)

	err := repo.DeleteByAuthor(ctx, "deletepolluser")
	if err != nil {
		t.Errorf("Failed to delete by author: %v", err)
		return
	}

	exists, _ := repo.Exists(ctx, tempPoll.ID)
	if exists {
		t.Error("Poll should have been deleted")
	}

	t.Logf("Successfully deleted polls by author 'deletepolluser'")
}

func TestPollDelete(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollRepository()

	// Create a poll to delete
	tempPoll := &models.Poll{
		ID:          5401,
		Type:        "poll",
		Title:       "Temporary poll for deletion test",
		Score:       10,
		Author:      "tempuser",
		PollOptions: []int{40, 41},
		Reply_Ids:   []int{},
		Created_At:  time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempPoll)

	err := repo.Delete(ctx, tempPoll.ID)
	if err != nil {
		t.Errorf("Failed to delete poll: %v", err)
		return
	}

	// Verify deletion
	exists, _ := repo.Exists(ctx, tempPoll.ID)
	if exists {
		t.Error("Poll should have been deleted")
	}

	t.Logf("Successfully deleted poll ID: %d", tempPoll.ID)
}
