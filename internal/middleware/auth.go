package middleware

import (
	"net/http"
	"strings"

	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey       = "user_id"
	RestaurantIDKey = "restaurant_id"
	UserRoleKey     = "role"
	UserEmailKey    = "email"
)

// RequireAuth validates JWT token and extracts user context
func RequireAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Store user context in Gin context
		c.Set(UserIDKey, claims.UserID)
		c.Set(RestaurantIDKey, claims.RestaurantID)
		c.Set(UserRoleKey, claims.Role)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

// RequireRole checks if the authenticated user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(UserRoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found in context"})
			c.Abort()
			return
		}

		role := userRole.(string)
		hasRole := false
		for _, requiredRole := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireKAMOrAdmin checks if the authenticated user is a KAM or Admin
func RequireKAMOrAdmin() gin.HandlerFunc {
	return RequireRole("KAM", "Admin")
}
