package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yeha-api/controllers"
	"yeha-api/middleware"
)

func PostRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	postController := controllers.PostController{DB: db}

	// Public Route: Get all posts (no authentication required)
	rg.GET("/posts", postController.GetAllPosts)

	// Protected Routes (require authentication)
	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware only to these routes
	{
		protected.POST("/", postController.CreatePost)      // Create Post
		protected.GET("/:id", postController.GetPostByID)   // Get Post by ID
		protected.PUT("/:id", postController.UpdatePost)    // Update Post
		protected.DELETE("/:id", postController.DeletePost) // Delete Post
	}
}
