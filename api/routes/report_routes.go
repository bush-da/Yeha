package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yeha-api/controllers"
	"yeha-api/middleware"
)

func ReportRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	reportController := controllers.ReportController{DB: db}

	// Routes for reports
	protected := rg.Group("/reports")
	protected.Use(middleware.AuthMiddleware())
	{
		// User can create reports
		protected.POST("/", reportController.CreateReport)

		// Admin routes
		protected.GET("/", reportController.GetReports)                   // View reports
		protected.POST("/:report_id/action", reportController.TakeAction) // Take action
		protected.DELETE("/user/:user_id", reportController.DeleteUser)   // Delete user
	}
}
