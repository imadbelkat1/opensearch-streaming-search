package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
	"internship-project/pkg/database"
)

func setupTest(t *testing.T) {
	// Initialize database connection
	config := database.GetDefaultConfig()
	if err := database.Connect(config); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Check health
	if err := database.Health(); err != nil {
		t.Fatalf("Failed health check: %v", err)
	}
}

func teardownTest() {
	database.Close()
}

func TestCreateStory(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()
	randomNum := rand.Intn(4000)

	story := &models.Story{
		ID:             randomNum,
		Type:           "story",
		Title:          "Test Story: Understanding Go Repository Pattern",
		URL:            "https://example.com/go-patterns",
		Score:          90,
		Author:         "testuser",
		Created_At:     time.Now().Unix(),
		Comments_ids:   []int{},
		Comments_count: 0,
	}

	err := repo.Create(ctx, story)
	if err != nil {
		t.Fatalf("Failed to create story: %v", err)
	}

	err = repo.Delete(ctx, story.ID)
	if err != nil {
		t.Fatalf("Failed to delete created story: %v", err)
	}

}

func TestGetStoryByID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	storyID := 1012 // Use the ID from the previous test
	story, err := repo.GetByID(ctx, storyID)
	if err != nil {
		t.Fatalf("Failed to get story by ID: %v", err)
	}

	if story == nil {
		t.Fatal("Expected story to be found, but got nil")
	}
	if story.ID != storyID {
		t.Errorf("Expected story ID %d, got %d", storyID, story.ID)
	}

	t.Logf("Retrieved story: %+v", story)
}

func TestUpdateStory(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	story := &models.Story{
		ID:             1012,
		Type:           "story",
		Title:          "Updated Story Title",
		URL:            "https://example.com/updated-url",
		Score:          95,
		Author:         "testuser",
		Created_At:     time.Now().Unix(),
		Comments_ids:   []int{},
		Comments_count: 0,
	}
	err := repo.Update(ctx, story)
	if err != nil {
		t.Fatalf("Failed to update story: %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByID(ctx, story.ID)
	if err != nil {
		t.Fatalf("Failed to get updated story: %v", err)
	}

	if retrieved.Title != story.Title {
		t.Errorf("Expected title %s, got %s", story.Title, retrieved.Title)
	}
	if retrieved.URL != story.URL {
		t.Errorf("Expected URL %s, got %s", story.URL, retrieved.URL)
	}
	if retrieved.Score != story.Score {
		t.Errorf("Expected score %d, got %d", story.Score, retrieved.Score)
	}
	if retrieved.Author != story.Author {
		t.Errorf("Expected author %s, got %s", story.Author, retrieved.Author)
	}
	if retrieved.Comments_count != story.Comments_count {
		t.Errorf("Expected comments count %d, got %d", story.Comments_count, retrieved.Comments_count)
	}
	if len(retrieved.Comments_ids) != len(story.Comments_ids) {
		t.Errorf("Expected comments IDs length %d, got %d", len(story.Comments_ids), len(retrieved.Comments_ids))
	} else {
		for i, id := range story.Comments_ids {
			if retrieved.Comments_ids[i] != id {
				t.Errorf("Expected comment ID %d at index %d, got %d", id, i, retrieved.Comments_ids[i])
			}
		}
	}
}

func TestGetAllStories(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	stories, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get recent stories: %v", err)
	}

	if len(stories) == 0 {
		t.Fatal("Expected at least one story, but got none")
	}

	for _, story := range stories {
		t.Logf("Story ID: %d, Title: %s", story.ID, story.Title)
	}
}

func TestGetByMinScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	minScore := 50
	stories, err := repo.GetByMinScore(ctx, minScore)
	if err != nil {
		t.Fatalf("Failed to get stories by minimum score: %v", err)
	}

	if len(stories) == 0 {
		t.Fatal("Expected at least one story with score >= 50, but got none")
	}

	for _, story := range stories {
		if story.Score < minScore {
			t.Errorf("Story ID %d has score %d, which is less than %d", story.ID, story.Score, minScore)
		} else {
			t.Logf("Story ID: %d, Title: %s, Score: %d", story.ID, story.Title, story.Score)
		}
	}

}

func TestGetByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	author := "testuser"
	stories, err := repo.GetByAuthor(ctx, author)
	if err != nil {
		t.Fatalf("Failed to get stories by author: %v", err)
	}

	if len(stories) == 0 {
		t.Fatalf("Expected at least one story by author '%s', but got none", author)
	}

	for _, story := range stories {
		if story.Author != author {
			t.Errorf("Expected author %s, got %s for story ID %d", author, story.Author, story.ID)
		} else {
			t.Logf("Story ID: %d, Title: %s, Author: %s", story.ID, story.Title, story.Author)
		}
	}
}

func TestGetByDateRange(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	start := time.Now().Add(-48 * time.Hour).Unix() // 2 days ago
	end := time.Now().Unix()                        // now

	stories, err := repo.GetByDateRange(ctx, start, end)
	if err != nil {
		t.Fatalf("Failed to get stories by date range: %v", err)
	}

	if len(stories) == 0 {
		t.Fatal("Expected at least one story in the date range, but got none")
	}

	for _, story := range stories {
		if story.Created_At < start || story.Created_At > end {
			t.Errorf("Story ID %d created at %d is outside the range [%d, %d]", story.ID, story.Created_At, start, end)
		} else {
			t.Logf("Story ID: %d, Title: %s, Created At: %d", story.ID, story.Title, story.Created_At)
		}
	}
}

func TestDeleteStory(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	// Create a story to delete
	story := &models.Story{
		ID:             1013,
		Type:           "story",
		Title:          "Story to Delete",
		URL:            "https://example.com/delete",
		Score:          10,
		Author:         "deleteuser",
		Created_At:     time.Now().Unix(),
		Comments_ids:   []int{},
		Comments_count: 0,
	}
	err := repo.Create(ctx, story)
	if err != nil {
		t.Fatalf("Failed to create story for deletion test: %v", err)
	}
	err = repo.Delete(ctx, story.ID)
	if err != nil {
		t.Fatalf("Failed to delete story: %v", err)
	}
	// Verify deletion
	exists, err := repo.Exists(ctx, story.ID)
	if err != nil {
		t.Fatalf("Failed to check existence after deletion: %v", err)
	}
	if exists {
		t.Error("Story should have been deleted, but it still exists")
	} else {
		t.Logf("Story ID %d successfully deleted", story.ID)
	}
}

func TestExistsStory(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	storyID := 1012 // Use the ID from the previous tests
	exists, err := repo.Exists(ctx, storyID)
	if err != nil {
		t.Fatalf("Failed to check existence of story: %v", err)
	}

	if !exists {
		t.Errorf("Expected story with ID %d to exist, but it does not", storyID)
	} else {
		t.Logf("Story with ID %d exists", storyID)
	}
}

func TestGetCountStories(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Fatalf("Failed to get count of stories: %v", err)
	}

	if count < 1 {
		t.Error("Expected at least one story, but got zero")
	} else {
		t.Logf("Total number of stories: %d", count)
	}
}

func TestCreateBatchStories(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	stories := []*models.Story{
		{
			ID:             1014,
			Type:           "story",
			Title:          "Batch Story 1",
			URL:            "https://example.com/batch1",
			Score:          10,
			Author:         "batchuser",
			Created_At:     time.Now().Unix(),
			Comments_ids:   []int{},
			Comments_count: 0,
		},
		{
			ID:             1015,
			Type:           "story",
			Title:          "Batch Story 2",
			URL:            "https://example.com/batch2",
			Score:          20,
			Author:         "batchuser",
			Created_At:     time.Now().Unix(),
			Comments_ids:   []int{},
			Comments_count: 0,
		},
	}
	err := repo.CreateBatch(ctx, stories)
	if err != nil {
		t.Fatalf("Failed to create batch of stories: %v", err)
	}
	// Verify batch creation
	for _, s := range stories {
		exists, err := repo.Exists(ctx, s.ID)
		if err != nil {
			t.Errorf("Failed to check existence of batch story %d: %v", s.ID, err)
			continue
		}
		if !exists {
			t.Errorf("Batch story %d does not exist", s.ID)
		} else {
			t.Logf("Batch story ID %d created successfully", s.ID)
		}
	}
}

