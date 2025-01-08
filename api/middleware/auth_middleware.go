package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Load the secret key from environment variable
var secretKey []byte

func init() {
	secretKeyStr := os.Getenv("JWT_SECRET_KEY")
	if secretKeyStr == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}
	secretKey = []byte(secretKeyStr)
}

// GenerateToken - Create JWT token with userID and isAdmin
func GenerateToken(userID uuid.UUID, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"userID":  userID.String(),
		"isAdmin": isAdmin, // Include isAdmin in token claims
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// AuthMiddleware - Middleware to validate JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Get user ID from claims
		userIDStr, ok := claims["userID"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		// Get isAdmin from claims
		isAdmin, ok := claims["isAdmin"].(bool) // Extract isAdmin claim
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin status in token"})
			c.Abort()
			return
		}

		// Store user ID and isAdmin in context
		c.Set("userID", userID)
		c.Set("isAdmin", isAdmin) // Store isAdmin in context
		c.Next()
	}
}

// IsAdmin - Helper function to check if the user is an admin
func IsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("isAdmin")
	if !exists {
		return false
	}

	admin, ok := isAdmin.(bool)
	return ok && admin
}
