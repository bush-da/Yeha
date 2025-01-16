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

// Get all follower relationships
func (fc *FollowController) GetAllFollowers(c *gin.Context) {
	var followers []models.Follower

	// Fetch all follower relationships from the database
	if err := fc.DB.Find(&followers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	// Return the raw list of followers
	c.JSON(http.StatusOK, gin.H{"followers": followers})
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

	// Parse the user ID to UUID format
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Initialize an empty slice for followers
	var followers []models.Follower

	// Preload the related User data (follower and followed)
	if err := fc.DB.Preload("Follower").Where("followed_id = ?", userUUID).Find(&followers).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Followers not found"})
		return
	}

	// Build the customized response
	var followerDetails []map[string]interface{}
	for _, follow := range followers {
		followerDetails = append(followerDetails, map[string]interface{}{
			"id":              follow.Follower.ID,
			"username":        follow.Follower.Username,
			"email":           follow.Follower.Email,
			"profile_picture": follow.Follower.ProfilePicture,
			"bio":             follow.Follower.Bio,
			"gender":          follow.Follower.Gender,
		})
	}

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"user_id":   userUUID,
		"followers": followerDetails,
	})
}

// Get all users a user is following
func (fc *FollowController) GetFollowing(c *gin.Context) {
	userID := c.Param("user_id")

	// Parse the user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve the list of users the current user is following
	var following []models.Follower
	if err := fc.DB.Preload("Followed").Where("follower_id = ?", userUUID).Find(&following).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Following not found"})
		return
	}

	// Build the customized response
	var followingDetails []map[string]interface{}
	for _, follow := range following {
		followingDetails = append(followingDetails, map[string]interface{}{
			"id":              follow.Followed.ID,
			"username":        follow.Followed.Username,
			"email":           follow.Followed.Email,
			"profile_picture": follow.Followed.ProfilePicture,
			"bio":             follow.Followed.Bio,
			"gender":          follow.Followed.Gender,
		})
	}

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"user_id":   userUUID,
		"following": followingDetails,
	})
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
