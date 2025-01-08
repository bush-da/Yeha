package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"yeha-api/models"
)

type PostController struct {
	DB *gorm.DB
}

type PostResponse struct {
	ID       uuid.UUID        `json:"id"`
	Title    string           `json:"title"`
	Contents []models.Content `json:"contents"`
	Tags     []models.Tag     `json:"tags"`
	Author   models.User      `json:"author"`
}

func (pc *PostController) CreatePost(c *gin.Context) {
	var post struct {
		Title    string           `json:"title"`
		Contents []models.Content `json:"contents"`
		Tags     []models.Tag     `json:"tags"`
	}

	// Bind the Post JSON (for title and other fields, contents, and tags)
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract userID from context (ensures post creation is tied to a logged-in user)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Create the Post
	postModel := models.Post{
		Title:    post.Title,
		AuthorID: userUUID,
	}

	if err := pc.DB.Create(&postModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Create and associate contents with the post
	for _, content := range post.Contents {
		content.PostID = postModel.ID
		if err := pc.DB.Create(&content).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create content"})
			return
		}
		// Append each created content to the post
		postModel.Contents = append(postModel.Contents, content)
	}

	// Handle dynamic tags and associate them with the post
	if len(post.Tags) > 0 {
		// Loop through the tags and associate them with the post
		for _, tag := range post.Tags {
			// Check if the tag already exists
			if err := pc.DB.Where("name = ?", tag.Name).First(&tag).Error; err != nil {
				// If the tag doesn't exist, create it
				if err := pc.DB.Create(&tag).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
					return
				}
			}

			// Associate the tag with the post
			if err := pc.DB.Model(&postModel).Association("Tags").Append(&tag); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags with post"})
				return
			}
		}
	}

	// Return the created post with associated contents and tags
	response := PostResponse{
		ID:       postModel.ID,
		Title:    postModel.Title,
		Contents: postModel.Contents,
		Tags:     postModel.Tags,
		Author:   postModel.Author, // Include author directly
	}

	c.JSON(http.StatusCreated, response)
}

// GetPostByID - Fetches a post by its ID

func (pc *PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	// Preload related data
	if err := pc.DB.Preload("Author"). // Load Author data
						Preload("Contents").             // Load Contents
						Preload("Comments").             // Load Comments
						Preload("Likes").                // Load Likes
						Preload("Tags").                 // Load Tags
						Preload("Contents.Post").        // Prevent nested empty post in content
						Preload("Contents.Post.Author"). // Load nested Author in Content
						First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost - Updates an existing post by ID
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	// Fetch the existing post with relationships
	if err := pc.DB.Preload("Contents").Preload("Tags").First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if the logged-in user is the owner of the post
	userID, _ := c.Get("userID")
	authorUUID, err := uuid.Parse(post.AuthorID.String())
	if err != nil || authorUUID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this post"})
		return
	}

	// Parse the request body into a temporary struct
	var requestBody struct {
		Title    string           `json:"title"`
		Contents []models.Content `json:"contents"`
		TagIDs   []uuid.UUID      `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the post title
	post.Title = requestBody.Title

	// Handle content updates
	var existingContents []models.Content
	pc.DB.Where("post_id = ?", post.ID).Find(&existingContents)

	// Map existing content by ID for quick lookup
	contentMap := make(map[uuid.UUID]models.Content)
	for _, content := range existingContents {
		contentMap[content.ID] = content
	}

	// Process contents in the request
	for _, content := range requestBody.Contents {
		if content.ID != uuid.Nil { // Check if content already exists
			if existingContent, found := contentMap[content.ID]; found {
				// Update existing content
				existingContent.FileURL = content.FileURL
				existingContent.Content = content.Content
				existingContent.Paragraph = content.Paragraph
				existingContent.ContentType = content.ContentType
				pc.DB.Save(&existingContent)   // Save updates
				delete(contentMap, content.ID) // Remove from map to track processed items
			} else {
				// Invalid ID provided, handle error
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
				return
			}
		} else {
			// Create new content
			content.PostID = post.ID
			pc.DB.Create(&content)
		}
	}

	// Delete any remaining contents not in the request
	for _, content := range contentMap {
		pc.DB.Delete(&content)
	}

	// Update tags
	var updatedTags []models.Tag
	if len(requestBody.TagIDs) > 0 {
		if err := pc.DB.Find(&updatedTags, requestBody.TagIDs).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tags provided"})
			return
		}
	}
	if len(updatedTags) > 0 {
		pc.DB.Model(&post).Association("Tags").Replace(updatedTags)
	} else {
		pc.DB.Model(&post).Association("Tags").Clear()
	}

	// Save the updated post
	if err := pc.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Reload the updated post with relationships, but avoid preloading "post" inside "contents"
	pc.DB.Preload("Author").
		Preload("Contents", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Post") // Ignore loading the Post field
		}).
		Preload("Tags").
		First(&post, "id = ?", id)

	c.JSON(http.StatusOK, post)
}

func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")

	// Find the post
	var post models.Post
	if err := pc.DB.Preload("Tags").First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Authorization check
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authorUUID, _ := uuid.Parse(post.AuthorID.String())
	if authorUUID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
		return
	}

	// Begin transaction
	tx := pc.DB.Begin()

	// Cleanup related entities
	if err := tx.Where("post_id = ?", post.ID).Delete(&models.Content{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contents"})
		return
	}
	if err := tx.Where("post_id = ?", post.ID).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comments"})
		return
	}
	if err := tx.Where("post_id = ?", post.ID).Delete(&models.Like{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete likes"})
		return
	}
	if err := tx.Where("post_id = ?", post.ID).Delete(&models.Report{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reports"})
		return
	}

	// Remove tag associations
	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear tags"})
		return
	}

	// Finally, delete the post itself
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// Commit transaction
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetAllPosts - Fetches all posts
func (pc *PostController) GetAllPosts(c *gin.Context) {
	var posts []models.Post
	if err := pc.DB.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
