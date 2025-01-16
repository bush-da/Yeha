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
	rg.GET("/:post_id/comments", commentController.GetComments)

	// Protected Routes (require authentication)
	protected := rg.Group("/:post_id/comment")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware only to these routes
	{
		protected.POST("/", commentController.CreateComment)              // Create Comment
		protected.DELETE("/:comment_id", commentController.DeleteComment) // Delete Comment
	}
}
