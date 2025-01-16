package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"yeha-api/models"
)

type TagController struct {
	DB *gorm.DB
}

// DeleteTagsByPostID - Deletes all tags associated with a specific post
func (tc *TagController) DeleteTagsByPostID(c *gin.Context) {
	postID := c.Param("post_id")
	var post models.Post

	// Find the post by ID
	if err := tc.DB.Preload("Tags").First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Remove the association between the post and tags
	if err := tc.DB.Model(&post).Association("Tags").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags removed successfully"})
}

// CreateTag - Creates a new tag
func (tc *TagController) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the new tag in the database
	if err := tc.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tag created successfully"})
}
