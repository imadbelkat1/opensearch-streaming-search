package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"internship-project/internal/models"
	"internship-project/internal/repository/postgres"
)

func TestCreateUser(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("testuser%d", randomNum)

	user := &models.User{
		Username:   username,
		Karma:      150,
		About:      "Test user for Go patterns",
		Created_At: time.Now().Unix(),
		Submitted:  []int{randomNum, randomNum + 1, randomNum + 2},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user was created
	exists, err := repo.UserExists(ctx, username)
	if err != nil {
		t.Fatalf("Failed to check user existence: %v", err)
	}
	if !exists {
		t.Errorf("User %s should exist after creation", username)
	}

	t.Logf("User %s created and deleted successfully", username)
}

func TestGetUserByUsername(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("testuser%d", randomNum)

	// Create a user first
	user := &models.User{
		Username:   username,
		Karma:      100,
		About:      "Test user for get by username",
		Created_At: time.Now().Unix(),
		Submitted:  []int{randomNum + 9, randomNum * 5, randomNum / 3},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for get test: %v", err)
	}

	// Test getting the user
	retrievedUser, err := repo.GetByIDString(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}

	if retrievedUser == nil {
		t.Fatal("Expected user to be found, but got nil")
	}
	if retrievedUser.Username != username {
		t.Errorf("Expected username %s, got %s", username, retrievedUser.Username)
	}
	if retrievedUser.Karma != user.Karma {
		t.Errorf("Expected karma %d, got %d", user.Karma, retrievedUser.Karma)
	}
	if retrievedUser.About != user.About {
		t.Errorf("Expected about %s, got %s", user.About, retrievedUser.About)
	}

	t.Logf("Retrieved user: ID=%d, Username=%s, Karma=%d",
		retrievedUser.ID, retrievedUser.Username, retrievedUser.Karma)
}

func TestUpdateUser(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("testuser%d", randomNum)

	// Create a user first
	user := &models.User{
		Username:   username,
		Karma:      100,
		About:      "Original about",
		Created_At: time.Now().Unix(),
		Submitted:  []int{randomNum + 6, randomNum - 594},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for update test: %v", err)
	}

	// Get the user to have the correct ID
	createdUser, err := repo.GetByIDString(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get created user: %v", err)
	}

	// Update the user
	createdUser.Karma = 200
	createdUser.About = "Updated about"
	createdUser.Submitted = []int{rand.Intn(4000), rand.Intn(4000), rand.Intn(4000), rand.Intn(4000)}

	err = repo.Update(ctx, createdUser)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify the update
	retrieved, err := repo.GetByIDString(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if retrieved.Karma != 200 {
		t.Errorf("Expected karma 200, got %d", retrieved.Karma)
	}
	if retrieved.About != "Updated about" {
		t.Errorf("Expected about 'Updated about', got %s", retrieved.About)
	}
	if len(retrieved.Submitted) != 4 {
		t.Errorf("Expected 4 submitted items, got %d", len(retrieved.Submitted))
	}

	t.Logf("User updated successfully: karma=%d, about=%s, submissions=%d",
		retrieved.Karma, retrieved.About, len(retrieved.Submitted))
}

func TestGetAllUsers(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()

	// Create test users
	testUsers := []*models.User{
		{
			Username:   "getalluser1",
			Karma:      100,
			About:      "Test user 1",
			Created_At: time.Now().Unix(),
			Submitted:  []int{1, 2},
		},
		{
			Username:   "getalluser2",
			Karma:      200,
			About:      "Test user 2",
			Created_At: time.Now().Unix(),
			Submitted:  []int{3, 4},
		},
	}

	for _, user := range testUsers {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create test user %s: %v", user.Username, err)
		}
	}

	// Test getting all users
	users, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(users) < 2 {
		t.Errorf("Expected at least 2 users, got %d", len(users))
	}

	t.Logf("Retrieved %d users", len(users))
}

func TestGetByMinKarma(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()

	// Create test users with different karma levels
	testUsers := []*models.User{
		{
			Username:   "lowkarmauser",
			Karma:      25,
			About:      "Low karma user",
			Created_At: time.Now().Unix(),
			Submitted:  []int{},
		},
		{
			Username:   "highkarmauser",
			Karma:      150,
			About:      "High karma user",
			Created_At: time.Now().Unix(),
			Submitted:  []int{},
		},
	}

	for _, user := range testUsers {
		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create test user %s: %v", user.Username, err)
		}
	}

	// Test getting users with minimum karma
	minKarma := 100
	users, err := repo.GetByMinKarma(ctx, minKarma)
	if err != nil {
		t.Fatalf("Failed to get users by minimum karma: %v", err)
	}

	for _, user := range users {
		if user.Karma < minKarma {
			t.Errorf("User %s has karma %d, which is less than %d", user.Username, user.Karma, minKarma)
		}
	}

	t.Logf("Retrieved %d users with karma >= %d", len(users), minKarma)
}

