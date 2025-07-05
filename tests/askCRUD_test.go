package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestCreateAsk(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()
	randomNum := rand.Intn(1000)

	ask := &models.Ask{
		ID:            randomNum,
		Type:          "ask",
		Title:         "Ask HN: How do you implement clean architecture in Go with modern practices?",
		Text:          "I'm looking for best practices and real-world examples of clean architecture in Go projects. Specifically interested in dependency injection, testing strategies, and project structure.",
		Score:         rand.Intn(100),
		Author:        "enhanced_curious_dev",
		Reply_ids:     []int{rand.Intn(100), rand.Intn(100), rand.Intn(100)},
		Replies_count: 3,
		Created_At:    time.Now().Unix(),
	}

	err := repo.Create(ctx, ask)
	if err != nil {
		t.Fatalf("Failed to create ask: %v", err)
	} else {
		t.Logf("Ask created successfully: %v", ask)
	}
}

func TestAskGetByID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()
	id := 852

	ask, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get ask by ID: %v", err)
	}

	if ask == nil {
		t.Fatalf("Expected ask to be found, but got nil")
	}

	t.Logf("Successfully retrieved ask: ID=%d, Title=%s, Replies=%d", ask.ID, ask.Title, ask.Replies_count)
}

func TestAskUpdate(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	ask := &models.Ask{
		ID:            852,
		Type:          "ask",
		Title:         "Updated: Best Practices for Clean Architecture in Go (2025 Edition)",
		Text:          "Updated text with more comprehensive details, examples, and modern Go patterns including generics, context usage, and microservice patterns.",
		Score:         45,
		Author:        "enhanced_curious_dev",
		Reply_ids:     []int{101, 102, 103, 104, 105},
		Replies_count: 5,
		Created_At:    time.Now().Unix(),
	}

	err := repo.Update(ctx, ask)
	if err != nil {
		t.Fatalf("Failed to update ask ID: %d", ask.ID)
	} else {
		t.Logf("Ask with ID %d updated successfully", ask.ID)
	}

	retrieved, err := repo.GetByID(ctx, ask.ID)
	if err != nil {
		t.Fatalf("Failed to get ask by ID %d", ask.ID)
	}

	if retrieved.Title != ask.Title {
		t.Fatalf("Failed to update title")
	}

	if retrieved.Text != ask.Text {
		t.Fatalf("Failed to update text")
	}

	if retrieved.Score != ask.Score {
		t.Fatalf("Failed to update score")
	}

	if len(retrieved.Reply_ids) != len(ask.Reply_ids) {
		t.Fatalf("Expected reply IDs length %d, got %d", len(ask.Reply_ids), len(retrieved.Reply_ids))
	}

	if retrieved.Replies_count != ask.Replies_count {
		t.Fatalf("Failed to update replies count")
	}
}

func TestAskUpdateScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	err := repo.UpdateScore(ctx, 852, 75)
	if err != nil {
		t.Errorf("Failed to update ask score: %v", err)
		return
	}

	// Verify score update
	retrieved, err := repo.GetByID(ctx, 852)
	if err != nil {
		t.Errorf("Failed to get ask after score update: %v", err)
		return
	}

	if retrieved.Score != 75 {
		t.Errorf("Expected score 75, got %d", retrieved.Score)
	}

	t.Logf("Successfully updated ask score to %d", retrieved.Score)
}

func TestAskUpdateRepliesCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	err := repo.UpdateRepliesCount(ctx, 852, 12)
	if err != nil {
		t.Errorf("Failed to update replies count: %v", err)
		return
	}

	// Verify replies count update
	retrieved, err := repo.GetByID(ctx, 852)
	if err != nil {
		t.Errorf("Failed to get ask after replies count update: %v", err)
		return
	}

	if retrieved.Replies_count != 12 {
		t.Errorf("Expected replies count 12, got %d", retrieved.Replies_count)
	}

	t.Logf("Successfully updated ask replies count to %d", retrieved.Replies_count)
}

func TestAskGetByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	asks, err := repo.GetByAuthor(ctx, "enhanced_curious_dev")
	if err != nil {
		t.Errorf("Failed to get asks by author: %v", err)
		return
	}

	t.Logf("Found %d asks by author 'enhanced_curious_dev'", len(asks))
	for _, ask := range asks {
		t.Logf("Ask: ID=%d, Title=%s, Score=%d", ask.ID, ask.Title, ask.Score)
	}
}

func TestAskGetByMinScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	asks, err := repo.GetByMinScore(ctx, 50)
	if err != nil {
		t.Errorf("Failed to get asks by min score: %v", err)
		return
	}

	for _, ask := range asks {
		if ask.Score < 50 {
			t.Errorf("Got ask with score %d, expected >= 50", ask.Score)
		}
	}

	t.Logf("Found %d asks with score >= 50", len(asks))
}

