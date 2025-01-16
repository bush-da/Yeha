package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yeha-api/controllers"
	"yeha-api/middleware"
)

func PostRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	postController := controllers.PostController{DB: db}
	tagController := controllers.TagController{DB: db}
	contentController := controllers.ContentController{DB: db}

	// Public Routes
	rg.GET("/", postController.GetAllPosts)
	rg.GET("/posts/:post_id", postController.GetPostByID)
	// the route for fetching all tags
	rg.GET("/tags", postController.GetAllTags)

	// Routes for content
	rg.GET("/contents/post/:post_id", contentController.GetContentsByPostID) // Get contents by post ID
	rg.GET("/contents/:content_id", contentController.GetContentByID)        // Get a content block by ID

	// Routes for tags
	rg.DELETE("/tags/post/:post_id", tagController.DeleteTagsByPostID)
	rg.POST("/tags", tagController.CreateTag)

	// Protected Routes
	protected := rg.Group("/posts")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/", postController.CreatePost)
		protected.PUT("/:post_id", postController.UpdatePost)
		protected.DELETE("/:post_id", postController.DeletePost)

		// Routes for content
		protected.POST("/contents", contentController.CreateContent)
		protected.PUT("/contents/:content_id", contentController.UpdateContent)
		protected.DELETE("/contents/:content_id", contentController.DeleteContent)
	}
}
