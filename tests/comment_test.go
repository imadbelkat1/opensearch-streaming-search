package tests

import (
	"context"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestCommentRepository(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewCommentRepository()

	// Create a test comment
	comment := &models.Comment{
		ID:         2001,
		StoryID:    1001,
		Type:       "comment",
		Text:       "This is a great article about Go patterns!",
		Author:     "commenter1",
		Created_At: time.Now().Unix(),
		Parent:     0,
		Replies:    []int{},
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, comment)
		if err != nil {
			t.Errorf("Failed to create comment: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		retrieved, err := repo.GetByID(ctx, comment.ID)
		if err != nil {
			t.Errorf("Failed to get comment: %v", err)
			return
		}

		if retrieved.ID != comment.ID {
			t.Errorf("Expected ID %d, got %d", comment.ID, retrieved.ID)
		}
		if retrieved.Text != comment.Text {
			t.Errorf("Expected text %s, got %s", comment.Text, retrieved.Text)
		}
		if retrieved.StoryID != comment.StoryID {
			t.Errorf("Expected story ID %d, got %d", comment.StoryID, retrieved.StoryID)
		}
	})

	t.Run("CreateReply", func(t *testing.T) {
		reply := &models.Comment{
			ID:         2002,
			StoryID:    1001,
			Type:       "comment",
			Text:       "I agree! The repository pattern is very useful.",
			Author:     "commenter2",
			Created_At: time.Now().Unix(),
			Parent:     comment.ID,
			Replies:    []int{},
		}

		err := repo.Create(ctx, reply)
		if err != nil {
			t.Errorf("Failed to create reply: %v", err)
			return
		}

		// Update parent comment with reply
		comment.Replies = []int{reply.ID}
		err = repo.Update(ctx, comment)
		if err != nil {
			t.Errorf("Failed to update parent comment: %v", err)
			return
		}

		// Verify parent has reply
		parent, err := repo.GetByID(ctx, comment.ID)
		if err != nil {
			t.Errorf("Failed to get parent comment: %v", err)
			return
		}

		if len(parent.Replies) != 1 || parent.Replies[0] != reply.ID {
			t.Errorf("Expected parent to have reply %d, got %v", reply.ID, parent.Replies)
		}

		// Cleanup reply
		_ = repo.Delete(ctx, reply.ID)
	})

	t.Run("Update", func(t *testing.T) {
		comment.Text = "Updated comment text"
		err := repo.Update(ctx, comment)
		if err != nil {
			t.Errorf("Failed to update comment: %v", err)
			return
		}

		// Verify update
		retrieved, err := repo.GetByID(ctx, comment.ID)
		if err != nil {
			t.Errorf("Failed to get updated comment: %v", err)
			return
		}

		if retrieved.Text != comment.Text {
			t.Errorf("Expected text %s, got %s", comment.Text, retrieved.Text)
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		comments, err := repo.GetByAuthor(ctx, "commenter1")
		if err != nil {
			t.Errorf("Failed to get comments by author: %v", err)
			return
		}

		if len(comments) == 0 {
			t.Error("Expected at least one comment by author 'commenter1'")
		}
	})

	t.Run("GetRecent", func(t *testing.T) {
		comments, err := repo.GetRecent(ctx, 5)
		if err != nil {
			t.Errorf("Failed to get recent comments: %v", err)
			return
		}

		if len(comments) == 0 {
			t.Error("Expected at least one recent comment")
		}
	})

	t.Run("GetByDateRange", func(t *testing.T) {
		start := time.Now().Add(-24 * time.Hour).Unix()
		end := time.Now().Add(24 * time.Hour).Unix()

		comments, err := repo.GetByDateRange(ctx, start, end)
		if err != nil {
			t.Errorf("Failed to get comments by date range: %v", err)
			return
		}

		if len(comments) == 0 {
			t.Error("Expected at least one comment in date range")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, comment.ID)
		if err != nil {
			t.Errorf("Failed to check existence: %v", err)
			return
		}

		if !exists {
			t.Error("Expected comment to exist")
		}
	})

	t.Run("GetCount", func(t *testing.T) {
		count, err := repo.GetCount(ctx)
		if err != nil {
			t.Errorf("Failed to get count: %v", err)
			return
		}

		if count < 1 {
			t.Error("Expected at least one comment")
		}
	})

	t.Run("DeleteByAuthor", func(t *testing.T) {
		// Create a comment to delete
		tempComment := &models.Comment{
			ID:         2003,
			StoryID:    1001,
			Type:       "comment",
			Text:       "Comment to delete",
			Author:     "deleteuser",
			Created_At: time.Now().Unix(),
			Parent:     0,
			Replies:    []int{},
		}

		_ = repo.Create(ctx, tempComment)

		err := repo.DeleteByAuthor(ctx, "deleteuser")
		if err != nil {
			t.Errorf("Failed to delete by author: %v", err)
			return
		}

		exists, _ := repo.Exists(ctx, tempComment.ID)
		if exists {
			t.Error("Comment should have been deleted")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, comment.ID)
		if err != nil {
			t.Errorf("Failed to delete comment: %v", err)
			return
		}

		// Verify deletion
		exists, _ := repo.Exists(ctx, comment.ID)
		if exists {
			t.Error("Comment should have been deleted")
		}
	})
}
