package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"yeha-api/models"
)

type FollowController struct {
	DB *gorm.DB
}

// Follow a user
func (fc *FollowController) FollowUser(c *gin.Context) {
	followingID := c.Param("user_id")
	var follower models.Follower

	// Convert followingID to UUID
	followingUUID, err := uuid.Parse(followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userID.(uuid.UUID) == followingUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself"})
		return
	}

	// Check if already following
	var existingFollow models.Follower
	if err := fc.DB.Where("follower_id = ? AND followed_id = ?", userID.(uuid.UUID), followingUUID).First(&existingFollow).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already following this user"})
		return
	}

	// Create follow record
	follower.FollowerID = userID.(uuid.UUID)
	follower.FollowedID = followingUUID // Fixed field

	if err := fc.DB.Create(&follower).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed user successfully", "follower": follower})
}

// Get all followers of a user
func (fc *FollowController) GetFollowers(c *gin.Context) {
	userID := c.Param("user_id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var followers []models.Follower
	if err := fc.DB.Where("followed_id = ?", userUUID).Find(&followers).Error; err != nil { // Fixed field
		c.JSON(http.StatusNotFound, gin.H{"error": "Followers not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"followers": followers})
}

// Get all users a user is following
func (fc *FollowController) GetFollowing(c *gin.Context) {
	userID := c.Param("user_id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var following []models.Follower
	if err := fc.DB.Where("follower_id = ?", userUUID).Find(&following).Error; err != nil { // FollowerID remains
		c.JSON(http.StatusNotFound, gin.H{"error": "Following not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": following})
}

// Unfollow a user
func (fc *FollowController) UnfollowUser(c *gin.Context) {
	followingID := c.Param("user_id")
	var follower models.Follower

	followingUUID, err := uuid.Parse(followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userID.(uuid.UUID) == followingUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot unfollow yourself"})
		return
	}

	if err := fc.DB.Where("follower_id = ? AND followed_id = ?", userID.(uuid.UUID), followingUUID).Delete(&follower).Error; err != nil { // Fixed field
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed user successfully"})
}
