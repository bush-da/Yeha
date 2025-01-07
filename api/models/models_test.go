package models

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Global variable for database connection
var db *gorm.DB

// Initialize the test database connection
func setup() {
	var err error

	// Fetch database credentials from environment variables or defaults
	user := os.Getenv("MYSQL_USER")
	if user == "" {
		log.Println("Warning: MYSQL_USER not set, using default.")
		user = "test_user"
	}
	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		log.Println("Warning: MYSQL_PASSWORD not set, using default.")
		password = "test_password"
	}
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		log.Println("Warning: MYSQL_HOST not set, using default.")
		host = "localhost"
	}
	dbName := os.Getenv("MYSQL_DB")
	if dbName == "" {
		log.Println("Warning: MYSQL_DB not set, using default.")
		dbName = "test_db"
	}

	// Create DSN and connect
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Content{}, &Like{}, &Tag{}, &Report{}, &Follower{}, &PostTag{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Clear existing data
	resetTables()
}

// Clear all tables in reverse dependency order
func resetTables() {
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	db.Exec("TRUNCATE TABLE reports")
	db.Exec("TRUNCATE TABLE post_tags")
	db.Exec("TRUNCATE TABLE tags")
	db.Exec("TRUNCATE TABLE comments")
	db.Exec("TRUNCATE TABLE likes")
	db.Exec("TRUNCATE TABLE posts")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
}

// Cleanup after tests
func teardown() {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

// Test User model
func TestUserModel(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		user     User
		expected string
		hasError bool
	}{
		{"Valid User", User{Username: "john_doe", Email: "john@example.com", Password: "password123"}, "john_doe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.Create(&tt.user)
			if tt.hasError {
				// Check for error when expected
				assert.NotNil(t, result.Error, "Expected an error for invalid user")
			} else {
				// Check for no error when valid
				assert.Nil(t, result.Error)
				assert.NotEmpty(t, tt.user.ID, "User ID should not be empty")
				assert.Equal(t, tt.expected, tt.user.Username)

				// Validate UUID format
				_, err := uuid.Parse(tt.user.ID.String()) // Convert UUID to string before parsing
				assert.Nil(t, err, "User ID should be a valid UUID")
			}
		})
	}
}

// Test Post model
func TestPostModel(t *testing.T) {
	setup()
	defer teardown()

	user := User{Username: "john_doe", Email: "john@example.com", Password: "password123"}
	db.Create(&user)

	post := Post{
		Title:    "First Post",
		AuthorID: user.ID,
	}

	result := db.Create(&post)
	assert.Nil(t, result.Error)
	assert.NotEmpty(t, post.ID, "Post ID should not be empty")

	var fetchedPost Post
	db.First(&fetchedPost, post.ID)
	assert.Equal(t, post.Title, fetchedPost.Title)
	assert.Equal(t, post.AuthorID, fetchedPost.AuthorID)
}

// Test Report model with valid data only
func TestReportModel(t *testing.T) {
	setup()
	defer teardown()

	// Create a user
	user := User{Username: "john_report", Email: "report@example.com", Password: "password123"}
	result := db.Create(&user)
	assert.Nil(t, result.Error, "Failed to create user")
	assert.NotEmpty(t, user.ID, "User ID should not be empty")

	// Create a post
	post := Post{Title: "Report Post", AuthorID: user.ID}
	result = db.Create(&post)
	assert.Nil(t, result.Error, "Failed to create post")
	assert.NotEmpty(t, post.ID, "Post ID should not be empty")

	// Create a comment
	comment := Comment{Content: "Test comment", AuthorID: user.ID, PostID: post.ID}
	result = db.Create(&comment)
	assert.Nil(t, result.Error, "Failed to create comment")
	assert.NotEmpty(t, comment.ID, "Comment ID should not be empty")

	// Create a valid report for a post
	report := Report{
		UserID:    user.ID,
		PostID:    &post.ID, // Set PostID, leave CommentID nil
		CommentID: nil,
		Reason:    "Spam",
	}
	result = db.Create(&report)
	assert.Nil(t, result.Error, "Failed to create report for post")
	assert.NotEmpty(t, report.ID, "Report ID should not be empty")

	// Create a valid report for a comment
	reportWithComment := Report{
		UserID:    user.ID,
		PostID:    nil,
		CommentID: &comment.ID, // Set CommentID, leave PostID nil
		Reason:    "Harassment",
	}
	result = db.Create(&reportWithComment)
	assert.Nil(t, result.Error, "Failed to create report for comment")
	assert.NotEmpty(t, reportWithComment.ID, "Report ID should not be empty")
}

