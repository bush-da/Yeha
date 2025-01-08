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
	rg.GET("/posts/:post_id/likes", likeController.GetLikes)

	// Protected Routes (require authentication)
	protected := rg.Group("/posts/:post_id/likes")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware only to these routes
	{
		protected.POST("/", likeController.CreateLike)           // Like a post
		protected.DELETE("/:like_id", likeController.DeleteLike) // Delete like
	}
}
