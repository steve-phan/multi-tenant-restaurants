package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/models"

	brevo "github.com/getbrevo/brevo-go/lib"
)

// EmailTemplateID constants for Brevo template IDs
// These should be configured in Brevo dashboard and updated here
const (
	TemplateRestaurantWelcome       int64 = 2
	TemplateUserInvitation          int64 = 3
	TemplatePasswordReset           int64 = 4
	TemplateOrderConfirmation       int64 = 5
	TemplateOrderStatusUpdate       int64 = 11 // Not implemented
	TemplateReservationConfirm      int64 = 6
	TemplateReservationStatusUpdate int64 = 10 // Not implemented
)

// EmailService handles email operations via Brevo
type EmailService struct {
	client      *brevo.APIClient
	config      *config.Config
	senderEmail string
	senderName  string
}

// NewEmailService creates a new EmailService instance
func NewEmailService(cfg *config.Config) *EmailService {
	// Configure Brevo API client
	configuration := brevo.NewConfiguration()
	configuration.AddDefaultHeader("api-key", cfg.BrevoAPIKey)

	client := brevo.NewAPIClient(configuration)

	return &EmailService{
		client:      client,
		config:      cfg,
		senderEmail: cfg.BrevoSenderEmail,
		senderName:  cfg.BrevoSenderName,
	}
}

