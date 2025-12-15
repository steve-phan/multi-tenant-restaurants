package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetTenantContext sets the PostgreSQL session variable for RLS
// This middleware must run after RequireAuth middleware
// Note: KAM and Admin users may not have a restaurant_id
func SetTenantContext(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get restaurant_id from context (set by auth middleware)
		restaurantIDValue, exists := c.Get(RestaurantIDKey)
		if !exists {
			c.JSON(500, gin.H{"error": "restaurant_id not found in context"})
			c.Abort()
			return
		}

		restaurantID, ok := restaurantIDValue.(uint)
		if !ok {
			c.JSON(500, gin.H{"error": "invalid restaurant_id type"})
			c.Abort()
			return
		}

		// Check if this is platform organization (KAM users)
		// For platform users, we still set the context but RLS policies handle them differently
		userRole, _ := c.Get(UserRoleKey)

		// Set the PostgreSQL role to restaurant_app_user for RLS policies to take effect
		// This ensures all queries run with the role that has RLS policies applied
		// Note: This must be done per-request, not at connection time (for migrations)
		db.Exec(`
			DO $$
			BEGIN
				IF EXISTS (SELECT FROM pg_roles WHERE rolname = 'restaurant_app_user') THEN
					SET ROLE restaurant_app_user;
				END IF;
			END $$;
		`)

		// Set the PostgreSQL session variable for RLS
		// This ensures all queries in this request are isolated to the tenant
		sql := fmt.Sprintf("SET app.current_restaurant = %d", restaurantID)
		if err := db.Exec(sql).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to set tenant context"})
			c.Abort()
			return
		}

		// Also set user role for RLS policies that check role
		if userRole != nil {
			roleSQL := fmt.Sprintf("SET app.current_user_role = '%s'", userRole.(string))
			_ = db.Exec(roleSQL).Error // Ignore error for role setting
		}

		// Mirror restaurant and role into request context to be accessible
		// by services/repositories that use context.Context directly.
		reqCtx := c.Request.Context()
		reqCtx = context.WithValue(reqCtx, RestaurantIDKey, restaurantID)
		if userRole != nil {
			reqCtx = context.WithValue(reqCtx, UserRoleKey, userRole.(string))
		}
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}
