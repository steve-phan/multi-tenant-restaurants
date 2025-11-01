package services

import (
	"errors"
	"strings"

	"restaurant-backend/internal/dto"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryService creates a new CategoryService instance
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(req *dto.CreateCategoryRequest, restaurantID uint) (*models.MenuCategory, error) {
	// Trim name
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	// Check if name already exists for this restaurant
	existing, _ := s.categoryRepo.GetByName(name)
	if existing != nil && existing.RestaurantID == restaurantID {
		return nil, errors.New("category name already taken")
	}

	category := &models.MenuCategory{
		RestaurantID: restaurantID,
		Name:         name,
		Description:  req.Description,
		DisplayOrder: req.DisplayOrder,
		IsActive:     req.IsActive,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates a category (only updates provided fields)
func (s *CategoryService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*models.MenuCategory, error) {
	// Verify category exists
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Build update map with only provided (non-nil) fields
	updates := make(map[string]interface{})

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return nil, errors.New("name cannot be empty")
		}
		updates["name"] = trimmed
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.DisplayOrder != nil {
		updates["display_order"] = *req.DisplayOrder
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Only update if there are fields to update
	if len(updates) == 0 {
		return category, nil // No changes
	}

	// Update the category
	if err := s.categoryRepo.Update(id, updates); err != nil {
		return nil, err
	}

	// Fetch and return updated category
	return s.categoryRepo.GetByID(id)
}