// SendRestaurantWelcomeEmail sends a welcome email to a newly activated restaurant
// Uses Brevo template ID: TemplateRestaurantWelcome
func (s *EmailService) SendRestaurantWelcomeEmail(
	ctx context.Context,
	restaurant *models.Restaurant,
	adminEmail string,
	tempPassword string,
) error {
	sender := brevo.SendSmtpEmailSender{
		Name:  s.senderName,
		Email: s.senderEmail,
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: adminEmail,
			Name:  restaurant.ContactName,
		},
	}

	// Template parameters
	params := map[string]interface{}{
		"contact_name":    restaurant.ContactName,
		"restaurant_name": restaurant.Name,
		"admin_email":     adminEmail,
		"temp_password":   tempPassword,
		"frontend_url":    s.config.FrontendURL,
	}

	emailRequest := brevo.SendSmtpEmail{
		Sender:     &sender,
		To:         to,
		TemplateId: TemplateRestaurantWelcome,
		Params:     params,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
} // GenerateSecurePassword generates a secure random password
// Format: 12 characters with uppercase, lowercase, numbers, and symbols
func GenerateSecurePassword() (string, error) {
	const (
		length    = 12
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		numbers   = "0123456789"
		symbols   = "!@#$%^&*"
		allChars  = uppercase + lowercase + numbers + symbols
	)

	password := make([]byte, length)

	// Ensure at least one character from each category
	categories := []string{uppercase, lowercase, numbers, symbols}
	for i, category := range categories {
		char, err := randomChar(category)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Fill remaining characters randomly
	for i := len(categories); i < length; i++ {
		char, err := randomChar(allChars)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Shuffle the password to avoid predictable patterns
	for i := length - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

// randomChar returns a random character from the given string
func randomChar(chars string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, err
	}
	return chars[n.Int64()], nil
}

// ExtractFirstName extracts the first name from a full name
func ExtractFirstName(fullName string) string {
	parts := strings.Fields(fullName)
	if len(parts) > 0 {
		return parts[0]
	}
	return fullName
}

// ExtractLastName extracts the last name from a full name
func ExtractLastName(fullName string) string {
	parts := strings.Fields(fullName)
	if len(parts) > 1 {
		return strings.Join(parts[1:], " ")
	}
	return ""
}

// OrderItem represents an item in an order for email template
type OrderItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

// SendUserInvitationEmail sends an invitation email to a new user
// Uses Brevo template ID: TemplateUserInvitation
func (s *EmailService) SendUserInvitationEmail(
	ctx context.Context,
	userEmail string,
	userFirstName string,
	restaurantName string,
	inviterName string,
	tempPassword string,
	userRole string,
) error {
	sender := brevo.SendSmtpEmailSender{
		Name:  s.senderName,
		Email: s.senderEmail,
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: userEmail,
			Name:  userFirstName,
		},
	}

	roleDescription := map[string]string{
		"Admin":  "as an administrator with full access to manage the restaurant",
		"Staff":  "as a staff member to help manage orders and operations",
		"Client": "to place orders and make reservations",
	}

	roleDesc, ok := roleDescription[userRole]
	if !ok {
		roleDesc = "to your restaurant"
	}

	// Template parameters
	params := map[string]interface{}{
		"user_first_name":  userFirstName,
		"inviter_name":     inviterName,
		"restaurant_name":  restaurantName,
		"user_email":       userEmail,
		"temp_password":    tempPassword,
		"user_role":        userRole,
		"role_description": roleDesc,
		"frontend_url":     s.config.FrontendURL,
	}

	emailRequest := brevo.SendSmtpEmail{
		Sender:     &sender,
		To:         to,
		TemplateId: TemplateUserInvitation,
		Params:     params,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send user invitation email: %w", err)
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email
// Uses Brevo template ID: TemplatePasswordReset
func (s *EmailService) SendPasswordResetEmail(
	ctx context.Context,
	userEmail string,
	userFirstName string,
	resetToken string,
	expirationHours int,
) error {
	sender := brevo.SendSmtpEmailSender{
		Name:  s.senderName,
		Email: s.senderEmail,
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: userEmail,
			Name:  userFirstName,
		},
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.config.FrontendURL, resetToken)

	// Template parameters
	params := map[string]interface{}{
		"user_first_name":  userFirstName,
		"reset_link":       resetLink,
		"reset_token":      resetToken,
		"expiration_hours": expirationHours,
	}

	emailRequest := brevo.SendSmtpEmail{
		Sender:     &sender,
		To:         to,
		TemplateId: TemplatePasswordReset,
		Params:     params,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

// SendOrderConfirmationEmail sends order confirmation email to customer
// Uses Brevo template ID: TemplateOrderConfirmation
func (s *EmailService) SendOrderConfirmationEmail(
	ctx context.Context,
	customerEmail string,
	customerName string,
	restaurantName string,
	orderID uint,
	items []OrderItem,
	subtotal float64,
	tax float64,
	deliveryFee float64,
	total float64,
	estimatedMinutes int,
	specialNotes string,
	restaurantPhone string,
	restaurantAddress string,
) error {
	sender := brevo.SendSmtpEmailSender{
		Name:  s.senderName,
		Email: s.senderEmail,
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: customerEmail,
			Name:  customerName,
		},
	}

	// Template parameters
	params := map[string]interface{}{
		"customer_name":      customerName,
		"restaurant_name":    restaurantName,
		"order_id":           orderID,
		"order_items":        items,
		"subtotal":           subtotal,
		"tax":                tax,
		"delivery_fee":       deliveryFee,
		"total":              total,
		"estimated_minutes":  estimatedMinutes,
		"special_notes":      specialNotes,
		"restaurant_phone":   restaurantPhone,
		"restaurant_address": restaurantAddress,
		"frontend_url":       s.config.FrontendURL,
	}

	emailRequest := brevo.SendSmtpEmail{
		Sender:     &sender,
		To:         to,
		TemplateId: TemplateOrderConfirmation,
		Params:     params,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send order confirmation email: %w", err)
	}

	return nil
}

// SendOrderStatusUpdateEmail sends order status update email
// Uses Brevo template ID: TemplateOrderStatusUpdate
func (s *EmailService) SendOrderStatusUpdateEmail(
	ctx context.Context,
	customerEmail string,
	customerName string,
	restaurantName string,
	orderID uint,
	status string,
	statusMessage string,
	statusEmoji string,
	estimatedMinutes int,
) error {
	sender := brevo.SendSmtpEmailSender{
		Name:  s.senderName,
		Email: s.senderEmail,
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: customerEmail,
			Name:  customerName,
		},
	}

	// Template parameters
	params := map[string]interface{}{
		"customer_name":     customerName,
		"restaurant_name":   restaurantName,
		"order_id":          orderID,
		"status":            status,
		"status_message":    statusMessage,
		"status_emoji":      statusEmoji,
		"estimated_minutes": estimatedMinutes,
		"frontend_url":      s.config.FrontendURL,
	}

	emailRequest := brevo.SendSmtpEmail{
		Sender:     &sender,
		To:         to,
		TemplateId: TemplateOrderStatusUpdate,
		Params:     params,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send order status update email: %w", err)
	}

	return nil
}

// SendReservationConfirmationEmail sends reservation confirmation email
// Uses Brevo template ID: TemplateReservationConfirm
func (s *EmailService) SendReservationConfirmationEmail(
	ctx context.Context,
	customerEmail string,
	customerName string,
	restaurantName string,
	reservationID uint,
	reservationDate string,
	reservationTime string,
	durationMinutes int,
	numberOfGuests int,
	tableNumber string,
	specialRequests string,
	restaurantAddress string,
	restaurantPhone string,
	confirmationCode string,
) error {
	sender := brevo.SendSmtpEmailSender{
		Name:  s.senderName,
		Email: s.senderEmail,
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: customerEmail,
			Name:  customerName,
		},
	}

	// Template parameters
	params := map[string]interface{}{
		"customer_name":      customerName,
		"restaurant_name":    restaurantName,
		"reservation_id":     reservationID,
		"reservation_date":   reservationDate,
		"reservation_time":   reservationTime,
		"duration_minutes":   durationMinutes,
		"number_of_guests":   numberOfGuests,
		"table_number":       tableNumber,
		"special_requests":   specialRequests,
		"restaurant_address": restaurantAddress,
		"restaurant_phone":   restaurantPhone,
		"confirmation_code":  confirmationCode,
		"frontend_url":       s.config.FrontendURL,
	}

	emailRequest := brevo.SendSmtpEmail{
		Sender:     &sender,
		To:         to,
		TemplateId: TemplateReservationConfirm,
		Params:     params,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send reservation confirmation email: %w", err)
	}

	return nil
}

// SendReservationStatusUpdateEmail sends reservation status update email
// Uses Brevo template ID: TemplateReservationStatusUpdate
func (s *EmailService) SendReservationStatusUpdateEmail(
	ctx context.Context,
	customerEmail string,
	customerName string,
	restaurantName string,
	reservationID uint,
	status string,
	statusMessage string,
	reservationDate string,
	reservationTime string,
	cancellationReason string,
) error {
	sender := brevo.SendSmtpEmailSender{
		Name:  s.senderName,
		Email: s.senderEmail,
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: customerEmail,
			Name:  customerName,
		},
	}

	// Template parameters
	params := map[string]interface{}{
		"customer_name":       customerName,
		"restaurant_name":     restaurantName,
		"reservation_id":      reservationID,
		"status":              status,
		"status_message":      statusMessage,
		"reservation_date":    reservationDate,
		"reservation_time":    reservationTime,
		"cancellation_reason": cancellationReason,
		"frontend_url":        s.config.FrontendURL,
	}

	emailRequest := brevo.SendSmtpEmail{
		Sender:     &sender,
		To:         to,
		TemplateId: TemplateReservationStatusUpdate,
		Params:     params,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send reservation status update email: %w", err)
	}

	return nil
}
