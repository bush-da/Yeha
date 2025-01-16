package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"yeha-api/models"
)

type CommentController struct {
	DB *gorm.DB
}

// Create a new comment
func (cc *CommentController) CreateComment(c *gin.Context) {
	postID := c.Param("post_id") // Extract post ID from URL parameter
	var comment models.Comment

	// Bind the incoming JSON payload to the comment model
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get the user ID from the JWT token (authentication)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert postID to uuid.UUID
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Assign the author of the comment
	comment.AuthorID = userID.(uuid.UUID) // Updated field name
	comment.PostID = postUUID

	if len(comment.Content) == 0 || len(comment.Content) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment must be between 1 and 500 characters"})
		return
	}

	// Save the comment
	if err := cc.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully", "comment": comment})
}

// Get all comments for a post
func (cc *CommentController) GetComments(c *gin.Context) {
	postID := c.Param("post_id") // Extract post ID from URL parameter

	var comments []models.Comment

	// Preload both Author and Post, and Post's Author
	if err := cc.DB.Preload("Author").Preload("Post").Preload("Post.Author").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comments not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// Delete a comment
func (cc *CommentController) DeleteComment(c *gin.Context) {
	isAdmin, _ := c.Get("isAdmin")

	commentID := c.Param("comment_id") // Extract comment ID from URL parameter
	var comment models.Comment

	// Find the comment by ID
	if err := cc.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Authorization check: Ensure the user is the author of the comment
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if comment.AuthorID != userID.(uuid.UUID) && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this comment"})
		return
	}

	// Delete the comment
	if err := cc.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
