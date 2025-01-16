package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
	"yeha-api/models"
)

type PostController struct {
	DB *gorm.DB
}

// Struct for Post Response
type PostResponse struct {
	ID       uuid.UUID        `json:"id"`
	Title    string           `json:"title"`
	Contents []models.Content `json:"contents"`
	Tags     []models.Tag     `json:"tags"`
	Author   models.User      `json:"author"`
}

// GetAllPosts - Fetches all posts with associated data
func (pc *PostController) GetAllPosts(c *gin.Context) {
	var posts []models.Post

	// Preload relationships
	if err := pc.DB.Preload("Contents").Preload("Tags").Preload("Author").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts: " + err.Error()})
		return
	}

	// Format responses
	var postResponses []PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, PostResponse{
			ID:       post.ID,
			Title:    post.Title,
			Contents: post.Contents,
			Tags:     post.Tags,
			Author:   post.Author,
		})
	}

	c.JSON(http.StatusOK, postResponses)
}

func (pc *PostController) GetAllTags(c *gin.Context) {
	var tags []struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		PostCount int       `json:"post_count"`
	}

	// Query to fetch tags and count the associated posts
	if err := pc.DB.Table("tags").
		Select("tags.id, tags.name, tags.created_at, tags.updated_at, COUNT(post_tags.post_id) AS post_count").
		Joins("LEFT JOIN post_tags ON post_tags.tag_id = tags.id").
		Group("tags.id").
		Scan(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	// Return the tags with their associated post counts
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func (pc *PostController) CreatePost(c *gin.Context) {
	var post struct {
		Title    string           `json:"title"`
		Contents []models.Content `json:"contents"`
		Tags     []models.Tag     `json:"tags"`
	}

	// Bind the Post JSON (for title, contents, and tags)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post: " + err.Error()})
		return
	}

	// Create and associate contents with the post
	for _, content := range post.Contents {
		if content.Content == "" && content.FileURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content or file URL must be provided"})
			return
		}
		content.PostID = postModel.ID
		if err := pc.DB.Create(&content).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create content: " + err.Error()})
			return
		}
		postModel.Contents = append(postModel.Contents, content)
	}

	// Handle dynamic tags and associate them with the post
	for _, tag := range post.Tags {
		var existingTag models.Tag
		if err := pc.DB.Where("name = ?", tag.Name).First(&existingTag).Error; err == gorm.ErrRecordNotFound {
			// Tag doesn't exist, create a new one
			if err := pc.DB.Create(&tag).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag: " + err.Error()})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tag: " + err.Error()})
			return
		} else {
			tag = existingTag
		}

		// Associate the tag with the post
		if err := pc.DB.Model(&postModel).Association("Tags").Append(&tag); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags with post: " + err.Error()})
			return
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

func (pc *PostController) GetPostByID(c *gin.Context) {
	id := c.Param("post_id")
	var post models.Post

	// Preload related data for the post
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

	// Fetch followers of the post's author
	var followers []models.Follower
	if err := pc.DB.Where("followed_id = ?", post.AuthorID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	// Extract follower IDs
	followerIDs := make([]uuid.UUID, 0, len(followers))
	for _, follower := range followers {
		followerIDs = append(followerIDs, follower.FollowerID)
	}

	// Construct response
	response := gin.H{
		"post": gin.H{
			"id":       post.ID,
			"title":    post.Title,
			"author":   post.Author,
			"contents": post.Contents,
			"comments": post.Comments,
			"likes":    len(post.Likes),
			"tags":     post.Tags,
		},
		"followers": followerIDs, // Include author follower IDs
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePost - Updates an existing post by ID
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("post_id")
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
		Tags     []struct {
			Name string `json:"name"`
		} `json:"tags"`
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

	// Handle dynamic tags
	var updatedTags []models.Tag
	for _, tagInput := range requestBody.Tags {
		var tag models.Tag
		if err := pc.DB.Where("name = ?", tagInput.Name).First(&tag).Error; err == gorm.ErrRecordNotFound {
			// Tag doesn't exist, create it
			tag = models.Tag{Name: tagInput.Name}
			if err := pc.DB.Create(&tag).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag: " + err.Error()})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tag: " + err.Error()})
			return
		}

		// Add the tag to the list of updated tags
		updatedTags = append(updatedTags, tag)
	}

	// Replace post's tags with updated ones
	if err := pc.DB.Model(&post).Association("Tags").Replace(updatedTags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags with post"})
		return
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

// DeletePost - Deletes a post
func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("post_id")
	var post models.Post

	// Preload Tags to handle post-tag relationships
	if err := pc.DB.Preload("Tags").First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Authorization check
	userID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	if post.AuthorID != userID && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Start a transaction
	tx := pc.DB.Begin()

	// Delete the post-tag relationships explicitly
	if err := tx.Exec("DELETE FROM post_tags WHERE post_id = ?", post.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post-tag relationships"})
		return
	}

	// Check and delete orphaned tags
	for _, tag := range post.Tags {
		var count int64
		tx.Table("post_tags").Where("tag_id = ?", tag.ID).Count(&count)
		if count == 0 {
			if err := tx.Delete(&models.Tag{}, "id = ?", tag.ID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete orphaned tags"})
				return
			}
		}
	}

	// Delete the post
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	// Commit transaction
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
