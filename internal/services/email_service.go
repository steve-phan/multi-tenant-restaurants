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
func (s *EmailService) SendRestaurantWelcomeEmail(
	ctx context.Context,
	restaurant *models.Restaurant,
	adminEmail string,
	tempPassword string,
) error {
	// Prepare email parameters
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

	subject := fmt.Sprintf("Welcome to Restaurant Platform - %s is now Active!", restaurant.Name)

	// HTML email body
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .credentials { background: white; padding: 20px; border-left: 4px solid #667eea; margin: 20px 0; border-radius: 4px; }
        .button { display: inline-block; background: #667eea; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #6b7280; font-size: 14px; }
        .warning { background: #fef3c7; border-left: 4px solid #f59e0b; padding: 15px; margin: 20px 0; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to Restaurant Platform!</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>Great news! Your restaurant <strong>"%s"</strong> has been successfully activated and is now live on our platform.</p>
            
            <div class="credentials">
                <h3>Your Admin Account Credentials</h3>
                <p><strong>Email:</strong> %s</p>
                <p><strong>Temporary Password:</strong> <code style="background: #f3f4f6; padding: 4px 8px; border-radius: 4px; font-family: monospace;">%s</code></p>
            </div>

            <div class="warning">
                <strong>⚠️ Important Security Notice:</strong><br>
                Please change your password immediately after your first login for security purposes.
            </div>

            <div style="text-align: center;">
                <a href="%s/login" class="button">Login to Your Dashboard</a>
            </div>

            <h3>What's Next?</h3>
            <ul>
                <li>Complete your restaurant profile</li>
                <li>Add your menu items</li>
                <li>Configure table settings</li>
                <li>Start accepting orders and reservations</li>
            </ul>

            <p>If you have any questions or need assistance, please don't hesitate to reach out to your Key Account Manager.</p>

            <p>Best regards,<br>
            <strong>The Restaurant Platform Team</strong></p>
        </div>
        <div class="footer">
            <p>This is an automated message. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
	`, restaurant.ContactName, restaurant.Name, adminEmail, tempPassword, s.config.FrontendURL)

	// Create email request
	emailRequest := brevo.SendSmtpEmail{
		Sender:      &sender,
		To:          to,
		Subject:     subject,
		HtmlContent: htmlContent,
	}

	// Send email
	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailRequest)
	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

// GenerateSecurePassword generates a secure random password
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
