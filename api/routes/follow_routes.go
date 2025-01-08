package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yeha-api/controllers"
	"yeha-api/middleware"
)

func FollowRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	followController := controllers.FollowController{DB: db}

	// Public Routes: Get followers of a user
	rg.GET("/users/:user_id/followers", followController.GetFollowers)

	// Public Routes: Get users a user is following
	rg.GET("/users/:user_id/following", followController.GetFollowing)

	// Protected Routes (require authentication)
	protected := rg.Group("/users/:user_id/follow")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to these routes
	{
		protected.POST("/", followController.FollowUser)     // Follow a user
		protected.DELETE("/", followController.UnfollowUser) // Unfollow a user
	}
}
