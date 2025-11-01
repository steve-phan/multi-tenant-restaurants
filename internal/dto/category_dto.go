package dto

// CreateCategoryRequest represents a category creation request
type CreateCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

// UpdateCategoryRequest represents a category update request
// All fields are optional (pointers) - only provided fields will be updated
type UpdateCategoryRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	DisplayOrder *int    `json:"display_order"`
	IsActive     *bool   `json:"is_active"`
}
