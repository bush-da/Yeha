package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"yeha-api/models"
)

type ContentController struct {
	DB *gorm.DB
}

type ContentResponse struct {
	ID        string `json:"id"`
	FileURL   string `json:"file_url"`
	Paragraph int    `json:"paragraph"`
	Content   string `json:"content"`
	PostID    string `json:"post_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetContentsByPostID - Fetches all contents associated with a specific post
func (cc *ContentController) GetContentsByPostID(c *gin.Context) {
	postID := c.Param("post_id")

	var contents []models.Content
	if err := cc.DB.Where("post_id = ?", postID).Find(&contents).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contents not found"})
		return
	}

	response := make([]ContentResponse, len(contents))
	for i, content := range contents {
		response[i] = ContentResponse{
			ID:        content.ID.String(),
			FileURL:   content.FileURL,
			Paragraph: content.Paragraph,
			Content:   content.Content,
			PostID:    content.PostID.String(),
			CreatedAt: content.CreatedAt.String(),
			UpdatedAt: content.UpdatedAt.String(),
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetContentByID - Fetches a content block by its ID
func (cc *ContentController) GetContentByID(c *gin.Context) {
	contentID := c.Param("content_id")

	var content models.Content
	if err := cc.DB.First(&content, "id = ?", contentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	response := ContentResponse{
		ID:        content.ID.String(),
		FileURL:   content.FileURL,
		Paragraph: content.Paragraph,
		Content:   content.Content,
		PostID:    content.PostID.String(),
		CreatedAt: content.CreatedAt.String(),
		UpdatedAt: content.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, response)
}

// CreateContent - Creates a new content block
func (cc *ContentController) CreateContent(c *gin.Context) {
	var content models.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cc.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create content"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Content created successfully"})
}

// UpdateContent - Updates an existing content block
func (cc *ContentController) UpdateContent(c *gin.Context) {
	contentID := c.Param("content_id")
	var existingContent models.Content

	// Find the content by ID
	if err := cc.DB.First(&existingContent, "id = ?", contentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	// Bind the incoming JSON payload
	var updatedData struct {
		FileURL   *string `json:"file_url,omitempty"`
		Content   *string `json:"content,omitempty"`
		Paragraph *int    `json:"paragraph,omitempty"` // Use `int` to match model field type
	}
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the content fields selectively
	if updatedData.FileURL != nil {
		existingContent.FileURL = *updatedData.FileURL
	}
	if updatedData.Content != nil {
		existingContent.Content = *updatedData.Content
	}
	if updatedData.Paragraph != nil {
		existingContent.Paragraph = *updatedData.Paragraph
	}

	// Save the updated content (GORM auto-updates `updated_at`)
	if err := cc.DB.Save(&existingContent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update content", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content updated successfully"})
}

// DeleteContent - Deletes a content block by its ID
func (cc *ContentController) DeleteContent(c *gin.Context) {
	contentID := c.Param("content_id")

	var content models.Content
	if err := cc.DB.First(&content, "id = ?", contentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	if err := cc.DB.Delete(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Content deleted successfully"})
}
