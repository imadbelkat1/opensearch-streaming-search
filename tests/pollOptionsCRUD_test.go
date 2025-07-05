package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestCreatePollOption(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()
	randomNum := rand.Intn(400)

	option := &models.PollOption{
		ID:         randomNum,
		Type:       "PollOption",
		PollID:     5001,
		Author:     "enhanced_poll_creator",
		OptionText: "Gin Framework",
		CreatedAt:  time.Now().Unix(),
		Votes:      rand.Intn(20),
	}

	err := repo.Create(ctx, option)
	if err != nil {
		t.Fatalf("Failed to create poll option: %v", err)
	} else {
		t.Logf("Poll option created successfully: %v", option)
	}
}

func TestPollOptionGetByID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()
	id := 153

	option, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get poll option by ID: %v", err)
	}

	if option == nil {
		t.Fatalf("Expected poll option to be found, but got nil")
	}

	t.Logf("Successfully retrieved poll option: ID=%d, Text=%s, Votes=%d", option.ID, option.OptionText, option.Votes)
}

func TestPollOptionUpdate(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	option := &models.PollOption{
		ID:         153,
		Type:       "PollOption",
		PollID:     5001,
		Author:     "updated_poll_creator",
		OptionText: "Updated: Gin Framework (Latest Version)",
		CreatedAt:  time.Now().Unix(),
		Votes:      15,
	}

	err := repo.Update(ctx, option)
	if err != nil {
		t.Fatalf("Failed to update poll option ID: %d", option.ID)
	} else {
		t.Logf("Poll option with ID %d updated successfully", option.ID)
	}

	retrieved, err := repo.GetByID(ctx, option.ID)
	if err != nil {
		t.Fatalf("Failed to get poll option by ID %d", option.ID)
	}

	if retrieved.OptionText != option.OptionText {
		t.Fatalf("Failed to update option text")
	}

	if retrieved.Votes != option.Votes {
		t.Fatalf("Failed to update votes")
	}

	if retrieved.PollID != option.PollID {
		t.Fatalf("Failed to update poll ID")
	}
}

func TestPollOptionUpdateVotes(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()
	randomNum := rand.Intn(20)

	err := repo.UpdateVotes(ctx, 153, randomNum)
	if err != nil {
		t.Errorf("Failed to update votes: %v", err)
		return
	}

	// Verify vote update
	count, err := repo.GetVoteCount(ctx, 340)
	if err != nil {
		t.Errorf("Failed to get vote count: %v", err)
		return
	}

	if count != 25 {
		t.Errorf("Expected vote count %d, got %d", randomNum, count)
	}

	t.Logf("Successfully updated votes to %d", count)
}

func TestPollOptionGetVoteCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	count, err := repo.GetVoteCount(ctx, 340)
	if err != nil {
		t.Errorf("Failed to get vote count: %v", err)
		return
	}

	t.Logf("Vote count for option ID 1: %d", count)
}

func TestPollOptionGetByPollID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	options, err := repo.GetByPollID(ctx, 5001)
	if err != nil {
		t.Errorf("Failed to get options by poll ID: %v", err)
		return
	}

	t.Logf("Found %d options for poll 5001", len(options))
	for _, option := range options {
		t.Logf("Option: ID=%d, Text=%s, Votes=%d", option.ID, option.OptionText, option.Votes)
	}
}

func TestPollOptionGetByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	options, err := repo.GetByAuthor(ctx, "enhanced_poll_creator")
	if err != nil {
		t.Errorf("Failed to get options by author: %v", err)
		return
	}

	t.Logf("Found %d options by author 'enhanced_poll_creator'", len(options))
	for _, option := range options {
		t.Logf("Option: ID=%d, Text=%s, Poll=%d", option.ID, option.OptionText, option.PollID)
	}
}

func TestPollOptionGetRecent(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	options, err := repo.GetRecent(ctx, 5)
	if err != nil {
		t.Errorf("Failed to get recent options: %v", err)
		return
	}

	if len(options) == 0 {
		t.Error("Expected at least one recent option")
	}

	t.Logf("Retrieved %d recent poll options", len(options))
}

func TestPollOptionGetByDateRange(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	start := time.Now().Add(-24 * time.Hour).Unix()
	end := time.Now().Add(24 * time.Hour).Unix()

	options, err := repo.GetByDateRange(ctx, start, end)
	if err != nil {
		t.Errorf("Failed to get options by date range: %v", err)
		return
	}

	t.Logf("Found %d options in date range", len(options))
}

func TestPollOptionGetAll(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	options, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get all options: %v", err)
		return
	}

	if len(options) == 0 {
		t.Error("Expected at least one option")
	}

	t.Logf("Retrieved %d total poll options", len(options))
}

func TestPollOptionExists(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	exists, err := repo.Exists(ctx, 340)
	if err != nil {
		t.Errorf("Failed to check option existence: %v", err)
		return
	}

	if !exists {
		t.Error("Expected option to exist")
	}

	t.Logf("Poll option exists: %v", exists)
}