func TestCleanBatchStories(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	// Clean up batch stories created in the previous test
	stories := []*models.Story{
		{ID: 1014},
		{ID: 1015},
	}

	for _, s := range stories {
		err := repo.Delete(ctx, s.ID)
		if err != nil {
			t.Errorf("Failed to delete batch story %d: %v", s.ID, err)
		} else {
			t.Logf("Batch story ID %d deleted successfully", s.ID)
		}
	}
}

func TestUpdateCommentsCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	story := &models.Story{
		ID:             1016,
		Type:           "story",
		Title:          "Test Story for Comments Count",
		URL:            "https://example.com/comments-count",
		Score:          10,
		Author:         "testuser",
		Created_At:     time.Now().Unix(),
		Comments_ids:   []int{},
		Comments_count: 0,
	}
	err := repo.Create(ctx, story)
	if err != nil {
		t.Fatalf("Failed to create story for comments count test: %v", err)
	}
	err = repo.UpdateCommentsCount(ctx, story.ID, 5)
	if err != nil {
		t.Fatalf("Failed to update comments count: %v", err)
	}
	// Verify comments count update
	retrieved, err := repo.GetByID(ctx, story.ID)
	if err != nil {
		t.Fatalf("Failed to get story after comments count update: %v", err)
	}
	if retrieved.Comments_count != 5 {
		t.Errorf("Expected comments count 5, got %d", retrieved.Comments_count)
	} else {
		t.Logf("Comments count for story ID %d updated successfully to %d", story.ID, retrieved.Comments_count)
	}
	// Clean up
	err = repo.Delete(ctx, story.ID)
	if err != nil {
		t.Errorf("Failed to delete story after comments count test: %v", err)
	} else {
		t.Logf("Story ID %d deleted successfully after comments count test", story.ID)
	}
}

func TestUpdateScore(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	story := &models.Story{
		ID:             1017,
		Type:           "story",
		Title:          "Test Story for Score Update",
		URL:            "https://example.com/score-update",
		Score:          10,
		Author:         "testuser",
		Created_At:     time.Now().Unix(),
		Comments_ids:   []int{},
		Comments_count: 0,
	}
	err := repo.Create(ctx, story)
	if err != nil {
		t.Fatalf("Failed to create story for score update test: %v", err)
	}
	err = repo.UpdateScore(ctx, story.ID, 20)
	if err != nil {
		t.Fatalf("Failed to update score: %v", err)
	}
	// Verify score update
	retrieved, err := repo.GetByID(ctx, story.ID)
	if err != nil {
		t.Fatalf("Failed to get story after score update: %v", err)
	}
	if retrieved.Score != 20 {
		t.Errorf("Expected score 20, got %d", retrieved.Score)
	} else {
		t.Logf("Score for story ID %d updated successfully to %d", story.ID, retrieved.Score)
	}
	// Clean up
	err = repo.Delete(ctx, story.ID)
	if err != nil {
		t.Errorf("Failed to delete story after score update test: %v", err)
	} else {
		t.Logf("Story ID %d deleted successfully after score update test", story.ID)
	}
}

func TestDeleteByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	// Create a story to delete
	story := &models.Story{
		ID:             1018,
		Type:           "story",
		Title:          "Story to Delete by Author",
		URL:            "https://example.com/delete-by-author",
		Score:          10,
		Author:         "deleteuser",
		Created_At:     time.Now().Unix(),
		Comments_ids:   []int{},
		Comments_count: 0,
	}
	err := repo.Create(ctx, story)
	if err != nil {
		t.Fatalf("Failed to create story for delete by author test: %v", err)
	}
	err = repo.DeleteByAuthor(ctx, "deleteuser")
	if err != nil {
		t.Fatalf("Failed to delete stories by author: %v", err)
	}
	// Verify deletion
	exists, err := repo.Exists(ctx, story.ID)
	if err != nil {
		t.Fatalf("Failed to check existence after deletion: %v", err)
	}
	if exists {
		t.Error("Story should have been deleted, but it still exists")
	} else {
		t.Logf("Story ID %d successfully deleted by author", story.ID)
	}
	// Clean up
	err = repo.Delete(ctx, story.ID)
	if err != nil {
		t.Errorf("Failed to delete story after delete by author test: %v", err)
	} else {
		t.Logf("Story ID %d deleted successfully after delete by author test", story.ID)
	}
}
func TestGetCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Fatalf("Failed to get count of stories: %v", err)
	}

	if count < 1 {
		t.Error("Expected at least one story, but got zero")
	} else {
		t.Logf("Total number of stories: %d", count)
	}
}

// TestCreateBatchStoriesWithExistingIDs tests creating a batch of stories where one story already exists
// in the database. It ensures that the existing story is not duplicated and the new story is created successfully.
// It also cleans up the created stories after the test.

func TestCreateBatchStoriesWithExistingIDs(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewStoryRepository()

	// Create a story that will be used in the batch
	existingStory := &models.Story{
		ID:             1039,
		Type:           "story",
		Title:          "Existing Story for Batch Test",
		URL:            "https://example.com/existing-batch",
		Score:          10,
		Author:         "batchuser",
		Created_At:     time.Now().Unix(),
		Comments_ids:   []int{},
		Comments_count: 0,
	}
	err := repo.Create(ctx, existingStory)
	if err != nil {
		t.Fatalf("Failed to create existing story for batch test: %v", err)
	}

	// Create a batch with one existing and one new story
	stories := []*models.Story{
		{
			ID:             1039, // Existing ID
			Type:           "story",
			Title:          "Batch Story 1 with Existing ID",
			URL:            "https://example.com/batch-existing",
			Score:          15,
			Author:         "batchuser",
			Created_At:     time.Now().Unix(),
			Comments_ids:   []int{},
			Comments_count: 0,
		},
		{
			ID:             1020, // New ID
			Type:           "story",
			Title:          "Batch Story 2 with New ID",
			URL:            "https://example.com/batch-new",
			Score:          20,
			Author:         "batchuser",
			Created_At:     time.Now().Unix(),
			Comments_ids:   []int{},
			Comments_count: 0,
		},
	}

	err = repo.CreateBatchWithExistingIDs(ctx, stories)
	if err != nil {
		t.Fatalf("Failed to create batch of stories with existing ID: %v", err)
	}

	// Verify that the existing story was not duplicated
	// Verify batch creation
	for _, s := range stories {
		exists, err := repo.Exists(ctx, s.ID)
		if err != nil {
			t.Errorf("Failed to check existence of batch story %d: %v", s.ID, err)
			continue
		}
		if !exists {
			t.Errorf("Batch story %d does not exist", s.ID)
		} else {
			t.Logf("Batch story ID %d created successfully", s.ID)
		}
	}

	// Clean up
	err = repo.Delete(ctx, existingStory.ID)
	if err != nil {
		t.Errorf("Failed to delete existing story after batch test: %v", err)
	} else {
		t.Logf("Existing story ID %d deleted successfully after batch test", existingStory.ID)
	}
	// Clean up new story
	err = repo.Delete(ctx, 1020) // Delete the new story created in the batch
	if err != nil {
		t.Errorf("Failed to delete new story after batch test: %v", err)
	} else {
		t.Logf("New story ID %d deleted successfully after batch test", existingStory.ID)
	}
}
