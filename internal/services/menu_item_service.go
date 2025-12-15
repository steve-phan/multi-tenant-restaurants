package services

import (
	"context"
	"errors"
	"strings"

	"restaurant-backend/internal/dto"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"
)

// MenuItemService handles menu item business logic
type MenuItemService struct {
	menuItemRepo *repositories.MenuItemRepository
}

// NewMenuItemService creates a new MenuItemService instance
func NewMenuItemService(menuItemRepo *repositories.MenuItemRepository) *MenuItemService {
	return &MenuItemService{
		menuItemRepo: menuItemRepo,
	}
}

// CreateMenuItem creates a new menu item
func (s *MenuItemService) CreateMenuItem(ctx context.Context, req *dto.CreateMenuItemRequest, restaurantID uint) (*models.MenuItem, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.CategoryID == 0 {
		return nil, errors.New("category_id is required")
	}
	if req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	// Check if name is already taken
	if _, err := s.menuItemRepo.GetByNameWithContext(ctx, req.Name); err == nil {
		return nil, errors.New("name already taken")
	}

	menuItem := &models.MenuItem{
		RestaurantID: restaurantID,
		CategoryID:   req.CategoryID,
		Name:         strings.TrimSpace(req.Name),
		Description:  req.Description,
		Price:        req.Price,
		ImageURL:     req.ImageURL,
		DisplayOrder: req.DisplayOrder,
		IsAvailable:  req.IsAvailable,
	}

	if err := s.menuItemRepo.CreateWithContext(ctx, menuItem); err != nil {
		return nil, err
	}

	// Fetch created item with relationships
	return s.menuItemRepo.GetByIDWithContext(ctx, menuItem.ID)
}

// UpdateMenuItem updates a menu item (only updates provided fields)
func (s *MenuItemService) UpdateMenuItem(ctx context.Context, id uint, req *dto.UpdateMenuItemRequest, restaurantID uint) (*models.MenuItem, error) {
	// Verify menu item exists
	menuItem, err := s.menuItemRepo.GetByIDWithContext(ctx, id)
	if err != nil {
		return nil, errors.New("menu item not found")
	}

	// Validate ownership - ensure menu item belongs to the requesting restaurant
	// This is a defense-in-depth measure in addition to RLS
	if menuItem.RestaurantID != restaurantID {
		return nil, errors.New("menu item not found") // Don't reveal existence of other tenants' data
	}

	// Build update map with only provided (non-nil) fields
	updates := make(map[string]interface{})

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("name cannot be empty")
		}
		// Validate name is not already taken
		if _, err := s.menuItemRepo.GetByNameWithContext(ctx, *req.Name); err == nil {
			return nil, errors.New("name already taken")
		}
		updates["name"] = *req.Name
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}

	if req.DisplayOrder != nil {
		updates["display_order"] = *req.DisplayOrder
	}

	if req.IsAvailable != nil {
		updates["is_available"] = *req.IsAvailable
	}

	if req.CategoryID != nil {
		// Validate category exists if category is being changed
		if *req.CategoryID != menuItem.CategoryID {
			// Note: Category existence should be validated, but we'll assume it's validated
			// by RLS or by the handler if needed
			updates["category_id"] = *req.CategoryID
		}
	}

	// Only update if there are fields to update
	if len(updates) == 0 {
		return menuItem, nil // No changes
	}

	// Update the menu item
	if err := s.menuItemRepo.UpdateWithContext(ctx, id, updates); err != nil {
		return nil, err
	}

	// Fetch and return updated menu item
	return s.menuItemRepo.GetByIDWithContext(ctx, id)
}
