package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"yeha-api/models"
)

type LikeController struct {
	DB *gorm.DB
}

// Create a like for a post
func (lc *LikeController) CreateLike(c *gin.Context) {
	postID := c.Param("post_id") // Extract post ID from URL parameter
	var like models.Like

	// Convert postID from string to uuid.UUID
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Get the user ID from the JWT token (authentication)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the user has already liked this post
	var existingLike models.Like
	if err := lc.DB.Where("author_id = ? AND post_id = ?", userID.(uuid.UUID), postUUID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already liked this post"})
		return
	}

	// Create a new like
	like.AuthorID = userID.(uuid.UUID)
	like.PostID = postUUID

	// Save the like
	if err := lc.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully", "like": like})
}

// Get all likes for a post
func (lc *LikeController) GetLikes(c *gin.Context) {
	postID := c.Param("post_id") // Extract post ID from URL parameter

	// Convert postID from string to uuid.UUID
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var likes []models.Like
	if err := lc.DB.Where("post_id = ?", postUUID).Find(&likes).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Likes not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes})
}

// Delete a like from a post
func (lc *LikeController) DeleteLike(c *gin.Context) {
	likeID := c.Param("like_id") // Extract like ID from URL parameter
	postID := c.Param("post_id")
	var like models.Like

	// Find the like by ID
	if err := lc.DB.First(&like, "id = ? AND post_id = ?", likeID, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	// Authorization check: Ensure the user is the one who liked the post
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if like.AuthorID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this like"})
		return
	}

	// Delete the like
	if err := lc.DB.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
}

// Check if a specific user liked a post
func (lc *LikeController) IsPostLikedByUser(c *gin.Context) {
	postID := c.Param("post_id") // Extract post ID from URL parameter
	userID := c.Param("user_id") // Extract user ID from URL parameter

	// Convert IDs from string to uuid.UUID
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if the like exists
	var like models.Like
	if err := lc.DB.Where("post_id = ? AND author_id = ?", postUUID, userUUID).First(&like).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"liked": false}) // User has not liked the post
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": true}) // User has liked the post
}
