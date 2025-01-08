package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"yeha-api/models"
)

// Load the secret key from environment variable
var secretKey []byte

// Initialize the secret key
func init() {
	secretKeyStr := os.Getenv("JWT_SECRET_KEY")
	if secretKeyStr == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}
	secretKey = []byte(secretKeyStr)
}

// GenerateToken - Creates JWT token with user ID and admin status
func generateToken(userID uuid.UUID, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"userID":  userID.String(),
		"isAdmin": isAdmin,                               // Add admin claim
		"exp":     time.Now().Add(72 * time.Hour).Unix(), // Expiry time
	}

	// Generate token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// UserController - Dependency Injection for database
type UserController struct {
	DB *gorm.DB
}

// NewUserController - Initializes a new user controller
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// RegisterUser - Handles user registration
func (uc *UserController) RegisterUser(c *gin.Context) {
	// Input validation structure
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := uc.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user object
	user := models.User{
		ID:             uuid.New(),
		Username:       input.Username,
		Email:          input.Email,
		Password:       string(hashedPassword),
		ProfilePicture: "",
		Bio:            "",
		IsAdmin:        false, // Default to non-admin
	}

	// Save user to database
	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// LoginUser - Handles user login
func (uc *UserController) LoginUser(c *gin.Context) {
	// Input validation structure
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := uc.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Avoid revealing whether email exists for security reasons
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID, user.IsAdmin) // Include admin status in token
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return success with token and user details
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"isAdmin":  user.IsAdmin, // Include isAdmin in response
		},
	})
}

// GetUserProfile - Fetch user profile by ID
func (uc *UserController) GetUserProfile(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find user
	var user models.User
	if err := uc.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return user profile
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserProfile - Update user profile
func (uc *UserController) UpdateUserProfile(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find user
	var user models.User
	if err := uc.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Input for update
	var input struct {
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
		Bio            string `json:"bio"`
		Gender         string `json:"gender"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	user.Username = input.Username
	user.ProfilePicture = input.ProfilePicture
	user.Bio = input.Bio
	user.Gender = input.Gender
	// user.UpdatedAt = time.Now()

	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
}

// DeleteUser - Delete user by ID
// func (uc *UserController) DeleteUser(c *gin.Context) {
// 	id := c.Param("id")

// 	// Parse UUID
// 	userID, err := uuid.Parse(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	// Delete user
// 	if err := uc.DB.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
//	}
//
// DeleteUser - Delete user by ID (Self-deletion allowed, Admin can delete any user)
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user ID from the JWT token (authentication)
	authUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the authenticated user is the same as the user to be deleted or an admin
	var authUser models.User
	if err := uc.DB.First(&authUser, "id = ?", authUserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if the user is trying to delete their own account or if the user is an admin
	if authUserID != userID.String() && !authUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this user"})
		return
	}

	// Delete user
	if err := uc.DB.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
