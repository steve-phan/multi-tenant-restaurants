package dto

// CreateUserDTO represents the data for creating a user
type CreateUserDTO struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Role        string `json:"role" binding:"required,oneof=Admin Staff Client"`
	Phone       string `json:"phone,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	Language    string `json:"language,omitempty"`
	Preferences string `json:"preferences,omitempty"` // JSON string
}

// UpdateUserDTO represents the data for updating a user
type UpdateUserDTO struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Role        string `json:"role,omitempty" binding:"omitempty,oneof=Admin Staff Client"`
	Phone       string `json:"phone,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	Language    string `json:"language,omitempty"`
	Preferences string `json:"preferences,omitempty"` // JSON string
}

// UpdateUserStatusDTO represents the data for updating user status
type UpdateUserStatusDTO struct {
	IsActive bool `json:"is_active" binding:"required"`
}

// UpdateProfileDTO represents the data for updating current user's profile
type UpdateProfileDTO struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	Language  string `json:"language,omitempty"`
}

// ChangePasswordDTO represents the data for changing password
type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// UpdatePreferencesDTO represents the data for updating user preferences
type UpdatePreferencesDTO struct {
	Preferences string `json:"preferences" binding:"required"` // JSON string
}