func TestPollOptionCountByPollID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	count, err := repo.CountByPollID(ctx, 5001)
	if err != nil {
		t.Errorf("Failed to count options: %v", err)
		return
	}

	t.Logf("Poll 5001 has %d options", count)
}

func TestPollOptionGetTopVoted(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	// Get top voted poll options by PollID and by limit (from top downwards)
	topOptions, err := repo.GetTopVoted(ctx, 5001, 2)
	if err != nil {
		t.Errorf("Failed to get top voted: %v", err)
		return
	}

	if len(topOptions) == 0 {
		t.Error("Expected at least one top voted option")
		return
	}

	t.Logf("Top %d voted options for poll 5001:", len(topOptions))
	for i, option := range topOptions {
		t.Logf("%d. %s - %d votes", i+1, option.OptionText, option.Votes)
	}
}

func TestPollOptionGetCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Errorf("Failed to get total count: %v", err)
		return
	}

	t.Logf("Total poll option count: %d", count)
}

func TestPollOptionCreateBatch(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	options := []*models.PollOption{
		{
			ID:         1201,
			Type:       "PollOption",
			PollID:     5201,
			Author:     "batch_option_creator",
			OptionText: "Option A - Batch Created",
			CreatedAt:  time.Now().Unix(),
			Votes:      0,
		},
		{
			ID:         1202,
			Type:       "PollOption",
			PollID:     5201,
			Author:     "batch_option_creator",
			OptionText: "Option B - Batch Created",
			CreatedAt:  time.Now().Unix(),
			Votes:      0,
		},
	}

	err := repo.CreateBatch(ctx, options)
	if err != nil {
		t.Errorf("Failed to create batch options: %v", err)
		return
	}

	// Verify batch creation
	for _, option := range options {
		exists, err := repo.Exists(ctx, option.ID)
		if err != nil {
			t.Errorf("Failed to check existence of batch option %d: %v", option.ID, err)
			continue
		}
		if !exists {
			t.Errorf("Batch option %d does not exist", option.ID)
		}
	}

	t.Logf("Successfully created %d poll options in batch", len(options))

	// Cleanup
	for _, option := range options {
		_ = repo.Delete(ctx, option.ID)
	}
}

func TestPollOptionDeleteByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	// Create an option to delete
	tempOption := &models.PollOption{
		ID:         1301,
		Type:       "PollOption",
		PollID:     5301,
		Author:     "deleteoptionuser",
		OptionText: "Option to Delete",
		CreatedAt:  time.Now().Unix(),
		Votes:      0,
	}

	_ = repo.Create(ctx, tempOption)

	err := repo.DeleteByAuthor(ctx, "deleteoptionuser")
	if err != nil {
		t.Errorf("Failed to delete by author: %v", err)
		return
	}

	exists, _ := repo.Exists(ctx, tempOption.ID)
	if exists {
		t.Error("Option should have been deleted")
	}

	t.Logf("Successfully deleted poll options by author 'deleteoptionuser'")
}

func TestPollOptionDeleteByPollID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	// Create options for a specific poll
	testOptions := []*models.PollOption{
		{
			ID:         1401,
			Type:       "PollOption",
			PollID:     6001,
			Author:     "test_creator",
			OptionText: "Test Option 1",
			CreatedAt:  time.Now().Unix(),
			Votes:      0,
		},
		{
			ID:         1402,
			Type:       "PollOption",
			PollID:     6001,
			Author:     "test_creator",
			OptionText: "Test Option 2",
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

	err := repo.DeleteByPollID(ctx, 6001)
	if err != nil {
		t.Errorf("Failed to delete options by poll ID: %v", err)
		return
	}

	// Verify deletion
	count, err := repo.CountByPollID(ctx, 6001)
	if err != nil {
		t.Errorf("Failed to count options after deletion: %v", err)
		return
	}

	if count != 0 {
		t.Errorf("Expected 0 options for poll 6001 after deletion, got %d", count)
	}

	t.Logf("Successfully deleted all options for poll ID 6001")
}

func TestPollOptionDelete(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewPollOptionRepository()

	// Create an option to delete
	tempOption := &models.PollOption{
		ID:         1501,
		Type:       "PollOption",
		PollID:     5501,
		Author:     "tempuser",
		OptionText: "Temporary option for deletion test",
		CreatedAt:  time.Now().Unix(),
		Votes:      0,
	}

	_ = repo.Create(ctx, tempOption)

	err := repo.Delete(ctx, tempOption.ID)
	if err != nil {
		t.Errorf("Failed to delete poll option: %v", err)
		return
	}

	// Verify deletion
	exists, _ := repo.Exists(ctx, tempOption.ID)
	if exists {
		t.Error("Poll option should have been deleted")
	}

	t.Logf("Successfully deleted poll option ID: %d", tempOption.ID)
}
