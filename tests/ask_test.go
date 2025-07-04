package tests

import (
	"context"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestAskRepository(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewAskRepository()

	// Create a test ask
	ask := &models.Ask{
		ID:            3001,
		Type:          "ask",
		Title:         "Ask HN: How do you implement clean architecture in Go?",
		Text:          "I'm looking for best practices and real-world examples of clean architecture in Go projects.",
		Score:         15,
		Author:        "curious_dev",
		Reply_ids:     []int{},
		Replies_count: 0,
		Created_At:    time.Now().Unix(),
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, ask)
		if err != nil {
			t.Errorf("Failed to create ask: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		retrieved, err := repo.GetByID(ctx, ask.ID)
		if err != nil {
			t.Errorf("Failed to get ask: %v", err)
			return
		}

		if retrieved.ID != ask.ID {
			t.Errorf("Expected ID %d, got %d", ask.ID, retrieved.ID)
		}
		if retrieved.Title != ask.Title {
			t.Errorf("Expected title %s, got %s", ask.Title, retrieved.Title)
		}
		if retrieved.Text != ask.Text {
			t.Errorf("Expected text %s, got %s", ask.Text, retrieved.Text)
		}
	})

	t.Run("Update", func(t *testing.T) {
		ask.Title = "Updated: Best Practices for Clean Architecture in Go"
		ask.Text = "Updated text with more details"
		ask.Reply_ids = []int{101, 102}

		err := repo.Update(ctx, ask)
		if err != nil {
			t.Errorf("Failed to update ask: %v", err)
			return
		}

		// Verify update
		retrieved, err := repo.GetByID(ctx, ask.ID)
		if err != nil {
			t.Errorf("Failed to get updated ask: %v", err)
			return
		}

		if retrieved.Title != ask.Title {
			t.Errorf("Expected title %s, got %s", ask.Title, retrieved.Title)
		}
		if retrieved.Text != ask.Text {
			t.Errorf("Expected text %s, got %s", ask.Text, retrieved.Text)
		}
		if len(retrieved.Reply_ids) != 2 {
			t.Errorf("Expected 2 reply IDs, got %d", len(retrieved.Reply_ids))
		}
	})

	t.Run("UpdateScore", func(t *testing.T) {
		err := repo.UpdateScore(ctx, ask.ID, 25)
		if err != nil {
			t.Errorf("Failed to update score: %v", err)
			return
		}

		// Verify score update
		retrieved, err := repo.GetByID(ctx, ask.ID)
		if err != nil {
			t.Errorf("Failed to get ask after score update: %v", err)
			return
		}

		if retrieved.Score != 25 {
			t.Errorf("Expected score 25, got %d", retrieved.Score)
		}
	})

	t.Run("UpdateRepliesCount", func(t *testing.T) {
		err := repo.UpdateRepliesCount(ctx, ask.ID, 5)
		if err != nil {
			t.Errorf("Failed to update replies count: %v", err)
			return
		}

		// Verify replies count update
		retrieved, err := repo.GetByID(ctx, ask.ID)
		if err != nil {
			t.Errorf("Failed to get ask after replies count update: %v", err)
			return
		}

		if retrieved.Replies_count != 5 {
			t.Errorf("Expected replies count 5, got %d", retrieved.Replies_count)
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		asks, err := repo.GetByAuthor(ctx, "curious_dev")
		if err != nil {
			t.Errorf("Failed to get asks by author: %v", err)
			return
		}

		if len(asks) == 0 {
			t.Error("Expected at least one ask by author 'curious_dev'")
		}
	})

	t.Run("GetByMinScore", func(t *testing.T) {
		asks, err := repo.GetByMinScore(ctx, 20)
		if err != nil {
			t.Errorf("Failed to get asks by score: %v", err)
			return
		}

		for _, a := range asks {
			if a.Score < 20 {
				t.Errorf("Got ask with score %d, expected >= 20", a.Score)
			}
		}
	})

	t.Run("GetRecent", func(t *testing.T) {
		asks, err := repo.GetRecent(ctx, 10)
		if err != nil {
			t.Errorf("Failed to get recent asks: %v", err)
			return
		}

		if len(asks) == 0 {
			t.Error("Expected at least one recent ask")
		}
	})

	t.Run("GetByDateRange", func(t *testing.T) {
		start := time.Now().Add(-24 * time.Hour).Unix()
		end := time.Now().Add(24 * time.Hour).Unix()

		asks, err := repo.GetByDateRange(ctx, start, end)
		if err != nil {
			t.Errorf("Failed to get asks by date range: %v", err)
			return
		}

		if len(asks) == 0 {
			t.Error("Expected at least one ask in date range")
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		asks, err := repo.GetAll(ctx)
		if err != nil {
			t.Errorf("Failed to get all asks: %v", err)
			return
		}

		if len(asks) == 0 {
			t.Error("Expected at least one ask")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, ask.ID)
		if err != nil {
			t.Errorf("Failed to check existence: %v", err)
			return
		}

		if !exists {
			t.Error("Expected ask to exist")
		}
	})

	t.Run("GetCount", func(t *testing.T) {
		count, err := repo.GetCount(ctx)
		if err != nil {
			t.Errorf("Failed to get count: %v", err)
			return
		}

		if count < 1 {
			t.Error("Expected at least one ask")
		}
	})

	t.Run("CreateBatch", func(t *testing.T) {
		asks := []*models.Ask{
			{
				ID:            3002,
				Type:          "ask",
				Title:         "Ask HN: Best Go testing practices?",
				Text:          "What are your favorite testing patterns?",
				Score:         5,
				Author:        "tester1",
				Reply_ids:     []int{},
				Replies_count: 0,
				Created_At:    time.Now().Unix(),
			},
			{
				ID:            3003,
				Type:          "ask",
				Title:         "Ask HN: Go performance tips?",
				Text:          "How do you optimize Go applications?",
				Score:         8,
				Author:        "tester2",
				Reply_ids:     []int{},
				Replies_count: 0,
				Created_At:    time.Now().Unix(),
			},
		}

		err := repo.CreateBatch(ctx, asks)
		if err != nil {
			t.Errorf("Failed to create batch: %v", err)
			return
		}

		// Verify batch creation
		for _, a := range asks {
			exists, err := repo.Exists(ctx, a.ID)
			if err != nil {
				t.Errorf("Failed to check existence of batch ask %d: %v", a.ID, err)
				continue
			}
			if !exists {
				t.Errorf("Batch ask %d does not exist", a.ID)
			}
		}

		// Cleanup batch asks
		for _, a := range asks {
			_ = repo.Delete(ctx, a.ID)
		}
	})

	t.Run("DeleteByAuthor", func(t *testing.T) {
		// Create an ask to delete
		tempAsk := &models.Ask{
			ID:            3004,
			Type:          "ask",
			Title:         "Ask to Delete",
			Text:          "This will be deleted",
			Score:         1,
			Author:        "deleteuser",
			Reply_ids:     []int{},
			Replies_count: 0,
			Created_At:    time.Now().Unix(),
		}

		_ = repo.Create(ctx, tempAsk)

		err := repo.DeleteByAuthor(ctx, "deleteuser")
		if err != nil {
			t.Errorf("Failed to delete by author: %v", err)
			return
		}

		exists, _ := repo.Exists(ctx, tempAsk.ID)
		if exists {
			t.Error("Ask should have been deleted")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, ask.ID)
		if err != nil {
			t.Errorf("Failed to delete ask: %v", err)
			return
		}

		// Verify deletion
		exists, _ := repo.Exists(ctx, ask.ID)
		if exists {
			t.Error("Ask should have been deleted")
		}
	})
}
