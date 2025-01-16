package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yeha-api/controllers"
	"yeha-api/middleware"
)

// UserRoutes defines routes related to user management
func UserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	// Public Routes
	rg.POST("/register", userController.RegisterUser)
	rg.POST("/login", userController.LoginUser)

	// Protected Routes
	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Add auth middleware here

	// Require authentication for the following routes
	protected.GET("/", userController.GetAllUsers) // New route to fetch all users
	protected.GET("/:id", userController.GetUserProfile)
	protected.PUT("/:id", userController.UpdateUserProfile)
	protected.PUT("/:id/password", userController.UpdatePassword)
	protected.DELETE("/:id", userController.DeleteUser)
}
