package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"yeha-api/models"
)

type ReportController struct {
	DB *gorm.DB
}

// Middleware to check admin status
func IsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("isAdmin")
	return exists && isAdmin.(bool)
}

// Report a Post or Comment
func (rc *ReportController) CreateReport(c *gin.Context) {
	var input struct {
		PostID    *uuid.UUID `json:"post_id"`
		CommentID *uuid.UUID `json:"comment_id"`
		Reason    string     `json:"reason" binding:"required"`
	}

	// Bind input
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding report input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get user ID from token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Validate report input
	if (input.PostID == nil && input.CommentID == nil) || (input.PostID != nil && input.CommentID != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide either PostID or CommentID, not both"})
		return
	}

	// Initialize report
	report := models.Report{
		UserID: userID.(uuid.UUID),
		Reason: input.Reason,
	}

	// Fetch AuthorID based on PostID or CommentID
	if input.PostID != nil {
		// Fetch post to get AuthorID
		var post models.Post
		if err := rc.DB.First(&post, "id = ?", input.PostID).Error; err != nil {
			log.Println("Error fetching post:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		report.PostID = input.PostID
		report.AuthorID = post.AuthorID // Assuming Post has a UserID field
	} else if input.CommentID != nil {
		// Fetch comment to get AuthorID
		var comment models.Comment
		if err := rc.DB.First(&comment, "id = ?", input.CommentID).Error; err != nil {
			log.Println("Error fetching comment:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		report.CommentID = input.CommentID
		report.AuthorID = comment.AuthorID // Assuming Comment has a UserID field
	}

	// Save report
	if err := rc.DB.Create(&report).Error; err != nil {
		log.Println("Error saving report:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit report"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Report submitted successfully", "report": report})
}

// Get all reports (Admin Only)
func (rc *ReportController) GetReports(c *gin.Context) {
	// if !IsAdmin(c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	// 	return
	// }

	var reports []models.Report
	if err := rc.DB.Preload("Post").Preload("Comment").Find(&reports).Error; err != nil {
		log.Println("Error fetching reports:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reports"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// Take action on a report (Admin Only)
func (rc *ReportController) TakeAction(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var input struct {
		Action string `json:"action" binding:"required"` // Accepts "delete" or "ignore"
	}

	// Parse report ID
	reportID := c.Param("report_id")
	id, err := uuid.Parse(reportID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	// Bind action input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	// Retrieve report
	var report models.Report
	if err := rc.DB.First(&report, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	// Handle actions using transactions
	tx := rc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	switch input.Action {
	case "delete":
		// Soft-delete related post or comment
		if report.PostID != nil {
			// Update the post_id in the report to NULL before deleting the post
			var post_id = report.PostID
			report.PostID = nil
			if err := tx.Save(&report).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
				return
			}

			// Now delete the post
			if err := tx.Delete(&models.Post{}, "id = ?", post_id).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
				return
			}
		} else if report.CommentID != nil {
			if err := tx.Delete(&models.Comment{}, "id = ?", report.CommentID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
				return
			}
		}
		report.ActionTaken = "deleted"

	case "ignore":
		report.ActionTaken = "ignored"

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	// Mark as reviewed
	report.Reviewed = true
	if err := tx.Save(&report).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Action completed successfully"})
}

// Delete a user (Admin Only)
func (rc *ReportController) DeleteUser(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Parse user ID
	userID := c.Param("user_id")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Soft-delete user
	if err := rc.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		log.Println("Error deleting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
