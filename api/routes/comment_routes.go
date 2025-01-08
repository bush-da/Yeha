package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yeha-api/controllers"
	"yeha-api/middleware"
)

func CommentRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	commentController := controllers.CommentController{DB: db}

	// Public Route: Get all comments for a post (no authentication required)
	rg.GET("/posts/:post_id/comments", commentController.GetComments)

	// Protected Routes (require authentication)
	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware only to these routes
	{
		protected.POST("/:post_id", commentController.CreateComment)      // Create Comment
		protected.DELETE("/:comment_id", commentController.DeleteComment) // Delete Comment
	}
}