func TestAskGetRecent(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	asks, err := repo.GetRecent(ctx, 600)
	if err != nil {
		t.Errorf("Failed to get recent asks: %v", err)
		return
	}

	if len(asks) == 0 {
		t.Error("Expected at least one recent ask")
	}

	t.Logf("Retrieved %d recent asks", len(asks))
}

func TestAskGetByDateRange(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	start := time.Now().Add(-24 * time.Hour).Unix()
	end := time.Now().Add(24 * time.Hour).Unix()

	asks, err := repo.GetByDateRange(ctx, start, end)
	if err != nil {
		t.Errorf("Failed to get asks by date range: %v", err)
		return
	}

	t.Logf("Found %d asks in date range", len(asks))
}

func TestAskGetAll(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	asks, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get all asks: %v", err)
		return
	}

	if len(asks) == 0 {
		t.Error("Expected at least one ask")
	}

	t.Logf("Retrieved %d total asks", len(asks))
}

func TestAskExists(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	exists, err := repo.Exists(ctx, 852)
	if err != nil {
		t.Errorf("Failed to check ask existence: %v", err)
		return
	}

	if !exists {
		t.Error("Expected ask to exist")
	}

	t.Logf("Ask exists: %v", exists)
}

func TestAskGetCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Errorf("Failed to get ask count: %v", err)
		return
	}

	t.Logf("Total ask count: %d", count)
}

func TestAskCreateBatch(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	asks := []*models.Ask{
		{
			ID:            3201,
			Type:          "ask",
			Title:         "Ask HN: Go vs Rust for systems programming in 2025?",
			Text:          "Which language do you prefer for systems programming and why? Looking for real-world experiences.",
			Score:         30,
			Author:        "batch_ask_creator",
			Reply_ids:     []int{},
			Replies_count: 0,
			Created_At:    time.Now().Unix(),
		},
		{
			ID:            3202,
			Type:          "ask",
			Title:         "Ask HN: Best practices for Go error handling?",
			Text:          "How do you handle errors effectively in large Go codebases? Any recommended patterns?",
			Score:         25,
			Author:        "batch_ask_creator",
			Reply_ids:     []int{},
			Replies_count: 0,
			Created_At:    time.Now().Unix(),
		},
	}

	err := repo.CreateBatch(ctx, asks)
	if err != nil {
		t.Errorf("Failed to create batch asks: %v", err)
		return
	}

	// Verify batch creation
	for _, ask := range asks {
		exists, err := repo.Exists(ctx, ask.ID)
		if err != nil {
			t.Errorf("Failed to check existence of batch ask %d: %v", ask.ID, err)
			continue
		}
		if !exists {
			t.Errorf("Batch ask %d does not exist", ask.ID)
		}
	}

	t.Logf("Successfully created %d asks in batch", len(asks))

	// Cleanup
	for _, ask := range asks {
		_ = repo.Delete(ctx, ask.ID)
	}
}

func TestAskDeleteByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	// Create an ask to delete
	tempAsk := &models.Ask{
		ID:            3301,
		Type:          "ask",
		Title:         "Ask to Delete",
		Text:          "This ask will be deleted",
		Score:         5,
		Author:        "deleteaskuser",
		Reply_ids:     []int{},
		Replies_count: 0,
		Created_At:    time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempAsk)

	err := repo.DeleteByAuthor(ctx, "deleteaskuser")
	if err != nil {
		t.Errorf("Failed to delete by author: %v", err)
		return
	}

	exists, _ := repo.Exists(ctx, tempAsk.ID)
	if exists {
		t.Error("Ask should have been deleted")
	}

	t.Logf("Successfully deleted asks by author 'deleteaskuser'")
}

func TestAskDelete(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	// Create an ask to delete
	tempAsk := &models.Ask{
		ID:            3401,
		Type:          "ask",
		Title:         "Temporary ask for deletion test",
		Text:          "This ask is created for testing deletion functionality",
		Score:         10,
		Author:        "tempuser",
		Reply_ids:     []int{},
		Replies_count: 0,
		Created_At:    time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempAsk)

	err := repo.Delete(ctx, tempAsk.ID)
	if err != nil {
		t.Errorf("Failed to delete ask: %v", err)
		return
	}

	// Verify deletion
	exists, _ := repo.Exists(ctx, tempAsk.ID)
	if exists {
		t.Error("Ask should have been deleted")
	}

	t.Logf("Successfully deleted ask ID: %d", tempAsk.ID)
}
