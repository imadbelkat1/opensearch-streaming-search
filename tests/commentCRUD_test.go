package tests

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestCreateComment(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()
	randomNum := rand.Intn(4000) + 2000

	comment := &models.Comment{
		ID:         randomNum,
		Type:       "comment",
		Text:       "Enhanced test comment with detailed content",
		Author:     "enhanced_testuser",
		Parent:     rand.Intn(1000),
		Replies:    []int{rand.Intn(4000), rand.Intn(4000), rand.Intn(4000)},
		Created_At: time.Now().Unix(),
	}

	err := repo.Create(ctx, comment)
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	} else {
		t.Logf("Comment created successfully: %v", comment)
	}
}

func TestCommentGetByID(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()
	id := 3018

	comment, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get comment by ID: %v", err)
	}

	if comment == nil {
		t.Fatalf("Expected comment to be found, but got nil")
	}

	t.Logf("Successfully retrieved comment: %v", comment)
}

func TestCommentUpdate(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	comment := &models.Comment{
		ID:         3018,
		Type:       "comment",
		Text:       "Enhanced updated test comment with more details",
		Author:     "Enhanced updated testuser",
		Parent:     452,
		Replies:    []int{123, 456, 789, 999},
		Created_At: time.Now().Unix(),
	}

	err := repo.Update(ctx, comment)
	if err != nil {
		t.Fatalf("Failed to update comment ID: %d", comment.ID)
	} else {
		t.Logf("Comment with ID %d updated successfully", comment.ID)
	}

	retrieved, err := repo.GetByID(ctx, comment.ID)
	if err != nil {
		t.Fatalf("Failed to get comment by ID %d", comment.ID)
	}

	if retrieved.Text != comment.Text {
		t.Fatalf("Failed to update Text")
	}

	if retrieved.Author != comment.Author {
		t.Fatalf("Failed to update author")
	}

	if retrieved.Parent != comment.Parent {
		t.Fatalf("Failed to update parent")
	}

	if len(retrieved.Replies) != len(comment.Replies) {
		t.Fatalf("Expected comment IDs length %d, got %d", len(comment.Replies), len(retrieved.Replies))
	}

	for i := 0; i < len(comment.Replies); i++ {
		if retrieved.Replies[i] != comment.Replies[i] {
			t.Fatalf("Expected %d at index %d, got %d", comment.Replies[i], i, retrieved.Replies[i])
		}
	}
}

func TestCommentGetByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	comments, err := repo.GetByAuthor(ctx, "enhanced_testuser")
	if err != nil {
		t.Errorf("Failed to get comments by author: %v", err)
		return
	}

	t.Logf("Found %d comments by author 'enhanced_testuser'", len(comments))
	for _, comment := range comments {
		t.Logf("Comment: ID=%d, Text=%s", comment.ID, comment.Text)
	}
}

func TestCommentGetRecent(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	comments, err := repo.GetRecent(ctx, 5)
	if err != nil {
		t.Errorf("Failed to get recent comments: %v", err)
		return
	}

	if len(comments) == 0 {
		t.Error("Expected at least one recent comment")
	}

	t.Logf("Retrieved %d recent comments", len(comments))
}

func TestCommentGetByDateRange(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	start := time.Now().Add(-24 * time.Hour).Unix()
	end := time.Now().Add(24 * time.Hour).Unix()

	comments, err := repo.GetByDateRange(ctx, start, end)
	if err != nil {
		t.Errorf("Failed to get comments by date range: %v", err)
		return
	}

	t.Logf("Found %d comments in date range", len(comments))
}

func TestCommentGetAll(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	comments, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get all comments: %v", err)
		return
	}

	if len(comments) == 0 {
		t.Error("Expected at least one comment")
	}

	t.Logf("Retrieved %d total comments", len(comments))
}

func TestCommentExists(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	exists, err := repo.Exists(ctx, 3018)
	if err != nil {
		t.Errorf("Failed to check comment existence: %v", err)
		return
	}

	if !exists {
		t.Error("Expected comment to exist")
	}

	t.Logf("Comment exists: %v", exists)
}

func TestCommentGetCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Errorf("Failed to get comment count: %v", err)
		return
	}

	t.Logf("Total comment count: %d", count)
}

func TestCommentDeleteByAuthor(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	// Create a comment to delete
	tempComment := &models.Comment{
		ID:         9999,
		Type:       "comment",
		Text:       "Comment to delete",
		Author:     "deletecommentuser",
		Parent:     1001,
		Replies:    []int{},
		Created_At: time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempComment)

	err := repo.DeleteByAuthor(ctx, "deletecommentuser")
	if err != nil {
		t.Errorf("Failed to delete by author: %v", err)
		return
	}

	exists, _ := repo.Exists(ctx, tempComment.ID)
	if exists {
		t.Error("Comment should have been deleted")
	}

	t.Logf("Successfully deleted comments by author 'deletecommentuser'")
}

func TestCommentDelete(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	// Create a comment to delete
	tempComment := &models.Comment{
		ID:         8888,
		Type:       "comment",
		Text:       "Temporary comment for deletion test",
		Author:     "tempuser",
		Parent:     1001,
		Replies:    []int{},
		Created_At: time.Now().Unix(),
	}

	_ = repo.Create(ctx, tempComment)

	err := repo.Delete(ctx, tempComment.ID)
	if err != nil {
		t.Errorf("Failed to delete comment: %v", err)
		return
	}

	// Verify deletion
	exists, _ := repo.Exists(ctx, tempComment.ID)
	if exists {
		t.Error("Comment should have been deleted")
	}

	t.Logf("Successfully deleted comment ID: %d", tempComment.ID)
}