func TestUpdateKarma(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("karmauser%d", randomNum)

	user := &models.User{
		Username:   username,
		Karma:      100,
		About:      "Test user for karma update",
		Created_At: time.Now().Unix(),
		Submitted:  []int{rand.Intn(4000), rand.Intn(4000), rand.Intn(4000), rand.Intn(4000)},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for karma update test: %v", err)
	}

	newKarma := 150
	err = repo.UpdateKarma(ctx, username, newKarma)
	if err != nil {
		t.Fatalf("Failed to update karma: %v", err)
	}

	retrieved, err := repo.GetByIDString(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get user after karma update: %v", err)
	}

	if retrieved.Karma != newKarma {
		t.Errorf("Expected karma %d, got %d", newKarma, retrieved.Karma)
	}

	t.Logf("Karma for user %s updated successfully from 100 to %d", username, retrieved.Karma)
}

func TestUpdateAbout(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("aboutuser%d", randomNum)

	user := &models.User{
		Username:   username,
		Karma:      50,
		About:      "Original about text",
		Created_At: time.Now().Unix(),
		Submitted:  []int{},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for about update test: %v", err)
	}

	newAbout := "Updated about text"
	err = repo.UpdateAbout(ctx, username, newAbout)
	if err != nil {
		t.Fatalf("Failed to update about: %v", err)
	}

	retrieved, err := repo.GetByIDString(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get user after about update: %v", err)
	}

	if retrieved.About != newAbout {
		t.Errorf("Expected about %s, got %s", newAbout, retrieved.About)
	}

	t.Logf("About for user %s updated successfully", username)
}

func TestAddRemoveSubmission(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("subuser%d", randomNum)

	user := &models.User{
		Username:   username,
		Karma:      75,
		About:      "Test user for submissions",
		Created_At: time.Now().Unix(),
		Submitted:  []int{},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for submission test: %v", err)
	}

	itemID := 12345
	err = repo.AddSubmission(ctx, username, itemID)
	if err != nil {
		t.Fatalf("Failed to add submission: %v", err)
	}

	retrieved, err := repo.GetByIDString(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get user after adding submission: %v", err)
	}

	found := false
	for _, id := range retrieved.Submitted {
		if id == itemID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected item ID %d to be in submitted list", itemID)
	}

	err = repo.RemoveSubmission(ctx, username, itemID)
	if err != nil {
		t.Fatalf("Failed to remove submission: %v", err)
	}

	retrieved, err = repo.GetByIDString(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get user after removing submission: %v", err)
	}

	for _, id := range retrieved.Submitted {
		if id == itemID {
			t.Errorf("Item ID %d should have been removed from submitted list", itemID)
		}
	}

	t.Logf("Successfully added and removed submission %d", itemID)
}

func TestGetSubmissionCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("subcountuser%d", randomNum)

	user := &models.User{
		Username:   username,
		Karma:      60,
		About:      "Test user for submission count",
		Created_At: time.Now().Unix(),
		Submitted:  []int{1, 2, 3},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for submission count test: %v", err)
	}

	count, err := repo.GetSubmissionCount(ctx, username)
	if err != nil {
		t.Fatalf("Failed to get submission count: %v", err)
	}

	if count != len(user.Submitted) {
		t.Errorf("Expected submission count %d, got %d", len(user.Submitted), count)
	}

	t.Logf("Submission count for user %s: %d", username, count)
}

func TestUserExists(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("existuser%d", randomNum)

	// Create a user first
	user := &models.User{
		Username:   username,
		Karma:      100,
		About:      "Test user for exists check",
		Created_At: time.Now().Unix(),
		Submitted:  []int{},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for exists test: %v", err)
	}

	exists, err := repo.UserExists(ctx, username)
	if err != nil {
		t.Fatalf("Failed to check user existence: %v", err)
	}

	if !exists {
		t.Errorf("Expected user %s to exist, but it does not", username)
	}

	// Test non-existent user
	exists, err = repo.UserExists(ctx, "nonexistentuser")
	if err != nil {
		t.Fatalf("Failed to check non-existent user: %v", err)
	}

	if exists {
		t.Error("Expected non-existent user to not exist, but it does")
	}

	t.Logf("User existence checks working correctly")
}

func TestCreateBatch(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()

	users := []*models.User{
		{
			Username:   fmt.Sprintf("batchuser%d", rand.Intn(4000)),
			Karma:      100,
			About:      "Batch user 1",
			Created_At: time.Now().Unix(),
			Submitted:  []int{1, 2},
		},
		{
			Username:   fmt.Sprintf("batchuser%d", rand.Intn(4000)),
			Karma:      150,
			About:      "Batch user 2",
			Created_At: time.Now().Unix(),
			Submitted:  []int{3, 4, 5},
		},
	}

	err := repo.CreateBatch(ctx, users)
	if err != nil {
		t.Fatalf("Failed to create batch of users: %v", err)
	}

	for _, u := range users {
		exists, err := repo.UserExists(ctx, u.Username)
		if err != nil {
			t.Errorf("Failed to check existence of batch user %s: %v", u.Username, err)
			continue
		}
		if !exists {
			t.Errorf("Batch user %s does not exist", u.Username)
		}
	}

	t.Logf("Created batch of %d users successfully", len(users))
}

func TestGetUserCount(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()

	// Create a test user to ensure count > 0
	user := &models.User{
		Username:   fmt.Sprintf("countuser%d", rand.Intn(4000)),
		Karma:      100,
		About:      "Test user for count",
		Created_At: time.Now().Unix(),
		Submitted:  []int{},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for count test: %v", err)
	}

	count, err := repo.GetCount(ctx)
	if err != nil {
		t.Fatalf("Failed to get count of users: %v", err)
	}

	if count < 1 {
		t.Error("Expected at least one user, but got zero")
	}

	t.Logf("Total number of users: %d", count)
}

func TestDeleteUser(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	ctx := context.Background()
	repo := postgres.NewUserRepository()
	randomNum := rand.Intn(4000)
	username := fmt.Sprintf("deleteuser%d", randomNum)

	user := &models.User{
		Username:   username,
		Karma:      25,
		About:      "User to delete",
		Created_At: time.Now().Unix(),
		Submitted:  []int{},
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user for deletion test: %v", err)
	}

	err = repo.Delete(ctx, username)
	if err != nil {
		t.Fatalf("Failed to delete user for deletion test: %v", err)
	}

	exists, err := repo.UserExists(ctx, username)
	if err != nil {
		t.Fatalf("Failed to check existence after deletion: %v", err)
	}

	if exists {
		t.Error("User should have been deleted, but it still exists")
	}

	t.Logf("User %s successfully deleted", username)
}
