package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yeha-api/controllers"
	"yeha-api/middleware"
)

func LikeRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	likeController := controllers.LikeController{DB: db}

	// Public Route: Get all likes for a post
	rg.GET("/:post_id/likes", likeController.GetLikes)

	// Public Route: Check if a user liked a specific post
	rg.GET("/:post_id/like/:user_id", likeController.IsPostLikedByUser)

	// Protected Routes (require authentication)
	protected := rg.Group("/:post_id/like")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware only to these routes
	{
		protected.POST("", likeController.CreateLike)            // Like a post
		protected.DELETE("/:like_id", likeController.DeleteLike) // Delete like
	}
}
