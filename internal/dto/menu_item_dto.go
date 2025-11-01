package dto

// CreateMenuItemRequest represents a menu item creation request
type CreateMenuItemRequest struct {
	CategoryID   uint    `json:"category_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" binding:"required,min=0"`
	ImageURL     string  `json:"image_url"`
	DisplayOrder int     `json:"display_order"`
	IsAvailable  bool    `json:"is_available"`
}

// UpdateMenuItemRequest represents a menu item update request
// All fields are optional (pointers) - only provided fields will be updated
type UpdateMenuItemRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price"`
	ImageURL     *string  `json:"image_url"`
	DisplayOrder *int     `json:"display_order"`
	IsAvailable  *bool    `json:"is_available"`
	CategoryID   *uint    `json:"category_id"`
}