// Test many-to-many PostTag relationship
func TestPostTagModel(t *testing.T) {
	setup()
	defer teardown()

	user := User{Username: "john_tag", Email: "tag@example.com", Password: "password123"}
	db.Create(&user)

	post := Post{Title: "Tagged Post", AuthorID: user.ID}
	db.Create(&post)

	tag := Tag{Name: "Test Tag"}
	db.Create(&tag)

	// Create association
	err := db.Model(&post).Association("Tags").Append(&tag)
	assert.Nil(t, err, "Failed to associate tag with post")

	// Verify the association
	var fetchedPost Post
	db.Preload("Tags").First(&fetchedPost, post.ID)
	assert.Equal(t, 1, len(fetchedPost.Tags), "Post should have one associated tag")
	assert.Equal(t, tag.Name, fetchedPost.Tags[0].Name)
}

// Test Comment model
func TestCommentModel(t *testing.T) {
	setup()
	defer teardown()

	user := User{Username: "commenter", Email: "comment@example.com", Password: "password123"}
	db.Create(&user)

	post := Post{Title: "Post with Comment", AuthorID: user.ID}
	db.Create(&post)

	comment := Comment{
		Content:  "Sample Comment",
		AuthorID: user.ID,
		PostID:   post.ID,
	}

	result := db.Create(&comment)
	assert.Nil(t, result.Error)
	assert.NotEmpty(t, comment.ID, "Comment ID should not be empty")

	var fetchedComment Comment
	db.First(&fetchedComment, comment.ID)
	assert.Equal(t, comment.Content, fetchedComment.Content)
	assert.Equal(t, comment.AuthorID, fetchedComment.AuthorID)
}

// Test Like model
func TestLikeModel(t *testing.T) {
	setup()
	defer teardown()

	// Create user
	user := User{Username: "liker", Email: "like@example.com", Password: "password123"}
	db.Create(&user)

	// Create post
	post := Post{Title: "Liked Post", AuthorID: user.ID}
	db.Create(&post)

	// Create like
	like := Like{AuthorID: user.ID, PostID: post.ID}
	result := db.Create(&like)
	assert.Nil(t, result.Error)

	// Fetch like to verify composite key match
	var fetchedLike Like
	db.First(&fetchedLike, "author_id = ? AND post_id = ?", user.ID, post.ID)
	assert.Equal(t, like.AuthorID, fetchedLike.AuthorID)
	assert.Equal(t, like.PostID, fetchedLike.PostID)
}

func TestFollowerModel(t *testing.T) {
	setup()
	defer teardown()

	// Create two users
	user1 := User{Username: "follower1", Email: "follower1@example.com", Password: "password123"}
	user2 := User{Username: "followed1", Email: "followed1@example.com", Password: "password123"}

	db.Create(&user1)
	db.Create(&user2)

	// User1 follows User2
	follower := Follower{FollowerID: user1.ID, FollowedID: user2.ID}
	result := db.Create(&follower)
	assert.Nil(t, result.Error, "Failed to create follower relationship")

	// Check if the relationship exists
	var fetchedFollower Follower
	err := db.Where("follower_id = ? AND followed_id = ?", user1.ID, user2.ID).First(&fetchedFollower).Error
	assert.Nil(t, err, "Follower relationship not found in database")
	assert.Equal(t, user1.ID, fetchedFollower.FollowerID, "Follower ID mismatch")
	assert.Equal(t, user2.ID, fetchedFollower.FollowedID, "Followed ID mismatch")

	// Count followers and followings
	var followerCount int64
	db.Model(&Follower{}).Where("followed_id = ?", user2.ID).Count(&followerCount)
	assert.Equal(t, int64(1), followerCount, "User2 should have 1 follower")

	var followingCount int64
	db.Model(&Follower{}).Where("follower_id = ?", user1.ID).Count(&followingCount)
	assert.Equal(t, int64(1), followingCount, "User1 should be following 1 user")

	// Test unfollowing
	db.Delete(&fetchedFollower)
	var afterUnfollowCount int64
	db.Model(&Follower{}).Where("followed_id = ?", user2.ID).Count(&afterUnfollowCount)
	assert.Equal(t, int64(0), afterUnfollowCount, "User2 should have 0 followers after unfollow")
}
