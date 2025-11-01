package services

import (
	"errors"

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
func (s *MenuItemService) CreateMenuItem(req *dto.CreateMenuItemRequest, restaurantID uint) (*models.MenuItem, error) {
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

	menuItem := &models.MenuItem{
		RestaurantID: restaurantID,
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		ImageURL:     req.ImageURL,
		DisplayOrder: req.DisplayOrder,
		IsAvailable:  req.IsAvailable,
	}

	if err := s.menuItemRepo.Create(menuItem); err != nil {
		return nil, err
	}

	// Fetch created item with relationships
	return s.menuItemRepo.GetByID(menuItem.ID)
}

// UpdateMenuItem updates a menu item (only updates provided fields)
func (s *MenuItemService) UpdateMenuItem(id uint, req *dto.UpdateMenuItemRequest) (*models.MenuItem, error) {
	// Verify menu item exists
	menuItem, err := s.menuItemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("menu item not found")
	}

	// Build update map with only provided (non-nil) fields
	updates := make(map[string]interface{})

	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("name cannot be empty")
		}
		updates["name"] = *req.Name
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return nil, errors.New("price cannot be negative")
		}
		updates["price"] = *req.Price
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
	if err := s.menuItemRepo.Update(id, updates); err != nil {
		return nil, err
	}

	// Fetch and return updated menu item
	return s.menuItemRepo.GetByID(id)
}
