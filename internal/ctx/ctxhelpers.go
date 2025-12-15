package ctx

import (
	"context"

	"restaurant-backend/internal/middleware"
)

// GetUserID returns the user ID from context if present
func GetUserID(ctx context.Context) (uint, bool) {
	if ctx == nil {
		return 0, false
	}
	v := ctx.Value(middleware.UserIDKey)
	if v == nil {
		return 0, false
	}
	uid, ok := v.(uint)
	return uid, ok
}

// GetRestaurantID returns the restaurant ID from context if present
func GetRestaurantID(ctx context.Context) (uint, bool) {
	if ctx == nil {
		return 0, false
	}
	v := ctx.Value(middleware.RestaurantIDKey)
	if v == nil {
		return 0, false
	}
	rid, ok := v.(uint)
	return rid, ok
}

// GetUserRole returns the user role from context if present
func GetUserRole(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	v := ctx.Value(middleware.UserRoleKey)
	if v == nil {
		return "", false
	}
	role, ok := v.(string)
	return role, ok
}

// GetUserEmail returns the user email from context if present
func GetUserEmail(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	v := ctx.Value(middleware.UserEmailKey)
	if v == nil {
		return "", false
	}
	email, ok := v.(string)
	return email, ok
}
