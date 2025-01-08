package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"yeha-api/database"
	//	"yeha-api/middleware"
	"yeha-api/routes" // Import routes
)

func main() {
	// Initialize database connection
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Apply middleware
	// r.Use(middleware.AuthMiddleware())

	// Group routes under "/api"
	api := r.Group("/api")

	// Register user routes under "/api/users"
	userGroup := api.Group("/users")
	routes.UserRoutes(userGroup, db)

	// Register post routes under "/api/posts"
	postGroup := api.Group("/posts")
	routes.PostRoutes(postGroup, db)

	// Register comment routes under "/api/comments"
	commentGroup := api.Group("/comments")
	routes.CommentRoutes(commentGroup, db)

	// Register like routes under "/api/likes"
	likeGroup := api.Group("/likes")
	routes.LikeRoutes(likeGroup, db)

	// Register follow routes under "/api/follow"
	followGroup := api.Group("/follow")
	routes.FollowRoutes(followGroup, db)

	// Report routes under "/api/report"
	reportGroup := api.Group("/report")
	routes.ReportRoutes(reportGroup, db)

	// Start the server
	log.Println("Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
