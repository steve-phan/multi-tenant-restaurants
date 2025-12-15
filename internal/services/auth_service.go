package services

import (
	"context"
	"errors"
	"time"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	db       *gorm.DB
	config   *config.Config
	userRepo *repositories.UserRepository
}

// NewAuthService creates a new AuthService instance
func NewAuthService(db *gorm.DB, cfg *config.Config, userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		db:       db,
		config:   cfg,
		userRepo: userRepo,
	}
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	RestaurantID uint   `json:"restaurant_id"` // Always present (KAMs belong to Platform Organization)
	Email        string `json:"email"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Use repository to load user (preloads Restaurant)
	user, err := s.userRepo.GetByEmailGlobalWithContext(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// RegisterRequest represents registration request payload
// Note: KAM role is NOT allowed here - must use CreateKAM endpoint
type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=8"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Role         string `json:"role" binding:"required,oneof=Admin Staff Client"`
	RestaurantID uint   `json:"restaurant_id" binding:"required"`
}

// Register creates a new user account (for restaurant users only)
// KAM users cannot be created via this endpoint - use CreateKAM endpoint
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	// KAM role is not allowed in regular registration
	if req.Role == "KAM" {
		return nil, errors.New("KAM users cannot be created via this endpoint. Use the KAM creation endpoint instead")
	}
	// Verify restaurant exists and is active
	var restaurant models.Restaurant
	if err := s.db.WithContext(ctx).First(&restaurant, req.RestaurantID).Error; err != nil {
		return nil, errors.New("restaurant not found")
	}

	if restaurant.Status != models.RestaurantStatusActive {
		return nil, errors.New("restaurant is not active")
	}

	// Check if user already exists (use repository)
	if existing, _ := s.userRepo.GetByEmailWithContext(ctx, req.Email, req.RestaurantID); existing != nil {
		return nil, errors.New("user with this email already exists in this restaurant")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		RestaurantID: req.RestaurantID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.userRepo.CreateWithContext(ctx, user); err != nil {
		return nil, err
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return user, nil
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.config.JWTExpiration) * time.Hour)

	claims := &JWTClaims{
		UserID:       user.ID,
		RestaurantID: user.RestaurantID, // Always present
		Email:        user.Email,
		Role:         user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
